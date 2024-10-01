/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"fmt"
	"os"
	"strings"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/google/go-github/v62/github"
	"github.com/stuttgart-things/kaeffken/models"
)

var (
	technologyDefaults string
	fluxAppDefaults    string
	apps               string
)

func RenderFluxApplication(defaultsPath, appDefaultsPath, appsPath string) (renderedTemplates map[string]string) {

	// CREATE GITHUB CLIENT
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(nil).WithAuthToken(gitHubToken)

	renderedTemplates = make(map[string]string)

	if strings.Contains(defaultsPath, "@") {
		technologyDefaults = GetFileContentFromFileInGitHubRepo(client, defaultsPath)
	} else {
		// READ YAML FILE FROM FS
		yamlFile, err := os.ReadFile(defaultsPath)
		if err != nil {
			log.Error("ERROR READING ", err)
		}
		technologyDefaults = string(yamlFile)
	}

	if strings.Contains(appDefaultsPath, "@") {
		fluxAppDefaults = GetFileContentFromFileInGitHubRepo(client, appDefaultsPath)
	} else {
		// READ YAML FILE FROM FS
		yamlFile, err := os.ReadFile(appDefaultsPath)
		if err != nil {
			log.Error("ERROR READING ", err)
		}
		fluxAppDefaults = string(yamlFile)
	}

	if strings.Contains(appsPath, "@") {
		apps = GetFileContentFromFileInGitHubRepo(client, appsPath)
	} else {
		// READ YAML FILE FROM FS
		yamlFile, err := os.ReadFile(appsPath)
		if err != nil {
			log.Error("ERROR READING ", err)
		}
		apps = string(yamlFile)
	}

	// READ FLUX DEFAULTS
	fluxDefaults, err := ReadYAMLFile[models.FluxDefaults](technologyDefaults)
	if err != nil {
		log.Error("ERROR READING ", err)
	}

	// READ APP DEFAULTS
	appDefaults, err := ReadYAMLFile[models.AppDefaults](fluxAppDefaults)
	if err != nil {
		log.Error("ERROR READING ", err)
	}

	// READ APP VALUES
	appValues, err := ReadYAMLFile[models.Apps](apps)
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
			appVariables := appValues.Variables

			// MERGE DEFAULT VARIABLES + VALUES
			variables := sthingsBase.MergeMaps(defaultVariables, appVariables)
			log.Info("MERGED VARS: ", variables)

			// CREATE VARIABLES TO BE SUBSTITUTED
			substituteValues := make(map[string]interface{})
			for _, variable := range variables {
				substituteValues[variable.Name] = variable.Value
			}

			// SET SUBSTITUTION SECRETS
			defaultSecrets := fluxDefaults.FluxAppDefaults[appkey].Secrets
			appSecrets := appValues.Secrets

			// MERGE DEFAULT VARIABLES + VALUES
			secrets := sthingsBase.MergeMaps(defaultSecrets, appSecrets)
			log.Info("MERGED SECRETS: ", secrets)

			// CREATE SECRETS TO BE SUBSTITUTED
			substituteSecrets := make(map[string]interface{})

			// SET NAME + KIND TO SECRET
			for _, secret := range secrets {
				substituteSecrets[secret.Name] = "Secret"
			}

			kustomization := models.Kustomization{
				APIVersion: appDefaults.FluxKustomization.CR.APIVersion,
				Kind:       appDefaults.FluxKustomization.CR.Kind,
				Metadata:   models.Metadata{Name: SetAppParameter(appValues.Name, appkey, "NOT-DEFINED"), Namespace: appDefaults.FluxKustomization.CR.Namespace},
				Spec: models.Spec{
					Interval:      SetAppParameter(appValues.Spec.Interval, fluxDefaults.FluxAppDefaults[appkey].Spec.Interval, appDefaults.FluxKustomization.Spec.Interval),
					RetryInterval: SetAppParameter(appValues.Spec.RetryInterval, fluxDefaults.FluxAppDefaults[appkey].Spec.RetryInterval, appDefaults.FluxKustomization.Spec.RetryInterval),
					Timeout:       SetAppParameter(appValues.Spec.Timeout, fluxDefaults.FluxAppDefaults[appkey].Spec.Timeout, appDefaults.FluxKustomization.Spec.Timeout),
					Path:          SetAppParameter("", fluxDefaults.FluxAppDefaults[appkey].Path, ""),
					SourceRef: models.SourceRef{
						Kind: SetAppParameter(appDefaults.FluxKustomization.Spec.SourceRef.Kind, fluxDefaults.FluxAppDefaults[appkey].Spec.SourceRef.Kind, appDefaults.FluxKustomization.Spec.SourceRef.Kind),
						Name: SetAppParameter(appDefaults.FluxKustomization.Spec.SourceRef.Name, fluxDefaults.FluxAppDefaults[appkey].Spec.SourceRef.Name, appDefaults.FluxKustomization.Spec.SourceRef.Name),
					},
					PostBuild: models.PostBuild{Substitute: substituteValues, SubstituteFrom: substituteSecrets},
				},
			}

			rendered, err := RenderTemplate(models.TemplateFluxKustomization, kustomization)
			if err != nil {
				log.Error("ERROR READING TEMPLATE ", err)
			}

			log.Info("TEMPLATE WAS RENDERED ", appkey)
			renderedTemplates[appkey] = rendered

			// SECRET RENDERING
			secretVariables := make(map[string]interface{})

			fmt.Println("MERGED SECRETS: ", secrets)
			for key, secret := range secrets {
				secretVariables["metaName"] = key
				secretVariables["metaNamespace"] = appDefaults.FluxKustomization.CR.Namespace

				keyValues := make(map[string]interface{})

				for _, secretValue := range secret.Data {
					fmt.Println("SECRET: ", secretValue)
					parts := strings.Split(secretValue, ":")
					if len(parts) != 2 {
						fmt.Println("Invalid secret format: ", secretValue)
						continue
					}
					keyValues[parts[0]] = parts[1]
				}
				secretVariables["Data"] = keyValues
			}

			renderedSecret, err := sthingsBase.RenderTemplateInline(models.K8sSecret, "missingkey=error", "{{", "}}", secretVariables)
			if err != nil {
				log.Error("ERROR WHILE TEMPLATING", err)
			}

			fmt.Println(string(renderedSecret))

			// ENCRIPT SECRET

		} else {
			log.Error("APP NOT FOUND! ", appkey)
		}
	}
	return
}
