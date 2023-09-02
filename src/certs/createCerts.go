// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/createCerts.go
// Original timestamp: 2023/08/25 16:30

package certs

import "crypto/rsa"

func Create(certconfigfile string) error {
	var privateKey *rsa.PrivateKey
	var err error
	var cert CertificateStruct

	if cert, err = LoadCertificateConfFile(certconfigfile); err != nil {
		return err
	}

	if err = createCertificateRootDirectories(); err != nil {
		return err
	}

	if privateKey, err = cert.createPrivateKey(); err != nil {
		return err
	}

	privateKey = nil
	return nil
}
