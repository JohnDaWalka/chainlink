// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package don_id_claimer

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

var DonIdClaimerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"capabilitiesRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimNextDonID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"donID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDonID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDONID\",\"inputs\":[{\"name\":\"donId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"syncDonIdWithCapReg\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516102be3803806102be83398101604081905261002f91610067565b6001600160a01b038116610056576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b0316608052610097565b60006020828403121561007957600080fd5b81516001600160a01b038116811461009057600080fd5b9392505050565b60805161020d6100b1600039600060a2015261020d6000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80636a6339c11461005c5780636e74336b14610072578063756e6f461461007b5780638e16ff101461008c578063957120d214610096575b600080fd5b6000545b60405190815260200160405180910390f35b61006060005481565b610060610089366004610167565b90565b61009461009e565b005b61006061014e565b60007f0000000000000000000000000000000000000000000000000000000000000000905060008173ffffffffffffffffffffffffffffffffffffffff1663fcdc8efe6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610110573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101349190610180565b90506101416001826101c3565b63ffffffff166000555050565b6000805461015d9060016101e7565b6000819055919050565b60006020828403121561017957600080fd5b5035919050565b60006020828403121561019257600080fd5b815163ffffffff811681146101a657600080fd5b9392505050565b634e487b7160e01b600052601160045260246000fd5b63ffffffff8281168282160390808211156101e0576101e06101ad565b5092915050565b808201808211156101fa576101fa6101ad565b9291505056fea164736f6c6343000818000a",
}

var DonIdClaimerABI = DonIdClaimerMetaData.ABI

var DonIdClaimerBin = DonIdClaimerMetaData.Bin

func DeployDonIdClaimer(auth *bind.TransactOpts, backend bind.ContractBackend, capabilitiesRegistry common.Address) (common.Address, *types.Transaction, *DonIdClaimer, error) {
	parsed, err := DonIdClaimerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DonIdClaimerBin), backend, capabilitiesRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DonIdClaimer{address: address, abi: *parsed, DonIdClaimerCaller: DonIdClaimerCaller{contract: contract}, DonIdClaimerTransactor: DonIdClaimerTransactor{contract: contract}, DonIdClaimerFilterer: DonIdClaimerFilterer{contract: contract}}, nil
}

type DonIdClaimer struct {
	address common.Address
	abi     abi.ABI
	DonIdClaimerCaller
	DonIdClaimerTransactor
	DonIdClaimerFilterer
}

type DonIdClaimerCaller struct {
	contract *bind.BoundContract
}

type DonIdClaimerTransactor struct {
	contract *bind.BoundContract
}

type DonIdClaimerFilterer struct {
	contract *bind.BoundContract
}

type DonIdClaimerSession struct {
	Contract     *DonIdClaimer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DonIdClaimerCallerSession struct {
	Contract *DonIdClaimerCaller
	CallOpts bind.CallOpts
}

type DonIdClaimerTransactorSession struct {
	Contract     *DonIdClaimerTransactor
	TransactOpts bind.TransactOpts
}

type DonIdClaimerRaw struct {
	Contract *DonIdClaimer
}

type DonIdClaimerCallerRaw struct {
	Contract *DonIdClaimerCaller
}

type DonIdClaimerTransactorRaw struct {
	Contract *DonIdClaimerTransactor
}

func NewDonIdClaimer(address common.Address, backend bind.ContractBackend) (*DonIdClaimer, error) {
	abi, err := abi.JSON(strings.NewReader(DonIdClaimerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDonIdClaimer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DonIdClaimer{address: address, abi: abi, DonIdClaimerCaller: DonIdClaimerCaller{contract: contract}, DonIdClaimerTransactor: DonIdClaimerTransactor{contract: contract}, DonIdClaimerFilterer: DonIdClaimerFilterer{contract: contract}}, nil
}

func NewDonIdClaimerCaller(address common.Address, caller bind.ContractCaller) (*DonIdClaimerCaller, error) {
	contract, err := bindDonIdClaimer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DonIdClaimerCaller{contract: contract}, nil
}

func NewDonIdClaimerTransactor(address common.Address, transactor bind.ContractTransactor) (*DonIdClaimerTransactor, error) {
	contract, err := bindDonIdClaimer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DonIdClaimerTransactor{contract: contract}, nil
}

func NewDonIdClaimerFilterer(address common.Address, filterer bind.ContractFilterer) (*DonIdClaimerFilterer, error) {
	contract, err := bindDonIdClaimer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DonIdClaimerFilterer{contract: contract}, nil
}

func bindDonIdClaimer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DonIdClaimerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DonIdClaimer *DonIdClaimerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DonIdClaimer.Contract.DonIdClaimerCaller.contract.Call(opts, result, method, params...)
}

func (_DonIdClaimer *DonIdClaimerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DonIdClaimer.Contract.DonIdClaimerTransactor.contract.Transfer(opts)
}

func (_DonIdClaimer *DonIdClaimerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DonIdClaimer.Contract.DonIdClaimerTransactor.contract.Transact(opts, method, params...)
}

func (_DonIdClaimer *DonIdClaimerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DonIdClaimer.Contract.contract.Call(opts, result, method, params...)
}

func (_DonIdClaimer *DonIdClaimerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DonIdClaimer.Contract.contract.Transfer(opts)
}

func (_DonIdClaimer *DonIdClaimerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DonIdClaimer.Contract.contract.Transact(opts, method, params...)
}

func (_DonIdClaimer *DonIdClaimerCaller) DonID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DonIdClaimer.contract.Call(opts, &out, "donID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DonIdClaimer *DonIdClaimerSession) DonID() (*big.Int, error) {
	return _DonIdClaimer.Contract.DonID(&_DonIdClaimer.CallOpts)
}

func (_DonIdClaimer *DonIdClaimerCallerSession) DonID() (*big.Int, error) {
	return _DonIdClaimer.Contract.DonID(&_DonIdClaimer.CallOpts)
}

func (_DonIdClaimer *DonIdClaimerCaller) GetDonID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DonIdClaimer.contract.Call(opts, &out, "getDonID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DonIdClaimer *DonIdClaimerSession) GetDonID() (*big.Int, error) {
	return _DonIdClaimer.Contract.GetDonID(&_DonIdClaimer.CallOpts)
}

func (_DonIdClaimer *DonIdClaimerCallerSession) GetDonID() (*big.Int, error) {
	return _DonIdClaimer.Contract.GetDonID(&_DonIdClaimer.CallOpts)
}

func (_DonIdClaimer *DonIdClaimerCaller) SetDONID(opts *bind.CallOpts, donId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DonIdClaimer.contract.Call(opts, &out, "setDONID", donId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DonIdClaimer *DonIdClaimerSession) SetDONID(donId *big.Int) (*big.Int, error) {
	return _DonIdClaimer.Contract.SetDONID(&_DonIdClaimer.CallOpts, donId)
}

func (_DonIdClaimer *DonIdClaimerCallerSession) SetDONID(donId *big.Int) (*big.Int, error) {
	return _DonIdClaimer.Contract.SetDONID(&_DonIdClaimer.CallOpts, donId)
}

func (_DonIdClaimer *DonIdClaimerTransactor) ClaimNextDonID(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DonIdClaimer.contract.Transact(opts, "claimNextDonID")
}

func (_DonIdClaimer *DonIdClaimerSession) ClaimNextDonID() (*types.Transaction, error) {
	return _DonIdClaimer.Contract.ClaimNextDonID(&_DonIdClaimer.TransactOpts)
}

func (_DonIdClaimer *DonIdClaimerTransactorSession) ClaimNextDonID() (*types.Transaction, error) {
	return _DonIdClaimer.Contract.ClaimNextDonID(&_DonIdClaimer.TransactOpts)
}

func (_DonIdClaimer *DonIdClaimerTransactor) SyncDonIdWithCapReg(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DonIdClaimer.contract.Transact(opts, "syncDonIdWithCapReg")
}

func (_DonIdClaimer *DonIdClaimerSession) SyncDonIdWithCapReg() (*types.Transaction, error) {
	return _DonIdClaimer.Contract.SyncDonIdWithCapReg(&_DonIdClaimer.TransactOpts)
}

func (_DonIdClaimer *DonIdClaimerTransactorSession) SyncDonIdWithCapReg() (*types.Transaction, error) {
	return _DonIdClaimer.Contract.SyncDonIdWithCapReg(&_DonIdClaimer.TransactOpts)
}

func (_DonIdClaimer *DonIdClaimer) Address() common.Address {
	return _DonIdClaimer.address
}

type DonIdClaimerInterface interface {
	DonID(opts *bind.CallOpts) (*big.Int, error)

	GetDonID(opts *bind.CallOpts) (*big.Int, error)

	SetDONID(opts *bind.CallOpts, donId *big.Int) (*big.Int, error)

	ClaimNextDonID(opts *bind.TransactOpts) (*types.Transaction, error)

	SyncDonIdWithCapReg(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
