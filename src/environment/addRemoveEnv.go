// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/environment/addRemoveEnv.go
// Original timestamp: 2023/09/15 08:23

package environment

import (
	"certificateManager/helpers"
	"fmt"
	"os"
	"strings"
)

func RemoveEnvFile(envfile string) error {
	if err := os.Remove(envfile); err != nil {
		return err
	}

	return nil
}

func AddEnvFile(envfile string) error {
	var env EnvironmentStruct
	var err error
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
	helpers.GetStringValFromPrompt("Enter the certificate root dir (where the PKI directories will sit): ", &env.CertificateRootDir)
	if !strings.HasPrefix(env.CertificateRootDir, "/") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.CertificatesConfigDir, helpers.Red("is not an absolute path"))}
	}
	helpers.GetStringValFromPrompt("Enter the rootCA directory name: ", &env.RootCAdir)
	if strings.HasPrefix(env.RootCAdir, "/") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.RootCAdir, helpers.Red("must be an absolute path"))}
	}
	helpers.GetStringValFromPrompt("Enter the servers certificate directory name: ", &env.ServerCertsDir)
	if strings.HasPrefix(env.RootCAdir, "/") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.ServerCertsDir, helpers.Red("must be an absolute path"))}
	}
	helpers.GetStringValFromPrompt("Enter the servers certificates config directory name: ", &env.CertificatesConfigDir)
	if strings.HasPrefix(env.RootCAdir, "/") {
		return EnvironmentStruct{}, helpers.CustomError{Message: fmt.Sprintf("%s %s\n", env.CertificatesConfigDir, helpers.Red("must be an absolute path"))}
	}
	env.RemoveDuplicates = true
	return EnvironmentStruct{}, nil
}
