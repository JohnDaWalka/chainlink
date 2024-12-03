package usdc_reader_tester

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

var USDCReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"MessageSent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"sourceDomain\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destinationDomain\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"recipient\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"destinationCaller\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sender\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"messageBody\",\"type\":\"bytes\"}],\"name\":\"emitMessageSent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061032c806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806362826f1814610030575b600080fd5b61004361003e366004610129565b610045565b005b600061008d8a8a8a87898c8c8a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506100d292505050565b90507f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036816040516100be9190610228565b60405180910390a150505050505050505050565b606088888888888888886040516020016100f3989796959493929190610279565b604051602081830303815290604052905098975050505050505050565b803563ffffffff8116811461012457600080fd5b919050565b60008060008060008060008060006101008a8c03121561014857600080fd5b6101518a610110565b985061015f60208b01610110565b975061016d60408b01610110565b965060608a0135955060808a0135945060a08a0135935060c08a013567ffffffffffffffff80821682146101a057600080fd5b90935060e08b013590808211156101b657600080fd5b818c0191508c601f8301126101ca57600080fd5b8135818111156101d957600080fd5b8d60208285010111156101eb57600080fd5b6020830194508093505050509295985092959850929598565b60005b8381101561021f578181015183820152602001610207565b50506000910152565b6020815260008251806020840152610247816040850160208701610204565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60007fffffffff00000000000000000000000000000000000000000000000000000000808b60e01b168352808a60e01b166004840152808960e01b166008840152507fffffffffffffffff0000000000000000000000000000000000000000000000008760c01b16600c830152856014830152846034830152836054830152825161030b816074850160208701610204565b91909101607401999850505050505050505056fea164736f6c6343000818000a",
}

var USDCReaderTesterABI = USDCReaderTesterMetaData.ABI

var USDCReaderTesterBin = USDCReaderTesterMetaData.Bin

func DeployUSDCReaderTester(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *generated_zks.Transaction, *USDCReaderTester, error) {
	parsed, err := USDCReaderTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated_zks.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated_zks.DeployContract(auth, parsed, common.FromHex(USDCReaderTesterZKBin), backend)
		contractReturn := &USDCReaderTester{address: address, abi: *parsed, USDCReaderTesterCaller: USDCReaderTesterCaller{contract: contractBind}, USDCReaderTesterTransactor: USDCReaderTesterTransactor{contract: contractBind}, USDCReaderTesterFilterer: USDCReaderTesterFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(USDCReaderTesterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated_zks.Transaction{Transaction: tx, Hash_zks: tx.Hash()}, &USDCReaderTester{address: address, abi: *parsed, USDCReaderTesterCaller: USDCReaderTesterCaller{contract: contract}, USDCReaderTesterTransactor: USDCReaderTesterTransactor{contract: contract}, USDCReaderTesterFilterer: USDCReaderTesterFilterer{contract: contract}}, nil
}

type USDCReaderTester struct {
	address common.Address
	abi     abi.ABI
	USDCReaderTesterCaller
	USDCReaderTesterTransactor
	USDCReaderTesterFilterer
}

type USDCReaderTesterCaller struct {
	contract *bind.BoundContract
}

type USDCReaderTesterTransactor struct {
	contract *bind.BoundContract
}

type USDCReaderTesterFilterer struct {
	contract *bind.BoundContract
}

type USDCReaderTesterSession struct {
	Contract     *USDCReaderTester
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type USDCReaderTesterCallerSession struct {
	Contract *USDCReaderTesterCaller
	CallOpts bind.CallOpts
}

type USDCReaderTesterTransactorSession struct {
	Contract     *USDCReaderTesterTransactor
	TransactOpts bind.TransactOpts
}

type USDCReaderTesterRaw struct {
	Contract *USDCReaderTester
}

type USDCReaderTesterCallerRaw struct {
	Contract *USDCReaderTesterCaller
}

type USDCReaderTesterTransactorRaw struct {
	Contract *USDCReaderTesterTransactor
}

func NewUSDCReaderTester(address common.Address, backend bind.ContractBackend) (*USDCReaderTester, error) {
	abi, err := abi.JSON(strings.NewReader(USDCReaderTesterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUSDCReaderTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTester{address: address, abi: abi, USDCReaderTesterCaller: USDCReaderTesterCaller{contract: contract}, USDCReaderTesterTransactor: USDCReaderTesterTransactor{contract: contract}, USDCReaderTesterFilterer: USDCReaderTesterFilterer{contract: contract}}, nil
}

func NewUSDCReaderTesterCaller(address common.Address, caller bind.ContractCaller) (*USDCReaderTesterCaller, error) {
	contract, err := bindUSDCReaderTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterCaller{contract: contract}, nil
}

func NewUSDCReaderTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*USDCReaderTesterTransactor, error) {
	contract, err := bindUSDCReaderTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterTransactor{contract: contract}, nil
}

func NewUSDCReaderTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*USDCReaderTesterFilterer, error) {
	contract, err := bindUSDCReaderTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterFilterer{contract: contract}, nil
}

func bindUSDCReaderTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := USDCReaderTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_USDCReaderTester *USDCReaderTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCReaderTester.Contract.USDCReaderTesterCaller.contract.Call(opts, result, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.USDCReaderTesterTransactor.contract.Transfer(opts)
}

func (_USDCReaderTester *USDCReaderTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.USDCReaderTesterTransactor.contract.Transact(opts, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCReaderTester.Contract.contract.Call(opts, result, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.contract.Transfer(opts)
}

func (_USDCReaderTester *USDCReaderTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.contract.Transact(opts, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterTransactor) EmitMessageSent(opts *bind.TransactOpts, version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.contract.Transact(opts, "emitMessageSent", version, sourceDomain, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

func (_USDCReaderTester *USDCReaderTesterSession) EmitMessageSent(version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.EmitMessageSent(&_USDCReaderTester.TransactOpts, version, sourceDomain, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

func (_USDCReaderTester *USDCReaderTesterTransactorSession) EmitMessageSent(version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.EmitMessageSent(&_USDCReaderTester.TransactOpts, version, sourceDomain, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

type USDCReaderTesterMessageSentIterator struct {
	Event *USDCReaderTesterMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCReaderTesterMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCReaderTesterMessageSent)
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
		it.Event = new(USDCReaderTesterMessageSent)
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

func (it *USDCReaderTesterMessageSentIterator) Error() error {
	return it.fail
}

func (it *USDCReaderTesterMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCReaderTesterMessageSent struct {
	Arg0 []byte
	Raw  types.Log
}

func (_USDCReaderTester *USDCReaderTesterFilterer) FilterMessageSent(opts *bind.FilterOpts) (*USDCReaderTesterMessageSentIterator, error) {

	logs, sub, err := _USDCReaderTester.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterMessageSentIterator{contract: _USDCReaderTester.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

func (_USDCReaderTester *USDCReaderTesterFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *USDCReaderTesterMessageSent) (event.Subscription, error) {

	logs, sub, err := _USDCReaderTester.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCReaderTesterMessageSent)
				if err := _USDCReaderTester.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

func (_USDCReaderTester *USDCReaderTesterFilterer) ParseMessageSent(log types.Log) (*USDCReaderTesterMessageSent, error) {
	event := new(USDCReaderTesterMessageSent)
	if err := _USDCReaderTester.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_USDCReaderTester *USDCReaderTester) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _USDCReaderTester.abi.Events["MessageSent"].ID:
		return _USDCReaderTester.ParseMessageSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (USDCReaderTesterMessageSent) Topic() common.Hash {
	return common.HexToHash("0x8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036")
}

func (_USDCReaderTester *USDCReaderTester) Address() common.Address {
	return _USDCReaderTester.address
}

type USDCReaderTesterInterface interface {
	EmitMessageSent(opts *bind.TransactOpts, version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error)

	FilterMessageSent(opts *bind.FilterOpts) (*USDCReaderTesterMessageSentIterator, error)

	WatchMessageSent(opts *bind.WatchOpts, sink chan<- *USDCReaderTesterMessageSent) (event.Subscription, error)

	ParseMessageSent(log types.Log) (*USDCReaderTesterMessageSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var USDCReaderTesterZKBin string = ("0x0000008003000039000000400030043f0000000100200190000000b20000c13d00000060021002700000003302200197000000040020008c000000ba0000413d000000000301043b0000003503300197000000360030009c000000ba0000c13d000001040020008c000000ba0000413d0000000003000416000000000003004b000000ba0000c13d0000000403100370000000000303043b000000330030009c000000ba0000213d0000002404100370000000000404043b000000330040009c000000ba0000213d0000004405100370000000000505043b000000330050009c000000ba0000213d000000c406100370000000000606043b000000370060009c000000ba0000213d000000e407100370000000000907043b000000370090009c000000ba0000213d0000002307900039000000000027004b000000ba0000813d0000000408900039000000000781034f000000000707043b000000370070009c000000ba0000213d00000000097900190000002409900039000000000029004b000000ba0000213d0000001f097000390000003f099001970000003f099000390000003f09900197000000380090009c000000bc0000813d0000008009900039000000400090043f0000002008800039000000000981034f000000800070043f0000003f0a7001980000001f0b70018f000000a008a00039000000460000613d000000a00c000039000000000d09034f00000000de0d043c000000000cec043600000000008c004b000000420000c13d00000000000b004b000000530000613d0000000009a9034f000000030ab00210000000000b080433000000000bab01cf000000000bab022f000000000909043b000001000aa000890000000009a9022f0000000009a901cf0000000009b9019f0000000000980435000000a0077000390000000000070435000000e008300210000000400700043d00000020037000390000000000830435000000e00440021000000024087000390000000000480435000000e00450021000000028057000390000000000450435000000c0046002100000002c057000390000000000450435000000a404100370000000000404043b000000340570003900000000004504350000006404100370000000000404043b000000540570003900000000004504350000008401100370000000000101043b000000740470003900000000001404350000009404700039000000800100043d000000000001004b0000007a0000613d00000000050000190000000006450019000000a008500039000000000808043300000000008604350000002005500039000000000015004b000000730000413d0000000004410019000000000004043500000074041000390000000000470435000000b3011000390000003f041001970000000001740019000000000041004b00000000040000390000000104004039000000370010009c000000bc0000213d0000000100400190000000bc0000c13d000000400010043f00000020040000390000000005410436000000000407043300000000004504350000004005100039000000000004004b000000980000613d000000000600001900000000075600190000000008360019000000000808043300000000008704350000002006600039000000000046004b000000910000413d0000001f034000390000003f023001970000000003540019000000000003043500000040022000390000006003200210000000390020009c0000003a03008041000000330010009c00000033010080410000004001100210000000000113019f0000000002000414000000330020009c0000003302008041000000c00220021000000000012100190000003b0110009a0000800d0200003900000001030000390000003c0400004100c700c20000040f0000000100200190000000ba0000613d0000000001000019000000c80001042e0000000001000416000000000001004b000000ba0000c13d0000002001000039000001000010044300000120000004430000003401000041000000c80001042e0000000001000019000000c9000104300000003d01000041000000000010043f0000004101000039000000040010043f0000003e01000041000000c900010430000000c5002104210000000102000039000000000001042d0000000002000019000000000001042d000000c700000432000000c80001042e000000c9000104300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffff0000000200000000000000000000000000000040000001000000000000000000ffffffff0000000000000000000000000000000000000000000000000000000062826f1800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff80000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffff000000000000000000000000fe000000000000000000000000000000000000000000000000000000000000008c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b0364e487b71000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00000000000000000000000000000000000000000000000000000000000000000")
