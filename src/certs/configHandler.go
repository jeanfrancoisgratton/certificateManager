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
func (c CertificateStruct) SaveCertificateConfFile(outputfile string) error {
	var env environment.EnvironmentStruct
	var err error

	if outputfile == "" {

	}
	// fetch environment
	if env, err = environment.LoadEnvironmentFile(); err != nil {
		return err
	}
	if outputfile == "" {
		outputfile = CertConfigFile
	}
	if _, err := os.Stat(filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(env.CertificateRootDir, env.CertificatesConfigDir), os.ModePerm)
	}

	jStream, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	rcFile := filepath.Join(env.CertificateRootDir, outputfile)
	//remove file if exists
	if _, err := os.Stat(rcFile); os.IsExist(err) {
		os.Remove(rcFile)
	}

	return os.WriteFile(rcFile, jStream, 0600)
}
