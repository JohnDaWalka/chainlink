package metadata

import (
	"encoding/json"
	"fmt"

	dstypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type VerifierProxyMetadata struct {
	FeeManagerAddress       string
	AccessControllerAddress string
	Verifiers               []string
}

type RecipientWeights struct {
	PayeeAddress  string
	SignerAddress string
	Weight        string
}
type VerifierConfig struct {
	Active       bool
	ConfigDigest types.ConfigDigest
	F            int
	Signers      []string
}
type VerifierMetadata struct {
	Active               bool
	VerifierProxyAddress string
	Configs              []VerifierConfig
}

// SerializedContractMetadata provides a generic container for contract metadata
// that can be serialized/deserialized to/from JSON
type SerializedContractMetadata struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

// Clone creates a copy of the SerializedContractMetadata
func (s SerializedContractMetadata) Clone() SerializedContractMetadata {
	contentCopy := make([]byte, len(s.Content))
	copy(contentCopy, s.Content)

	return SerializedContractMetadata{
		Type:    s.Type,
		Content: contentCopy,
	}
}

// ToVerifierMetadata converts the serialized metadata to VerifierMetadata
func (s SerializedContractMetadata) ToVerifierMetadata() (VerifierMetadata, error) {
	if s.Type != dstypes.Verifier.String() {
		return VerifierMetadata{}, fmt.Errorf("metadata is not of type verifier")
	}

	var metadata VerifierMetadata
	if err := json.Unmarshal(s.Content, &metadata); err != nil {
		return VerifierMetadata{}, err
	}

	return metadata, nil
}

// ToVerifierProxyMetadata converts the serialized metadata to VerifierProxyMetadata
func (s SerializedContractMetadata) ToVerifierProxyMetadata() (VerifierProxyMetadata, error) {
	if s.Type != dstypes.VerifierProxy.String() {
		return VerifierProxyMetadata{}, fmt.Errorf("metadata is not of type verifier_proxy")
	}

	var metadata VerifierProxyMetadata
	if err := json.Unmarshal(s.Content, &metadata); err != nil {
		return VerifierProxyMetadata{}, err
	}

	return metadata, nil
}

// NewVerifierMetadata creates a new SerializedContractMetadata from a VerifierMetadata
func NewVerifierMetadata(metadata VerifierMetadata) (SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return SerializedContractMetadata{}, err
	}

	return SerializedContractMetadata{
		Type:    dstypes.Verifier.String(),
		Content: content,
	}, nil
}

// NewVerifierProxyMetadata creates a new SerializedContractMetadata from a VerifierProxyMetadata
func NewVerifierProxyMetadata(metadata VerifierProxyMetadata) (SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return SerializedContractMetadata{}, err
	}

	return SerializedContractMetadata{
		Type: dstypes.VerifierProxy.String(), Content: content,
	}, nil
}
