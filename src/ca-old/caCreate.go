package ca_old

import (
	"cm/certs"
	"cm/helpers"
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

// func CreateRootCA(caconfig *config-old.CertificateStruct) error {
func CreateRootCA(privateKeySize int) error {
	// Generate a new private key for the CA
	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		return err
	}

	certs.CertConfig, err = certs.Json2Config()
	if err != nil {
		return err
	}
	// We cannot allow a certificate to last 0yr, 0mt, 0d, so we set a default value of 1 year
	if certs.CertConfig.Duration == 0 {
		certs.CertConfig.Duration = 1
	}
	// Create a new self-signed certificate template
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: certs.CertConfig.CommonName, Locality: []string{certs.CertConfig.Locality}, Country: []string{certs.CertConfig.Country}, Organization: []string{certs.CertConfig.Organization}, OrganizationalUnit: []string{certs.CertConfig.OrganizationalUnit}, Province: []string{certs.CertConfig.Province}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(certs.CertConfig.Duration, 0, 0),
		KeyUsage:              helpers.GetKeyUsageFromStrings(certs.CertConfig.KeyUsage),
		IsCA:                  certs.CertConfig.IsCA,
		BasicConstraintsValid: true,
		DNSNames:              certs.CertConfig.DNSNames,
		IPAddresses:           certs.CertConfig.IPAddresses,
		EmailAddresses:        certs.CertConfig.EmailAddresses,
	}
	if certs.CertConfig.IsCA {
		template.KeyUsage = helpers.ReindexKeyUsage(certs.CertConfig)
	}
	// Create the certificate using the template and the private key
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// Write the private key and certificate to files
	privateKeyFile, err := os.Create(filepath.Join(certs.CertConfig.CertificateDirectory, certs.CertConfig.CertificateName) + ".key")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return err
	}

	certFile, err := os.Create(filepath.Join(certs.CertConfig.CertificateDirectory, certs.CertConfig.CertificateName) + ".crt")
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
