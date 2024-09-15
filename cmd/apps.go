/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"

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

			// outputFormat, _ := cmd.LocalFlags().GetString("output")

			switch appKind {
			case "flux":
				renderedTemplates = RenderFlux()
			default:
				log.Error("Unknown app kind: ", appKind)
			}

			modules.HandleRenderOutput(renderedTemplates, "file", "/tmp")

		},
	}
)

func RenderFlux() (renderedTemplates map[string]string) {

	renderedTemplates = make(map[string]string)
	// appTemplatesFiles := make(map[string]string)

	// gitHubToken := os.Getenv("GITHUB_TOKEN")

	// GET APP TEMPLATE FILES FROM LOCAL FILES
	pathFluxDefaults := "/home/patrick/projects/kaeffken/tests/flux-defaults.yaml"
	pathFluxAppDefaults := "/home/patrick/projects/kaeffken/tests/app-defaults.yaml"
	pathAppValues := "/home/patrick/projects/kaeffken/tests/apps.yaml"

	// client := modules.CreateGithubClient(gitHubToken)
	// fmt.Println(client)

	owner, repo, branch, path, _ := modules.ParseGitHubURL("https://github.com/stuttgart-things/stuttgart-things.git@main:kaeffken/apps/flux/app-defaults.yaml")
	fmt.Println(owner, repo, branch, path)

	// fileContent := modules.GetFileContentFromFileInGitHubRepo(client, "stuttgart-things", "stuttgart-things", "main", "kaeffken/apps/flux/app-defaults.yaml")
	// fmt.Println(fileContent)

	// sthingsCli.GetFileContentFromGithubRepo

	fluxDefaults, err := modules.ReadYAMLFile[models.FluxDefaults](pathFluxDefaults)
	if err != nil {
		log.Error("Error reading ", pathFluxDefaults)
	}

	// Read app defaults
	appDefaults, err := modules.ReadYAMLFile[models.AppDefaults](pathFluxAppDefaults)
	if err != nil {
		log.Error("Error reading ", pathFluxAppDefaults)
	}

	// Read app values
	appValues, err := modules.ReadYAMLFile[models.Apps](pathAppValues)
	if err != nil {
		log.Error("Error reading ", pathAppValues)
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
				log.Error("Error reading template ", err)
			}

			log.Info("TEMPLATE WAS RENDERED ", appkey)
			renderedTemplates[appkey] = rendered
			// fmt.Println(rendered)

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
}
