// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/environment/addRemoveEnv.go
// Original timestamp: 2023/09/15 08:23

package environment

import (
	"certificateManager/helpers"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RemoveEnvFile(envfile string) error {
	if !strings.HasSuffix(envfile, ".json") {
		envfile += ".json"
	}
	if err := os.Remove(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", envfile)); err != nil {
		return err
	}

	fmt.Printf("%s removed succesfully\n", envfile)
	return nil
}

func AddEnvFile(envfile string) error {
	var env EnvironmentStruct
	var err error

	if !strings.HasSuffix(envfile, ".json") {
		envfile += ".json"
	}

	if env, err = prompt4EnvironmentValues(); err != nil {
		return err
	} else {
		err = env.SaveEnvironmentFile(envfile)
	}
	return err
}

func prompt4EnvironmentValues() (EnvironmentStruct, error) {
	var env EnvironmentStruct
	fmt.Println("The root dir value should be an absolute path, and all other values relative to it")
	env.CertificateRootDir = helpers.GetStringValFromPrompt("Enter the certificate root dir (where the PKI directories will sit): ")
	if !strings.HasPrefix(env.CertificateRootDir, "/") && !strings.HasPrefix(env.CertificateRootDir, "$HOME") && !strings.HasPrefix(env.CertificateRootDir, "~") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.CertificatesConfigDir, helpers.Red("is not an absolute path"))}
	}
	env.RootCAdir = helpers.GetStringValFromPrompt("Enter the rootCA directory name: ")
	if strings.HasPrefix(env.RootCAdir, "/") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.RootCAdir, helpers.Red("must be an absolute path"))}
	}
	env.ServerCertsDir = helpers.GetStringValFromPrompt("Enter the servers certificate directory name: ")
	if strings.HasPrefix(env.RootCAdir, "/") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.ServerCertsDir, helpers.Red("must be an absolute path"))}
	}
	env.CertificatesConfigDir = helpers.GetStringValFromPrompt("Enter the servers certificates config directory name: ")
	if strings.HasPrefix(env.RootCAdir, "/") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.CertificatesConfigDir, helpers.Red("must be an absolute path"))}
	}
	env.RemoveDuplicates = true
	return env, nil
}
