package main

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/cmd/environment"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/cmd/root"
)

func init() {
	root.RootCmd.AddCommand(environment.EnvironmentCmd)
}

func main() {
	if err := root.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
