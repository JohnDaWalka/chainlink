// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package message_transformer_offramp

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type InternalAny2EVMRampMessage struct {
	Header       InternalRampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     common.Address
	GasLimit     *big.Int
	TokenAmounts []InternalAny2EVMTokenTransfer
}

type InternalAny2EVMTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  common.Address
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

type InternalExecutionReport struct {
	SourceChainSelector uint64
	Messages            []InternalAny2EVMRampMessage
	OffchainTokenData   [][][]byte
	Proofs              [][32]byte
	ProofFlagBits       *big.Int
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type MultiOCR3BaseConfigInfo struct {
	ConfigDigest                   [32]byte
	F                              uint8
	N                              uint8
	IsSignatureVerificationEnabled bool
}

type MultiOCR3BaseOCRConfig struct {
	ConfigInfo   MultiOCR3BaseConfigInfo
	Signers      []common.Address
	Transmitters []common.Address
}

type MultiOCR3BaseOCRConfigArgs struct {
	ConfigDigest                   [32]byte
	OcrPluginType                  uint8
	F                              uint8
	IsSignatureVerificationEnabled bool
	Signers                        []common.Address
	Transmitters                   []common.Address
}

type OffRampDynamicConfig struct {
	FeeQuoter                               common.Address
	PermissionLessExecutionThresholdSeconds uint32
	IsRMNVerificationDisabled               bool
	MessageInterceptor                      common.Address
}

type OffRampGasLimitOverride struct {
	ReceiverExecutionGasLimit *big.Int
	TokenGasOverrides         []uint32
}

type OffRampSourceChainConfig struct {
	Router    common.Address
	IsEnabled bool
	MinSeqNr  uint64
	OnRamp    []byte
}

type OffRampSourceChainConfigArgs struct {
	Router              common.Address
	SourceChainSelector uint64
	IsEnabled           bool
	OnRamp              []byte
}

type OffRampStaticConfig struct {
	ChainSelector        uint64
	GasForCallExactCheck uint16
	RmnRemote            common.Address
	TokenAdminRegistry   common.Address
	NonceManager         common.Address
}

var MessageTransformerOffRampMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applySourceChainConfigUpdates\",\"inputs\":[{\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"commit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeSingleMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllSourceChainConfigs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExecutionState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestPriceSequenceNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMerkleRoot\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageTransformerAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSourceChainConfig\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"ocrConfig\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"components\":[{\"name\":\"configInfo\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"n\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyExecute\",\"inputs\":[{\"name\":\"reports\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\",\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"components\":[{\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOCR3Configs\",\"inputs\":[{\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AlreadyAttempted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitReportAccepted\",\"inputs\":[{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"F\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DynamicConfigSet\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecutionStateChanged\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"state\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RootRemoved\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedAlreadyExecutedMessage\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedReportExecution\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainConfigSet\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sourceConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainSelectorAdded\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaticConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CanOnlySelfCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CommitOnRampMismatch\",\"inputs\":[{\"name\":\"reportOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"configOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyReport\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ExecutionError\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ForkedChain\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[{\"name\":\"errorType\",\"type\":\"uint8\",\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\"}]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInterval\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionGasLimit\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"newLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionTokenGasOverride\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oldLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverride\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidMessageDestChainSelector\",\"inputs\":[{\"name\":\"messageDestChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidNewState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newState\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}]},{\"type\":\"error\",\"name\":\"InvalidOnRampUpdate\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LeavesCannotBeEmpty\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionGasAmountCountMismatch\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ManualExecutionGasLimitMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionNotYetEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageTransformError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MessageValidationError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonUniqueSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotACompatiblePool\",\"inputs\":[{\"name\":\"notPool\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OracleCannotBeZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReceiverError\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ReleaseOrMintBalanceMismatch\",\"inputs\":[{\"name\":\"amountReleased\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePre\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RootAlreadyCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"RootNotCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignaturesOutOfRegistration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SourceChainNotEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SourceChainSelectorMismatch\",\"inputs\":[{\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StaleCommitReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaticConfigCannotBeChanged\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TokenDataMismatch\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TokenHandlingError\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnexpectedTokenData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongMessageLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WrongNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroChainSelectorNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461088f57617bd9803803809161001d82856108c5565b8339810190808203610160811261088f5760a0811261088f5760405160a081016001600160401b038111828210176108945760405261005b836108e8565b815260208301519061ffff8216820361088f57602081019182526040840151936001600160a01b038516850361088f576040820194855261009e606082016108fc565b946060830195865260806100b38184016108fc565b84820190815295609f19011261088f57604051936100d0856108aa565b6100dc60a084016108fc565b855260c08301519363ffffffff8516850361088f576020860194855261010460e08501610910565b966040870197885261011961010086016108fc565b606088019081526101208601519095906001600160401b03811161088f5781018b601f8201121561088f5780519b6001600160401b038d11610894578c60051b91604051809e6020850161016d90836108c5565b81526020019281016020019082821161088f5760208101935b82851061078f57505050505061014061019f91016108fc565b98331561077e57600180546001600160a01b031916331790554660805284516001600160a01b031615801561076c575b801561075a575b6107385782516001600160401b0316156107495782516001600160401b0390811660a090815286516001600160a01b0390811660c0528351811660e0528451811661010052865161ffff90811661012052604080519751909416875296519096166020860152955185169084015251831660608301525190911660808201527fb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b79190a182516001600160a01b031615610738579151600480548351865160ff60c01b90151560c01b1663ffffffff60a01b60a09290921b919091166001600160a01b039485166001600160c81b0319909316831717179091558351600580549184166001600160a01b031990921691909117905560408051918252925163ffffffff166020820152935115159184019190915290511660608201529091907fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee90608090a16000915b81518310156106805760009260208160051b8401015160018060401b036020820151169081156106715780516001600160a01b031615610662578186526008602052604086206060820151916001820192610399845461091d565b610603578254600160a81b600160e81b031916600160a81b1783556040518581527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb990602090a15b805180159081156105d8575b506105c9578051906001600160401b0382116105b55761040d855461091d565b601f8111610570575b50602090601f83116001146104f857918091600080516020617bb98339815191529695949360019a9b9c926104ed575b5050600019600383901b1c191690881b1783555b60408101518254915160a089811b8a9003801960ff60a01b1990951693151590911b60ff60a01b1692909217929092169116178155610498846109da565b506104e26040519283926020845254888060a01b038116602085015260ff8160a01c1615156040850152888060401b039060a81c16606084015260808084015260a0830190610957565b0390a201919061033e565b015190503880610446565b858b52818b20919a601f198416905b8181106105585750916001999a9b8492600080516020617bb9833981519152989796958c951061053f575b505050811b01835561045a565b015160001960f88460031b161c19169055388080610532565b828d0151845560209c8d019c60019094019301610507565b858b5260208b20601f840160051c810191602085106105ab575b601f0160051c01905b8181106105a05750610416565b8b8155600101610593565b909150819061058a565b634e487b7160e01b8a52604160045260248afd5b6342bcdf7f60e11b8952600489fd5b9050602082012060405160208101908b8252602081526105f96040826108c5565b51902014386103ed565b825460a81c6001600160401b03166001141580610634575b156103e157632105803760e11b89526004859052602489fd5b5060405161064d816106468188610957565b03826108c5565b6020815191012081516020830120141561061b565b6342bcdf7f60e11b8652600486fd5b63c656089560e01b8652600486fd5b6001600160a01b0381161561073857600b8054600160401b600160e01b031916604092831b600160401b600160e01b03161790555161714b9081610a6e823960805181614367015260a0518181816105590152615582015260c0518181816105bc015281816136c601528181613c52015261550f015260e0518181816105f80152615f3f01526101005181818161063401526159c101526101205181818161058001528181612cfc015281816160730152616da90152f35b6342bcdf7f60e11b60005260046000fd5b63c656089560e01b60005260046000fd5b5081516001600160a01b0316156101d6565b5080516001600160a01b0316156101cf565b639b15e16f60e01b60005260046000fd5b84516001600160401b03811161088f5782016080818603601f19011261088f57604051906107bc826108aa565b60208101516001600160a01b038116810361088f5782526107df604082016108e8565b60208301526107f060608201610910565b604083015260808101516001600160401b03811161088f57602091010185601f8201121561088f5780516001600160401b0381116108945760405191610840601f8301601f1916602001846108c5565b818352876020838301011161088f5760005b82811061087a5750509181600060208096949581960101526060820152815201940193610186565b80602080928401015182828701015201610852565b600080fd5b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761089457604052565b601f909101601f19168101906001600160401b0382119082101761089457604052565b51906001600160401b038216820361088f57565b51906001600160a01b038216820361088f57565b5190811515820361088f57565b90600182811c9216801561094d575b602083101461093757565b634e487b7160e01b600052602260045260246000fd5b91607f169161092c565b600092918154916109678361091d565b80835292600181169081156109bd575060011461098357505050565b60009081526020812093945091925b8383106109a3575060209250010190565b600181602092949394548385870101520191019190610992565b915050602093945060ff929192191683830152151560051b010190565b80600052600760205260406000205415600014610a675760065468010000000000000000811015610894576001810180600655811015610a51577ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0181905560065460009182526007602052604090912055600190565b634e487b7160e01b600052603260045260246000fd5b5060009056fe6080604052600436101561001257600080fd5b60003560e01c806304666f9c1461016757806306285c69146101625780631056edee1461015d578063181f5a77146101585780633f4b04aa146101535780635215505b1461014e5780635e36480c146101495780635e7bb0081461014457806360987c201461013f5780637437ff9f1461013a57806379ba5097146101355780637edf52f41461013057806385572ffb1461012b5780638da5cb5b14610126578063c673e58414610121578063ccd37ba31461011c578063de5e0b9a14610117578063e9d68a8e14610112578063f2fde38b1461010d578063f58e03fc146101085763f716f99f1461010357600080fd5b611f07565b611dac565b611cb8565b611bec565b611b2f565b611a87565b6119be565b611880565b61180e565b611593565b6114aa565b6113bc565b611320565b6110d9565b610b04565b610966565b61080d565b610790565b6106c0565b610507565b6103a8565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6080810190811067ffffffffffffffff8211176101b757604052565b61016c565b60a0810190811067ffffffffffffffff8211176101b757604052565b6040810190811067ffffffffffffffff8211176101b757604052565b6060810190811067ffffffffffffffff8211176101b757604052565b60c0810190811067ffffffffffffffff8211176101b757604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176101b757604052565b6040519061027c60c08361022c565b565b6040519061027c60a08361022c565b6040519061027c6101008361022c565b6040519061027c60408361022c565b67ffffffffffffffff81116101b75760051b60200190565b73ffffffffffffffffffffffffffffffffffffffff8116036102e257565b600080fd5b67ffffffffffffffff8116036102e257565b359061027c826102e7565b801515036102e257565b359061027c82610304565b67ffffffffffffffff81116101b757601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b92919261035f82610319565b9161036d604051938461022c565b8294818452818301116102e2578281602093846000960137010152565b9080601f830112156102e2578160206103a593359101610353565b90565b346102e25760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760043567ffffffffffffffff81116102e257366023820112156102e257806004013590610403826102ac565b90610411604051928361022c565b8282526024602083019360051b820101903682116102e25760248101935b8285106104415761043f84612082565b005b843567ffffffffffffffff81116102e257820160807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc82360301126102e2576040519161048d8361019b565b602482013561049b816102c4565b835260448201356104ab816102e7565b602084015260648201356104be81610304565b604084015260848201359267ffffffffffffffff84116102e2576104ec60209493602486953692010161038a565b606082015281520194019361042f565b60009103126102e257565b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25761053e6123e0565b506106bc60405161054e816101bc565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815261ffff7f000000000000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016606082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016608082015260405191829182919091608073ffffffffffffffffffffffffffffffffffffffff8160a084019567ffffffffffffffff815116855261ffff6020820151166020860152826040820151166040860152826060820151166060860152015116910152565b0390f35b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e257602073ffffffffffffffffffffffffffffffffffffffff600b5460401c16604051908152f35b6040519061072460208361022c565b60008252565b60005b83811061073d5750506000910152565b818101518382015260200161072d565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f6020936107898151809281875287808801910161072a565b0116010190565b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e2576106bc60408051906107d1818361022c565b601182527f4f666652616d7020312e362e302d64657600000000000000000000000000000060208301525191829160208352602083019061074d565b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e257602067ffffffffffffffff600b5416604051908152f35b90608060606103a59373ffffffffffffffffffffffffffffffffffffffff815116845260208101511515602085015267ffffffffffffffff6040820151166040850152015191816060820152019061074d565b6040810160408252825180915260206060830193019060005b818110610946575050506020818303910152815180825260208201916020808360051b8301019401926000915b8383106108fb57505050505090565b9091929394602080610937837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086600196030187528951610853565b970193019301919392906108ec565b825167ffffffffffffffff168552602094850194909201916001016108bf565b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e2576006546109a1816102ac565b906109af604051928361022c565b8082527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06109dc826102ac565b0160005b818110610aa25750506109f281612432565b9060005b818110610a0e5750506106bc604051928392836108a6565b80610a46610a2d610a206001946153cc565b67ffffffffffffffff1690565b610a3783876124c2565b9067ffffffffffffffff169052565b610a86610a81610a67610a5984886124c2565b5167ffffffffffffffff1690565b67ffffffffffffffff166000526008602052604060002090565b6125e6565b610a9082876124c2565b52610a9b81866124c2565b50016109f6565b602090610aad61240b565b828287010152016109e0565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60041115610af257565b610ab9565b906004821015610af25752565b346102e25760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e2576020610b56600435610b44816102e7565b60243590610b51826102e7565b6126a5565b610b636040518092610af7565bf35b91908260a09103126102e257604051610b7d816101bc565b6080808294803584526020810135610b94816102e7565b60208501526040810135610ba7816102e7565b60408501526060810135610bba816102e7565b6060850152013591610bcb836102e7565b0152565b359061027c826102c4565b63ffffffff8116036102e257565b359061027c82610bda565b81601f820112156102e257803590610c0a826102ac565b92610c18604051948561022c565b82845260208085019360051b830101918183116102e25760208101935b838510610c4457505050505090565b843567ffffffffffffffff81116102e257820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082860301126102e25760405191610c90836101bc565b602082013567ffffffffffffffff81116102e257856020610cb39285010161038a565b83526040820135610cc3816102c4565b6020840152610cd460608301610be8565b604084015260808201359267ffffffffffffffff84116102e25760a083610d0288602080988198010161038a565b606084015201356080820152815201940193610c35565b919091610140818403126102e257610d2f61026d565b92610d3a8183610b65565b845260a082013567ffffffffffffffff81116102e25781610d5c91840161038a565b602085015260c082013567ffffffffffffffff81116102e25781610d8191840161038a565b6040850152610d9260e08301610bcf565b6060850152610100820135608085015261012082013567ffffffffffffffff81116102e257610dc19201610bf3565b60a0830152565b9080601f830112156102e2578135610ddf816102ac565b92610ded604051948561022c565b81845260208085019260051b820101918383116102e25760208201905b838210610e1957505050505090565b813567ffffffffffffffff81116102e257602091610e3c87848094880101610d19565b815201910190610e0a565b81601f820112156102e257803590610e5e826102ac565b92610e6c604051948561022c565b82845260208085019360051b830101918183116102e25760208101935b838510610e9857505050505090565b843567ffffffffffffffff81116102e257820183603f820112156102e2576020810135610ec4816102ac565b91610ed2604051938461022c565b8183526020808085019360051b83010101918683116102e25760408201905b838210610f0b575050509082525060209485019401610e89565b813567ffffffffffffffff81116102e257602091610f308a848080958901010161038a565b815201910190610ef1565b929190610f47816102ac565b93610f55604051958661022c565b602085838152019160051b81019283116102e257905b828210610f7757505050565b8135815260209182019101610f6b565b9080601f830112156102e2578160206103a593359101610f3b565b81601f820112156102e257803590610fb9826102ac565b92610fc7604051948561022c565b82845260208085019360051b830101918183116102e25760208101935b838510610ff357505050505090565b843567ffffffffffffffff81116102e257820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082860301126102e25761103a61027e565b91611047602083016102f9565b8352604082013567ffffffffffffffff81116102e25785602061106c92850101610dc8565b6020840152606082013567ffffffffffffffff81116102e25785602061109492850101610e47565b604084015260808201359267ffffffffffffffff84116102e25760a0836110c2886020809881980101610f87565b606084015201356080820152815201940193610fe4565b346102e25760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760043567ffffffffffffffff81116102e257611128903690600401610fa2565b60243567ffffffffffffffff81116102e257366023820112156102e257806004013591611154836102ac565b91611162604051938461022c565b8383526024602084019460051b820101903682116102e25760248101945b8286106111915761043f85856126ee565b853567ffffffffffffffff81116102e2578201366043820112156102e25760248101356111bd816102ac565b916111cb604051938461022c565b818352602060248185019360051b83010101903682116102e25760448101925b828410611205575050509082525060209586019501611180565b833567ffffffffffffffff81116102e25760249083010160407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082360301126102e25760405190611255826101d8565b60208101358252604081013567ffffffffffffffff81116102e257602091010136601f820112156102e25780359061128c826102ac565b9161129a604051938461022c565b80835260208084019160051b830101913683116102e257602001905b8282106112d557505050918160209384809401528152019301926111eb565b6020809183356112e481610bda565b8152019101906112b6565b9181601f840112156102e25782359167ffffffffffffffff83116102e2576020808501948460051b0101116102e257565b346102e25760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760043567ffffffffffffffff81116102e25761136f903690600401610d19565b60243567ffffffffffffffff81116102e25761138f9036906004016112ef565b916044359267ffffffffffffffff84116102e2576113b461043f9436906004016112ef565b939092612b85565b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e2576113f3612f39565b506106bc6040516114038161019b565b60ff60045473ffffffffffffffffffffffffffffffffffffffff8116835263ffffffff8160a01c16602084015260c01c161515604082015273ffffffffffffffffffffffffffffffffffffffff60055416606082015260405191829182919091606073ffffffffffffffffffffffffffffffffffffffff81608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760005473ffffffffffffffffffffffffffffffffffffffff81163303611569577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346102e25760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760006040516115d08161019b565b6004356115dc816102c4565b81526024356115ea81610bda565b60208201526044356115fb81610304565b604082015260643561160c816102c4565b60608201526116196140d4565b73ffffffffffffffffffffffffffffffffffffffff815116156117e6576117e0816116b773ffffffffffffffffffffffffffffffffffffffff7fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee94511673ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff00000000000000000000000000000000000000006004541617600455565b60208101516004547fffffffffffffff0000000000ffffffffffffffffffffffffffffffffffffffff77ffffffff000000000000000000000000000000000000000078ff0000000000000000000000000000000000000000000000006040860151151560c01b169360a01b169116171760045561178f61174e606083015173ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff00000000000000000000000000000000000000006005541617600555565b60405191829182919091606073ffffffffffffffffffffffffffffffffffffffff81608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b0390a180f35b6004827f8579befe000000000000000000000000000000000000000000000000000000008152fd5b346102e25760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760043567ffffffffffffffff81116102e2577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60a091360301126102e257600080fd5b346102e25760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e257602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b6004359060ff821682036102e257565b359060ff821682036102e257565b906020808351928381520192019060005b81811061190e5750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101611901565b906103a59160208152606082518051602084015260ff602082015116604084015260ff60408201511682840152015115156080820152604061198b602084015160c060a085015260e08401906118f0565b9201519060c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828503019101526118f0565b346102e25760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760ff6119f76118d2565b606060408051611a06816101f4565b611a0e612f39565b815282602082015201521660005260026020526106bc60406000206003611a7660405192611a3b846101f4565b611a4481612f5e565b8452604051611a6181611a5a8160028601612f97565b038261022c565b6020850152611a5a6040518094819301612f97565b60408201526040519182918261193a565b346102e25760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e257600435611ac2816102e7565b67ffffffffffffffff6024359116600052600a6020526040600020906000526020526020604060002054604051908152f35b906004916044116102e257565b9181601f840112156102e25782359167ffffffffffffffff83116102e257602083818601950101116102e257565b346102e25760c07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e257611b6736611af4565b60443567ffffffffffffffff81116102e257611b87903690600401611b01565b60649291923567ffffffffffffffff81116102e257611baa9036906004016112ef565b608435949167ffffffffffffffff86116102e257611bcf61043f9636906004016112ef565b94909360a43596613681565b9060206103a5928181520190610853565b346102e25760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25767ffffffffffffffff600435611c30816102e7565b611c3861240b565b501660005260086020526106bc60406000206001611ca760405192611c5c8461019b565b67ffffffffffffffff815473ffffffffffffffffffffffffffffffffffffffff8116865260ff8160a01c161515602087015260a81c166040850152611a5a6040518094819301612529565b606082015260405191829182611bdb565b346102e25760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25773ffffffffffffffffffffffffffffffffffffffff600435611d08816102c4565b611d106140d4565b16338114611d8257807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102e25760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e257611de436611af4565b60443567ffffffffffffffff81116102e257611e04903690600401611b01565b918282016020838203126102e25782359067ffffffffffffffff82116102e257611e2f918401610fa2565b604051906020611e3f818461022c565b600083527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810160005b818110611e915750505061043f9491611e81916143c1565b611e89613ce0565b928392614cf3565b60608582018401528201611e69565b9080601f830112156102e2578135611eb7816102ac565b92611ec5604051948561022c565b81845260208085019260051b8201019283116102e257602001905b828210611eed5750505090565b602080918335611efc816102c4565b815201910190611ee0565b346102e25760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102e25760043567ffffffffffffffff81116102e257366023820112156102e257806004013590611f62826102ac565b90611f70604051928361022c565b8282526024602083019360051b820101903682116102e25760248101935b828510611f9e5761043f84613cfc565b843567ffffffffffffffff81116102e257820160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc82360301126102e257611fe561026d565b9160248201358352611ff9604483016118e2565b602084015261200a606483016118e2565b604084015261201b6084830161030e565b606084015260a482013567ffffffffffffffff81116102e2576120449060243691850101611ea0565b608084015260c48201359267ffffffffffffffff84116102e257612072602094936024869536920101611ea0565b60a0820152815201940193611f8e565b61208a6140d4565b60005b81518110156123dc576120a081836124c2565b51906120b7602083015167ffffffffffffffff1690565b9167ffffffffffffffff83169081156123b2576121076120ee6120ee835173ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1690565b156122e45761212a8467ffffffffffffffff166000526008602052604060002090565b90606081015191600181019561214087546124d6565b612324576121c87ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9916121ad8475010000000000000000000000000000000000000000007fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff825416179055565b60405167ffffffffffffffff90911681529081906020820190565b0390a15b8251801590811561230e575b506122e4576122c56122846122db936122147f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b9660019a61419c565b61226a6122246040830151151590565b85547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690151560a01b74ff000000000000000000000000000000000000000016178555565b5173ffffffffffffffffffffffffffffffffffffffff1690565b829073ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055565b6122ce846168b2565b50604051918291826142c7565b0390a20161208d565b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b9050602084012061231d61411f565b14386121d8565b600167ffffffffffffffff612345845467ffffffffffffffff9060a81c1690565b16141580612393575b61235857506121cc565b7f420b006e0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff1660045260246000fd5b6000fd5b5061239d876125cb565b6020815191012084516020860120141561234e565b7fc65608950000000000000000000000000000000000000000000000000000000060005260046000fd5b5050565b604051906123ed826101bc565b60006080838281528260208201528260408201528260608201520152565b604051906124188261019b565b606080836000815260006020820152600060408201520152565b9061243c826102ac565b612449604051918261022c565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe061247782946102ac565b0190602036910137565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8051156124bd5760200190565b612481565b80518210156124bd5760209160051b010190565b90600182811c9216801561251f575b60208310146124f057565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f16916124e5565b60009291815491612539836124d6565b808352926001811690811561258f575060011461255557505050565b60009081526020812093945091925b838310612575575060209250010190565b600181602092949394548385870101520191019190612564565b905060209495507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091509291921683830152151560051b010190565b9061027c6125df9260405193848092612529565b038361022c565b90600160606040516125f78161019b565b610bcb819567ffffffffffffffff815473ffffffffffffffffffffffffffffffffffffffff8116855260ff8160a01c161515602086015260a81c1660408401526126476040518096819301612529565b038461022c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b908160051b918083046020149015171561269357565b61264e565b9190820391821161269357565b6126b182607f9261431f565b9116906801fffffffffffffffe67ffffffffffffffff83169260011b169180830460021490151715612693576003911c166004811015610af25790565b6126f6614365565b80518251810361293c5760005b8181106127165750509061027c916143c1565b61272081846124c2565b51602081019081515161273384886124c2565b51928351820361293c5790916000925b808410612757575050505050600101612703565b91949398612769848b989395986124c2565b5151986127778888516124c2565b5199806128d8575b5060a08a01988b60206127958b8d5151936124c2565b51015151036128805760005b8a515181101561286b576127dd6127d46127ca8f60206127c28f87936124c2565b5101516124c2565b5163ffffffff1690565b63ffffffff1690565b8b816127ee575b50506001016127a1565b6127d460406128018561280d94516124c2565b51015163ffffffff1690565b9081811061281c57508b6127e4565b8d51516040517f48e617b30000000000000000000000000000000000000000000000000000000081526004810191909152602481019390935260448301919091526064820152608490fd5b0390fd5b50985098509893949095600101929091612743565b61238f8b5161289c6060825192015167ffffffffffffffff1690565b7f70a193fd0000000000000000000000000000000000000000000000000000000060005260049190915267ffffffffffffffff16602452604490565b60808b015181101561277f5761238f908b6128fb885167ffffffffffffffff1690565b9051517f7531a8c60000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff909116600452602452604452606490565b7f83e3f5640000000000000000000000000000000000000000000000000000000060005260046000fd5b60405190612973826101d8565b60006020838281520152565b6040519061298e60208361022c565b600080835282815b8281106129a257505050565b6020906129ad612966565b82828501015201612996565b8051825267ffffffffffffffff60208201511660208301526080612a016129ef604084015160a0604087015260a086019061074d565b6060840151858203606087015261074d565b9101519160808183039101526020808351928381520192019060005b818110612a2a5750505090565b8251805173ffffffffffffffffffffffffffffffffffffffff1685526020908101518186015260409094019390920191600101612a1d565b9060206103a59281815201906129b9565b6040513d6000823e3d90fd5b3d15612aaa573d90612a9082610319565b91612a9e604051938461022c565b82523d6000602084013e565b606090565b9060206103a592818152019061074d565b81601f820112156102e2578051612ad681610319565b92612ae4604051948561022c565b818452602082840101116102e2576103a5916020808501910161072a565b90916060828403126102e2578151612b1981610304565b92602083015167ffffffffffffffff81116102e257604091612b3c918501612ac0565b92015190565b9293606092959461ffff612b7373ffffffffffffffffffffffffffffffffffffffff946080885260808801906129b9565b97166020860152604085015216910152565b9391929093303303612f0f57612b9a906148a6565b92612ba361297f565b9460a08501518051612eba575b5050505050805191612bcf6020845194015167ffffffffffffffff1690565b906020830151916040840192612bfd845192612be961027e565b97885267ffffffffffffffff166020880152565b60408601526060850152608084015273ffffffffffffffffffffffffffffffffffffffff612c4060055473ffffffffffffffffffffffffffffffffffffffff1690565b1680612e0b575b5051511580612dff575b8015612ddc575b8015612da6575b6123dc57612d259181612ca46120ee612c8a610a676020600097510167ffffffffffffffff90511690565b5473ffffffffffffffffffffffffffffffffffffffff1690565b9083612ccc6060608084015193015173ffffffffffffffffffffffffffffffffffffffff1690565b93604051968795869485937f3cf979830000000000000000000000000000000000000000000000000000000085527f00000000000000000000000000000000000000000000000000000000000000009060048601612b42565b03925af1908115612da157600090600092612d7a575b5015612d445750565b612867906040519182917f0a8d6e8c00000000000000000000000000000000000000000000000000000000835260048301612aaf565b9050612d9991503d806000833e612d91818361022c565b810190612b02565b509038612d3b565b612a73565b50612dd7612dd3612dce606084015173ffffffffffffffffffffffffffffffffffffffff1690565b614b68565b1590565b612c5f565b50606081015173ffffffffffffffffffffffffffffffffffffffff163b15612c58565b50608081015115612c51565b803b156102e257600060405180927f08d450a1000000000000000000000000000000000000000000000000000000008252818381612e4c8a60048301612a62565b03925af19081612e9f575b50612e9957612867612e67612a7f565b6040519182917f09c2532500000000000000000000000000000000000000000000000000000000835260048301612aaf565b38612c47565b80612eae6000612eb49361022c565b806104fc565b38612e57565b8596506020612f04960151612ee6606089015173ffffffffffffffffffffffffffffffffffffffff1690565b90612efe60208a510167ffffffffffffffff90511690565b92614a32565b903880808080612bb0565b7f371a73280000000000000000000000000000000000000000000000000000000060005260046000fd5b60405190612f468261019b565b60006060838281528260208201528260408201520152565b90604051612f6b8161019b565b606060ff600183958054855201548181166020850152818160081c16604085015260101c161515910152565b906020825491828152019160005260206000209060005b818110612fbb5750505090565b825473ffffffffffffffffffffffffffffffffffffffff16845260209093019260019283019201612fae565b9061027c6125df9260405193848092612f97565b35907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff821682036102e257565b81601f820112156102e25780359061303b826102ac565b92613049604051948561022c565b82845260208085019360061b830101918183116102e257602001925b828410613073575050505090565b6040848303126102e2576020604091825161308d816101d8565b8635613098816102e7565b81526130a5838801612ffb565b83820152815201930192613065565b81601f820112156102e2578035906130cb826102ac565b926130d9604051948561022c565b82845260208085019360051b830101918183116102e25760208101935b83851061310557505050505090565b843567ffffffffffffffff81116102e257820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082860301126102e25760405191613151836101bc565b602082013561315f816102e7565b835260408201359267ffffffffffffffff84116102e25760a08361318a88602080988198010161038a565b85840152606081013561319c816102e7565b604084015260808101356131af816102e7565b6060840152013560808201528152019401936130f6565b81601f820112156102e2578035906131dd826102ac565b926131eb604051948561022c565b82845260208085019360061b830101918183116102e257602001925b828410613215575050505090565b6040848303126102e2576020604091825161322f816101d8565b863581528287013583820152815201930192613207565b6020818303126102e25780359067ffffffffffffffff82116102e257016060818303126102e2576040519161327a836101f4565b813567ffffffffffffffff81116102e25782016040818303126102e257604051906132a4826101d8565b803567ffffffffffffffff81116102e257810183601f820112156102e25780356132cd816102ac565b916132db604051938461022c565b81835260208084019260061b820101908682116102e257602001915b81831061337657505050825260208101359067ffffffffffffffff82116102e25761332491849101613024565b60208201528352602082013567ffffffffffffffff81116102e2578161334b9184016130b4565b6020840152604082013567ffffffffffffffff81116102e25761336e92016131c6565b604082015290565b6040838803126102e25760206040918251613390816101d8565b853561339b816102c4565b81526133a8838701612ffb565b838201528152019201916132f7565b9080602083519182815201916020808360051b8301019401926000915b8383106133e357505050505090565b9091929394602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0856001950301865288519067ffffffffffffffff82511681526080806134408585015160a08786015260a085019061074d565b9367ffffffffffffffff604082015116604085015267ffffffffffffffff60608201511660608501520151910152970193019301919392906133d4565b9173ffffffffffffffffffffffffffffffffffffffff6134ab921683526060602084015260608301906133b7565b9060408183039101526020808351928381520192019060005b8181106134d15750505090565b82518051855260209081015181860152604090940193909201916001016134c4565b906020808351928381520192019060005b8181106135115750505090565b8251805167ffffffffffffffff1685526020908101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168186015260409094019390920191600101613504565b9190604081019083519160408252825180915260206060830193019060005b81811061359b5750505060206103a5939401519060208184039101526134f3565b8251805173ffffffffffffffffffffffffffffffffffffffff1686526020908101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16818701526040909501949092019160010161357a565b9060206103a592818152019061355b565b908160209103126102e257516103a581610304565b909161362e6103a59360408452604084019061074d565b916020818403910152612529565b67ffffffffffffffff60019116019067ffffffffffffffff821161269357565b90916136736103a5936040845260408401906133b7565b91602081840391015261355b565b92969395919097949761369682820182613246565b986136aa612dd360045460ff9060c01c1690565b613c29575b89518051511590811591613c1a575b50613ae8575b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16999860208a019860005b8a518051821015613a86578161371a916124c2565b518d61372e825167ffffffffffffffff1690565b6040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff00000000000000000000000000000000608083901b1660048201529091602090829060249082905afa908115612da157600091613a58575b50613a21576137a390614be8565b602082018051602081519101209060018301916137bf836125cb565b60208151910120036139eb5750508054604083015167ffffffffffffffff9081169160a81c1681148015906139c1575b613953575060808201519081156139295761384282613833613819865167ffffffffffffffff1690565b67ffffffffffffffff16600052600a602052604060002090565b90600052602052604060002090565b546138da5782916138bd6138d39261387061386b606060019998015167ffffffffffffffff1690565b61363c565b7fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff7cffffffffffffffff00000000000000000000000000000000000000000083549260a81b169116179055565b61383361381942945167ffffffffffffffff1690565b5501613705565b506138f061238f925167ffffffffffffffff1690565b7f32cf0cbf0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff16600452602452604490565b7f504570e30000000000000000000000000000000000000000000000000000000060005260046000fd5b8261238f9161397f606061396f845167ffffffffffffffff1690565b93015167ffffffffffffffff1690565b7fd5e0f0d60000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff92831660045290821660245216604452606490565b506139da610a20606085015167ffffffffffffffff1690565b67ffffffffffffffff8216116137ef565b516128676040519283927fb80d8fa900000000000000000000000000000000000000000000000000000000845260048401613617565b7ffdbd6a720000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff1660045260246000fd5b613a79915060203d8111613a7f575b613a71818361022c565b810190613602565b38613795565b503d613a67565b5050613ae29496989b507f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e461027c9b613ada949597999b51905190613ad06040519283928361365c565b0390a13691610f3b565b943691610f3b565b93615087565b613afe602086015b3567ffffffffffffffff1690565b600b5467ffffffffffffffff82811691161015613be557613b4e9067ffffffffffffffff167fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000600b541617600b55565b613b736120ee6120ee60045473ffffffffffffffffffffffffffffffffffffffff1690565b8a5190803b156102e257613bbb6000929183926040519485809481937f3937306f000000000000000000000000000000000000000000000000000000008352600483016135f1565b03925af18015612da157613bd0575b506136c4565b80612eae6000613bdf9361022c565b38613bca565b5060208a0151516136c4577f226111670000000000000000000000000000000000000000000000000000000060005260046000fd5b602091500151511515386136be565b60208a01518051613c3b575b506136af565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169060408c0151823b156102e257613cbc926000926040518095819482937f70a9089e000000000000000000000000000000000000000000000000000000008452306004850161347d565b03915afa8015612da15715613c355780612eae6000613cda9361022c565b38613c35565b60405190613cef60208361022c565b6000808352366020840137565b613d046140d4565b60005b81518110156123dc57613d1a81836124c2565b51906040820160ff613d2d825160ff1690565b16156140a357602083015160ff1692613d538460ff166000526002602052604060002090565b9160018301918254613d6e613d688260ff1690565b60ff1690565b61404f5750613db7613d836060830151151590565b84547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff1690151560101b62ff000016178455565b60a08101918251610100815111613f93578051156140205760038601613de5613ddf82612fe7565b8a616575565b6060840151613e93575b947fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f54794600294613e6f613e41613e8d9a96613e3a8760019f9c613e35613e859a8f61676a565b6152b4565b5160ff1690565b84547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff821617909455565b5190818555519060405195869501908886615348565b0390a16167f9565b01613d07565b97946002879395970196613eaf613ea989612fe7565b88616575565b608085015194610100865111613ff1578551613ed7613d68613ed28a5160ff1690565b6152a0565b1015613fc2578551845111613f9357613e6f613e417fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f54798613e3a8760019f613e35613e8d9f9a8f613f7b60029f613f75613e859f8f90613e358492613f3d845160ff1690565b90805490917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff90911660089190911b61ff0016179055565b8261665e565b505050979c9f50975050969a50505094509450613def565b7f367f56a200000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b7f367f56a200000000000000000000000000000000000000000000000000000000600052600360045260246000fd5b7f367f56a200000000000000000000000000000000000000000000000000000000600052600260045260246000fd5b7f367f56a200000000000000000000000000000000000000000000000000000000600052600560045260246000fd5b60101c60ff1661406a6140656060840151151590565b151590565b90151514613db7577f87f6037c0000000000000000000000000000000000000000000000000000000060005260ff861660045260246000fd5b7f367f56a20000000000000000000000000000000000000000000000000000000060005261238f6024906000600452565b73ffffffffffffffffffffffffffffffffffffffff6001541633036140f557565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b6040516020810190600082526020815261413a60408261022c565b51902090565b81811061414b575050565b60008155600101614140565b9190601f811161416657505050565b61027c926000526020600020906020601f840160051c83019310614192575b601f0160051c0190614140565b9091508190614185565b919091825167ffffffffffffffff81116101b7576141c4816141be84546124d6565b84614157565b6020601f8211600114614222578190614213939495600092614217575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055565b0151905038806141e1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082169061425584600052602060002090565b9160005b8181106142af57509583600195969710614278575b505050811b019055565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c1916905538808061426e565b9192602060018192868b015181550194019201614259565b90600160a06103a5936020815267ffffffffffffffff845473ffffffffffffffffffffffffffffffffffffffff8116602084015260ff81851c161515604084015260a81c166060820152608080820152019101612529565b9067ffffffffffffffff614361921660005260096020526701ffffffffffffff60406000209160071c1667ffffffffffffffff16600052602052604060002090565b5490565b7f00000000000000000000000000000000000000000000000000000000000000004681036143905750565b7f0f01ce85000000000000000000000000000000000000000000000000000000006000526004524660245260446000fd5b9190918051156144815782511592602091604051926143e0818561022c565b600084527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810160005b81811061445d5750505060005b8151811015614455578061443e614430600193856124c2565b518815614444578690615496565b01614417565b61444e83876124c2565b5190615496565b505050509050565b829060405161446b816101d8565b600081526060838201528282890101520161440a565b7fc2e5347d0000000000000000000000000000000000000000000000000000000060005260046000fd5b91908260a09103126102e2576040516144c3816101bc565b60808082948051845260208101516144da816102e7565b602085015260408101516144ed816102e7565b60408501526060810151614500816102e7565b6060850152015191610bcb836102e7565b519061027c826102c4565b519061027c82610bda565b81601f820112156102e25780519061453e826102ac565b9261454c604051948561022c565b82845260208085019360051b830101918183116102e25760208101935b83851061457857505050505090565b845167ffffffffffffffff81116102e257820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082860301126102e257604051916145c4836101bc565b602082015167ffffffffffffffff81116102e2578560206145e792850101612ac0565b835260408201516145f7816102c4565b60208401526146086060830161451c565b604084015260808201519267ffffffffffffffff84116102e25760a083614636886020809881980101612ac0565b606084015201516080820152815201940193614569565b6020818303126102e25780519067ffffffffffffffff82116102e25701610140818303126102e25761467d61026d565b9161468881836144ab565b835260a082015167ffffffffffffffff81116102e257816146aa918401612ac0565b602084015260c082015167ffffffffffffffff81116102e257816146cf918401612ac0565b60408401526146e060e08301614511565b6060840152610100820151608084015261012082015167ffffffffffffffff81116102e25761470f9201614527565b60a082015290565b9080602083519182815201916020808360051b8301019401926000915b83831061474357505050505090565b9091929394602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085600195030186528851906080806147d1614791855160a0865260a086019061074d565b73ffffffffffffffffffffffffffffffffffffffff87870151168786015263ffffffff60408701511660408601526060860151858203606087015261074d565b93015191015297019301930191939290614734565b6103a59167ffffffffffffffff6080835180518452826020820151166020850152826040820151166040850152826060820151166060850152015116608082015260a061485861484660208501516101408486015261014085019061074d565b604085015184820360c086015261074d565b606084015173ffffffffffffffffffffffffffffffffffffffff1660e0840152926080810151610100840152015190610120818403910152614717565b9060206103a59281815201906147e6565b600061494481926040516148b981610210565b6148c16123e0565b81526060602082015260606040820152836060820152836080820152606060a08201525061490e6120ee6120ee600b5473ffffffffffffffffffffffffffffffffffffffff9060401c1690565b906040519485809481937f4546c6e500000000000000000000000000000000000000000000000000000000835260048301614895565b03925af160009181614993575b506103a557612867614961612a7f565b6040519182917f828ebdfb00000000000000000000000000000000000000000000000000000000835260048301612aaf565b6149b19192503d806000833e6149a9818361022c565b81019061464d565b9038614951565b91908110156124bd5760051b0190565b356103a581610bda565b91908110156124bd5760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156102e257019081359167ffffffffffffffff83116102e25760200182360381136102e2579190565b909294919397968151967fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0614a7f614a698a6102ac565b99614a776040519b8c61022c565b808b526102ac565b0160005b818110614b5157505060005b8351811015614b445780614ad68c8a8a8a614ad0614ac9878d614ac2828f8f9d8f9e60019f81614af2575b5050506124c2565b51976149d2565b3691610353565b93615eb0565b614ae0828c6124c2565b52614aeb818b6124c2565b5001614a8f565b63ffffffff614b0a614b058585856149b8565b6149c8565b1615614aba57614b3a92614b2192614b05926149b8565b6040614b2d85856124c2565b51019063ffffffff169052565b8f8f908391614aba565b5096985050505050505050565b602090614b5c612966565b82828d01015201614a83565b614b927f85572ffb000000000000000000000000000000000000000000000000000000008261639e565b9081614bac575b81614ba2575090565b6103a5915061633e565b9050614bb781616278565b1590614b99565b614b927faff2afbf000000000000000000000000000000000000000000000000000000008261639e565b67ffffffffffffffff1680600052600860205260406000209060ff825460a01c1615614c12575090565b7fed053c590000000000000000000000000000000000000000000000000000000060005260045260246000fd5b608401908160841161269357565b60a001908160a01161269357565b9190820180921161269357565b60031115610af257565b6003821015610af25752565b9061027c604051614c8e816101d8565b602060ff829554818116845260081c169101614c72565b80548210156124bd5760005260206000200190600090565b60ff60019116019060ff821161269357565b60ff601b9116019060ff821161269357565b90606092604091835260208301370190565b6001600052600260205293614d277fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0612f5e565b93853594614d3485614c3f565b6060820190614d438251151590565b615059575b80360361502857508151878103614ff65750614d62614365565b60016000526003602052614dbe614db97fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c5b3373ffffffffffffffffffffffffffffffffffffffff16600052602052604060002090565b614c7e565b60026020820151614dce81614c68565b614dd781614c68565b149081614f81575b5015614f575751614e3e575b50505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef090614e21613af060019460200190565b6040805192835267ffffffffffffffff91909116602083015290a2565b614e5f613d68614e5a602085969799989a955194015160ff1690565b614cbd565b03614f2d578151835103614f0357614efb6000614e2194613af094614ea97f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09960019b3691610353565b60208151910120604051614ef281614ec689602083019586614ce1565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810183528261022c565b5190208a616400565b948394614deb565b7fa75d88af0000000000000000000000000000000000000000000000000000000060005260046000fd5b7f71253a250000000000000000000000000000000000000000000000000000000060005260046000fd5b7fda0f08e80000000000000000000000000000000000000000000000000000000060005260046000fd5b60016000526002602052614fee91506120ee90614fce90614fc860037fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e05b01915160ff1690565b90614ca5565b905473ffffffffffffffffffffffffffffffffffffffff9160031b1c1690565b331438614ddf565b7f93df584c00000000000000000000000000000000000000000000000000000000600052600452602487905260446000fd5b7f8e1192e1000000000000000000000000000000000000000000000000000000006000526004523660245260446000fd5b6150829061507c61507261506d875161267d565b614c4d565b61507c885161267d565b90614c5b565b614d48565b600080526002602052949093909290916150c07fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b612f5e565b948635956150cd83614c3f565b60608201906150dc8251151590565b61527d575b8036036150285750815188810361524b57506150fb614365565b600080526003602052615130614db97f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff614d94565b6002602082015161514081614c68565b61514981614c68565b149081615202575b5015614f575751615194575b5050505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef090614e21613af060009460200190565b6151b0613d68614e5a602087989a999b96975194015160ff1690565b03614f2d578351865103614f03576000967f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef096614e2195614ea96151f994613af0973691610353565b9483943861515d565b60008052600260205261524391506120ee90614fce90614fc860037fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b614fbf565b331438615151565b7f93df584c00000000000000000000000000000000000000000000000000000000600052600452602488905260446000fd5b61529b9061507c61529161506d895161267d565b61507c8a5161267d565b6150e1565b60ff166003029060ff821691820361269357565b81519167ffffffffffffffff83116101b7576801000000000000000083116101b757602090825484845580851061532b575b500190600052602060002060005b8381106153015750505050565b600190602073ffffffffffffffffffffffffffffffffffffffff85511694019381840155016152f4565b615342908460005285846000209182019101614140565b386152e6565b95949392909160ff61536d93168752602087015260a0604087015260a0860190612f97565b84810360608601526020808351928381520192019060005b8181106153a05750505090608061027c9294019060ff169052565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101615385565b6006548110156124bd5760066000527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f015490565b67ffffffffffffffff6103a5949381606094168352166020820152816040820152019061074d565b6040906103a593928152816020820152019061074d565b92919067ffffffffffffffff908160649516600452166024526004811015610af257604452565b9493926154806060936154919388526020880190610af7565b60806040870152608086019061074d565b930152565b906154a9825167ffffffffffffffff1690565b81516040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff00000000000000000000000000000000608084901b16600482015290151593919067ffffffffffffffff8216906020816024817f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff165afa908115612da157600091615d17575b50615cd4576020830191825151948515615c8957604085018051518703615c5f5761557f87612432565b957f00000000000000000000000000000000000000000000000000000000000000006155b560016155af87614be8565b016125cb565b6020815191012060405161561681614ec66020820194868b8767ffffffffffffffff6060929594938160808401977f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f85521660208401521660408201520152565b5190209067ffffffffffffffff1660005b8a8110615b91575050508060806060615647930151910151908886616a6c565b978815615b595760005b8881106156645750505050505050505050565b5a6156708289516124c2565b5180516060015161568b9067ffffffffffffffff16886126a5565b61569481610ae8565b8015908d8283159384615b46575b15615b025760608815615a8357506156c960206156bf898d6124c2565b5101519242612698565b6004546156de9060a01c63ffffffff166127d4565b108015615a70575b15615a38576156f5878b6124c2565b5151615a22575b8451608001516157159067ffffffffffffffff16610a20565b615943575b506157268689516124c2565b5160a0850151518151036158ec579361578d9695938c938f9661576c8e958c92615766615760606089510167ffffffffffffffff90511690565b89616a9f565b86616bb9565b9a908096615787606088510167ffffffffffffffff90511690565b90616b2a565b615882575b505061579d82610ae8565b60028203615820575b6001966158167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b9367ffffffffffffffff9351926158076157fe8b6157f6606088015167ffffffffffffffff1690565b96519b6124c2565b51985a90612698565b91604051958695169885615467565b0390a45b01615651565b9150919394925061583082610ae8565b60038203615844578b929493918a916157a6565b51606001517f926c5a3e0000000000000000000000000000000000000000000000000000000060005261238f919067ffffffffffffffff1689615440565b61588b84610ae8565b600384036157925790929495506158a3919350610ae8565b6158b3578b92918a913880615792565b849051516128676040519283927f2b11b8d900000000000000000000000000000000000000000000000000000000845260048401615429565b61238f8b615907606088510167ffffffffffffffff90511690565b7f1cfe6d8b0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff91821660045216602452604490565b61594c83610ae8565b615957575b3861571a565b83516080015167ffffffffffffffff16602080860151918c6159a660405194859384937fe0e03cae00000000000000000000000000000000000000000000000000000000855260048501615401565b0381600073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115612da157600091615a04575b5061595157505050505060019061581a565b615a1c915060203d8111613a7f57613a71818361022c565b386159f2565b615a2c878b6124c2565b515160808601526156fc565b7fa9cfc8620000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff8b1660045260246000fd5b50615a7a83610ae8565b600383146156e6565b915083615a8f84610ae8565b156156fc57506001959450615afa9250615ad791507f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe6512093510167ffffffffffffffff90511690565b6040805167ffffffffffffffff808c168252909216602083015290918291820190565b0390a161581a565b505050506001929150615afa615ad760607f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c93510167ffffffffffffffff90511690565b50615b5083610ae8565b600383146156a2565b7f7dd17a7e0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff841660045260246000fd5b615b9c818a516124c2565b5180516040015167ffffffffffffffff16838103615c28575080516020015167ffffffffffffffff16898103615beb575090615bda84600193616947565b615be4828d6124c2565b5201615627565b7f6c95f1eb0000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff808a166004521660245260446000fd5b7f38432a220000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff1660045260246000fd5b7f57e0e0830000000000000000000000000000000000000000000000000000000060005260046000fd5b61238f615c9e865167ffffffffffffffff1690565b7fced9e4960000000000000000000000000000000000000000000000000000000060005267ffffffffffffffff16600452602490565b5092915050613a215760405167ffffffffffffffff9190911681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d89493390602090a1565b615d30915060203d602011613a7f57613a71818361022c565b38615555565b908160209103126102e257516103a5816102c4565b906103a5916020815260e0615e40615e0d615d748551610100602087015261012086019061074d565b602086015167ffffffffffffffff166040860152604086015173ffffffffffffffffffffffffffffffffffffffff16606086015260608601516080860152615dd9608087015160a087019073ffffffffffffffffffffffffffffffffffffffff169052565b60a08601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08683030160c087015261074d565b60c08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0858303018486015261074d565b920151906101007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08285030191015261074d565b60409073ffffffffffffffffffffffffffffffffffffffff6103a59493168152816020820152019061074d565b908160209103126102e2575190565b91939293615ebc612966565b50602083015173ffffffffffffffffffffffffffffffffffffffff166040517fbbe4f6db00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152909590926020848060248101038173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa938415612da157600094616247575b5073ffffffffffffffffffffffffffffffffffffffff84169586158015616235575b6161f15761606f61609892614ec692615fbf615fb86127d460408c015163ffffffff1690565b8c89616d4b565b9690996080810151615fee6060835193015193615fda61028d565b96875267ffffffffffffffff166020870152565b73ffffffffffffffffffffffffffffffffffffffff8a166040860152606085015273ffffffffffffffffffffffffffffffffffffffff8d16608085015260a084015260c083015260e08201526040519283917f3907753700000000000000000000000000000000000000000000000000000000602084015260248301615d4b565b82857f000000000000000000000000000000000000000000000000000000000000000092616dff565b949091156161bc575080516020810361618a5750906160c1826020808a95518301019101615ea1565b9573ffffffffffffffffffffffffffffffffffffffff841603616113575b505050505061610b6160ef61029d565b73ffffffffffffffffffffffffffffffffffffffff9093168352565b602082015290565b6161269361612091612698565b91616d4b565b50908082108015616177575b61613e578084816160df565b61238f927fa966e21f00000000000000000000000000000000000000000000000000000000600052929190606493600452602452604452565b50826161838284612698565b1415616132565b7f78ef802400000000000000000000000000000000000000000000000000000000600052602060045260245260446000fd5b6128676040519283927f9fe2f95a00000000000000000000000000000000000000000000000000000000845260048401615e74565b7fae9b4ce90000000000000000000000000000000000000000000000000000000060005273ffffffffffffffffffffffffffffffffffffffff851660045260246000fd5b50616242612dd386614bbe565b615f92565b61626a91945060203d602011616271575b616262818361022c565b810190615d36565b9238615f70565b503d616258565b60405160208101917f01ffc9a70000000000000000000000000000000000000000000000000000000083527fffffffff000000000000000000000000000000000000000000000000000000006024830152602482526162d860448361022c565b6179185a10616314576020926000925191617530fa6000513d82616308575b5081616301575090565b9050151590565b602011159150386162f7565b7fea7f4b120000000000000000000000000000000000000000000000000000000060005260046000fd5b60405160208101917f01ffc9a70000000000000000000000000000000000000000000000000000000083527f01ffc9a7000000000000000000000000000000000000000000000000000000006024830152602482526162d860448361022c565b604051907fffffffff0000000000000000000000000000000000000000000000000000000060208301937f01ffc9a7000000000000000000000000000000000000000000000000000000008552166024830152602482526162d860448361022c565b919390926000948051946000965b86881061641f575050505050505050565b60208810156124bd5760206000616437878b1a614ccf565b6164418b876124c2565b51906164786164508d8a6124c2565b5160405193849389859094939260ff6060936080840197845216602083015260408201520152565b838052039060015afa15612da1576164cb614db96000516164a68960ff166000526003602052604060002090565b9073ffffffffffffffffffffffffffffffffffffffff16600052602052604060002090565b90600160208301516164dc81614c68565b6164e581614c68565b0361654b576165026164f8835160ff1690565b60ff600191161b90565b8116616521576165186164f86001935160ff1690565b1797019661640e565b7ff67bc7c40000000000000000000000000000000000000000000000000000000060005260046000fd5b7fca31867a0000000000000000000000000000000000000000000000000000000060005260046000fd5b91909160005b83518110156165e85760019060ff8316600052600360205260006165e16040822073ffffffffffffffffffffffffffffffffffffffff6165bb858a6124c2565b511673ffffffffffffffffffffffffffffffffffffffff16600052602052604060002090565b550161657b565b50509050565b815181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9190911617815590602001516003811015610af25781547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1660089190911b61ff0016179055565b919060005b81518110156165e85761667961226a82846124c2565b906166a2616698836164a68860ff166000526003602052604060002090565b5460081c60ff1690565b6166ab81614c68565b61673c5773ffffffffffffffffffffffffffffffffffffffff8216156167125761670c6001926167076166dc61029d565b60ff85168152916166f08660208501614c72565b6164a68960ff166000526003602052604060002090565b6165ee565b01616663565b7fd6c62c9b0000000000000000000000000000000000000000000000000000000060005260046000fd5b7f367f56a2000000000000000000000000000000000000000000000000000000006000526004805260246000fd5b919060005b81518110156165e85761678561226a82846124c2565b906167a4616698836164a68860ff166000526003602052604060002090565b6167ad81614c68565b61673c5773ffffffffffffffffffffffffffffffffffffffff821615616712576167f36001926167076167de61029d565b60ff85168152916166f0600260208501614c72565b0161676f565b60ff1680600052600260205260ff60016040600020015460101c1690801560001461687857501561684e577fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000600b5416600b55565b7f2f7b1ba20000000000000000000000000000000000000000000000000000000060005260046000fd5b6001146168825750565b61688857565b7ff718e9a00000000000000000000000000000000000000000000000000000000060005260046000fd5b8060005260076020526040600020541560001461693057600654680100000000000000008110156101b7576001810160065560006006548210156124bd57600690527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01819055600654906000526007602052604060002055600190565b50600090565b9060206103a5928181520190614717565b61413a81518051906169f8616973606086015173ffffffffffffffffffffffffffffffffffffffff1690565b614ec661698b606085015167ffffffffffffffff1690565b936169a56080808a015192015167ffffffffffffffff1690565b906040519586946020860198899367ffffffffffffffff60809473ffffffffffffffffffffffffffffffffffffffff82959998949960a089019a8952166020880152166040860152606085015216910152565b519020614ec66020840151602081519101209360a0604082015160208151910120910151604051616a3181614ec6602082019485616936565b51902090604051958694602086019889919260a093969594919660c08401976000855260208501526040840152606083015260808201520152565b9267ffffffffffffffff92616a8092616f25565b9116600052600a60205260406000209060005260205260406000205490565b607f8216906801fffffffffffffffe67ffffffffffffffff83169260011b16918083046002149015171561269357616b279167ffffffffffffffff616ae4858461431f565b921660005260096020526701ffffffffffffff60406000209460071c169160036001831b921b191617929067ffffffffffffffff16600052602052604060002090565b55565b9091607f83166801fffffffffffffffe67ffffffffffffffff82169160011b16908082046002149015171561269357616b63848461431f565b6004831015610af25767ffffffffffffffff616b279416600052600960205260036701ffffffffffffff60406000209660071c1693831b921b191617929067ffffffffffffffff16600052602052604060002090565b9190303b156102e25790616c0292916040519384937f60987c200000000000000000000000000000000000000000000000000000000085526060600486015260648501906147e6565b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc848203016024850152815180825260208201916020808360051b8301019401926000915b838310616cfa5750505050507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8382030160448401526020808351928381520192019060005b818110616cdb575050509080600092038183305af19081616cc6575b50616cbb57616cb6612a7f565b600391565b6002906103a5610715565b80612eae6000616cd59361022c565b38616ca9565b825163ffffffff16845285945060209384019390920191600101616c8d565b91939596509193602080616d38837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08660019603018752895161074d565b9701930193019092879695949293616c47565b6040517f70a0823100000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff9091166024820152919291616dce90616da58160448101614ec6565b84837f000000000000000000000000000000000000000000000000000000000000000092616dff565b929091156161bc575080516020810361618a575090616df9826020806103a595518301019101615ea1565b93612698565b939193616e0c6084610319565b94616e1a604051968761022c565b60848652616e286084610319565b947fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0602088019601368737833b15616efb575a90808210616ed1578291038060061c90031115616ea7576000918291825a9560208451940192f1905a9003923d9060848211616e9e575b6000908287523e929190565b60849150616e92565b7f37c3be290000000000000000000000000000000000000000000000000000000060005260046000fd5b7fafa32a2c0000000000000000000000000000000000000000000000000000000060005260046000fd5b7f0c3b563c0000000000000000000000000000000000000000000000000000000060005260046000fd5b80519282519084156170d657610101851115806170ca575b15616fe057818501947fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8601956101008711616fe05786156170ba57616f8287612432565b9660009586978795885b84811061701e5750505050507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe018095149384617014575b50508261700a575b505015616fe057616fdc916124c2565b5190565b7f09bde3390000000000000000000000000000000000000000000000000000000060005260046000fd5b1490503880616fcc565b1492503880616fc4565b6001811b828116036170ac57868a10156170975761704060018b019a856124c2565b51905b8c888c1015617083575061705b60018c019b866124c2565b515b818d11616fe05761707c828f9261707690600196617100565b926124c2565b5201616f8c565b60018d019c617091916124c2565b5161705d565b6170a560018c019b8d6124c2565b5190617043565b6170a56001890198846124c2565b505050509050616fdc91506124b0565b50610101821115616f3d565b7f11a6b2640000000000000000000000000000000000000000000000000000000060005260046000fd5b8181101561711257906103a591617117565b6103a5915b9060405190602082019260018452604083015260608201526060815261413a60808261022c56fea164736f6c634300081a000a49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b",
}

var MessageTransformerOffRampABI = MessageTransformerOffRampMetaData.ABI

var MessageTransformerOffRampBin = MessageTransformerOffRampMetaData.Bin

func DeployMessageTransformerOffRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OffRampStaticConfig, dynamicConfig OffRampDynamicConfig, sourceChainConfigs []OffRampSourceChainConfigArgs, messageTransformerAddr common.Address) (common.Address, *types.Transaction, *MessageTransformerOffRamp, error) {
	parsed, err := MessageTransformerOffRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageTransformerOffRampBin), backend, staticConfig, dynamicConfig, sourceChainConfigs, messageTransformerAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MessageTransformerOffRamp{address: address, abi: *parsed, MessageTransformerOffRampCaller: MessageTransformerOffRampCaller{contract: contract}, MessageTransformerOffRampTransactor: MessageTransformerOffRampTransactor{contract: contract}, MessageTransformerOffRampFilterer: MessageTransformerOffRampFilterer{contract: contract}}, nil
}

type MessageTransformerOffRamp struct {
	address common.Address
	abi     abi.ABI
	MessageTransformerOffRampCaller
	MessageTransformerOffRampTransactor
	MessageTransformerOffRampFilterer
}

type MessageTransformerOffRampCaller struct {
	contract *bind.BoundContract
}

type MessageTransformerOffRampTransactor struct {
	contract *bind.BoundContract
}

type MessageTransformerOffRampFilterer struct {
	contract *bind.BoundContract
}

type MessageTransformerOffRampSession struct {
	Contract     *MessageTransformerOffRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MessageTransformerOffRampCallerSession struct {
	Contract *MessageTransformerOffRampCaller
	CallOpts bind.CallOpts
}

type MessageTransformerOffRampTransactorSession struct {
	Contract     *MessageTransformerOffRampTransactor
	TransactOpts bind.TransactOpts
}

type MessageTransformerOffRampRaw struct {
	Contract *MessageTransformerOffRamp
}

type MessageTransformerOffRampCallerRaw struct {
	Contract *MessageTransformerOffRampCaller
}

type MessageTransformerOffRampTransactorRaw struct {
	Contract *MessageTransformerOffRampTransactor
}

func NewMessageTransformerOffRamp(address common.Address, backend bind.ContractBackend) (*MessageTransformerOffRamp, error) {
	abi, err := abi.JSON(strings.NewReader(MessageTransformerOffRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMessageTransformerOffRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRamp{address: address, abi: abi, MessageTransformerOffRampCaller: MessageTransformerOffRampCaller{contract: contract}, MessageTransformerOffRampTransactor: MessageTransformerOffRampTransactor{contract: contract}, MessageTransformerOffRampFilterer: MessageTransformerOffRampFilterer{contract: contract}}, nil
}

func NewMessageTransformerOffRampCaller(address common.Address, caller bind.ContractCaller) (*MessageTransformerOffRampCaller, error) {
	contract, err := bindMessageTransformerOffRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampCaller{contract: contract}, nil
}

func NewMessageTransformerOffRampTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageTransformerOffRampTransactor, error) {
	contract, err := bindMessageTransformerOffRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampTransactor{contract: contract}, nil
}

func NewMessageTransformerOffRampFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageTransformerOffRampFilterer, error) {
	contract, err := bindMessageTransformerOffRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampFilterer{contract: contract}, nil
}

func bindMessageTransformerOffRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MessageTransformerOffRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageTransformerOffRamp.Contract.MessageTransformerOffRampCaller.contract.Call(opts, result, method, params...)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.MessageTransformerOffRampTransactor.contract.Transfer(opts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.MessageTransformerOffRampTransactor.contract.Transact(opts, method, params...)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageTransformerOffRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.contract.Transfer(opts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.contract.Transact(opts, method, params...)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "ccipReceive", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _MessageTransformerOffRamp.Contract.CcipReceive(&_MessageTransformerOffRamp.CallOpts, arg0)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _MessageTransformerOffRamp.Contract.CcipReceive(&_MessageTransformerOffRamp.CallOpts, arg0)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getAllSourceChainConfigs")

	if err != nil {
		return *new([]uint64), *new([]OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)
	out1 := *abi.ConvertType(out[1], new([]OffRampSourceChainConfig)).(*[]OffRampSourceChainConfig)

	return out0, out1, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetAllSourceChainConfigs(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetAllSourceChainConfigs(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OffRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampDynamicConfig)).(*OffRampDynamicConfig)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetDynamicConfig(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetDynamicConfig(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getExecutionState", sourceChainSelector, sequenceNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _MessageTransformerOffRamp.Contract.GetExecutionState(&_MessageTransformerOffRamp.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _MessageTransformerOffRamp.Contract.GetExecutionState(&_MessageTransformerOffRamp.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getLatestPriceSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _MessageTransformerOffRamp.Contract.GetLatestPriceSequenceNumber(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _MessageTransformerOffRamp.Contract.GetLatestPriceSequenceNumber(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getMerkleRoot", sourceChainSelector, root)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _MessageTransformerOffRamp.Contract.GetMerkleRoot(&_MessageTransformerOffRamp.CallOpts, sourceChainSelector, root)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _MessageTransformerOffRamp.Contract.GetMerkleRoot(&_MessageTransformerOffRamp.CallOpts, sourceChainSelector, root)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetMessageTransformerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getMessageTransformerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetMessageTransformerAddress() (common.Address, error) {
	return _MessageTransformerOffRamp.Contract.GetMessageTransformerAddress(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetMessageTransformerAddress() (common.Address, error) {
	return _MessageTransformerOffRamp.Contract.GetMessageTransformerAddress(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampSourceChainConfig)).(*OffRampSourceChainConfig)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetSourceChainConfig(&_MessageTransformerOffRamp.CallOpts, sourceChainSelector)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetSourceChainConfig(&_MessageTransformerOffRamp.CallOpts, sourceChainSelector)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OffRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampStaticConfig)).(*OffRampStaticConfig)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetStaticConfig(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _MessageTransformerOffRamp.Contract.GetStaticConfig(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "latestConfigDetails", ocrPluginType)

	if err != nil {
		return *new(MultiOCR3BaseOCRConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOCRConfig)).(*MultiOCR3BaseOCRConfig)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _MessageTransformerOffRamp.Contract.LatestConfigDetails(&_MessageTransformerOffRamp.CallOpts, ocrPluginType)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _MessageTransformerOffRamp.Contract.LatestConfigDetails(&_MessageTransformerOffRamp.CallOpts, ocrPluginType)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) Owner() (common.Address, error) {
	return _MessageTransformerOffRamp.Contract.Owner(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) Owner() (common.Address, error) {
	return _MessageTransformerOffRamp.Contract.Owner(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MessageTransformerOffRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) TypeAndVersion() (string, error) {
	return _MessageTransformerOffRamp.Contract.TypeAndVersion(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampCallerSession) TypeAndVersion() (string, error) {
	return _MessageTransformerOffRamp.Contract.TypeAndVersion(&_MessageTransformerOffRamp.CallOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "acceptOwnership")
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.AcceptOwnership(&_MessageTransformerOffRamp.TransactOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.AcceptOwnership(&_MessageTransformerOffRamp.TransactOpts)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "applySourceChainConfigUpdates", sourceChainConfigUpdates)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.ApplySourceChainConfigUpdates(&_MessageTransformerOffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.ApplySourceChainConfigUpdates(&_MessageTransformerOffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) Commit(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "commit", reportContext, report, rs, ss, rawVs)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.Commit(&_MessageTransformerOffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.Commit(&_MessageTransformerOffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) Execute(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "execute", reportContext, report)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.Execute(&_MessageTransformerOffRamp.TransactOpts, reportContext, report)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.Execute(&_MessageTransformerOffRamp.TransactOpts, reportContext, report)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData, tokenGasOverrides)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.ExecuteSingleMessage(&_MessageTransformerOffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.ExecuteSingleMessage(&_MessageTransformerOffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "manuallyExecute", reports, gasLimitOverrides)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.ManuallyExecute(&_MessageTransformerOffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.ManuallyExecute(&_MessageTransformerOffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.SetDynamicConfig(&_MessageTransformerOffRamp.TransactOpts, dynamicConfig)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.SetDynamicConfig(&_MessageTransformerOffRamp.TransactOpts, dynamicConfig)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "setOCR3Configs", ocrConfigArgs)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.SetOCR3Configs(&_MessageTransformerOffRamp.TransactOpts, ocrConfigArgs)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.SetOCR3Configs(&_MessageTransformerOffRamp.TransactOpts, ocrConfigArgs)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.TransferOwnership(&_MessageTransformerOffRamp.TransactOpts, to)
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MessageTransformerOffRamp.Contract.TransferOwnership(&_MessageTransformerOffRamp.TransactOpts, to)
}

type MessageTransformerOffRampAlreadyAttemptedIterator struct {
	Event *MessageTransformerOffRampAlreadyAttempted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampAlreadyAttemptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampAlreadyAttempted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampAlreadyAttempted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampAlreadyAttemptedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampAlreadyAttemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampAlreadyAttempted struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterAlreadyAttempted(opts *bind.FilterOpts) (*MessageTransformerOffRampAlreadyAttemptedIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampAlreadyAttemptedIterator{contract: _MessageTransformerOffRamp.contract, event: "AlreadyAttempted", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampAlreadyAttempted) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampAlreadyAttempted)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseAlreadyAttempted(log types.Log) (*MessageTransformerOffRampAlreadyAttempted, error) {
	event := new(MessageTransformerOffRampAlreadyAttempted)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampCommitReportAcceptedIterator struct {
	Event *MessageTransformerOffRampCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampCommitReportAccepted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampCommitReportAccepted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampCommitReportAccepted struct {
	MerkleRoots  []InternalMerkleRoot
	PriceUpdates InternalPriceUpdates
	Raw          types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*MessageTransformerOffRampCommitReportAcceptedIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampCommitReportAcceptedIterator{contract: _MessageTransformerOffRamp.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampCommitReportAccepted)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseCommitReportAccepted(log types.Log) (*MessageTransformerOffRampCommitReportAccepted, error) {
	event := new(MessageTransformerOffRampCommitReportAccepted)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampConfigSetIterator struct {
	Event *MessageTransformerOffRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampConfigSet struct {
	OcrPluginType uint8
	ConfigDigest  [32]byte
	Signers       []common.Address
	Transmitters  []common.Address
	F             uint8
	Raw           types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*MessageTransformerOffRampConfigSetIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampConfigSetIterator{contract: _MessageTransformerOffRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampConfigSet)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseConfigSet(log types.Log) (*MessageTransformerOffRampConfigSet, error) {
	event := new(MessageTransformerOffRampConfigSet)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampDynamicConfigSetIterator struct {
	Event *MessageTransformerOffRampDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampDynamicConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampDynamicConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampDynamicConfigSet struct {
	DynamicConfig OffRampDynamicConfig
	Raw           types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts) (*MessageTransformerOffRampDynamicConfigSetIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampDynamicConfigSetIterator{contract: _MessageTransformerOffRamp.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampDynamicConfigSet) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampDynamicConfigSet)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseDynamicConfigSet(log types.Log) (*MessageTransformerOffRampDynamicConfigSet, error) {
	event := new(MessageTransformerOffRampDynamicConfigSet)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampExecutionStateChangedIterator struct {
	Event *MessageTransformerOffRampExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampExecutionStateChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampExecutionStateChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	MessageHash         [32]byte
	State               uint8
	ReturnData          []byte
	GasUsed             *big.Int
	Raw                 types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*MessageTransformerOffRampExecutionStateChangedIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampExecutionStateChangedIterator{contract: _MessageTransformerOffRamp.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampExecutionStateChanged)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseExecutionStateChanged(log types.Log) (*MessageTransformerOffRampExecutionStateChanged, error) {
	event := new(MessageTransformerOffRampExecutionStateChanged)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampOwnershipTransferRequestedIterator struct {
	Event *MessageTransformerOffRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOffRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampOwnershipTransferRequestedIterator{contract: _MessageTransformerOffRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampOwnershipTransferRequested)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*MessageTransformerOffRampOwnershipTransferRequested, error) {
	event := new(MessageTransformerOffRampOwnershipTransferRequested)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampOwnershipTransferredIterator struct {
	Event *MessageTransformerOffRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOffRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampOwnershipTransferredIterator{contract: _MessageTransformerOffRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampOwnershipTransferred)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseOwnershipTransferred(log types.Log) (*MessageTransformerOffRampOwnershipTransferred, error) {
	event := new(MessageTransformerOffRampOwnershipTransferred)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampRootRemovedIterator struct {
	Event *MessageTransformerOffRampRootRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampRootRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampRootRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampRootRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampRootRemovedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampRootRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampRootRemoved struct {
	Root [32]byte
	Raw  types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterRootRemoved(opts *bind.FilterOpts) (*MessageTransformerOffRampRootRemovedIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampRootRemovedIterator{contract: _MessageTransformerOffRamp.contract, event: "RootRemoved", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampRootRemoved) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampRootRemoved)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "RootRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseRootRemoved(log types.Log) (*MessageTransformerOffRampRootRemoved, error) {
	event := new(MessageTransformerOffRampRootRemoved)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "RootRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampSkippedAlreadyExecutedMessageIterator struct {
	Event *MessageTransformerOffRampSkippedAlreadyExecutedMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampSkippedAlreadyExecutedMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampSkippedAlreadyExecutedMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampSkippedAlreadyExecutedMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampSkippedAlreadyExecutedMessageIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampSkippedAlreadyExecutedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampSkippedAlreadyExecutedMessage struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*MessageTransformerOffRampSkippedAlreadyExecutedMessageIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampSkippedAlreadyExecutedMessageIterator{contract: _MessageTransformerOffRamp.contract, event: "SkippedAlreadyExecutedMessage", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSkippedAlreadyExecutedMessage) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampSkippedAlreadyExecutedMessage)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseSkippedAlreadyExecutedMessage(log types.Log) (*MessageTransformerOffRampSkippedAlreadyExecutedMessage, error) {
	event := new(MessageTransformerOffRampSkippedAlreadyExecutedMessage)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampSkippedReportExecutionIterator struct {
	Event *MessageTransformerOffRampSkippedReportExecution

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampSkippedReportExecutionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampSkippedReportExecution)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampSkippedReportExecution)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampSkippedReportExecutionIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampSkippedReportExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampSkippedReportExecution struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterSkippedReportExecution(opts *bind.FilterOpts) (*MessageTransformerOffRampSkippedReportExecutionIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampSkippedReportExecutionIterator{contract: _MessageTransformerOffRamp.contract, event: "SkippedReportExecution", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSkippedReportExecution) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampSkippedReportExecution)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseSkippedReportExecution(log types.Log) (*MessageTransformerOffRampSkippedReportExecution, error) {
	event := new(MessageTransformerOffRampSkippedReportExecution)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampSourceChainConfigSetIterator struct {
	Event *MessageTransformerOffRampSourceChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampSourceChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampSourceChainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampSourceChainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampSourceChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampSourceChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceConfig        OffRampSourceChainConfig
	Raw                 types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*MessageTransformerOffRampSourceChainConfigSetIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampSourceChainConfigSetIterator{contract: _MessageTransformerOffRamp.contract, event: "SourceChainConfigSet", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampSourceChainConfigSet)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseSourceChainConfigSet(log types.Log) (*MessageTransformerOffRampSourceChainConfigSet, error) {
	event := new(MessageTransformerOffRampSourceChainConfigSet)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampSourceChainSelectorAddedIterator struct {
	Event *MessageTransformerOffRampSourceChainSelectorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampSourceChainSelectorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampSourceChainSelectorAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampSourceChainSelectorAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampSourceChainSelectorAddedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampSourceChainSelectorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampSourceChainSelectorAdded struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*MessageTransformerOffRampSourceChainSelectorAddedIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampSourceChainSelectorAddedIterator{contract: _MessageTransformerOffRamp.contract, event: "SourceChainSelectorAdded", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSourceChainSelectorAdded) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampSourceChainSelectorAdded)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseSourceChainSelectorAdded(log types.Log) (*MessageTransformerOffRampSourceChainSelectorAdded, error) {
	event := new(MessageTransformerOffRampSourceChainSelectorAdded)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampStaticConfigSetIterator struct {
	Event *MessageTransformerOffRampStaticConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampStaticConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampStaticConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampStaticConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampStaticConfigSetIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampStaticConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampStaticConfigSet struct {
	StaticConfig OffRampStaticConfig
	Raw          types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterStaticConfigSet(opts *bind.FilterOpts) (*MessageTransformerOffRampStaticConfigSetIterator, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampStaticConfigSetIterator{contract: _MessageTransformerOffRamp.contract, event: "StaticConfigSet", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampStaticConfigSet) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampStaticConfigSet)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseStaticConfigSet(log types.Log) (*MessageTransformerOffRampStaticConfigSet, error) {
	event := new(MessageTransformerOffRampStaticConfigSet)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOffRampTransmittedIterator struct {
	Event *MessageTransformerOffRampTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOffRampTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOffRampTransmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(MessageTransformerOffRampTransmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *MessageTransformerOffRampTransmittedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOffRampTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOffRampTransmitted struct {
	OcrPluginType  uint8
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*MessageTransformerOffRampTransmittedIterator, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.FilterLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOffRampTransmittedIterator{contract: _MessageTransformerOffRamp.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampTransmitted, ocrPluginType []uint8) (event.Subscription, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _MessageTransformerOffRamp.contract.WatchLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOffRampTransmitted)
				if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRampFilterer) ParseTransmitted(log types.Log) (*MessageTransformerOffRampTransmitted, error) {
	event := new(MessageTransformerOffRampTransmitted)
	if err := _MessageTransformerOffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MessageTransformerOffRamp *MessageTransformerOffRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MessageTransformerOffRamp.abi.Events["AlreadyAttempted"].ID:
		return _MessageTransformerOffRamp.ParseAlreadyAttempted(log)
	case _MessageTransformerOffRamp.abi.Events["CommitReportAccepted"].ID:
		return _MessageTransformerOffRamp.ParseCommitReportAccepted(log)
	case _MessageTransformerOffRamp.abi.Events["ConfigSet"].ID:
		return _MessageTransformerOffRamp.ParseConfigSet(log)
	case _MessageTransformerOffRamp.abi.Events["DynamicConfigSet"].ID:
		return _MessageTransformerOffRamp.ParseDynamicConfigSet(log)
	case _MessageTransformerOffRamp.abi.Events["ExecutionStateChanged"].ID:
		return _MessageTransformerOffRamp.ParseExecutionStateChanged(log)
	case _MessageTransformerOffRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _MessageTransformerOffRamp.ParseOwnershipTransferRequested(log)
	case _MessageTransformerOffRamp.abi.Events["OwnershipTransferred"].ID:
		return _MessageTransformerOffRamp.ParseOwnershipTransferred(log)
	case _MessageTransformerOffRamp.abi.Events["RootRemoved"].ID:
		return _MessageTransformerOffRamp.ParseRootRemoved(log)
	case _MessageTransformerOffRamp.abi.Events["SkippedAlreadyExecutedMessage"].ID:
		return _MessageTransformerOffRamp.ParseSkippedAlreadyExecutedMessage(log)
	case _MessageTransformerOffRamp.abi.Events["SkippedReportExecution"].ID:
		return _MessageTransformerOffRamp.ParseSkippedReportExecution(log)
	case _MessageTransformerOffRamp.abi.Events["SourceChainConfigSet"].ID:
		return _MessageTransformerOffRamp.ParseSourceChainConfigSet(log)
	case _MessageTransformerOffRamp.abi.Events["SourceChainSelectorAdded"].ID:
		return _MessageTransformerOffRamp.ParseSourceChainSelectorAdded(log)
	case _MessageTransformerOffRamp.abi.Events["StaticConfigSet"].ID:
		return _MessageTransformerOffRamp.ParseStaticConfigSet(log)
	case _MessageTransformerOffRamp.abi.Events["Transmitted"].ID:
		return _MessageTransformerOffRamp.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MessageTransformerOffRampAlreadyAttempted) Topic() common.Hash {
	return common.HexToHash("0x3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120")
}

func (MessageTransformerOffRampCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4")
}

func (MessageTransformerOffRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (MessageTransformerOffRampDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0xcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee")
}

func (MessageTransformerOffRampExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b")
}

func (MessageTransformerOffRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MessageTransformerOffRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MessageTransformerOffRampRootRemoved) Topic() common.Hash {
	return common.HexToHash("0x202f1139a3e334b6056064c0e9b19fd07e44a88d8f6e5ded571b24cf8c371f12")
}

func (MessageTransformerOffRampSkippedAlreadyExecutedMessage) Topic() common.Hash {
	return common.HexToHash("0x3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c")
}

func (MessageTransformerOffRampSkippedReportExecution) Topic() common.Hash {
	return common.HexToHash("0xaab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d894933")
}

func (MessageTransformerOffRampSourceChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b")
}

func (MessageTransformerOffRampSourceChainSelectorAdded) Topic() common.Hash {
	return common.HexToHash("0xf4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9")
}

func (MessageTransformerOffRampStaticConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b7")
}

func (MessageTransformerOffRampTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_MessageTransformerOffRamp *MessageTransformerOffRamp) Address() common.Address {
	return _MessageTransformerOffRamp.address
}

type MessageTransformerOffRampInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error)

	GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error)

	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error)

	GetMessageTransformerAddress(opts *bind.CallOpts) (common.Address, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error)

	GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error)

	LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error)

	Commit(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Execute(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error)

	ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error)

	ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error)

	SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAlreadyAttempted(opts *bind.FilterOpts) (*MessageTransformerOffRampAlreadyAttemptedIterator, error)

	WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampAlreadyAttempted) (event.Subscription, error)

	ParseAlreadyAttempted(log types.Log) (*MessageTransformerOffRampAlreadyAttempted, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*MessageTransformerOffRampCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*MessageTransformerOffRampCommitReportAccepted, error)

	FilterConfigSet(opts *bind.FilterOpts) (*MessageTransformerOffRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*MessageTransformerOffRampConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts) (*MessageTransformerOffRampDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampDynamicConfigSet) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*MessageTransformerOffRampDynamicConfigSet, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*MessageTransformerOffRampExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*MessageTransformerOffRampExecutionStateChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOffRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MessageTransformerOffRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOffRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MessageTransformerOffRampOwnershipTransferred, error)

	FilterRootRemoved(opts *bind.FilterOpts) (*MessageTransformerOffRampRootRemovedIterator, error)

	WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampRootRemoved) (event.Subscription, error)

	ParseRootRemoved(log types.Log) (*MessageTransformerOffRampRootRemoved, error)

	FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*MessageTransformerOffRampSkippedAlreadyExecutedMessageIterator, error)

	WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSkippedAlreadyExecutedMessage) (event.Subscription, error)

	ParseSkippedAlreadyExecutedMessage(log types.Log) (*MessageTransformerOffRampSkippedAlreadyExecutedMessage, error)

	FilterSkippedReportExecution(opts *bind.FilterOpts) (*MessageTransformerOffRampSkippedReportExecutionIterator, error)

	WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSkippedReportExecution) (event.Subscription, error)

	ParseSkippedReportExecution(log types.Log) (*MessageTransformerOffRampSkippedReportExecution, error)

	FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*MessageTransformerOffRampSourceChainConfigSetIterator, error)

	WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error)

	ParseSourceChainConfigSet(log types.Log) (*MessageTransformerOffRampSourceChainConfigSet, error)

	FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*MessageTransformerOffRampSourceChainSelectorAddedIterator, error)

	WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampSourceChainSelectorAdded) (event.Subscription, error)

	ParseSourceChainSelectorAdded(log types.Log) (*MessageTransformerOffRampSourceChainSelectorAdded, error)

	FilterStaticConfigSet(opts *bind.FilterOpts) (*MessageTransformerOffRampStaticConfigSetIterator, error)

	WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampStaticConfigSet) (event.Subscription, error)

	ParseStaticConfigSet(log types.Log) (*MessageTransformerOffRampStaticConfigSet, error)

	FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*MessageTransformerOffRampTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *MessageTransformerOffRampTransmitted, ocrPluginType []uint8) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*MessageTransformerOffRampTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
