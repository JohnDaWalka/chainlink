// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package permissionless_feeds_consumer

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

var PermissionlessKeystoneFeedsConsumerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getPrice\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint224\",\"internalType\":\"uint224\"},{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FeedReceived\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"price\",\"type\":\"uint224\",\"indexed\":false,\"internalType\":\"uint224\"},{\"name\":\"timestamp\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6109fa806101576000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063805f213211610050578063805f21321461015e5780638da5cb5b14610171578063f2fde38b1461019957600080fd5b806301ffc9a71461007757806331d98b3f1461009f57806379ba509714610154575b600080fd5b61008a6100853660046106b8565b6101ac565b60405190151581526020015b60405180910390f35b61011b6100ad366004610701565b6000908152600260209081526040918290208251808401909352547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168084527c010000000000000000000000000000000000000000000000000000000090910463ffffffff169290910182905291565b604080517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909316835263ffffffff909116602083015201610096565b61015c610245565b005b61015c61016c366004610763565b610347565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610096565b61015c6101a73660046107cf565b61052c565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f805f213200000000000000000000000000000000000000000000000000000000148061023f57507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610355828401846108ac565b905060005b8151811015610524576040518060400160405280838381518110610380576103806109be565b6020026020010151602001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681526020018383815181106103c1576103c16109be565b60200260200101516040015163ffffffff16815250600260008484815181106103ec576103ec6109be565b602090810291909101810151518252818101929092526040016000208251929091015163ffffffff167c0100000000000000000000000000000000000000000000000000000000027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff909216919091179055815182908290811061046e5761046e6109be565b6020026020010151600001517f2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d8383815181106104ad576104ad6109be565b6020026020010151602001518484815181106104cb576104cb6109be565b6020026020010151604001516040516105149291907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216825263ffffffff16602082015260400190565b60405180910390a260010161035a565b505050505050565b610534610540565b61053d816105c3565b50565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105c1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102c2565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610642576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102c2565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156106ca57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146106fa57600080fd5b9392505050565b60006020828403121561071357600080fd5b5035919050565b60008083601f84011261072c57600080fd5b50813567ffffffffffffffff81111561074457600080fd5b60208301915083602082850101111561075c57600080fd5b9250929050565b6000806000806040858703121561077957600080fd5b843567ffffffffffffffff8082111561079157600080fd5b61079d8883890161071a565b909650945060208701359150808211156107b657600080fd5b506107c38782880161071a565b95989497509550505050565b6000602082840312156107e157600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146106fa57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff8111828210171561085757610857610805565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156108a4576108a4610805565b604052919050565b600060208083850312156108bf57600080fd5b823567ffffffffffffffff808211156108d757600080fd5b818501915085601f8301126108eb57600080fd5b8135818111156108fd576108fd610805565b61090b848260051b0161085d565b8181528481019250606091820284018501918883111561092a57600080fd5b938501935b828510156109b25780858a0312156109475760008081fd5b61094f610834565b85358152868601357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146109825760008081fd5b8188015260408681013563ffffffff8116811461099f5760008081fd5b908201528452938401939285019261092f565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c6343000818000a",
}

var PermissionlessKeystoneFeedsConsumerABI = PermissionlessKeystoneFeedsConsumerMetaData.ABI

var PermissionlessKeystoneFeedsConsumerBin = PermissionlessKeystoneFeedsConsumerMetaData.Bin

func DeployPermissionlessKeystoneFeedsConsumer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PermissionlessKeystoneFeedsConsumer, error) {
	parsed, err := PermissionlessKeystoneFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PermissionlessKeystoneFeedsConsumerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PermissionlessKeystoneFeedsConsumer{address: address, abi: *parsed, PermissionlessKeystoneFeedsConsumerCaller: PermissionlessKeystoneFeedsConsumerCaller{contract: contract}, PermissionlessKeystoneFeedsConsumerTransactor: PermissionlessKeystoneFeedsConsumerTransactor{contract: contract}, PermissionlessKeystoneFeedsConsumerFilterer: PermissionlessKeystoneFeedsConsumerFilterer{contract: contract}}, nil
}

type PermissionlessKeystoneFeedsConsumer struct {
	address common.Address
	abi     abi.ABI
	PermissionlessKeystoneFeedsConsumerCaller
	PermissionlessKeystoneFeedsConsumerTransactor
	PermissionlessKeystoneFeedsConsumerFilterer
}

type PermissionlessKeystoneFeedsConsumerCaller struct {
	contract *bind.BoundContract
}

type PermissionlessKeystoneFeedsConsumerTransactor struct {
	contract *bind.BoundContract
}

type PermissionlessKeystoneFeedsConsumerFilterer struct {
	contract *bind.BoundContract
}

type PermissionlessKeystoneFeedsConsumerSession struct {
	Contract     *PermissionlessKeystoneFeedsConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PermissionlessKeystoneFeedsConsumerCallerSession struct {
	Contract *PermissionlessKeystoneFeedsConsumerCaller
	CallOpts bind.CallOpts
}

type PermissionlessKeystoneFeedsConsumerTransactorSession struct {
	Contract     *PermissionlessKeystoneFeedsConsumerTransactor
	TransactOpts bind.TransactOpts
}

type PermissionlessKeystoneFeedsConsumerRaw struct {
	Contract *PermissionlessKeystoneFeedsConsumer
}

type PermissionlessKeystoneFeedsConsumerCallerRaw struct {
	Contract *PermissionlessKeystoneFeedsConsumerCaller
}

type PermissionlessKeystoneFeedsConsumerTransactorRaw struct {
	Contract *PermissionlessKeystoneFeedsConsumerTransactor
}

func NewPermissionlessKeystoneFeedsConsumer(address common.Address, backend bind.ContractBackend) (*PermissionlessKeystoneFeedsConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(PermissionlessKeystoneFeedsConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPermissionlessKeystoneFeedsConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PermissionlessKeystoneFeedsConsumer{address: address, abi: abi, PermissionlessKeystoneFeedsConsumerCaller: PermissionlessKeystoneFeedsConsumerCaller{contract: contract}, PermissionlessKeystoneFeedsConsumerTransactor: PermissionlessKeystoneFeedsConsumerTransactor{contract: contract}, PermissionlessKeystoneFeedsConsumerFilterer: PermissionlessKeystoneFeedsConsumerFilterer{contract: contract}}, nil
}

func NewPermissionlessKeystoneFeedsConsumerCaller(address common.Address, caller bind.ContractCaller) (*PermissionlessKeystoneFeedsConsumerCaller, error) {
	contract, err := bindPermissionlessKeystoneFeedsConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PermissionlessKeystoneFeedsConsumerCaller{contract: contract}, nil
}

func NewPermissionlessKeystoneFeedsConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*PermissionlessKeystoneFeedsConsumerTransactor, error) {
	contract, err := bindPermissionlessKeystoneFeedsConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PermissionlessKeystoneFeedsConsumerTransactor{contract: contract}, nil
}

func NewPermissionlessKeystoneFeedsConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*PermissionlessKeystoneFeedsConsumerFilterer, error) {
	contract, err := bindPermissionlessKeystoneFeedsConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PermissionlessKeystoneFeedsConsumerFilterer{contract: contract}, nil
}

func bindPermissionlessKeystoneFeedsConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PermissionlessKeystoneFeedsConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PermissionlessKeystoneFeedsConsumer.Contract.PermissionlessKeystoneFeedsConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.PermissionlessKeystoneFeedsConsumerTransactor.contract.Transfer(opts)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.PermissionlessKeystoneFeedsConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PermissionlessKeystoneFeedsConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.contract.Transfer(opts)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerCaller) GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error) {
	var out []interface{}
	err := _PermissionlessKeystoneFeedsConsumer.contract.Call(opts, &out, "getPrice", feedId)

	if err != nil {
		return *new(*big.Int), *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return out0, out1, err

}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.GetPrice(&_PermissionlessKeystoneFeedsConsumer.CallOpts, feedId)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerCallerSession) GetPrice(feedId [32]byte) (*big.Int, uint32, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.GetPrice(&_PermissionlessKeystoneFeedsConsumer.CallOpts, feedId)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PermissionlessKeystoneFeedsConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerSession) Owner() (common.Address, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.Owner(&_PermissionlessKeystoneFeedsConsumer.CallOpts)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerCallerSession) Owner() (common.Address, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.Owner(&_PermissionlessKeystoneFeedsConsumer.CallOpts)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _PermissionlessKeystoneFeedsConsumer.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.SupportsInterface(&_PermissionlessKeystoneFeedsConsumer.CallOpts, interfaceId)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.SupportsInterface(&_PermissionlessKeystoneFeedsConsumer.CallOpts, interfaceId)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.AcceptOwnership(&_PermissionlessKeystoneFeedsConsumer.TransactOpts)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.AcceptOwnership(&_PermissionlessKeystoneFeedsConsumer.TransactOpts)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactor) OnReport(opts *bind.TransactOpts, arg0 []byte, rawReport []byte) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.contract.Transact(opts, "onReport", arg0, rawReport)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerSession) OnReport(arg0 []byte, rawReport []byte) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.OnReport(&_PermissionlessKeystoneFeedsConsumer.TransactOpts, arg0, rawReport)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactorSession) OnReport(arg0 []byte, rawReport []byte) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.OnReport(&_PermissionlessKeystoneFeedsConsumer.TransactOpts, arg0, rawReport)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.TransferOwnership(&_PermissionlessKeystoneFeedsConsumer.TransactOpts, to)
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PermissionlessKeystoneFeedsConsumer.Contract.TransferOwnership(&_PermissionlessKeystoneFeedsConsumer.TransactOpts, to)
}

type PermissionlessKeystoneFeedsConsumerFeedReceivedIterator struct {
	Event *PermissionlessKeystoneFeedsConsumerFeedReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PermissionlessKeystoneFeedsConsumerFeedReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionlessKeystoneFeedsConsumerFeedReceived)
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
		it.Event = new(PermissionlessKeystoneFeedsConsumerFeedReceived)
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

func (it *PermissionlessKeystoneFeedsConsumerFeedReceivedIterator) Error() error {
	return it.fail
}

func (it *PermissionlessKeystoneFeedsConsumerFeedReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PermissionlessKeystoneFeedsConsumerFeedReceived struct {
	FeedId    [32]byte
	Price     *big.Int
	Timestamp uint32
	Raw       types.Log
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*PermissionlessKeystoneFeedsConsumerFeedReceivedIterator, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _PermissionlessKeystoneFeedsConsumer.contract.FilterLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return &PermissionlessKeystoneFeedsConsumerFeedReceivedIterator{contract: _PermissionlessKeystoneFeedsConsumer.contract, event: "FeedReceived", logs: logs, sub: sub}, nil
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *PermissionlessKeystoneFeedsConsumerFeedReceived, feedId [][32]byte) (event.Subscription, error) {

	var feedIdRule []interface{}
	for _, feedIdItem := range feedId {
		feedIdRule = append(feedIdRule, feedIdItem)
	}

	logs, sub, err := _PermissionlessKeystoneFeedsConsumer.contract.WatchLogs(opts, "FeedReceived", feedIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PermissionlessKeystoneFeedsConsumerFeedReceived)
				if err := _PermissionlessKeystoneFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
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

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) ParseFeedReceived(log types.Log) (*PermissionlessKeystoneFeedsConsumerFeedReceived, error) {
	event := new(PermissionlessKeystoneFeedsConsumerFeedReceived)
	if err := _PermissionlessKeystoneFeedsConsumer.contract.UnpackLog(event, "FeedReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PermissionlessKeystoneFeedsConsumerOwnershipTransferRequestedIterator struct {
	Event *PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PermissionlessKeystoneFeedsConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested)
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
		it.Event = new(PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested)
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

func (it *PermissionlessKeystoneFeedsConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *PermissionlessKeystoneFeedsConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessKeystoneFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PermissionlessKeystoneFeedsConsumerOwnershipTransferRequestedIterator{contract: _PermissionlessKeystoneFeedsConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessKeystoneFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested)
				if err := _PermissionlessKeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested, error) {
	event := new(PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested)
	if err := _PermissionlessKeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PermissionlessKeystoneFeedsConsumerOwnershipTransferredIterator struct {
	Event *PermissionlessKeystoneFeedsConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PermissionlessKeystoneFeedsConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PermissionlessKeystoneFeedsConsumerOwnershipTransferred)
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
		it.Event = new(PermissionlessKeystoneFeedsConsumerOwnershipTransferred)
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

func (it *PermissionlessKeystoneFeedsConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *PermissionlessKeystoneFeedsConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PermissionlessKeystoneFeedsConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessKeystoneFeedsConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PermissionlessKeystoneFeedsConsumerOwnershipTransferredIterator{contract: _PermissionlessKeystoneFeedsConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PermissionlessKeystoneFeedsConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PermissionlessKeystoneFeedsConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PermissionlessKeystoneFeedsConsumerOwnershipTransferred)
				if err := _PermissionlessKeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferred, error) {
	event := new(PermissionlessKeystoneFeedsConsumerOwnershipTransferred)
	if err := _PermissionlessKeystoneFeedsConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _PermissionlessKeystoneFeedsConsumer.abi.Events["FeedReceived"].ID:
		return _PermissionlessKeystoneFeedsConsumer.ParseFeedReceived(log)
	case _PermissionlessKeystoneFeedsConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _PermissionlessKeystoneFeedsConsumer.ParseOwnershipTransferRequested(log)
	case _PermissionlessKeystoneFeedsConsumer.abi.Events["OwnershipTransferred"].ID:
		return _PermissionlessKeystoneFeedsConsumer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (PermissionlessKeystoneFeedsConsumerFeedReceived) Topic() common.Hash {
	return common.HexToHash("0x2c30f5cb3caf4239d0f994ce539d7ef24817fa550169c388e3a110f02e40197d")
}

func (PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (PermissionlessKeystoneFeedsConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_PermissionlessKeystoneFeedsConsumer *PermissionlessKeystoneFeedsConsumer) Address() common.Address {
	return _PermissionlessKeystoneFeedsConsumer.address
}

type PermissionlessKeystoneFeedsConsumerInterface interface {
	GetPrice(opts *bind.CallOpts, feedId [32]byte) (*big.Int, uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnReport(opts *bind.TransactOpts, arg0 []byte, rawReport []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterFeedReceived(opts *bind.FilterOpts, feedId [][32]byte) (*PermissionlessKeystoneFeedsConsumerFeedReceivedIterator, error)

	WatchFeedReceived(opts *bind.WatchOpts, sink chan<- *PermissionlessKeystoneFeedsConsumerFeedReceived, feedId [][32]byte) (event.Subscription, error)

	ParseFeedReceived(log types.Log) (*PermissionlessKeystoneFeedsConsumerFeedReceived, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PermissionlessKeystoneFeedsConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*PermissionlessKeystoneFeedsConsumerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
