// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package offramp_with_message_transformer

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

var OffRampWithMessageTransformerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applySourceChainConfigUpdates\",\"inputs\":[{\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"commit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeSingleMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllSourceChainConfigs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExecutionState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestPriceSequenceNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMerkleRoot\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageTransformer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSourceChainConfig\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"ocrConfig\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"components\":[{\"name\":\"configInfo\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"n\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyExecute\",\"inputs\":[{\"name\":\"reports\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\",\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"components\":[{\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOCR3Configs\",\"inputs\":[{\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AlreadyAttempted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitReportAccepted\",\"inputs\":[{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"F\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DynamicConfigSet\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecutionStateChanged\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"state\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RootRemoved\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedAlreadyExecutedMessage\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedReportExecution\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainConfigSet\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sourceConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainSelectorAdded\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaticConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CanOnlySelfCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CommitOnRampMismatch\",\"inputs\":[{\"name\":\"reportOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"configOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyReport\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ExecutionError\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ForkedChain\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientGasForCallWithExact\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[{\"name\":\"errorType\",\"type\":\"uint8\",\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\"}]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInterval\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionGasLimit\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"newLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionTokenGasOverride\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oldLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverride\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidMessageDestChainSelector\",\"inputs\":[{\"name\":\"messageDestChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidNewState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newState\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}]},{\"type\":\"error\",\"name\":\"InvalidOnRampUpdate\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LeavesCannotBeEmpty\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionGasAmountCountMismatch\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ManualExecutionGasLimitMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionNotYetEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageTransformError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MessageValidationError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonUniqueSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotACompatiblePool\",\"inputs\":[{\"name\":\"notPool\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OracleCannotBeZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReceiverError\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ReleaseOrMintBalanceMismatch\",\"inputs\":[{\"name\":\"amountReleased\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePre\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RootAlreadyCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"RootNotCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignaturesOutOfRegistration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SourceChainNotEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SourceChainSelectorMismatch\",\"inputs\":[{\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StaleCommitReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaticConfigCannotBeChanged\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TokenDataMismatch\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TokenHandlingError\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnexpectedTokenData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongMessageLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WrongNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroChainSelectorNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461088f576167ba803803809161001d82856108c5565b8339810190808203610160811261088f5760a0811261088f5760405160a081016001600160401b038111828210176108945760405261005b836108e8565b815260208301519061ffff8216820361088f57602081019182526040840151936001600160a01b038516850361088f576040820194855261009e606082016108fc565b946060830195865260806100b38184016108fc565b84820190815295609f19011261088f57604051936100d0856108aa565b6100dc60a084016108fc565b855260c08301519363ffffffff8516850361088f576020860194855261010460e08501610910565b966040870197885261011961010086016108fc565b606088019081526101208601519095906001600160401b03811161088f5781018b601f8201121561088f5780519b6001600160401b038d11610894578c60051b91604051809e6020850161016d90836108c5565b81526020019281016020019082821161088f5760208101935b82851061078f57505050505061014061019f91016108fc565b98331561077e57600180546001600160a01b031916331790554660805284516001600160a01b031615801561076c575b801561075a575b6107385782516001600160401b0316156107495782516001600160401b0390811660a090815286516001600160a01b0390811660c0528351811660e0528451811661010052865161ffff90811661012052604080519751909416875296519096166020860152955185169084015251831660608301525190911660808201527fb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b79190a182516001600160a01b031615610738579151600480548351865160ff60c01b90151560c01b1663ffffffff60a01b60a09290921b919091166001600160a01b039485166001600160c81b0319909316831717179091558351600580549184166001600160a01b031990921691909117905560408051918252925163ffffffff166020820152935115159184019190915290511660608201529091907fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee90608090a16000915b81518310156106805760009260208160051b8401015160018060401b036020820151169081156106715780516001600160a01b031615610662578186526008602052604086206060820151916001820192610399845461091d565b610603578254600160a81b600160e81b031916600160a81b1783556040518581527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb990602090a15b805180159081156105d8575b506105c9578051906001600160401b0382116105b55761040d855461091d565b601f8111610570575b50602090601f83116001146104f85791809160008051602061679a8339815191529695949360019a9b9c926104ed575b5050600019600383901b1c191690881b1783555b60408101518254915160a089811b8a9003801960ff60a01b1990951693151590911b60ff60a01b1692909217929092169116178155610498846109da565b506104e26040519283926020845254888060a01b038116602085015260ff8160a01c1615156040850152888060401b039060a81c16606084015260808084015260a0830190610957565b0390a201919061033e565b015190503880610446565b858b52818b20919a601f198416905b8181106105585750916001999a9b849260008051602061679a833981519152989796958c951061053f575b505050811b01835561045a565b015160001960f88460031b161c19169055388080610532565b828d0151845560209c8d019c60019094019301610507565b858b5260208b20601f840160051c810191602085106105ab575b601f0160051c01905b8181106105a05750610416565b8b8155600101610593565b909150819061058a565b634e487b7160e01b8a52604160045260248afd5b6342bcdf7f60e11b8952600489fd5b9050602082012060405160208101908b8252602081526105f96040826108c5565b51902014386103ed565b825460a81c6001600160401b03166001141580610634575b156103e157632105803760e11b89526004859052602489fd5b5060405161064d816106468188610957565b03826108c5565b6020815191012081516020830120141561061b565b6342bcdf7f60e11b8652600486fd5b63c656089560e01b8652600486fd5b6001600160a01b0381161561073857600b8054600160401b600160e01b031916604092831b600160401b600160e01b031617905551615d2c9081610a6e823960805181613703015260a05181818161048f01526146d1015260c0518181816104e501528181612d3101528181613185015261466b015260e0518181816105140152614eae0152610100518181816105430152614a940152610120518181816104b6015281816124a301528181614fa10152615a610152f35b6342bcdf7f60e11b60005260046000fd5b63c656089560e01b60005260046000fd5b5081516001600160a01b0316156101d6565b5080516001600160a01b0316156101cf565b639b15e16f60e01b60005260046000fd5b84516001600160401b03811161088f5782016080818603601f19011261088f57604051906107bc826108aa565b60208101516001600160a01b038116810361088f5782526107df604082016108e8565b60208301526107f060608201610910565b604083015260808101516001600160401b03811161088f57602091010185601f8201121561088f5780516001600160401b0381116108945760405191610840601f8301601f1916602001846108c5565b818352876020838301011161088f5760005b82811061087a5750509181600060208096949581960101526060820152815201940193610186565b80602080928401015182828701015201610852565b600080fd5b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761089457604052565b601f909101601f19168101906001600160401b0382119082101761089457604052565b51906001600160401b038216820361088f57565b51906001600160a01b038216820361088f57565b5190811515820361088f57565b90600182811c9216801561094d575b602083101461093757565b634e487b7160e01b600052602260045260246000fd5b91607f169161092c565b600092918154916109678361091d565b80835292600181169081156109bd575060011461098357505050565b60009081526020812093945091925b8383106109a3575060209250010190565b600181602092949394548385870101520191019190610992565b915050602093945060ff929192191683830152151560051b010190565b80600052600760205260406000205415600014610a675760065468010000000000000000811015610894576001810180600655811015610a51577ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0181905560065460009182526007602052604090912055600190565b634e487b7160e01b600052603260045260246000fd5b5060009056fe6080604052600436101561001257600080fd5b60003560e01c806304666f9c1461016757806306285c691461016257806315777ab21461015d578063181f5a77146101585780633f4b04aa146101535780635215505b1461014e5780635e36480c146101495780635e7bb0081461014457806360987c201461013f5780637437ff9f1461013a57806379ba5097146101355780637edf52f41461013057806385572ffb1461012b5780638da5cb5b14610126578063c673e58414610121578063ccd37ba31461011c578063de5e0b9a14610117578063e9d68a8e14610112578063f2fde38b1461010d578063f58e03fc146101085763f716f99f1461010357600080fd5b61190a565b6117ed565b611762565b6116c3565b611627565b61159f565b6114f4565b61140c565b6113d6565b611210565b611190565b6110e7565b61106c565b610e67565b6108f9565b6107b4565b6106a7565b610648565b6105c1565b61045c565b61033c565b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761019d57604052565b61016c565b60a081019081106001600160401b0382111761019d57604052565b604081019081106001600160401b0382111761019d57604052565b606081019081106001600160401b0382111761019d57604052565b60c081019081106001600160401b0382111761019d57604052565b90601f801991011681019081106001600160401b0382111761019d57604052565b6040519061023e60c08361020e565b565b6040519061023e60a08361020e565b6040519061023e6101008361020e565b6040519061023e60408361020e565b6001600160401b03811161019d5760051b60200190565b6001600160a01b0381160361029657565b600080fd5b6001600160401b0381160361029657565b359061023e8261029b565b8015150361029657565b359061023e826102b7565b6001600160401b03811161019d57601f01601f191660200190565b9291926102f3826102cc565b91610301604051938461020e565b829481845281830111610296578281602093846000960137010152565b9080601f8301121561029657816020610339933591016102e7565b90565b34610296576020366003190112610296576004356001600160401b0381116102965736602382011215610296578060040135906103788261026e565b90610386604051928361020e565b8282526024602083019360051b820101903682116102965760248101935b8285106103b6576103b484611a45565b005b84356001600160401b0381116102965782016080602319823603011261029657604051916103e383610182565b60248201356103f181610285565b835260448201356104018161029b565b60208401526064820135610414816102b7565b60408401526084820135926001600160401b0384116102965761044160209493602486953692010161031e565b60608201528152019401936103a4565b600091031261029657565b3461029657600036600319011261029657610475611cf1565b506105bd604051610485816101a2565b6001600160401b037f000000000000000000000000000000000000000000000000000000000000000016815261ffff7f00000000000000000000000000000000000000000000000000000000000000001660208201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660408201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660608201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660808201526040519182918291909160806001600160a01b038160a08401956001600160401b03815116855261ffff6020820151166020860152826040820151166040860152826060820151166060860152015116910152565b0390f35b346102965760003660031901126102965760206001600160a01b03600b5460401c16604051908152f35b604051906105fa60208361020e565b60008252565b60005b8381106106135750506000910152565b8181015183820152602001610603565b9060209161063c81518092818552858086019101610600565b601f01601f1916010190565b34610296576000366003190112610296576105bd604080519061066b818361020e565b601182527f4f666652616d7020312e362e302d646576000000000000000000000000000000602083015251918291602083526020830190610623565b346102965760003660031901126102965760206001600160401b03600b5416604051908152f35b9060806060610339936001600160a01b0381511684526020810151151560208501526001600160401b0360408201511660408501520151918160608201520190610623565b6040810160408252825180915260206060830193019060005b818110610795575050506020818303910152815180825260208201916020808360051b8301019401926000915b83831061076857505050505090565b9091929394602080610786600193601f1986820301875289516106ce565b97019301930191939290610759565b82516001600160401b031685526020948501949092019160010161072c565b34610296576000366003190112610296576006546107d18161026e565b906107df604051928361020e565b808252601f196107ee8261026e565b0160005b8181106108b057505061080481611d43565b9060005b8181106108205750506105bd60405192839283610713565b8061085661083e610832600194614552565b6001600160401b031690565b6108488387611d9d565b906001600160401b03169052565b61089461088f6108766108698488611d9d565b516001600160401b031690565b6001600160401b03166000526008602052604060002090565b611e89565b61089e8287611d9d565b526108a98186611d9d565b5001610808565b6020906108bb611d1c565b828287010152016107f2565b634e487b7160e01b600052602160045260246000fd5b600411156108e757565b6108c7565b9060048210156108e75752565b3461029657604036600319011261029657602061092d60043561091b8161029b565b602435906109288261029b565b611f21565b61093a60405180926108ec565bf35b91908260a091031261029657604051610954816101a2565b608080829480358452602081013561096b8161029b565b6020850152604081013561097e8161029b565b604085015260608101356109918161029b565b60608501520135916109a28361029b565b0152565b359061023e82610285565b63ffffffff81160361029657565b359061023e826109b1565b81601f82011215610296578035906109e18261026e565b926109ef604051948561020e565b82845260208085019360051b830101918183116102965760208101935b838510610a1b57505050505090565b84356001600160401b03811161029657820160a0818503601f1901126102965760405191610a48836101a2565b60208201356001600160401b03811161029657856020610a6a9285010161031e565b83526040820135610a7a81610285565b6020840152610a8b606083016109bf565b60408401526080820135926001600160401b0384116102965760a083610ab888602080988198010161031e565b606084015201356080820152815201940193610a0c565b9190916101408184031261029657610ae561022f565b92610af0818361093c565b845260a08201356001600160401b0381116102965781610b1191840161031e565b602085015260c08201356001600160401b0381116102965781610b3591840161031e565b6040850152610b4660e083016109a6565b606085015261010082013560808501526101208201356001600160401b03811161029657610b7492016109ca565b60a0830152565b9080601f83011215610296578135610b928161026e565b92610ba0604051948561020e565b81845260208085019260051b820101918383116102965760208201905b838210610bcc57505050505090565b81356001600160401b03811161029657602091610bee87848094880101610acf565b815201910190610bbd565b81601f8201121561029657803590610c108261026e565b92610c1e604051948561020e565b82845260208085019360051b830101918183116102965760208101935b838510610c4a57505050505090565b84356001600160401b03811161029657820183603f82011215610296576020810135610c758161026e565b91610c83604051938461020e565b8183526020808085019360051b83010101918683116102965760408201905b838210610cbc575050509082525060209485019401610c3b565b81356001600160401b03811161029657602091610ce08a848080958901010161031e565b815201910190610ca2565b929190610cf78161026e565b93610d05604051958661020e565b602085838152019160051b810192831161029657905b828210610d2757505050565b8135815260209182019101610d1b565b9080601f830112156102965781602061033993359101610ceb565b81601f8201121561029657803590610d698261026e565b92610d77604051948561020e565b82845260208085019360051b830101918183116102965760208101935b838510610da357505050505090565b84356001600160401b03811161029657820160a0818503601f19011261029657610dcb610240565b91610dd8602083016102ac565b835260408201356001600160401b03811161029657856020610dfc92850101610b7b565b602084015260608201356001600160401b03811161029657856020610e2392850101610bf9565b60408401526080820135926001600160401b0384116102965760a083610e50886020809881980101610d37565b606084015201356080820152815201940193610d94565b34610296576040366003190112610296576004356001600160401b03811161029657610e97903690600401610d52565b6024356001600160401b038111610296573660238201121561029657806004013591610ec28361026e565b91610ed0604051938461020e565b8383526024602084019460051b820101903682116102965760248101945b828610610eff576103b48585611f69565b85356001600160401b03811161029657820136604382011215610296576024810135610f2a8161026e565b91610f38604051938461020e565b818352602060248185019360051b83010101903682116102965760448101925b828410610f72575050509082525060209586019501610eee565b83356001600160401b038111610296576024908301016040601f1982360301126102965760405190610fa3826101bd565b6020810135825260408101356001600160401b03811161029657602091010136601f8201121561029657803590610fd98261026e565b91610fe7604051938461020e565b80835260208084019160051b8301019136831161029657602001905b8282106110225750505091816020938480940152815201930192610f58565b602080918335611031816109b1565b815201910190611003565b9181601f84011215610296578235916001600160401b038311610296576020808501948460051b01011161029657565b34610296576060366003190112610296576004356001600160401b0381116102965761109c903690600401610acf565b6024356001600160401b038111610296576110bb90369060040161103c565b91604435926001600160401b038411610296576110df6103b494369060040161103c565b93909261237c565b3461029657600036600319011261029657611100612654565b506105bd60405161111081610182565b60ff6004546001600160a01b038116835263ffffffff8160a01c16602084015260c01c16151560408201526001600160a01b036005541660608201526040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b34610296576000366003190112610296576000546001600160a01b03811633036111ff576001600160a01b0319600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b63015aa1e360e11b60005260046000fd5b3461029657608036600319011261029657600060405161122f81610182565b60043561123b81610285565b8152602435611249816109b1565b602082015260443561125a816102b7565b604082015260643561126b81610285565b6060820152611278613500565b6001600160a01b03815116156113c7576113c1816112d76001600160a01b037fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9451166001600160a01b03166001600160a01b03196004541617600455565b60208101516004547fffffffffffffff0000000000ffffffffffffffffffffffffffffffffffffffff77ffffffff000000000000000000000000000000000000000078ff0000000000000000000000000000000000000000000000006040860151151560c01b169360a01b169116171760045561137d61136160608301516001600160a01b031690565b6001600160a01b03166001600160a01b03196005541617600555565b6040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b0390a180f35b6342bcdf7f60e11b8252600482fd5b34610296576020366003190112610296576004356001600160401b0381116102965760a090600319903603011261029657600080fd5b346102965760003660031901126102965760206001600160a01b0360015416604051908152f35b6004359060ff8216820361029657565b359060ff8216820361029657565b906020808351928381520192019060005b81811061146f5750505090565b82516001600160a01b0316845260209384019390920191600101611462565b906103399160208152606082518051602084015260ff602082015116604084015260ff6040820151168284015201511515608082015260406114df602084015160c060a085015260e0840190611451565b9201519060c0601f1982850301910152611451565b346102965760203660031901126102965760ff61150f611433565b60606040805161151e816101d8565b611526612654565b815282602082015201521660005260026020526105bd6040600020600361158e60405192611553846101d8565b61155c81612679565b84526040516115798161157281600286016126b2565b038261020e565b602085015261157260405180948193016126b2565b60408201526040519182918261148e565b34610296576040366003190112610296576004356115bc8161029b565b6001600160401b036024359116600052600a6020526040600020906000526020526020604060002054604051908152f35b9060049160441161029657565b9181601f84011215610296578235916001600160401b038311610296576020838186019501011161029657565b346102965760c036600319011261029657611641366115ed565b6044356001600160401b038111610296576116609036906004016115fa565b6064929192356001600160401b0381116102965761168290369060040161103c565b60843594916001600160401b038611610296576116a66103b496369060040161103c565b94909360a43596612cec565b9060206103399281815201906106ce565b34610296576020366003190112610296576001600160401b036004356116e88161029b565b6116f0611d1c565b501660005260086020526105bd604060002060016117516040519261171484610182565b6001600160401b0381546001600160a01b038116865260ff8160a01c161515602087015260a81c1660408501526115726040518094819301611deb565b6060820152604051918291826116b2565b34610296576020366003190112610296576001600160a01b0360043561178781610285565b61178f613500565b163381146117dc57806001600160a01b031960005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b636d6c4ee560e11b60005260046000fd5b3461029657606036600319011261029657611807366115ed565b6044356001600160401b038111610296576118269036906004016115fa565b91828201602083820312610296578235906001600160401b03821161029657611850918401610d52565b604051906020611860818461020e565b60008352601f19810160005b818110611894575050506103b4949161188491613744565b61188c6131fb565b928392613f1a565b6060858201840152820161186c565b9080601f830112156102965781356118ba8161026e565b926118c8604051948561020e565b81845260208085019260051b82010192831161029657602001905b8282106118f05750505090565b6020809183356118ff81610285565b8152019101906118e3565b34610296576020366003190112610296576004356001600160401b0381116102965736602382011215610296578060040135906119468261026e565b90611954604051928361020e565b8282526024602083019360051b820101903682116102965760248101935b828510611982576103b484613217565b84356001600160401b03811161029657820160c06023198236030112610296576119aa61022f565b91602482013583526119be60448301611443565b60208401526119cf60648301611443565b60408401526119e0608483016102c1565b606084015260a48201356001600160401b03811161029657611a0890602436918501016118a3565b608084015260c4820135926001600160401b03841161029657611a356020949360248695369201016118a3565b60a0820152815201940193611972565b611a4d613500565b60005b8151811015611ced57611a638183611d9d565b5190611a7960208301516001600160401b031690565b916001600160401b038316908115611cdc57611aae611aa2611aa283516001600160a01b031690565b6001600160a01b031690565b15611c4357611ad0846001600160401b03166000526008602052604060002090565b906060810151916001810195611ae68754611db1565b611c6a57611b597ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb991611b3f84750100000000000000000000000000000000000000000067ffffffffffffffff60a81b19825416179055565b6040516001600160401b0390911681529081906020820190565b0390a15b82518015908115611c54575b50611c4357611c24611c08611c3a93611ba57f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b9660019a6135a2565b611bfb611bb56040830151151590565b85547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690151560a01b74ff000000000000000000000000000000000000000016178555565b516001600160a01b031690565b82906001600160a01b03166001600160a01b0319825416179055565b611c2d8461557a565b5060405191829182613673565b0390a201611a50565b6342bcdf7f60e11b60005260046000fd5b90506020840120611c63613525565b1438611b69565b60016001600160401b03611c8984546001600160401b039060a81c1690565b16141580611cbd575b611c9c5750611b5d565b632105803760e11b6000526001600160401b031660045260246000fd5b6000fd5b50611cc787611e6e565b60208151910120845160208601201415611c92565b63c656089560e01b60005260046000fd5b5050565b60405190611cfe826101a2565b60006080838281528260208201528260408201528260608201520152565b60405190611d2982610182565b606080836000815260006020820152600060408201520152565b90611d4d8261026e565b611d5a604051918261020e565b8281528092611d6b601f199161026e565b0190602036910137565b634e487b7160e01b600052603260045260246000fd5b805115611d985760200190565b611d75565b8051821015611d985760209160051b010190565b90600182811c92168015611de1575b6020831014611dcb57565b634e487b7160e01b600052602260045260246000fd5b91607f1691611dc0565b60009291815491611dfb83611db1565b8083529260018116908115611e515750600114611e1757505050565b60009081526020812093945091925b838310611e37575060209250010190565b600181602092949394548385870101520191019190611e26565b915050602093945060ff929192191683830152151560051b010190565b9061023e611e829260405193848092611deb565b038361020e565b9060016060604051611e9a81610182565b6109a281956001600160401b0381546001600160a01b038116855260ff8160a01c161515602086015260a81c166040840152611edc6040518096819301611deb565b038461020e565b634e487b7160e01b600052601160045260246000fd5b908160051b9180830460201490151715611f0f57565b611ee3565b91908203918211611f0f57565b611f2d82607f926136bd565b9116906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611f0f576003911c1660048110156108e75790565b611f71613701565b8051825181036121685760005b818110611f915750509061023e91613744565b611f9b8184611d9d565b516020810190815151611fae8488611d9d565b5192835182036121685790916000925b808410611fd2575050505050600101611f7e565b91949398611fe4848b98939598611d9d565b515198611ff2888851611d9d565b51998061211f575b5060a08a01988b60206120108b8d515193611d9d565b51015151036120e25760005b8a51518110156120cd5761205861204f6120458f602061203d8f8793611d9d565b510151611d9d565b5163ffffffff1690565b63ffffffff1690565b8b81612069575b505060010161201c565b61204f604061207c856120889451611d9d565b51015163ffffffff1690565b9081811061209757508b61205f565b8d51516040516348e617b360e01b81526004810191909152602481019390935260448301919091526064820152608490fd5b0390fd5b50985098509893949095600101929091611fbe565b611cb98b516120fd606082519201516001600160401b031690565b6370a193fd60e01b6000526004919091526001600160401b0316602452604490565b60808b0151811015611ffa57611cb9908b61214188516001600160401b031690565b905151633a98d46360e11b6000526001600160401b03909116600452602452604452606490565b6320f8fd5960e21b60005260046000fd5b60405190612186826101bd565b60006020838281520152565b604051906121a160208361020e565b600080835282815b8281106121b557505050565b6020906121c0612179565b828285010152016121a9565b805182526001600160401b0360208201511660208301526080612213612201604084015160a0604087015260a0860190610623565b60608401518582036060870152610623565b9101519160808183039101526020808351928381520192019060005b81811061223c5750505090565b825180516001600160a01b03168552602090810151818601526040909401939092019160010161222f565b9060206103399281815201906121cc565b6040513d6000823e3d90fd5b3d156122af573d90612295826102cc565b916122a3604051938461020e565b82523d6000602084013e565b606090565b906020610339928181520190610623565b81601f820112156102965780516122db816102cc565b926122e9604051948561020e565b81845260208284010111610296576103399160208085019101610600565b909160608284031261029657815161231e816102b7565b9260208301516001600160401b038111610296576040916123409185016122c5565b92015190565b9293606092959461ffff61236a6001600160a01b03946080885260808801906121cc565b97166020860152604085015216910152565b93919290933033036126435761239190613b94565b9261239a612192565b9460a085015180516125fc575b50505050508051916123c5602084519401516001600160401b031690565b9060208301519160408401926123f28451926123df610240565b9788526001600160401b03166020880152565b6040860152606085015260808401526001600160a01b0361241b6005546001600160a01b031690565b168061257f575b5051511580612573575b801561255d575b8015612534575b611ced576124cc9181612471611aa2612464610876602060009751016001600160401b0390511690565b546001600160a01b031690565b908361248c606060808401519301516001600160a01b031690565b604051633cf9798360e01b815296879586948593917f00000000000000000000000000000000000000000000000000000000000000009060048601612346565b03925af190811561252f57600090600092612508575b50156124eb5750565b6040516302a35ba360e21b81529081906120c990600483016122b4565b905061252791503d806000833e61251f818361020e565b810190612307565b5090386124e2565b612278565b5061255861255461254f60608401516001600160a01b031690565b613ddb565b1590565b61243a565b5060608101516001600160a01b03163b15612433565b5060808101511561242c565b803b1561029657600060405180926308d450a160e01b82528183816125a78a60048301612267565b03925af190816125e1575b506125db576120c96125c2612284565b6040516309c2532560e01b8152918291600483016122b4565b38612422565b806125f060006125f69361020e565b80610451565b386125b2565b859650602061263896015161261b60608901516001600160a01b031690565b9061263260208a51016001600160401b0390511690565b92613cc2565b9038808080806123a7565b6306e34e6560e31b60005260046000fd5b6040519061266182610182565b60006060838281528260208201528260408201520152565b9060405161268681610182565b606060ff600183958054855201548181166020850152818160081c16604085015260101c161515910152565b906020825491828152019160005260206000209060005b8181106126d65750505090565b82546001600160a01b03168452602090930192600192830192016126c9565b9061023e611e8292604051938480926126b2565b35906001600160e01b038216820361029657565b81601f82011215610296578035906127348261026e565b92612742604051948561020e565b82845260208085019360061b8301019181831161029657602001925b82841061276c575050505090565b6040848303126102965760206040918251612786816101bd565b86356127918161029b565b815261279e838801612709565b8382015281520193019261275e565b81601f82011215610296578035906127c48261026e565b926127d2604051948561020e565b82845260208085019360051b830101918183116102965760208101935b8385106127fe57505050505090565b84356001600160401b03811161029657820160a0818503601f190112610296576040519161282b836101a2565b60208201356128398161029b565b83526040820135926001600160401b0384116102965760a08361286388602080988198010161031e565b8584015260608101356128758161029b565b604084015260808101356128888161029b565b6060840152013560808201528152019401936127ef565b81601f82011215610296578035906128b68261026e565b926128c4604051948561020e565b82845260208085019360061b8301019181831161029657602001925b8284106128ee575050505090565b6040848303126102965760206040918251612908816101bd565b8635815282870135838201528152019301926128e0565b602081830312610296578035906001600160401b03821161029657016060818303126102965760405191612952836101d8565b81356001600160401b038111610296578201604081830312610296576040519061297b826101bd565b80356001600160401b03811161029657810183601f820112156102965780356129a38161026e565b916129b1604051938461020e565b81835260208084019260061b8201019086821161029657602001915b818310612a495750505082526020810135906001600160401b038211610296576129f99184910161271d565b6020820152835260208201356001600160401b0381116102965781612a1f9184016127ad565b602084015260408201356001600160401b03811161029657612a41920161289f565b604082015290565b6040838803126102965760206040918251612a63816101bd565b8535612a6e81610285565b8152612a7b838701612709565b838201528152019201916129cd565b9080602083519182815201916020808360051b8301019401926000915b838310612ab657505050505090565b9091929394602080600192601f198582030186528851906001600160401b038251168152608080612af48585015160a08786015260a0850190610623565b936001600160401b0360408201511660408501526001600160401b036060820151166060850152015191015297019301930191939290612aa7565b916001600160a01b03612b5092168352606060208401526060830190612a8a565b9060408183039101526020808351928381520192019060005b818110612b765750505090565b8251805185526020908101518186015260409094019390920191600101612b69565b906020808351928381520192019060005b818110612bb65750505090565b825180516001600160401b031685526020908101516001600160e01b03168186015260409094019390920191600101612ba9565b9190604081019083519160408252825180915260206060830193019060005b818110612c2a57505050602061033993940151906020818403910152612b98565b825180516001600160a01b031686526020908101516001600160e01b03168187015260409095019490920191600101612c09565b906020610339928181520190612bea565b908160209103126102965751610339816102b7565b9091612c9b61033993604084526040840190610623565b916020818403910152611deb565b6001600160401b036001911601906001600160401b038211611f0f57565b9091612cde61033993604084526040840190612a8a565b916020818403910152612bea565b929693959190979497612d018282018261291f565b98612d1561255460045460ff9060c01c1690565b613169575b8951805151159081159161315a575b50613081575b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316999860208a019860005b8a51805182101561301f5781612d7891611d9d565b518d612d8b82516001600160401b031690565b604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b1660048201529091602090829060249082905afa90811561252f57600091612ff1575b50612fd457612dda90613e29565b60208201805160208151910120906001830191612df683611e6e565b6020815191012003612fb7575050805460408301516001600160401b039081169160a81c168114801590612f8f575b612f3d57506080820151908115612f2c57612e7682612e67612e4e86516001600160401b031690565b6001600160401b0316600052600a602052604060002090565b90600052602052604060002090565b54612ef8578291612edc612ef192612ea3612e9e60606001999801516001600160401b031690565b612ca9565b67ffffffffffffffff60a81b197cffffffffffffffff00000000000000000000000000000000000000000083549260a81b169116179055565b612e67612e4e4294516001600160401b031690565b5501612d63565b50612f0d611cb992516001600160401b031690565b6332cf0cbf60e01b6000526001600160401b0316600452602452604490565b63504570e360e01b60005260046000fd5b82611cb991612f676060612f5884516001600160401b031690565b9301516001600160401b031690565b636af0786b60e11b6000526001600160401b0392831660045290821660245216604452606490565b50612fa761083260608501516001600160401b031690565b6001600160401b03821611612e25565b516120c960405192839263b80d8fa960e01b845260048401612c84565b637edeb53960e11b6000526001600160401b031660045260246000fd5b613012915060203d8111613018575b61300a818361020e565b810190612c6f565b38612dcc565b503d613000565b505061307b9496989b507f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e461023e9b613073949597999b5190519061306960405192839283612cc7565b0390a13691610ceb565b943691610ceb565b93614214565b613096602086015b356001600160401b031690565b600b546001600160401b038281169116101561313e576130cc906001600160401b03166001600160401b0319600b541617600b55565b6130e4611aa2611aa26004546001600160a01b031690565b8a5190803b1561029657604051633937306f60e01b81529160009183918290849082906131149060048301612c5e565b03925af1801561252f57613129575b50612d2f565b806125f060006131389361020e565b38613123565b5060208a015151612d2f57632261116760e01b60005260046000fd5b60209150015151151538612d29565b60208a0151805161317b575b50612d1a565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169060408c0151823b1561029657604051633854844f60e11b8152926000928492839182916131d7913060048501612b2f565b03915afa801561252f571561317557806125f060006131f59361020e565b38613175565b6040519061320a60208361020e565b6000808352366020840137565b61321f613500565b60005b8151811015611ced576132358183611d9d565b51906040820160ff613248825160ff1690565b16156134ea57602083015160ff169261326e8460ff166000526002602052604060002090565b91600183019182546132896132838260ff1690565b60ff1690565b6134af57506132b661329e6060830151151590565b845462ff0000191690151560101b62ff000016178455565b60a081019182516101008151116134575780511561349957600386016132e46132de826126f5565b8a615328565b6060840151613374575b947fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f5479460029461335061334061336e9a966133398760019f9c6133346133669a8f615489565b614455565b5160ff1690565b845460ff191660ff821617909455565b51908185555190604051958695019088866144db565b0390a161550b565b01613222565b9794600287939597019661339061338a896126f5565b88615328565b6080850151946101008651116134835785516133b86132836133b38a5160ff1690565b614441565b101561346d578551845111613457576133506133407fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547986133398760019f61333461336e9f9a8f61343f60029f6134396133669f8f90613334849261341e845160ff1690565b908054909161ff001990911660089190911b61ff0016179055565b826153bc565b505050979c9f50975050969a505050945094506132ee565b631b3fab5160e11b600052600160045260246000fd5b631b3fab5160e11b600052600360045260246000fd5b631b3fab5160e11b600052600260045260246000fd5b631b3fab5160e11b600052600560045260246000fd5b60101c60ff166134ca6134c56060840151151590565b151590565b901515146132b6576321fd80df60e21b60005260ff861660045260246000fd5b631b3fab5160e11b600090815260045260246000fd5b6001600160a01b0360015416330361351457565b6315ae3a6f60e11b60005260046000fd5b6040516020810190600082526020815261354060408261020e565b51902090565b818110613551575050565b60008155600101613546565b9190601f811161356c57505050565b61023e926000526020600020906020601f840160051c83019310613598575b601f0160051c0190613546565b909150819061358b565b91909182516001600160401b03811161019d576135c9816135c38454611db1565b8461355d565b6020601f821160011461360a5781906135fb9394956000926135ff575b50508160011b916000199060031b1c19161790565b9055565b0151905038806135e6565b601f1982169061361f84600052602060002090565b9160005b81811061365b57509583600195969710613642575b505050811b019055565b015160001960f88460031b161c19169055388080613638565b9192602060018192868b015181550194019201613623565b90600160a061033993602081526001600160401b0384546001600160a01b038116602084015260ff81851c161515604084015260a81c166060820152608080820152019101611deb565b906001600160401b036136fd921660005260096020526701ffffffffffffff60406000209160071c166001600160401b0316600052602052604060002090565b5490565b7f000000000000000000000000000000000000000000000000000000000000000046810361372c5750565b630f01ce8560e01b6000526004524660245260446000fd5b9190918051156137e6578251159260209160405192613763818561020e565b60008452601f19810160005b8181106137c25750505060005b81518110156137ba57806137a361379560019385611d9d565b5188156137a957869061461a565b0161377c565b6137b38387611d9d565b519061461a565b505050509050565b82906040516137d0816101bd565b600081526060838201528282890101520161376f565b63c2e5347d60e01b60005260046000fd5b91908260a09103126102965760405161380f816101a2565b60808082948051845260208101516138268161029b565b602085015260408101516138398161029b565b6040850152606081015161384c8161029b565b60608501520151916109a28361029b565b519061023e82610285565b519061023e826109b1565b81601f820112156102965780519061388a8261026e565b92613898604051948561020e565b82845260208085019360051b830101918183116102965760208101935b8385106138c457505050505090565b84516001600160401b03811161029657820160a0818503601f19011261029657604051916138f1836101a2565b60208201516001600160401b03811161029657856020613913928501016122c5565b8352604082015161392381610285565b602084015261393460608301613868565b60408401526080820151926001600160401b0384116102965760a0836139618860208098819801016122c5565b6060840152015160808201528152019401936138b5565b602081830312610296578051906001600160401b038211610296570161014081830312610296576139a761022f565b916139b281836137f7565b835260a08201516001600160401b03811161029657816139d39184016122c5565b602084015260c08201516001600160401b03811161029657816139f79184016122c5565b6040840152613a0860e0830161385d565b606084015261010082015160808401526101208201516001600160401b03811161029657613a369201613873565b60a082015290565b9080602083519182815201916020808360051b8301019401926000915b838310613a6a57505050505090565b9091929394602080600192601f19858203018652885190608080613acd613a9a855160a0865260a0860190610623565b6001600160a01b0387870151168786015263ffffffff604087015116604086015260608601518582036060870152610623565b93015191015297019301930191939290613a5b565b610339916001600160401b036080835180518452826020820151166020850152826040820151166040850152826060820151166060850152015116608082015260a0613b53613b41602085015161014084860152610140850190610623565b604085015184820360c0860152610623565b60608401516001600160a01b031660e0840152926080810151610100840152015190610120818403910152613a3e565b906020610339928181520190613ae2565b6000613c0c8192604051613ba7816101f3565b613baf611cf1565b81526060602082015260606040820152836060820152836080820152606060a082015250613bef611aa2611aa2600b546001600160a01b039060401c1690565b90604051948580948193634546c6e560e01b835260048301613b83565b03925af160009181613c42575b50610339576120c9613c29612284565b60405163828ebdfb60e01b8152918291600483016122b4565b613c609192503d806000833e613c58818361020e565b810190613978565b9038613c19565b9190811015611d985760051b0190565b35610339816109b1565b9190811015611d985760051b81013590601e19813603018212156102965701908135916001600160401b038311610296576020018236038113610296579190565b90929491939796815196613cd58861026e565b97613ce3604051998a61020e565b808952613cf2601f199161026e565b0160005b818110613dc457505060005b8351811015613db75780613d498c8a8a8a613d43613d3c878d613d35828f8f9d8f9e60019f81613d65575b505050611d9d565b5197613c81565b36916102e7565b93614e5f565b613d53828c611d9d565b52613d5e818b611d9d565b5001613d02565b63ffffffff613d7d613d78858585613c67565b613c77565b1615613d2d57613dad92613d9492613d7892613c67565b6040613da08585611d9d565b51019063ffffffff169052565b8f8f908391613d2d565b5096985050505050505050565b602090613dcf612179565b82828d01015201613cf6565b613dec6385572ffb60e01b826151c2565b9081613e06575b81613dfc575090565b6103399150615194565b9050613e1181615119565b1590613df3565b613dec63aff2afbf60e01b826151c2565b6001600160401b031680600052600860205260406000209060ff825460a01c1615613e52575090565b63ed053c5960e01b60005260045260246000fd5b6084019081608411611f0f57565b60a001908160a011611f0f57565b91908201809211611f0f57565b600311156108e757565b60038210156108e75752565b9061023e604051613eb5816101bd565b602060ff829554818116845260081c169101613e99565b8054821015611d985760005260206000200190600090565b60ff60019116019060ff8211611f0f57565b60ff601b9116019060ff8211611f0f57565b90606092604091835260208301370190565b6001600052600260205293613f4e7fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0612679565b93853594613f5b85613e66565b6060820190613f6a8251151590565b6141e6575b8036036141ce575081518781036141b55750613f89613701565b60016000526003602052613fd8613fd37fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c5b336001600160a01b0316600052602052604060002090565b613ea5565b60026020820151613fe881613e8f565b613ff181613e8f565b14908161414d575b5015614121575b51614058575b50505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09061403c61308960019460200190565b604080519283526001600160401b0391909116602083015290a2565b614079613283614074602085969799989a955194015160ff1690565b613ee4565b036141105781518351036140ff576140f7600061403c94613089946140c37f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09960019b36916102e7565b602081519101206040516140ee816140e089602083019586613f08565b03601f19810183528261020e565b5190208a6151f2565b948394614006565b63a75d88af60e01b60005260046000fd5b6371253a2560e01b60005260046000fd5b72c11c11c11c11c11c11c11c11c11c11c11c11c133031561400057631b41e11d60e31b60005260046000fd5b600160005260026020526141ad9150611aa29061419a9061419460037fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e05b01915160ff1690565b90613ecc565b90546001600160a01b039160031b1c1690565b331438613ff9565b6324f7d61360e21b600052600452602487905260446000fd5b638e1192e160e01b6000526004523660245260446000fd5b61420f906142096141ff6141fa8751611ef9565b613e74565b6142098851611ef9565b90613e82565b613f6f565b6000805260026020529490939092909161424d7fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b612679565b9486359561425a83613e66565b60608201906142698251151590565b61441e575b8036036141ce575081518881036144055750614288613701565b6000805260036020526142bd613fd37f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff613fbb565b600260208201516142cd81613e8f565b6142d681613e8f565b1490816143bc575b5015614390575b51614322575b5050505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09061403c61308960009460200190565b61433e613283614074602087989a999b96975194015160ff1690565b036141105783518651036140ff576000967f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09661403c956140c3614387946130899736916102e7565b948394386142eb565b72c11c11c11c11c11c11c11c11c11c11c11c11c13303156142e557631b41e11d60e31b60005260046000fd5b6000805260026020526143fd9150611aa29061419a9061419460037fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b61418b565b3314386142de565b6324f7d61360e21b600052600452602488905260446000fd5b61443c906142096144326141fa8951611ef9565b6142098a51611ef9565b61426e565b60ff166003029060ff8216918203611f0f57565b8151916001600160401b03831161019d5768010000000000000000831161019d5760209082548484558085106144be575b500190600052602060002060005b8381106144a15750505050565b60019060206001600160a01b038551169401938184015501614494565b6144d5908460005285846000209182019101613546565b38614486565b95949392909160ff61450093168752602087015260a0604087015260a08601906126b2565b84810360608601526020808351928381520192019060005b8181106145335750505090608061023e9294019060ff169052565b82516001600160a01b0316845260209384019390920191600101614518565b600654811015611d985760066000527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f015490565b6001600160401b036103399493816060941683521660208201528160408201520190610623565b604090610339939281528160208201520190610623565b9291906001600160401b039081606495166004521660245260048110156108e757604452565b94939261460460609361461593885260208801906108ec565b608060408701526080860190610623565b930152565b9061462c82516001600160401b031690565b8151604051632cbc26bb60e01b815267ffffffffffffffff60801b608084901b1660048201529015159391906001600160401b038216906020816024817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa90811561252f57600091614d48575b50614d06576020830191825151948515614cd657604085018051518703614cc5576146ce87611d43565b957f000000000000000000000000000000000000000000000000000000000000000061470460016146fe87613e29565b01611e6e565b60208151910120604051614764816140e06020820194868b876001600160401b036060929594938160808401977f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f85521660208401521660408201520152565b519020906001600160401b031660005b8a8110614c2d575050508060806060614794930151910151908886615717565b978815614c0f5760005b8881106147b15750505050505050505050565b5a6147bd828951611d9d565b518051606001516147d7906001600160401b031688611f21565b6147e0816108dd565b8015908d8283159384614bfc575b15614bb95760608815614b3c5750614815602061480b898d611d9d565b5101519242611f14565b60045461482a9060a01c63ffffffff1661204f565b108015614b29575b15614b0b57614841878b611d9d565b5151614af5575b845160800151614860906001600160401b0316610832565b614a3d575b50614871868951611d9d565b5160a085015151815103614a0157936148d69695938c938f966148b68e958c926148b06148aa60608951016001600160401b0390511690565b89615749565b8661594a565b9a9080966148d060608851016001600160401b0390511690565b906157d1565b6149af575b50506148e6826108dd565b60028203614967575b60019661495d7f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b936001600160401b0393519261494e6149458b61493d60608801516001600160401b031690565b96519b611d9d565b51985a90611f14565b916040519586951698856145eb565b0390a45b0161479e565b91509193949250614977826108dd565b6003820361498b578b929493918a916148ef565b51606001516349362d1f60e11b600052611cb991906001600160401b0316896145c5565b6149b8846108dd565b600384036148db5790929495506149d09193506108dd565b6149e0578b92918a9138806148db565b5151604051632b11b8d960e01b81529081906120c9908790600484016145ae565b611cb98b614a1b60608851016001600160401b0390511690565b631cfe6d8b60e01b6000526001600160401b0391821660045216602452604490565b614a46836108dd565b614a51575b38614865565b8351608001516001600160401b0316602080860151918c614a8660405194859384936370701e5760e11b855260048501614587565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af190811561252f57600091614ad7575b50614a4b575050505050600190614961565b614aef915060203d81116130185761300a818361020e565b38614ac5565b614aff878b611d9d565b51516080860152614848565b6354e7e43160e11b6000526001600160401b038b1660045260246000fd5b50614b33836108dd565b60038314614832565b915083614b48846108dd565b1561484857506001959450614bb19250614b8f91507f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe651209351016001600160401b0390511690565b604080516001600160401b03808c168252909216602083015290918291820190565b0390a1614961565b505050506001929150614bb1614b8f60607f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c9351016001600160401b0390511690565b50614c06836108dd565b600383146147ee565b633ee8bd3f60e11b6000526001600160401b03841660045260246000fd5b614c38818a51611d9d565b518051604001516001600160401b0316838103614ca857508051602001516001600160401b0316898103614c85575090614c748460019361560f565b614c7e828d611d9d565b5201614774565b636c95f1eb60e01b6000526001600160401b03808a166004521660245260446000fd5b631c21951160e11b6000526001600160401b031660045260246000fd5b6357e0e08360e01b60005260046000fd5b611cb9614cea86516001600160401b031690565b63676cf24b60e11b6000526001600160401b0316600452602490565b5092915050612fd4576040516001600160401b039190911681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d89493390602090a1565b614d61915060203d6020116130185761300a818361020e565b386146a4565b90816020910312610296575161033981610285565b90610339916020815260e0614e1a614e05614da585516101006020870152610120860190610623565b60208601516001600160401b0316604086015260408601516001600160a01b0316606086015260608601516080860152614def608087015160a08701906001600160a01b03169052565b60a0860151858203601f190160c0870152610623565b60c0850151848203601f190184860152610623565b92015190610100601f1982850301910152610623565b6040906001600160a01b0361033994931681528160208201520190610623565b90816020910312610296575190565b91939293614e6b612179565b5060208301516001600160a01b031660405163bbe4f6db60e01b81526001600160a01b038216600482015290959092602084806024810103816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa93841561252f576000946150e8575b506001600160a01b03841695861580156150d6575b6150b857614f9d614fc6926140e092614f21614f1a61204f60408c015163ffffffff1690565b8c89615a29565b9690996080810151614f4f6060835193015193614f3c61024f565b9687526001600160401b03166020870152565b6001600160a01b038a16604086015260608501526001600160a01b038d16608085015260a084015260c083015260e0820152604051633907753760e01b602082015292839160248301614d7c565b82857f000000000000000000000000000000000000000000000000000000000000000092615ab7565b9490911561509c5750805160208103615083575090614fef826020808a95518301019101614e50565b956001600160a01b03841603615027575b505050505061501f61501061025f565b6001600160a01b039093168352565b602082015290565b61503a9361503491611f14565b91615a29565b50908082108015615070575b61505257808481615000565b63a966e21f60e01b6000908152600493909352602452604452606490fd5b508261507c8284611f14565b1415615046565b631e3be00960e21b600052602060045260245260446000fd5b6120c9604051928392634ff17cad60e11b845260048401614e30565b63ae9b4ce960e01b6000526001600160a01b03851660045260246000fd5b506150e361255486613e18565b614ef4565b61510b91945060203d602011615112575b615103818361020e565b810190614d67565b9238614edf565b503d6150f9565b60405160208101916301ffc9a760e01b835263ffffffff60e01b60248301526024825261514760448361020e565b6179185a10615183576020926000925191617530fa6000513d82615177575b5081615170575090565b9050151590565b60201115915038615166565b63753fa58960e11b60005260046000fd5b60405160208101916301ffc9a760e01b83526301ffc9a760e01b60248301526024825261514760448361020e565b6040519060208201926301ffc9a760e01b845263ffffffff60e01b1660248301526024825261514760448361020e565b919390926000948051946000965b868810615211575050505050505050565b6020881015611d985760206000615229878b1a613ef6565b6152338b87611d9d565b519061526a6152428d8a611d9d565b5160405193849389859094939260ff6060936080840197845216602083015260408201520152565b838052039060015afa1561252f576152b0613fd36000516152988960ff166000526003602052604060002090565b906001600160a01b0316600052602052604060002090565b90600160208301516152c181613e8f565b6152ca81613e8f565b03615317576152e76152dd835160ff1690565b60ff600191161b90565b8116615306576152fd6152dd6001935160ff1690565b17970196615200565b633d9ef1f160e21b60005260046000fd5b636518c33d60e11b60005260046000fd5b91909160005b83518110156153815760019060ff83166000526003602052600061537a604082206001600160a01b03615361858a611d9d565b51166001600160a01b0316600052602052604060002090565b550161532e565b50509050565b8151815460ff191660ff91909116178155906020015160038110156108e757815461ff00191660089190911b61ff0016179055565b919060005b8151811015615381576153d7611bfb8284611d9d565b906154006153f6836152988860ff166000526003602052604060002090565b5460081c60ff1690565b61540981613e8f565b615474576001600160a01b038216156154635761545d60019261545861542d61025f565b60ff85168152916154418660208501613e99565b6152988960ff166000526003602052604060002090565b615387565b016153c1565b63d6c62c9b60e01b60005260046000fd5b631b3fab5160e11b6000526004805260246000fd5b919060005b8151811015615381576154a4611bfb8284611d9d565b906154c36153f6836152988860ff166000526003602052604060002090565b6154cc81613e8f565b615474576001600160a01b03821615615463576155056001926154586154f061025f565b60ff8516815291615441600260208501613e99565b0161548e565b60ff1680600052600260205260ff60016040600020015460101c16908015600014615559575015615548576001600160401b0319600b5416600b55565b6317bd8dd160e11b60005260046000fd5b6001146155635750565b61556957565b6307b8c74d60e51b60005260046000fd5b806000526007602052604060002054156000146155f8576006546801000000000000000081101561019d57600181016006556000600654821015611d9857600690527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01819055600654906000526007602052604060002055600190565b50600090565b906020610339928181520190613a3e565b61354081518051906156a361562e60608601516001600160a01b031690565b6140e061564560608501516001600160401b031690565b9361565e6080808a01519201516001600160401b031690565b90604051958694602086019889936001600160401b036080946001600160a01b0382959998949960a089019a8952166020880152166040860152606085015216910152565b5190206140e06020840151602081519101209360a06040820151602081519101209101516040516156dc816140e06020820194856155fe565b51902090604051958694602086019889919260a093969594919660c08401976000855260208501526040840152606083015260808201520152565b926001600160401b039261572a92615b74565b9116600052600a60205260406000209060005260205260406000205490565b607f8216906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611f0f576157ce916001600160401b0361578c85846136bd565b921660005260096020526701ffffffffffffff60406000209460071c169160036001831b921b19161792906001600160401b0316600052602052604060002090565b55565b9091607f83166801fffffffffffffffe6001600160401b0382169160011b169080820460021490151715611f0f5761580984846136bd565b60048310156108e7576001600160401b036157ce9416600052600960205260036701ffffffffffffff60406000209660071c1693831b921b19161792906001600160401b0316600052602052604060002090565b9061587090606083526060830190613ae2565b8181036020830152825180825260208201916020808360051b8301019501926000915b8383106158e057505050505060408183039101526020808351928381520192019060005b8181106158c45750505090565b825163ffffffff168452602093840193909201916001016158b7565b90919293956020806158fe600193601f198682030187528a51610623565b98019301930191939290615893565b80516020909101516001600160e01b031981169291906004821061592f575050565b6001600160e01b031960049290920360031b82901b16169150565b90303b15610296576000916159736040519485938493630304c3e160e51b85526004850161585d565b038183305af19081615a14575b50615a095761598d612284565b9072c11c11c11c11c11c11c11c11c11c11c11c11c133146159af575b60039190565b6159c86159bb8361590d565b6001600160e01b03191690565b6337c3be2960e01b1480156159ee575b156159a957631d1ccf9f60e01b60005260046000fd5b506159fb6159bb8361590d565b632be8ca8b60e21b146159d8565b6002906103396105eb565b806125f06000615a239361020e565b38615980565b6040516370a0823160e01b60208201526001600160a01b039091166024820152919291615a8690615a5d81604481016140e0565b84837f000000000000000000000000000000000000000000000000000000000000000092615ab7565b9290911561509c5750805160208103615083575090615ab18260208061033995518301019101614e50565b93611f14565b939193615ac460846102cc565b94615ad2604051968761020e565b60848652615ae060846102cc565b602087019590601f1901368737833b15615b63575a90808210615b52578291038060061c90031115615b41576000918291825a9560208451940192f1905a9003923d9060848211615b38575b6000908287523e929190565b60849150615b2c565b6337c3be2960e01b60005260046000fd5b632be8ca8b60e21b60005260046000fd5b63030ed58f60e21b60005260046000fd5b8051928251908415615cd05761010185111580615cc4575b15615bf357818501946000198601956101008711615bf3578615615cb457615bb387611d43565b9660009586978795885b848110615c18575050505050600119018095149384615c0e575b505082615c04575b505015615bf357615bef91611d9d565b5190565b6309bde33960e01b60005260046000fd5b1490503880615bdf565b1492503880615bd7565b6001811b82811603615ca657868a1015615c9157615c3a60018b019a85611d9d565b51905b8c888c1015615c7d5750615c5560018c019b86611d9d565b515b818d11615bf357615c76828f92615c7090600196615ce1565b92611d9d565b5201615bbd565b60018d019c615c8b91611d9d565b51615c57565b615c9f60018c019b8d611d9d565b5190615c3d565b615c9f600189019884611d9d565b505050509050615bef9150611d8b565b50610101821115615b8c565b630469ac9960e21b60005260046000fd5b81811015615cf3579061033991615cf8565b610339915b9060405190602082019260018452604083015260608201526060815261354060808261020e56fea164736f6c634300081a000a49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b",
}

var OffRampWithMessageTransformerABI = OffRampWithMessageTransformerMetaData.ABI

var OffRampWithMessageTransformerBin = OffRampWithMessageTransformerMetaData.Bin

func DeployOffRampWithMessageTransformer(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OffRampStaticConfig, dynamicConfig OffRampDynamicConfig, sourceChainConfigs []OffRampSourceChainConfigArgs, messageTransformerAddr common.Address) (common.Address, *types.Transaction, *OffRampWithMessageTransformer, error) {
	parsed, err := OffRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OffRampWithMessageTransformerBin), backend, staticConfig, dynamicConfig, sourceChainConfigs, messageTransformerAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffRampWithMessageTransformer{address: address, abi: *parsed, OffRampWithMessageTransformerCaller: OffRampWithMessageTransformerCaller{contract: contract}, OffRampWithMessageTransformerTransactor: OffRampWithMessageTransformerTransactor{contract: contract}, OffRampWithMessageTransformerFilterer: OffRampWithMessageTransformerFilterer{contract: contract}}, nil
}

type OffRampWithMessageTransformer struct {
	address common.Address
	abi     abi.ABI
	OffRampWithMessageTransformerCaller
	OffRampWithMessageTransformerTransactor
	OffRampWithMessageTransformerFilterer
}

type OffRampWithMessageTransformerCaller struct {
	contract *bind.BoundContract
}

type OffRampWithMessageTransformerTransactor struct {
	contract *bind.BoundContract
}

type OffRampWithMessageTransformerFilterer struct {
	contract *bind.BoundContract
}

type OffRampWithMessageTransformerSession struct {
	Contract     *OffRampWithMessageTransformer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OffRampWithMessageTransformerCallerSession struct {
	Contract *OffRampWithMessageTransformerCaller
	CallOpts bind.CallOpts
}

type OffRampWithMessageTransformerTransactorSession struct {
	Contract     *OffRampWithMessageTransformerTransactor
	TransactOpts bind.TransactOpts
}

type OffRampWithMessageTransformerRaw struct {
	Contract *OffRampWithMessageTransformer
}

type OffRampWithMessageTransformerCallerRaw struct {
	Contract *OffRampWithMessageTransformerCaller
}

type OffRampWithMessageTransformerTransactorRaw struct {
	Contract *OffRampWithMessageTransformerTransactor
}

func NewOffRampWithMessageTransformer(address common.Address, backend bind.ContractBackend) (*OffRampWithMessageTransformer, error) {
	abi, err := abi.JSON(strings.NewReader(OffRampWithMessageTransformerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOffRampWithMessageTransformer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformer{address: address, abi: abi, OffRampWithMessageTransformerCaller: OffRampWithMessageTransformerCaller{contract: contract}, OffRampWithMessageTransformerTransactor: OffRampWithMessageTransformerTransactor{contract: contract}, OffRampWithMessageTransformerFilterer: OffRampWithMessageTransformerFilterer{contract: contract}}, nil
}

func NewOffRampWithMessageTransformerCaller(address common.Address, caller bind.ContractCaller) (*OffRampWithMessageTransformerCaller, error) {
	contract, err := bindOffRampWithMessageTransformer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerCaller{contract: contract}, nil
}

func NewOffRampWithMessageTransformerTransactor(address common.Address, transactor bind.ContractTransactor) (*OffRampWithMessageTransformerTransactor, error) {
	contract, err := bindOffRampWithMessageTransformer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerTransactor{contract: contract}, nil
}

func NewOffRampWithMessageTransformerFilterer(address common.Address, filterer bind.ContractFilterer) (*OffRampWithMessageTransformerFilterer, error) {
	contract, err := bindOffRampWithMessageTransformer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerFilterer{contract: contract}, nil
}

func bindOffRampWithMessageTransformer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OffRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRampWithMessageTransformer.Contract.OffRampWithMessageTransformerCaller.contract.Call(opts, result, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.OffRampWithMessageTransformerTransactor.contract.Transfer(opts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.OffRampWithMessageTransformerTransactor.contract.Transact(opts, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRampWithMessageTransformer.Contract.contract.Call(opts, result, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.contract.Transfer(opts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.contract.Transact(opts, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "ccipReceive", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRampWithMessageTransformer.Contract.CcipReceive(&_OffRampWithMessageTransformer.CallOpts, arg0)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRampWithMessageTransformer.Contract.CcipReceive(&_OffRampWithMessageTransformer.CallOpts, arg0)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getAllSourceChainConfigs")

	if err != nil {
		return *new([]uint64), *new([]OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)
	out1 := *abi.ConvertType(out[1], new([]OffRampSourceChainConfig)).(*[]OffRampSourceChainConfig)

	return out0, out1, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetAllSourceChainConfigs(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetAllSourceChainConfigs(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OffRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampDynamicConfig)).(*OffRampDynamicConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetDynamicConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetDynamicConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getExecutionState", sourceChainSelector, sequenceNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRampWithMessageTransformer.Contract.GetExecutionState(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRampWithMessageTransformer.Contract.GetExecutionState(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getLatestPriceSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRampWithMessageTransformer.Contract.GetLatestPriceSequenceNumber(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRampWithMessageTransformer.Contract.GetLatestPriceSequenceNumber(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getMerkleRoot", sourceChainSelector, root)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRampWithMessageTransformer.Contract.GetMerkleRoot(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, root)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRampWithMessageTransformer.Contract.GetMerkleRoot(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, root)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetMessageTransformer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getMessageTransformer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetMessageTransformer() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.GetMessageTransformer(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetMessageTransformer() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.GetMessageTransformer(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampSourceChainConfig)).(*OffRampSourceChainConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetSourceChainConfig(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetSourceChainConfig(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OffRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampStaticConfig)).(*OffRampStaticConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetStaticConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetStaticConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "latestConfigDetails", ocrPluginType)

	if err != nil {
		return *new(MultiOCR3BaseOCRConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOCRConfig)).(*MultiOCR3BaseOCRConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRampWithMessageTransformer.Contract.LatestConfigDetails(&_OffRampWithMessageTransformer.CallOpts, ocrPluginType)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRampWithMessageTransformer.Contract.LatestConfigDetails(&_OffRampWithMessageTransformer.CallOpts, ocrPluginType)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) Owner() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.Owner(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) Owner() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.Owner(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) TypeAndVersion() (string, error) {
	return _OffRampWithMessageTransformer.Contract.TypeAndVersion(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) TypeAndVersion() (string, error) {
	return _OffRampWithMessageTransformer.Contract.TypeAndVersion(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "acceptOwnership")
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.AcceptOwnership(&_OffRampWithMessageTransformer.TransactOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.AcceptOwnership(&_OffRampWithMessageTransformer.TransactOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "applySourceChainConfigUpdates", sourceChainConfigUpdates)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ApplySourceChainConfigUpdates(&_OffRampWithMessageTransformer.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ApplySourceChainConfigUpdates(&_OffRampWithMessageTransformer.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) Commit(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "commit", reportContext, report, rs, ss, rawVs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Commit(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Commit(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) Execute(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "execute", reportContext, report)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Execute(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Execute(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData, tokenGasOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ExecuteSingleMessage(&_OffRampWithMessageTransformer.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ExecuteSingleMessage(&_OffRampWithMessageTransformer.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "manuallyExecute", reports, gasLimitOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ManuallyExecute(&_OffRampWithMessageTransformer.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ManuallyExecute(&_OffRampWithMessageTransformer.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetDynamicConfig(&_OffRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetDynamicConfig(&_OffRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "setOCR3Configs", ocrConfigArgs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetOCR3Configs(&_OffRampWithMessageTransformer.TransactOpts, ocrConfigArgs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetOCR3Configs(&_OffRampWithMessageTransformer.TransactOpts, ocrConfigArgs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "transferOwnership", to)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.TransferOwnership(&_OffRampWithMessageTransformer.TransactOpts, to)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.TransferOwnership(&_OffRampWithMessageTransformer.TransactOpts, to)
}

type OffRampWithMessageTransformerAlreadyAttemptedIterator struct {
	Event *OffRampWithMessageTransformerAlreadyAttempted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerAlreadyAttemptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerAlreadyAttempted)
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
		it.Event = new(OffRampWithMessageTransformerAlreadyAttempted)
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

func (it *OffRampWithMessageTransformerAlreadyAttemptedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerAlreadyAttemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerAlreadyAttempted struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerAlreadyAttemptedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerAlreadyAttemptedIterator{contract: _OffRampWithMessageTransformer.contract, event: "AlreadyAttempted", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerAlreadyAttempted) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerAlreadyAttempted)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseAlreadyAttempted(log types.Log) (*OffRampWithMessageTransformerAlreadyAttempted, error) {
	event := new(OffRampWithMessageTransformerAlreadyAttempted)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerCommitReportAcceptedIterator struct {
	Event *OffRampWithMessageTransformerCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerCommitReportAccepted)
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
		it.Event = new(OffRampWithMessageTransformerCommitReportAccepted)
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

func (it *OffRampWithMessageTransformerCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerCommitReportAccepted struct {
	MerkleRoots  []InternalMerkleRoot
	PriceUpdates InternalPriceUpdates
	Raw          types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerCommitReportAcceptedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerCommitReportAcceptedIterator{contract: _OffRampWithMessageTransformer.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerCommitReportAccepted)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseCommitReportAccepted(log types.Log) (*OffRampWithMessageTransformerCommitReportAccepted, error) {
	event := new(OffRampWithMessageTransformerCommitReportAccepted)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerConfigSetIterator struct {
	Event *OffRampWithMessageTransformerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerConfigSet)
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

func (it *OffRampWithMessageTransformerConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerConfigSet struct {
	OcrPluginType uint8
	ConfigDigest  [32]byte
	Signers       []common.Address
	Transmitters  []common.Address
	F             uint8
	Raw           types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerConfigSetIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseConfigSet(log types.Log) (*OffRampWithMessageTransformerConfigSet, error) {
	event := new(OffRampWithMessageTransformerConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerDynamicConfigSetIterator struct {
	Event *OffRampWithMessageTransformerDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerDynamicConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerDynamicConfigSet)
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

func (it *OffRampWithMessageTransformerDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerDynamicConfigSet struct {
	DynamicConfig OffRampDynamicConfig
	Raw           types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerDynamicConfigSetIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerDynamicConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerDynamicConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerDynamicConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseDynamicConfigSet(log types.Log) (*OffRampWithMessageTransformerDynamicConfigSet, error) {
	event := new(OffRampWithMessageTransformerDynamicConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerExecutionStateChangedIterator struct {
	Event *OffRampWithMessageTransformerExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerExecutionStateChanged)
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
		it.Event = new(OffRampWithMessageTransformerExecutionStateChanged)
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

func (it *OffRampWithMessageTransformerExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	MessageHash         [32]byte
	State               uint8
	ReturnData          []byte
	GasUsed             *big.Int
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampWithMessageTransformerExecutionStateChangedIterator, error) {

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

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerExecutionStateChangedIterator{contract: _OffRampWithMessageTransformer.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerExecutionStateChanged)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseExecutionStateChanged(log types.Log) (*OffRampWithMessageTransformerExecutionStateChanged, error) {
	event := new(OffRampWithMessageTransformerExecutionStateChanged)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerOwnershipTransferRequestedIterator struct {
	Event *OffRampWithMessageTransformerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerOwnershipTransferRequested)
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
		it.Event = new(OffRampWithMessageTransformerOwnershipTransferRequested)
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

func (it *OffRampWithMessageTransformerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerOwnershipTransferRequestedIterator{contract: _OffRampWithMessageTransformer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerOwnershipTransferRequested)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferRequested, error) {
	event := new(OffRampWithMessageTransformerOwnershipTransferRequested)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerOwnershipTransferredIterator struct {
	Event *OffRampWithMessageTransformerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerOwnershipTransferred)
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
		it.Event = new(OffRampWithMessageTransformerOwnershipTransferred)
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

func (it *OffRampWithMessageTransformerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerOwnershipTransferredIterator{contract: _OffRampWithMessageTransformer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerOwnershipTransferred)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseOwnershipTransferred(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferred, error) {
	event := new(OffRampWithMessageTransformerOwnershipTransferred)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerRootRemovedIterator struct {
	Event *OffRampWithMessageTransformerRootRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerRootRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerRootRemoved)
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
		it.Event = new(OffRampWithMessageTransformerRootRemoved)
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

func (it *OffRampWithMessageTransformerRootRemovedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerRootRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerRootRemoved struct {
	Root [32]byte
	Raw  types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterRootRemoved(opts *bind.FilterOpts) (*OffRampWithMessageTransformerRootRemovedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerRootRemovedIterator{contract: _OffRampWithMessageTransformer.contract, event: "RootRemoved", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerRootRemoved) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerRootRemoved)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "RootRemoved", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseRootRemoved(log types.Log) (*OffRampWithMessageTransformerRootRemoved, error) {
	event := new(OffRampWithMessageTransformerRootRemoved)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "RootRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator struct {
	Event *OffRampWithMessageTransformerSkippedAlreadyExecutedMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
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
		it.Event = new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
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

func (it *OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSkippedAlreadyExecutedMessage struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator{contract: _OffRampWithMessageTransformer.contract, event: "SkippedAlreadyExecutedMessage", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedAlreadyExecutedMessage) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessage, error) {
	event := new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSkippedReportExecutionIterator struct {
	Event *OffRampWithMessageTransformerSkippedReportExecution

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSkippedReportExecutionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSkippedReportExecution)
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
		it.Event = new(OffRampWithMessageTransformerSkippedReportExecution)
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

func (it *OffRampWithMessageTransformerSkippedReportExecutionIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSkippedReportExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSkippedReportExecution struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedReportExecutionIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSkippedReportExecutionIterator{contract: _OffRampWithMessageTransformer.contract, event: "SkippedReportExecution", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedReportExecution) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSkippedReportExecution)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSkippedReportExecution(log types.Log) (*OffRampWithMessageTransformerSkippedReportExecution, error) {
	event := new(OffRampWithMessageTransformerSkippedReportExecution)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSourceChainConfigSetIterator struct {
	Event *OffRampWithMessageTransformerSourceChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSourceChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSourceChainConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerSourceChainConfigSet)
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

func (it *OffRampWithMessageTransformerSourceChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSourceChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceConfig        OffRampSourceChainConfig
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampWithMessageTransformerSourceChainConfigSetIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSourceChainConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "SourceChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSourceChainConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSourceChainConfigSet(log types.Log) (*OffRampWithMessageTransformerSourceChainConfigSet, error) {
	event := new(OffRampWithMessageTransformerSourceChainConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSourceChainSelectorAddedIterator struct {
	Event *OffRampWithMessageTransformerSourceChainSelectorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSourceChainSelectorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSourceChainSelectorAdded)
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
		it.Event = new(OffRampWithMessageTransformerSourceChainSelectorAdded)
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

func (it *OffRampWithMessageTransformerSourceChainSelectorAddedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSourceChainSelectorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSourceChainSelectorAdded struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSourceChainSelectorAddedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSourceChainSelectorAddedIterator{contract: _OffRampWithMessageTransformer.contract, event: "SourceChainSelectorAdded", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainSelectorAdded) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSourceChainSelectorAdded)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSourceChainSelectorAdded(log types.Log) (*OffRampWithMessageTransformerSourceChainSelectorAdded, error) {
	event := new(OffRampWithMessageTransformerSourceChainSelectorAdded)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerStaticConfigSetIterator struct {
	Event *OffRampWithMessageTransformerStaticConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerStaticConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerStaticConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerStaticConfigSet)
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

func (it *OffRampWithMessageTransformerStaticConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerStaticConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerStaticConfigSet struct {
	StaticConfig OffRampStaticConfig
	Raw          types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerStaticConfigSetIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerStaticConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "StaticConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerStaticConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerStaticConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseStaticConfigSet(log types.Log) (*OffRampWithMessageTransformerStaticConfigSet, error) {
	event := new(OffRampWithMessageTransformerStaticConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerTransmittedIterator struct {
	Event *OffRampWithMessageTransformerTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerTransmitted)
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
		it.Event = new(OffRampWithMessageTransformerTransmitted)
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

func (it *OffRampWithMessageTransformerTransmittedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerTransmitted struct {
	OcrPluginType  uint8
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampWithMessageTransformerTransmittedIterator, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerTransmittedIterator{contract: _OffRampWithMessageTransformer.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerTransmitted, ocrPluginType []uint8) (event.Subscription, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerTransmitted)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseTransmitted(log types.Log) (*OffRampWithMessageTransformerTransmitted, error) {
	event := new(OffRampWithMessageTransformerTransmitted)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OffRampWithMessageTransformer.abi.Events["AlreadyAttempted"].ID:
		return _OffRampWithMessageTransformer.ParseAlreadyAttempted(log)
	case _OffRampWithMessageTransformer.abi.Events["CommitReportAccepted"].ID:
		return _OffRampWithMessageTransformer.ParseCommitReportAccepted(log)
	case _OffRampWithMessageTransformer.abi.Events["ConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["DynamicConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseDynamicConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["ExecutionStateChanged"].ID:
		return _OffRampWithMessageTransformer.ParseExecutionStateChanged(log)
	case _OffRampWithMessageTransformer.abi.Events["OwnershipTransferRequested"].ID:
		return _OffRampWithMessageTransformer.ParseOwnershipTransferRequested(log)
	case _OffRampWithMessageTransformer.abi.Events["OwnershipTransferred"].ID:
		return _OffRampWithMessageTransformer.ParseOwnershipTransferred(log)
	case _OffRampWithMessageTransformer.abi.Events["RootRemoved"].ID:
		return _OffRampWithMessageTransformer.ParseRootRemoved(log)
	case _OffRampWithMessageTransformer.abi.Events["SkippedAlreadyExecutedMessage"].ID:
		return _OffRampWithMessageTransformer.ParseSkippedAlreadyExecutedMessage(log)
	case _OffRampWithMessageTransformer.abi.Events["SkippedReportExecution"].ID:
		return _OffRampWithMessageTransformer.ParseSkippedReportExecution(log)
	case _OffRampWithMessageTransformer.abi.Events["SourceChainConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseSourceChainConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["SourceChainSelectorAdded"].ID:
		return _OffRampWithMessageTransformer.ParseSourceChainSelectorAdded(log)
	case _OffRampWithMessageTransformer.abi.Events["StaticConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseStaticConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["Transmitted"].ID:
		return _OffRampWithMessageTransformer.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OffRampWithMessageTransformerAlreadyAttempted) Topic() common.Hash {
	return common.HexToHash("0x3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120")
}

func (OffRampWithMessageTransformerCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4")
}

func (OffRampWithMessageTransformerConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (OffRampWithMessageTransformerDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0xcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee")
}

func (OffRampWithMessageTransformerExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b")
}

func (OffRampWithMessageTransformerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OffRampWithMessageTransformerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OffRampWithMessageTransformerRootRemoved) Topic() common.Hash {
	return common.HexToHash("0x202f1139a3e334b6056064c0e9b19fd07e44a88d8f6e5ded571b24cf8c371f12")
}

func (OffRampWithMessageTransformerSkippedAlreadyExecutedMessage) Topic() common.Hash {
	return common.HexToHash("0x3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c")
}

func (OffRampWithMessageTransformerSkippedReportExecution) Topic() common.Hash {
	return common.HexToHash("0xaab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d894933")
}

func (OffRampWithMessageTransformerSourceChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b")
}

func (OffRampWithMessageTransformerSourceChainSelectorAdded) Topic() common.Hash {
	return common.HexToHash("0xf4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9")
}

func (OffRampWithMessageTransformerStaticConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b7")
}

func (OffRampWithMessageTransformerTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformer) Address() common.Address {
	return _OffRampWithMessageTransformer.address
}

type OffRampWithMessageTransformerInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error)

	GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error)

	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error)

	GetMessageTransformer(opts *bind.CallOpts) (common.Address, error)

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

	FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerAlreadyAttemptedIterator, error)

	WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerAlreadyAttempted) (event.Subscription, error)

	ParseAlreadyAttempted(log types.Log) (*OffRampWithMessageTransformerAlreadyAttempted, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*OffRampWithMessageTransformerCommitReportAccepted, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OffRampWithMessageTransformerConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerDynamicConfigSet) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*OffRampWithMessageTransformerDynamicConfigSet, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampWithMessageTransformerExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*OffRampWithMessageTransformerExecutionStateChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferred, error)

	FilterRootRemoved(opts *bind.FilterOpts) (*OffRampWithMessageTransformerRootRemovedIterator, error)

	WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerRootRemoved) (event.Subscription, error)

	ParseRootRemoved(log types.Log) (*OffRampWithMessageTransformerRootRemoved, error)

	FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator, error)

	WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedAlreadyExecutedMessage) (event.Subscription, error)

	ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessage, error)

	FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedReportExecutionIterator, error)

	WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedReportExecution) (event.Subscription, error)

	ParseSkippedReportExecution(log types.Log) (*OffRampWithMessageTransformerSkippedReportExecution, error)

	FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampWithMessageTransformerSourceChainConfigSetIterator, error)

	WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error)

	ParseSourceChainConfigSet(log types.Log) (*OffRampWithMessageTransformerSourceChainConfigSet, error)

	FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSourceChainSelectorAddedIterator, error)

	WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainSelectorAdded) (event.Subscription, error)

	ParseSourceChainSelectorAdded(log types.Log) (*OffRampWithMessageTransformerSourceChainSelectorAdded, error)

	FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerStaticConfigSetIterator, error)

	WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerStaticConfigSet) (event.Subscription, error)

	ParseStaticConfigSet(log types.Log) (*OffRampWithMessageTransformerStaticConfigSet, error)

	FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampWithMessageTransformerTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerTransmitted, ocrPluginType []uint8) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*OffRampWithMessageTransformerTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
