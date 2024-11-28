package ethereum

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

var DualAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"minAnswer_\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"maxAnswer_\",\"type\":\"int192\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"billingAccessController\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"requesterAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description_\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"secondaryProxy_\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"cutoffTime_\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxSyncIterations_\",\"type\":\"uint32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CalldataLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferPayeeToSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FMustBePositive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FaultyOracleFTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientGas\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidOnChainConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeftGasCannotExceedInitialGas\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxSyncIterations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MedianIsOutOfMinMaxRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NumObservationsOutOfBounds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByEOA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCurrentPayeeCanUpdate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyOwnerAndBillingAdminCanCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyOwnerAndRequesterCanCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPayeeCanWithdraw\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyProposedPayeesCanAccept\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PayeeAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSignerAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitterAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReportLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RoundNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooFewValuesToTrustMedian\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferRemainingFundsFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransmittersSizeNotEqualPayeeSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"cutoffTime\",\"type\":\"uint32\"}],\"name\":\"CutoffTimeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"oldLinkToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"newLinkToken\",\"type\":\"address\"}],\"name\":\"LinkTokenSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxSyncIterations\",\"type\":\"uint32\"}],\"name\":\"MaxSyncIterationsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationsTimestamp\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"primaryRoundId\",\"type\":\"uint32\"}],\"name\":\"PrimaryFeedUnlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"secondaryRoundId\",\"type\":\"uint32\"}],\"name\":\"SecondaryRoundIdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"ValidatorConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId_\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorConfig\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer_\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp_\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_cutoffTime\",\"type\":\"uint32\"}],\"name\":\"setCutoffTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maxSyncIterations\",\"type\":\"uint32\"}],\"name\":\"setMaxSyncIterations\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"newValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newGasLimit\",\"type\":\"uint32\"}],\"name\":\"setValidatorConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmitSecondary\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b50604051620064ad380380620064ad8339810160408190526200003591620005a1565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001a2565b50506001805460ff60a01b1916600160a01b1790555060ff851660e052601789810b60805288900b60a0526001600160a01b0380841660c05260148054918c166001600160a01b0319909216821790556040516000907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a908290a362000145876200024d565b6200015086620002c6565b6200015d60008062000341565b6012805463ffffffff838116640100000000026001600160401b03199092169085161717905560136200019185826200072b565b5050505050505050505050620007f7565b336001600160a01b03821603620001fc5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6015546001600160a01b039081169082168114620002c257601580546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291015b60405180910390a15b5050565b620002d062000428565b6010546001600160a01b039081169082168114620002c257601080546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349101620002b9565b6200034b62000428565b60408051808201909152600f546001600160a01b03808216808452600160a01b90920463ffffffff16602084015284161415806200039957508163ffffffff16816020015163ffffffff1614155b1562000423576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600f80546001600160c01b0319168417600160a01b830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6000546001600160a01b03163314620004845760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000083565b565b6001600160a01b03811681146200049c57600080fd5b50565b8051601781900b8114620004b257600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620004df57600080fd5b81516001600160401b0380821115620004fc57620004fc620004b7565b604051601f8301601f19908116603f01168101908282118183101715620005275762000527620004b7565b81604052838152602092508660208588010111156200054557600080fd5b600091505b838210156200056957858201830151818301840152908201906200054a565b6000602085830101528094505050505092915050565b8051620004b28162000486565b805163ffffffff81168114620004b257600080fd5b6000806000806000806000806000806101408b8d031215620005c257600080fd5b8a51620005cf8162000486565b9950620005df60208c016200049f565b9850620005ef60408c016200049f565b975060608b0151620006018162000486565b60808c0151909750620006148162000486565b60a08c015190965060ff811681146200062c57600080fd5b60c08c01519095506001600160401b038111156200064957600080fd5b620006578d828e01620004cd565b9450506200066860e08c016200057f565b9250620006796101008c016200058c565b91506200068a6101208c016200058c565b90509295989b9194979a5092959850565b600181811c90821680620006b057607f821691505b602082108103620006d157634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111562000423576000816000526020600020601f850160051c81016020861015620007025750805b601f850160051c820191505b8181101562000723578281556001016200070e565b505050505050565b81516001600160401b03811115620007475762000747620004b7565b6200075f816200075884546200069b565b84620006d7565b602080601f8311600181146200079757600084156200077e5750858301515b600019600386901b1c1916600185901b17855562000723565b600085815260208120601f198616915b82811015620007c857888601518255948401946001909101908401620007a7565b5085821015620007e75787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c05160e051615c606200084d6000396000610493015260006132c5015260008181610534015281816121af01526147300152600081816103c90152818161218701526147050152615c606000f3fe608060405234801561001057600080fd5b50600436106103365760003560e01c80639bd2c0b1116101b2578063c4c92b37116100f9578063e5fe4577116100a2578063eb5dcd6c1161007c578063eb5dcd6c146108ea578063f2fde38b146108fd578063fbffd2c114610910578063feaf968c1461092357600080fd5b8063e5fe45771461086f578063e76d5168146108b9578063eb457163146108d757600080fd5b8063dc7f0124116100d3578063dc7f01241461080f578063e3d0e71214610834578063e4902f821461084757600080fd5b8063c4c92b37146107cb578063d09dc339146107e9578063daffc4b5146107f157600080fd5b8063b17f2a6b1161015b578063b633620c11610135578063b633620c14610792578063ba0cb29e146107a5578063c1075329146107b857600080fd5b8063b17f2a6b14610759578063b1dc65a41461076c578063b5ab58dc1461077f57600080fd5b8063a118f2491161018c578063a118f24914610702578063afcb95d714610715578063b121e1471461074657600080fd5b80639bd2c0b11461067d5780639c849b30146106dc5780639e3ceeab146106ef57600080fd5b8063668a0f021161028157806381ff70481161022a5780638ac28d5a116102045780638ac28d5a146105be5780638da5cb5b146105d157806398e5b12a146106105780639a6fc8f51461063357600080fd5b806381ff7048146105735780638205bf6a146105a35780638823da6c146105ab57600080fd5b80637284e4161161025b5780637284e4161461055b57806379ba5097146105635780638038e4a11461056b57600080fd5b8063668a0f02146105075780636b14daf81461050f57806370da2f671461053257600080fd5b8063313ce567116102e357806354fd4d50116102bd57806354fd4d50146104d8578063643dc105146104df578063666cab8d146104f257600080fd5b8063313ce5671461048c5780634fb17470146104bd57806350d25bcd146104d057600080fd5b8063181f5a7711610314578063181f5a771461037e57806322adbc78146103c757806329937268146103f057600080fd5b8063055aae661461033b5780630a756983146103505780630eafb25b14610358575b600080fd5b61034e610349366004614f9e565b61092b565b005b61034e6109a7565b61036b610366366004614fdd565b610a26565b6040519081526020015b60405180910390f35b6103ba6040518060400160405280601481526020017f4475616c41676772656761746f7220312e302e3000000000000000000000000081525081565b604051610375919061505e565b7f000000000000000000000000000000000000000000000000000000000000000060170b61036b565b600d54600c546040805163ffffffff6e010000000000000000000000000000850481168252720100000000000000000000000000000000000085048116602083015276010000000000000000000000000000000000000000000085048116928201929092527a01000000000000000000000000000000000000000000000000000090930416606083015262ffffff16608082015260a001610375565b60405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610375565b61034e6104cb366004615071565b610b58565b61036b610dfe565b600661036b565b61034e6104ed3660046150aa565b610e2b565b6104fa6110c3565b6040516103759190615175565b61036b611132565b61052261051d366004615294565b611147565b6040519015158152602001610375565b7f000000000000000000000000000000000000000000000000000000000000000060170b61036b565b6103ba61117c565b61034e611205565b61034e611307565b600e54600b546040805163ffffffff80851682526401000000009094049093166020840152820152606001610375565b61036b61139b565b61034e6105b9366004614fdd565b6113e8565b61034e6105cc366004614fdd565b61149d565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610375565b610618611506565b60405169ffffffffffffffffffff9091168152602001610375565b6106466106413660046152e4565b611686565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a001610375565b604080518082018252600f5473ffffffffffffffffffffffffffffffffffffffff81168083527401000000000000000000000000000000000000000090910463ffffffff16602092830181905283519182529181019190915201610375565b61034e6106ea36600461535c565b61174f565b61034e6106fd366004614fdd565b61196f565b61034e610710366004614fdd565b611a20565b600b54600d546040805160008152602081019390935261010090910460081c63ffffffff1690820152606001610375565b61034e610754366004614fdd565b611ad4565b61034e610767366004614f9e565b611bcc565b61034e61077a3660046153c8565b611c37565b61036b61078d3660046154ad565b611c53565b61036b6107a03660046154ad565b611c8e565b61034e6107b33660046153c8565b611ce9565b61034e6107c63660046154c6565b611cfb565b60155473ffffffffffffffffffffffffffffffffffffffff166105eb565b61036b611fd3565b60105473ffffffffffffffffffffffffffffffffffffffff166105eb565b6001546105229074010000000000000000000000000000000000000000900460ff1681565b61034e6108423660046155bc565b612089565b61085a610855366004614fdd565b6128c2565b60405163ffffffff9091168152602001610375565b61087761298d565b6040805195865263ffffffff909416602086015260ff9092169284019290925260179190910b606083015267ffffffffffffffff16608082015260a001610375565b60145473ffffffffffffffffffffffffffffffffffffffff166105eb565b61034e6108e5366004615689565b612a3c565b61034e6108f8366004615071565b612b73565b61034e61090b366004614fdd565b612ccb565b61034e61091e366004614fdd565b612cdc565b610646612ced565b610933612d89565b601280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff1664010000000063ffffffff8416908102919091179091556040519081527fcba51f727ba38740aa888ce0cb33f68de587733f61d3fafa0d9fb2b29e7f829f906020015b60405180910390a150565b6109af612d89565b60015474010000000000000000000000000000000000000000900460ff1615610a2457600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b73ffffffffffffffffffffffffffffffffffffffff811660009081526003602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff169181019190915290610a9a5750600092915050565b600d546020820151600091760100000000000000000000000000000000000000000000900463ffffffff169060079060ff16601f8110610adc57610adc6156b7565b600881049190910154600d54610b12926007166004026101000a90910463ffffffff908116916601000000000000900416615715565b63ffffffff16610b229190615732565b610b3090633b9aca00615732565b905081604001516bffffffffffffffffffffffff1681610b509190615749565b949350505050565b610b60612d89565b60145473ffffffffffffffffffffffffffffffffffffffff908116908316819003610b8a57505050565b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff8416906370a0823190602401602060405180830381865afa158015610bf4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c18919061575c565b50610c21612e0a565b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009073ffffffffffffffffffffffffffffffffffffffff8316906370a0823190602401602060405180830381865afa158015610c8e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cb2919061575c565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152602482018390529192509083169063a9059cbb906044016020604051808303816000875af1158015610d2b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d4f9190615775565b610d85576040517f7725087a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff86811691821790925560405190918416907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a90600090a350505b5050565b600060116000610e0c61325e565b63ffffffff16815260208101919091526040016000205460170b919050565b60005473ffffffffffffffffffffffffffffffffffffffff16331480610ee857506015546040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690636b14daf890610ea79033906000903690600401615797565b602060405180830381865afa158015610ec4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ee89190615775565b610f1e576040517f91ed77c500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610f26612e0a565b600d80547fffffffffffffffffffff0000000000000000ffffffffffffffffffffffffffff166e01000000000000000000000000000063ffffffff8881169182027fffffffffffffffffffff00000000ffffffffffffffffffffffffffffffffffff16929092177201000000000000000000000000000000000000888416908102919091177fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff167601000000000000000000000000000000000000000000008885169081027fffff00000000ffffffffffffffffffffffffffffffffffffffffffffffffffff16919091177a01000000000000000000000000000000000000000000000000000094881694850217909455600c80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000001662ffffff871690811790915560408051938452602084019290925290820193909352606081019190915260808101919091527f0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f9060a00160405180910390a15050505050565b6060600680548060200260200160405190810160405280929190818152602001828054801561112857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116110fd575b5050505050905090565b600061113c61325e565b63ffffffff16905090565b600061115383836133c6565b80611173575073ffffffffffffffffffffffffffffffffffffffff831632145b90505b92915050565b60606013805461118b90615801565b80601f01602080910402602001604051908101604052809291908181526020018280546111b790615801565b80156111285780601f106111d957610100808354040283529160200191611128565b820191906000526020600020905b8154815290600101906020018083116111e757509395945050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461128b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61130f612d89565b60015474010000000000000000000000000000000000000000900460ff16610a2457600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b6000601160006113a961325e565b63ffffffff90811682526020820192909252604001600020547c0100000000000000000000000000000000000000000000000000000000900416919050565b6113f0612d89565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff161561149a5773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905590519182527f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1910161099c565b50565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152601660205260409020541633146114fd576040517f2ab4a3db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61149a8161341b565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906115c857506010546040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690636b14daf8906115859033906000903690600401615797565b602060405180830381865afa1580156115a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115c69190615775565b155b156115ff576040517f4cdc445800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d54600b546040805191825263ffffffff6101008404600881901c8216602085015260ff811684840152915164ffffffffff9092169366010000000000009004169133917f41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c9181900360600190a261167981600161584e565b63ffffffff169250505090565b600080600080600061169661325e565b63ffffffff168669ffffffffffffffffffff1611156116c357506000935083925082915081905080611746565b5050505063ffffffff82811660009081526011602090815260409182902082516060810184529054601781900b8083527801000000000000000000000000000000000000000000000000820486169383018490527c0100000000000000000000000000000000000000000000000000000000909104909416920182905284935090835b91939590929450565b611757612d89565b828114611790576040517f3d2f942900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b838110156119685760008585838181106117af576117af6156b7565b90506020020160208101906117c49190614fdd565b905060008484848181106117da576117da6156b7565b90506020020160208101906117ef9190614fdd565b73ffffffffffffffffffffffffffffffffffffffff80841660009081526016602052604090205491925016801580158161185557508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b1561188c576040517faeae062800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff848116600090815260166020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614611959578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b50505050806001019050611793565b5050505050565b611977612d89565b60105473ffffffffffffffffffffffffffffffffffffffff9081169082168114610dfa57601080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491015b60405180910390a15050565b611a28612d89565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff1661149a5773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905590519182527f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4910161099c565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260176020526040902054163314611b34576040517f6599cbbe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526016602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556017909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b611bd4612d89565b601280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff83169081179091556040519081527fb24a681ce3399a408a89fd0c2b59dfc24bdad592b1c7ec7671cf060596c1c4d19060200161099c565b611c4988888888888888886000613678565b5050505050505050565b6000611c5d61325e565b63ffffffff16821115611c7257506000919050565b5063ffffffff1660009081526011602052604090205460170b90565b6000611c9861325e565b63ffffffff16821115611cad57506000919050565b5063ffffffff9081166000908152601160205260409020547c010000000000000000000000000000000000000000000000000000000090041690565b611c4988888888888888886001613678565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590611dbc57506015546040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690636b14daf890611d799033906000903690600401615797565b602060405180830381865afa158015611d96573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611dba9190615775565b155b15611df3576040517f91ed77c500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000611dfd6139fc565b6014546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015291925060009173ffffffffffffffffffffffffffffffffffffffff909116906370a0823190602401602060405180830381865afa158015611e71573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e95919061575c565b905081811015611ed1576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60145473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85611f03611efd868661586b565b87613be7565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff909216600483015260248201526044016020604051808303816000875af1158015611f73573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f979190615775565b611fcd576040517f356680b700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050565b6014546040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152600091829173ffffffffffffffffffffffffffffffffffffffff909116906370a0823190602401602060405180830381865afa158015612046573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061206a919061575c565b905060006120766139fc565b9050612082818361587e565b9250505090565b612091612d89565b601f865111156120cd576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8451865114612108576040517f250a65b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b855161211585600361589e565b60ff161061214f576040517f20c9729a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61215b8460ff16613bfe565b604080517f010000000000000000000000000000000000000000000000000000000000000060208201527f0000000000000000000000000000000000000000000000000000000000000000821b60218201527f000000000000000000000000000000000000000000000000000000000000000090911b60398201526051016040516020818303038152906040528051906020012083805190602001201461222e576040517fa8811dc600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160c08101825267ffffffffffffffff8416815260ff86166020820152908101849052606081018290526080810187905260a08101869052600d80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000ff169055612299612e0a565b60055460005b81811015612388576000600582815481106122bc576122bc6156b7565b60009182526020822001546006805473ffffffffffffffffffffffffffffffffffffffff909216935090849081106122f6576122f66156b7565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff948516835260048252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905594168252600390529190912080547fffffffffffffffffffffffffffffffffffff00000000000000000000000000001690555060010161229f565b5061239560056000614e3c565b6123a160066000614e3c565b60005b8260800151518110156126b75760046000846080015183815181106123cb576123cb6156b7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615612436576040517f16c6131500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600460008560800151848151811061246b5761246b6156b7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281810192909252604001600090812083518154949093015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff931515939093167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093179190911790915560a08401518051600392919084908110612521576125216156b7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff161561258c576040517fd63d347400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180606001604052806001151581526020018260ff16815260200160006bffffffffffffffffffffffff16815250600360008560a0015184815181106125d6576125d66156b7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600020835181549385015194909201516bffffffffffffffffffffffff1662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff931515939093167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009094169390931791909117929092161790556001016123a4565b50608082015180516126d191600591602090910190614e5a565b5060a082015180516126eb91600691602090910190614e5a565b506020820151600d80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600e80547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff43811682029283178555908304811693600193909260009261277d92869290821691161761584e565b92506101000a81548163ffffffff021916908363ffffffff1602179055506127dc4630600e60009054906101000a900463ffffffff1663ffffffff1686608001518760a00151886020015189604001518a600001518b60600151613c38565b600b819055600e54608085015160a08601516020870151604080890151895160608b015192517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0598612842988b98919763ffffffff9091169691959094919391926158ba565b60405180910390a1600d546601000000000000900463ffffffff1660005b8460800151518110156128b55781600782601f8110612881576128816156b7565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff160217905550806001019050612860565b5050505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526003602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff1691810191909152906129365750600092915050565b6007816020015160ff16601f8110612950576129506156b7565b600881049190910154600d54612986926007166004026101000a90910463ffffffff908116916601000000000000900416615715565b9392505050565b6000808080803332146129cc576040517f74e2cd5100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050600b54600d5463ffffffff6601000000000000820481166000908152601160205260409020549296610100909204600881901c8216965064ffffffffff169450601783900b93507c010000000000000000000000000000000000000000000000000000000090920490911690565b612a44612d89565b60408051808201909152600f5473ffffffffffffffffffffffffffffffffffffffff8082168084527401000000000000000000000000000000000000000090920463ffffffff1660208401528416141580612aaf57508163ffffffff16816020015163ffffffff1614155b15612b6e5760408051808201825273ffffffffffffffffffffffffffffffffffffffff85811680835263ffffffff8681166020948501819052600f80547fffffffffffffffff00000000000000000000000000000000000000000000000016841774010000000000000000000000000000000000000000830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b73ffffffffffffffffffffffffffffffffffffffff828116600090815260166020526040902054163314612bd3576040517fb97d016a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81163303612c22576040517f79df0c6600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff808316600090815260176020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114612b6e5760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b612cd3612d89565b61149a81613ce4565b612ce4612d89565b61149a81613dd9565b600080600080600080612cfe61325e565b63ffffffff90811660008181526011602090815260409182902082516060810184529054601781900b8083527801000000000000000000000000000000000000000000000000820487169383018490527c0100000000000000000000000000000000000000000000000000000000909104909516920182905291999298509096509450879350915050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a24576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611282565b601454600d54604080516103e081019182905273ffffffffffffffffffffffffffffffffffffffff90931692660100000000000090920463ffffffff1691600091600790601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411612e555790505050505050905060006006805480602002602001604051908101604052809291908181526020018280548015612efd57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612ed2575b5050505050905060005b815181101561325057600060036000848481518110612f2857612f286156b7565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160029054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff169050600060036000858581518110612fae57612fae6156b7565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160026101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008483601f8110613035576130356156b7565b6020020151600d5490870363ffffffff9081169250760100000000000000000000000000000000000000000000909104168102633b9aca0002820180156132455760006016600087878151811061308e5761308e6156b7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff90811683529082019290925260409081016000205490517fa9059cbb00000000000000000000000000000000000000000000000000000000815290821660048201819052602482018590529250908a169063a9059cbb906044016020604051808303816000875af1158015613128573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061314c9190615775565b613182576040517f356680b700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b878786601f8110613195576131956156b7565b602002019063ffffffff16908163ffffffff16815250508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff168787815181106131ec576131ec6156b7565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c8560405161323b91815260200190565b60405180910390a4505b505050600101612f07565b50611968600783601f614ee4565b600d5460009063ffffffff660100000000000082048116916a010000000000000000000081049091169060ff7e010000000000000000000000000000000000000000000000000000000000009091041673ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633036133585760125463ffffffff838116600090815260116020526040902054429261333592908116917c010000000000000000000000000000000000000000000000000000000090041661584e565b63ffffffff16101561335157613349613e7a565b935050505090565b5092915050565b8163ffffffff168363ffffffff16036133be578080156133ae575063ffffffff8084166000908152601160205260409020547c010000000000000000000000000000000000000000000000000000000090041642145b156133be57613349600184615715565b509092915050565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff168061117357505060015474010000000000000000000000000000000000000000900460ff161592915050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600360209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046bffffffffffffffffffffffff169281019290925261348a575050565b600061349583610a26565b90508015612b6e5773ffffffffffffffffffffffffffffffffffffffff838116600090815260166020526040908190205460145491517fa9059cbb000000000000000000000000000000000000000000000000000000008152908316600482018190526024820185905292919091169063a9059cbb906044016020604051808303816000875af115801561352d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135519190615775565b613587576040517f356680b700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d60000160069054906101000a900463ffffffff166007846020015160ff16601f81106135b7576135b76156b7565b6008810491909101805460079092166004026101000a63ffffffff81810219909316939092169190910291909117905573ffffffffffffffffffffffffffffffffffffffff84811660008181526003602090815260409182902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601454915186815291841693851692917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c910160405180910390a450505050565b60005a90506136898a898887613f55565b60006136ca8a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061406a92505050565b6040805161012081018252600d5460ff808216835261010080830464ffffffffff1660208501526601000000000000830463ffffffff908116958501959095526a01000000000000000000008304851660608501526e01000000000000000000000000000083048516608085015272010000000000000000000000000000000000008304851660a08501527601000000000000000000000000000000000000000000008304851660c08501527a010000000000000000000000000000000000000000000000000000830490941660e08401527e010000000000000000000000000000000000000000000000000000000000009091041615159181019190915290915083156138b9576000806137de84614104565b9150915081156138b6578063ffffffff16836060015163ffffffff1610613831576040517ff803a2ca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d80547fffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffff166a010000000000000000000063ffffffff8416908102919091179091556040517f8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b90600090a26138ac838560000151876142a5565b50505050506139f1565b50505b602081810151908d01359064ffffffffff80831691161415806138df5750816101000151155b1561395257816020015164ffffffffff168164ffffffffff161161392f576040517ff803a2ca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61393f8d8d8d8d8d8d8d8d6143e7565b61394d828e3583868961463c565b613991565b600d54604051660100000000000090910463ffffffff16907fda2435684a37fba6f7841e49b59e6ad975e462bbebd28ec9da4ed9746a6992be90600090a25b600d80547fff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e010000000000000000000000000000000000000000000000000000000000008715150217905582516139ec908390866142a5565b505050505b505050505050505050565b6000806006805480602002602001604051908101604052809291908181526020018280548015613a6257602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613a37575b50508351600d54604080516103e08101918290529697509195660100000000000090910463ffffffff169450600093509150600790601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411613a9e5790505050505050905060005b83811015613b27578181601f8110613afe57613afe6156b7565b6020020151613b0d9084615715565b613b1d9063ffffffff1687615749565b9550600101613ae4565b50600d54613b5990760100000000000000000000000000000000000000000000900463ffffffff16633b9aca00615732565b613b639086615732565b945060005b83811015613bdf5760036000868381518110613b8657613b866156b7565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054613bd5906201000090046bffffffffffffffffffffffff1687615749565b9550600101613b68565b505050505090565b600081831015613bf8575081611176565b50919050565b6000811161149a576040517f39d1a4d000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808a8a8a8a8a8a8a8a8a604051602001613c5c99989796959493929190615950565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150505b9998505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603613d63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611282565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60155473ffffffffffffffffffffffffffffffffffffffff9081169082168114610dfa57601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129101611a14565b600d5460125460009163ffffffff6601000000000000909104811691640100000000900416815b63ffffffff811615613f365763ffffffff8216613ebe8285615715565b63ffffffff1614613f365760125463ffffffff8281166000908152601160205260409020544292613f1392908116917c010000000000000000000000000000000000000000000000000000000090041661584e565b63ffffffff161015613f26579392505050565b613f2f816159e5565b9050613ea1565b5050600d546a0100000000000000000000900463ffffffff1692915050565b3360009081526003602052604090205460ff16613f9e576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600b54843514613fda576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b613fe5838383614b40565b600d54613ff69060ff166001615a23565b60ff168214614031576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b808214611fcd576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160808101825260008082526020820152606091810182905281810191909152600080600080858060200190518101906140a79190615a4e565b93509350935093506140b98683614bbe565b815160408051602081018690526000910160408051918152928152825160808101845260179490940b845263ffffffff90961660208401525081019390935260608301525092915050565b600d54601254600091829163ffffffff660100000000000090920482169164010000000090910416815b63ffffffff8116156142975763ffffffff821661414b8285615715565b63ffffffff1603614188576040517f9bc973fd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff80821660009081526011602090815260409182902082516060810184529054601781900b82527801000000000000000000000000000000000000000000000000810485168284018190527c010000000000000000000000000000000000000000000000000000000090910485169382019390935290890151909216111561421d5750600096879650945050505050565b866020015163ffffffff16816020015163ffffffff1614801561427357506060870151805161424e90600290615b17565b8151811061425e5761425e6156b7565b602002602001015160170b816000015160170b145b1561428657506001969095509350505050565b50614290816159e5565b905061412e565b506000958695509350505050565b60008260170b12156142b657505050565b60006142dd633b9aca003a048560a0015163ffffffff16866080015163ffffffff16614c20565b90506010360260005a600c5490915060009061430a9063ffffffff8716908690869062ffffff1686614c46565b90506000670de0b6b3a764000077ffffffffffffffffffffffffffffffffffffffffffffffff881683023360009081526003602052604090205460e08b01519290910492506201000090046bffffffffffffffffffffffff9081169163ffffffff16633b9aca00028284010190811682111561438c5750505050505050505050565b33600090815260036020526040902080546bffffffffffffffffffffffff90921662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff909216919091179055505050505050505050565b600087876040516143f9929190615b52565b604051908190038120614410918b90602001615b62565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828252805160209182012083830190925260008084529083018190529092509060005b878110156145e157600060018587846020811061447d5761447d6156b7565b61448a91901a601b615a23565b8c8c8681811061449c5761449c6156b7565b905060200201358b8b878181106144b5576144b56156b7565b90506020020135604051600081526020016040526040516144f2949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015614514573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526004602090815290849020838501909452925460ff80821615158085526101009092041693830193909352909550925090506145c2576040517fcd2467c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b840193505080600101905061445e565b5081827e0101010101010101010101010101010101010101010101010101010101010116146128b5576040517f8044bb3300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601f826060015151111561467c576040517fff6c220500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b846000015160ff16826060015151116146c1576040517f5765bdd700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b64ffffffffff8316602086015260608201518051600091906146e590600290615b17565b815181106146f5576146f56156b7565b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b138061475857507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b135b1561478f576040517fca191b2600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040860180519061479f82615b78565b63ffffffff90811690915260408051606081018252601785900b815260208781015184168183019081524285168385019081528c85015186166000908152601190935293909120915182549151935185167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff949095167801000000000000000000000000000000000000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090921677ffffffffffffffffffffffffffffffffffffffffffffffff9091161717919091169190911790555081156148da5760408087015163ffffffff166060880181905290517f8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b90600090a25b85600d60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600e6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160000160126101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160000160166101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600001601a6101000a81548163ffffffff021916908363ffffffff16021790555061010082015181600001601e6101000a81548160ff021916908315150217905550905050856040015163ffffffff167fc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a823386602001518760600151886040015189600001518c8c604051614a8d989796959493929190615b9b565b60405180910390a2604080870151602080860151925163ffffffff9384168152600093909216917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271910160405180910390a3856040015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f42604051614b1f91815260200190565b60405180910390a3614b3886604001518260170b614c94565b505050505050565b6000614b4d826020615732565b614b58846020615732565b614b6486610144615749565b614b6e9190615749565b614b789190615749565b614b83906000615749565b9050368114611fcd576040517fb4d895d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081516020614bce9190615732565b614bd99060a0615749565b614be4906000615749565b905080835114612b6e576040517fd4e1416000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008383811015614c3357600285850304015b614c3d8184613be7565b95945050505050565b600081861015614c82576040517ffbf484ab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50633b9aca0094039190910101020290565b60408051808201909152600f5473ffffffffffffffffffffffffffffffffffffffff81168083527401000000000000000000000000000000000000000090910463ffffffff166020830152614ce857505050565b6000614cf5600185615715565b63ffffffff818116600081815260116020526040808220549051602481019390935260170b60448301819052928816606483015260848201879052929350909190614dc89060a401604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052602080820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fbeed9b510000000000000000000000000000000000000000000000000000000017905286519087015163ffffffff16611388614e02565b91505080614b38576040517f1c26714c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805a838110614e3257839003604081048103851015614e3257600080885160208a0160008a8af19250600191505b5094509492505050565b508054600082559060005260206000209081019061149a9190614f77565b828054828255906000526020600020908101928215614ed4579160200282015b82811115614ed457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614e7a565b50614ee0929150614f77565b5090565b600483019183908215614ed45791602002820160005b83821115614f3e57835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302614efa565b8015614f6e5782816101000a81549063ffffffff0219169055600401602081600301049283019260010302614f3e565b5050614ee09291505b5b80821115614ee05760008155600101614f78565b63ffffffff8116811461149a57600080fd5b600060208284031215614fb057600080fd5b813561298681614f8c565b73ffffffffffffffffffffffffffffffffffffffff8116811461149a57600080fd5b600060208284031215614fef57600080fd5b813561298681614fbb565b6000815180845260005b8181101561502057602081850181015186830182015201615004565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006111736020830184614ffa565b6000806040838503121561508457600080fd5b823561508f81614fbb565b9150602083013561509f81614fbb565b809150509250929050565b600080600080600060a086880312156150c257600080fd5b85356150cd81614f8c565b945060208601356150dd81614f8c565b935060408601356150ed81614f8c565b925060608601356150fd81614f8c565b9150608086013562ffffff8116811461511557600080fd5b809150509295509295909350565b60008151808452602080850194506020840160005b8381101561516a57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101615138565b509495945050505050565b6020815260006111736020830184615123565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156151fe576151fe615188565b604052919050565b600082601f83011261521757600080fd5b813567ffffffffffffffff81111561523157615231615188565b61526260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016151b7565b81815284602083860101111561527757600080fd5b816020850160208301376000918101602001919091529392505050565b600080604083850312156152a757600080fd5b82356152b281614fbb565b9150602083013567ffffffffffffffff8111156152ce57600080fd5b6152da85828601615206565b9150509250929050565b6000602082840312156152f657600080fd5b813569ffffffffffffffffffff8116811461298657600080fd5b60008083601f84011261532257600080fd5b50813567ffffffffffffffff81111561533a57600080fd5b6020830191508360208260051b850101111561535557600080fd5b9250929050565b6000806000806040858703121561537257600080fd5b843567ffffffffffffffff8082111561538a57600080fd5b61539688838901615310565b909650945060208701359150808211156153af57600080fd5b506153bc87828801615310565b95989497509550505050565b60008060008060008060008060e0898b0312156153e457600080fd5b606089018a8111156153f557600080fd5b8998503567ffffffffffffffff8082111561540f57600080fd5b818b0191508b601f83011261542357600080fd5b81358181111561543257600080fd5b8c602082850101111561544457600080fd5b6020830199508098505060808b013591508082111561546257600080fd5b61546e8c838d01615310565b909750955060a08b013591508082111561548757600080fd5b506154948b828c01615310565b999c989b50969995989497949560c00135949350505050565b6000602082840312156154bf57600080fd5b5035919050565b600080604083850312156154d957600080fd5b82356154e481614fbb565b946020939093013593505050565b600067ffffffffffffffff82111561550c5761550c615188565b5060051b60200190565b600082601f83011261552757600080fd5b8135602061553c615537836154f2565b6151b7565b8083825260208201915060208460051b87010193508684111561555e57600080fd5b602086015b8481101561558357803561557681614fbb565b8352918301918301615563565b509695505050505050565b803560ff8116811461559f57600080fd5b919050565b803567ffffffffffffffff8116811461559f57600080fd5b60008060008060008060c087890312156155d557600080fd5b863567ffffffffffffffff808211156155ed57600080fd5b6155f98a838b01615516565b9750602089013591508082111561560f57600080fd5b61561b8a838b01615516565b965061562960408a0161558e565b9550606089013591508082111561563f57600080fd5b61564b8a838b01615206565b945061565960808a016155a4565b935060a089013591508082111561566f57600080fd5b5061567c89828a01615206565b9150509295509295509295565b6000806040838503121561569c57600080fd5b82356156a781614fbb565b9150602083013561509f81614f8c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff828116828216039080821115613351576133516156e6565b8082028115828204841417611176576111766156e6565b80820180821115611176576111766156e6565b60006020828403121561576e57600080fd5b5051919050565b60006020828403121561578757600080fd5b8151801515811461298657600080fd5b73ffffffffffffffffffffffffffffffffffffffff8416815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b600181811c9082168061581557607f821691505b602082108103613bf8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b63ffffffff818116838216019080821115613351576133516156e6565b81810381811115611176576111766156e6565b8181036000831280158383131683831282161715613351576133516156e6565b60ff8181168382160290811690818114613351576133516156e6565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526158ea8184018a615123565b905082810360808401526158fe8189615123565b905060ff871660a084015282810360c084015261591b8187614ffa565b905067ffffffffffffffff851660e08401528281036101008401526159408185614ffa565b9c9b505050505050505050505050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526159978285018b615123565b915083820360808501526159ab828a615123565b915060ff881660a085015283820360c08501526159c88288614ffa565b90861660e085015283810361010085015290506159408185614ffa565b600063ffffffff8216806159fb576159fb6156e6565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b60ff8181168382160190811115611176576111766156e6565b8051601781900b811461559f57600080fd5b60008060008060808587031215615a6457600080fd5b8451615a6f81614f8c565b809450506020808601519350604086015167ffffffffffffffff811115615a9557600080fd5b8601601f81018813615aa657600080fd5b8051615ab4615537826154f2565b81815260059190911b8201830190838101908a831115615ad357600080fd5b928401925b82841015615af857615ae984615a3c565b82529284019290840190615ad8565b8096505050505050615b0c60608601615a3c565b905092959194509250565b600082615b4d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600063ffffffff808316818103615b9157615b916156e6565b6001019392505050565b600061010080830160178c60170b8552602073ffffffffffffffffffffffffffffffffffffffff8d16602087015263ffffffff8c1660408701528360608701528293508a518084526101208701945060208c01935060005b81811015615c11578451840b86529482019493820193600101615bf3565b50505050508281036080840152615c288188614ffa565b915050615c3a60a083018660170b9052565b8360c0830152613cd760e083018464ffffffffff16905256fea164736f6c6343000818000a",
}

var DualAggregatorABI = DualAggregatorMetaData.ABI

var DualAggregatorBin = DualAggregatorMetaData.Bin

func DeployDualAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, minAnswer_ *big.Int, maxAnswer_ *big.Int, billingAccessController common.Address, requesterAccessController common.Address, decimals_ uint8, description_ string, secondaryProxy_ common.Address, cutoffTime_ uint32, maxSyncIterations_ uint32) (common.Address, *types.Transaction, *DualAggregator, error) {
	parsed, err := DualAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DualAggregatorBin), backend, link, minAnswer_, maxAnswer_, billingAccessController, requesterAccessController, decimals_, description_, secondaryProxy_, cutoffTime_, maxSyncIterations_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DualAggregator{address: address, abi: *parsed, DualAggregatorCaller: DualAggregatorCaller{contract: contract}, DualAggregatorTransactor: DualAggregatorTransactor{contract: contract}, DualAggregatorFilterer: DualAggregatorFilterer{contract: contract}}, nil
}

type DualAggregator struct {
	address common.Address
	abi     abi.ABI
	DualAggregatorCaller
	DualAggregatorTransactor
	DualAggregatorFilterer
}

type DualAggregatorCaller struct {
	contract *bind.BoundContract
}

type DualAggregatorTransactor struct {
	contract *bind.BoundContract
}

type DualAggregatorFilterer struct {
	contract *bind.BoundContract
}

type DualAggregatorSession struct {
	Contract     *DualAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DualAggregatorCallerSession struct {
	Contract *DualAggregatorCaller
	CallOpts bind.CallOpts
}

type DualAggregatorTransactorSession struct {
	Contract     *DualAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type DualAggregatorRaw struct {
	Contract *DualAggregator
}

type DualAggregatorCallerRaw struct {
	Contract *DualAggregatorCaller
}

type DualAggregatorTransactorRaw struct {
	Contract *DualAggregatorTransactor
}

func NewDualAggregator(address common.Address, backend bind.ContractBackend) (*DualAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(DualAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDualAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DualAggregator{address: address, abi: abi, DualAggregatorCaller: DualAggregatorCaller{contract: contract}, DualAggregatorTransactor: DualAggregatorTransactor{contract: contract}, DualAggregatorFilterer: DualAggregatorFilterer{contract: contract}}, nil
}

func NewDualAggregatorCaller(address common.Address, caller bind.ContractCaller) (*DualAggregatorCaller, error) {
	contract, err := bindDualAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCaller{contract: contract}, nil
}

func NewDualAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*DualAggregatorTransactor, error) {
	contract, err := bindDualAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorTransactor{contract: contract}, nil
}

func NewDualAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*DualAggregatorFilterer, error) {
	contract, err := bindDualAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorFilterer{contract: contract}, nil
}

func bindDualAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DualAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DualAggregator *DualAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DualAggregator.Contract.DualAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_DualAggregator *DualAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.Contract.DualAggregatorTransactor.contract.Transfer(opts)
}

func (_DualAggregator *DualAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DualAggregator.Contract.DualAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_DualAggregator *DualAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DualAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_DualAggregator *DualAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.Contract.contract.Transfer(opts)
}

func (_DualAggregator *DualAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DualAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_DualAggregator *DualAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) CheckEnabled() (bool, error) {
	return _DualAggregator.Contract.CheckEnabled(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _DualAggregator.Contract.CheckEnabled(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) Decimals() (uint8, error) {
	return _DualAggregator.Contract.Decimals(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) Decimals() (uint8, error) {
	return _DualAggregator.Contract.Decimals(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) Description() (string, error) {
	return _DualAggregator.Contract.Description(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) Description() (string, error) {
	return _DualAggregator.Contract.Description(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetAnswer(&_DualAggregator.CallOpts, roundId)
}

func (_DualAggregator *DualAggregatorCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetAnswer(&_DualAggregator.CallOpts, roundId)
}

func (_DualAggregator *DualAggregatorCaller) GetBilling(opts *bind.CallOpts) (GetBilling,

	error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(GetBilling)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPriceGwei = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ReasonableGasPriceGwei = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ObservationPaymentGjuels = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.TransmissionPaymentGjuels = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.AccountingGas = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_DualAggregator *DualAggregatorSession) GetBilling() (GetBilling,

	error) {
	return _DualAggregator.Contract.GetBilling(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) GetBilling() (GetBilling,

	error) {
	return _DualAggregator.Contract.GetBilling(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) GetBillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getBillingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) GetBillingAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetBillingAccessController(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) GetBillingAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetBillingAccessController(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) GetLinkToken() (common.Address, error) {
	return _DualAggregator.Contract.GetLinkToken(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) GetLinkToken() (common.Address, error) {
	return _DualAggregator.Contract.GetLinkToken(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) GetRequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getRequesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) GetRequesterAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetRequesterAccessController(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) GetRequesterAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetRequesterAccessController(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getRoundData", roundId)

	outstruct := new(GetRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_DualAggregator *DualAggregatorSession) GetRoundData(roundId *big.Int) (GetRoundData,

	error) {
	return _DualAggregator.Contract.GetRoundData(&_DualAggregator.CallOpts, roundId)
}

func (_DualAggregator *DualAggregatorCallerSession) GetRoundData(roundId *big.Int) (GetRoundData,

	error) {
	return _DualAggregator.Contract.GetRoundData(&_DualAggregator.CallOpts, roundId)
}

func (_DualAggregator *DualAggregatorCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetTimestamp(&_DualAggregator.CallOpts, roundId)
}

func (_DualAggregator *DualAggregatorCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetTimestamp(&_DualAggregator.CallOpts, roundId)
}

func (_DualAggregator *DualAggregatorCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) GetTransmitters() ([]common.Address, error) {
	return _DualAggregator.Contract.GetTransmitters(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) GetTransmitters() ([]common.Address, error) {
	return _DualAggregator.Contract.GetTransmitters(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) GetValidatorConfig(opts *bind.CallOpts) (GetValidatorConfig,

	error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getValidatorConfig")

	outstruct := new(GetValidatorConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.GasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_DualAggregator *DualAggregatorSession) GetValidatorConfig() (GetValidatorConfig,

	error) {
	return _DualAggregator.Contract.GetValidatorConfig(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) GetValidatorConfig() (GetValidatorConfig,

	error) {
	return _DualAggregator.Contract.GetValidatorConfig(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _DualAggregator.Contract.HasAccess(&_DualAggregator.CallOpts, _user, _calldata)
}

func (_DualAggregator *DualAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _DualAggregator.Contract.HasAccess(&_DualAggregator.CallOpts, _user, _calldata)
}

func (_DualAggregator *DualAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.LatestAnswer(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.LatestAnswer(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_DualAggregator *DualAggregatorSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _DualAggregator.Contract.LatestConfigDetails(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _DualAggregator.Contract.LatestConfigDetails(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_DualAggregator *DualAggregatorSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _DualAggregator.Contract.LatestConfigDigestAndEpoch(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _DualAggregator.Contract.LatestConfigDigestAndEpoch(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) LatestRound() (*big.Int, error) {
	return _DualAggregator.Contract.LatestRound(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _DualAggregator.Contract.LatestRound(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_DualAggregator *DualAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DualAggregator.Contract.LatestRoundData(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DualAggregator.Contract.LatestRoundData(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _DualAggregator.Contract.LatestTimestamp(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _DualAggregator.Contract.LatestTimestamp(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (LatestTransmissionDetails,

	error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(LatestTransmissionDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Round = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.LatestAnswer = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LatestTimestamp = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_DualAggregator *DualAggregatorSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _DualAggregator.Contract.LatestTransmissionDetails(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _DualAggregator.Contract.LatestTransmissionDetails(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _DualAggregator.Contract.LinkAvailableForPayment(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _DualAggregator.Contract.LinkAvailableForPayment(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MaxAnswer(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MaxAnswer(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) MinAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MinAnswer(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MinAnswer(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "oracleObservationCount", transmitterAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _DualAggregator.Contract.OracleObservationCount(&_DualAggregator.CallOpts, transmitterAddress)
}

func (_DualAggregator *DualAggregatorCallerSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _DualAggregator.Contract.OracleObservationCount(&_DualAggregator.CallOpts, transmitterAddress)
}

func (_DualAggregator *DualAggregatorCaller) OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "owedPayment", transmitterAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _DualAggregator.Contract.OwedPayment(&_DualAggregator.CallOpts, transmitterAddress)
}

func (_DualAggregator *DualAggregatorCallerSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _DualAggregator.Contract.OwedPayment(&_DualAggregator.CallOpts, transmitterAddress)
}

func (_DualAggregator *DualAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) Owner() (common.Address, error) {
	return _DualAggregator.Contract.Owner(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) Owner() (common.Address, error) {
	return _DualAggregator.Contract.Owner(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) TypeAndVersion() (string, error) {
	return _DualAggregator.Contract.TypeAndVersion(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) TypeAndVersion() (string, error) {
	return _DualAggregator.Contract.TypeAndVersion(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DualAggregator *DualAggregatorSession) Version() (*big.Int, error) {
	return _DualAggregator.Contract.Version(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorCallerSession) Version() (*big.Int, error) {
	return _DualAggregator.Contract.Version(&_DualAggregator.CallOpts)
}

func (_DualAggregator *DualAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "acceptOwnership")
}

func (_DualAggregator *DualAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptOwnership(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptOwnership(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_DualAggregator *DualAggregatorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptPayeeship(&_DualAggregator.TransactOpts, transmitter)
}

func (_DualAggregator *DualAggregatorTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptPayeeship(&_DualAggregator.TransactOpts, transmitter)
}

func (_DualAggregator *DualAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "addAccess", _user)
}

func (_DualAggregator *DualAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AddAccess(&_DualAggregator.TransactOpts, _user)
}

func (_DualAggregator *DualAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AddAccess(&_DualAggregator.TransactOpts, _user)
}

func (_DualAggregator *DualAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "disableAccessCheck")
}

func (_DualAggregator *DualAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.DisableAccessCheck(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.DisableAccessCheck(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "enableAccessCheck")
}

func (_DualAggregator *DualAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.EnableAccessCheck(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.EnableAccessCheck(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "removeAccess", _user)
}

func (_DualAggregator *DualAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.RemoveAccess(&_DualAggregator.TransactOpts, _user)
}

func (_DualAggregator *DualAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.RemoveAccess(&_DualAggregator.TransactOpts, _user)
}

func (_DualAggregator *DualAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "requestNewRound")
}

func (_DualAggregator *DualAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _DualAggregator.Contract.RequestNewRound(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _DualAggregator.Contract.RequestNewRound(&_DualAggregator.TransactOpts)
}

func (_DualAggregator *DualAggregatorTransactor) SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setBilling", maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_DualAggregator *DualAggregatorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBilling(&_DualAggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBilling(&_DualAggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_DualAggregator *DualAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

func (_DualAggregator *DualAggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBillingAccessController(&_DualAggregator.TransactOpts, _billingAccessController)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBillingAccessController(&_DualAggregator.TransactOpts, _billingAccessController)
}

func (_DualAggregator *DualAggregatorTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_DualAggregator *DualAggregatorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetConfig(&_DualAggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetConfig(&_DualAggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_DualAggregator *DualAggregatorTransactor) SetCutoffTime(opts *bind.TransactOpts, _cutoffTime uint32) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setCutoffTime", _cutoffTime)
}

func (_DualAggregator *DualAggregatorSession) SetCutoffTime(_cutoffTime uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetCutoffTime(&_DualAggregator.TransactOpts, _cutoffTime)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetCutoffTime(_cutoffTime uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetCutoffTime(&_DualAggregator.TransactOpts, _cutoffTime)
}

func (_DualAggregator *DualAggregatorTransactor) SetLinkToken(opts *bind.TransactOpts, linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setLinkToken", linkToken, recipient)
}

func (_DualAggregator *DualAggregatorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetLinkToken(&_DualAggregator.TransactOpts, linkToken, recipient)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetLinkToken(&_DualAggregator.TransactOpts, linkToken, recipient)
}

func (_DualAggregator *DualAggregatorTransactor) SetMaxSyncIterations(opts *bind.TransactOpts, _maxSyncIterations uint32) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setMaxSyncIterations", _maxSyncIterations)
}

func (_DualAggregator *DualAggregatorSession) SetMaxSyncIterations(_maxSyncIterations uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetMaxSyncIterations(&_DualAggregator.TransactOpts, _maxSyncIterations)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetMaxSyncIterations(_maxSyncIterations uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetMaxSyncIterations(&_DualAggregator.TransactOpts, _maxSyncIterations)
}

func (_DualAggregator *DualAggregatorTransactor) SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setPayees", transmitters, payees)
}

func (_DualAggregator *DualAggregatorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetPayees(&_DualAggregator.TransactOpts, transmitters, payees)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetPayees(&_DualAggregator.TransactOpts, transmitters, payees)
}

func (_DualAggregator *DualAggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, requesterAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setRequesterAccessController", requesterAccessController)
}

func (_DualAggregator *DualAggregatorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetRequesterAccessController(&_DualAggregator.TransactOpts, requesterAccessController)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetRequesterAccessController(&_DualAggregator.TransactOpts, requesterAccessController)
}

func (_DualAggregator *DualAggregatorTransactor) SetValidatorConfig(opts *bind.TransactOpts, newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setValidatorConfig", newValidator, newGasLimit)
}

func (_DualAggregator *DualAggregatorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetValidatorConfig(&_DualAggregator.TransactOpts, newValidator, newGasLimit)
}

func (_DualAggregator *DualAggregatorTransactorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetValidatorConfig(&_DualAggregator.TransactOpts, newValidator, newGasLimit)
}

func (_DualAggregator *DualAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transferOwnership", to)
}

func (_DualAggregator *DualAggregatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferOwnership(&_DualAggregator.TransactOpts, to)
}

func (_DualAggregator *DualAggregatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferOwnership(&_DualAggregator.TransactOpts, to)
}

func (_DualAggregator *DualAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_DualAggregator *DualAggregatorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferPayeeship(&_DualAggregator.TransactOpts, transmitter, proposed)
}

func (_DualAggregator *DualAggregatorTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferPayeeship(&_DualAggregator.TransactOpts, transmitter, proposed)
}

func (_DualAggregator *DualAggregatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_DualAggregator *DualAggregatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.Transmit(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_DualAggregator *DualAggregatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.Transmit(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_DualAggregator *DualAggregatorTransactor) TransmitSecondary(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transmitSecondary", reportContext, report, rs, ss, rawVs)
}

func (_DualAggregator *DualAggregatorSession) TransmitSecondary(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransmitSecondary(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_DualAggregator *DualAggregatorTransactorSession) TransmitSecondary(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransmitSecondary(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_DualAggregator *DualAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "withdrawFunds", recipient, amount)
}

func (_DualAggregator *DualAggregatorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawFunds(&_DualAggregator.TransactOpts, recipient, amount)
}

func (_DualAggregator *DualAggregatorTransactorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawFunds(&_DualAggregator.TransactOpts, recipient, amount)
}

func (_DualAggregator *DualAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "withdrawPayment", transmitter)
}

func (_DualAggregator *DualAggregatorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawPayment(&_DualAggregator.TransactOpts, transmitter)
}

func (_DualAggregator *DualAggregatorTransactorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawPayment(&_DualAggregator.TransactOpts, transmitter)
}

type DualAggregatorAddedAccessIterator struct {
	Event *DualAggregatorAddedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorAddedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorAddedAccess)
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
		it.Event = new(DualAggregatorAddedAccess)
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

func (it *DualAggregatorAddedAccessIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*DualAggregatorAddedAccessIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorAddedAccessIterator{contract: _DualAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *DualAggregatorAddedAccess) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorAddedAccess)
				if err := _DualAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseAddedAccess(log types.Log) (*DualAggregatorAddedAccess, error) {
	event := new(DualAggregatorAddedAccess)
	if err := _DualAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorAnswerUpdatedIterator struct {
	Event *DualAggregatorAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorAnswerUpdated)
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
		it.Event = new(DualAggregatorAnswerUpdated)
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

func (it *DualAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*DualAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorAnswerUpdatedIterator{contract: _DualAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *DualAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorAnswerUpdated)
				if err := _DualAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*DualAggregatorAnswerUpdated, error) {
	event := new(DualAggregatorAnswerUpdated)
	if err := _DualAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorBillingAccessControllerSetIterator struct {
	Event *DualAggregatorBillingAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorBillingAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorBillingAccessControllerSet)
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
		it.Event = new(DualAggregatorBillingAccessControllerSet)
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

func (it *DualAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*DualAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorBillingAccessControllerSetIterator{contract: _DualAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorBillingAccessControllerSet)
				if err := _DualAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*DualAggregatorBillingAccessControllerSet, error) {
	event := new(DualAggregatorBillingAccessControllerSet)
	if err := _DualAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorBillingSetIterator struct {
	Event *DualAggregatorBillingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorBillingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorBillingSet)
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
		it.Event = new(DualAggregatorBillingSet)
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

func (it *DualAggregatorBillingSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorBillingSet struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
	Raw                       types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*DualAggregatorBillingSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorBillingSetIterator{contract: _DualAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorBillingSet)
				if err := _DualAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseBillingSet(log types.Log) (*DualAggregatorBillingSet, error) {
	event := new(DualAggregatorBillingSet)
	if err := _DualAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorCheckAccessDisabledIterator struct {
	Event *DualAggregatorCheckAccessDisabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorCheckAccessDisabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorCheckAccessDisabled)
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
		it.Event = new(DualAggregatorCheckAccessDisabled)
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

func (it *DualAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorCheckAccessDisabled struct {
	Raw types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*DualAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCheckAccessDisabledIterator{contract: _DualAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *DualAggregatorCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorCheckAccessDisabled)
				if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*DualAggregatorCheckAccessDisabled, error) {
	event := new(DualAggregatorCheckAccessDisabled)
	if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorCheckAccessEnabledIterator struct {
	Event *DualAggregatorCheckAccessEnabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorCheckAccessEnabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorCheckAccessEnabled)
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
		it.Event = new(DualAggregatorCheckAccessEnabled)
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

func (it *DualAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorCheckAccessEnabled struct {
	Raw types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*DualAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCheckAccessEnabledIterator{contract: _DualAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *DualAggregatorCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorCheckAccessEnabled)
				if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*DualAggregatorCheckAccessEnabled, error) {
	event := new(DualAggregatorCheckAccessEnabled)
	if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorConfigSetIterator struct {
	Event *DualAggregatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorConfigSet)
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
		it.Event = new(DualAggregatorConfigSet)
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

func (it *DualAggregatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*DualAggregatorConfigSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorConfigSetIterator{contract: _DualAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorConfigSet)
				if err := _DualAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseConfigSet(log types.Log) (*DualAggregatorConfigSet, error) {
	event := new(DualAggregatorConfigSet)
	if err := _DualAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorCutoffTimeSetIterator struct {
	Event *DualAggregatorCutoffTimeSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorCutoffTimeSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorCutoffTimeSet)
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
		it.Event = new(DualAggregatorCutoffTimeSet)
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

func (it *DualAggregatorCutoffTimeSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorCutoffTimeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorCutoffTimeSet struct {
	CutoffTime uint32
	Raw        types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterCutoffTimeSet(opts *bind.FilterOpts) (*DualAggregatorCutoffTimeSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "CutoffTimeSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCutoffTimeSetIterator{contract: _DualAggregator.contract, event: "CutoffTimeSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchCutoffTimeSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorCutoffTimeSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "CutoffTimeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorCutoffTimeSet)
				if err := _DualAggregator.contract.UnpackLog(event, "CutoffTimeSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseCutoffTimeSet(log types.Log) (*DualAggregatorCutoffTimeSet, error) {
	event := new(DualAggregatorCutoffTimeSet)
	if err := _DualAggregator.contract.UnpackLog(event, "CutoffTimeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorLinkTokenSetIterator struct {
	Event *DualAggregatorLinkTokenSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorLinkTokenSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorLinkTokenSet)
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
		it.Event = new(DualAggregatorLinkTokenSet)
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

func (it *DualAggregatorLinkTokenSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, oldLinkToken []common.Address, newLinkToken []common.Address) (*DualAggregatorLinkTokenSetIterator, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorLinkTokenSetIterator{contract: _DualAggregator.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorLinkTokenSet, oldLinkToken []common.Address, newLinkToken []common.Address) (event.Subscription, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorLinkTokenSet)
				if err := _DualAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseLinkTokenSet(log types.Log) (*DualAggregatorLinkTokenSet, error) {
	event := new(DualAggregatorLinkTokenSet)
	if err := _DualAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorMaxSyncIterationsSetIterator struct {
	Event *DualAggregatorMaxSyncIterationsSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorMaxSyncIterationsSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorMaxSyncIterationsSet)
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
		it.Event = new(DualAggregatorMaxSyncIterationsSet)
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

func (it *DualAggregatorMaxSyncIterationsSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorMaxSyncIterationsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorMaxSyncIterationsSet struct {
	MaxSyncIterations uint32
	Raw               types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterMaxSyncIterationsSet(opts *bind.FilterOpts) (*DualAggregatorMaxSyncIterationsSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "MaxSyncIterationsSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorMaxSyncIterationsSetIterator{contract: _DualAggregator.contract, event: "MaxSyncIterationsSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchMaxSyncIterationsSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorMaxSyncIterationsSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "MaxSyncIterationsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorMaxSyncIterationsSet)
				if err := _DualAggregator.contract.UnpackLog(event, "MaxSyncIterationsSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseMaxSyncIterationsSet(log types.Log) (*DualAggregatorMaxSyncIterationsSet, error) {
	event := new(DualAggregatorMaxSyncIterationsSet)
	if err := _DualAggregator.contract.UnpackLog(event, "MaxSyncIterationsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorNewRoundIterator struct {
	Event *DualAggregatorNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorNewRound)
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
		it.Event = new(DualAggregatorNewRound)
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

func (it *DualAggregatorNewRoundIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*DualAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorNewRoundIterator{contract: _DualAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *DualAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorNewRound)
				if err := _DualAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseNewRound(log types.Log) (*DualAggregatorNewRound, error) {
	event := new(DualAggregatorNewRound)
	if err := _DualAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorNewTransmissionIterator struct {
	Event *DualAggregatorNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorNewTransmission)
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
		it.Event = new(DualAggregatorNewTransmission)
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

func (it *DualAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorNewTransmission struct {
	AggregatorRoundId     uint32
	Answer                *big.Int
	Transmitter           common.Address
	ObservationsTimestamp uint32
	Observations          []*big.Int
	Observers             []byte
	JuelsPerFeeCoin       *big.Int
	ConfigDigest          [32]byte
	EpochAndRound         *big.Int
	Raw                   types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*DualAggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorNewTransmissionIterator{contract: _DualAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *DualAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorNewTransmission)
				if err := _DualAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseNewTransmission(log types.Log) (*DualAggregatorNewTransmission, error) {
	event := new(DualAggregatorNewTransmission)
	if err := _DualAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorOraclePaidIterator struct {
	Event *DualAggregatorOraclePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorOraclePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorOraclePaid)
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
		it.Event = new(DualAggregatorOraclePaid)
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

func (it *DualAggregatorOraclePaidIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*DualAggregatorOraclePaidIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorOraclePaidIterator{contract: _DualAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *DualAggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorOraclePaid)
				if err := _DualAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseOraclePaid(log types.Log) (*DualAggregatorOraclePaid, error) {
	event := new(DualAggregatorOraclePaid)
	if err := _DualAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorOwnershipTransferRequestedIterator struct {
	Event *DualAggregatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorOwnershipTransferRequested)
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
		it.Event = new(DualAggregatorOwnershipTransferRequested)
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

func (it *DualAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DualAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorOwnershipTransferRequestedIterator{contract: _DualAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorOwnershipTransferRequested)
				if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*DualAggregatorOwnershipTransferRequested, error) {
	event := new(DualAggregatorOwnershipTransferRequested)
	if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorOwnershipTransferredIterator struct {
	Event *DualAggregatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorOwnershipTransferred)
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
		it.Event = new(DualAggregatorOwnershipTransferred)
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

func (it *DualAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DualAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorOwnershipTransferredIterator{contract: _DualAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DualAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorOwnershipTransferred)
				if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*DualAggregatorOwnershipTransferred, error) {
	event := new(DualAggregatorOwnershipTransferred)
	if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorPayeeshipTransferRequestedIterator struct {
	Event *DualAggregatorPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorPayeeshipTransferRequested)
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
		it.Event = new(DualAggregatorPayeeshipTransferRequested)
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

func (it *DualAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*DualAggregatorPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorPayeeshipTransferRequestedIterator{contract: _DualAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorPayeeshipTransferRequested)
				if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*DualAggregatorPayeeshipTransferRequested, error) {
	event := new(DualAggregatorPayeeshipTransferRequested)
	if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorPayeeshipTransferredIterator struct {
	Event *DualAggregatorPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorPayeeshipTransferred)
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
		it.Event = new(DualAggregatorPayeeshipTransferred)
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

func (it *DualAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*DualAggregatorPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorPayeeshipTransferredIterator{contract: _DualAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *DualAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorPayeeshipTransferred)
				if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*DualAggregatorPayeeshipTransferred, error) {
	event := new(DualAggregatorPayeeshipTransferred)
	if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorPrimaryFeedUnlockedIterator struct {
	Event *DualAggregatorPrimaryFeedUnlocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorPrimaryFeedUnlockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorPrimaryFeedUnlocked)
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
		it.Event = new(DualAggregatorPrimaryFeedUnlocked)
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

func (it *DualAggregatorPrimaryFeedUnlockedIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorPrimaryFeedUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorPrimaryFeedUnlocked struct {
	PrimaryRoundId uint32
	Raw            types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterPrimaryFeedUnlocked(opts *bind.FilterOpts, primaryRoundId []uint32) (*DualAggregatorPrimaryFeedUnlockedIterator, error) {

	var primaryRoundIdRule []interface{}
	for _, primaryRoundIdItem := range primaryRoundId {
		primaryRoundIdRule = append(primaryRoundIdRule, primaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "PrimaryFeedUnlocked", primaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorPrimaryFeedUnlockedIterator{contract: _DualAggregator.contract, event: "PrimaryFeedUnlocked", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchPrimaryFeedUnlocked(opts *bind.WatchOpts, sink chan<- *DualAggregatorPrimaryFeedUnlocked, primaryRoundId []uint32) (event.Subscription, error) {

	var primaryRoundIdRule []interface{}
	for _, primaryRoundIdItem := range primaryRoundId {
		primaryRoundIdRule = append(primaryRoundIdRule, primaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "PrimaryFeedUnlocked", primaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorPrimaryFeedUnlocked)
				if err := _DualAggregator.contract.UnpackLog(event, "PrimaryFeedUnlocked", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParsePrimaryFeedUnlocked(log types.Log) (*DualAggregatorPrimaryFeedUnlocked, error) {
	event := new(DualAggregatorPrimaryFeedUnlocked)
	if err := _DualAggregator.contract.UnpackLog(event, "PrimaryFeedUnlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorRemovedAccessIterator struct {
	Event *DualAggregatorRemovedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorRemovedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorRemovedAccess)
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
		it.Event = new(DualAggregatorRemovedAccess)
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

func (it *DualAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*DualAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorRemovedAccessIterator{contract: _DualAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *DualAggregatorRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorRemovedAccess)
				if err := _DualAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseRemovedAccess(log types.Log) (*DualAggregatorRemovedAccess, error) {
	event := new(DualAggregatorRemovedAccess)
	if err := _DualAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorRequesterAccessControllerSetIterator struct {
	Event *DualAggregatorRequesterAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorRequesterAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorRequesterAccessControllerSet)
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
		it.Event = new(DualAggregatorRequesterAccessControllerSet)
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

func (it *DualAggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*DualAggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorRequesterAccessControllerSetIterator{contract: _DualAggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorRequesterAccessControllerSet)
				if err := _DualAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*DualAggregatorRequesterAccessControllerSet, error) {
	event := new(DualAggregatorRequesterAccessControllerSet)
	if err := _DualAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorRoundRequestedIterator struct {
	Event *DualAggregatorRoundRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorRoundRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorRoundRequested)
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
		it.Event = new(DualAggregatorRoundRequested)
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

func (it *DualAggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [32]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*DualAggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorRoundRequestedIterator{contract: _DualAggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorRoundRequested)
				if err := _DualAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseRoundRequested(log types.Log) (*DualAggregatorRoundRequested, error) {
	event := new(DualAggregatorRoundRequested)
	if err := _DualAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorSecondaryRoundIdUpdatedIterator struct {
	Event *DualAggregatorSecondaryRoundIdUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorSecondaryRoundIdUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorSecondaryRoundIdUpdated)
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
		it.Event = new(DualAggregatorSecondaryRoundIdUpdated)
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

func (it *DualAggregatorSecondaryRoundIdUpdatedIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorSecondaryRoundIdUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorSecondaryRoundIdUpdated struct {
	SecondaryRoundId uint32
	Raw              types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterSecondaryRoundIdUpdated(opts *bind.FilterOpts, secondaryRoundId []uint32) (*DualAggregatorSecondaryRoundIdUpdatedIterator, error) {

	var secondaryRoundIdRule []interface{}
	for _, secondaryRoundIdItem := range secondaryRoundId {
		secondaryRoundIdRule = append(secondaryRoundIdRule, secondaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "SecondaryRoundIdUpdated", secondaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorSecondaryRoundIdUpdatedIterator{contract: _DualAggregator.contract, event: "SecondaryRoundIdUpdated", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchSecondaryRoundIdUpdated(opts *bind.WatchOpts, sink chan<- *DualAggregatorSecondaryRoundIdUpdated, secondaryRoundId []uint32) (event.Subscription, error) {

	var secondaryRoundIdRule []interface{}
	for _, secondaryRoundIdItem := range secondaryRoundId {
		secondaryRoundIdRule = append(secondaryRoundIdRule, secondaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "SecondaryRoundIdUpdated", secondaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorSecondaryRoundIdUpdated)
				if err := _DualAggregator.contract.UnpackLog(event, "SecondaryRoundIdUpdated", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseSecondaryRoundIdUpdated(log types.Log) (*DualAggregatorSecondaryRoundIdUpdated, error) {
	event := new(DualAggregatorSecondaryRoundIdUpdated)
	if err := _DualAggregator.contract.UnpackLog(event, "SecondaryRoundIdUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorTransmittedIterator struct {
	Event *DualAggregatorTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorTransmitted)
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
		it.Event = new(DualAggregatorTransmitted)
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

func (it *DualAggregatorTransmittedIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterTransmitted(opts *bind.FilterOpts) (*DualAggregatorTransmittedIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorTransmittedIterator{contract: _DualAggregator.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *DualAggregatorTransmitted) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorTransmitted)
				if err := _DualAggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseTransmitted(log types.Log) (*DualAggregatorTransmitted, error) {
	event := new(DualAggregatorTransmitted)
	if err := _DualAggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DualAggregatorValidatorConfigSetIterator struct {
	Event *DualAggregatorValidatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DualAggregatorValidatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorValidatorConfigSet)
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
		it.Event = new(DualAggregatorValidatorConfigSet)
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

func (it *DualAggregatorValidatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *DualAggregatorValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DualAggregatorValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log
}

func (_DualAggregator *DualAggregatorFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*DualAggregatorValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorValidatorConfigSetIterator{contract: _DualAggregator.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

func (_DualAggregator *DualAggregatorFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DualAggregatorValidatorConfigSet)
				if err := _DualAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
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

func (_DualAggregator *DualAggregatorFilterer) ParseValidatorConfigSet(log types.Log) (*DualAggregatorValidatorConfigSet, error) {
	event := new(DualAggregatorValidatorConfigSet)
	if err := _DualAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetBilling struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}
type GetRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type GetValidatorConfig struct {
	Validator common.Address
	GasLimit  uint32
}
type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}
type LatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestTransmissionDetails struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}

func (_DualAggregator *DualAggregator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DualAggregator.abi.Events["AddedAccess"].ID:
		return _DualAggregator.ParseAddedAccess(log)
	case _DualAggregator.abi.Events["AnswerUpdated"].ID:
		return _DualAggregator.ParseAnswerUpdated(log)
	case _DualAggregator.abi.Events["BillingAccessControllerSet"].ID:
		return _DualAggregator.ParseBillingAccessControllerSet(log)
	case _DualAggregator.abi.Events["BillingSet"].ID:
		return _DualAggregator.ParseBillingSet(log)
	case _DualAggregator.abi.Events["CheckAccessDisabled"].ID:
		return _DualAggregator.ParseCheckAccessDisabled(log)
	case _DualAggregator.abi.Events["CheckAccessEnabled"].ID:
		return _DualAggregator.ParseCheckAccessEnabled(log)
	case _DualAggregator.abi.Events["ConfigSet"].ID:
		return _DualAggregator.ParseConfigSet(log)
	case _DualAggregator.abi.Events["CutoffTimeSet"].ID:
		return _DualAggregator.ParseCutoffTimeSet(log)
	case _DualAggregator.abi.Events["LinkTokenSet"].ID:
		return _DualAggregator.ParseLinkTokenSet(log)
	case _DualAggregator.abi.Events["MaxSyncIterationsSet"].ID:
		return _DualAggregator.ParseMaxSyncIterationsSet(log)
	case _DualAggregator.abi.Events["NewRound"].ID:
		return _DualAggregator.ParseNewRound(log)
	case _DualAggregator.abi.Events["NewTransmission"].ID:
		return _DualAggregator.ParseNewTransmission(log)
	case _DualAggregator.abi.Events["OraclePaid"].ID:
		return _DualAggregator.ParseOraclePaid(log)
	case _DualAggregator.abi.Events["OwnershipTransferRequested"].ID:
		return _DualAggregator.ParseOwnershipTransferRequested(log)
	case _DualAggregator.abi.Events["OwnershipTransferred"].ID:
		return _DualAggregator.ParseOwnershipTransferred(log)
	case _DualAggregator.abi.Events["PayeeshipTransferRequested"].ID:
		return _DualAggregator.ParsePayeeshipTransferRequested(log)
	case _DualAggregator.abi.Events["PayeeshipTransferred"].ID:
		return _DualAggregator.ParsePayeeshipTransferred(log)
	case _DualAggregator.abi.Events["PrimaryFeedUnlocked"].ID:
		return _DualAggregator.ParsePrimaryFeedUnlocked(log)
	case _DualAggregator.abi.Events["RemovedAccess"].ID:
		return _DualAggregator.ParseRemovedAccess(log)
	case _DualAggregator.abi.Events["RequesterAccessControllerSet"].ID:
		return _DualAggregator.ParseRequesterAccessControllerSet(log)
	case _DualAggregator.abi.Events["RoundRequested"].ID:
		return _DualAggregator.ParseRoundRequested(log)
	case _DualAggregator.abi.Events["SecondaryRoundIdUpdated"].ID:
		return _DualAggregator.ParseSecondaryRoundIdUpdated(log)
	case _DualAggregator.abi.Events["Transmitted"].ID:
		return _DualAggregator.ParseTransmitted(log)
	case _DualAggregator.abi.Events["ValidatorConfigSet"].ID:
		return _DualAggregator.ParseValidatorConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DualAggregatorAddedAccess) Topic() common.Hash {
	return common.HexToHash("0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4")
}

func (DualAggregatorAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (DualAggregatorBillingAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912")
}

func (DualAggregatorBillingSet) Topic() common.Hash {
	return common.HexToHash("0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f")
}

func (DualAggregatorCheckAccessDisabled) Topic() common.Hash {
	return common.HexToHash("0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638")
}

func (DualAggregatorCheckAccessEnabled) Topic() common.Hash {
	return common.HexToHash("0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480")
}

func (DualAggregatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (DualAggregatorCutoffTimeSet) Topic() common.Hash {
	return common.HexToHash("0xb24a681ce3399a408a89fd0c2b59dfc24bdad592b1c7ec7671cf060596c1c4d1")
}

func (DualAggregatorLinkTokenSet) Topic() common.Hash {
	return common.HexToHash("0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a")
}

func (DualAggregatorMaxSyncIterationsSet) Topic() common.Hash {
	return common.HexToHash("0xcba51f727ba38740aa888ce0cb33f68de587733f61d3fafa0d9fb2b29e7f829f")
}

func (DualAggregatorNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (DualAggregatorNewTransmission) Topic() common.Hash {
	return common.HexToHash("0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a")
}

func (DualAggregatorOraclePaid) Topic() common.Hash {
	return common.HexToHash("0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c")
}

func (DualAggregatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (DualAggregatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (DualAggregatorPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (DualAggregatorPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (DualAggregatorPrimaryFeedUnlocked) Topic() common.Hash {
	return common.HexToHash("0xda2435684a37fba6f7841e49b59e6ad975e462bbebd28ec9da4ed9746a6992be")
}

func (DualAggregatorRemovedAccess) Topic() common.Hash {
	return common.HexToHash("0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1")
}

func (DualAggregatorRequesterAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634")
}

func (DualAggregatorRoundRequested) Topic() common.Hash {
	return common.HexToHash("0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c")
}

func (DualAggregatorSecondaryRoundIdUpdated) Topic() common.Hash {
	return common.HexToHash("0x8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b")
}

func (DualAggregatorTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (DualAggregatorValidatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541")
}

func (_DualAggregator *DualAggregator) Address() common.Address {
	return _DualAggregator.address
}

type DualAggregatorInterface interface {
	CheckEnabled(opts *bind.CallOpts) (bool, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	GetBilling(opts *bind.CallOpts) (GetBilling,

		error)

	GetBillingAccessController(opts *bind.CallOpts) (common.Address, error)

	GetLinkToken(opts *bind.CallOpts) (common.Address, error)

	GetRequesterAccessController(opts *bind.CallOpts) (common.Address, error)

	GetRoundData(opts *bind.CallOpts, roundId *big.Int) (GetRoundData,

		error)

	GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	GetTransmitters(opts *bind.CallOpts) ([]common.Address, error)

	GetValidatorConfig(opts *bind.CallOpts) (GetValidatorConfig,

		error)

	HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	LatestTimestamp(opts *bind.CallOpts) (*big.Int, error)

	LatestTransmissionDetails(opts *bind.CallOpts) (LatestTransmissionDetails,

		error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	MaxAnswer(opts *bind.CallOpts) (*big.Int, error)

	MinAnswer(opts *bind.CallOpts) (*big.Int, error)

	OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error)

	OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error)

	DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error)

	EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error)

	RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error)

	RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error)

	SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error)

	SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetCutoffTime(opts *bind.TransactOpts, _cutoffTime uint32) (*types.Transaction, error)

	SetLinkToken(opts *bind.TransactOpts, linkToken common.Address, recipient common.Address) (*types.Transaction, error)

	SetMaxSyncIterations(opts *bind.TransactOpts, _maxSyncIterations uint32) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error)

	SetRequesterAccessController(opts *bind.TransactOpts, requesterAccessController common.Address) (*types.Transaction, error)

	SetValidatorConfig(opts *bind.TransactOpts, newValidator common.Address, newGasLimit uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	TransmitSecondary(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	FilterAddedAccess(opts *bind.FilterOpts) (*DualAggregatorAddedAccessIterator, error)

	WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *DualAggregatorAddedAccess) (event.Subscription, error)

	ParseAddedAccess(log types.Log) (*DualAggregatorAddedAccess, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*DualAggregatorAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *DualAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*DualAggregatorAnswerUpdated, error)

	FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*DualAggregatorBillingAccessControllerSetIterator, error)

	WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorBillingAccessControllerSet) (event.Subscription, error)

	ParseBillingAccessControllerSet(log types.Log) (*DualAggregatorBillingAccessControllerSet, error)

	FilterBillingSet(opts *bind.FilterOpts) (*DualAggregatorBillingSetIterator, error)

	WatchBillingSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorBillingSet) (event.Subscription, error)

	ParseBillingSet(log types.Log) (*DualAggregatorBillingSet, error)

	FilterCheckAccessDisabled(opts *bind.FilterOpts) (*DualAggregatorCheckAccessDisabledIterator, error)

	WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *DualAggregatorCheckAccessDisabled) (event.Subscription, error)

	ParseCheckAccessDisabled(log types.Log) (*DualAggregatorCheckAccessDisabled, error)

	FilterCheckAccessEnabled(opts *bind.FilterOpts) (*DualAggregatorCheckAccessEnabledIterator, error)

	WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *DualAggregatorCheckAccessEnabled) (event.Subscription, error)

	ParseCheckAccessEnabled(log types.Log) (*DualAggregatorCheckAccessEnabled, error)

	FilterConfigSet(opts *bind.FilterOpts) (*DualAggregatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*DualAggregatorConfigSet, error)

	FilterCutoffTimeSet(opts *bind.FilterOpts) (*DualAggregatorCutoffTimeSetIterator, error)

	WatchCutoffTimeSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorCutoffTimeSet) (event.Subscription, error)

	ParseCutoffTimeSet(log types.Log) (*DualAggregatorCutoffTimeSet, error)

	FilterLinkTokenSet(opts *bind.FilterOpts, oldLinkToken []common.Address, newLinkToken []common.Address) (*DualAggregatorLinkTokenSetIterator, error)

	WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorLinkTokenSet, oldLinkToken []common.Address, newLinkToken []common.Address) (event.Subscription, error)

	ParseLinkTokenSet(log types.Log) (*DualAggregatorLinkTokenSet, error)

	FilterMaxSyncIterationsSet(opts *bind.FilterOpts) (*DualAggregatorMaxSyncIterationsSetIterator, error)

	WatchMaxSyncIterationsSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorMaxSyncIterationsSet) (event.Subscription, error)

	ParseMaxSyncIterationsSet(log types.Log) (*DualAggregatorMaxSyncIterationsSet, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*DualAggregatorNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *DualAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*DualAggregatorNewRound, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*DualAggregatorNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *DualAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*DualAggregatorNewTransmission, error)

	FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*DualAggregatorOraclePaidIterator, error)

	WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *DualAggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error)

	ParseOraclePaid(log types.Log) (*DualAggregatorOraclePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DualAggregatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*DualAggregatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DualAggregatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DualAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*DualAggregatorOwnershipTransferred, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*DualAggregatorPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*DualAggregatorPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*DualAggregatorPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *DualAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*DualAggregatorPayeeshipTransferred, error)

	FilterPrimaryFeedUnlocked(opts *bind.FilterOpts, primaryRoundId []uint32) (*DualAggregatorPrimaryFeedUnlockedIterator, error)

	WatchPrimaryFeedUnlocked(opts *bind.WatchOpts, sink chan<- *DualAggregatorPrimaryFeedUnlocked, primaryRoundId []uint32) (event.Subscription, error)

	ParsePrimaryFeedUnlocked(log types.Log) (*DualAggregatorPrimaryFeedUnlocked, error)

	FilterRemovedAccess(opts *bind.FilterOpts) (*DualAggregatorRemovedAccessIterator, error)

	WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *DualAggregatorRemovedAccess) (event.Subscription, error)

	ParseRemovedAccess(log types.Log) (*DualAggregatorRemovedAccess, error)

	FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*DualAggregatorRequesterAccessControllerSetIterator, error)

	WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorRequesterAccessControllerSet) (event.Subscription, error)

	ParseRequesterAccessControllerSet(log types.Log) (*DualAggregatorRequesterAccessControllerSet, error)

	FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*DualAggregatorRoundRequestedIterator, error)

	WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorRoundRequested, requester []common.Address) (event.Subscription, error)

	ParseRoundRequested(log types.Log) (*DualAggregatorRoundRequested, error)

	FilterSecondaryRoundIdUpdated(opts *bind.FilterOpts, secondaryRoundId []uint32) (*DualAggregatorSecondaryRoundIdUpdatedIterator, error)

	WatchSecondaryRoundIdUpdated(opts *bind.WatchOpts, sink chan<- *DualAggregatorSecondaryRoundIdUpdated, secondaryRoundId []uint32) (event.Subscription, error)

	ParseSecondaryRoundIdUpdated(log types.Log) (*DualAggregatorSecondaryRoundIdUpdated, error)

	FilterTransmitted(opts *bind.FilterOpts) (*DualAggregatorTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *DualAggregatorTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*DualAggregatorTransmitted, error)

	FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*DualAggregatorValidatorConfigSetIterator, error)

	WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error)

	ParseValidatorConfigSet(log types.Log) (*DualAggregatorValidatorConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
