// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/cert/helpers.go
// Original time: 2023/06/16 16:37

package cert

import (
	"bufio"
	"certificateManager/environment"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var CertConfigFile = "defaultCertConfig.json"
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
		EmailAddresses:     []string{"cert@myorg.net", "cert@org,net"},
		Duration:           10,
		KeyUsage:           []string{"cert sign", "crl sign", "digital signature"},
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
	"EmailAddresses" : ["cert@myorg.net", "cert@org.net"], -> Email addresses responsible for this cert
	"Duration" : 10, -> CA duration, in years
	"KeyUsage" : ["Digital Signature", "Certificate Sign", "CRL Sign"], -> Certificate usage. This here are common values for CAs
	"DNSNames" : ["myorg.net","myorg.com","lan.myorg.net"], -> DNS names assigned to this cert
	"IPAddresses" : ["10.1.1.11", "127.0.0.1"], -> IP addresses assigned to this cert (never a good idea to assign IPs to a CA)
	"CertificateName" : "sample_cert", -> cert filename, no extension to the filename
	"IsCA": true, -> Are we creating a CA or a "normal" server cert ?
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

// createCertificateRootDirectories() : creates the directory structure needed to store all cert, keys, CA, CSR, etc
func createCertificateRootDirectories() error {
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}
	dirRange := []string{e.RootCAdir, e.CertificatesConfigDir, filepath.Join(e.ServerCertsDir, "private"), filepath.Join(e.ServerCertsDir, "csr"),
		filepath.Join(e.ServerCertsDir, "certs"), filepath.Join(e.ServerCertsDir, "java")}

	for _, directory := range dirRange {
		if err = os.MkdirAll(directory, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// check4DuplicateCert :
// We browse the index.txt database to see if the signature of thecertificate we are creating already exists
func (c CertificateStruct) check4DuplicateCert(ndxFilePath string) (bool, error) {
	lookoutstring := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s", c.Country, c.Province, c.Locality,
		c.Organization, c.OrganizationalUnit, c.CommonName)

	isDupe := false

	// open index.txt db
	indexfileHandle, err := os.Open(ndxFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer indexfileHandle.Close()

	// Scan the file to find a duplicate cert signature
	scanner := bufio.NewScanner(indexfileHandle)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, lookoutstring) && strings.HasPrefix(line, "V") {
			isDupe = true
			break
		}
	}
	return isDupe, nil
}
