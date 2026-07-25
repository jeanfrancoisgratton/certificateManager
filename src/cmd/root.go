// certificateManager
// src/cmd/root.go

package cmd

import (
	"certificateManager/cert"
	"certificateManager/environment"
	"fmt"
	"os"
	"runtime"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cm",
	Short:   "Certificate / PKI management tool",
	Version: hftx.White(fmt.Sprintf("1.8.0-%s 2026.07.25", runtime.GOARCH)),
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

	rootCmd.AddCommand(completionCmd, certCmd, envCmd)

	rootCmd.PersistentFlags().StringVarP(&environment.EnvConfigFile, "env", "e", "defaultEnv.json", "Default env configuration file; this is a per-user setting.")
	certCreateCmd.PersistentFlags().BoolVarP(&cert.CertJava, "java", "j", false, "Also create a Java Keystore (JKS).")
	certRevokeCmd.PersistentFlags().BoolVarP(&cert.CertRemoveFiles, "remove", "r", false, "Remove all artefacts from PKI.")
	certVerifyCmd.Flags().BoolVarP(&cert.CaVerifyVerbose, "verbose", "v", false, "Display the full output.")
	certVerifyCmd.Flags().BoolVarP(&cert.CaVerifyComments, "comments", "c", false, "Display the comments (if any) at the end of the configuration file.")
	certCreateCmd.Flags().IntVarP(&cert.CertPKsize, "keysize", "b", 4096, "Certificate private key size in bits.")
}
