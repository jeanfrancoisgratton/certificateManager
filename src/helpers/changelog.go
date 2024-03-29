// certificateManager
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/helpers/changelog.go
// Original timestamp: 2023/09/30 16:07

package helpers

import "fmt"

func Changelog() {
	//fmt.Printf("\x1b[2J")
	fmt.Printf("\x1bc")

	CenterPrint("CHANGELOG")
	fmt.Println()
	fmt.Println()

	fmt.Print(`
VERSION		DATE			COMMENT
-------		----			-------
1.24.00		2024.03.01		GO version bump
1.23.00		2023.12.07		GO version bump to 1.21.5, packages upgrades
1.22.00		2023.11.02		cm cert rm was broken (duplicated path)
1.21.00		2023.11.02		defaultEnv.json does not need to be specified anymore if it's the environment we use
1.20.06		2023.10.31		more explicit error message in cert verify
1.20.05		2023.10.16		moved to a saner version numbering scheme
1.205		2023.10.15		fixed wrong path for server's private keys
1.200		2023.10.13		go version bump, folded all environment directories into a single var for readability issues with filepath.Join()
1.100-0		2023.10.04		cm cert verify now works w/ the verbose flag; duplicate certificate creation is now prevented 
1.010-1		2023.10.03		Fixed typo in directory name
1.010		2023.10.03		Fixed issue when serial number was not incremented within certificate
1.001		2023.10.03		Minor changes: verbosity, doc update
1.000		2023.09.30		Completed prod-ready version
0.500		2023.06.03		server cert management
0.400		2023.04.22		config-old management
0.300		2023.04.20		ca-old edit, ca-old del
0.200		2023.04.20		ca-old create and ca-old verify
0.100		2023.04.16		near-config-old-aware
\n`)
}
