// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/config/configHelpers.go
// Original time: 2023/06/16 16:37

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

//type EnvironmentStruct struct {
//	CertificateRootDir    string `json:"CertificateRootDir"`
//	RootCAdir             string `json:"RootCAdir"`
//	ServerCertsDir        string `json:"ServerCertsDir"`
//	CertificatesConfigDir string `json:"CertificatesConfigDir"`
//	RemoveDuplicates      bool   `json:"RemoveDuplicates"`
//}

var EnvConfigFile = "defaultEnvConfig.json"

func (e EnvironmentStruct) Json2EnvironmentFile() (EnvironmentStruct, error) {
	var payload EnvironmentStruct
	var err error

	if !strings.HasSuffix(EnvConfigFile, ".json") {
		EnvConfigFile += ".json"
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", EnvConfigFile)
	jFile, err := os.ReadFile(rcFile)
	if err != nil {
		return EnvironmentStruct{}, err
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return EnvironmentStruct{}, err
	} else {
		return payload, nil
	}
}

func (e EnvironmentStruct) EnvironmentFile2Json(outputfile string) error {
	if outputfile == "" {
		outputfile = EnvConfigFile
	}
	jStream, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", outputfile)
	err = os.WriteFile(rcFile, jStream, 0600)

	return err
}
