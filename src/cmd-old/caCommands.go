// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/cmd/caCommands.go

package cmd

import (
	"cm/ca-old"
	"cm/helpers-old"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var privKeySize int

// caCmd represents the ca-old command
var caCmd = &cobra.Command{
	Use:   "ca-old",
	Short: "Root Certificate Authority management",
	Long:  `This is where you will manage (add/verify/delete) your rootCAs.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage ca-old {create|verify")
			os.Exit(0)
		}
	},
}

// Create a rootCA based on the config-old file as defined with the -c global flag
var caCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "rootCA creation command",
	//Long:    `This is where you will manage (add/remove) your rootCAs\' config-old files.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ca_old.CreateRootCA(privKeySize)
		if err != nil {
			fmt.Printf("%s", helpers_old.Red("Error while creating the root CA:"))
			fmt.Println(err)
		} else {
			fmt.Printf("A %v bits-keysize certificate %s has been created in %s\n", helpers_old.Green(strconv.Itoa(privKeySize)), helpers_old.Green(helpers_old.CertConfig.CertificateName), helpers_old.Green(helpers_old.CertConfig.CertificateDirectory))
		}
	},
}

// Verify a rootCA
var caVerifyCmd = &cobra.Command{
	Use:   "verify certificate_filename",
	Short: "verify the created CA certificate",
	Long:  `If you do not provide a filename extension (.crt or .pem), .crt is assumed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You need to provide the certificate name to be verified")
			os.Exit(0)
		}
		if !strings.HasSuffix(args[0], ".crt") && !strings.HasSuffix(args[0], ".pem") {
			args[0] += ".crt"
		}
		err := ca_old.VerifyCACertificate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Delete a root CA based on the config-old file as defined with the -c global flag
var caDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del", "remove"},
	Short:   "Removes a CA certificate",
	Long: `The configuration file describing the certificate should be present.
If not, empty or defaults values will be supplied, and the file will be created`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ca_old.RemoveCACertificate()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(caCmd)
	caCmd.AddCommand(caCreateCmd)
	caCmd.AddCommand(caVerifyCmd)
	caCmd.AddCommand(caDeleteCmd)

	caVerifyCmd.Flags().BoolVarP(&ca_old.CaVerifyVerbose, "verbose", "v", false, "display the full output")
	caVerifyCmd.Flags().BoolVarP(&ca_old.CaVerifyComments, "comments", "", false, "display the comments (if any) at the end of the configuration file")
	caCreateCmd.Flags().IntVarP(&privKeySize, "keysize", "b", 4096, "certificate private key size in bits")
}
