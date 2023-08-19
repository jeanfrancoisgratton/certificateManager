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
)

var EnvConfigFile = "defaultEnvConfig.json"

type EnvironmentStruct struct {
	CertificateRootDir    string `json:"CertificateRootDir"`
	RootCAdir             string `json:"RootCAdir"`
	ServerCertsDir        string `json:"ServerCertsDir"`
	CertificatesConfigDir string `json:"CertificatesConfigDir"`
	RemoveDuplicates      bool   `json:"RemoveDuplicates"`
}

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

func CreateSampleEnv() error {
	var err error
	e := EnvironmentStruct{filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "certificates"), "rootCA", "servers", "conf", true}
	//e := EnvironmentStruct{filepath.Join(os.Getenv("HOME"),".config","certificatemanager"),"certificates", "rootCA", "servers", "conf", true}

	if err = e.EnvironmentFile2Json("environmentSample.json"); err != nil {
		return err
	}

	exptext := `{
 "CertificateRootDir" : "$HOME/.config/certificatemanager/certificates  <-- absolute path, always",
 "RootCAdir" : "rootCA  <-- relative path to CertificateRootDir",
 "ServerCertsDir" : "servers  <-- relative path to CertificateRootDir",
 "CertificatesConfigDir" : "conf"  <-- relative path to CertificateRootDir,
 "RemoveDuplicates": true  <-- should always be set to true, there is no use-case yet to set it to false
}`
	expFile, err := os.Create(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "environmentSample-README.txt"))
	if err != nil {
		return err
	}
	defer expFile.Close()

	_, err = expFile.WriteString(exptext)
	if err != nil {
		return err
	}

	return nil
}
