// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/certs/certHelpers.go
// Original time: 2023/06/16 16:37

package certs

import (
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
func LoadCertificateFile() (CertificateStruct, error) {
	var payload CertificateStruct
	var err error

	if !strings.HasSuffix(CertConfigFile, ".json") {
		CertConfigFile += ".json"
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", CertConfigFile)
	jFile, err := os.ReadFile(rcFile)
	if err != nil {
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
	var sampleCertConfig = CertificateStruct{
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
	if err := sampleCertConfig.SaveCertificateFile("certificateSample.json"); err != nil {
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

	expFile, err := os.Create(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "certificateSample-README.txt"))
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
