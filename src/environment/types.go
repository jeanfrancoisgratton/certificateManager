// certificatemanager
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/environment/types.go
// Original timestamp: 2026/05/13 11:09:28

package environment

var EnvConfigFile string

// This structure holds the basic software config but is ignored when the software is invoked with the -s flag
// This is basically used when we store everything just like in my own internal devops/certificates/ repos

type EnvironmentStruct struct {
	CertificateRootDir    string `json:"CertificateRootDir"`
	RootCAdir             string `json:"RootCAdir"`
	ServerCertsDir        string `json:"ServerCertsDir"`
	CertificatesConfigDir string `json:"CertificatesConfigDir"`
	RemoveDuplicates      bool   `json:"RemoveDuplicates"`
}
