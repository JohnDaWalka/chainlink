package mock_usdc_token_messenger

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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated_zks"
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

var MockE2EUSDCTokenMessengerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burnToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"mintRecipient\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"destinationDomain\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"destinationTokenMessenger\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"destinationCaller\",\"type\":\"bytes32\"}],\"name\":\"DepositForBurn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DESTINATION_TOKEN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"destinationDomain\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"mintRecipient\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"burnToken\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"destinationCaller\",\"type\":\"bytes32\"}],\"name\":\"depositForBurnWithCaller\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"localMessageTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"localMessageTransmitterWithRelay\",\"outputs\":[{\"internalType\":\"contractIMessageTransmitterWithRelay\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBodyVersion\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_nonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b5060405161082d38038061082d83398101604081905261002f91610063565b63ffffffff909116608052600080546001600160401b03191660011790556001600160a01b031660a081905260c0526100b2565b6000806040838503121561007657600080fd5b825163ffffffff8116811461008a57600080fd5b60208401519092506001600160a01b03811681146100a757600080fd5b809150509250929050565b60805160a05160c0516107396100f4600039600081816101290152818161049b015261055b01526000607901526000818160fa01526102b801526107396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063a250c66a11610050578063a250c66a14610124578063f856ddb61461014b578063fb8406a91461015e57600080fd5b80632c121921146100775780637eccf63e146100c35780639cdbb181146100f0575b600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6000546100d79067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100ba565b60405163ffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100ba565b6100997f000000000000000000000000000000000000000000000000000000000000000081565b6100d761015936600461059e565b610193565b6101857f17c71eed51b181d8ae1908b4743526c6dbf099c201f158a1acd5f6718e82e8f681565b6040519081526020016100ba565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526044810186905260009073ffffffffffffffffffffffffffffffffffffffff8416906323b872dd906064016020604051808303816000875af115801561020f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102339190610612565b506040517f42966c680000000000000000000000000000000000000000000000000000000081526004810187905273ffffffffffffffffffffffffffffffffffffffff8416906342966c6890602401600060405180830381600087803b15801561029c57600080fd5b505af11580156102b0573d6000803e3d6000fd5b5050604080517f000000000000000000000000000000000000000000000000000000000000000060e01b7fffffffff0000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff8716602482015260448101889052606481018a9052336084808301919091528251808303909101815260a490910190915291506103779050867f17c71eed51b181d8ae1908b4743526c6dbf099c201f158a1acd5f6718e82e8f68584610457565b600080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff929092169182179055604080518981526020810188905263ffffffff8916918101919091527f17c71eed51b181d8ae1908b4743526c6dbf099c201f158a1acd5f6718e82e8f6606082015260808101859052339173ffffffffffffffffffffffffffffffffffffffff8716917f2fa9ca894982930190727e75500a97d8dc500233a5065e0f3126c48fbe0343c09060a00160405180910390a4505060005467ffffffffffffffff1695945050505050565b60008261051e576040517f0ba469bc00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690630ba469bc906104d49088908890879060040161069f565b6020604051808303816000875af11580156104f3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061051791906106cd565b9050610596565b6040517ff7259a7500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063f7259a75906104d49088908890889088906004016106f7565b949350505050565b600080600080600060a086880312156105b657600080fd5b85359450602086013563ffffffff811681146105d157600080fd5b935060408601359250606086013573ffffffffffffffffffffffffffffffffffffffff8116811461060157600080fd5b949793965091946080013592915050565b60006020828403121561062457600080fd5b8151801515811461063457600080fd5b9392505050565b6000815180845260005b8181101561066157602081850181015186830182015201610645565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b63ffffffff841681528260208201526060604082015260006106c4606083018461063b565b95945050505050565b6000602082840312156106df57600080fd5b815167ffffffffffffffff8116811461063457600080fd5b63ffffffff85168152836020820152826040820152608060608201526000610722608083018461063b565b969550505050505056fea164736f6c6343000818000a",
}

var MockE2EUSDCTokenMessengerABI = MockE2EUSDCTokenMessengerMetaData.ABI

var MockE2EUSDCTokenMessengerBin = MockE2EUSDCTokenMessengerMetaData.Bin

func DeployMockE2EUSDCTokenMessenger(auth *bind.TransactOpts, backend bind.ContractBackend, version uint32, transmitter common.Address) (common.Address, *generated_zks.Transaction, *MockE2EUSDCTokenMessenger, error) {
	parsed, err := MockE2EUSDCTokenMessengerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated_zks.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated_zks.DeployContract(auth, parsed, common.FromHex(MockE2EUSDCTokenMessengerZKBin), backend, version, transmitter)
		contractReturn := &MockE2EUSDCTokenMessenger{address: address, abi: *parsed, MockE2EUSDCTokenMessengerCaller: MockE2EUSDCTokenMessengerCaller{contract: contractBind}, MockE2EUSDCTokenMessengerTransactor: MockE2EUSDCTokenMessengerTransactor{contract: contractBind}, MockE2EUSDCTokenMessengerFilterer: MockE2EUSDCTokenMessengerFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockE2EUSDCTokenMessengerBin), backend, version, transmitter)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated_zks.Transaction{Transaction: tx, Hash_zks: tx.Hash()}, &MockE2EUSDCTokenMessenger{address: address, abi: *parsed, MockE2EUSDCTokenMessengerCaller: MockE2EUSDCTokenMessengerCaller{contract: contract}, MockE2EUSDCTokenMessengerTransactor: MockE2EUSDCTokenMessengerTransactor{contract: contract}, MockE2EUSDCTokenMessengerFilterer: MockE2EUSDCTokenMessengerFilterer{contract: contract}}, nil
}

type MockE2EUSDCTokenMessenger struct {
	address common.Address
	abi     abi.ABI
	MockE2EUSDCTokenMessengerCaller
	MockE2EUSDCTokenMessengerTransactor
	MockE2EUSDCTokenMessengerFilterer
}

type MockE2EUSDCTokenMessengerCaller struct {
	contract *bind.BoundContract
}

type MockE2EUSDCTokenMessengerTransactor struct {
	contract *bind.BoundContract
}

type MockE2EUSDCTokenMessengerFilterer struct {
	contract *bind.BoundContract
}

type MockE2EUSDCTokenMessengerSession struct {
	Contract     *MockE2EUSDCTokenMessenger
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockE2EUSDCTokenMessengerCallerSession struct {
	Contract *MockE2EUSDCTokenMessengerCaller
	CallOpts bind.CallOpts
}

type MockE2EUSDCTokenMessengerTransactorSession struct {
	Contract     *MockE2EUSDCTokenMessengerTransactor
	TransactOpts bind.TransactOpts
}

type MockE2EUSDCTokenMessengerRaw struct {
	Contract *MockE2EUSDCTokenMessenger
}

type MockE2EUSDCTokenMessengerCallerRaw struct {
	Contract *MockE2EUSDCTokenMessengerCaller
}

type MockE2EUSDCTokenMessengerTransactorRaw struct {
	Contract *MockE2EUSDCTokenMessengerTransactor
}

func NewMockE2EUSDCTokenMessenger(address common.Address, backend bind.ContractBackend) (*MockE2EUSDCTokenMessenger, error) {
	abi, err := abi.JSON(strings.NewReader(MockE2EUSDCTokenMessengerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockE2EUSDCTokenMessenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTokenMessenger{address: address, abi: abi, MockE2EUSDCTokenMessengerCaller: MockE2EUSDCTokenMessengerCaller{contract: contract}, MockE2EUSDCTokenMessengerTransactor: MockE2EUSDCTokenMessengerTransactor{contract: contract}, MockE2EUSDCTokenMessengerFilterer: MockE2EUSDCTokenMessengerFilterer{contract: contract}}, nil
}

func NewMockE2EUSDCTokenMessengerCaller(address common.Address, caller bind.ContractCaller) (*MockE2EUSDCTokenMessengerCaller, error) {
	contract, err := bindMockE2EUSDCTokenMessenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTokenMessengerCaller{contract: contract}, nil
}

func NewMockE2EUSDCTokenMessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*MockE2EUSDCTokenMessengerTransactor, error) {
	contract, err := bindMockE2EUSDCTokenMessenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTokenMessengerTransactor{contract: contract}, nil
}

func NewMockE2EUSDCTokenMessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*MockE2EUSDCTokenMessengerFilterer, error) {
	contract, err := bindMockE2EUSDCTokenMessenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTokenMessengerFilterer{contract: contract}, nil
}

func bindMockE2EUSDCTokenMessenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockE2EUSDCTokenMessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockE2EUSDCTokenMessenger.Contract.MockE2EUSDCTokenMessengerCaller.contract.Call(opts, result, method, params...)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockE2EUSDCTokenMessenger.Contract.MockE2EUSDCTokenMessengerTransactor.contract.Transfer(opts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockE2EUSDCTokenMessenger.Contract.MockE2EUSDCTokenMessengerTransactor.contract.Transact(opts, method, params...)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockE2EUSDCTokenMessenger.Contract.contract.Call(opts, result, method, params...)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockE2EUSDCTokenMessenger.Contract.contract.Transfer(opts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockE2EUSDCTokenMessenger.Contract.contract.Transact(opts, method, params...)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCaller) DESTINATIONTOKENMESSENGER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MockE2EUSDCTokenMessenger.contract.Call(opts, &out, "DESTINATION_TOKEN_MESSENGER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerSession) DESTINATIONTOKENMESSENGER() ([32]byte, error) {
	return _MockE2EUSDCTokenMessenger.Contract.DESTINATIONTOKENMESSENGER(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCallerSession) DESTINATIONTOKENMESSENGER() ([32]byte, error) {
	return _MockE2EUSDCTokenMessenger.Contract.DESTINATIONTOKENMESSENGER(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCaller) LocalMessageTransmitter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockE2EUSDCTokenMessenger.contract.Call(opts, &out, "localMessageTransmitter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerSession) LocalMessageTransmitter() (common.Address, error) {
	return _MockE2EUSDCTokenMessenger.Contract.LocalMessageTransmitter(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCallerSession) LocalMessageTransmitter() (common.Address, error) {
	return _MockE2EUSDCTokenMessenger.Contract.LocalMessageTransmitter(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCaller) LocalMessageTransmitterWithRelay(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockE2EUSDCTokenMessenger.contract.Call(opts, &out, "localMessageTransmitterWithRelay")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerSession) LocalMessageTransmitterWithRelay() (common.Address, error) {
	return _MockE2EUSDCTokenMessenger.Contract.LocalMessageTransmitterWithRelay(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCallerSession) LocalMessageTransmitterWithRelay() (common.Address, error) {
	return _MockE2EUSDCTokenMessenger.Contract.LocalMessageTransmitterWithRelay(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCaller) MessageBodyVersion(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _MockE2EUSDCTokenMessenger.contract.Call(opts, &out, "messageBodyVersion")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerSession) MessageBodyVersion() (uint32, error) {
	return _MockE2EUSDCTokenMessenger.Contract.MessageBodyVersion(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCallerSession) MessageBodyVersion() (uint32, error) {
	return _MockE2EUSDCTokenMessenger.Contract.MessageBodyVersion(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCaller) SNonce(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _MockE2EUSDCTokenMessenger.contract.Call(opts, &out, "s_nonce")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerSession) SNonce() (uint64, error) {
	return _MockE2EUSDCTokenMessenger.Contract.SNonce(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerCallerSession) SNonce() (uint64, error) {
	return _MockE2EUSDCTokenMessenger.Contract.SNonce(&_MockE2EUSDCTokenMessenger.CallOpts)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerTransactor) DepositForBurnWithCaller(opts *bind.TransactOpts, amount *big.Int, destinationDomain uint32, mintRecipient [32]byte, burnToken common.Address, destinationCaller [32]byte) (*types.Transaction, error) {
	return _MockE2EUSDCTokenMessenger.contract.Transact(opts, "depositForBurnWithCaller", amount, destinationDomain, mintRecipient, burnToken, destinationCaller)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerSession) DepositForBurnWithCaller(amount *big.Int, destinationDomain uint32, mintRecipient [32]byte, burnToken common.Address, destinationCaller [32]byte) (*types.Transaction, error) {
	return _MockE2EUSDCTokenMessenger.Contract.DepositForBurnWithCaller(&_MockE2EUSDCTokenMessenger.TransactOpts, amount, destinationDomain, mintRecipient, burnToken, destinationCaller)
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerTransactorSession) DepositForBurnWithCaller(amount *big.Int, destinationDomain uint32, mintRecipient [32]byte, burnToken common.Address, destinationCaller [32]byte) (*types.Transaction, error) {
	return _MockE2EUSDCTokenMessenger.Contract.DepositForBurnWithCaller(&_MockE2EUSDCTokenMessenger.TransactOpts, amount, destinationDomain, mintRecipient, burnToken, destinationCaller)
}

type MockE2EUSDCTokenMessengerDepositForBurnIterator struct {
	Event *MockE2EUSDCTokenMessengerDepositForBurn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MockE2EUSDCTokenMessengerDepositForBurnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockE2EUSDCTokenMessengerDepositForBurn)
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
		it.Event = new(MockE2EUSDCTokenMessengerDepositForBurn)
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

func (it *MockE2EUSDCTokenMessengerDepositForBurnIterator) Error() error {
	return it.fail
}

func (it *MockE2EUSDCTokenMessengerDepositForBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MockE2EUSDCTokenMessengerDepositForBurn struct {
	Nonce                     uint64
	BurnToken                 common.Address
	Amount                    *big.Int
	Depositor                 common.Address
	MintRecipient             [32]byte
	DestinationDomain         uint32
	DestinationTokenMessenger [32]byte
	DestinationCaller         [32]byte
	Raw                       types.Log
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerFilterer) FilterDepositForBurn(opts *bind.FilterOpts, nonce []uint64, burnToken []common.Address, depositor []common.Address) (*MockE2EUSDCTokenMessengerDepositForBurnIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var burnTokenRule []interface{}
	for _, burnTokenItem := range burnToken {
		burnTokenRule = append(burnTokenRule, burnTokenItem)
	}

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _MockE2EUSDCTokenMessenger.contract.FilterLogs(opts, "DepositForBurn", nonceRule, burnTokenRule, depositorRule)
	if err != nil {
		return nil, err
	}
	return &MockE2EUSDCTokenMessengerDepositForBurnIterator{contract: _MockE2EUSDCTokenMessenger.contract, event: "DepositForBurn", logs: logs, sub: sub}, nil
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerFilterer) WatchDepositForBurn(opts *bind.WatchOpts, sink chan<- *MockE2EUSDCTokenMessengerDepositForBurn, nonce []uint64, burnToken []common.Address, depositor []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var burnTokenRule []interface{}
	for _, burnTokenItem := range burnToken {
		burnTokenRule = append(burnTokenRule, burnTokenItem)
	}

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _MockE2EUSDCTokenMessenger.contract.WatchLogs(opts, "DepositForBurn", nonceRule, burnTokenRule, depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MockE2EUSDCTokenMessengerDepositForBurn)
				if err := _MockE2EUSDCTokenMessenger.contract.UnpackLog(event, "DepositForBurn", log); err != nil {
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

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessengerFilterer) ParseDepositForBurn(log types.Log) (*MockE2EUSDCTokenMessengerDepositForBurn, error) {
	event := new(MockE2EUSDCTokenMessengerDepositForBurn)
	if err := _MockE2EUSDCTokenMessenger.contract.UnpackLog(event, "DepositForBurn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessenger) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MockE2EUSDCTokenMessenger.abi.Events["DepositForBurn"].ID:
		return _MockE2EUSDCTokenMessenger.ParseDepositForBurn(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MockE2EUSDCTokenMessengerDepositForBurn) Topic() common.Hash {
	return common.HexToHash("0x2fa9ca894982930190727e75500a97d8dc500233a5065e0f3126c48fbe0343c0")
}

func (_MockE2EUSDCTokenMessenger *MockE2EUSDCTokenMessenger) Address() common.Address {
	return _MockE2EUSDCTokenMessenger.address
}

type MockE2EUSDCTokenMessengerInterface interface {
	DESTINATIONTOKENMESSENGER(opts *bind.CallOpts) ([32]byte, error)

	LocalMessageTransmitter(opts *bind.CallOpts) (common.Address, error)

	LocalMessageTransmitterWithRelay(opts *bind.CallOpts) (common.Address, error)

	MessageBodyVersion(opts *bind.CallOpts) (uint32, error)

	SNonce(opts *bind.CallOpts) (uint64, error)

	DepositForBurnWithCaller(opts *bind.TransactOpts, amount *big.Int, destinationDomain uint32, mintRecipient [32]byte, burnToken common.Address, destinationCaller [32]byte) (*types.Transaction, error)

	FilterDepositForBurn(opts *bind.FilterOpts, nonce []uint64, burnToken []common.Address, depositor []common.Address) (*MockE2EUSDCTokenMessengerDepositForBurnIterator, error)

	WatchDepositForBurn(opts *bind.WatchOpts, sink chan<- *MockE2EUSDCTokenMessengerDepositForBurn, nonce []uint64, burnToken []common.Address, depositor []common.Address) (event.Subscription, error)

	ParseDepositForBurn(log types.Log) (*MockE2EUSDCTokenMessengerDepositForBurn, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var MockE2EUSDCTokenMessengerZKBin string = ("0x0001000000000002000f0000000000020000006003100270000000b0033001970000000100200190000000250000c13d0000008002000039000000400020043f000000040030008c000001010000413d000000000201043b000000e002200270000000b60020009c0000005f0000213d000000ba0020009c0000006c0000613d000000bb0020009c000000770000613d000000bc0020009c000001010000c13d0000000001000416000000000001004b000001010000c13d0000000001000412000d00000001001d000c00000000003d0000800501000039000000440300003900000000040004150000000d0440008a0000000504400210000000c60200004102bd029a0000040f000000b001100197000000800010043f000000be01000041000002be0001042e0000000002000416000000000002004b000001010000c13d0000001f02300039000000b102200197000000e002200039000000400020043f0000001f0430018f000000b205300198000000e002500039000000360000613d000000e006000039000000000701034f000000007807043c0000000006860436000000000026004b000000320000c13d000000000004004b000000430000613d000000000151034f0000000304400210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000400030008c000001010000413d000000e00200043d000000b00020009c000001010000213d000001000100043d000000b30010009c000001010000213d000000800020043f000000000300041a000000b40330019700000001033001bf000000000030041b000000a00010043f000000c00010043f0000014000000443000001600020044300000020020000390000018000200443000001a0001004430000004003000039000001c000300443000001e000100443000001000020044300000003010000390000012000100443000000b501000041000002be0001042e000000b70020009c0000007f0000613d000000b80020009c000000900000613d000000b90020009c000001010000c13d0000000001000416000000000001004b000001010000c13d000000bd01000041000000800010043f000000be01000041000002be0001042e0000000001000416000000000001004b000001010000c13d0000000001000412000f00000001001d000e00200000003d0000800501000039000000440300003900000000040004150000000f0440008a000000890000013d0000000001000416000000000001004b000001010000c13d000000000100041a000000ca01100197000000800010043f000000be01000041000002be0001042e0000000001000416000000000001004b000001010000c13d0000000001000412000b00000001001d000a00400000003d0000800501000039000000440300003900000000040004150000000b0440008a0000000504400210000000c60200004102bd029a0000040f000000b301100197000000800010043f000000be01000041000002be0001042e000000a40030008c000001010000413d0000000002000416000000000002004b000001010000c13d0000000402100370000000000502043b0000002402100370000000000702043b000000b00070009c000001010000213d0000004402100370000000000602043b0000006402100370000000000202043b000000b30020009c000001010000213d000000b3022001970000008401100370000000000101043b000700000001001d000000bf01000041000000800010043f0000000001000411000000840010043f0000000001000410000000a40010043f000000c40050043f0000000001000414000000040020008c000000b40000c13d0000000003000031000000200030008c00000020040000390000000004034019000000e00000013d000500000007001d000600000006001d000800000005001d000000b00010009c000000b001008041000000c001100210000000c0011001c7000900000002001d02bd02b30000040f000000800a0000390000006003100270000000b003300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000080057001bf000000cc0000613d000000000801034f000000008908043c000000000a9a043600000000005a004b000000c80000c13d000000000006004b000000d90000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000000000003001f0000000100200190000001030000613d00000009020000290000000805000029000000060600002900000005070000290000001f01400039000000600110018f00000080011001bf000000400010043f000000200030008c000001010000413d000000800100043d000000000001004b0000000003000039000000010300c039000000000031004b000001010000c13d000500000007001d000600000006001d000800000005001d000000c1010000410000000000100443000900000002001d00000004002004430000000001000414000000b00010009c000000b001008041000000c001100210000000c2011001c7000080020200003902bd02b80000040f00000001002001900000025c0000613d000000000101043b000000000001004b00000009020000290000000803000029000001210000c13d0000000001000019000002bf000104300000001f0530018f000000b206300198000000400200043d00000000046200190000010e0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000010a0000c13d000000000005004b0000011b0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000000b00020009c000000b0020080410000004002200210000000000112019f000002bf00010430000000400400043d000000c3010000410000000001140436000300000001001d000000040140003900000000003104350000000001000414000000040020008c0000013a0000613d000000b00040009c000000b00300004100000000030440190000004003300210000000b00010009c000000b001008041000000c001100210000000000131019f000000c4011001c7000400000004001d02bd02b30000040f00000004040000290000006003100270000000b00030019d0000000100200190000001dd0000613d000000c50040009c000001420000413d000000cf01000041000000000010043f0000004101000039000000040010043f000000c401000041000002bf00010430000400000004001d000000400040043f000000c60100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000000b00010009c000000b001008041000000c001100210000000c7011001c7000080050200003902bd02b80000040f00000001002001900000025c0000613d000000000101043b0000000407000029000000840270003900000000030004110000000000320435000000640270003900000008030000290000000000320435000000440270003900000006030000290000000000320435000000e0011002100000000302000029000000000012043500000024017000390000000902000029000000000021043500000084010000390000000000170435000000c80070009c0000013c0000213d0000000501000029000100b00010019b000001440170003900000124027000390000010403700039000000e404700039000000c405700039000000c006700039000500000006001d000000400060043f0000000006070433000200000006001d000000070000006b000001ea0000c13d000000cb060000410000000507000029000000000067043500000001060000290000000000650435000000bd0500004100000000005404350000006004000039000000000043043500000002060000290000000000620435000000000006004b00000003050000290000018b0000613d000000000200001900000000031200190000000004520019000000000404043300000000004304350000002002200039000000000062004b000001840000413d00000000011600190000000000010435000000c601000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000000b00010009c000000b001008041000000c001100210000000c7011001c7000080050200003902bd02b80000040f00000001002001900000025c0000613d000000000201043b0000000001000414000000b302200197000000040020008c0000021a0000613d00000002030000290000001f03300039000000d0033001970000000504000029000000b00040009c000000b00400804100000040044002100000008403300039000000b00030009c000000b0030080410000006003300210000000000343019f000000b00010009c000000b001008041000000c001100210000000000131019f02bd02b30000040f0000006003100270000000b003300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000505700029000001c10000613d000000000801034f0000000509000029000000008a08043c0000000009a90436000000000059004b000001bd0000c13d000000000006004b000001ce0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000000000003001f00000001002001900000021e0000c13d0000001f0530018f000000b206300198000000400200043d00000000046200190000010e0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000001d80000c13d0000010e0000013d000000b0033001970000001f0530018f000000b206300198000000400200043d00000000046200190000010e0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000001e50000c13d0000010e0000013d000000c9060000410000000507000029000000000067043500000001060000290000000000650435000000bd05000041000000000054043500000007040000290000000000430435000000800300003900000000003204350000000206000029000000000061043500000004010000290000016401100039000000000006004b0000000305000029000002040000613d000000000200001900000000031200190000000004520019000000000404043300000000004304350000002002200039000000000062004b000001fd0000413d00000000011600190000000000010435000000c601000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000000b00010009c000000b001008041000000c001100210000000c7011001c7000080050200003902bd02b80000040f00000001002001900000025c0000613d000000000201043b0000000001000414000000b302200197000000040020008c0000025d0000c13d0000000003000031000000200030008c000000200400003900000000040340190000001f01400039000000600210018f0000000501200029000000000021004b00000000020000390000000102004039000000ca0010009c0000013c0000213d00000001002001900000013c0000c13d000000400010043f000000200030008c000000090600002900000008040000290000000607000029000001010000413d00000005020000290000000002020433000500000002001d000000ca0020009c000001010000213d000000000200041a000000b4022001970000000505000029000000000252019f000000000020041b0000008002100039000000070300002900000000003204350000006002100039000000bd030000410000000000320435000000400210003900000001030000290000000000320435000000200210003900000000007204350000000000410435000000b00010009c000000b00100804100000040011002100000000002000414000000b00020009c000000b002008041000000c002200210000000000112019f000000cc011001c70000800d020000390000000403000039000000cd04000041000000000700041102bd02b30000040f0000000100200190000001010000613d000000400100043d00000005020000290000000000210435000000b00010009c000000b0010080410000004001100210000000ce011001c7000002be0001042e000000000001042f00000002030000290000001f03300039000000d0033001970000000504000029000000b00040009c000000b0040080410000004004400210000000a403300039000000b00030009c000000b0030080410000006003300210000000000343019f000000b00010009c000000b001008041000000c001100210000000000131019f02bd02b30000040f0000006003100270000000b003300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000005057000290000027d0000613d000000000801034f0000000509000029000000008a08043c0000000009a90436000000000059004b000002790000c13d000000000006004b0000028a0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000000000003001f00000001002001900000021e0000c13d0000001f0530018f000000b206300198000000400200043d00000000046200190000010e0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000002940000c13d0000010e0000013d000000000001042f00000000050100190000000000200443000000040100003900000005024002700000000002020031000000000121043a0000002004400039000000000031004b0000029d0000413d000000b00030009c000000b00300804100000060013002100000000002000414000000b00020009c000000b002008041000000c002200210000000000112019f000000d1011001c7000000000205001902bd02b80000040f0000000100200190000002b20000613d000000000101043b000000000001042d000000000001042f000002b6002104210000000102000039000000000001042d0000000002000019000000000001042d000002bb002104230000000102000039000000000001042d0000000002000019000000000001042d000002bd00000432000002be0001042e000002bf0001043000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000020000000000000000000000000000010000000100000000000000000000000000000000000000000000000000000000000000000000000000a250c66900000000000000000000000000000000000000000000000000000000a250c66a00000000000000000000000000000000000000000000000000000000f856ddb600000000000000000000000000000000000000000000000000000000fb8406a9000000000000000000000000000000000000000000000000000000002c121921000000000000000000000000000000000000000000000000000000007eccf63e000000000000000000000000000000000000000000000000000000009cdbb18117c71eed51b181d8ae1908b4743526c6dbf099c201f158a1acd5f6718e82e8f6000000000000000000000000000000000000002000000080000000000000000023b872dd0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000640000008000000000000000001806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b83020000020000000000000000000000000000002400000000000000000000000042966c680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000000000000010000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e0200000200000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff3ff7259a7500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0ba469bc0000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000a00000000000000000000000002fa9ca894982930190727e75500a97d8dc500233a5065e0f3126c48fbe0343c000000000000000000000000000000000000000200000000000000000000000004e487b7100000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe002000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
