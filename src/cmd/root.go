// certificateManager
// src/cmd/root.go

package cmd

import (
	"certificateManager/cert"
	"certificateManager/environment"
	"fmt"
	"os"
	"runtime"

	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cm",
	Short:   "Certificate / PKI management tool",
	Version: hf.White(fmt.Sprintf("1.61.01-%s 2025.10.15", runtime.GOARCH)),
}

var clCmd = &cobra.Command{
	Use:     "changelog",
	Aliases: []string{"cl"},
	Short:   "Shows changelog",
	Run: func(cmd *cobra.Command, args []string) {
		changelog()
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

func changelog() {
	//fmt.Printf("\x1b[2J")
	fmt.Printf("\x1bc")

	fmt.Println("CHANGELOG")
	fmt.Println()
	fmt.Println()

	fmt.Print(`
VERSION		DATE			COMMENT
-------		----			-------
1.61.01		2025.10.15		GO version bump, doc update
1.61.00		2025.08.25		GO version bump, build deps update
1.60.00		2025.06.06		Updated to GO 1.24.4 to incorporate the crypto/x509 bugfixes
1.59.00		2025.01.13		CN will be using the certificate name if omitted
1.58.00		2024.11.08		Fixed missing config path issue
1.57.00		2024.07.09		Happy birthday Liliane Gratton xxx. Reverted to 1.55.00 and added better error handling
1.56.01		2024.06.28		Attempt at fixing the MacOS limitation issue
1.55.00		2024.06.02		Limit expiration date of certs to 800 days, thank you Apple
1.52.00		2024.05.11		More error handling fixing, GO version bump (1.22.3)
1.51.00		2024.05.10		Fixed error handling where customError was already the return value
1.50.00		2024.05.03		Most helper functions + custom errors moved to outside packages
1.25.01		2024.03.28		Minor: CERTDEFAULTS and CADEFAULTS can be in lowercase
1.25.00		2024.03.28		Moved all configs from $HOME/.config/ to $HOME/.config/JFG/
1.24.01		2024.03.28		UNPUBLISHED: Minor output fixes, version now shows GOARCH, GO version bump, dependencies updates
1.24.00		2024.03.01		GO version bump
1.23.00		2023.12.07		GO version bump to 1.21.5, packages upgrades
1.22.00		2023.11.02		cm cert rm was broken (duplicated path)
1.21.00		2023.11.02		defaultEnv.json does not need to be specified anymore if it's the environment we use
1.20.06		2023.10.31		more explicit error message in cert verify
1.20.05		2023.10.16		moved to a saner version numbering scheme: 1.20.05 is 1.205, actually
1.205		2023.10.15		fixed wrong path for server's private keys
1.200		2023.10.13		go version bump, folded all environment directories into a single var for readability issues with filepath.Join()
1.100-0		2023.10.04		cm cert verify now works w/ the verbose flag; duplicate certificate creation is now prevented 
1.010-1		2023.10.03		Fixed typo in directory name
1.010		2023.10.03		Fixed issue when serial number was not incremented within certificate
1.001		2023.10.03		Minor changes: verbosity, doc update
1.000		2023.09.30		Completed prod-ready version
0.500		2023.06.03		server cert management
0.400		2023.04.22		config-old management
0.300		2023.04.20		ca-old edit, ca-old del
0.200		2023.04.20		ca-old create and ca-old verify
0.100		2023.04.16		near-config-old-aware
\n`)
}
