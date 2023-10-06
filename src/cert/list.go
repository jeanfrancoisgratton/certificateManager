// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/list.go
// Original timestamp: 2023/08/20 18:43

package cert

import (
	"certificateManager/environment"
	"certificateManager/helpers"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"path/filepath"
	"strings"
)

func ListCertificates() error {
	// Read env file
	var err error
	var fileInfos []os.FileInfo
	env := environment.EnvironmentStruct{}

	// fetch environment
	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}
	// We need to preserve the current Certconfigfile name as it'll get overwritten down here
	oldCfg := CertConfigFile

	// list certificate files
	err = filepath.Walk(filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			fileInfos = append(fileInfos, info)
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("Number of certificates: %s\n", helpers.Green(fmt.Sprintf("%d", len(fileInfos))))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Cert name", "Common Name", "File size", "Modification time"})

	for _, fi := range fileInfos {
		var certDomain string
		CertConfigFile = fi.Name()
		if certDomain, err = fetchCN(filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir, fi.Name())); err != nil {
			CertConfigFile = oldCfg
			return err
		}
		t.AppendRow([]interface{}{helpers.Green(fi.Name()), certDomain, helpers.Green(helpers.SI(uint64(fi.Size()))), helpers.Green(fmt.Sprintf("%v", fi.ModTime().Format("2006/01/02 15:04:05")))})
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

func fetchCN(domain string) (string, error) {
	var domainCert CertificateStruct
	var err error

	if domainCert, err = LoadCertificateConfFile(domain); err != nil {
		return "", err
	}
	return domainCert.CommonName, nil
}
