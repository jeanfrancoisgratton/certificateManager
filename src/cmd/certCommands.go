// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/certCommands.go
// Original timestamp: 2023/08/20 18:40

package cmd

import (
	"certificateManager/certs"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

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

var certCreateCmd = &cobra.Command{
	Use: "create",
	//Aliases: []string{"ls"},
	Example: "cm cert create CONFIGFILE",
	Short:   "Creates a certificate, specifying the config file to use",
	Run: func(cmd *cobra.Command, args []string) {
		if err := certs.Create(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	},
}
