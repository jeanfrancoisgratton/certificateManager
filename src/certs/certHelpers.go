// certificateManager
// Écrit par J.F. Gratton <jean-francois@famillegratton.net>
// Orininal name: src/certs/certHelpers.go
// Original time: 2023/06/16 16:37

package certs

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var CertConfig = CertificateStruct{Duration: 1, KeyUsage: []string{"certs sign", "crl sign", "digital signature"}}
var CertConfigFile = "defaultCertConfig.json"

type CertificateStruct struct {
	Country              string   `json:"Country"`
	Province             string   `json:"Province"`
	Locality             string   `json:"Locality"`
	Organization         string   `json:"Organization"`
	OrganizationalUnit   string   `json:"OrganizationalUnit,omitempty"`
	CommonName           string   `json:"CommonName"`
	IsCA                 bool     `json:"IsCA,omitempty"`
	EmailAddresses       []string `json:"EmailAddresses,omitempty"`
	Duration             int      `json:"Duration"`
	KeyUsage             []string `json:"KeyUsage"`
	DNSNames             []string `json:"DNSNames,omitempty"`
	IPAddresses          []net.IP `json:"IPAddresses,omitempty"`
	CertificateDirectory string   `json:"CertificateDirectory"`
	CertificateName      string   `json:"CertificateName"`
	Comments             []string `json:"Comments,omitempty"`
}

func Json2CertConfig() (CertificateStruct, error) {
	var payload CertificateStruct
	var err error

	if !strings.HasSuffix(CertConfigFile, ".json") {
		CertConfigFile += ".json"
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "certificatemanager", CertConfigFile)
	jFile, err := os.ReadFile(rcFile)
	if err != nil {
		return CertificateStruct{}, err
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return CertificateStruct{}, err
	} else {
		return payload, nil
	}
}

func (c CertificateStruct) CertConfig2Json(outputfile string) error {
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
