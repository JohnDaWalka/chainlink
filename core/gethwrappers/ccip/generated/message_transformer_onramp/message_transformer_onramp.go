// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package message_transformer_onramp

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

var MessageTransformerOnRampMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"destChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowlistUpdates\",\"inputs\":[{\"name\":\"allowlistConfigArgsItems\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.AllowlistConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"addedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"removedAllowlistedSenders\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyDestChainConfigUpdates\",\"inputs\":[{\"name\":\"destChainConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structOnRamp.DestChainConfigArgs[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forwardFromRouter\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllowedSendersList\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configuredAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExpectedNextSequenceNumber\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageTransformerAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPoolBySourceToken\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPoolV1\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedTokens\",\"inputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFeeTokens\",\"inputs\":[{\"name\":\"feeTokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdminSet\",\"inputs\":[{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersAdded\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListSendersRemoved\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"senders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CCIPMessageSent\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeValueJuels\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOnRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reentrancyGuardEntered\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowlistAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DestChainConfigSet\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"router\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIRouter\"},{\"name\":\"allowlistEnabled\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTokenWithdrawn\",\"inputs\":[{\"name\":\"feeAggregator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feeToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotSendZeroTokens\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"GetSupportedTokensFunctionalityRemovedCheckAdminRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAllowListRequest\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDestChainConfig\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageTransformError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeCalledByRouter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAllowlistAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RouterMustSetOriginalSender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnsupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610100604052346105a15761401e8038038061001a816105db565b92833981019080820361016081126105a157608081126105a15761003c6105bc565b9161004681610600565b835260208101516001600160a01b03811681036105a1576020840190815261007060408301610614565b916040850192835260a061008660608301610614565b6060870190815294607f1901126105a15760405160a081016001600160401b038111828210176105a6576040526100bf60808301610614565b81526100cd60a08301610628565b602082019081526100e060c08401610614565b90604083019182526100f460e08501610614565b92606081019384526101096101008601610614565b608082019081526101208601519095906001600160401b0381116105a15781018b601f820112156105a1578051906001600160401b0382116105a6578160051b602001610155906105db565b9c8d838152602001926060028201602001918183116105a157602001925b828410610533575050505061014061018b9101610614565b98331561052257600180546001600160a01b0319163317905580516001600160401b0316158015610510575b80156104fe575b80156104ec575b6104bf57516001600160401b0316608081905295516001600160a01b0390811660a08190529751811660c08190529851811660e081905282519091161580156104da575b80156104d0575b6104bf57815160028054855160ff60a01b90151560a01b166001600160a01b039384166001600160a81b0319909216919091171790558451600380549183166001600160a01b03199283161790558651600480549184169183169190911790558751600580549190931691161790557fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1986101209860606102af6105bc565b8a8152602080820193845260408083019586529290910194855281519a8b5291516001600160a01b03908116928b019290925291518116918901919091529051811660608801529051811660808701529051151560a08601529051811660c08501529051811660e0840152905116610100820152a160005b82518110156104185761033a8184610635565b516001600160401b0361034d8386610635565b5151169081156104035760008281526006602090815260409182902081840151815494840151600160401b600160e81b03198616604883901b600160481b600160e81b031617901515851b68ff000000000000000016179182905583516001600160401b0390951685526001600160a01b031691840191909152811c60ff1615159082015260019291907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef590606090a201610327565b5063c35aa79d60e01b60005260045260246000fd5b506001600160a01b031680156104ae57600780546001600160a01b0319169190911790556040516139be908161066082396080518181816103ee01528181610c15015281816122690152612a53015260a0518181816122a2015281816127bf0152612a8c015260c05181818161113f015281816122de0152612ac8015260e05181818161231a01528181612b0401526131030152f35b6342bcdf7f60e11b60005260046000fd5b6306b7c75960e31b60005260046000fd5b5082511515610210565b5084516001600160a01b031615610209565b5088516001600160a01b0316156101c5565b5087516001600160a01b0316156101be565b5086516001600160a01b0316156101b7565b639b15e16f60e01b60005260046000fd5b6060848303126105a15760405190606082016001600160401b038111838210176105a65760405261056385610600565b82526020850151906001600160a01b03821682036105a1578260209283606095015261059160408801610628565b6040820152815201930192610173565b600080fd5b634e487b7160e01b600052604160045260246000fd5b60405190608082016001600160401b038111838210176105a657604052565b6040519190601f01601f191682016001600160401b038111838210176105a657604052565b51906001600160401b03821682036105a157565b51906001600160a01b03821682036105a157565b519081151582036105a157565b80518210156106495760209160051b010190565b634e487b7160e01b600052603260045260246000fdfe608080604052600436101561001357600080fd5b600090813560e01c90816306285c69146129ed575080631056edee1461299b578063181f5a771461291c57806320487ded146126e35780632716072b1461243357806327e936f11461202d57806348a98aa414611faa5780635cb80c5d14611ced5780636def4ce714611c5e5780637437ff9f14611b4157806379ba509714611a5c5780638da5cb5b14611a0a5780639041be3d1461195d578063972b46121461188f578063c9b146b3146114c6578063df0aa9e914610237578063f2fde38b1461014a5763fbca3b74146100e757600080fd5b346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757600490610121612cf9565b507f9e7177c8000000000000000000000000000000000000000000000000000000008152fd5b80fd5b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475773ffffffffffffffffffffffffffffffffffffffff610197612d4f565b61019f61362f565b1633811461020f57807fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788380a380f35b6004827fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b50346101475760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475761026f612cf9565b67ffffffffffffffff602435116114c25760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc602435360301126114c2576102b6612d72565b60025460ff8160a01c1661149a577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001760025567ffffffffffffffff8216835260066020526040832073ffffffffffffffffffffffffffffffffffffffff82161561147257805460ff8160401c16611404575b60481c73ffffffffffffffffffffffffffffffffffffffff1633036113dc5773ffffffffffffffffffffffffffffffffffffffff6003541680611368575b50805467ffffffffffffffff811667ffffffffffffffff811461133b579067ffffffffffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009493011692839116179055604051906103e082612bc3565b84825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602083015267ffffffffffffffff841660408301526060820152836080820152610441602480350160243560040161329b565b6104506004602435018061329b565b61045e6064602435016131a7565b936104736044602435016024356004016132ec565b9490507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06104b96104a387612d2a565b966104b16040519889612c18565b808852612d2a565b018a5b81811061132457505061050d93929161050191604051986104dc8a612bdf565b895273ffffffffffffffffffffffffffffffffffffffff8a1660208a0152369161336c565b6040870152369161336c565b606084015273ffffffffffffffffffffffffffffffffffffffff6020926040516105378582612c18565b88815260808601521660a084015260443560c08401528560e084015261010083015261056d6044602435016024356004016132ec565b61057981969296612d2a565b906105876040519283612c18565b8082528382018097368360061b8201116113205780915b8360061b820183106112e95750505050865b6105c46044602435016024356004016132ec565b905081101561093a576105d78183613258565b51906105f26105eb6004602435018061329b565b369161336c565b916105fb613340565b5085810151156109125773ffffffffffffffffffffffffffffffffffffffff610626818351166130a4565b169283158015610868575b61082557908a6106de92888873ffffffffffffffffffffffffffffffffffffffff8d838701518280895116926040519761066a89612bc3565b885267ffffffffffffffff87890196168652816040890191168152606088019283526080880193845267ffffffffffffffff6040519b8c998a997f9a4575b9000000000000000000000000000000000000000000000000000000008b5260048b01525160a060248b015260c48a0190612cb6565b965116604488015251166064860152516084850152511660a4830152038183885af191821561081a578b92610773575b50600193828492898061076c9651930151910151916040519361073085612bc3565b84528a8401526040830152606082015260405161074d8982612c18565b8c81526080820152610100890151906107668383613258565b52613258565b50016105b0565b91503d808c843e6107848184612c18565b82019187818403126108125780519067ffffffffffffffff821161081657019360408584031261081257604051916107bb83612bfc565b855167ffffffffffffffff811161080e57846107d89188016133a3565b8352888601519267ffffffffffffffff841161080e5761080061076c958795600199016133a3565b8a820152935091509361070e565b8d80fd5b8b80fd5b8c80fd5b6040513d8d823e3d90fd5b60248b73ffffffffffffffffffffffffffffffffffffffff8451167fbf16aab6000000000000000000000000000000000000000000000000000000008252600452fd5b506040517f01ffc9a70000000000000000000000000000000000000000000000000000000081527faff2afbf0000000000000000000000000000000000000000000000000000000060048201528781602481885afa908115610907578c916108d2575b5015610631565b90508781813d8311610900575b6108e98183612c18565b81010312610812576108fa90612e31565b386108cb565b503d6108df565b6040513d8e823e3d90fd5b60048a7f5cf04449000000000000000000000000000000000000000000000000000000008152fd5b50957f430d138c0000000000000000000000000000000000000000000000000000000081939497839773ffffffffffffffffffffffffffffffffffffffff600254169187610a2c8c6109fc6109936064602435016131a7565b9161010073ffffffffffffffffffffffffffffffffffffffff6109c060846024350160243560040161329b565b92909301519467ffffffffffffffff6040519e8f9d8e521660048d01521660248b015260443560448b015260c060648b015260c48a0191612e8e565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8883030160848901526133e5565b917ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8684030160a48701525191828152019190855b898282106112ae575050505082809103915afa9384156112a357828390849385976111b9575b5060e0890152156110d05750815b67ffffffffffffffff6080885101911690526080860152805b61010086015151811015610ae45780610ac960019286613258565b516080610adb836101008b0151613258565b51015201610aae565b50838180610b8588604051610af881612bdf565b610b006131c8565b8152838882015260606040820152606080820152606060808201528360a08201528360c08201528360e082015260606101008201525073ffffffffffffffffffffffffffffffffffffffff60075416906040519485809481937f8a06fadb000000000000000000000000000000000000000000000000000000008352600483016134d0565b03925af1839181610e17575b50610bdf5783610b9f61367a565b90610bdb6040519283927f828ebdfb00000000000000000000000000000000000000000000000000000000845260048401526024830190612cb6565b0390fd5b91604051848101907f130ac867e79e2789f923760a88743d292acdf7002139a588206e2260f73f7321825267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015267ffffffffffffffff8416606082015230608082015260808152610c5f60a082612c18565b51902073ffffffffffffffffffffffffffffffffffffffff8585015116845167ffffffffffffffff6080816060840151169201511673ffffffffffffffffffffffffffffffffffffffff60a08801511660c088015191604051938a850195865260408501526060840152608083015260a082015260a08152610ce260c082612c18565b51902060608501518681519101206040860151878151910120610100870151604051610d4881610d1c8c8201948d865260408301906133e5565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282612c18565b51902091608088015189815191012093604051958a870197885260408701526060860152608085015260a084015260c083015260e082015260e08152610d9061010082612c18565b51902082515267ffffffffffffffff60608351015116907f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f3267ffffffffffffffff60405192169180610de286826134d0565b0390a37fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff600254166002555151604051908152f35b9091503d8085833e610e298183612c18565b81019085818303126110cc5780519067ffffffffffffffff82116110c4570190818103916101a083126110c45760a060405193610e6585612bdf565b126110c457604051610e7681612bc3565b81518152610e858883016134bb565b88820152610e95604083016134bb565b6040820152610ea6606083016134bb565b6060820152610eb7608083016134bb565b60808201528352610eca60a08201613083565b8784015260c081015167ffffffffffffffff81116110c85782610eee9183016133a3565b604084015260e081015167ffffffffffffffff81116110c85782610f139183016133a3565b606084015261010081015167ffffffffffffffff81116110c85782610f399183016133a3565b6080840152610f4b6101208201613083565b60a084015261014081015160c084015261016081015160e08401526101808101519067ffffffffffffffff82116110c8570181601f820112156110c4578051610f9381612d2a565b92610fa16040519485612c18565b818452888085019260051b840101928184116110c057898101925b848410610fd55750505050506101008201529085610b91565b835167ffffffffffffffff81116110bc57820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082860301126110bc576040519061102182612bc3565b61102c8d8201613083565b8252604081015167ffffffffffffffff811161081657858e611050928401016133a3565b8d830152606081015167ffffffffffffffff811161081657858e611076928401016133a3565b60408301526080810151606083015260a08101519067ffffffffffffffff821161081657916110ac868f809694819601016133a3565b6080820152815201930192610fbc565b8a80fd5b8880fd5b8580fd5b8680fd5b8480fd5b73ffffffffffffffffffffffffffffffffffffffff604051917fea458c0c00000000000000000000000000000000000000000000000000000000835267ffffffffffffffff8816600484015216602482015283816044818673ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af19081156111ae578391611175575b50610a95565b90508381813d83116111a7575b61118c8183612c18565b810103126111a35761119d906134bb565b3861116f565b8280fd5b503d611182565b6040513d85823e3d90fd5b9650505090503d8083863e6111ce8186612c18565b8401906080858303126111a3578451946111e9858201612e31565b95604082015167ffffffffffffffff81116110c4578461120a9184016133a3565b9160608101519067ffffffffffffffff82116110c8570184601f820112156110c457805161123781612d2a565b956112456040519788612c18565b818752888088019260051b840101928184116110c057898101925b8484106112765750505050509590929538610a87565b835167ffffffffffffffff81116110bc578b91611298858480948701016133a3565b815201930192611260565b6040513d84823e3d90fd5b8351805173ffffffffffffffffffffffffffffffffffffffff1686528101518186015289975088965060409094019390920191600101610a61565b6040833603126110bc5786604091825161130281612bfc565b61130b86612d95565b8152828601358382015281520192019161059e565b8980fd5b60209061132f613340565b82828a010152016104bc565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b803b156110cc578460405180927fe0a0e5060000000000000000000000000000000000000000000000000000000082528183816113ae6024356004018b60048401612ecd565b03925af180156113d1571561037e57846113ca91959295612c18565b923861037e565b6040513d87823e3d90fd5b6004847f1c0a3529000000000000000000000000000000000000000000000000000000008152fd5b73ffffffffffffffffffffffffffffffffffffffff831660009081526002830160205260409020546103405760248573ffffffffffffffffffffffffffffffffffffffff857fd0d2597600000000000000000000000000000000000000000000000000000000835216600452fd5b6004847fa4ec7479000000000000000000000000000000000000000000000000000000008152fd5b6004847f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b5080fd5b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475760043567ffffffffffffffff81116114c257611516903690600401612db6565b73ffffffffffffffffffffffffffffffffffffffff600154163303611847575b919081907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8181360301915b84811015611843578060051b820135838112156110cc578201916080833603126110cc576040519461159286612b78565b61159b84612d15565b86526115a960208501612d42565b9660208701978852604085013567ffffffffffffffff81116111a3576115d290369087016131f3565b9460408801958652606081013567ffffffffffffffff811161183f576115fa913691016131f3565b60608801908152875167ffffffffffffffff1683526006602052604080842099518a547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169015159182901b68ff000000000000000016178a55909590815151611717575b5095976001019550815b855180518210156116a857906116a173ffffffffffffffffffffffffffffffffffffffff61169983600195613258565b511689613753565b5001611669565b505095909694506001929193519081516116c8575b505001939293611561565b61170d67ffffffffffffffff7fc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d42158692511692604051918291602083526020830190612de7565b0390a238806116bd565b9893959296919094979860001461180857600184019591875b865180518210156117ad5761175a8273ffffffffffffffffffffffffffffffffffffffff92613258565b51168015611776579061176f6001928a6136c2565b5001611730565b60248a67ffffffffffffffff8e51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b50509692955090929796937f330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc32816117fe67ffffffffffffffff8a51169251604051918291602083526020830190612de7565b0390a2388061165f565b60248767ffffffffffffffff8b51167f463258ff000000000000000000000000000000000000000000000000000000008252600452fd5b8380fd5b8380f35b73ffffffffffffffffffffffffffffffffffffffff60055416330315611536576004837f905d7d9b000000000000000000000000000000000000000000000000000000008152fd5b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475767ffffffffffffffff6118d0612cf9565b16808252600660205260ff604083205460401c16908252600660205260016040832001916040518093849160208254918281520191845260208420935b81811061194457505061192292500383612c18565b61194060405192839215158352604060208401526040830190612de7565b0390f35b845483526001948501948794506020909301920161190d565b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475767ffffffffffffffff61199e612cf9565b1681526006602052600167ffffffffffffffff604083205416019067ffffffffffffffff82116119dd5760208267ffffffffffffffff60405191168152f35b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526011600452fd5b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757805473ffffffffffffffffffffffffffffffffffffffff81163303611b19577fffffffffffffffffffffffff000000000000000000000000000000000000000060015491338284161760015516825573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b6004827f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757611b786131c8565b5060a0604051611b8781612bc3565b60ff60025473ffffffffffffffffffffffffffffffffffffffff81168352831c161515602082015273ffffffffffffffffffffffffffffffffffffffff60035416604082015273ffffffffffffffffffffffffffffffffffffffff60045416606082015273ffffffffffffffffffffffffffffffffffffffff600554166080820152611c5c604051809273ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565bf35b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757604060609167ffffffffffffffff611ca4612cf9565b1681526006602052205473ffffffffffffffffffffffffffffffffffffffff6040519167ffffffffffffffff8116835260ff8160401c161515602084015260481c166040820152f35b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475760043567ffffffffffffffff81116114c257611d3d903690600401612db6565b9073ffffffffffffffffffffffffffffffffffffffff6004541690835b83811015611fa65773ffffffffffffffffffffffffffffffffffffffff611d858260051b84016131a7565b1690604051917f70a08231000000000000000000000000000000000000000000000000000000008352306004840152602083602481845afa928315611f9b578793611f68575b5082611ddd575b506001915001611d5a565b8460405193611e7960208601957fa9059cbb00000000000000000000000000000000000000000000000000000000875283602482015282604482015260448152611e28606482612c18565b8a80604098895193611e3a8b86612c18565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082895af1611e7261367a565b90866138e5565b805180611eb5575b505060207f508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e9160019651908152a338611dd2565b8192949596935090602091810103126110c0576020611ed49101612e31565b15611ee55792919085903880611e81565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b9092506020813d8211611f93575b81611f8360209383612c18565b810103126110c857519138611dcb565b3d9150611f76565b6040513d89823e3d90fd5b8480f35b50346101475760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757611fe2612cf9565b506024359073ffffffffffffffffffffffffffffffffffffffff8216820361014757602061200f836130a4565b73ffffffffffffffffffffffffffffffffffffffff60405191168152f35b50346101475760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475760405161206981612bc3565b612071612d4f565b815260243580151581036111a3576020820190815260443573ffffffffffffffffffffffffffffffffffffffff8116810361183f57604083019081526120b5612d72565b90606084019182526084359273ffffffffffffffffffffffffffffffffffffffff841684036110c457608085019384526120ed61362f565b73ffffffffffffffffffffffffffffffffffffffff855116158015612414575b801561240a575b6123e2579273ffffffffffffffffffffffffffffffffffffffff859381809461012097827fc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f19a51167fffffffffffffffffffffffff000000000000000000000000000000000000000060025416176002555115157fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006002549260a01b1691161760025551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600354161760035551167fffffffffffffffffffffffff0000000000000000000000000000000000000000600454161760045551167fffffffffffffffffffffffff000000000000000000000000000000000000000060055416176005556123de6040519161225e83612b78565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016835273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604084015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016606084015261238e604051809473ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b608083019073ffffffffffffffffffffffffffffffffffffffff60808092828151168552602081015115156020860152826040820151166040860152826060820151166060860152015116910152565ba180f35b6004867f35be3ac8000000000000000000000000000000000000000000000000000000008152fd5b5080511515612114565b5073ffffffffffffffffffffffffffffffffffffffff8351161561210d565b50346101475760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610147576004359067ffffffffffffffff8211610147573660238301121561014757816004013561248f81612d2a565b9261249d6040519485612c18565b818452602460606020860193028201019036821161183f57602401915b81831061263b575050506124cc61362f565b805b8251811015612637576124e18184613258565b5167ffffffffffffffff6124f58386613258565b51511690811561260b57907fd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5606060019493838752600660205260ff604088206125cc604060208501519483547fffffff0000000000000000000000000000000000000000ffffffffffffffffff7cffffffffffffffffffffffffffffffffffffffff0000000000000000008860481b1691161784550151151582907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff68ff0000000000000000835492151560401b169116179055565b5473ffffffffffffffffffffffffffffffffffffffff6040519367ffffffffffffffff8316855216602084015260401c1615156040820152a2016124ce565b602484837fc35aa79d000000000000000000000000000000000000000000000000000000008252600452fd5b5080f35b60608336031261183f576040516060810181811067ffffffffffffffff8211176126b65760405261266b84612d15565b8152602084013573ffffffffffffffffffffffffffffffffffffffff811681036110c45791816060936020809401526126a660408701612d42565b60408201528152019201916124ba565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50346101475760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101475761271b612cf9565b60243567ffffffffffffffff81116111a35760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126111a3576040517f2cbc26bb00000000000000000000000000000000000000000000000000000000815277ffffffffffffffff000000000000000000000000000000008360801b16600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa9081156129115784916128d7575b506128a15761284e9160209173ffffffffffffffffffffffffffffffffffffffff60025416906040518095819482937fd8694ccd0000000000000000000000000000000000000000000000000000000084526004019060048401612ecd565b03915afa9081156112a357829161286b575b602082604051908152f35b90506020813d602011612899575b8161288660209383612c18565b810103126114c257602091505138612860565b3d9150612879565b60248367ffffffffffffffff847ffdbd6a7200000000000000000000000000000000000000000000000000000000835216600452fd5b90506020813d602011612909575b816128f260209383612c18565b8101031261183f5761290390612e31565b386127ef565b3d91506128e5565b6040513d86823e3d90fd5b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610147575061194060405161295d604082612c18565b601081527f4f6e52616d7020312e362e302d646576000000000000000000000000000000006020820152604051918291602083526020830190612cb6565b503461014757807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014757602073ffffffffffffffffffffffffffffffffffffffff60075416604051908152f35b9050346114c257817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126114c25780612a29606092612b78565b82815282602082015282604082015201526080604051612a4881612b78565b67ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602082015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016604082015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166060820152611c5c604051809273ffffffffffffffffffffffffffffffffffffffff6060809267ffffffffffffffff8151168552826020820151166020860152826040820151166040860152015116910152565b6080810190811067ffffffffffffffff821117612b9457604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60a0810190811067ffffffffffffffff821117612b9457604052565b610120810190811067ffffffffffffffff821117612b9457604052565b6040810190811067ffffffffffffffff821117612b9457604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117612b9457604052565b67ffffffffffffffff8111612b9457601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b838110612ca65750506000910152565b8181015183820152602001612c96565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f602093612cf281518092818752878088019101612c93565b0116010190565b6004359067ffffffffffffffff82168203612d1057565b600080fd5b359067ffffffffffffffff82168203612d1057565b67ffffffffffffffff8111612b945760051b60200190565b35908115158203612d1057565b6004359073ffffffffffffffffffffffffffffffffffffffff82168203612d1057565b6064359073ffffffffffffffffffffffffffffffffffffffff82168203612d1057565b359073ffffffffffffffffffffffffffffffffffffffff82168203612d1057565b9181601f84011215612d105782359167ffffffffffffffff8311612d10576020808501948460051b010111612d1057565b906020808351928381520192019060005b818110612e055750505090565b825173ffffffffffffffffffffffffffffffffffffffff16845260209384019390920191600101612df8565b51908115158203612d1057565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215612d1057016020813591019167ffffffffffffffff8211612d10578136038313612d1057565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b9067ffffffffffffffff9093929316815260406020820152612f43612f06612ef58580612e3e565b60a0604086015260e0850191612e8e565b612f136020860186612e3e565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858403016060860152612e8e565b9060408401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811215612d105784016020813591019267ffffffffffffffff8211612d10578160061b36038413612d10578281037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0016080840152818152602001929060005b818110613044575050506130118473ffffffffffffffffffffffffffffffffffffffff6130016060613041979801612d95565b1660a08401526080810190612e3e565b9160c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082860301910152612e8e565b90565b90919360408060019273ffffffffffffffffffffffffffffffffffffffff61306b89612d95565b16815260208881013590820152019501929101612fce565b519073ffffffffffffffffffffffffffffffffffffffff82168203612d1057565b73ffffffffffffffffffffffffffffffffffffffff604051917fbbe4f6db00000000000000000000000000000000000000000000000000000000835216600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa801561319b5760009061314e575b73ffffffffffffffffffffffffffffffffffffffff91501690565b506020813d602011613193575b8161316860209383612c18565b81010312612d105761318e73ffffffffffffffffffffffffffffffffffffffff91613083565b613133565b3d915061315b565b6040513d6000823e3d90fd5b3573ffffffffffffffffffffffffffffffffffffffff81168103612d105790565b604051906131d582612bc3565b60006080838281528260208201528260408201528260608201520152565b9080601f83011215612d1057813561320a81612d2a565b926132186040519485612c18565b81845260208085019260051b820101928311612d1057602001905b8282106132405750505090565b6020809161324d84612d95565b815201910190613233565b805182101561326c5760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612d10570180359067ffffffffffffffff8211612d1057602001918136038313612d1057565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215612d10570180359067ffffffffffffffff8211612d1057602001918160061b36038313612d1057565b6040519061334d82612bc3565b6060608083600081528260208201528260408201526000838201520152565b92919261337882612c59565b916133866040519384612c18565b829481845281830111612d10578281602093846000960137010152565b81601f82011215612d105780516133b981612c59565b926133c76040519485612c18565b81845260208284010111612d10576130419160208085019101612c93565b9080602083519182815201916020808360051b8301019401926000915b83831061341157505050505090565b90919293946020806134ac837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301875289519073ffffffffffffffffffffffffffffffffffffffff8251168152608061349161347f8685015160a08886015260a0850190612cb6565b60408501518482036040860152612cb6565b92606081015160608401520151906080818403910152612cb6565b97019301930191939290613402565b519067ffffffffffffffff82168203612d1057565b90613041916020815267ffffffffffffffff6080835180516020850152826020820151166040850152826040820151166060850152826060820151168285015201511660a082015273ffffffffffffffffffffffffffffffffffffffff60208301511660c08201526101006135c461358f61355c60408601516101a060e08701526101c0860190612cb6565b60608601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08683030185870152612cb6565b60808501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085830301610120860152612cb6565b9273ffffffffffffffffffffffffffffffffffffffff60a08201511661014084015260c081015161016084015260e08101516101808401520151906101a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828503019101526133e5565b73ffffffffffffffffffffffffffffffffffffffff60015416330361365057565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b3d156136a5573d9061368b82612c59565b916136996040519384612c18565b82523d6000602084013e565b606090565b805482101561326c5760005260206000200190600090565b600082815260018201602052604090205461374c5780549068010000000000000000821015612b9457826137356137008460018096018555846136aa565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b905580549260005201602052604060002055600190565b5050600090565b90600182019181600052826020526040600020548015156000146138dc577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116138ad578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116138ad57818103613876575b50505080548015613847577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019061380882826136aa565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b61389661388661370093866136aa565b90549060031b1c928392866136aa565b9055600052836020526040600020553880806137d0565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50505050600090565b9192901561396057508151156138f9575090565b3b156139025790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b8251909150156139735750805190602001fd5b610bdb906040519182917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401526024830190612cb656fea164736f6c634300081a000a",
}

var MessageTransformerOnRampABI = MessageTransformerOnRampMetaData.ABI

var MessageTransformerOnRampBin = MessageTransformerOnRampMetaData.Bin

func DeployMessageTransformerOnRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OnRampStaticConfig, dynamicConfig OnRampDynamicConfig, destChainConfigs []OnRampDestChainConfigArgs, messageTransformerAddr common.Address) (common.Address, *types.Transaction, *MessageTransformerOnRamp, error) {
	parsed, err := MessageTransformerOnRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageTransformerOnRampBin), backend, staticConfig, dynamicConfig, destChainConfigs, messageTransformerAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MessageTransformerOnRamp{address: address, abi: *parsed, MessageTransformerOnRampCaller: MessageTransformerOnRampCaller{contract: contract}, MessageTransformerOnRampTransactor: MessageTransformerOnRampTransactor{contract: contract}, MessageTransformerOnRampFilterer: MessageTransformerOnRampFilterer{contract: contract}}, nil
}

type MessageTransformerOnRamp struct {
	address common.Address
	abi     abi.ABI
	MessageTransformerOnRampCaller
	MessageTransformerOnRampTransactor
	MessageTransformerOnRampFilterer
}

type MessageTransformerOnRampCaller struct {
	contract *bind.BoundContract
}

type MessageTransformerOnRampTransactor struct {
	contract *bind.BoundContract
}

type MessageTransformerOnRampFilterer struct {
	contract *bind.BoundContract
}

type MessageTransformerOnRampSession struct {
	Contract     *MessageTransformerOnRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MessageTransformerOnRampCallerSession struct {
	Contract *MessageTransformerOnRampCaller
	CallOpts bind.CallOpts
}

type MessageTransformerOnRampTransactorSession struct {
	Contract     *MessageTransformerOnRampTransactor
	TransactOpts bind.TransactOpts
}

type MessageTransformerOnRampRaw struct {
	Contract *MessageTransformerOnRamp
}

type MessageTransformerOnRampCallerRaw struct {
	Contract *MessageTransformerOnRampCaller
}

type MessageTransformerOnRampTransactorRaw struct {
	Contract *MessageTransformerOnRampTransactor
}

func NewMessageTransformerOnRamp(address common.Address, backend bind.ContractBackend) (*MessageTransformerOnRamp, error) {
	abi, err := abi.JSON(strings.NewReader(MessageTransformerOnRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMessageTransformerOnRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRamp{address: address, abi: abi, MessageTransformerOnRampCaller: MessageTransformerOnRampCaller{contract: contract}, MessageTransformerOnRampTransactor: MessageTransformerOnRampTransactor{contract: contract}, MessageTransformerOnRampFilterer: MessageTransformerOnRampFilterer{contract: contract}}, nil
}

func NewMessageTransformerOnRampCaller(address common.Address, caller bind.ContractCaller) (*MessageTransformerOnRampCaller, error) {
	contract, err := bindMessageTransformerOnRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampCaller{contract: contract}, nil
}

func NewMessageTransformerOnRampTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageTransformerOnRampTransactor, error) {
	contract, err := bindMessageTransformerOnRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampTransactor{contract: contract}, nil
}

func NewMessageTransformerOnRampFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageTransformerOnRampFilterer, error) {
	contract, err := bindMessageTransformerOnRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampFilterer{contract: contract}, nil
}

func bindMessageTransformerOnRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MessageTransformerOnRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageTransformerOnRamp.Contract.MessageTransformerOnRampCaller.contract.Call(opts, result, method, params...)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.MessageTransformerOnRampTransactor.contract.Transfer(opts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.MessageTransformerOnRampTransactor.contract.Transact(opts, method, params...)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageTransformerOnRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.contract.Transfer(opts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.contract.Transact(opts, method, params...)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

	error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getAllowedSendersList", destChainSelector)

	outstruct := new(GetAllowedSendersList)
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsEnabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfiguredAddresses = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _MessageTransformerOnRamp.Contract.GetAllowedSendersList(&_MessageTransformerOnRamp.CallOpts, destChainSelector)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetAllowedSendersList(destChainSelector uint64) (GetAllowedSendersList,

	error) {
	return _MessageTransformerOnRamp.Contract.GetAllowedSendersList(&_MessageTransformerOnRamp.CallOpts, destChainSelector)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

	error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getDestChainConfig", destChainSelector)

	outstruct := new(GetDestChainConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequenceNumber = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.AllowlistEnabled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Router = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _MessageTransformerOnRamp.Contract.GetDestChainConfig(&_MessageTransformerOnRamp.CallOpts, destChainSelector)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetDestChainConfig(destChainSelector uint64) (GetDestChainConfig,

	error) {
	return _MessageTransformerOnRamp.Contract.GetDestChainConfig(&_MessageTransformerOnRamp.CallOpts, destChainSelector)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OnRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampDynamicConfig)).(*OnRampDynamicConfig)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _MessageTransformerOnRamp.Contract.GetDynamicConfig(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetDynamicConfig() (OnRampDynamicConfig, error) {
	return _MessageTransformerOnRamp.Contract.GetDynamicConfig(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _MessageTransformerOnRamp.Contract.GetExpectedNextSequenceNumber(&_MessageTransformerOnRamp.CallOpts, destChainSelector)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _MessageTransformerOnRamp.Contract.GetExpectedNextSequenceNumber(&_MessageTransformerOnRamp.CallOpts, destChainSelector)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getFee", destChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _MessageTransformerOnRamp.Contract.GetFee(&_MessageTransformerOnRamp.CallOpts, destChainSelector, message)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetFee(destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _MessageTransformerOnRamp.Contract.GetFee(&_MessageTransformerOnRamp.CallOpts, destChainSelector, message)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetMessageTransformerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getMessageTransformerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetMessageTransformerAddress() (common.Address, error) {
	return _MessageTransformerOnRamp.Contract.GetMessageTransformerAddress(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetMessageTransformerAddress() (common.Address, error) {
	return _MessageTransformerOnRamp.Contract.GetMessageTransformerAddress(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetPoolBySourceToken(opts *bind.CallOpts, arg0 uint64, sourceToken common.Address) (common.Address, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getPoolBySourceToken", arg0, sourceToken)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _MessageTransformerOnRamp.Contract.GetPoolBySourceToken(&_MessageTransformerOnRamp.CallOpts, arg0, sourceToken)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetPoolBySourceToken(arg0 uint64, sourceToken common.Address) (common.Address, error) {
	return _MessageTransformerOnRamp.Contract.GetPoolBySourceToken(&_MessageTransformerOnRamp.CallOpts, arg0, sourceToken)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetStaticConfig(opts *bind.CallOpts) (OnRampStaticConfig, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OnRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OnRampStaticConfig)).(*OnRampStaticConfig)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _MessageTransformerOnRamp.Contract.GetStaticConfig(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetStaticConfig() (OnRampStaticConfig, error) {
	return _MessageTransformerOnRamp.Contract.GetStaticConfig(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) GetSupportedTokens(opts *bind.CallOpts, arg0 uint64) ([]common.Address, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "getSupportedTokens", arg0)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _MessageTransformerOnRamp.Contract.GetSupportedTokens(&_MessageTransformerOnRamp.CallOpts, arg0)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) GetSupportedTokens(arg0 uint64) ([]common.Address, error) {
	return _MessageTransformerOnRamp.Contract.GetSupportedTokens(&_MessageTransformerOnRamp.CallOpts, arg0)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) Owner() (common.Address, error) {
	return _MessageTransformerOnRamp.Contract.Owner(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) Owner() (common.Address, error) {
	return _MessageTransformerOnRamp.Contract.Owner(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MessageTransformerOnRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) TypeAndVersion() (string, error) {
	return _MessageTransformerOnRamp.Contract.TypeAndVersion(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampCallerSession) TypeAndVersion() (string, error) {
	return _MessageTransformerOnRamp.Contract.TypeAndVersion(&_MessageTransformerOnRamp.CallOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.contract.Transact(opts, "acceptOwnership")
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.AcceptOwnership(&_MessageTransformerOnRamp.TransactOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.AcceptOwnership(&_MessageTransformerOnRamp.TransactOpts)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactor) ApplyAllowlistUpdates(opts *bind.TransactOpts, allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.contract.Transact(opts, "applyAllowlistUpdates", allowlistConfigArgsItems)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.ApplyAllowlistUpdates(&_MessageTransformerOnRamp.TransactOpts, allowlistConfigArgsItems)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorSession) ApplyAllowlistUpdates(allowlistConfigArgsItems []OnRampAllowlistConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.ApplyAllowlistUpdates(&_MessageTransformerOnRamp.TransactOpts, allowlistConfigArgsItems)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactor) ApplyDestChainConfigUpdates(opts *bind.TransactOpts, destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.contract.Transact(opts, "applyDestChainConfigUpdates", destChainConfigArgs)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.ApplyDestChainConfigUpdates(&_MessageTransformerOnRamp.TransactOpts, destChainConfigArgs)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorSession) ApplyDestChainConfigUpdates(destChainConfigArgs []OnRampDestChainConfigArgs) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.ApplyDestChainConfigUpdates(&_MessageTransformerOnRamp.TransactOpts, destChainConfigArgs)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactor) ForwardFromRouter(opts *bind.TransactOpts, destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.contract.Transact(opts, "forwardFromRouter", destChainSelector, message, feeTokenAmount, originalSender)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.ForwardFromRouter(&_MessageTransformerOnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorSession) ForwardFromRouter(destChainSelector uint64, message ClientEVM2AnyMessage, feeTokenAmount *big.Int, originalSender common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.ForwardFromRouter(&_MessageTransformerOnRamp.TransactOpts, destChainSelector, message, feeTokenAmount, originalSender)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.SetDynamicConfig(&_MessageTransformerOnRamp.TransactOpts, dynamicConfig)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorSession) SetDynamicConfig(dynamicConfig OnRampDynamicConfig) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.SetDynamicConfig(&_MessageTransformerOnRamp.TransactOpts, dynamicConfig)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.TransferOwnership(&_MessageTransformerOnRamp.TransactOpts, to)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.TransferOwnership(&_MessageTransformerOnRamp.TransactOpts, to)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactor) WithdrawFeeTokens(opts *bind.TransactOpts, feeTokens []common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.contract.Transact(opts, "withdrawFeeTokens", feeTokens)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.WithdrawFeeTokens(&_MessageTransformerOnRamp.TransactOpts, feeTokens)
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampTransactorSession) WithdrawFeeTokens(feeTokens []common.Address) (*types.Transaction, error) {
	return _MessageTransformerOnRamp.Contract.WithdrawFeeTokens(&_MessageTransformerOnRamp.TransactOpts, feeTokens)
}

type MessageTransformerOnRampAllowListAdminSetIterator struct {
	Event *MessageTransformerOnRampAllowListAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampAllowListAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampAllowListAdminSet)
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
		it.Event = new(MessageTransformerOnRampAllowListAdminSet)
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

func (it *MessageTransformerOnRampAllowListAdminSetIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampAllowListAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampAllowListAdminSet struct {
	AllowlistAdmin common.Address
	Raw            types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*MessageTransformerOnRampAllowListAdminSetIterator, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampAllowListAdminSetIterator{contract: _MessageTransformerOnRamp.contract, event: "AllowListAdminSet", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error) {

	var allowlistAdminRule []interface{}
	for _, allowlistAdminItem := range allowlistAdmin {
		allowlistAdminRule = append(allowlistAdminRule, allowlistAdminItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "AllowListAdminSet", allowlistAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampAllowListAdminSet)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseAllowListAdminSet(log types.Log) (*MessageTransformerOnRampAllowListAdminSet, error) {
	event := new(MessageTransformerOnRampAllowListAdminSet)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "AllowListAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampAllowListSendersAddedIterator struct {
	Event *MessageTransformerOnRampAllowListSendersAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampAllowListSendersAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampAllowListSendersAdded)
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
		it.Event = new(MessageTransformerOnRampAllowListSendersAdded)
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

func (it *MessageTransformerOnRampAllowListSendersAddedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampAllowListSendersAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampAllowListSendersAdded struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*MessageTransformerOnRampAllowListSendersAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampAllowListSendersAddedIterator{contract: _MessageTransformerOnRamp.contract, event: "AllowListSendersAdded", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "AllowListSendersAdded", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampAllowListSendersAdded)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseAllowListSendersAdded(log types.Log) (*MessageTransformerOnRampAllowListSendersAdded, error) {
	event := new(MessageTransformerOnRampAllowListSendersAdded)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "AllowListSendersAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampAllowListSendersRemovedIterator struct {
	Event *MessageTransformerOnRampAllowListSendersRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampAllowListSendersRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampAllowListSendersRemoved)
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
		it.Event = new(MessageTransformerOnRampAllowListSendersRemoved)
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

func (it *MessageTransformerOnRampAllowListSendersRemovedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampAllowListSendersRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampAllowListSendersRemoved struct {
	DestChainSelector uint64
	Senders           []common.Address
	Raw               types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*MessageTransformerOnRampAllowListSendersRemovedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampAllowListSendersRemovedIterator{contract: _MessageTransformerOnRamp.contract, event: "AllowListSendersRemoved", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "AllowListSendersRemoved", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampAllowListSendersRemoved)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseAllowListSendersRemoved(log types.Log) (*MessageTransformerOnRampAllowListSendersRemoved, error) {
	event := new(MessageTransformerOnRampAllowListSendersRemoved)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "AllowListSendersRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampCCIPMessageSentIterator struct {
	Event *MessageTransformerOnRampCCIPMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampCCIPMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampCCIPMessageSent)
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
		it.Event = new(MessageTransformerOnRampCCIPMessageSent)
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

func (it *MessageTransformerOnRampCCIPMessageSentIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampCCIPMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*MessageTransformerOnRampCCIPMessageSentIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampCCIPMessageSentIterator{contract: _MessageTransformerOnRamp.contract, event: "CCIPMessageSent", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampCCIPMessageSent)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseCCIPMessageSent(log types.Log) (*MessageTransformerOnRampCCIPMessageSent, error) {
	event := new(MessageTransformerOnRampCCIPMessageSent)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampConfigSetIterator struct {
	Event *MessageTransformerOnRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampConfigSet)
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
		it.Event = new(MessageTransformerOnRampConfigSet)
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

func (it *MessageTransformerOnRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampConfigSet struct {
	StaticConfig  OnRampStaticConfig
	DynamicConfig OnRampDynamicConfig
	Raw           types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*MessageTransformerOnRampConfigSetIterator, error) {

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampConfigSetIterator{contract: _MessageTransformerOnRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampConfigSet)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseConfigSet(log types.Log) (*MessageTransformerOnRampConfigSet, error) {
	event := new(MessageTransformerOnRampConfigSet)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampDestChainConfigSetIterator struct {
	Event *MessageTransformerOnRampDestChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampDestChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampDestChainConfigSet)
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
		it.Event = new(MessageTransformerOnRampDestChainConfigSet)
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

func (it *MessageTransformerOnRampDestChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampDestChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampDestChainConfigSet struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Router            common.Address
	AllowlistEnabled  bool
	Raw               types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*MessageTransformerOnRampDestChainConfigSetIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampDestChainConfigSetIterator{contract: _MessageTransformerOnRamp.contract, event: "DestChainConfigSet", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "DestChainConfigSet", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampDestChainConfigSet)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseDestChainConfigSet(log types.Log) (*MessageTransformerOnRampDestChainConfigSet, error) {
	event := new(MessageTransformerOnRampDestChainConfigSet)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "DestChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampFeeTokenWithdrawnIterator struct {
	Event *MessageTransformerOnRampFeeTokenWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampFeeTokenWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampFeeTokenWithdrawn)
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
		it.Event = new(MessageTransformerOnRampFeeTokenWithdrawn)
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

func (it *MessageTransformerOnRampFeeTokenWithdrawnIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampFeeTokenWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampFeeTokenWithdrawn struct {
	FeeAggregator common.Address
	FeeToken      common.Address
	Amount        *big.Int
	Raw           types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*MessageTransformerOnRampFeeTokenWithdrawnIterator, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampFeeTokenWithdrawnIterator{contract: _MessageTransformerOnRamp.contract, event: "FeeTokenWithdrawn", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error) {

	var feeAggregatorRule []interface{}
	for _, feeAggregatorItem := range feeAggregator {
		feeAggregatorRule = append(feeAggregatorRule, feeAggregatorItem)
	}
	var feeTokenRule []interface{}
	for _, feeTokenItem := range feeToken {
		feeTokenRule = append(feeTokenRule, feeTokenItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "FeeTokenWithdrawn", feeAggregatorRule, feeTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampFeeTokenWithdrawn)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseFeeTokenWithdrawn(log types.Log) (*MessageTransformerOnRampFeeTokenWithdrawn, error) {
	event := new(MessageTransformerOnRampFeeTokenWithdrawn)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "FeeTokenWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampOwnershipTransferRequestedIterator struct {
	Event *MessageTransformerOnRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampOwnershipTransferRequested)
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
		it.Event = new(MessageTransformerOnRampOwnershipTransferRequested)
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

func (it *MessageTransformerOnRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOnRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampOwnershipTransferRequestedIterator{contract: _MessageTransformerOnRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampOwnershipTransferRequested)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*MessageTransformerOnRampOwnershipTransferRequested, error) {
	event := new(MessageTransformerOnRampOwnershipTransferRequested)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MessageTransformerOnRampOwnershipTransferredIterator struct {
	Event *MessageTransformerOnRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MessageTransformerOnRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageTransformerOnRampOwnershipTransferred)
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
		it.Event = new(MessageTransformerOnRampOwnershipTransferred)
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

func (it *MessageTransformerOnRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MessageTransformerOnRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MessageTransformerOnRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOnRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MessageTransformerOnRampOwnershipTransferredIterator{contract: _MessageTransformerOnRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MessageTransformerOnRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MessageTransformerOnRampOwnershipTransferred)
				if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRampFilterer) ParseOwnershipTransferred(log types.Log) (*MessageTransformerOnRampOwnershipTransferred, error) {
	event := new(MessageTransformerOnRampOwnershipTransferred)
	if err := _MessageTransformerOnRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MessageTransformerOnRamp *MessageTransformerOnRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MessageTransformerOnRamp.abi.Events["AllowListAdminSet"].ID:
		return _MessageTransformerOnRamp.ParseAllowListAdminSet(log)
	case _MessageTransformerOnRamp.abi.Events["AllowListSendersAdded"].ID:
		return _MessageTransformerOnRamp.ParseAllowListSendersAdded(log)
	case _MessageTransformerOnRamp.abi.Events["AllowListSendersRemoved"].ID:
		return _MessageTransformerOnRamp.ParseAllowListSendersRemoved(log)
	case _MessageTransformerOnRamp.abi.Events["CCIPMessageSent"].ID:
		return _MessageTransformerOnRamp.ParseCCIPMessageSent(log)
	case _MessageTransformerOnRamp.abi.Events["ConfigSet"].ID:
		return _MessageTransformerOnRamp.ParseConfigSet(log)
	case _MessageTransformerOnRamp.abi.Events["DestChainConfigSet"].ID:
		return _MessageTransformerOnRamp.ParseDestChainConfigSet(log)
	case _MessageTransformerOnRamp.abi.Events["FeeTokenWithdrawn"].ID:
		return _MessageTransformerOnRamp.ParseFeeTokenWithdrawn(log)
	case _MessageTransformerOnRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _MessageTransformerOnRamp.ParseOwnershipTransferRequested(log)
	case _MessageTransformerOnRamp.abi.Events["OwnershipTransferred"].ID:
		return _MessageTransformerOnRamp.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MessageTransformerOnRampAllowListAdminSet) Topic() common.Hash {
	return common.HexToHash("0xb8c9b44ae5b5e3afb195f67391d9ff50cb904f9c0fa5fd520e497a97c1aa5a1e")
}

func (MessageTransformerOnRampAllowListSendersAdded) Topic() common.Hash {
	return common.HexToHash("0x330939f6eafe8bb516716892fe962ff19770570838686e6579dbc1cc51fc3281")
}

func (MessageTransformerOnRampAllowListSendersRemoved) Topic() common.Hash {
	return common.HexToHash("0xc237ec1921f855ccd5e9a5af9733f2d58943a5a8501ec5988e305d7a4d421586")
}

func (MessageTransformerOnRampCCIPMessageSent) Topic() common.Hash {
	return common.HexToHash("0x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32")
}

func (MessageTransformerOnRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc7372d2d886367d7bb1b0e0708a5436f2c91d6963de210eb2dc1ec2ecd6d21f1")
}

func (MessageTransformerOnRampDestChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0xd5ad72bc37dc7a80a8b9b9df20500046fd7341adb1be2258a540466fdd7dcef5")
}

func (MessageTransformerOnRampFeeTokenWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x508d7d183612c18fc339b42618912b9fa3239f631dd7ec0671f950200a0fa66e")
}

func (MessageTransformerOnRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MessageTransformerOnRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_MessageTransformerOnRamp *MessageTransformerOnRamp) Address() common.Address {
	return _MessageTransformerOnRamp.address
}

type MessageTransformerOnRampInterface interface {
	GetAllowedSendersList(opts *bind.CallOpts, destChainSelector uint64) (GetAllowedSendersList,

		error)

	GetDestChainConfig(opts *bind.CallOpts, destChainSelector uint64) (GetDestChainConfig,

		error)

	GetDynamicConfig(opts *bind.CallOpts) (OnRampDynamicConfig, error)

	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetFee(opts *bind.CallOpts, destChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetMessageTransformerAddress(opts *bind.CallOpts) (common.Address, error)

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

	FilterAllowListAdminSet(opts *bind.FilterOpts, allowlistAdmin []common.Address) (*MessageTransformerOnRampAllowListAdminSetIterator, error)

	WatchAllowListAdminSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampAllowListAdminSet, allowlistAdmin []common.Address) (event.Subscription, error)

	ParseAllowListAdminSet(log types.Log) (*MessageTransformerOnRampAllowListAdminSet, error)

	FilterAllowListSendersAdded(opts *bind.FilterOpts, destChainSelector []uint64) (*MessageTransformerOnRampAllowListSendersAddedIterator, error)

	WatchAllowListSendersAdded(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampAllowListSendersAdded, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersAdded(log types.Log) (*MessageTransformerOnRampAllowListSendersAdded, error)

	FilterAllowListSendersRemoved(opts *bind.FilterOpts, destChainSelector []uint64) (*MessageTransformerOnRampAllowListSendersRemovedIterator, error)

	WatchAllowListSendersRemoved(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampAllowListSendersRemoved, destChainSelector []uint64) (event.Subscription, error)

	ParseAllowListSendersRemoved(log types.Log) (*MessageTransformerOnRampAllowListSendersRemoved, error)

	FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*MessageTransformerOnRampCCIPMessageSentIterator, error)

	WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error)

	ParseCCIPMessageSent(log types.Log) (*MessageTransformerOnRampCCIPMessageSent, error)

	FilterConfigSet(opts *bind.FilterOpts) (*MessageTransformerOnRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*MessageTransformerOnRampConfigSet, error)

	FilterDestChainConfigSet(opts *bind.FilterOpts, destChainSelector []uint64) (*MessageTransformerOnRampDestChainConfigSetIterator, error)

	WatchDestChainConfigSet(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampDestChainConfigSet, destChainSelector []uint64) (event.Subscription, error)

	ParseDestChainConfigSet(log types.Log) (*MessageTransformerOnRampDestChainConfigSet, error)

	FilterFeeTokenWithdrawn(opts *bind.FilterOpts, feeAggregator []common.Address, feeToken []common.Address) (*MessageTransformerOnRampFeeTokenWithdrawnIterator, error)

	WatchFeeTokenWithdrawn(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampFeeTokenWithdrawn, feeAggregator []common.Address, feeToken []common.Address) (event.Subscription, error)

	ParseFeeTokenWithdrawn(log types.Log) (*MessageTransformerOnRampFeeTokenWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOnRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MessageTransformerOnRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MessageTransformerOnRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MessageTransformerOnRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MessageTransformerOnRampOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
