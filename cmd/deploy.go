/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/stuttgart-things/kaeffken/modules"

	"github.com/spf13/cobra"
)

var (
	values         = make(map[string]string)
	mandatoryFlags = []string{"repository", "branch", "clusterName", "envPath"}
)

// DEPLOYCMD REPRESENTS THE DEPLOY COMMAND
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "CREATE AND DEPLOY GITOPS APPS",
	Long:  `CREATE AND DEPLOY GITOPS APPLICATIONS.`,
	Run: func(cmd *cobra.Command, args []string) {

		// SET DEFAULTS
		values["logFilePath"] = logFilePath
		values["repository"] = gitRepository
		values["branch"] = gitBranch
		values["commitID"] = gitCommitID

		// READ FLAGS
		values["rootPath"], _ = cmd.LocalFlags().GetString("root")
		values["envPath"], _ = cmd.LocalFlags().GetString("env")
		values["clusterName"], _ = cmd.LocalFlags().GetString("name")
		values["clustersfileName"], _ = cmd.LocalFlags().GetString("clustersfile")

		// SET VARS
		values["clusterPath"] = values["rootPath"] + "/" + values["envPath"] + "/" + values["clusterName"]
		values["clusterFilePath"] = values["rootPath"] + "/" + values["envPath"] + "/" + values["clustersfileName"]

		// VERIFY / OUTPUT ALL VALUES
		if !modules.VerifyValues(values, mandatoryFlags) {
			log.Error("KAEFFKEN EXITED")
			os.Exit(3)
		}

		// LOAD CLUSTERFILE - DEFAULT IS <ROOT>/<ENV>/<LAB>/clusters.yaml
		repository, cloned := modules.CloneGitRepository(values)

		if !cloned {
			log.Error("GIT REPOSITORY CAN NOT BE CLONED: ", values["repository"])
			os.Exit(3)
		}

		// LOAD CLUSTERSFILE
		clustersFile := modules.LoadDataFromRepository(repository, values["clusterFilePath"])
		fmt.Println(clustersFile)

		// LOAD FLUX INFRA CATALOGUE
		infraCatalogue := modules.LoadDataFromRepository(repository, "clusters/config/infraCatalog.json")
		fmt.Println(infraCatalogue)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("root", "clusters", "cluster root path in repository")
	deployCmd.Flags().String("env", "labul/vsphere", "env path in repository")
	deployCmd.Flags().String("name", "", "cluster name")
	deployCmd.Flags().String("clustersfile", "clusters.yaml", "clustersfile name")
}
