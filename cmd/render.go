/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("render called")

		interactive, _ := cmd.LocalFlags().GetBool("interactive")
		fmt.Println("interactive:", interactive)

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().Bool("interactive", false, "interactive (prompted) survey. default: false")
}
