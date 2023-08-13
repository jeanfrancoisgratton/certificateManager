// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/helpers/keyUsage.go
// 4/29/23 17:01:24

package helpers

import (
	"certificateManager/certs"
	"crypto/x509"
	"strings"
)

// GetKeyUsageFromStrings() : converts a slice of strings into
// A x509.KeyUsage value. We use slices of strings because x509.KeyUsage
// Can hold multiple operations at once
func GetKeyUsageFromStrings(usageStrings []string) x509.KeyUsage {
	keyUsage := x509.KeyUsage(0)
	for _, usage := range usageStrings {
		switch strings.ToLower(usage) {
		case "digital signature":
			keyUsage |= x509.KeyUsageDigitalSignature
		case "content commitment":
			keyUsage |= x509.KeyUsageContentCommitment
		case "key encipherment":
			keyUsage |= x509.KeyUsageKeyEncipherment
		case "data encipherment":
			keyUsage |= x509.KeyUsageDataEncipherment
		case "key agreement":
			keyUsage |= x509.KeyUsageKeyAgreement
		case "certs sign", "certificate sign":
			keyUsage |= x509.KeyUsageCertSign
		case "crl sign", "crl":
			keyUsage |= x509.KeyUsageCRLSign
		case "encipheronly", "encipher":
			keyUsage |= x509.KeyUsageEncipherOnly
		case "decipheronly", "decipher":
			keyUsage |= x509.KeyUsageDecipherOnly
		}
	}
	return keyUsage
}

// GetStringsFromKeyUsage(): takes the x509.KeyUsage numerical value
// And converts it in a slice of human-readable strings,
// As KeyUsage can hold multiple operations at once.
func GetStringsFromKeyUsage(keyUsage x509.KeyUsage) []string {
	var usages []string

	if keyUsage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "digital signature")
	}
	if keyUsage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "content commitment")
	}
	if keyUsage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "key encipherment")
	}
	if keyUsage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "data encipherment")
	}
	if keyUsage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "key agreement")
	}
	if keyUsage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "certs sign")
	}
	if keyUsage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "crl sign")
	}
	if keyUsage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "encipher only")
	}
	if keyUsage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "decipher only")
	}
	return usages
}

// ReindexKeyUsage() : Ensures that the CertificateStruct.KeyUsage contains only unique values
func ReindexKeyUsage(cfg certs.CertificateStruct) x509.KeyUsage {
	org := cfg.KeyUsage
	// We append the CA-related usages
	org = append(org, "certs sign", "crl sign", "digital signature")

	// We map the new slices
	//[]string to map : https://kylewbanks.com/blog/creating-unique-slices-in-go
	s := make([]string, 0, len(org))
	m := make(map[string]bool)

	for _, value := range org {
		if _, ok := m[value]; !ok {
			m[value] = true
			s = append(s, value)
		}
	}
	return GetKeyUsageFromStrings(s)
}
