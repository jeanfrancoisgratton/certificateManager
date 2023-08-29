// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/indexFiles.go
// Original timestamp: 2023/08/26 12:26

package certs

import (
	"certificateManager/environment"
	"os"
	"path/filepath"
)

// writeAttributeFile() : this function is trivial in the sense that we simply ensure that index.attr and index.attr.old
// hold the same value, which is : "unique_subject = yes". This is used for a CA functionality that this software does
// not act upon
// Parameters:
// - none
// Returns : the eventual IO error, if any
func writeAttributeFile() error {
	e, err := environment.EnvironmentStruct.LoadEnvironmentFile(environment.EnvironmentStruct{})
	if err != nil {
		return err
	}

	// if the file does not exist, this means we are using a brand-new setup,
	if err := copyFile(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.attr"), filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.attr.old")); err != nil {
		return err
	} else {
		ffile, err := os.Create(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.attr"))
		if err != nil {
			return err
		}
		_, err = ffile.WriteString("unique_subject = yes")
		if err != nil {
			return err
		}
	}
	return nil
}

//_, err = ffile.WriteString("unique_subject = yes")

func writeIndexFile(ndxLine string) error {
	e, err := environment.EnvironmentStruct.LoadEnvironmentFile(environment.EnvironmentStruct{})
	if err != nil {
		return err
	}

	// if the file does not exist, this means we are using a brand-new setup,
	if err := copyFile(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt"), filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.old")); err != nil {
		return err
	} else {
		ffile, err := os.Create(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt"))
		if err != nil {
			return err
		}
		_, err = ffile.WriteString("unique_subject = yes")
		if err != nil {
			return err
		}
	}
	return nil
}
