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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reportOnRamp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"configOnRamp\",\"type\":\"bytes\"}],\"name\":\"CommitOnRampMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyBatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenGasOverride\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionTokenGasOverride\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"messageDestChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidMessageDestChainSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidOnRampUpdate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeavesCannotBeEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionGasAmountCountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"errorReason\",\"type\":\"bytes\"}],\"name\":\"MessageValidationError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReleased\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePre\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePost\",\"type\":\"uint256\"}],\"name\":\"ReleaseOrMintBalanceMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"RootAlreadyCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleCommitReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelectorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SkippedReportExecution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"}],\"name\":\"StaticConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllSourceChainConfigs\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReport[]\",\"name\":\"reports\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101406040523480156200001257600080fd5b5060405162006d0638038062006d06833981016040819052620000359162000936565b336000816200005757604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008a576200008a81620001dc565b50504660805260408301516001600160a01b03161580620000b6575060608301516001600160a01b0316155b80620000cd575060808301516001600160a01b0316155b15620000ec576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160401b0316600003620001185760405163c656089560e01b815260040160405180910390fd5b82516001600160401b0390811660a0908152604080860180516001600160a01b0390811660c05260608089018051831660e0526080808b0180518516610100526020808d01805161ffff9081166101205289518f51909c168c52905116908a0152945184169588019590955251821690860152905116908301527fb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b7910160405180910390a1620001c88262000256565b620001d38162000344565b50505062000cde565b336001600160a01b038216036200020657604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03166200027f576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516001600160a01b039889166001600160c01b03199097168717600160a01b63ffffffff958616021760ff60c01b1916600160c01b911515919091021790965560608089018051600580546001600160a01b031916918b169190911790558251968752935190921693850193909352935115159183019190915251909216908201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9060800160405180910390a150565b60005b8151811015620005d957600082828151811062000368576200036862000a16565b60200260200101519050600081602001519050806001600160401b0316600003620003a65760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b0316620003cf576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b03811660009081526008602052604090206060830151600182018054620003fd9062000a2c565b905060000362000460578154600160a81b600160e81b031916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1620004d1565b8154600160a81b90046001600160401b0316600114801590620004a35750805160208201206040516200049890600185019062000a68565b604051809103902014155b15620004d157604051632105803760e11b81526001600160401b038416600482015260240160405180910390fd5b80511580620005075750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b1562000526576040516342bcdf7f60e11b815260040160405180910390fd5b6001820162000536828262000b3b565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b02929092166001600160a81b0319909116171782556200058560066001600160401b038516620005dd565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b83604051620005c1919062000c07565b60405180910390a25050505080600101905062000347565b5050565b6000620005eb8383620005f4565b90505b92915050565b60008181526001830160205260408120546200063d57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620005ee565b506000620005ee565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b038111828210171562000681576200068162000646565b60405290565b60405160a081016001600160401b038111828210171562000681576200068162000646565b604051601f8201601f191681016001600160401b0381118282101715620006d757620006d762000646565b604052919050565b80516001600160401b0381168114620006f757600080fd5b919050565b6001600160a01b03811681146200071257600080fd5b50565b80518015158114620006f757600080fd5b6000608082840312156200073957600080fd5b620007436200065c565b905081516200075281620006fc565b8152602082015163ffffffff811681146200076c57600080fd5b60208201526200077f6040830162000715565b604082015260608201516200079481620006fc565b606082015292915050565b6000601f83601f840112620007b357600080fd5b825160206001600160401b0380831115620007d257620007d262000646565b8260051b620007e3838201620006ac565b9384528681018301938381019089861115620007fe57600080fd5b84890192505b8583101562000929578251848111156200081e5760008081fd5b89016080601f19828d038101821315620008385760008081fd5b620008426200065c565b888401516200085181620006fc565b8152604062000862858201620006df565b8a83015260606200087581870162000715565b838301529385015193898511156200088d5760008081fd5b84860195508f603f870112620008a557600094508485fd5b8a860151945089851115620008be57620008be62000646565b620008cf8b858f88011601620006ac565b93508484528f82868801011115620008e75760008081fd5b60005b8581101562000907578681018301518582018d01528b01620008ea565b5060009484018b01949094525091820152835250918401919084019062000804565b9998505050505050505050565b60008060008385036101408112156200094e57600080fd5b60a08112156200095d57600080fd5b506200096862000687565b6200097385620006df565b8152602085015161ffff811681146200098b57600080fd5b60208201526040850151620009a081620006fc565b60408201526060850151620009b581620006fc565b60608201526080850151620009ca81620006fc565b60808201529250620009e08560a0860162000726565b6101208501519092506001600160401b03811115620009fe57600080fd5b62000a0c868287016200079f565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b600181811c9082168062000a4157607f821691505b60208210810362000a6257634e487b7160e01b600052602260045260246000fd5b50919050565b600080835462000a788162000a2c565b6001828116801562000a93576001811462000aa95762000ada565b60ff198416875282151583028701945062000ada565b8760005260208060002060005b8581101562000ad15781548a82015290840190820162000ab6565b50505082870194505b50929695505050505050565b601f82111562000b36576000816000526020600020601f850160051c8101602086101562000b115750805b601f850160051c820191505b8181101562000b325782815560010162000b1d565b5050505b505050565b81516001600160401b0381111562000b575762000b5762000646565b62000b6f8162000b68845462000a2c565b8462000ae6565b602080601f83116001811462000ba7576000841562000b8e5750858301515b600019600386901b1c1916600185901b17855562000b32565b600085815260208120601f198616915b8281101562000bd85788860151825594840194600190910190840162000bb7565b508582101562000bf75787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b602080825282546001600160a01b0381168383015260a081901c60ff161515604084015260a81c6001600160401b0316606083015260808083015260018084018054600093929190849062000c5c8162000a2c565b8060a089015260c0600183166000811462000c80576001811462000c9d5762000ccf565b60ff19841660c08b015260c083151560051b8b0101945062000ccf565b85600052602060002060005b8481101562000cc65781548c820185015290880190890162000ca9565b8b0160c0019550505b50929998505050505050505050565b60805160a05160c05160e0516101005161012051615f8b62000d7b600039600081816101b001528181610ce801528181612ed1015261380b0152600081816102380152612aa10152600081816102090152612d490152600081816101da01528181610fcc0152818161117c01526124a10152600081816101810152818161264c015261270301526000818161195b015261198e0152615f8b6000f3fe608060405234801561001057600080fd5b506004361061012c5760003560e01c80637edf52f4116100ad578063de5e0b9a11610071578063de5e0b9a146104eb578063e9d68a8e146104fe578063f2fde38b1461051e578063f58e03fc14610531578063f716f99f1461054457600080fd5b80637edf52f41461044b57806385572ffb1461045e5780638da5cb5b1461046c578063c673e58414610487578063ccd37ba3146104a757600080fd5b80635e36480c116100f45780635e36480c146103405780635e7bb0081461036057806360987c20146103735780637437ff9f1461038657806379ba50971461044357600080fd5b806304666f9c1461013157806306285c6914610146578063181f5a77146102c65780633f4b04aa1461030f5780635215505b1461032a575b600080fd5b61014461013f366004613eaf565b610557565b005b6102686040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526040518060a001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160401b031681526020017f000000000000000000000000000000000000000000000000000000000000000061ffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b6040805182516001600160401b0316815260208084015161ffff1690820152828201516001600160a01b03908116928201929092526060808401518316908201526080928301519091169181019190915260a0015b60405180910390f35b6103026040518060400160405280601181526020017f4f666652616d7020312e362e302d64657600000000000000000000000000000081525081565b6040516102bd919061401d565b600b546040516001600160401b0390911681526020016102bd565b61033261056b565b6040516102bd929190614077565b61035361034e366004614118565b6107c6565b6040516102bd9190614175565b61014461036e3660046146de565b61081b565b61014461038136600461496d565b610aaf565b6103fc60408051608081018252600080825260208201819052918101829052606081019190915250604080516080810182526004546001600160a01b038082168352600160a01b820463ffffffff166020840152600160c01b90910460ff16151592820192909252600554909116606082015290565b6040516102bd919081516001600160a01b03908116825260208084015163ffffffff1690830152604080840151151590830152606092830151169181019190915260800190565b610144610d8a565b610144610459366004614a01565b610e0d565b61014461012c366004614a66565b6001546040516001600160a01b0390911681526020016102bd565b61049a610495366004614ab1565b610e1e565b6040516102bd9190614b11565b6104dd6104b5366004614b86565b6001600160401b03919091166000908152600a60209081526040808320938352929052205490565b6040519081526020016102bd565b6101446104f9366004614c02565b610f7c565b61051161050c366004614cb4565b61147f565b6040516102bd9190614ccf565b61014461052c366004614ce2565b61158b565b61014461053f366004614cff565b61159c565b610144610552366004614dba565b611605565b61055f611647565b61056881611674565b50565b606080600061057a60066118fd565b6001600160401b0381111561059157610591613ccf565b6040519080825280602002602001820160405280156105e257816020015b60408051608081018252600080825260208083018290529282015260608082015282526000199092019101816105af5790505b50905060006105f160066118fd565b6001600160401b0381111561060857610608613ccf565b604051908082528060200260200182016040528015610631578160200160208202803683370190505b50905060005b61064160066118fd565b8110156107bd57610653600682611907565b82828151811061066557610665614ef7565b60200260200101906001600160401b031690816001600160401b0316815250506008600083838151811061069b5761069b614ef7565b6020908102919091018101516001600160401b039081168352828201939093526040918201600020825160808101845281546001600160a01b038116825260ff600160a01b820416151593820193909352600160a81b9092049093169181019190915260018201805491929160608401919061071690614f0d565b80601f016020809104026020016040519081016040528092919081815260200182805461074290614f0d565b801561078f5780601f106107645761010080835404028352916020019161078f565b820191906000526020600020905b81548152906001019060200180831161077257829003601f168201915b5050505050815250508382815181106107aa576107aa614ef7565b6020908102919091010152600101610637565b50939092509050565b60006107d460016004614f5d565b60026107e1608085614f86565b6001600160401b03166107f49190614fac565b6107fe8585611913565b901c1660038111156108125761081261414b565b90505b92915050565b610823611958565b815181518114610846576040516320f8fd5960e21b815260040160405180910390fd5b60005b81811015610a9f57600084828151811061086557610865614ef7565b6020026020010151905060008160200151519050600085848151811061088d5761088d614ef7565b60200260200101519050805182146108b8576040516320f8fd5960e21b815260040160405180910390fd5b60005b82811015610a905760008282815181106108d7576108d7614ef7565b60200260200101516000015190506000856020015183815181106108fd576108fd614ef7565b6020026020010151905081600014610956578060800151821015610956578551815151604051633a98d46360e11b81526001600160401b0390921660048301526024820152604481018390526064015b60405180910390fd5b83838151811061096857610968614ef7565b602002602001015160200151518160a0015151146109b557805180516060909101516040516370a193fd60e01b815260048101929092526001600160401b0316602482015260440161094d565b60005b8160a0015151811015610a825760008585815181106109d9576109d9614ef7565b60200260200101516020015182815181106109f6576109f6614ef7565b602002602001015163ffffffff16905080600014610a795760008360a001518381518110610a2657610a26614ef7565b60200260200101516040015163ffffffff16905080821015610a77578351516040516348e617b360e01b8152600481019190915260248101849052604481018290526064810183905260840161094d565b505b506001016109b8565b5050508060010190506108bb565b50505050806001019050610849565b50610aaa83836119c0565b505050565b333014610acf576040516306e34e6560e31b815260040160405180910390fd5b6040805160008082526020820190925281610b0c565b6040805180820190915260008082526020820152815260200190600190039081610ae55790505b5060a08701515190915015610b4257610b3f8660a001518760200151886060015189600001516020015189898989611a83565b90505b6040805160a081018252875151815287516020908101516001600160401b03168183015288015181830152908701516060820152608081018290526005546001600160a01b03168015610c35576040516308d450a160e01b81526001600160a01b038216906308d450a190610bbb908590600401615070565b600060405180830381600087803b158015610bd557600080fd5b505af1925050508015610be6575060015b610c35573d808015610c14576040519150601f19603f3d011682016040523d82523d6000602084013e610c19565b606091505b50806040516309c2532560e01b815260040161094d919061401d565b604088015151158015610c4a57506080880151155b80610c61575060608801516001600160a01b03163b155b80610c8857506060880151610c86906001600160a01b03166385572ffb60e01b611c34565b155b15610c9557505050610d83565b87516020908101516001600160401b03166000908152600890915260408082205460808b015160608c01519251633cf9798360e01b815284936001600160a01b0390931692633cf9798392610d119289927f00000000000000000000000000000000000000000000000000000000000000009291600401615083565b6000604051808303816000875af1158015610d30573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610d5891908101906150bf565b509150915081610d7d57806040516302a35ba360e21b815260040161094d919061401d565b50505050505b5050505050565b6000546001600160a01b03163314610db55760405163015aa1e360e11b815260040160405180910390fd5b600180546001600160a01b0319808216339081179093556000805490911681556040516001600160a01b03909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610e15611647565b61056881611c50565b610e616040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c082015294855291820180548451818402810184019095528085529293858301939092830182828015610f0a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610eec575b5050505050815260200160038201805480602002602001604051908101604052809291908181526020018280548015610f6c57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610f4e575b5050505050815250509050919050565b6000610f8a8789018961536c565b6004805491925090600160c01b900460ff1661103457602082015151156110345760208201516040808401519051633854844f60e11b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016926370a9089e926110039230929190600401615594565b60006040518083038186803b15801561101b57600080fd5b505afa15801561102f573d6000803e3d6000fd5b505050505b8151515115158061104a57508151602001515115155b1561111557600b5460208b0135906001600160401b03808316911610156110ed57600b805467ffffffffffffffff19166001600160401b03831617905581548351604051633937306f60e01b81526001600160a01b0390921691633937306f916110b6916004016156a7565b600060405180830381600087803b1580156110d057600080fd5b505af11580156110e4573d6000803e3d6000fd5b50505050611113565b82602001515160000361111357604051632261116760e01b815260040160405180910390fd5b505b60005b8260200151518110156113cb5760008360200151828151811061113d5761113d614ef7565b60209081029190910101518051604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b166004820152919250906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632cbc26bb90602401602060405180830381865afa1580156111c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111e791906156ba565b1561121057604051637edeb53960e11b81526001600160401b038216600482015260240161094d565b600061121b82611d55565b90508060010160405161122e91906156d7565b60405180910390208360200151805190602001201461126b5782602001518160010160405163b80d8fa960e01b815260040161094d9291906157ca565b60408301518154600160a81b90046001600160401b0390811691161415806112ac575082606001516001600160401b031683604001516001600160401b0316115b156112f157825160408085015160608601519151636af0786b60e11b81526001600160401b03938416600482015290831660248201529116604482015260640161094d565b6080830151806113145760405163504570e360e01b815260040160405180910390fd5b83516001600160401b03166000908152600a602090815260408083208484529091529020541561136c5783516040516332cf0cbf60e01b81526001600160401b0390911660048201526024810182905260440161094d565b606084015161137c9060016157ef565b825467ffffffffffffffff60a81b1916600160a81b6001600160401b0392831602179092559251166000908152600a602090815260408083209483529390529190912042905550600101611118565b50602082015182516040517f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e492611403929091615816565b60405180910390a1610d7d60008b8b8b8b8b8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808f0282810182019093528e82529093508e92508d9182918501908490808284376000920191909152508c9250611da1915050565b60408051608080820183526000808352602080840182905283850182905260608085018190526001600160401b03878116845260088352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b90920490921694830194909452600184018054939492939184019161150b90614f0d565b80601f016020809104026020016040519081016040528092919081815260200182805461153790614f0d565b8015610f6c5780601f1061155957610100808354040283529160200191610f6c565b820191906000526020600020905b81548152906001019060200180831161156757505050919092525091949350505050565b611593611647565b6105688161209a565b6115dc6115ab8284018461583b565b60408051600080825260208201909252906115d6565b60608152602001906001900390816115c15790505b506119c0565b6040805160008082526020820190925290506115ff600185858585866000611da1565b50505050565b61160d611647565b60005b81518110156116435761163b82828151811061162e5761162e614ef7565b6020026020010151612113565b600101611610565b5050565b6001546001600160a01b03163314611672576040516315ae3a6f60e11b815260040160405180910390fd5b565b60005b815181101561164357600082828151811061169457611694614ef7565b60200260200101519050600081602001519050806001600160401b03166000036116d15760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b03166116f9576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b0381166000908152600860205260409020606083015160018201805461172590614f0d565b905060000361178757815467ffffffffffffffff60a81b1916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a16117f0565b8154600160a81b90046001600160401b03166001148015906117c75750805160208201206040516117bc9060018501906156d7565b604051809103902014155b156117f057604051632105803760e11b81526001600160401b038416600482015260240161094d565b805115806118255750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b15611843576040516342bcdf7f60e11b815260040160405180910390fd5b6001820161185182826158bf565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b029290921674ffffffffffffffffffffffffffffffffffffffffff19909116171782556118ac60066001600160401b03851661243d565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b836040516118e6919061597e565b60405180910390a250505050806001019050611677565b6000610815825490565b60006108128383612449565b6001600160401b0382166000908152600960205260408120816119376080856159cc565b6001600160401b031681526020810191909152604001600020549392505050565b467f00000000000000000000000000000000000000000000000000000000000000001461167257604051630f01ce8560e01b81527f0000000000000000000000000000000000000000000000000000000000000000600482015246602482015260440161094d565b81516000036119e25760405163c2e5347d60e01b815260040160405180910390fd5b80516040805160008082526020820190925291159181611a25565b6040805180820190915260008152606060208201528152602001906001900390816119fd5790505b50905060005b8451811015610d8357611a7b858281518110611a4957611a49614ef7565b602002602001015184611a7557858381518110611a6857611a68614ef7565b6020026020010151612473565b83612473565b600101611a2b565b606088516001600160401b03811115611a9e57611a9e613ccf565b604051908082528060200260200182016040528015611ae357816020015b6040805180820190915260008082526020820152815260200190600190039081611abc5790505b509050811560005b8a51811015611c265781611b8357848482818110611b0b57611b0b614ef7565b9050602002016020810190611b2091906159f2565b63ffffffff1615611b8357848482818110611b3d57611b3d614ef7565b9050602002016020810190611b5291906159f2565b8b8281518110611b6457611b64614ef7565b60200260200101516040019063ffffffff16908163ffffffff16815250505b611c018b8281518110611b9857611b98614ef7565b60200260200101518b8b8b8b8b87818110611bb557611bb5614ef7565b9050602002810190611bc79190615a0d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612d0e92505050565b838281518110611c1357611c13614ef7565b6020908102919091010152600101611aeb565b505098975050505050505050565b6000611c3f8361300e565b801561081257506108128383613041565b80516001600160a01b0316611c78576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516001600160a01b039889167fffffffffffffffff0000000000000000000000000000000000000000000000009097168717600160a01b63ffffffff958616021760ff60c01b1916600160c01b911515919091021790965560608089018051600580546001600160a01b031916918b169190911790558251968752935190921693850193909352935115159183019190915251909216908201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9060800160405180910390a150565b6001600160401b03811660009081526008602052604081208054600160a01b900460ff166108155760405163ed053c5960e01b81526001600160401b038416600482015260240161094d565b60ff87811660009081526002602090815260408083208151608081018352815481526001909101548086169382019390935261010083048516918101919091526201000090910490921615156060830152873590611e00876084615a53565b9050826060015115611e48578451611e19906020614fac565b8651611e26906020614fac565b611e319060a0615a53565b611e3b9190615a53565b611e459082615a53565b90505b368114611e7157604051638e1192e160e01b81526004810182905236602482015260440161094d565b5081518114611ea05781516040516324f7d61360e21b815260048101919091526024810182905260440161094d565b611ea8611958565b60ff808a1660009081526003602090815260408083203384528252808320815180830190925280548086168352939491939092840191610100909104166002811115611ef657611ef661414b565b6002811115611f0757611f0761414b565b9052509050600281602001516002811115611f2457611f2461414b565b148015611f785750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff1681548110611f6057611f60614ef7565b6000918252602090912001546001600160a01b031633145b611f9557604051631b41e11d60e31b815260040160405180910390fd5b50816060015115612045576020820151611fb0906001615a66565b60ff16855114611fd3576040516371253a2560e01b815260040160405180910390fd5b8351855114611ff55760405163a75d88af60e01b815260040160405180910390fd5b60008787604051612007929190615a7f565b60405190819003812061201e918b90602001615a8f565b6040516020818303038152906040528051906020012090506120438a828888886130cb565b505b6040805182815260208a8101356001600160401b03169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b336001600160a01b038216036120c357604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff1660000361213e576000604051631b3fab5160e11b815260040161094d9190615aa3565b60208082015160ff8082166000908152600290935260408320600181015492939092839216900361218f576060840151600182018054911515620100000262ff0000199092169190911790556121cb565b6060840151600182015460ff62010000909104161515901515146121cb576040516321fd80df60e21b815260ff8416600482015260240161094d565b60a0840151805161010010156121f7576001604051631b3fab5160e11b815260040161094d9190615aa3565b805160000361221c576005604051631b3fab5160e11b815260040161094d9190615aa3565b612282848460030180548060200260200160405190810160405280929190818152602001828054801561227857602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161225a575b505050505061327e565b8460600151156123b2576122f08484600201805480602002602001604051908101604052809291908181526020018280548015612278576020028201919060005260206000209081546001600160a01b0316815260019091019060200180831161225a57505050505061327e565b60808501518051610100101561231c576002604051631b3fab5160e11b815260040161094d9190615aa3565b604086015161232c906003615abd565b60ff16815111612352576003604051631b3fab5160e11b815260040161094d9190615aa3565b815181511015612378576001604051631b3fab5160e11b815260040161094d9190615aa3565b805160018401805461ff00191661010060ff8416021790556123a39060028601906020840190613c55565b506123b0858260016132e7565b505b6123be848260026132e7565b80516123d39060038501906020840190613c55565b5060408581015160018401805460ff191660ff8316179055865180855560a088015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f5479361242c9389939260028a01929190615ad9565b60405180910390a1610d8384613442565b600061081283836134c5565b600082600001828154811061246057612460614ef7565b9060005260206000200154905092915050565b81518151604051632cbc26bb60e01b8152608083901b67ffffffffffffffff60801b166004820152901515907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa1580156124f0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061251491906156ba565b1561258557801561254357604051637edeb53960e11b81526001600160401b038316600482015260240161094d565b6040516001600160401b03831681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d8949339060200160405180910390a150505050565b60208401515160008190036125bb57845160405163676cf24b60e11b81526001600160401b03909116600482015260240161094d565b84604001515181146125e0576040516357e0e08360e01b815260040160405180910390fd5b6000816001600160401b038111156125fa576125fa613ccf565b604051908082528060200260200182016040528015612623578160200160208202803683370190505b50905060007f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f857f000000000000000000000000000000000000000000000000000000000000000061267488611d55565b60010160405161268491906156d7565b6040519081900381206126bc949392916020019384526001600160401b03928316602085015291166040830152606082015260800190565b60405160208183030381529060405280519060200120905060005b838110156127f2576000886020015182815181106126f7576126f7614ef7565b602002602001015190507f00000000000000000000000000000000000000000000000000000000000000006001600160401b03168160000151604001516001600160401b03161461276e5780516040908101519051631c21951160e11b81526001600160401b03909116600482015260240161094d565b866001600160401b03168160000151602001516001600160401b0316146127c257805160200151604051636c95f1eb60e01b81526001600160401b03808a166004830152909116602482015260440161094d565b6127cc8184613514565b8483815181106127de576127de614ef7565b6020908102919091010152506001016126d7565b5050600061280a858389606001518a6080015161361c565b90508060000361283857604051633ee8bd3f60e11b81526001600160401b038616600482015260240161094d565b60005b83811015612d045760005a905060008960200151838151811061286057612860614ef7565b60200260200101519050600061287e898360000151606001516107c6565b905060008160038111156128945761289461414b565b14806128b1575060038160038111156128af576128af61414b565b145b61290757815160600151604080516001600160401b03808d16825290921660208301527f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c910160405180910390a1505050612cfc565b606088156129e6578a858151811061292157612921614ef7565b6020908102919091018101510151600454909150600090600160a01b900463ffffffff1661294f8842614f5d565b119050808061296f5750600383600381111561296d5761296d61414b565b145b612997576040516354e7e43160e11b81526001600160401b038c16600482015260240161094d565b8b86815181106129a9576129a9614ef7565b6020026020010151600001516000146129e0578b86815181106129ce576129ce614ef7565b60209081029190910101515160808501525b50612a52565b60008260038111156129fa576129fa61414b565b14612a5257825160600151604080516001600160401b03808e16825290921660208301527f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120910160405180910390a150505050612cfc565b8251608001516001600160401b031615612b28576000826003811115612a7a57612a7a61414b565b03612b285782516080015160208401516040516370701e5760e11b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263e0e03cae92612ad8928f929190600401615b8b565b6020604051808303816000875af1158015612af7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b1b91906156ba565b612b285750505050612cfc565b60008c604001518681518110612b4057612b40614ef7565b6020026020010151905080518460a001515114612b8a57835160600151604051631cfe6d8b60e01b81526001600160401b03808e166004830152909116602482015260440161094d565b612b9e8b8560000151606001516001613659565b600080612bac8684866136fe565b91509150612bc38d87600001516060015184613659565b8b15612c1a576003826003811115612bdd57612bdd61414b565b03612c1a576000856003811115612bf657612bf661414b565b14612c1a57855151604051632b11b8d960e01b815261094d91908390600401615bb7565b6002826003811115612c2e57612c2e61414b565b14612c6f576003826003811115612c4757612c4761414b565b14612c6f578551606001516040516349362d1f60e11b815261094d918f918590600401615bd0565b8560000151600001518660000151606001516001600160401b03168e6001600160401b03167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b8d8c81518110612cc757612cc7614ef7565b602002602001015186865a612cdc908f614f5d565b604051612cec9493929190615bf5565b60405180910390a4505050505050505b60010161283b565b5050505050505050565b6040805180820190915260008082526020820152602086015160405163bbe4f6db60e01b81526001600160a01b0380831660048301526000917f00000000000000000000000000000000000000000000000000000000000000009091169063bbe4f6db90602401602060405180830381865afa158015612d92573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612db69190615c2c565b90506001600160a01b0381161580612de55750612de36001600160a01b03821663aff2afbf60e01b611c34565b155b15612e0e5760405163ae9b4ce960e01b81526001600160a01b038216600482015260240161094d565b600080612e2688858c6040015163ffffffff166137b2565b915091506000806000612ef76040518061010001604052808e81526020018c6001600160401b031681526020018d6001600160a01b031681526020018f608001518152602001896001600160a01b031681526020018f6000015181526020018f6060015181526020018b815250604051602401612ea39190615c49565b60408051601f198184030181529190526020810180516001600160e01b0316633907753760e01b17905287867f000000000000000000000000000000000000000000000000000000000000000060846138b5565b92509250925082612f1f578582604051634ff17cad60e11b815260040161094d929190615d15565b8151602014612f4e578151604051631e3be00960e21b815260206004820152602481019190915260440161094d565b600082806020019051810190612f649190615d37565b9050866001600160a01b03168c6001600160a01b031614612fe0576000612f958d8a612f90868a614f5d565b6137b2565b50905086811080612faf575081612fac8883614f5d565b14155b15612fde5760405163a966e21f60e01b815260048101839052602481018890526044810182905260640161094d565b505b604080518082019091526001600160a01b039098168852602088015250949550505050505095945050505050565b6000613021826301ffc9a760e01b613041565b8015610815575061303a826001600160e01b0319613041565b1592915050565b6040516001600160e01b031982166024820152600090819060440160408051601f19818403018152919052602080820180516001600160e01b03166301ffc9a760e01b178152825192935060009283928392909183918a617530fa92503d915060005190508280156130b4575060208210155b80156130c05750600081115b979650505050505050565b8251600090815b81811015612d045760006001888684602081106130f1576130f1614ef7565b6130fe91901a601b615a66565b89858151811061311057613110614ef7565b602002602001015189868151811061312a5761312a614ef7565b602002602001015160405160008152602001604052604051613168949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561318a573d6000803e3d6000fd5b505060408051601f1981015160ff808e166000908152600360209081528582206001600160a01b038516835281528582208587019096528554808416865293975090955092939284019161010090041660028111156131eb576131eb61414b565b60028111156131fc576131fc61414b565b90525090506001816020015160028111156132195761321961414b565b1461323757604051636518c33d60e11b815260040160405180910390fd5b8051600160ff9091161b85161561326157604051633d9ef1f160e21b815260040160405180910390fd5b806000015160ff166001901b8517945050508060010190506130d2565b60005b8151811015610aaa5760ff8316600090815260036020526040812083519091908490849081106132b3576132b3614ef7565b6020908102919091018101516001600160a01b03168252810191909152604001600020805461ffff19169055600101613281565b60005b82518110156115ff57600083828151811061330757613307614ef7565b60200260200101519050600060028111156133245761332461414b565b60ff80871660009081526003602090815260408083206001600160a01b038716845290915290205461010090041660028111156133635761336361414b565b14613384576004604051631b3fab5160e11b815260040161094d9190615aa3565b6001600160a01b0381166133ab5760405163d6c62c9b60e01b815260040160405180910390fd5b60405180604001604052808360ff1681526020018460028111156133d1576133d161414b565b905260ff80871660009081526003602090815260408083206001600160a01b0387168452825290912083518154931660ff198416811782559184015190929091839161ffff19161761010083600281111561342e5761342e61414b565b0217905550905050508060010190506132ea565b60ff8181166000818152600260205260409020600101546201000090049091169061349a5780613485576040516317bd8dd160e11b815260040160405180910390fd5b600b805467ffffffffffffffff191690555050565b60001960ff831601611643578015611643576040516307b8c74d60e51b815260040160405180910390fd5b600081815260018301602052604081205461350c57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610815565b506000610815565b81518051606080850151908301516080808701519401516040516000958695889561357895919490939192916020019485526001600160a01b039390931660208501526001600160401b039182166040850152606084015216608082015260a00190565b604051602081830303815290604052805190602001208560200151805190602001208660400151805190602001208760a001516040516020016135bb9190615df1565b60408051601f198184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e0015b60405160208183030381529060405280519060200120905092915050565b60008061362a85858561398f565b6001600160401b0387166000908152600a6020908152604080832093835292905220549150505b949350505050565b60006002613668608085614f86565b6001600160401b031661367b9190614fac565b905060006136898585611913565b90508161369860016004614f5d565b901b1916818360038111156136af576136af61414b565b6001600160401b03871660009081526009602052604081209190921b929092179182916136dd6080886159cc565b6001600160401b031681526020810191909152604001600020555050505050565b604051630304c3e160e51b815260009060609030906360987c209061372b90889088908890600401615e88565b600060405180830381600087803b15801561374557600080fd5b505af1925050508015613756575060015b613795573d808015613784576040519150601f19603f3d011682016040523d82523d6000602084013e613789565b606091505b506003925090506137aa565b50506040805160208101909152600081526002905b935093915050565b6000806000806000613831886040516024016137dd91906001600160a01b0391909116815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166370a0823160e01b17905288887f000000000000000000000000000000000000000000000000000000000000000060846138b5565b92509250925082613859578682604051634ff17cad60e11b815260040161094d929190615d15565b6020825114613888578151604051631e3be00960e21b815260206004820152602481019190915260440161094d565b8180602001905181019061389c9190615d37565b6138a68288614f5d565b94509450505050935093915050565b6000606060008361ffff166001600160401b038111156138d7576138d7613ccf565b6040519080825280601f01601f191660200182016040528015613901576020820181803683370190505b509150863b61391b5763030ed58f60e21b60005260046000fd5b5a8581101561393557632be8ca8b60e21b60005260046000fd5b8590036040810481038710613955576337c3be2960e01b60005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d848111156139785750835b808352806000602085013e50955095509592505050565b82518251600091908183036139b757604051630469ac9960e21b815260040160405180910390fd5b61010182118015906139cb57506101018111155b6139e8576040516309bde33960e01b815260040160405180910390fd5b60001982820101610100811115613a12576040516309bde33960e01b815260040160405180910390fd5b80600003613a3f5786600081518110613a2d57613a2d614ef7565b60200260200101519350505050613c0d565b6000816001600160401b03811115613a5957613a59613ccf565b604051908082528060200260200182016040528015613a82578160200160208202803683370190505b50905060008080805b85811015613bac5760006001821b8b811603613ae65788851015613acf578c5160018601958e918110613ac057613ac0614ef7565b60200260200101519050613b08565b8551600185019487918110613ac057613ac0614ef7565b8b5160018401938d918110613afd57613afd614ef7565b602002602001015190505b600089861015613b38578d5160018701968f918110613b2957613b29614ef7565b60200260200101519050613b5a565b8651600186019588918110613b4f57613b4f614ef7565b602002602001015190505b82851115613b7b576040516309bde33960e01b815260040160405180910390fd5b613b858282613c14565b878481518110613b9757613b97614ef7565b60209081029190910101525050600101613a8b565b506001850382148015613bbe57508683145b8015613bc957508581145b613be6576040516309bde33960e01b815260040160405180910390fd5b836001860381518110613bfb57613bfb614ef7565b60200260200101519750505050505050505b9392505050565b6000818310613c2c57613c278284613c32565b610812565b61081283835b6040805160016020820152908101839052606081018290526000906080016135fe565b828054828255906000526020600020908101928215613caa579160200282015b82811115613caa57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613c75565b50613cb6929150613cba565b5090565b5b80821115613cb65760008155600101613cbb565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715613d0757613d07613ccf565b60405290565b60405160a081016001600160401b0381118282101715613d0757613d07613ccf565b60405160c081016001600160401b0381118282101715613d0757613d07613ccf565b604080519081016001600160401b0381118282101715613d0757613d07613ccf565b604051606081016001600160401b0381118282101715613d0757613d07613ccf565b604051601f8201601f191681016001600160401b0381118282101715613dbd57613dbd613ccf565b604052919050565b60006001600160401b03821115613dde57613dde613ccf565b5060051b60200190565b6001600160a01b038116811461056857600080fd5b80356001600160401b0381168114613e1457600080fd5b919050565b801515811461056857600080fd5b8035613e1481613e19565b60006001600160401b03821115613e4b57613e4b613ccf565b50601f01601f191660200190565b600082601f830112613e6a57600080fd5b8135613e7d613e7882613e32565b613d95565b818152846020838601011115613e9257600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215613ec257600080fd5b82356001600160401b0380821115613ed957600080fd5b818501915085601f830112613eed57600080fd5b8135613efb613e7882613dc5565b81815260059190911b83018401908481019088831115613f1a57600080fd5b8585015b83811015613fc057803585811115613f365760008081fd5b86016080818c03601f1901811315613f4e5760008081fd5b613f56613ce5565b89830135613f6381613de8565b81526040613f72848201613dfd565b8b830152606080850135613f8581613e19565b83830152928401359289841115613f9e57600091508182fd5b613fac8f8d86880101613e59565b908301525085525050918601918601613f1e565b5098975050505050505050565b60005b83811015613fe8578181015183820152602001613fd0565b50506000910152565b60008151808452614009816020860160208601613fcd565b601f01601f19169290920160200192915050565b6020815260006108126020830184613ff1565b6001600160a01b0381511682526020810151151560208301526001600160401b03604082015116604083015260006060820151608060608501526136516080850182613ff1565b604080825283519082018190526000906020906060840190828701845b828110156140b95781516001600160401b031684529284019290840190600101614094565b50505083810382850152845180825282820190600581901b8301840187850160005b8381101561410957601f198684030185526140f7838351614030565b948701949250908601906001016140db565b50909998505050505050505050565b6000806040838503121561412b57600080fd5b61413483613dfd565b915061414260208401613dfd565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b600481106141715761417161414b565b9052565b602081016108158284614161565b600060a0828403121561419557600080fd5b61419d613d0d565b9050813581526141af60208301613dfd565b60208201526141c060408301613dfd565b60408201526141d160608301613dfd565b60608201526141e260808301613dfd565b608082015292915050565b8035613e1481613de8565b803563ffffffff81168114613e1457600080fd5b600082601f83011261421d57600080fd5b8135602061422d613e7883613dc5565b82815260059290921b8401810191818101908684111561424c57600080fd5b8286015b8481101561431c5780356001600160401b03808211156142705760008081fd5b9088019060a0828b03601f190181131561428a5760008081fd5b614292613d0d565b87840135838111156142a45760008081fd5b6142b28d8a83880101613e59565b8252506040808501356142c481613de8565b828a015260606142d58682016141f8565b828401526080915081860135858111156142ef5760008081fd5b6142fd8f8c838a0101613e59565b9184019190915250919093013590830152508352918301918301614250565b509695505050505050565b6000610140828403121561433a57600080fd5b614342613d2f565b905061434e8383614183565b815260a08201356001600160401b038082111561436a57600080fd5b61437685838601613e59565b602084015260c084013591508082111561438f57600080fd5b61439b85838601613e59565b60408401526143ac60e085016141ed565b606084015261010084013560808401526101208401359150808211156143d157600080fd5b506143de8482850161420c565b60a08301525092915050565b600082601f8301126143fb57600080fd5b8135602061440b613e7883613dc5565b82815260059290921b8401810191818101908684111561442a57600080fd5b8286015b8481101561431c5780356001600160401b0381111561444d5760008081fd5b61445b8986838b0101614327565b84525091830191830161442e565b600082601f83011261447a57600080fd5b8135602061448a613e7883613dc5565b82815260059290921b840181019181810190868411156144a957600080fd5b8286015b8481101561431c5780356001600160401b03808211156144cc57600080fd5b818901915089603f8301126144e057600080fd5b858201356144f0613e7882613dc5565b81815260059190911b830160400190878101908c83111561451057600080fd5b604085015b838110156145495780358581111561452c57600080fd5b61453b8f6040838a0101613e59565b845250918901918901614515565b508752505050928401925083016144ad565b600082601f83011261456c57600080fd5b8135602061457c613e7883613dc5565b8083825260208201915060208460051b87010193508684111561459e57600080fd5b602086015b8481101561431c57803583529183019183016145a3565b600082601f8301126145cb57600080fd5b813560206145db613e7883613dc5565b82815260059290921b840181019181810190868411156145fa57600080fd5b8286015b8481101561431c5780356001600160401b038082111561461e5760008081fd5b9088019060a0828b03601f19018113156146385760008081fd5b614640613d0d565b61464b888501613dfd565b8152604080850135848111156146615760008081fd5b61466f8e8b838901016143ea565b8a84015250606080860135858111156146885760008081fd5b6146968f8c838a0101614469565b83850152506080915081860135858111156146b15760008081fd5b6146bf8f8c838a010161455b565b91840191909152509190930135908301525083529183019183016145fe565b600080604083850312156146f157600080fd5b6001600160401b038335111561470657600080fd5b61471384843585016145ba565b91506001600160401b036020840135111561472d57600080fd5b6020830135830184601f82011261474357600080fd5b614750613e788235613dc5565b81358082526020808301929160051b84010187101561476e57600080fd5b602083015b6020843560051b850101811015614914576001600160401b038135111561479957600080fd5b87603f8235860101126147ab57600080fd5b6147be613e786020833587010135613dc5565b81358501602081810135808452908301929160059190911b016040018a10156147e657600080fd5b604083358701015b83358701602081013560051b01604001811015614904576001600160401b038135111561481a57600080fd5b833587018135016040818d03603f1901121561483557600080fd5b61483d613d51565b604082013581526001600160401b036060830135111561485c57600080fd5b8c605f60608401358401011261487157600080fd5b6040606083013583010135614888613e7882613dc5565b808282526020820191508f60608460051b60608801358801010111156148ad57600080fd5b6060808601358601015b60608460051b6060880135880101018110156148e4576148d6816141f8565b8352602092830192016148b7565b5080602085015250505080855250506020830192506020810190506147ee565b5084525060209283019201614773565b508093505050509250929050565b60008083601f84011261493457600080fd5b5081356001600160401b0381111561494b57600080fd5b6020830191508360208260051b850101111561496657600080fd5b9250929050565b60008060008060006060868803121561498557600080fd5b85356001600160401b038082111561499c57600080fd5b6149a889838a01614327565b965060208801359150808211156149be57600080fd5b6149ca89838a01614922565b909650945060408801359150808211156149e357600080fd5b506149f088828901614922565b969995985093965092949392505050565b600060808284031215614a1357600080fd5b614a1b613ce5565b8235614a2681613de8565b8152614a34602084016141f8565b60208201526040830135614a4781613e19565b60408201526060830135614a5a81613de8565b60608201529392505050565b600060208284031215614a7857600080fd5b81356001600160401b03811115614a8e57600080fd5b820160a08185031215613c0d57600080fd5b803560ff81168114613e1457600080fd5b600060208284031215614ac357600080fd5b61081282614aa0565b60008151808452602080850194506020840160005b83811015614b065781516001600160a01b031687529582019590820190600101614ae1565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a0840152614b6060e0840182614acc565b90506040840151601f198483030160c0850152614b7d8282614acc565b95945050505050565b60008060408385031215614b9957600080fd5b614ba283613dfd565b946020939093013593505050565b806040810183101561081557600080fd5b60008083601f840112614bd357600080fd5b5081356001600160401b03811115614bea57600080fd5b60208301915083602082850101111561496657600080fd5b60008060008060008060008060c0898b031215614c1e57600080fd5b614c288a8a614bb0565b975060408901356001600160401b0380821115614c4457600080fd5b614c508c838d01614bc1565b909950975060608b0135915080821115614c6957600080fd5b614c758c838d01614922565b909750955060808b0135915080821115614c8e57600080fd5b50614c9b8b828c01614922565b999c989b50969995989497949560a00135949350505050565b600060208284031215614cc657600080fd5b61081282613dfd565b6020815260006108126020830184614030565b600060208284031215614cf457600080fd5b8135613c0d81613de8565b600080600060608486031215614d1457600080fd5b614d1e8585614bb0565b925060408401356001600160401b03811115614d3957600080fd5b614d4586828701614bc1565b9497909650939450505050565b600082601f830112614d6357600080fd5b81356020614d73613e7883613dc5565b8083825260208201915060208460051b870101935086841115614d9557600080fd5b602086015b8481101561431c578035614dad81613de8565b8352918301918301614d9a565b60006020808385031215614dcd57600080fd5b82356001600160401b0380821115614de457600080fd5b818501915085601f830112614df857600080fd5b8135614e06613e7882613dc5565b81815260059190911b83018401908481019088831115614e2557600080fd5b8585015b83811015613fc057803585811115614e4057600080fd5b860160c0818c03601f19011215614e575760008081fd5b614e5f613d2f565b8882013581526040614e72818401614aa0565b8a8301526060614e83818501614aa0565b8284015260809150614e96828501613e27565b9083015260a08381013589811115614eae5760008081fd5b614ebc8f8d83880101614d52565b838501525060c0840135915088821115614ed65760008081fd5b614ee48e8c84870101614d52565b9083015250845250918601918601614e29565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680614f2157607f821691505b602082108103614f4157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8181038181111561081557610815614f47565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b0380841680614fa057614fa0614f70565b92169190910692915050565b808202811582820484141761081557610815614f47565b80518252600060206001600160401b0381840151168185015260408084015160a06040870152614ff660a0870182613ff1565b90506060850151868203606088015261500f8282613ff1565b608087810151898303918a01919091528051808352908601935060009250908501905b8083101561506457835180516001600160a01b0316835286015186830152928501926001929092019190840190615032565b50979650505050505050565b6020815260006108126020830184614fc3565b6080815260006150966080830187614fc3565b61ffff9590951660208301525060408101929092526001600160a01b0316606090910152919050565b6000806000606084860312156150d457600080fd5b83516150df81613e19565b60208501519093506001600160401b038111156150fb57600080fd5b8401601f8101861361510c57600080fd5b805161511a613e7882613e32565b81815287602083850101111561512f57600080fd5b615140826020830160208601613fcd565b809450505050604084015190509250925092565b80356001600160e01b0381168114613e1457600080fd5b600082601f83011261517c57600080fd5b8135602061518c613e7883613dc5565b82815260069290921b840181019181810190868411156151ab57600080fd5b8286015b8481101561431c57604081890312156151c85760008081fd5b6151d0613d51565b6151d982613dfd565b81526151e6858301615154565b818601528352918301916040016151af565b600082601f83011261520957600080fd5b81356020615219613e7883613dc5565b82815260059290921b8401810191818101908684111561523857600080fd5b8286015b8481101561431c5780356001600160401b038082111561525c5760008081fd5b9088019060a0828b03601f19018113156152765760008081fd5b61527e613d0d565b615289888501613dfd565b81526040808501358481111561529f5760008081fd5b6152ad8e8b83890101613e59565b8a84015250606093506152c1848601613dfd565b9082015260806152d2858201613dfd565b9382019390935292013590820152835291830191830161523c565b600082601f8301126152fe57600080fd5b8135602061530e613e7883613dc5565b82815260069290921b8401810191818101908684111561532d57600080fd5b8286015b8481101561431c576040818903121561534a5760008081fd5b615352613d51565b813581528482013585820152835291830191604001615331565b6000602080838503121561537f57600080fd5b82356001600160401b038082111561539657600080fd5b90840190606082870312156153aa57600080fd5b6153b2613d73565b8235828111156153c157600080fd5b830160408189038113156153d457600080fd5b6153dc613d51565b8235858111156153eb57600080fd5b8301601f81018b136153fc57600080fd5b803561540a613e7882613dc5565b81815260069190911b8201890190898101908d83111561542957600080fd5b928a01925b828410156154795785848f0312156154465760008081fd5b61544e613d51565b843561545981613de8565b8152615466858d01615154565b818d0152825292850192908a019061542e565b84525050508287013591508482111561549157600080fd5b61549d8a83850161516b565b818801528352505082840135828111156154b657600080fd5b6154c2888286016151f8565b858301525060408301359350818411156154db57600080fd5b6154e7878585016152ed565b60408201529695505050505050565b600082825180855260208086019550808260051b84010181860160005b8481101561558757601f19868403018952815160a06001600160401b0380835116865286830151828888015261554b83880182613ff1565b60408581015184169089015260608086015190931692880192909252506080928301519290950191909152509783019790830190600101615513565b5090979650505050505050565b6001600160a01b0384168152600060206060818401526155b760608401866154f6565b83810360408581019190915285518083528387019284019060005b81811015614109578451805184528601518684015293850193918301916001016155d2565b805160408084528151848201819052600092602091908201906060870190855b8181101561564e57835180516001600160a01b031684528501516001600160e01b0316858401529284019291850191600101615617565b50508583015187820388850152805180835290840192506000918401905b8083101561506457835180516001600160401b031683528501516001600160e01b03168583015292840192600192909201919085019061566c565b60208152600061081260208301846155f7565b6000602082840312156156cc57600080fd5b8151613c0d81613e19565b60008083546156e581614f0d565b600182811680156156fd576001811461571257615741565b60ff1984168752821515830287019450615741565b8760005260208060002060005b858110156157385781548a82015290840190820161571f565b50505082870194505b50929695505050505050565b6000815461575a81614f0d565b8085526020600183811680156157775760018114615791576157bf565b60ff1985168884015283151560051b8801830195506157bf565b866000528260002060005b858110156157b75781548a820186015290830190840161579c565b890184019650505b505050505092915050565b6040815260006157dd6040830185613ff1565b8281036020840152614b7d818561574d565b6001600160401b0381811683821601908082111561580f5761580f614f47565b5092915050565b60408152600061582960408301856154f6565b8281036020840152614b7d81856155f7565b60006020828403121561584d57600080fd5b81356001600160401b0381111561586357600080fd5b613651848285016145ba565b601f821115610aaa576000816000526020600020601f850160051c810160208610156158985750805b601f850160051c820191505b818110156158b7578281556001016158a4565b505050505050565b81516001600160401b038111156158d8576158d8613ccf565b6158ec816158e68454614f0d565b8461586f565b602080601f83116001811461592157600084156159095750858301515b600019600386901b1c1916600185901b1785556158b7565b600085815260208120601f198616915b8281101561595057888601518255948401946001909101908401615931565b508582101561596e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60208152600082546001600160a01b038116602084015260ff8160a01c16151560408401526001600160401b038160a81c1660608401525060808083015261081260a083016001850161574d565b60006001600160401b03808416806159e6576159e6614f70565b92169190910492915050565b600060208284031215615a0457600080fd5b610812826141f8565b6000808335601e19843603018112615a2457600080fd5b8301803591506001600160401b03821115615a3e57600080fd5b60200191503681900382131561496657600080fd5b8082018082111561081557610815614f47565b60ff818116838216019081111561081557610815614f47565b8183823760009101908152919050565b828152604082602083013760600192915050565b6020810160068310615ab757615ab761414b565b91905290565b60ff818116838216029081169081811461580f5761580f614f47565b600060a0820160ff881683526020878185015260a0604085015281875480845260c0860191508860005282600020935060005b81811015615b315784546001600160a01b031683526001948501949284019201615b0c565b50508481036060860152865180825290820192508187019060005b81811015615b715782516001600160a01b031685529383019391830191600101615b4c565b50505060ff851660808501525090505b9695505050505050565b60006001600160401b03808616835280851660208401525060606040830152614b7d6060830184613ff1565b8281526040602082015260006136516040830184613ff1565b6001600160401b03848116825283166020820152606081016136516040830184614161565b848152615c056020820185614161565b608060408201526000615c1b6080830185613ff1565b905082606083015295945050505050565b600060208284031215615c3e57600080fd5b8151613c0d81613de8565b6020815260008251610100806020850152615c68610120850183613ff1565b91506020850151615c8460408601826001600160401b03169052565b5060408501516001600160a01b038116606086015250606085015160808501526080850151615cbe60a08601826001600160a01b03169052565b5060a0850151601f19808685030160c0870152615cdb8483613ff1565b935060c08701519150808685030160e0870152615cf88483613ff1565b935060e0870151915080868503018387015250615b818382613ff1565b6001600160a01b03831681526040602082015260006136516040830184613ff1565b600060208284031215615d4957600080fd5b5051919050565b600082825180855260208086019550808260051b84010181860160005b8481101561558757601f19868403018952815160a08151818652615d9382870182613ff1565b9150506001600160a01b03868301511686860152604063ffffffff8184015116818701525060608083015186830382880152615dcf8382613ff1565b6080948501519790940196909652505098840198925090830190600101615d6d565b6020815260006108126020830184615d50565b60008282518085526020808601955060208260051b8401016020860160005b8481101561558757601f19868403018952615e3f838351613ff1565b98840198925090830190600101615e23565b60008151808452602080850194506020840160005b83811015614b0657815163ffffffff1687529582019590820190600101615e66565b60608152600084518051606084015260208101516001600160401b0380821660808601528060408401511660a08601528060608401511660c08601528060808401511660e0860152505050602085015161014080610100850152615ef06101a0850183613ff1565b91506040870151605f198086850301610120870152615f0f8483613ff1565b935060608901519150615f2c838701836001600160a01b03169052565b608089015161016087015260a0890151925080868503016101808701525050615f558282615d50565b9150508281036020840152615f6a8186615e04565b90508281036040840152615b818185615e5156fea164736f6c6343000818000a",
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

func (_OffRamp *OffRampTransactor) Commit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "commit", reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampSession) Commit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactorSession) Commit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactor) Execute(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "execute", reportContext, report)
}

func (_OffRamp *OffRampSession) Execute(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Execute(&_OffRamp.TransactOpts, reportContext, report)
}

func (_OffRamp *OffRampTransactorSession) Execute(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
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

	Commit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Execute(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error)

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

var OffRampZKBin = ("0x0004000000000002002b000000000002000000600310027000000c6d0030019d00000c6d033001970003000000310355001d00000001035300020000000103550000000100200190001c00000003001d000000290000c13d0000008001000039001500000001001d000000400010043f000000040030008c0000004c0000413d0000001d0100035f000000000101043b000000e00210027000000c910020009c000000580000a13d00000c920020009c0000001c01000029000000e50000213d00000c990020009c0000077b0000a13d00000c9a0020009c00000b2e0000613d00000c9b0020009c00000a760000613d00000c9c0020009c0000004c0000c13d0000000001000416000000000001004b0000004c0000c13d0000000101000039000000000101041a00000c7201100197000000800010043f00000caf01000041000031b20001042e0000012003000039000000400030043f0000000001000416000000000001004b0000004c0000c13d0000001c080000290000001f0180003900000c6e011001970000012001100039000000400010043f0000001f0480018f00000c6f0580019800000120025000390000001d0700035f0000003d0000613d000000000607034f000000006106043c0000000003130436000000000023004b000000390000c13d000000000004004b0000004a0000613d000000000157034f0000000303400210000000000402043300000000043401cf000000000434022f000000000101043b0000010003300089000000000131022f00000000013101cf000000000141019f0000000000120435000001200080008c0000004e0000813d0000000001000019000031b300010430000000400100043d001d00000001001d00000c700010009c0000011a0000a13d00000d0a01000041000000000010043f0000004101000039000000040010043f00000c8501000041000031b30001043000000c9f0020009c0000001c01000029000001820000a13d00000ca00020009c0000076e0000a13d00000ca10020009c00000a250000613d00000ca20020009c000009e70000613d00000ca30020009c0000004c0000c13d0000001c0100002900120044001000940000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000201043b00000c710020009c0000004c0000213d00000023012000390000001c0010006c0000004c0000813d0000000403200039000c00000003001d0000001d0130035f000000000301043b00000c710030009c000000520000213d00000005043002100000003f0140003900000c730110019700000c700010009c000000520000213d0000008001100039000000400010043f000000800030043f000000240120003900000000020100190000000004140019000b00000004001d0000001c0040006c0000004c0000213d000000000003004b0000001d0700035f00000fb80000c13d0000002401700370000000000101043b001b00000001001d00000c710010009c0000004c0000213d0000001b0100002900000023011000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000300001900000c750300404100000c750010009c000000000302c019000000000003004b0000004c0000c13d0000001b0100002900000004011000390000001d0110035f000000000501043b00000c710050009c000000520000213d00000005045002100000003f0140003900000c7301100197000000400300043d0000000002130019001900000003001d000000000032004b0000000001000039000000010100403900000c710020009c000000520000213d0000000100100190000000520000c13d000000400020043f00000019010000290000000001510436001800000001001d0000001b0100002900000024031000390000000002340019001a00000002001d0000001c0020006c0000004c0000213d000000000005004b0000001d0100035f0000167d0000c13d00000cc9010000410000000000100443000000000100041200000004001004430000002400000443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000025790000613d000000000101043b001d00000001001d00000c77010000410000000000100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c78011001c70000800b0200003931b131ac0000040f0000000100200190000025790000613d000000000101043b0000001d0010006b0000257a0000c13d00000019010000290000000002010433000000800100043d000000000021004b000000180f000029000018580000c13d000000000001004b000018090000c13d0000008001000039000000190200002931b129970000040f0000000001000019000031b20001042e00000c930020009c000007a70000a13d00000c940020009c00000b4b0000613d00000c950020009c00000aa40000613d00000c960020009c0000004c0000c13d0000001c01000029000000240210008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000301043b00000c710030009c0000004c0000213d00000023013000390000001c0010006c0000004c0000813d00000004013000390000001d0110035f000000000601043b00000c710060009c000000520000213d00000005056002100000003f0150003900000c730110019700000c700010009c000000520000213d0000008001100039000000400010043f000000800060043f000000240430003900000000054500190000001c0050006c0000004c0000213d000000000006004b0000001d0100035f00000ec60000c13d0000000101000039000000000101041a00000c72011001970000000002000411000000000012004b00000e0a0000c13d000000800100043d000000000001004b00000b490000613d001100000000001d000001e80000013d0000001d010000290000008001100039000000400010043f000001200100043d00000c710010009c0000004c0000213d0000001d020000290000000001120436001b00000001001d000001400100043d00000c720010009c0000004c0000213d0000001b020000290000000000120435000001600100043d00000c720010009c0000004c0000213d0000001d020000290000004002200039001a00000002001d0000000000120435000001800100043d00000c720010009c0000004c0000213d0000001d020000290000006002200039001800000002001d0000000000120435000000400100043d001900000001001d00000c700010009c000000520000213d00000019010000290000008001100039000000400010043f000001a00100043d00000c720010009c0000004c0000213d00000019020000290000000001120436001700000001001d000001c00100043d00000c6d0010009c0000004c0000213d00000017020000290000000000120435000001e00100043d000000000001004b0000000002000039000000010200c039000000000021004b0000004c0000c13d00000019020000290000004002200039001600000002001d0000000000120435000002000100043d00000c720010009c0000004c0000213d00000019020000290000006002200039001500000002001d0000000000120435000002200400043d00000c710040009c0000004c0000213d0000001c0100002900000120021000390000013f01400039000000000021004b0000004c0000813d0000012003400039000000000603043300000c710060009c000000520000213d00000005056002100000003f0150003900000c7301100197000000400800043d0000000007180019001400000008001d000000000087004b0000000001000039000000010100403900000c710070009c000000520000213d0000000100100190000000520000c13d000000400070043f00000014010000290000000001610436001300000001001d00000140044000390000000005450019000000000025004b0000004c0000213d000000000006004b000016ff0000c13d0000000001000411000000000001004b000017600000c13d000000400100043d00000c8f0200004100000e0c0000013d00000ca60020009c000004c00000213d00000ca90020009c000007c80000613d00000caa0020009c0000004c0000c13d0000000001000416000000000001004b0000004c0000c13d31b1293c0000040f000000400100043d001d00000001001d31b1290b0000040f0000000001000412002b00000001001d002a00200000003d0000800501000039000000440300003900000000040004150000002b0440008a000000050440021000000cc90200004131b131890000040f00000c7102100197001b00000002001d0000001d010000290000000001210436001c00000001001d0000000001000412002900000001001d002800400000003d0000000004000415000000290440008a0000000504400210000080050100003900000cc902000041000000440300003931b131890000040f00000c72011001970000001c0200002900000000001204350000000001000412002700000001001d002600600000003d0000000004000415000000270440008a0000000504400210000080050100003900000cc902000041000000440300003931b131890000040f00000c72011001970000001d020000290000004002200039001a00000002001d00000000001204350000000001000412002500000001001d002400800000003d0000000004000415000000250440008a0000000504400210000080050100003900000cc902000041000000440300003931b131890000040f0000001d02000029000000600220003900000c72011001970000000000120435000000400100043d0000001b0300002900000000033104360000001c04000029000000000404043300000c720440019700000000004304350000001a03000029000000000303043300000c720330019700000040041000390000000000340435000000000202043300000c72022001970000006003100039000000000023043500000c6d0010009c00000c6d01008041000000400110021000000cf1011001c7000031b20001042e000000000001004b000018740000613d0000000b02000039000000000102041a00000cbb01100197000000000012041b0000001102000029001100010020003d000000800100043d000000110010006b00000b490000813d00000011010000290000000501100210000000a00110003900000000020104330000004001200039001000000001001d0000000001010433000000ff00100190000017f90000613d001200000002001d00000020012000390000000001010433000000ff0110018f001d00000001001d000000000010043f0000000201000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000120300002900000060053000390000000002050433000000000101043b001300000001001d0000000104100039000000000104041a000000ff00100190000002150000613d000000000002004b0000000002000039000000010200603900000cb2001001980000000001000039000000010100c039000000000121013f00000001001001900000021b0000c13d0000185b0000013d00000cb501100197000000000002004b00000cb6020000410000000002006019000000000112019f000000000014041b000000a001300039000d00000001001d00000000020104330000000031020434001700000003001d000001000010008c0000178c0000213d001900000005001d001800000002001d000f00000004001d000000400300043d000000000001004b000017ff0000613d00000013010000290000000301100039000000000201041a001b00000003001d001c00000002001d0000000002230436001a00000002001d000e00000001001d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001c05000029000000000005004b0000001a02000029000002480000613d000000000101043b0000001a020000290000000003000019000000000401041a00000c7204400197000000000242043600000001011000390000000103300039000000000053004b000002410000413d0000001b0300002900000000013200490000001f0110003900000d0e021001970000000001320019000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f0000000001030433000000000001004b000002850000613d0000000001000019001c00000001001d0000001d01000029000000000010043f0000000301000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001b0200002900000000020204330000001c03000029000000000032004b000025730000a13d00000005023002100000001a02200029000000000202043300000c7202200197000000000101043b000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000001041b0000001c0200002900000001022000390000001b010000290000000001010433000000000012004b0000000001020019000002590000413d00000019010000290000000001010433000000000001004b0000001301000029001400020010003d000003a20000613d0000001401000029000000000301041a000000400200043d001b00000002001d001c00000003001d0000000002320436001a00000002001d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001c05000029000000000005004b0000001a02000029000002aa0000613d000000000101043b0000001a020000290000000003000019000000000401041a00000c7204400197000000000242043600000001011000390000000103300039000000000053004b000002a30000413d0000001b0120006a0000001f0110003900000d0e021001970000001b01200029000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f0000001b010000290000000001010433000000000001004b000002e70000613d0000000001000019001c00000001001d0000001d01000029000000000010043f0000000301000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001b0200002900000000020204330000001c03000029000000000032004b000025730000a13d00000005023002100000001a02200029000000000202043300000c7202200197000000000101043b000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000001041b0000001c0200002900000001022000390000001b010000290000000001010433000000000012004b0000000001020019000002bb0000413d000000120100002900000080011000390000000001010433001600000001001d0000000014010434001500000001001d000001000040008c000018610000213d00000010010000290000000001010433000000fe0210018f000000550020008c0000000f03000029000000180200002900000b280000213d00000003011000c9000000ff0110018f000000000014004b000018670000a13d0000000001020433000000000014004b0000178c0000413d000000000103041a00000d0f0110019700000008024002100000ff000220018f000000000121019f000000000013041b0000001401000029000000000201041a000000000041041b001c00000004001d000000000024004b0000031e0000813d001b00000002001d0000001401000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000201043b0000001c012000290000001b02200029000000000021004b0000031e0000813d000000000001041b0000000101100039000000000021004b0000031a0000413d0000001401000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000200001900000016030000290000001c0600002900000000041200190000002003300039000000000503043300000c7205500197000000000054041b0000000102200039000000000062004b0000032d0000413d00000016010000290000000001010433000000000001004b000003a20000613d0000000002000019001b00000002001d000000050120021000000015011000290000000001010433001c00000001001d0000001d01000029000000000010043f0000000301000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001c0200002900000c7202200197000000000101043b001c00000002001d000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000101041a0000000801100270000000ff0110018f000000020010008c000017df0000213d000000000001004b0000166b0000c13d000000400300043d0000001c0000006b0000001b02000029000016760000613d00000c880030009c000000520000213d0000004001300039000000400010043f000000ff0120018f001a00000003001d00000000021304360000000101000039001900000002001d00000000001204350000001d01000029000000000010043f0000000301000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001c02000029000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000201041a00000d10022001970000001a030000290000000003030433000000ff0330018f000000000232019f000000000021041b00000019030000290000000003030433000000020030008c000017df0000213d00000d0f0220019700000008033002100000ff000330018f000000000223019f000000000021041b0000001b02000029000000010220003900000016010000290000000001010433000000000012004b0000033a0000413d00000018010000290000000001010433000000000001004b000004120000613d0000000002000019001b00000002001d000000050120021000000017011000290000000001010433001c00000001001d0000001d01000029000000000010043f0000000301000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001c0200002900000c7202200197000000000101043b001c00000002001d000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000101041a0000000801100270000000ff0110018f000000020010008c000017df0000213d000000000001004b0000166b0000c13d000000400300043d0000001c0000006b0000001b02000029000016760000613d00000c880030009c000000520000213d0000004001300039000000400010043f000000ff0120018f001a00000003001d00000000021304360000000201000039001900000002001d00000000001204350000001d01000029000000000010043f0000000301000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001c02000029000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000201041a00000d10022001970000001a030000290000000003030433000000ff0330018f000000000232019f000000000021041b00000019030000290000000003030433000000020030008c000017df0000213d00000d0f0220019700000008033002100000ff000330018f000000000223019f000000000021041b0000001b02000029000000010220003900000018010000290000000001010433000000000012004b000003a70000413d00000c710010009c000004130000a13d000000520000013d00000000010000190000000e03000029000000000203041a000000000013041b001c00000001001d000000000021004b0000042e0000813d001b00000002001d000000000030043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000201043b0000001c012000290000001b02200029000000000021004b0000000e030000290000042e0000813d000000000001041b0000000101100039000000000021004b0000042a0000413d000000000030043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001c06000029000000000006004b0000001805000029000004460000613d000000000200001900000000031200190000002005500039000000000405043300000c7204400197000000000043041b0000000102200039000000000062004b0000043e0000413d0000000f03000029000000000103041a00000d100110019700000010020000290000000002020433001c00ff002001930000001c011001af000000000013041b000000120100002900000000010104330000001302000029000000000012041b0000000d020000290000000002020433001b00000002001d000000400400043d0000004002400039000000a0030000390000000000320435000000200240003900000000001204350000001d010000290000000000140435001a00000004001d000000a0014000390000001402000029000000000302041a001900000003001d0000000000310435000000000020043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001a05000029000000c0025000390000001907000029000000000007004b0000047d0000613d000000000101043b00000000030000190000001b06000029000000000401041a00000c7204400197000000000242043600000001011000390000000103300039000000000073004b000004750000413d0000047e0000013d0000001b0600002900000000015200490000006003500039000000000013043500000000030604330000000001320436000000000003004b0000048d0000613d00000000020000190000002006600039000000000406043300000c720440019700000000014104360000000102200039000000000032004b000004860000413d00000080025000390000001c030000290000000000320435000000000151004900000c6d0010009c00000c6d01008041000000600110021000000c6d0050009c00000c6d050080410000004002500210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000121019f00000c83011001c70000800d02000039000000010300003900000cb90400004131b131a70000040f00000001002001900000004c0000613d0000001d01000029000000000010043f0000000201000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000000101100039000000000101041a0000001001100270000000ff0110018f0000001d02000029000000000002004b000001dd0000613d000000010020008c000001e30000c13d000000000001004b000001e30000613d000000400100043d00000cba0200004100000e0c0000013d00000ca70020009c000009d30000613d00000ca80020009c0000004c0000c13d0000001c01000029000000e40010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000006401100370000000000301043b00000c710030009c0000004c0000213d00000023013000390000001c0010006c0000004c0000813d00000004023000390000001d0120035f000000000101043b001b00000001001d00000c710010009c0000004c0000213d0000002403300039001900000003001d0000001b03300029001a00000003001d0000001c0030006c0000004c0000213d0000001d0100035f0000008401100370000000000301043b00000c710030009c0000004c0000213d00000023013000390000001c0010006c0000004c0000813d00000004013000390000001d0110035f000000000101043b001800000001001d00000c710010009c0000004c0000213d000000240130003900000018030000290000000503300210001700000001001d001500000003001d0000000003130019001600000003001d0000001c0030006c0000004c0000213d0000001d0100035f000000a401100370000000000301043b00000c710030009c0000004c0000213d00000023013000390000001c0010006c0000004c0000813d00000004013000390000001d0110035f000000000101043b001400000001001d00000c710010009c0000004c0000213d0000001b01000029000000200010008c0000004c0000413d000000240130003900000014030000290000000503300210001300000001001d001100000003001d0000000003130019001200000003001d0000001c0030006c0000004c0000213d00000020012000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d00000019011000290000001a0210006900000c740020009c0000004c0000213d000000600020008c0000004c0000413d000000e003000039000000400030043f0000001d0210035f000000000202043b00000c710020009c0000004c0000213d00000000021200190000001a0420006900000c740040009c0000004c0000213d000000400040008c0000004c0000413d0000012004000039000000400040043f0000001d0520035f000000000505043b00000c710050009c0000004c0000213d00000000052500190000001f065000390000001a0060006c0000004c0000813d0000001d0650035f000000000706043b00000c710070009c000000520000213d00000005067002100000003f0660003900000c730660019700000cd00060009c000000520000213d0000012006600039000000400060043f000001200070043f0000002005500039000000060670021000000000065600190000001a0060006c0000004c0000213d000000000007004b000019010000c13d000000e00040043f00000020042000390000001d0440035f000000000404043b00000c710040009c0000004c0000213d00000000052400190000001f025000390000001a06000029000000000062004b000000000400001900000c750400804100000c7502200197001c0c750060019b0000001c0620014f0000001c0020006c000000000200001900000c750200404100000c750060009c000000000204c019000000000002004b0000004c0000c13d0000001d0250035f000000000602043b00000c710060009c000000520000213d00000005026002100000003f0220003900000c7302200197000000400400043d0000000007240019000000000047004b0000000002000039000000010200403900000c710070009c000000520000213d0000000100200190000000520000c13d000000400070043f00000000006404350000002005500039000000060260021000000000065200190000001a0060006c0000004c0000213d000000000065004b0000058d0000813d00000000070400190000001a0250006900000c740020009c0000004c0000213d000000400020008c0000004c0000413d000000400800043d00000c880080009c000000520000213d0000004002800039000000400020043f0000001d0250035f000000000202043b00000c710020009c0000004c0000213d000000000228043600000020095000390000001d0990035f000000000909043b00000cf30090009c0000004c0000213d0000002007700039000000000092043500000000008704350000004005500039000000000065004b000005730000413d000001000040043f000000800030043f00000020021000390000001d0220035f000000000202043b00000c710020009c0000004c0000213d0000000002120019001000000002001d0000001f022000390000001a0020006c000000000300001900000c750300804100000c75022001970000001c0420014f0000001c0020006c000000000200001900000c750200404100000c750040009c000000000203c019000000000002004b0000004c0000c13d00000010030000290000001d0230035f000000000302043b00000c710030009c000000520000213d00000005043002100000003f0240003900000c7302200197000000400600043d0000000005260019000d00000006001d000000000065004b0000000002000039000000010200403900000c710050009c000000520000213d0000000100200190000000520000c13d000000400050043f0000000d020000290000000000320435000000100200002900000020052000390000000003540019000f00000003001d0000001a0030006c0000004c0000213d0000000f0050006c000006360000813d0000001a02000029000e0020002000920000000d090000290000001d0250035f000000000202043b00000c710020009c0000004c0000213d000000100d2000290000000e02d0006900000c740020009c0000004c0000213d000000a00020008c0000004c0000413d000000400a00043d00000cc500a0009c000000520000213d000000a002a00039000000400020043f0000002002d000390000001d0320035f000000000303043b00000c710030009c0000004c0000213d000000000c3a0436000000200b2000390000001d02b0035f000000000202043b00000c710020009c0000004c0000213d000000000fd200190000003f02f000390000001a0020006c000000000300001900000c750300804100000c75022001970000001c0420014f0000001c0020006c000000000200001900000c750200404100000c750040009c000000000203c019000000000002004b0000004c0000c13d0000002004f000390000001d0240035f000000000d02043b00000c7100d0009c000000520000213d0000001f02d0003900000d0e022001970000003f0220003900000d0e02200197000000400e00043d00000000032e00190000000000e3004b0000000002000039000000010200403900000c710030009c000000520000213d0000000100200190000000520000c13d0000004006f00039000000400030043f0000000002de043600000000036d00190000001a0030006c0000004c0000213d00000020034000390000001d0730035f00000d0e06d0019800000000046200190000060e0000613d000000000307034f000000000f020019000000003803043c000000000f8f043600000000004f004b0000060a0000c13d0000001f03d001900000061b0000613d000000000667034f0000000303300210000000000704043300000000073701cf000000000737022f000000000606043b0000010003300089000000000636022f00000000033601cf000000000373019f00000000003404350000000002d2001900000000000204350000000000ec04350000002002b000390000001d0320035f000000000303043b00000c710030009c0000004c0000213d0000004004a00039000000000034043500000020022000390000001d0320035f000000000303043b00000c710030009c0000004c0000213d00000020099000390000006004a00039000000000034043500000020022000390000001d0220035f000000000202043b0000008003a0003900000000002304350000000000a9043500000020055000390000000f0050006c000005c30000413d0000000d02000029000000a00020043f00000040021000390000001d0220035f000000000202043b00000c710020009c0000004c0000213d00000000031200190000001f013000390000001a0010006c000000000200001900000c750200804100000c75011001970000001c0410014f0000001c0010006c000000000100001900000c750100404100000c750040009c000000000102c019000000000001004b0000004c0000c13d0000001d0130035f000000000401043b00000c710040009c000000520000213d00000005014002100000003f0110003900000c7302100197000000400100043d0000000002210019000000000012004b0000000005000039000000010500403900000c710020009c000000520000213d0000000100500190000000520000c13d000000400020043f00000000004104350000002002300039000000060340021000000000032300190000001a0030006c0000004c0000213d000000000032004b0000067b0000813d00000000040100190000001a0520006900000c740050009c0000004c0000213d000000400050008c0000004c0000413d000000400500043d00000c880050009c000000520000213d00000020044000390000004006500039000000400060043f0000001d0620035f000000000606043b000000000665043600000020072000390000001d0770035f000000000707043b000000000076043500000000005404350000004002200039000000000032004b000006650000413d000000c00010043f0000000401000039000000000101041a00000cc300100198000006ef0000c13d000000a00100043d001d00000001001d0000000001010433000000000001004b000006ef0000613d00000cc90100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000c00200043d001c00000002001d00000cd102000041000000000020044300000c7201100197000f00000001001d0000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000000001004b0000004c0000613d000000400400043d00000024014000390000006002000039000000000021043500000cf4010000410000000000140435000000000100041000000c7202100197000000040140003900000000002104350000001d02000029000000000202043300000005032002100000000003340019001000000004001d000000640440003900000000002404350000008403300039000000000002004b000025e80000c13d00000000011300490000001002000029000000440220003900000000001204350000001c0100002900000000020104330000000001230436000000000002004b000006d20000613d00000000030000190000001c040000290000002004400039001c00000004001d0000000004040433000000005404043400000000044104360000000005050433000000000054043500000040011000390000000103300039000000000023004b000006c60000413d00000000020004140000000f03000029000000040030008c000006ea0000613d0000001003000029000000000131004900000c6d0010009c00000c6d01008041000000600110021000000c6d0030009c00000c6d030080410000004003300210000000000131019f00000c6d0020009c00000c6d02008041000000c002200210000000000121019f0000000f0200002931b131ac0000040f000000600310027000010c6d0030019d000300000001035500000001002001900000261d0000613d000000100100002900000c710010009c000000520000213d0000001001000029000000400010043f000000800100043d001c00000001001d0000000021010434001d00000002001d0000000001010433000000000001004b000006fb0000c13d0000001d0100002900000000010104330000000001010433000000000001004b000021d70000613d00000024010000390000000201100367000000000101043b00000c71011001970000000b02000039000000000302041a00000c7104300197000000000014004b000021d00000813d00000cbb03300197000000000113019f000000000012041b0000000401000039000000000101041a00000cd102000041000000000020044300000c7201100197001000000001001d0000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000000001004b0000004c0000613d000000400500043d00000cf60100004100000000001504350000000401500039000000200200003900000000002104350000001c010000290000000003010433000000240150003900000040020000390000000000210435000000640250003900000000040304330000000000420435001c00000005001d0000008402500039000000000004004b000007390000613d000000000500001900000020033000390000000006030433000000007606043400000c72066001970000000006620436000000000707043300000cf307700197000000000076043500000040022000390000000105500039000000000045004b0000072d0000413d00000000011200490000001d0300002900000000030304330000001c040000290000004404400039000000000014043500000000040304330000000001420436000000000004004b000007500000613d000000000200001900000020033000390000000005030433000000006505043400000c71055001970000000005510436000000000606043300000cf306600197000000000065043500000040011000390000000102200039000000000042004b000007440000413d00000000020004140000001003000029000000040030008c000007680000613d0000001c03000029000000000131004900000c6d0010009c00000c6d01008041000000600110021000000c6d0030009c00000c6d030080410000004003300210000000000131019f00000c6d0020009c00000c6d02008041000000c002200210000000000121019f000000100200002931b131a70000040f000000600310027000010c6d0030019d000300000001035500000001002001900000262a0000613d0000001c0100002900000c710010009c000000520000213d0000001c01000029000000400010043f000021d70000013d00000ca40020009c00000aba0000613d00000ca50020009c0000004c0000c13d0000000001000416000000000001004b0000004c0000c13d0000000b01000039000000000101041a00000c7101100197000000800010043f00000caf01000041000031b20001042e00000c9d0020009c00000b8c0000613d00000c9e0020009c0000004c0000c13d0000000001000416000000000001004b0000004c0000c13d31b1293c0000040f000000400100043d001d00000001001d31b1290b0000040f0000000401000039000000000101041a00000cc3001001980000000002000039000000010200c0390000001d0400002900000040034000390000000000230435000000a00210027000000c6d022001970000002003400039000000000023043500000c720110019700000000001404350000000501000039000000000101041a00000c7201100197000000600340003900000000020400190000000000130435000000400100043d001c00000001001d31b129280000040f0000001c02000029000000000121004900000c6d0010009c00000c6d01008041000000600110021000000c6d0020009c00000c6d020080410000004002200210000000000121019f000031b20001042e00000c970020009c00000cdd0000613d00000c980020009c0000004c0000c13d0000001c01000029000000440010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000101043b00000c710010009c0000004c0000213d000000000010043f0000000a01000039000000200010043f0000004002000039000000000100001931b131740000040f0000001d0200035f0000002402200370000000000202043b000000000020043f000000200010043f0000000001000019000000400200003931b131740000040f000000000101041a000000800010043f00000caf01000041000031b20001042e001b0024001000940000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000301043b00000c710030009c0000004c0000213d00000023013000390000001c0010006c0000004c0000813d00000004013000390000001d0110035f000000000601043b00000c710060009c000000520000213d00000005056002100000003f0150003900000c730110019700000c700010009c000000520000213d0000008001100039000000400010043f000000800060043f000000240430003900000000054500190000001c0050006c0000004c0000213d000000000006004b00000e120000c13d0000000101000039000000000101041a00000c72011001970000000002000411000000000012004b00000e0a0000c13d000000800100043d000000000001004b00000b490000613d001900000000001d00000019010000290000000501100210000000a001100039000000000201043300000020012000390000000001010433001c0c710010019c000018710000613d001800000002001d000000000102043300000c72001001980000177e0000613d0000001c01000029000000000010043f0000000801000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000401043b0000000101400039001d00000001001d000000000101041a000000010210019000000001071002700000007f0770618f0000001f0070008c00000000030000390000000103002039000000000331013f000000010030019000000b860000c13d000000180300002900000060033000390000000005030433000000000304041a000000000007004b001700200050003d001b00000004001d001a00000005001d000008430000613d00000c800330019700000c810030009c000008940000613d000000400300043d0000000005730436000000000002004b0000085a0000613d001400000005001d001500000003001d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c70000801002000039001600000007001d31b131ac0000040f000000160600002900000001002001900000004c0000613d000000000201043b000000000100001900000014050000290000000003510019000000000402041a000000000043043500000001022000390000002001100039000000000061004b0000083a0000413d00000015030000290000085d0000013d00000c860130019700000c81011001c7000000000014041b000000400100043d0000001c02000029000000000021043500000c6d0010009c00000c6d010080410000004001100210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c82011001c70000800d02000039000000010300003900000c870400004131b131a70000040f0000001a050000290000000100200190000008940000c13d0000004c0000013d00000d1001100197000000000015043500000020010000390000003f0110003900000d0e021001970000000001320019000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f00000c6d0050009c00000c6d050080410000004001500210000000000203043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f0000001a0300002900000001002001900000004c0000613d000000170200002900000c6d0020009c00000c6d020080410000004002200210000000000303043300000c6d0030009c00000c6d030080410000006003300210000000000223019f000000000101043b001600000001001d000000000100041400000c6d0010009c00000c6d01008041000000c001100210000000000121019f00000c83011001c7000080100200003931b131ac0000040f0000001a0500002900000001002001900000004c0000613d000000000101043b000000160010006b00001b600000c13d0000000001050433000000000001004b0000177e0000613d00000c6d0010009c00000c6d010080410000006001100210000000170200002900000c6d0020009c00000c6d020080410000004002200210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b001600000001001d000000400100043d00000020020000390000000002210436000000000002043500000c880010009c000000520000213d0000004003100039000000400030043f00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f0000001a0300002900000001002001900000004c0000613d000000000101043b000000160010006b0000177e0000613d000000000403043300000c710040009c000000520000213d0000001d01000029000000000101041a000000010010019000000001031002700000007f0330618f0000001f0030008c00000000020000390000000102002039000000000121013f00000001001001900000001b0600002900000b860000c13d000000200030008c001600000004001d000008f90000413d001500000003001d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d00000016040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b00000015010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000001b06000029000008f90000813d000000000002041b0000000102200039000000000012004b000008f50000413d000000200040008c000009250000413d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f0000001a0700002900000001002001900000004c0000613d000000160800002900000d0e02800198000000000101043b000009b10000613d000000010320008a00000005033002700000000003310019000000010430003900000020030000390000001b0600002900000000057300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000009110000c13d000000000082004b000009220000813d0000000302800210000000f80220018f00000d110220027f00000d110220016700000000037300190000000003030433000000000223016f000000000021041b000000010180021000000001011001bf000009310000013d000000000004004b000009300000613d000000030140021000000d110110027f00000d110110016700000017020000290000000002020433000000000112016f0000000102400210000000000121019f000009310000013d000000000100001900000018030000290000001d02000029000000000012041b00000040013000390000000001010433000000000001004b00000c89010000410000000001006019000000000206041a00000c8a02200197000000000112019f000000000203043300000c7202200197000000000121019f000000000016041b0000001c01000029000000000010043f0000000701000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000101041a000000000001004b0000096c0000c13d0000000601000039000000000101041a00000c710010009c000000520000213d00000001021000390000000603000039000000000023041b00000c8b0110009a0000001c02000029000000000021041b000000000103041a001a00000001001d000000000020043f0000000701000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001a02000029000000000021041b000000400600043d000000200100003900000000011604360000001b02000029000000000202041a00000080036000390000008004000039000000000043043500000c72032001970000000000310435000000a80120027000000c71011001970000006003600039000000000013043500000c8c002001980000000001000039000000010100c039000000400260003900000000001204350000001d01000029000000000101041a000000010210019000000001041002700000007f0440618f0000001f0040008c00000000030000390000000103002039000000000331013f000000010030019000000b860000c13d000000a0036000390000000000430435000000000002004b000009aa0000613d001a00000004001d001b00000006001d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001a07000029000000000007004b000009b60000613d0000001b06000029000000c002600039000000000301043b00000000010000190000000004210019000000000503041a000000000054043500000001033000390000002001100039000000000071004b000009a20000413d000009b80000013d00000d1001100197000000c0026000390000000000120435000000000004004b00000020010000390000000001006039000009b80000013d00000020030000390000001b06000029000000000082004b0000091a0000413d000009220000013d00000000010000190000001b0600002900000c6d0060009c00000c6d060080410000004002600210000000c00110003900000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000121019f00000c83011001c70000800d02000039000000020300003900000c8d040000410000001c0500002931b131a70000040f00000001002001900000004c0000613d0000001902000029001900010020003d000000800100043d000000190010006b000007f20000413d00000b490000013d0000000001000416000000000001004b0000004c0000c13d000000c001000039000000400010043f0000001101000039000000800010043f00000d0b01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e00200003931b129160000040f000000c00110008a00000c6d0010009c00000c6d01008041000000600110021000000d0c011001c7000031b20001042e000000440010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000101043b00000c710010009c0000004c0000213d0000001d0200035f0000002402200370000000000202043b001d00000002001d00000c710020009c0000004c0000213d000000000010043f0000000901000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001d020000290000000702200270000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001d030000290000000102300210000000000101043b000000000101041a0000007f0330019000000a1b0000613d000000ff0420018f00000000033400d9000000020030008c00000b280000c13d000000fe0220018f000000000121022f000000030110018f000000400200043d000000000012043500000c6d0020009c00000c6d02008041000000400120021000000cea011001c7000031b20001042e0000000001000416000000000001004b0000004c0000c13d0000000601000039000000000101041a001a00000001001d00000c710010009c000000520000213d0000001a0100002900000005021002100000003f0120003900000c730310019700000c700030009c000000520000213d0000008007300039000000400070043f0000001a01000029000000800010043f000000000001004b00000d3d0000c13d0000000003370019000000000073004b0000000001000039000000010100403900000c710030009c000000520000213d0000000100100190000000520000c13d000000400030043f0000001a010000290000000001170436001900000001001d0000001f0320018f000000000002004b00000a500000613d0000001c040000290000001d0140035f00000019040000290000000002240019000000001501043c0000000004540436000000000024004b00000a4c0000c13d000000000003004b0000001a0000006b00000d830000c13d000000400100043d000000400200003900000000032104360000000004070433000000400210003900000000004204350000006002100039000000000004004b00000a640000613d00000000050000190000002007700039000000000607043300000c710660019700000000026204360000000105500039000000000045004b00000a5d0000413d00000000041200490000000000430435000000800300043d0000000000320435000000050430021000000000044200190000002009400039000000000003004b00000f660000c13d000000000219004900000c6d0020009c00000c6d02008041000000600220021000000c6d0010009c00000c6d010080410000004001100210000000000112019f000031b20001042e000000840010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000010001000039000000400010043f0000001d0100035f0000000401100370000000000101043b00000c720010009c0000004c0000213d000000800010043f0000001d0200035f0000002402200370000000000202043b00000c6d0020009c0000004c0000213d000000a00020043f0000001d0300035f0000004403300370000000000303043b000000000003004b0000000004000039000000010400c039000000000043004b0000004c0000c13d000000c00030043f0000001d0400035f0000006404400370000000000404043b00000c720040009c0000004c0000213d000000e00040043f0000000105000039000000000505041a00000c72055001970000000006000411000000000056004b00000e840000c13d000000000001004b00000f990000c13d00000cc001000041000001000010043f00000cbd01000041000031b300010430000000240010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000601043b00000c720060009c0000004c0000213d0000000101000039000000000101041a00000c72011001970000000005000411000000000015004b00000d390000c13d000000000056004b00000d530000c13d00000cae01000041000000800010043f00000cac01000041000031b300010430000000840010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000006401100370000000000201043b00000c710020009c0000004c0000213d00000023012000390000001c0010006c0000004c0000813d00000004032000390000001d0130035f000000000101043b000f00000001001d00000c710010009c0000004c0000213d0000000f01000029000000200010008c0000004c0000413d0000002402200039000900000002001d0000000f02200029001b00000002001d0000001c0020006c0000004c0000213d00000020013000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d0000000901100029000b00000001001d0000001f011000390000001b0010006c0000004c0000813d0000000b020000290000001d0120035f000000000101043b00000c710010009c000000520000213d00000005031002100000003f0230003900000c730220019700000c700020009c000000520000213d0000008002200039000000400020043f000000800010043f0000000b04000029000000200540003900000000040500190000000005530019000a00000005001d0000001b0050006c0000004c0000213d000000000001004b0000134b0000c13d00000cc70020009c000000520000213d0000002001200039000000400010043f0000000000020435000000800100003931b129970000040f000000400100043d001d00000001001d00000cc70010009c000000520000213d0000001d010000290000002002100039001c00000002001d000000400020043f00000000000104350000000101000039000000000010043f0000000201000039000000200010043f000000400100043d00000c700010009c000000520000213d0000008002100039000000400020043f00000ceb02000041000000000202041a001a00000002001d000000000521043600000cec02000041000000000202041a0000000803200270000000ff0330018f00000040041000390000000000340435000000ff0320018f001800000005001d0000000000350435000000600410003900000cb2032001980000000001000039000000010100c039001900000004001d000000000014043500000004010000390000000201100367000000000101043b000000a50200008a0000000f0020006b000017810000a13d00000d0a01000041000000000010043f0000001101000039000000040010043f00000c8501000041000031b3000104300000000001000416000000000001004b0000004c0000c13d000000000100041a00000c72021001970000000006000411000000000026004b00000d350000c13d0000000102000039000000000302041a00000c7604300197000000000464019f000000000042041b00000c7601100197000000000010041b000000000100041400000c720530019700000c6d0010009c00000c6d01008041000000c00110021000000c83011001c70000800d02000039000000030300003900000cc20400004131b131a70000040f00000001002001900000004c0000613d0000000001000019000031b20001042e000000240010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000101043b00000c710010009c0000004c0000213d0000010002000039000000400020043f000000800000043f000000a00000043f000000c00000043f0000006002000039000000e00020043f000000000010043f0000000801000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000400200043d001d00000002001d00000c700020009c000000520000213d000000000101043b0000001d040000290000008002400039000000400020043f000000000201041a00000c72032001970000000006340436000000a80320027000000c71033001970000004005400039000000000035043500000c8c002001980000000002000039000000010200c03900000000002604350000000101100039000000000201041a000000010320019000000001072002700000007f0770618f0000001f0070008c00000000040000390000000104002039000000000442013f000000010040019000000d600000613d00000d0a01000041000000000010043f0000002201000039000000040010043f00000c8501000041000031b300010430000000640010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000101043b00000c710010009c0000004c0000213d00000004021000390000001c0120006900000c740010009c0000004c0000213d000001400010008c0000004c0000413d000001e003000039000000400030043f0000001d0120035f000000000101043b000001400010043f00000020012000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d000001600040043f00000020011000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d000001800040043f00000020011000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d000001a00040043f00000020011000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d000001c00040043f0000014004000039000000800040043f00000020041000390000001d0140035f000000000101043b00000c710010009c0000004c0000213d00000000062100190000001f016000390000001c0010006c0000004c0000813d0000001d0160035f000000000501043b00000c710050009c000000520000213d0000001f0150003900000d0e011001970000003f0110003900000d0e0710019700000cc40070009c000000520000213d0000002001600039000001e006700039000000400060043f000001e00050043f00000000061500190000001c0060006c0000004c0000213d0000001d0710035f00000d0e085001980000001f0950018f000002000680003900000be00000613d000002000a000039000000000b07034f00000000b10b043c000000000a1a043600000000006a004b00000bdc0000c13d000000000009004b00000bed0000613d000000000187034f0000000307900210000000000806043300000000087801cf000000000878022f000000000101043b0000010007700089000000000171022f00000000017101cf000000000181019f000000000016043500000200015000390000000000010435000000a00030043f00000020034000390000001d0130035f000000000101043b00000c710010009c0000004c0000213d00000000062100190000001f016000390000001c0010006c0000004c0000813d0000001d0160035f000000000401043b00000c710040009c000000520000213d0000001f0140003900000d0e011001970000003f0110003900000d0e01100197000000400500043d0000000007150019000000000057004b0000000001000039000000010100403900000c710070009c000000520000213d0000000100100190000000520000c13d0000002001600039000000400070043f000000000645043600000000071400190000001c0070006c0000004c0000213d0000001d0810035f00000d0e094001980000001f0a40018f000000000796001900000c1b0000613d000000000b08034f000000000106001900000000bc0b043c0000000001c10436000000000071004b00000c170000c13d00000000000a004b00000c280000613d000000000198034f0000000308a00210000000000907043300000000098901cf000000000989022f000000000101043b0000010008800089000000000181022f00000000018101cf000000000191019f000000000017043500000000014600190000000000010435000000c00050043f00000020033000390000001d0130035f000000000101043b00000c720010009c0000004c0000213d000000e00010043f00000020013000390000001d0110035f000000000101043b000001000010043f00000040013000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d0000000001210019001b00000001001d0000001f011000390000001c0010006c0000004c0000813d0000001b020000290000001d0120035f000000000601043b00000c710060009c000000520000213d00000005036002100000003f0130003900000c7301100197000000400400043d0000000002140019001900000004001d000000000042004b0000000001000039000000010100403900000c710020009c000000520000213d0000000100100190000000520000c13d000000400020043f000000190100002900000000006104350000001b0100002900000020041000390000000002430019001a00000002001d0000001c0020006c0000004c0000213d000000000006004b00001b660000c13d0000001901000029000001200010043f0000001d0100035f0000002401100370000000000101043b000c00000001001d00000c710010009c0000004c0000213d0000000c0100002900000023011000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000300001900000c750300404100000c750010009c000000000302c019000000000003004b0000004c0000c13d0000000c0100002900000004011000390000001d0110035f000000000101043b000700000001001d00000c710010009c0000004c0000213d0000000c01000029001000240010003d0000000701000029000000050110021000000010011000290000001c0010006c0000004c0000213d0000001d0100035f0000004401100370000000000201043b00000c710020009c0000004c0000213d00000023012000390000001c0010006c000000000300001900000c750300804100000c7501100197000000000001004b000000000400001900000c750400404100000c750010009c000000000403c019000000000004004b0000004c0000c13d00000004012000390000001d0110035f000000000101043b000b00000001001d00000c710010009c0000004c0000213d000300240020003d0000000b01000029000000050110021000000003011000290000001c0010006c0000004c0000213d0000000001000415000200000001001d000000400100043d000f00000001001d00000000010004100000000002000411000000000012004b00001cf60000c13d0000000f0100002900000cc70010009c000000520000213d0000000f020000290000002001200039000000400010043f0000000000020435000001200100043d000900000001001d0000000021010434000800000002001d000000000001004b00001cfe0000c13d000000400100043d001d00000001001d00000cc50010009c000000520000213d000000800100043d0000002002100039000000000202043300000000010104330000001d07000029000000a003700039000000a00400043d000000c00500043d000000400030043f00000080067000390000000f03000029001b00000006001d00000000003604350000006003700039001a00000003001d00000000005304350000004003700039001900000003001d0000000000430435000000000317043600000c7101200197001800000003001d00000000001304350000000501000039000000000101041a001c0c720010019c00001d340000c13d000000c00100043d0000000001010433000000000001004b00001db70000c13d0000000001000415000000200110008a001c000500100218000001000100043d000000000001004b002000000000003d002000010000603d00001dbb0000c13d000028e80000013d000000240010008c0000004c0000413d0000000001000416000000000001004b0000004c0000c13d0000001d0100035f0000000401100370000000000101043b000000ff0010008c0000004c0000213d0000016002000039000000400020043f000000e00000043f000001000000043f000001200000043f000001400000043f000000e002000039000000800020043f0000006002000039000000a00020043f000000c00020043f000000000010043f0000000201000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b001c00000001001d000000400100043d001d00000001001d00000cb00010009c000000520000213d0000001d020000290000006001200039000000400010043f00000cb10020009c000000520000213d0000001d05000029000000e002500039000000400020043f0000001c06000029000000000206041a00000000002104350000000102600039000000000202041a0000008003500039000000ff0420018f000000000043043500000cb2002001980000000003000039000000010300c039000000c00450003900000000003404350000000001150436001800000001001d000000a0015000390000000802200270000000ff0220018f00000000002104350000000201600039000000000301041a000000400200043d001b00000002001d001900000003001d0000000002320436001a00000002001d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001905000029000000000005004b000012c50000c13d0000001a04000029000012cf0000013d00000cc101000041000000800010043f00000cac01000041000031b30001043000000cab01000041000000800010043f00000cac01000041000031b30001043000000cd50030009c000000520000213d000000600400003900000000050000190000008001700039000000400010043f0000006001700039000000000041043500000040017000390000000000010435000000200170003900000000000104350000000000070435000000a0015000390000000000710435000000400700043d0000002005500039000000000025004b00000a390000813d00000c700070009c00000d410000a13d000000520000013d000000000100041a00000c7601100197000000000161019f000000000010041b000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c83011001c70000800d02000039000000030300003900000cad0400004100000b460000013d001b00000006001d001c00000005001d000000400500043d0000000004750436000000000003004b00000e880000613d001800000004001d001900000007001d001a00000005001d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001906000029000000000006004b00000000020000190000001a05000029000000180700002900000e8d0000613d000000000101043b00000000020000190000000003720019000000000401041a000000000043043500000001011000390000002002200039000000000062004b00000d7b0000413d00000e8d0000013d0000000003000019001b00000007001d0000000601000039000000000101041a000000000031004b000025730000a13d0000000601000039000000000010043f0000000001070433000000000031004b000025730000a13d00000c8b0130009a000000000101041a00000c71011001970000000504300210001c00000004001d000000190240002900000000001204350000000002070433000000000032004b000025730000a13d001d00000003001d000000000010043f0000000801000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000400500043d00000c700050009c0000001b07000029000000520000213d000000000101043b0000008002500039000000400020043f000000000201041a000000a80320027000000c71033001970000004004500039000000000034043500000c7203200197000000000335043600000c8c002001980000000002000039000000010200c03900000000002304350000000101100039000000000201041a000000010320019000000001082002700000007f0880618f0000001f0080008c00000000040000390000000104002039000000000442013f000000010040019000000b860000c13d000000400600043d0000000004860436000000000003004b00000de50000613d001500000004001d001600000008001d001700000006001d001800000005001d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001608000029000000000008004b00000deb0000613d000000000201043b00000000010000190000001b070000290000001805000029000000170600002900000015090000290000000003910019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00000ddd0000413d00000def0000013d00000d10012001970000000000140435000000000008004b0000002001000039000000000100603900000def0000013d00000000010000190000001b07000029000000180500002900000017060000290000003f0110003900000d0e021001970000000001620019000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f00000060015000390000000000610435000000800100043d0000001d03000029000000000031004b000025730000a13d0000001c01000029000000a0011000390000000000510435000000800100043d000000000031004b000025730000a13d00000001033000390000001a0030006c00000d850000413d00000a530000013d000000400100043d00000cab02000041000000000021043500000c6d0010009c00000c6d01008041000000400110021000000c90011001c7000031b300010430000000a00600003900000e1c0000013d0000000001910019000000000001043500000060018000390000000000a1043500000000068604360000002004400039000000000054004b000007e80000813d0000001d0140035f000000000101043b00000c710010009c0000004c0000213d00000000093100190000001b0190006900000c740010009c0000004c0000213d000000800010008c0000004c0000413d000000400800043d00000c700080009c000000520000213d0000008001800039000000400010043f00000024019000390000001d0a10035f000000000a0a043b00000c7200a0009c0000004c0000213d000000000aa8043600000020011000390000001d0b10035f000000000b0b043b00000c7100b0009c0000004c0000213d0000000000ba0435000000200a1000390000001d01a0035f000000000101043b000000000001004b000000000b000039000000010b00c0390000000000b1004b0000004c0000c13d000000400b80003900000000001b04350000002001a000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d000000000b9100190000004301b000390000001c0010006c000000000900001900000c750900804100000c7501100197000000000001004b000000000a00001900000c750a00404100000c750010009c000000000a09c01900000000000a004b0000004c0000c13d000000240cb000390000001d01c0035f000000000901043b00000c710090009c000000520000213d0000001f0190003900000d0e011001970000003f0110003900000d0e01100197000000400a00043d000000000d1a00190000000000ad004b0000000001000039000000010100403900000c7100d0009c000000520000213d0000000100100190000000520000c13d000000440bb000390000004000d0043f00000000019a0436000000000bb900190000001c00b0006c0000004c0000213d000000200bc000390000001d0db0035f00000d0e0e900198000000000ce1001900000e760000613d000000000f0d034f000000000b01001900000000f20f043c000000000b2b04360000000000cb004b00000e720000c13d0000001f0b90019000000e140000613d0000000002ed034f000000030bb00210000000000d0c0433000000000dbd01cf000000000dbd022f000000000202043b000001000bb000890000000002b2022f0000000002b201cf0000000002d2019f00000000002c043500000e140000013d00000cab01000041000001000010043f00000cbd01000041000031b30001043000000d10012001970000000000140435000000000007004b000000200200003900000000020060390000003f0220003900000d0e032001970000000002530019000000000032004b0000000003000039000000010300403900000c710020009c000000520000213d0000000100300190000000520000c13d0000000004050019000000400020043f0000001d05000029000000600350003900000000004304350000002004000039000000400200043d0000000004420436000000000505043300000c720550019700000000005404350000001b040000290000000004040433000000000004004b0000000004000039000000010400c039000000400520003900000000004504350000001c04000029000000000404043300000c7104400197000000600520003900000000004504350000000003030433000000800420003900000080050000390000000000540435000000a00420003900000000530304340000000000340435000000c004200039000000000003004b00000ec00000613d000000000600001900000000074600190000000008650019000000000808043300000000008704350000002006600039000000000036004b00000eb90000413d0000001f0530003900000d0e0150019700000000034300190000000000030435000000c0011000390000079f0000013d000000a00600003900000ecf0000013d000000a001700039000000000081043500000000067604360000002004400039000000000054004b0000001d0100035f0000010f0000813d000000000141034f000000000101043b00000c710010009c0000004c0000213d0000000008310019000000000182004900000c740010009c0000004c0000213d000000c00010008c0000004c0000413d000000400700043d00000cb30070009c000000520000213d000000c001700039000000400010043f00000024018000390000001d0110035f000000000101043b000000000917043600000044018000390000001d0a10035f000000000a0a043b000000ff00a0008c0000004c0000213d0000000000a9043500000020011000390000001d0910035f000000000909043b000000ff0090008c0000004c0000213d000000400a70003900000000009a043500000020091000390000001d0190035f000000000101043b000000000001004b000000000a000039000000010a00c0390000000000a1004b0000004c0000c13d000000600a70003900000000001a043500000020099000390000001d0190035f000000000101043b00000c710010009c0000004c0000213d000000000b8100190000004301b000390000001c0010006c000000000a00001900000c750a00804100000c7501100197000000000001004b000000000c00001900000c750c00404100000c750010009c000000000c0ac01900000000000c004b0000004c0000c13d0000002401b000390000001d0110035f000000000c01043b00000c7100c0009c000000520000213d000000050dc002100000003f01d0003900000c7301100197000000400a00043d000000000e1a00190000000000ae004b0000000001000039000000010100403900000c7100e0009c000000520000213d0000000100100190000000520000c13d0000004000e0043f0000000000ca0435000000440bb00039000000000cbd00190000001c00c0006c0000004c0000213d0000000000cb004b00000f2e0000813d000000000d0a00190000001d01b0035f000000000101043b00000c720010009c0000004c0000213d000000200dd0003900000000001d0435000000200bb000390000000000cb004b00000f250000413d00000080017000390000000000a1043500000020019000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d000000000981001900000043019000390000001c0010006c000000000800001900000c750800804100000c7501100197000000000001004b000000000a00001900000c750a00404100000c750010009c000000000a08c01900000000000a004b0000004c0000c13d00000024019000390000001d0110035f000000000a01043b00000c7100a0009c000000520000213d000000050ba002100000003f01b0003900000c7301100197000000400800043d000000000c18001900000000008c004b0000000001000039000000010100403900000c7100c0009c000000520000213d0000000100100190000000520000c13d0000004000c0043f0000000000a804350000004409900039000000000a9b00190000001c00a0006c0000004c0000213d0000000000a9004b00000ec80000813d000000000b0800190000001d0190035f000000000101043b00000c720010009c0000004c0000213d000000200bb0003900000000001b043500000020099000390000000000a9004b00000f5c0000413d00000ec80000013d000000800400003900000000060000190000000007020019000000000804001900000f730000013d000000000b9a001900000000000b04350000001f0aa0003900000d0e0aa0019700000000099a00190000000106600039000000000036004b00000a6d0000813d000000000a290049000000200aa0008a00000020077000390000000000a704350000002008800039000000000a08043300000000cb0a043400000c720bb00197000000000bb90436000000000c0c043300000000000c004b000000000c000039000000010c00c0390000000000cb0435000000400ba00039000000000b0b043300000c710bb00197000000400c9000390000000000bc0435000000600aa00039000000000a0a0433000000600b90003900000000004b0435000000800c90003900000000ba0a04340000000000ac0435000000a00990003900000000000a004b00000f6b0000613d000000000c000019000000000d9c0019000000000ecb0019000000000e0e04330000000000ed0435000000200cc000390000000000ac004b00000f910000413d00000f6b0000013d000000a00520021000000c7c0550019700000cbe06100197000000000565019f000000000003004b000000000600001900000c7d0600c041000000000565019f0000000406000039000000000706041a00000c7b07700197000000000575019f000000000056041b0000000505000039000000000605041a00000c7606600197000000000646019f000000000065041b000001000010043f000001200020043f000001400030043f000001600040043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cbf011001c70000800d02000039000000010300003900000c7e0400004100000b460000013d0000001c01000029001a00200010009200000fcb0000013d000000150400002900000020044000390000001005000029000000600150003900000000003104350000002001200039000000000117034f000000000101043b00000080025000390000000000120435001500000004001d00000000005404350000000d0200002900000020022000390000000b0020006c000000880000813d000d00000002001d000000000127034f000000000101043b00000c710010009c0000004c0000213d0000000c02100029001100000002001d0000001a0120006900000c740010009c0000004c0000213d000000a00010008c0000004c0000413d000000400100043d001000000001001d00000cc50010009c000000520000213d0000001001000029000000a001100039000000400010043f000000110100002900000020011000390000001d0210035f000000000202043b00000c710020009c0000004c0000213d00000010030000290000000002230436000e00000002001d00000020011000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d00000011021000290000003f012000390000001c0010006c000000000300001900000c750300804100000c7501100197000000000001004b000000000400001900000c750400404100000c750010009c000000000403c019000000000004004b0000004c0000c13d00000020012000390000001d0110035f000000000301043b00000c710030009c000000520000213d00000005043002100000003f0140003900000c7301100197000000400600043d0000000005160019000f00000006001d000000000065004b0000000001000039000000010100403900000c710050009c000000520000213d0000000100100190000000520000c13d000000400050043f0000000f010000290000000000310435001400400020003d0000001402400029001300000002001d0000001c0020006c0000004c0000213d0000001302000029000000140020006b0000001d0e00035f000011d80000813d00000014020000290000000f03000029000010250000013d000000160300002900000020033000390000001902000029000000a00120003900000018040000290000000000410435000000000023043500000017020000290000002002200039000000130020006c000011d80000813d001600000003001d001700000002001d00000000012e034f000000000101043b00000c710010009c0000004c0000213d00000014031000290000001c0130006900000c740010009c0000004c0000213d000001400010008c0000004c0000413d000000400100043d001900000001001d00000cb30010009c000000520000213d0000001901000029000000c002100039000000400020043f00000ce30010009c000000520000213d00000019010000290000016001100039000000400010043f0000001d0130035f000000000101043b000000000012043500000020013000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d0000001905000029000000e005500039000000000045043500000020011000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d00000019050000290000010005500039000000000045043500000020011000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d00000019050000290000012005500039000000000045043500000020011000390000001d0410035f000000000404043b00000c710040009c0000004c0000213d000000190500002900000000022504360000014005500039000000000045043500000020041000390000001d0140035f000000000101043b00000c710010009c0000004c0000213d000000000a3100190000001f01a000390000001c0010006c000000000500001900000c750500804100000c7501100197000000000001004b000000000600001900000c750600404100000c750010009c000000000605c019000000000006004b0000004c0000c13d0000001d01a0035f000000000801043b00000c710080009c000000520000213d0000001f0180003900000d0e011001970000003f0110003900000d0e01100197000000400900043d0000000005190019000000000095004b0000000001000039000000010100403900000c710050009c000000520000213d0000000100100190000000520000c13d0000002001a00039000000400050043f000000000a89043600000000051800190000001c0050006c0000004c0000213d0000001d0700035f000000000617034f00000d0e05800198000000000b5a0019000010950000613d000000000c06034f00000000010a001900000000cd0c043c0000000001d104360000000000b1004b000010910000c13d0000001f01800190000010a20000613d000000000556034f000000030110021000000000060b043300000000061601cf000000000616022f000000000505043b0000010001100089000000000515022f00000000011501cf000000000161019f00000000001b043500000000018a0019000000000001043500000000009204350000002002400039000000000127034f000000000101043b00000c710010009c0000004c0000213d00000000093100190000001f019000390000001c0010006c000000000400001900000c750400804100000c7501100197000000000001004b000000000500001900000c750500404100000c750010009c000000000504c019000000000005004b0000004c0000c13d0000001d0190035f000000000401043b00000c710040009c000000520000213d0000001f0140003900000d0e011001970000003f0110003900000d0e01100197000000400800043d0000000005180019000000000085004b0000000001000039000000010100403900000c710050009c000000520000213d0000000100100190000000520000c13d0000002001900039000000400050043f000000000948043600000000051400190000001c0050006c0000004c0000213d0000001d0700035f000000000617034f00000d0e05400198000000000a590019000010d90000613d000000000b06034f000000000109001900000000bc0b043c0000000001c104360000000000a1004b000010d50000c13d0000001f01400190000010e60000613d000000000556034f000000030110021000000000060a043300000000061601cf000000000616022f000000000505043b0000010001100089000000000515022f00000000011501cf000000000161019f00000000001a0435000000000149001900000000000104350000001901000029000000400110003900000000008104350000002002200039000000000127034f000000000101043b00000c720010009c0000004c0000213d00000019060000290000006004600039000000000014043500000020012000390000001d0110035f000000000101043b0000008004600039000000000014043500000040012000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d0000000001310019001b00000001001d0000001f011000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000300001900000c750300404100000c750010009c000000000302c019000000000003004b0000004c0000c13d0000001b020000290000001d0120035f000000000201043b00000c710020009c000000520000213d00000005042002100000003f0140003900000c7301100197000000400500043d0000000003150019001800000005001d000000000053004b0000000001000039000000010100403900000c710030009c000000520000213d0000000100100190000000520000c13d000000400030043f000000180100002900000000002104350000001b01000029000000200a1000390000000004a400190000001c0040006c0000004c0000213d00000000004a004b0000001d0e00035f0000101a0000813d0000001802000029000011380000013d00000020022000390000000001bf0019000000000001043500000060019000390000000000c104350000002001d0003900000000011e034f000000000101043b000000800390003900000000001304350000000000920435000000200aa0003900000000004a004b0000101a0000813d0000000001ae034f000000000101043b00000c710010009c0000004c0000213d0000001b0b1000290000001a01b0006900000c740010009c0000004c0000213d000000a00010008c0000004c0000413d000000400900043d00000cc50090009c000000520000213d000000200db000390000001d01d0035f000000a00c9000390000004000c0043f000000000101043b00000c710010009c0000004c0000213d0000000006b100190000003f016000390000001c0010006c000000000300001900000c750300804100000c7501100197000000000001004b000000000500001900000c750500404100000c750010009c000000000503c019000000000005004b0000004c0000c13d00000020086000390000001d0180035f000000000f01043b00000c7100f0009c000000520000213d0000001f01f0003900000d0e011001970000003f0110003900000d0e011001970000000001c1001900000c710010009c000000520000213d0000004003600039000000400010043f0000000000fc043500000000013f00190000001c0010006c0000004c0000213d00000020018000390000001d0310035f00000d0e05f00198000000c00e90003900000000065e0019000011770000613d000000000803034f00000000010e0019000000008708043c0000000001710436000000000061004b000011730000c13d0000001f01f00190000011840000613d000000000353034f0000000301100210000000000506043300000000051501cf000000000515022f000000000303043b0000010001100089000000000313022f00000000011301cf000000000151019f00000000001604350000000001ef001900000000000104350000000003c904360000002001d000390000001d0510035f000000000505043b00000c720050009c0000004c0000213d000000000053043500000020011000390000001d0310035f000000000303043b00000c6d0030009c0000004c0000213d00000040059000390000000000350435000000200d1000390000001d01d0035f000000000101043b00000c710010009c0000004c0000213d000000000fb100190000003f01f000390000001c0010006c000000000300001900000c750300804100000c7501100197000000000001004b000000000500001900000c750500404100000c750010009c000000000503c019000000000005004b0000004c0000c13d0000002006f000390000001d0160035f000000000b01043b00000c7100b0009c000000520000213d0000001f01b0003900000d0e011001970000003f0110003900000d0e01100197000000400c00043d00000000051c00190000000000c5004b0000000001000039000000010100403900000c710050009c000000520000213d0000000100100190000000520000c13d0000004001f00039000000400050043f000000000fbc043600000000011b00190000001c0010006c0000004c0000213d00000020016000390000001d0e00035f00000000051e034f00000d0e08b0019800000000068f0019000011ca0000613d000000000305034f00000000010f0019000000003703043c0000000001710436000000000061004b000011c60000c13d0000001f01b001900000112a0000613d000000000385034f0000000301100210000000000506043300000000051501cf000000000515022f000000000303043b0000010001100089000000000313022f00000000011301cf000000000151019f00000000001604350000112a0000013d0000000e010000290000000f0200002900000000002104350000001101000029000000600110003900000000011e034f000000000101043b00000c710010009c0000004c0000213d00000011031000290000003f013000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000400001900000c750400404100000c750010009c000000000402c019000000000004004b0000004c0000c13d0000002002300039001800000002001d0000001d0120035f000000000201043b00000c710020009c000000520000213d00000005052002100000003f0150003900000c7301100197000000400600043d0000000004160019001700000006001d000000000064004b0000000001000039000000010100403900000c710040009c000000520000213d0000000100100190000000520000c13d000000400040043f0000001701000029000000000021043500000040043000390000000002450019001900000002001d0000001c0020006c0000004c0000213d000000190040006c0000001d0300035f0000128b0000813d0000001702000029000012150000013d0000001b0200002900000020022000390000000000a204350000002004400039000000190040006c0000001d0300035f0000128b0000813d001b00000002001d000000000143034f000000000101043b00000c710010009c0000004c0000213d00000018091000290000003f019000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000300001900000c750300404100000c750010009c000000000302c019000000000003004b0000004c0000c13d00000020019000390000001d0110035f000000000201043b00000c710020009c000000520000213d00000005032002100000003f0130003900000c7301100197000000400a00043d00000000051a00190000000000a5004b0000000001000039000000010100403900000c710050009c000000520000213d0000000100100190000000520000c13d000000400050043f00000000002a0435000000400b900039000000000cb300190000001c00c0006c0000004c0000213d0000000000cb004b0000120e0000813d000000000d0a0019000012490000013d000000200dd000390000000001e3001900000000000104350000000000fd0435000000200bb000390000000000cb004b0000120e0000813d0000001d01b0035f000000000101043b00000c710010009c0000004c0000213d00000000039100190000005f013000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000500001900000c750500404100000c750010009c000000000502c019000000000005004b0000004c0000c13d00000040063000390000001d0160035f000000000e01043b00000c7100e0009c000000520000213d0000001f01e0003900000d0e011001970000003f0110003900000d0e01100197000000400f00043d00000000021f00190000000000f2004b0000000001000039000000010100403900000c710020009c000000520000213d0000000100100190000000520000c13d0000006001300039000000400020043f0000000003ef043600000000011e00190000001c0010006c0000004c0000213d00000020016000390000001d0810035f00000d0e05e0019800000000065300190000127d0000613d000000000208034f0000000001030019000000002702043c0000000001710436000000000061004b000012790000c13d0000001f01e00190000012420000613d000000000258034f0000000301100210000000000506043300000000051501cf000000000515022f000000000202043b0000010001100089000000000212022f00000000011201cf000000000151019f0000000000160435000012420000013d000000100100002900000040011000390000001702000029000000000021043500000011010000290000008002100039000000000123034f000000000101043b00000c710010009c0000004c0000213d00000011041000290000003f014000390000001c0010006c000000000300001900000c750300804100000c7501100197000000000001004b000000000500001900000c750500404100000c750010009c000000000503c019000000000005004b0000004c0000c13d00000020014000390000001d0110035f000000000501043b00000c710050009c000000520000213d00000005065002100000003f0160003900000c7301100197000000400300043d0000000008130019000000000038004b0000000001000039000000010100403900000c710080009c000000520000213d0000000100100190000000520000c13d0000004004400039000000400080043f000000000053043500000000054600190000001c0050006c0000004c0000213d000000000045004b0000001d0700035f00000fbb0000a13d0000000006030019000000000147034f000000000101043b000000200660003900000000001604350000002004400039000000000054004b000012bd0000413d00000fbb0000013d000000000101043b00000000020000190000001a04000029000000000301041a00000c7203300197000000000434043600000001011000390000000102200039000000000052004b000012c80000413d0000001b0140006a0000001f0110003900000d0e021001970000001b01200029000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f00000018010000290000001b0200002900000000002104350000001c010000290000000301100039000000000301041a000000400200043d001c00000002001d001b00000003001d0000000002320436001a00000002001d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001b0000006b000012f40000c13d0000001a04000029000012ff0000013d000000000101043b00000000020000190000001a040000290000001b05000029000000000301041a00000c7203300197000000000434043600000001011000390000000102200039000000000052004b000012f80000413d0000001c0140006a0000001f0110003900000d0e021001970000001c01200029000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f0000001d0400002900000040024000390000001c0100002900000000001204350000002003000039000000400100043d00000000033104360000000004040433000000006504043400000000005304350000000003060433000000ff0330018f0000004005100039000000000035043500000040034000390000000003030433000000ff0330018f0000006005100039000000000035043500000060034000390000000003030433000000000003004b0000000003000039000000010300c0390000008004100039000000000034043500000018030000290000000004030433000000a003100039000000c0050000390000000000530435000000e003100039000000000504043300000000005304350000010003100039000000000005004b000013380000613d00000000060000190000002004400039000000000704043300000c720770019700000000037304360000000106600039000000000056004b000013310000413d0000000004130049000000200540008a0000000004020433000000c002100039000000000052043500000000050404330000000002530436000000000005004b000013490000613d00000000030000190000002004400039000000000604043300000c720660019700000000026204360000000103300039000000000053004b000013420000413d000000000212004900000a6e0000013d000d00800000003d0000001b0100002900190020001000920000001d0700035f000013600000013d0000000d0400002900000020044000390000001105000029000000600350003900000000002304350000002001100039000000000117034f000000000101043b00000080025000390000000000120435000d00000004001d00000000005404350000000c0400002900000020044000390000000a0040006c0000193f0000813d000c00000004001d000000000147034f000000000101043b00000c710010009c0000004c0000213d0000000b02100029001200000002001d000000190120006900000c740010009c0000004c0000213d000000a00010008c0000004c0000413d000000400100043d001100000001001d00000cc50010009c000000520000213d0000001101000029000000a001100039000000400010043f000000120100002900000020011000390000001d0210035f000000000202043b00000c710020009c0000004c0000213d00000011030000290000000002230436000e00000002001d00000020011000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d00000012011000290000003f021000390000001b04000029000000000042004b000000000300001900000c750300804100000c750220019700000c7509400197000000000492013f000000000092004b000000000200001900000c750200404100000c750040009c000000000203c019000000000002004b0000004c0000c13d00000020021000390000001d0220035f000000000202043b00000c710020009c000000520000213d00000005032002100000003f0430003900000c7304400197000000400500043d0000000004450019001000000005001d000000000054004b0000000005000039000000010500403900000c710040009c000000520000213d0000000100500190000000520000c13d000000400040043f00000010040000290000000000240435001400400010003d0000001402300029001300000002001d0000001b0020006c0000004c0000213d0000001302000029000000140020006b0000001d0e00035f000015760000813d00000014020000290000001003000029000013bd0000013d000000150300002900000020033000390000001802000029000000a00120003900000017040000290000000000410435000000000023043500000016020000290000002002200039000000130020006c000015760000813d001500000003001d001600000002001d00000000012e034f000000000101043b00000c710010009c0000004c0000213d00000014011000290000001b0210006900000c740020009c0000004c0000213d000001400020008c0000004c0000413d000000400200043d001800000002001d00000cb30020009c000000520000213d0000001803000029000000c002300039000000400020043f00000ce30030009c000000520000213d00000018030000290000016003300039000000400030043f0000001d0310035f000000000303043b000000000032043500000020031000390000001d0430035f000000000404043b00000c710040009c0000004c0000213d0000001805000029000000e005500039000000000045043500000020033000390000001d0430035f000000000404043b00000c710040009c0000004c0000213d00000018050000290000010005500039000000000045043500000020033000390000001d0430035f000000000404043b00000c710040009c0000004c0000213d00000018050000290000012005500039000000000045043500000020033000390000001d0430035f000000000404043b00000c710040009c0000004c0000213d000000180500002900000000022504360000014005500039000000000045043500000020033000390000001d0430035f000000000404043b00000c710040009c0000004c0000213d000000000a1400190000001f04a000390000001b0040006c000000000500001900000c750500804100000c7504400197000000000694013f000000000094004b000000000400001900000c750400404100000c750060009c000000000405c019000000000004004b0000004c0000c13d0000001d04a0035f000000000704043b00000c710070009c000000520000213d0000001f0470003900000d0e044001970000003f0440003900000d0e04400197000000400800043d0000000004480019000000000084004b0000000005000039000000010500403900000c710040009c000000520000213d0000000100500190000000520000c13d0000002005a00039000000400040043f000000000a78043600000000045700190000001b0040006c0000004c0000213d0000001d0e00035f00000000055e034f00000d0e04700198000000000b4a00190000142e0000613d000000000c05034f000000000d0a001900000000c60c043c000000000d6d04360000000000bd004b0000142a0000c13d0000001f067001900000143b0000613d000000000445034f000000030560021000000000060b043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004b043500000000047a001900000000000404350000000000820435000000200230003900000000032e034f000000000303043b00000c710030009c0000004c0000213d00000000081300190000001f038000390000001b0030006c000000000400001900000c750400804100000c7503300197000000000593013f000000000093004b000000000300001900000c750300404100000c750050009c000000000304c019000000000003004b0000004c0000c13d0000001d0380035f000000000303043b00000c710030009c000000520000213d0000001f0430003900000d0e044001970000003f0440003900000d0e04400197000000400700043d0000000004470019000000000074004b0000000005000039000000010500403900000c710040009c000000520000213d0000000100500190000000520000c13d0000002005800039000000400040043f000000000837043600000000045300190000001b0040006c0000004c0000213d0000001d0d00035f00000000055d034f00000d0e04300198000000000a480019000014730000613d000000000b05034f000000000c08001900000000b60b043c000000000c6c04360000000000ac004b0000146f0000c13d0000001f06300190000014800000613d000000000445034f000000030560021000000000060a043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004a043500000000033800190000000000030435000000180300002900000040033000390000000000730435000000200220003900000000032d034f000000000303043b00000c720030009c0000004c0000213d00000018060000290000006004600039000000000034043500000020032000390000001d0330035f000000000303043b0000008004600039000000000034043500000040022000390000001d0220035f000000000202043b00000c710020009c0000004c0000213d0000000001120019001c00000001001d0000001f011000390000001b0010006c000000000200001900000c750200804100000c7501100197000000000391013f000000000091004b000000000100001900000c750100404100000c750030009c000000000102c019000000000001004b0000004c0000c13d0000001c020000290000001d0120035f000000000201043b00000c710020009c000000520000213d00000005032002100000003f0130003900000c7301100197000000400400043d0000000001140019001700000004001d000000000041004b0000000004000039000000010400403900000c710010009c000000520000213d0000000100400190000000520000c13d000000400010043f000000170100002900000000002104350000001c0100002900000020021000390000000003230019001a00000003001d0000001b0030006c0000004c0000213d0000001a0020006c0000001d0e00035f000013b20000813d000000170a000029000014d40000013d000000200aa000390000000001bf0019000000000001043500000060018000390000000000c104350000002001d0003900000000011e034f000000000101043b0000008003800039000000000013043500000000008a043500000020022000390000001a0020006c000013b20000813d00000000012e034f000000000101043b00000c710010009c0000004c0000213d0000001c0b1000290000001901b0006900000c740010009c0000004c0000213d000000a00010008c0000004c0000413d000000400800043d00000cc50080009c000000520000213d000000200db000390000001d01d0035f000000a00c8000390000004000c0043f000000000101043b00000c710010009c0000004c0000213d0000000005b100190000003f015000390000001b0010006c000000000400001900000c750400804100000c7501100197000000000691013f000000000091004b000000000100001900000c750100404100000c750060009c000000000104c019000000000001004b0000004c0000c13d00000020075000390000001d0170035f000000000f01043b00000c7100f0009c000000520000213d0000001f01f0003900000d0e011001970000003f0110003900000d0e011001970000000001c1001900000c710010009c000000520000213d0000004004500039000000400010043f0000000000fc043500000000014f00190000001b0010006c0000004c0000213d00000020017000390000001d0110035f00000d0e04f00198000000c0068000390000000005460019000015140000613d000000000701034f000000000e060019000000007307043c000000000e3e043600000000005e004b000015100000c13d0000001f07f00190000015210000613d000000000141034f0000000303700210000000000405043300000000043401cf000000000434022f000000000101043b0000010003300089000000000131022f00000000013101cf000000000141019f000000000015043500000000016f001900000000000104350000000004c804360000002001d000390000001d0310035f000000000503043b00000c720050009c0000004c0000213d000000000054043500000020011000390000001d0310035f000000000403043b00000c6d0040009c0000004c0000213d00000040038000390000000000430435000000200d1000390000001d01d0035f000000000101043b00000c710010009c0000004c0000213d000000000fb100190000003f01f000390000001b0010006c000000000300001900000c750300804100000c7501100197000000000491013f000000000091004b000000000100001900000c750100404100000c750040009c000000000103c019000000000001004b0000004c0000c13d0000002005f000390000001d0150035f000000000b01043b00000c7100b0009c000000520000213d0000001f01b0003900000d0e011001970000003f0110003900000d0e01100197000000400c00043d00000000041c00190000000000c4004b0000000001000039000000010100403900000c710040009c000000520000213d0000000100100190000000520000c13d0000004001f00039000000400040043f000000000fbc043600000000011b00190000001b0010006c0000004c0000213d00000020015000390000001d0e00035f00000000041e034f00000d0e07b0019800000000057f0019000015680000613d000000000104034f00000000060f0019000000001301043c0000000006360436000000000056004b000015640000c13d0000001f01b00190000014c60000613d000000000374034f0000000301100210000000000405043300000000041401cf000000000414022f000000000303043b0000010001100089000000000313022f00000000011301cf000000000141019f0000000000150435000014c60000013d0000000e01000029000000100200002900000000002104350000001201000029000000600110003900000000011e034f000000000101043b00000c710010009c0000004c0000213d00000012021000290000003f012000390000001b0010006c000000000300001900000c750300804100000c7501100197000000000491013f000000000091004b000000000100001900000c750100404100000c750040009c000000000103c019000000000001004b0000004c0000c13d0000002003200039001700000003001d0000001d0130035f000000000101043b00000c710010009c000000520000213d00000005041002100000003f0340003900000c7303300197000000400500043d0000000003350019001600000005001d000000000053004b0000000005000039000000010500403900000c710030009c000000520000213d0000000100500190000000520000c13d000000400030043f00000016030000290000000000130435000000400120003900000000020100190000000003140019001800000003001d0000001b0030006c0000004c0000213d0000000001020019000000180010006c0000001d0200035f000016300000813d0000001603000029000015b70000013d0000001a0300002900000020033000390000000000a304350000001c010000290000002001100039000000180010006c0000001d0200035f000016300000813d001a00000003001d001c00000001001d000000000112034f000000000101043b00000c710010009c0000004c0000213d00000017081000290000003f018000390000001b0010006c000000000200001900000c750200804100000c7501100197000000000491013f000000000091004b000000000100001900000c750100404100000c750040009c000000000102c019000000000001004b0000004c0000c13d00000020018000390000001d0110035f000000000101043b00000c710010009c000000520000213d00000005021002100000003f0420003900000c7304400197000000400a00043d00000000044a00190000000000a4004b0000000005000039000000010500403900000c710040009c000000520000213d0000000100500190000000520000c13d000000400040043f00000000001a0435000000400b800039000000000cb200190000001b00c0006c0000004c0000213d0000000000cb004b000015af0000813d000000000d0a0019000015ed0000013d000000200dd000390000000001e2001900000000000104350000000000fd0435000000200bb000390000000000cb004b000015af0000813d0000001d01b0035f000000000101043b00000c710010009c0000004c0000213d00000000028100190000005f012000390000001b0010006c000000000400001900000c750400804100000c7501100197000000000591013f000000000091004b000000000100001900000c750100404100000c750050009c000000000104c019000000000001004b0000004c0000c13d00000040052000390000001d0150035f000000000e01043b00000c7100e0009c000000520000213d0000001f01e0003900000d0e011001970000003f0110003900000d0e01100197000000400f00043d00000000011f00190000000000f1004b0000000004000039000000010400403900000c710010009c000000520000213d0000000100400190000000520000c13d0000006004200039000000400010043f0000000002ef043600000000014e00190000001b0010006c0000004c0000213d00000020015000390000001d0710035f00000d0e04e001980000000005420019000016220000613d000000000107034f0000000006020019000000001301043c0000000006360436000000000056004b0000161e0000c13d0000001f01e00190000015e60000613d000000000347034f0000000301100210000000000405043300000000041401cf000000000414022f000000000303043b0000010001100089000000000313022f00000000011301cf000000000141019f0000000000150435000015e60000013d000000110100002900000040011000390000001603000029000000000031043500000012010000290000008001100039000000000212034f000000000202043b00000c710020009c0000004c0000213d00000012032000290000003f023000390000001b0020006c000000000400001900000c750400804100000c7502200197000000000592013f000000000092004b000000000200001900000c750200404100000c750050009c000000000204c019000000000002004b0000004c0000c13d00000020023000390000001d0220035f000000000402043b00000c710040009c000000520000213d00000005054002100000003f0250003900000c7306200197000000400200043d0000000007620019000000000027004b0000000006000039000000010600403900000c710070009c000000520000213d0000000100600190000000520000c13d0000004003300039000000400070043f000000000042043500000000043500190000001b0040006c0000004c0000213d000000000034004b0000001d0700035f000013500000a13d0000000005020019000000000637034f000000000606043b000000200550003900000000006504350000002003300039000000000043004b000016630000413d000013500000013d000000400100043d00000cb702000041000000000021043500000004021000390000000403000039000000000032043500000c6d0010009c00000c6d01008041000000400110021000000c85011001c7000031b30001043000000cb801000041000000000013043500000c6d0030009c00000c6d03008041000000400130021000000c90011001c7000031b3000104300000001805000029000016840000013d000000000575043600000020033000390000001a0030006c0000001d0100035f000000ba0000813d000000000131034f000000000101043b00000c710010009c0000004c0000213d0000001b0610002900000043016000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000400001900000c750400404100000c750010009c000000000402c019000000000004004b0000004c0000c13d00000024016000390000001d0110035f000000000801043b00000c710080009c000000520000213d00000005098002100000003f0190003900000c7301100197000000400700043d0000000002170019000000000072004b0000000001000039000000010100403900000c710020009c000000520000213d0000000100100190000000520000c13d000000400020043f0000000000870435000000440860003900000000098900190000001c0090006c0000004c0000213d000000000098004b0000167f0000813d000000000a070019000016b60000013d000000200aa000390000000000dc04350000000000ba04350000002008800039000000000098004b0000167f0000813d0000001d0180035f000000000101043b00000c710010009c0000004c0000213d000000000d6100190000001201d0006900000c740010009c0000004c0000213d000000400010008c0000004c0000413d000000400b00043d00000c8800b0009c000000520000213d0000004001b00039000000400010043f0000004401d000390000001d0110035f000000000101043b000000000c1b04360000006401d000390000001d0110035f000000000101043b00000c710010009c0000004c0000213d000000000ed100190000006301e000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000400001900000c750400404100000c750010009c000000000402c019000000000004004b0000004c0000c13d0000004401e000390000001d0110035f000000000f01043b00000c7100f0009c000000520000213d0000000502f002100000003f0120003900000c7301100197000000400d00043d00000000041d00190000000000d4004b0000000001000039000000010100403900000c710040009c000000520000213d0000000100100190000000520000c13d000000400040043f0000000000fd0435000000640ee00039000000000fe200190000001c00f0006c0000004c0000213d0000000000fe004b000016b00000813d00000000020d00190000001d01e0035f000000000101043b00000c6d0010009c0000004c0000213d00000020022000390000000000120435000000200ee000390000000000fe004b000016f50000413d000016b00000013d0000001c01000029000001000110003900000014070000290000170c0000013d00000020077000390000000009a900190000002009900039000000000009043500000060098000390000000000a904350000000000870435000000000054004b0000017c0000813d000000004804043400000c710080009c0000004c0000213d0000000009380019000000000891004900000c740080009c0000004c0000213d000000800080008c0000004c0000413d000000400800043d00000c700080009c000000520000213d000000800a8000390000004000a0043f000000200a900039000000000a0a043300000c7200a0009c0000004c0000213d000000000aa80436000000400b900039000000000b0b043300000c7100b0009c0000004c0000213d0000000000ba0435000000600a900039000000000a0a043300000000000a004b000000000b000039000000010b00c0390000000000ba004b0000004c0000c13d000000400b8000390000000000ab0435000000800a900039000000000a0a043300000c7100a0009c0000004c0000213d000000000b9a00190000003f09b00039000000000029004b000000000a00001900000c750a00804100000c7509900197000000000009004b000000000c00001900000c750c00404100000c750090009c000000000c0ac01900000000000c004b0000004c0000c13d0000002009b00039000000000909043300000c710090009c000000520000213d0000001f0a90003900000d0e0aa001970000003f0aa0003900000d0e0ca00197000000400a00043d000000000cca00190000000000ac004b000000000d000039000000010d00403900000c7100c0009c000000520000213d0000000100d00190000000520000c13d0000004000c0043f000000000c9a0436000000400bb00039000000000db9001900000000002d004b0000004c0000213d000000000009004b000017030000613d000000000d000019000000000edc0019000000000fbd0019000000000f0f04330000000000fe0435000000200dd0003900000000009d004b000017580000413d000017030000013d0000000102000039000000000302041a00000c7603300197000000000113019f000000000012041b00000c77010000410000000000100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c78011001c70000800b0200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000800010043f0000001b01000029000000000101043300000c72001001980000177e0000613d0000001a01000029000000000101043300000c72001001980000177e0000613d0000001801000029000000000101043300000c72001001980000186d0000c13d000000400100043d00000cc00200004100000e0c0000013d0000000f02000029000000a402200039000000000003004b000017a10000613d0000001d0300002900000000040304330000000503400210000000000004004b000017920000c13d000000a004000039000017980000013d000000400100043d00000cb702000041000000000021043500000004021000390000000103000039000016700000013d00000d120030009c00000b280000213d00000000044300d9000000200040008c00000b280000c13d000000a00430003900000000053400190000000002250019000000000052004b00000000050000390000000105004039000000000034001a00000b280000413d000000010050019000000b280000c13d0000000003000031000000000023004b000017e50000c13d0000001a0010006b000017f10000c13d00000cc9010000410000000000100443000000000100041200000004001004430000002400000443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000025790000613d000000000101043b001700000001001d00000c77010000410000000000100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c78011001c70000800b0200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000170010006b000018f90000c13d000000000100041100000c7201100197000000000010043f00000cef01000041000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000400200043d00000c880020009c000000520000213d000000000101043b0000004003200039000000400030043f000000000301041a000000ff0130018f00000000011204360000000803300270000000ff0330018f000000020030008c000019430000a13d00000d0a01000041000000000010043f0000002101000039000000040010043f00000c8501000041000031b300010430000000400100043d0000002404100039000000000034043500000ced0300004100000000003104350000000403100039000000000023043500000c6d0010009c00000c6d01008041000000400110021000000cd4011001c7000031b300010430000000400200043d0000002403200039000000000013043500000cee01000041000000000012043500000004012000390000001a03000029000025810000013d000000400100043d00000cb702000041000000000021043500000004021000390000000000020435000016710000013d00000cb701000041000000000013043500000004013000390000000502000039000000000021043500000c6d0030009c00000c6d03008041000000400130021000000c85011001c7000031b30001043000000000020000190000180e0000013d0000000102200039000000000012004b000000e00000813d00000005032002100000000004f300190000000004040433000000a00330003900000000030304330000000047040434000000200530003900000000050504330000000065050434000000000075004b000018580000c13d000000000005004b0000180b0000613d0000000007000019000018200000013d0000000107700039000000000057004b0000180b0000813d0000000508700210000000000984001900000000088600190000000008080433000000000909043300000000b9090434000000000009004b0000182c0000613d000000800a800039000000000a0a04330000000000a9004b0000192d0000413d000000a009800039000000000909043300000000a9090434000000000b0b043300000000bc0b04340000000000c9004b0000191c0000c13d000000000009004b0000181d0000613d000000000c0000190000183a0000013d000000010cc0003900000000009c004b0000181d0000813d000000050ec00210000000000deb0019000000000d0d043300000c6d0dd00198000018370000613d000000000eea0019000000000e0e0433000000400ee00039000000000e0e043300000c6d0ee001970000000000ed004b000018370000813d00000000010804330000000001010433000000400200043d00000064032000390000000000d3043500000044032000390000000000e3043500000024032000390000000000c3043500000ce80300004100000000003204350000000403200039000000000013043500000c6d0020009c00000c6d02008041000000400120021000000ce9011001c7000031b300010430000000400100043d00000ce50200004100000e0c0000013d000000400100043d00000cb402000041000000000021043500000004021000390000001d03000029000016700000013d000000400100043d00000cb702000041000000000021043500000004021000390000000203000039000016700000013d000000400100043d00000cb702000041000000000021043500000004021000390000000303000039000016700000013d0000001d01000029000000000101043300000c7101100198000018770000c13d000000400100043d00000d0d0200004100000e0c0000013d000000400100043d00000cbc0200004100000e0c0000013d000000a00010043f0000001b03000029000000000103043300000c7201100197000000c00010043f0000001a04000029000000000104043300000c7201100197000000e00010043f0000001805000029000000000105043300000c7201100197000001000010043f0000001d01000029000000000101043300000c7101100197000000400200043d0000000001120436000000000303043300000c72033001970000000000310435000000000104043300000c720110019700000040032000390000000000130435000000000105043300000c72011001970000006003200039000000000013043500000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c79011001c70000800d02000039000000010300003900000c7a0400004131b131a70000040f00000001002001900000004c0000613d0000001901000029000000000101043300000c72011001980000177e0000613d0000000402000039000000000302041a00000c7b0330019700000017050000290000000004050433000000a00440021000000c7c04400197000000000343019f00000016060000290000000004060433000000000004004b00000c7d040000410000000004006019000000000343019f000000000313019f000000000032041b0000001507000029000000000207043300000c72022001970000000504000039000000000304041a00000c7603300197000000000223019f000000000024041b000000400200043d0000000001120436000000000305043300000c6d0330019700000000003104350000000001060433000000000001004b0000000001000039000000010100c03900000040032000390000000000130435000000000107043300000c72011001970000006003200039000000000013043500000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c79011001c70000800d02000039000000010300003900000c7e0400004131b131a70000040f00000001002001900000004c0000613d00000014010000290000000001010433000000000001004b0000197b0000c13d000000800100043d00000140000004430000016000100443000000a00100043d00000020020000390000018000200443000001a000100443000000c00100043d0000004003000039000001c000300443000001e0001004430000006001000039000000e00300043d00000200001004430000022000300443000001000100043d00000080030000390000024000300443000002600010044300000100002004430000000501000039000001200010044300000c8e01000041000031b20001042e000000400200043d0000002403200039000000000013043500000ce401000041000000000012043500000004012000390000001703000029000025810000013d00000140070000390000001a0850006900000c740080009c0000004c0000213d000000400080008c0000004c0000413d000000400800043d00000c880080009c000000520000213d0000004009800039000000400090043f0000001d0950035f000000000909043b00000c720090009c0000004c0000213d0000000009980436000000200a5000390000001d0aa0035f000000000a0a043b00000cf300a0009c0000004c0000213d0000000000a9043500000000078704360000004005500039000000000065004b000019020000413d000005430000013d000000000108043300000000020104330000006001100039000000000101043300000c7101100197000000400300043d0000002404300039000000000014043500000ce70100004100000000001304350000000401300039000000000021043500000c6d0030009c00000c6d03008041000000400130021000000cd4011001c7000031b300010430000000000103043300000000020804330000000002020433000000400300043d000000440430003900000000009404350000002404300039000000000024043500000ce602000041000000000023043500000c71011001970000000402300039000000000012043500000c6d0030009c00000c6d03008041000000400130021000000cd9011001c7000031b300010430000000400200043d00000cc70020009c000000520000213d00000af80000013d0000000000310435000025870000c13d0000000101000039000000000010043f0000000201000039000000200010043f0000000001020433001700ff0010019300000cf001000041000000000201041a000000170020006c000025730000a13d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001701100029000000000101041a00000c72011001970000000002000411000000000012004b000025870000c13d00000019010000290000000001010433000000000001004b00001c170000c13d00000024010000390000000201100367000000000101043b00000c7101100197000000400200043d000000200320003900000000001304350000001a01000029000000000012043500000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c7f011001c70000800d020000390000000203000039000000010500003900000cf20400004100000b460000013d0000000002000019001600000002001d00000005012002100000001301100029000000000201043300000020012000390000000001010433001c0c710010019c000018710000613d001800000002001d000000000102043300000c72001001980000177e0000613d0000001c01000029000000000010043f0000000801000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b001a00000001001d0000000101100039001d00000001001d000000000101041a000000010210019000000001051002700000007f0550618f0000001f0050008c00000000030000390000000103002039000000000032004b00000b860000c13d0000001803000029000000600330003900000000040304330000001a03000029000000000303041a001b00000005001d000000000005004b001900000004001d001700200040003d000019cd0000613d00000c800330019700000c810030009c00001a200000613d000000400300043d001500000003001d0000001b040000290000000003430436001200000003001d000000000002004b000019e40000613d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000201043b000000000100001900000012050000290000001b060000290000000003510019000000000402041a000000000043043500000001022000390000002001100039000000000061004b000019c50000413d000019e80000013d00000c860130019700000c81011001c70000001a02000029000000000012041b000000400100043d0000001c02000029000000000021043500000c6d0010009c00000c6d010080410000004001100210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c82011001c70000800d02000039000000010300003900000c870400004131b131a70000040f000000010020019000001a200000c13d0000004c0000013d00000d10011001970000001202000029000000000012043500000020010000390000003f0110003900000d0e021001970000001501200029000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f000000120100002900000c6d0010009c00000c6d0100804100000040011002100000001502000029000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000170200002900000c6d0020009c00000c6d0200804100000040022002100000001903000029000000000303043300000c6d0030009c00000c6d030080410000006003300210000000000223019f000000000101043b001b00000001001d000000000100041400000c6d0010009c00000c6d01008041000000c001100210000000000121019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001b0010006b00001b600000c13d00000019010000290000000001010433000000000001004b0000177e0000613d00000c6d0010009c00000c6d010080410000006001100210000000170200002900000c6d0020009c00000c6d020080410000004002200210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b001b00000001001d000000400100043d00000020020000390000000002210436000000000002043500000c880010009c000000520000213d0000004003100039000000400030043f00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001b0010006b0000177e0000613d00000019010000290000000001010433001b00000001001d00000c710010009c000000520000213d0000001d01000029000000000101041a000000010010019000000001031002700000007f0330618f0000001f0030008c00000000020000390000000102002039000000000121013f000000010010019000000b860000c13d001500000003001d000000200030008c00001a840000413d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001b030000290000001f023000390000000502200270000000200030008c0000000002004019000000000301043b00000015010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00001a840000813d000000000002041b0000000102200039000000000012004b00001a800000413d0000001b010000290000001f0010008c00001aa40000a13d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000200200008a0000001b02200180000000000101043b00001ab20000613d000000010320008a0000000503300270000000000331001900000001043000390000002003000039000000190600002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b00001a9c0000c13d00001ab30000013d0000001b0000006b00001aa90000613d0000001701000029000000000101043300001aaa0000013d00000000010000190000001b04000029000000030240021000000d110220027f00000d1102200167000000000121016f0000000102400210000000000121019f00001ac00000013d00000020030000390000001b05000029000000000052004b00001abe0000813d0000000302500210000000f80220018f00000d110220027f00000d110220016700000019033000290000000003030433000000000223016f000000000021041b000000010150021000000001011001bf0000001d02000029000000000012041b000000180400002900000040014000390000000001010433000000000001004b00000c890100004100000000010060190000001a03000029000000000203041a00000c8a02200197000000000112019f000000000204043300000c7202200197000000000121019f000000000013041b0000001c01000029000000000010043f0000000701000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000101041a000000000001004b00001afc0000c13d0000000601000039000000000101041a00000c710010009c000000520000213d00000001021000390000000603000039000000000023041b00000c8b0110009a0000001c02000029000000000021041b000000000103041a001b00000001001d000000000020043f0000000701000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001b02000029000000000021041b000000400500043d000000200100003900000000011504360000001a02000029000000000202041a00000080035000390000008004000039000000000043043500000c72032001970000000000310435000000a80120027000000c71011001970000006003500039000000000013043500000c8c002001980000000001000039000000010100c039000000400250003900000000001204350000001d01000029000000000101041a000000010210019000000001041002700000007f0440618f0000001f0040008c00000000030000390000000103002039000000000331013f000000010030019000000b860000c13d001a00000005001d000000a003500039001b00000004001d0000000000430435000000000002004b00001b3a0000613d0000001d01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001b0000006b00001b420000613d0000001a06000029000000c002600039000000000301043b00000000010000190000001b070000290000000004210019000000000503041a000000000054043500000001033000390000002001100039000000000071004b00001b320000413d00001b440000013d00000d10011001970000001a06000029000000c00260003900000000001204350000001b0000006b0000002001000039000000000100603900001b440000013d00000000010000190000001a0600002900000c6d0060009c00000c6d060080410000004002600210000000c00110003900000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000121019f00000c83011001c70000800d02000039000000020300003900000c8d040000410000001c0500002931b131a70000040f00000001002001900000004c0000613d0000001602000029000000010220003900000014010000290000000001010433000000000012004b0000197c0000413d000018e10000013d000000400100043d00000c8402000041000000000021043500000004021000390000001c03000029000016700000013d0000001c01000029000000200610008a000000190700002900001b780000013d000000200770003900000000019c0019000000000001043500000060018000390000000000b104350000002001a000390000001d0110035f000000000101043b00000080028000390000000000120435000000000087043500000020044000390000001a0040006c00000c5c0000813d0000001d0140035f000000000101043b00000c710010009c0000004c0000213d0000001b09100029000000000196004900000c740010009c0000004c0000213d000000a00010008c0000004c0000413d000000400800043d00000cc50080009c000000520000213d000000200a9000390000001d01a0035f000000a00b8000390000004000b0043f000000000101043b00000c710010009c0000004c0000213d000000000d9100190000003f01d000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000300001900000c750300404100000c750010009c000000000302c019000000000003004b0000004c0000c13d000000200ed000390000001d01e0035f000000000c01043b00000c7100c0009c000000520000213d0000001f01c0003900000d0e011001970000003f0110003900000d0e011001970000000001b1001900000c710010009c000000520000213d0000004002d00039000000400010043f0000000000cb043500000000012c00190000001c0010006c0000004c0000213d0000002001e000390000001d0210035f00000d0e03c00198000000c00d800039000000000e3d001900001bb70000613d000000000f02034f00000000010d001900000000f50f043c00000000015104360000000000e1004b00001bb30000c13d0000001f01c0019000001bc40000613d000000000232034f000000030110021000000000030e043300000000031301cf000000000313022f000000000202043b0000010001100089000000000212022f00000000011201cf000000000131019f00000000001e04350000000001dc001900000000000104350000000002b804360000002001a000390000001d0310035f000000000303043b00000c720030009c0000004c0000213d000000000032043500000020011000390000001d0210035f000000000202043b00000c6d0020009c0000004c0000213d00000040038000390000000000230435000000200a1000390000001d01a0035f000000000101043b00000c710010009c0000004c0000213d000000000c9100190000003f01c000390000001c0010006c000000000200001900000c750200804100000c7501100197000000000001004b000000000300001900000c750300404100000c750010009c000000000302c019000000000003004b0000004c0000c13d000000200dc000390000001d01d0035f000000000901043b00000c710090009c000000520000213d0000001f0190003900000d0e011001970000003f0110003900000d0e01100197000000400b00043d00000000031b00190000000000b3004b0000000001000039000000010100403900000c710030009c000000520000213d0000000100100190000000520000c13d0000004001c00039000000400030043f000000000c9b043600000000011900190000001c0010006c0000004c0000213d0000002001d000390000001d0310035f00000d0e0e900198000000000dec001900001c090000613d000000000203034f00000000010c0019000000002502043c00000000015104360000000000d1004b00001c050000c13d0000001f0190019000001b6a0000613d0000000002e3034f000000030110021000000000030d043300000000031301cf000000000313022f000000000202043b0000010001100089000000000212022f00000000011201cf000000000131019f00000000001d043500001b6a0000013d00000018010000290000000001010433000000ff0110018f000000ff0010008c00000b280000613d00000001011000390000001d020000290000000002020433000000000012004b0000258a0000c13d0000000f010000290000001f0110003900000d0e011001970000003f0110003900000d0e02100197000000400100043d0000000002210019000000000012004b0000000004000039000000010400403900000c710020009c000000520000213d0000000100400190000000520000c13d000000400020043f0000000f0200002900000000022104360000001b05000029000000000050007c0000004c0000213d0000000f0500002900000d0e045001980000001f0550018f00000009030000290000000206300367000000000342001900001c420000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b00001c3e0000c13d000000000005004b00001c4f0000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004304350000000f03200029000000000003043500000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000301043b000000400100043d00000020021000390000000000320435000000a003100039000000400410003900000004050000390000000205500367000000005605043c0000000004640436000000000034004b00001c6b0000c13d0000008004000039000000000041043500000cc50010009c000000520000213d000000400030043f00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b001800000001001d0000001d010000290000000001010433001700000001001d000000000001004b000019640000613d001b00000000001d001900000000001d0000001b010000290000001f0010008c000025730000213d0000001d0100002900000000010104330000001b0010006c000025730000a13d0000001b0100002900000005011002100000001c011000290000000001010433000000400200043d000000600320003900000000001304350000004003200039000000000013043500000020012000390000001b03000039000000000031043500000018010000290000000000120435000000000000043f00000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000cf1011001c7000000010200003931b131ac0000040f000000600310027000000c6d03300197000000200030008c00000020050000390000000005034019000000200450019000001cbd0000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b00001cb90000c13d0000001f0550019000001cca0000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f0000000000540435000100000003001f000300000001035500000001002001900000258d0000613d000000000100043d00000c7201100197000000000010043f00000cef01000041000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000400200043d00000c880020009c000000520000213d000000000101043b0000004003200039000000400030043f000000000301041a000000ff0130018f00000000021204360000000803300270000000ff0330018f000000020030008c000017df0000213d0000000000320435000000010030008c000025e20000c13d000000010110020f0000001900100180000025e50000c13d00190019001001b30000001b020000290000000102200039001b00000002001d000000170020006c00001c8f0000413d000019640000013d00000cc6010000410000000f02000029000000000012043500000c6d0020009c00000c6d02008041000000400120021000000c90011001c7000031b3000104300000000001000415000100000001001d0000000901000029000000000301043300000c710030009c000000520000213d00000005013002100000003f0210003900000c7302200197000000400500043d0000000004250019000f00000005001d000000000054004b0000000002000039000000010200403900000c710040009c000000520000213d0000000100200190000000520000c13d000000a00200043d000500000002001d000000e00200043d000600000002001d000000800200043d00000020022000390000000002020433000000400040043f0000000f040000290000000004340436001d00000004001d000000000003004b00001d2c0000613d0000000003000019000000400400043d00000c880040009c000000520000213d0000004005400039000000400050043f0000002005400039000000000005043500000000000404350000001d0530002900000000004504350000002003300039000000000013004b00001d1f0000413d00000009010000290000000001010433000000000001004b00001de70000c13d00000000010004150000000101100069000000000100000200000cb10000013d00000cd10100004100000000001004430000001c010000290000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000000001004b0000004c0000613d000000400500043d00000cde0100004100000000001504350000000401500039000000200200003900000000002104350000001d010000290000000002010433000000240150003900000000002104350000001802000029000000000202043300000c710220019700000044035000390000000000230435000000190200002900000000020204330000006403500039000000a0040000390000000000430435000000c40350003900000000420204340000000000230435001700000005001d000000e403500039000000000002004b00001d670000613d000000000500001900000000063500190000000007540019000000000707043300000000007604350000002005500039000000000025004b00001d600000413d000000000432001900000000000404350000001f0220003900000d0e02200197000000000232001900000000031200490000001a04000029000000000404043300000017050000290000008405500039000000000035043500000000430404340000000002320436000000000003004b00001d7e0000613d000000000500001900000000062500190000000007540019000000000707043300000000007604350000002005500039000000000035004b00001d770000413d000000000423001900000000000404350000001f0330003900000d0e03300197000000000423001900000000011400490000001b0200002900000000020204330000001703000029000000a403300039000000000013043500000000030204330000000001340436000000000003004b00001d990000613d000000000400001900000020022000390000000005020433000000006505043400000c720550019700000000055104360000000006060433000000000065043500000040011000390000000104400039000000000034004b00001d8e0000413d00000000020004140000001c03000029000000040030008c00001db10000613d0000001703000029000000000131004900000c6d0010009c00000c6d01008041000000600110021000000c6d0030009c00000c6d030080410000004003300210000000000131019f00000c6d0020009c00000c6d02008041000000c002200210000000000121019f0000001c0200002931b131a70000040f000000600310027000010c6d0030019d00030000000103550000000100200190000025990000613d000000170100002900000c710010009c000000520000213d0000001701000029000000400010043f00000cd00000013d0000000001000415000000210110008a001c000500100218002100000000003d000000e00100043d00000cd102000041000000000020044300000c72011001970000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000000001004b0000001c010000290000000501100270000000000100003f000000010100603f000028e80000613d000000e00300043d0000000001000415001c00000001001d000000400100043d000000200210003900000ccb040000410000000000420435000000240510003900000000004504350000002404000039000000000041043500000cb00010009c000000520000213d00000c72043001970000006003100039000000400030043f001700000004001d000000040040008c000025a00000c13d0000000001020433000000000010043f0000000103000031000025cb0000013d000000060300002900110c720030019b00040c710020019b001b00000000001d0000000b0000006b0000001b02000029001700050020021800001e010000613d0000001b030000290000000b0030006c000025730000813d000000170300002900000003023000290000000202200367000000000202043b00000c6d0020009c0000004c0000213d000000000002004b00001e010000613d00000017030000290000000801300029000000000101043300000040011000390000000000210435000000090100002900000000010104330000001b0010006c000025730000a13d0000001b02000029000000070020006c000025730000813d000000170200002900000008012000290000000001010433001900000001001d00000010012000290000000202000367000000000112034f000000000101043b00000000030000310000000c0430006a000000430440008a00000c750540019700000c7506100197000000000756013f000000000056004b000000000500001900000c7505004041000000000041004b000000000400001900000c750400804100000c750070009c000000000504c019000000000005004b0000004c0000c13d0000001004100029000000000142034f000000000101043b00000c710010009c0000004c0000213d0000000006130049000000200540003900000c750460019700000c7507500197000000000847013f000000000047004b000000000400001900000c7504004041000000000065004b000000000600001900000c750600204100000c750080009c000000000406c019000000000004004b0000004c0000c13d0000001f0410003900000d0e044001970000003f0440003900000d0e04400197000000400600043d0000000004460019001200000006001d000000000064004b0000000006000039000000010600403900000c710040009c000000520000213d0000000100600190000000520000c13d000000400040043f000000120400002900000000041404360000000006510019000000000036004b0000004c0000213d000000000352034f00000d0e05100198000000000254001900001e500000613d000000000603034f0000000007040019000000006806043c0000000007870436000000000027004b00001e4c0000c13d0000001f0610019000001e5d0000613d000000000353034f0000000305600210000000000602043300000000065601cf000000000656022f000000000303043b0000010005500089000000000353022f00000000035301cf000000000363019f0000000000320435000000000114001900000000000104350000000001000415000e00000001001d000000400100043d00000c880010009c000000520000213d0000004002100039000000400020043f000000200210003900000000000204350000000000010435000000190100002900000020011000390000000001010433000000400300043d00000cc802000041000000000023043500000c7202100197001c00000003001d0000000401300039001a00000002001d000000000021043500000cc90100004100000000001004430000000001000412000000040010044300000060010000390000002400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000025790000613d000000000201043b000000000100041400000c7202200197000000040020008c00001e8d0000c13d0000000103000031000000200030008c0000002004000039000000000403401900001eb60000013d0000001c0300002900000c6d0030009c00000c6d03008041000000400330021000000c6d0010009c00000c6d01008041000000c001100210000000000131019f00000c85011001c731b131ac0000040f000000600310027000000c6d03300197000000200030008c0000002004000039000000000403401900000020064001900000001c0560002900001ea50000613d000000000701034f0000001c08000029000000007907043c0000000008980436000000000058004b00001ea10000c13d0000001f0740019000001eb20000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000026370000613d0000001f01400039000000600210018f0000001c01200029000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f000000200030008c0000004c0000413d0000001c010000290000000001010433001c00000001001d00000c720010009c0000004c0000213d0000001c0000006b000026490000613d0000000001000415001600000001001d000000400100043d000000200210003900000ccb040000410000000000420435000000240310003900000000004304350000002403000039000000000031043500000cb00010009c000000520000213d0000006003100039000000400030043f0000001c03000029000000040030008c00001edf0000c13d0000000001020433000000000010043f000000010300003100001f090000013d00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f00000ccc011001c70000001c0200002931b131ac0000040f000000600310027000000c6d03300197000000200030008c00000020050000390000000005034019000000200450019000001ef70000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b00001ef30000c13d0000001f0550019000001f040000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f0000000000540435000100000003001f00030000000103550000000100200190000026460000613d000000000100043d000000200030008c000026460000413d000000000001004b000026460000613d000000400100043d000000200210003900000ccb040000410000000000420435000000240410003900000ccd0500004100000000005404350000002404000039000000000041043500000cb00010009c000000520000213d0000006004100039000000400040043f0000001c04000029000000040040008c00001f200000c13d0000000001020433000000000010043f00001f510000013d00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f00000ccc011001c70000001c0200002931b131ac0000040f000000600310027000000c6d03300197000000200030008c00000020050000390000000005034019000000200450019000001f380000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b00001f340000c13d0000001f0550019000001f450000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f00000000005404350003000000010355000100000003001f0000001f0030008c00000000010000390000000101002039000000000112016f0000000002000415000000230220008a0018000500200218000000010010008c00001f560000c13d000000000100043d0000000002000415000000220220008a0018000500200218000000000001004b000026460000c13d000000400100043d000000200210003900000ccb040000410000000000420435000000240410003900000cce0500004100000000005404350000002404000039000000000041043500000cb00010009c000000520000213d0000006004100039000000400040043f0000001c04000029000000040040008c00001f690000c13d0000000001020433000000000010043f00001f930000013d00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f00000ccc011001c70000001c0200002931b131ac0000040f000000600310027000000c6d03300197000000200030008c00000020050000390000000005034019000000200450019000001f810000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b00001f7d0000c13d0000001f0550019000001f8e0000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f0000000000540435000100000003001f00030000000103550000000100200190000026430000613d000000000100043d000000200030008c000026430000413d000000000001004b00000018010000290000000501100270000000000100003f000000010100c03f000000000100041500000016011000690000000001000002000026490000613d000000190100002900000040011000390000000001010433000000400400043d000000200340003900000ccf02000041000a00000003001d000000000023043500000024024000390000001103000029000000000032043500000024020000390000000000240435001400000004001d00000cb00040009c000000520000213d00000014030000290000006002300039001800000002001d000000400020043f00000cd00030009c000000520000213d00130c6d0010019b00000014040000290000012001400039000000400010043f000000840200003900000018030000290000000000230435000000800340003900000000020000310000000202200367001600000003001d000000002402043c0000000003430436000000000013004b00001fbf0000c13d00000cd10100004100000000001004430000001a010000290000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000000001004b0000265b0000613d0000000001000414000013880110008c0000265f0000413d00000006021002700000000001210049000000130010006c000026630000a13d0000000001000414000d00000001001d0000001a01000029000000040010008c00001fe30000c13d00000003010003670000000104000031000000000200001900001ff70000013d0000000a0100002900000c6d0010009c00000c6d0100804100000040011002100000001402000029000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000112019f0000001302000029000000c002200210000000000121019f0000001a0200002931b131a70000040f000000010220015f0003000000010355000000600310027000010c6d0030019d00000c6d043001970000000003000414000000840040008c000000840400803900000018050000290000000000450435000000e0064001900000001605600029000020050000613d000000000701034f0000001608000029000000007907043c0000000008980436000000000058004b000020010000c13d0000001f04400190000020120000613d000000000161034f0000000304400210000000000605043300000000064601cf000000000646022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000161019f000000000015043500000001002001900000264c0000c13d00000018010000290000000001010433000000200010008c000026700000c13d0000000d01300069001400130010007300000b280000413d000000400600043d00000cd50060009c000000520000213d00000016010000290000000001010433001300000001001d0000001902000029000000800120003900000000070104330000000004020433000000600120003900000000030104330000010001600039000000400010043f000000e00160003900000012020000290000000000210435000000c0026000390000000000320435000000a003600039000000000043043500000080046000390000001a0500002900000000005404350000006005600039000000000075043500000040076000390000001108000029000000000087043500000020096000390000000408000029000000000089043500000005080000290000000000860435000000400c00043d000000200ac0003900000cd608000041000d0000000a001d00000000008a04350000002408c00039000000200a0000390000000000a8043500000000060604330000004408c00039000001000a0000390000000000a80435000001440bc0003900000000a806043400000000008b043500190000000c001d0000016406c00039000000000008004b000020580000613d000000000b000019000000000c6b0019000000000dba0019000000000d0d04330000000000dc0435000000200bb0003900000000008b004b000020510000413d000000000a68001900000000000a0435000000000909043300000c7109900197000000190b000029000000640ab0003900000000009a0435000000000707043300000c72077001970000008409b0003900000000007904350000000005050433000000a407b000390000000000570435000000000404043300000c7204400197000000c405b0003900000000004504350000001f0480003900000d0e04400197000000e405b00039000000000303043300000120074000390000000000750435000000000664001900000000540304340000000003460436000000000004004b0000207d0000613d000000000600001900000000073600190000000008650019000000000808043300000000008704350000002006600039000000000046004b000020760000413d000000000534001900000000000504350000001f0440003900000d0e04400197000000000534001900000019040000290000000003450049000000440330008a00000000020204330000010404400039000000000034043500000000430204340000000002350436000000000003004b000020940000613d000000000500001900000000062500190000000007540019000000000707043300000000007604350000002005500039000000000035004b0000208d0000413d000000000423001900000000000404350000001f0330003900000d0e03300197000000000223001900000019040000290000000003420049000000440330008a00000000010104330000012404400039000000000034043500000000310104340000000002120436000000000001004b000020ab0000613d000000000400001900000000052400190000000006430019000000000606043300000000006504350000002004400039000000000014004b000020a40000413d00000000032100190000000000030435000000190400002900000000024200490000001f0110003900000d0e011001970000000001210019000000200210008a00000000002404350000001f0110003900000d0e011001970000000002410019000000000012004b00000000010000390000000101004039001800000002001d00000c710020009c000000520000213d0000000100100190000000520000c13d0000001801000029000000400010043f00000cb30010009c000000520000213d0000001803000029000000c001300039000000400010043f0000008402000039000000000323043600000000020000310000000202200367001600000003001d000000002402043c0000000003430436000000000013004b000020cb0000c13d00000cd10100004100000000001004430000001c010000290000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000000001004b0000265b0000613d0000000001000414000013880110008c0000265f0000413d00000006021002700000000001210049000000140010006c000026630000a13d0000000001000414001200000001001d0000001c01000029000000040010008c000020ef0000c13d000000030100036700000001040000310000000002000019000021050000013d0000000d0100002900000c6d0010009c00000c6d0100804100000040011002100000001902000029000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000140200002900000c6d0020009c00000c6d02008041000000c002200210000000000121019f0000001c0200002931b131a70000040f000000010220015f0003000000010355000000600310027000010c6d0030019d00000c6d043001970000000003000414000000840040008c000000840400803900000018050000290000000000450435000000e0064001900000001605600029000021130000613d000000000701034f0000001608000029000000007907043c0000000008980436000000000058004b0000210f0000c13d0000001f04400190000021200000613d000000000161034f0000000304400210000000000605043300000000064601cf000000000646022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000161019f00000000001504350000000100200190000026530000c13d00000018010000290000000001010433000000200010008c000026700000c13d0000001c02000029000000060120014f00000016020000290000000002020433001c00000002001d00000c7200100198000021b10000613d0000001201300069001400140010007300000b280000413d000000400300043d000000200230003900000ccf01000041000d00000002001d000000000012043500000024013000390000001102000029000000000021043500000024010000390000000000130435001600000003001d00000cb00030009c000000520000213d00000016020000290000006001200039001900000001001d000000400010043f00000cd00020009c000000520000213d00000016040000290000012001400039000000400010043f000000840200003900000019030000290000000000230435000000800340003900000000020000310000000202200367001800000003001d000000002402043c0000000003430436000000000013004b0000214d0000c13d00000cd10100004100000000001004430000001a010000290000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000025790000613d000000000101043b000000000001004b0000265b0000613d0000000001000414000013880110008c0000265f0000413d00000006021002700000000001210049000000140010006c000026630000a13d0000000001000414001200000001001d0000001a01000029000000040010008c000021710000c13d000000030100036700000001040000310000000002000019000021870000013d0000000d0100002900000c6d0010009c00000c6d0100804100000040011002100000001602000029000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000140200002900000c6d0020009c00000c6d02008041000000c002200210000000000121019f0000001a0200002931b131a70000040f000000010220015f0003000000010355000000600310027000010c6d0030019d00000c6d043001970000000003000414000000840040008c000000840400803900000019050000290000000000450435000000e0064001900000001805600029000021950000613d000000000701034f0000001808000029000000007907043c0000000008980436000000000058004b000021910000c13d0000001f04400190000021a20000613d000000000161034f0000000304400210000000000605043300000000064601cf000000000646022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000161019f00000000001504350000000100200190000026670000c13d00000019010000290000000001010433000000200010008c000026700000c13d0000001201300069000000140010006c00000b280000213d00000018010000290000000001010433000000130210006c000026780000413d0000001c0020006c000026780000c13d000000400100043d00000c880010009c000000520000213d0000004002100039000000400020043f00000020021000390000001c0300002900000000003204350000001a02000029000000000021043500000000020004150000000e0220006900000000020000020000000f0200002900000000020204330000001b0020006c000025730000a13d00000017030000290000001d0230002900000000001204350000000f0100002900000000010104330000001b0010006c000025730000a13d000000090100002900000000010104330000001b02000029001b00010020003d0000001b0010006b00001deb0000413d00001d300000013d000000a00100043d0000000001010433000000000001004b000021d70000c13d000000400100043d00000cf50200004100000e0c0000013d000000a00200043d0000000001020433000000000001004b0000233b0000613d000d00000000001d0000000d010000290000000501100210000000000112001900000020011000390000000001010433001000000001001d0000000021010434000c00000002001d000000400300043d00000cf702000041000000000023043500000c7101100197001c00000001001d0000008001100210001d00000003001d0000000402300039000000000012043500000cc90100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000025790000613d000000000201043b000000000100041400000c7202200197000000040020008c000022060000c13d0000000103000031000000200030008c000000200400003900000000040340190000222f0000013d0000001d0300002900000c6d0030009c00000c6d03008041000000400330021000000c6d0010009c00000c6d01008041000000c001100210000000000131019f00000c85011001c731b131ac0000040f000000600310027000000c6d03300197000000200030008c0000002004000039000000000403401900000020064001900000001d056000290000221e0000613d000000000701034f0000001d08000029000000007907043c0000000008980436000000000058004b0000221a0000c13d0000001f074001900000222b0000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000026880000613d0000001f01400039000000600210018f0000001d01200029000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f000000200030008c0000004c0000413d0000001d020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b0000004c0000c13d000000000002004b000026940000c13d0000001c01000029000000000010043f0000000801000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000e00000001001d000000000101041a00000c8c00100198000026960000613d0000000c010000290000000001010433000000200210003900000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000000e020000290000000102200039000a00000002001d000000000202041a000000010320019000000001052002700000007f0550618f000000400400043d000f00000004001d000000000101043b000b00000001001d001d00000005001d0000001f0050008c00000000010000390000000101002039000000000112013f000000010010019000000b860000c13d0000000f010000290000001d040000290000000001410436001c00000001001d000000000003004b0000229b0000613d0000000a01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001d0000006b000022a20000613d000000000201043b00000000010000190000001c03100029000000000402041a0000000000430435000000010220003900000020011000390000001d0010006c000022930000413d000022a30000013d00000d10012001970000001c0200002900000000001204350000001d0000006b00000020010000390000000001006039000022a30000013d00000000010000190000003f0110003900000d0e021001970000000f01200029000000000021004b0000000002000039000000010200403900000c710010009c000000520000213d0000000100200190000000520000c13d000000400010043f0000001c0100002900000c6d0010009c00000c6d0100804100000040011002100000000f02000029000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000000b0010006b000026990000c13d00000010020000290000006001200039001c00000001001d000000000101043300000c71011001970000004002200039000000000202043300000c71022001970000000e03000029000000000303041a000000a80330027000000c7103300197000000000023004b000026d90000c13d000000000013004b0000000002030019000026d90000213d000000100100002900000080011000390000000001010433001d00000001001d000000000001004b000026ea0000613d0000001001000029000000000101043300000c7101100197000000000010043f0000000a01000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001d02000029000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000101041a000000000001004b000026ed0000c13d0000001c01000029000000000101043300000c710110019700000c710010009c00000b280000613d0000000e03000029000000000203041a00000c8602200197000000a80110021000000c8a0110009a00000c8001100197000000000112019f000000000013041b00000cfb010000410000000000100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c78011001c70000800b0200003931b131ac0000040f0000000100200190000025790000613d000000000101043b001c00000001001d0000001001000029000000000101043300000c7101100197000000000010043f0000000a01000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001d02000029000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001c02000029000000000021041b0000000d03000029000d00010030003d000000a00200043d00000000010204330000000d0010006b000021dc0000413d000000800100043d001d00000001001d000000400100043d00000040030000390000000003310436001c00000003001d000000000702043300000005037002100000000003310019000000400810003900000000007804350000006003300039000000000007004b000026f90000c13d00000000021300490000001c0400002900000000002404350000001d020000290000000046020434000000400200003900000000022304360000004005300039000000000706043300000000007504350000006005300039000000000007004b000023630000613d00000000080000190000002006600039000000000906043300000000a909043400000c72099001970000000009950436000000000a0a043300000cf30aa001970000000000a9043500000040055000390000000108800039000000000078004b000023570000413d00000000040404330000000003350049000000000032043500000000030404330000000002350436000000000003004b000023770000613d000000000500001900000020044000390000000006040433000000007606043400000c71066001970000000006620436000000000707043300000cf307700197000000000076043500000040022000390000000105500039000000000035004b0000236b0000413d000000000212004900000c6d0020009c00000c6d02008041000000600220021000000c6d0010009c00000c6d010080410000004001100210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000121019f00000c83011001c70000800d02000039000000010300003900000cff0400004131b131a70000040f00000001002001900000004c0000613d00000015010000290000003f0110003900000c7301100197000000400300043d0000000002130019001c00000003001d000000000032004b0000000001000039000000010100403900000c710020009c000000520000213d0000000100100190000000520000c13d0000000001000031000000400020043f00000018020000290000001c030000290000000002230436001500000002001d000000160010006b0000004c0000213d0000000202000367000000180000006b000023ad0000613d0000001c030000290000001705000029000000000452034f000000000404043b000000200330003900000000004304350000002005500039001700000005001d000000160050006c000023a40000413d000000c403200370000000000303043b001600000003001d00000011030000290000003f0330003900000c7303300197000000400400043d0000000003340019001800000004001d000000000043004b0000000004000039000000010400403900000c710030009c000000520000213d0000000100400190000000520000c13d000000400030043f000000140300002900000018040000290000000003340436001100000003001d000000120010006b0000004c0000213d000000140000006b000023d00000613d00000018030000290000001305000029000000000452034f000000000404043b000000200330003900000000004304350000002005500039001300000005001d000000120050006c000023c70000413d000000000000043f0000000203000039000000200030043f000000400300043d00000c700030009c000000520000213d0000008004300039000000400040043f00000d0004000041000000000404041a001700000004001d000000000743043600000d0104000041000000000404041a0000000805400270000000ff0550018f00000040063000390000000000560435000000ff0540018f001400000007001d0000000000570435000000600530003900000cb2044001980000000003000039000000010300c039001300000005001d00000000003504350000000402200370000000000202043b000000a50300008a0000001b0030006b00000b280000213d0000001b03000029000000a403300039000000000004004b000024100000613d0000001c0400002900000000050404330000000504500210000000000005004b000023fe0000613d00000d120040009c00000b280000213d00000000055400d9000000200050008c00000b280000c13d000000180500002900000000060504330000000505600210000000000006004b000024060000613d00000000066500d9000000200060008c00000b280000c13d000000a00640003900000000046500190000000003340019000000000043004b00000000040000390000000104004039000000000065001a00000b280000413d000000010040019000000b280000c13d000000000031004b0000272c0000c13d000000170020006b000027330000c13d00000cc9010000410000000000100443000000000100041200000004001004430000002400000443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000025790000613d000000000101043b001d00000001001d00000c77010000410000000000100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c78011001c70000800b0200003931b131ac0000040f0000000100200190000025790000613d000000000101043b0000001d0010006b0000257a0000c13d000000000100041100000c7201100197000000000010043f00000d0201000041000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000400200043d00000c880020009c000000520000213d000000000101043b0000004003200039000000400030043f000000000301041a000000ff0130018f00000000011204360000000803300270000000ff0330018f000000030030008c000017df0000813d0000000000310435000000020030008c000025870000c13d000000000000043f0000000201000039000000200010043f0000000001020433001d00ff0010019300000d0301000041000000000201041a0000001d0020006c000025730000a13d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b0000001d01100029000000000101041a00000c72011001970000000002000411000000000012004b000025870000c13d00000013010000290000000001010433000000000001004b0000255c0000613d00000014010000290000000001010433000000ff0110018f000000ff0010008c00000b280000613d0000001c0200002900000000020204330000000101100039000000000012004b0000258a0000c13d000000400100043d00000018030000290000000003030433000000000032004b0000273c0000c13d0000001b020000290000001f0220003900000d0e022001970000003f0220003900000d0e022001970000000002210019000000000012004b0000000004000039000000010400403900000c710020009c000000520000213d0000000100400190000000520000c13d000000400020043f0000001b0200002900000000022104360000001a05000029000000000050007c0000004c0000213d0000001b0500002900000d0e045001980000001f0550018f0000001903000029000000020630036700000000034200190000249d0000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000024990000c13d000000000005004b000024aa0000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004304350000001b03200029000000000003043500000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000301043b000000400100043d00000020021000390000000000320435000000a003100039000000400410003900000004050000390000000205500367000000005605043c0000000004640436000000000034004b000024c60000c13d0000008004000039000000000041043500000cc50010009c000000520000213d000000400030043f00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b001a00000001001d0000001c010000290000000001010433001900000001001d000000000001004b0000255c0000613d001d00000000001d001b00000000001d0000001d010000290000001f0010008c000025730000213d0000001d01000029000000030110021000000016011001ef00000d070010009c00000b280000213d0000001c0200002900000000020204330000001d0020006c000025730000a13d000000180200002900000000020204330000001d0020006c000025730000a13d000000f8011002700000001b011000390000001d0200002900000005022002100000001103200029000000150220002900000000020204330000000003030433000000400400043d0000006005400039000000000035043500000040034000390000000000230435000000200240003900000000001204350000001a010000290000000000140435000000000000043f00000c6d0040009c00000c6d040080410000004001400210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000cf1011001c7000000010200003931b131ac0000040f000000600310027000000c6d03300197000000200030008c000000200500003900000000050340190000002004500190000025240000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b000025200000c13d0000001f05500190000025310000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f0000000000540435000100000003001f000300000001035500000001002001900000273e0000613d000000000100043d00000c7201100197000000000010043f00000d0201000041000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000400200043d00000c880020009c000000520000213d000000000101043b0000004003200039000000400030043f000000000301041a000000ff0130018f00000000021204360000000803300270000000ff0330018f000000020030008c000017df0000213d0000000000320435000000010030008c000025e20000c13d000000010110020f0000001b00100180000025e50000c13d001b001b001001b30000001d020000290000000102200039001d00000002001d000000190020006c000024ea0000413d00000024010000390000000201100367000000000101043b00000c7101100197000000400200043d000000200320003900000000001304350000001701000029000000000012043500000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c7f011001c70000800d02000039000000020300003900000cf204000041000000000500001900000b460000013d00000d0a01000041000000000010043f0000003201000039000000040010043f00000c8501000041000031b300010430000000000001042f000000400200043d0000002403200039000000000013043500000ce401000041000000000012043500000004012000390000001d03000029000000000031043500000c6d0020009c00000c6d02008041000000400120021000000cd4011001c7000031b300010430000000400100043d00000d040200004100000e0c0000013d000000400100043d00000d050200004100000e0c0000013d0000001f0530018f00000c6f06300198000000400200043d0000000004620019000028f80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000025940000c13d000028f80000013d31b1294f0000040f00000cdf02000041000000400300043d001d00000003001d00000000002304350000000002010019000028dc0000013d00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f00000ccc011001c7000000170200002931b131ac0000040f000000600310027000000c6d03300197000000200030008c000000200400003900000000040340190000001f0540018f0000002004400190000025b90000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b000025b50000c13d000000000005004b000025c60000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f0000000000540435000100000003001f00030000000103550000000100200190000028620000613d000000000100043d000000200030008c000028620000413d000000000001004b000028620000613d000000400100043d000000200210003900000ccb040000410000000000420435000000240410003900000ccd0500004100000000005404350000002404000039000000000041043500000cb00010009c000000520000213d0000006004100039000000400040043f0000001704000029000000040040008c0000274a0000c13d0000000001020433000000000010043f0000277c0000013d000000400100043d00000d080200004100000e0c0000013d000000400100043d00000d090200004100000e0c0000013d000000a0050000390000000007000019000026010000013d000000000ba9001900000000000b0435000000400b800039000000000b0b043300000c710bb00197000000400c3000390000000000bc0435000000600b800039000000000b0b043300000c710bb00197000000600c3000390000000000bc043500000080033000390000008008800039000000000808043300000000008304350000001f0390003900000d0e033001970000000003a300190000000107700039000000000027004b000006bc0000813d000000100830006a000000840880008a000000200440003900000000008404350000001d080000290000002008800039001d00000008001d000000000808043300000000a908043400000c71099001970000000009930436000000000a0a04330000000000590435000000a00c30003900000000b90a043400000000009c0435000000c00a300039000000000009004b000025eb0000613d000000000c000019000000000dac0019000000000ecb0019000000000e0e04330000000000ed0435000000200cc0003900000000009c004b000026150000413d000025eb0000013d00000c6d033001970000001f0530018f00000c6f06300198000000400200043d0000000004620019000028f80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000026250000c13d000028f80000013d00000c6d033001970000001f0530018f00000c6f06300198000000400200043d0000000004620019000028f80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000026320000c13d000028f80000013d0000001f0530018f00000c6f06300198000000400200043d0000000004620019000028f80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000263e0000c13d000028f80000013d00000018010000290000000501100270000000000100003f000000000100041500000016011000690000000001000002000000400100043d00000cdd0200004100001b620000013d000000400200043d001d00000002001d00000cda01000041000000000012043500000004012000390000001a02000029000026590000013d000000400200043d001d00000002001d00000cda01000041000000000012043500000004012000390000001c0200002900000018030000290000266e0000013d00000cdc01000041000000000010043f00000c9001000041000031b30001043000000cdb01000041000000000010043f00000c9001000041000031b30001043000000cd701000041000000000010043f00000c9001000041000031b300010430000000400200043d001d00000002001d00000cda01000041000000000012043500000004012000390000001a02000029000000190300002931b1315a0000040f000028de0000013d000000400200043d0000002403200039000000000013043500000cd301000041000000000012043500000004012000390000002003000039000025810000013d000000400200043d0000004403200039000000000013043500000024012000390000001303000029000000000031043500000cd801000041000000000012043500000004012000390000001c03000029000000000031043500000c6d0020009c00000c6d02008041000000400120021000000cd9011001c7000031b3000104300000001f0530018f00000c6f06300198000000400200043d0000000004620019000028f80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000268f0000c13d000028f80000013d00000cf80200004100001b620000013d000000400100043d00000cfe0200004100001b620000013d0000000c010000290000000002010433000000400500043d00000cf9010000410000000000150435000000040150003900000040030000390000000000310435000000440350003900000000420204340000000000230435001b00000005001d0000006403500039000000000002004b000026b00000613d000000000500001900000000063500190000000007540019000000000707043300000000007604350000002005500039000000000025004b000026a90000413d000000000432001900000000000404350000001f0220003900000d0e02200197000000000232001900000000011200490000001b03000029000000240330003900000000001304350000000a01000029000000000101041a000000010310019000000001041002700000007f0440618f001d00000004001d0000001f0040008c00000000040000390000000104002039000000000441013f000000010040019000000b860000c13d0000001d040000290000000002420436001c00000002001d000000000003004b000027940000613d0000000a01000029000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d0000001d0000006b0000279b0000c13d0000000001000019000027a40000013d00000010030000290000000003030433000000400400043d000000440540003900000000001504350000002401400039000000000021043500000cfd01000041000000000014043500000c71013001970000000402400039000000000012043500000c6d0040009c00000c6d04008041000000400140021000000cd9011001c7000031b300010430000000400100043d00000cfc0200004100000e0c0000013d00000010010000290000000001010433000000400200043d00000024032000390000001d04000029000000000043043500000cfa03000041000000000032043500000c710110019700000004032000390000000000130435000025820000013d000000a009000039000000000b000019000027120000013d0000000004ed001900000000000404350000004004c00039000000000404043300000c7104400197000000400530003900000000004504350000006004c00039000000000404043300000c71044001970000006005300039000000000045043500000080033000390000008004c00039000000000404043300000000004304350000001f03d0003900000d0e033001970000000003e30019000000010bb0003900000000007b004b000023490000813d0000000005130049000000600550008a000000200880003900000000005804350000002002200039000000000c02043300000000d50c043400000c71055001970000000005530436000000000d0d04330000000000950435000000a00530003900000000fd0d04340000000000d50435000000c00e30003900000000000d004b000026fc0000613d00000000050000190000000004e5001900000000065f00190000000006060433000000000064043500000020055000390000000000d5004b000027240000413d000026fc0000013d000000400200043d0000002404200039000000000014043500000ced0100004100000000001204350000000401200039000025810000013d000000400100043d0000002403100039000000000023043500000cee020000410000000000210435000000040210003900000017030000290000000000320435000017ec0000013d00000d060200004100000e0c0000013d0000001f0530018f00000c6f06300198000000400200043d0000000004620019000028f80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000027450000c13d000028f80000013d00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f00000ccc011001c7000000170200002931b131ac0000040f000000600310027000000c6d03300197000000200030008c000000200400003900000000040340190000001f0540018f0000002004400190000027630000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b0000275f0000c13d000000000005004b000027700000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f00000000005404350003000000010355000100000003001f0000001f0030008c00000000010000390000000101002039000000000112016f00000000020004150000001f0220008a0016000500200218000000010010008c000027810000c13d000000000100043d00000000020004150000001e0220008a0016000500200218000000000001004b000028620000c13d000000400100043d000000200210003900000ccb040000410000000000420435000000240410003900000ce00500004100000000005404350000002404000039000000000041043500000cb00010009c000000520000213d0000006004100039000000400040043f0000001704000029000000040040008c000027af0000c13d0000000001020433000000000010043f000027da0000013d00000d10011001970000001c0200002900000000001204350000001d0000006b00000020010000390000000001006039000027a40000013d000000000201043b00000000010000190000001c03100029000000000402041a0000000000430435000000010220003900000020011000390000001d0010006c0000279d0000413d0000001b030000290000001c02300069000000000112001900000c6d0010009c00000c6d01008041000000600110021000000c6d0030009c00000c6d030080410000004002300210000000000121019f000031b30001043000000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f00000ccc011001c7000000170200002931b131ac0000040f000000600310027000000c6d03300197000000200030008c000000200400003900000000040340190000001f0540018f0000002004400190000027c80000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b000027c40000c13d000000000005004b000027d50000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f0000000000540435000100000003001f000300000001035500000001002001900000285f0000613d000000000100043d000000200030008c0000285f0000413d000000000001004b00000016010000290000000501100270000000000100003f000000010100c03f00000000010004150000001c011000690000000001000002000028e80000613d000000800100043d0000002001100039000000000101043300000c7101100197000000000010043f0000000801000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000004c0000613d000000000101043b000000000201041a000001000100043d000000e00400043d000000400800043d00000ce1030000410000000003380436001c00000003001d0000000403800039000000800500003900000000005304350000001d030000290000000005030433000000840380003900000000005304350000001805000029000000000505043300000c7105500197000000a406800039000000000056043500000019050000290000000005050433000000c406800039000000a0070000390000000000760435000001240680003900000000750504340000000000560435001d00000008001d0000014406800039000000000005004b0000281d0000613d00000000080000190000000009680019000000000a870019000000000a0a04330000000000a904350000002008800039000000000058004b000028160000413d000000000765001900000000000704350000001f0550003900000d0e05500197000000000565001900000000063500490000001a0700002900000000070704330000001d08000029000000e408800039000000000068043500000000760704340000000005650436000000000006004b000028340000613d00000000080000190000000009580019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b0000282d0000413d00000c720440019700000c7202200197000000000756001900000000000704350000001f0660003900000d0e06600197000000000756001900000000033700490000001b0500002900000000050504330000001d060000290000010406600039000000000036043500000000060504330000000003670436000000000006004b000028510000613d000000000700001900000020055000390000000008050433000000009808043400000c720880019700000000088304360000000009090433000000000098043500000040033000390000000107700039000000000067004b000028460000413d0000001d0600002900000064056000390000000000450435000000440460003900000000001404350000002401600039000013880400003900000000004104350000000001000414000000040020008c000028660000c13d000000030100036700000001030000310000287a0000013d00000016010000290000000501100270000000000100003f00000000010004150000001c011000690000000001000002000028e80000013d0000001d04000029000000000343004900000c6d0030009c00000c6d03008041000000600330021000000c6d0040009c00000c6d040080410000004004400210000000000343019f00000c6d0010009c00000c6d01008041000000c001100210000000000131019f31b131a70000040f000000600310027000010c6d0030019d00000c6d0330019700030000000103550000000100200190000028ed0000613d00000d0e043001980000001f0530018f0000001d02400029000028840000613d000000000601034f0000001d07000029000000006806043c0000000007870436000000000027004b000028800000c13d000000000005004b000028910000613d000000000141034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f00000000001204350000001f0130003900000d0e011001970000001d02100029000000000012004b0000000001000039000000010100403900000c710020009c000000520000213d0000000100100190000000520000c13d000000400020043f00000c740030009c0000004c0000213d000000600030008c0000004c0000413d0000001d010000290000000001010433000000000001004b0000000004000039000000010400c039000000000041004b0000004c0000c13d0000001c04000029000000000404043300000c710040009c0000004c0000213d0000001d053000290000001d034000290000001f04300039000000000054004b000000000600001900000c750600804100000c750440019700000c7507500197000000000874013f000000000074004b000000000400001900000c750400404100000c750080009c000000000406c019000000000004004b0000004c0000c13d000000004303043400000c710030009c000000520000213d0000001f0630003900000d0e066001970000003f0660003900000d0e06600197000000000626001900000c710060009c000000520000213d000000400060043f00000000063204360000000007430019000000000057004b0000004c0000213d000000000003004b000028d40000613d000000000500001900000000076500190000000008450019000000000808043300000000008704350000002005500039000000000035004b000028cd0000413d00000000036300190000000000030435000000000001004b000028e80000c13d000000400300043d001d00000003001d00000ce2010000410000000000130435000000040130003931b129820000040f0000001d02000029000000000121004900000c6d0010009c00000c6d01008041000000600110021000000c6d0020009c00000c6d020080410000004002200210000000000121019f000031b3000104300000000001000415000000020110006900000000010000020000000001000019000031b20001042e0000001f0530018f00000c6f06300198000000400200043d0000000004620019000028f80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000028f40000c13d000000000005004b000029050000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f0000000000140435000000600130021000000c6d0020009c00000c6d020080410000004002200210000000000112019f000031b30001043000000d130010009c000029100000813d0000008001100039000000400010043f000000000001042d00000d0a01000041000000000010043f0000004101000039000000040010043f00000c8501000041000031b30001043000000000430104340000000001320436000000000003004b000029220000613d000000000200001900000000051200190000000006240019000000000606043300000000006504350000002002200039000000000032004b0000291b0000413d000000000213001900000000000204350000001f0230003900000d0e022001970000000001210019000000000001042d000000004302043400000c72033001970000000003310436000000000404043300000c6d04400197000000000043043500000040032000390000000003030433000000000003004b0000000003000039000000010300c039000000400410003900000000003404350000006002200039000000000202043300000c7202200197000000600310003900000000002304350000008001100039000000000001042d000000400100043d00000d130010009c000029490000813d0000008002100039000000400020043f0000006002100039000000000002043500000040021000390000000000020435000000200210003900000000000204350000000000010435000000000001042d00000d0a01000041000000000010043f0000004101000039000000040010043f00000c8501000041000031b30001043000000001020000320000297a0000613d00000d140020009c0000297c0000813d0000001f0120003900000d0e011001970000003f0110003900000d0e04100197000000400100043d0000000004410019000000000014004b0000000005000039000000010500403900000c710040009c0000297c0000213d00000001005001900000297c0000c13d000000400040043f000000000621043600000d0e032001980000001f0420018f000000000236001900000003050003670000296c0000613d000000000705034f000000007807043c0000000006860436000000000026004b000029680000c13d000000000004004b0000297b0000613d000000000335034f0000000304400210000000000502043300000000054501cf000000000545022f000000000303043b0000010004400089000000000343022f00000000034301cf000000000353019f0000000000320435000000000001042d0000006001000039000000000001042d00000d0a01000041000000000010043f0000004101000039000000040010043f00000c8501000041000031b30001043000000020030000390000000004310436000000003202043400000000002404350000004001100039000000000002004b000029910000613d000000000400001900000000051400190000000006430019000000000606043300000000006504350000002004400039000000000024004b0000298a0000413d000000000312001900000000000304350000001f0220003900000d0e022001970000000001120019000000000001042d001d000000000002000300000002001d000000400200043d000600000001001d0000000031010434000500000003001d000000000001004b000031580000613d000100000002001d00000d150020009c000030a20000813d00000003010000290000000021010434000200000002001d000400000001001d00000001020000290000002001200039000000400010043f000000000002043500000006010000290000000001010433000000000001004b000030990000613d0000000003000019000029b60000013d0000000703000029000000010330003900000006010000290000000001010433000000000013004b000030990000813d000000050130021000000005021000290000000002020433001100000002001d000000040000006b000700000003001d000029c40000613d00000003020000290000000002020433000000000032004b0000309c0000a13d00000002011000290000000001010433000029c50000013d0000000101000029001000000001001d0000000001010433000b00000001001d00000011010000290000000021010434001800000002001d000000400300043d00000cf702000041000000000023043500000c7102100197001c00000003001d0000000401300039001d00000002001d0000008002200210000000000021043500000cc90100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000b0000006b00000000030000390000000103006039000800000003001d0000000100200190000030ae0000613d000000000201043b000000000100041400000c7202200197000000040020008c000029f20000c13d0000000103000031000000200030008c000000200400003900000000040340190000001c0a00002900002a1c0000013d0000001c0300002900000c6d0030009c00000c6d03008041000000400330021000000c6d0010009c00000c6d01008041000000c001100210000000000131019f00000c85011001c731b131ac0000040f0000001c0a000029000000600310027000000c6d03300197000000200030008c00000020040000390000000004034019000000200640019000000000056a001900002a0b0000613d000000000701034f00000000080a0019000000007907043c0000000008980436000000000058004b00002a070000c13d0000001f0740019000002a180000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f000300000001035500000001002001900000310f0000613d0000001f01400039000000600110018f0000000002a10019000000000012004b00000000010000390000000101004039001a00000002001d00000c710020009c000030a20000213d0000000100100190000030a20000c13d0000001a01000029000000400010043f0000001f0030008c0000309a0000a13d00000000010a0433000000000001004b0000000002000039000000010200c039000000000021004b0000309a0000c13d000000000001004b00002a480000613d0000000b0000006b000031080000c13d0000001d010000290000001a02000029000000000012043500000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c82011001c70000800d02000039000000010300003900000d160400004131b131a70000040f0000000100200190000029b00000c13d0000309a0000013d000000180100002900000000010104330000000002010433000000000002004b0000312d0000613d00000011010000290000004001100039000a00000001001d00000000010104330000000001010433000000000012004b0000313a0000c13d001700000002001d00000c710020009c000030a20000213d000000170100002900000005011002100000003f0210003900000c73022001970000001a0220002900000c710020009c000030a20000213d000000400020043f0000001a0200002900000017030000290000000002320436000f00000002001d000000000001004b00002a6d0000613d0000000f04000029000000000214001900000000030000310000000203300367000000003503043c0000000004540436000000000024004b00002a690000c13d0000001f0010019000000cc90100004100000000001004430000000001000412000000040010044300000020010000390000002400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000030ae0000613d000000000101043b001c00000001001d0000001d01000029000000000010043f0000000801000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000400600043d000000000101043b000000000201041a00000c8c00200198000031420000613d0000000101100039000000000201041a000000010320019000000001052002700000007f0550618f0000001f0050008c00000000040000390000000104002039000000000043004b0000314c0000c13d0000000007560436000000000003004b00002abb0000613d001600000005001d001900000007001d001b00000006001d000000000010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c82011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d0000001608000029000000000008004b0000001b0600002900002ac20000613d000000000201043b0000000001000019000000200500008a00000019070000290000000003710019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00002ab30000413d00002ac50000013d00000d10012001970000000000170435000000000005004b00000020010000390000000001006039000000200500008a00002ac50000013d0000000001000019000000200500008a00000019070000290000003f01100039000000000251016f0000000001620019000000000021004b0000000002000039000000010200403900000c710010009c000030a20000213d0000000100200190000030a20000c13d000000400010043f00000c6d0070009c00000c6d070080410000004001700210000000000206043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000201043b000000400100043d0000008003100039000000000023043500000040021000390000001d0300002900000000003204350000001c0200002900000c71032001970000006002100039001300000003001d0000000000320435000000200210003900000d180300004100000000003204350000008003000039000000000031043500000cc50010009c000030a20000213d000000a003100039000000400030043f00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b001200000001001d0000000005000019000000180100002900000000020104330000000001020433000000000051004b0000309c0000a13d000000400100043d000000050350021000000020063000390000000002260019000000000802043300000000920804340000004003200039000000000303043300000c7103300197000000130030006c000030af0000c13d0000002003200039000000000303043300000c71033001970000001d0030006c000030b30000c13d001600000006001d001900000005001d0000006003800039000000000303043300000060042000390000000004040433000000000502043300000080022000390000000002020433000000800680003900000000060604330000008007100039000000000067043500000c7102200197000000a006100039000000000026043500000c71024001970000006004100039000000000024043500000c72023001970000004003100039000000000023043500000020021000390000000000520435000000a003000039000000000031043500000cb30010009c000030a20000213d001b00000009001d001c00000008001d000000c003100039000000400030043f00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d0000001b020000290000000002020433000000200320003900000c6d0030009c00000c6d030080410000004003300210000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000232019f000000000101043b001b00000001001d000000000100041400000c6d0010009c00000c6d01008041000000c001100210000000000121019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d0000001c0200002900000040022000390000000002020433000000200320003900000c6d0030009c00000c6d030080410000004003300210000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000232019f000000000101043b001500000001001d000000000100041400000c6d0010009c00000c6d01008041000000c001100210000000000121019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b001400000001001d0000001c01000029000000a0011000390000000003010433000000400100043d00000020041000390000002002000039001c00000004001d0000000000240435000000000403043300000005054002100000000006510019000000400510003900000000004504350000006006600039000000000004004b00002bd90000613d000000000700001900002ba20000013d00000000029a0019000000000002043500000080026000390000008006800039000000000606043300000000006204350000001f02a000390000000002e2016f00000000069200190000000107700039000000000047004b00002bda0000813d0000000008160049000000600880008a000000200550003900000000008504350000002003300039000000000803043300000000c9080434000000a002000039000000000b260436000000a00a60003900000000d909043400000000009a0435000000c00a600039000000000009004b00002bb90000613d000000000e000019000000000fae00190000000002ed0019000000000202043300000000002f0435000000200ee0003900000000009e004b00002bb20000413d0000000002a90019000000000002043500000000020c043300000c720220019700000000002b04350000004002800039000000000202043300000c6d02200197000000400b60003900000000002b04350000001f02900039000000200e00008a0000000002e2016f0000000002a200190000000009620049000000600a600039000000600b800039000000000b0b043300000000009a043500000000ba0b04340000000009a2043600000000000a004b00002b960000613d000000000c00001900000000029c0019000000000dcb0019000000000d0d04330000000000d20435000000200cc000390000000000ac004b00002bd10000413d00002b960000013d000000200e00008a0000000002160049000000200320008a00000000003104350000001f022000390000000002e2016f0000000003120019000000000023004b0000000004000039000000010400403900000c710030009c000030a20000213d0000000100400190000030a20000c13d000000400030043f0000001c0200002900000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000201043b000000400100043d000000c0031000390000000000230435000000a0021000390000001403000029000000000032043500000080021000390000001503000029000000000032043500000060021000390000001b030000290000000000320435000000400210003900000012030000290000000000320435000000c0020000390000000002210436000000000002043500000cb10010009c000030a20000213d000000e003100039000000400030043f00000c6d0020009c00000c6d020080410000004002200210000000000101043300000c6d0010009c00000c6d010080410000006001100210000000000121019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d0000001a0200002900000000020204330000001905000029000000000052004b0000309c0000a13d00000016030000290000001a02300029000000000101043b00000000001204350000000105500039000000170050006c00002b0c0000413d0000001a010000290000000006010433000000000006004b000031520000613d000001010060008c000030b70000213d000000110100002900000060011000390000000001010433000e00000001001d0000000021010434000d00000002001d000c00000001001d000001020010008c000030b70000813d0000000c01600029000900000001001d000000010710008a000001000070008c000030b70000213d000000000007004b00002cfa0000613d00000005017002100000003f0210003900000c7302200197000000400300043d0000000002230019001600000003001d000000000032004b0000000003000039000000010300403900000c710020009c000030a20000213d0000000100300190000030a20000c13d000000110300002900000080033000390000000003030433001300000003001d000000400020043f00000016020000290000000002720436001500000002001d000000000001004b00002c650000613d0000001504000029000000000214001900000000030000310000000203300367000000003503043c0000000004540436000000000024004b00002c610000c13d0000001f00100190000000000a000019001400000000001d00000000090000190000000008000019001200000006001d001100000007001d0000000101a0020f000000ff00a0008c0000000001002019000000130210017f000000000012004b00002c7c0000c13d000000000069004b00002c850000813d0000001a010000290000000001010433000000000091004b0000309c0000a13d000000050190021000000001099000390000000f0200002900002c8c0000013d0000000e010000290000000001010433000000140010006c0000309c0000a13d00000014020000290000000501200210001400010020003d0000000d0200002900002c8c0000013d00000016010000290000000001010433000000000081004b0000309c0000a13d00000005028002100000000108800039000000150100002900000000011200190000000001010433000000000069004b00002c980000813d0000001a020000290000000002020433000000000092004b0000309c0000a13d000000050290021000000001099000390000000f0300002900002c9f0000013d00000016020000290000000002020433000000000082004b0000309c0000a13d0000000503800210000000010880003900000015020000290000000000a8004b000030b70000213d00000000022300190000000004020433000000400200043d00000020032000390000000105000039000000000053043500000060052000390000004006200039000000000041004b001c00000008001d001b00000009001d00190000000a001d00002cc00000813d000000000016043500000000004504350000006001000039000000000012043500000c700020009c000030a20000213d0000008001200039000000400010043f00000c6d0030009c00000c6d030080410000004001300210000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000000200041400002cd10000013d000000000046043500000000001504350000006001000039000000000012043500000c700020009c000030a20000213d0000008001200039000000400010043f00000c6d0030009c00000c6d030080410000004001300210000000000202043300000c6d0020009c00000c6d020080410000006002200210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000120600002900000016020000290000000002020433000000190a0000290000000000a2004b00000011030000290000001c080000290000001b090000290000309c0000a13d0000000502a002100000001502200029000000000101043b0000000000120435000000010aa0003900000000003a004b00002c6c0000413d0000000901000029000000020110008a000000000018004b000030b70000c13d000000000069004b000030b70000c13d00000014020000290000000c0020006c000030b70000c13d00000016010000290000000001010433000000000081004b0000309c0000a13d0000000501800210000000150110002900002cfb0000013d0000000f010000290000000001010433001c00000001001d0000001d01000029000000000010043f0000000a01000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b0000001c02000029000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b000000000101041a000900000001001d000000000001004b000031550000613d001b00000000001d00002d3e0000013d00000000044300190000000000040435000000600410003900000000002404350000001f0230003900000d0e02200197000000a00220003900000c6d0020009c00000c6d02008041000000600220021000000c6d0010009c00000c6d010080410000004001100210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c70000800d02000039000000040300003900000d23040000410000001d0500002931b131a70000040f00000001002001900000309a0000613d0000001b020000290000000102200039001b00000002001d000000170020006c000029b00000813d0000000001000414001200000001001d0000001801000029000000000101043300000000020104330000001b03000029000000000032004b0000309c0000a13d0000000502300210001100000002001d001500200020003d00000015011000290000000001010433001c00000001001d0000000021010434001300000002001d00000060011000390000000001010433001900000001001d0000001d01000029000000000010043f0000000901000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000190200002900000c7102200197000000000101043b001600000002001d0000000702200270000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d00000016020000290000000102200210000000000101043b000000000101041a00000019030000290000007f0330019000002d790000613d000000ff0420018f00000000033400d9000000020030008c000030a80000c13d000000fe0220018f000000000121022f0000000302100190001900000002001d00002d990000613d000000030020008c00002e0e0000c13d0000000b0000006b00002d9c0000c13d0000001c0100002900000000010104330000006001100039000000000101043300000c7101100197000000400200043d000000200320003900000000001304350000001d01000029000000000012043500000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c7f011001c70000800d02000039000000010300003900000d1c0400004100002e240000013d0000000b0000006b000000600300003900002dcc0000613d000000100100002900000000010104330000001b0010006c0000309c0000a13d00000015020000290000001001200029001600000001001d000000000101043300000020011000390000000001010433001400000001001d00000cfb010000410000000000100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c78011001c70000800b0200003931b131ac0000040f0000000100200190000030ae0000613d000000000101043b000000090110006c000030a80000413d0000001902000029000000030020008c00002dbe0000613d0000000402000039000000000202041a000000a00220027000000c6d02200197000000000021004b000030f10000a13d000000100100002900000000010104330000001b0010006c0000309c0000a13d000000160100002900000000010104330000000001010433000000000001004b000000140300002900002dcb0000613d0000001c02000029000000800220003900000000001204350000001902000029001400000003001d000000000002004b00002e6b0000c13d0000001c0100002900000000010104330000008001100039000000000101043300000c710110019800002e6b0000613d00000013020000290000000002020433000000400500043d0000004403500039000000600400003900000000004304350000002403500039000000000013043500000d1e01000041000000000015043500000004015000390000001d030000290000000000310435000000640350003900000000160204340000000000630435001600000005001d0000008402500039000000000006004b00002df10000613d000000000300001900000000042300190000000005310019000000000505043300000000005404350000002003300039000000000063004b00002dea0000413d000e00000006001d0000000001260019000000000001043500000cc90100004100000000001004430000000001000412000000040010044300000080010000390000002400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cca011001c7000080050200003931b131ac0000040f0000000100200190000030ae0000613d000000000201043b000000000100041400000c7202200197000000040020008c000000160a00002900002e250000c13d0000000103000031000000200030008c0000002004000039000000000403401900002e560000013d0000001c0100002900000000010104330000006001100039000000000101043300000c7101100197000000400200043d000000200320003900000000001304350000001d01000029000000000012043500000c6d0020009c00000c6d020080410000004001200210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c7f011001c70000800d02000039000000010300003900000d1b0400004100002d360000013d0000000e030000290000001f0330003900000d0e0330019700000c6d00a0009c00000c6d0400004100000000040a40190000004004400210000000840330003900000c6d0030009c00000c6d030080410000006003300210000000000343019f00000c6d0010009c00000c6d01008041000000c001100210000000000131019f31b131a70000040f000000160a000029000000600310027000000c6d03300197000000200030008c00000020040000390000000004034019000000200640019000000000056a001900002e450000613d000000000701034f00000000080a0019000000007907043c0000000008980436000000000058004b00002e410000c13d0000001f0740019000002e520000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000030fc0000613d0000001f01400039000000600210018f0000000001a20019000000000021004b0000000002000039000000010200403900000c710010009c000030a20000213d0000000100200190000030a20000c13d000000400010043f000000200030008c0000309a0000413d00000000010a0433000000000001004b0000000002000039000000010200c039000000000021004b0000309a0000c13d000000000001004b00002d390000613d0000000a01000029000000000101043300000000020104330000001b0020006c0000309c0000a13d00000015021000290000001c0300002900000000010304330000006001100039000000000101043300000c710610019700000000040204330000000002040433000000a00530003900000000030504330000000003030433000000000023004b000030bf0000c13d000d00000005001d000e00000004001d001600000006001d00000001036002100000007f0110019000002e870000613d000000ff0230018f00000000011200d9000000020010008c000030a80000c13d001500000003001d0000001d01000029000000000010043f0000000901000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b00000016020000290000000702200270001600000002001d000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b000000000101041a000c00000001001d0000001d01000029000000000010043f0000000901000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b0000001602000029000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d0000001502000029000000fe0220018f000000030320020f00000d11033001670000000c0330017f000000010220020f000000000223019f000000000101043b000000000021041b00000cd101000041000000000010044300000000010004100000000400100443000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000cd2011001c7000080020200003931b131ac0000040f0000000100200190000030ae0000613d000000000101043b000000000001004b0000309a0000613d000000400800043d00000d200100004100000000071804360000000401800039000000600200003900000000002104350000001c010000290000000001010433000000003201043400000064048000390000000000240435000000000203043300000c7102200197000000840380003900000000002304350000004002100039000000000202043300000c7102200197000000a40380003900000000002304350000006002100039000000000202043300000c7102200197000000c40380003900000000002304350000008001100039000000000101043300000c7101100197000000e402800039000000000012043500000013010000290000000001010433000001040280003900000140030000390000000000320435000001a40480003900000000320104340000000000240435000001c401800039000000000002004b00002f0b0000613d000000000400001900000000051400190000000006430019000000000606043300000000006504350000002004400039000000000024004b00002f040000413d000000000312001900000000000304350000001f02200039000000200e00008a0000000002e2016f0000001c0300002900000040033000390000000003030433000001240480003900000160052000390000000000540435000000000112001900000000320304340000000001210436000000000002004b00002f230000613d000000000400001900000000051400190000000006430019000000000606043300000000006504350000002004400039000000000024004b00002f1c0000413d001600000007001d000000000312001900000000000304350000001c050000290000006003500039000000000303043300000c72033001970000014404800039000000000034043500000080035000390000000003030433000001640480003900000000003404350000001f022000390000000002e2016f00000000021200190000000001820049000000640110008a001500000008001d00000184048000390000000d030000290000000003030433000000000014043500000000040304330000000000420435000000050140021000000000011200190000002001100039000000000004004b00002f870000613d0000000005000019000000000602001900002f500000013d000000000a89001900000000000a043500000080011000390000008007700039000000000707043300000000007104350000001f019000390000000001e1016f00000000018100190000000105500039000000000045004b00002f870000813d0000000007210049000000200770008a000000200660003900000000007604350000002003300039000000000703043300000000b8070434000000a009000039000000000a910436000000a00910003900000000c80804340000000000890435000000c009100039000000000008004b00002f670000613d000000000d000019000000000e9d0019000000000fdc0019000000000f0f04330000000000fe0435000000200dd0003900000000008d004b00002f600000413d000000000c98001900000000000c0435000000000b0b043300000c720bb001970000000000ba0435000000400a700039000000000a0a043300000c6d0aa00197000000400b1000390000000000ab04350000001f08800039000000200e00008a0000000008e8016f00000000089800190000000009180049000000600a100039000000600b700039000000000b0b043300000000009a043500000000a90b04340000000008980436000000000009004b00002f440000613d000000000b000019000000000c8b0019000000000dba0019000000000d0d04330000000000dc0435000000200bb0003900000000009b004b00002f7f0000413d00002f440000013d000000150b0000290000000002b10049000000040220008a0000002403b0003900000000002304350000000e0c00002900000000020c04330000000000210435000000050320021000000000033100190000002005300039000000000002004b00002fb20000613d0000000003000019000000000401001900002f9f0000013d000000000756001900000000000704350000001f066000390000000006e6016f00000000055600190000000103300039000000000023004b00002fb20000813d0000000006150049000000200660008a00000020044000390000000000640435000000200cc0003900000000060c043300000000760604340000000005650436000000000006004b00002f970000613d00000000080000190000000009580019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b00002faa0000413d00002f970000013d0000000001b50049000000040110008a0000004402b000390000000000120435000000140600002900000000020604330000000001250436000000000002004b00002fc50000613d000000000300001900000016080000290000002006600039000000000406043300000c6d0440019700000000014104360000000103300039000000000023004b00002fbd0000413d00002fc60000013d000000160800002900000000020004140000000003000410000000040030008c00002fe00000613d0000000001b1004900000c6d0010009c00000c6d01008041000000600110021000000c6d00b0009c00000c6d0300004100000000030b40190000004003300210000000000131019f00000c6d0020009c00000c6d02008041000000c002200210000000000121019f000000000200041031b131a70000040f000000150b0000290000001608000029000000600310027000010c6d0030019d000300000001035500000001002001900000306b0000613d00000c7100b0009c000030a20000213d0000004000b0043f00000cc700b0009c000030a20000213d000000400080043f00000000000b043500000001050000390000000204000039001300000005001d0000001c0100002900000000010104330000006001100039000000000101043300000c7102100197001400000002001d00000001032002100000007f0110019000002ff70000613d000000ff0230018f00000000011200d9000000020010008c000030a80000c13d000d00000003001d000e00000004001d00150000000b001d001600000008001d0000001d01000029000000000010043f0000000901000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b00000014020000290000000702200270001400000002001d000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b000000000101041a000c00000001001d0000001d01000029000000000010043f0000000901000039000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d000000000101043b0000001402000029000000000020043f000000200010043f000000000100041400000c6d0010009c00000c6d01008041000000c00110021000000c7f011001c7000080100200003931b131ac0000040f00000001002001900000309a0000613d0000000d02000029000000fe0220018f0000000e0900002900000000032901cf000000030220020f00000d11022001670000000c0220017f000000000232019f000000000101043b000000000021041b000000190000006b000000130200002900000001022061bf0000001c01000029000000000101043300000000070104330000000100200190000030cc0000613d0000001a0200002900000000020204330000001b0020006c000000160a00002900000015080000290000309c0000a13d0000006001100039000000000101043300000011030000290000000f0230002900000000030204330000000002000414000000120220006b000030a80000413d00000c7106100197000000400100043d000000400410003900000080050000390000000000540435000000200410003900000000009404350000000000310435000000000308043300000080041000390000000000340435000000a004100039000000000003004b00002d1e0000613d0000000005000019000000000845001900000000095a0019000000000909043300000000009804350000002005500039000000000035004b000030630000413d00002d1e0000013d00000c6d02300198000000080500002900000080080000390000000304000039000000600b00003900002fe90000613d0000001f0320003900000c6e033001970000003f0330003900000d2103300197000000400b00043d00000000033b00190000000000b3004b0000000004000039000000010400403900000c710030009c000030a20000213d0000000100400190000030a20000c13d000000400030043f00000000082b043600000c6f042001980000000003480019000030890000613d000000000501034f0000000006080019000000005705043c0000000006760436000000000036004b000030850000c13d0000001f02200190000030960000613d000000000141034f0000000302200210000000000403043300000000042401cf000000000424022f000000000101043b0000010002200089000000000121022f00000000012101cf000000000141019f00000000001304350000000805000029000000030400003900002fe90000013d000000000001042d0000000001000019000031b30001043000000d0a01000041000000000010043f0000003201000039000000040010043f00000c8501000041000031b30001043000000d0a01000041000000000010043f0000004101000039000000040010043f00000c8501000041000031b30001043000000d0a01000041000000000010043f0000001101000039000000040010043f00000c8501000041000031b300010430000000000001042f00000d190200004100000000002104350000000402100039000030f60000013d0000002402100039000000000032043500000d1a02000041000030c30000013d000000400100043d00000d2502000041000000000021043500000c6d0010009c00000c6d01008041000000400110021000000c90011001c7000031b300010430000000400100043d0000002402100039000000000062043500000d1f02000041000000000021043500000004021000390000001d03000029000000000032043500000c6d0010009c00000c6d01008041000000400110021000000cd4011001c7000031b300010430000000400100043d00000024021000390000004003000039000000000032043500000d220200004100000000002104350000000402100039000000000072043500000015020000290000000002020433000000440310003900000000002304350000006403100039000000000002004b0000001607000029000030e40000613d000000000400001900000000053400190000000006470019000000000606043300000000006504350000002004400039000000000024004b000030dd0000413d0000001f0420003900000d0e0440019700000000023200190000000000020435000000640240003900000c6d0020009c00000c6d02008041000000600220021000000c6d0010009c00000c6d010080410000004001100210000000000112019f000031b300010430000000400100043d00000d1d02000041000000000021043500000004021000390000001d03000029000000000032043500000c6d0010009c00000c6d01008041000000400110021000000c85011001c7000031b3000104300000001f0530018f00000c6f06300198000000400200043d00000000046200190000311a0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000031030000c13d0000311a0000013d00000cf8010000410000001a03000029000000000013043500000004013000390000001d020000290000000000210435000031350000013d0000001f0530018f00000c6f06300198000000400200043d00000000046200190000311a0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000031160000c13d000000000005004b000031270000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f0000000000140435000000600130021000000c6d0020009c00000c6d020080410000004002200210000000000112019f000031b3000104300000001101000029000000000101043300000d27020000410000001a03000029000000000023043500000c71011001970000000402300039000000000012043500000c6d0030009c00000c6d03008041000000400130021000000c85011001c7000031b30001043000000d17010000410000001a02000029000000000012043500000c6d0020009c00000c6d02008041000000400120021000000c90011001c7000031b30001043000000cfe01000041000000000016043500000004016000390000001d02000029000000000021043500000c6d0060009c00000c6d06008041000000400160021000000c85011001c7000031b30001043000000d0a01000041000000000010043f0000002201000039000000040010043f00000c8501000041000031b300010430000000400100043d00000d2602000041000030b90000013d000000400100043d00000d2402000041000030f30000013d00000d28010000410000313c0000013d00000020041000390000004005000039000000000054043500000c720220019700000000002104350000004004100039000000003203043400000000002404350000006001100039000000000002004b0000316d0000613d000000000400001900000000051400190000000006430019000000000606043300000000006504350000002004400039000000000024004b000031660000413d000000000312001900000000000304350000001f0220003900000d0e022001970000000001120019000000000001042d000000000001042f00000c6d0010009c00000c6d01008041000000400110021000000c6d0020009c00000c6d020080410000006002200210000000000112019f000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000c83011001c7000080100200003931b131ac0000040f0000000100200190000031870000613d000000000101043b000000000001042d0000000001000019000031b30001043000000000050100190000000000200443000000050030008c000031970000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b0000318f0000413d00000c6d0030009c00000c6d030080410000006001300210000000000200041400000c6d0020009c00000c6d02008041000000c002200210000000000112019f00000d29011001c7000000000205001931b131ac0000040f0000000100200190000031a60000613d000000000101043b000000000001042d000000000001042f000031aa002104210000000102000039000000000001042d0000000002000019000000000001042d000031af002104230000000102000039000000000001042d0000000002000019000000000001042d000031b100000432000031b20001042e000031b30001043000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000000000000000000000000000ffffffffffffff7f000000000000000000000000000000000000000000000000ffffffffffffffff000000000000000000000000ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffff00000000000000000000000000000000000000009a8a0592ac89c5ad3bc6df8224c17b485976f597df104ee20d0df415241f670b02000002000000000000000000000000000000040000000000000000000000000200000000000000000000000000000000000080000000000000000000000000683eb52ee924eb817377cfa8f41f238f4bb7a877da5267869dfffbad85f564d8ffffffffffffff000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000cbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee0200000000000000000000000000000000000040000000000000000000000000000000ffffffffffffffff000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000002000000000000000000000000000000000000200000000000000000000000000200000000000000000000000000000000000000000000000000000000000000420b006e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000ffffff0000000000000000fffffffffffffffffffffffffffffffffffffffffff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9000000000000000000000000000000000000000000000000ffffffffffffffbf0000000000000000000000010000000000000000000000000000000000000000ffffffffffffffffffffff00000000000000000000000000000000000000000009addddcec1d7ba6ad726df49aeea3e93fb0c1037d551236841a60c0c883f2c10000000000000000000000ff000000000000000000000000000000000000000049f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b00000002000000000000000000000000000001800000010000000000000000009b15e16f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000060987c1f00000000000000000000000000000000000000000000000000000000c673e58300000000000000000000000000000000000000000000000000000000e9d68a8d00000000000000000000000000000000000000000000000000000000e9d68a8e00000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000f716f99f00000000000000000000000000000000000000000000000000000000c673e58400000000000000000000000000000000000000000000000000000000ccd37ba30000000000000000000000000000000000000000000000000000000079ba50960000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000007edf52f4000000000000000000000000000000000000000000000000000000008da5cb5b0000000000000000000000000000000000000000000000000000000060987c20000000000000000000000000000000000000000000000000000000007437ff9f00000000000000000000000000000000000000000000000000000000311cd512000000000000000000000000000000000000000000000000000000005215505a000000000000000000000000000000000000000000000000000000005215505b000000000000000000000000000000000000000000000000000000005e36480c000000000000000000000000000000000000000000000000000000005e7bb00800000000000000000000000000000000000000000000000000000000311cd513000000000000000000000000000000000000000000000000000000003f4b04aa00000000000000000000000000000000000000000000000000000000181f5a7600000000000000000000000000000000000000000000000000000000181f5a77000000000000000000000000000000000000000000000000000000002d04ab760000000000000000000000000000000000000000000000000000000004666f9c0000000000000000000000000000000000000000000000000000000006285c692b5c74de000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000800000000000000000000000000000000000000000000000000000000000000000ffffffffffffff9f000000000000000000000000000000000000000000000000ffffffffffffff1f0000000000000000000000000000000000000000000000000000000000ff0000000000000000000000000000000000000000000000000000ffffffffffffff3f87f6037c00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff0000000000000000000000000000000000000000000000000000000000010000367f56a200000000000000000000000000000000000000000000000000000000d6c62c9b00000000000000000000000000000000000000000000000000000000ab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547f718e9a000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffff00000000000000002f7b1ba2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000001000000000000000000ffffffffffffff0000000000ffffffffffffffffffffffffffffffffffffffff02000000000000000000000000000000000000800000010000000000000000008579befe0000000000000000000000000000000000000000000000000000000002b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e000000000000000ff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffe1f000000000000000000000000000000000000000000000000ffffffffffffff5f371a732800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffdfbbe4f6db00000000000000000000000000000000000000000000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e020000020000000000000000000000000000004400000000000000000000000001ffc9a7000000000000000000000000000000000000000000000000000000000000000000007530000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000000aff2afbf0000000000000000000000000000000000000000000000000000000070a0823100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffedf1806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b83020000020000000000000000000000000000002400000000000000000000000078ef8024000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffeff390775370000000000000000000000000000000000000000000000000000000037c3be2900000000000000000000000000000000000000000000000000000000a966e21f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000640000000000000000000000009fe2f95a00000000000000000000000000000000000000000000000000000000afa32a2c000000000000000000000000000000000000000000000000000000000c3b563c00000000000000000000000000000000000000000000000000000000ae9b4ce90000000000000000000000000000000000000000000000000000000008d450a10000000000000000000000000000000000000000000000000000000009c253250000000000000000000000000000000000000000000000000000000085572ffb000000000000000000000000000000000000000000000000000000003cf97983000000000000000000000000000000000000000000000000000000000a8d6e8c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffe9f0f01ce850000000000000000000000000000000000000000000000000000000083e3f564000000000000000000000000000000000000000000000000000000007531a8c60000000000000000000000000000000000000000000000000000000070a193fd0000000000000000000000000000000000000000000000000000000048e617b30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000840000000000000000000000000000000000000000000000000000000000000020000000000000000000000000e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e18e1192e10000000000000000000000000000000000000000000000000000000093df584c00000000000000000000000000000000000000000000000000000000a15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054ce90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e30000000000000000000000000000000000000080000000000000000000000000198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff70a9089e0000000000000000000000000000000000000000000000000000000022611167000000000000000000000000000000000000000000000000000000003937306f000000000000000000000000000000000000000000000000000000002cbc26bb00000000000000000000000000000000000000000000000000000000fdbd6a7200000000000000000000000000000000000000000000000000000000b80d8fa90000000000000000000000000000000000000000000000000000000032cf0cbf00000000000000000000000000000000000000000000000000000000796b89b91644bc98cd93958e4c9038275d622183e25ac5af08cc6b5d95539132504570e300000000000000000000000000000000000000000000000000000000d5e0f0d600000000000000000000000000000000000000000000000000000000ed053c590000000000000000000000000000000000000000000000000000000035c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4ac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077bac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077c3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92effac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077eda0f08e80000000000000000000000000000000000000000000000000000000071253a2500000000000000000000000000000000000000000000000000000000a75d88af00000000000000000000000000000000000000000000000000000000e4ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffca31867a00000000000000000000000000000000000000000000000000000000f67bc7c4000000000000000000000000000000000000000000000000000000004e487b71000000000000000000000000000000000000000000000000000000004f666652616d7020312e362e302d6465760000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000c656089500000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5f000000000000000000000000000000000000000000000000ffffffffffffff800000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000ffffffffffffffe0aab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d89493357e0e083000000000000000000000000000000000000000000000000000000002425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f38432a22000000000000000000000000000000000000000000000000000000006c95f1eb000000000000000000000000000000000000000000000000000000003b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120a9cfc86200000000000000000000000000000000000000000000000000000000e0e03cae000000000000000000000000000000000000000000000000000000001cfe6d8b0000000000000000000000000000000000000000000000000000000060987c200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003ffffffe02b11b8d90000000000000000000000000000000000000000000000000000000005665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b7dd17a7e0000000000000000000000000000000000000000000000000000000009bde3390000000000000000000000000000000000000000000000000000000011a6b26400000000000000000000000000000000000000000000000000000000ced9e49600000000000000000000000000000000000000000000000000000000c2e5347d0000000000000000000000000000000000000000000000000000000002000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

func DeployOffRampZK(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OffRampStaticConfig, dynamicConfig OffRampDynamicConfig,

	sourceChainConfigs []OffRampSourceChainConfigArgs) (common.
	Address, *generated.Transaction, *OffRamp, error) {
	parsed, err := OffRampMetaData.GetAbi()
	if err != nil {
		return common.Address{},
			nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	address, ethTx, contract, err := generated.
		DeployContract(auth, parsed,
			common.FromHex(OffRampZKBin), backend,
			staticConfig, dynamicConfig,
			sourceChainConfigs,
		)
	if err != nil {
		return common.Address{}, nil,
			nil, err
	}
	return address,

		ethTx, &OffRamp{address: address,
			abi: *parsed, OffRampCaller: OffRampCaller{contract: contract}, OffRampTransactor: OffRampTransactor{contract: contract}, OffRampFilterer: OffRampFilterer{
				contract: contract}}, nil
}
