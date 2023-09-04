// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/createCerts.go
// Original timestamp: 2023/08/25 16:30

package certs

import "certificateManager/helpers"

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

func populateCertificateStructure(cs *CertificateStruct) error {
	helpers.GetStringValFromPrompt("Please enter the certificate's country", &cs.Country)
}
