// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/config/configCreate.go
// Original time: 2023/06/16 16:53

package config

import (
	"bufio"
	"fmt"
	"os"
)

func ConfCreate() error {
	var err error
	environment := prompt4values()
	environment.Json2EnvironmentFile()
	return err
}

// TODO: add path validation routines
func prompt4values() EnvironmentStruct {
	var env EnvironmentStruct
	inputScanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Please enter the root directory where all files will be hosted (absolute path): ")
	inputScanner.Scan()
	env.CertificateRootDir = inputScanner.Text()

	fmt.Printf("Please enter the rootCA directory (relative path to %s): ", env.CertificateRootDir)
	inputScanner.Scan()
	env.RootCAdir = inputScanner.Text()

	fmt.Printf("Please enter the server certificates directory (relative path to %s): ", env.CertificateRootDir)
	inputScanner.Scan()
	env.ServerCertsDir = inputScanner.Text()

	fmt.Printf("Please enter the server certificates config files directory (relative path to %s): ", env.CertificateRootDir)
	inputScanner.Scan()
	env.CertificatesConfigDir = inputScanner.Text()
	return env
}
