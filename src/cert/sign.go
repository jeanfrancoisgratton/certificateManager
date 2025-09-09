// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/sign.go
// Original timestamp: 2023/09/11 10:28

package cert

import (
	"certificateManager/environment"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"software.sslmate.com/src/go-pkcs12"
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
func (c CertificateStruct) signCert(env environment.EnvironmentStruct) *cerr.CustomError {
	var csrBytes, caCertPEM, caKeyPEM []byte
	var csrRequest *x509.CertificateRequest
	var caCert *x509.Certificate
	var caKey *rsa.PrivateKey
	var err error

	// Ensure there is a single file in the CA directory and fetch its name
	caCertFiles, err := filepath.Glob(filepath.Join(env.RootCAdir, "*.crt"))
	if err != nil {
		return &cerr.CustomError{Title: "Error listing CA certificate files: ", Message: err.Error(), Fatality: cerr.Fatal}
	}
	if len(caCertFiles) != 1 {
		return &cerr.CustomError{Message: "Expected one CA certificate file, found " + hf.Red(string(len(caCertFiles)))}
	}
	baseFN := strings.TrimSuffix(filepath.Base(caCertFiles[0]), filepath.Ext(filepath.Base(caCertFiles[0])))

	// 1. Load the CA cert and key files
	if caCertPEM, err = os.ReadFile(filepath.Join(env.RootCAdir, baseFN+".crt")); err != nil {
		return &cerr.CustomError{Title: "Error reading CA certificate: ", Message: err.Error(), Fatality: cerr.Fatal}
	}
	if caKeyPEM, err = os.ReadFile(filepath.Join(env.RootCAdir, baseFN+".key")); err != nil {
		return &cerr.CustomError{Title: "Error reading CA private key: ", Message: err.Error(), Fatality: cerr.Fatal}
	}

	// 2. Parse the CA cert and key files
	caCertBlock, _ := pem.Decode(caCertPEM)
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caCertBlock == nil || caKeyBlock == nil {
		return &cerr.CustomError{Message: "Error PEM-decoding the CA certificate or its private key", Fatality: cerr.Fatal}
	}
	if caCert, err = x509.ParseCertificate(caCertBlock.Bytes); err != nil {
		return &cerr.CustomError{Title: "Error parsing CA certificate", Message: err.Error(), Fatality: cerr.Fatal}
	}
	if caKey, err = x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes); err != nil {
		return &cerr.CustomError{Title: "Error parsing CA private key", Message: err.Error(), Fatality: cerr.Fatal}
	}

	// 3. Load, decode and parse the CSR file
	if csrBytes, err = os.ReadFile(filepath.Join(env.ServerCertsDir, "csr", c.CertificateName+".csr")); err != nil {
		return &cerr.CustomError{Title: "Error reading CSR certificate: ", Message: err.Error(), Fatality: cerr.Fatal}
	}
	if csrBlock, _ := pem.Decode(csrBytes); csrBlock == nil {
		return &cerr.CustomError{Message: "Error PEM-decoding the CSR file", Fatality: cerr.Fatal}
	} else {
		if csrRequest, err = x509.ParseCertificateRequest(csrBlock.Bytes); err != nil {
			return &cerr.CustomError{Title: "Error parsing CSR certificate", Message: err.Error(), Fatality: cerr.Fatal}
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
		return &cerr.CustomError{Title: "Error creating certificate", Message: err.Error(), Fatality: cerr.Fatal}
	}

	// 6. Encode, save to disk
	certFile, err := os.Create(filepath.Join(env.ServerCertsDir, "certs", c.CertificateName+".crt"))
	if err != nil {
		return &cerr.CustomError{Title: "Error creating certificate file: ", Message: err.Error(), Fatality: cerr.Fatal}
	}
	defer certFile.Close()

	if err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return &cerr.CustomError{Title: "Error encoding certificate to PEM: ", Message: err.Error(), Fatality: cerr.Fatal}
	}

	// We also need to save the new certificate in the rootCA "newcerts" directory
	if err = os.Mkdir(filepath.Join(env.RootCAdir, "newcerts"), os.ModePerm); err != nil && !os.IsExist(err) {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}
	newcertFile, err := os.Create(filepath.Join(env.RootCAdir, "newcerts", fmt.Sprintf("%04X.pem", c.SerialNumber)))
	if err != nil {
		return &cerr.CustomError{Title: "Unable to create the certificate within root CA's PKI: ", Message: err.Error(), Fatality: cerr.Fatal}
	}
	defer newcertFile.Close()
	if pem.Encode(newcertFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return &cerr.CustomError{Title: "Error encoding certificate to PEM: ", Message: err.Error(), Fatality: cerr.Fatal}
	}

	if CertJava {
		return c.createJavaCert(env, caCert, caKey)
	}

	fmt.Printf("Certificate %s with a duration of %v years successfully created in %s\n",
		hf.White(c.CertificateName), hf.White(fmt.Sprintf("%v", c.Duration)), hf.White(filepath.Join(env.ServerCertsDir, "certs")))
	return nil
}

// createCA and signCert are very similar: one is for non-CA cert, the other (below) for CA cert
// I *could* fold both into a single function, with tons of "if c.IsCA{}" clauses, but it's not worth
// the readability headache that it'd bring
func (c CertificateStruct) createCA(env environment.EnvironmentStruct, privateKey *rsa.PrivateKey) *cerr.CustomError {
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
		return &cerr.CustomError{Title: "Unable to create CA certificate: ", Message: err.Error(), Fatality: cerr.Fatal}
	}

	cafile, err := os.Create(filepath.Join(env.RootCAdir, c.CertificateName+".crt"))
	if err != nil {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}
	defer cafile.Close()

	if err = pem.Encode(cafile, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}

	fmt.Printf("Root CA certificate %s with a duration of %v years successfully created in %s\n",
		hf.White(c.CertificateName), hf.White(fmt.Sprintf("%v", c.Duration)), hf.White(env.RootCAdir))
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
func (c CertificateStruct) createJavaCert(e environment.EnvironmentStruct, caCert *x509.Certificate, caKey *rsa.PrivateKey) *cerr.CustomError {
	var certPEM []byte
	var certBlock *pem.Block
	var err error
	var ce *cerr.CustomError
	var serverCert *x509.Certificate
	var serverKey *rsa.PrivateKey
	certPasswd := ""

	// Fetch the server's private key
	if serverKey, ce = c.getPrivateKey(e); ce != nil {
		return ce
	}

	// Load, decode and parse the current server cert
	if certPEM, err = os.ReadFile(filepath.Join(e.ServerCertsDir, "certs", c.CertificateName+".crt")); err != nil {
		return &cerr.CustomError{Title: "Error reading CA certificate", Message: err.Error(), Fatality: cerr.Fatal}
	}
	certBlock, _ = pem.Decode(certPEM)
	if serverCert, err = x509.ParseCertificate(certBlock.Bytes); err != nil {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}

	// PKCS#12 requires the file to be password-protected
	certPasswd = hf.GetPassword(fmt.Sprintf("Please provide a password for %s: ", c.CertificateName), false)

	// Remove outdated .p12 and .jks files, if present
	basename := filepath.Join(e.ServerCertsDir, "java", c.CertificateName)
	for _, fn := range []string{basename + ".p12", basename + ".jks"} {
		errfn := os.Remove(fn)
		if errfn != nil {
			if os.IsNotExist(errfn) {
				continue
			} else {
				return &cerr.CustomError{Title: fmt.Sprintf("Unable to remove %s : ", fn),
					Message: err.Error(), Fatality: cerr.Fatal}
			}
		}
	}

	// Convert cert to PKCS#12
	pkcs12Data, err := pkcs12.Encode(rand.Reader, serverKey, serverCert, []*x509.Certificate{caCert}, certPasswd)
	if err != nil {
		return &cerr.CustomError{Title: "Error encoding the certificate in PKCS#12: ",
			Message: err.Error(), Fatality: cerr.Fatal}
	}

	if err = os.WriteFile(filepath.Join(e.ServerCertsDir, "java", c.CertificateName+".p12"), pkcs12Data, 0644); err != nil {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}

	// No other way for now <sigh>
	cmd := exec.Command("keytool", "-importkeystore", "-srcstorepass", certPasswd,
		"-deststorepass", certPasswd,
		"-destkeystore", filepath.Join(e.ServerCertsDir, "java", c.CertificateName+".jks"),
		"-srckeystore", filepath.Join(e.ServerCertsDir, "java", c.CertificateName+".p12"),
		"-srcstoretype", "PKCS12")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return &cerr.CustomError{Title: "Keytool command failed: ", Message: err.Error(), Fatality: cerr.Fatal}
	}

	return nil
}
