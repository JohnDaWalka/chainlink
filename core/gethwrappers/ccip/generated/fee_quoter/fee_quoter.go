package fee_quoter

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

type AuthorizedCallersAuthorizedCallerArgs struct {
	AddedCallers   []common.Address
	RemovedCallers []common.Address
}

type ClientEVM2AnyMessage struct {
	Receiver     []byte
	Data         []byte
	TokenAmounts []ClientEVMTokenAmount
	FeeToken     common.Address
	ExtraArgs    []byte
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool
	MaxNumberOfTokensPerMsg           uint16
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
	DestGasOverhead                   uint32
	DestGasPerPayloadByte             uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	DefaultTokenFeeUSDCents           uint16
	DefaultTokenDestGasOverhead       uint32
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	NetworkFeeUSDCents                uint32
	GasPriceStalenessThreshold        uint32
	EnforceOutOfOrder                 bool
	ChainFamilySelector               [4]byte
}

type FeeQuoterDestChainConfigArgs struct {
	DestChainSelector uint64
	DestChainConfig   FeeQuoterDestChainConfig
}

type FeeQuoterPremiumMultiplierWeiPerEthArgs struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg            *big.Int
	LinkToken                    common.Address
	TokenPriceStalenessThreshold uint32
}

type FeeQuoterTokenPriceFeedConfig struct {
	DataFeedAddress common.Address
	TokenDecimals   uint8
	IsEnabled       bool
}

type FeeQuoterTokenPriceFeedUpdate struct {
	SourceToken common.Address
	FeedConfig  FeeQuoterTokenPriceFeedConfig
}

type FeeQuoterTokenTransferFeeConfig struct {
	MinFeeUSDCents    uint32
	MaxFeeUSDCents    uint32
	DeciBps           uint16
	DestGasOverhead   uint32
	DestBytesOverhead uint32
	IsEnabled         bool
}

type FeeQuoterTokenTransferFeeConfigArgs struct {
	DestChainSelector       uint64
	TokenTransferFeeConfigs []FeeQuoterTokenTransferFeeConfigSingleTokenArgs
}

type FeeQuoterTokenTransferFeeConfigRemoveArgs struct {
	DestChainSelector uint64
	Token             common.Address
}

type FeeQuoterTokenTransferFeeConfigSingleTokenArgs struct {
	Token                  common.Address
	TokenTransferFeeConfig FeeQuoterTokenTransferFeeConfig
}

type InternalEVM2AnyTokenTransfer struct {
	SourcePoolAddress common.Address
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalTimestampedPackedUint224 struct {
	Value     *big.Int
	Timestamp uint32
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type KeystoneFeedsPermissionHandlerPermission struct {
	Forwarder     common.Address
	WorkflowName  [10]byte
	ReportName    [2]byte
	WorkflowOwner common.Address
	IsAllowed     bool
}

var FeeQuoterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"tokenPriceStalenessThreshold\",\"type\":\"uint32\"}],\"internalType\":\"structFeeQuoter.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"priceUpdaters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokens\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeeds\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigSingleTokenArgs[]\",\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structFeeQuoter.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DataFeedValueOutOfUint224Range\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"DestinationChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExtraArgOutOfOrderExecutionMustBeTrue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"FeeTokenNotSupported\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"}],\"name\":\"InvalidDestBytesOverhead\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidDestChainConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minFeeUSDCents\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint256\"}],\"name\":\"InvalidFeeRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStaticConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgFeeJuels\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint256\"}],\"name\":\"MessageFeeTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MessageGasLimitTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"},{\"internalType\":\"bytes2\",\"name\":\"reportName\",\"type\":\"bytes2\"}],\"name\":\"ReportForwarderUnauthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SourceTokenDataTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timePassed\",\"type\":\"uint256\"}],\"name\":\"StaleGasPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenNotSupported\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnauthorizedCaller\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numberOfTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint256\"}],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"name\":\"DestChainConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"}],\"name\":\"FeeTokenRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"name\":\"PremiumMultiplierWeiPerEthUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"priceFeedConfig\",\"type\":\"tuple\"}],\"name\":\"PriceFeedPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"reportId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"},{\"internalType\":\"bytes2\",\"name\":\"reportName\",\"type\":\"bytes2\"},{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isAllowed\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structKeystoneFeedsPermissionHandler.Permission\",\"name\":\"permission\",\"type\":\"tuple\"}],\"name\":\"ReportPermissionSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"TokenTransferFeeConfigDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"name\":\"TokenTransferFeeConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerTokenUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChain\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"UsdPerUnitGasUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FEE_BASE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"KEYSTONE_PRICE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"addedCallers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"removedCallers\",\"type\":\"address[]\"}],\"internalType\":\"structAuthorizedCallers.AuthorizedCallerArgs\",\"name\":\"authorizedCallerArgs\",\"type\":\"tuple\"}],\"name\":\"applyAuthorizedCallerUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"destChainConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyDestChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"feeTokensToRemove\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"feeTokensToAdd\",\"type\":\"address[]\"}],\"name\":\"applyFeeTokensUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"internalType\":\"structFeeQuoter.PremiumMultiplierWeiPerEthArgs[]\",\"name\":\"premiumMultiplierWeiPerEthArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyPremiumMultiplierWeiPerEthUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigSingleTokenArgs[]\",\"name\":\"tokenTransferFeeConfigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigArgs[]\",\"name\":\"tokenTransferFeeConfigArgs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfigRemoveArgs[]\",\"name\":\"tokensToUseDefaultFeeConfigs\",\"type\":\"tuple[]\"}],\"name\":\"applyTokenTransferFeeConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fromTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"toToken\",\"type\":\"address\"}],\"name\":\"convertTokenAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedCallers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerMsgGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerPayloadByte\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destDataAvailabilityOverheadGas\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"destGasPerDataAvailabilityByte\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"destDataAvailabilityMultiplierBps\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"defaultTokenFeeUSDCents\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"defaultTokenDestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"defaultTxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasMultiplierWeiPerEth\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"networkFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasPriceStalenessThreshold\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enforceOutOfOrder\",\"type\":\"bool\"},{\"internalType\":\"bytes4\",\"name\":\"chainFamilySelector\",\"type\":\"bytes4\"}],\"internalType\":\"structFeeQuoter.DestChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestinationChainGasPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getPremiumMultiplierWeiPerEth\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"premiumMultiplierWeiPerEth\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96\",\"name\":\"maxFeeJuelsPerMsg\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"tokenPriceStalenessThreshold\",\"type\":\"uint32\"}],\"internalType\":\"structFeeQuoter.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getTokenAndGasPrices\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"tokenPrice\",\"type\":\"uint224\"},{\"internalType\":\"uint224\",\"name\":\"gasPriceValue\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenPriceFeedConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"getTokenPrices\",\"outputs\":[{\"components\":[{\"internalType\":\"uint224\",\"name\":\"value\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structInternal.TimestampedPackedUint224[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenTransferFeeConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"minFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxFeeUSDCents\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"deciBps\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"destGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destBytesOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenTransferFeeConfig\",\"name\":\"tokenTransferFeeConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getValidatedFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getValidatedTokenPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"onRampTokenTransfers\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"sourceTokenAmounts\",\"type\":\"tuple[]\"}],\"name\":\"processMessageArgs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"msgFeeJuels\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOutOfOrderExecution\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"convertedExtraArgs\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"destExecDataPerToken\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"},{\"internalType\":\"bytes2\",\"name\":\"reportName\",\"type\":\"bytes2\"},{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isAllowed\",\"type\":\"bool\"}],\"internalType\":\"structKeystoneFeedsPermissionHandler.Permission[]\",\"name\":\"permissions\",\"type\":\"tuple[]\"}],\"name\":\"setReportPermissions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"updatePrices\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"dataFeedAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedConfig\",\"name\":\"feedConfig\",\"type\":\"tuple\"}],\"internalType\":\"structFeeQuoter.TokenPriceFeedUpdate[]\",\"name\":\"tokenPriceFeedUpdates\",\"type\":\"tuple[]\"}],\"name\":\"updateTokenPriceFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162007a6838038062007a6883398101604081905262000034916200189c565b85336000816200005757604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008a576200008a81620001d0565b5050604080518082018252828152815160008152602080820190935291810191909152620000b8906200024a565b5060208701516001600160a01b03161580620000dc575086516001600160601b0316155b80620000f05750604087015163ffffffff16155b156200010f5760405163d794ef9560e01b815260040160405180910390fd5b6020878101516001600160a01b031660a05287516001600160601b031660805260408089015163ffffffff1660c052805160008152918201905262000155908662000399565b6200016084620004e1565b6200016b81620005d9565b620001768262000a45565b60408051600080825260208201909252620001c391859190620001bc565b6040805180820190915260008082526020820152815260200190600190039081620001945790505b5062000b11565b5050505050505062001b5a565b336001600160a01b03821603620001fa57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b602081015160005b8151811015620002da576000828281518110620002735762000273620019bb565b602090810291909101015190506200028d60028262000e97565b15620002d0576040516001600160a01b03821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b5060010162000252565b50815160005b815181101562000393576000828281518110620003015762000301620019bb565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200033f576040516342bcdf7f60e11b815260040160405180910390fd5b6200034c60028262000eb7565b506040516001600160a01b03821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a150600101620002e0565b50505050565b60005b82518110156200043a57620003d8838281518110620003bf57620003bf620019bb565b6020026020010151600b62000ece60201b90919060201c565b156200043157828181518110620003f357620003f3620019bb565b60200260200101516001600160a01b03167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f9160405160405180910390a25b6001016200039c565b5060005b8151811015620004dc576200047a828281518110620004615762000461620019bb565b6020026020010151600b62000eb760201b90919060201c565b15620004d357818181518110620004955762000495620019bb565b60200260200101516001600160a01b03167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba2360405160405180910390a25b6001016200043e565b505050565b60005b8151811015620005d5576000828281518110620005055762000505620019bb565b6020908102919091018101518051818301516001600160a01b0380831660008181526007875260409081902084518154868a0180518589018051949098166001600160a81b03199093168317600160a01b60ff928316021760ff60a81b1916600160a81b9415159490940293909317909355835190815291511697810197909752915115159186019190915292945090929091907fe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf9060600160405180910390a2505050806001019050620004e4565b5050565b60005b8151811015620005d5576000828281518110620005fd57620005fd620019bb565b6020026020010151905060008383815181106200061e576200061e620019bb565b6020026020010151600001519050600082602001519050816001600160401b03166000148062000657575061016081015163ffffffff16155b806200067957506102008101516001600160e01b031916630a04b54b60e21b14155b80620006995750806060015163ffffffff1681610160015163ffffffff16115b15620006c85760405163c35aa79d60e01b81526001600160401b03831660048201526024015b60405180910390fd5b6001600160401b038216600090815260096020526040812060010154600160a81b900460e01b6001600160e01b03191690036200074857816001600160401b03167f525e3d4e0c31cef19cf9426af8d2c0ddd2d576359ca26bed92aac5fadda46265826040516200073a9190620019d1565b60405180910390a26200078c565b816001600160401b03167f283b699f411baff8f1c29fe49f32a828c8151596244b8e7e4c164edd6569a83582604051620007839190620019d1565b60405180910390a25b8060096000846001600160401b03166001600160401b0316815260200190815260200160002060008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548161ffff021916908361ffff16021790555060408201518160000160036101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160076101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600b6101000a81548163ffffffff021916908363ffffffff16021790555060a082015181600001600f6101000a81548161ffff021916908361ffff16021790555060c08201518160000160116101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160000160156101000a81548161ffff021916908361ffff1602179055506101008201518160000160176101000a81548161ffff021916908361ffff1602179055506101208201518160000160196101000a81548161ffff021916908361ffff16021790555061014082015181600001601b6101000a81548163ffffffff021916908363ffffffff1602179055506101608201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101808201518160010160046101000a8154816001600160401b0302191690836001600160401b031602179055506101a082015181600101600c6101000a81548163ffffffff021916908363ffffffff1602179055506101c08201518160010160106101000a81548163ffffffff021916908363ffffffff1602179055506101e08201518160010160146101000a81548160ff0219169083151502179055506102008201518160010160156101000a81548163ffffffff021916908360e01c0217905550905050505050806001019050620005dc565b60005b8151811015620005d557600082828151811062000a695762000a69620019bb565b6020026020010151600001519050600083838151811062000a8e5762000a8e620019bb565b6020908102919091018101518101516001600160a01b03841660008181526008845260409081902080546001600160401b0319166001600160401b0385169081179091559051908152919350917fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d910160405180910390a2505060010162000a48565b60005b825181101562000dd157600083828151811062000b355762000b35620019bb565b6020026020010151905060008160000151905060005b82602001515181101562000dc25760008360200151828151811062000b745762000b74620019bb565b602002602001015160200151905060008460200151838151811062000b9d5762000b9d620019bb565b6020026020010151600001519050816020015163ffffffff16826000015163ffffffff161062000bf857815160208301516040516305a7b3d160e11b815263ffffffff928316600482015291166024820152604401620006bf565b602063ffffffff16826080015163ffffffff16101562000c495760808201516040516312766e0160e11b81526001600160a01b038316600482015263ffffffff9091166024820152604401620006bf565b6001600160401b0384166000818152600a602090815260408083206001600160a01b0386168085529083529281902086518154938801518389015160608a015160808b015160a08c01511515600160901b0260ff60901b1963ffffffff928316600160701b021664ffffffffff60701b199383166a01000000000000000000000263ffffffff60501b1961ffff90961668010000000000000000029590951665ffffffffffff60401b19968416640100000000026001600160401b0319909b16939097169290921798909817939093169390931717919091161792909217909155519091907f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b59062000daf908690600060c08201905063ffffffff80845116835280602085015116602084015261ffff60408501511660408401528060608501511660608401528060808501511660808401525060a0830151151560a083015292915050565b60405180910390a3505060010162000b4b565b50505080600101905062000b14565b5060005b8151811015620004dc57600082828151811062000df65762000df6620019bb565b6020026020010151600001519050600083838151811062000e1b5762000e1b620019bb565b6020908102919091018101518101516001600160401b0384166000818152600a845260408082206001600160a01b038516808452955280822080546001600160981b03191690555192945090917f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b9190a3505060010162000dd5565b600062000eae836001600160a01b03841662000ee5565b90505b92915050565b600062000eae836001600160a01b03841662000fe9565b600062000eae836001600160a01b0384166200103b565b6000818152600183016020526040812054801562000fde57600062000f0c60018362001b22565b855490915060009062000f229060019062001b22565b905081811462000f8e57600086600001828154811062000f465762000f46620019bb565b906000526020600020015490508087600001848154811062000f6c5762000f6c620019bb565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062000fa25762000fa262001b44565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062000eb1565b600091505062000eb1565b6000818152600183016020526040812054620010325750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000eb1565b50600062000eb1565b6000818152600183016020526040812054801562000fde5760006200106260018362001b22565b8554909150600090620010789060019062001b22565b905080821462000f8e57600086600001828154811062000f465762000f46620019bb565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b0381118282101715620010d757620010d76200109c565b60405290565b604080519081016001600160401b0381118282101715620010d757620010d76200109c565b60405160c081016001600160401b0381118282101715620010d757620010d76200109c565b60405161022081016001600160401b0381118282101715620010d757620010d76200109c565b604051601f8201601f191681016001600160401b03811182821017156200117857620011786200109c565b604052919050565b80516001600160a01b03811681146200119857600080fd5b919050565b805163ffffffff811681146200119857600080fd5b600060608284031215620011c557600080fd5b620011cf620010b2565b82519091506001600160601b0381168114620011ea57600080fd5b8152620011fa6020830162001180565b60208201526200120d604083016200119d565b604082015292915050565b60006001600160401b038211156200123457620012346200109c565b5060051b60200190565b600082601f8301126200125057600080fd5b8151602062001269620012638362001218565b6200114d565b8083825260208201915060208460051b8701019350868411156200128c57600080fd5b602086015b84811015620012b357620012a58162001180565b835291830191830162001291565b509695505050505050565b805180151581146200119857600080fd5b600082601f830112620012e157600080fd5b81516020620012f4620012638362001218565b82815260079290921b840181019181810190868411156200131457600080fd5b8286015b84811015620012b3578088036080811215620013345760008081fd5b6200133e620010dd565b620013498362001180565b8152606080601f1984011215620013605760008081fd5b6200136a620010b2565b92506200137987850162001180565b835260408085015160ff81168114620013925760008081fd5b84890152620013a3858301620012be565b90840152508086019190915283529183019160800162001318565b80516001600160401b03811681146200119857600080fd5b805161ffff811681146200119857600080fd5b600082601f830112620013fb57600080fd5b815160206200140e620012638362001218565b82815260059290921b840181019181810190868411156200142e57600080fd5b8286015b84811015620012b35780516001600160401b03808211156200145357600080fd5b908801906040601f19838c0381018213156200146e57600080fd5b62001478620010dd565b62001485898601620013be565b815282850151848111156200149957600080fd5b8086019550508c603f860112620014af57600080fd5b888501519350620014c4620012638562001218565b84815260e09094028501830193898101908e861115620014e357600080fd5b958401955b85871015620015bc57868f0360e08112156200150357600080fd5b6200150d620010dd565b620015188962001180565b815260c086830112156200152b57600080fd5b6200153562001102565b9150620015448d8a016200119d565b825262001553878a016200119d565b8d8301526200156560608a01620013d6565b878301526200157760808a016200119d565b60608301526200158a60a08a016200119d565b60808301526200159d60c08a01620012be565b60a0830152808d0191909152825260e09690960195908a0190620014e8565b828b01525087525050509284019250830162001432565b600082601f830112620015e557600080fd5b81516020620015f8620012638362001218565b82815260069290921b840181019181810190868411156200161857600080fd5b8286015b84811015620012b35760408189031215620016375760008081fd5b62001641620010dd565b6200164c8262001180565b81526200165b858301620013be565b818601528352918301916040016200161c565b80516001600160e01b0319811681146200119857600080fd5b600082601f8301126200169957600080fd5b81516020620016ac620012638362001218565b8281526102409283028501820192828201919087851115620016cd57600080fd5b8387015b858110156200188f5780890382811215620016ec5760008081fd5b620016f6620010dd565b6200170183620013be565b815261022080601f1984011215620017195760008081fd5b6200172362001127565b925062001732888501620012be565b8352604062001743818601620013d6565b898501526060620017568187016200119d565b82860152608091506200176b8287016200119d565b9085015260a06200177e8682016200119d565b8286015260c0915062001793828701620013d6565b9085015260e0620017a68682016200119d565b828601526101009150620017bc828701620013d6565b90850152610120620017d0868201620013d6565b828601526101409150620017e6828701620013d6565b90850152610160620017fa8682016200119d565b828601526101809150620018108287016200119d565b908501526101a062001824868201620013be565b828601526101c091506200183a8287016200119d565b908501526101e06200184e8682016200119d565b82860152610200915062001864828701620012be565b90850152620018758583016200166e565b9084015250808701919091528452928401928101620016d1565b5090979650505050505050565b6000806000806000806000610120888a031215620018b957600080fd5b620018c58989620011b2565b60608901519097506001600160401b0380821115620018e357600080fd5b620018f18b838c016200123e565b975060808a01519150808211156200190857600080fd5b620019168b838c016200123e565b965060a08a01519150808211156200192d57600080fd5b6200193b8b838c01620012cf565b955060c08a01519150808211156200195257600080fd5b620019608b838c01620013e9565b945060e08a01519150808211156200197757600080fd5b620019858b838c01620015d3565b93506101008a01519150808211156200199d57600080fd5b50620019ac8a828b0162001687565b91505092959891949750929550565b634e487b7160e01b600052603260045260246000fd5b81511515815261022081016020830151620019f2602084018261ffff169052565b50604083015162001a0b604084018263ffffffff169052565b50606083015162001a24606084018263ffffffff169052565b50608083015162001a3d608084018263ffffffff169052565b5060a083015162001a5460a084018261ffff169052565b5060c083015162001a6d60c084018263ffffffff169052565b5060e083015162001a8460e084018261ffff169052565b506101008381015161ffff9081169184019190915261012080850151909116908301526101408084015163ffffffff9081169184019190915261016080850151821690840152610180808501516001600160401b0316908401526101a0808501518216908401526101c080850151909116908301526101e080840151151590830152610200928301516001600160e01b031916929091019190915290565b8181038181111562000eb157634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c051615ebb62001bad6000396000818161032801526119170152600081816102ec0152818161103401526110940152600081816102b8015281816110bd015261112d0152615ebb6000f3fe608060405234801561001057600080fd5b50600436106101e45760003560e01c8063770e2dc41161010f578063bf78e03f116100a2578063d8694ccd11610071578063d8694ccd14610b04578063f2fde38b14610b17578063fbe3f77814610b2a578063ffdb4b3714610b3d57600080fd5b8063bf78e03f14610a03578063cdc73d5114610ae1578063d02641a014610ae9578063d63d3af214610afc57600080fd5b806382b49eb0116100de57806382b49eb0146108455780638da5cb5b146109b557806391a2749a146109dd578063a69c64c0146109f057600080fd5b8063770e2dc41461080457806379ba5097146108175780637afac3221461081f578063805f21321461083257600080fd5b80633937306f116101875780634ab35b0b116101565780634ab35b0b14610472578063514e8cff146104b25780636cb5f3dd146105555780636def4ce71461056857600080fd5b80633937306f1461040757806341ed29e71461041c578063430d138c1461042f57806345ac924d1461045257600080fd5b806306285c69116101c357806306285c691461028b578063181f5a77146103a15780632451a627146103ea578063325c868e146103ff57600080fd5b806241e5be146101e957806301ffc9a71461020f578063061877e314610232575b600080fd5b6101fc6101f7366004614568565b610b85565b6040519081526020015b60405180910390f35b61022261021d3660046145d4565b610bf3565b6040519015158152602001610206565b6102726102403660046145ef565b73ffffffffffffffffffffffffffffffffffffffff1660009081526008602052604090205467ffffffffffffffff1690565b60405167ffffffffffffffff9091168152602001610206565b610355604080516060810182526000808252602082018190529181019190915260405180606001604052807f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000063ffffffff16815250905090565b6040805182516bffffffffffffffffffffffff16815260208084015173ffffffffffffffffffffffffffffffffffffffff16908201529181015163ffffffff1690820152606001610206565b6103dd6040518060400160405280601381526020017f46656551756f74657220312e362e302d6465760000000000000000000000000081525081565b604051610206919061466e565b6103f2610d24565b6040516102069190614681565b6101fc602481565b61041a6104153660046146db565b610d35565b005b61041a61042a366004614887565b610fea565b61044261043d366004614a72565b61102c565b6040516102069493929190614b66565b610465610460366004614c05565b61123c565b6040516102069190614c47565b6104856104803660046145ef565b611305565b6040517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9091168152602001610206565b6105486104c0366004614cc2565b60408051808201909152600080825260208201525067ffffffffffffffff166000908152600560209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811683527c0100000000000000000000000000000000000000000000000000000000900463ffffffff169082015290565b6040516102069190614cdd565b61041a610563366004614d3e565b611310565b6107f7610576366004614cc2565b6040805161022081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101829052610160810182905261018081018290526101a081018290526101c081018290526101e081018290526102008101919091525067ffffffffffffffff908116600090815260096020908152604091829020825161022081018452815460ff8082161515835261ffff61010080840482169685019690965263ffffffff630100000084048116978501979097526701000000000000008304871660608501526b0100000000000000000000008304871660808501526f010000000000000000000000000000008304811660a0850152710100000000000000000000000000000000008304871660c08501527501000000000000000000000000000000000000000000808404821660e08087019190915277010000000000000000000000000000000000000000000000850483169786019790975279010000000000000000000000000000000000000000000000000084049091166101208501527b01000000000000000000000000000000000000000000000000000000909204861661014084015260019093015480861661016084015264010000000081049096166101808301526c01000000000000000000000000860485166101a083015270010000000000000000000000000000000086049094166101c082015274010000000000000000000000000000000000000000850490911615156101e08201527fffffffff0000000000000000000000000000000000000000000000000000000092909304901b1661020082015290565b6040516102069190614f5e565b61041a61081236600461515c565b611324565b61041a611336565b61041a61082d366004615476565b611404565b61041a6108403660046154da565b611416565b610955610853366004615546565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a08101919091525067ffffffffffffffff919091166000908152600a6020908152604080832073ffffffffffffffffffffffffffffffffffffffff94909416835292815290829020825160c081018452905463ffffffff8082168352640100000000820481169383019390935268010000000000000000810461ffff16938201939093526a01000000000000000000008304821660608201526e01000000000000000000000000000083049091166080820152720100000000000000000000000000000000000090910460ff16151560a082015290565b6040516102069190600060c08201905063ffffffff80845116835280602085015116602084015261ffff60408501511660408401528060608501511660608401528060808501511660808401525060a0830151151560a083015292915050565b60015460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610206565b61041a6109eb366004615570565b611852565b61041a6109fe366004615601565b611863565b610aa4610a113660046145ef565b60408051606080820183526000808352602080840182905292840181905273ffffffffffffffffffffffffffffffffffffffff9485168152600783528390208351918201845254938416815260ff74010000000000000000000000000000000000000000850481169282019290925275010000000000000000000000000000000000000000009093041615159082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff16815260208084015160ff169082015291810151151590820152606001610206565b6103f2611874565b610548610af73660046145ef565b611880565b6101fc601281565b6101fc610b123660046156c6565b611a35565b61041a610b253660046145ef565b611f6d565b61041a610b3836600461572a565b611f7e565b610b50610b4b36600461584a565b611f8f565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff938416815292909116602083015201610206565b6000610b9082612047565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16610bb785612047565b610bdf907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16856158a3565b610be991906158ba565b90505b9392505050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f805f2132000000000000000000000000000000000000000000000000000000001480610c8657507fffffffff0000000000000000000000000000000000000000000000000000000082167f9b645f4100000000000000000000000000000000000000000000000000000000145b80610cd257507fffffffff0000000000000000000000000000000000000000000000000000000082167f181f5a7700000000000000000000000000000000000000000000000000000000145b80610d1e57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b6060610d3060026120e1565b905090565b610d3d6120ee565b6000610d4982806158f5565b9050905060005b81811015610e93576000610d6484806158f5565b83818110610d7457610d7461595d565b905060400201803603810190610d8a91906159b8565b604080518082018252602080840180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff908116845263ffffffff42818116858701908152885173ffffffffffffffffffffffffffffffffffffffff9081166000908152600690975295889020965190519092167c010000000000000000000000000000000000000000000000000000000002919092161790935584519051935194955016927f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a92610e829290917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b60405180910390a250600101610d50565b506000610ea360208401846158f5565b9050905060005b81811015610fe4576000610ec160208601866158f5565b83818110610ed157610ed161595d565b905060400201803603810190610ee791906159f5565b604080518082018252602080840180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff908116845263ffffffff42818116858701908152885167ffffffffffffffff9081166000908152600590975295889020965190519092167c010000000000000000000000000000000000000000000000000000000002919092161790935584519051935194955016927fdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e92610fd39290917bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b60405180910390a250600101610eaa565b50505050565b610ff2612133565b60005b8151811015611028576110208282815181106110135761101361595d565b6020026020010151612184565b600101610ff5565b5050565b6000806060807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168c73ffffffffffffffffffffffffffffffffffffffff160361108d578a93506110bb565b6110b88c8c7f0000000000000000000000000000000000000000000000000000000000000000610b85565b93505b7f00000000000000000000000000000000000000000000000000000000000000006bffffffffffffffffffffffff1684111561115f576040517f6a92a483000000000000000000000000000000000000000000000000000000008152600481018590526bffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b67ffffffffffffffff8d1660009081526009602052604081206001015463ffffffff169061118e8c8c84612356565b9050806020015194506111a48f8b8b8b8b6124ff565b92508585611224836040805182516024820152602092830151151560448083019190915282518083039091018152606490910190915290810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f181dcf100000000000000000000000000000000000000000000000000000000017905290565b95509550955050509950995099509995505050505050565b60608160008167ffffffffffffffff81111561125a5761125a614716565b60405190808252806020026020018201604052801561129f57816020015b60408051808201909152600080825260208201528152602001906001900390816112785790505b50905060005b828110156112fc576112d78686838181106112c2576112c261595d565b9050602002016020810190610af791906145ef565b8282815181106112e9576112e961595d565b60209081029190910101526001016112a5565b50949350505050565b6000610d1e82612047565b611318612133565b61132181612882565b50565b61132c612133565b6110288282612d54565b60005473ffffffffffffffffffffffffffffffffffffffff163314611387576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61140c612133565b61102882826131ca565b600080600061145a87878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061331192505050565b92509250925061146c3383858461332c565b600061147a85870187615a18565b905060005b8151811015611847576000600760008484815181106114a0576114a061595d565b6020908102919091018101515173ffffffffffffffffffffffffffffffffffffffff908116835282820193909352604091820160002082516060810184529054938416815260ff740100000000000000000000000000000000000000008504811692820192909252750100000000000000000000000000000000000000000090930416151590820181905290915061159b578282815181106115445761154461595d565b6020908102919091010151516040517f06439c6b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401611156565b60006115e8601283602001518686815181106115b9576115b961595d565b6020026020010151602001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16613484565b9050600660008585815181106116005761160061595d565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff168484815181106116725761167261595d565b60200260200101516040015163ffffffff16101561169157505061183f565b6040518060400160405280827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020018585815181106116d2576116d261595d565b60200260200101516040015163ffffffff16815250600660008686815181106116fd576116fd61595d565b6020908102919091018101515173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040016000208251929091015163ffffffff167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff90921691909117905583518490849081106117955761179561595d565b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff167f52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a828686815181106117eb576117eb61595d565b6020026020010151604001516040516118349291907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216825263ffffffff16602082015260400190565b60405180910390a250505b60010161147f565b505050505050505050565b61185a612133565b61132181613547565b61186b612133565b611321816136d3565b6060610d30600b6120e1565b604080518082019091526000808252602082015273ffffffffffffffffffffffffffffffffffffffff82166000908152600660209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116835263ffffffff7c010000000000000000000000000000000000000000000000000000000090910481169183018290527f000000000000000000000000000000000000000000000000000000000000000016906119429042615adf565b101561194e5792915050565b73ffffffffffffffffffffffffffffffffffffffff80841660009081526007602090815260409182902082516060810184529054938416815260ff74010000000000000000000000000000000000000000850481169282019290925275010000000000000000000000000000000000000000009093041615801591830191909152806119ef5750805173ffffffffffffffffffffffffffffffffffffffff16155b156119fb575092915050565b6000611a06826137bd565b9050826020015163ffffffff16816020015163ffffffff161015611a2a5782611a2c565b805b95945050505050565b67ffffffffffffffff8083166000908152600960209081526040808320815161022081018352815460ff808216151580845261ffff61010080850482169886019890985263ffffffff630100000085048116978601979097526701000000000000008404871660608601526b0100000000000000000000008404871660808601526f010000000000000000000000000000008404811660a0860152710100000000000000000000000000000000008404871660c08601527501000000000000000000000000000000000000000000808504821660e08088019190915277010000000000000000000000000000000000000000000000860483169987019990995279010000000000000000000000000000000000000000000000000085049091166101208601527b01000000000000000000000000000000000000000000000000000000909304861661014085015260019094015480861661016085015264010000000081049098166101808401526c01000000000000000000000000880485166101a084015270010000000000000000000000000000000088049094166101c083015274010000000000000000000000000000000000000000870490931615156101e08201527fffffffff000000000000000000000000000000000000000000000000000000009290950490921b16610200840152909190611c6f576040517f99ac52f200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401611156565b611c8a611c8260808501606086016145ef565b600b9061394f565b611ce957611c9e60808401606085016145ef565b6040517f2502348c00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401611156565b6000611cf860408501856158f5565b9150611d54905082611d0d6020870187615af2565b905083611d1a8880615af2565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061397e92505050565b6000611d6e611d6960808701606088016145ef565b612047565b90506000611d8187856101c00151613a3b565b9050600080808515611dc157611db5878b611da260808d0160608e016145ef565b88611db060408f018f6158f5565b613b3b565b91945092509050611de1565b6101a0870151611dde9063ffffffff16662386f26fc100006158a3565b92505b61010087015160009061ffff1615611e2557611e22886dffffffffffffffffffffffffffff607088901c16611e1960208e018e615af2565b90508a86613e13565b90505b61018088015160009067ffffffffffffffff16611e4e611e4860808e018e615af2565b8c613ec3565b600001518563ffffffff168b60a0015161ffff168e8060200190611e729190615af2565b611e7d9291506158a3565b8c6080015163ffffffff16611e929190615b57565b611e9c9190615b57565b611ea69190615b57565b611ec0906dffffffffffffffffffffffffffff89166158a3565b611eca91906158a3565b9050867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168282600860008f6060016020810190611f0491906145ef565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002054611f3f9067ffffffffffffffff16896158a3565b611f499190615b57565b611f539190615b57565b611f5d91906158ba565b9c9b505050505050505050505050565b611f75612133565b61132181613f84565b611f86612133565b61132181614048565b67ffffffffffffffff8116600090815260096020526040812054819060ff16611ff0576040517f99ac52f200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401611156565b611ff984612047565b67ffffffffffffffff841660009081526009602052604090206001015461203b908590700100000000000000000000000000000000900463ffffffff16613a3b565b915091505b9250929050565b60008061205383611880565b9050806020015163ffffffff166000148061208b575080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16155b156120da576040517f06439c6b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401611156565b5192915050565b60606000610bec8361419a565b6120f960023361394f565b612131576040517fd86ad9cf000000000000000000000000000000000000000000000000000000008152336004820152602401611156565b565b60015473ffffffffffffffffffffffffffffffffffffffff163314612131576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061223d82600001518360600151846020015185604001516040805173ffffffffffffffffffffffffffffffffffffffff80871660208301528516918101919091527fffffffffffffffffffff00000000000000000000000000000000000000000000831660608201527fffff0000000000000000000000000000000000000000000000000000000000008216608082015260009060a001604051602081830303815290604052805190602001209050949350505050565b60808301516000828152600460205260409081902080549215157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00909316929092179091555190915081907f32a4ba3fa3351b11ad555d4c8ec70a744e8705607077a946807030d64b6ab1a39061234a908590600060a08201905073ffffffffffffffffffffffffffffffffffffffff8084511683527fffffffffffffffffffff0000000000000000000000000000000000000000000060208501511660208401527fffff00000000000000000000000000000000000000000000000000000000000060408501511660408401528060608501511660608401525060808301511515608083015292915050565b60405180910390a25050565b6040805180820190915260008082526020820152600083900361239757506040805180820190915267ffffffffffffffff8216815260006020820152610bec565b60006123a38486615b6a565b905060006123b48560048189615bb0565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509293505050507fffffffff0000000000000000000000000000000000000000000000000000000082167fe7e230f0000000000000000000000000000000000000000000000000000000000161245157808060200190518101906124489190615bda565b92505050610bec565b7f6859a837000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316016124cd576040518060400160405280828060200190518101906124b99190615c06565b815260006020909101529250610bec915050565b6040517f5247fdce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff808616600090815260096020526040902060010154606091750100000000000000000000000000000000000000000090910460e01b90859081111561254f5761254f614716565b60405190808252806020026020018201604052801561258257816020015b606081526020019060019003908161256d5790505b50915060005b858110156128775760008585838181106125a4576125a461595d565b6125ba92602060409092020190810191506145ef565b905060008888848181106125d0576125d061595d565b90506020028101906125e29190615c1f565b6125f0906040810190615af2565b91505060208111156126a05767ffffffffffffffff8a166000908152600a6020908152604080832073ffffffffffffffffffffffffffffffffffffffff861684529091529020546e010000000000000000000000000000900463ffffffff168111156126a0576040517f36f536ca00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401611156565b612710848a8a868181106126b6576126b661595d565b90506020028101906126c89190615c1f565b6126d6906020810190615af2565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506141f692505050565b67ffffffffffffffff8a166000908152600a6020908152604080832073ffffffffffffffffffffffffffffffffffffffff861684528252808320815160c081018352905463ffffffff8082168352640100000000820481169483019490945268010000000000000000810461ffff16928201929092526a01000000000000000000008204831660608201526e010000000000000000000000000000820490921660808301527201000000000000000000000000000000000000900460ff16151560a082018190529091906128225767ffffffffffffffff8c166000908152600960205260409020547b01000000000000000000000000000000000000000000000000000000900463ffffffff16612828565b81606001515b6040805163ffffffff831660208201529192500160405160208183030381529060405287868151811061285d5761285d61595d565b602002602001018190525050505050806001019050612588565b505095945050505050565b60005b81518110156110285760008282815181106128a2576128a261595d565b6020026020010151905060008383815181106128c0576128c061595d565b60200260200101516000015190506000826020015190508167ffffffffffffffff16600014806128f9575061016081015163ffffffff16155b8061294b57506102008101517fffffffff00000000000000000000000000000000000000000000000000000000167f2812d52c0000000000000000000000000000000000000000000000000000000014155b8061296a5750806060015163ffffffff1681610160015163ffffffff16115b156129ad576040517fc35aa79d00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff83166004820152602401611156565b67ffffffffffffffff82166000908152600960205260408120600101547501000000000000000000000000000000000000000000900460e01b7fffffffff00000000000000000000000000000000000000000000000000000000169003612a55578167ffffffffffffffff167f525e3d4e0c31cef19cf9426af8d2c0ddd2d576359ca26bed92aac5fadda4626582604051612a489190614f5e565b60405180910390a2612a98565b8167ffffffffffffffff167f283b699f411baff8f1c29fe49f32a828c8151596244b8e7e4c164edd6569a83582604051612a8f9190614f5e565b60405180910390a25b80600960008467ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548161ffff021916908361ffff16021790555060408201518160000160036101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160076101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600b6101000a81548163ffffffff021916908363ffffffff16021790555060a082015181600001600f6101000a81548161ffff021916908361ffff16021790555060c08201518160000160116101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160000160156101000a81548161ffff021916908361ffff1602179055506101008201518160000160176101000a81548161ffff021916908361ffff1602179055506101208201518160000160196101000a81548161ffff021916908361ffff16021790555061014082015181600001601b6101000a81548163ffffffff021916908363ffffffff1602179055506101608201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101808201518160010160046101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506101a082015181600101600c6101000a81548163ffffffff021916908363ffffffff1602179055506101c08201518160010160106101000a81548163ffffffff021916908363ffffffff1602179055506101e08201518160010160146101000a81548160ff0219169083151502179055506102008201518160010160156101000a81548163ffffffff021916908360e01c0217905550905050505050806001019050612885565b60005b82518110156130e1576000838281518110612d7457612d7461595d565b6020026020010151905060008160000151905060005b8260200151518110156130d357600083602001518281518110612daf57612daf61595d565b6020026020010151602001519050600084602001518381518110612dd557612dd561595d565b6020026020010151600001519050816020015163ffffffff16826000015163ffffffff1610612e4757815160208301516040517f0b4f67a200000000000000000000000000000000000000000000000000000000815263ffffffff928316600482015291166024820152604401611156565b602063ffffffff16826080015163ffffffff161015612ebc5760808201516040517f24ecdc0200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8316600482015263ffffffff9091166024820152604401611156565b67ffffffffffffffff84166000818152600a6020908152604080832073ffffffffffffffffffffffffffffffffffffffff86168085529083529281902086518154938801518389015160608a015160808b015160a08c015115157201000000000000000000000000000000000000027fffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffff63ffffffff9283166e01000000000000000000000000000002167fffffffffffffffffffffffffff0000000000ffffffffffffffffffffffffffff9383166a0100000000000000000000027fffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffff61ffff9096166801000000000000000002959095167fffffffffffffffffffffffffffffffffffff000000000000ffffffffffffffff968416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909b16939097169290921798909817939093169390931717919091161792909217909155519091907f94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b5906130c1908690600060c08201905063ffffffff80845116835280602085015116602084015261ffff60408501511660408401528060608501511660608401528060808501511660808401525060a0830151151560a083015292915050565b60405180910390a35050600101612d8a565b505050806001019050612d57565b5060005b81518110156131c55760008282815181106131025761310261595d565b602002602001015160000151905060008383815181106131245761312461595d565b60209081029190910181015181015167ffffffffffffffff84166000818152600a8452604080822073ffffffffffffffffffffffffffffffffffffffff8516808452955280822080547fffffffffffffffffffffffffff000000000000000000000000000000000000001690555192945090917f4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b9190a350506001016130e5565b505050565b60005b825181101561326d576132038382815181106131eb576131eb61595d565b6020026020010151600b61424890919063ffffffff16565b156132655782818151811061321a5761321a61595d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167f1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f9160405160405180910390a25b6001016131cd565b5060005b81518110156131c5576132a782828151811061328f5761328f61595d565b6020026020010151600b61426a90919063ffffffff16565b15613309578181815181106132be576132be61595d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167fdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba2360405160405180910390a25b600101613271565b6040810151604a820151605e90920151909260609290921c91565b6040805173ffffffffffffffffffffffffffffffffffffffff868116602080840191909152908616828401527fffffffffffffffffffff00000000000000000000000000000000000000000000851660608301527fffff00000000000000000000000000000000000000000000000000000000000084166080808401919091528351808403909101815260a09092018352815191810191909120600081815260049092529190205460ff1661347d576040517f097e17ff00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8087166004830152851660248201527fffffffffffffffffffff00000000000000000000000000000000000000000000841660448201527fffff00000000000000000000000000000000000000000000000000000000000083166064820152608401611156565b5050505050565b6000806134918486615c5d565b9050600060248260ff1611156134cb576134af602460ff8416615adf565b6134ba90600a615d96565b6134c490856158ba565b90506134f1565b6134d960ff83166024615adf565b6134e490600a615d96565b6134ee90856158a3565b90505b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811115611a2c576040517f10cb51d100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081015160005b81518110156135e257600082828151811061356c5761356c61595d565b6020026020010151905061358a81600261428c90919063ffffffff16565b156135d95760405173ffffffffffffffffffffffffffffffffffffffff821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b5060010161354f565b50815160005b8151811015610fe45760008282815181106136055761360561595d565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603613675576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61368060028261426a565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a1506001016135e8565b60005b81518110156110285760008282815181106136f3576136f361595d565b602002602001015160000151905060008383815181106137155761371561595d565b60209081029190910181015181015173ffffffffffffffffffffffffffffffffffffffff841660008181526008845260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff85169081179091559051908152919350917fbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d910160405180910390a250506001016136d6565b60408051808201909152600080825260208201526000826000015190506000808273ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015613828573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061384c9190615dbc565b50935050925050600082121561388e576040517f10cb51d100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061390d8473ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa1580156138de573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139029190615e0c565b876020015185613484565b604080518082019091527bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909116815263ffffffff909216602083015250949350505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610bec565b836040015163ffffffff168311156139d75760408085015190517f8693378900000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015260248101849052604401611156565b836020015161ffff16821115613a2c5760208401516040517fd88dddd60000000000000000000000000000000000000000000000000000000081526004810184905261ffff9091166024820152604401611156565b610fe4846102000151826141f6565b67ffffffffffffffff821660009081526005602090815260408083208151808301909252547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116825263ffffffff7c010000000000000000000000000000000000000000000000000000000090910481169282019290925290831615613b33576000816020015163ffffffff1642613ad09190615adf565b90508363ffffffff16811115613b31576040517ff08bcb3e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015263ffffffff8516602482015260448101829052606401611156565b505b519392505050565b6000808083815b81811015613e05576000878783818110613b5e57613b5e61595d565b905060400201803603810190613b749190615e29565b67ffffffffffffffff8c166000908152600a60209081526040808320845173ffffffffffffffffffffffffffffffffffffffff168452825291829020825160c081018452905463ffffffff8082168352640100000000820481169383019390935268010000000000000000810461ffff16938201939093526a01000000000000000000008304821660608201526e01000000000000000000000000000083049091166080820152720100000000000000000000000000000000000090910460ff16151560a0820181905291925090613c94576101208d0151613c619061ffff16662386f26fc100006158a3565b613c6b9088615b57565b96508c610140015186613c7e9190615e62565b9550613c8b602086615e62565b94505050613dfd565b604081015160009061ffff1615613d4d5760008c73ffffffffffffffffffffffffffffffffffffffff16846000015173ffffffffffffffffffffffffffffffffffffffff1614613cf0578351613ce990612047565b9050613cf3565b508a5b620186a0836040015161ffff16613d358660200151847bffffffffffffffffffffffffffffffffffffffffffffffffffffffff166142ae90919063ffffffff16565b613d3f91906158a3565b613d4991906158ba565b9150505b6060820151613d5c9088615e62565b9650816080015186613d6e9190615e62565b8251909650600090613d8d9063ffffffff16662386f26fc100006158a3565b905080821015613dac57613da1818a615b57565b985050505050613dfd565b6000836020015163ffffffff16662386f26fc10000613dcb91906158a3565b905080831115613deb57613ddf818b615b57565b99505050505050613dfd565b613df5838b615b57565b995050505050505b600101613b42565b505096509650969350505050565b60008063ffffffff8316613e29610120866158a3565b613e35876101e0615b57565b613e3f9190615b57565b613e499190615b57565b905060008760c0015163ffffffff168860e0015161ffff1683613e6c91906158a3565b613e769190615b57565b61010089015190915061ffff16613e9d6dffffffffffffffffffffffffffff8916836158a3565b613ea791906158a3565b613eb790655af3107a40006158a3565b98975050505050505050565b60408051808201909152600080825260208201526000613eef858585610160015163ffffffff16612356565b9050826060015163ffffffff1681600001511115613f39576040517f4c4fc93a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826101e001518015613f4d57508060200151155b15610be9576040517fee433e9900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603613fd3576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005b81518110156110285760008282815181106140685761406861595d565b60209081029190910181015180518183015173ffffffffffffffffffffffffffffffffffffffff80831660008181526007875260409081902084518154868a0180518589018051949098167fffffffffffffffffffffff00000000000000000000000000000000000000000090931683177401000000000000000000000000000000000000000060ff92831602177fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff1675010000000000000000000000000000000000000000009415159490940293909317909355835190815291511697810197909752915115159186019190915292945090929091907fe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf9060600160405180910390a250505080600101905061404b565b6060816000018054806020026020016040519081016040528092919081815260200182805480156141ea57602002820191906000526020600020905b8154815260200190600101908083116141d6575b50505050509050919050565b7fd7ed2ad4000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000831601611028576131c5816142eb565b6000610bec8373ffffffffffffffffffffffffffffffffffffffff841661439e565b6000610bec8373ffffffffffffffffffffffffffffffffffffffff8416614498565b6000610bec8373ffffffffffffffffffffffffffffffffffffffff84166144e7565b6000670de0b6b3a76400006142e1837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff86166158a3565b610bec91906158ba565b6000815160201461432a57816040517f8d666f60000000000000000000000000000000000000000000000000000000008152600401611156919061466e565b6000828060200190518101906143409190615c06565b905073ffffffffffffffffffffffffffffffffffffffff811180614365575061040081105b15610d1e57826040517f8d666f60000000000000000000000000000000000000000000000000000000008152600401611156919061466e565b600081815260018301602052604081205480156144875760006143c2600183615adf565b85549091506000906143d690600190615adf565b905080821461443b5760008660000182815481106143f6576143f661595d565b90600052602060002001549050808760000184815481106144195761441961595d565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061444c5761444c615e7f565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610d1e565b6000915050610d1e565b5092915050565b60008181526001830160205260408120546144df57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610d1e565b506000610d1e565b6000818152600183016020526040812054801561448757600061450b600183615adf565b855490915060009061451f90600190615adf565b905081811461443b5760008660000182815481106143f6576143f661595d565b803573ffffffffffffffffffffffffffffffffffffffff8116811461456357600080fd5b919050565b60008060006060848603121561457d57600080fd5b6145868461453f565b92506020840135915061459b6040850161453f565b90509250925092565b80357fffffffff000000000000000000000000000000000000000000000000000000008116811461456357600080fd5b6000602082840312156145e657600080fd5b610bec826145a4565b60006020828403121561460157600080fd5b610bec8261453f565b6000815180845260005b8181101561463057602081850181015186830182015201614614565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610bec602083018461460a565b6020808252825182820181905260009190848201906040850190845b818110156146cf57835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161469d565b50909695505050505050565b6000602082840312156146ed57600080fd5b813567ffffffffffffffff81111561470457600080fd5b820160408185031215610bec57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff8111828210171561476857614768614716565b60405290565b6040805190810167ffffffffffffffff8111828210171561476857614768614716565b604051610220810167ffffffffffffffff8111828210171561476857614768614716565b60405160c0810167ffffffffffffffff8111828210171561476857614768614716565b6040516060810167ffffffffffffffff8111828210171561476857614768614716565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561484257614842614716565b604052919050565b600067ffffffffffffffff82111561486457614864614716565b5060051b60200190565b801515811461132157600080fd5b80356145638161486e565b6000602080838503121561489a57600080fd5b823567ffffffffffffffff8111156148b157600080fd5b8301601f810185136148c257600080fd5b80356148d56148d08261484a565b6147fb565b81815260a091820283018401918482019190888411156148f457600080fd5b938501935b838510156149c75780858a0312156149115760008081fd5b614919614745565b6149228661453f565b8152868601357fffffffffffffffffffff00000000000000000000000000000000000000000000811681146149575760008081fd5b818801526040868101357fffff000000000000000000000000000000000000000000000000000000000000811681146149905760008081fd5b9082015260606149a187820161453f565b908201526080868101356149b48161486e565b90820152835293840193918501916148f9565b50979650505050505050565b803567ffffffffffffffff8116811461456357600080fd5b60008083601f8401126149fd57600080fd5b50813567ffffffffffffffff811115614a1557600080fd5b60208301915083602082850101111561204057600080fd5b60008083601f840112614a3f57600080fd5b50813567ffffffffffffffff811115614a5757600080fd5b6020830191508360208260051b850101111561204057600080fd5b600080600080600080600080600060c08a8c031215614a9057600080fd5b614a998a6149d3565b9850614aa760208b0161453f565b975060408a0135965060608a013567ffffffffffffffff80821115614acb57600080fd5b614ad78d838e016149eb565b909850965060808c0135915080821115614af057600080fd5b614afc8d838e01614a2d565b909650945060a08c0135915080821115614b1557600080fd5b818c0191508c601f830112614b2957600080fd5b813581811115614b3857600080fd5b8d60208260061b8501011115614b4d57600080fd5b6020830194508093505050509295985092959850929598565b848152600060208515158184015260806040840152614b88608084018661460a565b8381036060850152845180825282820190600581901b8301840184880160005b83811015614bf4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018552614be283835161460a565b94870194925090860190600101614ba8565b50909b9a5050505050505050505050565b60008060208385031215614c1857600080fd5b823567ffffffffffffffff811115614c2f57600080fd5b614c3b85828601614a2d565b90969095509350505050565b602080825282518282018190526000919060409081850190868401855b82811015614cb557614ca584835180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16825260209081015163ffffffff16910152565b9284019290850190600101614c64565b5091979650505050505050565b600060208284031215614cd457600080fd5b610bec826149d3565b81517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260208083015163ffffffff169082015260408101610d1e565b803561ffff8116811461456357600080fd5b803563ffffffff8116811461456357600080fd5b60006020808385031215614d5157600080fd5b823567ffffffffffffffff811115614d6857600080fd5b8301601f81018513614d7957600080fd5b8035614d876148d08261484a565b8181526102409182028301840191848201919088841115614da757600080fd5b938501935b838510156149c75784890381811215614dc55760008081fd5b614dcd61476e565b614dd6876149d3565b8152610220807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084011215614e0b5760008081fd5b614e13614791565b9250614e2089890161487c565b83526040614e2f818a01614d18565b8a8501526060614e40818b01614d2a565b8286015260809150614e53828b01614d2a565b9085015260a0614e648a8201614d2a565b8286015260c09150614e77828b01614d18565b9085015260e0614e888a8201614d2a565b828601526101009150614e9c828b01614d18565b90850152610120614eae8a8201614d18565b828601526101409150614ec2828b01614d18565b90850152610160614ed48a8201614d2a565b828601526101809150614ee8828b01614d2a565b908501526101a0614efa8a82016149d3565b828601526101c09150614f0e828b01614d2a565b908501526101e0614f208a8201614d2a565b828601526102009150614f34828b0161487c565b90850152614f438983016145a4565b90840152508088019190915283529384019391850191614dac565b81511515815261022081016020830151614f7e602084018261ffff169052565b506040830151614f96604084018263ffffffff169052565b506060830151614fae606084018263ffffffff169052565b506080830151614fc6608084018263ffffffff169052565b5060a0830151614fdc60a084018261ffff169052565b5060c0830151614ff460c084018263ffffffff169052565b5060e083015161500a60e084018261ffff169052565b506101008381015161ffff9081169184019190915261012080850151909116908301526101408084015163ffffffff90811691840191909152610160808501518216908401526101808085015167ffffffffffffffff16908401526101a0808501518216908401526101c080850151909116908301526101e080840151151590830152610200808401517fffffffff000000000000000000000000000000000000000000000000000000008116828501525b505092915050565b600082601f8301126150d557600080fd5b813560206150e56148d08361484a565b82815260069290921b8401810191818101908684111561510457600080fd5b8286015b8481101561515157604081890312156151215760008081fd5b61512961476e565b615132826149d3565b815261513f85830161453f565b81860152835291830191604001615108565b509695505050505050565b6000806040838503121561516f57600080fd5b67ffffffffffffffff8335111561518557600080fd5b83601f84358501011261519757600080fd5b6151a76148d0843585013561484a565b8335840180358083526020808401939260059290921b909101018610156151cd57600080fd5b602085358601015b85358601803560051b016020018110156153da5767ffffffffffffffff813511156151ff57600080fd5b8035863587010160407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828a0301121561523857600080fd5b61524061476e565b61524c602083016149d3565b815267ffffffffffffffff6040830135111561526757600080fd5b88603f60408401358401011261527c57600080fd5b6152926148d0602060408501358501013561484a565b6020604084810135850182810135808552928401939260e00201018b10156152b957600080fd5b6040808501358501015b6040858101358601602081013560e00201018110156153bb5760e0818d0312156152ec57600080fd5b6152f461476e565b6152fd8261453f565b815260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0838f0301121561533157600080fd5b6153396147b5565b61534560208401614d2a565b815261535360408401614d2a565b602082015261536460608401614d18565b604082015261537560808401614d2a565b606082015261538660a08401614d2a565b608082015261539860c084013561486e565b60c083013560a0820152602082810191909152908452929092019160e0016152c3565b50806020840152505080855250506020830192506020810190506151d5565b5092505067ffffffffffffffff602084013511156153f757600080fd5b61540784602085013585016150c4565b90509250929050565b600082601f83011261542157600080fd5b813560206154316148d08361484a565b8083825260208201915060208460051b87010193508684111561545357600080fd5b602086015b84811015615151576154698161453f565b8352918301918301615458565b6000806040838503121561548957600080fd5b823567ffffffffffffffff808211156154a157600080fd5b6154ad86838701615410565b935060208501359150808211156154c357600080fd5b506154d085828601615410565b9150509250929050565b600080600080604085870312156154f057600080fd5b843567ffffffffffffffff8082111561550857600080fd5b615514888389016149eb565b9096509450602087013591508082111561552d57600080fd5b5061553a878288016149eb565b95989497509550505050565b6000806040838503121561555957600080fd5b615562836149d3565b91506154076020840161453f565b60006020828403121561558257600080fd5b813567ffffffffffffffff8082111561559a57600080fd5b90830190604082860312156155ae57600080fd5b6155b661476e565b8235828111156155c557600080fd5b6155d187828601615410565b8252506020830135828111156155e657600080fd5b6155f287828601615410565b60208301525095945050505050565b6000602080838503121561561457600080fd5b823567ffffffffffffffff81111561562b57600080fd5b8301601f8101851361563c57600080fd5b803561564a6148d08261484a565b81815260069190911b8201830190838101908783111561566957600080fd5b928401925b828410156156bb57604084890312156156875760008081fd5b61568f61476e565b6156988561453f565b81526156a58686016149d3565b818701528252604093909301929084019061566e565b979650505050505050565b600080604083850312156156d957600080fd5b6156e2836149d3565b9150602083013567ffffffffffffffff8111156156fe57600080fd5b830160a0818603121561571057600080fd5b809150509250929050565b60ff8116811461132157600080fd5b6000602080838503121561573d57600080fd5b823567ffffffffffffffff81111561575457600080fd5b8301601f8101851361576557600080fd5b80356157736148d08261484a565b81815260079190911b8201830190838101908783111561579257600080fd5b928401925b828410156156bb5783880360808112156157b15760008081fd5b6157b961476e565b6157c28661453f565b81526060807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0840112156157f65760008081fd5b6157fe6147d8565b925061580b88880161453f565b835260408088013561581c8161571b565b848a0152908701359061582e8261486e565b8301528087019190915282526080939093019290840190615797565b6000806040838503121561585d57600080fd5b6158668361453f565b9150615407602084016149d3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610d1e57610d1e615874565b6000826158f0577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261592a57600080fd5b83018035915067ffffffffffffffff82111561594557600080fd5b6020019150600681901b360382131561204057600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461456357600080fd5b6000604082840312156159ca57600080fd5b6159d261476e565b6159db8361453f565b81526159e96020840161598c565b60208201529392505050565b600060408284031215615a0757600080fd5b615a0f61476e565b6159db836149d3565b60006020808385031215615a2b57600080fd5b823567ffffffffffffffff811115615a4257600080fd5b8301601f81018513615a5357600080fd5b8035615a616148d08261484a565b81815260609182028301840191848201919088841115615a8057600080fd5b938501935b838510156149c75780858a031215615a9d5760008081fd5b615aa56147d8565b615aae8661453f565b8152615abb87870161598c565b878201526040615acc818801614d2a565b9082015283529384019391850191615a85565b81810381811115610d1e57610d1e615874565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112615b2757600080fd5b83018035915067ffffffffffffffff821115615b4257600080fd5b60200191503681900382131561204057600080fd5b80820180821115610d1e57610d1e615874565b7fffffffff0000000000000000000000000000000000000000000000000000000081358181169160048510156150bc5760049490940360031b84901b1690921692915050565b60008085851115615bc057600080fd5b83861115615bcd57600080fd5b5050820193919092039150565b600060408284031215615bec57600080fd5b615bf461476e565b8251815260208301516159e98161486e565b600060208284031215615c1857600080fd5b5051919050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61833603018112615c5357600080fd5b9190910192915050565b60ff8181168382160190811115610d1e57610d1e615874565b600181815b80851115615ccf57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115615cb557615cb5615874565b80851615615cc257918102915b93841c9390800290615c7b565b509250929050565b600082615ce657506001610d1e565b81615cf357506000610d1e565b8160018114615d095760028114615d1357615d2f565b6001915050610d1e565b60ff841115615d2457615d24615874565b50506001821b610d1e565b5060208310610133831016604e8410600b8410161715615d52575081810a610d1e565b615d5c8383615c76565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115615d8e57615d8e615874565b029392505050565b6000610bec8383615cd7565b805169ffffffffffffffffffff8116811461456357600080fd5b600080600080600060a08688031215615dd457600080fd5b615ddd86615da2565b9450602086015193506040860151925060608601519150615e0060808701615da2565b90509295509295909350565b600060208284031215615e1e57600080fd5b8151610bec8161571b565b600060408284031215615e3b57600080fd5b615e4361476e565b615e4c8361453f565b8152602083013560208201528091505092915050565b63ffffffff81811683821601908082111561449157614491615874565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var FeeQuoterABI = FeeQuoterMetaData.ABI

var FeeQuoterBin = FeeQuoterMetaData.Bin

func DeployFeeQuoter(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig FeeQuoterStaticConfig, priceUpdaters []common.Address, feeTokens []common.Address, tokenPriceFeeds []FeeQuoterTokenPriceFeedUpdate, tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs, destChainConfigArgs []FeeQuoterDestChainConfigArgs) (common.Address, *generated.Transaction, *FeeQuoter, error) {
	parsed, err := FeeQuoterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(FeeQuoterZKBin), backend, staticConfig, priceUpdaters, feeTokens, tokenPriceFeeds, tokenTransferFeeConfigArgs, premiumMultiplierWeiPerEthArgs, destChainConfigArgs)
		contractReturn := &FeeQuoter{address: address, abi: *parsed, FeeQuoterCaller: FeeQuoterCaller{contract: contractBind}, FeeQuoterTransactor: FeeQuoterTransactor{contract: contractBind}, FeeQuoterFilterer: FeeQuoterFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeeQuoterBin), backend, staticConfig, priceUpdaters, feeTokens, tokenPriceFeeds, tokenTransferFeeConfigArgs, premiumMultiplierWeiPerEthArgs, destChainConfigArgs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, Hash_zks: tx.Hash()}, &FeeQuoter{address: address, abi: *parsed, FeeQuoterCaller: FeeQuoterCaller{contract: contract}, FeeQuoterTransactor: FeeQuoterTransactor{contract: contract}, FeeQuoterFilterer: FeeQuoterFilterer{contract: contract}}, nil
}

type FeeQuoter struct {
	address common.Address
	abi     abi.ABI
	FeeQuoterCaller
	FeeQuoterTransactor
	FeeQuoterFilterer
}

type FeeQuoterCaller struct {
	contract *bind.BoundContract
}

type FeeQuoterTransactor struct {
	contract *bind.BoundContract
}

type FeeQuoterFilterer struct {
	contract *bind.BoundContract
}

type FeeQuoterSession struct {
	Contract     *FeeQuoter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FeeQuoterCallerSession struct {
	Contract *FeeQuoterCaller
	CallOpts bind.CallOpts
}

type FeeQuoterTransactorSession struct {
	Contract     *FeeQuoterTransactor
	TransactOpts bind.TransactOpts
}

type FeeQuoterRaw struct {
	Contract *FeeQuoter
}

type FeeQuoterCallerRaw struct {
	Contract *FeeQuoterCaller
}

type FeeQuoterTransactorRaw struct {
	Contract *FeeQuoterTransactor
}

func NewFeeQuoter(address common.Address, backend bind.ContractBackend) (*FeeQuoter, error) {
	abi, err := abi.JSON(strings.NewReader(FeeQuoterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFeeQuoter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeQuoter{address: address, abi: abi, FeeQuoterCaller: FeeQuoterCaller{contract: contract}, FeeQuoterTransactor: FeeQuoterTransactor{contract: contract}, FeeQuoterFilterer: FeeQuoterFilterer{contract: contract}}, nil
}

func NewFeeQuoterCaller(address common.Address, caller bind.ContractCaller) (*FeeQuoterCaller, error) {
	contract, err := bindFeeQuoter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterCaller{contract: contract}, nil
}

func NewFeeQuoterTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeQuoterTransactor, error) {
	contract, err := bindFeeQuoter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterTransactor{contract: contract}, nil
}

func NewFeeQuoterFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeQuoterFilterer, error) {
	contract, err := bindFeeQuoter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterFilterer{contract: contract}, nil
}

func bindFeeQuoter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeQuoterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FeeQuoter *FeeQuoterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeQuoter.Contract.FeeQuoterCaller.contract.Call(opts, result, method, params...)
}

func (_FeeQuoter *FeeQuoterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeQuoter.Contract.FeeQuoterTransactor.contract.Transfer(opts)
}

func (_FeeQuoter *FeeQuoterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeQuoter.Contract.FeeQuoterTransactor.contract.Transact(opts, method, params...)
}

func (_FeeQuoter *FeeQuoterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeQuoter.Contract.contract.Call(opts, result, method, params...)
}

func (_FeeQuoter *FeeQuoterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeQuoter.Contract.contract.Transfer(opts)
}

func (_FeeQuoter *FeeQuoterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeQuoter.Contract.contract.Transact(opts, method, params...)
}

func (_FeeQuoter *FeeQuoterCaller) FEEBASEDECIMALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "FEE_BASE_DECIMALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) FEEBASEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.FEEBASEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) FEEBASEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.FEEBASEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) KEYSTONEPRICEDECIMALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "KEYSTONE_PRICE_DECIMALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) KEYSTONEPRICEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.KEYSTONEPRICEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) KEYSTONEPRICEDECIMALS() (*big.Int, error) {
	return _FeeQuoter.Contract.KEYSTONEPRICEDECIMALS(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "convertTokenAmount", fromToken, fromTokenAmount, toToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.ConvertTokenAmount(&_FeeQuoter.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_FeeQuoter *FeeQuoterCallerSession) ConvertTokenAmount(fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.ConvertTokenAmount(&_FeeQuoter.CallOpts, fromToken, fromTokenAmount, toToken)
}

func (_FeeQuoter *FeeQuoterCaller) GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getAllAuthorizedCallers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetAllAuthorizedCallers(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetAllAuthorizedCallers(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (FeeQuoterDestChainConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	if err != nil {
		return *new(FeeQuoterDestChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterDestChainConfig)).(*FeeQuoterDestChainConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetDestChainConfig(destChainSelector uint64) (FeeQuoterDestChainConfig, error) {
	return _FeeQuoter.Contract.GetDestChainConfig(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetDestChainConfig(destChainSelector uint64) (FeeQuoterDestChainConfig, error) {
	return _FeeQuoter.Contract.GetDestChainConfig(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCaller) GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getDestinationChainGasPrice", destChainSelector)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetDestinationChainGasPrice(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetDestinationChainGasPrice(destChainSelector uint64) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetDestinationChainGasPrice(&_FeeQuoter.CallOpts, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCaller) GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getFeeTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetFeeTokens() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetFeeTokens(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetFeeTokens() ([]common.Address, error) {
	return _FeeQuoter.Contract.GetFeeTokens(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getPremiumMultiplierWeiPerEth", token)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _FeeQuoter.Contract.GetPremiumMultiplierWeiPerEth(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetPremiumMultiplierWeiPerEth(token common.Address) (uint64, error) {
	return _FeeQuoter.Contract.GetPremiumMultiplierWeiPerEth(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetStaticConfig(opts *bind.CallOpts) (FeeQuoterStaticConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(FeeQuoterStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterStaticConfig)).(*FeeQuoterStaticConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetStaticConfig() (FeeQuoterStaticConfig, error) {
	return _FeeQuoter.Contract.GetStaticConfig(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetStaticConfig() (FeeQuoterStaticConfig, error) {
	return _FeeQuoter.Contract.GetStaticConfig(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenAndGasPrices", token, destChainSelector)

	outstruct := new(GetTokenAndGasPrices)
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenPrice = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.GasPriceValue = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _FeeQuoter.Contract.GetTokenAndGasPrices(&_FeeQuoter.CallOpts, token, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenAndGasPrices(token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

	error) {
	return _FeeQuoter.Contract.GetTokenAndGasPrices(&_FeeQuoter.CallOpts, token, destChainSelector)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenPrice", token)

	if err != nil {
		return *new(InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalTimestampedPackedUint224)).(*InternalTimestampedPackedUint224)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenPrice(token common.Address) (InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (FeeQuoterTokenPriceFeedConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenPriceFeedConfig", token)

	if err != nil {
		return *new(FeeQuoterTokenPriceFeedConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterTokenPriceFeedConfig)).(*FeeQuoterTokenPriceFeedConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenPriceFeedConfig(token common.Address) (FeeQuoterTokenPriceFeedConfig, error) {
	return _FeeQuoter.Contract.GetTokenPriceFeedConfig(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenPriceFeedConfig(token common.Address) (FeeQuoterTokenPriceFeedConfig, error) {
	return _FeeQuoter.Contract.GetTokenPriceFeedConfig(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenPrices", tokens)

	if err != nil {
		return *new([]InternalTimestampedPackedUint224), err
	}

	out0 := *abi.ConvertType(out[0], new([]InternalTimestampedPackedUint224)).(*[]InternalTimestampedPackedUint224)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrices(&_FeeQuoter.CallOpts, tokens)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenPrices(tokens []common.Address) ([]InternalTimestampedPackedUint224, error) {
	return _FeeQuoter.Contract.GetTokenPrices(&_FeeQuoter.CallOpts, tokens)
}

func (_FeeQuoter *FeeQuoterCaller) GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getTokenTransferFeeConfig", destChainSelector, token)

	if err != nil {
		return *new(FeeQuoterTokenTransferFeeConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeQuoterTokenTransferFeeConfig)).(*FeeQuoterTokenTransferFeeConfig)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetTokenTransferFeeConfig(destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error) {
	return _FeeQuoter.Contract.GetTokenTransferFeeConfig(&_FeeQuoter.CallOpts, destChainSelector, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetTokenTransferFeeConfig(destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error) {
	return _FeeQuoter.Contract.GetTokenTransferFeeConfig(&_FeeQuoter.CallOpts, destChainSelector, token)
}

func (_FeeQuoter *FeeQuoterCaller) GetValidatedFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getValidatedFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetValidatedFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedFee(&_FeeQuoter.CallOpts, destChainSelector, message)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetValidatedFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedFee(&_FeeQuoter.CallOpts, destChainSelector, message)
}

func (_FeeQuoter *FeeQuoterCaller) GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "getValidatedTokenPrice", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCallerSession) GetValidatedTokenPrice(token common.Address) (*big.Int, error) {
	return _FeeQuoter.Contract.GetValidatedTokenPrice(&_FeeQuoter.CallOpts, token)
}

func (_FeeQuoter *FeeQuoterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) Owner() (common.Address, error) {
	return _FeeQuoter.Contract.Owner(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) Owner() (common.Address, error) {
	return _FeeQuoter.Contract.Owner(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCaller) ProcessMessageArgs(opts *bind.CallOpts, destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

	error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "processMessageArgs", destChainSelector, feeToken, feeTokenAmount, extraArgs, onRampTokenTransfers, sourceTokenAmounts)

	outstruct := new(ProcessMessageArgs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MsgFeeJuels = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.IsOutOfOrderExecution = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.ConvertedExtraArgs = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.DestExecDataPerToken = *abi.ConvertType(out[3], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

func (_FeeQuoter *FeeQuoterSession) ProcessMessageArgs(destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

	error) {
	return _FeeQuoter.Contract.ProcessMessageArgs(&_FeeQuoter.CallOpts, destChainSelector, feeToken, feeTokenAmount, extraArgs, onRampTokenTransfers, sourceTokenAmounts)
}

func (_FeeQuoter *FeeQuoterCallerSession) ProcessMessageArgs(destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

	error) {
	return _FeeQuoter.Contract.ProcessMessageArgs(&_FeeQuoter.CallOpts, destChainSelector, feeToken, feeTokenAmount, extraArgs, onRampTokenTransfers, sourceTokenAmounts)
}

func (_FeeQuoter *FeeQuoterCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeQuoter.Contract.SupportsInterface(&_FeeQuoter.CallOpts, interfaceId)
}

func (_FeeQuoter *FeeQuoterCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeQuoter.Contract.SupportsInterface(&_FeeQuoter.CallOpts, interfaceId)
}

func (_FeeQuoter *FeeQuoterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeQuoter.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_FeeQuoter *FeeQuoterSession) TypeAndVersion() (string, error) {
	return _FeeQuoter.Contract.TypeAndVersion(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterCallerSession) TypeAndVersion() (string, error) {
	return _FeeQuoter.Contract.TypeAndVersion(&_FeeQuoter.CallOpts)
}

func (_FeeQuoter *FeeQuoterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "acceptOwnership")
}

func (_FeeQuoter *FeeQuoterSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeeQuoter.Contract.AcceptOwnership(&_FeeQuoter.TransactOpts)
}

func (_FeeQuoter *FeeQuoterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeeQuoter.Contract.AcceptOwnership(&_FeeQuoter.TransactOpts)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyAuthorizedCallerUpdates", authorizedCallerArgs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyAuthorizedCallerUpdates(&_FeeQuoter.TransactOpts, authorizedCallerArgs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyAuthorizedCallerUpdates(&_FeeQuoter.TransactOpts, authorizedCallerArgs)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyDestChainConfigUpdates(destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyDestChainConfigUpdates(&_FeeQuoter.TransactOpts, destChainConfigArgs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyDestChainConfigUpdates(&_FeeQuoter.TransactOpts, destChainConfigArgs)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyFeeTokensUpdates", feeTokensToRemove, feeTokensToAdd)
}

func (_FeeQuoter *FeeQuoterSession) ApplyFeeTokensUpdates(feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyFeeTokensUpdates(&_FeeQuoter.TransactOpts, feeTokensToRemove, feeTokensToAdd)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyFeeTokensUpdates(feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyFeeTokensUpdates(&_FeeQuoter.TransactOpts, feeTokensToRemove, feeTokensToAdd)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyPremiumMultiplierWeiPerEthUpdates", premiumMultiplierWeiPerEthArgs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_FeeQuoter.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyPremiumMultiplierWeiPerEthUpdates(&_FeeQuoter.TransactOpts, premiumMultiplierWeiPerEthArgs)
}

func (_FeeQuoter *FeeQuoterTransactor) ApplyTokenTransferFeeConfigUpdates(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "applyTokenTransferFeeConfigUpdates", tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_FeeQuoter *FeeQuoterSession) ApplyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyTokenTransferFeeConfigUpdates(&_FeeQuoter.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_FeeQuoter *FeeQuoterTransactorSession) ApplyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error) {
	return _FeeQuoter.Contract.ApplyTokenTransferFeeConfigUpdates(&_FeeQuoter.TransactOpts, tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs)
}

func (_FeeQuoter *FeeQuoterTransactor) OnReport(opts *bind.TransactOpts, metadata []byte, report []byte) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "onReport", metadata, report)
}

func (_FeeQuoter *FeeQuoterSession) OnReport(metadata []byte, report []byte) (*types.Transaction, error) {
	return _FeeQuoter.Contract.OnReport(&_FeeQuoter.TransactOpts, metadata, report)
}

func (_FeeQuoter *FeeQuoterTransactorSession) OnReport(metadata []byte, report []byte) (*types.Transaction, error) {
	return _FeeQuoter.Contract.OnReport(&_FeeQuoter.TransactOpts, metadata, report)
}

func (_FeeQuoter *FeeQuoterTransactor) SetReportPermissions(opts *bind.TransactOpts, permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "setReportPermissions", permissions)
}

func (_FeeQuoter *FeeQuoterSession) SetReportPermissions(permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error) {
	return _FeeQuoter.Contract.SetReportPermissions(&_FeeQuoter.TransactOpts, permissions)
}

func (_FeeQuoter *FeeQuoterTransactorSession) SetReportPermissions(permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error) {
	return _FeeQuoter.Contract.SetReportPermissions(&_FeeQuoter.TransactOpts, permissions)
}

func (_FeeQuoter *FeeQuoterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "transferOwnership", to)
}

func (_FeeQuoter *FeeQuoterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.TransferOwnership(&_FeeQuoter.TransactOpts, to)
}

func (_FeeQuoter *FeeQuoterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeeQuoter.Contract.TransferOwnership(&_FeeQuoter.TransactOpts, to)
}

func (_FeeQuoter *FeeQuoterTransactor) UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "updatePrices", priceUpdates)
}

func (_FeeQuoter *FeeQuoterSession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdatePrices(&_FeeQuoter.TransactOpts, priceUpdates)
}

func (_FeeQuoter *FeeQuoterTransactorSession) UpdatePrices(priceUpdates InternalPriceUpdates) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdatePrices(&_FeeQuoter.TransactOpts, priceUpdates)
}

func (_FeeQuoter *FeeQuoterTransactor) UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _FeeQuoter.contract.Transact(opts, "updateTokenPriceFeeds", tokenPriceFeedUpdates)
}

func (_FeeQuoter *FeeQuoterSession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdateTokenPriceFeeds(&_FeeQuoter.TransactOpts, tokenPriceFeedUpdates)
}

func (_FeeQuoter *FeeQuoterTransactorSession) UpdateTokenPriceFeeds(tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error) {
	return _FeeQuoter.Contract.UpdateTokenPriceFeeds(&_FeeQuoter.TransactOpts, tokenPriceFeedUpdates)
}

type FeeQuoterAuthorizedCallerAddedIterator struct {
	Event *FeeQuoterAuthorizedCallerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterAuthorizedCallerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterAuthorizedCallerAdded)
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
		it.Event = new(FeeQuoterAuthorizedCallerAdded)
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

func (it *FeeQuoterAuthorizedCallerAddedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterAuthorizedCallerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterAuthorizedCallerAdded struct {
	Caller common.Address
	Raw    types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerAddedIterator, error) {

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return &FeeQuoterAuthorizedCallerAddedIterator{contract: _FeeQuoter.contract, event: "AuthorizedCallerAdded", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerAdded) (event.Subscription, error) {

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterAuthorizedCallerAdded)
				if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseAuthorizedCallerAdded(log types.Log) (*FeeQuoterAuthorizedCallerAdded, error) {
	event := new(FeeQuoterAuthorizedCallerAdded)
	if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterAuthorizedCallerRemovedIterator struct {
	Event *FeeQuoterAuthorizedCallerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterAuthorizedCallerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterAuthorizedCallerRemoved)
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
		it.Event = new(FeeQuoterAuthorizedCallerRemoved)
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

func (it *FeeQuoterAuthorizedCallerRemovedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterAuthorizedCallerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterAuthorizedCallerRemoved struct {
	Caller common.Address
	Raw    types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerRemovedIterator, error) {

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return &FeeQuoterAuthorizedCallerRemovedIterator{contract: _FeeQuoter.contract, event: "AuthorizedCallerRemoved", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerRemoved) (event.Subscription, error) {

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterAuthorizedCallerRemoved)
				if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseAuthorizedCallerRemoved(log types.Log) (*FeeQuoterAuthorizedCallerRemoved, error) {
	event := new(FeeQuoterAuthorizedCallerRemoved)
	if err := _FeeQuoter.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterDestChainAddedIterator struct {
	Event *FeeQuoterDestChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterDestChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterDestChainAdded)
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
		it.Event = new(FeeQuoterDestChainAdded)
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

func (it *FeeQuoterDestChainAddedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterDestChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterDestChainAdded struct {
	DestChainSelector uint64
	DestChainConfig   FeeQuoterDestChainConfig
	Raw               types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterDestChainAddedIterator{contract: _FeeQuoter.contract, event: "DestChainAdded", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "DestChainAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterDestChainAdded)
				if err := _FeeQuoter.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseDestChainAdded(log types.Log) (*FeeQuoterDestChainAdded, error) {
	event := new(FeeQuoterDestChainAdded)
	if err := _FeeQuoter.contract.UnpackLog(event, "DestChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterDestChainConfigUpdatedIterator struct {
	Event *FeeQuoterDestChainConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterDestChainConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterDestChainConfigUpdated)
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
		it.Event = new(FeeQuoterDestChainConfigUpdated)
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

func (it *FeeQuoterDestChainConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterDestChainConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterDestChainConfigUpdated struct {
	DestChainSelector uint64
	DestChainConfig   FeeQuoterDestChainConfig
	Raw               types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterDestChainConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainConfigUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "DestChainConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterDestChainConfigUpdatedIterator{contract: _FeeQuoter.contract, event: "DestChainConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainConfigUpdated, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "DestChainConfigUpdated", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterDestChainConfigUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseDestChainConfigUpdated(log types.Log) (*FeeQuoterDestChainConfigUpdated, error) {
	event := new(FeeQuoterDestChainConfigUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "DestChainConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterFeeTokenAddedIterator struct {
	Event *FeeQuoterFeeTokenAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterFeeTokenAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterFeeTokenAdded)
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
		it.Event = new(FeeQuoterFeeTokenAdded)
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

func (it *FeeQuoterFeeTokenAddedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterFeeTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterFeeTokenAdded struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenAddedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterFeeTokenAddedIterator{contract: _FeeQuoter.contract, event: "FeeTokenAdded", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenAdded, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "FeeTokenAdded", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterFeeTokenAdded)
				if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseFeeTokenAdded(log types.Log) (*FeeQuoterFeeTokenAdded, error) {
	event := new(FeeQuoterFeeTokenAdded)
	if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterFeeTokenRemovedIterator struct {
	Event *FeeQuoterFeeTokenRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterFeeTokenRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterFeeTokenRemoved)
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
		it.Event = new(FeeQuoterFeeTokenRemoved)
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

func (it *FeeQuoterFeeTokenRemovedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterFeeTokenRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterFeeTokenRemoved struct {
	FeeToken common.Address
	Raw      types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenRemovedIterator, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterFeeTokenRemovedIterator{contract: _FeeQuoter.contract, event: "FeeTokenRemoved", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error) {

	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "FeeTokenRemoved", feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterFeeTokenRemoved)
				if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseFeeTokenRemoved(log types.Log) (*FeeQuoterFeeTokenRemoved, error) {
	event := new(FeeQuoterFeeTokenRemoved)
	if err := _FeeQuoter.contract.UnpackLog(event, "FeeTokenRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterOwnershipTransferRequestedIterator struct {
	Event *FeeQuoterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterOwnershipTransferRequested)
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
		it.Event = new(FeeQuoterOwnershipTransferRequested)
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

func (it *FeeQuoterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterOwnershipTransferRequestedIterator{contract: _FeeQuoter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterOwnershipTransferRequested)
				if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseOwnershipTransferRequested(log types.Log) (*FeeQuoterOwnershipTransferRequested, error) {
	event := new(FeeQuoterOwnershipTransferRequested)
	if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterOwnershipTransferredIterator struct {
	Event *FeeQuoterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterOwnershipTransferred)
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
		it.Event = new(FeeQuoterOwnershipTransferred)
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

func (it *FeeQuoterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterOwnershipTransferredIterator{contract: _FeeQuoter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterOwnershipTransferred)
				if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseOwnershipTransferred(log types.Log) (*FeeQuoterOwnershipTransferred, error) {
	event := new(FeeQuoterOwnershipTransferred)
	if err := _FeeQuoter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator struct {
	Event *FeeQuoterPremiumMultiplierWeiPerEthUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
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
		it.Event = new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
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

func (it *FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterPremiumMultiplierWeiPerEthUpdated struct {
	Token                      common.Address
	PremiumMultiplierWeiPerEth uint64
	Raw                        types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator{contract: _FeeQuoter.contract, event: "PremiumMultiplierWeiPerEthUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "PremiumMultiplierWeiPerEthUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*FeeQuoterPremiumMultiplierWeiPerEthUpdated, error) {
	event := new(FeeQuoterPremiumMultiplierWeiPerEthUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "PremiumMultiplierWeiPerEthUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterPriceFeedPerTokenUpdatedIterator struct {
	Event *FeeQuoterPriceFeedPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterPriceFeedPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterPriceFeedPerTokenUpdated)
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
		it.Event = new(FeeQuoterPriceFeedPerTokenUpdated)
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

func (it *FeeQuoterPriceFeedPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterPriceFeedPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterPriceFeedPerTokenUpdated struct {
	Token           common.Address
	PriceFeedConfig FeeQuoterTokenPriceFeedConfig
	Raw             types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPriceFeedPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterPriceFeedPerTokenUpdatedIterator{contract: _FeeQuoter.contract, event: "PriceFeedPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "PriceFeedPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterPriceFeedPerTokenUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParsePriceFeedPerTokenUpdated(log types.Log) (*FeeQuoterPriceFeedPerTokenUpdated, error) {
	event := new(FeeQuoterPriceFeedPerTokenUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "PriceFeedPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterReportPermissionSetIterator struct {
	Event *FeeQuoterReportPermissionSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterReportPermissionSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterReportPermissionSet)
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
		it.Event = new(FeeQuoterReportPermissionSet)
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

func (it *FeeQuoterReportPermissionSetIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterReportPermissionSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterReportPermissionSet struct {
	ReportId   [32]byte
	Permission KeystoneFeedsPermissionHandlerPermission
	Raw        types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterReportPermissionSet(opts *bind.FilterOpts, reportId [][32]byte) (*FeeQuoterReportPermissionSetIterator, error) {

	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "ReportPermissionSet", reportIdRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterReportPermissionSetIterator{contract: _FeeQuoter.contract, event: "ReportPermissionSet", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchReportPermissionSet(opts *bind.WatchOpts, sink chan<- *FeeQuoterReportPermissionSet, reportId [][32]byte) (event.Subscription, error) {

	var reportIdRule []interface{}
	for _, reportIdItem := range reportId {
		reportIdRule = append(reportIdRule, reportIdItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "ReportPermissionSet", reportIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterReportPermissionSet)
				if err := _FeeQuoter.contract.UnpackLog(event, "ReportPermissionSet", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseReportPermissionSet(log types.Log) (*FeeQuoterReportPermissionSet, error) {
	event := new(FeeQuoterReportPermissionSet)
	if err := _FeeQuoter.contract.UnpackLog(event, "ReportPermissionSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterTokenTransferFeeConfigDeletedIterator struct {
	Event *FeeQuoterTokenTransferFeeConfigDeleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterTokenTransferFeeConfigDeletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterTokenTransferFeeConfigDeleted)
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
		it.Event = new(FeeQuoterTokenTransferFeeConfigDeleted)
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

func (it *FeeQuoterTokenTransferFeeConfigDeletedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterTokenTransferFeeConfigDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterTokenTransferFeeConfigDeleted struct {
	DestChainSelector uint64
	Token             common.Address
	Raw               types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigDeletedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "TokenTransferFeeConfigDeleted", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterTokenTransferFeeConfigDeletedIterator{contract: _FeeQuoter.contract, event: "TokenTransferFeeConfigDeleted", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigDeleted, destChainSelector []uint64, token []common.Address) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "TokenTransferFeeConfigDeleted", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterTokenTransferFeeConfigDeleted)
				if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseTokenTransferFeeConfigDeleted(log types.Log) (*FeeQuoterTokenTransferFeeConfigDeleted, error) {
	event := new(FeeQuoterTokenTransferFeeConfigDeleted)
	if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterTokenTransferFeeConfigUpdatedIterator struct {
	Event *FeeQuoterTokenTransferFeeConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterTokenTransferFeeConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterTokenTransferFeeConfigUpdated)
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
		it.Event = new(FeeQuoterTokenTransferFeeConfigUpdated)
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

func (it *FeeQuoterTokenTransferFeeConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterTokenTransferFeeConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterTokenTransferFeeConfigUpdated struct {
	DestChainSelector      uint64
	Token                  common.Address
	TokenTransferFeeConfig FeeQuoterTokenTransferFeeConfig
	Raw                    types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterTokenTransferFeeConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigUpdatedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "TokenTransferFeeConfigUpdated", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterTokenTransferFeeConfigUpdatedIterator{contract: _FeeQuoter.contract, event: "TokenTransferFeeConfigUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchTokenTransferFeeConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigUpdated, destChainSelector []uint64, token []common.Address) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "TokenTransferFeeConfigUpdated", destChainSelectorRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterTokenTransferFeeConfigUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseTokenTransferFeeConfigUpdated(log types.Log) (*FeeQuoterTokenTransferFeeConfigUpdated, error) {
	event := new(FeeQuoterTokenTransferFeeConfigUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "TokenTransferFeeConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterUsdPerTokenUpdatedIterator struct {
	Event *FeeQuoterUsdPerTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterUsdPerTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterUsdPerTokenUpdated)
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
		it.Event = new(FeeQuoterUsdPerTokenUpdated)
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

func (it *FeeQuoterUsdPerTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterUsdPerTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterUsdPerTokenUpdated struct {
	Token     common.Address
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterUsdPerTokenUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterUsdPerTokenUpdatedIterator{contract: _FeeQuoter.contract, event: "UsdPerTokenUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerTokenUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "UsdPerTokenUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterUsdPerTokenUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseUsdPerTokenUpdated(log types.Log) (*FeeQuoterUsdPerTokenUpdated, error) {
	event := new(FeeQuoterUsdPerTokenUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FeeQuoterUsdPerUnitGasUpdatedIterator struct {
	Event *FeeQuoterUsdPerUnitGasUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FeeQuoterUsdPerUnitGasUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeQuoterUsdPerUnitGasUpdated)
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
		it.Event = new(FeeQuoterUsdPerUnitGasUpdated)
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

func (it *FeeQuoterUsdPerUnitGasUpdatedIterator) Error() error {
	return it.fail
}

func (it *FeeQuoterUsdPerUnitGasUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FeeQuoterUsdPerUnitGasUpdated struct {
	DestChain uint64
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log
}

func (_FeeQuoter *FeeQuoterFilterer) FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*FeeQuoterUsdPerUnitGasUpdatedIterator, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _FeeQuoter.contract.FilterLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return &FeeQuoterUsdPerUnitGasUpdatedIterator{contract: _FeeQuoter.contract, event: "UsdPerUnitGasUpdated", logs: logs, sub: sub}, nil
}

func (_FeeQuoter *FeeQuoterFilterer) WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}

	logs, sub, err := _FeeQuoter.contract.WatchLogs(opts, "UsdPerUnitGasUpdated", destChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FeeQuoterUsdPerUnitGasUpdated)
				if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
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

func (_FeeQuoter *FeeQuoterFilterer) ParseUsdPerUnitGasUpdated(log types.Log) (*FeeQuoterUsdPerUnitGasUpdated, error) {
	event := new(FeeQuoterUsdPerUnitGasUpdated)
	if err := _FeeQuoter.contract.UnpackLog(event, "UsdPerUnitGasUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetTokenAndGasPrices struct {
	TokenPrice    *big.Int
	GasPriceValue *big.Int
}
type ProcessMessageArgs struct {
	MsgFeeJuels           *big.Int
	IsOutOfOrderExecution bool
	ConvertedExtraArgs    []byte
	DestExecDataPerToken  [][]byte
}

func (_FeeQuoter *FeeQuoter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FeeQuoter.abi.Events["AuthorizedCallerAdded"].ID:
		return _FeeQuoter.ParseAuthorizedCallerAdded(log)
	case _FeeQuoter.abi.Events["AuthorizedCallerRemoved"].ID:
		return _FeeQuoter.ParseAuthorizedCallerRemoved(log)
	case _FeeQuoter.abi.Events["DestChainAdded"].ID:
		return _FeeQuoter.ParseDestChainAdded(log)
	case _FeeQuoter.abi.Events["DestChainConfigUpdated"].ID:
		return _FeeQuoter.ParseDestChainConfigUpdated(log)
	case _FeeQuoter.abi.Events["FeeTokenAdded"].ID:
		return _FeeQuoter.ParseFeeTokenAdded(log)
	case _FeeQuoter.abi.Events["FeeTokenRemoved"].ID:
		return _FeeQuoter.ParseFeeTokenRemoved(log)
	case _FeeQuoter.abi.Events["OwnershipTransferRequested"].ID:
		return _FeeQuoter.ParseOwnershipTransferRequested(log)
	case _FeeQuoter.abi.Events["OwnershipTransferred"].ID:
		return _FeeQuoter.ParseOwnershipTransferred(log)
	case _FeeQuoter.abi.Events["PremiumMultiplierWeiPerEthUpdated"].ID:
		return _FeeQuoter.ParsePremiumMultiplierWeiPerEthUpdated(log)
	case _FeeQuoter.abi.Events["PriceFeedPerTokenUpdated"].ID:
		return _FeeQuoter.ParsePriceFeedPerTokenUpdated(log)
	case _FeeQuoter.abi.Events["ReportPermissionSet"].ID:
		return _FeeQuoter.ParseReportPermissionSet(log)
	case _FeeQuoter.abi.Events["TokenTransferFeeConfigDeleted"].ID:
		return _FeeQuoter.ParseTokenTransferFeeConfigDeleted(log)
	case _FeeQuoter.abi.Events["TokenTransferFeeConfigUpdated"].ID:
		return _FeeQuoter.ParseTokenTransferFeeConfigUpdated(log)
	case _FeeQuoter.abi.Events["UsdPerTokenUpdated"].ID:
		return _FeeQuoter.ParseUsdPerTokenUpdated(log)
	case _FeeQuoter.abi.Events["UsdPerUnitGasUpdated"].ID:
		return _FeeQuoter.ParseUsdPerUnitGasUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FeeQuoterAuthorizedCallerAdded) Topic() common.Hash {
	return common.HexToHash("0xeb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef")
}

func (FeeQuoterAuthorizedCallerRemoved) Topic() common.Hash {
	return common.HexToHash("0xc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda77580")
}

func (FeeQuoterDestChainAdded) Topic() common.Hash {
	return common.HexToHash("0x525e3d4e0c31cef19cf9426af8d2c0ddd2d576359ca26bed92aac5fadda46265")
}

func (FeeQuoterDestChainConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x283b699f411baff8f1c29fe49f32a828c8151596244b8e7e4c164edd6569a835")
}

func (FeeQuoterFeeTokenAdded) Topic() common.Hash {
	return common.HexToHash("0xdf1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23")
}

func (FeeQuoterFeeTokenRemoved) Topic() common.Hash {
	return common.HexToHash("0x1795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f91")
}

func (FeeQuoterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FeeQuoterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FeeQuoterPremiumMultiplierWeiPerEthUpdated) Topic() common.Hash {
	return common.HexToHash("0xbb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d")
}

func (FeeQuoterPriceFeedPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0xe6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bf")
}

func (FeeQuoterReportPermissionSet) Topic() common.Hash {
	return common.HexToHash("0x32a4ba3fa3351b11ad555d4c8ec70a744e8705607077a946807030d64b6ab1a3")
}

func (FeeQuoterTokenTransferFeeConfigDeleted) Topic() common.Hash {
	return common.HexToHash("0x4de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b")
}

func (FeeQuoterTokenTransferFeeConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x94967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b5")
}

func (FeeQuoterUsdPerTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0x52f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a")
}

func (FeeQuoterUsdPerUnitGasUpdated) Topic() common.Hash {
	return common.HexToHash("0xdd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6e")
}

func (_FeeQuoter *FeeQuoter) Address() common.Address {
	return _FeeQuoter.address
}

type FeeQuoterInterface interface {
	FEEBASEDECIMALS(opts *bind.CallOpts) (*big.Int, error)

	KEYSTONEPRICEDECIMALS(opts *bind.CallOpts) (*big.Int, error)

	ConvertTokenAmount(opts *bind.CallOpts, fromToken common.Address, fromTokenAmount *big.Int, toToken common.Address) (*big.Int, error)

	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (FeeQuoterDestChainConfig, error)

	GetDestinationChainGasPrice(opts *bind.CallOpts, destChainSelector uint64) (InternalTimestampedPackedUint224, error)

	GetFeeTokens(opts *bind.CallOpts) ([]common.Address, error)

	GetPremiumMultiplierWeiPerEth(opts *bind.CallOpts, token common.Address) (uint64, error)

	GetStaticConfig(opts *bind.CallOpts) (FeeQuoterStaticConfig, error)

	GetTokenAndGasPrices(opts *bind.CallOpts, token common.Address, destChainSelector uint64) (GetTokenAndGasPrices,

		error)

	GetTokenPrice(opts *bind.CallOpts, token common.Address) (InternalTimestampedPackedUint224, error)

	GetTokenPriceFeedConfig(opts *bind.CallOpts, token common.Address) (FeeQuoterTokenPriceFeedConfig, error)

	GetTokenPrices(opts *bind.CallOpts, tokens []common.Address) ([]InternalTimestampedPackedUint224, error)

	GetTokenTransferFeeConfig(opts *bind.CallOpts, destChainSelector uint64, token common.Address) (FeeQuoterTokenTransferFeeConfig, error)

	GetValidatedFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetValidatedTokenPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	ProcessMessageArgs(opts *bind.CallOpts, destChainSelector uint64, feeToken common.Address, feeTokenAmount *big.Int, extraArgs []byte, onRampTokenTransfers []InternalEVM2AnyTokenTransfer, sourceTokenAmounts []ClientEVMTokenAmount) (ProcessMessageArgs,

		error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error)

	ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []FeeQuoterDestChainConfigArgs) (*types.Transaction, error)

	ApplyFeeTokensUpdates(opts *bind.TransactOpts, feeTokensToRemove []common.Address, feeTokensToAdd []common.Address) (*types.Transaction, error)

	ApplyPremiumMultiplierWeiPerEthUpdates(opts *bind.TransactOpts, premiumMultiplierWeiPerEthArgs []FeeQuoterPremiumMultiplierWeiPerEthArgs) (*types.Transaction, error)

	ApplyTokenTransferFeeConfigUpdates(opts *bind.TransactOpts, tokenTransferFeeConfigArgs []FeeQuoterTokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs []FeeQuoterTokenTransferFeeConfigRemoveArgs) (*types.Transaction, error)

	OnReport(opts *bind.TransactOpts, metadata []byte, report []byte) (*types.Transaction, error)

	SetReportPermissions(opts *bind.TransactOpts, permissions []KeystoneFeedsPermissionHandlerPermission) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdatePrices(opts *bind.TransactOpts, priceUpdates InternalPriceUpdates) (*types.Transaction, error)

	UpdateTokenPriceFeeds(opts *bind.TransactOpts, tokenPriceFeedUpdates []FeeQuoterTokenPriceFeedUpdate) (*types.Transaction, error)

	FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerAddedIterator, error)

	WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerAdded) (event.Subscription, error)

	ParseAuthorizedCallerAdded(log types.Log) (*FeeQuoterAuthorizedCallerAdded, error)

	FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*FeeQuoterAuthorizedCallerRemovedIterator, error)

	WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterAuthorizedCallerRemoved) (event.Subscription, error)

	ParseAuthorizedCallerRemoved(log types.Log) (*FeeQuoterAuthorizedCallerRemoved, error)

	FilterDestChainAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainAddedIterator, error)

	WatchDestChainAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainAdded(log types.Log) (*FeeQuoterDestChainAdded, error)

	FilterDestChainConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64) (*FeeQuoterDestChainConfigUpdatedIterator, error)

	WatchDestChainConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterDestChainConfigUpdated, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigUpdated(log types.Log) (*FeeQuoterDestChainConfigUpdated, error)

	FilterFeeTokenAdded(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenAddedIterator, error)

	WatchFeeTokenAdded(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenAdded, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenAdded(log types.Log) (*FeeQuoterFeeTokenAdded, error)

	FilterFeeTokenRemoved(opts *bind.FilterOpts, feeToken []common.Address) (*FeeQuoterFeeTokenRemovedIterator, error)

	WatchFeeTokenRemoved(opts *bind.WatchOpts, sink chan<- *FeeQuoterFeeTokenRemoved, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenRemoved(log types.Log) (*FeeQuoterFeeTokenRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FeeQuoterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeQuoterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeQuoterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FeeQuoterOwnershipTransferred, error)

	FilterPremiumMultiplierWeiPerEthUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPremiumMultiplierWeiPerEthUpdatedIterator, error)

	WatchPremiumMultiplierWeiPerEthUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPremiumMultiplierWeiPerEthUpdated, token []common.Address) (event.Subscription, error)

	ParsePremiumMultiplierWeiPerEthUpdated(log types.Log) (*FeeQuoterPremiumMultiplierWeiPerEthUpdated, error)

	FilterPriceFeedPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterPriceFeedPerTokenUpdatedIterator, error)

	WatchPriceFeedPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterPriceFeedPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParsePriceFeedPerTokenUpdated(log types.Log) (*FeeQuoterPriceFeedPerTokenUpdated, error)

	FilterReportPermissionSet(opts *bind.FilterOpts, reportId [][32]byte) (*FeeQuoterReportPermissionSetIterator, error)

	WatchReportPermissionSet(opts *bind.WatchOpts, sink chan<- *FeeQuoterReportPermissionSet, reportId [][32]byte) (event.Subscription, error)

	ParseReportPermissionSet(log types.Log) (*FeeQuoterReportPermissionSet, error)

	FilterTokenTransferFeeConfigDeleted(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigDeletedIterator, error)

	WatchTokenTransferFeeConfigDeleted(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigDeleted, destChainSelector []uint64, token []common.Address) (event.Subscription, error)

	ParseTokenTransferFeeConfigDeleted(log types.Log) (*FeeQuoterTokenTransferFeeConfigDeleted, error)

	FilterTokenTransferFeeConfigUpdated(opts *bind.FilterOpts, destChainSelector []uint64, token []common.Address) (*FeeQuoterTokenTransferFeeConfigUpdatedIterator, error)

	WatchTokenTransferFeeConfigUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterTokenTransferFeeConfigUpdated, destChainSelector []uint64, token []common.Address) (event.Subscription, error)

	ParseTokenTransferFeeConfigUpdated(log types.Log) (*FeeQuoterTokenTransferFeeConfigUpdated, error)

	FilterUsdPerTokenUpdated(opts *bind.FilterOpts, token []common.Address) (*FeeQuoterUsdPerTokenUpdatedIterator, error)

	WatchUsdPerTokenUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerTokenUpdated, token []common.Address) (event.Subscription, error)

	ParseUsdPerTokenUpdated(log types.Log) (*FeeQuoterUsdPerTokenUpdated, error)

	FilterUsdPerUnitGasUpdated(opts *bind.FilterOpts, destChain []uint64) (*FeeQuoterUsdPerUnitGasUpdatedIterator, error)

	WatchUsdPerUnitGasUpdated(opts *bind.WatchOpts, sink chan<- *FeeQuoterUsdPerUnitGasUpdated, destChain []uint64) (event.Subscription, error)

	ParseUsdPerUnitGasUpdated(log types.Log) (*FeeQuoterUsdPerUnitGasUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var FeeQuoterZKBin string = ("0x0003000000000002002e0000000000020002000000010355000000600310027000000a090030019d00000a090330019700000001002001900000009b0000c13d0000008002000039000000400020043f000000040030008c000000bc0000413d000000000401043b000000e00440027000000a500040009c000000be0000213d00000a660040009c0000020d0000a13d00000a670040009c000002310000213d00000a6d0040009c0000050a0000213d00000a700040009c000007190000613d00000a710040009c000000bc0000c13d000000c40030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b002600000002001d00000a0f0020009c000000bc0000213d0000002402100370000000000202043b002500000002001d00000a0e0020009c000000bc0000213d0000006402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000bc0000813d0000000404200039000000000441034f000000000404043b002400000004001d00000a0f0040009c000000bc0000213d0000002404200039002100000004001d002000240040002d000000200030006b000000bc0000213d0000008402100370000000000202043b002300000002001d00000a0f0020009c000000bc0000213d00000023020000290000002302200039000000000032004b000000bc0000813d00000023020000290000000402200039000000000221034f000000000202043b001f00000002001d00000a0f0020009c000000bc0000213d000000230200002900000024042000390000001f020000290000000502200210002200000004001d001e00000002001d0000000002420019000000000032004b000000bc0000213d000000a402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000bc0000813d0000000404200039000000000441034f000000000404043b001d00000004001d00000a0f0040009c000000bc0000213d001c00240020003d0000001d0200002900000006022002100000001c02200029000000000032004b000000bc0000213d0000004401100370000000000101043b001800000001001d00000a8b0100004100000000001004430000000001000412000000040010044300000020010000390000002400100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f000000010020019000001a910000613d000000000301043b000000250130014f00000a0e00100198000000810000613d00000025010000290000001802000029282020620000040f001800000001001d00000a8b010000410000000000100443000000000100041200000004001004430000002400000443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f000000010020019000001a910000613d000000000101043b00000a0d01100197000000180010006b000013100000a13d000000400200043d0000002403200039000000000013043500000ab201000041000000000012043500000004012000390000001803000029000013a20000013d000000e004000039000000400040043f0000000002000416000000000002004b000000bc0000c13d0000001f0230003900000a0a02200197000000e002200039000000400020043f0000001f0530018f00000a0b06300198000000e002600039000000ad0000613d000000000701034f000000007807043c0000000004840436000000000024004b000000a90000c13d000000000005004b000000ba0000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000001200030008c000000dc0000813d0000000001000019000028220001043000000a510040009c000002200000a13d00000a520040009c0000037b0000213d00000a580040009c000005190000213d00000a5b0040009c000007bf0000613d00000a5c0040009c000000bc0000c13d0000000001000416000000000001004b000000bc0000c13d0000000b02000039000000000102041a000000800010043f000000000020043f000000000001004b0000021e0000613d000000a00400003900000aa10200004100000000030000190000000005040019000000000402041a000000000445043600000001022000390000000103300039000000000013004b000000d40000413d00000ab60000013d000000400100043d002600000001001d00000a0c0010009c000000e60000a13d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000026010000290000006001100039000000400010043f000000e00100043d00000a0d0010009c000000bc0000213d00000026020000290000000001120436002500000001001d000001000100043d00000a0e0010009c000000bc0000213d00000025020000290000000000120435000001200100043d00000a090010009c000000bc0000213d00000026020000290000004002200039002400000002001d0000000000120435000001400200043d00000a0f0020009c000000bc0000213d000000e001300039000000ff04200039000000000014004b000000bc0000813d000000e004200039000000000404043300000a0f0040009c000000e00000213d00000005054002100000003f0650003900000a1006600197000000400700043d0000000006670019002300000007001d000000000076004b0000000007000039000000010700403900000a0f0060009c000000e00000213d0000000100700190000000e00000c13d000000400060043f0000002306000029000000000046043500000100022000390000000005250019000000000015004b000000bc0000213d000000000004004b000001240000613d0000002304000029000000002602043400000a0e0060009c000000bc0000213d00000020044000390000000000640435000000000052004b0000011d0000413d000001600200043d00000a0f0020009c000000bc0000213d0000001f04200039000000000034004b000000000500001900000a110500804100000a1104400197000000000004004b000000000600001900000a110600404100000a110040009c000000000605c019000000000006004b000000bc0000c13d000000e004200039000000000404043300000a0f0040009c000000e00000213d00000005054002100000003f0650003900000a1006600197000000400700043d0000000006670019002200000007001d000000000076004b0000000007000039000000010700403900000a0f0060009c000000e00000213d0000000100700190000000e00000c13d000000400060043f00000022060000290000000006460436001e00000006001d00000100022000390000000005250019000000000015004b000000bc0000213d000000000004004b000001560000613d0000002204000029000000002602043400000a0e0060009c000000bc0000213d00000020044000390000000000640435000000000052004b0000014f0000413d000001800200043d00000a0f0020009c000000bc0000213d0000001f04200039000000000034004b000000000500001900000a110500804100000a1104400197000000000004004b000000000600001900000a110600404100000a110040009c000000000605c019000000000006004b000000bc0000c13d000000e004200039000000000504043300000a0f0050009c000000e00000213d00000005045002100000003f0440003900000a1004400197000000400600043d0000000004460019001d00000006001d000000000064004b0000000006000039000000010600403900000a0f0040009c000000e00000213d0000000100600190000000e00000c13d000000400040043f0000001d040000290000000004540436001c00000004001d000001000220003900000007045002100000000004240019000000000014004b000000bc0000213d000000000005004b000012e40000c13d000001a00200043d00000a0f0020009c000000bc0000213d0000001f04200039000000000034004b000000000500001900000a110500804100000a1104400197000000000004004b000000000600001900000a110600404100000a110040009c000000000605c019000000000006004b000000bc0000c13d000000e004200039002100000004001d000000000704043300000a0f0070009c000000e00000213d00000005067002100000003f0460003900000a1005400197000000400400043d0000000005540019000000000045004b0000000008000039000000010800403900000a0f0050009c000000e00000213d0000000100800190000000e00000c13d000000400050043f002e00000004001d00000000007404350000010005200039002000000056001d000000200010006b000000bc0000213d000000000007004b000014190000c13d000001c00200043d00000a0f0020009c000000bc0000213d0000001f04200039000000000034004b000000000500001900000a110500404100000a1104400197000000000004004b000000000600001900000a110600204100000a110040009c000000000605c019000000000006004b000000bc0000613d000000e004200039000000000504043300000a0f0050009c000000e00000213d00000005045002100000003f0440003900000a1004400197000000400600043d0000000004460019001a00000006001d000000000064004b0000000006000039000000010600403900000a0f0040009c000000e00000213d0000000100600190000000e00000c13d000000400040043f0000001a04000029002d00000004001d0000000000540435000001000220003900000006045002100000000004240019000000000014004b000000bc0000213d000000000005004b000016f30000c13d000001e00200043d00000a0f0020009c000000bc0000213d0000001f04200039000000000034004b000000000300001900000a110300404100000a1104400197000000000004004b000000000500001900000a110500204100000a110040009c000000000503c019000000000005004b000000bc0000613d000000e003200039000000000403043300000a0f0040009c000000e00000213d00000005034002100000003f0330003900000a1003300197000000400500043d0000000003350019001900000005001d000000000053004b0000000005000039000000010500403900000a0f0030009c000000e00000213d0000000100500190000000e00000c13d000000400030043f00000019030000290000000003430436001800000003001d000001000220003900000240034000c90000000003230019000000000013004b000000bc0000213d000000000004004b000017c80000c13d000000400100043d001f00000001001d0000000001000411000000000001004b00001a920000c13d00000a4f010000410000001f02000029000000000012043500000a090020009c00000a0902008041000000400120021000000a4e011001c7000028220001043000000a720040009c000003960000a13d00000a730040009c0000047d0000213d00000a760040009c000005f00000613d00000a770040009c000000bc0000c13d0000000001000416000000000001004b000000bc0000c13d0000000202000039000000000102041a000000800010043f000000000020043f000000000001004b00000aac0000c13d000000200200003900000ab70000013d00000a5d0040009c000003dd0000a13d00000a5e0040009c000004a70000213d00000a610040009c000006040000613d00000a620040009c000000bc0000c13d0000000001000416000000000001004b000000bc0000c13d0000000101000039000000000101041a00000a0e01100197000000800010043f00000aa001000041000028210001042e00000a680040009c000005240000213d00000a6b0040009c000007e80000613d00000a6c0040009c000000bc0000c13d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000bc0000813d0000000404200039000000000441034f000000000504043b00000a0f0050009c000000e00000213d00000005045002100000003f0440003900000a100440019700000a7d0040009c000000e00000213d0000008004400039000000400040043f000000800050043f000000240220003900000240045000c90000000004240019000000000034004b000000bc0000213d000000000005004b00000c800000c13d0000000101000039000000000101041a00000a0e011001970000000002000411000000000012004b00000adf0000c13d000000800100043d000000000001004b00000aaa0000613d002600000000001d00000026010000290000000501100210000000a0011000390000000001010433000000001201043400000a0f0420019800000f090000613d00000000030104330000016001300039002500000001001d000000000101043300000a090110019800000f090000613d0000020002300039002400000002001d000000000202043300000a270220019700000a280020009c00000f090000c13d0000006002300039002200000002001d000000000202043300000a0902200197000000000021004b00000f090000213d002100000003001d000000000040043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c70000801002000039002300000004001d2820281b0000040f0000000100200190000000bc0000613d0000002306000029000000000101043b0000000101100039000000000101041a00000021050000290000000042050434000000000002004b0000000003000039000000010300c039000000400200043d0000000003320436002000000004001d00000000040404330000ffff0440018f00000000004304350000004003500039001f00000003001d000000000303043300000a0903300197000000400420003900000000003404350000002203000029000000000303043300000a0903300197000000600420003900000000003404350000008003500039001e00000003001d000000000303043300000a090330019700000080042000390000000000340435000000a003500039001d00000003001d00000000030304330000ffff0330018f000000a0042000390000000000340435000000c003500039001c00000003001d000000000303043300000a0903300197000000c0042000390000000000340435000000e003500039001b00000003001d00000000030304330000ffff0330018f000000e00420003900000000003404350000010003500039001a00000003001d00000000030304330000ffff0330018f000001000420003900000000003404350000012003500039001800000003001d00000000030304330000ffff0330018f000001200420003900000000003404350000014003500039001600000003001d000000000303043300000a0903300197000001400420003900000000003404350000002503000029000000000303043300000a0903300197000001600420003900000000003404350000018003500039001900000003001d000000000303043300000a0f0330019700000180042000390000000000340435000001a003500039001700000003001d000000000303043300000a0903300197000001a0042000390000000000340435000001c003500039001500000003001d000000000303043300000a0903300197000001c0042000390000000000340435000001e003500039001400000003001d0000000003030433000000000003004b0000000003000039000000010300c039000001e00420003900000000003404350000002403000029000000000303043300000a27033001970000020004200039000000000034043500000a090020009c00000a0902008041000000400220021000000a2900100198000002ff0000613d000000000100041400000a090010009c00000a0901008041000000c001100210000000000121019f00000aae011001c70000800d02000039000000020300003900000a2a04000041000003080000013d000000000100041400000a090010009c00000a0901008041000000c001100210000000000121019f00000aae011001c70000800d02000039000000020300003900000a2b040000410000000005060019282028160000040f0000000100200190000000bc0000613d0000002301000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d00000021020000290000000002020433000000000002004b000000000101043b000000000201041a00000a2c02200197000000010220c1bf00000020030000290000000003030433000000080330021000000a2d03300197000000000232019f0000001f030000290000000003030433000000180330021000000a2e03300197000000000232019f00000022030000290000000003030433000000380330021000000a2f03300197000000000232019f0000001e030000290000000003030433000000580330021000000a3003300197000000000232019f0000001d030000290000000003030433000000780330021000000a3103300197000000000232019f0000001c030000290000000003030433000000880330021000000a3203300197000000000232019f0000001b030000290000000003030433000000a80330021000000a3303300197000000000232019f0000001a030000290000000003030433000000b80330021000000a3403300197000000000232019f00000018030000290000000003030433000000c80330021000000a3503300197000000000232019f00000016030000290000000003030433000000d80330021000000a3603300197000000000232019f000000000021041b00000001011000390000002502000029000000000202043300000a0902200197000000000301041a00000a3703300197000000000223019f00000019030000290000000003030433000000200330021000000a3803300197000000000232019f00000017030000290000000003030433000000600330021000000a3903300197000000000232019f00000015030000290000000003030433000000800330021000000a3a03300197000000000232019f00000014030000290000000003030433000000000003004b00000a3b030000410000000003006019000000000232019f00000024030000290000000003030433000000380330027000000a2903300197000000000232019f000000000021041b0000002602000029002600010020003d000000800100043d000000260010006b000002610000413d00000aaa0000013d00000a530040009c000005c10000213d00000a560040009c000008080000613d00000a570040009c000000bc0000c13d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000601043b00000a0e0060009c000000bc0000213d0000000101000039000000000101041a00000a0e011001970000000005000411000000000015004b00000acb0000c13d000000000056004b00000acf0000c13d00000a8101000041000000800010043f00000a7f01000041000028220001043000000a780040009c00000a770000613d00000a790040009c000009ac0000613d00000a7a0040009c000000bc0000c13d0000000001000416000000000001004b000000bc0000c13d282023a40000040f000000400100043d002600000001001d28201f760000040f0000000001000412002c00000001001d002b00000000003d0000800501000039000000440300003900000000040004150000002c0440008a000000050440021000000a8b02000041282027f80000040f00000a0d02100197002400000002001d00000026010000290000000001210436002500000001001d0000000001000412002a00000001001d002900200000003d00000000040004150000002a0440008a0000000504400210000080050100003900000a8b020000410000004403000039282027f80000040f00000a0e01100197000000250200002900000000001204350000000001000412002800000001001d002700400000003d0000000004000415000000280440008a0000000504400210000080050100003900000a8b020000410000004403000039282027f80000040f0000002602000029000000400220003900000a09011001970000000000120435000000400100043d000000240300002900000000033104360000002504000029000000000404043300000a0e044001970000000000430435000000000202043300000a09022001970000004003100039000000000023043500000a090010009c00000a0901008041000000400110021000000abc011001c7000028210001042e00000a630040009c00000a8f0000613d00000a640040009c000009c00000613d00000a650040009c000000bc0000c13d000000440030008c000000bc0000413d0000000004000416000000000004004b000000bc0000c13d0000000404100370000000000604043b00000a0f0060009c000000bc0000213d0000002304600039000000000034004b000000bc0000813d0000000405600039000000000451034f000000000404043b00000a0f0040009c000000bc0000213d00000000064600190000002406600039000000000036004b000000bc0000213d0000002406100370000000000606043b00000a0f0060009c000000bc0000213d0000002307600039000000000037004b000000bc0000813d002400040060003d0000002407100360000000000707043b002500000007001d00000a0f0070009c000000bc0000213d0000002406600039002300000006001d002600250060002d000000260030006b000000bc0000213d0000001f0340003900000abd033001970000003f0330003900000abd0330019700000a7d0030009c000000e00000213d0000008003300039000000400030043f0000002003500039000000000331034f000000800040043f00000abd054001980000001f0640018f000000a0015000390000041f0000613d000000a007000039000000000803034f000000008908043c0000000007970436000000000017004b0000041b0000c13d000000000006004b0000042c0000613d000000000353034f0000000305600210000000000601043300000000065601cf000000000656022f000000000303043b0000010005500089000000000353022f00000000035301cf000000000363019f0000000000310435000000a0014000390000000000010435000000de0100043d002200000001001d00000aa303100197000000400100043d0000008004100039000000ca0500043d000000c00600043d0000000000340435002100000006001d00000a220360019700000060041000390000000000340435000000000221043600000060045002700000004003100039002000000004001d0000000000430435000000000300041100000a0e03300197000000000032043500000aa40010009c000000e00000213d000000a003100039000000400030043f00000a090020009c00000a09020080410000004002200210000000000101043300000a090010009c00000a09010080410000006001100210000000000121019f000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a20011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000010043f0000000401000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000ff00100190000010c40000c13d000000400200043d002600000002001d00000aa801000041000000000012043500000004012000390000000002000411000000200300002900000021040000290000002205000029282027800000040f0000002602000029000000000121004900000a090010009c00000a0901008041000000600110021000000a090020009c00000a09020080410000004002200210000000000121019f000028220001043000000a740040009c000006450000613d00000a750040009c000000bc0000c13d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000101043b002000000001001d00000a0f0010009c000000bc0000213d000000200130006a00000a120010009c000000bc0000213d000000440010008c000000bc0000413d0000000001000411000000000010043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000000001004b00000ae70000c13d000000400100043d00000ab802000041000000000021043500000004021000390000000003000411000005ea0000013d00000a5f0040009c0000064c0000613d00000a600040009c000000bc0000c13d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000bc0000813d0000000404200039000000000441034f000000000504043b00000a0f0050009c000000e00000213d00000005045002100000003f0440003900000a100440019700000a7d0040009c000000e00000213d0000008004400039000000400040043f000000800050043f000000240220003900000006045002100000000004240019000000000034004b000000bc0000213d000000000005004b00000bf00000c13d0000000101000039000000000101041a00000a0e011001970000000002000411000000000012004b00000adf0000c13d000000800100043d000000000001004b00000aaa0000613d002600000000001d00000026010000290000000501100210000000a001100039000000000101043300000020021000390000000002020433002400000002001d000000000101043300000a0e01100197002500000001001d000000000010043f0000000801000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000240200002900000a0f02200197000000000101043b000000000301041a00000a3e03300197000000000323019f000000000031041b000000400100043d000000000021043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a1b011001c70000800d02000039000000020300003900000a3f040000410000002505000029282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d000000800100043d000000260010006b000004d50000413d00000aaa0000013d00000a6e0040009c000008ae0000613d00000a6f0040009c000000bc0000c13d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000101043b00000a0e0010009c000000bc0000213d282025cf0000040f00000a880000013d00000a590040009c000008f20000613d00000a5a0040009c000000bc0000c13d0000000001000416000000000001004b000000bc0000c13d0000001201000039000000800010043f00000aa001000041000028210001042e00000a690040009c000009050000613d00000a6a0040009c000000bc0000c13d000000440030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b002500000002001d00000a0f0020009c000000bc0000213d00000025020000290000002302200039000000000032004b000000bc0000813d00000025020000290000000402200039000000000221034f000000000202043b00000a0f0020009c000000e00000213d00000005052002100000003f0450003900000a100440019700000a7d0040009c000000e00000213d0000008004400039000000400040043f000000800020043f00000025040000290000002404400039002400000045001d000000240030006b000000bc0000213d000000000002004b00000d160000c13d0000002402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000000500001900000a110500404100000a1104400197000000000004004b000000000600001900000a110600204100000a110040009c000000000605c019000000000006004b000000bc0000613d0000000404200039000000000441034f000000000504043b00000a0f0050009c000000e00000213d00000005045002100000003f0440003900000a1004400197000000400600043d0000000004460019001f00000006001d000000000064004b0000000006000039000000010600403900000a0f0040009c000000e00000213d0000000100600190000000e00000c13d000000400040043f0000001f040000290000000004540436001e00000004001d000000240220003900000006045002100000000004240019000000000034004b000000bc0000213d000000000005004b000010410000c13d0000000101000039000000000101041a00000a0e011001970000000002000411000000000012004b00000adf0000c13d000000800100043d000000000001004b000011df0000c13d0000001f010000290000000001010433000000000001004b00000aaa0000613d002600000000001d000000260100002900000005011002100000001e01100029000000000101043300000020021000390000000002020433002500000002001d000000000101043300000a0f01100197002400000001001d000000000010043f0000000a01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000250200002900000a0e02200197000000000101043b002500000002001d000000000020043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000001041b000000000100041400000a090010009c00000a0901008041000000c00110021000000a20011001c70000800d02000039000000030300003900000a4b0400004100000024050000290000002506000029282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d0000001f010000290000000001010433000000260010006b000005860000413d00000aaa0000013d00000a540040009c000009380000613d00000a550040009c000000bc0000c13d000000440030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b002600000002001d00000a0e0020009c000000bc0000213d0000002401100370000000000101043b002500000001001d00000a0f0010009c000000bc0000213d0000002501000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000ff0010019000000bd60000c13d000000400100043d00000a7c02000041000000000021043500000004021000390000002503000029000000000032043500000a090010009c00000a0901008041000000400110021000000a3d011001c700002822000104300000000001000416000000000001004b000000bc0000c13d000000c001000039000000400010043f0000001301000039000000800010043f00000aba01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e00200003928201fb50000040f000000c00110008a00000a090010009c00000a0901008041000000600110021000000abb011001c7000028210001042e000000440030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b00000a0f0020009c000000bc0000213d0000002401100370000000000101043b002600000001001d00000a0e0010009c000000bc0000213d0000014001000039000000400010043f000000800000043f000000a00000043f000000c00000043f000000e00000043f000001000000043f000001200000043f000000000020043f0000000a01000039000000200010043f00000040020000390000000001000019282027e30000040f0000002602000029282023940000040f002600000001001d000000400100043d002500000001001d28201f8c0000040f0000002601000029000000000101041a00000a88001001980000000002000039000000010200c0390000002504000029000000a0034000390000000000230435000000700210027000000a090220019700000080034000390000000000230435000000500210027000000a09022001970000006003400039000000000023043500000040021002700000ffff0220018f00000040034000390000000000230435000000200210027000000a09022001970000002003400039000000000023043500000a090110019700000000001404350000000002040019000000400100043d002600000001001d282020350000040f00000abd0000013d0000000001000416000000000001004b000000bc0000c13d0000002401000039000000800010043f00000aa001000041000028210001042e000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b00000a0f0020009c000000bc0000213d000000000423004900000a120040009c000000bc0000213d000000440040008c000000bc0000413d000000c005000039000000400050043f0000000404200039000000000641034f000000000606043b00000a0f0060009c000000bc0000213d00000000062600190000002307600039000000000037004b000000bc0000813d0000000407600039000000000771034f000000000807043b00000a0f0080009c000000e00000213d00000005078002100000003f0970003900000a100990019700000a140090009c000000e00000213d000000c009900039000000400090043f000000c00080043f00000024066000390000000007670019000000000037004b000000bc0000213d000000000008004b000006810000613d000000000861034f000000000808043b00000a0e0080009c000000bc0000213d000000200550003900000000008504350000002006600039000000000076004b000006780000413d000000c005000039000000800050043f0000002004400039000000000441034f000000000404043b00000a0f0040009c000000bc0000213d00000000022400190000002304200039000000000034004b000000000500001900000a110500804100000a1104400197000000000004004b000000000600001900000a110600404100000a110040009c000000000605c019000000000006004b000000bc0000c13d0000000404200039000000000441034f000000000404043b00000a0f0040009c000000e00000213d00000005054002100000003f0650003900000a1006600197000000400700043d0000000006670019002300000007001d000000000076004b0000000007000039000000010700403900000a0f0060009c000000e00000213d0000000100700190000000e00000c13d000000400060043f00000023060000290000000004460436002200000004001d00000024022000390000000004250019000000000034004b000000bc0000213d000000000042004b000006bb0000813d0000002303000029000000000521034f000000000505043b00000a0e0050009c000000bc0000213d000000200330003900000000005304350000002002200039000000000042004b000006b20000413d0000002301000029000000a00010043f0000000101000039000000000101041a00000a0e011001970000000002000411000000000012004b00000adf0000c13d00000023010000290000000001010433000000000001004b0000127d0000c13d000000800100043d002300000001001d0000000021010434002400000002001d000000000001004b00000aaa0000613d002600000000001d000000260100002900000005011002100000002401100029000000000101043300000a0e0110019800001b800000613d002500000001001d000000000010043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000000001004b000007000000c13d0000000201000039000000000101041a00000a0f0010009c000000e00000213d00000001021000390000000203000039000000000023041b00000a1d0110009a0000002502000029000000000021041b000000000103041a002200000001001d000000000020043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002202000029000000000021041b000000400100043d0000002502000029000000000021043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a1b011001c70000800d02000039000000010300003900000a1e04000041282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d00000023010000290000000001010433000000260010006b000006ce0000413d00000aaa0000013d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000bc0000813d0000000404200039000000000441034f000000000504043b00000a0f0050009c000000e00000213d00000005045002100000003f0440003900000a100440019700000a7d0040009c000000e00000213d0000008004400039000000400040043f000000800050043f0000002402200039000000a0045000c90000000004240019000000000034004b000000bc0000213d000000000005004b00000c0b0000c13d0000000101000039000000000101041a00000a0e011001970000000002000411000000000012004b00000adf0000c13d000000800100043d000000000001004b00000aaa0000613d002600000000001d00000026010000290000000501100210000000a00110003900000000040104330000004001400039002300000001001d000000000101043300000aa3021001970000000053040434002500000004001d0000006001400039002400000001001d0000000004010433002200000005001d0000000005050433000000400100043d0000008006100039000000000026043500000a22025001970000006005100039000000000025043500000a0e024001970000004004100039000000000024043500000a0e03300197000000200210003900000000003204350000008003000039000000000031043500000aa40010009c000000e00000213d000000a003100039000000400030043f00000a090020009c00000a09020080410000004002200210000000000101043300000a090010009c00000a09010080410000006001100210000000000121019f000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a20011001c700008010020000392820281b0000040f00000001002001900000002503000029000000bc0000613d000000000101043b0000008002300039002000000002001d0000000002020433001f00000002001d002100000001001d000000000010043f0000000401000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000025040000290000000100200190000000bc0000613d000000000101043b000000000201041a00000abe022001970000001f0000006b000000010220c1bf000000000021041b000000000104043300000a0e01100197000000400200043d00000000011204360000002203000029000000000303043300000a220330019700000000003104350000002301000029000000000101043300000aa301100197000000400320003900000000001304350000002401000029000000000101043300000a0e011001970000006003200039000000000013043500000020010000290000000001010433000000000001004b0000000001000039000000010100c0390000008003200039000000000013043500000a090020009c00000a09020080410000004001200210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000ab5011001c70000800d02000039000000020300003900000ab6040000410000002105000029282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d000000800100043d000000260010006b000007430000413d00000aaa0000013d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000101043b002600000001001d00000a0e0010009c000000bc0000213d282023a40000040f0000002601000029000000000010043f0000000701000039000000200010043f00000040020000390000000001000019282027e30000040f002600000001001d000000400100043d002500000001001d28201f760000040f0000002601000029000000000101041a00000a8d001001980000000002000039000000010200c039000000250400002900000040034000390000000000230435000000a002100270000000ff0220018f0000002003400039000000000023043500000a0e0110019700000000001404350000000002040019000000400100043d002600000001001d282020530000040f00000abd0000013d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000101043b002600000001001d00000a0f0010009c000000bc0000213d282023d30000040f0000002601000029000000000010043f0000000501000039000000200010043f00000040020000390000000001000019282027e30000040f002600000001001d000000400100043d002500000001001d28201f810000040f0000002601000029000000000101041a00000025040000290000002002400039000000e003100270000000000032043500000a160110019700000000001404350000000001040019000008fc0000013d000000440030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b002600000002001d00000a0f0020009c000000bc0000213d0000002401100370000000000101043b002500000001001d00000a0f0010009c000000bc0000213d000000250130006a00000a120010009c000000bc0000213d000000a40010008c000000bc0000413d0000002601000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400200043d00000a150020009c000000e00000213d000000000101043b0000022003200039000000400030043f000000000301041a000000d80430027000000a09044001970000014005200039002400000005001d0000000000450435000000c8043002700000ffff0440018f0000012005200039002300000005001d0000000000450435000000b8043002700000ffff0440018f0000010005200039001e00000005001d0000000000450435000000a8043002700000ffff0440018f000000e005200039001500000005001d0000000000450435000000880430027000000a0904400197000000c005200039001600000005001d000000000045043500000078043002700000ffff0440018f000000a005200039001800000005001d0000000000450435000000580430027000000a09044001970000008005200039001700000005001d0000000000450435000000380430027000000a09044001970000006005200039001b00000005001d0000000000450435000000180430027000000a09044001970000004005200039002200000005001d000000000045043500000008043002700000ffff0440018f0000002005200039002100000005001d0000000000450435000000ff033001900000000004000039000000010400c03900000000004204350000000101100039000000000101041a000001600520003900000a0904100197001900000005001d0000000000450435000000380410021000000a27044001970000020005200039002000000005001d000000000045043500000a23001001980000000004000039000000010400c039000001e005200039001a00000005001d0000000000450435000000800410027000000a0904400197000001c005200039001f00000005001d0000000000450435000000600410027000000a0904400197000001a005200039001d00000005001d00000000004504350000018002200039000000200110027000000a0f01100197001c00000002001d0000000000120435000000000003004b00000da80000613d00000025010000290000006401100039001400000001001d0000000201100367000000000101043b00000a0e0010009c000000bc0000213d000000000010043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d0000000202000367000000000101043b000000000101041a000000000001004b00000f0f0000c13d0000001401200360000000000101043b00000a0e0010009c000000bc0000213d000000400200043d00000a9f0300004100000000003204350000000403200039000000000013043500000a090020009c00000a0902008041000000400120021000000a3d011001c70000282200010430000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000bc0000813d0000000404200039000000000141034f000000000101043b002000000001001d00000a0f0010009c000000bc0000213d001f00240020003d000000200100002900000005021002100000001f01200029000000000031004b000000bc0000213d0000003f0120003900000a100310019700000a7d0030009c000000e00000213d0000008001300039000000400010043f0000002004000029000000800040043f000000000004004b00000c3e0000c13d00000020020000390000000003210436000000800200043d00000000002304350000004004100039000000000002004b000008e90000613d000000800300003900000000050000190000000006010019000000000704001900000020033000390000000004030433000000008404043400000a160440019700000000004704350000006004600039000000000608043300000a0906600197000000000064043500000040047000390000000105500039000000000025004b0000000006070019000008da0000413d000000000214004900000a090020009c00000a0902008041000000600220021000000a090010009c00000a09010080410000004001100210000000000112019f000028210001042e000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000101043b00000a0e0010009c000000bc0000213d282024360000040f000000400200043d002600000002001d28201fd70000040f000000260100002900000a090010009c00000a0901008041000000400110021000000a7b011001c7000028210001042e000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000101043b002600000001001d00000a0f0010009c000000bc0000213d28201fa90000040f00000200021000390000000000020435000001e0021000390000000000020435000001c0021000390000000000020435000001a00210003900000000000204350000018002100039000000000002043500000160021000390000000000020435000001400210003900000000000204350000012002100039000000000002043500000100021000390000000000020435000000e0021000390000000000020435000000c0021000390000000000020435000000a0021000390000000000020435000000800210003900000000000204350000006002100039000000000002043500000040021000390000000000020435000000000101043600000000000104350000002601000029282023b50000040f282023e20000040f0000000002010019000000400100043d002600000001001d28201fde0000040f00000abd0000013d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000bc0000813d0000000404200039000000000441034f000000000504043b00000a0f0050009c000000e00000213d00000005045002100000003f0440003900000a100440019700000a7d0040009c000000e00000213d0000008004400039000000400040043f000000800050043f000000240220003900000007045002100000000004240019000000000034004b000000bc0000213d000000000005004b00000c4f0000c13d0000000101000039000000000101041a00000a0e011001970000000002000411000000000012004b00000adf0000c13d000000800100043d000000000001004b00000aaa0000613d002600000000001d00000026010000290000000501100210000000a001100039000000000101043300000020021000390000000002020433002400000002001d000000000101043300000a0e01100197002500000001001d000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000201041a00000a2202200197000000240600002900000020036000390000000004030433000000a00440021000000a2304400197000000000242019f00000040046000390000000005040433000000000005004b00000a24050000410000000005006019000000000252019f000000000506043300000a0e05500197000000000252019f000000000021041b000000400100043d00000000025104360000000003030433000000ff0330018f00000000003204350000000002040433000000000002004b0000000002000039000000010200c0390000004003100039000000000023043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a25011001c70000800d02000039000000020300003900000a26040000410000002505000029282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d000000800100043d000000260010006b000009620000413d00000aaa0000013d000000240030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000401100370000000000101043b00000a0e0010009c000000bc0000213d000000000010043f0000000801000039000000200010043f00000040020000390000000001000019282027e30000040f000000000101041a00000a0f01100197000000800010043f00000aa001000041000028210001042e000000440030008c000000bc0000413d0000000004000416000000000004004b000000bc0000c13d0000000404100370000000000404043b00000a0f0040009c000000bc0000213d0000002305400039000000000035004b000000bc0000813d0000000405400039000000000551034f000000000605043b00000a0f0060009c000000e00000213d00000005056002100000003f0750003900000a100770019700000a7d0070009c000000e00000213d0000008007700039000000400070043f000000800060043f00000024044000390000000005450019000000000035004b000000bc0000213d000000000006004b000009e80000613d000000000641034f000000000606043b00000a0e0060009c000000bc0000213d000000200220003900000000006204350000002004400039000000000054004b000009df0000413d0000002402100370000000000202043b00000a0f0020009c000000bc0000213d0000002304200039000000000034004b000000000500001900000a110500804100000a1104400197000000000004004b000000000600001900000a110600404100000a110040009c000000000605c019000000000006004b000000bc0000c13d0000000404200039000000000441034f000000000404043b00000a0f0040009c000000e00000213d00000005054002100000003f0650003900000a1006600197000000400700043d0000000006670019002300000007001d000000000076004b0000000007000039000000010700403900000a0f0060009c000000e00000213d0000000100700190000000e00000c13d000000400060043f00000023060000290000000006460436002200000006001d00000024022000390000000005250019000000000035004b000000bc0000213d000000000004004b00000a1e0000613d0000002303000029000000000421034f000000000404043b00000a0e0040009c000000bc0000213d000000200330003900000000004304350000002002200039000000000052004b00000a150000413d0000000101000039000000000101041a00000a0e011001970000000002000411000000000012004b00000adf0000c13d000000800100043d000000000001004b0000105f0000c13d00000023010000290000000001010433000000000001004b00000aaa0000613d002600000000001d00000a330000013d0000002602000029002600010020003d00000023010000290000000001010433000000260010006b00000aaa0000813d000000260100002900000005011002100000002201100029002400000001001d000000000101043300000a0e01100197002500000001001d000000000010043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000000001004b00000a2d0000c13d0000000b03000039000000000103041a00000a0f0010009c000000e00000213d0000000102100039000000000023041b00000a1f0110009a0000002502000029000000000021041b000000000103041a002100000001001d000000000020043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002102000029000000000021041b00000023010000290000000001010433000000260010006c00001e580000a13d00000024010000290000000001010433000000000200041400000a0e0510019700000a090020009c00000a0902008041000000c00120021000000a20011001c70000800d02000039000000020300003900000a2104000041282028160000040f000000010020019000000a2d0000c13d000000bc0000013d000000640030008c000000bc0000413d0000000002000416000000000002004b000000bc0000c13d0000000402100370000000000402043b00000a0e0040009c000000bc0000213d0000004402100370000000000302043b00000a0e0030009c000000bc0000213d0000002401100370000000000201043b0000000001040019282020620000040f000000400200043d000000000012043500000a090020009c00000a0902008041000000400120021000000a9a011001c7000028210001042e0000000001000416000000000001004b000000bc0000c13d000000000100041a00000a0e021001970000000006000411000000000026004b00000ac70000c13d0000000102000039000000000302041a00000a1704300197000000000464019f000000000042041b00000a1701100197000000000010041b000000000100041400000a0e0530019700000a090010009c00000a0901008041000000c00110021000000a20011001c70000800d02000039000000030300003900000aad04000041282028160000040f0000000100200190000000bc0000613d0000000001000019000028210001042e000000a00400003900000ab90200004100000000030000190000000005040019000000000402041a000000000445043600000001022000390000000103300039000000000013004b00000aaf0000413d000000600250008a000000800100003928201f970000040f000000400100043d002600000001001d000000800200003928201fc70000040f0000002602000029000000000121004900000a090010009c00000a0901008041000000600110021000000a090020009c00000a09020080410000004002200210000000000121019f000028210001042e00000aac01000041000000800010043f00000a7f01000041000028220001043000000a7e01000041000000800010043f00000a7f010000410000282200010430000000000100041a00000a1701100197000000000161019f000000000010041b000000000100041400000a090010009c00000a0901008041000000c00110021000000a20011001c70000800d02000039000000030300003900000a8004000041282028160000040f0000000100200190000000bc0000613d00000aaa0000013d000000400100043d00000a7e02000041000000000021043500000a090010009c00000a0901008041000000400110021000000a4e011001c700002822000104300000002004000029002600040040003d00000002020003670000002601200360000000000301043b00000000010000310000000004410049000000230440008a00000a110540019700000a1106300197000000000756013f000000000056004b000000000500001900000a1105004041000000000043004b000000000400001900000a110400804100000a110070009c000000000504c019000000000005004b000000bc0000c13d0000002603300029000000000232034f000000000202043b001d00000002001d00000a0f0020009c000000bc0000213d0000001d02000029000000060220021000000000012100490000002002300039000000000012004b000000000300001900000a110300204100000a110110019700000a1102200197000000000412013f000000000012004b000000000100001900000a110100404100000a110040009c000000000103c019000000000001004b000000bc0000c13d00000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f000000010020019000001a910000613d000000000101043b001e00000001001d001f0a090010019b0000001d0000006b00000fbc0000c13d0000002601000029001d00200010003d00000002020003670000001d01200360000000000301043b0000000001000031000000200410006a000000230440008a00000a110540019700000a1106300197000000000756013f000000000056004b000000000500001900000a1105004041000000000043004b000000000400001900000a110400804100000a110070009c000000000504c019000000000005004b000000bc0000c13d0000002603300029000000000232034f000000000202043b001c00000002001d00000a0f0020009c000000bc0000213d0000001c02000029000000060220021000000000012100490000002002300039000000000012004b000000000300001900000a110300204100000a110110019700000a1102200197000000000412013f000000000012004b000000000100001900000a110100404100000a110040009c000000000103c019000000000001004b000000bc0000c13d0000001c0000006b00000aaa0000613d000000000900001900000002010003670000001d02100360000000000302043b0000000002000031000000200420006a000000230440008a00000a110540019700000a1106300197000000000756013f000000000056004b000000000500001900000a1105004041000000000043004b000000000400001900000a110400804100000a110070009c000000000504c019000000000005004b000000bc0000c13d0000002604300029000000000341034f000000000303043b00000a0f0030009c000000bc0000213d00000006053002100000000005520049000000200440003900000a110650019700000a1107400197000000000867013f000000000067004b000000000600001900000a1106004041000000000054004b000000000500001900000a110500204100000a110080009c000000000605c019000000000006004b000000bc0000c13d000000000039004b00001e580000813d00000006039002100000000003340019000000000232004900000a120020009c000000bc0000213d000000400020008c000000bc0000413d000000400400043d00000a130040009c000000e00000213d0000004002400039000000400020043f000000000231034f000000000202043b00000a0f0020009c000000bc0000213d00000000052404360000002002300039000000000121034f000000000101043b00000a160010009c000000bc0000213d002500000005001d0000000000150435000000400300043d00000a130030009c000000e00000213d0000004002300039000000400020043f002200000003001d00000000021304360000001f01000029002100000002001d0000000000120435000000000104043300000a0f01100197000000000010043f0000000501000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c70000801002000039002300000009001d002400000004001d2820281b0000040f000000240400002900000001002001900000002505000029000000bc0000613d0000002202000029000000000202043300000a160220019700000021030000290000000003030433000000e003300210000000000223019f000000000101043b000000000021041b00000000020404330000000001050433000000400300043d00000020043000390000001e05000029000000000054043500000a1601100197000000000013043500000a090030009c00000a09030080410000004001300210000000000300041400000a090030009c00000a0903008041000000c003300210000000000113019f00000a19011001c700000a0f052001970000800d02000039000000020300003900000ab704000041282028160000040f00000023090000290000000100200190000000bc0000613d00000001099000390000001c0090006c00000b520000413d00000aaa0000013d0000002601000029282025cf0000040f0000002502000029000000000020043f0000000902000039000000200020043f002600000001001d00000040020000390000000001000019282027e30000040f0000000101100039000000000101041a000000800110027000000a090210019700000025010000292820278d0000040f000000400200043d000000200320003900000000001304350000002601000029000000000012043500000a090020009c00000a0902008041000000400120021000000a7b011001c7000028210001042e000000a005000039000000000623004900000a120060009c000000bc0000213d000000400060008c000000bc0000413d000000400600043d00000a130060009c000000e00000213d0000004007600039000000400070043f000000000721034f000000000707043b00000a0e0070009c000000bc0000213d00000000077604360000002008200039000000000881034f000000000808043b00000a0f0080009c000000bc0000213d000000000087043500000000056504360000004002200039000000000042004b00000bf10000413d000004cb0000013d000000a005000039000000000623004900000a120060009c000000bc0000213d000000a00060008c000000bc0000413d000000400600043d00000aa40060009c000000e00000213d000000a007600039000000400070043f000000000721034f000000000707043b00000a0e0070009c000000bc0000213d00000000087604360000002007200039000000000971034f000000000909043b00000ab300900198000000bc0000c13d00000000009804350000002007700039000000000871034f000000000808043b00000ab400800198000000bc0000c13d000000400960003900000000008904350000002007700039000000000871034f000000000808043b00000a0e0080009c000000bc0000213d000000600960003900000000008904350000002007700039000000000771034f000000000707043b000000000007004b0000000008000039000000010800c039000000000087004b000000bc0000c13d000000800860003900000000007804350000000005650436000000a002200039000000000042004b00000c0c0000413d000007390000013d00000a140030009c000000e00000213d00000000030000190000004004100039000000400040043f000000200410003900000000000404350000000000010435000000a00430003900000000001404350000002003300039000000000023004b00000dae0000813d000000400100043d00000a130010009c00000c410000a13d000000e00000013d000000a005000039000000000623004900000a120060009c000000bc0000213d000000800060008c000000bc0000413d000000400600043d00000a130060009c000000e00000213d0000004007600039000000400070043f000000000721034f000000000707043b00000a0e0070009c000000bc0000213d0000000007760436000000400800043d00000a0c0080009c000000e00000213d0000006009800039000000400090043f0000002009200039000000000a91034f000000000a0a043b00000a0e00a0009c000000bc0000213d000000000aa804360000002009900039000000000b91034f000000000b0b043b000000ff00b0008c000000bc0000213d0000000000ba04350000002009900039000000000991034f000000000909043b000000000009004b000000000a000039000000010a00c0390000000000a9004b000000bc0000c13d000000400a80003900000000009a0435000000000087043500000000056504360000008002200039000000000042004b00000c500000413d000009580000013d000000a005000039000000000623004900000a120060009c000000bc0000213d000002400060008c000000bc0000413d000000400600043d00000a130060009c000000e00000213d0000004007600039000000400070043f000000000721034f000000000707043b00000a0f0070009c000000bc0000213d0000000007760436000000400800043d00000a150080009c000000e00000213d0000022009800039000000400090043f0000002009200039000000000a91034f000000000a0a043b00000000000a004b000000000b000039000000010b00c0390000000000ba004b000000bc0000c13d000000000aa804360000002009900039000000000b91034f000000000b0b043b0000ffff00b0008c000000bc0000213d0000000000ba04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000000400b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000000600b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000000800b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b0000ffff00a0008c000000bc0000213d000000a00b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000000c00b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b0000ffff00a0008c000000bc0000213d000000e00b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b0000ffff00a0008c000000bc0000213d000001000b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b0000ffff00a0008c000000bc0000213d000001200b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000001400b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000001600b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0f00a0009c000000bc0000213d000001800b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000001a00b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000a0900a0009c000000bc0000213d000001c00b8000390000000000ab04350000002009900039000000000a91034f000000000a0a043b00000000000a004b000000000b000039000000010b00c0390000000000ba004b000000bc0000c13d000001e00b8000390000000000ab04350000002009900039000000000991034f000000000909043b00000a1600900198000000bc0000c13d000002000a80003900000000009a0435000000000087043500000000056504360000024002200039000000000042004b00000c810000413d000002570000013d000000a006000039002300240030009200000d1f0000013d00000026020000290000000000a2043500000000068604360000002004400039000000240040006c0000054b0000813d000000000241034f000000000202043b00000a0f0020009c000000bc0000213d0000002502200029000000230520006900000a120050009c000000bc0000213d000000400050008c000000bc0000413d000000400800043d00000a130080009c000000e00000213d0000004005800039000000400050043f0000002405200039000000000751034f000000000707043b00000a0f0070009c000000bc0000213d0000000007780436002600000007001d0000002005500039000000000551034f000000000505043b00000a0f0050009c000000bc0000213d00000000022500190000004305200039000000000035004b000000000700001900000a110700804100000a1105500197000000000005004b000000000900001900000a110900404100000a110050009c000000000907c019000000000009004b000000bc0000c13d0000002405200039000000000551034f000000000c05043b00000a0f00c0009c000000e00000213d0000000505c002100000003f0550003900000a1005500197000000400a00043d000000000b5a00190000000000ab004b0000000005000039000000010500403900000a0f00b0009c000000e00000213d0000000100500190000000e00000c13d0000004000b0043f0000000000ca0435000000440b200039000000e002c000c9000000000cb2001900000000003c004b000000bc0000213d0000000000cb004b00000d190000813d000000000d0a00190000000002b3004900000a120020009c000000bc0000213d000000e00020008c000000bc0000413d000000400e00043d00000a1300e0009c000000e00000213d0000004002e00039000000400020043f0000000002b1034f000000000202043b00000a0e0020009c000000bc0000213d000000000f2e0436000000400200043d00000a140020009c000000e00000213d000000c005200039000000400050043f0000002005b00039000000000751034f000000000707043b00000a090070009c000000bc0000213d00000000077204360000002005500039000000000951034f000000000909043b00000a090090009c000000bc0000213d00000000009704350000002005500039000000000751034f000000000707043b0000ffff0070008c000000bc0000213d000000400920003900000000007904350000002005500039000000000751034f000000000707043b00000a090070009c000000bc0000213d000000600920003900000000007904350000002005500039000000000751034f000000000707043b00000a090070009c000000bc0000213d000000800920003900000000007904350000002005500039000000000551034f000000000505043b000000000005004b0000000007000039000000010700c039000000000075004b000000bc0000c13d000000200dd00039000000a007200039000000000057043500000000002f04350000000000ed0435000000e00bb000390000000000cb004b00000d620000413d00000d190000013d000000400100043d00000a7c02000041000000000021043500000004021000390000002603000029000005ea0000013d0000000002000019002300000002001d0000000502200210002200000002001d0000001f012000290000000201100367000000000101043b002600000001001d00000a0e0010009c000000bc0000213d000000400100043d00000a130010009c000000e00000213d0000004002100039000000400020043f0000002002100039000000000002043500000000000104350000002601000029000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400300043d00000a130030009c000000e00000213d000000000101043b0000004002300039000000400020043f000000000101041a00000a1602100197002500000003001d0000000002230436000000e001100270002100000002001d002400000001001d000000000012043500000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f000000010020019000001a910000613d000000000101043b00240024001000740000175c0000413d00000a8b0100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f000000010020019000001a910000613d000000000101043b00000a0901100197000000240010006b00000dff0000813d0000002302000029000000250500002900000ef70000013d0000002601000029000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000001002001900000002505000029000000bc0000613d000000400200043d00000a0c0020009c000000e00000213d000000000101043b0000006003200039000000400030043f000000000301041a000000400120003900000a8d003001980000000004000039000000010400c039000000000041043500000a0e013001970000000006120436000000a003300270000000ff0330018f000000000036043500000ef60000613d000000000001004b00000ef60000613d000000400100043d00000a130010009c000000e00000213d0000004003100039000000400030043f0000002003100039000000000003043500000000000104350000000002020433000000400b00043d00000a8e0100004100000000051b0436000000000100041400000a0e02200197000000040020008c002600000006001d00000e370000c13d0000000103000031000000a00030008c000000a004000039000000000403401900000e660000013d001d00000005001d00000a0900b0009c00000a090300004100000000030b4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c7001e00000002001d00240000000b001d2820281b0000040f000000240b000029000000600310027000000a0903300197000000a00030008c000000a0040000390000000004034019000000e00640019000000000056b001900000e530000613d000000000701034f00000000080b0019000000007907043c0000000008980436000000000058004b00000e4f0000c13d0000001f0740019000000e600000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f0000000100200190000013a40000613d00000026060000290000001e020000290000001d050000290000001f01400039000001e00110018f000000000ab1001900000000001a004b0000000001000039000000010100403900000a0f00a0009c000000e00000213d0000000100100190000000e00000c13d0000004000a0043f000000a00030008c000000bc0000413d00000000010b043300000a8f0010009c000000bc0000213d0000008001b00039000000000101043300000a8f0010009c000000bc0000213d000000000805043300000a110080009c000012760000813d0000006001b00039000000000701043300000a900100004100000000001a04350000000001000414000000040020008c000000200400003900000eb40000613d001d00000008001d001e00000007001d00000a0900a0009c00000a090300004100000000030a4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c700240000000a001d2820281b0000040f000000240a000029000000600310027000000a0903300197000000200030008c00000020040000390000000004034019000000200640019000000000056a001900000ea10000613d000000000701034f00000000080a0019000000007907043c0000000008980436000000000058004b00000e9d0000c13d0000001f0740019000000eae0000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f0000000100200190000013b00000613d00000026060000290000001e070000290000001d080000290000001f01400039000000600110018f0000000001a1001900000a0f0010009c000000e00000213d000000400010043f000000200030008c000000bc0000413d00000000020a0433000000ff0020008c000000bc0000213d0000000003060433000000ff0330018f0000000003230019000000ff0030008c0000175c0000213d000000240230008c00000ed80000213d000000010200003900000ed10000613d00000024033000890000000a04000039000000010030019000000000054400a9000000010400603900000000022400a90000000103300272000000000405001900000eca0000c13d000000000008004b00000f040000613d00000000038200a900000000048300d9000000000024004b00000ee60000613d0000175c0000013d0000004d0020008c0000175c0000213d00000001030000390000000a04000039000000010020019000000000054400a9000000010400603900000000033400a90000000102200272000000000405001900000edc0000c13d000000000003004b00001f160000613d00000000033800d900000a160030009c0000172c0000213d00000a130010009c0000002505000029000000e00000213d0000004002100039000000400020043f00000a09027001970000002004100039000000000024043500000000003104350000002103000029000000000303043300000a0903300197000000000032004b00000000050180190000002302000029000000800100043d000000000021004b00001e580000a13d0000002201000029000000a0011000390000000000510435000000800100043d000000000021004b00001e580000a13d0000000102200039000000200020006c00000daf0000413d0000105d0000013d000000000300001900000a130010009c000000250500002900000eeb0000a13d000000e00000013d000000400100043d00000a3c02000041000000000021043500000004021000390000000000420435000005eb0000013d00000025010000290000000405100039000000140100002900120020001000920000001201200360000000000401043b0000000003000031001300000005001d00000000015300490000001f0110008a00000a110510019700000a1106400197000000000756013f000000000056004b000000000600001900000a1106004041000000000014004b000000000800001900000a110800804100000a110070009c000000000608c019000000000006004b000000bc0000c13d0000001304400029000000000642034f000000000606043b001100000006001d00000a0f0060009c000000bc0000213d0000001106000029000000060660021000000000066300490000002004400039000000000064004b000000000700001900000a110700204100000a110660019700000a1104400197000000000864013f000000000064004b000000000400001900000a110400404100000a110080009c000000000407c019000000000004004b000000bc0000c13d0000001204000029000000200640008a000000000462034f000000000404043b00000a1107400197000000000857013f000000000057004b000000000700001900000a1107004041000000000014004b000000000900001900000a110900804100000a110080009c000000000709c019000000000007004b000000bc0000c13d0000001307400029000000000472034f000000000404043b00000a0f0040009c000000bc0000213d00000000084300490000002007700039000000000087004b000000000900001900000a110900204100000a110880019700000a1107700197000000000a87013f000000000087004b000000000700001900000a110700404100000a1100a0009c000000000709c019000000000007004b000000bc0000c13d000000200660008a000000000662034f000000000606043b00000a1107600197000000000857013f000000000057004b000000000500001900000a1105004041000000000016004b000000000100001900000a110100804100000a110080009c000000000501c019000000000005004b000000bc0000c13d0000001301600029000000000512034f000000000605043b00000a0f0060009c000000bc0000213d0000000005630049000000200710003900000a110150019700000a1108700197000000000918013f000000000018004b000000000100001900000a1101004041000000000057004b000000000500001900000a110500204100000a110090009c000000000105c019000000000001004b000000bc0000c13d0000001f0160003900000abd011001970000003f0110003900000abd05100197000000400100043d0000000005510019000000000015004b0000000008000039000000010800403900000a0f0050009c000000e00000213d0000000100800190000000e00000c13d000000400050043f00000000056104360000000008760019000000000038004b000000bc0000213d000000000772034f00000abd086001980000001f0960018f000000000385001900000fa10000613d000000000a07034f000000000b05001900000000ac0a043c000000000bcb043600000000003b004b00000f9d0000c13d000000000009004b00000fae0000613d000000000787034f0000000308900210000000000903043300000000098901cf000000000989022f000000000707043b0000010008800089000000000787022f00000000078701cf000000000797019f0000000000730435000000000365001900000000000304350000002203000029000000000303043300000a0903300197000000000043004b000015e40000813d000000400100043d0000002402100039000000000042043500000a9e0200004100000000002104350000000402100039000015f00000013d000000000900001900000002010003670000002602100360000000000302043b0000000002000031000000200420006a000000230440008a00000a110540019700000a1106300197000000000756013f000000000056004b000000000500001900000a1105004041000000000043004b000000000400001900000a110400804100000a110070009c000000000504c019000000000005004b000000bc0000c13d0000002604300029000000000341034f000000000303043b00000a0f0030009c000000bc0000213d00000006053002100000000005520049000000200440003900000a110650019700000a1107400197000000000867013f000000000067004b000000000600001900000a1106004041000000000054004b000000000500001900000a110500204100000a110080009c000000000605c019000000000006004b000000bc0000c13d000000000039004b00001e580000813d00000006039002100000000003340019000000000232004900000a120020009c000000bc0000213d000000400020008c000000bc0000413d000000400400043d00000a130040009c000000e00000213d0000004002400039000000400020043f000000000231034f000000000202043b00000a0e0020009c000000bc0000213d00000000052404360000002002300039000000000121034f000000000101043b00000a160010009c000000bc0000213d002500000005001d0000000000150435000000400300043d00000a130030009c000000e00000213d0000004002300039000000400020043f002200000003001d00000000021304360000001f01000029002100000002001d0000000000120435000000000104043300000a0e01100197000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c70000801002000039002300000009001d002400000004001d2820281b0000040f000000240400002900000001002001900000002505000029000000bc0000613d0000002202000029000000000202043300000a160220019700000021030000290000000003030433000000e003300210000000000223019f000000000101043b000000000021041b00000000020404330000000001050433000000400300043d00000020043000390000001e05000029000000000054043500000a1601100197000000000013043500000a090030009c00000a09030080410000004001300210000000000300041400000a090030009c00000a0903008041000000c003300210000000000113019f00000a19011001c700000a0e052001970000800d02000039000000020300003900000aa704000041282028160000040f00000023090000290000000100200190000000bc0000613d00000001099000390000001d0090006c00000fbd0000413d00000b230000013d0000001f05000029000000000623004900000a120060009c000000bc0000213d000000400060008c000000bc0000413d000000400600043d00000a130060009c000000e00000213d0000004007600039000000400070043f000000000721034f000000000707043b00000a0f0070009c000000bc0000213d00000000077604360000002008200039000000000881034f000000000808043b00000a0e0080009c000000bc0000213d0000002005500039000000000087043500000000006504350000004002200039000000000042004b000010420000413d000005780000013d000000400100043d000008d00000013d0000000002000019000010660000013d00000026020000290000000102200039000000800100043d000000000012004b00000a270000813d002600000002001d0000000501200210000000a001100039002400000001001d000000000101043300000a0e01100197002500000001001d000000000010043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000301041a000000000003004b000010610000613d0000000b01000039000000000201041a000000000002004b0000175c0000613d000000010130008a000000000023004b0000109f0000613d000000000012004b00001e580000a13d00000aa90130009a00000aa90220009a000000000202041a000000000021041b000000000020043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c70000801002000039002100000003001d2820281b0000040f00000021030000290000000100200190000000bc0000613d000000000101043b000000000031041b0000000b01000039000000000301041a000000000003004b00001b2c0000613d000000010130008a00000aa90230009a000000000002041b0000000b02000039000000000012041b0000002501000029000000000010043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000001041b000000800100043d000000260010006c00001e580000a13d00000024010000290000000001010433000000000200041400000a0e0510019700000a090020009c00000a0902008041000000c00120021000000a20011001c70000800d02000039000000020300003900000aab04000041282028160000040f0000000100200190000010610000c13d000000bc0000013d0000002501000029000000200010008c000000bc0000413d000000240100002900000020021000390000000201000367000000000221034f000000000202043b00000a0f0020009c000000bc0000213d00000023022000290000001f03200039000000260030006c000000bc0000813d000000000321034f000000000403043b00000a0f0040009c000000e00000213d00000005034002100000003f0330003900000a1003300197000000400500043d0000000003350019002500000005001d000000000053004b0000000005000039000000010500403900000a0f0030009c000000e00000213d0000000100500190000000e00000c13d000000400030043f00000025030000290000000003430436002100000003001d000000200220003900000060034000c90000000003230019000000260030006c000000bc0000213d000000000004004b00000aaa0000613d0000002504000029000000260520006900000a120050009c000000bc0000213d000000600050008c000000bc0000413d000000400500043d00000a0c0050009c000000e00000213d0000006006500039000000400060043f000000000621034f000000000606043b00000a0e0060009c000000bc0000213d00000000076504360000002006200039000000000861034f000000000808043b00000a160080009c000000bc0000213d00000000008704350000002006600039000000000661034f000000000606043b00000a090060009c000000bc0000213d00000020044000390000004007500039000000000067043500000000005404350000006002200039000000000032004b000010ef0000413d00000025010000290000000001010433000000000001004b00000aaa0000613d002600000000001d0000111c0000013d0000002602000029002600010020003d00000025010000290000000001010433000000260010006b00000aaa0000813d000000260100002900000005011002100000002101100029002400000001001d0000000001010433000000000101043300000a0e01100197000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400200043d00000a0c0020009c000000e00000213d000000000101043b0000006003200039000000400030043f000000000101041a000000400320003900000a8d001001980000000004000039000000010400c039000000000043043500000a0e031001970000000002320436000000a003100270000000ff0430018f0000000000420435000017210000613d00000025020000290000000002020433000000260020006c00001e580000a13d00000aa50110019700000aa60010009c0000175c0000213d000000240100002900000000010104330000002002100039000000000202043300000a16022001970000001203300039000000ff0430018f000000240340008c000011630000213d00000001030000390000115c0000613d00000024044000890000000a05000039000000010040019000000000065500a9000000010500603900000000033500a900000001044002720000000005060019000011550000c13d000000000002004b000011dd0000613d00000000042300a900000000022400d9000000000032004b000011710000613d0000175c0000013d0000004d0030008c0000175c0000213d00000001040000390000000a05000039000000010030019000000000065500a9000000010500603900000000044500a900000001033002720000000005060019000011670000c13d000000000004004b00001f160000613d00000000044200d9002200000004001d00000a160040009c0000172b0000213d00000040021000390000000002020433002300000002001d000000000101043300000a0e01100197000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000230200002900000a0902200197000000000101043b000000000101041a000000e001100270000000000012004b000011160000413d00000025010000290000000001010433000000260010006c00001e580000a13d000000400100043d002300000001001d00000a130010009c000000e00000213d000000240100002900000000010104330000004001100039000000000101043300000023030000290000004002300039000000400020043f0000002202000029000000000223043600000a0901100197002000000002001d000000000012043500000025010000290000000001010433000000260010006c00001e580000a13d00000024010000290000000001010433000000000101043300000a0e01100197000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d0000002302000029000000000202043300000a160220019700000020030000290000000003030433000000e003300210000000000223019f000000000101043b000000000021041b00000025010000290000000001010433000000260010006c00001e580000a13d0000002401000029000000000101043300000000020104330000004001100039000000000101043300000a0901100197000000400300043d000000200430003900000000001404350000002201000029000000000013043500000a090030009c00000a09030080410000004001300210000000000300041400000a090030009c00000a0903008041000000c003300210000000000113019f00000a0e0520019700000a19011001c70000800d02000039000000020300003900000aa704000041282028160000040f0000000100200190000011160000c13d000000bc0000013d002200000000001d000011740000013d0000000002000019000011e50000013d0000001d020000290000000102200039000000000012004b000005810000813d001d00000002001d0000000502200210000000a00220003900000000030204330000002002300039002000000002001d00000000020204330000000004020433000000000004004b000011e10000613d000000000103043300250a0f0010019b0000000004000019000000050140021000000000011200190000002001100039000000000201043300000020012000390000000006010433000000005106043400000a0901100197000000000305043300000a0903300197000000000031004b00001d950000813d002300000005001d002100000004001d000000000102043300260a0e0010019b002200000006001d0000008001600039002400000001001d000000000101043300000a09011001970000001f0010008c0000139b0000a13d0000002501000029000000000010043f0000000a01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000001002001900000002603000029000000bc0000613d000000000101043b000000000030043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000026060000290000000100200190000000bc0000613d00000023090000290000000002090433000000200220021000000a4202200197000000000101043b000000000301041a00000a4303300197000000000223019f000000220a0000290000004003a000390000000004030433000000400440021000000a4404400197000000000242019f0000006004a000390000000005040433000000500550021000000a4505500197000000000252019f00000024080000290000000005080433000000700550021000000a4605500197000000000252019f000000a005a000390000000007050433000000000007004b00000a47070000410000000007006019000000000272019f00000000070a043300000a0907700197000000000272019f000000000021041b000000400100043d0000000002710436000000000709043300000a0907700197000000000072043500000000020304330000ffff0220018f00000040031000390000000000230435000000000204043300000a090220019700000060031000390000000000230435000000000208043300000a0902200197000000800310003900000000002304350000000002050433000000000002004b0000000002000039000000010200c039000000a003100039000000000023043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a48011001c70000800d02000039000000030300003900000a49040000410000002505000029282028160000040f0000000100200190000000bc0000613d00000021040000290000000104400039000000200100002900000000020104330000000001020433000000000014004b000011f20000413d000000800100043d000011e10000013d00000a910100004100000000001a043500000a0900a0009c00000a090a0080410000004001a0021000000a4e011001c700002822000104300000000002000019000012850000013d0000002502000029000000010220003900000023010000290000000001010433000000000012004b000006c70000813d002500000002001d00000005012002100000002201100029000000000101043300000a0e01100197002600000001001d000000000010043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a002400000001001d000000000001004b0000127f0000613d0000000201000039000000000201041a000000000002004b0000175c0000613d0000002403000029000000010130008a000000000032004b000012be0000613d000000000012004b00001e580000a13d00000a1a0130009a00000a1a0220009a000000000202041a000000000021041b000000000020043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002402000029000000000021041b0000000201000039000000000301041a000000000003004b00001b2c0000613d000000010130008a00000a1a0230009a000000000002041b0000000202000039000000000012041b0000002601000029000000000010043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000001041b000000400100043d0000002602000029000000000021043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a1b011001c70000800d02000039000000010300003900000a1c04000041282028160000040f00000001002001900000127f0000c13d000000bc0000013d0000001c05000029000000000621004900000a120060009c000000bc0000213d000000800060008c000000bc0000413d000000400600043d00000a130060009c000000e00000213d0000004007600039000000400070043f000000009702043400000a0e0070009c000000bc0000213d0000000007760436000000400800043d00000a0c0080009c000000e00000213d000000600a8000390000004000a0043f000000000909043300000a0e0090009c000000bc0000213d0000000009980436000000400a200039000000000a0a0433000000ff00a0008c000000bc0000213d0000000000a9043500000060092000390000000009090433000000000009004b000000000a000039000000010a00c0390000000000a9004b000000bc0000c13d000000400a80003900000000009a0435000000000087043500000000056504360000008002200039000000000042004b000012e50000413d000001810000013d0000002601000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000400200043d00000a130020009c000000e00000213d0000000101100039000000000101041a0000004003200039000000400030043f000000200320003900000000000304350000000000020435000000240000006b000013ce0000c13d000000400200043d001700000002001d00000a130020009c000000e00000213d00000a090110019700000017030000290000004002300039000000400020043f0000000000130435001500000000001d000000000100001900000017020000290000002002200039001600000002001d00000000001204350000002601000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000001e020000290000003f0220003900000a1002200197000000400300043d0000000002230019001b00000003001d000000000032004b0000000003000039000000010300403900000a0f0020009c000000e00000213d0000000100300190000000e00000c13d0000000101100039000000000101041a000000400020043f0000001f020000290000001b030000290000000003230436001a00000003001d000000000002004b0000149c0000c13d000000400400043d000000200240003900000a97010000410000000000120435000000170100002900000000010104330000002403400039000000000013043500000016010000290000000001010433000000000001004b0000000001000039000000010100c039000000440340003900000000001304350000004401000039000000000014043500000a7d0040009c000000e00000213d0000008001400039000000400010043f000000c003400039000000800500003900000000005304350000001503000029000000010330018f000000a0054000390000000000350435000000180300002900000000003104350000010003400039000000000504043300000000005304350000012003400039000000000005004b000013890000613d000000000600001900000000073600190000000008260019000000000808043300000000008704350000002006600039000000000056004b000013820000413d000000000235001900000000000204350000001f0550003900000abd0550019700000000033500190000000005130049000000e00440003900000000005404350000001b0400002900000000040404330000000000430435000000050540021000000000055300190000002007500039000000000004004b000016c10000c13d0000000002170049000008ea0000013d000000400200043d0000002403200039000000000013043500000a4a01000041000000000012043500000004012000390000002603000029000000000031043500001d9c0000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000013bb0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000013ab0000c13d000013bb0000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000013bb0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000013b70000c13d000000000005004b000013c80000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f0000000000140435000000600130021000000a090020009c00000a09020080410000004002200210000000000112019f000028220001043000000024020000290000000301200210000000200110008900000a270310021f000000040020008c00000a270300804100000002050003670000002101500360000000000401043b000000040620008c000000bc0000413d00000024010000290000001b0110003900000abd011001970000003f0110003900000abd02100197000000400100043d0000000002210019000000000012004b0000000007000039000000010700403900000a0f0020009c000000e00000213d0000000100700190000000e00000c13d000000400020043f00000000026104360000002009000029000000000090007c000000bc0000213d00000021070000290000000407700039000000000775034f00000abd086001980000001f0660018f0000000005820019000013f90000613d000000000907034f000000000a020019000000009b09043c000000000aba043600000000005a004b000013f50000c13d000000000334016f000000000006004b000014070000613d000000000487034f0000000306600210000000000705043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f000000000045043500000024041000290000001c04400039000000000004043500000a970030009c000015f60000613d00000a980030009c00001ded0000c13d000000000101043300000a120010009c000000bc0000213d000000200010008c000000bc0000413d000000400100043d001700000001001d00000a130010009c000000e00000213d00000000010204330000132f0000013d001f00c00030003d000014200000013d00000020044000390000000000a904350000000000840435000000200050006c000001aa0000813d000000005205043400000a0f0020009c000000bc0000213d00000021022000290000001f0620006900000a120060009c000000bc0000213d000000400060008c000000bc0000413d000000400800043d00000a130080009c000000e00000213d0000004006800039000000400060043f0000002006200039000000000606043300000a0f0060009c000000bc0000213d00000000096804360000004006200039000000000606043300000a0f0060009c000000bc0000213d00000000022600190000003f06200039000000000016004b000000000700001900000a110700804100000a1106600197000000000006004b000000000a00001900000a110a00404100000a110060009c000000000a07c01900000000000a004b000000bc0000c13d0000002006200039000000000c06043300000a0f00c0009c000000e00000213d0000000506c002100000003f0660003900000a1006600197000000400a00043d00000000066a00190000000000a6004b0000000007000039000000010700403900000a0f0060009c000000e00000213d0000000100700190000000e00000c13d000000400060043f0000000000ca0435000000400b200039000000e002c000c9000000000cb2001900000000001c004b000000bc0000213d0000000000cb004b0000141b0000813d000000000d0a00190000000002b1004900000a120020009c000000bc0000213d000000e00020008c000000bc0000413d000000400e00043d00000a1300e0009c000000e00000213d0000004002e00039000000400020043f00000000620b043400000a0e0020009c000000bc0000213d000000000f2e0436000000400200043d00000a140020009c000000e00000213d000000c007200039000000400070043f000000000606043300000a090060009c000000bc0000213d00000000066204360000004007b00039000000000707043300000a090070009c000000bc0000213d00000000007604350000006006b0003900000000060604330000ffff0060008c000000bc0000213d000000400720003900000000006704350000008006b00039000000000606043300000a090060009c000000bc0000213d00000060072000390000000000670435000000a006b00039000000000606043300000a090060009c000000bc0000213d00000080072000390000000000670435000000c006b000390000000006060433000000000006004b0000000007000039000000010700c039000000000076004b000000bc0000c13d000000200dd00039000000a007200039000000000067043500000000002f04350000000000ed0435000000e00bb000390000000000cb004b0000145e0000413d0000141b0000013d000000600200003900000000030000190000001a050000290000001e06000029000000000453001900000000002404350000002003300039000000000063004b000014a00000413d00190a290010019b00000000040000190000001d0040006c00001e580000813d002400000004001d00000006014002100000001c011000290000000202000367000000000112034f000000000101043b002100000001001d00000a0e0010009c000000bc0000213d00000024010000290000000501100210002000000001001d001e00220010002d0000001e01200360000000000101043b0000000003000031000000220430006a0000009f0440008a00000a110540019700000a1106100197000000000756013f000000000056004b000000000500001900000a1105004041000000000041004b000000000600001900000a110600804100000a110070009c000000000506c019000000000005004b000000bc0000c13d00000022051000290000004006500039000000000662034f000000000606043b00000000075300490000001f0770008a00000a110870019700000a1109600197000000000a89013f000000000089004b000000000800001900000a1108004041000000000076004b000000000700001900000a110700804100000a1100a0009c000000000807c019000000000008004b000000bc0000c13d0000000005560019000000000652034f000000000606043b002500000006001d00000a0f0060009c000000bc0000213d000000250630006a0000002005500039000000000065004b000000000700001900000a110700204100000a110660019700000a1105500197000000000865013f000000000065004b000000000500001900000a110500404100000a110080009c000000000507c019000000000005004b000000bc0000c13d0000002505000029000000200050008c000015190000a13d0000002601000029000000000010043f0000000a01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002102000029000000000020043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000700110027000000a0901100197000000250010006b0000172e0000213d0000000003000031000000230130006a00000002020003670000001e05200360000000c30410008a000000000105043b000000000041004b000000000500001900000a110500804100000a110440019700000a1106100197000000000746013f000000000046004b000000000400001900000a110400404100000a110070009c000000000405c019000000000004004b000000bc0000c13d00000022011000290000002004100039000000000442034f000000000404043b00000000051300490000001f0550008a00000a110650019700000a1107400197000000000867013f000000000067004b000000000600001900000a1106004041000000000054004b000000000500001900000a110500804100000a110080009c000000000605c019000000000006004b000000bc0000c13d0000000001140019000000000412034f000000000504043b00000a0f0050009c000000bc0000213d0000000004530049000000200610003900000a110140019700000a1107600197000000000817013f000000000017004b000000000100001900000a1101004041000000000046004b000000000400001900000a110400204100000a110080009c000000000104c019000000000001004b000000bc0000c13d0000001f0150003900000abd011001970000003f0110003900000abd04100197000000400100043d0000000004410019000000000014004b0000000007000039000000010700403900000a0f0040009c000000e00000213d0000000100700190000000e00000c13d000000400040043f00000000045104360000000007650019000000000037004b000000bc0000213d000000000362034f00000abd065001980000000002640019000015690000613d000000000703034f0000000008040019000000007907043c0000000008980436000000000028004b000015650000c13d0000001f07500190000015760000613d000000000363034f0000000306700210000000000702043300000000076701cf000000000767022f000000000303043b0000010006600089000000000363022f00000000036301cf000000000373019f000000000032043500000000025400190000000000020435000000190200002900000ab00020009c000015820000c13d0000000002010433000000200020008c000016e00000c13d0000000002040433000004000220008a00000a830020009c000017340000813d0000002601000029000000000010043f0000000a01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002102000029000000000020043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400200043d00000a140020009c000000e00000213d000000000101043b000000c003200039000000400030043f000000000101041a000000a00320003900000a88001001980000000004000039000000010400c0390000000000430435000000700310027000000a09033001970000008004200039000000000034043500000040031002700000ffff0330018f00000040042000390000000000340435000000200310027000000a09033001970000002004200039000000000034043500000a090310019700000000003204350000006003200039000000500110027000000a09021001970000000000230435000015cc0000c13d0000002601000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000d80110027000000a0902100197000000400100043d000000200310003900000000002304350000002002000039000000000021043500000a130010009c000000e00000213d0000004002100039000000400020043f0000001b0300002900000000020304330000002404000029000000000042004b00001e580000a13d00000020050000290000001a0250002900000000001204350000000001030433000000000041004b00001e580000a13d00000001044000390000001f0040006c000014a70000413d0000135d0000013d000000210300002900000000030304330000ffff0330018f000000110030006c0000160d0000813d000000400100043d0000002402100039000000000032043500000a9d02000041000000000021043500000004021000390000001103000029000000000032043500000a090010009c00000a0901008041000000400110021000000a41011001c70000282200010430000000000301043300000a120030009c000000bc0000213d000000400030008c000000bc0000413d000000400300043d001700000003001d00000a130030009c000000e00000213d00000017040000290000004003400039000000400030043f0000000002020433000000000024043500000040011000390000000001010433000000000001004b0000000002000039000000010200c039001500000002001d000000000021004b000000bc0000c13d000013350000013d0000002003000029000000000303043300000a270330019700000a280030009c000016190000c13d0000000003010433000000200030008c000016e00000c13d0000000003050433000004000330008a00000a830030009c0000170c0000813d0000001401200360000000000101043b00000a0e0010009c000000bc0000213d282025cf0000040f0000001f020000290000000002020433002200000002001d0000002602000029000000000020043f0000000502000039000000200020043f000600000001001d000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400200043d002100000002001d00000a130020009c000000e00000213d000000220200002900200a090020019c000000000101043b00000021030000290000004002300039000000400020043f000000000101041a0000002002300039000000e004100270002200000004001d000000000042043500000a1601100197000100000001001d00000000001304350000174e0000c13d000000110000006b000018500000c13d0000001d01000029000000000101043300000a090110019700220a89001000d5001d00000000001d001f00000000001d0000000002000031000000130320006a0000001401000029000000400610008a0000000201000367000000000461034f0000001f0530008a000000000404043b0000001e0300002900000000030304330000ffff03300190002600000000001d000017620000c13d000000000054004b000000000300001900000a110300804100000a110750019700000a1108400197000000000978013f000000000078004b000000000800001900000a110800404100000a110090009c000000000803c019000000000008004b000000bc0000c13d000000170300002900000000030304330000001304400029000000000841034f000000000808043b00000a0f0080009c000000bc0000213d00000000098200490000002004400039000000000094004b000000000a00001900000a110a00204100000a110990019700000a1104400197000000000b94013f000000000094004b000000000400001900000a110400404100000a1100b0009c00000000040ac019000000000004004b000000bc0000c13d000000180400002900000000040404330000ffff0940018f00000000048900a9000000000008004b000016850000613d00000a960880019700000a960a40019700000000088a00d9000000000089004b0000175c0000c13d0000006006600039000000000661034f000000000606043b00000a1108600197000000000978013f000000000078004b000000000700001900000a1107004041000000000056004b000000000500001900000a110500804100000a110090009c000000000705c019000000000007004b000000bc0000c13d0000001306600029000000000561034f000000000505043b00000a0f0050009c000000bc0000213d0000000007520049000000200860003900000a110670019700000a1109800197000000000a69013f000000000069004b000000000600001900000a1106004041000000000078004b000000000700001900000a110700204100000a1100a0009c000000000607c019000000000006004b000000bc0000c13d000000400600043d00000a130060009c000000e00000213d0000004007600039000000400070043f000000200760003900000000000704350000000000060435000000400700043d00000a130070009c000000e00000213d000000190600002900000000060604330000004009700039000000400090043f000000200970003900000000000904350000000000070435000000000005004b00001da40000c13d000000400500043d00000a130050009c000000e00000213d00000a090260019700001de80000013d000000000500001900000000060300190000001b0d000029000016cd0000013d000000000987001900000000000904350000001f0880003900000abd0880019700000000078700190000000105500039000000000045004b000013990000813d0000000008370049000000200880008a00000020066000390000000000860435000000200dd0003900000000080d043300000000980804340000000007870436000000000008004b000016c50000613d000000000a000019000000000b7a0019000000000ca90019000000000c0c04330000000000cb0435000000200aa0003900000000008a004b000016d80000413d000016c50000013d000000400400043d002600000004001d00000a82020000410000000000240435000000040240003900000020030000390000000000320435000000240240003928201fb50000040f0000002602000029000000000121004900000a090010009c00000a090100804100000a090020009c00000a090200804100000060011002100000004002200210000000000121019f00002822000104300000001a05000029000000000621004900000a120060009c000000bc0000213d000000400060008c000000bc0000413d000000400600043d00000a130060009c000000e00000213d0000004007600039000000400070043f000000008702043400000a0e0070009c000000bc0000213d0000000007760436000000000808043300000a0f0080009c000000bc0000213d0000002005500039000000000087043500000000006504350000004002200039000000000042004b000016f40000413d000001d50000013d000000400200043d00000a820300004100000000003204350000000403200039000000200400003900000000004304350000000001010433000000240320003900000000001304350000004403200039000000000001004b000017480000613d000000000400001900000000063400190000000007540019000000000707043300000000007604350000002004400039000000000014004b000017190000413d000017480000013d00000025010000290000002602000029282023c60000040f00000000010104330000000001010433000000400200043d00000a9203000041000000000032043500000a0e01100197000008a70000013d000000400100043d00000a910200004100000ae10000013d000000400100043d00000ab102000041000000000021043500000004021000390000002103000029000005ea0000013d000000400200043d00000a820300004100000000003204350000000403200039000000200500003900000000005304350000000001010433000000240320003900000000001304350000004403200039000000000001004b000017480000613d000000000500001900000000063500190000000007450019000000000707043300000000007604350000002005500039000000000015004b000017410000413d0000001f0410003900000abd04400197000000000131001900000000000104350000004401400039000004750000013d00000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f000000010020019000001a910000613d000000000101043b000000220110006c000017b60000813d00000aaa01000041000000000010043f0000001101000039000000040010043f00000a3d010000410000282200010430000000000054004b000000000700001900000a110700804100000a110850019700000a1109400197000000000a89013f000000000089004b000000000800001900000a110800404100000a1100a0009c000000000807c019000000000008004b000000bc0000c13d0000001308400029000000000781034f000000000707043b00000a0f0070009c000000bc0000213d00000000097200490000002008800039000000000098004b000000000a00001900000a110a00204100000a110990019700000a1108800197000000000b98013f000000000098004b000000000800001900000a110800404100000a1100b0009c00000000080ac019000000000008004b000000bc0000c13d000000110900002900000120089000c9000000000009004b0000178a0000613d00000011098000fa000001200090008c0000175c0000c13d000001e007700039000000000087001a0000175c0000413d00000000078700190000001f0070002a0000175c0000413d0000001f0870002a000000150700002900000000070704330000ffff0970018f00000000078900a9000017990000613d00000000088700d9000000000098004b0000175c0000c13d0000001608000029000000000808043300000a0908800197000000000078001a0000175c0000413d000000000878001a002600000000001d000016570000613d0000000107000029000000700970027000000000079800a900000000088700d9000000000098004b0000175c0000c13d000000000007004b002600000000001d000016570000613d00000000083700a900000000077800d9000000000037004b0000175c0000c13d000000000008004b002600000000001d000016570000613d00260a95008000d500000026038000f900000a950030009c000016570000613d0000175c0000013d000000200010006c0000184b0000a13d000000400200043d0000004403200039000000000013043500000024012000390000002003000029000000000031043500000a8601000041000000000012043500000004012000390000002603000029000000000031043500000a090020009c00000a0902008041000000400120021000000a87011001c700002822000104300000001804000029000000000521004900000a120050009c000000bc0000213d000002400050008c000000bc0000413d000000400500043d00000a130050009c000000e00000213d0000004006500039000000400060043f000000008602043400000a0f0060009c000000bc0000213d0000000006650436000000400700043d00000a150070009c000000e00000213d0000022009700039000000400090043f0000000008080433000000000008004b0000000009000039000000010900c039000000000098004b000000bc0000c13d0000000008870436000000400920003900000000090904330000ffff0090008c000000bc0000213d00000000009804350000006008200039000000000808043300000a090080009c000000bc0000213d000000400970003900000000008904350000008008200039000000000808043300000a090080009c000000bc0000213d00000060097000390000000000890435000000a008200039000000000808043300000a090080009c000000bc0000213d00000080097000390000000000890435000000c00820003900000000080804330000ffff0080008c000000bc0000213d000000a0097000390000000000890435000000e008200039000000000808043300000a090080009c000000bc0000213d000000c0097000390000000000890435000001000820003900000000080804330000ffff0080008c000000bc0000213d000000e0097000390000000000890435000001200820003900000000080804330000ffff0080008c000000bc0000213d00000100097000390000000000890435000001400820003900000000080804330000ffff0080008c000000bc0000213d000001200970003900000000008904350000016008200039000000000808043300000a090080009c000000bc0000213d000001400970003900000000008904350000018008200039000000000808043300000a090080009c000000bc0000213d00000160097000390000000000890435000001a008200039000000000808043300000a0f0080009c000000bc0000213d00000180097000390000000000890435000001c008200039000000000808043300000a090080009c000000bc0000213d000001a0097000390000000000890435000001e008200039000000000808043300000a090080009c000000bc0000213d000001c009700039000000000089043500000200082000390000000008080433000000000008004b0000000009000039000000010900c039000000000098004b000000bc0000c13d000001e00970003900000000008904350000022008200039000000000808043300000a1600800198000000bc0000c13d00000200097000390000000000890435000000000076043500000000045404360000024002200039000000000032004b000017c90000413d000002000000013d0000002101000029000000000101043300010a160010019b000000110000006b000016440000613d00000002010003670000001402100360000000000202043b000800000002001d00000a0e0020009c000000bc0000213d0000001202100360000000000302043b0000000002000031000000250420006a000000230440008a00000a110540019700000a1106300197000000000756013f000000000056004b000000000500001900000a1105004041000000000043004b000000000400001900000a110400804100000a110070009c000000000504c019000000000005004b000000bc0000c13d0000001303300029000000000131034f000000000101043b000c00000001001d00000a0f0010009c000000bc0000213d0000000c0100002900000006011002100000000001120049000000200530003900000a110210019700000a1103500197000000000423013f000000000023004b000000000200001900000a1102004041000b00000005001d000000000015004b000000000100001900000a110100204100000a110040009c000000000201c019000000000002004b000000bc0000c13d0000000c0000006b001d00000000001d001f00000000001d002200000000001d0000164a0000613d002100000000001d001f00000000001d001d00000000001d002200000000001d0000188f0000013d00000021020000290000000102200039002100000002001d0000000c0020006c0000164a0000813d000000210100002900000006011002100000000b01100029000000000210007900000a120020009c000000bc0000213d000000400020008c000000bc0000413d000000400200043d002500000002001d00000a130020009c000000e00000213d00000025020000290000004002200039000000400020043f0000000202000367000000000312034f000000000303043b00000a0e0030009c000000bc0000213d000000250400002900000000033404360000002001100039000000000112034f000000000101043b001000000003001d00000000001304350000002601000029000000000010043f0000000a01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002502000029000000000202043300000a0e02200197000000000020043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400200043d002000000002001d00000a140020009c000000e00000213d000000000101043b0000002004000029000000c002400039000000400020043f000000000201041a000000a00140003900000a88002001980000000003000039000000010300c039000000000031043500000a09012001970000000003140436000000700120027000000a09011001970000008005400039001200000005001d0000000000150435000000200120027000000a0901100197000d00000003001d0000000000130435000000500120027000000a09011001970000006003400039000f00000003001d0000000000130435000000400440003900000040022002700000ffff0320018f000e00000004001d0000000000340435000018f80000613d000000000003004b0000190e0000613d0000002502000029000000000202043300000a0e04200197002500000004001d000000080040006c000019100000c13d000000060000006b0000198f0000613d00000010010000290000000001010433000000060300002900001a5d0000013d000000230100002900000000010104330000ffff0110018f00000a89011000d1000000220010002a0000175c0000413d0000001d0200002900000a09022001970000002403000029000000000303043300000a09033001970000000002230019001d00000002001d00000a090020009c0000175c0000213d0000001f0200002900000a090220019700000a8a0020009c0000175c0000213d002200220010002d001f00200020003d0000188a0000013d000000000200001900001a690000013d000000400100043d00000a130010009c000000e00000213d0000004002100039000000400020043f0000002002100039000000000002043500000000000104350000002501000029000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400200043d000a00000002001d00000a130020009c000000e00000213d000000000101043b0000000a030000290000004002300039000000400020043f000000000101041a00000a16021001970000000002230436000000e001100270000900000002001d000700000001001d000000000012043500000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f000000010020019000001a910000613d000000000101043b00070007001000740000175c0000413d00000a8b0100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f000000010020019000001a910000613d000000000101043b00000a0901100197000000070010006b00001a530000413d0000002501000029000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000400200043d00000a0c0020009c000000e00000213d000000000101043b0000006003200039000000400030043f000000000301041a000000400120003900000a8d003001980000000004000039000000010400c039000000000041043500000a0e013001970000000004120436000000a003300270000000ff0330018f000400000004001d000000000034043500001a530000613d000000000001004b00001a530000613d000000400100043d00000a130010009c000000e00000213d0000004003100039000000400030043f0000002003100039000000000003043500000000000104350000000002020433000000400300043d00000a8e01000041000700000003001d0000000001130436000200000001001d000000000100041400000a0e02200197000300000002001d000000040020008c000019910000c13d0000000103000031000000a00030008c000000a0040000390000000004034019000019ba0000013d000000000200001900001a650000013d000000070200002900000a090020009c00000a0902008041000000400220021000000a090010009c00000a0901008041000000c001100210000000000121019f00000a4e011001c700000003020000292820281b0000040f000000600310027000000a0903300197000000a00030008c000000a0040000390000000004034019000000e0064001900000000705600029000019aa0000613d000000000701034f0000000708000029000000007907043c0000000008980436000000000058004b000019a60000c13d0000001f07400190000019b70000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f000000010020019000001df00000613d0000001f01400039000001e00110018f0000000702100029000000000012004b00000000010000390000000101004039000500000002001d00000a0f0020009c000000e00000213d0000000100100190000000e00000c13d0000000501000029000000400010043f000000a00030008c000000bc0000413d0000000701000029000000000101043300000a8f0010009c000000bc0000213d00000007010000290000008001100039000000000101043300000a8f0010009c000000bc0000213d00000002010000290000000001010433000200000001001d00000a110010009c00001dfc0000813d000000070100002900000060011000390000000001010433000700000001001d00000a90010000410000000502000029000000000012043500000000010004140000000302000029000000040020008c000000200400003900001a0c0000613d000000050200002900000a090020009c00000a0902008041000000400220021000000a090010009c00000a0901008041000000c001100210000000000121019f00000a4e011001c700000003020000292820281b0000040f000000600310027000000a0903300197000000200030008c0000002004000039000000000403401900000020064001900000000505600029000019fc0000613d000000000701034f0000000508000029000000007907043c0000000008980436000000000058004b000019f80000c13d0000001f0740019000001a090000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f000000010020019000001dff0000613d0000001f01400039000000600110018f000000050110002900000a0f0010009c000000e00000213d000000400010043f000000200030008c000000bc0000413d00000005020000290000000002020433000000ff0020008c000000bc0000213d00000004030000290000000003030433000000ff0330018f0000000003230019000000ff0030008c0000175c0000213d000000240230008c00001a320000213d000000010200003900001a2b0000613d00000024033000890000000a04000039000000010030019000000000054400a9000000010400603900000000022400a90000000103300272000000000405001900001a240000c13d000000020000006b00001a8d0000613d00000002032000b900000002043000fa000000000024004b00001a400000613d0000175c0000013d0000004d0020008c0000175c0000213d00000001030000390000000a04000039000000010020019000000000054400a9000000010400603900000000033400a90000000102200272000000000405001900001a360000c13d000000000003004b00001f160000613d00000002033000f900000a160030009c0000172c0000213d00000a130010009c000000e00000213d0000004002100039000000400020043f000000070200002900000a09022001970000002004100039000000000024043500000000003104350000000903000029000000000303043300000a0903300197000000000032004b00000000020100190000000a02004029000a00000002001d000900200020003d0000000901000029000000000101043300000a090010019800001da10000613d0000000a01000029000000000101043300000a160310019800001da10000613d0000001001000029000000000101043300000000023100a900000000033200d9000000000013004b0000175c0000c13d0000000f0100002900000000010104330000000e0300002900000000030304330000ffff0330018f00000a930220012a00000000023200a900000a940220012a0000001d0300002900000a090330019700000a09011001970000000001310019001d00000001001d00000a090010009c0000175c0000213d0000001f0100002900000a09011001970000001203000029000000000303043300000a09033001970000000001130019001f00000001001d00000a090010009c0000175c0000213d0000002001000029000000000101043300000a090110019700000a89011000d1000000000012004b00001a850000413d0000000d01000029000000000101043300000a090110019700000a89011000d1000000000012004b00001a890000a13d000000220010002a0000175c0000413d002200220010002d0000188a0000013d000000220020002a0000175c0000413d002200220020002d0000188a0000013d000000000300001900000a130010009c00001a440000a13d000000e00000013d000000000001042f0000000102000039000000000302041a00000a1703300197000000000113019f000000000012041b0000001f0100002900000a180010009c000000e00000213d0000001f010000290000002002100039001b00000002001d000000400020043f0000000000010435000000400100043d001700000001001d00000a130010009c000000e00000213d00000017030000290000004001300039000000400010043f00000020013000390000001f020000290000000000210435000000230100002900000000001304350000000001020433000000000001004b00001ac10000c13d00000023010000290000000001010433000000000001004b00001b320000c13d0000002501000029000000000101043300000a0e0110019800001abe0000613d0000002602000029000000000202043300000a0d0020019800001abe0000613d0000002402000029000000000202043300000a090020019800001b870000c13d000000400100043d00000a4d0200004100000ae10000013d002100000000001d00001ac90000013d0000002102000029002100010020003d0000001f010000290000000001010433000000210010006b00001b830000813d000000210100002900000005011002100000001b01100029000000000101043300000a0e01100197002000000001001d000000000010043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a002300000001001d000000000001004b00001ac30000613d0000000201000039000000000201041a000000000002004b0000175c0000613d0000002303000029000000010130008a000000000032004b00001b050000613d000000000012004b00001e580000a13d000000230100002900000a1a0110009a00000a1a0220009a000000000202041a000000000021041b000000000020043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002302000029000000000021041b0000000201000039000000000101041a002300000001001d000000000001004b00001b2c0000613d0000002301000029000000010110008a000000230200002900000a1a0220009a000000000002041b0000000202000039000000000012041b0000002001000029000000000010043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000001041b000000400100043d0000002002000029000000000021043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a1b011001c70000800d02000039000000010300003900000a1c04000041282028160000040f000000010020019000001ac30000c13d000000bc0000013d00000aaa01000041000000000010043f0000003101000039000000040010043f00000a3d0100004100002822000104300000002301000029001f00200010003d002000000000001d000000200100002900000005011002100000001f01100029000000000101043300210a0e0010019c00001b800000613d0000002101000029000000000010043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000000001004b00001b670000c13d0000000201000039000000000101041a00000a0f0010009c000000e00000213d00000001021000390000000203000039000000000023041b00000a1d0110009a0000002102000029000000000021041b000000000103041a001b00000001001d000000000020043f0000000301000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000001b02000029000000000021041b000000400100043d0000002102000029000000000021043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a1b011001c70000800d02000039000000010300003900000a1e04000041282028160000040f0000000100200190000000bc0000613d0000002002000029002000010020003d00000023010000290000000001010433000000200010006b00001b350000413d00001ab20000013d000000400100043d00000aa20200004100000ae10000013d00000017010000290000000001010433002300000001001d00001aae0000013d000000a00010043f0000002601000029000000000101043300000a0d01100197000000800010043f0000002401000029000000000101043300000a0901100197000000c00010043f000000400100043d00000a180010009c000000e00000213d0000002002100039000000400020043f000000000001043500000022010000290000000001010433000000000001004b00001e0b0000c13d0000001d010000290000000001010433000000000001004b00001be90000613d002600000000001d000000260100002900000005011002100000001c01100029000000000101043300000020021000390000000002020433002400000002001d000000000101043300000a0e01100197002500000001001d000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000201041a00000a2202200197000000240600002900000020036000390000000004030433000000a00440021000000a2304400197000000000242019f00000040046000390000000005040433000000000005004b00000a24050000410000000005006019000000000252019f000000000506043300000a0e05500197000000000252019f000000000021041b000000400100043d00000000025104360000000003030433000000ff0330018f00000000003204350000000002040433000000000002004b0000000002000039000000010200c0390000004003100039000000000023043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a25011001c70000800d02000039000000020300003900000a26040000410000002505000029282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d0000001d010000290000000001010433000000260010006b00001b9f0000413d00000019010000290000000001010433000000000001004b00001cb10000613d002400000000001d0000002401000029000000050110021000000018011000290000000001010433000000001201043400250a0f0020019c00001e5e0000613d0000000001010433002600000001001d0000016001100039002200000001001d000000000101043300000a090110019800001e5e0000613d00000026020000290000020002200039002300000002001d000000000202043300000a270220019700000a280020009c00001e5e0000c13d00000026020000290000006002200039002100000002001d000000000202043300000a0902200197000000000021004b00001e5e0000213d0000002501000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000000101100039000000000101041a002000000001001d000000400100043d001f00000001001d000000260200002928201f1f0000040f0000001f02000029000000000121004900000a090020009c00000a090200804100000a090010009c00000a090100804100000040022002100000006001100210000000000121019f000000200200002900000a290020019800001c350000613d000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a20011001c70000800d02000039000000020300003900000a2a0400004100001c3e0000013d000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a20011001c70000800d02000039000000020300003900000a2b040000410000002505000029282028160000040f0000000100200190000000bc0000613d0000002501000029000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d00000026040000290000000032040434000000000002004b000000000101043b000000000201041a00000a2c02200197000000010220c1bf0000000003030433000000080330021000000a2d03300197000000000232019f00000040034000390000000003030433000000180330021000000a2e03300197000000000232019f00000021030000290000000003030433000000380330021000000a2f03300197000000000232019f00000080034000390000000003030433000000580330021000000a3003300197000000000232019f000000a0034000390000000003030433000000780330021000000a3103300197000000000232019f000000c0034000390000000003030433000000880330021000000a3203300197000000000232019f000000e0034000390000000003030433000000a80330021000000a3303300197000000000232019f00000100034000390000000003030433000000b80330021000000a3403300197000000000232019f00000120034000390000000003030433000000c80330021000000a3503300197000000000232019f00000140034000390000000003030433000000d80330021000000a3603300197000000000232019f000000000021041b00000001011000390000002202000029000000000202043300000a0902200197000000000301041a00000a3703300197000000000223019f00000180034000390000000003030433000000200330021000000a3803300197000000000232019f000001a0034000390000000003030433000000600330021000000a3903300197000000000232019f000001c0034000390000000003030433000000800330021000000a3a03300197000000000232019f000001e0034000390000000003030433000000000003004b00000a3b030000410000000003006019000000000232019f00000023030000290000000003030433000000380330027000000a2903300197000000000232019f000000000021041b0000002402000029002400010020003d00000019010000290000000001010433000000240010006b00001bee0000413d001a002d0000002d0000001a010000290000000001010433000000000001004b00001ced0000613d0000001a01000029002300200010003d002600000000001d000000260100002900000005011002100000002301100029000000000101043300000020021000390000000002020433002400000002001d000000000101043300000a0e01100197002500000001001d000000000010043f0000000801000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000240200002900000a0f02200197000000000101043b000000000301041a00000a3e03300197000000000323019f000000000031041b000000400100043d000000000021043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a1b011001c70000800d02000039000000020300003900000a3f040000410000002505000029282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d0000001a010000290000000001010433000000260010006b00001cb80000413d000000400100043d001c00000001001d00000a180010009c000000e00000213d0000001c010000290000002002100039001b00000002001d000000400020043f00000000000104350000002e01000029001d00000001001d0000000021010434001e00000002001d000000000001004b00001ea00000613d001f00000000001d00001d020000013d0000001f02000029001f00010020003d0000001f0010006b00001e610000813d0000001f0200002900000005022002100000001e0220002900000000030204330000002002300039002000000002001d00000000020204330000000004020433000000000004004b00001cfe0000613d000000000103043300230a0f0010019b002600000000001d0000002601000029000000050110021000000000011200190000002001100039000000000201043300000020012000390000000001010433002400000001001d000000003101043400000a0901100197002200000003001d000000000303043300000a0903300197000000000031004b00001d950000813d000000000102043300250a0e0010019b00000024010000290000008001100039002100000001001d000000000101043300000a09011001970000001f0010008c00001eb00000a13d0000002301000029000000000010043f0000000a01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002502000029000000000020043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d00000022080000290000000002080433000000200220021000000a4202200197000000000101043b000000000301041a00000a4303300197000000000223019f000000240900002900000040039000390000000004030433000000400440021000000a4404400197000000000242019f00000060049000390000000005040433000000500550021000000a4505500197000000000252019f00000021070000290000000005070433000000700550021000000a4605500197000000000252019f000000a0059000390000000006050433000000000006004b00000a47060000410000000006006019000000000262019f000000000609043300000a0906600197000000000262019f000000000021041b000000400100043d0000000002610436000000000608043300000a0906600197000000000062043500000000020304330000ffff0220018f00000040031000390000000000230435000000000204043300000a090220019700000060031000390000000000230435000000000207043300000a0902200197000000800310003900000000002304350000000002050433000000000002004b0000000002000039000000010200c039000000a003100039000000000023043500000a090010009c00000a09010080410000004001100210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a48011001c70000800d02000039000000030300003900000a490400004100000023050000290000002506000029282028160000040f0000000100200190000000bc0000613d0000002603000029002600010030003d000000200100002900000000020104330000000001020433000000260010006b00001d0f0000413d0000001d01000029000000000101043300001cfe0000013d000000400200043d0000002404200039000000000034043500000a400300004100000000003204350000000403200039000000000013043500000a090020009c00000a0902008041000000400120021000000a41011001c70000282200010430000000400100043d00000a9202000041000005e70000013d000000000681034f0000000307500210000000200770008900000a270970021f000000040b50008c00000a2709008041000000000a06043b000000bc0000413d0000001b0650003900000abd066001970000003f0660003900000abd07600197000000400600043d0000000007760019000000000067004b000000000c000039000000010c00403900000a0f0070009c000000e00000213d0000000100c00190000000e00000c13d000000400070043f0000000007b60436000000000c85001900000000002c004b000000bc0000213d0000000402800039000000000821034f00000abd0cb001980000001f0bb0018f0000000002c7001900001dca0000613d000000000d08034f000000000e07001900000000df0d043c000000000efe043600000000002e004b00001dc60000c13d00000000099a016f00000000000b004b00001dd80000613d0000000008c8034f000000030ab00210000000000b020433000000000bab01cf000000000bab022f000000000808043b000001000aa000890000000008a8022f0000000008a801cf0000000008b8019f000000000082043500000000026500190000001c02200039000000000002043500000a970090009c00001eb80000613d00000a980090009c00001ded0000c13d000000000206043300000a120020009c000000bc0000213d000000200020008c000000bc0000413d000000400500043d00000a130050009c000000e00000213d00000000020704330000004006500039000000400060043f0000000000250435000000000600001900001ecb0000013d000000400100043d00000aaf0200004100000ae10000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000013bb0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00001df70000c13d000013bb0000013d00000a91010000410000000502000029000002070000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000013bb0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00001e060000c13d000013bb0000013d002600000000001d00001e130000013d0000002602000029002600010020003d00000022010000290000000001010433000000260010006b00001b9a0000813d000000260100002900000005011002100000001e01100029002400000001001d000000000101043300000a0e01100197002500000001001d000000000010043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a000000000001004b00001e0d0000c13d0000000b01000039000000000101041a00000a0f0010009c000000e00000213d00000001021000390000000b03000039000000000023041b00000a1f0110009a0000002502000029000000000021041b000000000103041a002300000001001d000000000020043f0000000c01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b0000002302000029000000000021041b00000022010000290000000001010433000000260010006c00001e580000a13d00000024010000290000000001010433000000000200041400000a0e0510019700000a090020009c00000a0902008041000000c00120021000000a20011001c70000800d02000039000000020300003900000a2104000041282028160000040f000000010020019000001e0d0000c13d000000bc0000013d00000aaa01000041000000000010043f0000003201000039000000040010043f00000a3d010000410000282200010430000000400100043d00000a3c02000041000005e70000013d0000001c010000290000000001010433000000000001004b00001ea00000613d002600000000001d000000260100002900000005011002100000001b01100029000000000101043300000020021000390000000002020433002500000002001d000000000101043300000a0f01100197002400000001001d000000000010043f0000000a01000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000250200002900000a0e02200197000000000101043b002500000002001d000000000020043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000001041b000000000100041400000a090010009c00000a0901008041000000c00110021000000a20011001c70000800d02000039000000030300003900000a4b0400004100000024050000290000002506000029282028160000040f0000000100200190000000bc0000613d0000002602000029002600010020003d0000001c010000290000000001010433000000260010006b00001e660000413d000000800100043d00000140000004430000016000100443000000a00100043d00000020020000390000018000200443000001a000100443000000c00100043d0000004003000039000001c000300443000001e00010044300000100002004430000000301000039000001200010044300000a4c01000041000028210001042e000000400200043d0000002403200039000000000013043500000a4a01000041000000000012043500000004012000390000002503000029000013a20000013d000000000206043300000a120020009c000000bc0000213d000000400020008c000000bc0000413d000000400500043d00000a130050009c000000e00000213d0000004002500039000000400020043f0000000002070433000000000025043500000040066000390000000006060433000000000006004b0000000007000039000000010700c039000000000076004b000000bc0000c13d000000200550003900000000006504350000001b05000029000000000505043300000a0905500197000000000052004b00001ed50000a13d000000400100043d00000a9c0200004100000ae10000013d000000000006004b00001ede0000c13d0000001a050000290000000005050433000000000005004b00001ede0000613d000000400100043d00000a9b0200004100000ae10000013d00000a09033001970000001d033000290000000003430019000000000032001a0000175c0000413d000000010400002900000a9904400198002500000000001d00001ef60000613d000000000332001900000000024300a900000000044200d9000000000034004b0000175c0000c13d000000000002004b002500000000001d00001ef60000613d0000001c03000029000000000303043300000a0f0330019700250000002300ad00000025022000f9000000000032004b0000175c0000c13d0000001401100360000000000101043b00000a0e0010009c000000bc0000213d000000000010043f0000000801000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000000bc0000613d000000000101043b000000000101041a00000a0f0210019700000022012000b9000000220000006b00001f0f0000613d00000022031000fa000000000023004b0000175c0000c13d000000250010002a0000175c0000413d0000002501100029000000260010002a0000175c0000413d000000060000006b00001f1c0000c13d00000aaa01000041000000000010043f0000001201000039000000040010043f00000a3d010000410000282200010430000000260110002900000006011000fa00000a880000013d0000000043020434000000000003004b0000000003000039000000010300c039000000000331043600000000040404330000ffff0440018f00000000004304350000004003200039000000000303043300000a0903300197000000400410003900000000003404350000006003200039000000000303043300000a0903300197000000600410003900000000003404350000008003200039000000000303043300000a090330019700000080041000390000000000340435000000a00320003900000000030304330000ffff0330018f000000a0041000390000000000340435000000c003200039000000000303043300000a0903300197000000c0041000390000000000340435000000e00320003900000000030304330000ffff0330018f000000e0041000390000000000340435000001000320003900000000030304330000ffff0330018f00000100041000390000000000340435000001200320003900000000030304330000ffff0330018f000001200410003900000000003404350000014003200039000000000303043300000a0903300197000001400410003900000000003404350000016003200039000000000303043300000a0903300197000001600410003900000000003404350000018003200039000000000303043300000a0f0330019700000180041000390000000000340435000001a003200039000000000303043300000a0903300197000001a0041000390000000000340435000001c003200039000000000303043300000a0903300197000001c0041000390000000000340435000001e0032000390000000003030433000000000003004b0000000003000039000000010300c039000001e00410003900000000003404350000020002200039000000000202043300000a2702200197000002000310003900000000002304350000022001100039000000000001042d00000abf0010009c00001f7b0000813d0000006001100039000000400010043f000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000ac00010009c00001f860000813d0000004001100039000000400010043f000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000ac10010009c00001f910000813d000000c001100039000000400010043f000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d0100004100002822000104300000001f0220003900000abd022001970000000001120019000000000021004b0000000002000039000000010200403900000a0f0010009c00001fa30000213d000000010020019000001fa30000c13d000000400010043f000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d010000410000282200010430000000400100043d00000ac20010009c00001faf0000813d0000022002100039000000400020043f000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000000430104340000000001320436000000000003004b00001fc10000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b00001fba0000413d000000000231001900000000000204350000001f0230003900000abd022001970000000001210019000000000001042d00000020030000390000000004310436000000000302043300000000003404350000004001100039000000000003004b00001fd60000613d00000000040000190000002002200039000000000502043300000a0e0550019700000000015104360000000104400039000000000034004b00001fcf0000413d000000000001042d000000003101043400000a16011001970000000001120436000000000203043300000a09022001970000000000210435000000000001042d0000000043020434000000000003004b0000000003000039000000010300c039000000000331043600000000040404330000ffff0440018f00000000004304350000004003200039000000000303043300000a0903300197000000400410003900000000003404350000006003200039000000000303043300000a0903300197000000600410003900000000003404350000008003200039000000000303043300000a090330019700000080041000390000000000340435000000a00320003900000000030304330000ffff0330018f000000a0041000390000000000340435000000c003200039000000000303043300000a0903300197000000c0041000390000000000340435000000e00320003900000000030304330000ffff0330018f000000e0041000390000000000340435000001000320003900000000030304330000ffff0330018f00000100041000390000000000340435000001200320003900000000030304330000ffff0330018f000001200410003900000000003404350000014003200039000000000303043300000a0903300197000001400410003900000000003404350000016003200039000000000303043300000a0903300197000001600410003900000000003404350000018003200039000000000303043300000a0f0330019700000180041000390000000000340435000001a003200039000000000303043300000a0903300197000001a0041000390000000000340435000001c003200039000000000303043300000a0903300197000001c0041000390000000000340435000001e0032000390000000003030433000000000003004b0000000003000039000000010300c039000001e00410003900000000003404350000020002200039000000000202043300000a2702200197000002000310003900000000002304350000022001100039000000000001042d000000004302043400000a09033001970000000003310436000000000404043300000a09044001970000000000430435000000400320003900000000030304330000ffff0330018f000000400410003900000000003404350000006003200039000000000303043300000a0903300197000000600410003900000000003404350000008003200039000000000303043300000a090330019700000080041000390000000000340435000000a0022000390000000002020433000000000002004b0000000002000039000000010200c039000000a0031000390000000000230435000000c001100039000000000001042d000000004302043400000a0e0330019700000000033104360000000004040433000000ff0440018f000000000043043500000040022000390000000002020433000000000002004b0000000002000039000000010200c039000000400310003900000000002304350000006001100039000000000001042d0009000000000002000500000003001d000900000002001d000000400200043d00000ac00020009c0000231e0000813d0000004003200039000000400030043f00000020032000390000000000030435000000000002043500000a0e01100197000400000001001d000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000023240000613d000000400300043d00000a130030009c0000231e0000213d000000000101043b0000004002300039000000400020043f000000000101041a00000a1602100197000800000003001d0000000002230436000000e001100270000600000002001d000700000001001d000000000012043500000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f0000000100200190000023260000613d000000000101043b0007000700100074000023270000413d00000a8b0100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f0000000100200190000023260000613d000000000101043b00000a0901100197000000070010006b000020b00000813d00000008050000290000000601000029000000000101043300000a0900100198000021b50000c13d0000232d0000013d0000000401000029000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000001002001900000000805000029000023240000613d000000400200043d00000a0c0020009c0000231e0000213d000000000101043b0000006003200039000000400030043f000000000301041a000000400120003900000a8d003001980000000004000039000000010400c039000000000041043500000a0e013001970000000006120436000000a003300270000000ff0330018f0000000000360435000020e80000613d000000000001004b000020e80000613d000000400100043d00000a130010009c0000231e0000213d0000004003100039000000400030043f0000002003100039000000000003043500000000000104350000000002020433000000400c00043d00000a8e0100004100000000051c0436000000000100041400000a0e02200197000000040020008c000700000006001d000020ed0000c13d0000000103000031000000a00030008c000000a00400003900000000040340190000211d0000013d0000000601000029000000000101043300000a0900100198000021b50000c13d0000232d0000013d000100000005001d00000a0900c0009c00000a090300004100000000030c4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c7000200000002001d00030000000c001d2820281b0000040f000000030c000029000000600310027000000a0903300197000000a00030008c000000a00400003900000000040340190000001f0640018f000000e00740019000000000057c00190000210a0000613d000000000801034f00000000090c0019000000008a08043c0000000009a90436000000000059004b000021060000c13d000000000006004b000021170000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000023520000613d0000000706000029000000020200002900000001050000290000001f01400039000001e00110018f000000000bc1001900000000001b004b0000000001000039000000010100403900000a0f00b0009c0000231e0000213d00000001001001900000231e0000c13d0000004000b0043f000000a00030008c000023240000413d00000000010c043300000a8f0010009c000023240000213d0000008001c00039000000000101043300000a8f0010009c000023240000213d000000000805043300000a110080009c0000233e0000813d0000006001c00039000000000701043300000a900100004100000000001b04350000000001000414000000040020008c0000213d0000c13d00000020040000390000216d0000013d000100000008001d000200000007001d00000a0900b0009c00000a090300004100000000030b4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c700030000000b001d2820281b0000040f000000030b000029000000600310027000000a0903300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b00190000215a0000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b000021560000c13d000000000006004b000021670000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00000001002001900000235e0000613d0000000706000029000000020700002900000001080000290000001f01400039000000600110018f0000000001b1001900000a0f0010009c0000231e0000213d000000400010043f000000200030008c000023240000413d00000000020b0433000000ff0020008c000023240000213d0000000003060433000000ff0330018f0000000003230019000000ff0030008c000023270000213d000000240230008c000021820000213d000021910000c13d00000001020000390000219b0000013d0000004d0020008c000023270000213d0000000a040000390000000103000039000000010020019000000000054400a9000000010400603900000000033400a900000001022002720000000004050019000021860000c13d000000000003004b0000234c0000613d00000000033800d9000021a10000013d0000000a0400003900000001020000390000002403300089000000010030019000000000054400a9000000010400603900000000022400a900000001033002720000000004050019000021940000c13d000000000008004b000023150000613d00000000038200a900000000048300d9000000000024004b000023270000c13d00000a160030009c000023450000213d00000a130010009c00000008050000290000231e0000213d0000004002100039000000400020043f00000a09027001970000002004100039000000000024043500000000003104350000000603000029000000000303043300000a0903300197000000000032004b00000000050180190000002001500039000000000101043300000a09001001980000232d0000613d000000000105043300000a16011001980000232d0000613d00080009001000bd000000090000006b000021bf0000613d000000090300002900000008023000f9000000000012004b000023270000c13d000000400100043d00000a130010009c0000231e0000213d0000004002100039000000400020043f000000200210003900000000000204350000000000010435000000050100002900000a0e01100197000500000001001d000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000023240000613d000000400300043d00000a130030009c0000231e0000213d000000000101043b0000004002300039000000400020043f000000000101041a00000a1602100197000900000003001d0000000002230436000000e001100270000600000002001d000700000001001d000000000012043500000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f0000000100200190000023260000613d000000000101043b0007000700100074000023270000413d00000a8b0100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f0000000100200190000023260000613d000000000101043b00000a0901100197000000070010006b0000220b0000813d00000009050000290000000601000029000000000101043300000a0900100198000023100000c13d000023330000013d0000000501000029000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000001002001900000000905000029000023240000613d000000400200043d00000a0c0020009c0000231e0000213d000000000101043b0000006003200039000000400030043f000000000301041a000000400120003900000a8d003001980000000004000039000000010400c039000000000041043500000a0e013001970000000006120436000000a003300270000000ff0330018f0000000000360435000022430000613d000000000001004b000022430000613d000000400100043d00000a130010009c0000231e0000213d0000004003100039000000400030043f0000002003100039000000000003043500000000000104350000000002020433000000400c00043d00000a8e0100004100000000051c0436000000000100041400000a0e02200197000000040020008c000700000006001d000022480000c13d0000000103000031000000a00030008c000000a0040000390000000004034019000022780000013d0000000601000029000000000101043300000a0900100198000023100000c13d000023330000013d000200000005001d00000a0900c0009c00000a090300004100000000030c4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c7000300000002001d00040000000c001d2820281b0000040f000000040c000029000000600310027000000a0903300197000000a00030008c000000a00400003900000000040340190000001f0640018f000000e00740019000000000057c0019000022650000613d000000000801034f00000000090c0019000000008a08043c0000000009a90436000000000059004b000022610000c13d000000000006004b000022720000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00000001002001900000236a0000613d0000000706000029000000030200002900000002050000290000001f01400039000001e00110018f000000000bc1001900000000001b004b0000000001000039000000010100403900000a0f00b0009c0000231e0000213d00000001001001900000231e0000c13d0000004000b0043f000000a00030008c000023240000413d00000000010c043300000a8f0010009c000023240000213d0000008001c00039000000000101043300000a8f0010009c000023240000213d000000000805043300000a110080009c0000233e0000813d0000006001c00039000000000701043300000a900100004100000000001b04350000000001000414000000040020008c000022980000c13d0000002004000039000022c80000013d000200000008001d000300000007001d00000a0900b0009c00000a090300004100000000030b4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c700040000000b001d2820281b0000040f000000040b000029000000600310027000000a0903300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b0019000022b50000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b000022b10000c13d000000000006004b000022c20000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000023760000613d0000000706000029000000030700002900000002080000290000001f01400039000000600110018f0000000001b1001900000a0f0010009c0000231e0000213d000000400010043f000000200030008c000023240000413d00000000020b0433000000ff0020008c000023240000213d0000000003060433000000ff0330018f0000000003230019000000ff0030008c000023270000213d000000240230008c000022dd0000213d000022ec0000c13d0000000102000039000022f60000013d0000004d0020008c000023270000213d0000000a040000390000000103000039000000010020019000000000054400a9000000010400603900000000033400a900000001022002720000000004050019000022e10000c13d000000000003004b0000234c0000613d00000000033800d9000022fc0000013d0000000a0400003900000001020000390000002403300089000000010030019000000000054400a9000000010400603900000000022400a900000001033002720000000004050019000022ef0000c13d000000000008004b0000231a0000613d00000000038200a900000000048300d9000000000024004b000023270000c13d00000a160030009c000023450000213d00000a130010009c00000009050000290000231e0000213d0000004002100039000000400020043f00000a09027001970000002004100039000000000024043500000000003104350000000603000029000000000303043300000a0903300197000000000032004b00000000050180190000002001500039000000000101043300000a0900100198000023330000613d000000000105043300000a1601100198000023330000613d00000008011000f9000000000001042d000000000300001900000a130010009c0000000805000029000021a60000a13d0000231e0000013d000000000300001900000a130010009c0000000905000029000023010000a13d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000000010000190000282200010430000000000001042f00000aaa01000041000000000010043f0000001101000039000000040010043f00000a3d010000410000282200010430000000400100043d00000a9202000041000000000021043500000004021000390000000403000029000023380000013d000000400100043d00000a9202000041000000000021043500000004021000390000000503000029000000000032043500000a090010009c00000a0901008041000000400110021000000a3d011001c7000028220001043000000a910100004100000000001b043500000a0900b0009c00000a090b0080410000004001b0021000000a4e011001c7000028220001043000000a9102000041000000000021043500000a090010009c00000a0901008041000000400110021000000a4e011001c7000028220001043000000aaa01000041000000000010043f0000001201000039000000040010043f00000a3d0100004100002822000104300000001f0530018f00000a0b06300198000000400200043d0000000004620019000023810000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000023590000c13d000023810000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000023810000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000023650000c13d000023810000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000023810000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000023710000c13d000023810000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000023810000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000237d0000c13d000000000005004b0000238e0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f0000000000140435000000600130021000000a090020009c00000a09020080410000004002200210000000000112019f000028220001043000000a0e02200197000000000020043f000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000023a20000613d000000000101043b000000000001042d00000000010000190000282200010430000000400100043d00000abf0010009c000023af0000813d0000006002100039000000400020043f00000040021000390000000000020435000000200210003900000000000204350000000000010435000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000a0f01100197000000000010043f0000000901000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000023c40000613d000000000101043b000000000001042d000000000100001900002822000104300000000003010433000000000023004b000023cd0000a13d000000050220021000000000012100190000002001100039000000000001042d00000aaa01000041000000000010043f0000003201000039000000040010043f00000a3d010000410000282200010430000000400100043d00000ac00010009c000023dc0000813d0000004002100039000000400020043f000000200210003900000000000204350000000000010435000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d0100004100002822000104300000000002010019000000400100043d00000ac20010009c000024300000813d0000022003100039000000400030043f000000000302041a000000d80430027000000a090440019700000140051000390000000000450435000000c8043002700000ffff0440018f00000120051000390000000000450435000000b8043002700000ffff0440018f00000100051000390000000000450435000000a8043002700000ffff0440018f000000e0051000390000000000450435000000880430027000000a0904400197000000c005100039000000000045043500000078043002700000ffff0440018f000000a0051000390000000000450435000000580430027000000a090440019700000080051000390000000000450435000000380430027000000a090440019700000060051000390000000000450435000000180430027000000a09044001970000004005100039000000000045043500000008043002700000ffff0440018f00000020051000390000000000450435000000ff003001900000000003000039000000010300c03900000000003104350000000102200039000000000202041a000001600310003900000a09042001970000000000430435000000380320021000000a27033001970000020004100039000000000034043500000a23002001980000000003000039000000010300c039000001e0041000390000000000340435000000800320027000000a0903300197000001c0041000390000000000340435000000600320027000000a0903300197000001a0041000390000000000340435000000200220027000000a0f0220019700000180031000390000000000230435000000000001042d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d0100004100002822000104300006000000000002000000400200043d00000ac00020009c000025820000813d0000004003200039000000400030043f00000020032000390000000000030435000000000002043500000a0e01100197000400000001001d000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000025880000613d000000400300043d00000a130030009c000025820000213d000000000101043b0000004002300039000000400020043f000000000101041a00000a1602100197000600000003001d0000000002230436000000e001100270000300000002001d000500000001001d000000000012043500000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f00000001002001900000258a0000613d000000000101043b00050005001000740000258b0000413d00000a8b0100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f00000001002001900000258a0000613d000000000101043b00000a0901100197000000050010006b0000247e0000813d0000000605000029000024b60000013d0000000401000029000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000001002001900000000605000029000025880000613d000000400200043d00000a0c0020009c000025820000213d000000000101043b0000006003200039000000400030043f000000000301041a000000400120003900000a8d003001980000000004000039000000010400c039000000000041043500000a0e013001970000000006120436000000a003300270000000ff0330018f0000000000360435000024b60000613d000000000001004b000024b60000613d000000400100043d00000a130010009c000025820000213d0000004003100039000000400030043f0000002003100039000000000003043500000000000104350000000002020433000000400c00043d00000a8e0100004100000000051c0436000000000100041400000a0e02200197000000040020008c000500000006001d000024b80000c13d0000000103000031000000a00030008c000000a0040000390000000004034019000024e80000013d0000000001050019000000000001042d000100000005001d00000a0900c0009c00000a090300004100000000030c4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c7000200000002001d00040000000c001d2820281b0000040f000000040c000029000000600310027000000a0903300197000000a00030008c000000a00400003900000000040340190000001f0640018f000000e00740019000000000057c0019000024d50000613d000000000801034f00000000090c0019000000008a08043c0000000009a90436000000000059004b000024d10000c13d000000000006004b000024e20000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00000001002001900000259f0000613d0000000506000029000000020200002900000001050000290000001f01400039000001e00110018f000000000bc1001900000000001b004b0000000001000039000000010100403900000a0f00b0009c000025820000213d0000000100100190000025820000c13d0000004000b0043f000000a00030008c000025880000413d00000000010c043300000a8f0010009c000025880000213d0000008001c00039000000000101043300000a8f0010009c000025880000213d000000000805043300000a110080009c000025910000813d0000006001c00039000000000701043300000a900100004100000000001b04350000000001000414000000040020008c000025080000c13d0000002004000039000025380000013d000100000008001d000200000007001d00000a0900b0009c00000a090300004100000000030b4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c700040000000b001d2820281b0000040f000000040b000029000000600310027000000a0903300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b0019000025250000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b000025210000c13d000000000006004b000025320000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000025ab0000613d0000000506000029000000020700002900000001080000290000001f01400039000000600110018f0000000001b1001900000a0f0010009c000025820000213d000000400010043f000000200030008c000025880000413d00000000020b0433000000ff0020008c000025880000213d0000000003060433000000ff0330018f0000000003230019000000ff0030008c0000258b0000213d000000240230008c0000254d0000213d0000255c0000c13d0000000102000039000025660000013d0000004d0020008c0000258b0000213d0000000a040000390000000103000039000000010020019000000000054400a9000000010400603900000000033400a900000001022002720000000004050019000025510000c13d000000000003004b000025c90000613d00000000033800d90000256c0000013d0000000a0400003900000001020000390000002403300089000000010030019000000000054400a9000000010400603900000000022400a9000000010330027200000000040500190000255f0000c13d000000000008004b0000257e0000613d00000000038200a900000000048300d9000000000024004b0000258b0000c13d00000a160030009c000025980000213d00000a130010009c0000000605000029000025820000213d0000004002100039000000400020043f00000a09027001970000002004100039000000000024043500000000003104350000000303000029000000000303043300000a0903300197000000000032004b00000000050180190000000001050019000000000001042d000000000300001900000a130010009c0000000605000029000025710000a13d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000000010000190000282200010430000000000001042f00000aaa01000041000000000010043f0000001101000039000000040010043f00000a3d01000041000028220001043000000a910100004100000000001b043500000a0900b0009c00000a090b0080410000004001b0021000000a4e011001c7000028220001043000000a9102000041000000000021043500000a090010009c00000a0901008041000000400110021000000a4e011001c700002822000104300000001f0530018f00000a0b06300198000000400200043d0000000004620019000025b60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000025a60000c13d000025b60000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000025b60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000025b20000c13d000000000005004b000025c30000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f0000000000140435000000600130021000000a090020009c00000a09020080410000004002200210000000000112019f000028220001043000000aaa01000041000000000010043f0000001201000039000000040010043f00000a3d0100004100002822000104300007000000000002000000400200043d00000ac00020009c000027280000813d0000004003200039000000400030043f00000020032000390000000000030435000000000002043500000a0e01100197000400000001001d000000000010043f0000000601000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f00000001002001900000272e0000613d000000400300043d00000a130030009c000027280000213d000000000101043b0000004002300039000000400020043f000000000101041a00000a1602100197000700000003001d0000000002230436000000e001100270000500000002001d000600000001001d000000000012043500000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f0000000100200190000027300000613d000000000101043b00060006001000740000273c0000413d00000a8b0100004100000000001004430000000001000412000000040010044300000040010000390000002400100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a8c011001c700008005020000392820281b0000040f0000000100200190000027300000613d000000000101043b00000a0901100197000000060010006b0000261b0000813d00000007050000290000000501000029000000000101043300000a0900100198000027200000c13d000027310000013d0000000401000029000000000010043f0000000701000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f000000010020019000000007050000290000272e0000613d000000400200043d00000a0c0020009c000027280000213d000000000101043b0000006003200039000000400030043f000000000301041a000000400120003900000a8d003001980000000004000039000000010400c039000000000041043500000a0e013001970000000006120436000000a003300270000000ff0330018f0000000000360435000026530000613d000000000001004b000026530000613d000000400100043d00000a130010009c000027280000213d0000004003100039000000400030043f0000002003100039000000000003043500000000000104350000000002020433000000400c00043d00000a8e0100004100000000051c0436000000000100041400000a0e02200197000000040020008c000600000006001d000026580000c13d0000000103000031000000a00030008c000000a0040000390000000004034019000026880000013d0000000501000029000000000101043300000a0900100198000027200000c13d000027310000013d000100000005001d00000a0900c0009c00000a090300004100000000030c4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c7000200000002001d00030000000c001d2820281b0000040f000000030c000029000000600310027000000a0903300197000000a00030008c000000a00400003900000000040340190000001f0640018f000000e00740019000000000057c0019000026750000613d000000000801034f00000000090c0019000000008a08043c0000000009a90436000000000059004b000026710000c13d000000000006004b000026820000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000027500000613d0000000606000029000000020200002900000001050000290000001f01400039000001e00110018f000000000bc1001900000000001b004b0000000001000039000000010100403900000a0f00b0009c000027280000213d0000000100100190000027280000c13d0000004000b0043f000000a00030008c0000272e0000413d00000000010c043300000a8f0010009c0000272e0000213d0000008001c00039000000000101043300000a8f0010009c0000272e0000213d000000000805043300000a110080009c000027420000813d0000006001c00039000000000701043300000a900100004100000000001b04350000000001000414000000040020008c000026a80000c13d0000002004000039000026d80000013d000100000008001d000200000007001d00000a0900b0009c00000a090300004100000000030b4019000000400330021000000a090010009c00000a0901008041000000c001100210000000000131019f00000a4e011001c700030000000b001d2820281b0000040f000000030b000029000000600310027000000a0903300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b0019000026c50000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b000026c10000c13d000000000006004b000026d20000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00000001002001900000275c0000613d0000000606000029000000020700002900000001080000290000001f01400039000000600110018f0000000001b1001900000a0f0010009c000027280000213d000000400010043f000000200030008c0000272e0000413d00000000020b0433000000ff0020008c0000272e0000213d0000000003060433000000ff0330018f0000000003230019000000ff0030008c0000273c0000213d000000240230008c000026ed0000213d000026fc0000c13d0000000102000039000027060000013d0000004d0020008c0000273c0000213d0000000a040000390000000103000039000000010020019000000000054400a9000000010400603900000000033400a900000001022002720000000004050019000026f10000c13d000000000003004b0000277a0000613d00000000033800d90000270c0000013d0000000a0400003900000001020000390000002403300089000000010030019000000000054400a9000000010400603900000000022400a900000001033002720000000004050019000026ff0000c13d000000000008004b000027240000613d00000000038200a900000000048300d9000000000024004b0000273c0000c13d00000a160030009c000027490000213d00000a130010009c0000000705000029000027280000213d0000004002100039000000400020043f00000a09027001970000002004100039000000000024043500000000003104350000000503000029000000000303043300000a0903300197000000000032004b00000000050180190000002001500039000000000101043300000a0900100198000027310000613d000000000105043300000a1601100198000027310000613d000000000001042d000000000300001900000a130010009c0000000705000029000027110000a13d00000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d01000041000028220001043000000000010000190000282200010430000000000001042f000000400100043d00000a9202000041000000000021043500000004021000390000000403000029000000000032043500000a090010009c00000a0901008041000000400110021000000a3d011001c7000028220001043000000aaa01000041000000000010043f0000001101000039000000040010043f00000a3d01000041000028220001043000000a910100004100000000001b043500000a0900b0009c00000a090b0080410000004001b0021000000a4e011001c7000028220001043000000a9102000041000000000021043500000a090010009c00000a0901008041000000400110021000000a4e011001c700002822000104300000001f0530018f00000a0b06300198000000400200043d0000000004620019000027670000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000027570000c13d000027670000013d0000001f0530018f00000a0b06300198000000400200043d0000000004620019000027670000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000027630000c13d000000000005004b000027740000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f0000000000140435000000600130021000000a090020009c00000a09020080410000004002200210000000000112019f000028220001043000000aaa01000041000000000010043f0000001201000039000000040010043f00000a3d01000041000028220001043000000aa3055001970000006006100039000000000056043500000a22044001970000004005100039000000000045043500000a0e033001970000002004100039000000000034043500000a0e0220019700000000002104350000008001100039000000000001042d0004000000000002000400000002001d00000a0f01100197000100000001001d000000000010043f0000000501000039000000200010043f000000000100041400000a090010009c00000a0901008041000000c00110021000000a19011001c700008010020000392820281b0000040f0000000100200190000027c40000613d000000400300043d00000ac00030009c000027c60000813d000000000101043b0000004002300039000000400020043f000000000101041a0000002002300039000000e004100270000000000042043500000a16011001970000000000130435000000040200002900000a0902200198000027c30000613d000300000004001d000400000002001d000200000003001d00000a84010000410000000000100443000000000100041400000a090010009c00000a0901008041000000c00110021000000a85011001c70000800b020000392820281b0000040f0000000100200190000027cc0000613d000000000101043b000000030110006c0000000404000029000027cd0000413d000000000041004b0000000202000029000027d30000213d000000000102043300000a1601100197000000000001042d0000000001000019000028220001043000000aaa01000041000000000010043f0000004101000039000000040010043f00000a3d010000410000282200010430000000000001042f00000aaa01000041000000000010043f0000001101000039000000040010043f00000a3d010000410000282200010430000000400200043d000000440320003900000000001304350000002401200039000000000041043500000a8601000041000000000012043500000004012000390000000103000029000000000031043500000a090020009c00000a0902008041000000400120021000000a87011001c70000282200010430000000000001042f00000a090010009c00000a0901008041000000400110021000000a090020009c00000a09020080410000006002200210000000000112019f000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000a20011001c700008010020000392820281b0000040f0000000100200190000027f60000613d000000000101043b000000000001042d0000000001000019000028220001043000000000050100190000000000200443000000050030008c000028060000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b000027fe0000413d00000a090030009c00000a09030080410000006001300210000000000200041400000a090020009c00000a0902008041000000c002200210000000000112019f00000ac3011001c700000000020500192820281b0000040f0000000100200190000028150000613d000000000101043b000000000001042d000000000001042f00002819002104210000000102000039000000000001042d0000000002000019000000000001042d0000281e002104230000000102000039000000000001042d0000000002000019000000000001042d0000282000000432000028210001042e0000282200010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000000000000000000000000000ffffffffffffff9f0000000000000000000000000000000000000000ffffffffffffffffffffffff000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe080000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffbf000000000000000000000000000000000000000000000000ffffffffffffff3f000000000000000000000000000000000000000000000000fffffffffffffddf00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffdf0200000000000000000000000000000000000040000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a5330200000000000000000000000000000000000020000000000000000000000000c3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda77580bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a532eb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdeffe8a4859c7bd88fc0f24184464406785daae8e84cb1860cc4a4eff72e05fe2470200000000000000000000000000000000000000000000000000000000000000df1b1bd32a69711488d71554706bb130b1fc63a5fa1a2cd85e8440f84065ba23ffffffffffffffffffff000000000000000000000000000000000000000000000000000000000000000000ff000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000200000000000000000000000000000000000060000000000000000000000000e6a7a17d710bf0b2cd05e5397dc6f97a5da4ee79e31e234bf5f965ee2bd9a5bfffffffff000000000000000000000000000000000000000000000000000000002812d52c0000000000000000000000000000000000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000283b699f411baff8f1c29fe49f32a828c8151596244b8e7e4c164edd6569a835525e3d4e0c31cef19cf9426af8d2c0ddd2d576359ca26bed92aac5fadda46265ff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffff0000000000000000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000000000ffffffff0000000000000000000000000000000000000000000000000000ffff0000000000000000000000000000000000000000000000000000ffffffff0000000000000000000000000000000000000000000000000000ffff00000000000000000000000000000000000000000000000000000000ffff00000000000000000000000000000000000000000000000000000000ffff0000000000000000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000000000000000ffffffffffffff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000c35aa79d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000bb77da6f7210cdd16904228a9360133d1d7dfff99b1bc75f128da5b53e28f97d0b4f67a2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000ffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000ffff0000000000000000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000000000ffffffff0000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000002000000000000000000000000000000000000c000000000000000000000000094967ae9ea7729ad4f54021c1981765d2b1d954f7c92fbec340aa0a54f46b8b524ecdc02000000000000000000000000000000000000000000000000000000004de5b1bcbca6018c11303a2c3f4a4b4f22a1c741d8c4ba430d246ac06c5ddf8b0000000200000000000000000000000000000100000001000000000000000000d794ef950000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000009b15e16f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000079ba509600000000000000000000000000000000000000000000000000000000bf78e03e00000000000000000000000000000000000000000000000000000000d8694ccc00000000000000000000000000000000000000000000000000000000fbe3f77700000000000000000000000000000000000000000000000000000000fbe3f77800000000000000000000000000000000000000000000000000000000ffdb4b3700000000000000000000000000000000000000000000000000000000d8694ccd00000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000d026419f00000000000000000000000000000000000000000000000000000000d02641a000000000000000000000000000000000000000000000000000000000d63d3af200000000000000000000000000000000000000000000000000000000bf78e03f00000000000000000000000000000000000000000000000000000000cdc73d510000000000000000000000000000000000000000000000000000000082b49eaf0000000000000000000000000000000000000000000000000000000091a274990000000000000000000000000000000000000000000000000000000091a2749a00000000000000000000000000000000000000000000000000000000a69c64c00000000000000000000000000000000000000000000000000000000082b49eb0000000000000000000000000000000000000000000000000000000008da5cb5b0000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000007afac32200000000000000000000000000000000000000000000000000000000805f21320000000000000000000000000000000000000000000000000000000041ed29e600000000000000000000000000000000000000000000000000000000514e8cfe000000000000000000000000000000000000000000000000000000006def4ce6000000000000000000000000000000000000000000000000000000006def4ce700000000000000000000000000000000000000000000000000000000770e2dc400000000000000000000000000000000000000000000000000000000514e8cff000000000000000000000000000000000000000000000000000000006cb5f3dd0000000000000000000000000000000000000000000000000000000045ac924c0000000000000000000000000000000000000000000000000000000045ac924d000000000000000000000000000000000000000000000000000000004ab35b0b0000000000000000000000000000000000000000000000000000000041ed29e700000000000000000000000000000000000000000000000000000000430d138c00000000000000000000000000000000000000000000000000000000181f5a7600000000000000000000000000000000000000000000000000000000325c868d00000000000000000000000000000000000000000000000000000000325c868e000000000000000000000000000000000000000000000000000000003937306f00000000000000000000000000000000000000000000000000000000181f5a77000000000000000000000000000000000000000000000000000000002451a627000000000000000000000000000000000000000000000000000000000041e5be00000000000000000000000000000000000000000000000000000000061877e30000000000000000000000000000000000000000000000000000000006285c69000000000000000000000000000000000000004000000000000000000000000099ac52f200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff7f2b5c74de000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000008d666f6000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffffffffffffffffffffffffffc00796b89b91644bc98cd93958e4c9038275d622183e25ac5af08cc6b5d955391320200000200000000000000000000000000000004000000000000000000000000f08bcb3e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000ff000000000000000000000000000000000000000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000ffffffdf310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e020000020000000000000000000000000000004400000000000000000000000000000000000000000000ff000000000000000000000000000000000000000000feaf968c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffff313ce5670000000000000000000000000000000000000000000000000000000010cb51d10000000000000000000000000000000000000000000000000000000006439c6b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000000186a000000000000000000000000000000000000000000000000000005af3107a400000000000000000000000000000000000ffffffffffffffffffffffffffffffff181dcf100000000000000000000000000000000000000000000000000000000097a657c900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffff0000000000000000000000000000000000000020000000000000000000000000ee433e99000000000000000000000000000000000000000000000000000000004c4fc93a00000000000000000000000000000000000000000000000000000000d88dddd60000000000000000000000000000000000000000000000000000000086933789000000000000000000000000000000000000000000000000000000002502348c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000008000000000000000000175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db98579befe00000000000000000000000000000000000000000000000000000000ffff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5f0000000000000000000000fe00000000000000000000000000000000000000000000000000000000000000ed000000000000000000000000000000000000000052f50aa6d1a95a4595361ecf953d095f125d442e4673716dede699e049de148a097e17ff00000000000000000000000000000000000000000000000000000000fe8a4859c7bd88fc0f24184464406785daae8e84cb1860cc4a4eff72e05fe2484e487b71000000000000000000000000000000000000000000000000000000001795838dc8ab2ffc5f431a1729a6afa0b587f982f7b2be0b9d7187a1ef547f9102b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e002000000000000000000000000000000000002200000000000000000000000005247fdce00000000000000000000000000000000000000000000000000000000000000000000002812d52c00000000000000000000000000000000000000000036f536ca000000000000000000000000000000000000000000000000000000006a92a4830000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff02000000000000000000000000000000000000a000000000000000000000000032a4ba3fa3351b11ad555d4c8ec70a744e8705607077a946807030d64b6ab1a3dd84a3fa9ef9409f550d54d6affec7e9c480c878c6ab27b78912a03e1b371c6ed86ad9cf00000000000000000000000000000000000000000000000000000000405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace46656551756f74657220312e362e302d646576000000000000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000060000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000ffffffffffffffa0000000000000000000000000000000000000000000000000ffffffffffffffc0000000000000000000000000000000000000000000000000ffffffffffffff40000000000000000000000000000000000000000000000000fffffffffffffde002000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
