/*
Copyright © 2023 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	values = make(map[string]string)
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "create and deploy gitops apps",
	Long:  `create and deploy gitops applications.`,
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

		// OUTPUT ALL VALUES
		for key, value := range values {
			log.Info(strings.ToUpper(key)+": ", value)
		}

		// VERIFY NEEDED VALUES

		// LOAD CLUSTERFILE - DEFAULT IS <ROOT>/<ENV>/<LAB>/clusters.yaml
		repository, cloned := sthingsCli.CloneGitRepository(values["repository"], values["branch"], values["commitID"], nil)
		fmt.Println(repository)

		if !cloned {
			log.Error("GIT REPOSITORY CAN NOT BE CLONED: ", values["repository"])
			os.Exit(3)
		}
		// LOAD CLUSTERSFILE
		fileList, directoryList := sthingsCli.GetFileListFromGitRepository("clusters/labul/", repository)
		fmt.Println(fileList, directoryList)

		if sthingsBase.CheckForStringInSlice(fileList, values["clustersfileName"]) {
			clusterFile := sthingsCli.ReadFileContentFromGitRepo(repository, "clusters/labul/"+values["clustersfileName"])
			fmt.Println(clusterFile)
		} else {
			log.Error("CLUSTERFILE DOES NOT EXIST IN REPOSITORY: ", gitRepository+":"+"clusters/labul/"+values["clustersfileName"])
			os.Exit(3)
		}

		// LOAD FLUX INFRA CATALOGUE
		fileList, directoryList = sthingsCli.GetFileListFromGitRepository("clusters/labul/", repository)
		fmt.Println(fileList, directoryList)

		if sthingsBase.CheckForStringInSlice(fileList, values["clustersfileName"]) {
			fmt.Println("LETS GOOO")
			infraCatalog := sthingsCli.ReadFileContentFromGitRepo(repository, "clusters/config/"+"infraCatalog.json")
			fmt.Println(infraCatalog)
		} else {
			log.Error("CLUSTERFILE DOES NOT EXIST IN REPOSITORY: ", gitRepository+":"+"clusters/config/"+"infraCatalog.json")
			os.Exit(3)
		}

	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("root", "clusters", "cluster root path in repository")
	deployCmd.Flags().String("env", "labul/vsphere", "env path in repository")
	deployCmd.Flags().String("name", "", "cluster name")
	deployCmd.Flags().String("clustersfile", "clusters.yaml", "clustersfile name")
}
