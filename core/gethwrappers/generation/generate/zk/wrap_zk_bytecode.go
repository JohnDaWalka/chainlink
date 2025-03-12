package main

import (
	"fmt"
	"os"
)

const template = `// Code generated - DO NOT EDIT.
package %s

import "github.com/ethereum/go-ethereum/common"

var ZkBytecode = common.Hex2Bytes("%s")
`

func main() {
	srcFile := os.Args[1]
	dstFile := os.Args[2]
	pkgName := os.Args[3]

	fmt.Printf("Generating zk bytecode binding for %s\n", pkgName)

	bytecode, err := os.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}

	content := []byte(fmt.Sprintf(template, os.Args[3], string(bytecode)[2:]))

	err = os.WriteFile(dstFile, content, 0600)
	if err != nil {
		panic(err)
	}
}
