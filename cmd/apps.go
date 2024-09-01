/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	appsCmd = &cobra.Command{
		Use:   "apps",
		Short: "apps will output the current build information",
		Long: `Print the apps information. For example:
	sthings apps`,

		Run: func(_ *cobra.Command, _ []string) {
			Test()
		},
	}
)

func Test() {
	fmt.Println("HELLO")
}

func init() {
	rootCmd.AddCommand(appsCmd)
}
