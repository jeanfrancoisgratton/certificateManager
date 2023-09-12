// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/signCertificate.go
// Original timestamp: 2023/09/11 10:28

package certs

import (
	"certificateManager/environment"
	"certificateManager/helpers"
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

// signCert: Sign the certificate against the root CA currently held in our custom PKI
// Steps:
// 1. Load the CA cert and private key
// 2. Parse them
// 3. Load & parse the CSR file
// 4. Populate a x509 cert template with the CertificateStruct values
// 5. Create (sign) the certificate
// 6. Save to disk
func (c CertificateStruct) signCert(env environment.EnvironmentStruct) error {
	var caCertBytes, caKeyBytes, csrBytes []byte
	var csrRequest *x509.CertificateRequest
	var caCert *x509.Certificate

	var caKey *rsa.PrivateKey
	var err error

	// Ensure there is a single file in CA/certs directory and fetch its name
	caCertFiles, err := filepath.Glob(filepath.Join(env.CertificateRootDir, env.RootCAdir, "*.crt"))
	if err != nil {
		return CustomError{Message: "Error listing CA certificate files: " + err.Error()}
	}
	if len(caCertFiles) != 1 {
		return CustomError{Message: "Expected one CA certificate file, found " + helpers.Red(string(len(caCertFiles)))}
	}
	baseFN := filepath.Base(caCertFiles[0])

	// 1. Load the CA cert and key files
	if caCertBytes, err = os.ReadFile(filepath.Join(env.CertificateRootDir, env.RootCAdir, baseFN+".crt")); err != nil {
		return CustomError{Message: "Error reading CA certificate: " + err.Error()}
	}
	if caKeyBytes, err = os.ReadFile(filepath.Join(env.CertificateRootDir, env.RootCAdir, baseFN+".key")); err != nil {
		return CustomError{Message: "Error reading CA private key: " + err.Error()}
	}

	// 2. Parse the cert and key files
	if caCert, err = x509.ParseCertificate(caCertBytes); err != nil {
		return err
	}
	if caKey, err = x509.ParsePKCS1PrivateKey(caKeyBytes); err != nil {
		return err
	}

	// 3. Load and parse the CSR file
	if csrBytes, err = os.ReadFile(filepath.Join(env.CertificateRootDir, env.ServerCertsDir, "csr", c.CertificateName+".csr")); err != nil {
		return err
	}
	if csrRequest, err = x509.ParseCertificateRequest(csrBytes); err != nil {
		return nil
	}

	// 4. Populate x509 template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: c.CommonName, Locality: []string{c.Locality}, Country: []string{c.Country}, Organization: []string{c.Organization}, OrganizationalUnit: []string{c.OrganizationalUnit}, Province: []string{c.Province}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(c.Duration, 0, 0),
		KeyUsage:              getKeyUsageFromStrings(c.KeyUsage),
		IsCA:                  c.IsCA,
		BasicConstraintsValid: true,
		DNSNames:              c.DNSNames,
		IPAddresses:           c.IPAddresses,
		EmailAddresses:        c.EmailAddresses,
	}
	if c.IsCA {
		template.KeyUsage = reindexKeyUsage(c)
	}
	// 5. Create (sign) the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, csrRequest.PublicKey, caKey)
	if err != nil {
		return err
	}

	// 6. Encode, save to disk
	certpath := ""
	if c.IsCA {
		certpath = filepath.Join(env.CertificateRootDir, env.RootCAdir, c.CertificateName+".crt")
	} else {
		certpath = filepath.Join(env.CertificateRootDir, env.ServerCertsDir, c.CertificateName+".crt")
	}
	certFile, err := os.Create(certpath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	if err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	//certificate, err := x509.CreateCertificate(rand.Reader, &c, caCert, )
	return nil
}
