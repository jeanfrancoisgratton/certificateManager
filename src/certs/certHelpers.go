// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/certs/certHelpers.go
// Original time: 2023/06/16 16:37

package certs

import (
	"certificateManager/environment"
	"crypto/x509"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var CertConfig = CertificateStruct{Duration: 1, KeyUsage: []string{"certs sign", "crl sign", "digital signature"}}
var CertConfigFile = "defaultCertConfig.json"

var CreateSingleCert bool

// This is the full data structure for an SSL certificate and CA
type CertificateStruct struct {
	Country              string   `json:"Country"`
	Province             string   `json:"Province"`
	Locality             string   `json:"Locality"`
	Organization         string   `json:"Organization"`
	OrganizationalUnit   string   `json:"OrganizationalUnit,omitempty"`
	CommonName           string   `json:"CommonName"`
	IsCA                 bool     `json:"IsCA,omitempty"`
	EmailAddresses       []string `json:"EmailAddresses,omitempty"`
	Duration             int      `json:"Duration"`
	KeyUsage             []string `json:"KeyUsage"`
	DNSNames             []string `json:"DNSNames,omitempty"`
	IPAddresses          []net.IP `json:"IPAddresses,omitempty"`
	CertificateDirectory string   `json:"CertificateDirectory"`
	CertificateName      string   `json:"CertificateName"`
	Comments             []string `json:"Comments,omitempty"`
}

// Loads the certificate config from the certificate file
func (c CertificateStruct) LoadCertificateConfFile(certfile string) (CertificateStruct, error) {
	var payload CertificateStruct
	var err error
	var jFile []byte
	var rcFile string

	// some quick, ugly hacks here to allow to load a cert file ad-hoc

	if !strings.HasSuffix(CertConfigFile, ".json") {
		CertConfigFile += ".json"
	}
	if len(certfile) > 0 {
		rcFile = certfile
	} else {
		rcFile = filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", CertConfigFile)
	}

	if jFile, err = os.ReadFile(rcFile); err != nil {
		return CertificateStruct{}, err
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return CertificateStruct{}, err
	} else {
		return payload, nil
	}
}

// Save a data structure into a certificate file in the directory defined in the JSON environment config file
func (c CertificateStruct) SaveCertificateFile(outputfile string) error {
	if outputfile == "" {
		outputfile = CertConfigFile
	}
	jStream, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", outputfile)
	err = os.WriteFile(rcFile, jStream, 0600)

	return err
}

// Dispatch both he Explanation file and the Sample cert
func CreateSampleCertificate() error {
	if err := createSampleCert(); err != nil {
		return err
	}
	if err := createExplanationfile(); err != nil {
		return err
	}
	return nil
}

// Create the sample certificate config file
func createSampleCert() error {
	sampleCertConfig := CertificateStruct{
		Country:              "CA",
		Province:             "Quebec",
		Locality:             "Blainville",
		Organization:         "myorg.net",
		OrganizationalUnit:   "myorg",
		CommonName:           "myorg.net root CA",
		EmailAddresses:       []string{"certs@myorg.net", "certificates@myorg.net"},
		Duration:             10,
		KeyUsage:             []string{"certs sign", "crl sign", "digital signature"},
		DNSNames:             []string{"myorg.net", "myorg.com", "lan.myorg.net"},
		IPAddresses:          []net.IP{net.ParseIP("10.1.1.11"), net.ParseIP("127.0.0.1")},
		CertificateDirectory: "/tmp",
		CertificateName:      "sample_cert",
		IsCA:                 true,
		Comments: []string{"To see which values to put in the Usage field, see https://pkg.go.dev/crypto/x509#KeyUsage",
			"Strip off 'KeyUsage' from the const name and there you go.",
			"",
			"Please note that this field offers no functionality and is strictly here for documentation purposes"},
	}
	if err := sampleCertConfig.SaveCertificateFile("samplee.certconfig.json"); err != nil {
		return err
	}
	return nil
}

// Create an explanation (.txt) file
func createExplanationfile() error {
	expText := `{
	"Country" : "CA", -> This is the country of origin for the certificate
	"Provice" : "Quebec", -> State of province of origin
	"Locality" : "Blainville", -> City of origin
	"Organization" : "myorg.net", -> Organization of origin
	"OrganizationalUnit" : "myorg", -> Sub-organization of origin
	"CommonName" : "myorg.net root CA", -> The name of the certificate
	"EmailAddresses" : ["certs@myorg.net", "certificates@myorg.net"], -> Array of email addresses responsible for this certs
	"Duration" : 10, -> CA duration, in years
	"KeyUsage" : ["Digital Signature", "Certificate Sign", "CRL Sign"], -> Certificate usage. This here are common values for CAs
	"DNSNames" : ["myorg.net","myorg.com","lan.myorg.net"], -> DNS names assigned to this certs
	"IPAddresses" : ["10.1.1.11", "127.0.0.1"], -> IP addresses assigned to this certs (never a good idea to assign IPs to a CA)
	"CertificateDirectory" : "/tmp/", -> directory where to write the certs
	"CertificateName" : "sample_cert", -> certs filename, no extension to the filename
	"IsCA": true, -> Are we creating a CA or a "normal" server certs ?
	"Comments": ["To see which values to put in the KeyUsage field, see https://pkg.go.dev/crypto/x509#KeyUsage", "Strip off 'KeyUsage' from the const name and there you go.", "", "Please note that this field offers no functionality and is strictly here for documentation purposes"] -> Those won't appear in the certificate file
}`

	expFile, err := os.Create(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "sample-README.txt"))
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
	e, err := environment.EnvironmentStruct.LoadEnvironmentFile(environment.EnvironmentStruct{})
	if err != nil {
		return err
	}
	dirRange := []string{filepath.Join(e.CertificateRootDir, e.CertificatesConfigDir),
		filepath.Join(e.CertificateRootDir, e.RootCAdir, "certs"), filepath.Join(e.CertificateRootDir, e.RootCAdir, "newcerts"), filepath.Join(e.CertificateRootDir, e.RootCAdir, "private"),
		filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "private"), filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "csr"), filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "certs"), filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java")}

	for _, directory := range dirRange {
		if err = os.MkdirAll(directory, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// getKeyUsageFromStrings() : converts a slice of strings into
// A x509.KeyUsage value. We use slices of strings because x509.KeyUsage
// Can hold multiple operations at once
func getKeyUsageFromStrings(usageStrings []string) x509.KeyUsage {
	keyUsage := x509.KeyUsage(0)
	for _, usage := range usageStrings {
		switch strings.ToLower(usage) {
		case "digital signature":
			keyUsage |= x509.KeyUsageDigitalSignature
		case "content commitment":
			keyUsage |= x509.KeyUsageContentCommitment
		case "key encipherment":
			keyUsage |= x509.KeyUsageKeyEncipherment
		case "data encipherment":
			keyUsage |= x509.KeyUsageDataEncipherment
		case "key agreement":
			keyUsage |= x509.KeyUsageKeyAgreement
		case "certs sign", "certificate sign":
			keyUsage |= x509.KeyUsageCertSign
		case "crl sign", "crl":
			keyUsage |= x509.KeyUsageCRLSign
		case "encipheronly", "encipher":
			keyUsage |= x509.KeyUsageEncipherOnly
		case "decipheronly", "decipher":
			keyUsage |= x509.KeyUsageDecipherOnly
		}
	}
	return keyUsage
}

// getStringsFromKeyUsage(): takes the x509.KeyUsage numerical value
// And converts it in a slice of human-readable strings,
// As KeyUsage can hold multiple operations at once.
func getStringsFromKeyUsage(keyUsage x509.KeyUsage) []string {
	var usages []string

	if keyUsage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "digital signature")
	}
	if keyUsage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "content commitment")
	}
	if keyUsage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "key encipherment")
	}
	if keyUsage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "data encipherment")
	}
	if keyUsage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "key agreement")
	}
	if keyUsage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "certs sign")
	}
	if keyUsage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "crl sign")
	}
	if keyUsage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "encipher only")
	}
	if keyUsage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "decipher only")
	}
	return usages
}

// reindexKeyUsage() : Ensures that the CertificateStruct.KeyUsage contains only unique values
func reindexKeyUsage(cfg CertificateStruct) x509.KeyUsage {
	org := cfg.KeyUsage
	// We append the CA-related usages
	org = append(org, "certs sign", "crl sign", "digital signature")

	// We map the new slices
	//[]string to map : https://kylewbanks.com/blog/creating-unique-slices-in-go
	s := make([]string, 0, len(org))
	m := make(map[string]bool)

	for _, value := range org {
		if _, ok := m[value]; !ok {
			m[value] = true
			s = append(s, value)
		}
	}
	return getKeyUsageFromStrings(s)
}
