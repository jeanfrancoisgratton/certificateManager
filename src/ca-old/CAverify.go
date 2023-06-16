// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/ca-old/verifyRootCAcert.go
// 4/18/23 05:37:11

package ca_old

import (
	"cm/certs"
	"cm/helpers"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var CaVerifyVerbose = false
var CaVerifyComments = false

func VerifyCACertificate(certFilePath string) error {
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
		ku := helpers.GetStringsFromKeyUsage(parsedCert.KeyUsage)
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

	if CaVerifyComments {
		cfg, err := certs.Json2Config()
		if err != nil {
			return err
		}
		if len(cfg.Comments) > 0 {
			fmt.Println("\n\nComments (part of the config, but NOT of the certificate itself)\n----------------------------------------------------------------")
			for _, cm := range cfg.Comments {
				fmt.Printf("\t• %s\n", cm)
			}
			fmt.Println()
		}
	}
	return nil
}
