package main

import (
	"certificateManager/cmd"
	"certificateManager/environment"
	"fmt"
	"os"
	"path/filepath"
)

//func main() {
//	// First we ensure that the user has a config directory
//	if err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager"), os.ModePerm); err != nil {
//		fmt.Println(err)
//		os.Exit(-1)
//	}
//	// Second, we create a sample certificate file if none exists
//	_, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "certificatemanagerConf.json"))
//	if err != nil {
//		config.CreateTemplatedCertificate()
//	}
//	// Third, we create a sample certificate in that config directory
//	if err := certs.TemplateConfigCreate(); err != nil {
//		fmt.Println(err)
//		os.Exit(-1)
//	}
//	cmd.Execute()
//}

func main() {
	// First, we need to create a configuration directory. This is a per-user config dir
	if err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager"), os.FileMode(0644)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Second, we create the sample environment file with an explanation file
	if err := environment.CreateSampleEnv(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := environment.CreateExplanationFile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Third, we create the sample certificate config, and its explanation file

	// We then launch the command loop
	cmd.Execute()
}
