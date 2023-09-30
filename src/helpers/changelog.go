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
1.000		2023.09.30		Completed prod-ready version
0.500		2023.06.03		server certs management
0.400		2023.04.22		config-old management
0.300		2023.04.20		ca-old edit, ca-old del
0.200		2023.04.20		ca-old create and ca-old verify
0.100		2023.04.16		near-config-old-aware
\n`)
}
