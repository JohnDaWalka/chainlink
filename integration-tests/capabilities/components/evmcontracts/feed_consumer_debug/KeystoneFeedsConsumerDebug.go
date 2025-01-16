// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feeds_consumer_debug

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

// Reference imports to suppress errors if they are not otherwise used.
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

// FeedsConsumerDebugMetaData contains all meta data concerning the FeedsConsumerDebug contract.
var FeedsConsumerDebugMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"}],\"name\":\"UnauthorizedWorkflowName\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"}],\"name\":\"UnauthorizedWorkflowOwner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint224\",\"name\":\"price\",\"type\":\"uint224\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"FeedReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllFeeds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint224[]\",\"name\":\"\",\"type\":\"uint224[]\"},{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_allowedSendersList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_allowedWorkflowOwnersList\",\"type\":\"address[]\"},{\"internalType\":\"bytes10[]\",\"name\":\"_allowedWorkflowNamesList\",\"type\":\"bytes10[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61165b806101576000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063805f21321161005b578063805f21321461018b5780638da5cb5b1461019e578063e3401711146101c6578063f2fde38b146101d957600080fd5b806301ffc9a71461008d57806312876798146100b557806331d98b3f146100cc57806379ba509714610181575b600080fd5b6100a061009b3660046110c0565b6101ec565b60405190151581526020015b60405180910390f35b6100bd610285565b6040516100ac9392919061114b565b6101486100da3660046111f7565b6000908152600360209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168084527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169290910182905291565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909316835263ffffffff9091166020830152016100ac565b6101896104c1565b005b610189610199366004611259565b6105c3565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ac565b6101896101d436600461130a565b610979565b6101896101e73660046113a4565b610db9565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f805f213200000000000000000000000000000000000000000000000000000000148061027f57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b6002546060908190819060008167ffffffffffffffff8111156102aa576102aa6113da565b6040519080825280602002602001820160405280156102d3578160200160208202803683370190505b50905060008267ffffffffffffffff8111156102f1576102f16113da565b60405190808252806020026020018201604052801561031a578160200160208202803683370190505b50905060008367ffffffffffffffff811115610338576103386113da565b604051908082528060200260200182016040528015610361578160200160208202803683370190505b50905060005b848110156104b35760006002828154811061038457610384611409565b6000918252602080832090910154808352600382526040928390208351808501909452547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116845263ffffffff7c010000000000000000000000000000000000000000000000000000000090910416918301919091528651909250829087908590811061041057610410611409565b602002602001018181525050806000015185848151811061043357610433611409565b60200260200101907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681525050806020015184848151811061049457610494611409565b63ffffffff909216602092830291909101909101525050600101610367565b509196909550909350915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610547576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3360009081526005602052604090205460ff1661060e576040517f3fcc3f1700000000000000000000000000000000000000000000000000000000815233600482015260240161053e565b60008061065086868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610dcd92505050565b7fffffffffffffffffffff000000000000000000000000000000000000000000008216600090815260096020526040902054919350915060ff166106e4576040517f4b942f800000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffff000000000000000000000000000000000000000000008316600482015260240161053e565b73ffffffffffffffffffffffffffffffffffffffff811660009081526007602052604090205460ff1661075b576040517fbf24162300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161053e565b6000610769848601866114b0565b905060005b815181101561096f57604051806040016040528083838151811061079457610794611409565b6020026020010151602001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020018383815181106107d5576107d5611409565b60200260200101516040015163ffffffff168152506003600084848151811061080057610800611409565b602090810291909101810151518252818101929092526040016000208251929091015163ffffffff167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909216919091179055815160029083908390811061088557610885611409565b602090810291909101810151518254600181018455600093845291909220015581518290829081106108b9576108b9611409565b6020026020010151600001517f2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d8383815181106108f8576108f8611409565b60200260200101516020015184848151811061091657610916611409565b60200260200101516040015160405161095f9291907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216825263ffffffff16602082015260400190565b60405180910390a260010161076e565b5050505050505050565b610981610de3565b60005b60045463ffffffff82161015610a225760006005600060048463ffffffff16815481106109b3576109b3611409565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610a1b816115c2565b9050610984565b5060005b63ffffffff8116861115610aca5760016005600089898563ffffffff16818110610a5257610a52611409565b9050602002016020810190610a6791906113a4565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610ac3816115c2565b9050610a26565b50610ad760048787610f5b565b5060005b60065463ffffffff82161015610b795760006007600060068463ffffffff1681548110610b0a57610b0a611409565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610b72816115c2565b9050610adb565b5060005b63ffffffff8116841115610c215760016007600087878563ffffffff16818110610ba957610ba9611409565b9050602002016020810190610bbe91906113a4565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610c1a816115c2565b9050610b7d565b50610c2e60068585610f5b565b5060005b60085463ffffffff82161015610cef5760006009600060088463ffffffff1681548110610c6157610c61611409565b600091825260208083206003808404909101549206600a026101000a90910460b01b7fffffffffffffffffffff00000000000000000000000000000000000000000000168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610ce8816115c2565b9050610c32565b5060005b63ffffffff8116821115610da35760016009600085858563ffffffff16818110610d1f57610d1f611409565b9050602002016020810190610d34919061160c565b7fffffffffffffffffffff00000000000000000000000000000000000000000000168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610d9c816115c2565b9050610cf3565b50610db060088383610fe3565b50505050505050565b610dc1610de3565b610dca81610e66565b50565b6040810151604a90910151909160609190911c90565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e64576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161053e565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610ee5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161053e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610fd3579160200282015b82811115610fd35781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610f7b565b50610fdf9291506110ab565b5090565b82805482825590600052602060002090600201600390048101928215610fd35791602002820160005b8382111561106c57833575ffffffffffffffffffffffffffffffffffffffffffff191683826101000a81548169ffffffffffffffffffff021916908360b01c02179055509260200192600a0160208160090104928301926001030261100c565b80156110a25782816101000a81549069ffffffffffffffffffff0219169055600a0160208160090104928301926001030261106c565b5050610fdf9291505b5b80821115610fdf57600081556001016110ac565b6000602082840312156110d257600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461110257600080fd5b9392505050565b60008151808452602080850194506020840160005b8381101561114057815163ffffffff168752958201959082019060010161111e565b509495945050505050565b606080825284519082018190526000906020906080840190828801845b8281101561118457815184529284019290840190600101611168565b5050508381038285015285518082528683019183019060005b818110156111d75783517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161119d565b505084810360408601526111eb8187611109565b98975050505050505050565b60006020828403121561120957600080fd5b5035919050565b60008083601f84011261122257600080fd5b50813567ffffffffffffffff81111561123a57600080fd5b60208301915083602082850101111561125257600080fd5b9250929050565b6000806000806040858703121561126f57600080fd5b843567ffffffffffffffff8082111561128757600080fd5b61129388838901611210565b909650945060208701359150808211156112ac57600080fd5b506112b987828801611210565b95989497509550505050565b60008083601f8401126112d757600080fd5b50813567ffffffffffffffff8111156112ef57600080fd5b6020830191508360208260051b850101111561125257600080fd5b6000806000806000806060878903121561132357600080fd5b863567ffffffffffffffff8082111561133b57600080fd5b6113478a838b016112c5565b9098509650602089013591508082111561136057600080fd5b61136c8a838b016112c5565b9096509450604089013591508082111561138557600080fd5b5061139289828a016112c5565b979a9699509497509295939492505050565b6000602082840312156113b657600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461110257600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6040516060810167ffffffffffffffff8111828210171561145b5761145b6113da565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156114a8576114a86113da565b604052919050565b600060208083850312156114c357600080fd5b823567ffffffffffffffff808211156114db57600080fd5b818501915085601f8301126114ef57600080fd5b813581811115611501576115016113da565b61150f848260051b01611461565b8181528481019250606091820284018501918883111561152e57600080fd5b938501935b828510156115b65780858a03121561154b5760008081fd5b611553611438565b85358152868601357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146115865760008081fd5b8188015260408681013563ffffffff811681146115a35760008081fd5b9082015284529384019392850192611533565b50979650505050505050565b600063ffffffff808316818103611602577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b60006020828403121561161e57600080fd5b81357fffffffffffffffffffff000000000000000000000000000000000000000000008116811461110257600080fdfea164736f6c6343000818000a",
}

// FeedsConsumerDebugABI is the input ABI used to generate the binding from.
// Deprecated: Use FeedsConsumerDebugMetaData.ABI instead.
var FeedsConsumerDebugABI = FeedsConsumerDebugMetaData.ABI

// FeedsConsumerDebugBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FeedsConsumerDebugMetaData.Bin instead.
var FeedsConsumerDebugBin = FeedsConsumerDebugMetaData.Bin

// DeployFeedsConsumerDebug deploys a new Ethereum contract, binding an instance of FeedsConsumerDebug to it.
func DeployFeedsConsumerDebug(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FeedsConsumerDebug, error) {
	parsed, err := FeedsConsumerDebugMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeedsConsumerDebugBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FeedsConsumerDebug{FeedsConsumerDebugCaller: FeedsConsumerDebugCaller{contract: contract}, FeedsConsumerDebugTransactor: FeedsConsumerDebugTransactor{contract: contract}, FeedsConsumerDebugFilterer: FeedsConsumerDebugFilterer{contract: contract}}, nil
}

// FeedsConsumerDebug is an auto generated Go binding around an Ethereum contract.
type FeedsConsumerDebug struct {
	FeedsConsumerDebugCaller     // Read-only binding to the contract
	FeedsConsumerDebugTransactor // Write-only binding to the contract
	FeedsConsumerDebugFilterer   // Log filterer for contract events
}

// FeedsConsumerDebugCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeedsConsumerDebugCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeedsConsumerDebugTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeedsConsumerDebugTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeedsConsumerDebugFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeedsConsumerDebugFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeedsConsumerDebugSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeedsConsumerDebugSession struct {
	Contract     *FeedsConsumerDebug // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// FeedsConsumerDebugCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeedsConsumerDebugCallerSession struct {
	Contract *FeedsConsumerDebugCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// FeedsConsumerDebugTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeedsConsumerDebugTransactorSession struct {
	Contract     *FeedsConsumerDebugTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// FeedsConsumerDebugRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeedsConsumerDebugRaw struct {
	Contract *FeedsConsumerDebug // Generic contract binding to access the raw methods on
}

// FeedsConsumerDebugCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeedsConsumerDebugCallerRaw struct {
	Contract *FeedsConsumerDebugCaller // Generic read-only contract binding to access the raw methods on
}

// FeedsConsumerDebugTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeedsConsumerDebugTransactorRaw struct {
	Contract *FeedsConsumerDebugTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeedsConsumerDebug creates a new instance of FeedsConsumerDebug, bound to a specific deployed contract.
func NewFeedsConsumerDebug(address common.Address, backend bind.ContractBackend) (*FeedsConsumerDebug, error) {
	contract, err := bindFeedsConsumerDebug(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebug{FeedsConsumerDebugCaller: FeedsConsumerDebugCaller{contract: contract}, FeedsConsumerDebugTransactor: FeedsConsumerDebugTransactor{contract: contract}, FeedsConsumerDebugFilterer: FeedsConsumerDebugFilterer{contract: contract}}, nil
}

// NewFeedsConsumerDebugCaller creates a new read-only instance of FeedsConsumerDebug, bound to a specific deployed contract.
func NewFeedsConsumerDebugCaller(address common.Address, caller bind.ContractCaller) (*FeedsConsumerDebugCaller, error) {
	contract, err := bindFeedsConsumerDebug(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebugCaller{contract: contract}, nil
}

// NewFeedsConsumerDebugTransactor creates a new write-only instance of FeedsConsumerDebug, bound to a specific deployed contract.
func NewFeedsConsumerDebugTransactor(address common.Address, transactor bind.ContractTransactor) (*FeedsConsumerDebugTransactor, error) {
	contract, err := bindFeedsConsumerDebug(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebugTransactor{contract: contract}, nil
}

// NewFeedsConsumerDebugFilterer creates a new log filterer instance of FeedsConsumerDebug, bound to a specific deployed contract.
func NewFeedsConsumerDebugFilterer(address common.Address, filterer bind.ContractFilterer) (*FeedsConsumerDebugFilterer, error) {
	contract, err := bindFeedsConsumerDebug(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebugFilterer{contract: contract}, nil
}

// bindFeedsConsumerDebug binds a generic wrapper to an already deployed contract.
func bindFeedsConsumerDebug(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeedsConsumerDebugMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeedsConsumerDebug *FeedsConsumerDebugRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeedsConsumerDebug.Contract.FeedsConsumerDebugCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeedsConsumerDebug *FeedsConsumerDebugRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.FeedsConsumerDebugTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeedsConsumerDebug *FeedsConsumerDebugRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.FeedsConsumerDebugTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeedsConsumerDebug *FeedsConsumerDebugCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeedsConsumerDebug.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.contract.Transact(opts, method, params...)
}

// GetAllFeeds is a free data retrieval call binding the contract method 0x12876798.
//
// Solidity: function getAllFeeds() view returns(bytes32[], uint224[], uint32[])
func (_FeedsConsumerDebug *FeedsConsumerDebugCaller) GetAllFeeds(opts *bind.CallOpts) ([][32]byte, []*big.Int, []uint32, error) {
	var out []interface{}
	err := _FeedsConsumerDebug.contract.Call(opts, &out, "getAllFeeds")

	if err != nil {
		return *new([][32]byte), *new([]*big.Int), *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	out2 := *abi.ConvertType(out[2], new([]uint32)).(*[]uint32)

	return out0, out1, out2, err

}

// GetAllFeeds is a free data retrieval call binding the contract method 0x12876798.
//
// Solidity: function getAllFeeds() view returns(bytes32[], uint224[], uint32[])
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) GetAllFeeds() ([][32]byte, []*big.Int, []uint32, error) {
	return _FeedsConsumerDebug.Contract.GetAllFeeds(&_FeedsConsumerDebug.CallOpts)
}

// GetAllFeeds is a free data retrieval call binding the contract method 0x12876798.
//
// Solidity: function getAllFeeds() view returns(bytes32[], uint224[], uint32[])
func (_FeedsConsumerDebug *FeedsConsumerDebugCallerSession) GetAllFeeds() ([][32]byte, []*big.Int, []uint32, error) {
	return _FeedsConsumerDebug.Contract.GetAllFeeds(&_FeedsConsumerDebug.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x31d98b3f.
//
// Solidity: function getPrice(bytes32 feedId) view returns(uint224, uint32)
func (_FeedsConsumerDebug *FeedsConsumerDebugCaller) GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error) {
	var out []interface{}
	err := _FeedsConsumerDebug.contract.Call(opts, &out, "getPrice", feedId)

	if err != nil {
		return *new(*big.Int), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return out0, out1, err

}

// GetPrice is a free data retrieval call binding the contract method 0x31d98b3f.
//
// Solidity: function getPrice(bytes32 feedId) view returns(uint224, uint32)
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _FeedsConsumerDebug.Contract.GetPrice(&_FeedsConsumerDebug.CallOpts, feedId)
}

// GetPrice is a free data retrieval call binding the contract method 0x31d98b3f.
//
// Solidity: function getPrice(bytes32 feedId) view returns(uint224, uint32)
func (_FeedsConsumerDebug *FeedsConsumerDebugCallerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _FeedsConsumerDebug.Contract.GetPrice(&_FeedsConsumerDebug.CallOpts, feedId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeedsConsumerDebug *FeedsConsumerDebugCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeedsConsumerDebug.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) Owner() (common.Address, error) {
	return _FeedsConsumerDebug.Contract.Owner(&_FeedsConsumerDebug.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeedsConsumerDebug *FeedsConsumerDebugCallerSession) Owner() (common.Address, error) {
	return _FeedsConsumerDebug.Contract.Owner(&_FeedsConsumerDebug.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_FeedsConsumerDebug *FeedsConsumerDebugCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _FeedsConsumerDebug.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeedsConsumerDebug.Contract.SupportsInterface(&_FeedsConsumerDebug.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_FeedsConsumerDebug *FeedsConsumerDebugCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeedsConsumerDebug.Contract.SupportsInterface(&_FeedsConsumerDebug.CallOpts, interfaceId)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeedsConsumerDebug.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.AcceptOwnership(&_FeedsConsumerDebug.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.AcceptOwnership(&_FeedsConsumerDebug.TransactOpts)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes metadata, bytes rawReport) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactor) OnReport(opts *bind.TransactOpts, metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _FeedsConsumerDebug.contract.Transact(opts, "onReport", metadata, rawReport)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes metadata, bytes rawReport) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.OnReport(&_FeedsConsumerDebug.TransactOpts, metadata, rawReport)
}

// OnReport is a paid mutator transaction binding the contract method 0x805f2132.
//
// Solidity: function onReport(bytes metadata, bytes rawReport) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactorSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.OnReport(&_FeedsConsumerDebug.TransactOpts, metadata, rawReport)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3401711.
//
// Solidity: function setConfig(address[] _allowedSendersList, address[] _allowedWorkflowOwnersList, bytes10[] _allowedWorkflowNamesList) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactor) SetConfig(opts *bind.TransactOpts, _allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _FeedsConsumerDebug.contract.Transact(opts, "setConfig", _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3401711.
//
// Solidity: function setConfig(address[] _allowedSendersList, address[] _allowedWorkflowOwnersList, bytes10[] _allowedWorkflowNamesList) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) SetConfig(_allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.SetConfig(&_FeedsConsumerDebug.TransactOpts, _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3401711.
//
// Solidity: function setConfig(address[] _allowedSendersList, address[] _allowedWorkflowOwnersList, bytes10[] _allowedWorkflowNamesList) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactorSession) SetConfig(_allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.SetConfig(&_FeedsConsumerDebug.TransactOpts, _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FeedsConsumerDebug.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.TransferOwnership(&_FeedsConsumerDebug.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_FeedsConsumerDebug *FeedsConsumerDebugTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FeedsConsumerDebug.Contract.TransferOwnership(&_FeedsConsumerDebug.TransactOpts, to)
}

// FeedsConsumerDebugFeedReceivedIterator is returned from FilterFeedReceived and is used to iterate over the raw logs and unpacked data for FeedReceived events raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugFeedReceivedIterator struct {
	Event *FeedsConsumerDebugFeedReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FeedsConsumerDebugFeedReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeedsConsumerDebugFeedReceived)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FeedsConsumerDebugFeedReceived)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FeedsConsumerDebugFeedReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeedsConsumerDebugFeedReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeedsConsumerDebugFeedReceived represents a FeedReceived event raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugFeedReceived struct {
	FeedId    [32]byte
	Price     *big.Int
	Timestamp uint32
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeedReceived is a free log retrieval operation binding the contract event 0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d.
//
// Solidity: event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*FeedsConsumerDebugFeedReceivedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _FeedsConsumerDebug.contract.FilterLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebugFeedReceivedIterator{contract: _FeedsConsumerDebug.contract, event: "FeedReceived", logs: logs, sub: sub}, nil
}

// WatchFeedReceived is a free log subscription operation binding the contract event 0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d.
//
// Solidity: event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *FeedsConsumerDebugFeedReceived, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _FeedsConsumerDebug.contract.WatchLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeedsConsumerDebugFeedReceived)
				if err := _FeedsConsumerDebug.contract.UnpackLog(event, "FeedReceived", log); err != nil {
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

// ParseFeedReceived is a log parse operation binding the contract event 0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d.
//
// Solidity: event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) ParseFeedReceived(log types.Log) (*FeedsConsumerDebugFeedReceived, error) {
	event := new(FeedsConsumerDebugFeedReceived)
	if err := _FeedsConsumerDebug.contract.UnpackLog(event, "FeedReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeedsConsumerDebugOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugOwnershipTransferRequestedIterator struct {
	Event *FeedsConsumerDebugOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FeedsConsumerDebugOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeedsConsumerDebugOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FeedsConsumerDebugOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FeedsConsumerDebugOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeedsConsumerDebugOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeedsConsumerDebugOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeedsConsumerDebugOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeedsConsumerDebug.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebugOwnershipTransferRequestedIterator{contract: _FeedsConsumerDebug.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FeedsConsumerDebugOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeedsConsumerDebug.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeedsConsumerDebugOwnershipTransferRequested)
				if err := _FeedsConsumerDebug.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) ParseOwnershipTransferRequested(log types.Log) (*FeedsConsumerDebugOwnershipTransferRequested, error) {
	event := new(FeedsConsumerDebugOwnershipTransferRequested)
	if err := _FeedsConsumerDebug.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeedsConsumerDebugOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugOwnershipTransferredIterator struct {
	Event *FeedsConsumerDebugOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FeedsConsumerDebugOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeedsConsumerDebugOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FeedsConsumerDebugOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FeedsConsumerDebugOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeedsConsumerDebugOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeedsConsumerDebugOwnershipTransferred represents a OwnershipTransferred event raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeedsConsumerDebugOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeedsConsumerDebug.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebugOwnershipTransferredIterator{contract: _FeedsConsumerDebug.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeedsConsumerDebugOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeedsConsumerDebug.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeedsConsumerDebugOwnershipTransferred)
				if err := _FeedsConsumerDebug.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) ParseOwnershipTransferred(log types.Log) (*FeedsConsumerDebugOwnershipTransferred, error) {
	event := new(FeedsConsumerDebugOwnershipTransferred)
	if err := _FeedsConsumerDebug.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
