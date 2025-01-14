// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package onramp_with_message_transformer

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

var OnRampWithMessageTransformerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"destChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowlistUpdates\",\"inputs\":[{\"name\":\"allowlistConfigArgsItems\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.AllowlistConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"addedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"removedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyDestChainConfigUpdates\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forwardFromRouter\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllowedSendersList\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configuredAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExpectedNextSequenceNumber\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageTransformer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPoolBySourceToken\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPoolV1\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedTokens\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFeeTokens\",\"inputs\":[{\"name\":\"feeTokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdminSet\",\"inputs\":[{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersAdded\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersRemoved\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CCIPMessageSent\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeValueJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigSet\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTokenWithdrawn\",\"inputs\":[{\"name\":\"feeAggregator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feeToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotSendZeroTokens\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"GetSupportedTokensFunctionalityRemovedCheckAdminRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAllowListRequest\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageTransformError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeCalledByRouter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAllowlistAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RouterMustSetOriginalSender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnsupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610100604052346105a1576140358038038061001a816105db565b92833981019080820361016081126105a157608081126105a15761003c6105bc565b9161004681610600565b835260208101516001600160a01b03811681036105a1576020840190815261007060408301610614565b916040850192835260a061008660608301610614565b6060870190815294607f1901126105a15760405160a081016001600160401b038111828210176105a6576040526100bf60808301610614565b81526100cd60a08301610628565b602082019081526100e060c08401610614565b90604083019182526100f460e08501610614565b92606081019384526101096101008601610614565b608082019081526101208601519095906001600160401b0381116105a15781018b601f820112156105a1578051906001600160401b0382116105a6578160051b602001610155906105db565b9c8d838152602001926060028201602001918183116105a157602001925b828410610533575050505061014061018b9101610614565b98331561052257600180546001600160a01b0319163317905580516001600160401b0316158015610510575b80156104fe575b80156104ec575b6104bf57516001600160401b0316608081905295516001600160a01b0390811660a08190529751811660c08190529851811660e081905282519091161580156104da575b80156104d0575b6104bf57815160028054855160ff60a01b90151560a01b166001600160a01b039384166001600160a81b0319909216919091171790558451600380549183166001600160a01b03199283161790558651600480549184169183169190911790558751600580549190931691161790557fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1986101209860606102af6105bc565b8a8152602080820193845260408083019586529290910194855281519a8b5291516001600160a01b03908116928b019290925291518116918901919091529051811660608801529051811660808701529051151560a08601529051811660c08501529051811660e0840152905116610100820152a160005b82518110156104185761033a8184610635565b516001600160401b0361034d8386610635565b5151169081156104035760008281526006602090815260409182902081840151815494840151600160401b600160e81b03198616604883901b600160481b600160e81b031617901515851b68ff000000000000000016179182905583516001600160401b0390951685526001600160a01b031691840191909152811c60ff1615159082015260019291907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef590606090a201610327565b5063c35aa79d60e01b60005260045260246000fd5b506001600160a01b031680156104ae57600780546001600160a01b0319169190911790556040516139d5908161066082396080518181816103ee01528181610b2501528181611ed401526126be015260a051818181611f0d0152818161242a01526126f7015260c051818181610d9601528181611f490152612733015260e051818181611f850152818161276f0152612d6e0152f35b6342bcdf7f60e11b60005260046000fd5b6306b7c75960e31b60005260046000fd5b5082511515610210565b5084516001600160a01b031615610209565b5088516001600160a01b0316156101c5565b5087516001600160a01b0316156101be565b5086516001600160a01b0316156101b7565b639b15e16f60e01b60005260046000fd5b6060848303126105a15760405190606082016001600160401b038111838210176105a65760405261056385610600565b82526020850151906001600160a01b03821682036105a1578260209283606095015261059160408801610628565b6040820152815201930192610173565b600080fd5b634e487b7160e01b600052604160045260246000fd5b60405190608082016001600160401b038111838210176105a657604052565b6040519190601f01601f191682016001600160401b038111838210176105a657604052565b51906001600160401b03821682036105a157565b51906001600160a01b03821682036105a157565b519081151582036105a157565b80518210156106495760209160051b010190565b634e487b7160e01b600052603260045260246000fdfe608080604052600436101561001357600080fd5b600090813560e01c90816306285c69146126585750806315777ab214612606578063181f5a771461258757806320487ded1461234e5780632716072b1461209e57806327e936f114611c9857806348a98aa414611c155780635cb80c5d146119585780636def4ce7146118c95780637437ff9f146117ac57806379ba5097146116c75780638da5cb5b146116755780639041be3d146115c8578063972b4612146114fa578063c9b146b314611131578063df0aa9e914610237578063f2fde38b1461014a5763fbca3b74146100e757600080fd5b346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757600490610121612964565b507f9e7177c8000000000000000000000000000000000000000000000000000000008152fd5b80fd5b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475773ffffffffffffffffffffffffffffffffffffffff6101976129ba565b61019f61329a565b1633811461020f57807fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788380a380f35b6004827fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b50346101475760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475761026f612964565b67ffffffffffffffff6024351161112d5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc6024353603011261112d576102b66129dd565b60025460ff8160a01c16611105577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001760025567ffffffffffffffff8216835260066020526040832073ffffffffffffffffffffffffffffffffffffffff8216156110dd57805460ff8160401c1661106f575b60481c73ffffffffffffffffffffffffffffffffffffffff1633036110475773ffffffffffffffffffffffffffffffffffffffff6003541680610fcf575b50805467ffffffffffffffff811667ffffffffffffffff8114610fa2579067ffffffffffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009493011692839116179055604051906103e08261282e565b84825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602083015267ffffffffffffffff8416604083015260608201528360808201526104416024803501602435600401612f06565b61045060046024350180612f06565b61045e606460243501612e12565b93610473604460243501602435600401612f57565b9490507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06104b96104a387612995565b966104b16040519889612883565b808852612995565b018a5b818110610f8b57505061050d93929161050191604051986104dc8a61284a565b895273ffffffffffffffffffffffffffffffffffffffff8a1660208a01523691612fd7565b60408701523691612fd7565b606084015273ffffffffffffffffffffffffffffffffffffffff6020926040516105378582612883565b88815260808601521660a084015260443560c08401528560e084015261010083015261056d604460243501602435600401612f57565b61057981969296612995565b906105876040519283612883565b8082528382018097368360061b820111610f875780915b8360061b82018310610f505750505050865b6105c4604460243501602435600401612f57565b905081101561093a576105d78183612ec3565b51906105f26105eb60046024350180612f06565b3691612fd7565b916105fb612fab565b5085810151156109125773ffffffffffffffffffffffffffffffffffffffff61062681835116612d0f565b169283158015610868575b61082557908a6106de92888873ffffffffffffffffffffffffffffffffffffffff8d838701518280895116926040519761066a8961282e565b885267ffffffffffffffff87890196168652816040890191168152606088019283526080880193845267ffffffffffffffff6040519b8c998a997f9a4575b9000000000000000000000000000000000000000000000000000000008b5260048b01525160a060248b015260c48a0190612921565b965116604488015251166064860152516084850152511660a4830152038183885af191821561081a578b92610773575b50600193828492898061076c965193015191015191604051936107308561282e565b84528a8401526040830152606082015260405161074d8982612883565b8c81526080820152610100890151906107668383612ec3565b52612ec3565b50016105b0565b91503d808c843e6107848184612883565b82019187818403126108125780519067ffffffffffffffff821161081657019360408584031261081257604051916107bb83612867565b855167ffffffffffffffff811161080e57846107d891880161300e565b8352888601519267ffffffffffffffff841161080e5761080061076c9587956001990161300e565b8a820152935091509361070e565b8d80fd5b8b80fd5b8c80fd5b6040513d8d823e3d90fd5b60248b73ffffffffffffffffffffffffffffffffffffffff8451167fbf16aab6000000000000000000000000000000000000000000000000000000008252600452fd5b506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf0000000000000000000000000000000000000000000000000000000060048201528781602481885afa908115610907578c916108d2575b5015610631565b90508781813d8311610900575b6108e98183612883565b81010312610812576108fa90612a9c565b386108cb565b503d6108df565b6040513d8e823e3d90fd5b60048a7f5cf04449000000000000000000000000000000000000000000000000000000008152fd5b50957f430d138c0000000000000000000000000000000000000000000000000000000081939497839773ffffffffffffffffffffffffffffffffffffffff600254169187610a2c8c6109fc610993606460243501612e12565b9161010073ffffffffffffffffffffffffffffffffffffffff6109c0608460243501602435600401612f06565b92909301519467ffffffffffffffff6040519e8f9d8e521660048d01521660248b015260443560448b015260c060648b015260c48a0191612af9565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc888303016084890152613050565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8684030160a48701525191828152019190855b89828210610f15575050505082809103915afa938415610f0a5782839084938597610e10575b5060e089015215610d275750815b67ffffffffffffffff6080885101911690526080860152805b61010086015151811015610ae45780610ac960019286612ec3565b516080610adb836101008b0151612ec3565b51015201610aae565b5083610aef86613315565b91604051848101907f130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015267ffffffffffffffff8416606082015230608082015260808152610b6f60a082612883565b51902073ffffffffffffffffffffffffffffffffffffffff8585015116845167ffffffffffffffff6080816060840151169201511673ffffffffffffffffffffffffffffffffffffffff60a08801511660c088015191604051938a850195865260408501526060840152608083015260a082015260a08152610bf260c082612883565b51902060608501518681519101206040860151878151910120610100870151604051610c5881610c2c8c8201948d86526040830190613050565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282612883565b51902091608088015189815191012093604051958a870197885260408701526060860152608085015260a084015260c083015260e082015260e08152610ca061010082612883565b51902082515267ffffffffffffffff60608351015116907f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f3267ffffffffffffffff60405192169180610cf2868261313b565b0390a37fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff600254166002555151604051908152f35b73ffffffffffffffffffffffffffffffffffffffff604051917fea458c0c00000000000000000000000000000000000000000000000000000000835267ffffffffffffffff8816600484015216602482015283816044818673ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af1908115610e05578391610dcc575b50610a95565b90508381813d8311610dfe575b610de38183612883565b81010312610dfa57610df490613126565b38610dc6565b8280fd5b503d610dd9565b6040513d85823e3d90fd5b9650505090503d8083863e610e258186612883565b840190608085830312610dfa57845194610e40858201612a9c565b95604082015167ffffffffffffffff8111610f025784610e6191840161300e565b9160608101519067ffffffffffffffff8211610f06570184601f82011215610f02578051610e8e81612995565b95610e9c6040519788612883565b818752888088019260051b84010192818411610efe57898101925b848410610ecd5750505050509590929538610a87565b835167ffffffffffffffff8111610efa578b91610eef8584809487010161300e565b815201930192610eb7565b8a80fd5b8880fd5b8580fd5b8680fd5b6040513d84823e3d90fd5b8351805173ffffffffffffffffffffffffffffffffffffffff1686528101518186015289975088965060409094019390920191600101610a61565b604083360312610efa57866040918251610f6981612867565b610f7286612a00565b8152828601358382015281520192019161059e565b8980fd5b602090610f96612fab565b82828a010152016104bc565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b803b15611043578460405180927fe0a0e5060000000000000000000000000000000000000000000000000000000082528183816110156024356004018b60048401612b38565b03925af18015611038571561037e578461103191959295612883565b923861037e565b6040513d87823e3d90fd5b8480fd5b6004847f1c0a3529000000000000000000000000000000000000000000000000000000008152fd5b73ffffffffffffffffffffffffffffffffffffffff831660009081526002830160205260409020546103405760248573ffffffffffffffffffffffffffffffffffffffff857fd0d2597600000000000000000000000000000000000000000000000000000000835216600452fd5b6004847fa4ec7479000000000000000000000000000000000000000000000000000000008152fd5b6004847f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b5080fd5b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475760043567ffffffffffffffff811161112d57611181903690600401612a21565b73ffffffffffffffffffffffffffffffffffffffff6001541633036114b2575b919081907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8181360301915b848110156114ae578060051b820135838112156110435782019160808336031261104357604051946111fd866127e3565b61120684612980565b8652611214602085016129ad565b9660208701978852604085013567ffffffffffffffff8111610dfa5761123d9036908701612e5e565b9460408801958652606081013567ffffffffffffffff81116114aa5761126591369101612e5e565b60608801908152875167ffffffffffffffff1683526006602052604080842099518a547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169015159182901b68ff000000000000000016178a55909590815151611382575b5095976001019550815b85518051821015611313579061130c73ffffffffffffffffffffffffffffffffffffffff61130483600195612ec3565b51168961376a565b50016112d4565b50509590969450600192919351908151611333575b5050019392936111cc565b61137867ffffffffffffffff7fc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d42158692511692604051918291602083526020830190612a52565b0390a23880611328565b9893959296919094979860001461147357600184019591875b86518051821015611418576113c58273ffffffffffffffffffffffffffffffffffffffff92612ec3565b511680156113e157906113da6001928a6136d9565b500161139b565b60248a67ffffffffffffffff8e51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b50509692955090929796937f330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc328161146967ffffffffffffffff8a51169251604051918291602083526020830190612a52565b0390a238806112ca565b60248767ffffffffffffffff8b51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b8380fd5b8380f35b73ffffffffffffffffffffffffffffffffffffffff600554163303156111a1576004837f905d7d9b000000000000000000000000000000000000000000000000000000008152fd5b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475767ffffffffffffffff61153b612964565b16808252600660205260ff604083205460401c16908252600660205260016040832001916040518093849160208254918281520191845260208420935b8181106115af57505061158d92500383612883565b6115ab60405192839215158352604060208401526040830190612a52565b0390f35b8454835260019485019487945060209093019201611578565b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475767ffffffffffffffff611609612964565b1681526006602052600167ffffffffffffffff604083205416019067ffffffffffffffff82116116485760208267ffffffffffffffff60405191168152f35b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526011600452fd5b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757805473ffffffffffffffffffffffffffffffffffffffff81163303611784577fffffffffffffffffffffffff000000000000000000000000000000000000000060015491338284161760015516825573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b6004827f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610147576117e3612e33565b5060a06040516117f28161282e565b60ff60025473ffffffffffffffffffffffffffffffffffffffff81168352831c161515602082015273ffffffffffffffffffffffffffffffffffffffff60035416604082015273ffffffffffffffffffffffffffffffffffffffff60045416606082015273ffffffffffffffffffffffffffffffffffffffff6005541660808201526118c7604051809273ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565bf35b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757604060609167ffffffffffffffff61190f612964565b1681526006602052205473ffffffffffffffffffffffffffffffffffffffff6040519167ffffffffffffffff8116835260ff8160401c161515602084015260481c166040820152f35b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475760043567ffffffffffffffff811161112d576119a8903690600401612a21565b9073ffffffffffffffffffffffffffffffffffffffff6004541690835b83811015611c115773ffffffffffffffffffffffffffffffffffffffff6119f08260051b8401612e12565b1690604051917f70a08231000000000000000000000000000000000000000000000000000000008352306004840152602083602481845afa928315611c06578793611bd3575b5082611a48575b5060019150016119c5565b8460405193611ae460208601957fa9059cbb00000000000000000000000000000000000000000000000000000000875283602482015282604482015260448152611a93606482612883565b8a80604098895193611aa58b86612883565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082895af1611add6132e5565b90866138fc565b805180611b20575b505060207f508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e9160019651908152a338611a3d565b819294959693509060209181010312610efe576020611b3f9101612a9c565b15611b505792919085903880611aec565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b9092506020813d8211611bfe575b81611bee60209383612883565b81010312610f0657519138611a36565b3d9150611be1565b6040513d89823e3d90fd5b8480f35b50346101475760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757611c4d612964565b506024359073ffffffffffffffffffffffffffffffffffffffff82168203610147576020611c7a83612d0f565b73ffffffffffffffffffffffffffffffffffffffff60405191168152f35b50346101475760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757604051611cd48161282e565b611cdc6129ba565b81526024358015158103610dfa576020820190815260443573ffffffffffffffffffffffffffffffffffffffff811681036114aa5760408301908152611d206129dd565b90606084019182526084359273ffffffffffffffffffffffffffffffffffffffff84168403610f025760808501938452611d5861329a565b73ffffffffffffffffffffffffffffffffffffffff85511615801561207f575b8015612075575b61204d579273ffffffffffffffffffffffffffffffffffffffff859381809461012097827fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f19a51167fffffffffffffffffffffffff000000000000000000000000000000000000000060025416176002555115157fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006002549260a01b1691161760025551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600354161760035551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600454161760045551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600554161760055561204960405191611ec9836127e3565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016835273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604084015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166060840152611ff9604051809473ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b608083019073ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565ba180f35b6004867f35be3ac8000000000000000000000000000000000000000000000000000000008152fd5b5080511515611d7f565b5073ffffffffffffffffffffffffffffffffffffffff83511615611d78565b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610147576004359067ffffffffffffffff821161014757366023830112156101475781600401356120fa81612995565b926121086040519485612883565b81845260246060602086019302820101903682116114aa57602401915b8183106122a65750505061213761329a565b805b82518110156122a25761214c8184612ec3565b5167ffffffffffffffff6121608386612ec3565b51511690811561227657907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5606060019493838752600660205260ff60408820612237604060208501519483547fffffff0000000000000000000000000000000000000000ffffffffffffffffff7cffffffffffffffffffffffffffffffffffffffff0000000000000000008860481b1691161784550151151582907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff68ff0000000000000000835492151560401b169116179055565b5473ffffffffffffffffffffffffffffffffffffffff6040519367ffffffffffffffff8316855216602084015260401c1615156040820152a201612139565b602484837fc35aa79d000000000000000000000000000000000000000000000000000000008252600452fd5b5080f35b6060833603126114aa576040516060810181811067ffffffffffffffff821117612321576040526122d684612980565b8152602084013573ffffffffffffffffffffffffffffffffffffffff81168103610f02579181606093602080940152612311604087016129ad565b6040820152815201920191612125565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50346101475760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757612386612964565b60243567ffffffffffffffff8111610dfa5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8236030112610dfa576040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008360801b16600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa90811561257c578491612542575b5061250c576124b99160209173ffffffffffffffffffffffffffffffffffffffff60025416906040518095819482937fd8694ccd0000000000000000000000000000000000000000000000000000000084526004019060048401612b38565b03915afa908115610f0a5782916124d6575b602082604051908152f35b90506020813d602011612504575b816124f160209383612883565b8101031261112d576020915051386124cb565b3d91506124e4565b60248367ffffffffffffffff847ffdbd6a7200000000000000000000000000000000000000000000000000000000835216600452fd5b90506020813d602011612574575b8161255d60209383612883565b810103126114aa5761256e90612a9c565b3861245a565b3d9150612550565b6040513d86823e3d90fd5b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757506115ab6040516125c8604082612883565b601081527f4f6e52616d7020312e362e302d646576000000000000000000000000000000006020820152604051918291602083526020830190612921565b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757602073ffffffffffffffffffffffffffffffffffffffff60075416604051908152f35b90503461112d57817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261112d57806126946060926127e3565b828152826020820152826040820152015260806040516126b3816127e3565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660608201526118c7604051809273ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b6080810190811067ffffffffffffffff8211176127ff57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60a0810190811067ffffffffffffffff8211176127ff57604052565b610120810190811067ffffffffffffffff8211176127ff57604052565b6040810190811067ffffffffffffffff8211176127ff57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176127ff57604052565b67ffffffffffffffff81116127ff57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b8381106129115750506000910152565b8181015183820152602001612901565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f60209361295d815180928187528780880191016128fe565b0116010190565b6004359067ffffffffffffffff8216820361297b57565b600080fd5b359067ffffffffffffffff8216820361297b57565b67ffffffffffffffff81116127ff5760051b60200190565b3590811515820361297b57565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361297b57565b6064359073ffffffffffffffffffffffffffffffffffffffff8216820361297b57565b359073ffffffffffffffffffffffffffffffffffffffff8216820361297b57565b9181601f8401121561297b5782359167ffffffffffffffff831161297b576020808501948460051b01011161297b57565b906020808351928381520192019060005b818110612a705750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101612a63565b5190811515820361297b57565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18236030181121561297b57016020813591019167ffffffffffffffff821161297b57813603831361297b57565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9067ffffffffffffffff9093929316815260406020820152612bae612b71612b608580612aa9565b60a0604086015260e0850191612af9565b612b7e6020860186612aa9565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858403016060860152612af9565b9060408401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18536030181121561297b5784016020813591019267ffffffffffffffff821161297b578160061b3603841361297b578281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0016080840152818152602001929060005b818110612caf57505050612c7c8473ffffffffffffffffffffffffffffffffffffffff612c6c6060612cac979801612a00565b1660a08401526080810190612aa9565b9160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082860301910152612af9565b90565b90919360408060019273ffffffffffffffffffffffffffffffffffffffff612cd689612a00565b16815260208881013590820152019501929101612c39565b519073ffffffffffffffffffffffffffffffffffffffff8216820361297b57565b73ffffffffffffffffffffffffffffffffffffffff604051917fbbe4f6db00000000000000000000000000000000000000000000000000000000835216600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa8015612e0657600090612db9575b73ffffffffffffffffffffffffffffffffffffffff91501690565b506020813d602011612dfe575b81612dd360209383612883565b8101031261297b57612df973ffffffffffffffffffffffffffffffffffffffff91612cee565b612d9e565b3d9150612dc6565b6040513d6000823e3d90fd5b3573ffffffffffffffffffffffffffffffffffffffff8116810361297b5790565b60405190612e408261282e565b60006080838281528260208201528260408201528260608201520152565b9080601f8301121561297b578135612e7581612995565b92612e836040519485612883565b81845260208085019260051b82010192831161297b57602001905b828210612eab5750505090565b60208091612eb884612a00565b815201910190612e9e565b8051821015612ed75760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561297b570180359067ffffffffffffffff821161297b5760200191813603831361297b57565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561297b570180359067ffffffffffffffff821161297b57602001918160061b3603831361297b57565b60405190612fb88261282e565b6060608083600081528260208201528260408201526000838201520152565b929192612fe3826128c4565b91612ff16040519384612883565b82948184528183011161297b578281602093846000960137010152565b81601f8201121561297b578051613024816128c4565b926130326040519485612883565b8184526020828401011161297b57612cac91602080850191016128fe565b9080602083519182815201916020808360051b8301019401926000915b83831061307c57505050505090565b9091929394602080613117837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301875289519073ffffffffffffffffffffffffffffffffffffffff825116815260806130fc6130ea8685015160a08886015260a0850190612921565b60408501518482036040860152612921565b92606081015160608401520151906080818403910152612921565b9701930193019193929061306d565b519067ffffffffffffffff8216820361297b57565b90612cac916020815267ffffffffffffffff6080835180516020850152826020820151166040850152826040820151166060850152826060820151168285015201511660a082015273ffffffffffffffffffffffffffffffffffffffff60208301511660c082015261010061322f6131fa6131c760408601516101a060e08701526101c0860190612921565b60608601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08683030185870152612921565b60808501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085830301610120860152612921565b9273ffffffffffffffffffffffffffffffffffffffff60a08201511661014084015260c081015161016084015260e08101516101808401520151906101a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082850301910152613050565b73ffffffffffffffffffffffffffffffffffffffff6001541633036132bb57565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b3d15613310573d906132f6826128c4565b916133046040519384612883565b82523d6000602084013e565b606090565b60006133b681926040516133288161284a565b613330612e33565b815283602082015260606040820152606080820152606060808201528360a08201528360c08201528360e082015260606101008201525073ffffffffffffffffffffffffffffffffffffffff60075416906040519485809481937f8a06fadb0000000000000000000000000000000000000000000000000000000083526004830161313b565b03925af18091600091613413575b5090612cac5761340f6133d56132e5565b6040519182917f828ebdfb000000000000000000000000000000000000000000000000000000008352602060048401526024830190612921565b0390fd5b3d8083833e6134228183612883565b810190602081830312610dfa5780519067ffffffffffffffff82116114aa570191828203926101a0841261112d5760a06040519461345f8661284a565b1261112d576040516134708161282e565b8151815261348060208301613126565b602082015261349160408301613126565b60408201526134a260608301613126565b60608201526134b360808301613126565b608082015284526134c660a08201612cee565b602085015260c081015167ffffffffffffffff8111610dfa57836134eb91830161300e565b604085015260e081015167ffffffffffffffff8111610dfa578361351091830161300e565b606085015261010081015167ffffffffffffffff8111610dfa578361353691830161300e565b60808501526135486101208201612cee565b60a085015261014081015160c085015261016081015160e08501526101808101519067ffffffffffffffff8211610dfa570182601f8201121561112d5780519161359183612995565b9361359f6040519586612883565b83855260208086019460051b84010192818411610dfa5760208101945b8486106135d557505050505050610100820152386133c4565b855167ffffffffffffffff811161104357820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828603011261104357604051906136218261282e565b61362d60208201612cee565b8252604081015167ffffffffffffffff8111610f06578560206136529284010161300e565b6020830152606081015167ffffffffffffffff8111610f065785602061367a9284010161300e565b60408301526080810151606083015260a08101519067ffffffffffffffff8211610f0657916136b18660208096948196010161300e565b60808201528152019501946135bc565b8054821015612ed75760005260206000200190600090565b600082815260018201602052604090205461376357805490680100000000000000008210156127ff578261374c6137178460018096018555846136c1565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905580549260005201602052604060002055600190565b5050600090565b90600182019181600052826020526040600020548015156000146138f3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116138c4578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116138c45781810361388d575b5050508054801561385e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019061381f82826136c1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6138ad61389d61371793866136c1565b90549060031b1c928392866136c1565b9055600052836020526040600020553880806137e7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50505050600090565b919290156139775750815115613910575090565b3b156139195790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b82519091501561398a5750805190602001fd5b61340f906040519182917f08c379a000000000000000000000000000000000000000000000000000000000835260206004840152602483019061292156fea164736f6c634300081a000a",
}

var OnRampWithMessageTransformerABI = OnRampWithMessageTransformerMetaData.ABI

var OnRampWithMessageTransformerBin = OnRampWithMessageTransformerMetaData.Bin

func DeployOnRampWithMessageTransformer(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OnRampStaticConfig, dynamicConfig OnRampDynamicConfig, destChainConfigs []OnRampDestChainConfigArgs, messageTransformerAddr common.Address) (common.Address, *types.Transaction, *OnRampWithMessageTransformer, error) {
	parsed, err := OnRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OnRampWithMessageTransformerBin), backend, staticConfig, dynamicConfig, destChainConfigs, messageTransformerAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OnRampWithMessageTransformer{address: address, abi: *parsed, OnRampWithMessageTransformerCaller: OnRampWithMessageTransformerCaller{contract: contract}, OnRampWithMessageTransformerTransactor: OnRampWithMessageTransformerTransactor{contract: contract}, OnRampWithMessageTransformerFilterer: OnRampWithMessageTransformerFilterer{contract: contract}}, nil
}

type OnRampWithMessageTransformer struct {
	address common.Address
	abi     abi.ABI
	OnRampWithMessageTransformerCaller
	OnRampWithMessageTransformerTransactor
	OnRampWithMessageTransformerFilterer
}

type OnRampWithMessageTransformerCaller struct {
	contract *bind.BoundContract
}

type OnRampWithMessageTransformerTransactor struct {
	contract *bind.BoundContract
}

type OnRampWithMessageTransformerFilterer struct {
	contract *bind.BoundContract
}

type OnRampWithMessageTransformerSession struct {
	Contract     *OnRampWithMessageTransformer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OnRampWithMessageTransformerCallerSession struct {
	Contract *OnRampWithMessageTransformerCaller
	CallOpts bind.CallOpts
}

type OnRampWithMessageTransformerTransactorSession struct {
	Contract     *OnRampWithMessageTransformerTransactor
	TransactOpts bind.TransactOpts
}

type OnRampWithMessageTransformerRaw struct {
	Contract *OnRampWithMessageTransformer
}

type OnRampWithMessageTransformerCallerRaw struct {
	Contract *OnRampWithMessageTransformerCaller
}

type OnRampWithMessageTransformerTransactorRaw struct {
	Contract *OnRampWithMessageTransformerTransactor
}

func NewOnRampWithMessageTransformer(address common.Address, backend bind.ContractBackend) (*OnRampWithMessageTransformer, error) {
	abi, err := abi.JSON(strings.NewReader(OnRampWithMessageTransformerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOnRampWithMessageTransformer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformer{address: address, abi: abi, OnRampWithMessageTransformerCaller: OnRampWithMessageTransformerCaller{contract: contract}, OnRampWithMessageTransformerTransactor: OnRampWithMessageTransformerTransactor{contract: contract}, OnRampWithMessageTransformerFilterer: OnRampWithMessageTransformerFilterer{contract: contract}}, nil
}

func NewOnRampWithMessageTransformerCaller(address common.Address, caller bind.ContractCaller) (*OnRampWithMessageTransformerCaller, error) {
	contract, err := bindOnRampWithMessageTransformer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerCaller{contract: contract}, nil
}

func NewOnRampWithMessageTransformerTransactor(address common.Address, transactor bind.ContractTransactor) (*OnRampWithMessageTransformerTransactor, error) {
	contract, err := bindOnRampWithMessageTransformer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerTransactor{contract: contract}, nil
}

func NewOnRampWithMessageTransformerFilterer(address common.Address, filterer bind.ContractFilterer) (*OnRampWithMessageTransformerFilterer, error) {
	contract, err := bindOnRampWithMessageTransformer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerFilterer{contract: contract}, nil
}

func bindOnRampWithMessageTransformer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OnRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRampWithMessageTransformer.Contract.OnRampWithMessageTransformerCaller.contract.Call(opts, result, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.OnRampWithMessageTransformerTransactor.contract.Transfer(opts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.OnRampWithMessageTransformerTransactor.contract.Transact(opts, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OnRampWithMessageTransformer.Contract.contract.Call(opts, result, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.contract.Transfer(opts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.contract.Transact(opts, method, params...)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

	error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getAllowedSendersList", destChainSelector)

	outstruct := new(GetAllowedSendersList)
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsEnabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfiguredAddresses = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetAllowedSendersList(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetAllowedSendersList(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

	error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	outstruct := new(GetDestChainConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequenceNumber = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.AllowlistEnabled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Router = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetDestChainConfig(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _OnRampWithMessageTransformer.Contract.GetDestChainConfig(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OnRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampDynamicConfig)).(*OnRampDynamicConfig)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetDynamicConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetDynamicConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRampWithMessageTransformer.Contract.GetExpectedNextSequenceNumber(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _OnRampWithMessageTransformer.Contract.GetExpectedNextSequenceNumber(&_OnRampWithMessageTransformer.CallOpts, destChainSelector)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRampWithMessageTransformer.Contract.GetFee(&_OnRampWithMessageTransformer.CallOpts, destChainSelector, message)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _OnRampWithMessageTransformer.Contract.GetFee(&_OnRampWithMessageTransformer.CallOpts, destChainSelector, message)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetMessageTransformer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getMessageTransformer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetMessageTransformer() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetMessageTransformer(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetMessageTransformer() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetMessageTransformer(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getPoolBySourceToken", arg0, sourceToken)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetPoolBySourceToken(&_OnRampWithMessageTransformer.CallOpts, arg0, sourceToken)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetPoolBySourceToken(&_OnRampWithMessageTransformer.CallOpts, arg0, sourceToken)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetStaticConfig(opts *bind.CallOpts) (OnRampStaticConfig, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OnRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampStaticConfig)).(*OnRampStaticConfig)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetStaticConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _OnRampWithMessageTransformer.Contract.GetStaticConfig(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "getSupportedTokens", arg0)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetSupportedTokens(&_OnRampWithMessageTransformer.CallOpts, arg0)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.GetSupportedTokens(&_OnRampWithMessageTransformer.CallOpts, arg0)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) Owner() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.Owner(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) Owner() (common.Address, error) {
	return _OnRampWithMessageTransformer.Contract.Owner(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OnRampWithMessageTransformer.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) TypeAndVersion() (string, error) {
	return _OnRampWithMessageTransformer.Contract.TypeAndVersion(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerCallerSession) TypeAndVersion() (string, error) {
	return _OnRampWithMessageTransformer.Contract.TypeAndVersion(&_OnRampWithMessageTransformer.CallOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "acceptOwnership")
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.AcceptOwnership(&_OnRampWithMessageTransformer.TransactOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.AcceptOwnership(&_OnRampWithMessageTransformer.TransactOpts)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) ApplyAllowlistUpdates(opts *bind.TransactOpts, allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "applyAllowlistUpdates", allowlistConfigArgsItems)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyAllowlistUpdates(&_OnRampWithMessageTransformer.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyAllowlistUpdates(&_OnRampWithMessageTransformer.TransactOpts, allowlistConfigArgsItems)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyDestChainConfigUpdates(&_OnRampWithMessageTransformer.TransactOpts, destChainConfigArgs)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ApplyDestChainConfigUpdates(&_OnRampWithMessageTransformer.TransactOpts, destChainConfigArgs)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "forwardFromRouter", destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ForwardFromRouter(&_OnRampWithMessageTransformer.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.ForwardFromRouter(&_OnRampWithMessageTransformer.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.SetDynamicConfig(&_OnRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.SetDynamicConfig(&_OnRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "transferOwnership", to)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.TransferOwnership(&_OnRampWithMessageTransformer.TransactOpts, to)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.TransferOwnership(&_OnRampWithMessageTransformer.TransactOpts, to)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactor) WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.contract.Transact(opts, "withdrawFeeTokens", feeTokens)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.WithdrawFeeTokens(&_OnRampWithMessageTransformer.TransactOpts, feeTokens)
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerTransactorSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _OnRampWithMessageTransformer.Contract.WithdrawFeeTokens(&_OnRampWithMessageTransformer.TransactOpts, feeTokens)
}

type OnRampWithMessageTransformerAllowListAdminSetIterator struct {
	Event *OnRampWithMessageTransformerAllowListAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerAllowListAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerAllowListAdminSet)
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
		it.Event = new(OnRampWithMessageTransformerAllowListAdminSet)
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

func (it *OnRampWithMessageTransformerAllowListAdminSetIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerAllowListAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerAllowListAdminSet struct {
	AllowlistAdmin common.Address
	Raw            types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampWithMessageTransformerAllowListAdminSetIterator, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerAllowListAdminSetIterator{contract: _OnRampWithMessageTransformer.contract, event: "AllowListAdminSet", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerAllowListAdminSet)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseAllowListAdminSet(log types.Log) (*OnRampWithMessageTransformerAllowListAdminSet, error) {
	event := new(OnRampWithMessageTransformerAllowListAdminSet)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerAllowListSendersAddedIterator struct {
	Event *OnRampWithMessageTransformerAllowListSendersAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerAllowListSendersAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerAllowListSendersAdded)
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
		it.Event = new(OnRampWithMessageTransformerAllowListSendersAdded)
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

func (it *OnRampWithMessageTransformerAllowListSendersAddedIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerAllowListSendersAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerAllowListSendersAdded struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerAllowListSendersAddedIterator{contract: _OnRampWithMessageTransformer.contract, event: "AllowListSendersAdded", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerAllowListSendersAdded)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseAllowListSendersAdded(log types.Log) (*OnRampWithMessageTransformerAllowListSendersAdded, error) {
	event := new(OnRampWithMessageTransformerAllowListSendersAdded)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerAllowListSendersRemovedIterator struct {
	Event *OnRampWithMessageTransformerAllowListSendersRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerAllowListSendersRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerAllowListSendersRemoved)
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
		it.Event = new(OnRampWithMessageTransformerAllowListSendersRemoved)
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

func (it *OnRampWithMessageTransformerAllowListSendersRemovedIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerAllowListSendersRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerAllowListSendersRemoved struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersRemovedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerAllowListSendersRemovedIterator{contract: _OnRampWithMessageTransformer.contract, event: "AllowListSendersRemoved", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerAllowListSendersRemoved)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseAllowListSendersRemoved(log types.Log) (*OnRampWithMessageTransformerAllowListSendersRemoved, error) {
	event := new(OnRampWithMessageTransformerAllowListSendersRemoved)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerCCIPMessageSentIterator struct {
	Event *OnRampWithMessageTransformerCCIPMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerCCIPMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerCCIPMessageSent)
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
		it.Event = new(OnRampWithMessageTransformerCCIPMessageSent)
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

func (it *OnRampWithMessageTransformerCCIPMessageSentIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerCCIPMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampWithMessageTransformerCCIPMessageSentIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerCCIPMessageSentIterator{contract: _OnRampWithMessageTransformer.contract, event: "CCIPMessageSent", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerCCIPMessageSent)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseCCIPMessageSent(log types.Log) (*OnRampWithMessageTransformerCCIPMessageSent, error) {
	event := new(OnRampWithMessageTransformerCCIPMessageSent)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerConfigSetIterator struct {
	Event *OnRampWithMessageTransformerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerConfigSet)
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
		it.Event = new(OnRampWithMessageTransformerConfigSet)
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

func (it *OnRampWithMessageTransformerConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerConfigSet struct {
	StaticConfig  OnRampStaticConfig
	DynamicConfig OnRampDynamicConfig
	Raw           types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OnRampWithMessageTransformerConfigSetIterator, error) {

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerConfigSetIterator{contract: _OnRampWithMessageTransformer.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerConfigSet) (event.Subscription, error) {

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerConfigSet)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseConfigSet(log types.Log) (*OnRampWithMessageTransformerConfigSet, error) {
	event := new(OnRampWithMessageTransformerConfigSet)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerDestChainConfigSetIterator struct {
	Event *OnRampWithMessageTransformerDestChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerDestChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerDestChainConfigSet)
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
		it.Event = new(OnRampWithMessageTransformerDestChainConfigSet)
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

func (it *OnRampWithMessageTransformerDestChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerDestChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerDestChainConfigSet struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Router            common.Address
	AllowlistEnabled  bool
	Raw               types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerDestChainConfigSetIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerDestChainConfigSetIterator{contract: _OnRampWithMessageTransformer.contract, event: "DestChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerDestChainConfigSet)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseDestChainConfigSet(log types.Log) (*OnRampWithMessageTransformerDestChainConfigSet, error) {
	event := new(OnRampWithMessageTransformerDestChainConfigSet)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerFeeTokenWithdrawnIterator struct {
	Event *OnRampWithMessageTransformerFeeTokenWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerFeeTokenWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerFeeTokenWithdrawn)
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
		it.Event = new(OnRampWithMessageTransformerFeeTokenWithdrawn)
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

func (it *OnRampWithMessageTransformerFeeTokenWithdrawnIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerFeeTokenWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerFeeTokenWithdrawn struct {
	FeeAggregator common.Address
	FeeToken      common.Address
	Amount        *big.Int
	Raw           types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampWithMessageTransformerFeeTokenWithdrawnIterator, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerFeeTokenWithdrawnIterator{contract: _OnRampWithMessageTransformer.contract, event: "FeeTokenWithdrawn", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerFeeTokenWithdrawn)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseFeeTokenWithdrawn(log types.Log) (*OnRampWithMessageTransformerFeeTokenWithdrawn, error) {
	event := new(OnRampWithMessageTransformerFeeTokenWithdrawn)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerOwnershipTransferRequestedIterator struct {
	Event *OnRampWithMessageTransformerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerOwnershipTransferRequested)
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
		it.Event = new(OnRampWithMessageTransformerOwnershipTransferRequested)
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

func (it *OnRampWithMessageTransformerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerOwnershipTransferRequestedIterator{contract: _OnRampWithMessageTransformer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerOwnershipTransferRequested)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseOwnershipTransferRequested(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferRequested, error) {
	event := new(OnRampWithMessageTransformerOwnershipTransferRequested)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OnRampWithMessageTransformerOwnershipTransferredIterator struct {
	Event *OnRampWithMessageTransformerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OnRampWithMessageTransformerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnRampWithMessageTransformerOwnershipTransferred)
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
		it.Event = new(OnRampWithMessageTransformerOwnershipTransferred)
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

func (it *OnRampWithMessageTransformerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OnRampWithMessageTransformerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OnRampWithMessageTransformerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OnRampWithMessageTransformerOwnershipTransferredIterator{contract: _OnRampWithMessageTransformer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OnRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OnRampWithMessageTransformerOwnershipTransferred)
				if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformerFilterer) ParseOwnershipTransferred(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferred, error) {
	event := new(OnRampWithMessageTransformerOwnershipTransferred)
	if err := _OnRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OnRampWithMessageTransformer.abi.Events["AllowListAdminSet"].ID:
		return _OnRampWithMessageTransformer.ParseAllowListAdminSet(log)
	case _OnRampWithMessageTransformer.abi.Events["AllowListSendersAdded"].ID:
		return _OnRampWithMessageTransformer.ParseAllowListSendersAdded(log)
	case _OnRampWithMessageTransformer.abi.Events["AllowListSendersRemoved"].ID:
		return _OnRampWithMessageTransformer.ParseAllowListSendersRemoved(log)
	case _OnRampWithMessageTransformer.abi.Events["CCIPMessageSent"].ID:
		return _OnRampWithMessageTransformer.ParseCCIPMessageSent(log)
	case _OnRampWithMessageTransformer.abi.Events["ConfigSet"].ID:
		return _OnRampWithMessageTransformer.ParseConfigSet(log)
	case _OnRampWithMessageTransformer.abi.Events["DestChainConfigSet"].ID:
		return _OnRampWithMessageTransformer.ParseDestChainConfigSet(log)
	case _OnRampWithMessageTransformer.abi.Events["FeeTokenWithdrawn"].ID:
		return _OnRampWithMessageTransformer.ParseFeeTokenWithdrawn(log)
	case _OnRampWithMessageTransformer.abi.Events["OwnershipTransferRequested"].ID:
		return _OnRampWithMessageTransformer.ParseOwnershipTransferRequested(log)
	case _OnRampWithMessageTransformer.abi.Events["OwnershipTransferred"].ID:
		return _OnRampWithMessageTransformer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OnRampWithMessageTransformerAllowListAdminSet) Topic() common.Hash {
	return common.HexToHash("0xb8c9b44ae5b5e3afb195f67391d9ff50cb904f9c0fa5fd520e497a97c1aa5a1e")
}

func (OnRampWithMessageTransformerAllowListSendersAdded) Topic() common.Hash {
	return common.HexToHash("0x330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc3281")
}

func (OnRampWithMessageTransformerAllowListSendersRemoved) Topic() common.Hash {
	return common.HexToHash("0xc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d421586")
}

func (OnRampWithMessageTransformerCCIPMessageSent) Topic() common.Hash {
	return common.HexToHash("0x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32")
}

func (OnRampWithMessageTransformerConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1")
}

func (OnRampWithMessageTransformerDestChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0xd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5")
}

func (OnRampWithMessageTransformerFeeTokenWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e")
}

func (OnRampWithMessageTransformerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OnRampWithMessageTransformerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OnRampWithMessageTransformer *OnRampWithMessageTransformer) Address() common.Address {
	return _OnRampWithMessageTransformer.address
}

type OnRampWithMessageTransformerInterface interface {
	GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

		error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

		error)

	GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error)

	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetMessageTransformer(opts *bind.CallOpts) (common.Address, error)

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

	FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*OnRampWithMessageTransformerAllowListAdminSetIterator, error)

	WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error)

	ParseAllowListAdminSet(log types.Log) (*OnRampWithMessageTransformerAllowListAdminSet, error)

	FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersAddedIterator, error)

	WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersAdded(log types.Log) (*OnRampWithMessageTransformerAllowListSendersAdded, error)

	FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerAllowListSendersRemovedIterator, error)

	WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersRemoved(log types.Log) (*OnRampWithMessageTransformerAllowListSendersRemoved, error)

	FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*OnRampWithMessageTransformerCCIPMessageSentIterator, error)

	WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error)

	ParseCCIPMessageSent(log types.Log) (*OnRampWithMessageTransformerCCIPMessageSent, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OnRampWithMessageTransformerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OnRampWithMessageTransformerConfigSet, error)

	FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*OnRampWithMessageTransformerDestChainConfigSetIterator, error)

	WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigSet(log types.Log) (*OnRampWithMessageTransformerDestChainConfigSet, error)

	FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*OnRampWithMessageTransformerFeeTokenWithdrawnIterator, error)

	WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenWithdrawn(log types.Log) (*OnRampWithMessageTransformerFeeTokenWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OnRampWithMessageTransformerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OnRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OnRampWithMessageTransformerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
