package ca_old

import (
	"cm/helpers-old"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

/*
The createCACert function takes a RootCAconfig struct as input, which contains the following fields:

CommonName: The common name to use for the root certificate.
ValidForYears: The number of years that the root certificate should be valid for.
DNSNames: A list of DNS names to include in the root certificate.
IPAddresses: A list of IP addresses to include in the root certificate.
KeyFilePath: The path to write the private key to.
CertFilePath: The path to write the root certificate to.
The function generates a new RSA private key
*/

// func CreateRootCA(caconfig *config-old.CertConfigStruct) error {
func CreateRootCA(privateKeySize int) error {
	// Generate a new private key for the CA
	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return err
	}

	helpers_old.CertConfig, err = helpers_old.Json2Config()
	if err != nil {
		return err
	}
	// We cannot allow a certificate to last 0yr, 0mt, 0d, so we set a default value of 1 year
	if helpers_old.CertConfig.Duration == 0 {
		helpers_old.CertConfig.Duration = 1
	}
	// Create a new self-signed certificate template
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: helpers_old.CertConfig.CommonName, Locality: []string{helpers_old.CertConfig.Locality}, Country: []string{helpers_old.CertConfig.Country}, Organization: []string{helpers_old.CertConfig.Organization}, OrganizationalUnit: []string{helpers_old.CertConfig.OrganizationalUnit}, Province: []string{helpers_old.CertConfig.Province}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(helpers_old.CertConfig.Duration, 0, 0),
		KeyUsage:              helpers_old.GetKeyUsageFromStrings(helpers_old.CertConfig.KeyUsage),
		IsCA:                  helpers_old.CertConfig.IsCA,
		BasicConstraintsValid: true,
		DNSNames:              helpers_old.CertConfig.DNSNames,
		IPAddresses:           helpers_old.CertConfig.IPAddresses,
		EmailAddresses:        helpers_old.CertConfig.EmailAddresses,
	}
	if helpers_old.CertConfig.IsCA {
		template.KeyUsage = helpers_old.ReindexKeyUsage(helpers_old.CertConfig)
	}
	// Create the certificate using the template and the private key
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// Write the private key and certificate to files
	privateKeyFile, err := os.Create(filepath.Join(helpers_old.CertConfig.CertificateDirectory, helpers_old.CertConfig.CertificateName) + ".key")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return err
	}

	certFile, err := os.Create(filepath.Join(helpers_old.CertConfig.CertificateDirectory, helpers_old.CertConfig.CertificateName) + ".crt")
	if err != nil {
		return err
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		return err
	}

	return nil
}
