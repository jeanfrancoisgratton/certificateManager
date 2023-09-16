// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/certPrivateKeys.go
// Original timestamp: 2023/08/25 18:58

package certs

import (
	"certificateManager/environment"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

// createPrivateKey : creates either a CA private key, or a standard cert private key
// Parameters:
// - filename (string): the name of the certificate appended with ".key"
// - pkrootdir (string) : corresponds to CertificateRootDir + RootCAdir + "private/" + filename + ".key" for a CA, or
// - CertificateRootDir + ServerCertsDir + "private/ + filename + ".key" for a standard cert
// Returns:
// - A pointer to the private key
// - the error code, if any
func (c CertificateStruct) createPrivateKey() (*rsa.PrivateKey, error) {
	var pk *rsa.PrivateKey
	var err error = nil
	var pkFile *os.File
	var pkfile string
	var env environment.EnvironmentStruct

	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return nil, err
	}

	if pk, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return nil, err
	}

	// rootCA keys are not stored at the same place as other SSL keys
	if c.IsCA {
		pkfile = filepath.Join(env.CertificateRootDir, env.RootCAdir, c.CertificateName+".key")
	} else {
		pkfile = filepath.Join(env.CertificateRootDir, env.ServerCertsDir, "private", c.CertificateName+".key")
	}

	if pkFile, err = os.Create(pkfile); err != nil {
		return nil, err
	}
	defer pkFile.Close()

	pkBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)}
	if err = pem.Encode(pkFile, pkBlock); err != nil {
		return nil, err
	}
	return pk, err
}
