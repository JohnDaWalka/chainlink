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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"}],\"name\":\"UnauthorizedWorkflowName\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"}],\"name\":\"UnauthorizedWorkflowOwner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"DebugEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint224\",\"name\":\"price\",\"type\":\"uint224\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"FeedReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllFeeds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint224[]\",\"name\":\"\",\"type\":\"uint224[]\"},{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_allowedSendersList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_allowedWorkflowOwnersList\",\"type\":\"address[]\"},{\"internalType\":\"bytes10[]\",\"name\":\"_allowedWorkflowNamesList\",\"type\":\"bytes10[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6117f7806101576000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063805f21321161005b578063805f21321461018b5780638da5cb5b1461019e578063e3401711146101c6578063f2fde38b146101d957600080fd5b806301ffc9a71461008d57806312876798146100b557806331d98b3f146100cc57806379ba509714610181575b600080fd5b6100a061009b36600461125c565b6101ec565b60405190151581526020015b60405180910390f35b6100bd610285565b6040516100ac939291906112e7565b6101486100da366004611393565b6000908152600360209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168084527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169290910182905291565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909316835263ffffffff9091166020830152016100ac565b6101896104c1565b005b6101896101993660046113f5565b6105c3565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ac565b6101896101d43660046114a6565b610b15565b6101896101e7366004611540565b610f55565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f805f213200000000000000000000000000000000000000000000000000000000148061027f57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b6002546060908190819060008167ffffffffffffffff8111156102aa576102aa611576565b6040519080825280602002602001820160405280156102d3578160200160208202803683370190505b50905060008267ffffffffffffffff8111156102f1576102f1611576565b60405190808252806020026020018201604052801561031a578160200160208202803683370190505b50905060008367ffffffffffffffff81111561033857610338611576565b604051908082528060200260200182016040528015610361578160200160208202803683370190505b50905060005b848110156104b357600060028281548110610384576103846115a5565b6000918252602080832090910154808352600382526040928390208351808501909452547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116845263ffffffff7c0100000000000000000000000000000000000000000000000000000000909104169183019190915286519092508290879085908110610410576104106115a5565b6020026020010181815250508060000151858481518110610433576104336115a5565b60200260200101907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815250508060200151848481518110610494576104946115a5565b63ffffffff909216602092830291909101909101525050600101610367565b509196909550909350915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610547576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3360009081526005602052604090205460ff16610675577f56f074d292557f2e3c567d982816e0fb5b72100ff196892f8fbd23b8a90736796040516106399060208082526012908201527f556e617574686f72697a656453656e6465720000000000000000000000000000604082015260600190565b60405180910390a16040517f3fcc3f1700000000000000000000000000000000000000000000000000000000815233600482015260240161053e565b6000806106b786868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610f6992505050565b7fffffffffffffffffffff000000000000000000000000000000000000000000008216600090815260096020526040902054919350915060ff166107b2577f56f074d292557f2e3c567d982816e0fb5b72100ff196892f8fbd23b8a90736796040516107549060208082526018908201527f556e617574686f72697a6564576f726b666c6f774e616d650000000000000000604082015260600190565b60405180910390a16040517f4b942f800000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffff000000000000000000000000000000000000000000008316600482015260240161053e565b73ffffffffffffffffffffffffffffffffffffffff811660009081526007602052604090205460ff16610890577f56f074d292557f2e3c567d982816e0fb5b72100ff196892f8fbd23b8a907367960405161083e9060208082526019908201527f556e617574686f72697a6564576f726b666c6f774f776e657200000000000000604082015260600190565b60405180910390a16040517fbf24162300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161053e565b600061089e8486018661164c565b90507f56f074d292557f2e3c567d982816e0fb5b72100ff196892f8fbd23b8a90736796040516108ff906020808252600f908201527f6465636f64696e6720776f726b65640000000000000000000000000000000000604082015260600190565b60405180910390a160005b8151811015610b0b576040518060400160405280838381518110610930576109306115a5565b6020026020010151602001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168152602001838381518110610971576109716115a5565b60200260200101516040015163ffffffff168152506003600084848151811061099c5761099c6115a5565b602090810291909101810151518252818101929092526040016000208251929091015163ffffffff167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092169190911790558151600290839083908110610a2157610a216115a5565b60209081029190910181015151825460018101845560009384529190922001558151829082908110610a5557610a556115a5565b6020026020010151600001517f2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d838381518110610a9457610a946115a5565b602002602001015160200151848481518110610ab257610ab26115a5565b602002602001015160400151604051610afb9291907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216825263ffffffff16602082015260400190565b60405180910390a260010161090a565b5050505050505050565b610b1d610f7f565b60005b60045463ffffffff82161015610bbe5760006005600060048463ffffffff1681548110610b4f57610b4f6115a5565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610bb78161175e565b9050610b20565b5060005b63ffffffff8116861115610c665760016005600089898563ffffffff16818110610bee57610bee6115a5565b9050602002016020810190610c039190611540565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610c5f8161175e565b9050610bc2565b50610c73600487876110f7565b5060005b60065463ffffffff82161015610d155760006007600060068463ffffffff1681548110610ca657610ca66115a5565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610d0e8161175e565b9050610c77565b5060005b63ffffffff8116841115610dbd5760016007600087878563ffffffff16818110610d4557610d456115a5565b9050602002016020810190610d5a9190611540565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610db68161175e565b9050610d19565b50610dca600685856110f7565b5060005b60085463ffffffff82161015610e8b5760006009600060088463ffffffff1681548110610dfd57610dfd6115a5565b600091825260208083206003808404909101549206600a026101000a90910460b01b7fffffffffffffffffffff00000000000000000000000000000000000000000000168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610e848161175e565b9050610dce565b5060005b63ffffffff8116821115610f3f5760016009600085858563ffffffff16818110610ebb57610ebb6115a5565b9050602002016020810190610ed091906117a8565b7fffffffffffffffffffff00000000000000000000000000000000000000000000168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610f388161175e565b9050610e8f565b50610f4c6008838361117f565b50505050505050565b610f5d610f7f565b610f6681611002565b50565b6040810151604a90910151909160609190911c90565b60005473ffffffffffffffffffffffffffffffffffffffff163314611000576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161053e565b565b3373ffffffffffffffffffffffffffffffffffffffff821603611081576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161053e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821561116f579160200282015b8281111561116f5781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190611117565b5061117b929150611247565b5090565b8280548282559060005260206000209060020160039004810192821561116f5791602002820160005b8382111561120857833575ffffffffffffffffffffffffffffffffffffffffffff191683826101000a81548169ffffffffffffffffffff021916908360b01c02179055509260200192600a016020816009010492830192600103026111a8565b801561123e5782816101000a81549069ffffffffffffffffffff0219169055600a01602081600901049283019260010302611208565b505061117b9291505b5b8082111561117b5760008155600101611248565b60006020828403121561126e57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461129e57600080fd5b9392505050565b60008151808452602080850194506020840160005b838110156112dc57815163ffffffff16875295820195908201906001016112ba565b509495945050505050565b606080825284519082018190526000906020906080840190828801845b8281101561132057815184529284019290840190600101611304565b5050508381038285015285518082528683019183019060005b818110156113735783517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611339565b5050848103604086015261138781876112a5565b98975050505050505050565b6000602082840312156113a557600080fd5b5035919050565b60008083601f8401126113be57600080fd5b50813567ffffffffffffffff8111156113d657600080fd5b6020830191508360208285010111156113ee57600080fd5b9250929050565b6000806000806040858703121561140b57600080fd5b843567ffffffffffffffff8082111561142357600080fd5b61142f888389016113ac565b9096509450602087013591508082111561144857600080fd5b50611455878288016113ac565b95989497509550505050565b60008083601f84011261147357600080fd5b50813567ffffffffffffffff81111561148b57600080fd5b6020830191508360208260051b85010111156113ee57600080fd5b600080600080600080606087890312156114bf57600080fd5b863567ffffffffffffffff808211156114d757600080fd5b6114e38a838b01611461565b909850965060208901359150808211156114fc57600080fd5b6115088a838b01611461565b9096509450604089013591508082111561152157600080fd5b5061152e89828a01611461565b979a9699509497509295939492505050565b60006020828403121561155257600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461129e57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6040516060810167ffffffffffffffff811182821017156115f7576115f7611576565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561164457611644611576565b604052919050565b6000602080838503121561165f57600080fd5b823567ffffffffffffffff8082111561167757600080fd5b818501915085601f83011261168b57600080fd5b81358181111561169d5761169d611576565b6116ab848260051b016115fd565b818152848101925060609182028401850191888311156116ca57600080fd5b938501935b828510156117525780858a0312156116e75760008081fd5b6116ef6115d4565b85358152868601357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146117225760008081fd5b8188015260408681013563ffffffff8116811461173f5760008081fd5b90820152845293840193928501926116cf565b50979650505050505050565b600063ffffffff80831681810361179e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b6000602082840312156117ba57600080fd5b81357fffffffffffffffffffff000000000000000000000000000000000000000000008116811461129e57600080fdfea164736f6c6343000818000a",
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

// FeedsConsumerDebugDebugEventIterator is returned from FilterDebugEvent and is used to iterate over the raw logs and unpacked data for DebugEvent events raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugDebugEventIterator struct {
	Event *FeedsConsumerDebugDebugEvent // Event containing the contract specifics and raw log

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
func (it *FeedsConsumerDebugDebugEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeedsConsumerDebugDebugEvent)
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
		it.Event = new(FeedsConsumerDebugDebugEvent)
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
func (it *FeedsConsumerDebugDebugEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeedsConsumerDebugDebugEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeedsConsumerDebugDebugEvent represents a DebugEvent event raised by the FeedsConsumerDebug contract.
type FeedsConsumerDebugDebugEvent struct {
	Reason string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDebugEvent is a free log retrieval operation binding the contract event 0x56f074d292557f2e3c567d982816e0fb5b72100ff196892f8fbd23b8a9073679.
//
// Solidity: event DebugEvent(string reason)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) FilterDebugEvent(opts *bind.FilterOpts) (*FeedsConsumerDebugDebugEventIterator, error) {

	logs, sub, err := _FeedsConsumerDebug.contract.FilterLogs(opts, "DebugEvent")
	if err != nil {
		return nil, err
	}
	return &FeedsConsumerDebugDebugEventIterator{contract: _FeedsConsumerDebug.contract, event: "DebugEvent", logs: logs, sub: sub}, nil
}

// WatchDebugEvent is a free log subscription operation binding the contract event 0x56f074d292557f2e3c567d982816e0fb5b72100ff196892f8fbd23b8a9073679.
//
// Solidity: event DebugEvent(string reason)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) WatchDebugEvent(opts *bind.WatchOpts, sink chan<- *FeedsConsumerDebugDebugEvent) (event.Subscription, error) {

	logs, sub, err := _FeedsConsumerDebug.contract.WatchLogs(opts, "DebugEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeedsConsumerDebugDebugEvent)
				if err := _FeedsConsumerDebug.contract.UnpackLog(event, "DebugEvent", log); err != nil {
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

// ParseDebugEvent is a log parse operation binding the contract event 0x56f074d292557f2e3c567d982816e0fb5b72100ff196892f8fbd23b8a9073679.
//
// Solidity: event DebugEvent(string reason)
func (_FeedsConsumerDebug *FeedsConsumerDebugFilterer) ParseDebugEvent(log types.Log) (*FeedsConsumerDebugDebugEvent, error) {
	event := new(FeedsConsumerDebugDebugEvent)
	if err := _FeedsConsumerDebug.contract.UnpackLog(event, "DebugEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
