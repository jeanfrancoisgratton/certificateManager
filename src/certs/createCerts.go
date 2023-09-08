// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/createCerts.go
// Original timestamp: 2023/08/25 16:30

package certs

import (
	"certificateManager/environment"
	"certificateManager/helpers"
	"fmt"
	"path/filepath"
)

// Create() : Create a certificate file. "Plain" SSL cert, or CA cert
// Parameters:
// - none
// Returns:
// error code
// UPDATE:		I got rid of the certconfigfile string param, as I will be prompting for args if -f is unset
// UPDATE 2:	New switch, -f... we will load/create CSR, private keys, etc, as before, but if -f is set,
//
//	I will use that filename as the basic config file for the following operations

// Workflow :
// 1. Create the directory structure
// 2. Fetch the current serial number
// 3. Populate the cert structure with user-defined values
// 4. Generate private key
// 5. Generate CSR
// 6. Populate the Certificate structure with default values
// 7. Update index.txt, index.attr.txt, serial
// 8. Sign certificate

func Create(certconfigfile string) error {
	//var privateKey *rsa.PrivateKey
	var err error
	var env environment.EnvironmentStruct
	var cert CertificateStruct

	// 1. Create the directory structure to hold all of those files
	if err = createCertificateRootDirectories(); err != nil {
		return err
	}

	// 2. Get the current serial number
	if cert.SerialNumber, err = getSerialNumber(); err != nil {
		return err
	}

	// 3. Populate the certificate structure with user-provided values
	fmt.Printf("An example of a certificate can be found at %s\n", helpers.Green(filepath.Join(env.CertificatesConfigDir, "sample_cert.json")))
	if err = populateCertificateStructure(&cert); err != nil {
		return err
	}

	//// We need to create a private key
	//// Destination is either ServerCertsDir/private or RootCAdir/private
	//if privateKey, err = cert.createPrivateKey(); err != nil {
	//	return err
	//}

	//	privateKey = nil
	return nil
}

// This is a beyond-f*ckin' ugly method, only there because I want to ship this software ASAP
// and won't bother (for now) for a better solution
func populateCertificateStructure(cs *CertificateStruct) error {
	//var err error

	helpers.GetBoolValFromPrompt("Is this certificate a CA certificate ?", &cs.IsCA)
	helpers.GetStringValFromPrompt("Please enter the certificate's country (C): ", &cs.Country)
	helpers.GetStringValFromPrompt("Please enter the certificate's province (ST): ", &cs.Province)
	helpers.GetStringValFromPrompt("Please enter the certificate's locality (L): ", &cs.Locality)
	helpers.GetStringValFromPrompt("Please enter the certificate's organization (O): ", &cs.Organization)
	helpers.GetStringValFromPrompt("Please enter the certificate's organizational unit (OU): ", &cs.OrganizationalUnit)
	helpers.GetStringValFromPrompt("Please enter the certificate's common name (CN):  ", &cs.CommonName)
	helpers.GetStringSliceFromPrompt("Please enter all email addresses you want to include: ", &cs.EmailAddresses)
	helpers.GetIntValFromPrompt("Please enter the certificate lifespan (duration): ", &cs.Duration)
	// Key usage is glitchy, suboptimal....
	cs.KeyUsage = helpers.GetKeyUsage()
	helpers.GetStringSliceFromPrompt("Please enter all DNS names this cert is tied to: ", &cs.DNSNames)
	helpers.GetStringValFromPrompt("Please enter the certificate's common name (CN):  ", &cs.CommonName)

	return nil
}
