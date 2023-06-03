// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/cmd/caCommands.go

package cmd

import (
	"cm/ca"
	"cm/helpers"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var privKeySize int

// caCmd represents the ca command
var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "Root Certificate Authority management",
	Long:  `This is where you will manage (add/verify/delete) your rootCAs.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage ca {create|verify")
			os.Exit(0)
		}
	},
}

// Create a rootCA based on the config file as defined with the -c global flag
var caCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "rootCA creation command",
	//Long:    `This is where you will manage (add/remove) your rootCAs\' config files.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ca.CreateRootCA(privKeySize)
		if err != nil {
			fmt.Printf("%s", helpers.Red("Error while creating the root CA:"))
			fmt.Println(err)
		} else {
			fmt.Printf("A %v bits-keysize certificate %s has been created in %s\n", helpers.Green(strconv.Itoa(privKeySize)), helpers.Green(helpers.CertConfig.CertificateName), helpers.Green(helpers.CertConfig.CertificateDirectory))
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
		err := ca.VerifyCACertificate(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Delete a root CA based on the config file as defined with the -c global flag
var caDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del", "remove"},
	Short:   "Removes a CA certificate",
	Long: `The configuration file describing the certificate should be present.
If not, empty or defaults values will be supplied, and the file will be created`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ca.RemoveCACertificate()
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

	caVerifyCmd.Flags().BoolVarP(&ca.CaVerifyVerbose, "verbose", "v", false, "display the full output")
	caVerifyCmd.Flags().BoolVarP(&ca.CaVerifyComments, "comments", "", false, "display the comments (if any) at the end of the configuration file")
	caCreateCmd.Flags().IntVarP(&privKeySize, "keysize", "b", 4096, "certificate private key size in bits")
}
