// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/create.go
// Original timestamp: 2023/08/25 16:30

package cert

import (
	"certificateManager/environment"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"path/filepath"

	cerr "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
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
// 8. Save/update the certificate config file in the config directory

func Create(certconfigfile string) *cerr.CustomError {
	var privateKey *rsa.PrivateKey
	var err *cerr.CustomError
	var env environment.EnvironmentStruct
	var certconfig CertificateStruct
	isDupe := false
	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}

	// 1. Create the directory structure to hold all of those files
	if err = createCertificateRootDirectories(); err != nil {
		return err
	}

	// 2a. Populate the certificate structure with user-provided values or a file
	if certconfigfile == "" {
		fmt.Printf("An example of a certificate can be found at %s\n", hftx.Green(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "certificatemanager", "sampleCert.json")))
		if err = populateCertificateStructure(&certconfig); err != nil {
			return err
		}
	} else {
		if certconfig, err = LoadCertificateConfFile(certconfigfile); err != nil {
			return err
		}
	}

	// 2b. Check if the proposed certificate already exists
	// There is no reason in a well-behaved PKI to allow duplicates. I offer the possibility just because there might
	// be use-cases that I am not aware of
	if env.RemoveDuplicates {
		if isDupe, err = certconfig.check4DuplicateCert(filepath.Join(env.RootCAdir, "index.txt")); err != nil {
			return err
		}
		if isDupe {
			return &cerr.CustomError{Title: "Certificate creation aborted", Message: "That certificate is already present in index.txt", Fatality: cerr.Warning}
		}
	}

	// Corner case : as EmailAddress is part of a certificate signature (ie: this is part on how
	// We differentiate the registered certconfig), we need to have a value in this field.
	if len(certconfig.EmailAddresses) == 0 {
		certconfig.EmailAddresses = []string{"none"}
	}

	// 3. Get the current serial number
	if certconfig.SerialNumber, err = getSerialNumber(); err != nil {
		return err
	}
	certconfig.SerialNumber++

	// 4. Generate a private key
	// Destination is either ServerCertsDir/private or RootCAdir
	if privateKey, err = certconfig.createPrivateKey(); err != nil {
		return err
	}

	// 5. Generate the CSR (if not a CA certconfig)
	if !certconfig.IsCA {
		if err = certconfig.generateCSR(env, privateKey); err != nil {
			return err
		}
	}

	// 6. Generate the certificate, also sign it if non-CA certconfig
	if certconfig.IsCA {
		if err := certconfig.createCA(env, privateKey); err != nil {
			return err
		}
	} else {
		if err := certconfig.signCert(env); err != nil {
			return err
		}
	}

	// 7. Update serial, index.txt.attr and index.txt
	// serial
	if err = setSerialNumber(certconfig.SerialNumber); err != nil {
		return err
	}
	// index.txt.attr
	if err = writeAttributeFile(); err != nil {
		return err
	}
	// index.txt
	if err = writeIndexFile(certconfig); err != nil {
		return err
	}

	// 8. Save JSON config file
	if err = certconfig.SaveCertificateConfFile(""); err != nil {
		return err
	}

	if !certconfig.IsCA {
		fmt.Printf("Certificate %s has been created.\n", hftx.Green(certconfig.CertificateName))
	}
	return nil
}

// This is a beyond ugly method, only there because I want to ship this software ASAP
// and won't bother (for now) for a better solution
func populateCertificateStructure(cs *CertificateStruct) *cerr.CustomError {
	var err *cerr.CustomError
	//var ips []string
	fmt.Println("Entries with multiple values (ip addresses, emails, key usage are separated with ENTER, with another ENTER pressed at the end.")
	cs.CertificateName = hf.GetStringValFromPrompt(fmt.Sprintf("Please enter the certificate's %s: ", hftx.Green("name")))
	cs.CommonName = hf.GetStringValFromPrompt(fmt.Sprintf("Please enter the %s (CN) -> defaults to cert name if omitted: ", hftx.Green("common name")))
	if cs.CommonName == "" {
		cs.CommonName = cs.CertificateName
	}
	cs.IsCA = hf.GetBoolValFromPrompt(fmt.Sprintf("[any values not starting with T,t or 1 will be treated as FALSE] Is this certificate a %s ? ", hftx.Green("CA certificate")))
	cs.Country = hf.GetStringValFromPrompt(fmt.Sprintf("Please enter the certificate's %s (C): ", hftx.Green("country")))
	cs.Province = hf.GetStringValFromPrompt(fmt.Sprintf("Please enter the certificate's %s (ST): ", hftx.Green("province/state")))
	cs.Locality = hf.GetStringValFromPrompt(fmt.Sprintf("Please enter the certificate's %s (L): ", hftx.Green("locality")))
	cs.Organization = hf.GetStringValFromPrompt(fmt.Sprintf("Please enter the certificate's %s (O): ", hftx.Green("organization")))
	cs.OrganizationalUnit = hf.GetStringValFromPrompt(fmt.Sprintf("Please enter the certificate's %s (OU): ", hftx.Green("organizational unit")))
	cs.EmailAddresses = hf.GetStringSliceFromPrompt(fmt.Sprintf("Please enter the certificate's %s: ", hftx.Green("email address")))
	// Here we have to deal with a new security "feature" by Apple that does not allow for CA to last longer than +- 825 days ("official" documentation varies, here)
	// Therefore we won't allow CAs to expire past 800 days
	if cs.IsCA {
		fmt.Printf("%s: Apple introduced a 'feature' where CA expiring after 825 days are no longer deemed valid. We will limit their duration to %s years\n",
			hftx.Yellow("WARNING:"), hftx.Red("2"))
	}
	if cs.Duration = hf.GetIntValFromPrompt(fmt.Sprintf("\nPlease enter the certificate's lifespan (%s) in YEARS, ENTER is 1: ", hftx.Green("duration"))); cs.Duration <= 1 {
		cs.Duration = 1
	}
	if cs.IsCA && cs.Duration > 2 {
		cs.Duration = 2
	}

	// Key usage is glitchy, suboptimal....
	fmt.Printf("Please enter the %s intended for this certificate:\n", hftx.Green("key usage"))
	cs.KeyUsage = getKeyUsage()
	cs.DNSNames = hf.GetStringSliceFromPrompt(fmt.Sprintf("Please enter all %s this cert is tied to: ", hftx.Green("DNS names")))
	ips := hf.GetStringSliceFromPrompt(fmt.Sprintf("\nPlease enter the certificate's %s: ", hftx.Green("IP address(es)")))
	if len(ips) > 0 {
		for _, val := range ips {
			cs.IPAddresses = append(cs.IPAddresses, net.ParseIP(val))
		}
	} else {
		cs.IPAddresses = []net.IP{}
	}
	if cs.SerialNumber, err = getSerialNumber(); err != nil {
		return err
	}
	cs.SerialNumber++

	cs.Comments = hf.GetStringSliceFromPrompt(fmt.Sprintf("\nPlease enter optional %s: ", hftx.Green("comments")))
	return nil
}
