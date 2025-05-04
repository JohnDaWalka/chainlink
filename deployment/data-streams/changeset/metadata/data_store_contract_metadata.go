package metadata

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	dstypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

type DataStreamsMutableDataStore = ds.MutableDataStore[SerializedContractMetadata, ds.DefaultMetadata]
type DataStreamsDataStore = ds.DataStore[SerializedContractMetadata, ds.DefaultMetadata]

type RewardManagerMetadata struct {
	FeeManagerAddress string
	RecipientWeights  []RecipientWeights
}

// TokenType represents the type of token used for fees
type TokenType string

const (
	// Native represents the native token of the blockchain (e.g., ETH on Ethereum)
	Native TokenType = "Native"
	// Link represents the LINK token
	Link TokenType = "Link"
)

// String returns the string representation of the TokenType
func (t TokenType) String() string {
	return string(t)
}

type FeeToken struct {
	TokenType TokenType
	Address   string
	Surcharge string
}

type StreamDiscounts struct {
	Stream       string
	DiscountType string
	TokenType    TokenType
	Value        string
}

type SubscriberDiscount struct {
	SubscriberAddress string
	SubscriberName    string
	StreamDiscounts   []StreamDiscounts
}

type ConfiguratorConfig struct{}

type DonId = string
type ConfiguratorMetadata struct {
	Configs         map[DonId]ConfiguratorConfig
	DeploymentBlock uint64
}
type ChannelConfigStoreMetadata struct {
	DeploymentBlock uint64
}

type FeeManagerMetadata struct {
	FeeTokens            []FeeToken
	RewardManagerAddress string
	VerifierProxyAddress string
}

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
		return VerifierMetadata{}, errors.New("metadata is not of type verifier")
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
		return VerifierProxyMetadata{}, errors.New("metadata is not of type verifier_proxy")
	}

	var metadata VerifierProxyMetadata
	if err := json.Unmarshal(s.Content, &metadata); err != nil {
		return VerifierProxyMetadata{}, err
	}

	return metadata, nil
}

// ToRewardManagerMetadata converts the serialized metadata to RewardManagerMetadata
func (s SerializedContractMetadata) ToRewardManagerMetadata() (RewardManagerMetadata, error) {
	if s.Type != dstypes.RewardManager.String() {
		return RewardManagerMetadata{}, errors.New("metadata is not of type reward_manager")
	}

	var metadata RewardManagerMetadata
	if err := json.Unmarshal(s.Content, &metadata); err != nil {
		return RewardManagerMetadata{}, err
	}

	return metadata, nil
}

// ToFeeManagerMetadata converts the serialized metadata to FeeManagerMetadata
func (s SerializedContractMetadata) ToFeeManagerMetadata() (FeeManagerMetadata, error) {
	if s.Type != dstypes.FeeManager.String() {
		return FeeManagerMetadata{}, errors.New("metadata is not of type fee_manager")
	}

	var metadata FeeManagerMetadata
	if err := json.Unmarshal(s.Content, &metadata); err != nil {
		return FeeManagerMetadata{}, err
	}

	return metadata, nil
}

// NewVerifierMetadata creates a new SerializedContractMetadata from a VerifierMetadata
func NewVerifierMetadata(metadata VerifierMetadata) (*SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return &SerializedContractMetadata{}, err
	}

	return &SerializedContractMetadata{
		Type:    dstypes.Verifier.String(),
		Content: content,
	}, nil
}

// NewVerifierProxyMetadata creates a new SerializedContractMetadata from a VerifierProxyMetadata
func NewVerifierProxyMetadata(metadata VerifierProxyMetadata) (*SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return &SerializedContractMetadata{}, err
	}

	return &SerializedContractMetadata{
		Type: dstypes.VerifierProxy.String(), Content: content,
	}, nil
}

// NewRewardManagerMetadata creates a new SerializedContractMetadata from a RewardManagerMetadata
func NewRewardManagerMetadata(metadata RewardManagerMetadata) (*SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return &SerializedContractMetadata{}, err
	}

	return &SerializedContractMetadata{
		Type:    dstypes.RewardManager.String(),
		Content: content,
	}, nil
}

// NewFeeManagerMetadata creates a new SerializedContractMetadata from a FeeManagerMetadata
func NewFeeManagerMetadata(metadata FeeManagerMetadata) (*SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return &SerializedContractMetadata{}, err
	}

	return &SerializedContractMetadata{
		Type:    dstypes.FeeManager.String(),
		Content: content,
	}, nil
}

// NewSerializedContractMetadata serializes any contract metadata
func NewSerializedContractMetadata[T any](metadata GenericContractMetadata[T]) (*SerializedContractMetadata, error) {
	content, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return &SerializedContractMetadata{
		Content: content,
	}, nil
}

// DeserializeMetadata converts a SerializedContractMetadata to a GenericContractMetadata of type T
func DeserializeMetadata[T any](serialized SerializedContractMetadata) (*GenericContractMetadata[T], error) {
	if len(serialized.Content) == 0 {
		return nil, fmt.Errorf("empty content in serialized metadata")
	}

	var result GenericContractMetadata[T]
	if err := json.Unmarshal(serialized.Content, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contract metadata: %w", err)
	}

	return &result, nil
}

type ContractMetadata struct {
	DeployBlock uint64 `json:"deployBlock"`
}

// GenericContractMetadata is a generic container for any view type
// Use as Content for SerializedContractMetadata
type GenericContractMetadata[T any] struct {
	Metadata ContractMetadata `json:"metadata"`
	// View is intended to be populated with the contract's view usually after state change to have an off chain representation of the contract
	View T `json:"view"`
}
