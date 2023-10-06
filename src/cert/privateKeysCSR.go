// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/certPrivateKeys.go
// Original timestamp: 2023/08/25 18:58

// Manages private Keys and CSRs

package cert

import (
	"certificateManager/environment"
	"certificateManager/helpers"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
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

	if pk, err = rsa.GenerateKey(rand.Reader, CertPKsize); err != nil {
		return nil, err
	}

	// rootCA keys are not stored at the same place as other SSL keys
	if c.IsCA {
		if err = os.MkdirAll(filepath.Join(env.CertificateRootDir, env.RootCAdir), os.ModePerm); err != nil {
			return nil, err
		}
		pkfile = filepath.Join(env.CertificateRootDir, env.RootCAdir, c.CertificateName+".key")
	} else {
		if err = os.MkdirAll(filepath.Join(env.CertificateRootDir, env.ServerCertsDir, "private"), os.ModePerm); err != nil {
			return nil, err
		}
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

func (c CertificateStruct) getPrivateKey(env environment.EnvironmentStruct) (*rsa.PrivateKey, error) {
	var err error
	var pKeyFile []byte
	var pkey *rsa.PrivateKey
	keyDir := env.CertificateRootDir

	// root CAs store their key somewhere else
	if c.IsCA {
		keyDir = filepath.Join(keyDir, env.RootCAdir)
	} else {
		keyDir = filepath.Join(keyDir, env.ServerCertsDir, "private")
	}
	// Load keyfile
	if pKeyFile, err = os.ReadFile(filepath.Join(keyDir, c.CertificateName+".key")); err != nil {
		return nil, helpers.CustomError{Message: "Error reading the private key: " + err.Error()}
	}

	// Decode keyfile
	if key, _ := pem.Decode(pKeyFile); key == nil {
		return nil, helpers.CustomError{Message: "Unable to PEM-decode the private key"}
	} else {
		// Parse keyfile
		if pkey, err = x509.ParsePKCS1PrivateKey(key.Bytes); err != nil {
			return nil, err
		}
	}
	return pkey, nil
}

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
