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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cm",
	Short:   "Certificate / PKI management tool",
	Version: helpers.White("1.24.00-0 (2024.03.01)"),
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
	certCreateCmd.PersistentFlags().BoolVarP(&cert.CertJava, "java", "j", false, "Also create a Java Keystore (JKS).")
	certRevokeCmd.PersistentFlags().BoolVarP(&cert.CertRemoveFiles, "remove", "r", false, "Remove all artefacts from PKI.")
	certVerifyCmd.Flags().BoolVarP(&cert.CaVerifyVerbose, "verbose", "v", false, "Display the full output.")
	certVerifyCmd.Flags().BoolVarP(&cert.CaVerifyComments, "comments", "c", false, "Display the comments (if any) at the end of the configuration file.")
	certCreateCmd.Flags().IntVarP(&cert.CertPKsize, "keysize", "b", 4096, "Certificate private key size in bits.")
}
