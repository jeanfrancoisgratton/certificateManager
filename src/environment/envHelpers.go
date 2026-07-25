// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/environment/envHelpers.go
// Original timestamp: 2023/08/19 10:02

package environment

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	cerr "github.com/jeanfrancoisgratton/customError/v3"
)

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
