// certificatemanager
// Written by J.F.Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.07.24 22:20:27
// Original filename : src/env/types.go

package environment

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
