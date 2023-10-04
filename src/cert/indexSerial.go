// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cert/indexSerial.go
// Original timestamp: 2023/08/26 12:26

// Manages the index.txt* and serial files

package cert

import (
	"bufio"
	"certificateManager/environment"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}

	if _, err = os.Stat(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt")); os.IsNotExist(err) {
		if filedesc, err = os.Create(filepath.Join(e.CertificateRootDir, e.RootCAdir, "index.txt")); err != nil {
			return err
		}
		defer filedesc.Close()
		newline := fmt.Sprintf("V\t%sZ\t%s\tunknown\t/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s/emailAddress=%s", time.Now().UTC().Format("060102150405"),
			fmt.Sprintf("%04x", c.SerialNumber), c.Country, c.Province, c.Locality, c.Organization, c.OrganizationalUnit, c.CommonName, c.EmailAddresses[0])
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
	string2replace := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s/emailaddress=%s", c.Country, c.Province,
		c.Locality, c.Organization, c.OrganizationalUnit, c.CommonName, c.EmailAddresses[0])

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
	newline := fmt.Sprintf("V\t%sZ\t%s\tunknown\t%s", time.Now().UTC().Format("060102150405"),
		fmt.Sprintf("%04X", c.SerialNumber), string2replace)
	_, err = fmt.Fprintln(of, newline)
	if err != nil {
		return err
	}
	if err := os.Rename(filepath.Join(sourcedir, "index.txt.tmp"), filepath.Join(sourcedir, "index.txt")); err != nil {
		return err
	}
	return nil
}

// getSerialNumber() : returns the current serial number on file (typically in CertificateRootDir/RootCAdir/serial)
// Parameters:
// - none
// Returns:
// - uint64 representing the decimal value of the serial number, or zero if error
// - the error code
func getSerialNumber() (uint64, error) {
	// We need the environment file
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return 0, err
	}
	serialPath := filepath.Join(e.CertificateRootDir, e.RootCAdir, "serial")

	// if the serial file does not exist, this means we are using a brand new setup,
	// thus the serial # is 1
	_, err = os.Stat(serialPath)
	if os.IsNotExist(err) {
		return 0, nil
	}
	// Read serial from file
	content, err := os.ReadFile(serialPath)
	if err != nil {
		return 0, err
	}

	// Convert content to a string and remove any leading/trailing whitespace
	hexString := strings.TrimSpace(string(content))
	// Corner case: file (serialPath) exists, but is of zero byte length
	if hexString == "" {
		hexString = "0"
	}

	// Convert hexadecimal string to a uint64
	decimalValue, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		decimalValue = 0
	}
	return decimalValue, err
}

// setSerialNumber() : Sets the serial value on file (typically in CertificateRootDir/RootCAdir/serial)
// We will also keep a backup of the serial file
func setSerialNumber(serialNo uint64) error {
	// We need the environment file
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}

	ffile, err := os.Create(filepath.Join(e.CertificateRootDir, e.RootCAdir, "serial"))
	if err != nil {
		return err
	}
	defer ffile.Close()

	_, err = ffile.WriteString(fmt.Sprintf("%04X\n", serialNo))
	if err != nil {
		return err
	}
	return nil
}
