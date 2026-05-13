// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/cert/helpers.go
// Original time: 2023/06/16 16:37

package cert

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"certificateManager/environment"
	cerr "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

var CertConfigFile = "defaultCertConfig.json"
var CertJava = false
var CertRemoveFiles = false

// This is the full data structure for an SSL certificate and CA
type CertificateStruct struct {
	Country            string   `json:"Country"`
	Province           string   `json:"Province"`
	Locality           string   `json:"Locality"`
	Organization       string   `json:"Organization"`
	OrganizationalUnit string   `json:"OrganizationalUnit,omitempty"`
	CommonName         string   `json:"CommonName"`
	IsCA               bool     `json:"IsCA"`
	EmailAddresses     []string `json:"EmailAddresses,omitempty"`
	Duration           int      `json:"Duration"`
	KeyUsage           []string `json:"KeyUsage"`
	DNSNames           []string `json:"DNSNames,omitempty"`
	IPAddresses        []net.IP `json:"IPAddresses,omitempty"`
	CertificateName    string   `json:"CertificateName"`
	SerialNumber       uint64   `json:"SerialNumber"`
	Comments           []string `json:"Comments,omitempty"`
}

// Create the sample certificate config file

// createCertificateRootDirectories() : creates the directory structure needed to store all cert, keys, CA, CSR, etc
func createCertificateRootDirectories() *cerr.CustomError {
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}
	dirRange := []string{e.RootCAdir, e.CertificatesConfigDir, filepath.Join(e.ServerCertsDir, "private"), filepath.Join(e.ServerCertsDir, "csr"),
		filepath.Join(e.ServerCertsDir, "certs"), filepath.Join(e.ServerCertsDir, "java")}

	for _, directory := range dirRange {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
		}
	}
	return nil
}

// check4DuplicateCert :
// We browse the index.txt database to see if the signature of thecertificate we are creating already exists
func (c CertificateStruct) check4DuplicateCert(ndxFilePath string) (bool, *cerr.CustomError) {
	lookoutstring := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s", c.Country, c.Province, c.Locality,
		c.Organization, c.OrganizationalUnit, c.CommonName)

	isDupe := false

	// open index.txt db
	indexfileHandle, err := os.Open(ndxFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}
	defer indexfileHandle.Close()

	// Scan the file to find a duplicate cert signature
	scanner := bufio.NewScanner(indexfileHandle)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, lookoutstring) && strings.HasPrefix(line, "V") {
			isDupe = true
			break
		}
	}
	return isDupe, nil
}

func getKeyUsage() []string {
	var keys []string
	inputScanner := bufio.NewScanner(os.Stdin)
	ku := []string{"decipher only", "encipher only", "crl sign", "cert sign", "key agreement",
		"data encipherment", "key encipherment", "content commitment", "digital signature", "CADEFAULTS", "CERTDEFAULTS"}
	inputs := []string{}

	fmt.Println("The valid key usage values are:")
	for i, j := range ku {
		if i%5 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("'%s' ", hftx.White(j))
	}
	fmt.Println("\nCADEFAULTS is a catch-all for default values of a root CA (lowercase accepted)")
	fmt.Println("CERTDEFAULTS is a catch-all for default values of a standard certificate (lowercase accepted)")
	for {
		input := ""
		fmt.Print("Please enter values from the above list, press ENTER to end : ")
		inputScanner.Scan()
		input = inputScanner.Text()
		if input == "" {
			break
		}
		if strings.ToLower(input) == "cadefaults" || strings.ToLower(input) == "certdefaults" {
			if strings.ToLower(input) == "cadefaults" {
				inputs = append(inputs, "digital signature", "cert sign", "crl sign")
			} else {
				inputs = append(inputs, "digital signature", "key encipherment",
					"data encipherment", "key agreement")
			}
		} else {
			if valueInList(input, ku) {
				inputs = append(inputs, input)
			}
		}
	}
	// if the array is empty, we return a default value
	if len(inputs) == 0 {
		keys = []string{"digital signature"}
		return keys
	}
	// now we need to ensure that we do not have any duplicates
	s := make([]string, 0, len(inputs))
	m := make(map[string]bool)

	for _, value := range inputs {
		if _, ok := m[value]; !ok {
			m[value] = true
			s = append(s, value)
		}
	}
	keys = s
	return keys
}

func valueInList(in string, list []string) bool {
	for _, x := range list {
		if x == in {
			return true
		}
	}
	return false
}
