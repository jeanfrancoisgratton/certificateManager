// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/certificates.go
// Original timestamp: 2023/09/15 13:40

package cmd

import (
	"certificateManager/cert"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	"github.com/spf13/cobra"
	"os"
)

var certCmd = &cobra.Command{
	Use:     "cert",
	Example: "cm cert { {create | delete | verify } } certificate_name | list }",
	Short:   "Certificate sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to specify one of the following subcommand: add | delete | verify | list")
		os.Exit(0)
	},
}

// Lists all certificates in conf directory as per $HOME/.config/JFG/certificatemanager/*.json (defined with -e flag)
var certlistCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "cm cert list",
	Short:   "Lists all certificates in defined rootDir",
	Run: func(cmd *cobra.Command, args []string) {
		var err *cerr.CustomError
		if err = cert.ListCertificates(); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	},
}

// Verify a given certificate
var certVerifyCmd = &cobra.Command{
	Use:     "verify",
	Example: "cm cert verify FILENAME",
	Short:   "Verifies a certificate, as per the provided filename",
	Long:    "Please note: the certificate filename does not need to be within the current PKI structure.",
	Run: func(cmd *cobra.Command, args []string) {
		var err *cerr.CustomError
		if err = cert.Verify(args); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	},
}

// The biggie: create a CA or "normal" SSL certificate
var certCreateCmd = &cobra.Command{
	Use:     "create",
	Example: "cm cert create [CERTICATE_CONFIG_FILE]",
	Short:   "Creates a certificate, specifying (or not) the config file to use",
	Run: func(cmd *cobra.Command, args []string) {
		certname := ""
		if len(args) != 0 {
			certname = args[0]
		}
		var err *cerr.CustomError
		if err = cert.Create(certname); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	},
}

var certRevokeCmd = &cobra.Command{
	Use:     "revoke",
	Aliases: []string{"rm", "remove"},
	Example: "cm cert revoke [CERTICATE_CONFIG_FILE]",
	Short:   "Revokes (deletes) a certificate, specifying the config file to use",
	Run: func(cmd *cobra.Command, args []string) {
		certname := ""
		if len(args) != 0 {
			certname = args[0]
		}
		var err *cerr.CustomError
		if err = cert.RevokeCertificate(certname); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	},
}
