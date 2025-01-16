// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feeds_consumer

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

var KeystoneFeedsConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"}],\"name\":\"UnauthorizedWorkflowName\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"}],\"name\":\"UnauthorizedWorkflowOwner\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint224\",\"name\":\"price\",\"type\":\"uint224\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"FeedReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllFeeds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint224[]\",\"name\":\"\",\"type\":\"uint224[]\"},{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"feedId\",\"type\":\"bytes32\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint224\",\"name\":\"\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_allowedSendersList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_allowedWorkflowOwnersList\",\"type\":\"address[]\"},{\"internalType\":\"bytes10[]\",\"name\":\"_allowedWorkflowNamesList\",\"type\":\"bytes10[]\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61165b806101576000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063805f21321161005b578063805f21321461018b5780638da5cb5b1461019e578063e3401711146101c6578063f2fde38b146101d957600080fd5b806301ffc9a71461008d57806312876798146100b557806331d98b3f146100cc57806379ba509714610181575b600080fd5b6100a061009b3660046110c0565b6101ec565b60405190151581526020015b60405180910390f35b6100bd610285565b6040516100ac9392919061114b565b6101486100da3660046111f7565b6000908152600360209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168084527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169290910182905291565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909316835263ffffffff9091166020830152016100ac565b6101896104c1565b005b610189610199366004611259565b6105c3565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ac565b6101896101d436600461130a565b610979565b6101896101e73660046113a4565b610db9565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f805f213200000000000000000000000000000000000000000000000000000000148061027f57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b6002546060908190819060008167ffffffffffffffff8111156102aa576102aa6113da565b6040519080825280602002602001820160405280156102d3578160200160208202803683370190505b50905060008267ffffffffffffffff8111156102f1576102f16113da565b60405190808252806020026020018201604052801561031a578160200160208202803683370190505b50905060008367ffffffffffffffff811115610338576103386113da565b604051908082528060200260200182016040528015610361578160200160208202803683370190505b50905060005b848110156104b35760006002828154811061038457610384611409565b6000918252602080832090910154808352600382526040928390208351808501909452547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116845263ffffffff7c010000000000000000000000000000000000000000000000000000000090910416918301919091528651909250829087908590811061041057610410611409565b602002602001018181525050806000015185848151811061043357610433611409565b60200260200101907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681525050806020015184848151811061049457610494611409565b63ffffffff909216602092830291909101909101525050600101610367565b509196909550909350915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610547576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3360009081526005602052604090205460ff1661060e576040517f3fcc3f1700000000000000000000000000000000000000000000000000000000815233600482015260240161053e565b60008061065086868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610dcd92505050565b7fffffffffffffffffffff000000000000000000000000000000000000000000008216600090815260096020526040902054919350915060ff166106e4576040517f4b942f800000000000000000000000000000000000000000000000000000000081527fffffffffffffffffffff000000000000000000000000000000000000000000008316600482015260240161053e565b73ffffffffffffffffffffffffffffffffffffffff811660009081526007602052604090205460ff1661075b576040517fbf24162300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161053e565b6000610769848601866114b0565b905060005b815181101561096f57604051806040016040528083838151811061079457610794611409565b6020026020010151602001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020018383815181106107d5576107d5611409565b60200260200101516040015163ffffffff168152506003600084848151811061080057610800611409565b602090810291909101810151518252818101929092526040016000208251929091015163ffffffff167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909216919091179055815160029083908390811061088557610885611409565b602090810291909101810151518254600181018455600093845291909220015581518290829081106108b9576108b9611409565b6020026020010151600001517f2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d8383815181106108f8576108f8611409565b60200260200101516020015184848151811061091657610916611409565b60200260200101516040015160405161095f9291907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216825263ffffffff16602082015260400190565b60405180910390a260010161076e565b5050505050505050565b610981610de3565b60005b60045463ffffffff82161015610a225760006005600060048463ffffffff16815481106109b3576109b3611409565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610a1b816115c2565b9050610984565b5060005b63ffffffff8116861115610aca5760016005600089898563ffffffff16818110610a5257610a52611409565b9050602002016020810190610a6791906113a4565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610ac3816115c2565b9050610a26565b50610ad760048787610f5b565b5060005b60065463ffffffff82161015610b795760006007600060068463ffffffff1681548110610b0a57610b0a611409565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610b72816115c2565b9050610adb565b5060005b63ffffffff8116841115610c215760016007600087878563ffffffff16818110610ba957610ba9611409565b9050602002016020810190610bbe91906113a4565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610c1a816115c2565b9050610b7d565b50610c2e60068585610f5b565b5060005b60085463ffffffff82161015610cef5760006009600060088463ffffffff1681548110610c6157610c61611409565b600091825260208083206003808404909101549206600a026101000a90910460b01b7fffffffffffffffffffff00000000000000000000000000000000000000000000168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610ce8816115c2565b9050610c32565b5060005b63ffffffff8116821115610da35760016009600085858563ffffffff16818110610d1f57610d1f611409565b9050602002016020810190610d34919061160c565b7fffffffffffffffffffff00000000000000000000000000000000000000000000168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055610d9c816115c2565b9050610cf3565b50610db060088383610fe3565b50505050505050565b610dc1610de3565b610dca81610e66565b50565b6040810151604a90910151909160609190911c90565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e64576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161053e565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610ee5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161053e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610fd3579160200282015b82811115610fd35781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610f7b565b50610fdf9291506110ab565b5090565b82805482825590600052602060002090600201600390048101928215610fd35791602002820160005b8382111561106c57833575ffffffffffffffffffffffffffffffffffffffffffff191683826101000a81548169ffffffffffffffffffff021916908360b01c02179055509260200192600a0160208160090104928301926001030261100c565b80156110a25782816101000a81549069ffffffffffffffffffff0219169055600a0160208160090104928301926001030261106c565b5050610fdf9291505b5b80821115610fdf57600081556001016110ac565b6000602082840312156110d257600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461110257600080fd5b9392505050565b60008151808452602080850194506020840160005b8381101561114057815163ffffffff168752958201959082019060010161111e565b509495945050505050565b606080825284519082018190526000906020906080840190828801845b8281101561118457815184529284019290840190600101611168565b5050508381038285015285518082528683019183019060005b818110156111d75783517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161119d565b505084810360408601526111eb8187611109565b98975050505050505050565b60006020828403121561120957600080fd5b5035919050565b60008083601f84011261122257600080fd5b50813567ffffffffffffffff81111561123a57600080fd5b60208301915083602082850101111561125257600080fd5b9250929050565b6000806000806040858703121561126f57600080fd5b843567ffffffffffffffff8082111561128757600080fd5b61129388838901611210565b909650945060208701359150808211156112ac57600080fd5b506112b987828801611210565b95989497509550505050565b60008083601f8401126112d757600080fd5b50813567ffffffffffffffff8111156112ef57600080fd5b6020830191508360208260051b850101111561125257600080fd5b6000806000806000806060878903121561132357600080fd5b863567ffffffffffffffff8082111561133b57600080fd5b6113478a838b016112c5565b9098509650602089013591508082111561136057600080fd5b61136c8a838b016112c5565b9096509450604089013591508082111561138557600080fd5b5061139289828a016112c5565b979a9699509497509295939492505050565b6000602082840312156113b657600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461110257600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6040516060810167ffffffffffffffff8111828210171561145b5761145b6113da565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156114a8576114a86113da565b604052919050565b600060208083850312156114c357600080fd5b823567ffffffffffffffff808211156114db57600080fd5b818501915085601f8301126114ef57600080fd5b813581811115611501576115016113da565b61150f848260051b01611461565b8181528481019250606091820284018501918883111561152e57600080fd5b938501935b828510156115b65780858a03121561154b5760008081fd5b611553611438565b85358152868601357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146115865760008081fd5b8188015260408681013563ffffffff811681146115a35760008081fd5b9082015284529384019392850192611533565b50979650505050505050565b600063ffffffff808316818103611602577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b60006020828403121561161e57600080fd5b81357fffffffffffffffffffff000000000000000000000000000000000000000000008116811461110257600080fdfea164736f6c6343000818000a",
}

var KeystoneFeedsConsumerABI = KeystoneFeedsConsumerMetaData.ABI

var KeystoneFeedsConsumerBin = KeystoneFeedsConsumerMetaData.Bin

func DeployKeystoneFeedsConsumer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeystoneFeedsConsumer, error) {
	parsed, err := KeystoneFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneFeedsConsumerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeystoneFeedsConsumer{address: address, abi: *parsed, KeystoneFeedsConsumerCaller: KeystoneFeedsConsumerCaller{contract: contract}, KeystoneFeedsConsumerTransactor: KeystoneFeedsConsumerTransactor{contract: contract}, KeystoneFeedsConsumerFilterer: KeystoneFeedsConsumerFilterer{contract: contract}}, nil
}

type KeystoneFeedsConsumer struct {
	address common.Address
	abi     abi.ABI
	KeystoneFeedsConsumerCaller
	KeystoneFeedsConsumerTransactor
	KeystoneFeedsConsumerFilterer
}

type KeystoneFeedsConsumerCaller struct {
	contract *bind.BoundContract
}

type KeystoneFeedsConsumerTransactor struct {
	contract *bind.BoundContract
}

type KeystoneFeedsConsumerFilterer struct {
	contract *bind.BoundContract
}

type KeystoneFeedsConsumerSession struct {
	Contract     *KeystoneFeedsConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeystoneFeedsConsumerCallerSession struct {
	Contract *KeystoneFeedsConsumerCaller
	CallOpts bind.CallOpts
}

type KeystoneFeedsConsumerTransactorSession struct {
	Contract     *KeystoneFeedsConsumerTransactor
	TransactOpts bind.TransactOpts
}

type KeystoneFeedsConsumerRaw struct {
	Contract *KeystoneFeedsConsumer
}

type KeystoneFeedsConsumerCallerRaw struct {
	Contract *KeystoneFeedsConsumerCaller
}

type KeystoneFeedsConsumerTransactorRaw struct {
	Contract *KeystoneFeedsConsumerTransactor
}

func NewKeystoneFeedsConsumer(address common.Address, backend bind.ContractBackend) (*KeystoneFeedsConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(KeystoneFeedsConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeystoneFeedsConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumer{address: address, abi: abi, KeystoneFeedsConsumerCaller: KeystoneFeedsConsumerCaller{contract: contract}, KeystoneFeedsConsumerTransactor: KeystoneFeedsConsumerTransactor{contract: contract}, KeystoneFeedsConsumerFilterer: KeystoneFeedsConsumerFilterer{contract: contract}}, nil
}

func NewKeystoneFeedsConsumerCaller(address common.Address, caller bind.ContractCaller) (*KeystoneFeedsConsumerCaller, error) {
	contract, err := bindKeystoneFeedsConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerCaller{contract: contract}, nil
}

func NewKeystoneFeedsConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*KeystoneFeedsConsumerTransactor, error) {
	contract, err := bindKeystoneFeedsConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerTransactor{contract: contract}, nil
}

func NewKeystoneFeedsConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*KeystoneFeedsConsumerFilterer, error) {
	contract, err := bindKeystoneFeedsConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerFilterer{contract: contract}, nil
}

func bindKeystoneFeedsConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeystoneFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneFeedsConsumer.Contract.KeystoneFeedsConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.KeystoneFeedsConsumerTransactor.contract.Transfer(opts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.KeystoneFeedsConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneFeedsConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.contract.Transfer(opts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCaller) GetAllFeeds(opts *bind.CallOpts) ([][32]byte, []*big.Int, []uint32, error) {
	var out []interface{}
	err := _KeystoneFeedsConsumer.contract.Call(opts, &out, "getAllFeeds")

	if err != nil {
		return *new([][32]byte), *new([]*big.Int), *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	out2 := *abi.ConvertType(out[2], new([]uint32)).(*[]uint32)

	return out0, out1, out2, err

}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) GetAllFeeds() ([][32]byte, []*big.Int, []uint32, error) {
	return _KeystoneFeedsConsumer.Contract.GetAllFeeds(&_KeystoneFeedsConsumer.CallOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerSession) GetAllFeeds() ([][32]byte, []*big.Int, []uint32, error) {
	return _KeystoneFeedsConsumer.Contract.GetAllFeeds(&_KeystoneFeedsConsumer.CallOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCaller) GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error) {
	var out []interface{}
	err := _KeystoneFeedsConsumer.contract.Call(opts, &out, "getPrice", feedId)

	if err != nil {
		return *new(*big.Int), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return out0, out1, err

}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _KeystoneFeedsConsumer.Contract.GetPrice(&_KeystoneFeedsConsumer.CallOpts, feedId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _KeystoneFeedsConsumer.Contract.GetPrice(&_KeystoneFeedsConsumer.CallOpts, feedId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneFeedsConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) Owner() (common.Address, error) {
	return _KeystoneFeedsConsumer.Contract.Owner(&_KeystoneFeedsConsumer.CallOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerSession) Owner() (common.Address, error) {
	return _KeystoneFeedsConsumer.Contract.Owner(&_KeystoneFeedsConsumer.CallOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _KeystoneFeedsConsumer.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _KeystoneFeedsConsumer.Contract.SupportsInterface(&_KeystoneFeedsConsumer.CallOpts, interfaceId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _KeystoneFeedsConsumer.Contract.SupportsInterface(&_KeystoneFeedsConsumer.CallOpts, interfaceId)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.AcceptOwnership(&_KeystoneFeedsConsumer.TransactOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.AcceptOwnership(&_KeystoneFeedsConsumer.TransactOpts)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) OnReport(opts *bind.TransactOpts, metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "onReport", metadata, rawReport)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.OnReport(&_KeystoneFeedsConsumer.TransactOpts, metadata, rawReport)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) OnReport(metadata []byte, rawReport []byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.OnReport(&_KeystoneFeedsConsumer.TransactOpts, metadata, rawReport)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) SetConfig(opts *bind.TransactOpts, _allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "setConfig", _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) SetConfig(_allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.SetConfig(&_KeystoneFeedsConsumer.TransactOpts, _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) SetConfig(_allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.SetConfig(&_KeystoneFeedsConsumer.TransactOpts, _allowedSendersList, _allowedWorkflowOwnersList, _allowedWorkflowNamesList)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.TransferOwnership(&_KeystoneFeedsConsumer.TransactOpts, to)
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneFeedsConsumer.Contract.TransferOwnership(&_KeystoneFeedsConsumer.TransactOpts, to)
}

type KeystoneFeedsConsumerFeedReceivedIterator struct {
	Event *KeystoneFeedsConsumerFeedReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneFeedsConsumerFeedReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneFeedsConsumerFeedReceived)
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
		it.Event = new(KeystoneFeedsConsumerFeedReceived)
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

func (it *KeystoneFeedsConsumerFeedReceivedIterator) Error() error {
	return it.fail
}

func (it *KeystoneFeedsConsumerFeedReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneFeedsConsumerFeedReceived struct {
	FeedId    [32]byte
	Price     *big.Int
	Timestamp uint32
	Raw       types.Log
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*KeystoneFeedsConsumerFeedReceivedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.FilterLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerFeedReceivedIterator{contract: _KeystoneFeedsConsumer.contract, event: "FeedReceived", logs: logs, sub: sub}, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerFeedReceived, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.WatchLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneFeedsConsumerFeedReceived)
				if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
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

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) ParseFeedReceived(log types.Log) (*KeystoneFeedsConsumerFeedReceived, error) {
	event := new(KeystoneFeedsConsumerFeedReceived)
	if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneFeedsConsumerOwnershipTransferRequestedIterator struct {
	Event *KeystoneFeedsConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneFeedsConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneFeedsConsumerOwnershipTransferRequested)
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
		it.Event = new(KeystoneFeedsConsumerOwnershipTransferRequested)
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

func (it *KeystoneFeedsConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeystoneFeedsConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneFeedsConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerOwnershipTransferRequestedIterator{contract: _KeystoneFeedsConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneFeedsConsumerOwnershipTransferRequested)
				if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferRequested, error) {
	event := new(KeystoneFeedsConsumerOwnershipTransferRequested)
	if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneFeedsConsumerOwnershipTransferredIterator struct {
	Event *KeystoneFeedsConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneFeedsConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneFeedsConsumerOwnershipTransferred)
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
		it.Event = new(KeystoneFeedsConsumerOwnershipTransferred)
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

func (it *KeystoneFeedsConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeystoneFeedsConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneFeedsConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneFeedsConsumerOwnershipTransferredIterator{contract: _KeystoneFeedsConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneFeedsConsumerOwnershipTransferred)
				if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferred, error) {
	event := new(KeystoneFeedsConsumerOwnershipTransferred)
	if err := _KeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeystoneFeedsConsumer.abi.Events["FeedReceived"].ID:
		return _KeystoneFeedsConsumer.ParseFeedReceived(log)
	case _KeystoneFeedsConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _KeystoneFeedsConsumer.ParseOwnershipTransferRequested(log)
	case _KeystoneFeedsConsumer.abi.Events["OwnershipTransferred"].ID:
		return _KeystoneFeedsConsumer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeystoneFeedsConsumerFeedReceived) Topic() common.Hash {
	return common.HexToHash("0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d")
}

func (KeystoneFeedsConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeystoneFeedsConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_KeystoneFeedsConsumer *KeystoneFeedsConsumer) Address() common.Address {
	return _KeystoneFeedsConsumer.address
}

type KeystoneFeedsConsumerInterface interface {
	GetAllFeeds(opts *bind.CallOpts) ([][32]byte, []*big.Int, []uint32, error)

	GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnReport(opts *bind.TransactOpts, metadata []byte, rawReport []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _allowedSendersList []common.Address, _allowedWorkflowOwnersList []common.Address, _allowedWorkflowNamesList [][10]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*KeystoneFeedsConsumerFeedReceivedIterator, error)

	WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerFeedReceived, feedId [][32]byte) (event.Subscription, error)

	ParseFeedReceived(log types.Log) (*KeystoneFeedsConsumerFeedReceived, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneFeedsConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneFeedsConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeystoneFeedsConsumerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
