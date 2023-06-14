/*
©2023 J.F.Gratton (jean-francois@famillegratton.net)
*/
package main

import (
	"cm/cmd"
	"cm/config-old"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager"), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := config_old.TemplateConfigCreate(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	cmd.Execute()
}
