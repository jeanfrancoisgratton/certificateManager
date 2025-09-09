// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/environment/envHelpers.go
// Original timestamp: 2023/08/19 10:02

package environment

import (
	"encoding/json"
	cerr "github.com/jeanfrancoisgratton/customError"
	"os"
	"path/filepath"
	"strings"
)

var EnvConfigFile string

// This structure holds the basic software config but is ignored when the software is invoked with the -s flag
// This is basically used when we store everything just like in my own internal gitea devops/certificates/ repos
type EnvironmentStruct struct {
	CertificateRootDir    string `json:"CertificateRootDir"`
	RootCAdir             string `json:"RootCAdir"`
	ServerCertsDir        string `json:"ServerCertsDir"`
	CertificatesConfigDir string `json:"CertificatesConfigDir"`
	RemoveDuplicates      bool   `json:"RemoveDuplicates"`
}

// Load the JSON environment file in the user's .config/certificatemanager directory, and store it into a data type (struct)
func LoadEnvironmentFile() (EnvironmentStruct, *cerr.CustomError) {
	var payload EnvironmentStruct
	var err error

	if !strings.HasSuffix(EnvConfigFile, ".json") {
		EnvConfigFile += ".json"
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "certificatemanager", EnvConfigFile)
	jFile, err := os.ReadFile(rcFile)
	if err != nil {
		return EnvironmentStruct{}, &cerr.CustomError{Title: err.Error(), Fatality: cerr.Fatal}
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return EnvironmentStruct{}, &cerr.CustomError{Title: err.Error(), Fatality: cerr.Fatal}
	} else {
		return payload, nil
	}
}

// Save the above structure into a JSON file in the user's .config/certificatemanager directory
func (e EnvironmentStruct) SaveEnvironmentFile(outputfile string) *cerr.CustomError {
	if outputfile == "" {
		outputfile = EnvConfigFile
	}
	jStream, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return &cerr.CustomError{Title: err.Error(), Fatality: cerr.Fatal}
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "certificatemanager", outputfile)
	if err = os.WriteFile(rcFile, jStream, 0600); err != nil {
		return &cerr.CustomError{Title: err.Error(), Fatality: cerr.Fatal}
	}

	return nil
}
