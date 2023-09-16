// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/createCerts.go
// Original timestamp: 2023/08/25 16:30

package certs

import (
	"certificateManager/environment"
	"certificateManager/helpers"
	"crypto/rsa"
	"fmt"
	"net"
	"path/filepath"
)

// Create() : Create a certificate file. "Plain" SSL cert, or CA cert
// Parameters:
// - none
// Returns:
// error code
// UPDATE:		I got rid of the certconfigfile string param, as I will be prompting for args if -f is unset
// UPDATE 2:	New switch, -f... we will load/create CSR, private keys, etc., as before, but if -f is set,
//
//	I will use that filename as the basic config file for the following operations

// Workflow :
// 1. Create the directory structure
// 2. Populate the cert structure with user-defined values
// 3. Fetch the current serial number, increment it
// 4. Generate private key
// 5. Generate CSR
// 6. Sign certificate
// 7. Update index.txt, index.attr.txt, serial

func Create(certconfigfile string) error {
	var privateKey *rsa.PrivateKey
	var err error
	var env environment.EnvironmentStruct
	var cert CertificateStruct
	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}
	// 1. Create the directory structure to hold all of those files
	if err = createCertificateRootDirectories(); err != nil {
		return err
	}

	// 2. Populate the certificate structure with user-provided values or a file
	if certconfigfile == "" {
		fmt.Printf("An example of a certificate can be found at %s\n", helpers.Green(filepath.Join(env.CertificatesConfigDir, "sample_cert.json")))
		if err = populateCertificateStructure(&cert); err != nil {
			return err
		}
	} else {
		if cert, err = LoadCertificateConfFile(certconfigfile); err != nil {
			return err
		}
	}

	// 3. Get the current serial number
	if cert.SerialNumber, err = getSerialNumber(); err != nil {
		return err
	} else {
		cert.SerialNumber++
	}

	// 4. Generate a private key
	// Destination is either ServerCertsDir/private or RootCAdir
	if privateKey, err = cert.createPrivateKey(); err != nil {
		return err
	}

	// 5. Generate the CSR (if not a CA cert)
	// 6. Sign the certificate
	if !cert.IsCA {
		if err = cert.generateCSR(env, privateKey); err != nil {
			return err
		}
		if err := cert.signCert(env); err != nil {
			return err
		}
	}

	// 7. Update serial, index.txt.attr and index.txt
	// serial
	if err = setSerialNumber(cert.SerialNumber); err != nil {
		return err
	}
	// index.txt.attr
	if err = writeAttributeFile(); err != nil {
		return err
	}
	// index.txt
	if err = writeIndexFile(cert); err != nil {
		return err
	}

	return nil
}

// This is a beyond ugly method, only there because I want to ship this software ASAP
// and won't bother (for now) for a better solution
func populateCertificateStructure(cs *CertificateStruct) error {
	var err error
	//var ips []string
	fmt.Println("Entries with multiple values (ip addresses, emails, key usage are separated with ENTER, with another ENTER pressed at the end.\n")
	cs.CertificateName = helpers.GetStringValFromPrompt("Please enter the certificate's name: ")
	cs.IsCA = helpers.GetBoolValFromPrompt("Is this certificate a CA certificate ? ")
	cs.Country = helpers.GetStringValFromPrompt("Please enter the certificate's country (C): ")
	cs.Province = helpers.GetStringValFromPrompt("Please enter the certificate's province (ST): ")
	cs.Locality = helpers.GetStringValFromPrompt("Please enter the certificate's locality (L): ")
	cs.Organization = helpers.GetStringValFromPrompt("Please enter the certificate's organization (O): ")
	cs.OrganizationalUnit = helpers.GetStringValFromPrompt("Please enter the certificate's organizational unit (OU): ")
	cs.CommonName = helpers.GetStringValFromPrompt("Please enter the certificate's common name (CN):  ")
	cs.EmailAddresses = helpers.GetStringSliceFromPrompt("Please enter all email addresses you want to include: ")
	cs.Duration = helpers.GetIntValFromPrompt("Please enter the certificate lifespan (duration): ")
	// Key usage is glitchy, suboptimal....
	cs.KeyUsage = helpers.GetKeyUsage()
	cs.DNSNames = helpers.GetStringSliceFromPrompt("Please enter all DNS names this cert is tied to: ")
	ips := helpers.GetStringSliceFromPrompt("Please enter the certificate's IP address(es): ")
	if len(ips) > 0 {
		for _, val := range ips {
			cs.IPAddresses = append(cs.IPAddresses, net.ParseIP(val))
		}
	} else {
		cs.IPAddresses = []net.IP{}
	}
	if cs.SerialNumber, err = getSerialNumber(); err != nil {
		return err
	} else {
		cs.SerialNumber++
	}
	cs.Comments = helpers.GetStringSliceFromPrompt("Please enter optional comments: ")
	return nil
}
