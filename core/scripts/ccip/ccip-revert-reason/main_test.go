package main

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/handler"
)

func TestRevertReason(t *testing.T) {
	errorCodeString := "ae236d9c000000000000000000000000000000000000000000000000c09c614ab4cba0de"

	decodedError, err := handler.DecodeErrorStringFromABI(errorCodeString)
	if err != nil {
		fmt.Printf("Error decoding error string: %v\n", err)
		return
	}

	fmt.Println(decodedError)
}
