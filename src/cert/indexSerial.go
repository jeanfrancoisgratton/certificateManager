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

	cerr "github.com/jeanfrancoisgratton/customError/v3"
)

// writeAttributeFile() : this function is trivial in the sense that we simply ensure that  index.attr.old
// holds the same value every run, which is : "unique_subject = yes". This is used for a CA functionality that this software does
// not act upon
// Parameters:
// - none
// Returns : the eventual IO error, if any
func writeAttributeFile() *cerr.CustomError {
	e, ce := environment.LoadEnvironmentFile()
	if ce != nil {
		return ce
	}

	ffile, err := os.Create(filepath.Join(e.RootCAdir, "index.txt.attr"))
	if err != nil {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}
	_, err = ffile.WriteString("unique_subject = yes")
	if err != nil {
		return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
	}
	return nil
}

func writeIndexFile(c CertificateStruct) *cerr.CustomError {
	var filedesc *os.File
	var err error
	e, ce := environment.LoadEnvironmentFile()
	if ce != nil {
		return ce
	}

	if _, err = os.Stat(filepath.Join(e.RootCAdir, "index.txt")); os.IsNotExist(err) {
		if filedesc, err = os.Create(filepath.Join(e.RootCAdir, "index.txt")); err != nil {
			return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
		}
		defer filedesc.Close()
		newline := fmt.Sprintf("V\t%sZ\t%s\tunknown\t/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s/emailAddress=%s", time.Now().UTC().Format("060102150405"),
			fmt.Sprintf("%04x", c.SerialNumber), c.Country, c.Province, c.Locality, c.Organization, c.OrganizationalUnit, c.CommonName, c.EmailAddresses[0])
		if _, err = filedesc.WriteString(newline); err != nil {
			return &cerr.CustomError{Message: err.Error(), Fatality: cerr.Fatal}
		}
	} else {
		if err := replaceStringInIndex(c, e.RootCAdir); err != nil {
			return err
		}
	}
	return nil
}

func replaceStringInIndex(c CertificateStruct, sourcedir string) *cerr.CustomError {
	string2replace := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s/emailaddress=%s", c.Country, c.Province,
		c.Locality, c.Organization, c.OrganizationalUnit, c.CommonName, c.EmailAddresses[0])

	sf, ferr := os.Open(filepath.Join(sourcedir, "index.txt"))
	if ferr != nil {
		return &cerr.CustomError{Message: ferr.Error(), Fatality: cerr.Fatal}
	}
	defer sf.Close()
	of, ferr := os.Create(filepath.Join(sourcedir, "index.txt.tmp"))
	if ferr != nil {
		return &cerr.CustomError{Message: ferr.Error(), Fatality: cerr.Fatal}
	}
	defer of.Close()

	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, string2replace) {
			continue // we will not write this line in the target file
		}
		_, ferr := fmt.Fprintln(of, line)
		if ferr != nil {
			return &cerr.CustomError{Message: ferr.Error(), Fatality: cerr.Fatal}
		}
	}
	newline := fmt.Sprintf("V\t%sZ\t%s\tunknown\t%s", time.Now().UTC().Format("060102150405"),
		fmt.Sprintf("%04X", c.SerialNumber), string2replace)
	_, ferr = fmt.Fprintln(of, newline)
	if ferr != nil {
		return &cerr.CustomError{Message: ferr.Error(), Fatality: cerr.Fatal}
	}
	if ferr := os.Rename(filepath.Join(sourcedir, "index.txt.tmp"), filepath.Join(sourcedir, "index.txt")); ferr != nil {
		return &cerr.CustomError{Message: ferr.Error(), Fatality: cerr.Fatal}
	}
	return nil
}

// getSerialNumber() : returns the current serial number on file (typically in CertificateRootDir/RootCAdir/serial)
// Parameters:
// - none
// Returns:
// - uint64 representing the decimal value of the serial number, or zero if error
// - the error code
func getSerialNumber() (uint64, *cerr.CustomError) {
	var serr error

	// We need the environment file
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return 0, err
	}
	serialPath := filepath.Join(e.RootCAdir, "serial")

	// if the serial file does not exist, this means we are using a brand new setup,
	// thus the serial # is 1
	_, serr = os.Stat(serialPath)
	if os.IsNotExist(serr) {
		return 0, nil
	}
	// Read serial from file
	content, serr := os.ReadFile(serialPath)
	if serr != nil {
		return 0, &cerr.CustomError{Message: serr.Error(), Fatality: cerr.Fatal}
	}

	// Convert content to a string and remove any leading/trailing whitespace
	hexString := strings.TrimSpace(string(content))
	// Corner case: file (serialPath) exists, but is of zero byte length
	if hexString == "" {
		hexString = "0"
	}

	// Convert hexadecimal string to a uint64
	decimalValue, perr := strconv.ParseUint(hexString, 16, 64)
	if perr != nil {
		decimalValue = 0
	}
	return decimalValue, nil
}

// setSerialNumber() : Sets the serial value on file (typically in CertificateRootDir/RootCAdir/serial)
// We will also keep a backup of the serial file
func setSerialNumber(serialNo uint64) *cerr.CustomError {
	// We need the environment file
	e, err := environment.LoadEnvironmentFile()
	if err != nil {
		return err
	}

	ffile, ce := os.Create(filepath.Join(e.RootCAdir, "serial"))
	if ce != nil {
		return &cerr.CustomError{Title: "Error creating the serial file", Message: err.Error(), Fatality: cerr.Fatal}
	}
	defer ffile.Close()

	_, ce = ffile.WriteString(fmt.Sprintf("%04X\n", serialNo))
	if ce != nil {
		return &cerr.CustomError{Title: "Error writing the serial file", Message: err.Error(), Fatality: cerr.Fatal}
	}
	return nil
}
