// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package offramp

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

var OffRampMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applySourceChainConfigUpdates\",\"inputs\":[{\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"commit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeSingleMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllSourceChainConfigs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExecutionState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestPriceSequenceNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMerkleRoot\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSourceChainConfig\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"ocrConfig\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"components\":[{\"name\":\"configInfo\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"n\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyExecute\",\"inputs\":[{\"name\":\"reports\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\",\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"components\":[{\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOCR3Configs\",\"inputs\":[{\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AlreadyAttempted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitReportAccepted\",\"inputs\":[{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"F\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DynamicConfigSet\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecutionStateChanged\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"state\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RootRemoved\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedAlreadyExecutedMessage\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedReportExecution\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainConfigSet\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sourceConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainSelectorAdded\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaticConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CanOnlySelfCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CommitOnRampMismatch\",\"inputs\":[{\"name\":\"reportOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"configOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyReport\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ExecutionError\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ForkedChain\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[{\"name\":\"errorType\",\"type\":\"uint8\",\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\"}]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInterval\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionGasLimit\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"newLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionTokenGasOverride\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oldLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverride\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidMessageDestChainSelector\",\"inputs\":[{\"name\":\"messageDestChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidNewState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newState\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}]},{\"type\":\"error\",\"name\":\"InvalidOnRampUpdate\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LeavesCannotBeEmpty\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionGasAmountCountMismatch\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ManualExecutionGasLimitMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionNotYetEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageValidationError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonUniqueSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotACompatiblePool\",\"inputs\":[{\"name\":\"notPool\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OracleCannotBeZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReceiverError\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ReleaseOrMintBalanceMismatch\",\"inputs\":[{\"name\":\"amountReleased\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePre\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RootAlreadyCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"RootNotCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignaturesOutOfRegistration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SourceChainNotEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SourceChainSelectorMismatch\",\"inputs\":[{\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StaleCommitReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaticConfigCannotBeChanged\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TokenDataMismatch\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TokenHandlingError\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnexpectedTokenData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongMessageLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WrongNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroChainSelectorNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461084857616351803803809161001d828561087e565b833981019080820361014081126108485760a08112610848576040519060a082016001600160401b0381118382101761084d5760405261005c836108a1565b825260208301519261ffff84168403610848576020830193845260408101516001600160a01b0381168103610848576040840190815261009e606083016108b5565b946060850195865260806100b38185016108b5565b86820190815294609f19011261084857604051946100d086610863565b6100dc60a085016108b5565b865260c08401519463ffffffff86168603610848576020870195865261010460e086016108c9565b976040880198895261011961010087016108b5565b6060890190815261012087015190966001600160401b03821161084857018a601f820112156108485780519a6001600160401b038c1161084d578b60051b916020806040519e8f9061016d8388018361087e565b81520193820101908282116108485760208101935b828510610748575050505050331561073757600180546001600160a01b031916331790554660805284516001600160a01b0316158015610725575b8015610713575b6106f15782516001600160401b0316156107025782516001600160401b0390811660a090815286516001600160a01b0390811660c0528351811660e0528451811661010052865161ffff90811661012052604080519751909416875296519096166020860152955185169084015251831660608301525190911660808201527fb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b79190a182516001600160a01b0316156106f1579151600480548351865160ff60c01b90151560c01b1663ffffffff60a01b60a09290921b919091166001600160a01b039485166001600160c81b0319909316831717179091558351600580549184166001600160a01b031990921691909117905560408051918252925163ffffffff166020820152935115159184019190915290511660608201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee90608090a160005b815181101561066b576020600582901b8301810151908101516001600160401b031690600090821561065c5780516001600160a01b03161561064d57828252600860205260408220906060810151600183019361038585546108d6565b6105ee578354600160a81b600160e81b031916600160a81b1784556040518681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb990602090a15b815180159081156105c3575b506105b4578151916001600160401b0383116105a0576103f986546108d6565b601f811161055b575b50602091601f84116001146104e257926001989796949281926000805160206163318339815191529795926104d7575b5050600019600383901b1c191690881b1783555b60408101518254915160a089811b8a9003801960ff60a01b1990951693151590911b60ff60a01b169290921792909216911617815561048484610993565b506104ce6040519283926020845254888060a01b038116602085015260ff8160a01c1615156040850152888060401b039060a81c16606084015260808084015260a0830190610910565b0390a201610328565b015190503880610432565b9190601f198416878452828420935b8181106105435750926001999897959392859260008051602061633183398151915298968c951061052a575b505050811b018355610446565b015160001960f88460031b161c1916905538808061051d565b929360206001819287860151815501950193016104f1565b86835260208320601f850160051c81019160208610610596575b601f0160051c01905b81811061058b5750610402565b83815560010161057e565b9091508190610575565b634e487b7160e01b82526041600452602482fd5b6342bcdf7f60e11b8152600490fd5b905060208301206040516020810190838252602081526105e460408261087e565b51902014386103d9565b835460a81c6001600160401b0316600114158061061f575b156103cd57632105803760e11b81526004869052602490fd5b50604051610638816106318189610910565b038261087e565b60208151910120825160208401201415610606565b6342bcdf7f60e11b8252600482fd5b63c656089560e01b8252600482fd5b60405161590a9081610a278239608051816136a4015260a0518181816104700152614203015260c0518181816104c601528181612cd201528181613126015261419d015260e0518181816104f501526149e001526101005181818161052401526145c60152610120518181816104970152818161244c01528181614ad3015261563f0152f35b6342bcdf7f60e11b60005260046000fd5b63c656089560e01b60005260046000fd5b5081516001600160a01b0316156101c4565b5080516001600160a01b0316156101bd565b639b15e16f60e01b60005260046000fd5b84516001600160401b0381116108485782016080818603601f190112610848576040519061077582610863565b60208101516001600160a01b0381168103610848578252610798604082016108a1565b60208301526107a9606082016108c9565b604083015260808101516001600160401b03811161084857602091010185601f820112156108485780516001600160401b03811161084d57604051916107f9601f8301601f19166020018461087e565b81835287602083830101116108485760005b8281106108335750509181600060208096949581960101526060820152815201940193610182565b8060208092840101518282870101520161080b565b600080fd5b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761084d57604052565b601f909101601f19168101906001600160401b0382119082101761084d57604052565b51906001600160401b038216820361084857565b51906001600160a01b038216820361084857565b5190811515820361084857565b90600182811c92168015610906575b60208310146108f057565b634e487b7160e01b600052602260045260246000fd5b91607f16916108e5565b60009291815491610920836108d6565b8083529260018116908115610976575060011461093c57505050565b60009081526020812093945091925b83831061095c575060209250010190565b60018160209294939454838587010152019101919061094b565b915050602093945060ff929192191683830152151560051b010190565b80600052600760205260406000205415600014610a20576006546801000000000000000081101561084d576001810180600655811015610a0a577ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0181905560065460009182526007602052604090912055600190565b634e487b7160e01b600052603260045260246000fd5b5060009056fe6080604052600436101561001257600080fd5b60003560e01c806304666f9c1461015757806306285c6914610152578063181f5a771461014d5780633f4b04aa146101485780635215505b146101435780635e36480c1461013e5780635e7bb0081461013957806360987c20146101345780637437ff9f1461012f57806379ba50971461012a5780637edf52f41461012557806385572ffb146101205780638da5cb5b1461011b578063c673e58414610116578063ccd37ba314610111578063de5e0b9a1461010c578063e9d68a8e14610107578063f2fde38b14610102578063f58e03fc146100fd5763f716f99f146100f857600080fd5b6118b6565b611799565b61170e565b611673565b6115d7565b611553565b6114a8565b6113c0565b61138a565b6111c4565b611144565b61109b565b611020565b610e1b565b6108b0565b61076b565b61065e565b6105ff565b61043d565b61031f565b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761018d57604052565b61015c565b60a081019081106001600160401b0382111761018d57604052565b604081019081106001600160401b0382111761018d57604052565b606081019081106001600160401b0382111761018d57604052565b90601f801991011681019081106001600160401b0382111761018d57604052565b6040519061021360c0836101e3565b565b6040519061021360a0836101e3565b60405190610213610100836101e3565b604051906102136040836101e3565b6001600160401b03811161018d5760051b60200190565b6001600160a01b0381160361026b57565b600080fd5b600435906001600160401b038216820361026b57565b35906001600160401b038216820361026b57565b8015150361026b57565b35906102138261029a565b6001600160401b03811161018d57601f01601f191660200190565b9291926102d6826102af565b916102e460405193846101e3565b82948184528183011161026b578281602093846000960137010152565b9080601f8301121561026b5781602061031c933591016102ca565b90565b3461026b57602036600319011261026b576004356001600160401b03811161026b573660238201121561026b5780600401359061035b82610243565b9061036960405192836101e3565b8282526024602083019360051b8201019036821161026b5760248101935b82851061039957610397846119f1565b005b84356001600160401b03811161026b5782016080602319823603011261026b57604051916103c683610172565b60248201356103d48161025a565b83526103e260448301610286565b602084015260648201356103f58161029a565b60408401526084820135926001600160401b03841161026b57610422602094936024869536920101610301565b6060820152815201940193610387565b600091031261026b57565b3461026b57600036600319011261026b57610456611c9d565b5061059e60405161046681610192565b6001600160401b037f000000000000000000000000000000000000000000000000000000000000000016815261ffff7f00000000000000000000000000000000000000000000000000000000000000001660208201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660408201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660608201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660808201526040519182918291909160806001600160a01b038160a08401956001600160401b03815116855261ffff6020820151166020860152826040820151166040860152826060820151166060860152015116910152565b0390f35b604051906105b16020836101e3565b60008252565b60005b8381106105ca5750506000910152565b81810151838201526020016105ba565b906020916105f3815180928185528580860191016105b7565b601f01601f1916010190565b3461026b57600036600319011261026b5761059e604080519061062281836101e3565b601182527f4f666652616d7020312e362e302d6465760000000000000000000000000000006020830152519182916020835260208301906105da565b3461026b57600036600319011261026b5760206001600160401b03600b5416604051908152f35b906080606061031c936001600160a01b0381511684526020810151151560208501526001600160401b03604082015116604085015201519181606082015201906105da565b6040810160408252825180915260206060830193019060005b81811061074c575050506020818303910152815180825260208201916020808360051b8301019401926000915b83831061071f57505050505090565b909192939460208061073d600193601f198682030187528951610685565b97019301930191939290610710565b82516001600160401b03168552602094850194909201916001016106e3565b3461026b57600036600319011261026b5760065461078881610243565b9061079660405192836101e3565b808252601f196107a582610243565b0160005b8181106108675750506107bb81611cef565b9060005b8181106107d757505061059e604051928392836106ca565b8061080d6107f56107e9600194614084565b6001600160401b031690565b6107ff8387611d49565b906001600160401b03169052565b61084b61084661082d6108208488611d49565b516001600160401b031690565b6001600160401b03166000526008602052604060002090565b611e35565b6108558287611d49565b526108608186611d49565b50016107bf565b602090610872611cc8565b828287010152016107a9565b634e487b7160e01b600052602160045260246000fd5b6004111561089e57565b61087e565b90600482101561089e5752565b3461026b57604036600319011261026b576108c9610270565b602435906001600160401b038216820361026b576020916108e991611ed1565b6108f660405180926108a3565bf35b91908260a091031261026b5760405161091081610192565b60806109558183958035855261092860208201610286565b602086015261093960408201610286565b604086015261094a60608201610286565b606086015201610286565b910152565b35906102138261025a565b63ffffffff81160361026b57565b359061021382610965565b81601f8201121561026b5780359061099582610243565b926109a360405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b8385106109cf57505050505090565b84356001600160401b03811161026b57820160a0818503601f19011261026b57604051916109fc83610192565b60208201356001600160401b03811161026b57856020610a1e92850101610301565b83526040820135610a2e8161025a565b6020840152610a3f60608301610973565b60408401526080820135926001600160401b03841161026b5760a083610a6c886020809881980101610301565b6060840152013560808201528152019401936109c0565b9190916101408184031261026b57610a99610204565b92610aa481836108f8565b845260a08201356001600160401b03811161026b5781610ac5918401610301565b602085015260c08201356001600160401b03811161026b5781610ae9918401610301565b6040850152610afa60e0830161095a565b606085015261010082013560808501526101208201356001600160401b03811161026b57610b28920161097e565b60a0830152565b9080601f8301121561026b578135610b4681610243565b92610b5460405194856101e3565b81845260208085019260051b8201019183831161026b5760208201905b838210610b8057505050505090565b81356001600160401b03811161026b57602091610ba287848094880101610a83565b815201910190610b71565b81601f8201121561026b57803590610bc482610243565b92610bd260405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b838510610bfe57505050505090565b84356001600160401b03811161026b57820183603f8201121561026b576020810135610c2981610243565b91610c3760405193846101e3565b8183526020808085019360051b830101019186831161026b5760408201905b838210610c70575050509082525060209485019401610bef565b81356001600160401b03811161026b57602091610c948a8480809589010101610301565b815201910190610c56565b929190610cab81610243565b93610cb960405195866101e3565b602085838152019160051b810192831161026b57905b828210610cdb57505050565b8135815260209182019101610ccf565b9080601f8301121561026b5781602061031c93359101610c9f565b81601f8201121561026b57803590610d1d82610243565b92610d2b60405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b838510610d5757505050505090565b84356001600160401b03811161026b57820160a0818503601f19011261026b57610d7f610215565b91610d8c60208301610286565b835260408201356001600160401b03811161026b57856020610db092850101610b2f565b602084015260608201356001600160401b03811161026b57856020610dd792850101610bad565b60408401526080820135926001600160401b03841161026b5760a083610e04886020809881980101610ceb565b606084015201356080820152815201940193610d48565b3461026b57604036600319011261026b576004356001600160401b03811161026b57610e4b903690600401610d06565b6024356001600160401b03811161026b573660238201121561026b57806004013591610e7683610243565b91610e8460405193846101e3565b8383526024602084019460051b8201019036821161026b5760248101945b828610610eb3576103978585611f19565b85356001600160401b03811161026b5782013660438201121561026b576024810135610ede81610243565b91610eec60405193846101e3565b818352602060248185019360051b830101019036821161026b5760448101925b828410610f26575050509082525060209586019501610ea2565b83356001600160401b03811161026b576024908301016040601f19823603011261026b5760405190610f57826101ad565b6020810135825260408101356001600160401b03811161026b57602091010136601f8201121561026b57803590610f8d82610243565b91610f9b60405193846101e3565b80835260208084019160051b8301019136831161026b57602001905b828210610fd65750505091816020938480940152815201930192610f0c565b602080918335610fe581610965565b815201910190610fb7565b9181601f8401121561026b578235916001600160401b03831161026b576020808501948460051b01011161026b57565b3461026b57606036600319011261026b576004356001600160401b03811161026b57611050903690600401610a83565b6024356001600160401b03811161026b5761106f903690600401610ff0565b91604435926001600160401b03841161026b57611093610397943690600401610ff0565b939092612325565b3461026b57600036600319011261026b576110b46125fd565b5061059e6040516110c481610172565b60ff6004546001600160a01b038116835263ffffffff8160a01c16602084015260c01c16151560408201526001600160a01b036005541660608201526040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b3461026b57600036600319011261026b576000546001600160a01b03811633036111b3576001600160a01b0319600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b63015aa1e360e11b60005260046000fd5b3461026b57608036600319011261026b5760006040516111e381610172565b6004356111ef8161025a565b81526024356111fd81610965565b602082015260443561120e8161029a565b604082015260643561121f8161025a565b606082015261122c6134a1565b6001600160a01b038151161561137b576113758161128b6001600160a01b037fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9451166001600160a01b03166001600160a01b03196004541617600455565b60208101516004547fffffffffffffff0000000000ffffffffffffffffffffffffffffffffffffffff77ffffffff000000000000000000000000000000000000000078ff0000000000000000000000000000000000000000000000006040860151151560c01b169360a01b169116171760045561133161131560608301516001600160a01b031690565b6001600160a01b03166001600160a01b03196005541617600555565b6040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b0390a180f35b6342bcdf7f60e11b8252600482fd5b3461026b57602036600319011261026b576004356001600160401b03811161026b5760a090600319903603011261026b57600080fd5b3461026b57600036600319011261026b5760206001600160a01b0360015416604051908152f35b6004359060ff8216820361026b57565b359060ff8216820361026b57565b906020808351928381520192019060005b8181106114235750505090565b82516001600160a01b0316845260209384019390920191600101611416565b9061031c9160208152606082518051602084015260ff602082015116604084015260ff604082015116828401520151151560808201526040611493602084015160c060a085015260e0840190611405565b9201519060c0601f1982850301910152611405565b3461026b57602036600319011261026b5760ff6114c36113e7565b6060604080516114d2816101c8565b6114da6125fd565b8152826020820152015216600052600260205261059e6040600020600361154260405192611507846101c8565b61151081612622565b845260405161152d81611526816002860161265b565b03826101e3565b6020850152611526604051809481930161265b565b604082015260405191829182611442565b3461026b57604036600319011261026b5761156c610270565b6001600160401b036024359116600052600a6020526040600020906000526020526020604060002054604051908152f35b9060049160441161026b57565b9181601f8401121561026b578235916001600160401b03831161026b576020838186019501011161026b57565b3461026b5760c036600319011261026b576115f13661159d565b6044356001600160401b03811161026b576116109036906004016115aa565b6064929192356001600160401b03811161026b57611632903690600401610ff0565b60843594916001600160401b03861161026b57611656610397963690600401610ff0565b94909360a43596612c8d565b90602061031c928181520190610685565b3461026b57602036600319011261026b576001600160401b03611694610270565b61169c611cc8565b5016600052600860205261059e604060002060016116fd604051926116c084610172565b6001600160401b0381546001600160a01b038116865260ff8160a01c161515602087015260a81c1660408501526115266040518094819301611d97565b606082015260405191829182611662565b3461026b57602036600319011261026b576001600160a01b036004356117338161025a565b61173b6134a1565b1633811461178857806001600160a01b031960005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b636d6c4ee560e11b60005260046000fd5b3461026b57606036600319011261026b576117b33661159d565b6044356001600160401b03811161026b576117d29036906004016115aa565b9182820160208382031261026b578235906001600160401b03821161026b576117fc918401610d06565b60405190602061180c81846101e3565b60008352601f19810160005b818110611840575050506103979491611830916136e5565b61183861319c565b928392613a95565b60608582018401528201611818565b9080601f8301121561026b57813561186681610243565b9261187460405194856101e3565b81845260208085019260051b82010192831161026b57602001905b82821061189c5750505090565b6020809183356118ab8161025a565b81520191019061188f565b3461026b57602036600319011261026b576004356001600160401b03811161026b573660238201121561026b578060040135906118f282610243565b9061190060405192836101e3565b8282526024602083019360051b8201019036821161026b5760248101935b82851061192e57610397846131b8565b84356001600160401b03811161026b57820160c0602319823603011261026b57611956610204565b916024820135835261196a604483016113f7565b602084015261197b606483016113f7565b604084015261198c608483016102a4565b606084015260a48201356001600160401b03811161026b576119b4906024369185010161184f565b608084015260c4820135926001600160401b03841161026b576119e160209493602486953692010161184f565b60a082015281520194019361191e565b6119f96134a1565b60005b8151811015611c9957611a0f8183611d49565b5190611a2560208301516001600160401b031690565b916001600160401b038316908115611c8857611a5a611a4e611a4e83516001600160a01b031690565b6001600160a01b031690565b15611bef57611a7c846001600160401b03166000526008602052604060002090565b906060810151916001810195611a928754611d5d565b611c1657611b057ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb991611aeb84750100000000000000000000000000000000000000000067ffffffffffffffff60a81b19825416179055565b6040516001600160401b0390911681529081906020820190565b0390a15b82518015908115611c00575b50611bef57611bd0611bb4611be693611b517f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b9660019a613543565b611ba7611b616040830151151590565b85547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690151560a01b74ff000000000000000000000000000000000000000016178555565b516001600160a01b031690565b82906001600160a01b03166001600160a01b0319825416179055565b611bd9846150ac565b5060405191829182613614565b0390a2016119fc565b6342bcdf7f60e11b60005260046000fd5b90506020840120611c0f6134c6565b1438611b15565b60016001600160401b03611c3584546001600160401b039060a81c1690565b16141580611c69575b611c485750611b09565b632105803760e11b6000526001600160401b031660045260246000fd5b6000fd5b50611c7387611e1a565b60208151910120845160208601201415611c3e565b63c656089560e01b60005260046000fd5b5050565b60405190611caa82610192565b60006080838281528260208201528260408201528260608201520152565b60405190611cd582610172565b606080836000815260006020820152600060408201520152565b90611cf982610243565b611d0660405191826101e3565b8281528092611d17601f1991610243565b0190602036910137565b634e487b7160e01b600052603260045260246000fd5b805115611d445760200190565b611d21565b8051821015611d445760209160051b010190565b90600182811c92168015611d8d575b6020831014611d7757565b634e487b7160e01b600052602260045260246000fd5b91607f1691611d6c565b60009291815491611da783611d5d565b8083529260018116908115611dfd5750600114611dc357505050565b60009081526020812093945091925b838310611de3575060209250010190565b600181602092949394548385870101520191019190611dd2565b915050602093945060ff929192191683830152151560051b010190565b90610213611e2e9260405193848092611d97565b03836101e3565b9060016060604051611e4681610172565b611e8f81956001600160401b0381546001600160a01b038116855260ff8160a01c161515602086015260a81c166040840152611e886040518096819301611d97565b03846101e3565b0152565b634e487b7160e01b600052601160045260246000fd5b908160051b9180830460201490151715611ebf57565b611e93565b91908203918211611ebf57565b611edd82607f9261365e565b9116906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611ebf576003911c16600481101561089e5790565b611f216136a2565b8051825181036121185760005b818110611f4157505090610213916136e5565b611f4b8184611d49565b516020810190815151611f5e8488611d49565b5192835182036121185790916000925b808410611f82575050505050600101611f2e565b91949398611f94848b98939598611d49565b515198611fa2888851611d49565b5199806120cf575b5060a08a01988b6020611fc08b8d515193611d49565b51015151036120925760005b8a515181101561207d57612008611fff611ff58f6020611fed8f8793611d49565b510151611d49565b5163ffffffff1690565b63ffffffff1690565b8b81612019575b5050600101611fcc565b611fff604061202c856120389451611d49565b51015163ffffffff1690565b9081811061204757508b61200f565b8d51516040516348e617b360e01b81526004810191909152602481019390935260448301919091526064820152608490fd5b0390fd5b50985098509893949095600101929091611f6e565b611c658b516120ad606082519201516001600160401b031690565b6370a193fd60e01b6000526004919091526001600160401b0316602452604490565b60808b0151811015611faa57611c65908b6120f188516001600160401b031690565b905151633a98d46360e11b6000526001600160401b03909116600452602452604452606490565b6320f8fd5960e21b60005260046000fd5b60405190612136826101ad565b60006020838281520152565b604051906121516020836101e3565b600080835282815b82811061216557505050565b602090612170612129565b82828501015201612159565b805182526001600160401b03602082015116602083015260806121c36121b1604084015160a0604087015260a08601906105da565b606084015185820360608701526105da565b9101519160808183039101526020808351928381520192019060005b8181106121ec5750505090565b825180516001600160a01b0316855260209081015181860152604090940193909201916001016121df565b90602061031c92818152019061217c565b6040513d6000823e3d90fd5b3d1561225f573d90612245826102af565b9161225360405193846101e3565b82523d6000602084013e565b606090565b90602061031c9281815201906105da565b909160608284031261026b57815161228c8161029a565b9260208301516001600160401b03811161026b5783019080601f8301121561026b578151916122ba836102af565b916122c860405193846101e3565b8383526020848301011161026b576040926122e991602080850191016105b7565b92015190565b9293606092959461ffff6123136001600160a01b039460808852608088019061217c565b97166020860152604085015216910152565b93919290933033036125ec5761233a90613798565b92612343612142565b9460a085015180516125a5575b505050505080519161236e602084519401516001600160401b031690565b90602083015191604084019261239b845192612388610215565b9788526001600160401b03166020880152565b6040860152606085015260808401526001600160a01b036123c46005546001600160a01b031690565b1680612528575b505151158061251c575b8015612506575b80156124dd575b611c9957612475918161241a611a4e61240d61082d602060009751016001600160401b0390511690565b546001600160a01b031690565b9083612435606060808401519301516001600160a01b031690565b604051633cf9798360e01b815296879586948593917f000000000000000000000000000000000000000000000000000000000000000090600486016122ef565b03925af19081156124d8576000906000926124b1575b50156124945750565b6040516302a35ba360e21b81529081906120799060048301612264565b90506124d091503d806000833e6124c881836101e3565b810190612275565b50903861248b565b612228565b506125016124fd6124f860608401516001600160a01b031690565b613956565b1590565b6123e3565b5060608101516001600160a01b03163b156123dc565b506080810151156123d5565b803b1561026b57600060405180926308d450a160e01b82528183816125508a60048301612217565b03925af1908161258a575b506125845761207961256b612234565b6040516309c2532560e01b815291829160048301612264565b386123cb565b80612599600061259f936101e3565b80610432565b3861255b565b85965060206125e19601516125c460608901516001600160a01b031690565b906125db60208a51016001600160401b0390511690565b9261383d565b903880808080612350565b6306e34e6560e31b60005260046000fd5b6040519061260a82610172565b60006060838281528260208201528260408201520152565b9060405161262f81610172565b606060ff600183958054855201548181166020850152818160081c16604085015260101c161515910152565b906020825491828152019160005260206000209060005b81811061267f5750505090565b82546001600160a01b0316845260209093019260019283019201612672565b90610213611e2e926040519384809261265b565b35906001600160e01b038216820361026b57565b81601f8201121561026b578035906126dd82610243565b926126eb60405194856101e3565b82845260208085019360061b8301019181831161026b57602001925b828410612715575050505090565b60408483031261026b576020604091825161272f816101ad565b61273887610286565b81526127458388016126b2565b83820152815201930192612707565b81601f8201121561026b5780359061276b82610243565b9261277960405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b8385106127a557505050505090565b84356001600160401b03811161026b57820160a0818503601f19011261026b57604051916127d283610192565b6127de60208301610286565b83526040820135926001600160401b03841161026b5760a083612808886020809881980101610301565b8584015261281860608201610286565b604084015261282960808201610286565b606084015201356080820152815201940193612796565b81601f8201121561026b5780359061285782610243565b9261286560405194856101e3565b82845260208085019360061b8301019181831161026b57602001925b82841061288f575050505090565b60408483031261026b57602060409182516128a9816101ad565b863581528287013583820152815201930192612881565b60208183031261026b578035906001600160401b03821161026b570160608183031261026b57604051916128f3836101c8565b81356001600160401b03811161026b57820160408183031261026b576040519061291c826101ad565b80356001600160401b03811161026b57810183601f8201121561026b57803561294481610243565b9161295260405193846101e3565b81835260208084019260061b8201019086821161026b57602001915b8183106129ea5750505082526020810135906001600160401b03821161026b5761299a918491016126c6565b6020820152835260208201356001600160401b03811161026b57816129c0918401612754565b602084015260408201356001600160401b03811161026b576129e29201612840565b604082015290565b60408388031261026b5760206040918251612a04816101ad565b8535612a0f8161025a565b8152612a1c8387016126b2565b8382015281520192019161296e565b9080602083519182815201916020808360051b8301019401926000915b838310612a5757505050505090565b9091929394602080600192601f198582030186528851906001600160401b038251168152608080612a958585015160a08786015260a08501906105da565b936001600160401b0360408201511660408501526001600160401b036060820151166060850152015191015297019301930191939290612a48565b916001600160a01b03612af192168352606060208401526060830190612a2b565b9060408183039101526020808351928381520192019060005b818110612b175750505090565b8251805185526020908101518186015260409094019390920191600101612b0a565b906020808351928381520192019060005b818110612b575750505090565b825180516001600160401b031685526020908101516001600160e01b03168186015260409094019390920191600101612b4a565b9190604081019083519160408252825180915260206060830193019060005b818110612bcb57505050602061031c93940151906020818403910152612b39565b825180516001600160a01b031686526020908101516001600160e01b03168187015260409095019490920191600101612baa565b90602061031c928181520190612b8b565b9081602091031261026b575161031c8161029a565b9091612c3c61031c936040845260408401906105da565b916020818403910152611d97565b6001600160401b036001911601906001600160401b038211611ebf57565b9091612c7f61031c93604084526040840190612a2b565b916020818403910152612b8b565b929693959190979497612ca2828201826128c0565b98612cb66124fd60045460ff9060c01c1690565b61310a575b895180515115908115916130fb575b50613022575b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316999860208a019860005b8a518051821015612fc05781612d1991611d49565b518d612d2c82516001600160401b031690565b604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b1660048201529091602090829060249082905afa9081156124d857600091612f92575b50612f7557612d7b906139a4565b60208201805160208151910120906001830191612d9783611e1a565b6020815191012003612f58575050805460408301516001600160401b039081169160a81c168114801590612f30575b612ede57506080820151908115612ecd57612e1782612e08612def86516001600160401b031690565b6001600160401b0316600052600a602052604060002090565b90600052602052604060002090565b54612e99578291612e7d612e9292612e44612e3f60606001999801516001600160401b031690565b612c4a565b67ffffffffffffffff60a81b197cffffffffffffffff00000000000000000000000000000000000000000083549260a81b169116179055565b612e08612def4294516001600160401b031690565b5501612d04565b50612eae611c6592516001600160401b031690565b6332cf0cbf60e01b6000526001600160401b0316600452602452604490565b63504570e360e01b60005260046000fd5b82611c6591612f086060612ef984516001600160401b031690565b9301516001600160401b031690565b636af0786b60e11b6000526001600160401b0392831660045290821660245216604452606490565b50612f486107e960608501516001600160401b031690565b6001600160401b03821611612dc6565b5161207960405192839263b80d8fa960e01b845260048401612c25565b637edeb53960e11b6000526001600160401b031660045260246000fd5b612fb3915060203d8111612fb9575b612fab81836101e3565b810190612c10565b38612d6d565b503d612fa1565b505061301c9496989b507f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e46102139b613014949597999b5190519061300a60405192839283612c68565b0390a13691610c9f565b943691610c9f565b93613d73565b613037602086015b356001600160401b031690565b600b546001600160401b03828116911610156130df5761306d906001600160401b03166001600160401b0319600b541617600b55565b613085611a4e611a4e6004546001600160a01b031690565b8a5190803b1561026b57604051633937306f60e01b81529160009183918290849082906130b59060048301612bff565b03925af180156124d8576130ca575b50612cd0565b8061259960006130d9936101e3565b386130c4565b5060208a015151612cd057632261116760e01b60005260046000fd5b60209150015151151538612cca565b60208a0151805161311c575b50612cbb565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169060408c0151823b1561026b57604051633854844f60e11b815292600092849283918291613178913060048501612ad0565b03915afa80156124d8571561311657806125996000613196936101e3565b38613116565b604051906131ab6020836101e3565b6000808352366020840137565b6131c06134a1565b60005b8151811015611c99576131d68183611d49565b51906040820160ff6131e9825160ff1690565b161561348b57602083015160ff169261320f8460ff166000526002602052604060002090565b916001830191825461322a6132248260ff1690565b60ff1690565b613450575061325761323f6060830151151590565b845462ff0000191690151560101b62ff000016178455565b60a081019182516101008151116133f85780511561343a576003860161328561327f8261269e565b8a614e5a565b6060840151613315575b947fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547946002946132f16132e161330f9a966132da8760019f9c6132d56133079a8f614fbb565b613f87565b5160ff1690565b845460ff191660ff821617909455565b519081855551906040519586950190888661400d565b0390a161503d565b016131c3565b9794600287939597019661333161332b8961269e565b88614e5a565b6080850151946101008651116134245785516133596132246133548a5160ff1690565b613f73565b101561340e5785518451116133f8576132f16132e17fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547986132da8760019f6132d561330f9f9a8f6133e060029f6133da6133079f8f906132d584926133bf845160ff1690565b908054909161ff001990911660089190911b61ff0016179055565b82614eee565b505050979c9f50975050969a5050509450945061328f565b631b3fab5160e11b600052600160045260246000fd5b631b3fab5160e11b600052600360045260246000fd5b631b3fab5160e11b600052600260045260246000fd5b631b3fab5160e11b600052600560045260246000fd5b60101c60ff1661346b6134666060840151151590565b151590565b90151514613257576321fd80df60e21b60005260ff861660045260246000fd5b631b3fab5160e11b600090815260045260246000fd5b6001600160a01b036001541633036134b557565b6315ae3a6f60e11b60005260046000fd5b604051602081019060008252602081526134e16040826101e3565b51902090565b8181106134f2575050565b600081556001016134e7565b9190601f811161350d57505050565b610213926000526020600020906020601f840160051c83019310613539575b601f0160051c01906134e7565b909150819061352c565b91909182516001600160401b03811161018d5761356a816135648454611d5d565b846134fe565b6020601f82116001146135ab57819061359c9394956000926135a0575b50508160011b916000199060031b1c19161790565b9055565b015190503880613587565b601f198216906135c084600052602060002090565b9160005b8181106135fc575095836001959697106135e3575b505050811b019055565b015160001960f88460031b161c191690553880806135d9565b9192602060018192868b0151815501940192016135c4565b90600160a061031c93602081526001600160401b0384546001600160a01b038116602084015260ff81851c161515604084015260a81c166060820152608080820152019101611d97565b906001600160401b0361369e921660005260096020526701ffffffffffffff60406000209160071c166001600160401b0316600052602052604060002090565b5490565b7f00000000000000000000000000000000000000000000000000000000000000004681036136cd5750565b630f01ce8560e01b6000526004524660245260446000fd5b91909180511561378757825115926020916040519261370481856101e3565b60008452601f19810160005b8181106137635750505060005b815181101561375b578061374461373660019385611d49565b51881561374a57869061414c565b0161371d565b6137548387611d49565b519061414c565b505050509050565b8290604051613771816101ad565b6000815260608382015282828901015201613710565b63c2e5347d60e01b60005260046000fd5b60405160c081018181106001600160401b0382111761018d5760609160a0916040526137c2611c9d565b815282602082015282604082015260008382015260006080820152015290565b9190811015611d445760051b0190565b3561031c81610965565b9190811015611d445760051b81013590601e198136030182121561026b5701908135916001600160401b03831161026b57602001823603811361026b579190565b9092949193979681519661385088610243565b9761385e604051998a6101e3565b80895261386d601f1991610243565b0160005b81811061393f57505060005b835181101561393257806138c48c8a8a8a6138be6138b7878d6138b0828f8f9d8f9e60019f816138e0575b505050611d49565b51976137fc565b36916102ca565b93614991565b6138ce828c611d49565b526138d9818b611d49565b500161387d565b63ffffffff6138f86138f38585856137e2565b6137f2565b16156138a8576139289261390f926138f3926137e2565b604061391b8585611d49565b51019063ffffffff169052565b8f8f9083916138a8565b5096985050505050505050565b60209061394a612129565b82828d01015201613871565b6139676385572ffb60e01b82614cf4565b9081613981575b81613977575090565b61031c9150614cc6565b905061398c81614c4b565b159061396e565b61396763aff2afbf60e01b82614cf4565b6001600160401b031680600052600860205260406000209060ff825460a01c16156139cd575090565b63ed053c5960e01b60005260045260246000fd5b6084019081608411611ebf57565b60a001908160a011611ebf57565b91908201809211611ebf57565b6003111561089e57565b600382101561089e5752565b90610213604051613a30816101ad565b602060ff829554818116845260081c169101613a14565b8054821015611d445760005260206000200190600090565b60ff60019116019060ff8211611ebf57565b60ff601b9116019060ff8211611ebf57565b90606092604091835260208301370190565b6001600052600260205293613ac97fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0612622565b93853594613ad6856139e1565b6060820190613ae58251151590565b613d45575b803603613d2d57508151878103613d145750613b046136a2565b60016000526003602052613b53613b4e7fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c5b336001600160a01b0316600052602052604060002090565b613a20565b60026020820151613b6381613a0a565b613b6c81613a0a565b149081613cac575b5015613c9b5751613bd2575b50505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef090613bb661302a60019460200190565b604080519283526001600160401b0391909116602083015290a2565b613bf3613224613bee602085969799989a955194015160ff1690565b613a5f565b03613c8a578151835103613c7957613c716000613bb69461302a94613c3d7f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09960019b36916102ca565b60208151910120604051613c6881613c5a89602083019586613a83565b03601f1981018352826101e3565b5190208a614d24565b948394613b80565b63a75d88af60e01b60005260046000fd5b6371253a2560e01b60005260046000fd5b631b41e11d60e31b60005260046000fd5b60016000526002602052613d0c9150611a4e90613cf990613cf360037fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e05b01915160ff1690565b90613a47565b90546001600160a01b039160031b1c1690565b331438613b74565b6324f7d61360e21b600052600452602487905260446000fd5b638e1192e160e01b6000526004523660245260446000fd5b613d6e90613d68613d5e613d598751611ea9565b6139ef565b613d688851611ea9565b906139fd565b613aea565b60008052600260205294909390929091613dac7fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b612622565b94863595613db9836139e1565b6060820190613dc88251151590565b613f50575b803603613d2d57508151888103613f375750613de76136a2565b600080526003602052613e1c613b4e7f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff613b36565b60026020820151613e2c81613a0a565b613e3581613a0a565b149081613eee575b5015613c9b5751613e80575b5050505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef090613bb661302a60009460200190565b613e9c613224613bee602087989a999b96975194015160ff1690565b03613c8a578351865103613c79576000967f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef096613bb695613c3d613ee59461302a9736916102ca565b94839438613e49565b600080526002602052613f2f9150611a4e90613cf990613cf360037fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b613cea565b331438613e3d565b6324f7d61360e21b600052600452602488905260446000fd5b613f6e90613d68613f64613d598951611ea9565b613d688a51611ea9565b613dcd565b60ff166003029060ff8216918203611ebf57565b8151916001600160401b03831161018d5768010000000000000000831161018d576020908254848455808510613ff0575b500190600052602060002060005b838110613fd35750505050565b60019060206001600160a01b038551169401938184015501613fc6565b6140079084600052858460002091820191016134e7565b38613fb8565b95949392909160ff61403293168752602087015260a0604087015260a086019061265b565b84810360608601526020808351928381520192019060005b818110614065575050509060806102139294019060ff169052565b82516001600160a01b031684526020938401939092019160010161404a565b600654811015611d445760066000527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f015490565b6001600160401b0361031c94938160609416835216602082015281604082015201906105da565b60409061031c9392815281602082015201906105da565b9291906001600160401b0390816064951660045216602452600481101561089e57604452565b94939261413660609361414793885260208801906108a3565b6080604087015260808601906105da565b930152565b9061415e82516001600160401b031690565b8151604051632cbc26bb60e01b815267ffffffffffffffff60801b608084901b1660048201529015159391906001600160401b038216906020816024817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9081156124d85760009161487a575b50614838576020830191825151948515614808576040850180515187036147f75761420087611cef565b957f00000000000000000000000000000000000000000000000000000000000000006142366001614230876139a4565b01611e1a565b6020815191012060405161429681613c5a6020820194868b876001600160401b036060929594938160808401977f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f85521660208401521660408201520152565b519020906001600160401b031660005b8a811061475f5750505080608060606142c69301519101519088866152ed565b9788156147415760005b8881106142e35750505050505050505050565b5a6142ef828951611d49565b51805160600151614309906001600160401b031688611ed1565b61431281610894565b8015908d828315938461472e575b156146eb576060881561466e5750614347602061433d898d611d49565b5101519242611ec4565b60045461435c9060a01c63ffffffff16611fff565b10801561465b575b1561463d57614373878b611d49565b5151614627575b845160800151614392906001600160401b03166107e9565b61456f575b506143a3868951611d49565b5160a08501515181510361453357936144089695938c938f966143e88e958c926143e26143dc60608951016001600160401b0390511690565b8961531f565b866154c6565b9a90809661440260608851016001600160401b0390511690565b906153a7565b6144e1575b505061441882610894565b60028203614499575b60019661448f7f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b936001600160401b039351926144806144778b61446f60608801516001600160401b031690565b96519b611d49565b51985a90611ec4565b9160405195869516988561411d565b0390a45b016142d0565b915091939492506144a982610894565b600382036144bd578b929493918a91614421565b51606001516349362d1f60e11b600052611c6591906001600160401b0316896140f7565b6144ea84610894565b6003840361440d579092949550614502919350610894565b614512578b92918a91388061440d565b5151604051632b11b8d960e01b8152908190612079908790600484016140e0565b611c658b61454d60608851016001600160401b0390511690565b631cfe6d8b60e01b6000526001600160401b0391821660045216602452604490565b61457883610894565b614583575b38614397565b8351608001516001600160401b0316602080860151918c6145b860405194859384936370701e5760e11b8552600485016140b9565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af19081156124d857600091614609575b5061457d575050505050600190614493565b614621915060203d8111612fb957612fab81836101e3565b386145f7565b614631878b611d49565b5151608086015261437a565b6354e7e43160e11b6000526001600160401b038b1660045260246000fd5b5061466583610894565b60038314614364565b91508361467a84610894565b1561437a575060019594506146e392506146c191507f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe651209351016001600160401b0390511690565b604080516001600160401b03808c168252909216602083015290918291820190565b0390a1614493565b5050505060019291506146e36146c160607f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c9351016001600160401b0390511690565b5061473883610894565b60038314614320565b633ee8bd3f60e11b6000526001600160401b03841660045260246000fd5b61476a818a51611d49565b518051604001516001600160401b03168381036147da57508051602001516001600160401b03168981036147b75750906147a6846001936151e5565b6147b0828d611d49565b52016142a6565b636c95f1eb60e01b6000526001600160401b03808a166004521660245260446000fd5b631c21951160e11b6000526001600160401b031660045260246000fd5b6357e0e08360e01b60005260046000fd5b611c6561481c86516001600160401b031690565b63676cf24b60e11b6000526001600160401b0316600452602490565b5092915050612f75576040516001600160401b039190911681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d89493390602090a1565b614893915060203d602011612fb957612fab81836101e3565b386141d6565b9081602091031261026b575161031c8161025a565b9061031c916020815260e061494c6149376148d7855161010060208701526101208601906105da565b60208601516001600160401b0316604086015260408601516001600160a01b0316606086015260608601516080860152614921608087015160a08701906001600160a01b03169052565b60a0860151858203601f190160c08701526105da565b60c0850151848203601f1901848601526105da565b92015190610100601f19828503019101526105da565b6040906001600160a01b0361031c949316815281602082015201906105da565b9081602091031261026b575190565b9193929361499d612129565b5060208301516001600160a01b031660405163bbe4f6db60e01b81526001600160a01b038216600482015290959092602084806024810103816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9384156124d857600094614c1a575b506001600160a01b0384169586158015614c08575b614bea57614acf614af892613c5a92614a53614a4c611fff60408c015163ffffffff1690565b8c89615607565b9690996080810151614a816060835193015193614a6e610224565b9687526001600160401b03166020870152565b6001600160a01b038a16604086015260608501526001600160a01b038d16608085015260a084015260c083015260e0820152604051633907753760e01b6020820152928391602483016148ae565b82857f000000000000000000000000000000000000000000000000000000000000000092615695565b94909115614bce5750805160208103614bb5575090614b21826020808a95518301019101614982565b956001600160a01b03841603614b59575b5050505050614b51614b42610234565b6001600160a01b039093168352565b602082015290565b614b6c93614b6691611ec4565b91615607565b50908082108015614ba2575b614b8457808481614b32565b63a966e21f60e01b6000908152600493909352602452604452606490fd5b5082614bae8284611ec4565b1415614b78565b631e3be00960e21b600052602060045260245260446000fd5b612079604051928392634ff17cad60e11b845260048401614962565b63ae9b4ce960e01b6000526001600160a01b03851660045260246000fd5b50614c156124fd86613993565b614a26565b614c3d91945060203d602011614c44575b614c3581836101e3565b810190614899565b9238614a11565b503d614c2b565b60405160208101916301ffc9a760e01b835263ffffffff60e01b602483015260248252614c796044836101e3565b6179185a10614cb5576020926000925191617530fa6000513d82614ca9575b5081614ca2575090565b9050151590565b60201115915038614c98565b63753fa58960e11b60005260046000fd5b60405160208101916301ffc9a760e01b83526301ffc9a760e01b602483015260248252614c796044836101e3565b6040519060208201926301ffc9a760e01b845263ffffffff60e01b16602483015260248252614c796044836101e3565b919390926000948051946000965b868810614d43575050505050505050565b6020881015611d445760206000614d5b878b1a613a71565b614d658b87611d49565b5190614d9c614d748d8a611d49565b5160405193849389859094939260ff6060936080840197845216602083015260408201520152565b838052039060015afa156124d857614de2613b4e600051614dca8960ff166000526003602052604060002090565b906001600160a01b0316600052602052604060002090565b9060016020830151614df381613a0a565b614dfc81613a0a565b03614e4957614e19614e0f835160ff1690565b60ff600191161b90565b8116614e3857614e2f614e0f6001935160ff1690565b17970196614d32565b633d9ef1f160e21b60005260046000fd5b636518c33d60e11b60005260046000fd5b91909160005b8351811015614eb35760019060ff831660005260036020526000614eac604082206001600160a01b03614e93858a611d49565b51166001600160a01b0316600052602052604060002090565b5501614e60565b50509050565b8151815460ff191660ff919091161781559060200151600381101561089e57815461ff00191660089190911b61ff0016179055565b919060005b8151811015614eb357614f09611ba78284611d49565b90614f32614f2883614dca8860ff166000526003602052604060002090565b5460081c60ff1690565b614f3b81613a0a565b614fa6576001600160a01b03821615614f9557614f8f600192614f8a614f5f610234565b60ff8516815291614f738660208501613a14565b614dca8960ff166000526003602052604060002090565b614eb9565b01614ef3565b63d6c62c9b60e01b60005260046000fd5b631b3fab5160e11b6000526004805260246000fd5b919060005b8151811015614eb357614fd6611ba78284611d49565b90614ff5614f2883614dca8860ff166000526003602052604060002090565b614ffe81613a0a565b614fa6576001600160a01b03821615614f9557615037600192614f8a615022610234565b60ff8516815291614f73600260208501613a14565b01614fc0565b60ff1680600052600260205260ff60016040600020015460101c1690801560001461508b57501561507a576001600160401b0319600b5416600b55565b6317bd8dd160e11b60005260046000fd5b6001146150955750565b61509b57565b6307b8c74d60e51b60005260046000fd5b8060005260076020526040600020541560001461512a576006546801000000000000000081101561018d57600181016006556000600654821015611d4457600690527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01819055600654906000526007602052604060002055600190565b50600090565b9080602083519182815201916020808360051b8301019401926000915b83831061515c57505050505090565b9091929394602080600192601f198582030186528851906080806151bf61518c855160a0865260a08601906105da565b6001600160a01b0387870151168786015263ffffffff6040870151166040860152606086015185820360608701526105da565b9301519101529701930193019193929061514d565b90602061031c928181520190615130565b6134e1815180519061527961520460608601516001600160a01b031690565b613c5a61521b60608501516001600160401b031690565b936152346080808a01519201516001600160401b031690565b90604051958694602086019889936001600160401b036080946001600160a01b0382959998949960a089019a8952166020880152166040860152606085015216910152565b519020613c5a6020840151602081519101209360a06040820151602081519101209101516040516152b281613c5a6020820194856151d4565b51902090604051958694602086019889919260a093969594919660c08401976000855260208501526040840152606083015260808201520152565b926001600160401b039261530092615752565b9116600052600a60205260406000209060005260205260406000205490565b607f8216906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611ebf576153a4916001600160401b03615362858461365e565b921660005260096020526701ffffffffffffff60406000209460071c169160036001831b921b19161792906001600160401b0316600052602052604060002090565b55565b9091607f83166801fffffffffffffffe6001600160401b0382169160011b169080820460021490151715611ebf576153df848461365e565b600483101561089e576001600160401b036153a49416600052600960205260036701ffffffffffffff60406000209660071c1693831b921b19161792906001600160401b0316600052602052604060002090565b9080602083519182815201916020808360051b8301019401926000915b83831061545f57505050505090565b909192939460208061547d600193601f1986820301875289516105da565b97019301930191939290615450565b906020808351928381520192019060005b8181106154aa5750505090565b825163ffffffff1684526020938401939092019160010161549d565b91606092303b1561026b576155c860a0926155b66000956155a46040519889978897630304c3e160e51b89528260048a01526001600160401b0360808251805160648d01528260208201511660848d01528260408201511660a48d015282868201511660c48d015201511660e48a01526155856155706155598b61014061010460208701519201526101a48d01906105da565b60408401518c8203606319016101248e01526105da565b938201516001600160a01b03166101448b0152565b60808101516101648a0152015187820360631901610184890152615130565b85810360031901602487015290615433565b8381036003190160448501529061548c565b038183305af190816155f2575b506155e7576155e2612234565b600391565b60029061031c6105a2565b806125996000615601936101e3565b386155d5565b6040516370a0823160e01b60208201526001600160a01b0390911660248201529192916156649061563b8160448101613c5a565b84837f000000000000000000000000000000000000000000000000000000000000000092615695565b92909115614bce5750805160208103614bb557509061568f8260208061031c95518301019101614982565b93611ec4565b9391936156a260846102af565b946156b060405196876101e3565b608486526156be60846102af565b602087019590601f1901368737833b15615741575a90808210615730578291038060061c9003111561571f576000918291825a9560208451940192f1905a9003923d9060848211615716575b6000908287523e929190565b6084915061570a565b6337c3be2960e01b60005260046000fd5b632be8ca8b60e21b60005260046000fd5b63030ed58f60e21b60005260046000fd5b80519282519084156158ae57610101851115806158a2575b156157d1578185019460001986019561010087116157d15786156158925761579187611cef565b9660009586978795885b8481106157f65750505050506001190180951493846157ec575b5050826157e2575b5050156157d1576157cd91611d49565b5190565b6309bde33960e01b60005260046000fd5b14905038806157bd565b14925038806157b5565b6001811b8281160361588457868a101561586f5761581860018b019a85611d49565b51905b8c888c101561585b575061583360018c019b86611d49565b515b818d116157d157615854828f9261584e906001966158bf565b92611d49565b520161579b565b60018d019c61586991611d49565b51615835565b61587d60018c019b8d611d49565b519061581b565b61587d600189019884611d49565b5050505090506157cd9150611d37565b5061010182111561576a565b630469ac9960e21b60005260046000fd5b818110156158d1579061031c916158d6565b61031c915b906040519060208201926001845260408301526060820152606081526134e16080826101e356fea164736f6c634300081a000a49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b",
}

var OffRampABI = OffRampMetaData.ABI

var OffRampBin = OffRampMetaData.Bin

func DeployOffRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OffRampStaticConfig, dynamicConfig OffRampDynamicConfig, sourceChainConfigs []OffRampSourceChainConfigArgs) (common.Address, *types.Transaction, *OffRamp, error) {
	parsed, err := OffRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OffRampBin), backend, staticConfig, dynamicConfig, sourceChainConfigs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffRamp{address: address, abi: *parsed, OffRampCaller: OffRampCaller{contract: contract}, OffRampTransactor: OffRampTransactor{contract: contract}, OffRampFilterer: OffRampFilterer{contract: contract}}, nil
}

type OffRamp struct {
	address common.Address
	abi     abi.ABI
	OffRampCaller
	OffRampTransactor
	OffRampFilterer
}

type OffRampCaller struct {
	contract *bind.BoundContract
}

type OffRampTransactor struct {
	contract *bind.BoundContract
}

type OffRampFilterer struct {
	contract *bind.BoundContract
}

type OffRampSession struct {
	Contract     *OffRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OffRampCallerSession struct {
	Contract *OffRampCaller
	CallOpts bind.CallOpts
}

type OffRampTransactorSession struct {
	Contract     *OffRampTransactor
	TransactOpts bind.TransactOpts
}

type OffRampRaw struct {
	Contract *OffRamp
}

type OffRampCallerRaw struct {
	Contract *OffRampCaller
}

type OffRampTransactorRaw struct {
	Contract *OffRampTransactor
}

func NewOffRamp(address common.Address, backend bind.ContractBackend) (*OffRamp, error) {
	abi, err := abi.JSON(strings.NewReader(OffRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOffRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffRamp{address: address, abi: abi, OffRampCaller: OffRampCaller{contract: contract}, OffRampTransactor: OffRampTransactor{contract: contract}, OffRampFilterer: OffRampFilterer{contract: contract}}, nil
}

func NewOffRampCaller(address common.Address, caller bind.ContractCaller) (*OffRampCaller, error) {
	contract, err := bindOffRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampCaller{contract: contract}, nil
}

func NewOffRampTransactor(address common.Address, transactor bind.ContractTransactor) (*OffRampTransactor, error) {
	contract, err := bindOffRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampTransactor{contract: contract}, nil
}

func NewOffRampFilterer(address common.Address, filterer bind.ContractFilterer) (*OffRampFilterer, error) {
	contract, err := bindOffRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffRampFilterer{contract: contract}, nil
}

func bindOffRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OffRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OffRamp *OffRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRamp.Contract.OffRampCaller.contract.Call(opts, result, method, params...)
}

func (_OffRamp *OffRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRamp.Contract.OffRampTransactor.contract.Transfer(opts)
}

func (_OffRamp *OffRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRamp.Contract.OffRampTransactor.contract.Transact(opts, method, params...)
}

func (_OffRamp *OffRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_OffRamp *OffRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRamp.Contract.contract.Transfer(opts)
}

func (_OffRamp *OffRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRamp.Contract.contract.Transact(opts, method, params...)
}

func (_OffRamp *OffRampCaller) CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "ccipReceive", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_OffRamp *OffRampSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRamp.Contract.CcipReceive(&_OffRamp.CallOpts, arg0)
}

func (_OffRamp *OffRampCallerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRamp.Contract.CcipReceive(&_OffRamp.CallOpts, arg0)
}

func (_OffRamp *OffRampCaller) GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getAllSourceChainConfigs")

	if err != nil {
		return *new([]uint64), *new([]OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)
	out1 := *abi.ConvertType(out[1], new([]OffRampSourceChainConfig)).(*[]OffRampSourceChainConfig)

	return out0, out1, err

}

func (_OffRamp *OffRampSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetAllSourceChainConfigs(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetAllSourceChainConfigs(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OffRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampDynamicConfig)).(*OffRampDynamicConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRamp.Contract.GetDynamicConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRamp.Contract.GetDynamicConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getExecutionState", sourceChainSelector, sequenceNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_OffRamp *OffRampSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRamp.Contract.GetExecutionState(&_OffRamp.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRamp *OffRampCallerSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRamp.Contract.GetExecutionState(&_OffRamp.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRamp *OffRampCaller) GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getLatestPriceSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OffRamp *OffRampSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRamp.Contract.GetLatestPriceSequenceNumber(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRamp.Contract.GetLatestPriceSequenceNumber(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getMerkleRoot", sourceChainSelector, root)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OffRamp *OffRampSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRamp.Contract.GetMerkleRoot(&_OffRamp.CallOpts, sourceChainSelector, root)
}

func (_OffRamp *OffRampCallerSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRamp.Contract.GetMerkleRoot(&_OffRamp.CallOpts, sourceChainSelector, root)
}

func (_OffRamp *OffRampCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampSourceChainConfig)).(*OffRampSourceChainConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetSourceChainConfig(&_OffRamp.CallOpts, sourceChainSelector)
}

func (_OffRamp *OffRampCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetSourceChainConfig(&_OffRamp.CallOpts, sourceChainSelector)
}

func (_OffRamp *OffRampCaller) GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OffRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampStaticConfig)).(*OffRampStaticConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRamp.Contract.GetStaticConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRamp.Contract.GetStaticConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "latestConfigDetails", ocrPluginType)

	if err != nil {
		return *new(MultiOCR3BaseOCRConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOCRConfig)).(*MultiOCR3BaseOCRConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRamp.Contract.LatestConfigDetails(&_OffRamp.CallOpts, ocrPluginType)
}

func (_OffRamp *OffRampCallerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRamp.Contract.LatestConfigDetails(&_OffRamp.CallOpts, ocrPluginType)
}

func (_OffRamp *OffRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OffRamp *OffRampSession) Owner() (common.Address, error) {
	return _OffRamp.Contract.Owner(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) Owner() (common.Address, error) {
	return _OffRamp.Contract.Owner(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OffRamp *OffRampSession) TypeAndVersion() (string, error) {
	return _OffRamp.Contract.TypeAndVersion(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) TypeAndVersion() (string, error) {
	return _OffRamp.Contract.TypeAndVersion(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "acceptOwnership")
}

func (_OffRamp *OffRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRamp.Contract.AcceptOwnership(&_OffRamp.TransactOpts)
}

func (_OffRamp *OffRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRamp.Contract.AcceptOwnership(&_OffRamp.TransactOpts)
}

func (_OffRamp *OffRampTransactor) ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "applySourceChainConfigUpdates", sourceChainConfigUpdates)
}

func (_OffRamp *OffRampSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.ApplySourceChainConfigUpdates(&_OffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRamp *OffRampTransactorSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.ApplySourceChainConfigUpdates(&_OffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRamp *OffRampTransactor) Commit(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "commit", reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactorSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactor) Execute(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "execute", reportContext, report)
}

func (_OffRamp *OffRampSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Execute(&_OffRamp.TransactOpts, reportContext, report)
}

func (_OffRamp *OffRampTransactorSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Execute(&_OffRamp.TransactOpts, reportContext, report)
}

func (_OffRamp *OffRampTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampTransactorSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampTransactor) ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "manuallyExecute", reports, gasLimitOverrides)
}

func (_OffRamp *OffRampSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.Contract.ManuallyExecute(&_OffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRamp *OffRampTransactorSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.Contract.ManuallyExecute(&_OffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRamp *OffRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OffRamp *OffRampSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRamp.Contract.SetDynamicConfig(&_OffRamp.TransactOpts, dynamicConfig)
}

func (_OffRamp *OffRampTransactorSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRamp.Contract.SetDynamicConfig(&_OffRamp.TransactOpts, dynamicConfig)
}

func (_OffRamp *OffRampTransactor) SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "setOCR3Configs", ocrConfigArgs)
}

func (_OffRamp *OffRampSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.SetOCR3Configs(&_OffRamp.TransactOpts, ocrConfigArgs)
}

func (_OffRamp *OffRampTransactorSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.SetOCR3Configs(&_OffRamp.TransactOpts, ocrConfigArgs)
}

func (_OffRamp *OffRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_OffRamp *OffRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRamp.Contract.TransferOwnership(&_OffRamp.TransactOpts, to)
}

func (_OffRamp *OffRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRamp.Contract.TransferOwnership(&_OffRamp.TransactOpts, to)
}

type OffRampAlreadyAttemptedIterator struct {
	Event *OffRampAlreadyAttempted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampAlreadyAttemptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampAlreadyAttempted)
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
		it.Event = new(OffRampAlreadyAttempted)
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

func (it *OffRampAlreadyAttemptedIterator) Error() error {
	return it.fail
}

func (it *OffRampAlreadyAttemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampAlreadyAttempted struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampAlreadyAttemptedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return &OffRampAlreadyAttemptedIterator{contract: _OffRamp.contract, event: "AlreadyAttempted", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampAlreadyAttempted) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampAlreadyAttempted)
				if err := _OffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseAlreadyAttempted(log types.Log) (*OffRampAlreadyAttempted, error) {
	event := new(OffRampAlreadyAttempted)
	if err := _OffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampCommitReportAcceptedIterator struct {
	Event *OffRampCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampCommitReportAccepted)
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
		it.Event = new(OffRampCommitReportAccepted)
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

func (it *OffRampCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *OffRampCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampCommitReportAccepted struct {
	MerkleRoots  []InternalMerkleRoot
	PriceUpdates InternalPriceUpdates
	Raw          types.Log
}

func (_OffRamp *OffRampFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampCommitReportAcceptedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &OffRampCommitReportAcceptedIterator{contract: _OffRamp.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampCommitReportAccepted)
				if err := _OffRamp.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseCommitReportAccepted(log types.Log) (*OffRampCommitReportAccepted, error) {
	event := new(OffRampCommitReportAccepted)
	if err := _OffRamp.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampConfigSetIterator struct {
	Event *OffRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampConfigSet)
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
		it.Event = new(OffRampConfigSet)
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

func (it *OffRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampConfigSet struct {
	OcrPluginType uint8
	ConfigDigest  [32]byte
	Signers       []common.Address
	Transmitters  []common.Address
	F             uint8
	Raw           types.Log
}

func (_OffRamp *OffRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OffRampConfigSetIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampConfigSetIterator{contract: _OffRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseConfigSet(log types.Log) (*OffRampConfigSet, error) {
	event := new(OffRampConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampDynamicConfigSetIterator struct {
	Event *OffRampDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampDynamicConfigSet)
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
		it.Event = new(OffRampDynamicConfigSet)
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

func (it *OffRampDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampDynamicConfigSet struct {
	DynamicConfig OffRampDynamicConfig
	Raw           types.Log
}

func (_OffRamp *OffRampFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampDynamicConfigSetIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampDynamicConfigSetIterator{contract: _OffRamp.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampDynamicConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampDynamicConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseDynamicConfigSet(log types.Log) (*OffRampDynamicConfigSet, error) {
	event := new(OffRampDynamicConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampExecutionStateChangedIterator struct {
	Event *OffRampExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampExecutionStateChanged)
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
		it.Event = new(OffRampExecutionStateChanged)
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

func (it *OffRampExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *OffRampExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	MessageHash         [32]byte
	State               uint8
	ReturnData          []byte
	GasUsed             *big.Int
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampExecutionStateChangedIterator, error) {

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

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &OffRampExecutionStateChangedIterator{contract: _OffRamp.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampExecutionStateChanged)
				if err := _OffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseExecutionStateChanged(log types.Log) (*OffRampExecutionStateChanged, error) {
	event := new(OffRampExecutionStateChanged)
	if err := _OffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampOwnershipTransferRequestedIterator struct {
	Event *OffRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampOwnershipTransferRequested)
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
		it.Event = new(OffRampOwnershipTransferRequested)
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

func (it *OffRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OffRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRamp *OffRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampOwnershipTransferRequestedIterator{contract: _OffRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampOwnershipTransferRequested)
				if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffRampOwnershipTransferRequested, error) {
	event := new(OffRampOwnershipTransferRequested)
	if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampOwnershipTransferredIterator struct {
	Event *OffRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampOwnershipTransferred)
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
		it.Event = new(OffRampOwnershipTransferred)
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

func (it *OffRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OffRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRamp *OffRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampOwnershipTransferredIterator{contract: _OffRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampOwnershipTransferred)
				if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseOwnershipTransferred(log types.Log) (*OffRampOwnershipTransferred, error) {
	event := new(OffRampOwnershipTransferred)
	if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampRootRemovedIterator struct {
	Event *OffRampRootRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampRootRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampRootRemoved)
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
		it.Event = new(OffRampRootRemoved)
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

func (it *OffRampRootRemovedIterator) Error() error {
	return it.fail
}

func (it *OffRampRootRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampRootRemoved struct {
	Root [32]byte
	Raw  types.Log
}

func (_OffRamp *OffRampFilterer) FilterRootRemoved(opts *bind.FilterOpts) (*OffRampRootRemovedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return &OffRampRootRemovedIterator{contract: _OffRamp.contract, event: "RootRemoved", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampRootRemoved) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampRootRemoved)
				if err := _OffRamp.contract.UnpackLog(event, "RootRemoved", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseRootRemoved(log types.Log) (*OffRampRootRemoved, error) {
	event := new(OffRampRootRemoved)
	if err := _OffRamp.contract.UnpackLog(event, "RootRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampSkippedAlreadyExecutedMessageIterator struct {
	Event *OffRampSkippedAlreadyExecutedMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSkippedAlreadyExecutedMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSkippedAlreadyExecutedMessage)
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
		it.Event = new(OffRampSkippedAlreadyExecutedMessage)
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

func (it *OffRampSkippedAlreadyExecutedMessageIterator) Error() error {
	return it.fail
}

func (it *OffRampSkippedAlreadyExecutedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSkippedAlreadyExecutedMessage struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampSkippedAlreadyExecutedMessageIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return &OffRampSkippedAlreadyExecutedMessageIterator{contract: _OffRamp.contract, event: "SkippedAlreadyExecutedMessage", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampSkippedAlreadyExecutedMessage) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSkippedAlreadyExecutedMessage)
				if err := _OffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampSkippedAlreadyExecutedMessage, error) {
	event := new(OffRampSkippedAlreadyExecutedMessage)
	if err := _OffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampSkippedReportExecutionIterator struct {
	Event *OffRampSkippedReportExecution

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSkippedReportExecutionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSkippedReportExecution)
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
		it.Event = new(OffRampSkippedReportExecution)
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

func (it *OffRampSkippedReportExecutionIterator) Error() error {
	return it.fail
}

func (it *OffRampSkippedReportExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSkippedReportExecution struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampSkippedReportExecutionIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return &OffRampSkippedReportExecutionIterator{contract: _OffRamp.contract, event: "SkippedReportExecution", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampSkippedReportExecution) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSkippedReportExecution)
				if err := _OffRamp.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSkippedReportExecution(log types.Log) (*OffRampSkippedReportExecution, error) {
	event := new(OffRampSkippedReportExecution)
	if err := _OffRamp.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampSourceChainConfigSetIterator struct {
	Event *OffRampSourceChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSourceChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSourceChainConfigSet)
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
		it.Event = new(OffRampSourceChainConfigSet)
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

func (it *OffRampSourceChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampSourceChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceConfig        OffRampSourceChainConfig
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampSourceChainConfigSetIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OffRampSourceChainConfigSetIterator{contract: _OffRamp.contract, event: "SourceChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSourceChainConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSourceChainConfigSet(log types.Log) (*OffRampSourceChainConfigSet, error) {
	event := new(OffRampSourceChainConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampSourceChainSelectorAddedIterator struct {
	Event *OffRampSourceChainSelectorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSourceChainSelectorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSourceChainSelectorAdded)
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
		it.Event = new(OffRampSourceChainSelectorAdded)
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

func (it *OffRampSourceChainSelectorAddedIterator) Error() error {
	return it.fail
}

func (it *OffRampSourceChainSelectorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSourceChainSelectorAdded struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampSourceChainSelectorAddedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return &OffRampSourceChainSelectorAddedIterator{contract: _OffRamp.contract, event: "SourceChainSelectorAdded", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainSelectorAdded) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSourceChainSelectorAdded)
				if err := _OffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSourceChainSelectorAdded(log types.Log) (*OffRampSourceChainSelectorAdded, error) {
	event := new(OffRampSourceChainSelectorAdded)
	if err := _OffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampStaticConfigSetIterator struct {
	Event *OffRampStaticConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampStaticConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampStaticConfigSet)
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
		it.Event = new(OffRampStaticConfigSet)
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

func (it *OffRampStaticConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampStaticConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampStaticConfigSet struct {
	StaticConfig OffRampStaticConfig
	Raw          types.Log
}

func (_OffRamp *OffRampFilterer) FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampStaticConfigSetIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampStaticConfigSetIterator{contract: _OffRamp.contract, event: "StaticConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampStaticConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampStaticConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseStaticConfigSet(log types.Log) (*OffRampStaticConfigSet, error) {
	event := new(OffRampStaticConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampTransmittedIterator struct {
	Event *OffRampTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampTransmitted)
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
		it.Event = new(OffRampTransmitted)
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

func (it *OffRampTransmittedIterator) Error() error {
	return it.fail
}

func (it *OffRampTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampTransmitted struct {
	OcrPluginType  uint8
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_OffRamp *OffRampFilterer) FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampTransmittedIterator, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return &OffRampTransmittedIterator{contract: _OffRamp.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampTransmitted, ocrPluginType []uint8) (event.Subscription, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampTransmitted)
				if err := _OffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseTransmitted(log types.Log) (*OffRampTransmitted, error) {
	event := new(OffRampTransmitted)
	if err := _OffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OffRamp *OffRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OffRamp.abi.Events["AlreadyAttempted"].ID:
		return _OffRamp.ParseAlreadyAttempted(log)
	case _OffRamp.abi.Events["CommitReportAccepted"].ID:
		return _OffRamp.ParseCommitReportAccepted(log)
	case _OffRamp.abi.Events["ConfigSet"].ID:
		return _OffRamp.ParseConfigSet(log)
	case _OffRamp.abi.Events["DynamicConfigSet"].ID:
		return _OffRamp.ParseDynamicConfigSet(log)
	case _OffRamp.abi.Events["ExecutionStateChanged"].ID:
		return _OffRamp.ParseExecutionStateChanged(log)
	case _OffRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _OffRamp.ParseOwnershipTransferRequested(log)
	case _OffRamp.abi.Events["OwnershipTransferred"].ID:
		return _OffRamp.ParseOwnershipTransferred(log)
	case _OffRamp.abi.Events["RootRemoved"].ID:
		return _OffRamp.ParseRootRemoved(log)
	case _OffRamp.abi.Events["SkippedAlreadyExecutedMessage"].ID:
		return _OffRamp.ParseSkippedAlreadyExecutedMessage(log)
	case _OffRamp.abi.Events["SkippedReportExecution"].ID:
		return _OffRamp.ParseSkippedReportExecution(log)
	case _OffRamp.abi.Events["SourceChainConfigSet"].ID:
		return _OffRamp.ParseSourceChainConfigSet(log)
	case _OffRamp.abi.Events["SourceChainSelectorAdded"].ID:
		return _OffRamp.ParseSourceChainSelectorAdded(log)
	case _OffRamp.abi.Events["StaticConfigSet"].ID:
		return _OffRamp.ParseStaticConfigSet(log)
	case _OffRamp.abi.Events["Transmitted"].ID:
		return _OffRamp.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OffRampAlreadyAttempted) Topic() common.Hash {
	return common.HexToHash("0x3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120")
}

func (OffRampCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4")
}

func (OffRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (OffRampDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0xcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee")
}

func (OffRampExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b")
}

func (OffRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OffRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OffRampRootRemoved) Topic() common.Hash {
	return common.HexToHash("0x202f1139a3e334b6056064c0e9b19fd07e44a88d8f6e5ded571b24cf8c371f12")
}

func (OffRampSkippedAlreadyExecutedMessage) Topic() common.Hash {
	return common.HexToHash("0x3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c")
}

func (OffRampSkippedReportExecution) Topic() common.Hash {
	return common.HexToHash("0xaab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d894933")
}

func (OffRampSourceChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b")
}

func (OffRampSourceChainSelectorAdded) Topic() common.Hash {
	return common.HexToHash("0xf4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9")
}

func (OffRampStaticConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b7")
}

func (OffRampTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_OffRamp *OffRamp) Address() common.Address {
	return _OffRamp.address
}

type OffRampInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error)

	GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error)

	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error)

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

	FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampAlreadyAttemptedIterator, error)

	WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampAlreadyAttempted) (event.Subscription, error)

	ParseAlreadyAttempted(log types.Log) (*OffRampAlreadyAttempted, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*OffRampCommitReportAccepted, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OffRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OffRampConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampDynamicConfigSet) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*OffRampDynamicConfigSet, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*OffRampExecutionStateChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OffRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OffRampOwnershipTransferred, error)

	FilterRootRemoved(opts *bind.FilterOpts) (*OffRampRootRemovedIterator, error)

	WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampRootRemoved) (event.Subscription, error)

	ParseRootRemoved(log types.Log) (*OffRampRootRemoved, error)

	FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampSkippedAlreadyExecutedMessageIterator, error)

	WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampSkippedAlreadyExecutedMessage) (event.Subscription, error)

	ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampSkippedAlreadyExecutedMessage, error)

	FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampSkippedReportExecutionIterator, error)

	WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampSkippedReportExecution) (event.Subscription, error)

	ParseSkippedReportExecution(log types.Log) (*OffRampSkippedReportExecution, error)

	FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampSourceChainConfigSetIterator, error)

	WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error)

	ParseSourceChainConfigSet(log types.Log) (*OffRampSourceChainConfigSet, error)

	FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampSourceChainSelectorAddedIterator, error)

	WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainSelectorAdded) (event.Subscription, error)

	ParseSourceChainSelectorAdded(log types.Log) (*OffRampSourceChainSelectorAdded, error)

	FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampStaticConfigSetIterator, error)

	WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampStaticConfigSet) (event.Subscription, error)

	ParseStaticConfigSet(log types.Log) (*OffRampStaticConfigSet, error)

	FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampTransmitted, ocrPluginType []uint8) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*OffRampTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
