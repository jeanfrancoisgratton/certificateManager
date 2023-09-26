// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/certificates.go
// Original timestamp: 2023/09/15 13:40

package cmd

import (
	"certificateManager/certs"
	"fmt"
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
	Example: "cm cert create [CERTICATE_CONFIG_FILE]",
	Short:   "Creates a certificate, specifying (or not) the config file to use",
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
		if err := certs.Revoke(certname); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	},
}
