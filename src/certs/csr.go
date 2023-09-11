// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/csr.go
// Original timestamp: 2023/09/09 09:25

package certs

import (
	"certificateManager/environment"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"
	"path/filepath"
)

// generateCSR : generate a certificate signing request, and save it to disk
func (c CertificateStruct) generateCSR(env environment.EnvironmentStruct, privateK *rsa.PrivateKey) error {
	var err error
	var csrFile *os.File
	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   c.CommonName,
			Organization: []string{c.Organization},
		},
	}

	certRequest, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privateK)
	if err != nil {
		return err
	}

	if csrFile, err = os.Create(filepath.Join(env.CertificateRootDir, env.ServerCertsDir, "csr", c.CertificateName+".csr")); err != nil {
		return err
	}
	defer csrFile.Close()

	if err = pem.Encode(csrFile, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: certRequest}); err != nil {
		return nil
	}

	return nil
}
