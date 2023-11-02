// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/sign.go
// Original timestamp: 2023/09/11 10:28

package cert

import (
	"certificateManager/environment"
//	"software.sslmate.com/src/go-pkcs12"
	"certificateManager/helpers"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// signCert: Sign the certificate against the root CA currently held in our custom PKI
// Steps:
// 0. Ensure that there is a single CA (.crt) file in the CA directory
// 1. Load the CA cert and private key
// 2. Decode and parse them
// 3. Load & parse the CSR file
// 4. Populate a x509 cert template with the CertificateStruct values
// 5. Sign (create) the certificate
// 6. Save to disk
func (c CertificateStruct) signCert(env environment.EnvironmentStruct) error {
	var csrBytes, caCertPEM, caKeyPEM []byte
	var csrRequest *x509.CertificateRequest
	var caCert *x509.Certificate
	var caKey *rsa.PrivateKey
	var err error

	// Ensure there is a single file in the CA directory and fetch its name
	caCertFiles, err := filepath.Glob(filepath.Join(env.RootCAdir, "*.crt"))
	if err != nil {
		return helpers.CustomError{Message: "Error listing CA certificate files: " + err.Error()}
	}
	if len(caCertFiles) != 1 {
		return helpers.CustomError{Message: "Expected one CA certificate file, found " + helpers.Red(string(len(caCertFiles)))}
	}
	baseFN := strings.TrimSuffix(filepath.Base(caCertFiles[0]), filepath.Ext(filepath.Base(caCertFiles[0])))

	// 1. Load the CA cert and key files
	if caCertPEM, err = os.ReadFile(filepath.Join(env.RootCAdir, baseFN+".crt")); err != nil {
		return helpers.CustomError{Message: "Error reading CA certificate: " + err.Error()}
	}
	if caKeyPEM, err = os.ReadFile(filepath.Join(env.RootCAdir, baseFN+".key")); err != nil {
		return helpers.CustomError{Message: "Error reading CA private key: " + err.Error()}
	}

	// 2. Parse the CA cert and key files
	caCertBlock, _ := pem.Decode(caCertPEM)
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caCertBlock == nil || caKeyBlock == nil {
		return helpers.CustomError{Message: "Error PEM-decoding the CA certificate or its private key"}
	}
	if caCert, err = x509.ParseCertificate(caCertBlock.Bytes); err != nil {
		return err
	}
	if caKey, err = x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes); err != nil {
		return err
	}

	// 3. Load, decode and parse the CSR file
	if csrBytes, err = os.ReadFile(filepath.Join(env.ServerCertsDir, "csr", c.CertificateName+".csr")); err != nil {
		return err
	}
	if csrBlock, _ := pem.Decode(csrBytes); csrBlock == nil {
		return helpers.CustomError{"Error PEM-decoding the CSR file"}
	} else {
		if csrRequest, err = x509.ParseCertificateRequest(csrBlock.Bytes); err != nil {
			return err
		}
	}

	// 4. Populate x509 template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(int64(c.SerialNumber)),
		Subject:               pkix.Name{CommonName: c.CommonName, Locality: []string{c.Locality}, Country: []string{c.Country}, Organization: []string{c.Organization}, OrganizationalUnit: []string{c.OrganizationalUnit}, Province: []string{c.Province}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(c.Duration, 0, 0),
		KeyUsage:              getKeyUsageFromStrings(c.KeyUsage),
		IsCA:                  c.IsCA,
		BasicConstraintsValid: true,
		DNSNames:              c.DNSNames,
		IPAddresses:           c.IPAddresses,
		EmailAddresses:        c.EmailAddresses,
	}
	//// why the next ???
	//if c.IsCA {
	//	template.KeyUsage = reindexKeyUsage(c)
	//}

	// 5. Create (sign) the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, csrRequest.PublicKey, caKey)
	if err != nil {
		return err
	}

	// 6. Encode, save to disk
	certFile, err := os.Create(filepath.Join(env.ServerCertsDir, "certs", c.CertificateName+".crt"))
	if err != nil {
		return err
	}
	defer certFile.Close()

	if err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	// We also need to save the new certificate in the rootCA "newcerts" directory
	if err = os.Mkdir(filepath.Join(env.RootCAdir, "newcerts"), os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	newcertFile, err := os.Create(filepath.Join(env.RootCAdir, "newcerts", fmt.Sprintf("%04X.pem", c.SerialNumber)))
	if err != nil {
		return helpers.CustomError{Message: "Unable to create the certificate within root CA's PKI: " + err.Error()}
	}
	defer newcertFile.Close()
	if pem.Encode(newcertFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	if CertJava {
		return c.createJavaCert(env, caCert, caKey)
	}

	fmt.Printf("Certificate %s with a duration of %v years successfully created in %s\n",
		helpers.White(c.CertificateName), helpers.White(fmt.Sprintf("%v", c.Duration)), helpers.White(filepath.Join(env.ServerCertsDir, "certs")))
	return nil
}

// createCA and signCert are very similar: one is for non-CA cert, the other (below) for CA cert
// I *could* fold both into a single function, with tons of "if c.IsCA{}" clauses, but it's not worth
// the readability headache that it'd bring
func (c CertificateStruct) createCA(env environment.EnvironmentStruct, privateKey *rsa.PrivateKey) error {
	var caBytes []byte
	var err error

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: c.CommonName, Locality: []string{c.Locality}, Country: []string{c.Country}, Organization: []string{c.Organization}, OrganizationalUnit: []string{c.OrganizationalUnit}, Province: []string{c.Province}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(c.Duration, 0, 0),
		KeyUsage:              getKeyUsageFromStrings(c.KeyUsage),
		IsCA:                  c.IsCA,
		BasicConstraintsValid: true,
		DNSNames:              c.DNSNames,
		IPAddresses:           c.IPAddresses,
		EmailAddresses:        c.EmailAddresses,
	}
	if caBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey); err != nil {
		return err
	}

	cafile, err := os.Create(filepath.Join(env.RootCAdir, c.CertificateName+".crt"))
	if err != nil {
		return err
	}
	defer cafile.Close()

	if err = pem.Encode(cafile, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
		return err
	}

	fmt.Printf("Root CA certificate %s with a duration of %v years successfully created in %s\n",
		helpers.White(c.CertificateName), helpers.White(fmt.Sprintf("%v", c.Duration)), helpers.White(env.RootCAdir))
	return nil
}

// createJavaCert:
// Much software still use the Java Keystore (JKS) format, which has been deemed obsolete for some time.
// The process is thus, so far:
// 1. Load the CA cert file, key, server cert file, key
// NEW step: remove old (outdated) .p12 and .jks files, if present
// 2. Convert the server .crt to PKCS#12 (.p12) format
// 3. Convert the .p12 file to .JKS
// I'll keep the .p12 file in storage, just in case that whatever software needing a JKS comes to its senses
// And asks for a .p12 instead
// All files will be stored in the java/ directory

// SIGNATURE: (environment, parsed cacert, parsed cacert key) returns error
func (c CertificateStruct) createJavaCert(e environment.EnvironmentStruct, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	var certPEM []byte
	var certBlock *pem.Block
	var err error
	var serverCert *x509.Certificate
	var serverKey *rsa.PrivateKey
	certPasswd := ""

	// Fetch the server's private key
	if serverKey, err = c.getPrivateKey(e); err != nil {
		return err
	}

	// Load, decode and parse the current server cert
	if certPEM, err = os.ReadFile(filepath.Join(e.ServerCertsDir, "certs", c.CertificateName+".crt")); err != nil {
		return helpers.CustomError{Message: "Error reading CA certificate: " + err.Error()}
	}
	certBlock, _ = pem.Decode(certPEM)
	if serverCert, err = x509.ParseCertificate(certBlock.Bytes); err != nil {
		return err
	}

	// PKCS#12 requires the file to be password-protected
	certPasswd = helpers.GetPassword("Please provide a password for this Java certificate: ")

	// Remove outdated .p12 and .jks files, if present
	//basename := filepath.Join(e.ServerCertsDir, "java", c.CertificateName)
	//for _, fn := range []string{basename + ".p12", basename + ".jks"} {
	//	errfn := os.Remove(fn)
	//	if errfn != nil {
	//		if os.IsNotExist(errfn) {
	//			continue
	//		} else {
	//			return helpers.CustomError{Message: fmt.Sprintf("Unable to remove %s : ", fn) + err.Error()}
	//		}
	//	}
	//}
	// Convert cert to PKCS#12
	pkcs12Data, err := pkcs12.Encode(rand.Reader, serverKey, serverCert, []*x509.Certificate{caCert}, certPasswd)
	if err != nil {
		return helpers.CustomError{Message: "Error encoding the certificate in PKCS#12: " + err.Error()}
	}

	if err = os.WriteFile(filepath.Join(e.ServerCertsDir, "java", c.CertificateName+".p12"), pkcs12Data, 0644); err != nil {
		return err
	}

	// FOLLOWING COMMENTED CODE IS THERE IN CASE I FIND A SOLUTION TO REPLACE cmd:= exec.Command(), a bit below

	// Now, it's a bit stupid, but I need to re-decode the P12 file to then encode it in JKS
	//if p12Data, err = os.ReadFile(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java", c.CertificateName+".p12")); err != nil {
	//	return err
	//}
	//if p12, _, err := pkcs12.Decode(p12Data, certPasswd); err != nil {
	//	return err
	//}
	//
	//// Create Keystore, encode
	//if jksFile, err = os.Create(filepath.Join(e.CertificateRootDir, e.ServerCertsDir, "java", c.CertificateName+".jks")); err != nil {
	//	return err
	//}
	//defer jksFile.Close()

	// No other way for now <sigh>
	cmd := exec.Command("keytool", "-importkeystore", "-srcstorepass", certPasswd,
		"-deststorepass", certPasswd,
		"-destkeystore", filepath.Join(e.ServerCertsDir, "java", c.CertificateName+".jks"),
		"-srckeystore", filepath.Join(e.ServerCertsDir, "java", c.CertificateName+".p12"),
		"-srcstoretype", "PKCS12")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return helpers.CustomError{Message: "Keytool command failed: " + err.Error()}
	}

	return nil
}
