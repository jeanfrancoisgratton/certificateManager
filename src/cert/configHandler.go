// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/configHandler.go
// Original timestamp: 2023/08/26 09:21

package cert

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

	// why is that....
	//if !strings.HasSuffix(CertConfigFile, ".json") {
	//	CertConfigFile += ".json"
	//}
	if certfile != "" {
		if !strings.HasSuffix(certfile, ".json") {
			certfile += ".json"
		}
		rcFile = filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir, filepath.Base(certfile))
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
	//basedir := filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir)
	//basedir := ""

	if outfile == "" {
		// fetch environment
		if env, err = environment.LoadEnvironmentFile(); err != nil {
			return err
		}
		outfile = filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir, c.CertificateName+".json")
	}

	if !strings.HasSuffix(outfile, ".json") {
		outfile += ".json"
	}
	// why the next ?
	//if outfile == "" {
	//	outfile = CertConfigFile
	//}

	jStream, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	//rcFile := filepath.Join(basedir, outfile)

	// Check if the file exists
	if _, err := os.Stat(outfile); !os.IsNotExist(err) {
		// Remove the file if it exists
		if err := os.Remove(outfile); err != nil {
			return err
		}
	}

	// Create the file
	file, err := os.Create(outfile)
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
