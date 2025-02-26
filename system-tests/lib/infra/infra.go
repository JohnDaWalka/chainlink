package infra

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func Host(nodeIndex int, donName string, infraDetails types.InfraDetails) string {
	if infraDetails.InfraType == types.InfraType_CRIB {
		return fmt.Sprintf("%s-ccip-%d.main.stage.cldev.sh", infraDetails.Namespace, nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}
