// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/config/configTemplateCreate.go
// Original time: 2023/06/17 13:57

package config

import (
	"os"
	"path/filepath"
)

func CreateTemplatedCertificate() error {
	if err := createCertificateExplanations(); err != nil {
		return err
	}
	if err := createSampleCertificate(); err != nil {
		return err
	}
	return nil
}

func createCertificateExplanations() error {
	exptext := `{
 "CertificateRootDir" : "$HOME/.config/certificatemanager/certificates  <-- absolute path, always",
 "RootCAdir" : "rootCA  <-- relative path to CertificateRootDir",
 "ServerCertsDir" : "servers  <-- relative path to CertificateRootDir",
 "CertificatesConfigDir" : "configs"  <-- relative path to CertificateRootDir,
 "RemoveDuplicates": true  <-- should always be set to true
}`
	expFile, err := os.Create(filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", "certificatemanagerConf-README.txt"))
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

func createSampleCertificate() error {
	var sampleEnvConfig = EnvironmentStruct{
		CertificateRootDir:    "$HOME/.config/certificatemanager/certificates",
		RootCAdir:             "rootCA",
		ServerCertsDir:        "servers",
		CertificatesConfigDir: "configs",
		RemoveDuplicates:      true,
	}
	if err := sampleEnvConfig.EnvironmentFile2Json("certificatemanagerConf.json"); err != nil {
		return err
	}
	return nil
}
