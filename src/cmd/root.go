// certificateManager
// src/cmd/root.go

package cmd

import (
	"certificateManager/certs"
	"certificateManager/environment"
	"certificateManager/helpers"
	"github.com/spf13/cobra"
	"os"
)

var version = "0.100-0 (2023.08.13)"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cm",
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

	rootCmd.DisableAutoGenTag = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(clCmd)
	rootCmd.AddCommand(certCmd)
	rootCmd.AddCommand(envCmd)
	certCmd.AddCommand(certlistCmd)
	certCmd.AddCommand(certVerifyCmd)
	certCmd.AddCommand(certCreateCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envRmCmd)
	envCmd.AddCommand(envAddCmd)
	envCmd.AddCommand(envInfoCmd)

	rootCmd.PersistentFlags().StringVarP(&environment.EnvConfigFile, "env", "e", "defaultEnv.json", "Default environment configuration file; this is a per-user setting.")
	//rootCmd.PersistentFlags().StringVarP(&helpers.CertificatesRootDir, "rootdir", "r", ".", "Certificate root dir; all other directories are relative to this one.")

	certCmd.PersistentFlags().BoolVarP(&certs.CreateSingleCert, "single", "s", false, "Create a certificate while ignoring a given directory structure.")

	certCreateCmd.PersistentFlags().StringVarP(&certs.CertName, "file", "f", "", "JSON file holding the certificate config.")
}
