package message_hasher

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

type ClientEVMExtraArgsV1 struct {
	GasLimit *big.Int
}

type ClientEVMExtraArgsV2 struct {
	GasLimit                 *big.Int
	AllowOutOfOrderExecution bool
}

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

var MessageHasherMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"decodeEVMExtraArgsV1\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMExtraArgsV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"name\":\"decodeEVMExtraArgsV2\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"internalType\":\"structClient.EVMExtraArgsV2\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"name\":\"encodeAny2EVMTokenAmountsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmount\",\"type\":\"tuple[]\"}],\"name\":\"encodeEVM2AnyTokenAmountsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMExtraArgsV1\",\"name\":\"extraArgs\",\"type\":\"tuple\"}],\"name\":\"encodeEVMExtraArgsV1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"internalType\":\"structClient.EVMExtraArgsV2\",\"name\":\"extraArgs\",\"type\":\"tuple\"}],\"name\":\"encodeEVMExtraArgsV2\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafDomainSeparator\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"metaDataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fixedSizeFieldsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"senderHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tokenAmountsHash\",\"type\":\"bytes32\"}],\"name\":\"encodeFinalHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"name\":\"encodeFixedSizeFieldsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"any2EVMMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"onRampHash\",\"type\":\"bytes32\"}],\"name\":\"encodeMetadataHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611053806100206000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063bf0619ad11610076578063c7ca9a181161005b578063c7ca9a1814610273578063e04767b814610286578063e733d2091461029957600080fd5b8063bf0619ad146101c9578063c63641bd1461021c57600080fd5b806394b6624b116100a757806394b6624b14610168578063ae5663d71461017b578063b17df7141461018e57600080fd5b80633ec7c377146100c35780638503839d14610147575b600080fd5b6101316100d136600461065b565b60408051602081019690965273ffffffffffffffffffffffffffffffffffffffff949094168585015267ffffffffffffffff928316606086015260808501919091521660a0808401919091528151808403909101815260c0909201905290565b60405161013e9190610716565b60405180910390f35b61015a610155366004610a5b565b6102ac565b60405190815260200161013e565b610131610176366004610b65565b610343565b610131610189366004610cd5565b61036c565b6101ba61019c366004610d12565b60408051602080820183526000909152815190810190915290815290565b6040519051815260200161013e565b6101316101d7366004610d2b565b604080516020810197909752868101959095526060860193909352608085019190915260a084015260c0808401919091528151808403909101815260e0909201905290565b61025661022a366004610d7e565b604080518082019091526000808252602082015250604080518082019091529182521515602082015290565b60408051825181526020928301511515928101929092520161013e565b610131610281366004610daa565b61037f565b610131610294366004610dfe565b610390565b6101316102a7366004610e42565b6103e3565b600061033c837f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f8560000151602001518660000151604001518680519060200120604051602001610321949392919093845267ffffffffffffffff928316602085015291166040830152606082015260800190565b604051602081830303815290604052805190602001206103ee565b9392505050565b6060816040516020016103569190610e84565b6040516020818303038152906040529050919050565b6060816040516020016103569190610f71565b606061038a82610521565b92915050565b6060848484846040516020016103ca949392919093845267ffffffffffffffff928316602085015291166040830152606082015260800190565b6040516020818303038152906040529050949350505050565b606061038a826105e3565b815180516060808501519083015160808087015194015160405160009586958895610460959194909391929160200194855273ffffffffffffffffffffffffffffffffffffffff93909316602085015267ffffffffffffffff9182166040850152606084015216608082015260a00190565b604051602081830303815290604052805190602001208560200151805190602001208660400151805190602001208760a001516040516020016104a39190610f71565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e00160405160208183030381529060405280519060200120905092915050565b604051815160248201526020820151151560448201526060907f181dcf1000000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b604051815160248201526060907f97a657c90000000000000000000000000000000000000000000000000000000090604401610560565b803573ffffffffffffffffffffffffffffffffffffffff8116811461063e57600080fd5b919050565b803567ffffffffffffffff8116811461063e57600080fd5b600080600080600060a0868803121561067357600080fd5b853594506106836020870161061a565b935061069160408701610643565b9250606086013591506106a660808701610643565b90509295509295909350565b6000815180845260005b818110156106d8576020818501810151868301820152016106bc565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061033c60208301846106b2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff8111828210171561077b5761077b610729565b60405290565b60405160c0810167ffffffffffffffff8111828210171561077b5761077b610729565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156107eb576107eb610729565b604052919050565b600060a0828403121561080557600080fd5b61080d610758565b90508135815261081f60208301610643565b602082015261083060408301610643565b604082015261084160608301610643565b606082015261085260808301610643565b608082015292915050565b600082601f83011261086e57600080fd5b813567ffffffffffffffff81111561088857610888610729565b6108b960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016107a4565b8181528460208386010111156108ce57600080fd5b816020850160208301376000918101602001919091529392505050565b600067ffffffffffffffff82111561090557610905610729565b5060051b60200190565b600082601f83011261092057600080fd5b81356020610935610930836108eb565b6107a4565b82815260059290921b8401810191818101908684111561095457600080fd5b8286015b84811015610a5057803567ffffffffffffffff808211156109795760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156109b25760008081fd5b6109ba610758565b87840135838111156109cc5760008081fd5b6109da8d8a8388010161085d565b82525060406109ea81860161061a565b8983015260608086013563ffffffff81168114610a075760008081fd5b808385015250608091508186013585811115610a235760008081fd5b610a318f8c838a010161085d565b9184019190915250919093013590830152508352918301918301610958565b509695505050505050565b60008060408385031215610a6e57600080fd5b823567ffffffffffffffff80821115610a8657600080fd5b908401906101408287031215610a9b57600080fd5b610aa3610781565b610aad87846107f3565b815260a083013582811115610ac157600080fd5b610acd8882860161085d565b60208301525060c083013582811115610ae557600080fd5b610af18882860161085d565b604083015250610b0360e0840161061a565b6060820152610100830135608082015261012083013582811115610b2657600080fd5b610b328882860161090f565b60a08301525093506020850135915080821115610b4e57600080fd5b50610b5b8582860161085d565b9150509250929050565b60006020808385031215610b7857600080fd5b823567ffffffffffffffff80821115610b9057600080fd5b818501915085601f830112610ba457600080fd5b8135610bb2610930826108eb565b81815260059190911b83018401908481019088831115610bd157600080fd5b8585015b83811015610cc857803585811115610bec57600080fd5b860160a0818c037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0011215610c215760008081fd5b610c29610758565b610c3489830161061a565b815260408083013588811115610c4a5760008081fd5b610c588e8c8387010161085d565b8b8401525060608084013589811115610c715760008081fd5b610c7f8f8d8388010161085d565b83850152506080915081840135818401525060a083013588811115610ca45760008081fd5b610cb28e8c8387010161085d565b9183019190915250845250918601918601610bd5565b5098975050505050505050565b600060208284031215610ce757600080fd5b813567ffffffffffffffff811115610cfe57600080fd5b610d0a8482850161090f565b949350505050565b600060208284031215610d2457600080fd5b5035919050565b60008060008060008060c08789031215610d4457600080fd5b505084359660208601359650604086013595606081013595506080810135945060a0013592509050565b8035801515811461063e57600080fd5b60008060408385031215610d9157600080fd5b82359150610da160208401610d6e565b90509250929050565b600060408284031215610dbc57600080fd5b6040516040810181811067ffffffffffffffff82111715610ddf57610ddf610729565b60405282358152610df260208401610d6e565b60208201529392505050565b60008060008060808587031215610e1457600080fd5b84359350610e2460208601610643565b9250610e3260408601610643565b9396929550929360600135925050565b600060208284031215610e5457600080fd5b6040516020810181811067ffffffffffffffff82111715610e7757610e77610729565b6040529135825250919050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b83811015610f63577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0898403018552815160a073ffffffffffffffffffffffffffffffffffffffff825116855288820151818a870152610f0f828701826106b2565b9150508782015185820389870152610f2782826106b2565b915050606080830151818701525060808083015192508582038187015250610f4f81836106b2565b968901969450505090860190600101610ead565b509098975050505050505050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b83811015610f63577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0898403018552815160a08151818652610fde828701826106b2565b91505073ffffffffffffffffffffffffffffffffffffffff89830151168986015263ffffffff8883015116888601526060808301518683038288015261102483826106b2565b6080948501519790940196909652505094870194925090860190600101610f9a56fea164736f6c6343000818000a",
}

var MessageHasherABI = MessageHasherMetaData.ABI

var MessageHasherBin = MessageHasherMetaData.Bin

func DeployMessageHasher(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *generated.Transaction, *MessageHasher, error) {
	parsed, err := MessageHasherMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(MessageHasherZKBin), backend)
		contractReturn := &MessageHasher{address: address, abi: *parsed, MessageHasherCaller: MessageHasherCaller{contract: contractBind}, MessageHasherTransactor: MessageHasherTransactor{contract: contractBind}, MessageHasherFilterer: MessageHasherFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageHasherBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, HashZks: tx.Hash()}, &MessageHasher{address: address, abi: *parsed, MessageHasherCaller: MessageHasherCaller{contract: contract}, MessageHasherTransactor: MessageHasherTransactor{contract: contract}, MessageHasherFilterer: MessageHasherFilterer{contract: contract}}, nil
}

type MessageHasher struct {
	address common.Address
	abi     abi.ABI
	MessageHasherCaller
	MessageHasherTransactor
	MessageHasherFilterer
}

type MessageHasherCaller struct {
	contract *bind.BoundContract
}

type MessageHasherTransactor struct {
	contract *bind.BoundContract
}

type MessageHasherFilterer struct {
	contract *bind.BoundContract
}

type MessageHasherSession struct {
	Contract     *MessageHasher
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MessageHasherCallerSession struct {
	Contract *MessageHasherCaller
	CallOpts bind.CallOpts
}

type MessageHasherTransactorSession struct {
	Contract     *MessageHasherTransactor
	TransactOpts bind.TransactOpts
}

type MessageHasherRaw struct {
	Contract *MessageHasher
}

type MessageHasherCallerRaw struct {
	Contract *MessageHasherCaller
}

type MessageHasherTransactorRaw struct {
	Contract *MessageHasherTransactor
}

func NewMessageHasher(address common.Address, backend bind.ContractBackend) (*MessageHasher, error) {
	abi, err := abi.JSON(strings.NewReader(MessageHasherABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMessageHasher(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageHasher{address: address, abi: abi, MessageHasherCaller: MessageHasherCaller{contract: contract}, MessageHasherTransactor: MessageHasherTransactor{contract: contract}, MessageHasherFilterer: MessageHasherFilterer{contract: contract}}, nil
}

func NewMessageHasherCaller(address common.Address, caller bind.ContractCaller) (*MessageHasherCaller, error) {
	contract, err := bindMessageHasher(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageHasherCaller{contract: contract}, nil
}

func NewMessageHasherTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageHasherTransactor, error) {
	contract, err := bindMessageHasher(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageHasherTransactor{contract: contract}, nil
}

func NewMessageHasherFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageHasherFilterer, error) {
	contract, err := bindMessageHasher(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageHasherFilterer{contract: contract}, nil
}

func bindMessageHasher(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MessageHasherMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MessageHasher *MessageHasherRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageHasher.Contract.MessageHasherCaller.contract.Call(opts, result, method, params...)
}

func (_MessageHasher *MessageHasherRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageHasher.Contract.MessageHasherTransactor.contract.Transfer(opts)
}

func (_MessageHasher *MessageHasherRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageHasher.Contract.MessageHasherTransactor.contract.Transact(opts, method, params...)
}

func (_MessageHasher *MessageHasherCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageHasher.Contract.contract.Call(opts, result, method, params...)
}

func (_MessageHasher *MessageHasherTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageHasher.Contract.contract.Transfer(opts)
}

func (_MessageHasher *MessageHasherTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageHasher.Contract.contract.Transact(opts, method, params...)
}

func (_MessageHasher *MessageHasherCaller) DecodeEVMExtraArgsV1(opts *bind.CallOpts, gasLimit *big.Int) (ClientEVMExtraArgsV1, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "decodeEVMExtraArgsV1", gasLimit)

	if err != nil {
		return *new(ClientEVMExtraArgsV1), err
	}

	out0 := *abi.ConvertType(out[0], new(ClientEVMExtraArgsV1)).(*ClientEVMExtraArgsV1)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) DecodeEVMExtraArgsV1(gasLimit *big.Int) (ClientEVMExtraArgsV1, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV1(&_MessageHasher.CallOpts, gasLimit)
}

func (_MessageHasher *MessageHasherCallerSession) DecodeEVMExtraArgsV1(gasLimit *big.Int) (ClientEVMExtraArgsV1, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV1(&_MessageHasher.CallOpts, gasLimit)
}

func (_MessageHasher *MessageHasherCaller) DecodeEVMExtraArgsV2(opts *bind.CallOpts, gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "decodeEVMExtraArgsV2", gasLimit, allowOutOfOrderExecution)

	if err != nil {
		return *new(ClientEVMExtraArgsV2), err
	}

	out0 := *abi.ConvertType(out[0], new(ClientEVMExtraArgsV2)).(*ClientEVMExtraArgsV2)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) DecodeEVMExtraArgsV2(gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV2(&_MessageHasher.CallOpts, gasLimit, allowOutOfOrderExecution)
}

func (_MessageHasher *MessageHasherCallerSession) DecodeEVMExtraArgsV2(gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV2(&_MessageHasher.CallOpts, gasLimit, allowOutOfOrderExecution)
}

func (_MessageHasher *MessageHasherCaller) EncodeAny2EVMTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeAny2EVMTokenAmountsHashPreimage", tokenAmounts)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeAny2EVMTokenAmountsHashPreimage(tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeAny2EVMTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmounts)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeAny2EVMTokenAmountsHashPreimage(tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeAny2EVMTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmounts)
}

func (_MessageHasher *MessageHasherCaller) EncodeEVM2AnyTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeEVM2AnyTokenAmountsHashPreimage", tokenAmount)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeEVM2AnyTokenAmountsHashPreimage(tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVM2AnyTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmount)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeEVM2AnyTokenAmountsHashPreimage(tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVM2AnyTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmount)
}

func (_MessageHasher *MessageHasherCaller) EncodeEVMExtraArgsV1(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV1) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeEVMExtraArgsV1", extraArgs)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeEVMExtraArgsV1(extraArgs ClientEVMExtraArgsV1) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV1(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeEVMExtraArgsV1(extraArgs ClientEVMExtraArgsV1) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV1(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCaller) EncodeEVMExtraArgsV2(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV2) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeEVMExtraArgsV2", extraArgs)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeEVMExtraArgsV2(extraArgs ClientEVMExtraArgsV2) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV2(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeEVMExtraArgsV2(extraArgs ClientEVMExtraArgsV2) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV2(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCaller) EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeFinalHashPreimage", leafDomainSeparator, metaDataHash, fixedSizeFieldsHash, senderHash, dataHash, tokenAmountsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, metaDataHash, fixedSizeFieldsHash, senderHash, dataHash, tokenAmountsHash)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, metaDataHash, fixedSizeFieldsHash, senderHash, dataHash, tokenAmountsHash)
}

func (_MessageHasher *MessageHasherCaller) EncodeFixedSizeFieldsHashPreimage(opts *bind.CallOpts, messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeFixedSizeFieldsHashPreimage", messageId, receiver, sequenceNumber, gasLimit, nonce)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeFixedSizeFieldsHashPreimage(messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFixedSizeFieldsHashPreimage(&_MessageHasher.CallOpts, messageId, receiver, sequenceNumber, gasLimit, nonce)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeFixedSizeFieldsHashPreimage(messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFixedSizeFieldsHashPreimage(&_MessageHasher.CallOpts, messageId, receiver, sequenceNumber, gasLimit, nonce)
}

func (_MessageHasher *MessageHasherCaller) EncodeMetadataHashPreimage(opts *bind.CallOpts, any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeMetadataHashPreimage", any2EVMMessageHash, sourceChainSelector, destChainSelector, onRampHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeMetadataHashPreimage(any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeMetadataHashPreimage(&_MessageHasher.CallOpts, any2EVMMessageHash, sourceChainSelector, destChainSelector, onRampHash)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeMetadataHashPreimage(any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeMetadataHashPreimage(&_MessageHasher.CallOpts, any2EVMMessageHash, sourceChainSelector, destChainSelector, onRampHash)
}

func (_MessageHasher *MessageHasherCaller) Hash(opts *bind.CallOpts, message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "hash", message, onRamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) Hash(message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error) {
	return _MessageHasher.Contract.Hash(&_MessageHasher.CallOpts, message, onRamp)
}

func (_MessageHasher *MessageHasherCallerSession) Hash(message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error) {
	return _MessageHasher.Contract.Hash(&_MessageHasher.CallOpts, message, onRamp)
}

func (_MessageHasher *MessageHasher) Address() common.Address {
	return _MessageHasher.address
}

type MessageHasherInterface interface {
	DecodeEVMExtraArgsV1(opts *bind.CallOpts, gasLimit *big.Int) (ClientEVMExtraArgsV1, error)

	DecodeEVMExtraArgsV2(opts *bind.CallOpts, gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error)

	EncodeAny2EVMTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error)

	EncodeEVM2AnyTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error)

	EncodeEVMExtraArgsV1(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV1) ([]byte, error)

	EncodeEVMExtraArgsV2(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV2) ([]byte, error)

	EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error)

	EncodeFixedSizeFieldsHashPreimage(opts *bind.CallOpts, messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error)

	EncodeMetadataHashPreimage(opts *bind.CallOpts, any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error)

	Hash(opts *bind.CallOpts, message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error)

	Address() common.Address
}

var MessageHasherZKBin = ("0x0x00020000000000020006000000000002000000000f01034f00010000000f035500000000030f00190000006003300270000001860030019d0000008004000039000000400040043f0000000100200190000000320000c13d0000018601300197000600000001001d000000040010008c000002990000413d00000000020f043b000000e003200270000001880030009c0000003a0000213d0000018f0030009c00050000000f03530000005d0000a13d000001900030009c000001bc0000613d000001910030009c000001dc0000613d000001920030009c000002990000c13d0000000601000029000000240010008c000002990000413d0000000001000416000000000001004b000002990000c13d0000000401f00370000000000101043b000600000001001d000000a001000039000000400010043f000000800000043f0614041d0000040f0000000602000029000000a00020043f000000400100043d0000000000210435000001860010009c000001860100804100000040011002100000019b011001c7000006150001042e0000000001000416000000000001004b000002990000c13d0000002001000039000001000010044300000120000004430000018701000041000006150001042e000001890030009c000001960000a13d0000018a0030009c000001ef0000613d0000018b0030009c0000021a0000613d0000018c0030009c000002990000c13d0000000601000029000000240010008c000002990000413d0000000001000416000000000001004b000002990000c13d0000000401f00370000000000101043b000000800010043f0000019502000041000000c00020043f000000c40010043f0000002401000039000000a00010043f0000010001000039000000400010043f0000002001000039000001000010043f000000a0010000390000012002000039061404000000040f000001000110008a000001860010009c0000018601008041000000600110021000000196011001c7000006150001042e000001930030009c0000023e0000613d000001940030009c000002990000c13d0000000601000029000000440010008c000002990000413d0000000002000416000000000002004b000002990000c13d0000000402f00370000000000402043b000001970040009c000002990000213d00000006024000690000019e0020009c000002990000213d000001440020008c000002990000413d000001e005000039000000400050043f000000040340003900000000023f034f000000000202043b000001400020043f000000240240003900000000062f034f000000000606043b000001970060009c000002990000213d000001600060043f000000200220003900000000062f034f000000000606043b000001970060009c000002990000213d000001800060043f000000200220003900000000062f034f000000000606043b000001970060009c000002990000213d000001a00060043f000000200220003900000000062f034f000000000606043b000001970060009c000002990000213d000001c00060043f0000014006000039000000800060043f000000200720003900000000027f034f000000000202043b000001970020009c000002990000213d00000000094200190000002302900039000000060020006c000002990000813d000000040a9000390000000002af034f000000000802043b000001a20080009c000001d60000813d0000001f02800039000001a8022001970000003f02200039000001a802200197000001a30020009c000001d60000213d000001e002200039000000400020043f000001e00080043f00000000028900190000002402200039000000060020006c000002990000213d0000002002a00039000000000a2f034f000001a80b8001980000001f0c80018f0000020009b00039000000b70000613d000002000d000039000000000e0a034f00000000e20e043c000000000d2d043600000000009d004b000000b30000c13d00000000000c004b000000c40000613d0000000002ba034f000000030ac00210000000000b090433000000000bab01cf000000000bab022f000000000202043b000001000aa000890000000002a2022f0000000002a201cf0000000002b2019f000000000029043500000200028000390000000000020435000000a00050043f0000002005700039000000050100035f000000000251034f000000000202043b000001970020009c000002990000213d00000000094200190000002302900039000000060020006c000002990000813d000000040a9000390000000002a1034f000000000402043b000001970040009c000001d60000213d0000001f02400039000001a8022001970000003f02200039000001a802200197000000400700043d0000000002270019000000000072004b00000000080000390000000108004039000001970020009c000001d60000213d0000000100800190000001d60000c13d000000400020043f000000000847043600000000024900190000002402200039000000060020006c000002990000213d0000002002a00039000000000921034f000001a80a4001980000001f0b40018f0000000006a80019000000f50000613d000000000c09034f000000000208001900000000cd0c043c0000000002d20436000000000062004b000000f10000c13d00000000000b004b000001020000613d0000000002a9034f0000000309b00210000000000a060433000000000a9a01cf000000000a9a022f000000000202043b0000010009900089000000000292022f00000000029201cf0000000002a2019f000000000026043500000000024800190000000000020435000000c00070043f0000002004500039000000050100035f000000000241034f000000000202043b000001a00020009c000002990000213d000000e00020043f0000002002400039000000000221034f000000000202043b000001000020043f0000004002400039000000000121034f000000000101043b000001970010009c000002990000213d000000000131001900000006020000290614049a0000040f000001200010043f00000024010000390000000101100367000000000101043b000001970010009c000002990000213d00000000020000310000000401100039061404500000040f000000800200043d00000040032000390000000003030433000600000003001d00000020022000390000000002020433000500000002001d0000000012010434061405fa0000040f00000005020000290000019702200197000000400400043d000500000004001d0000004003400039000000000023043500000006020000290000019702200197000000600340003900000000002304350000008002400039000000000012043500000080010000390000000002140436000600000002001d000001a40100004100000000001204350000000001040019061404120000040f000000050100002900000000020104330000000601000029061405fa0000040f000000800300043d00000080023000390000000006020433000000000203043300000060033000390000000004030433000500000001001d000001000500043d000000400700043d000600000007001d000000e00100043d000001a003100197000001970440019700000197066001970000002001700039000400000001001d0614059a0000040f00000006030000290000000002310049000000200120008a000000000013043500000000010300190614043e0000040f000000060100002900000000020104330000000401000029061405fa0000040f000400000001001d000000a00100043d0000000012010434061405fa0000040f000300000001001d000000c00100043d0000000012010434061405fa0000040f000200000001001d000001200200043d000000400100043d000600000001001d0000002001100039000100000001001d061405a80000040f00000006030000290000000002310049000000200120008a000000000013043500000000010300190614043e0000040f000000060100002900000000020104330000000101000029061405fa0000040f000000400400043d000100000004001d000000400240003900000005030000290000000000320435000000600240003900000004030000290000000000320435000000800240003900000003030000290000000000320435000000a00240003900000002030000290000000000320435000000c0024000390000000000120435000000c0010000390000000001140436000600000001001d00000000000104350000000001040019061404330000040f000000010100002900000000020104330000000601000029061405fa0000040f000000400200043d0000000000120435000001860020009c000001860200804100000040012002100000019b011001c7000006150001042e0000018d0030009c000002680000613d0000018e0030009c000002990000c13d0000000601000029000000440010008c000002990000413d0000000001000416000000000001004b000002990000c13d061405860000040f000600000001001d000000400100043d000500000001001d061404280000040f0000000502000029000000200120003900000000000104350000000000020435000000400100043d000500000001001d061404280000040f000000060000006b0000000001000039000000010100c03900000004020000390000000102200367000000000202043b00000005040000290000002003400039000000000013043500000000002404350000000002040019000000400100043d000600000001001d061405910000040f0000000602000029000002110000013d00000006010000290003002400100094000002990000413d0000000002000416000000000002004b000002990000c13d0000000402f00370000000000102043b000400000001001d000001970010009c000002990000213d00000004010000290000002302100039000000060020006c000002990000813d0000000401000029000000040210003900000000022f034f000000000302043b000001970030009c000001d60000213d00000005023002100000003f042000390000019c044001970000019d0040009c0000028f0000a13d000001a501000041000000000010043f0000004101000039000000040010043f000001a60100004100000616000104300000000601000029000000240010008c000002990000413d0000000001000416000000000001004b000002990000c13d0000000401f00370000000000101043b000001970010009c000002990000213d000000040110003900000006020000290614049a0000040f0000000002010019000000400100043d000600000001001d0000002001100039061405a80000040f000002040000013d0000000601000029000000440010008c000002990000413d0000000001000416000000000001004b000002990000c13d000000c001000039000000400010043f0000000401f00370000000000101043b000000800010043f061405860000040f000000a00010043f0000019901000041000000400300043d000600000003001d0000002002300039000000000012043500000080020000390000002401300039061405910000040f00000006030000290000000002310049000000200120008a000000000013043500000000010300190614043e0000040f0000002001000039000000400200043d000500000002001d00000000021204360000000601000029061404000000040f00000005020000290000000001210049000001860010009c00000186010080410000006001100210000001860020009c00000186020080410000004002200210000000000121019f000006150001042e0000000601000029000000840010008c000002990000413d0000000001000416000000000001004b000002990000c13d0000002401f00370000000000201043b000001970020009c000002990000213d0000004401f00370000000000301043b000001970030009c000002990000213d0000000401f00370000000000101043b000000a00010043f000000c00020043f000000e00030043f0000006401f00370000000000101043b000001000010043f000000800040043f0000012001000039000000400010043f0000002001000039000001200010043f00000140020000390000000001040019061404000000040f000001200110008a000001860010009c0000018601008041000000600110021000000198011001c7000006150001042e0000000601000029000000a40010008c000002990000413d0000000002000416000000000002004b000002990000c13d0000002402f00370000000000202043b000001a00020009c000002990000213d0000004403f00370000000000303043b000001970030009c000002990000213d0000008404f00370000000000404043b000001970040009c000002990000213d0000006405f00370000000000505043b0000000401f00370000000000101043b000000a00010043f000000c00020043f000000e00030043f000001000050043f000001200040043f000000a001000039000000800010043f0000014001000039000000400010043f0000002001000039000001400010043f00000080010000390000016002000039061404000000040f000001400110008a000001860010009c00000186010080410000006001100210000001a7011001c7000006150001042e0000000601000029000000c40010008c000002990000413d0000000001000416000000000001004b000002990000c13d0000000401f00370000000000101043b000000a00010043f0000002401f00370000000000101043b000000c00010043f0000004401f00370000000000101043b000000e00010043f0000006401f00370000000000101043b000001000010043f0000008401f00370000000000101043b000001200010043f000000a401f00370000000000101043b000001400010043f000000c001000039000000800010043f0000016001000039000000400010043f0000002001000039000001600010043f00000080010000390000018002000039061404000000040f000001600110008a000001860010009c000001860100804100000060011002100000019a011001c7000006150001042e00000080014000390000000004010019000000400010043f000000800030043f000000040100002900000024051000390000000002520019000200000002001d000000060020006c0000029b0000a13d00000000010000190000061600010430000000000003004b000002b50000c13d00000000010400190000002002100039000000200300003900000000003204350000004002100039000000800f00043d0000000000f2043500000060021000390000000503f00210000000000723001900000000000f004b000003aa0000c13d0000000002170049000000200320008a0000000000310435000500000001001d0614043e0000040f000000400200043d000600000002001d000000200100003900000000021204360000000501000029061404000000040f000001ba0000013d000000a001000039000000200700008a000002c10000013d0000000002ac0019000000000002043500000080029000390000000000b2043500000000019104360000002005500039000000020050006c000000200700008a000003fe0000813d00000000025f034f000000000202043b000001970020009c000002990000213d000000040a2000290000000302a000690000019e0020009c000002990000213d000000a00020008c000002990000413d000000400900043d0000019f0090009c000001d60000213d000000a002900039000000400020043f0000002402a0003900000000032f034f000000000303043b000001a00030009c000002990000213d000000000c390436000000200b2000390000000002bf034f000000000202043b000001970020009c000002990000213d0000000003a200190000004302300039000000060020006c0000000004000019000001a104008041000001a102200197000000000002004b0000000006000019000001a106004041000001a10020009c000000000604c019000000000006004b000002990000c13d000000240430003900000000024f034f000000000d02043b0000019700d0009c000001d60000213d0000001f02d00039000000000272016f0000003f02200039000000000272016f000000400e00043d00000000062e00190000000000e6004b00000000020000390000000102004039000001970060009c000001d60000213d0000000100200190000001d60000c13d000000400060043f0000000002de04360000000003d300190000004403300039000000060030006c000002990000213d000000200340003900000000083f034f00000000067d017000000000036200190000030b0000613d000000000408034f000000000f020019000000004704043c000000000f7f043600000000003f004b000003070000c13d0000001f04d00190000003180000613d000000000668034f0000000304400210000000000703043300000000074701cf000000000747022f000000000606043b0000010004400089000000000646022f00000000044601cf000000000474019f00000000004304350000000002d2001900000000000204350000000000ec0435000000200bb00039000000050f00035f0000000002bf034f000000000202043b000001970020009c000002990000213d0000000003a200190000004302300039000000060020006c0000000004000019000001a104008041000001a102200197000000000002004b0000000006000019000001a106004041000001a10020009c000000000604c019000000000006004b000002990000c13d000000240430003900000000024f034f000000000c02043b0000019700c0009c000001d60000213d0000001f02c00039000001a8022001970000003f02200039000001a802200197000000400d00043d00000000062d00190000000000d6004b00000000020000390000000102004039000001970060009c000001d60000213d0000000100200190000001d60000c13d000000400060043f000000000ecd04360000000002c300190000004402200039000000060020006c000002990000213d000000200240003900000000042f034f000001a806c0019800000000036e0019000003510000613d000000000804034f00000000020e0019000000008708043c0000000002720436000000000032004b0000034d0000c13d0000001f02c001900000035e0000613d000000000464034f0000000302200210000000000603043300000000062601cf000000000626022f000000000404043b0000010002200089000000000424022f00000000022401cf000000000262019f00000000002304350000000002ce0019000000000002043500000040029000390000000000d204350000002002b0003900000000022f034f000000000202043b000000600390003900000000002304350000004002b0003900000000022f034f000000000202043b000001970020009c000002990000213d0000000003a200190000004302300039000000060020006c0000000004000019000001a104008041000001a102200197000000000002004b0000000006000019000001a106004041000001a10020009c000000000604c019000000000006004b000002990000c13d000000240430003900000000024f034f000000000a02043b0000019700a0009c000001d60000213d0000001f02a00039000001a8022001970000003f02200039000001a802200197000000400b00043d00000000062b00190000000000b6004b00000000020000390000000102004039000001970060009c000001d60000213d0000000100200190000001d60000c13d000000400060043f000000000cab04360000000002a300190000004402200039000000060020006c000002990000213d000000200240003900000000042f034f000001a806a0019800000000036c00190000039c0000613d000000000804034f00000000020c0019000000008708043c0000000002720436000000000032004b000003980000c13d0000001f02a00190000002b80000613d000000000464034f0000000302200210000000000603043300000000062601cf000000000626022f000000000404043b0000010002200089000000000424022f00000000022401cf000000000262019f0000000000230435000002b80000013d000000a00300003900000000050000190000000006030019000003b60000013d0000001f09800039000001a80990019700000000088700190000000000080435000000000797001900000001055000390000000000f5004b000002a90000813d0000000008170049000000600880008a0000000002820436000000006806043400000000a9080434000001a0099001970000000009970436000000000a0a04330000000000390435000000a00c70003900000000b90a043400000000009c0435000000c00a700039000000000009004b000003cd0000613d000000000c000019000000000dac0019000000000ecb0019000000000e0e04330000000000ed0435000000200cc0003900000000009c004b000003c60000413d000000000b9a001900000000000b04350000001f09900039000001a80990019700000000099a0019000000400a800039000000000a0a0433000000000b790049000000400c7000390000000000bc043500000000ba0a04340000000009a9043600000000000a004b000003e30000613d000000000c000019000000000d9c0019000000000ecb0019000000000e0e04330000000000ed0435000000200cc000390000000000ac004b000003dc0000413d000000000ba9001900000000000b0435000000600b800039000000000b0b0433000000600c7000390000000000bc04350000001f0aa00039000001a80aa00197000000000aa900190000008008800039000000000808043300000000097a004900000080077000390000000000970435000000009808043400000000078a0436000000000008004b000003ae0000613d000000000a000019000000000b7a0019000000000ca90019000000000c0c04330000000000cb0435000000200aa0003900000000008a004b000003f60000413d000003ae0000013d000000400400043d0000029d0000013d00000000430104340000000001320436000000000003004b0000040c0000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b000004050000413d000000000231001900000000000204350000001f02300039000001a8022001970000000001210019000000000001042d000001a90010009c000004170000813d000000a001100039000000400010043f000000000001042d000001a501000041000000000010043f0000004101000039000000040010043f000001a6010000410000061600010430000001aa0010009c000004220000813d0000002001100039000000400010043f000000000001042d000001a501000041000000000010043f0000004101000039000000040010043f000001a6010000410000061600010430000001ab0010009c0000042d0000813d0000004001100039000000400010043f000000000001042d000001a501000041000000000010043f0000004101000039000000040010043f000001a6010000410000061600010430000001ac0010009c000004380000813d000000e001100039000000400010043f000000000001042d000001a501000041000000000010043f0000004101000039000000040010043f000001a60100004100000616000104300000001f02200039000001a8022001970000000001120019000000000021004b00000000020000390000000102004039000001970010009c0000044a0000213d00000001002001900000044a0000c13d000000400010043f000000000001042d000001a501000041000000000010043f0000004101000039000000040010043f000001a601000041000006160001043000000000030100190000001f01300039000000000021004b0000000004000019000001a104004041000001a105200197000001a101100197000000000651013f000000000051004b0000000001000019000001a101002041000001a10060009c000000000104c019000000000001004b000004980000613d0000000106000367000000000136034f000000000401043b000001a20040009c000004920000813d0000001f01400039000001a8011001970000003f01100039000001a805100197000000400100043d0000000005510019000000000015004b00000000080000390000000108004039000001970050009c000004920000213d0000000100800190000004920000c13d000000400050043f000000000541043600000020033000390000000008430019000000000028004b000004980000213d000000000336034f000001a8064001980000001f0740018f0000000002650019000004820000613d000000000803034f0000000009050019000000008a08043c0000000009a90436000000000029004b0000047e0000c13d000000000007004b0000048f0000613d000000000363034f0000000306700210000000000702043300000000076701cf000000000767022f000000000303043b0000010006600089000000000363022f00000000036301cf000000000373019f000000000032043500000000024500190000000000020435000000000001042d000001a501000041000000000010043f0000004101000039000000040010043f000001a6010000410000061600010430000000000100001900000616000104300006000000000002000500000001001d0000001f01100039000000000021004b0000000003000019000001a103004041000001a10a200197000001a1011001970000000004a1013f0000000000a1004b0000000001000019000001a101002041000001a10040009c000000000103c019000000000001004b0000057e0000613d00000001050003670000000501500360000000000101043b000001a20010009c000005800000813d00000005031002100000003f043000390000019c04400197000000400700043d0000000006470019000100000007001d000000000076004b00000000040000390000000104004039000001970060009c000005800000213d0000000100400190000005800000c13d000000400060043f0000000104000029000000000014043500000005010000290000002006100039000400000036001d000000040020006b0000057e0000213d000000040060006c0000057c0000813d0003002000200092000000010400002900020000000a001d000004d90000013d000000060400002900000020044000390000000001cf001900000000000104350000006001b000390000000000e104350000002001d00039000000000115034f000000000101043b0000008003b0003900000000001304350000000000b404350000002006600039000000040060006c0000057c0000813d000600000004001d000000000165034f000000000101043b000001970010009c0000057e0000213d000000050c1000290000000301c000690000019e0010009c0000057e0000213d000000a00010008c0000057e0000413d000000400b00043d0000019f00b0009c000005800000213d000000200dc000390000000001d5034f000000a00eb000390000004000e0043f000000000101043b000001970010009c0000057e0000213d0000000003c100190000003f01300039000000000021004b0000000004000019000001a104008041000001a1011001970000000007a1013f0000000000a1004b0000000001000019000001a101004041000001a10070009c000000000104c019000000000001004b0000057e0000c13d0000002007300039000000000175034f000000000f01043b0000019700f0009c000005800000213d0000001f01f00039000001a8011001970000003f01100039000001a8011001970000000001e10019000001970010009c000005800000213d000000400010043f0000000000fe04350000000001f300190000004001100039000000000021004b0000057e0000213d0000002001700039000000000915034f000001a801f00198000000c004b0003900000000081400190000051a0000613d000000000709034f0000000003040019000000007a07043c0000000003a30436000000000083004b000005160000c13d0000001f03f00190000005270000613d000000000119034f0000000303300210000000000708043300000000073701cf000000000737022f000000000101043b0000010003300089000000000131022f00000000013101cf000000000171019f00000000001804350000000001f4001900000000000104350000000003eb04360000002001d00039000000000415034f000000000404043b000001a00040009c000000020a0000290000057e0000213d00000000004304350000002001100039000000000315034f000000000303043b000001860030009c0000057e0000213d0000004004b000390000000000340435000000200d1000390000000001d5034f000000000101043b000001970010009c0000057e0000213d0000000003c100190000003f01300039000000000021004b0000000004000019000001a104008041000001a1011001970000000007a1013f0000000000a1004b0000000001000019000001a101004041000001a10070009c000000000104c019000000000001004b0000057e0000c13d0000002008300039000000000185034f000000000c01043b0000019700c0009c000005800000213d0000001f01c00039000001a8011001970000003f01100039000001a801100197000000400e00043d00000000011e00190000000000e1004b00000000040000390000000104004039000001970010009c000005800000213d0000000100400190000005800000c13d000000400010043f000000000fce04360000000001c300190000004001100039000000000021004b0000057e0000213d0000002001800039000000000115034f000001a807c0019800000000037f00190000056e0000613d000000000801034f00000000040f0019000000008908043c0000000004940436000000000034004b0000056a0000c13d0000001f04c00190000004ca0000613d000000000171034f0000000304400210000000000703043300000000074701cf000000000747022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000171019f0000000000130435000004ca0000013d0000000101000029000000000001042d00000000010000190000061600010430000001a501000041000000000010043f0000004101000039000000040010043f000001a601000041000006160001043000000024010000390000000101100367000000000101043b000000000001004b0000000002000039000000010200c039000000000021004b0000058f0000c13d000000000001042d00000000010000190000061600010430000000003202043400000000022104360000000003030433000000000003004b0000000003000039000000010300c03900000000003204350000004001100039000000000001042d00000197066001970000008007100039000000000067043500000060061000390000000000560435000001970440019700000040051000390000000000450435000001a003300197000000200410003900000000003404350000000000210435000000a001100039000000000001042d000100000000000200000000040100190000002003000039000000000334043600000000060204330000000000630435000100000004001d0000004005400039000000050360021000000000033500190000000001060019000000000006004b000005f80000613d0000000008000019000005c30000013d0000000006ba0019000000000006043500000080033000390000008006900039000000000606043300000000006304350000001f03b00039000001a80330019700000000033a00190000000108800039000000000018004b000005f80000813d000000010930006a000000400990008a00000000059504360000002002200039000000000902043300000000da090434000000a006000039000000000c630436000000a00b30003900000000ea0a04340000000000ab0435000000c00b30003900000000000a004b000005d90000613d000000000f0000190000000006bf00190000000007fe001900000000070704330000000000760435000000200ff000390000000000af004b000005d20000413d0000000006ab0019000000000006043500000000060d0433000001a00660019700000000006c0435000000400690003900000000060604330000018606600197000000400730003900000000006704350000001f06a00039000001a80660019700000000066b00190000000007360049000000600a300039000000600b900039000000000b0b043300000000007a043500000000cb0b0434000000000ab6043600000000000b004b000005b70000613d000000000d0000190000000006ad00190000000007dc001900000000070704330000000000760435000000200dd000390000000000bd004b000005f00000413d000005b70000013d0000000001030019000000000001042d000001860010009c00000186010080410000004001100210000001860020009c00000186020080410000006002200210000000000112019f0000000002000414000001860020009c0000018602008041000000c002200210000000000112019f000001ad011001c700008010020000390614060f0000040f00000001002001900000060d0000613d000000000101043b000000000001042d0000000001000019000006160001043000000612002104230000000102000039000000000001042d0000000002000019000000000001042d0000061400000432000006150001042e0000061600010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff000000020000000000000000000000000000004000000100000000000000000000000000000000000000000000000000000000000000000000000000bf0619ac00000000000000000000000000000000000000000000000000000000c7ca9a1700000000000000000000000000000000000000000000000000000000c7ca9a1800000000000000000000000000000000000000000000000000000000e04767b800000000000000000000000000000000000000000000000000000000e733d20900000000000000000000000000000000000000000000000000000000bf0619ad00000000000000000000000000000000000000000000000000000000c63641bd0000000000000000000000000000000000000000000000000000000094b6624a0000000000000000000000000000000000000000000000000000000094b6624b00000000000000000000000000000000000000000000000000000000ae5663d700000000000000000000000000000000000000000000000000000000b17df714000000000000000000000000000000000000000000000000000000003ec7c377000000000000000000000000000000000000000000000000000000008503839d97a657c9000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000000000000000001200000000000000000181dcf1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000200000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffff7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff5f000000000000000000000000ffffffffffffffffffffffffffffffffffffffff80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000fffffffffffffe1f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f4e487b710000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000000000001400000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffff60000000000000000000000000000000000000000000000000ffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffffc0000000000000000000000000000000000000000000000000ffffffffffffff2002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
