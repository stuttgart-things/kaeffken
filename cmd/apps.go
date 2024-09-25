/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"github.com/stuttgart-things/kaeffken/modules"

	"github.com/spf13/cobra"
)

var (
	appsCmd = &cobra.Command{
		Use:   "apps",
		Short: "render apps configs",
		Long:  `render apps configs for different app kinds`,
		Run: func(cmd *cobra.Command, args []string) {

			renderedTemplates := make(map[string]string)

			appKind, _ := cmd.LocalFlags().GetString("kind")
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
				renderedTemplates = modules.RenderFluxApplication(defaultsPath, appDefaultsPath, appsPath)
			default:
				log.Error("UNKNOWN APP KIND: ", appKind)
			}

			// HANDLE OUTPUT
			filesList := modules.HandleRenderOutput(renderedTemplates, outputFormat, outputDir, clusterPath)

			// CREATE PULL REQUEST
			if createPullRequest && outputFormat != "stdout" {
				modules.CreateGitHubPullRequest(token, gitOwner, gitOwner, "kaeffken@sthings.com", gitRepo, "test-commit", filesList)
			}
		},
	}
)

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
