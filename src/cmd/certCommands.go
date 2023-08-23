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
