// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/certs/createIndexFiles.go
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

// writeAttributeFile() : this function is trivial in the sense that we simply ensure that  index.attr.old
// holds the same value every run, which is : "unique_subject = yes". This is used for a CA functionality that this software does
// not act upon
// Parameters:
// - none
// Returns : the eventual IO error, if any
func writeAttributeFile() error {
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}

	ffile, err := os.Create(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt.attr"))
	if err != nil {
		return err
	}
	_, err = ffile.WriteString("unique_subject = yes")
	if err != nil {
		return err
	}
	return nil
}

func writeIndexFile(c CertificateStruct) error {
	var filedesc *os.File
	//string2replace := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/CN=%s/emailAddress=%s", c.Country, c.Province,
	//	c.Locality, c.Organization, c.CommonName, c.EmailAddresses[0])
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}

	if _, err = os.Stat(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt")); os.IsNotExist(err) {
		if filedesc, err = os.Create(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt")); err != nil {
			return err
		}
		defer filedesc.Close()
		newline := fmt.Sprintf("V\t%sZ\t\t%s\t%s\tunknown\t%s/C=%s/ST=%s/L=%s/O=%s/CN=%s/emailAddress=%s", time.Now().UTC().Format("060102150405"),
			fmt.Sprintf("%04x", c.SerialNumber), c.Country, c.Province, c.Locality, c.Organization, c.CommonName, c.EmailAddresses[0])
		if _, err = filedesc.WriteString(newline); err != nil {
			return err
		}
	} else {
		if err := replaceStringInIndex(c, filepath.Join(e.CertificateRootDir, e.RootCAdir)); err != nil {
			return err
		}
	}
	return nil
}

func replaceStringInIndex(c CertificateStruct, sourcedir string) error {
	string2replace := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/CN=%s", c.Country, c.Province,
		c.Locality, c.Organization, c.CommonName)
	sf, err := os.Open(filepath.Join(sourcedir, "index.txt"))
	if err != nil {
		return err
	}
	defer sf.Close()
	of, err := os.Create(filepath.Join(sourcedir, "index.txt.tmp"))
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
	if err := os.Rename(filepath.Join(sourcedir, "index.txt.tmp"), filepath.Join(sourcedir, "index.txt")); err != nil {
		return err
	}
	return nil
}
