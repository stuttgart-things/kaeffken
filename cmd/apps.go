/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/stuttgart-things/kaeffken/models"
	"github.com/stuttgart-things/kaeffken/modules"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

var (
	appsCmd = &cobra.Command{
		Use:   "apps",
		Short: "render apps configs",
		Long:  `render apps configs for different app kinds`,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(token)

			appKind, _ := cmd.LocalFlags().GetString("kind")
			renderedTemplates := make(map[string]string)
			outputFormat, _ := cmd.LocalFlags().GetString("output")
			outputDir, _ := cmd.LocalFlags().GetString("outputDir")
			defaultsPath, _ := cmd.LocalFlags().GetString("defaults")
			appDefaultsPath, _ := cmd.LocalFlags().GetString("appDefaults")
			appsPath, _ := cmd.LocalFlags().GetString("apps")
			createPullRequest, _ := cmd.LocalFlags().GetBool("pr")

			log.Info("DEFAULTS LOADED FROM: ", defaultsPath)
			log.Info("APP-DEFAULTS LOADED FROM: ", appDefaultsPath)
			log.Info("APPS DECLARED AT: ", appsPath)

			switch appKind {

			case "flux":
				renderedTemplates = RenderFlux(defaultsPath, appDefaultsPath, appsPath)
			default:
				log.Error("UNKNOWN APP KIND: ", appKind)
			}

			// HANDLE OUTPUT
			filesList := modules.HandleRenderOutput(renderedTemplates, outputFormat, outputDir, clusterPath)

			// CREATE PULL REQUEST
			if createPullRequest {
				modules.CreateGitHubPullRequest(token, gitOwner, gitOwner, "kaeffken@sthings.com", gitRepo, "test-commit", filesList)
			}
		},
	}
)

func RenderFlux(defaultsPath, appDefaultsPath, appsPath string) (renderedTemplates map[string]string) {

	var technologyDefaults string
	var fluxAppDefaults string
	var apps string

	// CREATE GITHUB CLIENT
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	client = github.NewClient(nil).WithAuthToken(gitHubToken)

	renderedTemplates = make(map[string]string)

	if strings.Contains(defaultsPath, "@") {
		technologyDefaults = modules.GetFileContentFromFileInGitHubRepo(client, defaultsPath)
	} else {
		// READ YAML FILE FROM FS
		yamlFile, err := os.ReadFile(defaultsPath)
		if err != nil {
			log.Error("ERROR READING ", err)
		}
		technologyDefaults = string(yamlFile)
	}

	if strings.Contains(appDefaultsPath, "@") {
		fluxAppDefaults = modules.GetFileContentFromFileInGitHubRepo(client, appDefaultsPath)
	} else {
		// READ YAML FILE FROM FS
		yamlFile, err := os.ReadFile(appDefaultsPath)
		if err != nil {
			log.Error("ERROR READING ", err)
		}
		fluxAppDefaults = string(yamlFile)
	}

	if strings.Contains(appsPath, "@") {
		apps = modules.GetFileContentFromFileInGitHubRepo(client, appsPath)
	} else {
		// READ YAML FILE FROM FS
		yamlFile, err := os.ReadFile(appsPath)
		if err != nil {
			log.Error("ERROR READING ", err)
		}
		apps = string(yamlFile)
	}

	fmt.Println(technologyDefaults)
	fmt.Println(fluxAppDefaults)
	fmt.Println(apps)

	// sthingsCli.GetFileContentFromGithubRepo

	// READ FLUX DEFAULTS
	fluxDefaults, err := modules.ReadYAMLFile[models.FluxDefaults](technologyDefaults)
	if err != nil {
		log.Error("ERROR READING ", err)
	}

	// READ APP DEFAULTS
	appDefaults, err := modules.ReadYAMLFile[models.AppDefaults](fluxAppDefaults)
	if err != nil {
		log.Error("ERROR READING ", err)
	}

	// READ APP VALUES
	appValues, err := modules.ReadYAMLFile[models.Apps](apps)
	if err != nil {
		log.Error("ERROR READING ", err)
	}

	log.Info("FLUX DEFAULT: ", fluxDefaults)
	log.Info("FLUX APP DEFAULTS: ", appDefaults)
	log.Info("APP VALUES: ", appValues)

	for appkey, appValues := range appValues.Flux {

		// CHECK IF APP(KEY) EXISTS IN fluxDefaults
		if _, ok := fluxDefaults.FluxAppDefaults[appkey]; ok {

			log.Info("FOUND FLUX APP! ", appkey)

			// SET SUBSTITUTION VARIABLES
			defaultVariables := fluxDefaults.FluxAppDefaults[appkey].Variables
			flagVariables := appValues.Variables

			// MERGE DEFAULT VARIABLES + VALUES
			variables := sthingsBase.MergeMaps(defaultVariables, flagVariables)
			log.Info("MERGED VARS: ", variables)

			substituteValues := make(map[string]interface{})
			for _, variable := range variables {
				substituteValues[variable.Name] = variable.Value
			}

			kustomization := models.Kustomization{
				APIVersion: appDefaults.FluxKustomization.CR.APIVersion,
				Kind:       appDefaults.FluxKustomization.CR.Kind,
				Metadata:   models.Metadata{Name: modules.SetAppParameter(appValues.Name, appkey, "NOT-DEFINED"), Namespace: appDefaults.FluxKustomization.CR.Namespace},
				Spec: models.Spec{
					Interval:      modules.SetAppParameter(appValues.Spec.Interval, fluxDefaults.FluxAppDefaults[appkey].Spec.Interval, appDefaults.FluxKustomization.Spec.Interval),
					RetryInterval: modules.SetAppParameter(appValues.Spec.RetryInterval, fluxDefaults.FluxAppDefaults[appkey].Spec.RetryInterval, appDefaults.FluxKustomization.Spec.RetryInterval),
					Timeout:       modules.SetAppParameter(appValues.Spec.Timeout, fluxDefaults.FluxAppDefaults[appkey].Spec.Timeout, appDefaults.FluxKustomization.Spec.Timeout),
					Path:          modules.SetAppParameter("", fluxDefaults.FluxAppDefaults[appkey].Path, ""),
					SourceRef: models.SourceRef{
						Kind: modules.SetAppParameter(appDefaults.FluxKustomization.Spec.SourceRef.Kind, fluxDefaults.FluxAppDefaults[appkey].Spec.SourceRef.Kind, appDefaults.FluxKustomization.Spec.SourceRef.Kind),
						Name: modules.SetAppParameter(appDefaults.FluxKustomization.Spec.SourceRef.Name, fluxDefaults.FluxAppDefaults[appkey].Spec.SourceRef.Name, appDefaults.FluxKustomization.Spec.SourceRef.Name),
					},
					PostBuild: models.PostBuild{Substitute: substituteValues},
				},
			}

			rendered, err := modules.RenderTemplate(models.TemplateFluxKustomization, kustomization)
			if err != nil {
				log.Error("ERROR READING TEMPLATE ", err)
			}

			log.Info("TEMPLATE WAS RENDERED ", appkey)
			renderedTemplates[appkey] = rendered

		} else {
			log.Error("APP NOT FOUND! ", appkey)
		}
	}
	return
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.Flags().String("kind", "flux", "app kind: flux|")
	appsCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	appsCmd.Flags().String("outputDir", "/tmp", "output directory")
	appsCmd.Flags().String("defaults", "https://github.com/stuttgart-things/stuttgart-things@main:kaeffken/apps/flux/flux-defaults.yaml", "default values for technology")
	appsCmd.Flags().String("appDefaults", "https://github.com/stuttgart-things/stuttgart-things@main:kaeffken/apps/flux/app-defaults.yaml", "app default values")
	appsCmd.Flags().String("apps", "https://github.com/stuttgart-things/stuttgart-things@main:kaeffken/apps/flux/apps.yaml", "defined apps")
	appsCmd.Flags().Bool("pr", false, "create pull request")
}
