/*
©2023 J.F.Gratton (jean-francois@famillegratton.net)
*/
package main

import (
	"cm/certs"
	"cm/cmd"
	"cm/config"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// First we ensure that the user has a config directory
	if err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager"), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// Second, we create a configuration file if none exists
	_, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "certificatemanagerConf.json"))
	if err != nil {
		config.ConfCreate()
	}
	// Third, we create a sample certificate in that config directory
	if err := certs.TemplateConfigCreate(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	cmd.Execute()
}
