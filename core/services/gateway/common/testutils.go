package common

import (
	"crypto/ecdsa"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TestNode struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
}

func NewTestNodes(t *testing.T, n int) []TestNode {
	nodes := make([]TestNode, n)
	for i := 0; i < n; i++ {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		address := strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		nodes[i] = TestNode{Address: address, PrivateKey: privateKey}
	}
	return nodes
}

// Sign Gateway messages
// signatures are over the following data:
// 1. MessageId aligned to 128 bytes
// 2. Method aligned to 64 bytes
// 3. DonId aligned to 64 bytes
// 4. Receiver (in hex) aligned to 42 bytes
// 5. Payload (raw bytes before parsing)
func Sign(m *gateway.Message, privateKey *ecdsa.PrivateKey) error {
	if m == nil {
		return errors.New("nil message")
	}
	rawData := gateway.GetRawMessageBody(&m.Body)

	signature, err := SignData(privateKey, rawData...)
	if err != nil {
		return err
	}
	m.Signature = utils.StringToHex(string(signature))
	m.Body.Sender = strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	return nil
}
