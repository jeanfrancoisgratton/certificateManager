// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/cmd/configCommands.go

package cmd

import (
	"cm/config"
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config-old command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration file management",
	Long:  `This is where you can create a templated file, edit/delete an existing config file, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

var configCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"generate"},
	Short:   "Generate a configuration file",
	//Long:  `This is where you can create a templated file, edit/delete an existing config file, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.ConfCreate()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var configTemplateCmd = &cobra.Command{
	Use: "template",
	//Aliases: []string{"update"},
	Short: "Create a template (blank) file",
	Long: `This is where you (re)create a templated file in case that the original has been deleted.

That file will be created in your home directory, under the .config-old/certificatemanager directory,
alongside with an explicative text file`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.TemplateConfigCreate()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configCreateCmd)
	configCmd.AddCommand(configTemplateCmd)
}
