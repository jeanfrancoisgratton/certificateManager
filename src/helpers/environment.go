// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/helpers/creds.go
// 4/29/23 17:00:45

package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func Json2Config() (CertConfigStruct, error) {
	var payload CertConfigStruct
	var err error

	if !strings.HasSuffix(CertConfigFile, ".json") {
		CertConfigFile += ".json"
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", CertConfigFile)
	jFile, err := os.ReadFile(rcFile)
	if err != nil {
		return CertConfigStruct{}, err
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return CertConfigStruct{}, err
	} else {
		return payload, nil
	}
}

func (c CertConfigStruct) Config2Json(outputfile string) error {
	if outputfile == "" {
		outputfile = CertConfigFile
	}
	jStream, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", outputfile)
	err = os.WriteFile(rcFile, jStream, 0600)

	return err
}
