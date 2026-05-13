// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/list.go
// Original timestamp: 2023/08/20 18:43

package cert

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"certificateManager/environment"
	cerr "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ListCertificates() *cerr.CustomError {
	// Read env file
	var err error
	var cErr *cerr.CustomError
	var fileInfos []os.FileInfo
	env := environment.EnvironmentStruct{}

	// fetch environment
	if env, cErr = environment.LoadEnvironmentFile(); cErr != nil {
		return cErr
	}
	// We need to preserve the current Certconfigfile name as it'll get overwritten down here
	oldCfg := CertConfigFile

	// list certificate files
	err = filepath.Walk(env.CertificatesConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return &cerr.CustomError{Message: err.Error()}
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			fileInfos = append(fileInfos, info)
		}
		return nil
	})

	if err != nil {
		return &cerr.CustomError{Message: err.Error()}
	}

	fmt.Printf("Number of certificates: %s\n", hftx.Green(fmt.Sprintf("%d", len(fileInfos))))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Cert name", "Common Name", "File size", "Modification time"})

	for _, fi := range fileInfos {
		var certDomain string
		var cErr *cerr.CustomError
		CertConfigFile = fi.Name()
		if certDomain, cErr = fetchCN(filepath.Join(env.CertificatesConfigDir, fi.Name())); cErr != nil {
			CertConfigFile = oldCfg
			return cErr
		}
		t.AppendRow([]interface{}{hftx.Green(fi.Name()), certDomain, hftx.Green(hf.SI(uint64(fi.Size()))), hftx.Green(fmt.Sprintf("%v", fi.ModTime().Format("2006/01/02 15:04:05")))})
	}
	t.SortBy([]table.SortBy{
		{Name: "Cert name", Mode: table.Asc},
		{Name: "File size", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}

func fetchCN(domain string) (string, *cerr.CustomError) {
	var domainCert CertificateStruct
	var err *cerr.CustomError

	if domainCert, err = LoadCertificateConfFile(domain); err != nil {
		return "", err
	}
	return domainCert.CommonName, nil
}
