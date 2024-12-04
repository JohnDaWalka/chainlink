package report_codec

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

type IRMNRemoteSignature struct {
	R [32]byte
	S [32]byte
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

type InternalExecutionReport struct {
	SourceChainSelector uint64
	Messages            []InternalAny2EVMRampMessage
	OffchainTokenData   [][][]byte
	Proofs              [][32]byte
	ProofFlagBits       *big.Int
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type OffRampCommitReport struct {
	PriceUpdates  InternalPriceUpdates
	MerkleRoots   []InternalMerkleRoot
	RmnSignatures []IRMNRemoteSignature
}

var ReportCodecMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportDecoded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structInternal.ExecutionReport[]\",\"name\":\"report\",\"type\":\"tuple[]\"}],\"name\":\"ExecuteReportDecoded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"decodeCommitReport\",\"outputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"}],\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"decodeExecuteReport\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReport[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506113c3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636fb349561461003b578063f816ec6014610064575b600080fd5b61004e610049366004610231565b610084565b60405161005b91906104ee565b60405180910390f35b610077610072366004610231565b6100a0565b60405161005b9190610833565b60608180602001905181019061009a9190610e6d565b92915050565b6040805160a08101825260608082018181526080830182905282526020808301829052928201528251909161009a918401810190840161122d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff8111828210171561012d5761012d6100db565b60405290565b60405160c0810167ffffffffffffffff8111828210171561012d5761012d6100db565b6040805190810167ffffffffffffffff8111828210171561012d5761012d6100db565b6040516060810167ffffffffffffffff8111828210171561012d5761012d6100db565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156101e3576101e36100db565b604052919050565b600067ffffffffffffffff821115610205576102056100db565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006020828403121561024357600080fd5b813567ffffffffffffffff81111561025a57600080fd5b8201601f8101841361026b57600080fd5b803561027e610279826101eb565b61019c565b81815285602083850101111561029357600080fd5b81602084016020830137600091810160200191909152949350505050565b60005b838110156102cc5781810151838201526020016102b4565b50506000910152565b600081518084526102ed8160208601602086016102b1565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600082825180855260208086019550808260051b84010181860160005b848110156103eb577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a08151818652610380828701826102d5565b91505073ffffffffffffffffffffffffffffffffffffffff868301511686860152604063ffffffff81840151168187015250606080830151868303828801526103c983826102d5565b608094850151979094019690965250509884019892509083019060010161033c565b5090979650505050505050565b6000828251808552602080860195506005818360051b8501018287016000805b868110156104a3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe088850381018c5283518051808752908801908887019080891b88018a01865b8281101561048c57858a830301845261047a8286516102d5565b948c0194938c01939150600101610460565b509e8a019e97505050938701935050600101610418565b50919998505050505050505050565b60008151808452602080850194506020840160005b838110156104e3578151875295820195908201906001016104c7565b509495945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156106d6577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452815160a0860167ffffffffffffffff8083511688528883015160a08a8a015282815180855260c08b01915060c08160051b8c010194508b8301925060005b8181101561067f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff408c87030183528351805180518852868f820151168f890152866040820151166040890152866060820151166060890152866080820151166080890152508d81015161014060a08901526106026101408901826102d5565b9050604082015188820360c08a015261061b82826102d5565b915050606082015161064560e08a018273ffffffffffffffffffffffffffffffffffffffff169052565b50608082015161010089015260a0820151915087810361012089015261066b818361031f565b97505050928c0192918c0191600101610582565b50505050506040820151878203604089015261069b82826103f8565b915050606082015187820360608901526106b582826104b2565b60809384015198909301979097525094509285019290850190600101610515565b5092979650505050505050565b60008151808452602080850194506020840160005b838110156104e3578151805167ffffffffffffffff1688528301517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1683880152604090960195908201906001016106f8565b600082825180855260208086019550808260051b84010181860160005b848110156103eb577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a067ffffffffffffffff8083511686528683015182888801526107bb838801826102d5565b60408581015184169089015260608086015190931692880192909252506080928301519290950191909152509783019790830190600101610764565b60008151808452602080850194506020840160005b838110156104e357815180518852830151838801526040909601959082019060010161080c565b602080825282516060838301528051604060808501819052815160c086018190526000949392840191859160e08801905b808410156108c1578451805173ffffffffffffffffffffffffffffffffffffffff1683528701517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1687830152938601936001939093019290820190610864565b50938501518785037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff800160a0890152936108fb81866106e3565b9450505050508185015191507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08085830301604086015261093c8284610747565b925060408601519150808584030160608601525061095a82826107f7565b95945050505050565b600067ffffffffffffffff82111561097d5761097d6100db565b5060051b60200190565b805167ffffffffffffffff8116811461099f57600080fd5b919050565b600060a082840312156109b657600080fd5b6109be61010a565b9050815181526109d060208301610987565b60208201526109e160408301610987565b60408201526109f260608301610987565b6060820152610a0360808301610987565b608082015292915050565b600082601f830112610a1f57600080fd5b8151610a2d610279826101eb565b818152846020838601011115610a4257600080fd5b610a538260208301602087016102b1565b949350505050565b805173ffffffffffffffffffffffffffffffffffffffff8116811461099f57600080fd5b600082601f830112610a9057600080fd5b81516020610aa061027983610963565b82815260059290921b84018101918181019086841115610abf57600080fd5b8286015b84811015610bbb57805167ffffffffffffffff80821115610ae45760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610b1d5760008081fd5b610b2561010a565b8784015183811115610b375760008081fd5b610b458d8a83880101610a0e565b8252506040610b55818601610a5b565b8983015260608086015163ffffffff81168114610b725760008081fd5b808385015250608091508186015185811115610b8e5760008081fd5b610b9c8f8c838a0101610a0e565b9184019190915250919093015190830152508352918301918301610ac3565b509695505050505050565b600082601f830112610bd757600080fd5b81516020610be761027983610963565b82815260059290921b84018101918181019086841115610c0657600080fd5b8286015b84811015610bbb57805167ffffffffffffffff80821115610c2b5760008081fd5b8189019150610140807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610c655760008081fd5b610c6d610133565b610c798c8986016109a4565b815260c084015183811115610c8e5760008081fd5b610c9c8d8a83880101610a0e565b898301525060e084015183811115610cb45760008081fd5b610cc28d8a83880101610a0e565b604083015250610cd56101008501610a5b565b60608201526101208401516080820152908301519082821115610cf85760008081fd5b610d068c8984870101610a7f565b60a08201528652505050918301918301610c0a565b600082601f830112610d2c57600080fd5b81516020610d3c61027983610963565b82815260059290921b84018101918181019086841115610d5b57600080fd5b8286015b84811015610bbb57805167ffffffffffffffff80821115610d7f57600080fd5b818901915089603f830112610d9357600080fd5b85820151610da361027982610963565b81815260059190911b830160400190878101908c831115610dc357600080fd5b604085015b83811015610dfc57805185811115610ddf57600080fd5b610dee8f6040838a0101610a0e565b845250918901918901610dc8565b50875250505092840192508301610d5f565b600082601f830112610e1f57600080fd5b81516020610e2f61027983610963565b8083825260208201915060208460051b870101935086841115610e5157600080fd5b602086015b84811015610bbb5780518352918301918301610e56565b60006020808385031215610e8057600080fd5b825167ffffffffffffffff80821115610e9857600080fd5b818501915085601f830112610eac57600080fd5b8151610eba61027982610963565b81815260059190911b83018401908481019088831115610ed957600080fd5b8585015b83811015610fd357805185811115610ef457600080fd5b860160a0818c037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0011215610f295760008081fd5b610f3161010a565b610f3c898301610987565b815260408083015188811115610f525760008081fd5b610f608e8c83870101610bc6565b8b8401525060608084015189811115610f795760008081fd5b610f878f8d83880101610d1b565b8385015250608091508184015189811115610fa25760008081fd5b610fb08f8d83880101610e0e565b918401919091525060a09290920151918101919091528352918601918601610edd565b5098975050505050505050565b80517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461099f57600080fd5b600082601f83011261101d57600080fd5b8151602061102d61027983610963565b82815260069290921b8401810191818101908684111561104c57600080fd5b8286015b84811015610bbb57604081890312156110695760008081fd5b611071610156565b61107a82610987565b8152611087858301610fe0565b81860152835291830191604001611050565b600082601f8301126110aa57600080fd5b815160206110ba61027983610963565b82815260059290921b840181019181810190868411156110d957600080fd5b8286015b84811015610bbb57805167ffffffffffffffff808211156110fe5760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156111375760008081fd5b61113f61010a565b61114a888501610987565b8152604080850151848111156111605760008081fd5b61116e8e8b83890101610a0e565b8a8401525060609350611182848601610987565b908201526080611193858201610987565b938201939093529201519082015283529183019183016110dd565b600082601f8301126111bf57600080fd5b815160206111cf61027983610963565b82815260069290921b840181019181810190868411156111ee57600080fd5b8286015b84811015610bbb576040818903121561120b5760008081fd5b611213610156565b8151815284820151858201528352918301916040016111f2565b6000602080838503121561124057600080fd5b825167ffffffffffffffff8082111561125857600080fd5b908401906060828703121561126c57600080fd5b611274610179565b82518281111561128357600080fd5b8301604081890381131561129657600080fd5b61129e610156565b8251858111156112ad57600080fd5b8301601f81018b136112be57600080fd5b80516112cc61027982610963565b81815260069190911b8201890190898101908d8311156112eb57600080fd5b928a01925b828410156113395785848f0312156113085760008081fd5b611310610156565b61131985610a5b565b81526113268c8601610fe0565b818d0152825292850192908a01906112f0565b84525050508287015191508482111561135157600080fd5b61135d8a83850161100c565b8188015283525050828401518281111561137657600080fd5b61138288828601611099565b8583015250604083015193508184111561139b57600080fd5b6113a7878585016111ae565b6040820152969550505050505056fea164736f6c6343000818000a",
}

var ReportCodecABI = ReportCodecMetaData.ABI

var ReportCodecBin = ReportCodecMetaData.Bin

func DeployReportCodec(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *generated.Transaction, *ReportCodec, error) {
	parsed, err := ReportCodecMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(ReportCodecZKBin), backend)
		contractReturn := &ReportCodec{address: address, abi: *parsed, ReportCodecCaller: ReportCodecCaller{contract: contractBind}, ReportCodecTransactor: ReportCodecTransactor{contract: contractBind}, ReportCodecFilterer: ReportCodecFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReportCodecBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, HashZks: tx.Hash()}, &ReportCodec{address: address, abi: *parsed, ReportCodecCaller: ReportCodecCaller{contract: contract}, ReportCodecTransactor: ReportCodecTransactor{contract: contract}, ReportCodecFilterer: ReportCodecFilterer{contract: contract}}, nil
}

type ReportCodec struct {
	address common.Address
	abi     abi.ABI
	ReportCodecCaller
	ReportCodecTransactor
	ReportCodecFilterer
}

type ReportCodecCaller struct {
	contract *bind.BoundContract
}

type ReportCodecTransactor struct {
	contract *bind.BoundContract
}

type ReportCodecFilterer struct {
	contract *bind.BoundContract
}

type ReportCodecSession struct {
	Contract     *ReportCodec
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ReportCodecCallerSession struct {
	Contract *ReportCodecCaller
	CallOpts bind.CallOpts
}

type ReportCodecTransactorSession struct {
	Contract     *ReportCodecTransactor
	TransactOpts bind.TransactOpts
}

type ReportCodecRaw struct {
	Contract *ReportCodec
}

type ReportCodecCallerRaw struct {
	Contract *ReportCodecCaller
}

type ReportCodecTransactorRaw struct {
	Contract *ReportCodecTransactor
}

func NewReportCodec(address common.Address, backend bind.ContractBackend) (*ReportCodec, error) {
	abi, err := abi.JSON(strings.NewReader(ReportCodecABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindReportCodec(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReportCodec{address: address, abi: abi, ReportCodecCaller: ReportCodecCaller{contract: contract}, ReportCodecTransactor: ReportCodecTransactor{contract: contract}, ReportCodecFilterer: ReportCodecFilterer{contract: contract}}, nil
}

func NewReportCodecCaller(address common.Address, caller bind.ContractCaller) (*ReportCodecCaller, error) {
	contract, err := bindReportCodec(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReportCodecCaller{contract: contract}, nil
}

func NewReportCodecTransactor(address common.Address, transactor bind.ContractTransactor) (*ReportCodecTransactor, error) {
	contract, err := bindReportCodec(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReportCodecTransactor{contract: contract}, nil
}

func NewReportCodecFilterer(address common.Address, filterer bind.ContractFilterer) (*ReportCodecFilterer, error) {
	contract, err := bindReportCodec(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReportCodecFilterer{contract: contract}, nil
}

func bindReportCodec(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReportCodecMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ReportCodec *ReportCodecRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReportCodec.Contract.ReportCodecCaller.contract.Call(opts, result, method, params...)
}

func (_ReportCodec *ReportCodecRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReportCodec.Contract.ReportCodecTransactor.contract.Transfer(opts)
}

func (_ReportCodec *ReportCodecRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReportCodec.Contract.ReportCodecTransactor.contract.Transact(opts, method, params...)
}

func (_ReportCodec *ReportCodecCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReportCodec.Contract.contract.Call(opts, result, method, params...)
}

func (_ReportCodec *ReportCodecTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReportCodec.Contract.contract.Transfer(opts)
}

func (_ReportCodec *ReportCodecTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReportCodec.Contract.contract.Transact(opts, method, params...)
}

func (_ReportCodec *ReportCodecCaller) DecodeCommitReport(opts *bind.CallOpts, report []byte) (OffRampCommitReport, error) {
	var out []interface{}
	err := _ReportCodec.contract.Call(opts, &out, "decodeCommitReport", report)

	if err != nil {
		return *new(OffRampCommitReport), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampCommitReport)).(*OffRampCommitReport)

	return out0, err

}

func (_ReportCodec *ReportCodecSession) DecodeCommitReport(report []byte) (OffRampCommitReport, error) {
	return _ReportCodec.Contract.DecodeCommitReport(&_ReportCodec.CallOpts, report)
}

func (_ReportCodec *ReportCodecCallerSession) DecodeCommitReport(report []byte) (OffRampCommitReport, error) {
	return _ReportCodec.Contract.DecodeCommitReport(&_ReportCodec.CallOpts, report)
}

func (_ReportCodec *ReportCodecCaller) DecodeExecuteReport(opts *bind.CallOpts, report []byte) ([]InternalExecutionReport, error) {
	var out []interface{}
	err := _ReportCodec.contract.Call(opts, &out, "decodeExecuteReport", report)

	if err != nil {
		return *new([]InternalExecutionReport), err
	}

	out0 := *abi.ConvertType(out[0], new([]InternalExecutionReport)).(*[]InternalExecutionReport)

	return out0, err

}

func (_ReportCodec *ReportCodecSession) DecodeExecuteReport(report []byte) ([]InternalExecutionReport, error) {
	return _ReportCodec.Contract.DecodeExecuteReport(&_ReportCodec.CallOpts, report)
}

func (_ReportCodec *ReportCodecCallerSession) DecodeExecuteReport(report []byte) ([]InternalExecutionReport, error) {
	return _ReportCodec.Contract.DecodeExecuteReport(&_ReportCodec.CallOpts, report)
}

type ReportCodecCommitReportDecodedIterator struct {
	Event *ReportCodecCommitReportDecoded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ReportCodecCommitReportDecodedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReportCodecCommitReportDecoded)
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
		it.Event = new(ReportCodecCommitReportDecoded)
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

func (it *ReportCodecCommitReportDecodedIterator) Error() error {
	return it.fail
}

func (it *ReportCodecCommitReportDecodedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ReportCodecCommitReportDecoded struct {
	Report OffRampCommitReport
	Raw    types.Log
}

func (_ReportCodec *ReportCodecFilterer) FilterCommitReportDecoded(opts *bind.FilterOpts) (*ReportCodecCommitReportDecodedIterator, error) {

	logs, sub, err := _ReportCodec.contract.FilterLogs(opts, "CommitReportDecoded")
	if err != nil {
		return nil, err
	}
	return &ReportCodecCommitReportDecodedIterator{contract: _ReportCodec.contract, event: "CommitReportDecoded", logs: logs, sub: sub}, nil
}

func (_ReportCodec *ReportCodecFilterer) WatchCommitReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecCommitReportDecoded) (event.Subscription, error) {

	logs, sub, err := _ReportCodec.contract.WatchLogs(opts, "CommitReportDecoded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ReportCodecCommitReportDecoded)
				if err := _ReportCodec.contract.UnpackLog(event, "CommitReportDecoded", log); err != nil {
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

func (_ReportCodec *ReportCodecFilterer) ParseCommitReportDecoded(log types.Log) (*ReportCodecCommitReportDecoded, error) {
	event := new(ReportCodecCommitReportDecoded)
	if err := _ReportCodec.contract.UnpackLog(event, "CommitReportDecoded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ReportCodecExecuteReportDecodedIterator struct {
	Event *ReportCodecExecuteReportDecoded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ReportCodecExecuteReportDecodedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReportCodecExecuteReportDecoded)
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
		it.Event = new(ReportCodecExecuteReportDecoded)
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

func (it *ReportCodecExecuteReportDecodedIterator) Error() error {
	return it.fail
}

func (it *ReportCodecExecuteReportDecodedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ReportCodecExecuteReportDecoded struct {
	Report []InternalExecutionReport
	Raw    types.Log
}

func (_ReportCodec *ReportCodecFilterer) FilterExecuteReportDecoded(opts *bind.FilterOpts) (*ReportCodecExecuteReportDecodedIterator, error) {

	logs, sub, err := _ReportCodec.contract.FilterLogs(opts, "ExecuteReportDecoded")
	if err != nil {
		return nil, err
	}
	return &ReportCodecExecuteReportDecodedIterator{contract: _ReportCodec.contract, event: "ExecuteReportDecoded", logs: logs, sub: sub}, nil
}

func (_ReportCodec *ReportCodecFilterer) WatchExecuteReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecExecuteReportDecoded) (event.Subscription, error) {

	logs, sub, err := _ReportCodec.contract.WatchLogs(opts, "ExecuteReportDecoded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ReportCodecExecuteReportDecoded)
				if err := _ReportCodec.contract.UnpackLog(event, "ExecuteReportDecoded", log); err != nil {
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

func (_ReportCodec *ReportCodecFilterer) ParseExecuteReportDecoded(log types.Log) (*ReportCodecExecuteReportDecoded, error) {
	event := new(ReportCodecExecuteReportDecoded)
	if err := _ReportCodec.contract.UnpackLog(event, "ExecuteReportDecoded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ReportCodec *ReportCodec) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ReportCodec.abi.Events["CommitReportDecoded"].ID:
		return _ReportCodec.ParseCommitReportDecoded(log)
	case _ReportCodec.abi.Events["ExecuteReportDecoded"].ID:
		return _ReportCodec.ParseExecuteReportDecoded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ReportCodecCommitReportDecoded) Topic() common.Hash {
	return common.HexToHash("0x31a4e1cb25733cdb9679561cd59cdc238d70a7d486f8bfc1f13242efd60fc29d")
}

func (ReportCodecExecuteReportDecoded) Topic() common.Hash {
	return common.HexToHash("0x9467c8093a35a72f74398d5b6e351d67dc82eddc378efc6177eafb4fc7a01d39")
}

func (_ReportCodec *ReportCodec) Address() common.Address {
	return _ReportCodec.address
}

type ReportCodecInterface interface {
	DecodeCommitReport(opts *bind.CallOpts, report []byte) (OffRampCommitReport, error)

	DecodeExecuteReport(opts *bind.CallOpts, report []byte) ([]InternalExecutionReport, error)

	FilterCommitReportDecoded(opts *bind.FilterOpts) (*ReportCodecCommitReportDecodedIterator, error)

	WatchCommitReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecCommitReportDecoded) (event.Subscription, error)

	ParseCommitReportDecoded(log types.Log) (*ReportCodecCommitReportDecoded, error)

	FilterExecuteReportDecoded(opts *bind.FilterOpts) (*ReportCodecExecuteReportDecodedIterator, error)

	WatchExecuteReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecExecuteReportDecoded) (event.Subscription, error)

	ParseExecuteReportDecoded(log types.Log) (*ReportCodecExecuteReportDecoded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var ReportCodecZKBin = ("0x000e0000000000020000000100200190000000890000c13d0000006002100270000001f904200197000001a003000039000000400030043f000000040040008c000000b70000413d000000000201043b000000e002200270000001fb0020009c000000930000613d000001fc0020009c000000b70000c13d000000240040008c000000b70000413d0000000002000416000000000002004b000000b70000c13d0000000402100370000000000602043b000001fd0060009c000000b70000213d0000002302600039000000000042004b000000b70000813d0000000407600039000000000271034f000000000502043b000002070050009c000000aa0000813d0000001f025000390000020c022001970000003f022000390000020c02200197000001fe0020009c000000aa0000213d000001a002200039000000400020043f000001a00050043f00000000025600190000002402200039000000000042004b000000b70000213d0000002002700039000000000421034f0000020c065001980000001f0750018f000001c001600039000000390000613d000001c008000039000000000904034f000000009209043c0000000008280436000000000018004b000000350000c13d000000000007004b000000460000613d000000000264034f0000000304700210000000000601043300000000064601cf000000000646022f000000000202043b0000010004400089000000000242022f00000000024201cf000000000262019f0000000000210435000001c0015000390000000000010435000000e00030043f000001a00300043d000001800030043f000002010030009c000000b70000213d000000200030008c000000b70000413d000001c00100043d000001200010043f000001fd0010009c000000b70000213d000001c002300039000001df03100039000000000023004b0000000004000019000002020400804100000202022001970000020203300197000000000523013f000000000023004b00000000020000190000020202004041000002020050009c000000000204c019000000000002004b000000b70000c13d000001c0011000390000000001010433000001000010043f000001fd0010009c000000aa0000213d00000005021002100000003f022000390000020302200197000000400300043d0000000002230019000000000032004b00000000040000390000000104004039000001fd0020009c000000aa0000213d0000000100400190000000aa0000c13d000000400020043f000000a00030043f000001600030043f0000000000130435000000a00100043d0000002001100039000000a00010043f000001400010043f000000e00100043d000001800200043d0000000005120019000001200200043d00000000031200190000004004300039000001000200043d000000050220021000000000022400190000002005500039000000000052004b000000b70000213d000000c00040043f000000d40000013d0000008001000039000000400010043f0000000001000416000000000001004b000000b70000c13d000000200100003900000100001004430000012000000443000001fa01000041000007df0001042e000000240040008c000000b70000413d0000000002000416000000000002004b000000b70000c13d0000000402100370000000000502043b000001fd0050009c000000b70000213d0000002302500039000000000042004b000000b70000813d0000000406500039000000000261034f000000000302043b000001fd0030009c000000aa0000213d0000001f073000390000020c077001970000003f077000390000020c07700197000001fe0070009c000000b00000a13d0000020a01000041000000000010043f0000004101000039000000040010043f0000020b01000041000007e000010430000001a007700039000000400070043f000001a00030043f00000000053500190000002405500039000000000045004b000004540000a13d0000000001000019000007e000010430000000030400002900000060014000390000000000310435000001200100043d000000800200043d0000000001120019000000e00200043d0000000001210019000000c001100039000000000101043300000080024000390000000000120435000000a00100043d0000000000410435000000a00100043d0000002001100039000000a00010043f000000c00100043d0000002004100039000000c00040043f000001200200043d000000e00100043d0000000003210019000001000200043d000000050220021000000000023200190000004002200039000000000024004b0000064e0000813d0000000002040433000000800020043f000001fd0020009c000000b70000213d00000000033200190000000001310049000001800200043d0000000001210019000000200110008a000002010010009c000000b70000213d000000a00010008c000000b70000413d000000400100043d000300000001001d000002000010009c000000aa0000213d0000000301000029000000a001100039000000400010043f00000040013000390000000001010433000001fd0010009c000000b70000213d00000003020000290000000001120436000100000001001d000000800100043d000001200200043d0000000001120019000000e00200043d000000000121001900000060031000390000000003030433000b00000003001d000001fd0030009c000000b70000213d0000000b01100029000001800300043d000000000232001900000020022000390000005f03100039000000000023004b0000000004000019000002020400804100000202022001970000020203300197000000000523013f000000000023004b00000000020000190000020202004041000002020050009c000000000204c019000000000002004b000000b70000c13d00000040011000390000000003010433000001fd0030009c000000aa0000213d00000005013002100000003f021000390000020302200197000000400400043d0000000002240019000200000004001d000000000042004b00000000040000390000000104004039000001fd0020009c000000aa0000213d0000000100400190000000aa0000c13d000000400020043f00000002020000290000000000320435000000e00300043d0000000b02300029000000800400043d0000000002420019000001200400043d0000000004420019000001800200043d0000000002320019000800600010003d00000008014000290000002002200039000000000021004b000000b70000213d0000006002400039000000000012004b000003360000813d0000000b01000029000400e00010003d000500c00010003d000600a00010003d000700800010003d0000000205000029000001490000013d000000090500002900000020055000390000000d02000029000000a0012000390000000c0300002900000000003104350000000000250435000000800100043d0000000b01100029000001200200043d0000000001210019000000e00300043d000000000431001900000008014000290000000a02000029000000000012004b000003360000813d000900000005001d0000000021020434000a00000002001d000001fd0010009c000000b70000213d00000000044100190000000002430049000001800300043d0000000002320019000000400220008a000002010020009c000000b70000213d000001400020008c000000b70000413d000000400200043d000d00000002001d000002080020009c000000aa0000213d0000000d02000029000000c003200039000000400030043f000002090020009c000000aa0000213d0000000d020000290000016002200039000000400020043f0000006002400039000000000202043300000000002304350000000702100029000000800400043d0000000002420019000001200400043d0000000002420019000000e00400043d00000000024200190000000002020433000001fd0020009c000000b70000213d0000000d04000029000000e0044000390000000000240435000000060e100029000000800200043d0000000002e20019000001200400043d0000000002420019000000e00400043d00000000024200190000000002020433000001fd0020009c000000b70000213d0000000d0400002900000100044000390000000000240435000000050f100029000000800200043d0000000002f20019000001200400043d0000000002420019000000e00400043d00000000024200190000000002020433000001fd0020009c000000b70000213d0000000d04000029000001200440003900000000002404350000000402100029000000800400043d0000000002420019000001200400043d0000000002420019000000e00400043d00000000024200190000000002020433000001fd0020009c000000b70000213d0000000d040000290000000003340436000001400440003900000000002404350000000b01100029000000800200043d0000000002120019000001200400043d0000000004420019000000e00200043d000000000424001900000100054000390000000005050433000001fd0050009c000000b70000213d0000000005540019000001800400043d0000000002420019000000200620003900000202026001970000007f045000390000020207400197000000000827013f000000000027004b00000000020000190000020202004041000000000064004b00000000040000190000020204008041000002020080009c000000000204c019000000000002004b000000b70000c13d00000060025000390000000004020433000001fd0040009c000000aa0000213d0000001f024000390000020c022001970000003f022000390000020c02200197000000400700043d0000000002270019000000000072004b00000000080000390000000108004039000001fd0020009c000000aa0000213d0000000100800190000000aa0000c13d000000400020043f000000000847043600000080055000390000000002540019000000000062004b000000b70000213d000000000004004b000001d90000613d000000000600001900000000028600190000000009560019000000000909043300000000009204350000002006600039000000000046004b000001d20000413d000000000248001900000000000204350000000000730435000000800200043d0000000002120019000001200300043d0000000003320019000000e00200043d000000000323001900000120043000390000000004040433000001fd0040009c000000b70000213d0000000004430019000001800300043d0000000002320019000000200520003900000202025001970000007f034000390000020206300197000000000726013f000000000026004b00000000020000190000020202004041000000000053004b00000000030000190000020203008041000002020070009c000000000203c019000000000002004b000000b70000c13d00000060024000390000000003020433000001fd0030009c000000aa0000213d0000001f023000390000020c022001970000003f022000390000020c02200197000000400600043d0000000002260019000000000062004b00000000070000390000000107004039000001fd0020009c000000aa0000213d0000000100700190000000aa0000c13d000000400020043f000000000736043600000080044000390000000002430019000000000052004b000000b70000213d000000000003004b000002190000613d000000000500001900000000027500190000000008450019000000000808043300000000008204350000002005500039000000000035004b000002120000413d000000000237001900000000000204350000000d0200002900000040022000390000000000620435000000800200043d0000000002210019000001200300043d0000000002320019000000e00300043d000000000232001900000140022000390000000002020433000002050020009c000000b70000213d0000000d0400002900000060034000390000000000230435000000800200043d0000000002210019000001200300043d0000000002320019000000e00300043d00000000023200190000016002200039000000800340003900000000020204330000000000230435000000800200043d0000000002120019000001200300043d0000000003320019000000e00200043d000000000423001900000180034000390000000003030433000001fd0030009c000000b70000213d0000000004340019000001800500043d000000000252001900000020022000390000007f05400039000000000025004b0000000006000019000002020600804100000202022001970000020205500197000000000725013f000000000025004b00000000020000190000020202004041000002020070009c000000000206c019000000000002004b000000b70000c13d00000060024000390000000004020433000001fd0040009c000000aa0000213d00000005054002100000003f025000390000020302200197000000400600043d0000000002260019000c00000006001d000000000062004b00000000060000390000000106004039000001fd0020009c000000aa0000213d0000000100600190000000aa0000c13d000000400020043f0000000c0200002900000000004204350000000007130019000000e00400043d0000000002470019000000800100043d0000000002120019000001200600043d0000000008620019000001800200043d0000000002420019000000000985001900000020022000390000008009900039000000000029004b000000b70000213d0000000002410019000000200170003900000000021200190000000006620019000000600550003900000000025600190000008008800039000000000028004b000001380000813d000e000000f3001d000000000ee3001900000100097000390000000c030000290000029a0000013d000000200330003900000000024a0019000000000002043500000060026000390000000000c204350000000002b90019000000800400043d0000000002420019000001200400043d0000000002420019000000e00400043d00000000024200190000000002020433000000800460003900000000002404350000000000630435000000800200043d0000000002120019000001200400043d0000000002420019000000e00400043d00000000064200190000000002560019000000000028004b000001380000813d000000008b080434000001fd00b0009c000000b70000213d000001800200043d0000000007420019000000000a6b00190000000002a70049000000400220008a000002010020009c000000b70000213d000000a00020008c000000b70000413d000000400600043d000002000060009c000000aa0000213d000000a004600039000000400040043f0000006002a000390000000002020433000001fd0020009c000000b70000213d000000200aa00039000000000a2a0019000000200c7000390000005f02a000390000020207200197000002020dc00197000000000fd7013f0000000000d7004b000000000700001900000202070040410000000000c2004b000000000200001900000202020080410000020200f0009c000000000702c019000000000007004b000000b70000c13d0000004002a000390000000007020433000001fd0070009c000000aa0000213d0000001f027000390000020c022001970000003f022000390000020c022001970000000002420019000001fd0020009c000000aa0000213d000000400020043f0000000000740435000000600aa000390000000002a700190000000000c2004b000000b70000213d000000c00c600039000000000007004b000002dc0000613d000000000d0000190000000002cd0019000000000fad0019000000000f0f04330000000000f20435000000200dd0003900000000007d004b000002d50000413d0000000002c70019000000000002043500000000044604360000000002be0019000000800700043d0000000002720019000001200700043d0000000002720019000000e00700043d00000000027200190000000002020433000002050020009c000000b70000213d00000000002404350000000e02b00029000000800400043d0000000002420019000001200400043d0000000002420019000000e00400043d00000000024200190000000002020433000001f90020009c000000b70000213d0000004004600039000000000024043500000000021b0019000000800400043d0000000002420019000001200400043d0000000004420019000000e00200043d0000000004240019000000c0074000390000000007070433000001fd0070009c000000b70000213d0000000004470019000001800700043d0000000002720019000000200720003900000202027001970000007f0a400039000002020ca00197000000000d2c013f00000000002c004b0000000002000019000002020200404100000000007a004b000000000a000019000002020a0080410000020200d0009c00000000020ac019000000000002004b000000b70000c13d000000200d4000390000004002d000390000000004020433000001fd0040009c000000aa0000213d0000001f024000390000020c022001970000003f022000390000020c02200197000000400c00043d00000000022c00190000000000c2004b000000000a000039000000010a004039000001fd0020009c000000aa0000213d0000000100a00190000000aa0000c13d000000400020043f000000000a4c0436000000600dd000390000000002d40019000000000072004b000000b70000213d000000000004004b000002810000613d00000000070000190000000002a70019000000000fd70019000000000f0f04330000000000f204350000002007700039000000000047004b0000032e0000413d000002810000013d000000010100002900000002020000290000000000210435000000800100043d000001200200043d0000000001120019000000e00200043d000000000121001900000080031000390000000003030433000001fd0030009c000000b70000213d0000000001310019000001800400043d000000000242001900000020022000390000005f04100039000000000024004b0000000005000019000002020500804100000202022001970000020204400197000000000624013f000000000024004b00000000020000190000020202004041000002020060009c000000000205c019000000000002004b000000b70000c13d00000040011000390000000004010433000001fd0040009c000000aa0000213d00000005014002100000003f021000390000020302200197000000400500043d0000000002250019000d00000005001d000000000052004b00000000050000390000000105004039000001fd0020009c000000aa0000213d0000000100500190000000aa0000c13d000000400020043f0000000d020000290000000000420435000000e00700043d0000000002370019000000800400043d0000000002420019000001200400043d0000000009420019000001800200043d0000000002720019000e00600010003d0000000e019000290000002002200039000000000021004b000000b70000213d0000006005900039000000000015004b000004030000813d0000000d06000029000003850000013d00000020066000390000000000760435000000800100043d0000000001310019000001200200043d0000000001210019000000e00700043d00000000097100190000000e01900029000000000015004b000004030000813d0000000058050434000001fd0080009c000000b70000213d0000000001980019000001800200043d000000000227001900000020022000390000007f07100039000000000027004b0000000009000019000002020900804100000202022001970000020207700197000000000a27013f000000000027004b000000000200001900000202020040410000020200a0009c000000000209c019000000000002004b000000b70000c13d00000060011000390000000001010433000001fd0010009c000000aa0000213d00000005091002100000003f029000390000020302200197000000400700043d0000000002270019000000000072004b000000000a000039000000010a004039000001fd0020009c000000aa0000213d0000000100a00190000000aa0000c13d000000400020043f00000000001704350000000008380019000000e00c00043d00000000018c0019000000800200043d0000000001210019000001200200043d000000000d210019000001800100043d0000000002c10019000000800990003900000000019d00190000002002200039000000000021004b000000b70000213d000000800ad0003900000000001a004b0000037a0000813d000000000b070019000003cc0000013d000000200bb000390000000001c1001900000000000104350000000000fb0435000000800100043d0000000001810019000001200200043d0000000001210019000000e00c00043d000000000dc1001900000000019d001900000000001a004b0000037a0000813d00000000a10a0434000001fd0010009c000000b70000213d000000000dd10019000001800100043d00000000011c0019000000200e1000390000020201e001970000009f02d00039000002020c200197000000000f1c013f00000000001c004b000000000100001900000202010040410000000000e2004b000000000200001900000202020080410000020200f0009c000000000102c019000000000001004b000000b70000c13d0000008001d00039000000000c010433000001fd00c0009c000000aa0000213d0000001f01c000390000020c011001970000003f011000390000020c01100197000000400f00043d00000000011f00190000000000f1004b00000000020000390000000102004039000001fd0010009c000000aa0000213d0000000100200190000000aa0000c13d000000400010043f0000000001cf0436000000a00dd000390000000002dc00190000000000e2004b000000b70000213d00000000000c004b000003bf0000613d000000000e00001900000000021e00190000000004de001900000000040404330000000000420435000000200ee000390000000000ce004b000003fb0000413d000003bf0000013d000000030100002900000040011000390000000d020000290000000000210435000000800100043d000001200200043d0000000001120019000000e00200043d0000000003210019000000a0013000390000000001010433000001fd0010009c000000b70000213d0000000003130019000001800400043d000000000242001900000020022000390000005f04300039000000000024004b0000000005000019000002020500804100000202022001970000020204400197000000000624013f000000000024004b00000000020000190000020202004041000002020060009c000000000205c019000000000002004b000000b70000c13d00000040023000390000000005020433000001fd0050009c000000aa0000213d00000005045002100000003f024000390000020302200197000000400300043d0000000002230019000000000032004b00000000060000390000000106004039000001fd0020009c000000aa0000213d0000000100600190000000aa0000c13d000000400020043f0000000000530435000000e00200043d000000800500043d0000000005250019000001800600043d00000000072600190000000006150019000001200200043d0000000006260019000000600440003900000000086400190000002007700039000000000078004b000000b70000213d0000000001140019000000000415001900000000022400190000006004600039000000000024004b000000b90000813d0000000005030019000000200550003900000000420404340000000000250435000000800200043d0000000002120019000001200600043d0000000002620019000000e00600043d0000000002620019000000000024004b000004480000413d000000b90000013d0000002004600039000000000441034f0000020c053001980000001f0630018f000001c001500039000004600000613d000001c007000039000000000804034f000000008908043c0000000007970436000000000017004b0000045c0000c13d000000000006004b0000046d0000613d000000000454034f0000000305600210000000000601043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f0000000000410435000001c0013000390000000000010435000000400300043d000001ff0030009c000000aa0000213d0000006004300039000000400040043f000002000030009c000000aa0000213d000000a001300039000000400010043f000000600100003900000000001404350000004005300039000000000015043500000020053000390000000000150435000000000043043500000080033000390000000000130435000001a00900043d000002010090009c000000b70000213d000000200090008c000000b70000413d000001c00500043d000001fd0050009c000000b70000213d000001a004900039000001a0065000390000000003640049000002010030009c000000b70000213d000000600030008c000000b70000413d000000400300043d000001ff0030009c000000aa0000213d0000006008300039000000400080043f000001c0075000390000000007070433000001fd0070009c000000b70000213d000000000a6700190000000004a40049000002010040009c000000b70000213d000000400040008c000000b70000413d000002000030009c000000aa0000213d000000a004300039000e00000004001d000000400040043f0000002004a000390000000007040433000001fd0070009c000000b70000213d000001c004900039000000000ca700190000003f07c00039000000000047004b000000000b000019000002020b0080410000020207700197000d02020040019b0000000d0d70014f0000000d0070006c000000000700001900000202070040410000020200d0009c00000000070bc019000000000007004b000000b70000c13d0000002007c00039000000000e070433000001fd00e0009c000000aa0000213d0000000507e002100000003f0770003900000203077001970000000e07700029000001fd0070009c000000aa0000213d000000400070043f0000000e070000290000000000e704350000004007c00039000000060be00210000000000d7b001900000000004d004b000000b70000213d00000000000e004b0000078f0000c13d0000000e0700002900000000007804350000004007a000390000000007070433000001fd0070009c000000b70000213d000000000aa700190000003f07a00039000000000047004b000000000b000019000002020b00804100000202077001970000000d0c70014f0000000d0070006c000000000700001900000202070040410000020200c0009c00000000070bc019000000000007004b000000b70000c13d0000002007a00039000000000c070433000001fd00c0009c000000aa0000213d0000000507c002100000003f077000390000020307700197000000400b00043d00000000077b00190000000000b7004b000000000d000039000000010d004039000001fd0070009c000000aa0000213d0000000100d00190000000aa0000c13d000000400070043f0000000000cb04350000004007a00039000000060cc00210000000000c7c001900000000004c004b000000b70000213d0000000000c7004b000005140000813d000000000d0b0019000000000e0700190000000007740049000002010070009c000000b70000213d000000400070008c000000b70000413d000000400700043d000002040070009c000000aa0000213d000000400f7000390000004000f0043f000000000f0e0433000001fd00f0009c000000b70000213d000000000ff70436000000600aa00039000000000a0a04330000020600a0009c000000b70000213d000000200dd000390000000000af043500000000007d04350000004007e000390000000000c7004b000000000a0e0019000004fa0000413d0000000007830436000700000007001d00000080073000390000000000b70435000001e0075000390000000007070433000001fd0070009c000000b70000213d0000000007670019000a00000007001d0000003f07700039000000000047004b0000000008000019000002020800804100000202077001970000000d0a70014f0000000d0070006c000000000700001900000202070040410000020200a0009c000000000708c019000000000007004b000000b70000c13d0000000a0700002900000020077000390000000007070433000001fd0070009c000000aa0000213d00000005087002100000003f0a800039000002030aa00197000000400b00043d000000000aab001900060000000b001d0000000000ba004b000000000b000039000000010b004039000001fd00a0009c000000aa0000213d0000000100b00190000000aa0000c13d0000004000a0043f000000060a00002900000000007a04350000000a07000029000e00400070003d0000000e07800029000900000007001d000000000047004b000000b70000213d00000009080000290000000e0080006b000005af0000813d000801800090003d000c00060000002d0000000e080000290000000087080434000e00000008001d000001fd0070009c000000b70000213d0000000a0f7000290000000807f00069000002010070009c000000b70000213d000000a00070008c000000b70000413d000000400800043d000002000080009c000000aa0000213d000000a007800039000000400070043f0000004007f000390000000007070433000001fd0070009c000000b70000213d0000000007780436000b00000007001d0000006007f000390000000007070433000001fd0070009c000000b70000213d0000000009f700190000005f07900039000000000047004b000000000a000019000002020a00804100000202077001970000000d0b70014f0000000d0070006c000000000700001900000202070040410000020200b0009c00000000070ac019000000000007004b000000b70000c13d0000004007900039000000000a070433000001fd00a0009c000000aa0000213d0000001f07a000390000020c077001970000003f077000390000020c07700197000000400d00043d00000000077d00190000000000d7004b000000000b000039000000010b004039000001fd0070009c000000aa0000213d0000000100b00190000000aa0000c13d000000400070043f0000000007ad04360000006009900039000000000b9a001900000000004b004b000000b70000213d00000000000a004b000005940000613d000000000b000019000000000c7b0019000000000e9b0019000000000e0e04330000000000ec0435000000200bb000390000000000ab004b0000058d0000413d0000000007a7001900000000000704350000000b070000290000000000d704350000008007f000390000000007070433000001fd0070009c000000b70000213d00000040098000390000000000790435000000a007f000390000000007070433000001fd0070009c000000b70000213d0000000c0a000029000000200aa0003900000060098000390000000000790435000000c007f00039000000000707043300000080098000390000000000790435000c0000000a001d00000000008a043500000009080000290000000e0080006b0000054b0000413d00000007070000290000000608000029000000000087043500000200055000390000000005050433000001fd0050009c000000b70000213d00000000056500190000003f06500039000000000046004b0000000007000019000002020700804100000202066001970000000d0860014f0000000d0060006c00000000060000190000020206004041000002020080009c000000000607c019000000000006004b000000b70000c13d00000020065000390000000007060433000001fd0070009c000000aa0000213d00000005067002100000003f066000390000020308600197000000400600043d0000000008860019000000000068004b00000000090000390000000109004039000001fd0080009c000000aa0000213d0000000100900190000000aa0000c13d000000400080043f0000000000760435000000060770021000000040055000390000000007750019000000000047004b000000b70000213d000000000075004b000005f10000813d00000000080600190000000009540049000002010090009c000000b70000213d000000400090008c000000b70000413d000000400900043d000002040090009c000000aa0000213d0000002008800039000000400a9000390000004000a0043f00000000ba050434000000000aa90436000000000b0b04330000000000ba043500000000009804350000004005500039000000000075004b000005de0000413d0000004004300039000d00000004001d00000000006404350000002006000039000000400400043d0000000006640436000000000303043300000000001604350000008007400039000000400800003900000000160304340000000000870435000000c00340003900000000070604330000000000730435000e00000004001d000000e003400039000000000007004b000006110000613d00000000080000190000002006600039000000000906043300000000a909043400000205099001970000000009930436000000000a0a0433000002060aa001970000000000a9043500000040033000390000000108800039000000000078004b000006050000413d0000000e040000290000000006430049000000800760008a0000000006010433000000a001400039000000000071043500000000070604330000000001730436000000000007004b000006280000613d0000000003000019000000200660003900000000080604330000000098080434000001fd08800197000000000881043600000000090904330000020609900197000000000098043500000040011000390000000103300039000000000073004b0000061c0000413d0000000e040000290000000003410049000000200630008a000000070300002900000000030304330000004007400039000000000067043500000000070304330000000000710435000000050670021000000000066100190000002006600039000000000007004b000007aa0000c13d0000000e030000290000000001360049000000200110008a0000000d0200002900000000020204330000006003300039000000000013043500000000030204330000000001360436000000000003004b0000064c0000613d000000000500001900000020022000390000000004020433000000006404043400000000044104360000000006060433000000000064043500000040011000390000000105500039000000000035004b000006420000413d0000000e020000290000065c0000013d000000400300043d00000020010000390000000001130436000001600200043d00000000020204330000000000210435000500000003001d000800400030003d00000005012002100000000801100029000400000002001d000000000002004b000006650000c13d00000005020000290000000001210049000001f90010009c000001f9010080410000006001100210000001f90020009c000001f9020080410000004002200210000000000121019f000007df0001042e000001400600043d00000000040000190000000003010019000006740000013d00000007060000290000008002700039000000090300002900000080033000390000000003030433000000000032043500000006040000290000000104400039000000040040006c00000000030100190000065b0000813d000600000004001d000000050130006a000000400110008a00000008020000290000000002120436000800000002001d0000000061060434000700000006001d000900000001001d0000000021010434000001fd0110019700000000011304360000000004020433000000a0020000390000000000210435000000a00130003900000000020404330000000000210435000000c0053000390000000501200210000000000b510019000a00000002001d000000000002004b000b00000003001d0000073a0000613d0000000002000019000006960000013d0000000c0200002900000001022000390000000a0020006c0000000b030000290000000e040000290000000d050000290000073a0000813d000c00000002001d00000000013b0049000000c00110008a0000000005150436000d00000005001d0000002004400039000e00000004001d00000000010404330000000032010434000000005402043400000000044b04360000000005050433000001fd05500197000000000054043500000040042000390000000004040433000001fd044001970000004005b00039000000000045043500000060042000390000000004040433000001fd044001970000006005b00039000000000045043500000080022000390000000002020433000001fd022001970000008004b0003900000000002404350000000002030433000000a003b00039000001400400003900000000004304350000014004b00039000000005302043400000000003404350000016004b00039000000000003004b000006c50000613d000000000600001900000000024600190000000008650019000000000808043300000000008204350000002006600039000000000036004b000006be0000413d000000000243001900000000000204350000001f023000390000020c022001970000000002420019000000400310003900000000030304330000000004b20049000000c005b00039000000000045043500000000540304340000000003420436000000000004004b000006db0000613d000000000600001900000000023600190000000008650019000000000808043300000000008204350000002006600039000000000046004b000006d40000413d00000000023400190000000000020435000000600210003900000000020204330000020502200197000000e005b000390000000000250435000000800210003900000000020204330000010005b0003900000000002504350000001f024000390000020c02200197000000000a3200190000000002ba00490000012004b00039000000a00110003900000000030104330000000000240435000000000403043300000000004a0435000000050140021000000000011a0019000000200b100039000000000004004b0000068f0000613d000000000500001900000000010a0019000007040000013d000000000268001900000000000204350000008002b000390000008007900039000000000707043300000000007204350000001f028000390000020c02200197000000000b6200190000000105500039000000000045004b0000068f0000813d0000000002ab0049000000200220008a000000200110003900000000002104350000002003300039000000000903043300000000e2090434000000a006000039000000000d6b0436000000a006b0003900000000c80204340000000000860435000000c006b00039000000000008004b0000071b0000613d00000000020000190000000007620019000000000f2c0019000000000f0f04330000000000f704350000002002200039000000000082004b000007140000413d0000000002680019000000000002043500000000020e0433000002050220019700000000002d043500000040029000390000000002020433000001f9022001970000004007b0003900000000002704350000001f028000390000020c0220019700000000026200190000000006b200490000006007b0003900000060089000390000000008080433000000000067043500000000c80804340000000006820436000000000008004b000006f80000613d00000000020000190000000007620019000000000d2c0019000000000d0d04330000000000d704350000002002200039000000000082004b000007320000413d000006f80000013d00000009010000290000004001100039000000000101043300000000023b004900000040033000390000000000230435000000000301043300000000003b0435000000050230021000000000022b00190000002005200039000e00000003001d000000000003004b0000077a0000613d0000000004000019000000000a0b00190000074f0000013d00000001044000390000000e0040006c00000000050d00190000077b0000813d0000000002b50049000000200220008a000000200aa0003900000000002a043500000020011000390000000009010433000000000c0904330000000000c504350000000502c002100000000002250019000000200d20003900000000000c004b0000074b0000613d000000000e000019000000000f050019000007670000013d000000000268001900000000000204350000001f028000390000020c02200197000000000d620019000000010ee000390000000000ce004b0000074b0000813d00000000025d0049000000200220008a000000200ff0003900000000002f043500000020099000390000000002090433000000003802043400000000068d0436000000000008004b0000075f0000613d00000000020000190000000007620019000000000d230019000000000d0d04330000000000d704350000002002200039000000000082004b000007720000413d0000075f0000013d000000000d0500190000000901000029000000600110003900000000030104330000000b0700002900000000017d004900000060027000390000000000120435000000000403043300000000014d0436000000000004004b000006690000613d000000000500001900000007060000290000002003300039000000000203043300000000012104360000000105500039000000000045004b000007880000413d0000066a0000013d000000c00e300039000000000f0700190000000007740049000002010070009c000000b70000213d000000400070008c000000b70000413d000000400700043d000002040070009c000000aa0000213d000000400b7000390000004000b0043f000000000b0f04330000020500b0009c000000b70000213d000000000bb70436000000600cc00039000000000c0c04330000020600c0009c000000b70000213d0000000000cb0435000000000e7e04360000004007f000390000000000d7004b000000000c0f0019000007900000413d000004cc0000013d000000a0080000390000000009000019000000000a010019000007c40000013d0000000004dc001900000000000404350000004004b000390000000004040433000001fd04400197000000400560003900000000004504350000006004b000390000000004040433000001fd044001970000006005600039000000000045043500000080046000390000008005b00039000000000505043300000000005404350000001f04c000390000020c044001970000000006d400190000000109900039000000000079004b000006360000813d000000000b160049000000200bb0008a000000200aa000390000000000ba04350000002003300039000000000b03043300000000dc0b0434000001fd0cc00197000000000cc60436000000000d0d043300000000008c0435000000a00f60003900000000ec0d04340000000000cf0435000000c00d60003900000000000c004b000007ae0000613d000000000f0000190000000005df00190000000004fe001900000000040404330000000000450435000000200ff000390000000000cf004b000007d60000413d000007ae0000013d000007de00000432000007df0001042e000007e00001043000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffff000000020000000000000000000000000000004000000100000000000000000000000000000000000000000000000000000000000000000000000000f816ec60000000000000000000000000000000000000000000000000000000006fb34956000000000000000000000000000000000000000000000000ffffffffffffffff000000000000000000000000000000000000000000000000fffffffffffffe5f000000000000000000000000000000000000000000000000ffffffffffffff9f000000000000000000000000000000000000000000000000ffffffffffffff5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffffbf000000000000000000000000ffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000ffffffffffffff3f000000000000000000000000000000000000000000000000fffffffffffffe9f4e487b71000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0")
