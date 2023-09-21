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
	"strings"
	"time"
)

// signCert: Sign the certificate against the root CA currently held in our custom PKI
// Steps:
// 0. Ensure that there is a single CA (.crt) file in the CA directory
// 1. Load the CA cert and private key
// 2. Decode and parse them
// 3. Load & parse the CSR file
// 4. Populate a x509 cert template with the CertificateStruct values
// 5. Sign (create) the certificate
// 6. Save to disk
func (c CertificateStruct) signCert(env environment.EnvironmentStruct) error {
	var csrBytes, caCertPEM, caKeyPEM []byte
	var csrRequest *x509.CertificateRequest
	var caCert *x509.Certificate
	var caKey *rsa.PrivateKey
	var err error

	// Ensure there is a single file in the CA directory and fetch its name
	caCertFiles, err := filepath.Glob(filepath.Join(env.CertificateRootDir, env.RootCAdir, "*.crt"))
	if err != nil {
		return helpers.CustomError{Message: "Error listing CA certificate files: " + err.Error()}
	}
	if len(caCertFiles) != 1 {
		return helpers.CustomError{Message: "Expected one CA certificate file, found " + helpers.Red(string(len(caCertFiles)))}
	}
	baseFN := strings.TrimSuffix(filepath.Base(caCertFiles[0]), filepath.Ext(filepath.Base(caCertFiles[0])))

	// 1. Load the CA cert and key files
	if caCertPEM, err = os.ReadFile(filepath.Join(env.CertificateRootDir, env.RootCAdir, baseFN+".crt")); err != nil {
		return helpers.CustomError{Message: "Error reading CA certificate: " + err.Error()}
	}
	if caKeyPEM, err = os.ReadFile(filepath.Join(env.CertificateRootDir, env.RootCAdir, baseFN+".key")); err != nil {
		return helpers.CustomError{Message: "Error reading CA private key: " + err.Error()}
	}

	// 2. Parse the cert and key files
	caCertBlock, _ := pem.Decode(caCertPEM)
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caCertBlock == nil || caKeyBlock == nil {
		return helpers.CustomError{Message: "Error PEM-decoding the CA certificate or its private key"}
	}
	if caCert, err = x509.ParseCertificate(caCertBlock.Bytes); err != nil {
		return err
	}
	if caKey, err = x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes); err != nil {
		return err
	}

	// 3. Load, decode and parse the CSR file
	if csrBytes, err = os.ReadFile(filepath.Join(env.CertificateRootDir, env.ServerCertsDir, "csr", c.CertificateName+".csr")); err != nil {
		return err
	}
	if csrBlock, _ := pem.Decode(csrBytes); csrBlock == nil {
		return helpers.CustomError{"Error PEM-decoding the CSR file"}
	} else {
		if csrRequest, err = x509.ParseCertificateRequest(csrBlock.Bytes); err != nil {
			return err
		}
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
	//// why the next ???
	//if c.IsCA {
	//	template.KeyUsage = reindexKeyUsage(c)
	//}

	// 5. Create (sign) the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, csrRequest.PublicKey, caKey)
	if err != nil {
		return err
	}

	// 6. Encode, save to disk
	certFile, err := os.Create(filepath.Join(env.CertificateRootDir, env.ServerCertsDir, "certs", c.CertificateName+".crt"))
	if err != nil {
		return err
	}
	defer certFile.Close()

	if err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	return nil
}

// createCA and signCert are very similar: one is for non-CA certs, the other (below) for CA certs
// I *could* fold both into a single function, with tons of "if c.IsCA{}" clauses, but it's not worth
// the readability headache that it'd bring
func (c CertificateStruct) createCA(env environment.EnvironmentStruct, privateKey *rsa.PrivateKey) error {
	var cabytes []byte
	var err error

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
	if cabytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey); err != nil {
		return err
	}

	cafile, err := os.Create(filepath.Join(env.CertificateRootDir, env.RootCAdir, c.CertificateName+".crt"))
	if err != nil {
		return err
	}
	defer cafile.Close()

	if err = pem.Encode(cafile, &pem.Block{Type: "CERTIFICATE", Bytes: cabytes}); err != nil {
		return err
	}
	return nil
}
