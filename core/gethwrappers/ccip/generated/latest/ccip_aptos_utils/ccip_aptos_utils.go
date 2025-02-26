// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ccip_aptos_utils

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

type AptosUtilsAny2AptosRampMessage struct {
	Header       AptosUtilsRampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     [32]byte
	GasLimit     *big.Int
	TokenAmounts []AptosUtilsAny2AptosTokenTransfer
}

type AptosUtilsAny2AptosTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  [32]byte
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

type AptosUtilsCommitReport struct {
	PriceUpdates         AptosUtilsPriceUpdates
	BlessedMerkleRoots   []AptosUtilsMerkleRoot
	UnblessedMerkleRoots []AptosUtilsMerkleRoot
	RmnSignatures        []AptosUtilsRMNSignature
	OfframpAddress       [32]byte
}

type AptosUtilsEVMExtraArgsV1 struct {
	GasLimit *big.Int
}

type AptosUtilsEVMExtraArgsV2 struct {
	GasLimit                 *big.Int
	AllowOutOfOrderExecution bool
}

type AptosUtilsExecutionReport struct {
	SourceChainSelector uint64
	Message             AptosUtilsAny2AptosRampMessage
	OffchainTokenData   [][]byte
	Proofs              [][32]byte
}

type AptosUtilsGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type AptosUtilsMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSequenceNumber   uint64
	MaxSequenceNumber   uint64
	MerkleRoot          [32]byte
}

type AptosUtilsPriceUpdates struct {
	TokenPriceUpdates []AptosUtilsTokenPriceUpdate
	GasPriceUpdates   []AptosUtilsGasPriceUpdate
}

type AptosUtilsRMNSignature struct {
	R [32]byte
	S [32]byte
}

type AptosUtilsRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type AptosUtilsSVMExtraArgsV1 struct {
	ComputeUnits             uint32
	AccountIsWritableBitmap  uint64
	AllowOutOfOrderExecution bool
	TokenReceiver            [32]byte
	Accounts                 [][32]byte
}

type AptosUtilsTokenPriceUpdate struct {
	SourceToken [32]byte
	UsdPerToken *big.Int
}

var AptosUtilsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"exposeCommitReport\",\"inputs\":[{\"name\":\"commitReport\",\"type\":\"tuple\",\"internalType\":\"structAptosUtils.CommitReport\",\"components\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structAptosUtils.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structAptosUtils.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"usdPerToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structAptosUtils.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"blessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structAptosUtils.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"unblessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structAptosUtils.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"rmnSignatures\",\"type\":\"tuple[]\",\"internalType\":\"structAptosUtils.RMNSignature[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"offrampAddress\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeEVMExtraArgsV1\",\"inputs\":[{\"name\":\"evmExtraArgsV1\",\"type\":\"tuple\",\"internalType\":\"structAptosUtils.EVMExtraArgsV1\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeEVMExtraArgsV2\",\"inputs\":[{\"name\":\"evmExtraArgsV2\",\"type\":\"tuple\",\"internalType\":\"structAptosUtils.EVMExtraArgsV2\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeExecutionReport\",\"inputs\":[{\"name\":\"executionReport\",\"type\":\"tuple[]\",\"internalType\":\"structAptosUtils.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structAptosUtils.Any2AptosRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structAptosUtils.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structAptosUtils.Any2AptosTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposeSVMExtraArgsV1\",\"inputs\":[{\"name\":\"svmExtraArgsV1\",\"type\":\"tuple\",\"internalType\":\"structAptosUtils.SVMExtraArgsV1\",\"components\":[{\"name\":\"computeUnits\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accountIsWritableBitmap\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"tokenReceiver\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"accounts\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"}]",
}

var AptosUtilsABI = AptosUtilsMetaData.ABI

type AptosUtils struct {
	address common.Address
	abi     abi.ABI
	AptosUtilsCaller
	AptosUtilsTransactor
	AptosUtilsFilterer
}

type AptosUtilsCaller struct {
	contract *bind.BoundContract
}

type AptosUtilsTransactor struct {
	contract *bind.BoundContract
}

type AptosUtilsFilterer struct {
	contract *bind.BoundContract
}

type AptosUtilsSession struct {
	Contract     *AptosUtils
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AptosUtilsCallerSession struct {
	Contract *AptosUtilsCaller
	CallOpts bind.CallOpts
}

type AptosUtilsTransactorSession struct {
	Contract     *AptosUtilsTransactor
	TransactOpts bind.TransactOpts
}

type AptosUtilsRaw struct {
	Contract *AptosUtils
}

type AptosUtilsCallerRaw struct {
	Contract *AptosUtilsCaller
}

type AptosUtilsTransactorRaw struct {
	Contract *AptosUtilsTransactor
}

func NewAptosUtils(address common.Address, backend bind.ContractBackend) (*AptosUtils, error) {
	abi, err := abi.JSON(strings.NewReader(AptosUtilsABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAptosUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AptosUtils{address: address, abi: abi, AptosUtilsCaller: AptosUtilsCaller{contract: contract}, AptosUtilsTransactor: AptosUtilsTransactor{contract: contract}, AptosUtilsFilterer: AptosUtilsFilterer{contract: contract}}, nil
}

func NewAptosUtilsCaller(address common.Address, caller bind.ContractCaller) (*AptosUtilsCaller, error) {
	contract, err := bindAptosUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AptosUtilsCaller{contract: contract}, nil
}

func NewAptosUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*AptosUtilsTransactor, error) {
	contract, err := bindAptosUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AptosUtilsTransactor{contract: contract}, nil
}

func NewAptosUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*AptosUtilsFilterer, error) {
	contract, err := bindAptosUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AptosUtilsFilterer{contract: contract}, nil
}

func bindAptosUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AptosUtilsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AptosUtils *AptosUtilsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AptosUtils.Contract.AptosUtilsCaller.contract.Call(opts, result, method, params...)
}

func (_AptosUtils *AptosUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AptosUtils.Contract.AptosUtilsTransactor.contract.Transfer(opts)
}

func (_AptosUtils *AptosUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AptosUtils.Contract.AptosUtilsTransactor.contract.Transact(opts, method, params...)
}

func (_AptosUtils *AptosUtilsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AptosUtils.Contract.contract.Call(opts, result, method, params...)
}

func (_AptosUtils *AptosUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AptosUtils.Contract.contract.Transfer(opts)
}

func (_AptosUtils *AptosUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AptosUtils.Contract.contract.Transact(opts, method, params...)
}

func (_AptosUtils *AptosUtilsCaller) ExposeCommitReport(opts *bind.CallOpts, commitReport AptosUtilsCommitReport) ([]byte, error) {
	var out []interface{}
	err := _AptosUtils.contract.Call(opts, &out, "exposeCommitReport", commitReport)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AptosUtils *AptosUtilsSession) ExposeCommitReport(commitReport AptosUtilsCommitReport) ([]byte, error) {
	return _AptosUtils.Contract.ExposeCommitReport(&_AptosUtils.CallOpts, commitReport)
}

func (_AptosUtils *AptosUtilsCallerSession) ExposeCommitReport(commitReport AptosUtilsCommitReport) ([]byte, error) {
	return _AptosUtils.Contract.ExposeCommitReport(&_AptosUtils.CallOpts, commitReport)
}

func (_AptosUtils *AptosUtilsCaller) ExposeEVMExtraArgsV1(opts *bind.CallOpts, evmExtraArgsV1 AptosUtilsEVMExtraArgsV1) ([]byte, error) {
	var out []interface{}
	err := _AptosUtils.contract.Call(opts, &out, "exposeEVMExtraArgsV1", evmExtraArgsV1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AptosUtils *AptosUtilsSession) ExposeEVMExtraArgsV1(evmExtraArgsV1 AptosUtilsEVMExtraArgsV1) ([]byte, error) {
	return _AptosUtils.Contract.ExposeEVMExtraArgsV1(&_AptosUtils.CallOpts, evmExtraArgsV1)
}

func (_AptosUtils *AptosUtilsCallerSession) ExposeEVMExtraArgsV1(evmExtraArgsV1 AptosUtilsEVMExtraArgsV1) ([]byte, error) {
	return _AptosUtils.Contract.ExposeEVMExtraArgsV1(&_AptosUtils.CallOpts, evmExtraArgsV1)
}

func (_AptosUtils *AptosUtilsCaller) ExposeEVMExtraArgsV2(opts *bind.CallOpts, evmExtraArgsV2 AptosUtilsEVMExtraArgsV2) ([]byte, error) {
	var out []interface{}
	err := _AptosUtils.contract.Call(opts, &out, "exposeEVMExtraArgsV2", evmExtraArgsV2)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AptosUtils *AptosUtilsSession) ExposeEVMExtraArgsV2(evmExtraArgsV2 AptosUtilsEVMExtraArgsV2) ([]byte, error) {
	return _AptosUtils.Contract.ExposeEVMExtraArgsV2(&_AptosUtils.CallOpts, evmExtraArgsV2)
}

func (_AptosUtils *AptosUtilsCallerSession) ExposeEVMExtraArgsV2(evmExtraArgsV2 AptosUtilsEVMExtraArgsV2) ([]byte, error) {
	return _AptosUtils.Contract.ExposeEVMExtraArgsV2(&_AptosUtils.CallOpts, evmExtraArgsV2)
}

func (_AptosUtils *AptosUtilsCaller) ExposeExecutionReport(opts *bind.CallOpts, executionReport []AptosUtilsExecutionReport) ([]byte, error) {
	var out []interface{}
	err := _AptosUtils.contract.Call(opts, &out, "exposeExecutionReport", executionReport)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AptosUtils *AptosUtilsSession) ExposeExecutionReport(executionReport []AptosUtilsExecutionReport) ([]byte, error) {
	return _AptosUtils.Contract.ExposeExecutionReport(&_AptosUtils.CallOpts, executionReport)
}

func (_AptosUtils *AptosUtilsCallerSession) ExposeExecutionReport(executionReport []AptosUtilsExecutionReport) ([]byte, error) {
	return _AptosUtils.Contract.ExposeExecutionReport(&_AptosUtils.CallOpts, executionReport)
}

func (_AptosUtils *AptosUtilsCaller) ExposeSVMExtraArgsV1(opts *bind.CallOpts, svmExtraArgsV1 AptosUtilsSVMExtraArgsV1) ([]byte, error) {
	var out []interface{}
	err := _AptosUtils.contract.Call(opts, &out, "exposeSVMExtraArgsV1", svmExtraArgsV1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AptosUtils *AptosUtilsSession) ExposeSVMExtraArgsV1(svmExtraArgsV1 AptosUtilsSVMExtraArgsV1) ([]byte, error) {
	return _AptosUtils.Contract.ExposeSVMExtraArgsV1(&_AptosUtils.CallOpts, svmExtraArgsV1)
}

func (_AptosUtils *AptosUtilsCallerSession) ExposeSVMExtraArgsV1(svmExtraArgsV1 AptosUtilsSVMExtraArgsV1) ([]byte, error) {
	return _AptosUtils.Contract.ExposeSVMExtraArgsV1(&_AptosUtils.CallOpts, svmExtraArgsV1)
}

func (_AptosUtils *AptosUtils) Address() common.Address {
	return _AptosUtils.address
}

type AptosUtilsInterface interface {
	ExposeCommitReport(opts *bind.CallOpts, commitReport AptosUtilsCommitReport) ([]byte, error)

	ExposeEVMExtraArgsV1(opts *bind.CallOpts, evmExtraArgsV1 AptosUtilsEVMExtraArgsV1) ([]byte, error)

	ExposeEVMExtraArgsV2(opts *bind.CallOpts, evmExtraArgsV2 AptosUtilsEVMExtraArgsV2) ([]byte, error)

	ExposeExecutionReport(opts *bind.CallOpts, executionReport []AptosUtilsExecutionReport) ([]byte, error)

	ExposeSVMExtraArgsV1(opts *bind.CallOpts, svmExtraArgsV1 AptosUtilsSVMExtraArgsV1) ([]byte, error)

	Address() common.Address
}
