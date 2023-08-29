// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/certs/certsHelpers.go
// Original time: 2023/06/16 16:37

package certs

import (
	"certificateManager/environment"
	"encoding/json"
	"io"
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
	Country            string   `json:"Country"`
	Province           string   `json:"Province"`
	Locality           string   `json:"Locality"`
	Organization       string   `json:"Organization"`
	OrganizationalUnit string   `json:"OrganizationalUnit,omitempty"`
	CommonName         string   `json:"CommonName"`
	IsCA               bool     `json:"IsCA,omitempty"`
	EmailAddresses     []string `json:"EmailAddresses,omitempty"`
	Duration           int      `json:"Duration"`
	KeyUsage           []string `json:"KeyUsage"`
	DNSNames           []string `json:"DNSNames,omitempty"`
	IPAddresses        []net.IP `json:"IPAddresses,omitempty"`
	CertificateName    string   `json:"CertificateName"`
	Comments           []string `json:"Comments,omitempty"`
}

// Loads the certificate config from the certificate file
func (c CertificateStruct) LoadCertificateConfFile(certfile string) (CertificateStruct, error) {
	var payload CertificateStruct
	var err error
	var jFile []byte
	var rcFile string
	env := environment.EnvironmentStruct{}

	// fetch environment
	if env, err = environment.EnvironmentStruct.LoadEnvironmentFile(environment.EnvironmentStruct{}); err != nil {
		return CertificateStruct{}, err
	}

	if !strings.HasSuffix(CertConfigFile, ".json") {
		CertConfigFile += ".json"
	}
	if len(certfile) > 0 {
		rcFile = certfile
	} else {
		rcFile = filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir, CertConfigFile)
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
	var env environment.EnvironmentStruct
	var err error

	// fetch environment
	if env, err = environment.EnvironmentStruct.LoadEnvironmentFile(environment.EnvironmentStruct{}); err != nil {
		return err
	}
	if outputfile == "" {
		outputfile = CertConfigFile
	}
	jStream, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	rcFile := filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir, outputfile)
	return os.WriteFile(rcFile, jStream, 0600)
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
		Country:            "CA",
		Province:           "Quebec",
		Locality:           "Blainville",
		Organization:       "myorg.net",
		OrganizationalUnit: "myorg",
		CommonName:         "myorg.net root CA",
		EmailAddresses:     []string{"certs@myorg.net", "certificates@myorg.net"},
		Duration:           10,
		KeyUsage:           []string{"certs sign", "crl sign", "digital signature"},
		DNSNames:           []string{"myorg.net", "myorg.com", "lan.myorg.net"},
		IPAddresses:        []net.IP{net.ParseIP("10.1.1.11"), net.ParseIP("127.0.0.1")},
		CertificateName:    "sample_cert",
		IsCA:               true,
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

// copyFile() : if original file already exists, we copy it onto "destfile"
// Parameters :
// - source file, dest file (string)
// Returns :
// - error code if IO error
func copyFile(source, dest string) error {
	// if the file does not exist, this means we are using a brand-new setup,
	_, err := os.Stat(source)
	if os.IsExist(err) {
		sfile, err := os.Open(source)
		if err != nil {
			return err
		}
		defer sfile.Close()

		dfile, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer dfile.Close()
		_, err = io.Copy(dfile, sfile)
		if err != nil {
			return err
		}
	}
	return nil
}
