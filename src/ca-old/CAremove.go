// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/ca-old/removeRootCA.go
// 4/22/23 08:55:20

package ca_old

import (
	"cm/certs"
	"os"
	"path/filepath"
)

// This is a stub, really, before we get to the actual removal in branch 0.600
func RemoveCACertificate() error {
	cfg, err := certs.Json2CertConfig()

	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(cfg.CertificateDirectory, cfg.CertificateName, ".key"))
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(cfg.CertificateDirectory, cfg.CertificateName, ".crt"))

	return err
}
