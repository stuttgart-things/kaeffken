/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

var (
	gitRepository string
	logFilePath   string
	gitBranch     string
	gitCommitID   string
	log           = sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)
)

var rootCmd = &cobra.Command{
	Use:   "kaeffken",
	Short: "KAEFFKEN CLI",
	Long:  `KAEFFKEN CLI - GITOPS CLUSTER MANAGEMENT CLI`,
}

func Execute(defCmd string) {
	var cmdFound bool
	cmd := rootCmd.Commands()

	for _, a := range cmd {
		for _, b := range os.Args[1:] {
			if a.Name() == b {
				cmdFound = true
				break
			}
		}
	}
	if !cmdFound {
		args := append([]string{defCmd}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&gitRepository, "git", "https://github.com/stuttgart-things/stuttgart-things.git", "source git repository")
	rootCmd.PersistentFlags().StringVar(&logFilePath, "log", "/tmp/kaeffken.log", "log file path")
	// rootCmd.PersistentFlags().StringVar(&gitUser, "gitUser", "git/data/github:username", "git user")
	rootCmd.PersistentFlags().StringVar(&gitBranch, "branch", "main", "git branch")
	rootCmd.PersistentFlags().StringVar(&gitCommitID, "commitID", "", "git commit id")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
