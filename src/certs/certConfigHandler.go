// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/certConfigHandler.go
// Original timestamp: 2023/08/26 09:21

package certs

// getCertificateConfig() : fetches an SSL certificate configuration and returns it in a data structure
// Params:
// - name (string) : the name of the certificate
// Returns:
// - the certificate config in a CertificateStruct
// - the error code, if any
func getCertificateConfig(name string) (CertificateStruct, error) {
	var c, sslCertConfig CertificateStruct
	var err error

	if sslCertConfig, err = c.LoadCertificateConfFile(name); err != nil {
		return CertificateStruct{}, err
	}
	return sslCertConfig, nil
}
