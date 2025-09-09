// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/environment/addRemoveEnv.go
// Original timestamp: 2023/09/15 08:23

package environment

import (
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"os"
	"path/filepath"
	"strings"
)

func RemoveEnvFile(envfile string) *cerr.CustomError {
	if !strings.HasSuffix(envfile, ".json") {
		envfile += ".json"
	}
	if err := os.Remove(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "certificatemanager", envfile)); err != nil {
		return &cerr.CustomError{Title: "Error removing " + envfile, Message: err.Error()}
	}

	fmt.Printf("%s removed succesfully\n", envfile)
	return nil
}

func AddEnvFile(envfile string) *cerr.CustomError {
	if !strings.HasSuffix(envfile, ".json") {
		envfile += ".json"
	}

	if env, err := prompt4EnvironmentValues(); err != nil {
		return err
	} else {
		return env.SaveEnvironmentFile(envfile)
	}
}

func prompt4EnvironmentValues() (EnvironmentStruct, *cerr.CustomError) {
	var env EnvironmentStruct
	fmt.Println("The root dir value should be an absolute path, and all other values relative to it")

	env.CertificateRootDir = hf.GetStringValFromPrompt("Enter the certificate root dir (where the PKI directories will sit): ")
	if !strings.HasPrefix(env.CertificateRootDir, "/") && !strings.HasPrefix(env.CertificateRootDir, "$HOME") && !strings.HasPrefix(env.CertificateRootDir, "~") {
		ce := &cerr.CustomError{Title: "Directory error", Message: fmt.Sprintf("%s is not an absolute path\n", env.CertificatesConfigDir), Fatality: cerr.Fatal}
		return EnvironmentStruct{}, ce
	}

	env.RootCAdir = hf.GetStringValFromPrompt("Enter the rootCA directory name: ")
	if strings.HasPrefix(env.RootCAdir, "/") {
		ce := &cerr.CustomError{Title: "Directory error", Message: fmt.Sprintf("%s must be an absolute path\n", env.RootCAdir), Fatality: cerr.Fatal}
		return EnvironmentStruct{}, ce
	} else {
		env.RootCAdir = filepath.Join(env.CertificateRootDir, env.RootCAdir)
	}

	env.ServerCertsDir = hf.GetStringValFromPrompt("Enter the servers certificate directory name: ")
	if strings.HasPrefix(env.ServerCertsDir, "/") {
		ce := &cerr.CustomError{Title: "Directory error", Message: fmt.Sprintf("%s must be a relative path\n", env.ServerCertsDir), Fatality: cerr.Fatal}
		return EnvironmentStruct{}, ce
	} else {
		env.ServerCertsDir = filepath.Join(env.CertificateRootDir, env.ServerCertsDir)
	}

	env.CertificatesConfigDir = hf.GetStringValFromPrompt("Enter the servers certificates config directory name: ")
	if strings.HasPrefix(env.CertificatesConfigDir, "/") {
		ce := &cerr.CustomError{Title: "Directory error", Message: fmt.Sprintf("%s must be a relative path\n", env.CertificatesConfigDir), Fatality: cerr.Fatal}
		return EnvironmentStruct{}, ce
	} else {
		env.CertificatesConfigDir = filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir)
	}
	env.RemoveDuplicates = true
	return env, nil
}
