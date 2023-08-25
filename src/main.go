package main

import (
	"certificateManager/certs"
	"certificateManager/cmd"
	"certificateManager/environment"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var err error
	// First, we need to create a configuration directory. This is a per-user config dir
	if err = os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager"), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Second, we create the sample environment file with an explanation file
	if err = environment.CreateSampleEnv(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Third, we create the sample certificate config, and its explanation file
	if err = certs.CreateSampleCertificate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// We then launch the command loop
	cmd.Execute()
}
