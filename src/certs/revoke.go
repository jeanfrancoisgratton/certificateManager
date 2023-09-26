// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/revoke.go
// Original timestamp: 2023/09/25 18:54

package certs

import (
	"certificateManager/environment"
	"fmt"
	"path/filepath"
	"strings"
)

// Revoke:
// We do not implement a full CRL-CDP mechanism to revoke a certificate, as this tool is mainly intended
// For local infrasctures, local PKIs
// Instead, we proceed this way:
// 1. The certificate is removed from index.txt, within the PKI. The cert is looked according to its signature (CN, O, OU, etc)
// from its config file
// 2. The certificate is removed from the PKI's rootCA/newcerts directory
// 3. Optionally, the CSR, private key and certificate and config file are also removed.
// The default setting is to leave them

func Revoke(certname string) error {
	var err error
	e := environment.EnvironmentStruct{}
	c := CertificateStruct{}
	targetString := ""

	if e, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}

	if !strings.HasSuffix(certname, ".json") {
		certname += ".json"
	}

	if c, err = LoadCertificateConfFile(filepath.Join(e.CertificateRootDir, e.CertificatesConfigDir, certname)); err != nil {
		return err
	}

	if err = removeFromIndex(c.Country, c.Province, c.Locality, c.Organization,
		c.OrganizationalUnit, c.CommonName, c.EmailAddresses); err != nil {
		return err
	}

	return nil
}

func removeFromIndex(country string, province string, locality string, org string, ou string,
	cn string, email []string) error {

	targetString := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s/emailAddress=%s",
		country, province, locality, org, ou, cn, email)

	return nil
}
