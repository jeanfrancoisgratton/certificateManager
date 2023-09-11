// certificateManager : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/helpers/certsCreateHelpers.go
// 4/29/23 17:36:16

package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetIntValFromPrompt(prompt string, value *int) {
	var err error
	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s [%d]: ", prompt, *value)
	inputScanner.Scan()
	nval := inputScanner.Text()

	if nval != "" {
		*value, err = strconv.Atoi(nval)
		if err != nil {
			*value = 1
		}
	}
}

func GetBoolValFromPrompt(prompt string, value *bool) {
	fmt.Printf("%s (any values not starting with T or t will be treated as FALSE): ", prompt)
	bval := ""
	*value = false

	fmt.Scanln(&bval)
	if strings.HasSuffix(strings.ToLower(bval), "t") {
		*value = true
	}
}

func GetStringValFromPrompt(prompt string, value *string) {
	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s [%s]: ", prompt, *value)
	inputScanner.Scan()
	nval := inputScanner.Text()

	if nval != "" {
		*value = nval
	}
}

func GetStringSliceFromPrompt(prompt string, valuesPointer *[]string) {
	slice := *valuesPointer
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s\n", prompt)
	for i := range slice {
		fmt.Println("A value of '' is an empty string to be inserted in the slice (array)")
		fmt.Println("A value of '.' means that we're done and may exit the current loop")
		scanner.Scan()
		input := scanner.Text()
		if input == "" {
			slice[i] = ""
		}
		if input == "." {
			*valuesPointer = slice
			return
		}
	}
	*valuesPointer = slice
}

func GetKeyUsage(keys *[]string) {
	inputScanner := bufio.NewScanner(os.Stdin)
	ku := []string{"decipher only", "encipher only", "crl sign", "certs sign", "key agreement",
		"data encipherment", "key encipherment", "content commitment", "digital signature"}
	inputs := []string{}

	fmt.Println("The valid key usage values are:")
	for i, j := range ku {
		if i%5 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("'%s' ", White(j))
	}
	fmt.Println()
	for {
		input := ""
		fmt.Print("Please enter a value from the above list, just press ENTER to end : ")
		inputScanner.Scan()
		input = inputScanner.Text()
		if input == "" {
			break
		}
		if valueInList(input, ku) {
			inputs = append(inputs, input)
		}
	}
	// if the array is empty, we return a default value
	if len(inputs) == 0 {
		*keys = []string{"digital signature"}
		return
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
	*keys = s
}

func valueInList(in string, list []string) bool {
	for _, x := range list {
		if x == in {
			return true
		}
	}
	return false
}
