// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/revoke.go
// Original timestamp: 2023/09/25 18:54

package certs

import (
	"certificateManager/environment"
	"certificateManager/helpers"
	"fmt"
	"path/filepath"
	"strings"
	"time"
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
	targetLine := fmt.Sprintf("V\t%sZ\t", time.Now().UTC().Format("060102150405"))

	if e, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}

	if !strings.HasSuffix(certname, ".json") {
		certname += ".json"
	}

	if c, err = LoadCertificateConfFile(filepath.Join(e.CertificateRootDir, e.CertificatesConfigDir, certname)); err != nil {
		return err
	}
	//targetLine += fmt.Sprintf("%s\tunknown\t", fm)
	if err = removeFromIndex(c.Country, c.Province, c.Locality, c.Organization,
		c.OrganizationalUnit, c.CommonName, c.EmailAddresses); err != nil {
		return helpers.CustomError{Message: "Unable to remove entry from index.txt: " + err.Error()}
	}

	return nil
}

// removeFromIndex: remove entry from index.txt, and unlink (delete) the certificate from newcerts/
// We scan the index.txt file, tracking an entry that contains "targetString"
// If the entry is found, we extract the serial number (the 4 digit hex number, 3rd field of the line)
// We then rewrite index.txt without that entry
// If found, we then unlink newcerts/$SERIAL_NUM.pem
func removeFromIndex(country string, province string, locality string, org string, ou string,
	cn string, email []string) error {
	indexValue := ""

	targetString := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s/emailAddress=%s",
		country, province, locality, org, ou, cn, email)

	return nil
}

// ALTERNATE APPROACH : RECOPY LINE, WITH R PREFIX INSTEAD OF V
// WE WOULD NOT NEED TO UNLINK FILE FROM newcerts/ and srv/{certs,private,csr}
