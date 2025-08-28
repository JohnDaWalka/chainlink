package vault

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
)

func TestAggregator_Valid_Signatures(t *testing.T) {
	signers := []string{
		"d6da96fe596705b32bc3a0e11cdefad77feaad79000000000000000000000000",
		"327aa349c9718cd36c877d1e90458fe1929768ad000000000000000000000000",
		"e9bf394856d73402b30e160d0e05c847796f0e29000000000000000000000000",
		"efd5bdb6c3256f04489a6ca32654d547297f48b9000000000000000000000000",
	}
	nodes := []capabilities.Node{}
	for _, s := range signers {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		nodes = append(nodes, capabilities.Node{Signer: [32]byte(b)})
	}
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	sig1, err := hex.DecodeString("d1067844e2849b404d903730c4cae19f090d53a578a1e8dc16ecbdc0285c1f186599108abbe0073b78bc148a6504907474ed3a6881df917e6d142cff70acfb5900")
	require.NoError(t, err)
	sig2, err := hex.DecodeString("c7517c188d297093a6f602046fad7feafe19454ee9dc269b19c8e6c01268037d1f7b423eeecbc495dd2d9a65e106bc3eab849ddfd74a10cbd4ad50c7d953bd4b01")
	require.NoError(t, err)

	rm := json.RawMessage([]byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`))
	sor := vault.SignedOCRResponse{
		Payload: rm,
		Context: ctx,
		Signatures: [][]byte{
			sig1,
			sig2,
		},
	}
	rawResp, err := json.Marshal(sor)
	require.NoError(t, err)

	currResp := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vault.MethodSecretsGet,
		Result:  (*json.RawMessage)(&rawResp),
	}
	ar := &activeRequest{
		responses: []*jsonrpc.Response[json.RawMessage]{currResp},
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), ar, currResp)
	require.NoError(t, err)
	assert.Equal(t, currResp, resp)
}

func mustRandom(length int) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return randomBytes
}

func newMessage(t *testing.T) *jsonrpc.Response[json.RawMessage] {
	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	rm := json.RawMessage([]byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`))
	sor := vault.SignedOCRResponse{
		Payload: rm,
		Context: ctx,
		Signatures: [][]byte{
			mustRandom(65),
			mustRandom(65),
		},
	}
	rawResp, err := json.Marshal(sor)
	require.NoError(t, err)

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vault.MethodSecretsGet,
		Result:  (*json.RawMessage)(&rawResp),
	}
}

func TestAggregator_Valid_FallsBackToQuorum(t *testing.T) {
	// No valid signers
	nodes := []capabilities.Node{}
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	currResp := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vault.MethodSecretsGet,
		Result:  (*json.RawMessage)(nil),
		Error: &jsonrpc.WireError{
			Code:    123,
			Message: "some error",
		},
	}
	ar := &activeRequest{
		responses: []*jsonrpc.Response[json.RawMessage]{currResp, currResp, currResp},
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), ar, currResp)
	require.NoError(t, err)
	assert.Equal(t, currResp, resp)
}

func TestAggregator_Valid_FallsBackToQuorum_ExcludesSignaturesInSha(t *testing.T) {
	// No valid signers
	nodes := []capabilities.Node{}
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	oldResp1 := newMessage(t)
	oldResp2 := newMessage(t)
	currResp := newMessage(t)
	ar := &activeRequest{
		responses: []*jsonrpc.Response[json.RawMessage]{oldResp1, oldResp2, currResp},
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), ar, currResp)
	require.NoError(t, err)
	assert.Equal(t, currResp, resp)
}

func TestAggregator_InsufficientResponses(t *testing.T) {
	mcr := &mockCapabilitiesRegistry{F: 1}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	rm := json.RawMessage([]byte(`{}`))
	currResp := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vault.MethodSecretsGet,
		Result:  &rm,
	}
	ar := &activeRequest{
		responses: []*jsonrpc.Response[json.RawMessage]{currResp},
	}
	_, err := agg.Aggregate(t.Context(), logger.Test(t), ar, currResp)
	require.ErrorContains(t, err, "insufficient valid responses to reach quorum")
}
