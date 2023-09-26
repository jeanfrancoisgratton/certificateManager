// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/helpers/createHelpers.go
// 4/29/23 17:36:16

package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetIntValFromPrompt(prompt string) int {
	var err error
	value := 0
	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s: ", prompt)
	inputScanner.Scan()
	nval := inputScanner.Text()

	if nval != "" {
		value, err = strconv.Atoi(nval)
		if err != nil {
			value = 1
		}
	}
	return value
}

func GetBoolValFromPrompt(prompt string) bool {
	fmt.Printf("%s(any values not starting with T,t or 1 will be treated as FALSE): ", prompt)
	bval := ""
	var value = false

	fmt.Scanln(&bval)
	if strings.HasPrefix(strings.ToLower(bval), "t") || bval == "1" {
		value = true
	}
	return value
}

func GetStringValFromPrompt(prompt string) string {
	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s: ", prompt)
	inputScanner.Scan()
	nval := inputScanner.Text()
	value := ""

	if nval != "" {
		value = nval
	}
	return value
}

func GetStringSliceFromPrompt(prompt string) []string {
	slice := []string{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s\n", prompt)
	for {
		fmt.Println("Just press enter to end the loop")
		scanner.Scan()
		input := scanner.Text()
		if input == "" {
			break
		} else {
			slice = append(slice, input)
		}
	}
	return slice
}

func GetKeyUsage() []string {
	var keys []string
	inputScanner := bufio.NewScanner(os.Stdin)
	ku := []string{"decipher only", "encipher only", "crl sign", "certs sign", "key agreement",
		"data encipherment", "key encipherment", "content commitment", "digital signature", "CADEFAULTS"}
	inputs := []string{}

	fmt.Println("The valid key usage values are:")
	for i, j := range ku {
		if i%5 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("'%s' ", White(j))
	}
	fmt.Println("CADEFAULTS is a catch-all for default values of a root CA")
	for {
		input := ""
		fmt.Print("Please enter a value from the above list, just press ENTER to end : ")
		inputScanner.Scan()
		input = inputScanner.Text()
		if input == "" {
			break
		}
		if input == "CADEFAULTS" {
			inputs = append(inputs, "digital signature", "certs sign", "crl sign")
		} else {
			if valueInList(input, ku) {
				inputs = append(inputs, input)
			}
		}
	}
	// if the array is empty, we return a default value
	if len(inputs) == 0 {
		keys = []string{"digital signature"}
		return keys
	}
	// now we need to ensure that we do not have any duplicates
	s := make([]string, 0, len(inputs))
	m := make(map[string]bool)

	for _, value := range inputs {
		if _, ok := m[value]; !ok {
			m[value] = true
			s = append(s, value)
		}
	}
	keys = s
	return keys
}

func valueInList(in string, list []string) bool {
	for _, x := range list {
		if x == in {
			return true
		}
	}
	return false
}
