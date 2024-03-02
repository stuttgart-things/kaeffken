/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"os"

	"github.com/stuttgart-things/kaeffken/modules"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var (
	sourceFile string
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt secrets",
	Long:  `encrypt secret files`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		source, _ := cmd.LocalFlags().GetString("source")
		ageKey, _ := cmd.LocalFlags().GetString("age")
		outputFormat, _ := cmd.LocalFlags().GetString("output")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")

		// CHECK FOR LOCAL FILE SOURCE
		sourceExists, _ := sthingsBase.VerifyIfFileOrDirExists(source, "file")
		if sourceExists {
			log.Info("SOURCE SECRET FOUND : ", source)
			sourceFile = sthingsBase.ReadFileToVariable(source)
		} else {
			log.Error("LOCAL SECRET NOT FOUND : ", source)
			os.Exit(3)
		}

		// CHECK FOR AGE KEY - IF EMPTY CREATE A NEW ONE
		if ageKey == "" {
			log.Warn("AGE KEY EMPTY, WILL CREATE ONE")
			identity := sthingsCli.GenerateAgeIdentitdy()
			ageKey = identity.Recipient().String()
			log.Info("GENERATED AGE KEY: ", ageKey)
		}

		// ENCRYPT SECRET WITH SOPS
		encryptedSecret := sthingsCli.EncryptStore(ageKey, sourceFile)

		// OUTPUT ENCRYPTED SECRET
		modules.HandleOutput(outputFormat, destinationPath, encryptedSecret)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().String("source", "", "source/path of secret file")
	encryptCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	encryptCmd.Flags().String("destination", "", "path to output (if output file)")
	encryptCmd.Flags().String("age", "", "age key")

}
