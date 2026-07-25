// certificatemanager
// Written by J.F.Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.07.24 22:32:20
// Original filename : src/cert/types.go

package cert

import "net"

var CertPKsize int

var CertConfigFile = "defaultCertConfig.json"
var CertJava = false
var CertRemoveFiles = false
var CaVerifyVerbose = false
var CaVerifyComments = false

// This is the full data structure for an SSL certificate and CA

type CertificateStruct struct {
	Country            string   `json:"Country"`
	Province           string   `json:"Province"`
	Locality           string   `json:"Locality"`
	Organization       string   `json:"Organization"`
	OrganizationalUnit string   `json:"OrganizationalUnit,omitempty"`
	CommonName         string   `json:"CommonName"`
	IsCA               bool     `json:"IsCA"`
	EmailAddresses     []string `json:"EmailAddresses,omitempty"`
	Duration           int      `json:"Duration"`
	KeyUsage           []string `json:"KeyUsage"`
	DNSNames           []string `json:"DNSNames,omitempty"`
	IPAddresses        []net.IP `json:"IPAddresses,omitempty"`
	CertificateName    string   `json:"CertificateName"`
	SerialNumber       uint64   `json:"SerialNumber"`
	Comments           []string `json:"Comments,omitempty"`
}
