// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/certs/helpers.go
// Original time: 2023/06/16 16:37

package certs

import (
	"certificateManager/environment"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// var CertConfig = CertificateStruct{Duration: 1, KeyUsage: []string{"certs sign", "crl sign", "digital signature"}}
var CertConfigFile = "defaultCertConfig.json"
var CertName = ""
var CertJava = false
var CertRemoveFiles = false

// This is the full data structure for an SSL certificate and CA
type CertificateStruct struct {
	Country            string   `json:"Country"`
	Province           string   `json:"Province"`
	Locality           string   `json:"Locality"`
	Organization       string   `json:"Organization"`
	OrganizationalUnit string   `json:"OrganizationalUnit,omitempty"`
	CommonName         string   `json:"CommonName"`
	IsCA               bool     `json:"IsCA"`
	EmailAddresses     []string `json:"EmailAddresses,omitempty"`
	Duration           int      `json:"Duration"`
	KeyUsage           []string `json:"KeyUsage"`
	DNSNames           []string `json:"DNSNames,omitempty"`
	IPAddresses        []net.IP `json:"IPAddresses,omitempty"`
	CertificateName    string   `json:"CertificateName"`
	SerialNumber       uint64   `json:"SerialNumber"`
	Comments           []string `json:"Comments,omitempty"`
}

// Create the sample certificate config file
func CreateSampleCert() error {
	sampleCertConfig := CertificateStruct{
		Country:            "CA",
		Province:           "Quebec",
		Locality:           "Blainville",
		Organization:       "myorg.net",
		OrganizationalUnit: "myorg",
		CommonName:         "myorg.net root CA",
		EmailAddresses:     []string{"certs@myorg.net", "certs@org,net"},
		Duration:           10,
		KeyUsage:           []string{"certs sign", "crl sign", "digital signature"},
		DNSNames:           []string{"myorg.net", "myorg.com", "lan.myorg.net"},
		IPAddresses:        []net.IP{net.ParseIP("10.0.0.1"), net.ParseIP("127.0.0.1")},
		CertificateName:    "sampleCert",
		IsCA:               true,
		SerialNumber:       1,
		Comments: []string{"To see which values to put in the KeyUsage field, see https://pkg.go.dev/crypto/x509#KeyUsage",
			"Strip off 'KeyUsage' from the const name and there you go.",
			"",
			"Please note that this field offers no functionality and is strictly here for documentation purposes"},
	}
	if !strings.HasSuffix(sampleCertConfig.CertificateName, ".json") {
		sampleCertConfig.CertificateName += ".json"
	}
	if err := sampleCertConfig.SaveCertificateConfFile(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", sampleCertConfig.CertificateName)); err != nil {
		return err
	}
	return nil
}

// Create an explanation (.txt) file
func CreateExplanationfile() error {
	expText := `{
	"Country" : "CA", -> This is the country of origin for the certificate
	"Provice" : "Quebec", -> State of province of origin
	"Locality" : "Blainville", -> City of origin
	"Organization" : "myorg.net", -> Organization of origin
	"OrganizationalUnit" : "myorg", -> Sub-organization of origin
	"CommonName" : "myorg.net root CA", -> The name of the certificate
	"EmailAddresses" : ["certs@myorg.net", "certs@org.net"], -> Email addresses responsible for this cert
	"Duration" : 10, -> CA duration, in years
	"KeyUsage" : ["Digital Signature", "Certificate Sign", "CRL Sign"], -> Certificate usage. This here are common values for CAs
	"DNSNames" : ["myorg.net","myorg.com","lan.myorg.net"], -> DNS names assigned to this certs
	"IPAddresses" : ["10.1.1.11", "127.0.0.1"], -> IP addresses assigned to this certs (never a good idea to assign IPs to a CA)
	"CertificateName" : "sample_cert", -> certs filename, no extension to the filename
	"IsCA": true, -> Are we creating a CA or a "normal" server certs ?
	"SerialNumber": this is an unsigned int64, handled by the software; put here any positive value
	"Comments": ["To see which values to put in the KeyUsage field, see https://pkg.go.dev/crypto/x509#KeyUsage", "Strip off 'KeyUsage' from the const name and there you go.", "", "Please note that this field offers no functionality and is strictly here for documentation purposes"] -> Those won't appear in the certificate file
}`

	expFile, err := os.Create(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "sampleCert-README.txt"))
	if err != nil {
		return err
	}
	defer expFile.Close()

	_, err = expFile.WriteString(expText)
	if err != nil {
		return err
	}

	return nil
}

// createCertificateRootDirectories() : creates the directory structure needed to store all certs, keys, CA, CSR, etc
func createCertificateRootDirectories() error {
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}
	dirRange := []string{filepath.Join(e.CertificateRootDir, e.CertificatesConfigDir),
		filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "private"), filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "csr"),
		filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "certs"), filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java")}

	for _, directory := range dirRange {
		if err = os.MkdirAll(directory, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
