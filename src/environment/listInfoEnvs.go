// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/environment/listInfoEnvs.go
// Original timestamp: 2023/09/13 16:01

package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cerr "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ListEnvironments(envdir string) *cerr.CustomError {
	var err error
	var dirFH *os.File
	var finfo, fileInfos []os.FileInfo

	// list environment files
	if envdir == "" {
		envdir = filepath.Join(os.Getenv("HOME"), ".config", "JFG", "certificatemanager")
	}
	if dirFH, err = os.Open(envdir); err != nil {
		ce := &cerr.CustomError{Title: "Unable to read config directory"}
		return ce
	}

	if fileInfos, err = dirFH.Readdir(0); err != nil {
		ce := &cerr.CustomError{Title: "Unable to read files in config directory"}
		return ce
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") && !strings.HasPrefix(info.Name(), "sample") {
			finfo = append(finfo, info)
		}
	}

	if err != nil {
		ce := &cerr.CustomError{Title: "Undefined error", Message: err.Error()}
		return ce
	}

	fmt.Printf("Number of environment files: %s\n", hftx.Green(fmt.Sprintf("%d", len(finfo))))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Environment file", "File size", "Modification time"})

	for _, fi := range finfo {
		t.AppendRow([]interface{}{hftx.Green(fi.Name()), hftx.Green(hf.SI(uint64(fi.Size()))),
			hftx.Green(fmt.Sprintf("%v", fi.ModTime().Format("2006/01/02 15:04:05")))})
	}
	t.SortBy([]table.SortBy{
		{Name: "Environment file", Mode: table.Asc},
		{Name: "File size", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}

func ExplainEnvFile(envfiles []string) *cerr.CustomError {
	oldEnvFile := EnvConfigFile

	fmt.Println("Paths are relative to Certificate root dir's path")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Environment file", "Certificate root dir", "CA dir", "Server certificates dir", "Certificates config dir"})

	for _, envfile := range envfiles {
		if !strings.HasSuffix(envfile, ".json") {
			envfile += ".json"
		}
		EnvConfigFile = envfile

		if e, err := LoadEnvironmentFile(); err != nil {
			EnvConfigFile = oldEnvFile
			return err
		} else {
			t.AppendRow([]interface{}{hftx.Green(envfile), hftx.Green(e.CertificateRootDir), hftx.Green(filepath.Base(e.RootCAdir)),
				hftx.Green(filepath.Base(e.ServerCertsDir)), hftx.Green(filepath.Base(e.CertificatesConfigDir))})
		}

	}
	t.SortBy([]table.SortBy{
		{Name: "Environment file", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	EnvConfigFile = oldEnvFile
	return nil
}
