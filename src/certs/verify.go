// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/verify.go
// Original timestamp: 2023/08/23 15:22

package certs

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

var CaVerifyVerbose = false
var CaVerifyComments = false

func Verify(certFilePaths []string) error {
	//var e environment.EnvironmentStruct
	//var cert CertificateStruct
	//var err error

	//if e, err = environment.LoadEnvironmentFile(); err != nil {
	//	return err
	//}

	for _, path := range certFilePaths {
		//if err := verifyCert(filepath.Join(e.CertificateRootDir, path)); err != nil {
		if err := verifyCert(path); err != nil {
			return err
		}
	}
	return nil
}

func verifyCert(certFilePath string) error {

	if !strings.HasSuffix(certFilePath, ".crt") {
		certFilePath += ".crt"
	}
	// Read the certificate file
	certPEMBlock, err := os.ReadFile(certFilePath)
	if err != nil {
		// we let caVerifyCmd() deal with the error
		return err
	}

	// Decode the PEM block into a certificate
	cert, _ := pem.Decode(certPEMBlock)
	if cert == nil {
		return fmt.Errorf("failed to decode certificate PEM block")
	}

	// Parse the certificate
	parsedCert, err := x509.ParseCertificate(cert.Bytes)
	if err != nil {
		return err
	}

	// Print certificate information
	fmt.Printf("Certificate: %s\n---\n", certFilePath)
	fmt.Printf("   Is this a Certificate Authority (root CA) ? ")
	if parsedCert.IsCA {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
	if CaVerifyVerbose {
		fmt.Printf("   Data:\n%s\n", string(certPEMBlock))
	}
	fmt.Printf("   Serial Number: %v\n\n", parsedCert.SerialNumber)
	fmt.Printf("   Subject: %v\n", parsedCert.Subject)
	fmt.Printf("   Issuer: %v\n", parsedCert.Issuer)
	fmt.Printf("\n   Not Before: %v\n", parsedCert.NotBefore)
	fmt.Printf("   Not After : %v\n", parsedCert.NotAfter)
	if len(parsedCert.IPAddresses) > 0 {
		fmt.Println("\n   IP Address(es):")
		for _, ipa := range parsedCert.IPAddresses {
			fmt.Printf("\t• %s\n", ipa)
		}
	}
	if len(parsedCert.EmailAddresses) > 0 {
		fmt.Println("\n   Email Address(es):")
		for _, email := range parsedCert.EmailAddresses {
			fmt.Printf("\t• %s\n", email)
		}
	}
	if len(parsedCert.URIs) > 0 {
		fmt.Println("\n   URIs:")
		for _, uri := range parsedCert.URIs {
			fmt.Printf("\t• %v\n", uri)
		}
	}
	if CaVerifyVerbose {
		fmt.Printf("   Signature Algorithm: %v\n", parsedCert.SignatureAlgorithm)
		fmt.Printf("   Signature: %v\n", parsedCert.Signature)
	}

	// x509v3 EXTENSIONS:
	fmt.Println("\n   x509v3 extensions\n   -----------------")

	if parsedCert.KeyUsage != 0 {
		fmt.Printf("\n   x509v3 Key usage:\n")
		ku := getStringsFromKeyUsage(parsedCert.KeyUsage)
		for _, k := range ku {
			fmt.Printf("\t• %s\n", k)
		}
	}

	// Print X509v3 Subject Alternative Name
	if len(parsedCert.DNSNames) > 0 {
		fmt.Printf("\n   x509v3 Subject Alternative Names (SAN):\n")
		for _, dns := range parsedCert.DNSNames {
			fmt.Printf("\t• %s\n", dns)
		}
	}

	// TODO: fix this
	if CaVerifyComments {
		var c CertificateStruct
		var err error
		if c, err = LoadCertificateConfFile(""); err != nil {
			return err
		}
		if len(c.Comments) > 0 {
			fmt.Println("\n\nComments (part of the config, but NOT of the certificate itself)\n----------------------------------------------------------------")
			for _, cm := range c.Comments {
				fmt.Printf("\t• %s\n", cm)
			}
			fmt.Println()
		}
	}
	return nil
}
