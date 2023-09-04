// certificateManager
// src/cmd/root.go

package cmd

import (
	"certificateManager/certs"
	"certificateManager/environment"
	"certificateManager/helpers"
	"fmt"
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

var certCmd = &cobra.Command{
	Use:     "cert",
	Example: "cm cert { {create | delete} } certificate_name | list }",
	Short:   "Certificate sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.Changelog()
	},
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.Changelog()
	},
}

// Lists all certificates in conf directory as per $HOME/.config/certificatemanager/*.json (defined with -e flag)
var certlistCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "cm cert list",
	Short:   "Lists all certificates in defined rootDir",
	Run: func(cmd *cobra.Command, args []string) {
		if err := certs.ListCertificates(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	},
}

// Verify a given certificate
var certVerifyCmd = &cobra.Command{
	Use: "verify",
	//Aliases: []string{"ls"},
	Example: "cm cert verify FILENAME",
	Short:   "Verifies a certificate, as per the provided filename",
	Run: func(cmd *cobra.Command, args []string) {
		if err := certs.Verify(args); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	},
}

// The biggie: create a CA or "normal" SSL certificate
var certCreateCmd = &cobra.Command{
	Use: "create",
	//Aliases: []string{"ls"},
	Example: "cm cert create CERTICATE_CONFIG_FILE",
	Short:   "Creates a certificate, specifying the config file to use",
	Run: func(cmd *cobra.Command, args []string) {
		certname := ""
		if len(args) != 0 {
			certname = args[0]
		}
		if err := certs.Create(certname); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
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
	rootCmd.AddCommand(clCmd)
	rootCmd.AddCommand(certCmd)
	rootCmd.AddCommand(envCmd)
	certCmd.AddCommand(certlistCmd)
	certCmd.AddCommand(certVerifyCmd)
	certCmd.AddCommand(certCreateCmd)

	rootCmd.PersistentFlags().StringVarP(&environment.EnvConfigFile, "env", "e", "defaultEnv.json", "Default environment configuration file; this is a per-user setting.")
	rootCmd.PersistentFlags().StringVarP(&helpers.CertificatesRootDir, "rootdir", "r", ".", "Certificate root dir; all other directories are relative to this one.")

	certCmd.PersistentFlags().BoolVarP(&certs.CreateSingleCert, "single", "s", false, "Create a certificate while ignoring a given directory structure.")

	certCreateCmd.PersistentFlags().StringVarP(&certs.CertName, "file", "f", "", "JSON file holding the certificate config.")
}
