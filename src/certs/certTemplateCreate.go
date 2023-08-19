// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/config-old/configTemplateCreate.go
// Original time: 2023/06/03 07:01

package certs

import (
	"net"
	"os"
	"path/filepath"
)

func CreateSampleCert() error {
	if err := createExplanationfile(); err != nil {
		return err
	}
	if err := createSampleCert(); err != nil {
		return err
	}
	return nil
}

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
	if err := sampleCertConfig.CertConfig2Json("certificateSample.json"); err != nil {
		return err
	}
	return nil
}
