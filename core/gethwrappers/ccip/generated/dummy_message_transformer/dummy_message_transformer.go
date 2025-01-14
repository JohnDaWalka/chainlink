// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dummy_message_transformer

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

type InternalAny2EVMRampMessage struct {
	Header       InternalRampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     common.Address
	GasLimit     *big.Int
	TokenAmounts []InternalAny2EVMTokenTransfer
}

type InternalAny2EVMTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  common.Address
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

type InternalEVM2AnyRampMessage struct {
	Header         InternalRampMessageHeader
	Sender         common.Address
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       common.Address
	FeeTokenAmount *big.Int
	FeeValueJuels  *big.Int
	TokenAmounts   []InternalEVM2AnyTokenTransfer
}

type InternalEVM2AnyTokenTransfer struct {
	SourcePoolAddress common.Address
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

var DummyMessageTransformerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"transformInboundMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transformOutboundMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeValueJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeValueJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"MessageTransformError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]",
	Bin: "0x60808060405234601557610c2a908161001b8239f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c80634546c6e51461057257638a06fadb1461003257600080fd5b3461056d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261056d5760043567ffffffffffffffff811161056d576101a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc823603011261056d576100a8610a11565b906100b63682600401610a8b565b82526100c460a48201610b5a565b6020830190815260c482013567ffffffffffffffff811161056d576100ef9060043691850101610ae9565b916040840192835260e481013567ffffffffffffffff811161056d5761011b9060043691840101610ae9565b926060850193845261010482013567ffffffffffffffff811161056d576101489060043691850101610ae9565b6080860190815261015c6101248401610b5a565b9460a0870195865260c0870192610144850135845260e088019461016481013586526101848101359067ffffffffffffffff821161056d5701953660238801121561056d5760048701356101b76101b282610b7b565b610a32565b9760206004818b858152019360051b830101019036821161056d576024819c9b9a9998979695949c01925b82841061046b575050505061010088019586526101fd610a11565b610205610bf2565b8152602081016000905260408101606090526060810160609052608081016060905260a081016000905260c081016000905260e0810160009052610100016060905260405198899860208a5260208a019051906102999167ffffffffffffffff6080809280518552826020820151166020860152826040820151166040860152826060820151166060860152015116910152565b5173ffffffffffffffffffffffffffffffffffffffff1660c08901525160e088016101a090526101c088016102cd91610b93565b9051908781037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00161010089015261030491610b93565b9051908681037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00161012088015261033b91610b93565b935173ffffffffffffffffffffffffffffffffffffffff16610140860152516101608501525161018084015251908281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0016101a0840152815180825260208201918160051b810160200193602001926000915b8383106103bd5786860387f35b919395509193602080610459837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301875289519073ffffffffffffffffffffffffffffffffffffffff8251168152608061043e61042c8685015160a08886015260a0850190610b93565b60408501518482036040860152610b93565b92606081015160608401520151906080818403910152610b93565b970193019301909286959492936103b0565b839c9495969798999a9b9c3567ffffffffffffffff811161056d5760049083010160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0823603011261056d576104c06109f1565b916104cd60208301610b5a565b8352604082013567ffffffffffffffff811161056d576104f39060203691850101610ae9565b6020840152606082013567ffffffffffffffff811161056d5761051c9060203691850101610ae9565b60408401526080820135606084015260a08201359267ffffffffffffffff841161056d576105536020949385809536920101610ae9565b60808201528152019301929b9a999897969594939b6101e2565b600080fd5b3461056d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261056d5760043567ffffffffffffffff811161056d576101407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc823603011261056d576105e86109a2565b906105f63682600401610a8b565b825260a481013567ffffffffffffffff811161056d5761061c9060043691840101610ae9565b6020830190815260c482013567ffffffffffffffff811161056d576106479060043691850101610ae9565b916040840192835261065b60e48201610b5a565b9260608501938452608085019161010481013583526101248101359067ffffffffffffffff821161056d5701923660238501121561056d5760048401356106a46101b282610b7b565b94602060048188858152019360051b830101019036821161056d5760248101925b8284106108c6576107758a8a8a8a73ffffffffffffffffffffffffffffffffffffffff6107ba8c6107898d60a08901968752606060a06107036109a2565b61070b610bf2565b8152826020820152826040820152600083820152600060808201520152604051998a9960208b5260208b01905167ffffffffffffffff6080809280518552826020820151166020860152826040820151166040860152826060820151166060860152015116910152565b5161014060c08a0152610160890190610b93565b90517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08883030160e0890152610b93565b9351166101008501525161012084015251907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083820301610140840152815180825260208201916020808360051b8301019401926000915b83831061081f5786860387f35b919395509193602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085600195030186528851906080806108ae61086e855160a0865260a0860190610b93565b73ffffffffffffffffffffffffffffffffffffffff87870151168786015263ffffffff604087015116604086015260608601518582036060870152610b93565b93015191015297019301930190928695949293610812565b833567ffffffffffffffff811161056d5760049083010160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0823603011261056d576109116109f1565b91602082013567ffffffffffffffff811161056d576109369060203691850101610ae9565b835261094460408301610b5a565b6020840152606082013563ffffffff8116810361056d57604084015260808201359267ffffffffffffffff841161056d5760a06020949361098b8695863691840101610ae9565b6060840152013560808201528152019301926106c5565b6040519060c0820182811067ffffffffffffffff8211176109c257604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519060a0820182811067ffffffffffffffff8211176109c257604052565b60405190610120820182811067ffffffffffffffff8211176109c257604052565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff8211176109c257604052565b359067ffffffffffffffff8216820361056d57565b91908260a091031261056d57610ae26080610aa46109f1565b9380358552610ab560208201610a76565b6020860152610ac660408201610a76565b6040860152610ad760608201610a76565b606086015201610a76565b6080830152565b81601f8201121561056d5780359067ffffffffffffffff82116109c257610b3760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f85011601610a32565b928284526020838301011161056d57816000926020809301838601378301015290565b359073ffffffffffffffffffffffffffffffffffffffff8216820361056d57565b67ffffffffffffffff81116109c25760051b60200190565b919082519283825260005b848110610bdd5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b80602080928401015182828601015201610b9e565b610bfa6109f1565b90600082526000602083015260006040830152600060608301526000608083015256fea164736f6c634300081a000a",
}

var DummyMessageTransformerABI = DummyMessageTransformerMetaData.ABI

var DummyMessageTransformerBin = DummyMessageTransformerMetaData.Bin

func DeployDummyMessageTransformer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DummyMessageTransformer, error) {
	parsed, err := DummyMessageTransformerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DummyMessageTransformerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DummyMessageTransformer{address: address, abi: *parsed, DummyMessageTransformerCaller: DummyMessageTransformerCaller{contract: contract}, DummyMessageTransformerTransactor: DummyMessageTransformerTransactor{contract: contract}, DummyMessageTransformerFilterer: DummyMessageTransformerFilterer{contract: contract}}, nil
}

type DummyMessageTransformer struct {
	address common.Address
	abi     abi.ABI
	DummyMessageTransformerCaller
	DummyMessageTransformerTransactor
	DummyMessageTransformerFilterer
}

type DummyMessageTransformerCaller struct {
	contract *bind.BoundContract
}

type DummyMessageTransformerTransactor struct {
	contract *bind.BoundContract
}

type DummyMessageTransformerFilterer struct {
	contract *bind.BoundContract
}

type DummyMessageTransformerSession struct {
	Contract     *DummyMessageTransformer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DummyMessageTransformerCallerSession struct {
	Contract *DummyMessageTransformerCaller
	CallOpts bind.CallOpts
}

type DummyMessageTransformerTransactorSession struct {
	Contract     *DummyMessageTransformerTransactor
	TransactOpts bind.TransactOpts
}

type DummyMessageTransformerRaw struct {
	Contract *DummyMessageTransformer
}

type DummyMessageTransformerCallerRaw struct {
	Contract *DummyMessageTransformerCaller
}

type DummyMessageTransformerTransactorRaw struct {
	Contract *DummyMessageTransformerTransactor
}

func NewDummyMessageTransformer(address common.Address, backend bind.ContractBackend) (*DummyMessageTransformer, error) {
	abi, err := abi.JSON(strings.NewReader(DummyMessageTransformerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDummyMessageTransformer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DummyMessageTransformer{address: address, abi: abi, DummyMessageTransformerCaller: DummyMessageTransformerCaller{contract: contract}, DummyMessageTransformerTransactor: DummyMessageTransformerTransactor{contract: contract}, DummyMessageTransformerFilterer: DummyMessageTransformerFilterer{contract: contract}}, nil
}

func NewDummyMessageTransformerCaller(address common.Address, caller bind.ContractCaller) (*DummyMessageTransformerCaller, error) {
	contract, err := bindDummyMessageTransformer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DummyMessageTransformerCaller{contract: contract}, nil
}

func NewDummyMessageTransformerTransactor(address common.Address, transactor bind.ContractTransactor) (*DummyMessageTransformerTransactor, error) {
	contract, err := bindDummyMessageTransformer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DummyMessageTransformerTransactor{contract: contract}, nil
}

func NewDummyMessageTransformerFilterer(address common.Address, filterer bind.ContractFilterer) (*DummyMessageTransformerFilterer, error) {
	contract, err := bindDummyMessageTransformer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DummyMessageTransformerFilterer{contract: contract}, nil
}

func bindDummyMessageTransformer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DummyMessageTransformerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DummyMessageTransformer *DummyMessageTransformerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyMessageTransformer.Contract.DummyMessageTransformerCaller.contract.Call(opts, result, method, params...)
}

func (_DummyMessageTransformer *DummyMessageTransformerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyMessageTransformer.Contract.DummyMessageTransformerTransactor.contract.Transfer(opts)
}

func (_DummyMessageTransformer *DummyMessageTransformerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyMessageTransformer.Contract.DummyMessageTransformerTransactor.contract.Transact(opts, method, params...)
}

func (_DummyMessageTransformer *DummyMessageTransformerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyMessageTransformer.Contract.contract.Call(opts, result, method, params...)
}

func (_DummyMessageTransformer *DummyMessageTransformerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyMessageTransformer.Contract.contract.Transfer(opts)
}

func (_DummyMessageTransformer *DummyMessageTransformerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyMessageTransformer.Contract.contract.Transact(opts, method, params...)
}

func (_DummyMessageTransformer *DummyMessageTransformerCaller) TransformInboundMessage(opts *bind.CallOpts, message InternalAny2EVMRampMessage) (InternalAny2EVMRampMessage, error) {
	var out []interface{}
	err := _DummyMessageTransformer.contract.Call(opts, &out, "transformInboundMessage", message)

	if err != nil {
		return *new(InternalAny2EVMRampMessage), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalAny2EVMRampMessage)).(*InternalAny2EVMRampMessage)

	return out0, err

}

func (_DummyMessageTransformer *DummyMessageTransformerSession) TransformInboundMessage(message InternalAny2EVMRampMessage) (InternalAny2EVMRampMessage, error) {
	return _DummyMessageTransformer.Contract.TransformInboundMessage(&_DummyMessageTransformer.CallOpts, message)
}

func (_DummyMessageTransformer *DummyMessageTransformerCallerSession) TransformInboundMessage(message InternalAny2EVMRampMessage) (InternalAny2EVMRampMessage, error) {
	return _DummyMessageTransformer.Contract.TransformInboundMessage(&_DummyMessageTransformer.CallOpts, message)
}

func (_DummyMessageTransformer *DummyMessageTransformerCaller) TransformOutboundMessage(opts *bind.CallOpts, message InternalEVM2AnyRampMessage) (InternalEVM2AnyRampMessage, error) {
	var out []interface{}
	err := _DummyMessageTransformer.contract.Call(opts, &out, "transformOutboundMessage", message)

	if err != nil {
		return *new(InternalEVM2AnyRampMessage), err
	}

	out0 := *abi.ConvertType(out[0], new(InternalEVM2AnyRampMessage)).(*InternalEVM2AnyRampMessage)

	return out0, err

}

func (_DummyMessageTransformer *DummyMessageTransformerSession) TransformOutboundMessage(message InternalEVM2AnyRampMessage) (InternalEVM2AnyRampMessage, error) {
	return _DummyMessageTransformer.Contract.TransformOutboundMessage(&_DummyMessageTransformer.CallOpts, message)
}

func (_DummyMessageTransformer *DummyMessageTransformerCallerSession) TransformOutboundMessage(message InternalEVM2AnyRampMessage) (InternalEVM2AnyRampMessage, error) {
	return _DummyMessageTransformer.Contract.TransformOutboundMessage(&_DummyMessageTransformer.CallOpts, message)
}

func (_DummyMessageTransformer *DummyMessageTransformer) Address() common.Address {
	return _DummyMessageTransformer.address
}

type DummyMessageTransformerInterface interface {
	TransformInboundMessage(opts *bind.CallOpts, message InternalAny2EVMRampMessage) (InternalAny2EVMRampMessage, error)

	TransformOutboundMessage(opts *bind.CallOpts, message InternalEVM2AnyRampMessage) (InternalEVM2AnyRampMessage, error)

	Address() common.Address
}
