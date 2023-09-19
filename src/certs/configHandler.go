// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/configHandler.go
// Original timestamp: 2023/08/26 09:21

package certs

import (
	"certificateManager/environment"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// LoadCertificateConfFile :
// Loads the certificate config from the certificate file
func LoadCertificateConfFile(certfile string) (CertificateStruct, error) {
	var payload CertificateStruct
	var err error
	var jFile []byte
	var rcFile string
	env := environment.EnvironmentStruct{}

	// fetch environment
	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return CertificateStruct{}, err
	}

	if !strings.HasSuffix(CertConfigFile, ".json") {
		CertConfigFile += ".json"
	}
	if len(certfile) > 0 {
		rcFile = certfile
	} else {
		rcFile = filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir, CertConfigFile)
	}

	if jFile, err = os.ReadFile(rcFile); err != nil {
		return CertificateStruct{}, err
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return CertificateStruct{}, err
	} else {
		return payload, nil
	}
}

// SaveCertificateConfFile :
// Save a data structure into a certificate file in the directory defined in the JSON environment config file
func (c CertificateStruct) SaveCertificateConfFile(outfile string) error {
	var env environment.EnvironmentStruct
	var err error
	basedir := ""

	if outfile == "" {
		// fetch environment
		if env, err = environment.LoadEnvironmentFile(); err != nil {
			return err
		}
		basedir = filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir)
	}

	// why the next ?
	//if outfile == "" {
	//	outfile = CertConfigFile
	//}

	jStream, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	rcFile := filepath.Join(basedir, outfile)

	// Check if the file exists
	if _, err := os.Stat(rcFile); !os.IsNotExist(err) {
		// Remove the file if it exists
		if err := os.Remove(rcFile); err != nil {
			return err
		}
	}

	// Create the file
	file, err := os.Create(rcFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jStream)
	if err != nil {
		return err
	}

	return nil
}
