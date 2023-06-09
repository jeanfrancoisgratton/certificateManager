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

// caCmd represents the cert command
var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Server certificates management",
	Long:  `This is where you will manage (add/verify/delete) your server certificates.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage ca {create|verify")
			os.Exit(0)
		}
	},
}

// Create a server certificate based on the config file as defined with the -c global flag
var certCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Server certificate creation command",
	Run: func(cmd *cobra.Command, args []string) {
		err := ca.CreateRootCA(privKeySize)
		if err != nil {
			fmt.Printf("%s", helpers.Red("Error while creating the certificate:"))
			fmt.Println(err)
		} else {
			fmt.Printf("A %v bits-keysize certificate %s has been created in %s\n", helpers.Green(strconv.Itoa(privKeySize)), helpers.Green(helpers.CertConfig.CertificateName), helpers.Green(helpers.CertConfig.CertificateDirectory))
		}
	},
}

// Verify a server certificate
var certVerifyCmd = &cobra.Command{
	Use:   "verify certificate_filename",
	Short: "Verify the created CA certificate",
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

// Delete a cert based on the config file as defined with the -c global flag
var certDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del", "remove"},
	Short:   "Removes a certificate",
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
	rootCmd.AddCommand(certCmd)
	certCmd.AddCommand(certCreateCmd)
	certCmd.AddCommand(certVerifyCmd)
	certCmd.AddCommand(certDeleteCmd)

	certVerifyCmd.Flags().BoolVarP(&ca.CaVerifyVerbose, "verbose", "v", false, "display the full output")
	certVerifyCmd.Flags().BoolVarP(&ca.CaVerifyComments, "comments", "", false, "display the comments (if any) at the end of the configuration file")
	certCreateCmd.Flags().IntVarP(&privKeySize, "keysize", "b", 4096, "certificate private key size in bits")
}
