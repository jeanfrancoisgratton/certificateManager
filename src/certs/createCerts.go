// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/createCerts.go
// Original timestamp: 2023/08/25 16:30

package certs

import "certificateManager/certs"

func Create() error {
	if err := certs.CreateCertificateRootDirectories(); err != nil {
		return err
	}

	return nil
}
