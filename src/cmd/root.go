// certificateManager
// src/cmd/root.go

package cmd

import (
	"certificateManager/cert"
	"certificateManager/environment"
	"certificateManager/helpers"
	"github.com/spf13/cobra"
	"os"
)

//var version = "1.001-0 (2023.10.03)"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cm",
	Short:   "Certificate / PKI management tool",
	Version: helpers.White("1.001-0 (2023.10.03)"),
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
	certCmd.AddCommand(certRevokeCmd)

	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envRmCmd)
	envCmd.AddCommand(envAddCmd)
	envCmd.AddCommand(envInfoCmd)

	rootCmd.PersistentFlags().StringVarP(&environment.EnvConfigFile, "env", "e", "defaultEnv.json", "Default environment configuration file; this is a per-user setting.")
	certCreateCmd.PersistentFlags().StringVarP(&certs.CertName, "file", "f", "", "JSON file holding the certificate config.")
	certCreateCmd.PersistentFlags().BoolVarP(&certs.CertJava, "java", "j", false, "Also create a Java Keystore (JKS).")
	certRevokeCmd.PersistentFlags().BoolVarP(&certs.CertRemoveFiles, "remove", "r", false, "Remove all artefacts from PKI.")
	certVerifyCmd.Flags().BoolVarP(&certs.CaVerifyVerbose, "verbose", "v", false, "Display the full output.")
	certVerifyCmd.Flags().BoolVarP(&certs.CaVerifyComments, "comments", "", false, "Display the comments (if any) at the end of the configuration file.")
	certCreateCmd.Flags().IntVarP(&certs.CertPKsize, "keysize", "b", 4096, "Certificate private key size in bits.")
}
