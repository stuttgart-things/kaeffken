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
		Short: "apps will render app configs",
		Long: `Print the apps information. For example:
	sthings apps`,

		Run: func(_ *cobra.Command, _ []string) {
			Flux()
		},
	}
)

func Flux() {
	fmt.Println("HELLO")

	// Read flux defaults
	fluxDefaults, err := modules.ReadYAMLFile[models.FluxDefaults]("/home/patrick/projects/kaeffken/tests/flux-defaults.yaml")
	if err != nil {
		log.Fatalf("Error reading flux-defaults.yaml: %v", err)
	}
	fmt.Printf("Flux Defaults: %+v\n", fluxDefaults)

	// Read app defaults
	appDefaults, err := modules.ReadYAMLFile[models.AppDefaults]("/home/patrick/projects/kaeffken/tests/app-defaults.yaml")
	if err != nil {
		log.Fatalf("Error reading flux-defaults.yaml: %v", err)
	}
	fmt.Printf("App Defaults: %+v\n", appDefaults)

	// Read app values
	appValues, err := modules.ReadYAMLFile[models.Apps]("/home/patrick/projects/kaeffken/tests/apps.yaml")
	if err != nil {
		log.Fatalf("Error reading flux-defaults.yaml: %v", err)
	}
	fmt.Printf("App Defaults: %+v\n", appValues)

	for key, appValues := range appValues.Flux {

		// CHECK IF APP(KEY) EXISTS IN fluxDefaults
		if _, ok := fluxDefaults.FluxAppDefaults[key]; ok {

			fmt.Println("FOUND THE APP!")

			// SET SUBSTITUTION VARIABLES
			defaultVariables := fluxDefaults.FluxAppDefaults[key].Variables
			flagVariables := appValues.Variables

			// MERGE DEFAULT VARIABLES + VALUES
			variables := sthingsBase.MergeMaps(defaultVariables, flagVariables)
			fmt.Println("MERGED VARS: ", variables)

			substituteValues := make(map[string]interface{})
			for _, variable := range variables {
				substituteValues[variable.Name] = variable.Value
			}

			fmt.Println("DEFAULT SPPPEEECC", appDefaults.FluxKustomization.Spec)
			fmt.Println("APP SPPPEEECC", appValues.Spec)

			kustomization := models.Kustomization{
				APIVersion: appDefaults.FluxKustomization.CR.APIVersion,
				Kind:       appDefaults.FluxKustomization.CR.Kind,
				Metadata:   models.Metadata{Name: modules.SetAppParameter(appValues.Name, key, "NOT-DEFINED"), Namespace: appDefaults.FluxKustomization.CR.Namespace},
				Spec: models.Spec{
					Interval:      modules.SetAppParameter(appValues.Spec.Interval, fluxDefaults.FluxAppDefaults[key].Spec.Interval, appDefaults.FluxKustomization.Spec.Interval),
					RetryInterval: modules.SetAppParameter(appDefaults.FluxKustomization.Spec.RetryInterval, fluxDefaults.FluxAppDefaults[key].Spec.RetryInterval, appDefaults.FluxKustomization.Spec.RetryInterval),
					Timeout:       modules.SetAppParameter(appDefaults.FluxKustomization.Spec.Timeout, fluxDefaults.FluxAppDefaults[key].Spec.Timeout, appDefaults.FluxKustomization.Spec.Timeout),
					SourceRef:     models.SourceRef{Kind: modules.SetAppParameter(appDefaults.FluxKustomization.Spec.SourceRef.Kind, fluxDefaults.FluxAppDefaults[key].Spec.SourceRef.Kind, appDefaults.FluxKustomization.Spec.SourceRef.Kind), Name: modules.SetAppParameter(appDefaults.FluxKustomization.Spec.SourceRef.Name, fluxDefaults.FluxAppDefaults[key].Spec.SourceRef.Name, appDefaults.FluxKustomization.Spec.SourceRef.Name)},
					PostBuild:     models.PostBuild{Substitute: substituteValues},
				},
			}

			rendered, err := modules.RenderTemplate(models.TemplateFluxKustomization, kustomization)
			if err != nil {
				log.Fatalf("Error rendering template: %v", err)
			}

			fmt.Println(rendered)

		} else {
			fmt.Println("APP NOT FOUND!")
		}
	}

}

func init() {
	rootCmd.AddCommand(appsCmd)
}
