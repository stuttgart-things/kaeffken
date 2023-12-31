/*
Copyright © 2023 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "create and deploy gitops apps",
	Long:  `create and deploy gitops applications.`,
	Run: func(cmd *cobra.Command, args []string) {

		// GET ENV VARS, CLUSTER VARS AND APP CATALOUGE
		// repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
		// profileFile = sthingsCli.ReadFileContentFromGitRepo(repo, profile)
		fmt.Println("deploy called")

	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
