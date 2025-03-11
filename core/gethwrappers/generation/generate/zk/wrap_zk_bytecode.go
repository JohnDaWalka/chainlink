package main

import (
	"fmt"
	"os"
)

func main() {
	bytecode, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading source file: %v\n", err)
		os.Exit(1)
	}

	template := fmt.Sprintf(`// Generated code - DO NOT EDIT.
package %s

import "github.com/ethereum/go-ethereum/common"

var ZkBytecode = common.Hex2Bytes("%s")
`, os.Args[3], string(bytecode)[2:])

	err = os.WriteFile(os.Args[2], []byte(template), 0600)
	if err != nil {
		fmt.Printf("Error writing destination file: %v\n", err)
		os.Exit(1)
	}
}
