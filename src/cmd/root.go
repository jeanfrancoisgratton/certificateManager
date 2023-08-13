// certificateManager
// src/cmd/root.go

package cmd

import (
	"certificateManager/helpers"
	"github.com/spf13/cobra"
	"os"
)

var version = "0.100-0 (2023.08.13)"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "certificateManager",
	Short:   "Add a short description here",
	Long:    "Add a long description here",
	Version: version,
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
}

func init() {

	rootCmd.AddCommand(clCmd)
}
