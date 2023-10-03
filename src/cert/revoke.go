// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/revoke.go
// Original timestamp: 2023/09/25 18:54

package certs

import (
	"bufio"
	"certificateManager/environment"
	"certificateManager/helpers"
	"fmt"
	"os"
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

	if e, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}

	if !strings.HasSuffix(certname, ".json") {
		certname += ".json"
	}

	// We need to load the cert's config in order to find the info needed to remove it from the index DB
	if c, err = LoadCertificateConfFile(filepath.Join(e.CertificateRootDir, e.CertificatesConfigDir, certname)); err != nil {
		return err
	}

	if err = putRevokeFlag(e, certname, c.Country, c.Province, c.Locality, c.Organization,
		c.OrganizationalUnit, c.CommonName); err != nil {
		return err
	}

	fmt.Printf("Certificate %s has sucessfully been %s\n", c.CertificateName, helpers.Green("revoked"))
	return nil
}

// putRevokeFlag: remove entry from index.txt, and unlink (delete) the certificate from newcerts/
// We scan the index.txt file, tracking an entry that contains "targetString"
// If the entry is found, we extract the serial number (the 4 digit hex number, 3rd field of the line)
// We then rewrite index.txt without that entry
// If found, we then unlink newcerts/$SERIAL_NUM.pem
func putRevokeFlag(e environment.EnvironmentStruct, certname string, country string, province string,
	locality string, org string, ou string, cn string) error {
	serialField := ""
	certfilename := ""

	// Remove extension from filename
	dotPos := strings.LastIndex(certname, ".")
	if dotPos >= 0 {
		certfilename = certname[:dotPos]
	} else {
		certfilename = certname
	}

	// This is the substring we use to track the correct certificate
	targetString := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s", country, province, locality, org, ou, cn)

	// Open in and out files
	inFile, err := os.Open(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt"))
	if err != nil {
		return helpers.CustomError{Message: "Unable to open index.txt: " + err.Error()}
	}
	defer inFile.Close()
	outFile, err := os.Create(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.new"))
	if err != nil {
		return helpers.CustomError{Message: "Unable to create temp index file: " + err.Error()}
	}
	defer outFile.Close()

	// Now we scan the infile to find the substring in parameters
	scanner := bufio.NewScanner(inFile)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Check if the line contains the search string
		if strings.Contains(line, targetString) {
			// Extract the 3rd field (certificate number)
			if len(fields) >= 3 {
				fields[0] = "R"
				serialField = fields[2]

				// Write the modified line to the output file
				newLine := strings.Join(fields, "\t") + "\n"
				outFile.WriteString(newLine)
			}
		} else {
			// Write unmodified lines to the output file
			outFile.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return helpers.CustomError{Message: "Unable to read index.txt: " + err.Error()}
	}

	if CertRemoveFiles {
		os.Remove(filepath.Join(e.CertificateRootDir, e.RootCAdir, "newcerts", strings.ToUpper(serialField)+".pem"))
		os.Remove(filepath.Join(e.CertificateRootDir, e.CertificatesConfigDir, certfilename+".json"))
		os.Remove(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "cert", certfilename+".crt"))
		os.Remove(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "csr", certfilename+".csr"))
		os.Remove(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "private", certfilename+".key"))
		if _, err = os.Stat(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java", certfilename+".p12")); err != nil && os.IsNotExist(err) {
			// nop
		} else {
			os.Remove(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java", certfilename+".p12"))
		}
		if _, err = os.Stat(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java", certfilename+".jks")); err != nil && os.IsNotExist(err) {
			// nop
		} else {
			os.Remove(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java", certfilename+".jks"))
		}
	}

	// Rename new file to index.txt
	os.Rename(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.new"), filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt"))

	return nil
}
