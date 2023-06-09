// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/cmd/root.go

package cmd

import (
	"cm/helpers"
	"github.com/spf13/cobra"
	"os"
)

var version = "0.400 (2023.04.22)"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cm {ca|config|cert}",
	Short:   "A rootCA and server certificates management tool",
	Version: version,
	Long:    `This tools allows you to manipulate your custom root CAs and all certificates signed against that rootCA.`,
}

var clCmd = &cobra.Command{
	Use:     "changelog",
	Aliases: []string{"cl"},
	Short:   "Shows changelog",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.Changelog()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	//fmt.Println("Remote Host is ==>", hosts.ConnectURI)
}

func init() {
	rootCmd.AddCommand(clCmd)
	rootCmd.PersistentFlags().StringVarP(&helpers.CertConfigFile, "config", "c", "defaultCertConfig.json", "certificate configuration file.")
}
