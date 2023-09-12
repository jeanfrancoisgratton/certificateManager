// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/indexFiles.go
// Original timestamp: 2023/08/26 12:26

package certs

import (
	"bufio"
	"certificateManager/environment"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// writeAttributeFile() : this function is trivial in the sense that we simply ensure that index.attr and index.attr.old
// hold the same value, which is : "unique_subject = yes". This is used for a CA functionality that this software does
// not act upon
// Parameters:
// - none
// Returns : the eventual IO error, if any
func writeAttributeFile() error {
	e, err := environment.LoadEnvironmentFile()
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

func writeIndexFile(c CertificateStruct) error {
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}

	// if the file does not exist, this means we are using a brand-new setup,
	if err := copyFile(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt"), filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.old")); err != nil {
		return err
	} else {
		if err := inPlaceReplace(c, filepath.Join(e.CertificateRootDir, e.RootCAdir)); err != nil {
			return err
		}
	}
	return nil
}

func inPlaceReplace(c CertificateStruct, sourcedir string) error {
	string2replace := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/%s/CN=%s", c.Country, c.Province,
		c.Locality, c.Organization, c.CommonName)
	sf, err := os.Open(filepath.Join(sourcedir, "index.txt.old"))
	if err != nil {
		return err
	}
	defer sf.Close()
	of, err := os.Create(filepath.Join(sourcedir, "index.txt"))
	if err != nil {
		return err
	}
	defer of.Close()

	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, string2replace) {
			continue // we will not write this line in the target file
		}

		_, err := fmt.Fprintln(of, line)
		if err != nil {
			return err
		}
	}
	newline := fmt.Sprintf("V\t%sZ\t\t%s\t%s\tunknown\t%s", time.Now().UTC().Format("060102150405"),
		c.SerialNumber, string2replace)
	_, err = fmt.Fprintln(of, newline)
	if err != nil {
		return err
	}
	return nil
}
