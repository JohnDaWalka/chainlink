package ping_pong_demo

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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

var PingPongDemoMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"feeToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isOutOfOrder\",\"type\":\"bool\"}],\"name\":\"OutOfOrderExecutionChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pingPongCount\",\"type\":\"uint256\"}],\"name\":\"Ping\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pingPongCount\",\"type\":\"uint256\"}],\"name\":\"Pong\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCounterpartAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCounterpartChainSelector\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOutOfOrderExecution\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"counterpartChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"counterpartAddress\",\"type\":\"address\"}],\"name\":\"setCounterpart\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setCounterpartAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"setCounterpartChainSelector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"outOfOrderExecution\",\"type\":\"bool\"}],\"name\":\"setOutOfOrderExecution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"pause\",\"type\":\"bool\"}],\"name\":\"setPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startPingPong\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620014a8380380620014a88339810160408190526200003491620001ff565b336000836001600160a01b03811662000067576040516335fdcccd60e21b81526000600482015260240160405180910390fd5b6001600160a01b0390811660805282166200009557604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b0384811691909117909155811615620000c857620000c8816200016c565b50506002805460ff60a01b19169055600380546001600160a01b0319166001600160a01b0383811691821790925560405163095ea7b360e01b8152918416600483015260001960248301529063095ea7b3906044016020604051808303816000875af11580156200013d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016391906200023e565b50505062000269565b336001600160a01b038216036200019657604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620001fc57600080fd5b50565b600080604083850312156200021357600080fd5b82516200022081620001e6565b60208401519092506200023381620001e6565b809150509250929050565b6000602082840312156200025157600080fd5b815180151581146200026257600080fd5b9392505050565b6080516112156200029360003960008181610295015281816106520152610a5901526112156000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c80638da5cb5b116100b2578063b187bd2611610081578063bee518a411610066578063bee518a4146102ef578063ca709a251461032d578063f2fde38b1461034b57600080fd5b8063b187bd26146102b9578063b5a11011146102dc57600080fd5b80638da5cb5b1461023f5780639d2aede51461025d578063ae90de5514610270578063b0f479a11461029357600080fd5b80632874d8bf11610109578063665ed537116100ee578063665ed5371461021157806379ba50971461022457806385572ffb1461022c57600080fd5b80632874d8bf146101ca5780632b6e5d63146101d257600080fd5b806301ffc9a71461013b57806316c38b3c14610163578063181f5a77146101785780631892b906146101b7575b600080fd5b61014e610149366004610c29565b61035e565b60405190151581526020015b60405180910390f35b610176610171366004610c72565b6103f7565b005b604080518082018252601281527f50696e67506f6e6744656d6f20312e352e3000000000000000000000000000006020820152905161015a9190610cf8565b6101766101c5366004610d28565b610449565b6101766104a4565b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161015a565b61017661021f366004610c72565b6104e0565b61017661056c565b61017661023a366004610d43565b61063a565b60015473ffffffffffffffffffffffffffffffffffffffff166101ec565b61017661026b366004610da2565b6106c3565b60035474010000000000000000000000000000000000000000900460ff1661014e565b7f00000000000000000000000000000000000000000000000000000000000000006101ec565b60025474010000000000000000000000000000000000000000900460ff1661014e565b6101766102ea366004610dbd565b610712565b60015474010000000000000000000000000000000000000000900467ffffffffffffffff1660405167ffffffffffffffff909116815260200161015a565b60035473ffffffffffffffffffffffffffffffffffffffff166101ec565b610176610359366004610da2565b6107b4565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f85572ffb0000000000000000000000000000000000000000000000000000000014806103f157507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b6103ff6107c5565b6002805491151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff909216919091179055565b6104516107c5565b6001805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff909216919091179055565b6104ac6107c5565b600280547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556104de6001610816565b565b6104e86107c5565b6003805482151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff9091161790556040517f05a3fef9935c9013a24c6193df2240d34fcf6b0ebf8786b85efe8401d696cdd99061056190831515815260200190565b60405180910390a150565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105bd576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146106af576040517fd7f7333400000000000000000000000000000000000000000000000000000000815233600482015260240160405180910390fd5b6106c06106bb82610ff3565b610b0f565b50565b6106cb6107c5565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b61071a6107c5565b6001805467ffffffffffffffff90931674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff909316929092179091556002805473ffffffffffffffffffffffffffffffffffffffff9092167fffffffffffffffffffffffff0000000000000000000000000000000000000000909216919091179055565b6107bc6107c5565b6106c081610b65565b60015473ffffffffffffffffffffffffffffffffffffffff1633146104de576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600116600103610859576040518181527f48257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f9060200160405180910390a161088d565b6040518181527f58b69f57828e6962d216502094c54f6562f3bf082ba758966c3454f9e37b15259060200160405180910390a15b6040805160a0810190915260025473ffffffffffffffffffffffffffffffffffffffff1660c08201526000908060e081016040516020818303038152906040528152602001836040516020016108e591815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190528152602001600060405190808252806020026020018201604052801561095f57816020015b60408051808201909152600080825260208201528152602001906001900390816109385790505b50815260035473ffffffffffffffffffffffffffffffffffffffff811660208084019190915260408051808201825262030d408082527401000000000000000000000000000000000000000090940460ff16151590830190815281516024810194909452511515604480850191909152815180850390910181526064909301815290820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f181dcf1000000000000000000000000000000000000000000000000000000000179052909101526001546040517f96f4e9f90000000000000000000000000000000000000000000000000000000081529192507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16916396f4e9f991610ac7917401000000000000000000000000000000000000000090910467ffffffffffffffff169085906004016110a0565b6020604051808303816000875af1158015610ae6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b0a91906111b5565b505050565b60008160600151806020019051810190610b2991906111b5565b60025490915074010000000000000000000000000000000000000000900460ff16610b6157610b61610b5c8260016111ce565b610816565b5050565b3373ffffffffffffffffffffffffffffffffffffffff821603610bb4576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215610c3b57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610c6b57600080fd5b9392505050565b600060208284031215610c8457600080fd5b81358015158114610c6b57600080fd5b6000815180845260005b81811015610cba57602081850181015186830182015201610c9e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610c6b6020830184610c94565b803567ffffffffffffffff81168114610d2357600080fd5b919050565b600060208284031215610d3a57600080fd5b610c6b82610d0b565b600060208284031215610d5557600080fd5b813567ffffffffffffffff811115610d6c57600080fd5b820160a08185031215610c6b57600080fd5b803573ffffffffffffffffffffffffffffffffffffffff81168114610d2357600080fd5b600060208284031215610db457600080fd5b610c6b82610d7e565b60008060408385031215610dd057600080fd5b610dd983610d0b565b9150610de760208401610d7e565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715610e4257610e42610df0565b60405290565b60405160a0810167ffffffffffffffff81118282101715610e4257610e42610df0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610eb257610eb2610df0565b604052919050565b600082601f830112610ecb57600080fd5b813567ffffffffffffffff811115610ee557610ee5610df0565b610f1660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610e6b565b818152846020838601011115610f2b57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f830112610f5957600080fd5b8135602067ffffffffffffffff821115610f7557610f75610df0565b610f83818360051b01610e6b565b82815260069290921b84018101918181019086841115610fa257600080fd5b8286015b84811015610fe85760408189031215610fbf5760008081fd5b610fc7610e1f565b610fd082610d7e565b81528185013585820152835291830191604001610fa6565b509695505050505050565b600060a0823603121561100557600080fd5b61100d610e48565b8235815261101d60208401610d0b565b6020820152604083013567ffffffffffffffff8082111561103d57600080fd5b61104936838701610eba565b6040840152606085013591508082111561106257600080fd5b61106e36838701610eba565b6060840152608085013591508082111561108757600080fd5b5061109436828601610f48565b60808301525092915050565b6000604067ffffffffffffffff851683526020604081850152845160a060408601526110cf60e0860182610c94565b9050818601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08087840301606088015261110a8383610c94565b6040890151888203830160808a01528051808352908601945060009350908501905b8084101561116b578451805173ffffffffffffffffffffffffffffffffffffffff1683528601518683015293850193600193909301929086019061112c565b50606089015173ffffffffffffffffffffffffffffffffffffffff1660a08901526080890151888203830160c08a015295506111a78187610c94565b9a9950505050505050505050565b6000602082840312156111c757600080fd5b5051919050565b808201808211156103f1577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea164736f6c6343000818000a",
}

var PingPongDemoABI = PingPongDemoMetaData.ABI

var PingPongDemoBin = PingPongDemoMetaData.Bin

func DeployPingPongDemo(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, feeToken common.Address) (common.Address, *generated.Transaction, *PingPongDemo, error) {
	parsed, err := PingPongDemoMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(PingPongDemoZKBin), backend, router, feeToken)
		contractReturn := &PingPongDemo{address: address, abi: *parsed, PingPongDemoCaller: PingPongDemoCaller{contract: contractBind}, PingPongDemoTransactor: PingPongDemoTransactor{contract: contractBind}, PingPongDemoFilterer: PingPongDemoFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PingPongDemoBin), backend, router, feeToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, Hash_zks: tx.Hash()}, &PingPongDemo{address: address, abi: *parsed, PingPongDemoCaller: PingPongDemoCaller{contract: contract}, PingPongDemoTransactor: PingPongDemoTransactor{contract: contract}, PingPongDemoFilterer: PingPongDemoFilterer{contract: contract}}, nil
}

type PingPongDemo struct {
	address common.Address
	abi     abi.ABI
	PingPongDemoCaller
	PingPongDemoTransactor
	PingPongDemoFilterer
}

type PingPongDemoCaller struct {
	contract *bind.BoundContract
}

type PingPongDemoTransactor struct {
	contract *bind.BoundContract
}

type PingPongDemoFilterer struct {
	contract *bind.BoundContract
}

type PingPongDemoSession struct {
	Contract     *PingPongDemo
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PingPongDemoCallerSession struct {
	Contract *PingPongDemoCaller
	CallOpts bind.CallOpts
}

type PingPongDemoTransactorSession struct {
	Contract     *PingPongDemoTransactor
	TransactOpts bind.TransactOpts
}

type PingPongDemoRaw struct {
	Contract *PingPongDemo
}

type PingPongDemoCallerRaw struct {
	Contract *PingPongDemoCaller
}

type PingPongDemoTransactorRaw struct {
	Contract *PingPongDemoTransactor
}

func NewPingPongDemo(address common.Address, backend bind.ContractBackend) (*PingPongDemo, error) {
	abi, err := abi.JSON(strings.NewReader(PingPongDemoABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPingPongDemo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PingPongDemo{address: address, abi: abi, PingPongDemoCaller: PingPongDemoCaller{contract: contract}, PingPongDemoTransactor: PingPongDemoTransactor{contract: contract}, PingPongDemoFilterer: PingPongDemoFilterer{contract: contract}}, nil
}

func NewPingPongDemoCaller(address common.Address, caller bind.ContractCaller) (*PingPongDemoCaller, error) {
	contract, err := bindPingPongDemo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoCaller{contract: contract}, nil
}

func NewPingPongDemoTransactor(address common.Address, transactor bind.ContractTransactor) (*PingPongDemoTransactor, error) {
	contract, err := bindPingPongDemo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoTransactor{contract: contract}, nil
}

func NewPingPongDemoFilterer(address common.Address, filterer bind.ContractFilterer) (*PingPongDemoFilterer, error) {
	contract, err := bindPingPongDemo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoFilterer{contract: contract}, nil
}

func bindPingPongDemo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PingPongDemoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_PingPongDemo *PingPongDemoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PingPongDemo.Contract.PingPongDemoCaller.contract.Call(opts, result, method, params...)
}

func (_PingPongDemo *PingPongDemoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.Contract.PingPongDemoTransactor.contract.Transfer(opts)
}

func (_PingPongDemo *PingPongDemoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PingPongDemo.Contract.PingPongDemoTransactor.contract.Transact(opts, method, params...)
}

func (_PingPongDemo *PingPongDemoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PingPongDemo.Contract.contract.Call(opts, result, method, params...)
}

func (_PingPongDemo *PingPongDemoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.Contract.contract.Transfer(opts)
}

func (_PingPongDemo *PingPongDemoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PingPongDemo.Contract.contract.Transact(opts, method, params...)
}

func (_PingPongDemo *PingPongDemoCaller) GetCounterpartAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getCounterpartAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetCounterpartAddress() (common.Address, error) {
	return _PingPongDemo.Contract.GetCounterpartAddress(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetCounterpartAddress() (common.Address, error) {
	return _PingPongDemo.Contract.GetCounterpartAddress(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetCounterpartChainSelector(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getCounterpartChainSelector")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetCounterpartChainSelector() (uint64, error) {
	return _PingPongDemo.Contract.GetCounterpartChainSelector(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetCounterpartChainSelector() (uint64, error) {
	return _PingPongDemo.Contract.GetCounterpartChainSelector(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetFeeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getFeeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetFeeToken() (common.Address, error) {
	return _PingPongDemo.Contract.GetFeeToken(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetFeeToken() (common.Address, error) {
	return _PingPongDemo.Contract.GetFeeToken(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetOutOfOrderExecution(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getOutOfOrderExecution")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetOutOfOrderExecution() (bool, error) {
	return _PingPongDemo.Contract.GetOutOfOrderExecution(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetOutOfOrderExecution() (bool, error) {
	return _PingPongDemo.Contract.GetOutOfOrderExecution(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetRouter() (common.Address, error) {
	return _PingPongDemo.Contract.GetRouter(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetRouter() (common.Address, error) {
	return _PingPongDemo.Contract.GetRouter(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) IsPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "isPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) IsPaused() (bool, error) {
	return _PingPongDemo.Contract.IsPaused(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) IsPaused() (bool, error) {
	return _PingPongDemo.Contract.IsPaused(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) Owner() (common.Address, error) {
	return _PingPongDemo.Contract.Owner(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) Owner() (common.Address, error) {
	return _PingPongDemo.Contract.Owner(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PingPongDemo.Contract.SupportsInterface(&_PingPongDemo.CallOpts, interfaceId)
}

func (_PingPongDemo *PingPongDemoCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PingPongDemo.Contract.SupportsInterface(&_PingPongDemo.CallOpts, interfaceId)
}

func (_PingPongDemo *PingPongDemoCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) TypeAndVersion() (string, error) {
	return _PingPongDemo.Contract.TypeAndVersion(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) TypeAndVersion() (string, error) {
	return _PingPongDemo.Contract.TypeAndVersion(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "acceptOwnership")
}

func (_PingPongDemo *PingPongDemoSession) AcceptOwnership() (*types.Transaction, error) {
	return _PingPongDemo.Contract.AcceptOwnership(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PingPongDemo.Contract.AcceptOwnership(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "ccipReceive", message)
}

func (_PingPongDemo *PingPongDemoSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.Contract.CcipReceive(&_PingPongDemo.TransactOpts, message)
}

func (_PingPongDemo *PingPongDemoTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.Contract.CcipReceive(&_PingPongDemo.TransactOpts, message)
}

func (_PingPongDemo *PingPongDemoTransactor) SetCounterpart(opts *bind.TransactOpts, counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "setCounterpart", counterpartChainSelector, counterpartAddress)
}

func (_PingPongDemo *PingPongDemoSession) SetCounterpart(counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpart(&_PingPongDemo.TransactOpts, counterpartChainSelector, counterpartAddress)
}

func (_PingPongDemo *PingPongDemoTransactorSession) SetCounterpart(counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpart(&_PingPongDemo.TransactOpts, counterpartChainSelector, counterpartAddress)
}

func (_PingPongDemo *PingPongDemoTransactor) SetCounterpartAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "setCounterpartAddress", addr)
}

func (_PingPongDemo *PingPongDemoSession) SetCounterpartAddress(addr common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpartAddress(&_PingPongDemo.TransactOpts, addr)
}

func (_PingPongDemo *PingPongDemoTransactorSession) SetCounterpartAddress(addr common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpartAddress(&_PingPongDemo.TransactOpts, addr)
}

func (_PingPongDemo *PingPongDemoTransactor) SetCounterpartChainSelector(opts *bind.TransactOpts, chainSelector uint64) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "setCounterpartChainSelector", chainSelector)
}

func (_PingPongDemo *PingPongDemoSession) SetCounterpartChainSelector(chainSelector uint64) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpartChainSelector(&_PingPongDemo.TransactOpts, chainSelector)
}

func (_PingPongDemo *PingPongDemoTransactorSession) SetCounterpartChainSelector(chainSelector uint64) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpartChainSelector(&_PingPongDemo.TransactOpts, chainSelector)
}

func (_PingPongDemo *PingPongDemoTransactor) SetOutOfOrderExecution(opts *bind.TransactOpts, outOfOrderExecution bool) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "setOutOfOrderExecution", outOfOrderExecution)
}

func (_PingPongDemo *PingPongDemoSession) SetOutOfOrderExecution(outOfOrderExecution bool) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetOutOfOrderExecution(&_PingPongDemo.TransactOpts, outOfOrderExecution)
}

func (_PingPongDemo *PingPongDemoTransactorSession) SetOutOfOrderExecution(outOfOrderExecution bool) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetOutOfOrderExecution(&_PingPongDemo.TransactOpts, outOfOrderExecution)
}

func (_PingPongDemo *PingPongDemoTransactor) SetPaused(opts *bind.TransactOpts, pause bool) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "setPaused", pause)
}

func (_PingPongDemo *PingPongDemoSession) SetPaused(pause bool) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetPaused(&_PingPongDemo.TransactOpts, pause)
}

func (_PingPongDemo *PingPongDemoTransactorSession) SetPaused(pause bool) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetPaused(&_PingPongDemo.TransactOpts, pause)
}

func (_PingPongDemo *PingPongDemoTransactor) StartPingPong(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "startPingPong")
}

func (_PingPongDemo *PingPongDemoSession) StartPingPong() (*types.Transaction, error) {
	return _PingPongDemo.Contract.StartPingPong(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactorSession) StartPingPong() (*types.Transaction, error) {
	return _PingPongDemo.Contract.StartPingPong(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "transferOwnership", to)
}

func (_PingPongDemo *PingPongDemoSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.TransferOwnership(&_PingPongDemo.TransactOpts, to)
}

func (_PingPongDemo *PingPongDemoTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.TransferOwnership(&_PingPongDemo.TransactOpts, to)
}

type PingPongDemoOutOfOrderExecutionChangeIterator struct {
	Event *PingPongDemoOutOfOrderExecutionChange

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoOutOfOrderExecutionChangeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoOutOfOrderExecutionChange)
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
		it.Event = new(PingPongDemoOutOfOrderExecutionChange)
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

func (it *PingPongDemoOutOfOrderExecutionChangeIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoOutOfOrderExecutionChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoOutOfOrderExecutionChange struct {
	IsOutOfOrder bool
	Raw          types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterOutOfOrderExecutionChange(opts *bind.FilterOpts) (*PingPongDemoOutOfOrderExecutionChangeIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "OutOfOrderExecutionChange")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoOutOfOrderExecutionChangeIterator{contract: _PingPongDemo.contract, event: "OutOfOrderExecutionChange", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchOutOfOrderExecutionChange(opts *bind.WatchOpts, sink chan<- *PingPongDemoOutOfOrderExecutionChange) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "OutOfOrderExecutionChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoOutOfOrderExecutionChange)
				if err := _PingPongDemo.contract.UnpackLog(event, "OutOfOrderExecutionChange", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseOutOfOrderExecutionChange(log types.Log) (*PingPongDemoOutOfOrderExecutionChange, error) {
	event := new(PingPongDemoOutOfOrderExecutionChange)
	if err := _PingPongDemo.contract.UnpackLog(event, "OutOfOrderExecutionChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoOwnershipTransferRequestedIterator struct {
	Event *PingPongDemoOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoOwnershipTransferRequested)
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
		it.Event = new(PingPongDemoOwnershipTransferRequested)
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

func (it *PingPongDemoOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoOwnershipTransferRequestedIterator{contract: _PingPongDemo.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoOwnershipTransferRequested)
				if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseOwnershipTransferRequested(log types.Log) (*PingPongDemoOwnershipTransferRequested, error) {
	event := new(PingPongDemoOwnershipTransferRequested)
	if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoOwnershipTransferredIterator struct {
	Event *PingPongDemoOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoOwnershipTransferred)
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
		it.Event = new(PingPongDemoOwnershipTransferred)
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

func (it *PingPongDemoOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoOwnershipTransferredIterator{contract: _PingPongDemo.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoOwnershipTransferred)
				if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseOwnershipTransferred(log types.Log) (*PingPongDemoOwnershipTransferred, error) {
	event := new(PingPongDemoOwnershipTransferred)
	if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoPingIterator struct {
	Event *PingPongDemoPing

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoPingIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoPing)
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
		it.Event = new(PingPongDemoPing)
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

func (it *PingPongDemoPingIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoPingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoPing struct {
	PingPongCount *big.Int
	Raw           types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterPing(opts *bind.FilterOpts) (*PingPongDemoPingIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "Ping")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoPingIterator{contract: _PingPongDemo.contract, event: "Ping", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchPing(opts *bind.WatchOpts, sink chan<- *PingPongDemoPing) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "Ping")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoPing)
				if err := _PingPongDemo.contract.UnpackLog(event, "Ping", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParsePing(log types.Log) (*PingPongDemoPing, error) {
	event := new(PingPongDemoPing)
	if err := _PingPongDemo.contract.UnpackLog(event, "Ping", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoPongIterator struct {
	Event *PingPongDemoPong

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoPongIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoPong)
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
		it.Event = new(PingPongDemoPong)
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

func (it *PingPongDemoPongIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoPongIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoPong struct {
	PingPongCount *big.Int
	Raw           types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterPong(opts *bind.FilterOpts) (*PingPongDemoPongIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "Pong")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoPongIterator{contract: _PingPongDemo.contract, event: "Pong", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchPong(opts *bind.WatchOpts, sink chan<- *PingPongDemoPong) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "Pong")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoPong)
				if err := _PingPongDemo.contract.UnpackLog(event, "Pong", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParsePong(log types.Log) (*PingPongDemoPong, error) {
	event := new(PingPongDemoPong)
	if err := _PingPongDemo.contract.UnpackLog(event, "Pong", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_PingPongDemo *PingPongDemo) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _PingPongDemo.abi.Events["OutOfOrderExecutionChange"].ID:
		return _PingPongDemo.ParseOutOfOrderExecutionChange(log)
	case _PingPongDemo.abi.Events["OwnershipTransferRequested"].ID:
		return _PingPongDemo.ParseOwnershipTransferRequested(log)
	case _PingPongDemo.abi.Events["OwnershipTransferred"].ID:
		return _PingPongDemo.ParseOwnershipTransferred(log)
	case _PingPongDemo.abi.Events["Ping"].ID:
		return _PingPongDemo.ParsePing(log)
	case _PingPongDemo.abi.Events["Pong"].ID:
		return _PingPongDemo.ParsePong(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (PingPongDemoOutOfOrderExecutionChange) Topic() common.Hash {
	return common.HexToHash("0x05a3fef9935c9013a24c6193df2240d34fcf6b0ebf8786b85efe8401d696cdd9")
}

func (PingPongDemoOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (PingPongDemoOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (PingPongDemoPing) Topic() common.Hash {
	return common.HexToHash("0x48257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f")
}

func (PingPongDemoPong) Topic() common.Hash {
	return common.HexToHash("0x58b69f57828e6962d216502094c54f6562f3bf082ba758966c3454f9e37b1525")
}

func (_PingPongDemo *PingPongDemo) Address() common.Address {
	return _PingPongDemo.address
}

type PingPongDemoInterface interface {
	GetCounterpartAddress(opts *bind.CallOpts) (common.Address, error)

	GetCounterpartChainSelector(opts *bind.CallOpts) (uint64, error)

	GetFeeToken(opts *bind.CallOpts) (common.Address, error)

	GetOutOfOrderExecution(opts *bind.CallOpts) (bool, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	IsPaused(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	SetCounterpart(opts *bind.TransactOpts, counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error)

	SetCounterpartAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error)

	SetCounterpartChainSelector(opts *bind.TransactOpts, chainSelector uint64) (*types.Transaction, error)

	SetOutOfOrderExecution(opts *bind.TransactOpts, outOfOrderExecution bool) (*types.Transaction, error)

	SetPaused(opts *bind.TransactOpts, pause bool) (*types.Transaction, error)

	StartPingPong(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOutOfOrderExecutionChange(opts *bind.FilterOpts) (*PingPongDemoOutOfOrderExecutionChangeIterator, error)

	WatchOutOfOrderExecutionChange(opts *bind.WatchOpts, sink chan<- *PingPongDemoOutOfOrderExecutionChange) (event.Subscription, error)

	ParseOutOfOrderExecutionChange(log types.Log) (*PingPongDemoOutOfOrderExecutionChange, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*PingPongDemoOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*PingPongDemoOwnershipTransferred, error)

	FilterPing(opts *bind.FilterOpts) (*PingPongDemoPingIterator, error)

	WatchPing(opts *bind.WatchOpts, sink chan<- *PingPongDemoPing) (event.Subscription, error)

	ParsePing(log types.Log) (*PingPongDemoPing, error)

	FilterPong(opts *bind.FilterOpts) (*PingPongDemoPongIterator, error)

	WatchPong(opts *bind.WatchOpts, sink chan<- *PingPongDemoPong) (event.Subscription, error)

	ParsePong(log types.Log) (*PingPongDemoPong, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var PingPongDemoZKBin string = ("0x00030000000000020005000000000002000200000001035500000060031002700000017b0030019d0000017b0330019700000001002001900000012f0000c13d0000008002000039000000400020043f000000040030008c000004640000413d000000000201043b000000e002200270000001890020009c000001620000213d000001960020009c000001810000a13d000001970020009c000001da0000a13d000001980020009c000002640000613d000001990020009c000002420000613d0000019a0020009c000004640000c13d000000240030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000401100370000000000101043b000300000001001d000001830010009c000004640000213d000000030130006a000001aa0010009c000004640000213d000000a40010008c000004640000413d000001ab01000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000017b0010009c0000017b01008041000000c001100210000001ac011001c7000080050200003905e805e30000040f00000001002001900000040a0000613d000000400200043d000000000101043b0000017e011001970000000003000411000000000013004b0000040b0000c13d000001ad0020009c000004d30000213d000000a001200039000000400010043f000000030500002900000004035000390000000201000367000000000331034f000000000303043b00000000043204360000002403500039000000000531034f000000000505043b000001830050009c000004640000213d00000000005404350000002005300039000000000351034f000000000303043b000001830030009c000004640000213d00000003093000290000002306900039000001ae076001970000000003000031000001ae04300197000000000847013f000000000047004b0000000007000019000001ae07004041000000000036004b0000000006000019000001ae06008041000001ae0080009c000000000706c019000000000007004b000004640000c13d000000040a9000390000000006a1034f000000000606043b000001830060009c000004d30000213d0000001f07600039000001c5077001970000003f07700039000001c508700197000000400700043d0000000008870019000000000078004b000000000b000039000000010b004039000001830080009c000004d30000213d0000000100b00190000004d30000c13d000000400080043f000000000867043600000000096900190000002409900039000000000039004b000004640000213d0000002009a00039000000000a91034f000001c50b6001980000001f0c60018f0000000009b80019000000860000613d000000000d0a034f000000000e08001900000000df0d043c000000000efe043600000000009e004b000000820000c13d00000000000c004b000000930000613d000000000aba034f000000030bc00210000000000c090433000000000cbc01cf000000000cbc022f000000000a0a043b000001000bb00089000000000aba022f000000000aba01cf000000000aca019f0000000000a9043500000000066800190000000000060435000000400620003900000000007604350000002006500039000000000561034f000000000505043b000001830050009c000004640000213d00000003095000290000002305900039000000000035004b0000000007000019000001ae07008041000001ae05500197000000000845013f000000000045004b0000000005000019000001ae05004041000001ae0080009c000000000507c019000000000005004b000004640000c13d000000040a9000390000000005a1034f000000000505043b000001830050009c000004d30000213d0000001f07500039000001c5077001970000003f07700039000001c508700197000000400700043d0000000008870019000000000078004b000000000b000039000000010b004039000001830080009c000004d30000213d0000000100b00190000004d30000c13d000000400080043f000000000857043600000000095900190000002409900039000000000039004b000004640000213d0000002009a00039000000000a91034f000001c50b5001980000001f0c50018f0000000009b80019000000ce0000613d000000000d0a034f000000000e08001900000000df0d043c000000000efe043600000000009e004b000000ca0000c13d00000000000c004b000000db0000613d000000000aba034f000000030bc00210000000000c090433000000000cbc01cf000000000cbc022f000000000a0a043b000001000bb00089000000000aba022f000000000aba01cf000000000aca019f0000000000a9043500000000055800190000000000050435000000600520003900000000007504350000002006600039000000000661034f000000000606043b000001830060009c000004640000213d00000003066000290000002307600039000000000037004b0000000008000019000001ae08008041000001ae07700197000000000947013f000000000047004b0000000004000019000001ae04004041000001ae0090009c000000000408c019000000000004004b000004640000c13d0000000404600039000000000441034f000000000804043b000001830080009c000004d30000213d00000005048002100000003f04400039000001af07400197000000400400043d0000000007740019000000000047004b00000000090000390000000109004039000001830070009c000004d30000213d0000000100900190000004d30000c13d000000400070043f0000000000840435000000240660003900000006078002100000000007670019000000000037004b000004640000213d000000000008004b000004720000c13d0000008001200039000000000041043500000000010504330000000012010434000001aa0020009c000004640000213d000000200020008c000004640000413d0000000202000039000000000202041a000300000002001d000001a900200198000004230000c13d0000000001010433000000010010003a000005a30000413d0000000102100039000000400100043d00000000002104350000017b0010009c0000017b010080410000004001100210000200000002001d00000001002001900000048c0000c13d00000000020004140000017b0020009c0000017b02008041000000c002200210000000000112019f000001b1011001c70000800d020000390000000103000039000001b304000041000004950000013d000000a004000039000000400040043f0000000002000416000000000002004b000004640000c13d0000001f023000390000017c02200197000000a002200039000000400020043f0000001f0530018f0000017d06300198000000a002600039000001410000613d000000000701034f000000007807043c0000000004840436000000000024004b0000013d0000c13d000000000005004b0000014e0000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000400030008c000004640000413d000000a00300043d0000017e0030009c000004640000213d000000c00200043d0000017e0020009c000004640000213d0000017e01300198000002070000c13d000000400100043d00000187020000410000000000210435000000040210003900000000000204350000017b0010009c0000017b01008041000000400110021000000188011001c7000005ea000104300000018a0020009c0000019c0000a13d0000018b0020009c000001e30000a13d0000018c0020009c000002880000613d0000018d0020009c0000025b0000613d0000018e0020009c000004640000c13d000000240030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000401100370000000000601043b0000017e0060009c000004640000213d0000000101000039000000000101041a0000017e011001970000000005000411000000000015004b0000039d0000c13d000000000056004b000004140000c13d000001a401000041000000800010043f000001a501000041000005ea000104300000019d0020009c000001b20000213d000001a00020009c000002130000613d000001a10020009c000004640000c13d000000240030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000401100370000000000201043b000000000002004b0000000001000039000000010100c039000300000002001d000000000012004b000004640000c13d05e805bb0000040f000000030000006b0000000001000019000001bc0100c0410000000202000039000000000302041a0000018003300197000001c70000013d000001910020009c000001cb0000213d000001940020009c000002240000613d000001950020009c000004640000c13d000000240030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000401100370000000000101043b000300000001001d0000017e0010009c000004640000213d05e805bb0000040f0000000201000039000000000201041a0000017f022001970000000303000029000002030000013d0000019e0020009c000002290000613d0000019f0020009c000004640000c13d000000240030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000401100370000000000101043b000300000001001d000001830010009c000004640000213d05e805bb0000040f0000000301000029000000a001100210000001a7011001970000000102000039000000000302041a000001a803300197000000000113019f000000000012041b0000000001000019000005e90001042e000001920020009c0000023d0000613d000001930020009c000004640000c13d0000000001000416000000000001004b000004640000c13d0000000001000412000500000001001d000400000000003d0000000001000415000000050110008a000000050110021005e805cb0000040f000002600000013d0000019b0020009c000002920000613d0000019c0020009c000004640000c13d0000000001000416000000000001004b000004640000c13d00000002010000390000025f0000013d0000018f0020009c000003740000613d000001900020009c000004640000c13d000000440030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000402100370000000000202043b000300000002001d000001830020009c000004640000213d0000002401100370000000000101043b000200000001001d0000017e0010009c000004640000213d05e805bb0000040f0000000301000029000000a001100210000001a7011001970000000102000039000000000302041a000001a803300197000000000113019f000000000012041b0000000201000039000000000201041a0000017f022001970000000203000029000000000232019f000000000021041b0000000001000019000005e90001042e000000800030043f000000400b00043d0000000003000411000000000003004b0000037f0000c13d000001850100004100000000001b04350000017b00b0009c0000017b0b0080410000004001b0021000000186011001c7000005ea00010430000000240030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000401100370000000000101043b000001c200100198000004640000c13d000001c30010009c00000000020000390000000102006039000001c40010009c00000001022061bf000000800020043f000001a601000041000005e90001042e0000000001000416000000000001004b000004640000c13d00000001010000390000025f0000013d0000000001000416000000000001004b000004640000c13d000000c001000039000000400010043f0000001201000039000000800010043f000001c001000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e00200003905e805a90000040f000000c00110008a0000017b0010009c0000017b010080410000006001100210000001c1011001c7000005e90001042e0000000001000416000000000001004b000004640000c13d0000000301000039000003780000013d0000000001000416000000000001004b000004640000c13d000000000100041a0000017e021001970000000006000411000000000026004b000003e80000c13d0000000102000039000000000302041a0000017f04300197000000000464019f000000000042041b0000017f01100197000000000010041b00000000010004140000017e053001970000017b0010009c0000017b01008041000000c001100210000001a2011001c70000800d020000390000000303000039000001bb04000041000004200000013d0000000001000416000000000001004b000004640000c13d0000000301000039000000000101041a0000017e01100197000000800010043f000001a601000041000005e90001042e000000240030008c000004640000413d0000000002000416000000000002004b000004640000c13d0000000401100370000000000101043b000000000001004b0000000002000039000000010200c039000000000021004b000004640000c13d0000000102000039000000000202041a0000017e022001970000000003000411000000000023004b0000039d0000c13d0000000302000039000000000302041a0000018003300197000000000001004b0000000004000019000001bc0400c041000000000343019f000000000032041b000000800010043f00000000010004140000017b0010009c0000017b01008041000000c001100210000001bd011001c70000800d020000390000000103000039000001be04000041000004200000013d0000000001000416000000000001004b000004640000c13d0000000101000039000000000101041a000000a0011002700000018301100197000000800010043f000001a601000041000005e90001042e0000000001000416000000000001004b000004640000c13d0000000103000039000000000103041a0000017e011001970000000002000411000000000012004b0000039d0000c13d0000000201000039000000000201041a000300000002001d0000018002200197000000000021041b000000800030043f00000000010004140000017b0010009c0000017b01008041000000c001100210000001bd011001c70000800d02000039000001b20400004105e805de0000040f0000000100200190000004640000613d00000003010000290000017e02100197000000400100043d0000002003100039000000000023043500000020020000390000000000210435000001b00010009c000004d30000213d0000004004100039000000400040043f0000006003100039000000010500003900000000005304350000000000240435000001b40010009c000004d30000213d0000008005100039000000400050043f000001ad0010009c000004d30000213d000000a002100039000000400020043f0000000000050435000000400300043d000001b00030009c000004d30000213d0000000302000039000000000202041a0000004006300039000000400060043f000001a9002001980000000006000039000000010600c03900000020073000390000000000670435000001b5060000410000000000630435000000400600043d0000002008600039000001b60900004100000000009804350000000003030433000000240860003900000000003804350000000003070433000000000003004b0000000003000039000000010300c0390000004407600039000000000037043500000044030000390000000000360435000001b40060009c000004d30000213d0000008007600039000000400070043f000001b70060009c000004d30000213d0000012003600039000000400030043f00000000001704350000017e01200197000000e0026000390000000000120435000000c003600039000000000053043500000100016000390000000000610435000000a00560003900000000004504350000000104000039000000000404041a000000400900043d000000240690003900000040080000390000000000860435000001b8060000410000000000690435000000a00440027000000183044001970000000406900039000000000046043500000000040704330000004406900039000000a0070000390000000000760435000000e40890003900000000760404340000000000680435000300000009001d0000010404900039000000000006004b000003110000613d00000000080000190000000009480019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b0000030a0000413d000000000764001900000000000704350000001f06600039000001c506600197000000000464001900000003070000290000000006740049000000440660008a00000000050504330000006407700039000000000067043500000000650504340000000004540436000000000005004b000003280000613d000000000700001900000000084700190000000009760019000000000909043300000000009804350000002007700039000000000057004b000003210000413d000000000654001900000000000604350000001f05500039000001c505500197000000000654001900000003070000290000000004760049000000440540008a00000000040304330000008403700039000000000053043500000000050404330000000003560436000000000005004b000003430000613d00000000060000190000002004400039000000000704043300000000870704340000017e0770019700000000077304360000000008080433000000000087043500000040033000390000000106600039000000000056004b000003380000413d00000000020204330000017e022001970000000305000029000000a40450003900000000002404350000000002530049000000440220008a0000000001010433000000c404500039000000000024043500000000160104340000000005630436000000000006004b000003590000613d000000000200001900000000035200190000000004210019000000000404043300000000004304350000002002200039000000000062004b000003520000413d000200000005001d000100000006001d00000000016500190000000000010435000001ab01000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000017b0010009c0000017b01008041000000c001100210000001ac011001c7000080050200003905e805e30000040f00000001002001900000040a0000613d000000000201043b00000000010004140000017e02200197000000040020008c000004250000c13d0000000104000031000000200040008c0000002004008039000004570000013d0000000001000416000000000001004b000004640000c13d0000000201000039000000000101041a000001a9001001980000000001000039000000010100c039000000800010043f000001a601000041000005e90001042e0000017e022001970000000105000039000000000405041a0000017f04400197000000000334019f000000000035041b0000000203000039000000000403041a0000018004400197000000000043041b0000000303000039000000000403041a0000017f04400197000000000424019f000000000043041b0000002403b00039000000010400008a0000000000430435000001810300004100000000003b04350000000403b0003900000000001304350000000001000414000000040020008c000003a10000c13d0000000103000031000000200030008c00000020040000390000000004034019000003cc0000013d000001bf01000041000000800010043f000001a501000041000005ea000104300000017b00b0009c0000017b0300004100000000030b401900000040033002100000017b0010009c0000017b01008041000000c001100210000000000131019f00000182011001c700030000000b001d05e805de0000040f000000030b00002900000060031002700000017b03300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b0019000003bc0000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b000003b80000c13d000000000006004b000003c90000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000003ec0000613d0000001f01400039000000600210018f0000000001b20019000000000021004b00000000020000390000000102004039000001830010009c000004d30000213d0000000100200190000004d30000c13d000000400010043f000000200030008c000004640000413d00000000010b0433000000000001004b0000000002000039000000010200c039000000000021004b000004640000c13d000000800100043d0000014000000443000001600010044300000020010000390000010000100443000000010100003900000120001004430000018401000041000005e90001042e000001ba01000041000000800010043f000001a501000041000005ea000104300000001f0530018f0000017d06300198000000400200043d0000000004620019000003f70000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000003f30000c13d000000000005004b000004040000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f000000000014043500000060013002100000017b0020009c0000017b020080410000004002200210000000000112019f000005ea00010430000000000001042f00000187010000410000000000120435000000040120003900000000003104350000017b0020009c0000017b02008041000000400120021000000188011001c7000005ea00010430000000000100041a0000017f01100197000000000161019f000000000010041b00000000010004140000017b0010009c0000017b01008041000000c001100210000001a2011001c70000800d020000390000000303000039000001a30400004105e805de0000040f0000000100200190000004640000613d0000000001000019000005e90001042e00000001030000290000001f03300039000001c5033001970000000305000029000000000353004900000002033000290000017b0050009c0000017b04000041000000000405401900000040044002100000017b0030009c0000017b030080410000006003300210000000000343019f0000017b0010009c0000017b01008041000000c001100210000000000131019f05e805de0000040f00000060031002700000017b03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000305700029000004470000613d000000000801034f0000000309000029000000008a08043c0000000009a90436000000000059004b000004430000c13d000000000006004b000004540000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000004660000613d0000001f01400039000000600210018f0000000301200029000000000021004b00000000020000390000000102004039000001830010009c000004d30000213d0000000100200190000004d30000c13d000000400010043f000000200040008c000004230000813d0000000001000019000005ea000104300000001f0530018f0000017d06300198000000400200043d0000000004620019000003f70000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000046d0000c13d000003f70000013d00000000080400190000000009630049000001aa0090009c000004640000213d000000400090008c000004640000413d000000400900043d000001b00090009c000004d30000213d000000400a9000390000004000a0043f000000000a61034f000000000a0a043b0000017e00a0009c000004640000213d0000002008800039000000000aa90436000000200b600039000000000bb1034f000000000b0b043b0000000000ba043500000000009804350000004006600039000000000076004b000004730000413d0000010c0000013d00000000020004140000017b0020009c0000017b02008041000000c002200210000000000112019f000001b1011001c70000800d020000390000000103000039000001b20400004105e805de0000040f0000000100200190000004640000613d00000003010000290000017e02100197000000400100043d0000002003100039000000000023043500000020020000390000000000210435000001b00010009c000004d30000213d0000004004100039000000400040043f0000006003100039000000020500002900000000005304350000000000240435000001b40010009c000004d30000213d0000008005100039000000400050043f000001ad0010009c000004d30000213d000000a002100039000000400020043f0000000000050435000000400300043d000001b00030009c000004d30000213d0000000302000039000000000202041a0000004006300039000000400060043f000001a9002001980000000006000039000000010600c03900000020073000390000000000670435000001b5060000410000000000630435000000400600043d0000002008600039000001b60900004100000000009804350000000003030433000000240860003900000000003804350000000003070433000000000003004b0000000003000039000000010300c0390000004407600039000000000037043500000044030000390000000000360435000001b40060009c000004d30000213d0000008007600039000000400070043f000001b70060009c000004d90000a13d000001b901000041000000000010043f0000004101000039000000040010043f0000018801000041000005ea000104300000012003600039000000400030043f00000000001704350000017e01200197000000e0026000390000000000120435000000c003600039000000000053043500000100016000390000000000610435000000a00560003900000000004504350000000104000039000000000404041a000000400900043d000000240690003900000040080000390000000000860435000001b8060000410000000000690435000000a00440027000000183044001970000000406900039000000000046043500000000040704330000004406900039000000a0070000390000000000760435000000e40690003900000000740404340000000000460435000300000009001d0000010406900039000000000004004b000005040000613d00000000080000190000000009680019000000000a870019000000000a0a04330000000000a904350000002008800039000000000048004b000004fd0000413d000000000746001900000000000704350000001f04400039000001c504400197000000000446001900000003070000290000000006740049000000440660008a00000000050504330000006407700039000000000067043500000000650504340000000004540436000000000005004b0000051b0000613d000000000700001900000000084700190000000009760019000000000909043300000000009804350000002007700039000000000057004b000005140000413d000000000654001900000000000604350000001f05500039000001c505500197000000000654001900000003070000290000000004760049000000440540008a00000000040304330000008403700039000000000053043500000000050404330000000003560436000000000005004b000005360000613d00000000060000190000002004400039000000000704043300000000870704340000017e0770019700000000077304360000000008080433000000000087043500000040033000390000000106600039000000000056004b0000052b0000413d00000000020204330000017e022001970000000305000029000000a40450003900000000002404350000000002530049000000440220008a0000000001010433000000c404500039000000000024043500000000420104340000000001230436000000000002004b0000054c0000613d000000000300001900000000051300190000000006340019000000000606043300000000006504350000002003300039000000000023004b000005450000413d0000000003210019000000000003043500000000030004140000000004000411000000040040008c000005570000c13d0000000103000031000000200030008c00000020040000390000000004034019000005890000013d0000001f02200039000001c5022001970000000304000029000000000242004900000000011200190000017b0010009c0000017b0100804100000060011002100000017b0040009c0000017b0200004100000000020440190000004002200210000000000121019f0000017b0030009c0000017b03008041000000c002300210000000000112019f000000000200041105e805de0000040f00000060031002700000017b03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000305700029000005790000613d000000000801034f0000000309000029000000008a08043c0000000009a90436000000000059004b000005750000c13d000000000006004b000005860000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000005970000613d0000001f01400039000000600210018f0000000301200029000000000021004b00000000020000390000000102004039000001830010009c000004d30000213d0000000100200190000004d30000c13d000000400010043f000000200030008c000004230000813d000004640000013d0000001f0530018f0000017d06300198000000400200043d0000000004620019000003f70000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000059e0000c13d000003f70000013d000001b901000041000000000010043f0000001101000039000000040010043f0000018801000041000005ea0001043000000000430104340000000001320436000000000003004b000005b50000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b000005ae0000413d000000000231001900000000000204350000001f02300039000001c5022001970000000001210019000000000001042d0000000101000039000000000101041a0000017e011001970000000002000411000000000012004b000005c20000c13d000000000001042d000000400100043d000001bf0200004100000000002104350000017b0010009c0000017b01008041000000400110021000000186011001c7000005ea00010430000000000001042f000001ab0200004100000000002004430000000501100270000000000201003100000004002004430000000101010031000000240010044300000000010004140000017b0010009c0000017b01008041000000c001100210000001ac011001c7000080050200003905e805e30000040f0000000100200190000005dd0000613d000000000101043b000000000001042d000000000001042f000005e1002104210000000102000039000000000001042d0000000002000019000000000001042d000005e6002104230000000102000039000000000001042d0000000002000019000000000001042d000005e800000432000005e90001042e000005ea00010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000ffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff095ea7b3000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff00000002000000000000000000000000000000800000010000000000000000009b15e16f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000d7f73334000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000000000000000000000008da5cb5a00000000000000000000000000000000000000000000000000000000b187bd2500000000000000000000000000000000000000000000000000000000bee518a300000000000000000000000000000000000000000000000000000000bee518a400000000000000000000000000000000000000000000000000000000ca709a2500000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000b187bd2600000000000000000000000000000000000000000000000000000000b5a1101100000000000000000000000000000000000000000000000000000000ae90de5400000000000000000000000000000000000000000000000000000000ae90de5500000000000000000000000000000000000000000000000000000000b0f479a1000000000000000000000000000000000000000000000000000000008da5cb5b000000000000000000000000000000000000000000000000000000009d2aede5000000000000000000000000000000000000000000000000000000002874d8be00000000000000000000000000000000000000000000000000000000665ed53600000000000000000000000000000000000000000000000000000000665ed5370000000000000000000000000000000000000000000000000000000079ba50970000000000000000000000000000000000000000000000000000000085572ffb000000000000000000000000000000000000000000000000000000002874d8bf000000000000000000000000000000000000000000000000000000002b6e5d6300000000000000000000000000000000000000000000000000000000181f5a7600000000000000000000000000000000000000000000000000000000181f5a77000000000000000000000000000000000000000000000000000000001892b9060000000000000000000000000000000000000000000000000000000001ffc9a70000000000000000000000000000000000000000000000000000000016c38b3c0200000000000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000000000000000000000000000000000000000002000000080000000000000000000000000ffffffffffffffff0000000000000000000000000000000000000000ffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff0000000000000000000000ff00000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e0200000200000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5f80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffffbf020000000000000000000000000000000000002000000000000000000000000048257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f58b69f57828e6962d216502094c54f6562f3bf082ba758966c3454f9e37b1525000000000000000000000000000000000000000000000000ffffffffffffff7f0000000000000000000000000000000000000000000000000000000000030d40181dcf1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffedf96f4e9f9000000000000000000000000000000000000000000000000000000004e487b710000000000000000000000000000000000000000000000000000000002b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e00000000000000000000000010000000000000000000000000000000000000000020000000000000000000000000000000000002000000080000000000000000005a3fef9935c9013a24c6193df2240d34fcf6b0ebf8786b85efe8401d696cdd92b5c74de0000000000000000000000000000000000000000000000000000000050696e67506f6e6744656d6f20312e352e3000000000000000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff01ffc9a70000000000000000000000000000000000000000000000000000000085572ffb00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00000000000000000000000000000000000000000000000000000000000000000")
