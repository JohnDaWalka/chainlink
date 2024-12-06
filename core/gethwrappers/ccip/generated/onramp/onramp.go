package onramp

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

type ClientEVM2AnyMessage struct {
	Receiver     []byte
	Data         []byte
	TokenAmounts []ClientEVMTokenAmount
	FeeToken     common.Address
	ExtraArgs    []byte
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
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

type OnRampAllowlistConfigArgs struct {
	DestChainSelector         uint64
	AllowlistEnabled          bool
	AddedAllowlistedSenders   []common.Address
	RemovedAllowlistedSenders []common.Address
}

type OnRampDestChainConfigArgs struct {
	DestChainSelector uint64
	Router            common.Address
	AllowlistEnabled  bool
}

type OnRampDynamicConfig struct {
	FeeQuoter              common.Address
	ReentrancyGuardEntered bool
	MessageInterceptor     common.Address
	FeeAggregator          common.Address
	AllowlistAdmin         common.Address
}

type OnRampStaticConfig struct {
	ChainSelector      uint64
	RmnRemote          common.Address
	NonceManager       common.Address
	TokenAdminRegistry common.Address
}

var OnRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"internalType\":\"structOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeAggregator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowlistAdmin\",\"type\":\"address\"}],\"internalType\":\"structOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowlistEnabled\",\"type\":\"bool\"}],\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotSendZeroTokens\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GetSupportedTokensFunctionalityRemovedCheckAdminRegistry\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidAllowListRequest\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidDestChainConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeCalledByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAllowlistAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustSetOriginalSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"UnsupportedToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"allowlistAdmin\",\"type\":\"address\"}],\"name\":\"AllowListAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"AllowListSendersAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"AllowListSendersRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeValueJuels\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPMessageSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOnRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeAggregator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowlistAdmin\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowlistEnabled\",\"type\":\"bool\"}],\"name\":\"DestChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeAggregator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeTokenWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowlistEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"addedAllowlistedSenders\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"removedAllowlistedSenders\",\"type\":\"address[]\"}],\"internalType\":\"structOnRamp.AllowlistConfigArgs[]\",\"name\":\"allowlistConfigArgsItems\",\"type\":\"tuple[]\"}],\"name\":\"applyAllowlistUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowlistEnabled\",\"type\":\"bool\"}],\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyDestChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"}],\"name\":\"forwardFromRouter\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getAllowedSendersList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"configuredAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getDestChainConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"allowlistEnabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeAggregator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowlistAdmin\",\"type\":\"address\"}],\"internalType\":\"structOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"contractIERC20\",\"name\":\"sourceToken\",\"type\":\"address\"}],\"name\":\"getPoolBySourceToken\",\"outputs\":[{\"internalType\":\"contractIPoolV1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"internalType\":\"structOnRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"getSupportedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeAggregator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowlistAdmin\",\"type\":\"address\"}],\"internalType\":\"structOnRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"feeTokens\",\"type\":\"address[]\"}],\"name\":\"withdrawFeeTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200409f3803806200409f833981016040819052620000359162000709565b336000816200005757604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008a576200008a816200014f565b505082516001600160401b03161580620000af575060208301516001600160a01b0316155b80620000c6575060408301516001600160a01b0316155b80620000dd575060608301516001600160a01b0316155b15620000fc576040516306b7c75960e31b815260040160405180910390fd5b82516001600160401b031660805260208301516001600160a01b0390811660a0526040840151811660c05260608401511660e0526200013b82620001c9565b620001468162000378565b5050506200080a565b336001600160a01b038216036200017957604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03161580620001ec575060608101516001600160a01b0316155b80620001f9575080602001515b1562000218576040516306b7c75960e31b815260040160405180910390fd5b8051600280546020808501511515600160a01b026001600160a81b03199092166001600160a01b039485161791909117909155604080840151600380549185166001600160a01b0319928316179055606080860151600480549187169184169190911790556080808701516005805491881691909416179092558251808301845291516001600160401b0316825260a05185169382019390935260c05184168183015260e05190931691830191909152517fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1916200036d91849082516001600160401b031681526020808401516001600160a01b03908116828401526040858101518216818501526060958601518216868501528451821660808086019190915292850151151560a0850152840151811660c084015293830151841660e0830152909101519091166101008201526101200190565b60405180910390a150565b60005b8151811015620004cb5760008282815181106200039c576200039c620007f4565b602002602001015190506000838381518110620003bd57620003bd620007f4565b6020026020010151600001519050806001600160401b0316600003620004055760405163c35aa79d60e01b81526001600160401b038216600482015260240160405180910390fd5b6001600160401b0381811660008181526006602090815260409182902086820151815488850151600160401b600160e81b031990911669010000000000000000006001600160a01b0390931692830260ff60401b19161768010000000000000000911515820217808455855197811688529387019190915260ff920491909116151591840191909152917fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef59060600160405180910390a25050508060010190506200037b565b5050565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b03811182821017156200050a576200050a620004cf565b60405290565b604051601f8201601f191681016001600160401b03811182821017156200053b576200053b620004cf565b604052919050565b80516001600160401b03811681146200055b57600080fd5b919050565b6001600160a01b03811681146200057657600080fd5b50565b805180151581146200055b57600080fd5b600060a082840312156200059d57600080fd5b60405160a081016001600160401b0381118282101715620005c257620005c2620004cf565b80604052508091508251620005d78162000560565b8152620005e76020840162000579565b60208201526040830151620005fc8162000560565b60408201526060830151620006118162000560565b60608201526080830151620006268162000560565b6080919091015292915050565b600082601f8301126200064557600080fd5b815160206001600160401b03821115620006635762000663620004cf565b62000673818360051b0162000510565b828152606092830285018201928282019190878511156200069357600080fd5b8387015b85811015620006fc5781818a031215620006b15760008081fd5b620006bb620004e5565b620006c68262000543565b815285820151620006d78162000560565b818701526040620006ea83820162000579565b90820152845292840192810162000697565b5090979650505050505050565b60008060008385036101408112156200072157600080fd5b60808112156200073057600080fd5b50604051608081016001600160401b038082118383101715620007575762000757620004cf565b81604052620007668762000543565b8352602087015191506200077a8262000560565b81602084015260408701519150620007928262000560565b81604084015260608701519150620007aa8262000560565b816060840152829550620007c288608089016200058a565b9450610120870151925080831115620007da57600080fd5b5050620007ea8682870162000633565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b60805160a05160c05160e05161381c62000883600039600081816101fc015281816107670152611ac40152600081816101c0015281816114b90152611a9d015260008181610184015281816105a50152611a7301526000818161015401528181611040015281816115d60152611a4f015261381c6000f3fe608060405234801561001057600080fd5b506004361061011b5760003560e01c80637437ff9f116100b2578063972b461211610081578063df0aa9e911610066578063df0aa9e9146104fb578063f2fde38b1461050e578063fbca3b741461052157600080fd5b8063972b4612146104c7578063c9b146b3146104e857600080fd5b80637437ff9f146103cb57806379ba5097146104755780638da5cb5b1461047d5780639041be3d1461049b57600080fd5b806327e936f1116100ee57806327e936f1146102ce57806348a98aa4146102e15780635cb80c5d146103195780636def4ce71461032c57600080fd5b806306285c6914610120578063181f5a771461024f57806320487ded146102985780632716072b146102b9575b600080fd5b61023960408051608081018252600080825260208201819052918101829052606081019190915260405180608001604052807f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16815250905090565b6040516102469190612589565b60405180910390f35b61028b6040518060400160405280601081526020017f4f6e52616d7020312e362e302d6465760000000000000000000000000000000081525081565b604051610246919061264e565b6102ab6102a636600461268f565b610541565b604051908152602001610246565b6102cc6102c73660046127fd565b6106fa565b005b6102cc6102dc3660046128eb565b61070e565b6102f46102ef366004612983565b61071f565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610246565b6102cc610327366004612a08565b6107d4565b61038f61033a366004612a4a565b67ffffffffffffffff9081166000908152600660205260409020549081169168010000000000000000820460ff16916901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1690565b6040805167ffffffffffffffff9094168452911515602084015273ffffffffffffffffffffffffffffffffffffffff1690820152606001610246565b6104686040805160a081018252600080825260208201819052918101829052606081018290526080810191909152506040805160a08101825260025473ffffffffffffffffffffffffffffffffffffffff80821683527401000000000000000000000000000000000000000090910460ff161515602083015260035481169282019290925260045482166060820152600554909116608082015290565b6040516102469190612a67565b6102cc610956565b60015473ffffffffffffffffffffffffffffffffffffffff166102f4565b6104ae6104a9366004612a4a565b610a24565b60405167ffffffffffffffff9091168152602001610246565b6104da6104d5366004612a4a565b610a4d565b604051610246929190612b12565b6102cc6104f6366004612a08565b610a91565b6102ab610509366004612b2d565b610dad565b6102cc61051c366004612b99565b6116bb565b61053461052f366004612a4a565b6116cc565b6040516102469190612bb6565b6040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff00000000000000000000000000000000608084901b16600482015260009073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690632cbc26bb90602401602060405180830381865afa1580156105ec573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106109190612bc9565b15610658576040517ffdbd6a7200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b6002546040517fd8694ccd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063d8694ccd906106b09086908690600401612ce8565b602060405180830381865afa1580156106cd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106f19190612e31565b90505b92915050565b610702611700565b61070b81611753565b50565b610716611700565b61070b816118f6565b6040517fbbe4f6db00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82811660048301526000917f00000000000000000000000000000000000000000000000000000000000000009091169063bbe4f6db90602401602060405180830381865afa1580156107b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106f19190612e4a565b60045473ffffffffffffffffffffffffffffffffffffffff1660005b8281101561095057600084848381811061080c5761080c612e67565b90506020020160208101906108219190612b99565b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015290915060009073ffffffffffffffffffffffffffffffffffffffff8316906370a0823190602401602060405180830381865afa158015610891573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108b59190612e31565b90508015610946576108de73ffffffffffffffffffffffffffffffffffffffff83168583611b26565b8173ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e8360405161093d91815260200190565b60405180910390a35b50506001016107f0565b50505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146109a7576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b67ffffffffffffffff80821660009081526006602052604081205490916106f491166001612ec5565b67ffffffffffffffff8116600090815260066020526040812080546060916801000000000000000090910460ff1690610a8890600101611bb3565b91509150915091565b60015473ffffffffffffffffffffffffffffffffffffffff163314610b015760055473ffffffffffffffffffffffffffffffffffffffff163314610b01576040517f905d7d9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610da8576000838383818110610b2057610b20612e67565b9050602002810190610b329190612ee6565b610b3b90612f97565b805167ffffffffffffffff1660009081526006602090815260409182902090830151815490151568010000000000000000027fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff90911617815590820151519192509015610cfb57816020015115610cba5760005b826040015151811015610c6a57600083604001518281518110610bd457610bd4612e67565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610c535783516040517f463258ff00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161064f565b610c606001840182611bc7565b5050600101610baf565b50816000015167ffffffffffffffff167f330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc32818360400151604051610cad9190612bb6565b60405180910390a2610cfb565b81516040517f463258ff00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161064f565b60005b826060015151811015610d4757610d3e83606001518281518110610d2457610d24612e67565b602002602001015183600101611be990919063ffffffff16565b50600101610cfe565b5060608201515115610d9e57816000015167ffffffffffffffff167fc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d4215868360600151604051610d959190612bb6565b60405180910390a25b5050600101610b04565b505050565b60025460009074010000000000000000000000000000000000000000900460ff1615610e05576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000017905567ffffffffffffffff8516600090815260066020526040902073ffffffffffffffffffffffffffffffffffffffff8316610eaa576040517fa4ec747900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805468010000000000000000900460ff1615610f1b57610ecd6001820184611c0b565b610f1b576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416600482015260240161064f565b80546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314610f78576040517f1c0a352900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60035473ffffffffffffffffffffffffffffffffffffffff16801561101e576040517fe0a0e50600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82169063e0a0e50690610feb908a908a90600401612ce8565b600060405180830381600087803b15801561100557600080fd5b505af1158015611019573d6000803e3d6000fd5b505050505b50604080516101c081019091526000610120820181815267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811661014085015289811661016085015284549293928392916101808401918791879161108c9116613048565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905567ffffffffffffffff168152602001600067ffffffffffffffff1681525081526020018573ffffffffffffffffffffffffffffffffffffffff168152602001878060200190611100919061306f565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250602001611144888061306f565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250938552505060408051602081810183529381529284019290925250016111a06080890160608a01612b99565b73ffffffffffffffffffffffffffffffffffffffff168152602001868152602001600081526020018780604001906111d891906130d4565b905067ffffffffffffffff8111156111f2576111f26126df565b60405190808252806020026020018201604052801561126b57816020015b6112586040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001606081526020016060815260200160008152602001606081525090565b8152602001906001900390816112105790505b5090529050600061127f60408801886130d4565b808060200260200160405190810160405280939291908181526020016000905b828210156112cb576112bc6040830286013681900381019061313c565b8152602001906001019061129f565b5050505050905060005b6112e260408901896130d4565b905081101561137c5761135282828151811061130057611300612e67565b60209081029190910101518a6113168b8061306f565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c9250611c3a915050565b836101000151828151811061136957611369612e67565b60209081029190910101526001016112d5565b50600254600090606090819073ffffffffffffffffffffffffffffffffffffffff1663430d138c8c6113b360808e018e8601612b99565b8c8e80608001906113c4919061306f565b8b61010001518b6040518863ffffffff1660e01b81526004016113ed9796959493929190613259565b600060405180830381865afa15801561140a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611450919081019061338b565b60e0890193909352909450925090508261152b576040517fea458c0c00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8c16600482015273ffffffffffffffffffffffffffffffffffffffff89811660248301527f0000000000000000000000000000000000000000000000000000000000000000169063ea458c0c906044016020604051808303816000875af1158015611502573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611526919061347e565b61152e565b60005b855167ffffffffffffffff909116608091820152850182905260005b856101000151518110156115a05781818151811061156a5761156a612e67565b6020026020010151866101000151828151811061158957611589612e67565b60209081029190910101516080015260010161154a565b50604080517f130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321602082015267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811692820192909252908c16606082015230608082015261163090869060a00160405160208183030381529060405280519060200120611f51565b85515284516060015160405167ffffffffffffffff918216918d16907f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f329061167990899061349b565b60405180910390a35050600280547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff169055505051519150505b949350505050565b6116c3611700565b61070b816120a3565b60606040517f9e7177c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015473ffffffffffffffffffffffffffffffffffffffff163314611751576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60005b81518110156118f257600082828151811061177357611773612e67565b60200260200101519050600083838151811061179157611791612e67565b60200260200101516000015190508067ffffffffffffffff166000036117ef576040517fc35aa79d00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161064f565b67ffffffffffffffff818116600081815260066020908152604091829020868201518154888501517fffffff000000000000000000000000000000000000000000ffffffffffffffff909116690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9093169283027fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff161768010000000000000000911515820217808455855197811688529387019190915260ff920491909116151591840191909152917fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef59060600160405180910390a2505050806001019050611756565b5050565b805173ffffffffffffffffffffffffffffffffffffffff1615806119325750606081015173ffffffffffffffffffffffffffffffffffffffff16155b8061193e575080602001515b15611975576040517f35be3ac800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160028054602080850151151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090921673ffffffffffffffffffffffffffffffffffffffff9485161791909117909155604080840151600380549185167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316179055606080860151600480549187169184169190911790556080808701516005805491881691909416179092558251918201835267ffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001682527f00000000000000000000000000000000000000000000000000000000000000008516938201939093527f00000000000000000000000000000000000000000000000000000000000000008416818301527f000000000000000000000000000000000000000000000000000000000000000090931691830191909152517fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f191611b1b9184906135f3565b60405180910390a150565b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052610da8908490612167565b60606000611bc083612273565b9392505050565b60006106f18373ffffffffffffffffffffffffffffffffffffffff84166122cf565b60006106f18373ffffffffffffffffffffffffffffffffffffffff841661231e565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156106f1565b611c826040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001606081526020016060815260200160008152602001606081525090565b8460200151600003611cc0576040517f5cf0444900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000611cd085876000015161071f565b905073ffffffffffffffffffffffffffffffffffffffff81161580611da057506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf00000000000000000000000000000000000000000000000000000000600482015273ffffffffffffffffffffffffffffffffffffffff8216906301ffc9a790602401602060405180830381865afa158015611d7a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d9e9190612bc9565b155b15611df25785516040517fbf16aab600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161064f565b60008173ffffffffffffffffffffffffffffffffffffffff16639a4575b96040518060a001604052808881526020018967ffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018a6020015181526020018a6000015173ffffffffffffffffffffffffffffffffffffffff168152506040518263ffffffff1660e01b8152600401611e9191906136a1565b6000604051808303816000875af1158015611eb0573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611ef69190810190613717565b6040805160a08101825273ffffffffffffffffffffffffffffffffffffffff90941684528151602080860191909152918201518482015288820151606085015280519182019052600081526080830152509050949350505050565b60008060001b8284602001518560000151606001518660000151608001518760a001518860c00151604051602001611fcf95949392919073ffffffffffffffffffffffffffffffffffffffff958616815267ffffffffffffffff94851660208201529290931660408301529092166060830152608082015260a00190565b6040516020818303038152906040528051906020012085606001518051906020012086604001518051906020012087610100015160405160200161201391906137a8565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201206080808d0151805190840120928501999099529183019690965260608201949094529485019190915260a084015260c083015260e08201526101000160405160208183030381529060405280519060200120905092915050565b3373ffffffffffffffffffffffffffffffffffffffff8216036120f2576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006121c9826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166124189092919063ffffffff16565b805190915015610da857808060200190518101906121e79190612bc9565b610da8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f74207375636365656400000000000000000000000000000000000000000000606482015260840161064f565b6060816000018054806020026020016040519081016040528092919081815260200182805480156122c357602002820191906000526020600020905b8154815260200190600101908083116122af575b50505050509050919050565b6000818152600183016020526040812054612316575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556106f4565b5060006106f4565b600081815260018301602052604081205480156124075760006123426001836137bb565b8554909150600090612356906001906137bb565b90508082146123bb57600086600001828154811061237657612376612e67565b906000526020600020015490508087600001848154811061239957612399612e67565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806123cc576123cc6137ce565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506106f4565b60009150506106f4565b5092915050565b60606116b38484600085856000808673ffffffffffffffffffffffffffffffffffffffff16858760405161244c91906137fd565b60006040518083038185875af1925050503d8060008114612489576040519150601f19603f3d011682016040523d82523d6000602084013e61248e565b606091505b509150915061249f878383876124aa565b979650505050505050565b606083156125405782516000036125395773ffffffffffffffffffffffffffffffffffffffff85163b612539576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161064f565b50816116b3565b6116b383838151156125555781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064f919061264e565b608081016106f4828467ffffffffffffffff8151168252602081015173ffffffffffffffffffffffffffffffffffffffff808216602085015280604084015116604085015280606084015116606085015250505050565b60005b838110156125fb5781810151838201526020016125e3565b50506000910152565b6000815180845261261c8160208601602086016125e0565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006106f16020830184612604565b67ffffffffffffffff8116811461070b57600080fd5b600060a0828403121561268957600080fd5b50919050565b600080604083850312156126a257600080fd5b82356126ad81612661565b9150602083013567ffffffffffffffff8111156126c957600080fd5b6126d585828601612677565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff81118282101715612731576127316126df565b60405290565b6040805190810167ffffffffffffffff81118282101715612731576127316126df565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156127a1576127a16126df565b604052919050565b600067ffffffffffffffff8211156127c3576127c36126df565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff8116811461070b57600080fd5b801515811461070b57600080fd5b6000602080838503121561281057600080fd5b823567ffffffffffffffff81111561282757600080fd5b8301601f8101851361283857600080fd5b803561284b612846826127a9565b61275a565b8181526060918202830184019184820191908884111561286a57600080fd5b938501935b838510156128cf5780858a0312156128875760008081fd5b61288f61270e565b853561289a81612661565b8152858701356128a9816127cd565b818801526040868101356128bc816127ef565b908201528352938401939185019161286f565b50979650505050505050565b80356128e6816127cd565b919050565b600060a082840312156128fd57600080fd5b60405160a0810181811067ffffffffffffffff82111715612920576129206126df565b604052823561292e816127cd565b8152602083013561293e816127ef565b60208201526040830135612951816127cd565b60408201526060830135612964816127cd565b60608201526080830135612977816127cd565b60808201529392505050565b6000806040838503121561299657600080fd5b82356129a181612661565b915060208301356129b1816127cd565b809150509250929050565b60008083601f8401126129ce57600080fd5b50813567ffffffffffffffff8111156129e657600080fd5b6020830191508360208260051b8501011115612a0157600080fd5b9250929050565b60008060208385031215612a1b57600080fd5b823567ffffffffffffffff811115612a3257600080fd5b612a3e858286016129bc565b90969095509350505050565b600060208284031215612a5c57600080fd5b8135611bc081612661565b60a081016106f4828473ffffffffffffffffffffffffffffffffffffffff808251168352602082015115156020840152806040830151166040840152806060830151166060840152806080830151166080840152505050565b60008151808452602080850194506020840160005b83811015612b0757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101612ad5565b509495945050505050565b82151581526040602082015260006116b36040830184612ac0565b60008060008060808587031215612b4357600080fd5b8435612b4e81612661565b9350602085013567ffffffffffffffff811115612b6a57600080fd5b612b7687828801612677565b935050604085013591506060850135612b8e816127cd565b939692955090935050565b600060208284031215612bab57600080fd5b8135611bc0816127cd565b6020815260006106f16020830184612ac0565b600060208284031215612bdb57600080fd5b8151611bc0816127ef565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112612c1b57600080fd5b830160208101925035905067ffffffffffffffff811115612c3b57600080fd5b803603821315612a0157600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b8183526000602080850194508260005b85811015612b07578135612cb6816127cd565b73ffffffffffffffffffffffffffffffffffffffff168752818301358388015260409687019690910190600101612ca3565b600067ffffffffffffffff808516835260406020840152612d098485612be6565b60a06040860152612d1e60e086018284612c4a565b915050612d2e6020860186612be6565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc080878503016060880152612d64848385612c4a565b9350604088013592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1883603018312612d9d57600080fd5b60209288019283019235915084821115612db657600080fd5b8160061b3603831315612dc857600080fd5b80878503016080880152612ddd848385612c93565b9450612deb606089016128db565b73ffffffffffffffffffffffffffffffffffffffff811660a08901529350612e166080890189612be6565b94509250808786030160c0880152505061249f838383612c4a565b600060208284031215612e4357600080fd5b5051919050565b600060208284031215612e5c57600080fd5b8151611bc0816127cd565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff81811683821601908082111561241157612411612e96565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81833603018112612f1a57600080fd5b9190910192915050565b600082601f830112612f3557600080fd5b81356020612f45612846836127a9565b8083825260208201915060208460051b870101935086841115612f6757600080fd5b602086015b84811015612f8c578035612f7f816127cd565b8352918301918301612f6c565b509695505050505050565b600060808236031215612fa957600080fd5b6040516080810167ffffffffffffffff8282108183111715612fcd57612fcd6126df565b8160405284359150612fde82612661565b908252602084013590612ff0826127ef565b816020840152604085013591508082111561300a57600080fd5b61301636838701612f24565b6040840152606085013591508082111561302f57600080fd5b5061303c36828601612f24565b60608301525092915050565b600067ffffffffffffffff80831681810361306557613065612e96565b6001019392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126130a457600080fd5b83018035915067ffffffffffffffff8211156130bf57600080fd5b602001915036819003821315612a0157600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261310957600080fd5b83018035915067ffffffffffffffff82111561312457600080fd5b6020019150600681901b3603821315612a0157600080fd5b60006040828403121561314e57600080fd5b613156612737565b8235613161816127cd565b81526020928301359281019290925250919050565b600082825180855260208086019550808260051b84010181860160005b8481101561324c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a073ffffffffffffffffffffffffffffffffffffffff82511685528582015181878701526131f582870182612604565b9150506040808301518683038288015261320f8382612604565b925050506060808301518187015250608080830151925085820381870152506132388183612604565b9a86019a9450505090830190600101613193565b5090979650505050505050565b67ffffffffffffffff881681526000602073ffffffffffffffffffffffffffffffffffffffff808a1682850152604089604086015260c060608601526132a360c08601898b612c4a565b85810360808701526132b58189613176565b86810360a0880152875180825285890192509085019060005b818110156132f55783518051871684528701518784015292860192918401916001016132ce565b50909e9d5050505050505050505050505050565b600082601f83011261331a57600080fd5b815167ffffffffffffffff811115613334576133346126df565b61336560207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161275a565b81815284602083860101111561337a57600080fd5b6116b38260208301602087016125e0565b600080600080608085870312156133a157600080fd5b845193506020808601516133b4816127ef565b604087015190945067ffffffffffffffff808211156133d257600080fd5b6133de89838a01613309565b945060608801519150808211156133f457600080fd5b818801915088601f83011261340857600080fd5b8151613416612846826127a9565b81815260059190911b8301840190848101908b83111561343557600080fd5b8585015b8381101561346d578051858111156134515760008081fd5b61345f8e89838a0101613309565b845250918601918601613439565b50989b979a50959850505050505050565b60006020828403121561349057600080fd5b8151611bc081612661565b602081526134ec60208201835180518252602081015167ffffffffffffffff808216602085015280604084015116604085015280606084015116606085015280608084015116608085015250505050565b6000602083015161351560c084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060408301516101a08060e08501526135326101c0850183612604565b915060608501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06101008187860301818801526135708584612604565b945060808801519250818786030161012088015261358e8584612604565b945060a088015192506135ba61014088018473ffffffffffffffffffffffffffffffffffffffff169052565b60c088015161016088015260e08801516101808801528701518685039091018387015290506135e98382613176565b9695505050505050565b610120810161364b828567ffffffffffffffff8151168252602081015173ffffffffffffffffffffffffffffffffffffffff808216602085015280604084015116604085015280606084015116606085015250505050565b825173ffffffffffffffffffffffffffffffffffffffff9081166080848101919091526020850151151560a08501526040850151821660c08501526060850151821660e085015284015116610100830152611bc0565b602081526000825160a060208401526136bd60c0840182612604565b905067ffffffffffffffff6020850151166040840152604084015173ffffffffffffffffffffffffffffffffffffffff8082166060860152606086015160808601528060808701511660a086015250508091505092915050565b60006020828403121561372957600080fd5b815167ffffffffffffffff8082111561374157600080fd5b908301906040828603121561375557600080fd5b61375d612737565b82518281111561376c57600080fd5b61377887828601613309565b82525060208301518281111561378d57600080fd5b61379987828601613309565b60208301525095945050505050565b6020815260006106f16020830184613176565b818103818111156106f4576106f4612e96565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b60008251612f1a8184602087016125e056fea164736f6c6343000818000a",
}

var OnRampABI = OnRampMetaData.ABI

var OnRampBin = OnRampMetaData.Bin

func DeployOnRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OnRampStaticConfig, dynamicConfig OnRampDynamicConfig, destChainConfigArgs []OnRampDestChainConfigArgs) (common.Address, *generated.Transaction, *OnRamp, error) {
	parsed, err := OnRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(OnRampZKBin), backend, staticConfig, dynamicConfig, destChainConfigArgs)
		contractReturn := &OnRamp{address: address, abi: *parsed, OnRampCaller: OnRampCaller{contract: contractBind}, OnRampTransactor: OnRampTransactor{contract: contractBind}, OnRampFilterer: OnRampFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OnRampBin), backend, staticConfig, dynamicConfig, destChainConfigArgs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, HashZks: tx.Hash()}, &OnRamp{address: address, abi: *parsed, OnRampCaller: OnRampCaller{contract: contract}, OnRampTransactor: OnRampTransactor{contract: contract}, OnRampFilterer: OnRampFilterer{contract: contract}}, nil
}

type OnRamp struct {
	address common.Address
	abi     abi.ABI
	OnRampCaller
	OnRampTransactor
	OnRampFilterer
}

type OnRampCaller struct {
	contract *bind.BoundContract
}

type OnRampTransactor struct {
	contract *bind.BoundContract
}

type OnRampFilterer struct {
	contract *bind.BoundContract
}

type OnRampSession struct {
	Contract     *OnRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OnRampCallerSession struct {
	Contract *OnRampCaller
	CallOpts bind.CallOpts
}

type OnRampTransactorSession struct {
	Contract     *OnRampTransactor
	TransactOpts bind.TransactOpts
}

type OnRampRaw struct {
	Contract *OnRamp
}

type OnRampCallerRaw struct {
	Contract *OnRampCaller
}

type OnRampTransactorRaw struct {
	Contract *OnRampTransactor
}

func NewOnRamp(address common.Address, backend bind.ContractBackend) (*OnRamp, error) {
	abi, err := abi.JSON(strings.NewReader(OnRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOnRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnRamp{address: address, abi: abi, OnRampCaller: OnRampCaller{contract: contract}, OnRampTransactor: OnRampTransactor{contract: contract}, OnRampFilterer: OnRampFilterer{contract: contract}}, nil
}

func NewOnRampCaller(address common.Address, caller bind.ContractCaller) (*OnRampCaller, error) {
	contract, err := bindOnRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampCaller{contract: contract}, nil
}

func NewOnRampTransactor(address common.Address, transactor bind.ContractTransactor) (*OnRampTransactor, error) {
	contract, err := bindOnRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampTransactor{contract: contract}, nil
}

func NewOnRampFilterer(address common.Address, filterer bind.ContractFilterer) (*OnRampFilterer, error) {
	contract, err := bindOnRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnRampFilterer{contract: contract}, nil
}

func bindOnRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OnRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OnRamp *OnRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRamp.Contract.OnRampCaller.contract.Call(opts, result, method, params...)
}

func (_OnRamp *OnRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRamp.Contract.OnRampTransactor.contract.Transfer(opts)
}

func (_OnRamp *OnRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRamp.Contract.OnRampTransactor.contract.Transact(opts, method, params...)
}

func (_OnRamp *OnRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_OnRamp *OnRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRamp.Contract.contract.Transfer(opts)
}

func (_OnRamp *OnRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRamp.Contract.contract.Transact(opts, method, params...)
}

func (_OnRamp *OnRampCaller) GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

	error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getAllowedSendersList", destChainSelector)

	outstruct := new(GetAllowedSendersList)
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsEnabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfiguredAddresses = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_OnRamp *OnRampSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRamp.Contract.GetAllowedSendersList(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCallerSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRamp.Contract.GetAllowedSendersList(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

	error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	outstruct := new(GetDestChainConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequenceNumber = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.AllowlistEnabled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Router = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_OnRamp *OnRampSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRamp.Contract.GetDestChainConfig(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCallerSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRamp.Contract.GetDestChainConfig(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCaller) GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OnRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampDynamicConfig)).(*OnRampDynamicConfig)

	return out0, err

}

func (_OnRamp *OnRampSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRamp.Contract.GetDynamicConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRamp.Contract.GetDynamicConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OnRamp *OnRampSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRamp.Contract.GetExpectedNextSequenceNumber(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRamp.Contract.GetExpectedNextSequenceNumber(&_OnRamp.CallOpts, destChainSelector)
}

func (_OnRamp *OnRampCaller) GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OnRamp *OnRampSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRamp.Contract.GetFee(&_OnRamp.CallOpts, destChainSelector, message)
}

func (_OnRamp *OnRampCallerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRamp.Contract.GetFee(&_OnRamp.CallOpts, destChainSelector, message)
}

func (_OnRamp *OnRampCaller) GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getPoolBySourceToken", arg0, sourceToken)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRamp *OnRampSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRamp.Contract.GetPoolBySourceToken(&_OnRamp.CallOpts, arg0, sourceToken)
}

func (_OnRamp *OnRampCallerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRamp.Contract.GetPoolBySourceToken(&_OnRamp.CallOpts, arg0, sourceToken)
}

func (_OnRamp *OnRampCaller) GetStaticConfig(opts *bind.CallOpts) (OnRampStaticConfig, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OnRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampStaticConfig)).(*OnRampStaticConfig)

	return out0, err

}

func (_OnRamp *OnRampSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRamp.Contract.GetStaticConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRamp.Contract.GetStaticConfig(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCaller) GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "getSupportedTokens", arg0)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OnRamp *OnRampSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRamp.Contract.GetSupportedTokens(&_OnRamp.CallOpts, arg0)
}

func (_OnRamp *OnRampCallerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRamp.Contract.GetSupportedTokens(&_OnRamp.CallOpts, arg0)
}

func (_OnRamp *OnRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRamp *OnRampSession) Owner() (common.Address, error) {
	return _OnRamp.Contract.Owner(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) Owner() (common.Address, error) {
	return _OnRamp.Contract.Owner(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OnRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OnRamp *OnRampSession) TypeAndVersion() (string, error) {
	return _OnRamp.Contract.TypeAndVersion(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampCallerSession) TypeAndVersion() (string, error) {
	return _OnRamp.Contract.TypeAndVersion(&_OnRamp.CallOpts)
}

func (_OnRamp *OnRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "acceptOwnership")
}

func (_OnRamp *OnRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRamp.Contract.AcceptOwnership(&_OnRamp.TransactOpts)
}

func (_OnRamp *OnRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRamp.Contract.AcceptOwnership(&_OnRamp.TransactOpts)
}

func (_OnRamp *OnRampTransactor) ApplyAllowlistUpdates(opts *bind.TransactOpts, allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "applyAllowlistUpdates", allowlistConfigArgsItems)
}

func (_OnRamp *OnRampSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyAllowlistUpdates(&_OnRamp.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRamp *OnRampTransactorSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyAllowlistUpdates(&_OnRamp.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRamp *OnRampTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_OnRamp *OnRampSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyDestChainConfigUpdates(&_OnRamp.TransactOpts, destChainConfigArgs)
}

func (_OnRamp *OnRampTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRamp.Contract.ApplyDestChainConfigUpdates(&_OnRamp.TransactOpts, destChainConfigArgs)
}

func (_OnRamp *OnRampTransactor) ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "forwardFromRouter", destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRamp *OnRampSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.ForwardFromRouter(&_OnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRamp *OnRampTransactorSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.ForwardFromRouter(&_OnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRamp *OnRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OnRamp *OnRampSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRamp.Contract.SetDynamicConfig(&_OnRamp.TransactOpts, dynamicConfig)
}

func (_OnRamp *OnRampTransactorSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRamp.Contract.SetDynamicConfig(&_OnRamp.TransactOpts, dynamicConfig)
}

func (_OnRamp *OnRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_OnRamp *OnRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.TransferOwnership(&_OnRamp.TransactOpts, to)
}

func (_OnRamp *OnRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.TransferOwnership(&_OnRamp.TransactOpts, to)
}

func (_OnRamp *OnRampTransactor) WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRamp.contract.Transact(opts, "withdrawFeeTokens", feeTokens)
}

func (_OnRamp *OnRampSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.WithdrawFeeTokens(&_OnRamp.TransactOpts, feeTokens)
}

func (_OnRamp *OnRampTransactorSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRamp.Contract.WithdrawFeeTokens(&_OnRamp.TransactOpts, feeTokens)
}

type OnRampAllowListAdminSetIterator struct {
	Event *OnRampAllowListAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampAllowListAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampAllowListAdminSet)
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
		it.Event = new(OnRampAllowListAdminSet)
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

func (it *OnRampAllowListAdminSetIterator) Error() error {
	return it.fail
}

func (it *OnRampAllowListAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampAllowListAdminSet struct {
	AllowlistAdmin common.Address
	Raw            types.Log
}

func (_OnRamp *OnRampFilterer) FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampAllowListAdminSetIterator, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return &OnRampAllowListAdminSetIterator{contract: _OnRamp.contract, event: "AllowListAdminSet", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampAllowListAdminSet)
				if err := _OnRamp.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseAllowListAdminSet(log types.Log) (*OnRampAllowListAdminSet, error) {
	event := new(OnRampAllowListAdminSet)
	if err := _OnRamp.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampAllowListSendersAddedIterator struct {
	Event *OnRampAllowListSendersAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampAllowListSendersAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampAllowListSendersAdded)
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
		it.Event = new(OnRampAllowListSendersAdded)
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

func (it *OnRampAllowListSendersAddedIterator) Error() error {
	return it.fail
}

func (it *OnRampAllowListSendersAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampAllowListSendersAdded struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampAllowListSendersAddedIterator{contract: _OnRamp.contract, event: "AllowListSendersAdded", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampAllowListSendersAdded)
				if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseAllowListSendersAdded(log types.Log) (*OnRampAllowListSendersAdded, error) {
	event := new(OnRampAllowListSendersAdded)
	if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampAllowListSendersRemovedIterator struct {
	Event *OnRampAllowListSendersRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampAllowListSendersRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampAllowListSendersRemoved)
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
		it.Event = new(OnRampAllowListSendersRemoved)
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

func (it *OnRampAllowListSendersRemovedIterator) Error() error {
	return it.fail
}

func (it *OnRampAllowListSendersRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampAllowListSendersRemoved struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersRemovedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampAllowListSendersRemovedIterator{contract: _OnRamp.contract, event: "AllowListSendersRemoved", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampAllowListSendersRemoved)
				if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseAllowListSendersRemoved(log types.Log) (*OnRampAllowListSendersRemoved, error) {
	event := new(OnRampAllowListSendersRemoved)
	if err := _OnRamp.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampCCIPMessageSentIterator struct {
	Event *OnRampCCIPMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampCCIPMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampCCIPMessageSent)
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
		it.Event = new(OnRampCCIPMessageSent)
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

func (it *OnRampCCIPMessageSentIterator) Error() error {
	return it.fail
}

func (it *OnRampCCIPMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampCCIPMessageSentIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &OnRampCCIPMessageSentIterator{contract: _OnRamp.contract, event: "CCIPMessageSent", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampCCIPMessageSent)
				if err := _OnRamp.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseCCIPMessageSent(log types.Log) (*OnRampCCIPMessageSent, error) {
	event := new(OnRampCCIPMessageSent)
	if err := _OnRamp.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampConfigSetIterator struct {
	Event *OnRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampConfigSet)
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
		it.Event = new(OnRampConfigSet)
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

func (it *OnRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampConfigSet struct {
	StaticConfig  OnRampStaticConfig
	DynamicConfig OnRampDynamicConfig
	Raw           types.Log
}

func (_OnRamp *OnRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OnRampConfigSetIterator, error) {

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OnRampConfigSetIterator{contract: _OnRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampConfigSet)
				if err := _OnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseConfigSet(log types.Log) (*OnRampConfigSet, error) {
	event := new(OnRampConfigSet)
	if err := _OnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampDestChainConfigSetIterator struct {
	Event *OnRampDestChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampDestChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampDestChainConfigSet)
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
		it.Event = new(OnRampDestChainConfigSet)
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

func (it *OnRampDestChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampDestChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampDestChainConfigSet struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Router            common.Address
	AllowlistEnabled  bool
	Raw               types.Log
}

func (_OnRamp *OnRampFilterer) FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampDestChainConfigSetIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampDestChainConfigSetIterator{contract: _OnRamp.contract, event: "DestChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampDestChainConfigSet)
				if err := _OnRamp.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseDestChainConfigSet(log types.Log) (*OnRampDestChainConfigSet, error) {
	event := new(OnRampDestChainConfigSet)
	if err := _OnRamp.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampFeeTokenWithdrawnIterator struct {
	Event *OnRampFeeTokenWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampFeeTokenWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampFeeTokenWithdrawn)
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
		it.Event = new(OnRampFeeTokenWithdrawn)
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

func (it *OnRampFeeTokenWithdrawnIterator) Error() error {
	return it.fail
}

func (it *OnRampFeeTokenWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampFeeTokenWithdrawn struct {
	FeeAggregator common.Address
	FeeToken      common.Address
	Amount        *big.Int
	Raw           types.Log
}

func (_OnRamp *OnRampFilterer) FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampFeeTokenWithdrawnIterator, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &OnRampFeeTokenWithdrawnIterator{contract: _OnRamp.contract, event: "FeeTokenWithdrawn", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampFeeTokenWithdrawn)
				if err := _OnRamp.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseFeeTokenWithdrawn(log types.Log) (*OnRampFeeTokenWithdrawn, error) {
	event := new(OnRampFeeTokenWithdrawn)
	if err := _OnRamp.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampOwnershipTransferRequestedIterator struct {
	Event *OnRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampOwnershipTransferRequested)
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
		it.Event = new(OnRampOwnershipTransferRequested)
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

func (it *OnRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OnRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRamp *OnRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampOwnershipTransferRequestedIterator{contract: _OnRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampOwnershipTransferRequested)
				if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*OnRampOwnershipTransferRequested, error) {
	event := new(OnRampOwnershipTransferRequested)
	if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampOwnershipTransferredIterator struct {
	Event *OnRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampOwnershipTransferred)
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
		it.Event = new(OnRampOwnershipTransferred)
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

func (it *OnRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OnRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRamp *OnRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampOwnershipTransferredIterator{contract: _OnRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OnRamp *OnRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampOwnershipTransferred)
				if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OnRamp *OnRampFilterer) ParseOwnershipTransferred(log types.Log) (*OnRampOwnershipTransferred, error) {
	event := new(OnRampOwnershipTransferred)
	if err := _OnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllowedSendersList struct {
	IsEnabled           bool
	ConfiguredAddresses []common.Address
}
type GetDestChainConfig struct {
	SequenceNumber   uint64
	AllowlistEnabled bool
	Router           common.Address
}

func (_OnRamp *OnRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OnRamp.abi.Events["AllowListAdminSet"].ID:
		return _OnRamp.ParseAllowListAdminSet(log)
	case _OnRamp.abi.Events["AllowListSendersAdded"].ID:
		return _OnRamp.ParseAllowListSendersAdded(log)
	case _OnRamp.abi.Events["AllowListSendersRemoved"].ID:
		return _OnRamp.ParseAllowListSendersRemoved(log)
	case _OnRamp.abi.Events["CCIPMessageSent"].ID:
		return _OnRamp.ParseCCIPMessageSent(log)
	case _OnRamp.abi.Events["ConfigSet"].ID:
		return _OnRamp.ParseConfigSet(log)
	case _OnRamp.abi.Events["DestChainConfigSet"].ID:
		return _OnRamp.ParseDestChainConfigSet(log)
	case _OnRamp.abi.Events["FeeTokenWithdrawn"].ID:
		return _OnRamp.ParseFeeTokenWithdrawn(log)
	case _OnRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _OnRamp.ParseOwnershipTransferRequested(log)
	case _OnRamp.abi.Events["OwnershipTransferred"].ID:
		return _OnRamp.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OnRampAllowListAdminSet) Topic() common.Hash {
	return common.HexToHash("0xb8c9b44ae5b5e3afb195f67391d9ff50cb904f9c0fa5fd520e497a97c1aa5a1e")
}

func (OnRampAllowListSendersAdded) Topic() common.Hash {
	return common.HexToHash("0x330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc3281")
}

func (OnRampAllowListSendersRemoved) Topic() common.Hash {
	return common.HexToHash("0xc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d421586")
}

func (OnRampCCIPMessageSent) Topic() common.Hash {
	return common.HexToHash("0x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32")
}

func (OnRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1")
}

func (OnRampDestChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0xd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5")
}

func (OnRampFeeTokenWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e")
}

func (OnRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OnRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OnRamp *OnRamp) Address() common.Address {
	return _OnRamp.address
}

type OnRampInterface interface {
	GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

		error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

		error)

	GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error)

	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error)

	GetStaticConfig(opts *bind.CallOpts) (OnRampStaticConfig, error)

	GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAllowlistUpdates(opts *bind.TransactOpts, allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error)

	ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error)

	ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OnRampDynamicConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error)

	FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampAllowListAdminSetIterator, error)

	WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error)

	ParseAllowListAdminSet(log types.Log) (*OnRampAllowListAdminSet, error)

	FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersAddedIterator, error)

	WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersAdded(log types.Log) (*OnRampAllowListSendersAdded, error)

	FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampAllowListSendersRemovedIterator, error)

	WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersRemoved(log types.Log) (*OnRampAllowListSendersRemoved, error)

	FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampCCIPMessageSentIterator, error)

	WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error)

	ParseCCIPMessageSent(log types.Log) (*OnRampCCIPMessageSent, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OnRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OnRampConfigSet, error)

	FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampDestChainConfigSetIterator, error)

	WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigSet(log types.Log) (*OnRampDestChainConfigSet, error)

	FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampFeeTokenWithdrawnIterator, error)

	WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenWithdrawn(log types.Log) (*OnRampFeeTokenWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OnRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OnRampOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var OnRampZKBin = ("0x0x0004000000000002001c00000000000200000000030100190000006004300270000005180340019700030000003103550002000000010355000005180040019d0000000100200190000000290000c13d0000008002000039000000400020043f000000040030008c0000004a0000413d000000000201043b000000e002200270000005300020009c000000550000a13d000005310020009c000000ce0000a13d000005320020009c000001bd0000a13d000005330020009c000006780000613d000005340020009c0000040c0000613d000005350020009c0000004a0000c13d000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000401100370000000000101043b0000051c0010009c0000004a0000213d0000054801000041000000800010043f00000549010000410000145e000104300000010004000039000000400040043f0000000002000416000000000002004b0000004a0000c13d0000001f0230003900000519022001970000010002200039000000400020043f0000001f0530018f0000051a0630019800000100026000390000003b0000613d000000000701034f000000007807043c0000000004840436000000000024004b000000370000c13d000000000005004b000000480000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000001400030008c0000004c0000813d00000000010000190000145e00010430000000400800043d0000051b0080009c000000710000a13d0000057401000041000000000010043f0000004101000039000000040010043f00000554010000410000145e000104300000053d0020009c000000f10000213d000005430020009c0000012a0000213d000005460020009c000004210000613d000005470020009c0000004a0000c13d0000000001000416000000000001004b0000004a0000c13d000000c001000039000000400010043f0000001001000039000000800010043f0000058c01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e002000039145c11e90000040f000000c00110008a000005180010009c000005180100804100000060011002100000058d011001c70000145d0001042e0000008001800039000000400010043f000001000100043d0000051c0010009c0000004a0000213d0000000007180436000001200100043d0000051d0010009c0000004a0000213d0000000000170435000001400100043d0000051d0010009c0000004a0000213d00000040098000390000000000190435000001600100043d0000051d0010009c0000004a0000213d000000600a80003900000000001a0435000000400100043d0000051e0010009c0000004f0000213d000000a002100039000000400020043f000001800200043d0000051d0020009c0000004a0000213d0000000002210436000001a00400043d000000000004004b0000000005000039000000010500c039000000000054004b0000004a0000c13d0000000000420435000001c00500043d0000051d0050009c0000004a0000213d00000040041000390000000000540435000001e00600043d0000051d0060009c0000004a0000213d00000060051000390000000000650435000002000b00043d0000051d00b0009c0000004a0000213d0000008006100039001400000006001d0000000000b60435000002200b00043d0000051c00b0009c0000004a0000213d00000100033000390000011f06b00039000000000036004b0000004a0000813d0000010006b00039000000000d0604330000051c00d0009c0000004f0000213d0000000506d002100000003f066000390000051f06600197000000400e00043d000000000c6e001900100000000e001d0000000000ec004b000000000600003900000001060040390000051c00c0009c0000004f0000213d00000001006001900000004f0000c13d0000004000c0043f00000010060000290000000006d60436000f00000006001d000001200bb000390000006006d000c9000000000cb6001900000000003c004b0000004a0000213d00000000000d004b00000cc50000c13d0000000003000411000000000003004b00000ce60000c13d000000400100043d0000052f0200004100000c1e0000013d000005380020009c000001060000213d0000053b0020009c000003da0000613d0000053c0020009c0000004a0000c13d0000000001000416000000000001004b0000004a0000c13d000000000100041a0000051d021001970000000006000411000000000026004b000006d80000c13d0000000102000039000000000302041a0000052204300197000000000464019f000000000042041b0000052201100197000000000010041b00000000010004140000051d05300197000005180010009c0000051801008041000000c0011002100000054b011001c70000800d0200003900000003030000390000057804000041145c14520000040f00000001002001900000004a0000613d00000000010000190000145d0001042e0000053e0020009c0000019e0000213d000005410020009c0000045f0000613d000005420020009c0000004a0000c13d000000440030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000402100370000000000202043b0000051c0020009c0000004a0000213d0000002401100370000000000101043b0000051d0010009c0000004a0000213d145c133e0000040f000006e10000013d000005390020009c000004030000613d0000053a0020009c0000004a0000c13d000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000401100370000000000101043b0000051c0010009c0000004a0000213d000000000010043f0000000601000039000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b000000000101041a0000051c011001970000051c0010009c000006e00000c13d0000057401000041000000000010043f0000001101000039000000040010043f00000554010000410000145e00010430000005440020009c0000051e0000613d000005450020009c0000004a0000c13d000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000402100370000000000202043b0000051c0020009c0000004a0000213d0000002304200039000000000034004b0000004a0000813d0000000404200039000000000441034f000000000504043b0000051c0050009c0000004f0000213d00000005045002100000003f044000390000051f044001970000051b0040009c0000004f0000213d0000008004400039000000400040043f000000800050043f000000240220003900000060045000c90000000004240019000000000034004b0000004a0000213d000000000005004b000007840000c13d0000000101000039000000000101041a0000051d011001970000000002000411000000000012004b000007ad0000c13d000000800100043d000000000001004b000000ef0000613d001400000000001d00000014010000290000000501100210000000a0011000390000000001010433001300000001001d0000000021010434001200000002001d0000051c0110019800000dc10000613d001100000001001d000000000010043f0000000601000039000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d0000001202000029000000000202043300000048032002100000052703300197000000000101043b000000000401041a0000052805400197000000000335019f000000130500002900000040055000390000000005050433000000000005004b0000000005000019000005290500c041000000000353019f000000000031041b0000000001000039000000010100c039000000400300043d000000400530003900000000001504350000051d01200197000000200230003900000000001204350000051c014001970000000000130435000005180030009c000005180300804100000040013002100000000002000414000005180020009c0000051802008041000000c002200210000000000112019f0000052a011001c70000800d0200003900000002030000390000052b040000410000001105000029145c14520000040f00000001002001900000004a0000613d0000001402000029001400010020003d000000800100043d000000140010006b000001580000413d000000ef0000013d0000053f0020009c000005510000613d000005400020009c0000004a0000c13d000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000401100370000000000101043b0000051c0010009c0000004a0000213d000000000010043f0000000601000039000000200010043f00000040020000390000000001000019145c14240000040f000000000101041a0000051c02100197000000800020043f00000552001001980000000002000039000000010200c039000000a00020043f00000048011002700000051d01100197000000c00010043f0000057a010000410000145d0001042e000005360020009c000006aa0000613d000005370020009c0000004a0000c13d000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000402100370000000000202043b000900000002001d0000051c0020009c0000004a0000213d00000009020000290000002302200039000000000032004b0000004a0000813d00000009020000290000000402200039000000000121034f000000000101043b000800000001001d0000051c0010009c0000004a0000213d0000000901000029000a00240010003d000000080100002900000005011002100000000a01100029000000000031004b0000004a0000213d0000000101000039000000000101041a0000051d021001970000000001000411000000000021004b000001e80000613d0000000502000039000000000202041a0000051d02200197000000000021004b000007a90000c13d000000080000006b000000ef0000613d000b00000000001d000002050000013d0000000001230049000005180010009c00000518010080410000006001100210000005180020009c00000518020080410000004002200210000000000121019f0000000002000414000005180020009c0000051802008041000000c002200210000000000121019f0000054b011001c70000800d0200003900000002030000390000057504000041145c14520000040f00000001002001900000004a0000613d0000000b020000290000000102200039000b00000002001d000000080020006c000000ef0000813d0000000b0100002900000005011002100000000a021000290000000201000367000000000221034f000000000302043b0000000002000031000000090420006a000000a30440008a0000055c054001970000055c06300197000000000756013f000000000056004b00000000050000190000055c05004041000000000043004b00000000040000190000055c040080410000055c0070009c000000000504c019000000000005004b0000004a0000c13d0000000a033000290000000004320049000005200040009c0000004a0000213d000000800040008c0000004a0000413d000000400400043d000c00000004001d0000051b0040009c0000004f0000213d0000000c040000290000008004400039000000400040043f000000000431034f000000000404043b0000051c0040009c0000004a0000213d0000000c05000029000000000c4504360000002004300039000000000541034f000000000505043b000000000005004b0000000006000039000000010600c039000000000065004b0000004a0000c13d00000000005c04350000002005400039000000000451034f000000000404043b0000051c0040009c0000004a0000213d00000000073400190000001f04700039000000000024004b00000000060000190000055c060080410000055c084001970000055c04200197000000000948013f000000000048004b00000000080000190000055c080040410000055c0090009c000000000806c019000000000008004b0000004a0000c13d000000000671034f000000000806043b0000051c0080009c0000004f0000213d00000005098002100000003f069000390000051f0a600197000000400600043d000000000aa6001900000000006a004b000000000b000039000000010b0040390000051c00a0009c0000004f0000213d0000000100b001900000004f0000c13d0000004000a0043f000000000086043500000020077000390000000008790019000000000028004b0000004a0000213d000000000087004b0000026d0000813d0000000009060019000000000a71034f000000000a0a043b0000051d00a0009c0000004a0000213d00000020099000390000000000a904350000002007700039000000000087004b000002640000413d0000000c070000290000004007700039001000000007001d00000000006704350000002005500039000000000551034f000000000505043b0000051c0050009c0000004a0000213d00000000053500190000001f03500039000000000023004b00000000060000190000055c060080410000055c03300197000000000743013f000000000043004b00000000030000190000055c030040410000055c0070009c000000000306c019000000000003004b0000004a0000c13d000000000351034f000000000403043b0000051c0040009c0000004f0000213d00000005064002100000003f036000390000051f07300197000000400300043d0000000007730019000000000037004b000000000800003900000001080040390000051c0070009c0000004f0000213d00000001008001900000004f0000c13d000000400070043f000000000043043500000020045000390000000005460019000000000025004b0000004a0000213d000000000054004b000002a60000813d0000000002030019000000000641034f000000000606043b0000051d0060009c0000004a0000213d000000200220003900000000006204350000002004400039000000000054004b0000029d0000413d0000000c010000290000006002100039000f00000002001d000000000032043500000000010104330000051c01100197000000000010043f0000000601000039000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c7000080100200003900140000000c001d145c14570000040f000000140300002900000001002001900000004a0000613d000000000201043b000000000102041a00000570011001970000000003030433000000000003004b00000529040000410000000004006019000000000141019f000000000012041b000000100100002900000000010104330000000004010433000000000004004b001400020020003d0000033d0000613d000000000003004b00000c240000613d001300010020003d0000000003000019000002d50000013d00000012030000290000000103300039000000100100002900000000010104330000000002010433000000000023004b000003150000813d001200000003001d00000005023002100000000001210019000000200110003900000000010104330000051d0110019800000c240000613d001100000001001d000000000010043f0000001401000029000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b000000000101041a000000000001004b000002ce0000c13d0000001302000029000000000102041a0000051c0010009c0000004f0000213d000e00000001001d0000000101100039000000000012041b000000000020043f0000000001000414000005180010009c0000051801008041000000c00110021000000571011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b0000000e011000290000001102000029000000000021041b0000001301000029000000000101041a000e00000001001d000000000020043f0000001401000029000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b0000000e02000029000000000021041b000002ce0000013d0000000c020000290000000005020433000000400200043d000000200300003900000000033204360000000004010433000000000043043500000040032000390000051c05500197000000000004004b000003280000613d0000000006000019000000200110003900000000070104330000051d0770019700000000037304360000000106600039000000000046004b000003210000413d0000000001230049000005180010009c00000518010080410000006001100210000005180020009c00000518020080410000004002200210000000000121019f0000000002000414000005180020009c0000051802008041000000c002200210000000000121019f0000054b011001c70000800d0200003900000002030000390000057204000041145c14520000040f00000001002001900000033e0000c13d0000004a0000013d001300010020003d0000000f0100002900000000010104330000000002010433000000000002004b000002000000613d00000000030000190000034e0000013d000000000101043b000000000001041b000000110300002900000001033000390000000f0100002900000000010104330000000002010433000000000023004b000003c40000813d001100000003001d00000005023002100000000001210019000000200110003900000000010104330000051d01100197001000000001001d000000000010043f0000001401000029000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b000000000201041a000000000002004b000003470000613d0000001301000029000000000301041a000000000003004b000001240000613d000000000032004b001200000002001d000003a50000613d000e00000003001d000000000010043f0000000001000414000005180010009c0000051801008041000000c00110021000000571011001c70000801002000039145c14570000040f00000001002001900000004a0000613d0000001202000029000d000100200092000000000101043b0000001303000029000000000203041a0000000d0020006c000011a90000a13d0000000e02000029000000010220008a0000000001120019000000000101041a000e00000001001d000000000030043f0000000001000414000005180010009c0000051801008041000000c00110021000000571011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b0000000d011000290000000e02000029000000000021041b000000000020043f0000001401000029000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b0000001202000029000000000021041b0000001301000029000000000201041a001200000002001d000000000002004b00000cb20000613d000000000010043f0000000001000414000005180010009c0000051801008041000000c00110021000000571011001c70000801002000039145c14570000040f00000001002001900000004a0000613d0000001202000029000000010220008a000000000101043b0000000001210019000000000001041b0000001301000029000000000021041b0000001001000029000000000010043f0000001401000029000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f0000000100200190000003450000c13d0000004a0000013d000000000002004b000002000000613d0000000c020000290000000005020433000000400200043d000000200300003900000000033204360000000004010433000000000043043500000040032000390000051c05500197000000000004004b000001ec0000613d0000000006000019000000200110003900000000070104330000051d0770019700000000037304360000000106600039000000000046004b000003d20000413d000001ec0000013d0000000001000416000000000001004b0000004a0000c13d0000012001000039000000400010043f000000800000043f000000a00000043f000000c00000043f000000e00000043f000001000000043f145c11cc0000040f0000000201000039000000000101041a0000051d02100197000001200020043f0000054e001001980000000001000039000000010100c039000001400010043f0000000301000039000000000101041a0000051d01100197000001600010043f0000000401000039000000000101041a0000051d01100197000001800010043f0000000501000039000000000101041a0000051d01100197000001a00010043f000000400200043d001400000002001d0000012001000039145c11fb0000040f0000001401000029000005180010009c0000051801008041000000400110021000000579011001c70000145d0001042e0000000001000416000000000001004b0000004a0000c13d0000000101000039000000000101041a0000051d01100197000000800010043f00000576010000410000145d0001042e000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000401100370000000000601043b0000051d0060009c0000004a0000213d0000000101000039000000000101041a0000051d011001970000000005000411000000000015004b000006dc0000c13d000000000056004b000006e80000c13d0000054d01000041000000800010043f00000549010000410000145e000104300000000001000416000000000001004b0000004a0000c13d000000800000043f000000a00000043f000000c00000043f000000e00000043f0000018001000039000000400010043f0000000001000412001c00000001001d001b00000000003d0000800501000039000000440300003900000000040004150000001c0440008a00000005044002100000055a02000041145c14390000040f0000051c01100197000001000010043f0000000001000412001a00000001001d001900200000003d00000000040004150000001a0440008a000000050440021000008005010000390000055a020000410000004403000039145c14390000040f0000051d01100197000001200010043f0000000001000412001800000001001d001700400000003d0000000004000415000000180440008a000000050440021000008005010000390000055a020000410000004403000039145c14390000040f0000051d01100197000001400010043f0000000001000412001600000001001d001500600000003d0000000004000415000000160440008a000000050440021000008005010000390000055a020000410000004403000039145c14390000040f0000051d01100197000001600010043f00000100010000390000018002000039145c11bb0000040f0000058e010000410000145d0001042e000000a40030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000012002000039000000400020043f0000000402100370000000000202043b0000051d0020009c0000004a0000213d000000800020043f0000002403100370000000000503043b000000000005004b0000000003000039000000010300c039000000000035004b0000004a0000c13d000000a00050043f0000004403100370000000000303043b0000051d0030009c0000004a0000213d000000c00030043f0000006404100370000000000404043b0000051d0040009c0000004a0000213d000000e00040043f0000008401100370000000000101043b0000051d0010009c0000004a0000213d000001000010043f0000000106000039000000000606041a0000051d066001970000000007000411000000000067004b00000bbb0000c13d000000000005004b00000bdf0000c13d000000000002004b00000bdf0000613d000000000004004b00000bdf0000613d0000051d011001970000000205000039000000000605041a0000052306600197000000000226019f000000000025041b0000000302000039000000000502041a0000052205500197000000000335019f000000000032041b0000000402000039000000000302041a0000052203300197000000000343019f000000000032041b0000000502000039000000000302041a0000052203300197000000000113019f000000000012041b000001a001000039000000400010043f0000055a0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000101043b0000051c01100197000001200010043f0000055a01000041000000000010044300000000010004120000000400100443000000200100003900000024001004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000101043b0000051d01100197000001400010043f0000055a01000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000101043b0000051d01100197000001600010043f0000055a01000041000000000010044300000000010004120000000400100443000000600100003900000024001004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000101043b0000051d01100197000001800010043f000001200100043d0000051c02100197000000400100043d0000000002210436000001400300043d0000051d033001970000000000320435000001600200043d0000051d0220019700000040031000390000000000230435000001800200043d0000051d0220019700000060031000390000000000230435000000800200043d0000051d0220019700000080031000390000000000230435000000a00200043d000000000002004b0000000002000039000000010200c039000000a0031000390000000000230435000000c00200043d0000051d02200197000000c0031000390000000000230435000000e00200043d0000051d02200197000000e0031000390000000000230435000001000200043d0000051d0220019700000100031000390000000000230435000005180010009c000005180100804100000040011002100000000002000414000005180020009c0000051802008041000000c002200210000000000112019f00000524011001c70000800d0200003900000001030000390000052504000041000000ec0000013d000000440030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000402100370000000000202043b001400000002001d0000051c0020009c0000004a0000213d0000002401100370000000000101043b001300000001001d0000051c0010009c0000004a0000213d000000130130006a000005200010009c0000004a0000213d000000a40010008c0000004a0000413d0000058701000041000000800010043f000000140100002900000080011002100000058801100197000000840010043f0000055a01000041000000000010044300000000010004120000000400100443000000200100003900000024001004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000201043b00000000010004140000051d02200197000000040020008c000006f50000c13d0000000103000031000000200030008c000000200400003900000000040340190000071b0000013d000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000402100370000000000202043b0000051c0020009c0000004a0000213d0000002304200039000000000034004b0000004a0000813d0000000404200039000000000141034f000000000101043b001000000001001d0000051c0010009c0000004a0000213d000f00240020003d000000100100002900000005011002100000000f01100029000000000031004b0000004a0000213d000000100000006b000000ef0000613d0000000401000039000000000101041a0012051d0010019b001400000000001d000005880000013d000000400100043d0000000000a10435000005180010009c000005180100804100000040011002100000000002000414000005180020009c0000051802008041000000c002200210000000000121019f00000571011001c70000800d020000390000000303000039000005820400004100000012050000290000000006090019145c14520000040f00000001002001900000004a0000613d00000014020000290000000102200039001400000002001d000000100020006c000000ef0000813d000000140100002900000005011002100000000f011000290000000201100367000000000901043b0000051d0090009c0000004a0000213d000000400a00043d0000057b0100004100000000001a04350000000401a00039000000000200041000000000002104350000000001000414000000040090008c001300000009001d0000059e0000c13d0000000103000031000000200030008c00000020040000390000000004034019000005cc0000013d0000051800a0009c000005180200004100000000020a40190000004002200210000005180010009c0000051801008041000000c001100210000000000121019f00000554011001c7000000000209001900110000000a001d145c14570000040f000000110a000029000000000301001900000060033002700000051803300197000000200030008c00000020040000390000000004034019000000200640019000000000056a0019000005ba0000613d000000000701034f00000000080a0019000000007907043c0000000008980436000000000058004b000005b60000c13d0000001f07400190000005c70000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000000130900002900000c470000613d0000001f01400039000000600210018f0000000001a20019000000000021004b000000000200003900000001020040390000051c0010009c0000004f0000213d00000001002001900000004f0000c13d000000400010043f000000200040008c0000004a0000413d000000000a0a043300000000000a004b000005830000613d00000044021000390000000000a2043500000020041000390000057c020000410000000000240435000000240210003900000012050000290000000000520435000000440200003900000000002104350000051b0010009c0000004f0000213d000000800b1000390000004000b0043f000005640010009c0000004f0000213d000000c002100039000000400020043f000000200200003900000000002b0435000000a0021000390000057d05000041000000000052043500000000050104330000000001000414000000040090008c00000001020000390000060e0000613d000005180040009c00000518040080410000004002400210000005180050009c00000518050080410000006003500210000000000223019f000005180010009c0000051801008041000000c001100210000000000112019f000000000209001900110000000a001d000e0000000b001d145c14520000040f0000000e0b000029000000110a000029000000130900002900030000000103550000006001100270000105180010019d0000051803100197000000000003004b000000800c000039000000600d0000390000063a0000613d0000051c0030009c0000004f0000213d0000001f013000390000058f011001970000003f011000390000058f01100197000000400d00043d00000000011d00190000000000d1004b000000000400003900000001040040390000051c0010009c0000004f0000213d00000001004001900000004f0000c13d000000400010043f000000000c3d04360000058f0430019800000000014c001900000003050003670000062d0000613d000000000605034f00000000070c0019000000006806043c0000000007870436000000000017004b000006290000c13d0000001f033001900000063a0000613d000000000445034f0000000303300210000000000501043300000000053501cf000000000535022f000000000404043b0000010003300089000000000434022f00000000033401cf000000000353019f000000000031043500000000030d0433000000010020019000000c310000613d000000000003004b000006580000c13d000e0000000d001d000d0000000c001d00110000000a001d0000055601000041000000000010044300000004009004430000000001000414000005180010009c0000051801008041000000c00110021000000557011001c70000800202000039145c14570000040f0000000100200190000011a80000613d000000000101043b000000000001004b0000000e0100002900000ca10000613d0000000003010433000000000003004b0000001309000029000000110a0000290000000d0c000029000005700000613d000005200030009c0000004a0000213d000000200030008c0000004a0000413d00000000020c0433000000000002004b0000000001000039000000010100c039000000000012004b0000004a0000c13d000000400100043d000000000002004b000005710000c13d00000064021000390000057f03000041000000000032043500000044021000390000058003000041000000000032043500000024021000390000002a0300003900000000003204350000057e020000410000000000210435000000040210003900000020030000390000000000320435000005180010009c0000051801008041000000400110021000000581011001c70000145e00010430000000840030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000402100370000000000202043b001400000002001d0000051c0020009c0000004a0000213d0000002402100370000000000202043b0000051c0020009c0000004a0000213d0000000002230049000005200020009c0000004a0000213d000000a40020008c0000004a0000413d0000006401100370000000000101043b001300000001001d0000051d0010009c0000004a0000213d0000000202000039000000000102041a0000054e001001980000073f0000c13d000005500110019700000551011001c7000000000012041b0000001401000029000000000010043f0000000601000039000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000201043b000000130000006b000007b00000c13d000000400100043d0000056e0200004100000c1e0000013d000000240030008c0000004a0000413d0000000002000416000000000002004b0000004a0000c13d0000000401100370000000000101043b0000051c0010009c0000004a0000213d000000000010043f0000000601000039000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b000000000201041a001300000002001d0000000101100039000000000301041a000000400200043d001400000002001d001100000003001d0000000002320436001200000002001d000000000010043f0000000001000414000005180010009c0000051801008041000000c00110021000000571011001c70000801002000039145c14570000040f00000001002001900000004a0000613d0000001105000029000000000005004b000007430000c13d00000012040000290000074c0000013d0000057701000041000000800010043f00000549010000410000145e000104300000054a01000041000000800010043f00000549010000410000145e000104300000000101100039000000400200043d0000000000120435000005180020009c0000051802008041000000400120021000000567011001c70000145d0001042e000000000100041a0000052201100197000000000161019f000000000010041b0000000001000414000005180010009c0000051801008041000000c0011002100000054b011001c70000800d0200003900000003030000390000054c04000041000000ec0000013d000005180010009c0000051801008041000000c00110021000000589011001c7145c14570000040f000000000301001900000060033002700000051803300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000080057001bf0000070a0000613d0000008008000039000000000901034f000000009a09043c0000000008a80436000000000058004b000007060000c13d000000000006004b000007170000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000007660000613d0000001f01400039000000600110018f00000080021001bf001200000002001d000000400020043f000000200030008c0000004a0000413d000000800200043d000000000002004b0000000003000039000000010300c039000000000032004b0000004a0000c13d00000084011001bf000000000002004b00000bbf0000c13d0000000202000039000000000202041a001100000002001d0000058b0200004100000012030000290000000000230435000000130200002900000004032000390000001402000029145c12210000040f000000000300041400000011020000290000051d02200197000000040020008c00000be30000c13d0000000103000031000000200030008c0000002004000039000000000403401900000c100000013d0000054f01000041000000800010043f00000549010000410000145e00010430000000000101043b00000000020000190000001204000029000000000301041a000000000434043600000001011000390000000102200039000000000052004b000007460000413d00000014010000290000000002140049145c11d70000040f000000400300043d001200000003001d000000200130003900000040020000390000000000210435000000130100002900000552001001980000000001000039000000010100c039000000000013043500000040023000390000001401000029145c12130000040f00000012020000290000000001210049000005180010009c0000051801008041000005180020009c000005180200804100000060011002100000004002200210000000000121019f0000145d0001042e0000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000076d0000c13d000000000005004b0000077e0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000005180020009c00000518020080410000004002200210000000000112019f0000145e00010430000000a0050000390000000006230049000005200060009c0000004a0000213d000000600060008c0000004a0000413d000000400600043d000005210060009c0000004f0000213d0000006007600039000000400070043f000000000721034f000000000707043b0000051c0070009c0000004a0000213d00000000087604360000002007200039000000000971034f000000000909043b0000051d0090009c0000004a0000213d00000000009804350000002007700039000000000771034f000000000707043b000000000007004b0000000008000039000000010800c039000000000087004b0000004a0000c13d0000004008600039000000000078043500000000056504360000006002200039000000000042004b000007850000413d0000014e0000013d0000056f01000041000000800010043f00000549010000410000145e00010430000000400100043d0000054a0200004100000c1e0000013d000000000102041a001200000001001d0000055200100198001100000002001d00000bc70000c13d000000120100002900000048011002700000051d011001970000000002000411000000000012004b00000c1c0000c13d0000000301000039000000000101041a0012051d0010019c00000c5f0000c13d0000001101000029000000000101041a0000051c021001970000051c0020009c000001240000613d000005590210019700000001011000390010051c0010019b00000010012001af0000001102000029000000000012041b000000400100043d001200000001001d0000051e0010009c0000004f0000213d0000001202000029000000a001200039000000400010043f0000000001020436001100000001001d0000055a0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000101043b00000012040000290000006002400039000000100300002900000000003204350000004002400039000000140300002900000000003204350000051c01100197000000110200002900000000001204350000008001400039000000000001043500000002020003670000002401200370000000000301043b0000002404300039000000000142034f000000000801043b0000000405300039000000000100003100000000065100490000001f0660008a0000055c076001970000055c09800197000000000a79013f000000000079004b00000000090000190000055c09004041000000000068004b000000000b0000190000055c0b0080410000055c00a0009c00000000090bc019000000000009004b0000004a0000c13d0000000008580019000000000982034f000000000b09043b0000051c00b0009c0000004a0000213d0000000009b10049000000200f8000390000055c089001970000055c0af00197000000000c8a013f00000000008a004b00000000080000190000055c0800404100000000009f004b00000000090000190000055c090020410000055c00c0009c000000000809c019000000000008004b0000004a0000c13d000000200440008a000000000842034f000000000808043b0000055c09800197000000000a79013f000000000079004b00000000090000190000055c09004041000000000068004b000000000c0000190000055c0c0080410000055c00a0009c00000000090cc019000000000009004b0000004a0000c13d0000000009580019000000000892034f000000000808043b0000051c0080009c0000004a0000213d000000000a810049000000200e9000390000055c09a001970000055c0ce00197000000000d9c013f00000000009c004b00000000090000190000055c090040410000000000ae004b000000000a0000190000055c0a0020410000055c00d0009c00000000090ac019000000000009004b0000004a0000c13d0000006004400039000000000942034f000000000909043b001100000009001d0000051d0090009c0000004a0000213d000000200440008a000000000442034f000000000404043b0010055c0040019b000000100a70014f000000100070006b00000000070000190000055c07004041000000000064004b00000000060000190000055c060080410000055c00a0009c000000000706c019000000000007004b0000004a0000c13d0000000005540019000000000652034f000000000606043b000f00000006001d0000051c0060009c0000004a0000213d0000000f06000029000e0006006002180000000e0610006a00000020055000390000055c076001970000055c0a500197000000000c7a013f00000000007a004b00000000070000190000055c07004041000000000065004b00000000060000190000055c060020410000055c00c0009c000000000706c019000000000007004b0000004a0000c13d0000000f06000029000000050c6002100000003f06c00039000b051f0060019b000000400700043d0000000b06700029000d00000007001d000000000076004b000000000700003900000001070040390000051c0060009c0000004f0000213d00000001007001900000004f0000c13d000000400060043f0000000f060000290000000d070000290000000007670436000000000006004b000008910000613d000000600d0000390000000006000019000000400a00043d0000051e00a0009c0000004f0000213d000000a009a00039000000400090043f0000008009a000390000000000d904350000004009a000390000000000d904350000002009a000390000000000d904350000006009a00039000000000009043500000000000a043500000000096700190000000000a9043500000020066000390000000000c6004b0000087e0000413d000000400600043d000800000006001d0000055d0060009c0000004f0000213d00000008070000290000012006700039000000400060043f00000012060000290000000006670436000700000006001d000000130700002900000000007604350000001f06b000390000058f066001970000003f066000390000058f06600197000000400c00043d00000000066c00190000000000c6004b000000000700003900000001070040390000051c0060009c0000004f0000213d00000001007001900000004f0000c13d000000400060043f0000000006bc0436001200000006001d0000000006fb0019000000000016004b0000004a0000213d000c000000f203530000058f0db001980000001f0fb0018f0000001207d00029000008bb0000613d0000000c0a00035f000000120600002900000000a90a043c0000000006960436000000000076004b000008b70000c13d00000000000f004b000008c80000613d0000000c06d0035f0000000309f00210000000000a070433000000000a9a01cf000000000a9a022f000000000606043b0000010009900089000000000696022f00000000069601cf0000000006a6019f00000000006704350000001206b00029000000000006043500000008060000290000004006600039000600000006001d0000000000c604350000001f068000390000058f066001970000003f066000390000058f06600197000000400b00043d00000000066b00190000000000b6004b000000000700003900000001070040390000051c0060009c0000004f0000213d00000001007001900000004f0000c13d000000400060043f000000000c8b04360000000006e80019000000000016004b0000004a0000213d0000000006e2034f0000058f0d8001980000001f0e80018f0000000007dc0019000008eb0000613d000000000a06034f000000000f0c001900000000a90a043c000000000f9f043600000000007f004b000008e70000c13d00000000000e004b000008f80000613d0000000006d6034f0000000309e00210000000000a070433000000000a9a01cf000000000a9a022f000000000606043b0000010009900089000000000696022f00000000069601cf0000000006a6019f000000000067043500000000068c0019000000000006043500000008060000290000006006600039000500000006001d0000000000b60435000000400600043d0000055e0060009c0000004f0000213d0000000007310049000000230770008a000000000074004b00000000080000190000055c080080410000055c07700197000000100970014f000000100070006b00000000070000190000055c070040410000055c0090009c000000000708c0190000002008600039000000400080043f000000000006043500000011080000290000051d088001970000000809000029000000a00a90003900030000000a001d00000000008a04350000008008900039000400000008001d00000000006804350000004406200370000000000606043b0000010008900039000c00000008001d0000000d0a0000290000000000a80435000000c008900039000100000008001d0000000000680435000000e006900039000200000006001d0000000000060435000000000007004b0000004a0000c13d000000400700043d0000000b06700029000b00000007001d000000000076004b000000000700003900000001070040390000051c0060009c0000004f0000213d00000001007001900000004f0000c13d000000400060043f0000000b060000290000000f070000290000000006760436000d00000006001d0000000e06500029000000000016004b0000004a0000213d000000000065004b000009530000813d0000000d070000290000000008510049000005200080009c0000004a0000213d000000400080008c0000004a0000413d000000400800043d0000055f0080009c0000004f0000213d0000004009800039000000400090043f000000000952034f000000000909043b0000051d0090009c0000004a0000213d0000000009980436000000200a500039000000000aa2034f000000000a0a043b0000000000a9043500000000078704360000004005500039000000000065004b0000093c0000413d00000013050000290009051d0050019b0000001405000029000a051c0050019b001300000000001d00000004053000390000000006450019000000000462034f000000000404043b0000051c0040009c0000004a0000213d000000060740021000000000077100490000002006600039000000000076004b00000000080000190000055c080020410000055c077001970000055c06600197000000000976013f000000000076004b00000000060000190000055c060040410000055c0090009c000000000608c019000000000006004b0000004a0000c13d000000130040006b00000dce0000813d0000000b040000290000000004040433000000130040006c000011a90000a13d0000000004310049000000000352034f000000000303043b000000230440008a0000055c064001970000055c07300197000000000867013f000000000067004b00000000060000190000055c06004041000000000043004b00000000040000190000055c040080410000055c0080009c000000000604c019000000000006004b0000004a0000c13d00000013040000290000000506400210000e00000006001d0000000d046000290000000004040433001100000004001d0000000004530019000000000342034f000000000303043b0000051c0030009c0000004a0000213d000000000631004900000020054000390000055c046001970000055c07500197000000000847013f000000000047004b00000000040000190000055c04004041000000000065004b00000000060000190000055c060020410000055c0080009c000000000406c019000000000004004b0000004a0000c13d0000001f043000390000058f044001970000003f044000390000058f04400197000000400600043d0000000004460019001000000006001d000000000064004b000000000600003900000001060040390000051c0040009c0000004f0000213d00000001006001900000004f0000c13d000000400040043f000000100400002900000000043404360000000006530019000000000016004b0000004a0000213d000000000252034f0000058f053001980000000001540019000009bd0000613d000000000602034f0000000007040019000000006806043c0000000007870436000000000017004b000009b90000c13d0000001f06300190000009ca0000613d000000000252034f0000000305600210000000000601043300000000065601cf000000000656022f000000000202043b0000010005500089000000000252022f00000000025201cf000000000262019f000000000021043500000000013400190000000000010435000000400100043d0000051e0010009c0000004f0000213d000000a002100039000000400020043f0000008002100039000000600300003900000000003204350000004002100039000000000032043500000020021000390000000000320435000000600210003900000000000204350000000000010435000000400100043d001200000001001d00000011010000290000002001100039000f00000001001d0000000001010433000000000001004b00000e520000613d0000001101000029000000000101043300000568020000410000001203000029000000000023043500000004023000390000051d0110019700000000001204350000055a01000041000000000010044300000000010004120000000400100443000000600100003900000024001004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000201043b00000000010004140000051d02200197000000040020008c00000a040000c13d0000000103000031000000200030008c0000002004000039000000000403401900000a2e0000013d0000001203000029000005180030009c00000518030080410000004003300210000005180010009c0000051801008041000000c001100210000000000131019f00000554011001c7145c14570000040f000000000301001900000060033002700000051803300197000000200030008c000000200400003900000000040340190000002006400190000000120560002900000a1d0000613d000000000701034f0000001208000029000000007907043c0000000008980436000000000058004b00000a190000c13d0000001f0740019000000a2a0000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000e680000613d0000001f01400039000000600110018f0000001202100029000000000012004b00000000010000390000000101004039001400000002001d0000051c0020009c0000004f0000213d00000001001001900000004f0000c13d0000001401000029000000400010043f000000200030008c0000004a0000413d00000012010000290000000001010433001200000001001d0000051d0010009c0000004a0000213d000000120000006b00000e5b0000613d00000569010000410000001402000029000000000012043500000004012000390000056a02000041000000000021043500000000010004140000001202000029000000040020008c000000200400003900000a7a0000613d0000001402000029000005180020009c00000518020080410000004002200210000005180010009c0000051801008041000000c001100210000000000121019f00000554011001c70000001202000029145c14570000040f000000000301001900000060033002700000051803300197000000200030008c000000200400003900000000040340190000002006400190000000140560002900000a690000613d000000000701034f0000001408000029000000007907043c0000000008980436000000000058004b00000a650000c13d0000001f0740019000000a760000613d000000000661034f0000000307700210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000e740000613d0000001f01400039000000600110018f00000014041000290000051c0040009c0000004f0000213d000000400040043f000000200030008c0000004a0000413d00000014010000290000000001010433000000000001004b0000000002000039000000010200c039000000000021004b0000004a0000c13d000000000001004b00000e5a0000613d0000051e0040009c0000004f0000213d0000000f01000029000000000501043300000011010000290000000001010433000000a002400039000000400020043f0000051d02100197000000800140003900000000002104350000006002400039000000000052043500000040054000390000000906000029000000000065043500000020064000390000000a07000029000000000076043500000010070000290000000000740435000000400900043d0000056b07000041000000000079043500000004079000390000002008000039000000000087043500000000040404330000002407900039000000a0080000390000000000870435000000c40890003900000000740404340000000000480435001400000009001d000000e408900039000000000004004b00000ab90000613d0000000009000019000000000a890019000000000b970019000000000b0b04330000000000ba04350000002009900039000000000049004b00000ab20000413d0000000007840019000000000007043500000000060604330000051c0660019700000014080000290000004407800039000000000067043500000000050504330000051d0550019700000064068000390000000000560435000000000202043300000084058000390000000000250435000000a40280003900000000010104330000051d01100197000000000012043500000000010004140000001202000029000000040020008c00000ad10000c13d000000030100036700000ae90000013d0000001f024000390000058f02200197000000e402200039000005180020009c000005180200804100000060022002100000001403000029000005180030009c00000518030080410000004003300210000000000232019f000005180010009c0000051801008041000000c001100210000000000121019f0000001202000029145c14520000040f00000000030100190000006003300270000105180030019d00000518033001970003000000010355000000010020019000000e800000613d0000058f04300198000000140240002900000af20000613d000000000501034f0000001406000029000000005705043c0000000006760436000000000026004b00000aee0000c13d0000001f0530019000000aff0000613d000000000141034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f00000000001204350000001f013000390000058f021001970000001401200029000000000021004b000000000200003900000001020040390000051c0010009c0000004f0000213d00000001002001900000004f0000c13d000000400010043f000005200030009c0000004a0000213d000000200030008c0000004a0000413d000000140200002900000000020204330000051c0020009c0000004a0000213d0000000004230049000005200040009c0000004a0000213d000000400040008c0000004a0000413d0000055f0010009c0000004f0000213d00000014042000290000004005100039000000400050043f00000000620404340000051c0020009c0000004a0000213d000000000742001900000014023000290000055c032001970000001f087000390000055c09800197000000000a39013f000000000039004b00000000090000190000055c09004041000000000028004b00000000080000190000055c080080410000055c00a0009c000000000908c019000000000009004b0000004a0000c13d00000000870704340000051c0070009c0000004f0000213d0000001f097000390000058f099001970000003f099000390000058f0990019700000000095900190000051c0090009c0000004f0000213d000000400090043f00000000007504350000000009870019000000000029004b0000004a0000213d0000006009100039000000000007004b00000b490000613d000000000a000019000000000b9a0019000000000c8a0019000000000c0c04330000000000cb0435000000200aa0003900000000007a004b00000b420000413d00000000079700190000000000070435000000000551043600000000060604330000051c0060009c0000004a0000213d00000000044600190000001f06400039000000000026004b00000000070000190000055c070080410000055c06600197000000000836013f000000000036004b00000000030000190000055c030040410000055c0080009c000000000307c019000000000003004b0000004a0000c13d00000000640404340000051c0040009c0000004f0000213d0000001f034000390000058f033001970000003f033000390000058f07300197000000400300043d0000000007730019000000000037004b000000000800003900000001080040390000051c0070009c0000004f0000213d00000001008001900000004f0000c13d000000400070043f00000000074304360000000008640019000000000028004b0000004a0000213d000000000004004b00000b7c0000613d000000000200001900000000087200190000000009620019000000000909043300000000009804350000002002200039000000000042004b00000b750000413d000000000247001900000000000204350000000000350435000000400200043d0000051e0020009c0000004f0000213d00000000010104330000000f040000290000000004040433000000a005200039000000400050043f00000060052000390000000000450435000000400420003900000000003404350000002003200039000000000013043500000012010000290000000000120435000000400100043d0000055e0010009c0000004f0000213d0000002003100039000000400030043f0000000000010435000000800320003900000000001304350000000c0100002900000000010104330000000003010433000000130030006c000011a90000a13d0000000e03100029000000200330003900000000002304350000000001010433000000130010006c000011a90000a13d00000002020003670000002401200370000000000301043b0000004401300039000000000112034f000000000401043b00000000010000310000000005310049000000230550008a0000055c065001970000055c07400197000000000867013f000000000067004b00000000060000190000055c06002041000000000054004b00000000050000190000055c050040410000055c0080009c000000000605c0190000001305000029001300010050003d000000000006004b000009580000c13d0000004a0000013d0000054a01000041000001200010043f00000585010000410000145e000104300000058a020000410000001203000029000000000023043500000014020000290000000000210435000000400130021000000554011001c70000145e000104300000001301000029000000000010043f0000000201200039000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b000000000101041a000000000001004b000007b50000c13d000000400100043d0000055302000041000000000021043500000004021000390000001303000029000000000032043500000dc60000013d0000052d01000041000001200010043f00000585010000410000145e0001043000000012040000290000000001410049000005180010009c000005180100804100000060011002100000004004400210000000000141019f000005180030009c0000051803008041000000c003300210000000000131019f145c14570000040f000000000301001900000060033002700000051803300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000120570002900000bff0000613d000000000801034f0000001209000029000000008a08043c0000000009a90436000000000059004b00000bfb0000c13d000000000006004b00000c0c0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000c530000613d0000001f01400039000000600110018f0000001201100029000000400010043f000000200030008c0000004a0000413d000000120200002900000000020204330000000000210435000000400110021000000567011001c70000145d0001042e000000400100043d00000555020000410000000000210435000005180010009c000005180100804100000040011002100000052e011001c70000145e000104300000000c010000290000000001010433000000400200043d000005730300004100000000003204350000051c0110019700000004032000390000000000130435000005180020009c0000051802008041000000400120021000000554011001c70000145e0001043000000000010b0019000000000003004b00000c990000c13d000000400400043d001400000004001d0000057e0200004100000000002404350000000403400039000000200200003900000000002304350000002402400039145c11e90000040f00000014020000290000000001210049000005180010009c0000051801008041000005180020009c000005180200804100000060011002100000004002200210000000000121019f0000145e000104300000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000c4e0000c13d000007710000013d0000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000c5a0000c13d000007710000013d00000556010000410000000000100443000000120100002900000004001004430000000001000414000005180010009c0000051801008041000000c00110021000000557011001c70000800202000039145c14570000040f0000000100200190000011a80000613d000000000101043b000000000001004b0000004a0000613d000000400200043d0000055801000041000000000012043500000024010000390000000201100367000000000101043b0000000403100039001000000002001d00000004012000390000001402000029145c12210000040f00000000020004140000001203000029000000040030008c00000c930000613d00000010030000290000000001310049000005180010009c00000518010080410000006001100210000005180030009c00000518030080410000004003300210000000000131019f000005180020009c0000051802008041000000c002200210000000000121019f0000001202000029145c14520000040f00000000030100190000006003300270000105180030019d0003000000010355000000010020019000000cb80000613d00000010010000290000051c0010009c0000004f0000213d0000001001000029000000400010043f000007bf0000013d0000051800c0009c000005180c0080410000004002c00210000005180030009c00000518030080410000006001300210000000000121019f0000145e00010430000000400100043d00000044021000390000058303000041000000000032043500000024021000390000001d0300003900000000003204350000057e020000410000000000210435000000040210003900000020030000390000000000320435000005180010009c0000051801008041000000400110021000000584011001c70000145e000104300000057401000041000000000010043f0000003101000039000000040010043f00000554010000410000145e0001043000000518033001970000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000cc00000c13d000007710000013d0000000f0d0000290000000006b30049000005200060009c0000004a0000213d000000600060008c0000004a0000413d000000400e00043d0000052100e0009c0000004f0000213d0000006006e00039000000400060043f00000000f60b04340000051c0060009c0000004a0000213d00000000066e0436000000000f0f04330000051d00f0009c0000004a0000213d0000000000f604350000004006b00039000000000f06043300000000000f004b0000000006000039000000010600c03900000000006f004b0000004a0000c13d0000004006e000390000000000f60435000000000ded0436000000600bb000390000000000cb004b00000cc60000413d000000c80000013d0000000106000039000000000b06041a000005220bb0019700000000033b019f000000000036041b00000000030804330000051c0330019800000dcb0000613d00000000060704330000051d0060019800000dcb0000613d00000000060904330000051d0060019800000dcb0000613d00000000060a04330000051d0060019800000dcb0000613d000000800030043f00000000060704330000051d07600197000000a00070043f00000000060904330000051d08600197000000c00080043f00000000060a04330000051d09600197000000e00090043f00000000060104330000051d0b60019800000dcb0000613d00000000060504330000051d0a60019800000dcb0000613d0000000006020433000000000006004b00000dcb0000c13d0000000206000039000000000c06041a000005230cc00197000000000bbc019f0000000000b6041b00000000060404330000051d06600197000000030b000039000000000c0b041a000005220cc0019700000000066c019f00000000006b041b000000040b00003900000000060b041a00000522066001970000000006a6019f00000000006b041b000000140600002900000000060604330000051d06600197000000050a000039000000000b0a041a000005220bb0019700000000066b019f00000000006a041b000000400a00043d0000051b00a0009c0000004f0000213d0000008006a00039000000400060043f0000006006a0003900000000009604350000004009a0003900000000008904350000002008a00039000000000078043500000000003a0435000000400700043d000000000337043600000000080804330000051d08800197000000000083043500000000030904330000051d033001970000004008700039000000000038043500000000030604330000051d033001970000006006700039000000000036043500000000010104330000051d01100197000000800370003900000000001304350000000001020433000000000001004b0000000001000039000000010100c039000000a002700039000000000012043500000000010404330000051d01100197000000c002700039000000000012043500000000010504330000051d01100197000000e0027000390000000000120435000000140100002900000000010104330000051d0110019700000100027000390000000000120435000005180070009c000005180700804100000040017002100000000002000414000005180020009c0000051802008041000000c002200210000000000112019f00000524011001c70000800d0200003900000001030000390000052504000041145c14520000040f00000001002001900000004a0000613d00000010010000290000000001010433000000000001004b00000dad0000613d0000000002000019001200000002001d00000005012002100000000f011000290000000001010433001400000001001d0000000021010434001300000002001d0000051c0110019800000dc10000613d001100000001001d000000000010043f0000000601000039000000200010043f0000000001000414000005180010009c0000051801008041000000c00110021000000526011001c70000801002000039145c14570000040f00000001002001900000004a0000613d0000001302000029000000000202043300000048032002100000052703300197000000000101043b000000000401041a0000052805400197000000000335019f000000140500002900000040055000390000000005050433000000000005004b0000000005000019000005290500c041000000000353019f000000000031041b0000000001000039000000010100c039000000400300043d000000400530003900000000001504350000051d01200197000000200230003900000000001204350000051c014001970000000000130435000005180030009c000005180300804100000040013002100000000002000414000005180020009c0000051802008041000000c002200210000000000112019f0000052a011001c70000800d0200003900000002030000390000052b040000410000001105000029145c14520000040f00000001002001900000004a0000613d0000001202000029000000010220003900000010010000290000000001010433000000000012004b00000d670000413d000000800100043d00000140000004430000016000100443000000a00100043d00000020020000390000018000200443000001a000100443000000c00100043d0000004003000039000001c000300443000001e0001004430000006001000039000000e00300043d000002000010044300000220003004430000010000200443000000040100003900000120001004430000052c010000410000145d0001042e000000400100043d0000058602000041000000000021043500000004021000390000000000020435000005180010009c0000051801008041000000400110021000000554011001c70000145e00010430000000400100043d0000052d0200004100000c1e0000013d0000006407300039000000000672034f0000000204000039000000000404041a000000000606043b0000051d0060009c0000004a0000213d00000000083100490000002003700039000000000332034f000000000303043b000000230780008a0000055c087001970000055c09300197000000000a89013f000000000089004b00000000080000190000055c08004041000000000073004b00000000070000190000055c070080410000055c00a0009c000000000807c019000000000008004b0000004a0000c13d0000000005530019000000000352034f000000000303043b0000051c0030009c0000004a0000213d000000000131004900000020055000390000055c071001970000055c08500197000000000978013f000000000078004b00000000070000190000055c07004041000000000015004b00000000010000190000055c010020410000055c0090009c000000000701c019000000000007004b0000004a0000c13d0000000c010000290000000001010433000000400900043d0000002407900039000000000067043500000560060000410000000006690436001100000006001d00000004069000390000000a0700002900000000007604350000004406200370000000000606043b0000006407900039000000c008000039000000000087043500000044079000390000000000670435000000000652034f000000c40290003900000000003204350000058f073001980000001f0830018f001400000009001d000000e402900039000000000572001900000e1c0000613d000000000906034f000000000a020019000000009b09043c000000000aba043600000000005a004b00000e180000c13d000000000008004b00000e290000613d000000000676034f0000000307800210000000000805043300000000087801cf000000000878022f000000000606043b0000010007700089000000000676022f00000000067601cf000000000686019f00000000006504350013051d0040019b000000000432001900000000000404350000001f033000390000058f03300197000000e0043000390000001405000029000000840550003900000000004504350000000002320019145c13c40000040f0000001403000029000000a4023000390000000003310049000000040330008a00000000003204350000000b0200002900000000020204330000000001210436000000000002004b00000e4b0000613d00000000030000190000000d050000290000000054050434000d00000005001d00000000540404340000051d0440019700000000044104360000000005050433000000000054043500000040011000390000000103300039000000000023004b00000e3f0000413d00000000020004140000001303000029000000040030008c00000e8c0000c13d0000000301000367000000010300003100000ea20000013d0000056d0100004100000012020000290000000000120435000005180020009c000005180200804100000040012002100000052e011001c70000145e00010430001400000004001d000000110100002900000000010104330000056c02000041000000140300002900000000002304350000051d0110019700000004023000390000000000120435000005180030009c0000051803008041000000400130021000000554011001c70000145e000104300000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000e6f0000c13d000007710000013d0000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000e7b0000c13d000007710000013d0000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000e870000c13d000007710000013d00000014030000290000000001310049000005180010009c00000518010080410000006001100210000005180030009c00000518030080410000004003300210000000000131019f000005180020009c0000051802008041000000c002200210000000000121019f0000001302000029145c14570000040f00000000030100190000006003300270000105180030019d00000518033001970003000000010355000000010020019000000f900000613d0000058f043001980000001f0530018f000000140240002900000eac0000613d000000000601034f0000001407000029000000006806043c0000000007870436000000000027004b00000ea80000c13d000000000005004b00000eb90000613d000000000141034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f00000000001204350000001f013000390000058f011001970000001402100029000000000012004b00000000010000390000000101004039001200000002001d0000051c0020009c0000004f0000213d00000001001001900000004f0000c13d0000001201000029000000400010043f000005200030009c0000004a0000213d000000800030008c0000004a0000413d00000011010000290000000001010433000000000001004b0000000002000039000000010200c039000000000021004b0000004a0000c13d0000001402000029000000400220003900000000020204330000051c0020009c0000004a0000213d000000140420002900000014023000290000055c032001970000001f054000390000055c06500197000000000736013f000000000036004b00000000060000190000055c06004041000000000025004b00000000050000190000055c050080410000055c0070009c000000000605c019000000000006004b0000004a0000c13d00000000650404340000051c0050009c0000004f0000213d0000001f045000390000058f044001970000003f044000390000058f0440019700000012074000290000051c0070009c0000004f0000213d00000014040000290000000004040433000000400070043f000000120700002900000000075704360000000008650019000000000028004b0000004a0000213d000000000005004b00000f020000613d00000000080000190000000009780019000000000a680019000000000a0a04330000000000a904350000002008800039000000000058004b00000efb0000413d000000000575001900000000000504350000001405000029000000600550003900000000050504330000051c0050009c0000004a0000213d00000014055000290000001f06500039000000000026004b00000000070000190000055c070080410000055c06600197000000000836013f000000000036004b00000000060000190000055c060040410000055c0080009c000000000607c019000000000006004b0000004a0000c13d00000000670504340000051c0070009c0000004f0000213d00000005087002100000003f098000390000051f09900197000000400a00043d00000000099a001900140000000a001d0000000000a9004b000000000a000039000000010a0040390000051c0090009c0000004f0000213d0000000100a001900000004f0000c13d000000400090043f000000140900002900000000007904350000000007860019000000000027004b0000004a0000213d000000000076004b00000f680000813d000000140800002900000000690604340000051c0090009c0000004a0000213d000000000b5900190000003f09b00039000000000029004b000000000a0000190000055c0a0080410000055c09900197000000000c39013f000000000039004b00000000090000190000055c090040410000055c00c0009c00000000090ac019000000000009004b0000004a0000c13d0000002009b0003900000000090904330000051c0090009c0000004f0000213d0000001f0a9000390000058f0aa001970000003f0aa000390000058f0ca00197000000400a00043d000000000cca00190000000000ac004b000000000d000039000000010d0040390000051c00c0009c0000004f0000213d0000000100d001900000004f0000c13d0000004000c0043f000000000c9a0436000000400bb00039000000000db9001900000000002d004b0000004a0000213d000000000009004b00000f620000613d000000000d000019000000000ecd0019000000000fbd0019000000000f0f04330000000000fe0435000000200dd0003900000000009d004b00000f5b0000413d000000200880003900000000099c001900000000000904350000000000a80435000000000076004b00000f300000413d00000002020000290000000000420435000000000001004b000000000100001900000fd80000c13d000000400300043d00000024013000390000000902000029000000000021043500000561010000410000000000130435001300000003001d00000004013000390000000a0200002900000000002104350000055a01000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000201043b00000000010004140000051d02200197000000040020008c00000f9c0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000fc70000013d0000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000f970000c13d000007710000013d0000001303000029000005180030009c00000518030080410000004003300210000005180010009c0000051801008041000000c001100210000000000131019f00000562011001c7145c14520000040f000000000301001900000060033002700000051803300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000130570002900000fb60000613d000000000801034f0000001309000029000000008a08043c0000000009a90436000000000059004b00000fb20000c13d000000000006004b00000fc30000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000011af0000613d0000001f01400039000000600210018f0000001301200029000000000021004b000000000200003900000001020040390000051c0010009c0000004f0000213d00000001002001900000004f0000c13d000000400010043f000000200030008c0000004a0000413d000000130100002900000000010104330000051c0010009c0000004a0000213d00000008020000290000000002020433000000800220003900000000001204350000000401000029000000120200002900000000002104350000000c0100002900000000010104330000000002010433000000000002004b00000ff70000613d000000000200001900000014030000290000000003030433000000000023004b000011a90000a13d0000000503200210000000200330003900000000011300190000001403300029000000000303043300000000010104330000008001100039000000000031043500000001022000390000000c0100002900000000010104330000000003010433000000000032004b00000fe50000413d000000400100043d001400000001001d00000020021000390000056301000041001300000002001d00000000001204350000055a0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f0000000100200190000011a80000613d000000000101043b000000140400002900000080024000390000000003000410000000000032043500000060024000390000000a0300002900000000003204350000051c0110019700000040024000390000000000120435000000800100003900000000001404350000051e0040009c0000004f0000213d0000001402000029000000a001200039000000400010043f0000001301000029000005180010009c000005180100804100000040011002100000000002020433000005180020009c00000518020080410000006002200210000000000112019f0000000002000414000005180020009c0000051802008041000000c002200210000000000112019f0000054b011001c70000801002000039145c14570000040f00000001002001900000004a0000613d0000000802000029000000000202043300000080032000390000006002200039000000000101043b001200000001001d00000007010000290000000001010433000000000202043300000000030304330000000304000029000000000404043300000001050000290000000005050433000000400700043d000000a00670003900000000005604350000051d04400197000000800570003900000000004504350000051c03300197000000600470003900000000003404350000051c02200197000000400370003900000000002304350000051d0210019700000020017000390000000000210435000000a0020000390000000000270435001400000007001d000005640070009c0000004f0000213d0000001403000029000000c002300039001300000002001d000000400020043f000005180010009c000005180100804100000040011002100000000002030433000005180020009c00000518020080410000006002200210000000000112019f0000000002000414000005180020009c0000051802008041000000c002200210000000000112019f0000054b011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000050200002900000000020204330000002003200039000005180030009c000005180300804100000040033002100000000002020433000005180020009c00000518020080410000006002200210000000000232019f000000000101043b001100000001001d0000000001000414000005180010009c0000051801008041000000c001100210000000000121019f0000054b011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000060200002900000000020204330000002003200039000005180030009c000005180300804100000040033002100000000002020433000005180020009c00000518020080410000006002200210000000000232019f000000000101043b001000000001001d0000000001000414000005180010009c0000051801008041000000c001100210000000000121019f0000054b011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b000f00000001001d0000000c0100002900000000010104330000001404000029000000e0034000390000002002000039000e00000003001d00000000002304350000010002400039145c13c40000040f000000140110006a000000e00210008a00000013030000290000000000230435000000a10110008a0000058f021001970000000001320019000000000021004b000000000200003900000001020040390000051c0010009c0000004f0000213d00000001002001900000004f0000c13d000000400010043f0000000e01000029000005180010009c0000051801008041000000400110021000000013020000290000000002020433000005180020009c00000518020080410000006002200210000000000112019f0000000002000414000005180020009c0000051802008041000000c002200210000000000112019f0000054b011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000040200002900000000020204330000002003200039000005180030009c000005180300804100000040033002100000000002020433000005180020009c00000518020080410000006002200210000000000232019f000000000101043b001400000001001d0000000001000414000005180010009c0000051801008041000000c001100210000000000121019f0000054b011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000201043b000000400100043d000000e0031000390000000000230435000000c00210003900000014030000290000000000320435000000a0021000390000000f030000290000000000320435000000800210003900000010030000290000000000320435000000600210003900000011030000290000000000320435000000400210003900000012030000290000000000320435000000e00200003900000000022104360000000000020435000005650010009c0000004f0000213d0000010003100039000000400030043f000005180020009c000005180200804100000040022002100000000001010433000005180010009c00000518010080410000006001100210000000000121019f0000000002000414000005180020009c0000051802008041000000c002200210000000000112019f0000054b011001c70000801002000039145c14570000040f00000001002001900000004a0000613d000000000101043b000000080300002900000000020304330000000000120435000000000103043300000060011000390000000001010433000000400600043d0000002002000039000000000226043600000000030304330000000054030434000000000042043500000000020504330000051c0220019700000040046000390000000000240435000000400230003900000000020204330000051c0220019700000060046000390000000000240435000000600230003900000000020204330000051c0220019700000080046000390000000000240435000000800230003900000000020204330000051c02200197000000a0036000390000000000230435000000070200002900000000020204330000051d02200197000000c003600039000000000023043500000006020000290000000002020433000000e003600039000001a0040000390000000000430435000001c005600039000000004302043400000000003504350000000008060019000001e002600039000000000003004b000011400000613d000000000500001900000000062500190000000007540019000000000707043300000000007604350000002005500039000000000035004b000011390000413d000000000423001900000000000404350000001f033000390000058f03300197000000050400002900000000040404330000010005800039000001c0063000390000000000650435000000000223001900000000430404340000000002320436000000000003004b000011560000613d000000000500001900000000062500190000000007540019000000000707043300000000007604350000002005500039000000000035004b0000114f0000413d0014051c0010019b000000000123001900000000000104350000001f013000390000058f01100197000000000121001900000000040800190000000002410049000000200220008a000000040300002900000000030304330000012004400039000000000024043500000000320304340000000001210436000000000002004b0000116f0000613d000000000400001900000000051400190000000006430019000000000606043300000000006504350000002004400039000000000024004b000011680000413d00000000031200190000000000030435000000030300002900000000030304330000051d03300197001300000008001d00000140048000390000000000340435000000010300002900000000030304330000016004800039000000000034043500000002030000290000000003030433000001800480003900000000003404350000001f022000390000058f0220019700000000021200190000000001820049000000200310008a0000000c010000290000000001010433000001a0048000390000000000340435145c13c40000040f00000013020000290000000001210049000005180020009c00000518020080410000004002200210000005180010009c00000518010080410000006001100210000000000121019f0000000002000414000005180020009c0000051802008041000000c002200210000000000121019f0000054b011001c70000800d02000039000000030300003900000566040000410000000a050000290000001406000029145c14520000040f00000001002001900000004a0000613d0000000202000039000000000102041a0000055001100197000000000012041b000000080100002900000000010104330000000001010433000006e10000013d000000000001042f0000057401000041000000000010043f0000003201000039000000040010043f00000554010000410000145e000104300000001f0530018f0000051a06300198000000400200043d0000000004620019000007710000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000011b60000c13d000007710000013d00000000430104340000051c03300197000000000332043600000000040404330000051d044001970000000000430435000000400310003900000000030304330000051d03300197000000400420003900000000003404350000006002200039000000600110003900000000010104330000051d011001970000000000120435000000000001042d000005900010009c000011d10000813d000000a001100039000000400010043f000000000001042d0000057401000041000000000010043f0000004101000039000000040010043f00000554010000410000145e000104300000001f022000390000058f022001970000000001120019000000000021004b000000000200003900000001020040390000051c0010009c000011e30000213d0000000100200190000011e30000c13d000000400010043f000000000001042d0000057401000041000000000010043f0000004101000039000000040010043f00000554010000410000145e0001043000000000430104340000000001320436000000000003004b000011f50000613d000000000200001900000000051200190000000006240019000000000606043300000000006504350000002002200039000000000032004b000011ee0000413d000000000213001900000000000204350000001f023000390000058f022001970000000001210019000000000001042d00000000430104340000051d0330019700000000033204360000000004040433000000000004004b0000000004000039000000010400c0390000000000430435000000400310003900000000030304330000051d0330019700000040042000390000000000340435000000600310003900000000030304330000051d03300197000000600420003900000000003404350000008002200039000000800110003900000000010104330000051d011001970000000000120435000000000001042d000000000301001900000000040304330000000001420436000000000004004b000012200000613d0000000002000019000000200330003900000000050304330000051d0550019700000000015104360000000102200039000000000042004b000012190000413d000000000001042d00010000000000020000002004100039000000400500003900000000005404350000051c0220019700000000002104350000000202000367000000000432034f000000000404043b000100000000003500000000063000790000001f0660008a0000055c076001970000055c08400197000000000978013f000000000078004b00000000080000190000055c08002041000000000064004b000000000a0000190000055c0a0040410000055c0090009c00000000080ac019000000000008004b0000133c0000613d0000000004340019000000000842034f000000000908043b0000051c0090009c0000133c0000213d00000020044000390000000108900069000000000084004b000000000a0000190000055c0a0020410000055c088001970000055c0b400197000000000c8b013f00000000008b004b00000000080000190000055c080040410000055c00c0009c00000000080ac019000000000008004b0000133c0000c13d0000004008100039000000a00a0000390000000000a80435000000e0081000390000000000980435000000000a42034f0000058f0c9001980000001f0d90018f000001000b1000390000000008cb00190000125f0000613d000000000e0a034f000000000f0b001900000000e40e043c000000000f4f043600000000008f004b0000125b0000c13d00000000000d004b0000126c0000613d0000000004ca034f000000030ad00210000000000c080433000000000cac01cf000000000cac022f000000000404043b000001000aa000890000000004a4022f0000000004a401cf0000000004c4019f00000000004804350000000004b900190000000000040435000000200a3000390000000004a2034f000000000804043b0000055c04800197000000000c74013f000000000074004b00000000040000190000055c04004041000000000068004b000000000d0000190000055c0d0080410000055c00c0009c00000000040dc019000000000004004b0000133c0000c13d0000000004380019000000000842034f000000000808043b0000051c0080009c0000133c0000213d000000200c400039000000010480006900000000004c004b000000000d0000190000055c0d0020410000055c044001970000055c0ec00197000000000f4e013f00000000004e004b00000000040000190000055c040040410000055c00f0009c00000000040dc019000000000004004b0000133c0000c13d0000001f049000390000058f044001970000000009b40019000000c004400039000000600b10003900000000004b0435000000000cc2034f00000000048904360000058f0d8001980000001f0e80018f0000000009d40019000012a30000613d000000000f0c034f000000000b04001900000000f50f043c000000000b5b043600000000009b004b0000129f0000c13d00000000000e004b000012b00000613d0000000005dc034f000000030be00210000000000c090433000000000cbc01cf000000000cbc022f000000000505043b000001000bb000890000000005b5022f0000000005b501cf0000000005c5019f0000000000590435000000000548001900000000000504350000002009a00039000000000592034f000000000a05043b0000055c05a00197000000000b75013f000000000075004b00000000050000190000055c0500404100000000006a004b000000000c0000190000055c0c0080410000055c00b0009c00000000050cc019000000000005004b0000133c0000c13d000000000b3a00190000000005b2034f000000000a05043b0000051c00a0009c0000133c0000213d000000200cb000390000000605a00210000000010550006900000000005c004b000000000b0000190000055c0b0020410000055c055001970000055c0dc00197000000000e5d013f00000000005d004b00000000050000190000055c050040410000055c00e0009c00000000050bc019000000000005004b0000133c0000c13d0000001f058000390000058f0550019700000000044500190000000005140049000000400550008a000000800810003900000000005804350000000008a4043600000000000a004b000012ef0000613d000000000b0000190000000004c2034f000000000404043b0000051d0040009c0000133c0000213d00000000044804360000002005c00039000000000552034f000000000505043b0000000000540435000000400cc000390000004008800039000000010bb000390000000000ab004b000012e10000413d0000002004900039000000000542034f000000000905043b0000051d0090009c0000133c0000213d000000a00510003900000000009504350000002004400039000000000442034f000000000904043b0000055c04900197000000000574013f000000000074004b00000000040000190000055c04004041000000000069004b00000000060000190000055c060080410000055c0050009c000000000406c019000000000004004b0000133c0000c13d0000000004390019000000000342034f000000000303043b0000051c0030009c0000133c0000213d00000020064000390000000104300069000000000046004b00000000050000190000055c050020410000055c044001970000055c07600197000000000947013f000000000047004b00000000040000190000055c040040410000055c0090009c000000000405c019000000000004004b0000133c0000c13d0000000004180049000000400440008a000000c0011000390000000000410435000000000562034f00000000013804360000058f063001980000001f0730018f0000000002610019000013290000613d000000000805034f0000000004010019000000008908043c0000000004940436000000000024004b000013250000c13d000000000007004b000013360000613d000000000465034f0000000305700210000000000602043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f0000000000420435000000000231001900000000000204350000001f023000390000058f022001970000000001210019000000000001042d00000000010000190000145e000104300001000000000002000000400300043d000005680200004100000000002304350000051d01100197000100000003001d000000040230003900000000001204350000055a01000041000000000010044300000000010004120000000400100443000000600100003900000024001004430000000001000414000005180010009c0000051801008041000000c0011002100000055b011001c70000800502000039145c14570000040f00000001002001900000139f0000613d000000000201043b00000000010004140000051d02200197000000040020008c000013600000c13d0000000103000031000000200030008c00000020040000390000000004034019000000010b0000290000138c0000013d0000000103000029000005180030009c00000518030080410000004003300210000005180010009c0000051801008041000000c001100210000000000131019f00000554011001c7145c14570000040f000000010b000029000000000301001900000060033002700000051803300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b00190000137b0000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b000013770000c13d000000000006004b000013880000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000013a60000613d0000001f01400039000000600210018f0000000001b20019000000000021004b000000000200003900000001020040390000051c0010009c000013a00000213d0000000100200190000013a00000c13d000000400010043f0000001f0030008c0000139d0000a13d00000000010b04330000051d0010009c0000139d0000213d000000000001042d00000000010000190000145e00010430000000000001042f0000057401000041000000000010043f0000004101000039000000040010043f00000554010000410000145e000104300000001f0530018f0000051a06300198000000400200043d0000000004620019000013b10000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000013ad0000c13d000000000005004b000013be0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000005180020009c00000518020080410000004002200210000000000112019f0000145e0001043000000000040104330000000000420435000000050340021000000000033200190000002003300039000000000004004b000014210000613d000000a00500003900000000070000190000000008020019000013d70000013d000000000a39001900000000000a04350000001f099000390000058f0990019700000000033900190000000107700039000000000047004b000014210000813d0000000009230049000000200990008a000000200880003900000000009804350000002001100039000000000901043300000000ba0904340000051d0aa00197000000000aa30436000000000b0b043300000000005a0435000000a00d30003900000000ca0b04340000000000ad0435000000c00b30003900000000000a004b000013f00000613d000000000d000019000000000ebd0019000000000fdc0019000000000f0f04330000000000fe0435000000200dd000390000000000ad004b000013e90000413d000000000cba001900000000000c04350000001f0aa000390000058f0aa00197000000000aba0019000000400b900039000000000b0b0433000000000c3a0049000000400d3000390000000000cd043500000000cb0b0434000000000aba043600000000000b004b000014060000613d000000000d000019000000000ead0019000000000fdc0019000000000f0f04330000000000fe0435000000200dd000390000000000bd004b000013ff0000413d000000000cab001900000000000c0435000000600c900039000000000c0c0433000000600d3000390000000000cd04350000001f0bb000390000058f0bb00197000000000bab001900000080099000390000000009090433000000000a3b004900000080033000390000000000a3043500000000a909043400000000039b0436000000000009004b000013cf0000613d000000000b000019000000000c3b0019000000000dba0019000000000d0d04330000000000dc0435000000200bb0003900000000009b004b000014190000413d000013cf0000013d0000000001030019000000000001042d000000000001042f000005180010009c00000518010080410000004001100210000005180020009c00000518020080410000006002200210000000000112019f0000000002000414000005180020009c0000051802008041000000c002200210000000000112019f0000054b011001c70000801002000039145c14570000040f0000000100200190000014370000613d000000000101043b000000000001042d00000000010000190000145e0001043000000000050100190000000000200443000000040100003900000005024002700000000002020031000000000121043a0000002004400039000000000031004b0000143c0000413d000005180030009c000005180300804100000060013002100000000002000414000005180020009c0000051802008041000000c002200210000000000112019f00000591011001c70000000002050019145c14570000040f0000000100200190000014510000613d000000000101043b000000000001042d000000000001042f00001455002104210000000102000039000000000001042d0000000002000019000000000001042d0000145a002104230000000102000039000000000001042d0000000002000019000000000001042d0000145c000004320000145d0001042e0000145e00010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000000000000000000000000000ffffffffffffff7f000000000000000000000000000000000000000000000000ffffffffffffffff000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff9fffffffffffffffffffffffff0000000000000000000000000000000000000000ffffffffffffffffffffff0000000000000000000000000000000000000000000200000000000000000000000000000000000120000000000000000000000000c7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f10200000000000000000000000000000000000040000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000ffffff000000000000000000000000000000000000000000ffffffffffffffff00000000000000000000000000000000000000000000000100000000000000000200000000000000000000000000000000000060000000000000000000000000d5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5000000020000000000000000000000000000014000000100000000000000000035be3ac80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000009b15e16f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007437ff9e00000000000000000000000000000000000000000000000000000000972b461100000000000000000000000000000000000000000000000000000000df0aa9e800000000000000000000000000000000000000000000000000000000df0aa9e900000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000fbca3b7400000000000000000000000000000000000000000000000000000000972b461200000000000000000000000000000000000000000000000000000000c9b146b3000000000000000000000000000000000000000000000000000000008da5cb5a000000000000000000000000000000000000000000000000000000008da5cb5b000000000000000000000000000000000000000000000000000000009041be3d000000000000000000000000000000000000000000000000000000007437ff9f0000000000000000000000000000000000000000000000000000000079ba50970000000000000000000000000000000000000000000000000000000027e936f0000000000000000000000000000000000000000000000000000000005cb80c5c000000000000000000000000000000000000000000000000000000005cb80c5d000000000000000000000000000000000000000000000000000000006def4ce70000000000000000000000000000000000000000000000000000000027e936f10000000000000000000000000000000000000000000000000000000048a98aa40000000000000000000000000000000000000000000000000000000020487dec0000000000000000000000000000000000000000000000000000000020487ded000000000000000000000000000000000000000000000000000000002716072b0000000000000000000000000000000000000000000000000000000006285c6900000000000000000000000000000000000000000000000000000000181f5a779e7177c80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000008000000000000000002b5c74de000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000ff00000000000000000000000000000000000000003ee5aeb500000000000000000000000000000000000000000000000000000000ffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff00000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000ff0000000000000000d0d259760000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000240000000000000000000000001c0a3529000000000000000000000000000000000000000000000000000000001806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b830200000200000000000000000000000000000024000000000000000000000000e0a0e50600000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e02000002000000000000000000000000000000440000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffedf000000000000000000000000000000000000000000000000ffffffffffffffdf000000000000000000000000000000000000000000000000ffffffffffffffbf430d138c00000000000000000000000000000000000000000000000000000000ea458c0c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321000000000000000000000000000000000000000000000000ffffffffffffff3f000000000000000000000000000000000000000000000000fffffffffffffeff192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f320000000000000000000000000000000000000020000000000000000000000000bbe4f6db0000000000000000000000000000000000000000000000000000000001ffc9a700000000000000000000000000000000000000000000000000000000aff2afbf000000000000000000000000000000000000000000000000000000009a4575b900000000000000000000000000000000000000000000000000000000bf16aab6000000000000000000000000000000000000000000000000000000005cf0444900000000000000000000000000000000000000000000000000000000a4ec747900000000000000000000000000000000000000000000000000000000905d7d9b00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff0200000000000000000000000000000000000020000000000000000000000000330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc3281463258ff000000000000000000000000000000000000000000000000000000004e487b7100000000000000000000000000000000000000000000000000000000c237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d421586000000000000000000000000000000000000002000000080000000000000000002b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000006000000080000000000000000070a0823100000000000000000000000000000000000000000000000000000000a9059cbb000000000000000000000000000000000000000000000000000000005361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656408c379a0000000000000000000000000000000000000000000000000000000006f742073756363656564000000000000000000000000000000000000000000005361666545524332303a204552433230206f7065726174696f6e20646964206e0000000000000000000000000000000000000084000000000000000000000000508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000004000001200000000000000000c35aa79d000000000000000000000000000000000000000000000000000000002cbc26bb000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff000000000000000000000000000000000000000000000000000000000000000000000024000000800000000000000000fdbd6a7200000000000000000000000000000000000000000000000000000000d8694ccd000000000000000000000000000000000000000000000000000000004f6e52616d7020312e362e302d646576000000000000000000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000080000001800000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffff6002000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
