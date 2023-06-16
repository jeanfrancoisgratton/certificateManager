// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/ca-old/editRootCA.go
// 4/20/23 18:10:04

package config_old

import (
	"cm/certs"
	"cm/helpers"
	"fmt"
	"net"
	"os"
)

func CreateConfig() error {
	var err error
	certs.CertConfig, err = certs.Json2Config()
	if err != nil {
		if !os.IsNotExist(err) {
			//if err != os.ErrNotExist {
			return err
		}
	}
	err = prompt4values(&certs.CertConfig)
	if err != nil {
		return err
	}

	// we now need to reinject the config-old in a json
	err = certs.CertConfig.Config2Json(helpers.CertConfigFile)
	if err != nil {
		return err
	}
	return nil
}

func prompt4values(cfg *certs.CertificateStruct) error {
	fmt.Println(`
You will now be prompted to provide values to all of the fields that should
be part of your config-old file. If a prompt shows a value between [brackets],
this means that this value is either already present, or a suggested default
value that can be accepted by just pressing ENTER.`)

	// this is beyond ugly....
	fmt.Print("\nIs this a certificate authority (CA) ?\nPlease enter 't' or 'f') :")
	fmt.Print("Any other answer will be treated as false: ")
	fmt.Scanln(&cfg.IsCA)

	helpers.GetDuration("Please enter the certification duration, in years.\nAn invalid duration would default to 1 year", &cfg.Duration)
	helpers.GetStringValFromPrompt("Please enter the certificate name", &cfg.CertificateName)
	helpers.GetStringValFromPrompt("Please enter the certificate rootdir", &cfg.CertificateDirectory)
	helpers.GetStringValFromPrompt("Please enter the country (C)", &cfg.Country)
	helpers.GetStringValFromPrompt("Please enter the province/state (ST)", &cfg.Province)
	helpers.GetStringValFromPrompt("Please enter the locality (L)", &cfg.Locality)
	helpers.GetStringValFromPrompt("Please enter the organization (O)", &cfg.Organization)
	helpers.GetStringValFromPrompt("Please enter the organizational unit (OU)", &cfg.OrganizationalUnit)
	helpers.GetStringValFromPrompt("Please enter the common name (CN)", &cfg.CommonName)

	// A non-CA certs should not have KeyUsage
	if cfg.IsCA {
		cfg.KeyUsage = helpers.GetKeyUsageFromPrompt()
	} else {
		cfg.KeyUsage = []string{}
	}
	helpers.GetStringSliceFromPrompt("Please enter the email address(es) to be included in this certicate", &cfg.EmailAddresses)
	helpers.GetStringSliceFromPrompt("Please enter the DNS name(s) to be included in this certicate", &cfg.DNSNames)
	helpers.GetStringSliceFromPrompt("Please enter the comments to be included in this certicate\n(Note: those are for documentation purposes only, not part of the certs)", &cfg.Comments)

	// Still need net.IP...
	netip := []string{}
	if len(cfg.IPAddresses) > 0 {
		for _, x := range cfg.IPAddresses {
			netip = append(netip, x.String())
		}
	}
	helpers.GetStringSliceFromPrompt("Please enter the IP address(es) to be included in this certicate\n(Note: this is NOT recommended in a CA)", &netip)
	cfg.IPAddresses = []net.IP{}
	if len(netip) > 0 {
		for _, x := range netip {
			cfg.IPAddresses = append(cfg.IPAddresses, net.ParseIP(x))
		}
	}

	return nil
}
