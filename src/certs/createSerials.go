// certificateManager
// Écrit par J.F.Gratton (jean-francois@famillegratton.net)
// createSerials.go, jfgratton : 2023-08-26

package certs

import (
	"certificateManager/environment"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
	_, err = ffile.WriteString(fmt.Sprintf("%04X\n", serialNo))
	if err != nil {
		return err
	}
	return nil
}
