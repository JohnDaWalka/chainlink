package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/genwrapper"
	zksyncwrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/zksync"
)

var (
	rootDir = "../../../"
)

func main() {
	project := os.Args[1]
	className := os.Args[2]
	pkgName := os.Args[3]

	var outDirSuffix string
	if len(os.Args) >= 5 {
		outDirSuffix = os.Args[4]
	}

	// Once vrf is moved to its own subfolder we can delete this rootDir override.
	if project == "vrf" || project == "automation" {
		rootDir = "../../"
	}

	abiPath := rootDir + "contracts/solc/" + project + "/" + className + "/" + className + ".sol/" + className + ".abi.json"
	binPath := rootDir + "contracts/solc/" + project + "/" + className + "/" + className + ".sol/" + className + ".bin"

	genwrapper.GenWrapper(abiPath, binPath, className, pkgName, outDirSuffix)

	if pkgName == "link_token" {
		zksolcBinPath := rootDir + "contracts/zkout/" + className + ".sol/" + className + ".json"

		zksyncwrapper.WrapZksyncDeploy(zksolcBinPath, className, pkgName, outDirSuffix)
	}
}
