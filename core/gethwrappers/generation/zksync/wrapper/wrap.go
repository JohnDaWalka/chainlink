package main

import (
	"os"
	"path/filepath"

	zksyncwrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/zksync"
)

const rootDir = "../../../"

func main() {
	className := os.Args[2]
	pkgName := os.Args[3]

	outDirSuffix := ""
	if len(os.Args) > 4 {
		outDirSuffix = os.Args[4]
	}

	zksolcBinPath := filepath.Join(rootDir, "contracts", "zkout", className+".sol", className+".json")
	outPath := filepath.Join("generated", outDirSuffix, pkgName, pkgName+"_zksync.go")

	zksyncwrapper.WrapZksyncDeploy(zksolcBinPath, className, pkgName, outPath)
}
