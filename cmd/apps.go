/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"

	"github.com/stuttgart-things/kaeffken/models"
	"github.com/stuttgart-things/kaeffken/modules"

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
	AppDefaults, err := modules.ReadYAMLFile[models.AppDefaults]("/home/patrick/projects/kaeffken/tests/app-defaults.yaml")
	if err != nil {
		log.Fatalf("Error reading flux-defaults.yaml: %v", err)
	}
	fmt.Printf("App Defaults: %+v\n", AppDefaults)

	// Read app values
	AppValues, err := modules.ReadYAMLFile[models.Apps]("/home/patrick/projects/kaeffken/tests/apps.yaml")
	if err != nil {
		log.Fatalf("Error reading flux-defaults.yaml: %v", err)
	}
	fmt.Printf("App Defaults: %+v\n", AppValues)

}

func init() {
	rootCmd.AddCommand(appsCmd)
}
