// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/configHandler.go
// Original timestamp: 2023/08/26 09:21

package cert

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"certificateManager/environment"
	cerr "github.com/jeanfrancoisgratton/customError/v3"
)

// LoadCertificateConfFile :
// Loads the certificate config from the certificate file
func LoadCertificateConfFile(certfile string) (CertificateStruct, *cerr.CustomError) {
	var payload CertificateStruct
	var err *cerr.CustomError
	var jFile []byte
	var rcFile string
	env := environment.EnvironmentStruct{}

	// fetch environment
	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return CertificateStruct{}, err
	}

	if certfile != "" {
		if !strings.HasSuffix(certfile, ".json") {
			certfile += ".json"
		}
		//rcFile = filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir, filepath.Base(certfile))
		rcFile = filepath.Join(env.CertificatesConfigDir, filepath.Base(certfile))
	}

	var readErr error
	if jFile, readErr = os.ReadFile(rcFile); readErr != nil {
		return CertificateStruct{}, &cerr.CustomError{Title: "Unable to read certificate config file",
			Message: readErr.Error(), Fatality: cerr.Fatal}
	}
	if readErr = json.Unmarshal(jFile, &payload); readErr != nil {
		return CertificateStruct{}, &cerr.CustomError{Title: "Unable to unmarshal JSON",
			Message: readErr.Error(), Fatality: cerr.Fatal}
	}
	return payload, nil
}

// SaveCertificateConfFile :
// Save a data structure into a certificate file in the directory defined in the JSON environment config file
func (c CertificateStruct) SaveCertificateConfFile(outfile string) *cerr.CustomError {
	var env environment.EnvironmentStruct
	var err error
	var ce *cerr.CustomError

	if outfile == "" {
		// fetch environment
		if env, ce = environment.LoadEnvironmentFile(); ce != nil {
			return ce
		}
		outfile = filepath.Join(env.CertificatesConfigDir, c.CertificateName+".json")
	}

	if !strings.HasSuffix(outfile, ".json") {
		outfile += ".json"
	}

	jStream, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return &cerr.CustomError{Title: "Unable to marshal JSON",
			Message: err.Error(), Fatality: cerr.Fatal}
	}

	// Check if the file exists
	if _, err := os.Stat(outfile); !os.IsNotExist(err) {
		// Remove the file if it exists
		if err := os.Remove(outfile); err != nil {
			return &cerr.CustomError{Title: "Unable to remove existing file",
				Message: err.Error(), Fatality: cerr.Fatal}
		}
	}

	// Create the file
	file, err := os.Create(outfile)
	if err != nil {
		return &cerr.CustomError{Title: "Unable to create the file",
			Message: err.Error(), Fatality: cerr.Fatal}
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jStream)
	if err != nil {
		return &cerr.CustomError{Title: "Unable to write the file",
			Message: err.Error(), Fatality: cerr.Fatal}
	}
	return nil
}
