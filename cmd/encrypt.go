/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/stuttgart-things/kaeffken/models"
	"github.com/stuttgart-things/kaeffken/modules"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var (
	sourceFile      string
	secretTemplates = make(map[string]string)
	keyValues       = make(map[string]interface{})
	variables       = make(map[string]interface{})
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt secrets",
	Long:  `encrypt secret files`,
	Run: func(cmd *cobra.Command, args []string) {

		// DEFAULTS
		secretTemplates["k8s"] = models.K8sSecret
		encryptedSecrets := make(map[string]string)

		// FLAGS
		source, _ := cmd.LocalFlags().GetString("source")
		ageKey, _ := cmd.LocalFlags().GetString("age")
		outputFormat, _ := cmd.LocalFlags().GetString("output")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")
		template, _ := cmd.LocalFlags().GetString("template")
		metaName, _ := cmd.LocalFlags().GetString("name")
		metaNamespace, _ := cmd.LocalFlags().GetString("namespace")
		templateValues, _ := cmd.Flags().GetStringSlice("values")
		branchName, _ := cmd.LocalFlags().GetString("branch")
		createPullRequest, _ := cmd.LocalFlags().GetBool("pr")

		if branchName == "" {
			log.Warn("BRANCH NAME EMPTY")
			branchName = "kaeffken-secret-" + timestamp
			log.Info("BRANCH NAME GENERATED: ", branchName)
		}

		// CHECK IF TEMPLATE IS SET
		if template != "" {
			log.Info("SECRET TEMPLATE: ", template)
			if _, exists := secretTemplates[template]; exists {
				log.Info("SECRET TEMPLATE FOUND: ", template)
			} else {
				log.Error("SECRET TEMPLATE NOT FOUND: ", template)
				os.Exit(3)
			}
		}

		// READ VALUES (IF DEFINED)
		if len(templateValues) > 0 {
			keyValues = sthingsCli.VerifyReadKeyValues(templateValues, log, true)
			variables["metaName"] = metaName
			variables["metaNamespace"] = metaNamespace
			variables["Data"] = keyValues
			log.Info("VARIABLES: ", variables)

			// RENDER TEMPLATE (IF DEFINED)
			renderedTemplate, err := sthingsBase.RenderTemplateInline(secretTemplates[template], "missingkey=error", "{{", "}}", variables)
			if err != nil {
				log.Error("ERROR WHILE TEMPLATING", err)
			}

			fmt.Println(string(renderedTemplate))
			sourceFile = string(renderedTemplate)

		} else {
			log.Warn("NO VALUES DEFINED")
		}

		// CHECK FOR LOCAL FILE SOURCE
		if source != "" && template == "" {
			log.Warn("NO SOURCE DEFINED")
			sourceExists, _ := sthingsBase.VerifyIfFileOrDirExists(source, "file")
			if sourceExists {
				log.Info("SOURCE SECRET FOUND : ", source)
				sourceFile = sthingsBase.ReadFileToVariable(source)
			} else {
				log.Error("LOCAL SECRET NOT FOUND : ", source)
				os.Exit(3)
			}
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

		encryptedSecrets[metaName] = encryptedSecret

		// HANDLE OUTPUT
		filesList := modules.HandleRenderOutput(encryptedSecrets, outputFormat, destinationPath, clusterPath)

		fmt.Println(filesList)
		// CREATE PULL REQUEST
		if createPullRequest && outputFormat != "stdout" {
			modules.CreateGitHubPullRequest(token, gitOwner, gitOwner, "kaeffken@sthings.com", gitRepo, "test-commit", branchName, filesList)
		}

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().String("template", "", "render a template and encrypt it")
	encryptCmd.Flags().String("source", "", "source/path of to be encrpted/secret file")
	encryptCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	encryptCmd.Flags().String("destination", "", "path to output (if output file)")
	encryptCmd.Flags().String("age", "", "age key")
	encryptCmd.Flags().StringSlice("values", []string{}, "templating values")
	encryptCmd.Flags().String("name", "secret", "meta name")
	encryptCmd.Flags().String("namespace", "default", "meta namespace")
	encryptCmd.Flags().Bool("pr", false, "create pull request")
	encryptCmd.Flags().String("branch", "", "name of to be created branch")
}
