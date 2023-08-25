// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/certsCreate.go
// Original timestamp: 2023/08/25 16:30

package certs

func Create() error {
	if err := createCertificateRootDirectories(); err != nil {
		return err
	}

	return nil
}
