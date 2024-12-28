/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package cmd

import (
	"os"

	"github.com/getsops/sops/v3/decrypt"
	"github.com/stuttgart-things/kaeffken/modules"

	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt secrets",
	Long:  `encrypt secret files`,

	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		source, _ := cmd.LocalFlags().GetString("source")
		ageKey, _ := cmd.LocalFlags().GetString("key")
		fileFormat, _ := cmd.LocalFlags().GetString("format")
		outputFormat, _ := cmd.LocalFlags().GetString("output")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")

		log.Info("SOURCE: ", source)
		log.Info("OUTPUT: ", outputFormat)
		log.Info("DESTINATION: ", destinationPath)
		log.Info("FORMAT: ", fileFormat)

		// CHECK IF AGE KEY IS SET
		if ageKey != "" {
			os.Setenv("SOPS_AGE_KEY", ageKey)
			log.Info("USING AGE KEY: ", ageKey)
		}

		if ageKey == "" && os.Getenv("SOPS_AGE_KEY") == "" {
			log.Warn("SOPS_AGE_KEY NOT SET")
			log.Error("AGE KEY NOT SET")
		}

		decryptedFile, err := decrypt.File(source, fileFormat)
		if err != nil {
			log.Error("FAILED TO DECRYPT: ", err)
		}

		secretsMap := modules.CreateSecretsMap(decryptedFile, nil)
		log.Info("CREATED SECRETS MAP: ", secretsMap)

		// HANDLE OUTPUT
		modules.HandleOutput(outputFormat, destinationPath, string(decryptedFile))
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().String("source", "", "source/path of to be decrypted/secret file")
	decryptCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	decryptCmd.Flags().String("key", "", "sops age key")
	decryptCmd.Flags().String("format", "yaml", "sops file format/extension")
	decryptCmd.Flags().String("destination", "", "path to output (if output file)")
}
