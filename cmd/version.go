/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

var (
	date       = "unknown"
	commit     = "unknown"
	output     = "yaml"
	version    = "unset"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "version will output the current build information",
		Long: `Print the version information. For example:
	sthings version`,

		Run: func(_ *cobra.Command, _ []string) {
			PrintBanner()
		},
	}
)

// https://fsymbols.com/generators/carty/
const banner = `

██╗░░██╗░█████╗░███████╗███████╗███████╗██╗░░██╗███████╗███╗░░██╗
██║░██╔╝██╔══██╗██╔════╝██╔════╝██╔════╝██║░██╔╝██╔════╝████╗░██║
█████═╝░███████║█████╗░░█████╗░░█████╗░░█████═╝░█████╗░░██╔██╗██║
██╔═██╗░██╔══██║██╔══╝░░██╔══╝░░██╔══╝░░██╔═██╗░██╔══╝░░██║╚████║
██║░╚██╗██║░░██║███████╗██║░░░░░██║░░░░░██║░╚██╗███████╗██║░╚███║
╚═╝░░╚═╝╚═╝░░╚═╝╚══════╝╚═╝░░░░░╚═╝░░░░░╚═╝░░╚═╝╚══════╝╚═╝░░╚══╝

`

// OUTPUT BANNER + VERSION OUTPUT
func PrintBanner() string {
	color.Cyan(banner)
	resp := goVersion.FuncWithOutput(false, version, commit, date, output)
	color.Magenta(resp + "\n")
	return resp
}

func init() {
	versionCmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of 'yaml' or 'json'.")
	rootCmd.AddCommand(versionCmd)
}
