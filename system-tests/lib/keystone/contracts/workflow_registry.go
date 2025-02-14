package contracts

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	libnet "github.com/smartcontractkit/chainlink/system-tests/lib/net"
)

func RegisterWorkflow(t *testing.T, sc *seth.Client, workflowRegistryAddr common.Address, donID uint32, workflowName, binaryURL, configURL string) {
	require.NotEmpty(t, binaryURL)
	workFlowData, err := libnet.DownloadAndDecodeBase64(binaryURL)
	require.NoError(t, err, "failed to download and decode workflow binary")

	var configData []byte
	if configURL != "" {
		configData, err = libnet.Download(configURL)
		require.NoError(t, err, "failed to download workflow config")
	}

	// use non-encoded workflow name
	workflowID, idErr := generateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, "")
	require.NoError(t, idErr, "failed to generate workflow ID")

	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	require.NoError(t, err, "failed to create workflow registry instance")

	// use non-encoded workflow name
	_, decodeErr := sc.Decode(workflowRegistryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), binaryURL, configURL, ""))
	require.NoError(t, decodeErr, "failed to register workflow")
}

func generateWorkflowIDFromStrings(owner string, name string, workflow []byte, config []byte, secretsURL string) (string, error) {
	ownerWithoutPrefix := owner
	if strings.HasPrefix(owner, "0x") {
		ownerWithoutPrefix = owner[2:]
	}

	ownerb, err := hex.DecodeString(ownerWithoutPrefix)
	if err != nil {
		return "", err
	}

	wid, err := pkgworkflows.GenerateWorkflowID(ownerb, name, workflow, config, secretsURL)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(wid[:]), nil
}
