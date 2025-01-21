package capabilities_test

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
)

func TestDebugWorkflow(t *testing.T) {
	pkey := os.Getenv("PRIVATE_KEY")
	require.NotEmpty(t, pkey)

	sc, err := seth.NewClientBuilder().
		WithRpcUrl("ws://localhost:8545").
		WithPrivateKeys([]string{pkey}).
		WithGethWrappersFolders([]string{"../../core/gethwrappers"}).
		Build()
	require.NoError(t, err)

	workflowName := "abcdefgasd"
	var workflowNameBytes [10]byte
	var HashTruncateName = func(name string) string {
		// Compute SHA-256 hash of the input string
		hash := sha256.Sum256([]byte(name))

		// Encode as hex to ensure UTF8
		var hashBytes []byte = hash[:]
		resultHex := hex.EncodeToString(hashBytes)

		// Truncate to 10 bytes
		truncated := []byte(resultHex)[:10]
		return string(truncated)
	}

	truncated := HashTruncateName(workflowName)
	fmt.Println("Truncated name: ", truncated)

	copy(workflowNameBytes[:], []byte(truncated))

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(common.HexToAddress("0xa513E6E4b8f2a923D98304ec87F64353C4D5C853"), sc.Client)
	require.NoError(t, err)

	tx, err := feedsConsumerInstance.SetConfig(
		sc.NewTXOpts(),
		[]common.Address{sc.MustGetRootKeyAddress()},
		[]common.Address{sc.MustGetRootKeyAddress()},
		[][10]byte{workflowNameBytes},
	)
	_, decodedErr := sc.Decode(tx, err)
	require.NoError(t, decodedErr)

	metadataB64 := "AOarLZa5Pc3v1orPvDdPFQhzKrjgYkuhYsMK9qC6Mr1kMGRmMjI5NTAx85/W5RqtiPb0zmq4gnJ5z/+5ImYAAQ"
	metaDecoded, err := base64.RawStdEncoding.DecodeString(metadataB64)
	require.NoError(t, err)

	reportB64 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQGL/ohAcABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC7XBYsgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZ4+Txw"
	reportDecoded, err := base64.RawStdEncoding.DecodeString(reportB64)
	require.NoError(t, err)

	tx, err = feedsConsumerInstance.OnReport(sc.NewTXOpts(), metaDecoded, reportDecoded)
	_, decodedErr = sc.Decode(tx, err)
	require.NoError(t, decodedErr)
}
