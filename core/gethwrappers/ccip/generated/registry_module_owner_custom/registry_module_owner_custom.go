package registry_module_owner_custom

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

var RegistryModuleOwnerCustomMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"CanOnlySelfRegister\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"RequiredRoleNotFound\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"administrator\",\"type\":\"address\"}],\"name\":\"AdministratorRegistered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"registerAccessControlDefaultAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"registerAdminViaGetCCIPAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"registerAdminViaOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161064a38038061064a83398101604081905261002f91610067565b6001600160a01b03811661005657604051639fabe1c160e01b815260040160405180910390fd5b6001600160a01b0316608052610097565b60006020828403121561007957600080fd5b81516001600160a01b038116811461009057600080fd5b9392505050565b6080516105986100b260003960006103db01526105986000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063181f5a771461005157806369c0081e146100a357806396ea2f7a146100b8578063ff12c354146100cb575b600080fd5b61008d6040518060400160405280601f81526020017f52656769737472794d6f64756c654f776e6572437573746f6d20312e362e300081525081565b60405161009a9190610480565b60405180910390f35b6100b66100b136600461050f565b6100de565b005b6100b66100c636600461050f565b610255565b6100b66100d936600461050f565b6102d0565b60008173ffffffffffffffffffffffffffffffffffffffff1663a217fddf6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561012b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061014f9190610533565b6040517f91d148540000000000000000000000000000000000000000000000000000000081526004810182905233602482015290915073ffffffffffffffffffffffffffffffffffffffff8316906391d1485490604401602060405180830381865afa1580156101c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e7919061054c565b610247576040517f86e0b3440000000000000000000000000000000000000000000000000000000081523360048201526024810182905273ffffffffffffffffffffffffffffffffffffffff831660448201526064015b60405180910390fd5b610251823361031f565b5050565b6102cd818273ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156102a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c8919061056e565b61031f565b50565b6102cd818273ffffffffffffffffffffffffffffffffffffffff16638fd6a6ac6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156102a4573d6000803e3d6000fd5b73ffffffffffffffffffffffffffffffffffffffff8116331461038e576040517fc454d18200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff80831660048301528316602482015260440161023e565b6040517fe677ae3700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff838116600483015282811660248301527f0000000000000000000000000000000000000000000000000000000000000000169063e677ae3790604401600060405180830381600087803b15801561041f57600080fd5b505af1158015610433573d6000803e3d6000fd5b505060405173ffffffffffffffffffffffffffffffffffffffff8085169350851691507f09590fb70af4b833346363965e043a9339e8c7d378b8a2b903c75c277faec4f990600090a35050565b60006020808352835180602085015260005b818110156104ae57858101830151858201604001528201610492565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff811681146102cd57600080fd5b60006020828403121561052157600080fd5b813561052c816104ed565b9392505050565b60006020828403121561054557600080fd5b5051919050565b60006020828403121561055e57600080fd5b8151801515811461052c57600080fd5b60006020828403121561058057600080fd5b815161052c816104ed56fea164736f6c6343000818000a",
}

var RegistryModuleOwnerCustomABI = RegistryModuleOwnerCustomMetaData.ABI

var RegistryModuleOwnerCustomBin = RegistryModuleOwnerCustomMetaData.Bin

func DeployRegistryModuleOwnerCustom(auth *bind.TransactOpts, backend bind.ContractBackend, tokenAdminRegistry common.Address) (common.Address, *types.Transaction, *RegistryModuleOwnerCustom, error) {
	parsed, err := RegistryModuleOwnerCustomMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegistryModuleOwnerCustomBin), backend, tokenAdminRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RegistryModuleOwnerCustom{address: address, abi: *parsed, RegistryModuleOwnerCustomCaller: RegistryModuleOwnerCustomCaller{contract: contract}, RegistryModuleOwnerCustomTransactor: RegistryModuleOwnerCustomTransactor{contract: contract}, RegistryModuleOwnerCustomFilterer: RegistryModuleOwnerCustomFilterer{contract: contract}}, nil
}

type RegistryModuleOwnerCustom struct {
	address common.Address
	abi     abi.ABI
	RegistryModuleOwnerCustomCaller
	RegistryModuleOwnerCustomTransactor
	RegistryModuleOwnerCustomFilterer
}

type RegistryModuleOwnerCustomCaller struct {
	contract *bind.BoundContract
}

type RegistryModuleOwnerCustomTransactor struct {
	contract *bind.BoundContract
}

type RegistryModuleOwnerCustomFilterer struct {
	contract *bind.BoundContract
}

type RegistryModuleOwnerCustomSession struct {
	Contract     *RegistryModuleOwnerCustom
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RegistryModuleOwnerCustomCallerSession struct {
	Contract *RegistryModuleOwnerCustomCaller
	CallOpts bind.CallOpts
}

type RegistryModuleOwnerCustomTransactorSession struct {
	Contract     *RegistryModuleOwnerCustomTransactor
	TransactOpts bind.TransactOpts
}

type RegistryModuleOwnerCustomRaw struct {
	Contract *RegistryModuleOwnerCustom
}

type RegistryModuleOwnerCustomCallerRaw struct {
	Contract *RegistryModuleOwnerCustomCaller
}

type RegistryModuleOwnerCustomTransactorRaw struct {
	Contract *RegistryModuleOwnerCustomTransactor
}

func NewRegistryModuleOwnerCustom(address common.Address, backend bind.ContractBackend) (*RegistryModuleOwnerCustom, error) {
	abi, err := abi.JSON(strings.NewReader(RegistryModuleOwnerCustomABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRegistryModuleOwnerCustom(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustom{address: address, abi: abi, RegistryModuleOwnerCustomCaller: RegistryModuleOwnerCustomCaller{contract: contract}, RegistryModuleOwnerCustomTransactor: RegistryModuleOwnerCustomTransactor{contract: contract}, RegistryModuleOwnerCustomFilterer: RegistryModuleOwnerCustomFilterer{contract: contract}}, nil
}

func NewRegistryModuleOwnerCustomCaller(address common.Address, caller bind.ContractCaller) (*RegistryModuleOwnerCustomCaller, error) {
	contract, err := bindRegistryModuleOwnerCustom(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomCaller{contract: contract}, nil
}

func NewRegistryModuleOwnerCustomTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryModuleOwnerCustomTransactor, error) {
	contract, err := bindRegistryModuleOwnerCustom(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomTransactor{contract: contract}, nil
}

func NewRegistryModuleOwnerCustomFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryModuleOwnerCustomFilterer, error) {
	contract, err := bindRegistryModuleOwnerCustom(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomFilterer{contract: contract}, nil
}

func bindRegistryModuleOwnerCustom(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistryModuleOwnerCustomMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RegistryModuleOwnerCustom.Contract.RegistryModuleOwnerCustomCaller.contract.Call(opts, result, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegistryModuleOwnerCustomTransactor.contract.Transfer(opts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegistryModuleOwnerCustomTransactor.contract.Transact(opts, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RegistryModuleOwnerCustom.Contract.contract.Call(opts, result, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.contract.Transfer(opts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.contract.Transact(opts, method, params...)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RegistryModuleOwnerCustom.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) TypeAndVersion() (string, error) {
	return _RegistryModuleOwnerCustom.Contract.TypeAndVersion(&_RegistryModuleOwnerCustom.CallOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomCallerSession) TypeAndVersion() (string, error) {
	return _RegistryModuleOwnerCustom.Contract.TypeAndVersion(&_RegistryModuleOwnerCustom.CallOpts)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactor) RegisterAccessControlDefaultAdmin(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.contract.Transact(opts, "registerAccessControlDefaultAdmin", token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) RegisterAccessControlDefaultAdmin(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAccessControlDefaultAdmin(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorSession) RegisterAccessControlDefaultAdmin(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAccessControlDefaultAdmin(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactor) RegisterAdminViaGetCCIPAdmin(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.contract.Transact(opts, "registerAdminViaGetCCIPAdmin", token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) RegisterAdminViaGetCCIPAdmin(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaGetCCIPAdmin(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorSession) RegisterAdminViaGetCCIPAdmin(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaGetCCIPAdmin(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactor) RegisterAdminViaOwner(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.contract.Transact(opts, "registerAdminViaOwner", token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomSession) RegisterAdminViaOwner(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaOwner(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomTransactorSession) RegisterAdminViaOwner(token common.Address) (*types.Transaction, error) {
	return _RegistryModuleOwnerCustom.Contract.RegisterAdminViaOwner(&_RegistryModuleOwnerCustom.TransactOpts, token)
}

type RegistryModuleOwnerCustomAdministratorRegisteredIterator struct {
	Event *RegistryModuleOwnerCustomAdministratorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RegistryModuleOwnerCustomAdministratorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryModuleOwnerCustomAdministratorRegistered)
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
		it.Event = new(RegistryModuleOwnerCustomAdministratorRegistered)
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

func (it *RegistryModuleOwnerCustomAdministratorRegisteredIterator) Error() error {
	return it.fail
}

func (it *RegistryModuleOwnerCustomAdministratorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RegistryModuleOwnerCustomAdministratorRegistered struct {
	Token         common.Address
	Administrator common.Address
	Raw           types.Log
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) FilterAdministratorRegistered(opts *bind.FilterOpts, token []common.Address, administrator []common.Address) (*RegistryModuleOwnerCustomAdministratorRegisteredIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var administratorRule []interface{}
	for _, administratorItem := range administrator {
		administratorRule = append(administratorRule, administratorItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.FilterLogs(opts, "AdministratorRegistered", tokenRule, administratorRule)
	if err != nil {
		return nil, err
	}
	return &RegistryModuleOwnerCustomAdministratorRegisteredIterator{contract: _RegistryModuleOwnerCustom.contract, event: "AdministratorRegistered", logs: logs, sub: sub}, nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) WatchAdministratorRegistered(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomAdministratorRegistered, token []common.Address, administrator []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var administratorRule []interface{}
	for _, administratorItem := range administrator {
		administratorRule = append(administratorRule, administratorItem)
	}

	logs, sub, err := _RegistryModuleOwnerCustom.contract.WatchLogs(opts, "AdministratorRegistered", tokenRule, administratorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RegistryModuleOwnerCustomAdministratorRegistered)
				if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "AdministratorRegistered", log); err != nil {
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

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustomFilterer) ParseAdministratorRegistered(log types.Log) (*RegistryModuleOwnerCustomAdministratorRegistered, error) {
	event := new(RegistryModuleOwnerCustomAdministratorRegistered)
	if err := _RegistryModuleOwnerCustom.contract.UnpackLog(event, "AdministratorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustom) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RegistryModuleOwnerCustom.abi.Events["AdministratorRegistered"].ID:
		return _RegistryModuleOwnerCustom.ParseAdministratorRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RegistryModuleOwnerCustomAdministratorRegistered) Topic() common.Hash {
	return common.HexToHash("0x09590fb70af4b833346363965e043a9339e8c7d378b8a2b903c75c277faec4f9")
}

func (_RegistryModuleOwnerCustom *RegistryModuleOwnerCustom) Address() common.Address {
	return _RegistryModuleOwnerCustom.address
}

type RegistryModuleOwnerCustomInterface interface {
	TypeAndVersion(opts *bind.CallOpts) (string, error)

	RegisterAccessControlDefaultAdmin(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	RegisterAdminViaGetCCIPAdmin(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	RegisterAdminViaOwner(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	FilterAdministratorRegistered(opts *bind.FilterOpts, token []common.Address, administrator []common.Address) (*RegistryModuleOwnerCustomAdministratorRegisteredIterator, error)

	WatchAdministratorRegistered(opts *bind.WatchOpts, sink chan<- *RegistryModuleOwnerCustomAdministratorRegistered, token []common.Address, administrator []common.Address) (event.Subscription, error)

	ParseAdministratorRegistered(log types.Log) (*RegistryModuleOwnerCustomAdministratorRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var RegistryModuleOwnerCustomZKBin = ("0x00010000000000020004000000000002000000600310027000000094033001970000000100200190000000260000c13d0000008002000039000000400020043f000000040030008c000001420000413d000000000201043b000000e0022002700000009b0020009c000000530000213d0000009e0020009c000000970000613d0000009f0020009c000001420000c13d000000240030008c000001420000413d0000000002000416000000000002004b000001420000c13d0000000401100370000000000101043b000000970010009c000001420000213d0000009705100197000000a302000041000000800020043f0000000002000414000000040050008c000000c40000c13d0000000003000031000000200030008c00000020040000390000000004034019000000ed0000013d0000000002000416000000000002004b000001420000c13d0000001f023000390000009502200197000000a002200039000000400020043f0000001f0430018f0000009605300198000000a002500039000000370000613d000000a006000039000000000701034f000000007807043c0000000006860436000000000026004b000000330000c13d000000000004004b000000440000613d000000000151034f0000000304400210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000200030008c000001420000413d000000a00100043d000000970010009c000001420000213d0000009701100198000000bb0000c13d000000400100043d00000099020000410000000000210435000000940010009c000000940100804100000040011002100000009a011001c70000024e000104300000009c0020009c000000a70000613d0000009d0020009c000001420000c13d000000240030008c000001420000413d0000000002000416000000000002004b000001420000c13d0000000401100370000000000101043b000000970010009c000001420000213d0000009702100197000000a003000041000000800030043f0000000003000414000000040020008c000000b60000613d000400000001001d000000940030009c0000009403008041000000c001300210000000a1011001c7024c02470000040f000000800a00003900000060031002700000009403300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000080057001bf0000007b0000613d000000000801034f000000008908043c000000000a9a043600000000005a004b000000770000c13d000000000006004b000000880000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000000000003001f0000000100200190000001380000c13d0000001f0530018f0000009606300198000000400200043d0000000004620019000001a10000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000000920000c13d000001a10000013d0000000001000416000000000001004b000001420000c13d000000c001000039000000400010043f0000001f01000039000000800010043f000000a802000041000000a00020043f0000002003000039000000c00030043f000000e00010043f000001000020043f0000011f0000043f000000a9010000410000024d0001042e000000240030008c000001420000413d0000000002000416000000000002004b000001420000c13d0000000401100370000000000101043b000000970010009c000001420000213d0000009702100197000000a203000041000000800030043f0000000003000414000000040020008c000001130000c13d0000000003000031000000200030008c00000020040000390000000004034019000001390000013d000000800010043f00000140000004430000016000100443000000200100003900000100001004430000000101000039000001200010044300000098010000410000024d0001042e000400000001001d000000940020009c0000009402008041000000c001200210000000a1011001c7000300000005001d0000000002050019024c02470000040f00000060031002700000009403300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000080057001bf000000db0000613d0000008008000039000000000901034f000000009a09043c0000000008a80436000000000058004b000000d70000c13d000000000006004b000000e80000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000000000003001f0000000100200190000001470000613d000000040100002900000003050000290000001f02400039000000600420018f000000800b4001bf0000004000b0043f000000200030008c000001420000413d000000800600043d000000a40200004100000000002b043500000084024001bf0000000000620435000000a403400039000000000200041100000000002304350000000003000414000000040050008c0000015f0000c13d00000000044b0019000000400040043f00000000070b0433000000000007004b0000000003000039000000010300c039000000000037004b000001420000c13d000000000007004b000001440000c13d000000a601000041000000000014043500000004014001bf000000000021043500000044024000390000000000520435000000240240003900000000006204350000004001400210000000a7011001c70000024e00010430000400000001001d000000940030009c0000009403008041000000c001300210000000a1011001c7024c02470000040f00000060031002700000009403300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000080057001bf000001280000613d0000008008000039000000000901034f000000009a09043c0000000008a80436000000000058004b000001240000c13d000000000006004b000001350000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000000000003001f0000000100200190000001530000613d00000004010000290000001f02400039000000600220018f00000080022001bf000000400020043f000000200030008c000001420000413d000000800200043d000000970020009c000001440000a13d00000000010000190000024e00010430024c01b40000040f00000000010000190000024d0001042e0000001f0530018f0000009606300198000000400200043d0000000004620019000001a10000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000014e0000c13d000001a10000013d0000001f0530018f0000009606300198000000400200043d0000000004620019000001a10000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000015a0000c13d000001a10000013d000100000006001d000400000001001d000000940030009c0000009403008041000000c0013002100000004002b00210000000000121019f000000a5011001c7000300000005001d000000000205001900020000000b001d024c02470000040f000000020b00002900000060031002700000009403300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b00190000017b0000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b000001770000c13d000000000006004b000001880000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000000000003001f0000000100200190000001960000613d0000001f01400039000000600110018f0000000004b10019000000400040043f000000200030008c0000000401000029000000030500002900000001060000290000000002000411000001000000813d000001420000013d0000001f0530018f0000009606300198000000400200043d0000000004620019000001a10000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000019d0000c13d000000000005004b000001ae0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000000940020009c00000094020080410000004002200210000000000112019f0000024e000104300003000000000002000200000001001d000300970020019b0000000001000411000000030010006b0000020d0000c13d000000ab0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000000940010009c0000009401008041000000c001100210000000ac011001c70000800502000039024c02470000040f00000001002001900000020a0000613d000000000101043b000000ad0200004100000000002004430000009701100197000100000001001d00000004001004430000000001000414000000940010009c0000009401008041000000c001100210000000ae011001c70000800202000039024c02470000040f00000001002001900000020a0000613d000000000101043b000000000001004b0000020b0000613d000000400400043d000000240140003900000003020000290000000000210435000000af010000410000000000140435000000020100002900000097051001970000000401400039000000000051043500000000010004140000000102000029000000040020008c000001fa0000613d000000940040009c000000940300004100000000030440190000004003300210000000940010009c0000009401008041000000c001100210000000000131019f000000a5011001c7000200000005001d000100000004001d024c02420000040f000000010400002900000002050000290000006003100270000000940030019d0000000100200190000002220000613d000000b00040009c0000021c0000813d000000400040043f0000000001000414000000940010009c0000009401008041000000c001100210000000b1011001c70000800d020000390000000303000039000000b2040000410000000306000029024c02420000040f00000001002001900000020b0000613d000000000001042d000000000001042f00000000010000190000024e0001043000000002010000290000009701100197000000400200043d00000024032000390000000000130435000000aa010000410000000000120435000000040120003900000003030000290000000000310435000000940020009c00000094020080410000004001200210000000a5011001c70000024e00010430000000b301000041000000000010043f0000004101000039000000040010043f000000b4010000410000024e0001043000000094033001970000001f0530018f0000009606300198000000400200043d00000000046200190000022e0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000022a0000c13d000000000005004b0000023b0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000000940020009c00000094020080410000004002200210000000000112019f0000024e00010430000000000001042f00000245002104210000000102000039000000000001042d0000000002000019000000000001042d0000024a002104230000000102000039000000000001042d0000000002000019000000000001042d0000024c000004320000024d0001042e0000024e00010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffff00000002000000000000000000000000000000800000010000000000000000009fabe1c10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000096ea2f790000000000000000000000000000000000000000000000000000000096ea2f7a00000000000000000000000000000000000000000000000000000000ff12c35400000000000000000000000000000000000000000000000000000000181f5a770000000000000000000000000000000000000000000000000000000069c0081e8fd6a6ac0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000008000000000000000008da5cb5b00000000000000000000000000000000000000000000000000000000a217fddf0000000000000000000000000000000000000000000000000000000091d1485400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004400000000000000000000000086e0b34400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006400000000000000000000000052656769737472794d6f64756c654f776e6572437573746f6d20312e362e30000000000000000000000000000000000000000060000000c00000000000000000c454d18200000000000000000000000000000000000000000000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e02000002000000000000000000000000000000440000000000000000000000001806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b830200000200000000000000000000000000000024000000000000000000000000e677ae37000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000020000000000000000000000000000000000000000000000000000000000000009590fb70af4b833346363965e043a9339e8c7d378b8a2b903c75c277faec4f94e487b71000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000")

func DeployRegistryModuleOwnerCustomZK(auth *bind.TransactOpts, backend bind.ContractBackend, tokenAdminRegistry common.Address) (common.Address, *generated.
	Transaction, *RegistryModuleOwnerCustom, error) {
	parsed, err :=
		RegistryModuleOwnerCustomMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.
				Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	address, ethTx, contract,
		err := generated.DeployContract(auth, parsed, common.FromHex(RegistryModuleOwnerCustomZKBin),
		backend, tokenAdminRegistry,
	)
	if err !=
		nil {
		return common.
				Address{}, nil,
			nil, err
	}
	return address, ethTx,
		&RegistryModuleOwnerCustom{address: address, abi: *parsed, RegistryModuleOwnerCustomCaller: RegistryModuleOwnerCustomCaller{contract: contract}, RegistryModuleOwnerCustomTransactor: RegistryModuleOwnerCustomTransactor{contract: contract}, RegistryModuleOwnerCustomFilterer: RegistryModuleOwnerCustomFilterer{contract: contract}}, nil
}
