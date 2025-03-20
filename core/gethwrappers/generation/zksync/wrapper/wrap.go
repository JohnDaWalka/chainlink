package main

import (
	"os"

	zksyncwrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/zksync"
)

const rootDir = "../../../"

func main() {
	className := os.Args[1]
	pkgName := os.Args[2]

	var outDirSuffix string
	if len(os.Args) >= 4 {
		outDirSuffix = os.Args[3]
	}

	zksolcBinPath := rootDir + "contracts/zkout/" + className + ".sol/" + className + ".json"

	zksyncwrapper.WrapZksyncDeploy(zksolcBinPath, className, pkgName, outDirSuffix)
}
