package rmn_remote

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

type IRMNRemoteSignature struct {
	R [32]byte
	S [32]byte
}

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type RMNRemoteConfig struct {
	RmnHomeContractConfigDigest [32]byte
	Signers                     []RMNRemoteSigner
	F                           uint64
}

type RMNRemoteSigner struct {
	OnchainPublicKey common.Address
	NodeIndex        uint64
}

var RMNRemoteMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"AlreadyCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOnchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignerOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"NotCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfOrderSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ThresholdNotMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroValueNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Cursed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Uncursed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCursedSubjects\",\"outputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalChainSelector\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReportDigestHeader\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"digestHeader\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersionedConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offrampAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"verify\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620020ff380380620020ff833981016040819052620000349162000142565b336000816200005657604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038481169190911790915581161562000089576200008981620000c8565b5050806001600160401b0316600003620000b65760405163273e150360e21b815260040160405180910390fd5b6001600160401b031660805262000174565b336001600160a01b03821603620000f257604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200015557600080fd5b81516001600160401b03811681146200016d57600080fd5b9392505050565b608051611f68620001976000396000818161027a0152610a2c0152611f686000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806370a9089e11610097578063d881e09211610066578063d881e09214610257578063eaa83ddd1461026c578063f2fde38b146102a4578063f8bb876e146102b757600080fd5b806370a9089e1461020157806379ba5097146102145780638da5cb5b1461021c5780639a19b3291461024457600080fd5b8063397796f7116100d3578063397796f7146101a557806362eed415146101ad5780636509a954146101c05780636d2d3993146101ee57600080fd5b8063181f5a7714610105578063198f0f77146101575780631add205f1461016c5780632cbc26bb14610182575b600080fd5b6101416040518060400160405280601381526020017f524d4e52656d6f746520312e362e302d6465760000000000000000000000000081525081565b60405161014e9190611389565b60405180910390f35b61016a61016536600461139c565b6102ca565b005b6101746106c4565b60405161014e9291906113d7565b6101956101903660046114b5565b6107bc565b604051901515815260200161014e565b610195610819565b61016a6101bb3660046114b5565b610893565b6040517f9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf53815260200161014e565b61016a6101fc3660046114b5565b610907565b61016a61020f36600461153e565b610977565b61016a610cd2565b60015460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014e565b61016a6102523660046116bd565b610da0565b61025f610ea6565b60405161014e919061175a565b60405167ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016815260200161014e565b61016a6102b23660046117c0565b610eb2565b61016a6102c53660046116bd565b610ec6565b6102d2610fb8565b803561030a576040517f9cf8540c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015b61031a60208301836117dd565b90508110156103ea5761033060208301836117dd565b8281811061034057610340611845565b90506040020160200160208101906103589190611895565b67ffffffffffffffff1661036f60208401846117dd565b61037a6001856118e1565b81811061038957610389611845565b90506040020160200160208101906103a19190611895565b67ffffffffffffffff16106103e2576040517f4485151700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60010161030d565b506103fb6060820160408301611895565b6104069060026118f4565b610411906001611920565b67ffffffffffffffff1661042860208301836117dd565b90501015610462576040517f014c502000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6003545b80156104f45760086000600361047d6001856118e1565b8154811061048d5761048d611845565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556104ed81611941565b9050610466565b5060005b61050560208301836117dd565b905081101561063a576008600061051f60208501856117dd565b8481811061052f5761052f611845565b61054592602060409092020190810191506117c0565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040016000205460ff16156105a6576040517f28cae27d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600860006105b960208601866117dd565b858181106105c9576105c9611845565b6105df92602060409092020190810191506117c0565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790556001016104f8565b508060026106488282611a2f565b5050600580546000919082906106639063ffffffff16611b6a565b91906101000a81548163ffffffff021916908363ffffffff160217905590508063ffffffff167f7f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c836040516106b89190611b8d565b60405180910390a25050565b6040805160608082018352600080835260208301919091529181018290526005546040805160608101825260028054825260038054845160208281028201810190965281815263ffffffff9096169592948593818601939092909160009084015b82821015610793576000848152602090819020604080518082019091529084015473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000900467ffffffffffffffff1681830152825260019092019101610725565b505050908252506002919091015467ffffffffffffffff16602090910152919491935090915050565b60006107c8600661100b565b6000036107d757506000919050565b6107e2600683611015565b80610813575061081360067f0100000000000000000000000000000100000000000000000000000000000000611015565b92915050565b6000610825600661100b565b6000036108325750600090565b61085d60067f0100000000000000000000000000000000000000000000000000000000000000611015565b8061088e575061088e60067f0100000000000000000000000000000100000000000000000000000000000000611015565b905090565b6040805160018082528183019092526000916020808301908036833701905050905081816000815181106108c9576108c9611845565b7fffffffffffffffffffffffffffffffff000000000000000000000000000000009092166020928302919091019091015261090381610ec6565b5050565b60408051600180825281830190925260009160208083019080368337019050509050818160008151811061093d5761093d611845565b7fffffffffffffffffffffffffffffffff000000000000000000000000000000009092166020928302919091019091015261090381610da0565b60055463ffffffff166000036109b9576040517face124bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004546109d19067ffffffffffffffff166001611920565b67ffffffffffffffff16811015610a14576040517f59fa4a9300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160c08101825246815267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166020820152309181019190915273ffffffffffffffffffffffffffffffffffffffff8616606082015260025460808201526000907f9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf539060a08101610ab08789611c97565b9052604051610ac3929190602001611df7565b60405160208183030381529060405280519060200120905060008060005b84811015610cc757600184601b888885818110610b0057610b00611845565b90506040020160000135898986818110610b1c57610b1c611845565b9050604002016020013560405160008152602001604052604051610b5c949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610b7e573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015192505073ffffffffffffffffffffffffffffffffffffffff8216610bf6576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1610610c5b576040517fbbe15e7f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526008602052604090205460ff16610cba576040517faaaa914100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9091508190600101610ae1565b505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610d23576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610da8610fb8565b60005b8151811015610e6b57610de1828281518110610dc957610dc9611845565b6020026020010151600661105390919063ffffffff16565b610e6357818181518110610df757610df7611845565b60200260200101516040517f73281fa1000000000000000000000000000000000000000000000000000000008152600401610e5a91907fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b60405180910390fd5b600101610dab565b507f0676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba181604051610e9b919061175a565b60405180910390a150565b606061088e6006611081565b610eba610fb8565b610ec38161108e565b50565b610ece610fb8565b60005b8151811015610f8857610f07828281518110610eef57610eef611845565b6020026020010151600661115290919063ffffffff16565b610f8057818181518110610f1d57610f1d611845565b60200260200101516040517f19d5c79b000000000000000000000000000000000000000000000000000000008152600401610e5a91907fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b600101610ed1565b507f1716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f7481604051610e9b919061175a565b60015473ffffffffffffffffffffffffffffffffffffffff163314611009576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000610813825490565b7fffffffffffffffffffffffffffffffff000000000000000000000000000000008116600090815260018301602052604081205415155b9392505050565b600061104c837fffffffffffffffffffffffffffffffff000000000000000000000000000000008416611180565b6060600061104c8361127a565b3373ffffffffffffffffffffffffffffffffffffffff8216036110dd576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061104c837fffffffffffffffffffffffffffffffff0000000000000000000000000000000084166112d6565b600081815260018301602052604081205480156112695760006111a46001836118e1565b85549091506000906111b8906001906118e1565b905080821461121d5760008660000182815481106111d8576111d8611845565b90600052602060002001549050808760000184815481106111fb576111fb611845565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061122e5761122e611f2c565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610813565b6000915050610813565b5092915050565b6060816000018054806020026020016040519081016040528092919081815260200182805480156112ca57602002820191906000526020600020905b8154815260200190600101908083116112b6575b50505050509050919050565b600081815260018301602052604081205461131d57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610813565b506000610813565b6000815180845260005b8181101561134b5760208185018101518683018201520161132f565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061104c6020830184611325565b6000602082840312156113ae57600080fd5b813567ffffffffffffffff8111156113c557600080fd5b82016060818503121561104c57600080fd5b63ffffffff831681526040602080830182905283518383015283810151606080850152805160a085018190526000939291820190849060c08701905b8083101561145c578351805173ffffffffffffffffffffffffffffffffffffffff16835285015167ffffffffffffffff1685830152928401926001929092019190850190611413565b50604088015167ffffffffffffffff81166080890152945098975050505050505050565b80357fffffffffffffffffffffffffffffffff00000000000000000000000000000000811681146114b057600080fd5b919050565b6000602082840312156114c757600080fd5b61104c82611480565b73ffffffffffffffffffffffffffffffffffffffff81168114610ec357600080fd5b60008083601f84011261150457600080fd5b50813567ffffffffffffffff81111561151c57600080fd5b6020830191508360208260061b850101111561153757600080fd5b9250929050565b60008060008060006060868803121561155657600080fd5b8535611561816114d0565b9450602086013567ffffffffffffffff8082111561157e57600080fd5b818801915088601f83011261159257600080fd5b8135818111156115a157600080fd5b8960208260051b85010111156115b657600080fd5b6020830196508095505060408801359150808211156115d457600080fd5b506115e1888289016114f2565b969995985093965092949392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715611644576116446115f2565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611691576116916115f2565b604052919050565b600067ffffffffffffffff8211156116b3576116b36115f2565b5060051b60200190565b600060208083850312156116d057600080fd5b823567ffffffffffffffff8111156116e757600080fd5b8301601f810185136116f857600080fd5b803561170b61170682611699565b61164a565b81815260059190911b8201830190838101908783111561172a57600080fd5b928401925b8284101561174f5761174084611480565b8252928401929084019061172f565b979650505050505050565b6020808252825182820181905260009190848201906040850190845b818110156117b45783517fffffffffffffffffffffffffffffffff000000000000000000000000000000001683529284019291840191600101611776565b50909695505050505050565b6000602082840312156117d257600080fd5b813561104c816114d0565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261181257600080fd5b83018035915067ffffffffffffffff82111561182d57600080fd5b6020019150600681901b360382131561153757600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b67ffffffffffffffff81168114610ec357600080fd5b80356114b081611874565b6000602082840312156118a757600080fd5b813561104c81611874565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115610813576108136118b2565b67ffffffffffffffff818116838216028082169190828114611918576119186118b2565b505092915050565b67ffffffffffffffff818116838216019080821115611273576112736118b2565b600081611950576119506118b2565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b6000813561081381611874565b813561198e816114d0565b73ffffffffffffffffffffffffffffffffffffffff811690508154817fffffffffffffffffffffffff0000000000000000000000000000000000000000821617835560208401356119de81611874565b7bffffffffffffffff00000000000000000000000000000000000000008160a01b16837fffffffff000000000000000000000000000000000000000000000000000000008416171784555050505050565b81358155600180820160208401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112611a6d57600080fd5b8401803567ffffffffffffffff811115611a8657600080fd5b6020820191508060061b3603821315611a9e57600080fd5b68010000000000000000811115611ab757611ab76115f2565b825481845580821015611aec576000848152602081208381019083015b80821015611ae85782825590870190611ad4565b5050505b50600092835260208320925b81811015611b1c57611b0a8385611983565b92840192604092909201918401611af8565b5050505050610903611b3060408401611976565b6002830167ffffffffffffffff82167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008254161781555050565b600063ffffffff808316818103611b8357611b836118b2565b6001019392505050565b6000602080835260808301843582850152818501357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1863603018112611bd257600080fd5b8501828101903567ffffffffffffffff80821115611bef57600080fd5b8160061b3603831315611c0157600080fd5b6040606060408901528483865260a089019050849550600094505b83851015611c6c578535611c2f816114d0565b73ffffffffffffffffffffffffffffffffffffffff16815285870135611c5481611874565b83168188015294810194600194909401938101611c1c565b611c7860408b0161188a565b67ffffffffffffffff811660608b015296509998505050505050505050565b6000611ca561170684611699565b80848252602080830192508560051b850136811115611cc357600080fd5b855b81811015611deb57803567ffffffffffffffff80821115611ce65760008081fd5b818901915060a08236031215611cfc5760008081fd5b611d04611621565b8235611d0f81611874565b81528286013582811115611d235760008081fd5b8301601f3681830112611d365760008081fd5b813584811115611d4857611d486115f2565b611d77897fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848401160161164a565b94508085523689828501011115611d9057600091508182fd5b808984018a8701376000898287010152505050818682015260409150611db782840161188a565b8282015260609150611dca82840161188a565b91810191909152608091820135918101919091528552938201938201611cc5565b50919695505050505050565b60006040848352602060408185015261010084018551604086015281860151606067ffffffffffffffff808316606089015260408901519250608073ffffffffffffffffffffffffffffffffffffffff80851660808b015260608b0151945060a081861660a08c015260808c015160c08c015260a08c0151955060c060e08c015286915085518088526101209750878c019250878160051b8d01019750888701965060005b81811015611f19577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee08d8a030184528751868151168a528a810151848c8c0152611ee8858c0182611325565b828e015189168c8f01528983015189168a8d0152918701519a87019a909a5298509689019692890192600101611e9c565b50969d9c50505050505050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var RMNRemoteABI = RMNRemoteMetaData.ABI

var RMNRemoteBin = RMNRemoteMetaData.Bin

func DeployRMNRemote(auth *bind.TransactOpts, backend bind.ContractBackend, localChainSelector uint64) (common.Address, *generated_zks.Transaction, *RMNRemote, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated_zks.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated_zks.DeployContract(auth, parsed, common.FromHex(RMNRemoteZKBin), backend, localChainSelector)
		contractReturn := &RMNRemote{address: address, abi: *parsed, RMNRemoteCaller: RMNRemoteCaller{contract: contractBind}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contractBind}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNRemoteBin), backend, localChainSelector)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated_zks.Transaction{Transaction: tx, Hash_zks: tx.Hash()}, &RMNRemote{address: address, abi: *parsed, RMNRemoteCaller: RMNRemoteCaller{contract: contract}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contract}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contract}}, nil
}

type RMNRemote struct {
	address common.Address
	abi     abi.ABI
	RMNRemoteCaller
	RMNRemoteTransactor
	RMNRemoteFilterer
}

type RMNRemoteCaller struct {
	contract *bind.BoundContract
}

type RMNRemoteTransactor struct {
	contract *bind.BoundContract
}

type RMNRemoteFilterer struct {
	contract *bind.BoundContract
}

type RMNRemoteSession struct {
	Contract     *RMNRemote
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNRemoteCallerSession struct {
	Contract *RMNRemoteCaller
	CallOpts bind.CallOpts
}

type RMNRemoteTransactorSession struct {
	Contract     *RMNRemoteTransactor
	TransactOpts bind.TransactOpts
}

type RMNRemoteRaw struct {
	Contract *RMNRemote
}

type RMNRemoteCallerRaw struct {
	Contract *RMNRemoteCaller
}

type RMNRemoteTransactorRaw struct {
	Contract *RMNRemoteTransactor
}

func NewRMNRemote(address common.Address, backend bind.ContractBackend) (*RMNRemote, error) {
	abi, err := abi.JSON(strings.NewReader(RMNRemoteABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNRemote(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNRemote{address: address, abi: abi, RMNRemoteCaller: RMNRemoteCaller{contract: contract}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contract}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contract}}, nil
}

func NewRMNRemoteCaller(address common.Address, caller bind.ContractCaller) (*RMNRemoteCaller, error) {
	contract, err := bindRMNRemote(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteCaller{contract: contract}, nil
}

func NewRMNRemoteTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNRemoteTransactor, error) {
	contract, err := bindRMNRemote(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteTransactor{contract: contract}, nil
}

func NewRMNRemoteFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNRemoteFilterer, error) {
	contract, err := bindRMNRemote(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteFilterer{contract: contract}, nil
}

func bindRMNRemote(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNRemote *RMNRemoteRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNRemote.Contract.RMNRemoteCaller.contract.Call(opts, result, method, params...)
}

func (_RMNRemote *RMNRemoteRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.Contract.RMNRemoteTransactor.contract.Transfer(opts)
}

func (_RMNRemote *RMNRemoteRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNRemote.Contract.RMNRemoteTransactor.contract.Transact(opts, method, params...)
}

func (_RMNRemote *RMNRemoteCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNRemote.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNRemote *RMNRemoteTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.Contract.contract.Transfer(opts)
}

func (_RMNRemote *RMNRemoteTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNRemote.Contract.contract.Transact(opts, method, params...)
}

func (_RMNRemote *RMNRemoteCaller) GetCursedSubjects(opts *bind.CallOpts) ([][16]byte, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getCursedSubjects")

	if err != nil {
		return *new([][16]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][16]byte)).(*[][16]byte)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) GetCursedSubjects() ([][16]byte, error) {
	return _RMNRemote.Contract.GetCursedSubjects(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetCursedSubjects() ([][16]byte, error) {
	return _RMNRemote.Contract.GetCursedSubjects(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) GetLocalChainSelector(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getLocalChainSelector")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) GetLocalChainSelector() (uint64, error) {
	return _RMNRemote.Contract.GetLocalChainSelector(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetLocalChainSelector() (uint64, error) {
	return _RMNRemote.Contract.GetLocalChainSelector(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) GetReportDigestHeader(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getReportDigestHeader")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) GetReportDigestHeader() ([32]byte, error) {
	return _RMNRemote.Contract.GetReportDigestHeader(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetReportDigestHeader() ([32]byte, error) {
	return _RMNRemote.Contract.GetReportDigestHeader(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) GetVersionedConfig(opts *bind.CallOpts) (GetVersionedConfig,

	error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "getVersionedConfig")

	outstruct := new(GetVersionedConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Version = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Config = *abi.ConvertType(out[1], new(RMNRemoteConfig)).(*RMNRemoteConfig)

	return *outstruct, err

}

func (_RMNRemote *RMNRemoteSession) GetVersionedConfig() (GetVersionedConfig,

	error) {
	return _RMNRemote.Contract.GetVersionedConfig(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) GetVersionedConfig() (GetVersionedConfig,

	error) {
	return _RMNRemote.Contract.GetVersionedConfig(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isCursed", subject)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsCursed(subject [16]byte) (bool, error) {
	return _RMNRemote.Contract.IsCursed(&_RMNRemote.CallOpts, subject)
}

func (_RMNRemote *RMNRemoteCallerSession) IsCursed(subject [16]byte) (bool, error) {
	return _RMNRemote.Contract.IsCursed(&_RMNRemote.CallOpts, subject)
}

func (_RMNRemote *RMNRemoteCaller) IsCursed0(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isCursed0")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsCursed0() (bool, error) {
	return _RMNRemote.Contract.IsCursed0(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) IsCursed0() (bool, error) {
	return _RMNRemote.Contract.IsCursed0(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) Owner() (common.Address, error) {
	return _RMNRemote.Contract.Owner(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) Owner() (common.Address, error) {
	return _RMNRemote.Contract.Owner(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) TypeAndVersion() (string, error) {
	return _RMNRemote.Contract.TypeAndVersion(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCallerSession) TypeAndVersion() (string, error) {
	return _RMNRemote.Contract.TypeAndVersion(&_RMNRemote.CallOpts)
}

func (_RMNRemote *RMNRemoteCaller) Verify(opts *bind.CallOpts, offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "verify", offrampAddress, merkleRoots, signatures)

	if err != nil {
		return err
	}

	return err

}

func (_RMNRemote *RMNRemoteSession) Verify(offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, offrampAddress, merkleRoots, signatures)
}

func (_RMNRemote *RMNRemoteCallerSession) Verify(offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error {
	return _RMNRemote.Contract.Verify(&_RMNRemote.CallOpts, offrampAddress, merkleRoots, signatures)
}

func (_RMNRemote *RMNRemoteTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "acceptOwnership")
}

func (_RMNRemote *RMNRemoteSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNRemote.Contract.AcceptOwnership(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNRemote.Contract.AcceptOwnership(&_RMNRemote.TransactOpts)
}

func (_RMNRemote *RMNRemoteTransactor) Curse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "curse", subject)
}

func (_RMNRemote *RMNRemoteSession) Curse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactorSession) Curse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactor) Curse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "curse0", subjects)
}

func (_RMNRemote *RMNRemoteSession) Curse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse0(&_RMNRemote.TransactOpts, subjects)
}

func (_RMNRemote *RMNRemoteTransactorSession) Curse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Curse0(&_RMNRemote.TransactOpts, subjects)
}

func (_RMNRemote *RMNRemoteTransactor) SetConfig(opts *bind.TransactOpts, newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "setConfig", newConfig)
}

func (_RMNRemote *RMNRemoteSession) SetConfig(newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.Contract.SetConfig(&_RMNRemote.TransactOpts, newConfig)
}

func (_RMNRemote *RMNRemoteTransactorSession) SetConfig(newConfig RMNRemoteConfig) (*types.Transaction, error) {
	return _RMNRemote.Contract.SetConfig(&_RMNRemote.TransactOpts, newConfig)
}

func (_RMNRemote *RMNRemoteTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNRemote *RMNRemoteSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNRemote.Contract.TransferOwnership(&_RMNRemote.TransactOpts, to)
}

func (_RMNRemote *RMNRemoteTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNRemote.Contract.TransferOwnership(&_RMNRemote.TransactOpts, to)
}

func (_RMNRemote *RMNRemoteTransactor) Uncurse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "uncurse", subject)
}

func (_RMNRemote *RMNRemoteSession) Uncurse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactorSession) Uncurse(subject [16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse(&_RMNRemote.TransactOpts, subject)
}

func (_RMNRemote *RMNRemoteTransactor) Uncurse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.contract.Transact(opts, "uncurse0", subjects)
}

func (_RMNRemote *RMNRemoteSession) Uncurse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse0(&_RMNRemote.TransactOpts, subjects)
}

func (_RMNRemote *RMNRemoteTransactorSession) Uncurse0(subjects [][16]byte) (*types.Transaction, error) {
	return _RMNRemote.Contract.Uncurse0(&_RMNRemote.TransactOpts, subjects)
}

type RMNRemoteConfigSetIterator struct {
	Event *RMNRemoteConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteConfigSet)
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
		it.Event = new(RMNRemoteConfigSet)
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

func (it *RMNRemoteConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteConfigSet struct {
	Version uint32
	Config  RMNRemoteConfig
	Raw     types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterConfigSet(opts *bind.FilterOpts, version []uint32) (*RMNRemoteConfigSetIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "ConfigSet", versionRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteConfigSetIterator{contract: _RMNRemote.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNRemoteConfigSet, version []uint32) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "ConfigSet", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteConfigSet)
				if err := _RMNRemote.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseConfigSet(log types.Log) (*RMNRemoteConfigSet, error) {
	event := new(RMNRemoteConfigSet)
	if err := _RMNRemote.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteCursedIterator struct {
	Event *RMNRemoteCursed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteCursedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteCursed)
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
		it.Event = new(RMNRemoteCursed)
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

func (it *RMNRemoteCursedIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteCursedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteCursed struct {
	Subjects [][16]byte
	Raw      types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterCursed(opts *bind.FilterOpts) (*RMNRemoteCursedIterator, error) {

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "Cursed")
	if err != nil {
		return nil, err
	}
	return &RMNRemoteCursedIterator{contract: _RMNRemote.contract, event: "Cursed", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchCursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteCursed) (event.Subscription, error) {

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "Cursed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteCursed)
				if err := _RMNRemote.contract.UnpackLog(event, "Cursed", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseCursed(log types.Log) (*RMNRemoteCursed, error) {
	event := new(RMNRemoteCursed)
	if err := _RMNRemote.contract.UnpackLog(event, "Cursed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteOwnershipTransferRequestedIterator struct {
	Event *RMNRemoteOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteOwnershipTransferRequested)
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
		it.Event = new(RMNRemoteOwnershipTransferRequested)
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

func (it *RMNRemoteOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteOwnershipTransferRequestedIterator{contract: _RMNRemote.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteOwnershipTransferRequested)
				if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNRemoteOwnershipTransferRequested, error) {
	event := new(RMNRemoteOwnershipTransferRequested)
	if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteOwnershipTransferredIterator struct {
	Event *RMNRemoteOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteOwnershipTransferred)
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
		it.Event = new(RMNRemoteOwnershipTransferred)
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

func (it *RMNRemoteOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNRemoteOwnershipTransferredIterator{contract: _RMNRemote.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteOwnershipTransferred)
				if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseOwnershipTransferred(log types.Log) (*RMNRemoteOwnershipTransferred, error) {
	event := new(RMNRemoteOwnershipTransferred)
	if err := _RMNRemote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNRemoteUncursedIterator struct {
	Event *RMNRemoteUncursed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNRemoteUncursedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNRemoteUncursed)
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
		it.Event = new(RMNRemoteUncursed)
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

func (it *RMNRemoteUncursedIterator) Error() error {
	return it.fail
}

func (it *RMNRemoteUncursedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNRemoteUncursed struct {
	Subjects [][16]byte
	Raw      types.Log
}

func (_RMNRemote *RMNRemoteFilterer) FilterUncursed(opts *bind.FilterOpts) (*RMNRemoteUncursedIterator, error) {

	logs, sub, err := _RMNRemote.contract.FilterLogs(opts, "Uncursed")
	if err != nil {
		return nil, err
	}
	return &RMNRemoteUncursedIterator{contract: _RMNRemote.contract, event: "Uncursed", logs: logs, sub: sub}, nil
}

func (_RMNRemote *RMNRemoteFilterer) WatchUncursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteUncursed) (event.Subscription, error) {

	logs, sub, err := _RMNRemote.contract.WatchLogs(opts, "Uncursed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNRemoteUncursed)
				if err := _RMNRemote.contract.UnpackLog(event, "Uncursed", log); err != nil {
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

func (_RMNRemote *RMNRemoteFilterer) ParseUncursed(log types.Log) (*RMNRemoteUncursed, error) {
	event := new(RMNRemoteUncursed)
	if err := _RMNRemote.contract.UnpackLog(event, "Uncursed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetVersionedConfig struct {
	Version uint32
	Config  RMNRemoteConfig
}

func (_RMNRemote *RMNRemote) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNRemote.abi.Events["ConfigSet"].ID:
		return _RMNRemote.ParseConfigSet(log)
	case _RMNRemote.abi.Events["Cursed"].ID:
		return _RMNRemote.ParseCursed(log)
	case _RMNRemote.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNRemote.ParseOwnershipTransferRequested(log)
	case _RMNRemote.abi.Events["OwnershipTransferred"].ID:
		return _RMNRemote.ParseOwnershipTransferred(log)
	case _RMNRemote.abi.Events["Uncursed"].ID:
		return _RMNRemote.ParseUncursed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNRemoteConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c")
}

func (RMNRemoteCursed) Topic() common.Hash {
	return common.HexToHash("0x1716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f74")
}

func (RMNRemoteOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNRemoteOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (RMNRemoteUncursed) Topic() common.Hash {
	return common.HexToHash("0x0676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba1")
}

func (_RMNRemote *RMNRemote) Address() common.Address {
	return _RMNRemote.address
}

type RMNRemoteInterface interface {
	GetCursedSubjects(opts *bind.CallOpts) ([][16]byte, error)

	GetLocalChainSelector(opts *bind.CallOpts) (uint64, error)

	GetReportDigestHeader(opts *bind.CallOpts) ([32]byte, error)

	GetVersionedConfig(opts *bind.CallOpts) (GetVersionedConfig,

		error)

	IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error)

	IsCursed0(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Verify(opts *bind.CallOpts, offrampAddress common.Address, merkleRoots []InternalMerkleRoot, signatures []IRMNRemoteSignature) error

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Curse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error)

	Curse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newConfig RMNRemoteConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Uncurse(opts *bind.TransactOpts, subject [16]byte) (*types.Transaction, error)

	Uncurse0(opts *bind.TransactOpts, subjects [][16]byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts, version []uint32) (*RMNRemoteConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNRemoteConfigSet, version []uint32) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*RMNRemoteConfigSet, error)

	FilterCursed(opts *bind.FilterOpts) (*RMNRemoteCursedIterator, error)

	WatchCursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteCursed) (event.Subscription, error)

	ParseCursed(log types.Log) (*RMNRemoteCursed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNRemoteOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNRemoteOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNRemoteOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNRemoteOwnershipTransferred, error)

	FilterUncursed(opts *bind.FilterOpts) (*RMNRemoteUncursedIterator, error)

	WatchUncursed(opts *bind.WatchOpts, sink chan<- *RMNRemoteUncursed) (event.Subscription, error)

	ParseUncursed(log types.Log) (*RMNRemoteUncursed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var RMNRemoteZKBin string = ("0x0002000000000002000c00000000000200010000000103550000006003100270000002140030019d00000214033001970000000100200190000000370000c13d0000008002000039000000400020043f000000040030008c0000030e0000413d000000000401043b000000e0044002700000021d0040009c000000600000213d000002290040009c000000850000213d0000022f0040009c000001070000213d000002320040009c0000026d0000613d000002330040009c0000030e0000c13d000000240030008c0000030e0000413d0000000002000416000000000002004b0000030e0000c13d0000000402100370000000000202043b000600000002001d000002170020009c0000030e0000213d000000060530006a0000024f0050009c0000030e0000213d000000640050008c0000030e0000413d0000000104000039000000000204041a00000237022001970000000006000411000000000026004b0000033f0000c13d00000006020000290000000406200039000000000261034f000000000202043b000000000002004b000003ba0000c13d0000021a01000041000000800010043f00000241010000410000084d000104300000000002000416000000000002004b0000030e0000c13d0000001f023000390000021502200197000000a002200039000000400020043f0000001f0430018f0000021605300198000000a002500039000000480000613d000000a006000039000000000701034f000000007807043c0000000006860436000000000026004b000000440000c13d000000000004004b000000550000613d000000000151034f0000000304400210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000200030008c0000030e0000413d000000a00100043d000002170010009c0000030e0000213d0000000003000411000000000003004b000001210000c13d000000400100043d0000021c020000410000012a0000013d0000021e0040009c000000f00000213d000002240040009c000001300000213d000002270040009c000002810000613d000002280040009c0000030e0000c13d0000000001000416000000000001004b0000030e0000c13d000000000100041a00000237021001970000000006000411000000000026004b000003100000c13d0000000102000039000000000302041a0000021804300197000000000464019f000000000042041b0000021801100197000000000010041b00000000010004140000023705300197000002140010009c0000021401008041000000c0011002100000023d011001c70000800d0200003900000003030000390000024a04000041084b08410000040f00000001002001900000030e0000613d00000000010000190000084c0001042e0000022a0040009c000001810000213d0000022d0040009c000002b50000613d0000022e0040009c0000030e0000c13d000000240030008c0000030e0000413d0000000002000416000000000002004b0000030e0000c13d0000000401100370000000000101043b00000236001001980000030e0000c13d000000c002000039000000400020043f0000000102000039000000800020043f0000023801100197000000a00010043f000000000102041a00000237011001970000000002000411000000000012004b000003140000c13d0000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000023801100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000101041a000000000001004b000004a20000c13d0000000603000039000000000103041a000002170010009c000002ef0000213d0000000102100039000000000023041b0000023c0110009a0000000902000029000000000021041b000000000103041a000800000001001d000000000020043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b0000000802000029000000000021041b0000000a020000290000000102200039000000800100043d000000000012004b000000a00000413d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b000000e60000613d0000008004000039000000000500001900000020044000390000000006040433000002380660019700000000026204360000000105500039000000000035004b000000df0000413d0000000002120049000002140020009c00000214020080410000006002200210000002140010009c00000214010080410000004001100210000000000112019f00000000020004140000025c0000013d0000021f0040009c0000020d0000213d000002220040009c000002cd0000613d000002230040009c0000030e0000c13d0000000001000416000000000001004b0000030e0000c13d0000000001000412000c00000001001d000b00000000003d0000800501000039000000440300003900000000040004150000000c0440008a00000005044002100000024202000041084b08230000040f0000021701100197000000800010043f00000243010000410000084c0001042e000002300040009c000002d80000613d000002310040009c0000030e0000c13d000000240030008c0000030e0000413d0000000002000416000000000002004b0000030e0000c13d0000000401100370000000000101043b00000236001001980000030e0000c13d0000000602000039000000000202041a000000000002004b0000000002000019000003430000c13d000000010120018f000000400200043d0000000000120435000002140020009c000002140200804100000040012002100000025e011001c70000084c0001042e00000217001001980000000102000039000000000402041a0000021804400197000000000334019f000000000032041b000002650000c13d000000400100043d0000021a020000410000000000210435000002140010009c000002140100804100000040011002100000021b011001c70000084d00010430000002250040009c000002f50000613d000002260040009c0000030e0000c13d000000240030008c0000030e0000413d0000000002000416000000000002004b0000030e0000c13d0000000402100370000000000202043b000002170020009c0000030e0000213d0000002304200039000000000034004b0000030e0000813d0000000404200039000000000441034f000000000504043b000002170050009c000002ef0000213d00000005045002100000003f064000390000023406600197000002350060009c000002ef0000213d0000008006600039000000400060043f000000800050043f00000024022000390000000004240019000000000034004b0000030e0000213d000000000005004b0000015d0000613d0000008003000039000000000521034f000000000505043b00000236005001980000030e0000c13d000000200330003900000000005304350000002002200039000000000042004b000001540000413d0000000101000039000000000101041a00000237011001970000000002000411000000000012004b000004090000c13d000000800100043d000000000001004b0000040c0000c13d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b000001770000613d0000008004000039000000000500001900000020044000390000000006040433000002380660019700000000026204360000000105500039000000000035004b000001700000413d0000000002120049000002140020009c00000214020080410000006002200210000002140010009c00000214010080410000004001100210000000000112019f0000000002000414000002040000013d0000022b0040009c000002fe0000613d0000022c0040009c0000030e0000c13d000000240030008c0000030e0000413d0000000002000416000000000002004b0000030e0000c13d0000000401100370000000000101043b00000236001001980000030e0000c13d000000c002000039000000400020043f0000000102000039000000800020043f0000023801100197000000a00010043f000000000102041a00000237011001970000000002000411000000000012004b000003140000c13d0000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000023801100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000301041a000000000003004b0000045e0000613d0000000601000039000000000201041a000000000002004b000004b60000613d000000010130008a000000000023004b000001d20000613d000000000012004b000006e70000a13d000002460130009a000002460220009a000000000202041a000000000021041b000000000020043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039000800000003001d084b08460000040f00000001002001900000030e0000613d000000000101043b0000000802000029000000000021041b0000000601000039000000000301041a000000000003004b000004650000613d000000010130008a000002460230009a000000000002041b0000000602000039000000000012041b0000000901000029000000000010043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000001041b0000000a020000290000000102200039000000800100043d000000000012004b0000019a0000413d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b000001fb0000613d0000008004000039000000000500001900000020044000390000000006040433000002380660019700000000026204360000000105500039000000000035004b000001f40000413d0000000002120049000002140020009c00000214020080410000006002200210000002140010009c00000214010080410000004001100210000000000112019f0000000002000414000002140020009c0000021402008041000000c002200210000000000121019f0000023d011001c70000800d0200003900000001030000390000024804000041000002640000013d000002200040009c000003050000613d000002210040009c0000030e0000c13d000000240030008c0000030e0000413d0000000004000416000000000004004b0000030e0000c13d0000000404100370000000000404043b000002170040009c0000030e0000213d0000002305400039000000000035004b0000030e0000813d0000000405400039000000000551034f000000000605043b000002170060009c000002ef0000213d00000005056002100000003f075000390000023407700197000002350070009c000002ef0000213d0000008007700039000000400070043f000000800060043f00000024044000390000000005450019000000000035004b0000030e0000213d000000000006004b000002390000613d000000000341034f000000000303043b00000236003001980000030e0000c13d000000200220003900000000003204350000002004400039000000000054004b000002300000413d0000000101000039000000000101041a00000237011001970000000002000411000000000012004b000004090000c13d000000800100043d000000000001004b0000046b0000c13d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b000002530000613d0000008004000039000000000500001900000020044000390000000006040433000002380660019700000000026204360000000105500039000000000035004b0000024c0000413d0000000002120049000002140020009c00000214020080410000006002200210000002140010009c00000214010080410000004001100210000000000112019f0000000002000414000002140020009c0000021402008041000000c002200210000000000121019f0000023d011001c70000800d0200003900000001030000390000023e04000041000000800000013d000000800010043f0000014000000443000001600010044300000020010000390000010000100443000001200020044300000219010000410000084c0001042e0000000001000416000000000001004b0000030e0000c13d000000c001000039000000400010043f0000001301000039000000800010043f0000026f01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e002000039084b07f30000040f000000c00110008a000002140010009c0000021401008041000000600110021000000270011001c70000084c0001042e000000640030008c0000030e0000413d0000000002000416000000000002004b0000030e0000c13d0000000402100370000000000202043b000a00000002001d000002370020009c0000030e0000213d0000002402100370000000000202043b000002170020009c0000030e0000213d0000002304200039000000000034004b0000030e0000813d0000000404200039000000000441034f000000000404043b000002170040009c0000030e0000213d000000050440021000000000024200190000002402200039000000000032004b0000030e0000213d0000004402100370000000000202043b000002170020009c0000030e0000213d0000002304200039000000000034004b0000030e0000813d000600040020003d0000000601100360000000000101043b000002170010009c0000030e0000213d000500240020003d00000006021002100000000502200029000000000032004b0000030e0000213d0000000502000039000000000202041a0000021400200198000004b10000c13d0000025701000041000000800010043f00000241010000410000084d000104300000000001000416000000000001004b0000030e0000c13d0000000601000039000000000101041a000000000001004b00000000010000190000035f0000613d0000025a01000041000000000010043f0000000701000039000000200010043f0000025b01000041000000000101041a000000000001004b0000035e0000c13d0000025c01000041000000000010043f0000025d01000041000000000101041a000000000001004b0000000001000039000000010100c0390000035f0000013d0000000001000416000000000001004b0000030e0000c13d0000000602000039000000000102041a000000800010043f000000000020043f000000000001004b000003180000c13d0000002002000039000003230000013d0000000001000416000000000001004b0000030e0000c13d000000800000043f0000006002000039000000a00020043f000000c00000043f0000000501000039000000000301041a0000014004000039000000400040043f0000000201000039000000000101041a000000e00010043f0000000306000039000000000506041a0000025f0050009c000002ef0000813d00000005015002100000003f011000390000023401100197000002600010009c000003630000a13d0000025801000041000000000010043f0000004101000039000000040010043f0000023b010000410000084d000104300000000001000416000000000001004b0000030e0000c13d0000000101000039000000000101041a0000023701100197000000800010043f00000243010000410000084c0001042e0000000001000416000000000001004b0000030e0000c13d0000025201000041000000800010043f00000243010000410000084c0001042e000000240030008c0000030e0000413d0000000002000416000000000002004b0000030e0000c13d0000000401100370000000000601043b000002370060009c000003330000a13d00000000010000190000084d000104300000024901000041000000800010043f00000241010000410000084d000104300000024501000041000000c00010043f00000259010000410000084d00010430000000a004000039000002440200004100000000030000190000000005040019000000000402041a000000000445043600000001022000390000000103300039000000000013004b0000031b0000413d000000600250008a0000008001000039084b07e10000040f000000400100043d000a00000001001d0000008002000039084b08050000040f0000000a020000290000000001210049000002140010009c00000214010080410000006001100210000002140020009c00000214020080410000004002200210000000000121019f0000084c0001042e0000000101000039000000000101041a00000237011001970000000005000411000000000015004b0000033f0000c13d000000000056004b000003ad0000c13d0000024001000041000000800010043f00000241010000410000084d000104300000024501000041000000800010043f00000241010000410000084d000104300000023801100197000000000010043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000101041a000000000001004b000004070000c13d0000025c01000041000000000010043f0000000701000039000000200010043f0000025d01000041000000000101041a000000000001004b0000000002000039000000010200c039000001190000013d0000000101000039000000010110018f000000800010043f00000243010000410000084c0001042e0000014001100039000000400010043f000001400050043f000000000060043f000000000005004b0000037d0000613d000001600600003900000261070000410000000008000019000002620010009c000002ef0000213d0000004009100039000000400090043f000000000907041a000000a00a900270000002170aa00197000000200b1000390000000000ab04350000023709900197000000000091043500000000061604360000000107700039000000400100043d0000000108800039000000000058004b0000036c0000413d0000021403300197000001000040043f0000000404000039000000000404041a0000021704400197000001200040043f00000020041000390000004005000039000000000054043500000000003104350000004003100039000000e00400043d00000000004304350000006004100039000001000300043d0000000000240435000000a00210003900000000040304330000000000420435000000c002100039000000000004004b000003a00000613d00000000050000190000002003300039000000000603043300000000760604340000023706600197000000000662043600000000070704330000021707700197000000000076043500000040022000390000000105500039000000000045004b000003940000413d000001200300043d0000021703300197000000800410003900000000003404350000000002120049000002140020009c00000214020080410000006002200210000002140010009c00000214010080410000004001100210000000000112019f0000084c0001042e000000000100041a0000021801100197000000000161019f000000000010041b0000000001000414000002140010009c0000021401008041000000c0011002100000023d011001c70000800d0200003900000003030000390000023f04000041000000800000013d000400000002001d000900000006001d000800200060003d0000000802100360000000000202043b000000230850008a0000025105200197000302510080019b000000030650014f000000030050006c00000000050000190000025105004041000200000008001d000000000082004b00000000070000190000025107008041000002510060009c000000000507c019000000000005004b0000030e0000c13d0000000906200029000000000561034f000000000505043b000002170050009c0000030e0000213d000000060750021000000000077300490000002008600039000000000078004b0000000009000019000002510900204100000251077001970000025108800197000000000a78013f000000000078004b000000000700001900000251070040410000025100a0009c000000000709c019000000000007004b0000030e0000c13d000000020050008c000003f50000413d00000006074002100000000008670019000000000781034f000000000707043b000002170070009c0000030e0000213d0000004008800039000000000881034f000000000808043b000002170080009c0000030e0000213d000000000087004b000004bc0000813d0000000104400039000000000054004b000003e50000413d00000008040000290000002004400039000000000441034f000000000404043b000002170040009c0000030e0000213d000000010640021000000264046001970000026506600197000000000064004b000004b60000c13d00000001044001bf000000000045004b0000062b0000813d0000026e01000041000000800010043f00000241010000410000084d000104300000000102000039000001190000013d000000400100043d00000245020000410000012a0000013d0000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000023801100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000301041a000000000003004b0000045e0000613d0000000601000039000000000201041a000000000002004b000004b60000613d000000010130008a000000000023004b000004450000613d000000000012004b000006e70000a13d000002460130009a000002460220009a000000000202041a000000000021041b000000000020043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039000800000003001d084b08460000040f00000001002001900000030e0000613d000000000101043b0000000802000029000000000021041b0000000601000039000000000301041a000000000003004b000004650000613d000000010130008a000002460230009a000000000002041b0000000602000039000000000012041b0000000901000029000000000010043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000001041b0000000a020000290000000102200039000000800100043d000000000012004b0000040d0000413d000001660000013d00000080010000390000000a02000029084b08150000040f0000000001010433000000400200043d0000024703000041000004a80000013d0000025801000041000000000010043f0000003101000039000000040010043f0000023b010000410000084d000104300000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000023801100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000101041a000000000001004b000004a20000c13d0000000603000039000000000103041a000002170010009c000002ef0000213d0000000102100039000000000023041b0000023c0110009a0000000902000029000000000021041b000000000103041a000800000001001d000000000020043f0000000701000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b0000000802000029000000000021041b0000000a020000290000000102200039000000800100043d000000000012004b0000046c0000413d000002420000013d00000080010000390000000a02000029084b08150000040f0000000001010433000000400200043d0000023a030000410000000000320435000002380110019700000004032000390000000000130435000002140020009c000002140200804100000040012002100000023b011001c70000084d000104300000000402000039000000000202041a0000021702200197000002170020009c000004c00000c13d0000025801000041000000000010043f0000001101000039000000040010043f0000023b010000410000084d000104300000026301000041000000800010043f00000241010000410000084d00010430000000000021004b000006ed0000a13d0000000201000039000000000101041a000900000001001d0000014001000039000000400010043f0000024c0100004100000000001004430000000001000414000002140010009c0000021401008041000000c0011002100000024d011001c70000800b02000039084b08460000040f0000000100200190000006f10000613d000000000101043b000000800010043f000002420100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000002140010009c0000021401008041000000c0011002100000024e011001c70000800502000039084b08460000040f0000000100200190000006f10000613d000000000101043b0000021701100197000000a00010043f0000000001000410000000c00010043f0000000a01000029000000e00010043f0000000901000029000001000010043f00000001020003670000002401200370000000000101043b000900000001001d0000000401100039000000000112034f000000000101043b000002170010009c000002ef0000213d00000005031002100000003f043000390000023404400197000000400500043d0000000004450019000400000005001d000000000054004b00000000050000390000000105004039000002170040009c000002ef0000213d0000000100500190000002ef0000c13d000000400040043f0000000404000029000000000014043500000009010000290000002405100039000800000053001d000000080050006c000005820000813d0000000001000031000a00000001001d00070024001000920000000409000029000000000152034f000000000101043b000002170010009c0000030e0000213d000000090410002900000007014000690000024f0010009c0000030e0000213d000000a00010008c0000030e0000413d000000400a00043d0000025000a0009c000002ef0000213d000000a001a00039000000400010043f0000002401400039000000000312034f000000000303043b000002170030009c0000030e0000213d000000000c3a0436000000200b1000390000000001b2034f000000000101043b000002170010009c0000030e0000213d000000000441001900000043014000390000000a06000029000000000061004b0000000003000019000002510300804100000251011001970000025106600197000000000861013f000000000061004b00000000010000190000025101004041000002510080009c000000000103c019000000000001004b0000030e0000c13d0000002403400039000000000132034f000000000d01043b0000021700d0009c000002ef0000213d0000001f01d0003900000271011001970000003f011000390000027101100197000000400e00043d00000000061e00190000000000e6004b00000000010000390000000101004039000002170060009c000002ef0000213d0000000100100190000002ef0000c13d000000400060043f0000000001de04360000000004d4001900000044044000390000000a0040006c0000030e0000213d0000002003300039000000000832034f0000027106d0019800000000046100190000055a0000613d000000000308034f000000000f010019000000003703043c000000000f7f043600000000004f004b000005560000c13d0000001f03d00190000005670000613d000000000668034f0000000303300210000000000704043300000000073701cf000000000737022f000000000606043b0000010003300089000000000636022f00000000033601cf000000000373019f00000000003404350000000001d1001900000000000104350000000000ec04350000002001b00039000000000312034f000000000303043b000002170030009c0000030e0000213d0000004004a0003900000000003404350000002001100039000000000312034f000000000303043b000002170030009c0000030e0000213d00000020099000390000006004a0003900000000003404350000002001100039000000000112034f000000000101043b0000008003a0003900000000001304350000000000a904350000002005500039000000080050006c0000050d0000413d0000000401000029000001200010043f000000400200043d0000004001200039000000400300003900000000003104350000002003200039000002520100004100000000001304350000006001200039000000800400043d0000000000410435000000a00100043d000002170110019700000080042000390000000000140435000000c00100043d0000023701100197000000a0042000390000000000140435000000e00100043d0000023701100197000000c0042000390000000000140435000000e001200039000001000400043d00000000004104350000010001200039000000c005000039000001200400043d0000000000510435000001200120003900000000050404330000000000510435000001400620003900000005015002100000000007610019000000000005004b000007880000c13d0000000001270049000000200410008a00000000004204350000001f0110003900000271041001970000000001240019000000000041004b00000000040000390000000104004039000002170010009c000002ef0000213d0000000100400190000002ef0000c13d000000400010043f000002140030009c000002140300804100000040013002100000000002020433000002140020009c00000214020080410000006002200210000000000112019f0000000002000414000002140020009c0000021402008041000000c002200210000000000112019f0000023d011001c70000801002000039084b08460000040f00000001002001900000030e0000613d00000001020003670000000603200360000000000101043b000800000001001d000000000103043b000000000001004b000000830000613d000900000000001d000a00000000001d0000000a01000029000000060110021000000005011000290000002003100039000000000332034f000000000112034f000000000101043b000000000203043b000000400300043d000000600430003900000000002404350000004002300039000000000012043500000020013000390000001b02000039000000000021043500000008010000290000000000130435000000000000043f000002140030009c000002140300804100000040013002100000000002000414000002140020009c0000021402008041000000c002200210000000000112019f00000253011001c70000000102000039084b08460000040f00000060031002700000021403300197000000200030008c000000200500003900000000050340190000002004500190000005fd0000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b000005f90000c13d0000001f055001900000060a0000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f00000000005404350000000100200190000007ba0000613d000000000100043d0000023702100198000007d80000613d000000090020006b000007db0000813d000900000002001d000000000020043f0000000801000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000101041a000000ff00100190000007de0000613d00000001020003670000000601200360000000000101043b0000000a03000029000a00010030003d0000000a0010006b000900090000002d000005d20000413d000000830000013d0000000304000039000000000404041a000000000004004b000006500000613d0000000601000029000700240010003d000a0001004000920000000301000039000000000101041a0000000a0010006b000006e70000813d000002660140009a000000000101041a0000023701100197000000000010043f0000000801000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f00000001002001900000030e0000613d000000000101043b000000000201041a0000027202200197000000000021041b0000000a04000029000000000004004b000006310000c13d000000000300003100000001010003670000000702100360000000000202043b000000230400008a00050006004000720000000504300029000000000042004b0000000005000019000002510500804100000251044001970000025106200197000000000746013f000000000046004b00000000040000190000025104004041000002510070009c000000000405c019000000000004004b0000030e0000c13d0000000604000029000100440040003d000000000a0000190000000905200029000000000451034f000000000404043b000002170040009c0000030e0000213d00000006064002100000000006630049000000200550003900000251076001970000025108500197000000000978013f000000000078004b00000000070000190000025107004041000000000065004b00000000060000190000025106002041000002510090009c000000000706c019000000000007004b0000030e0000c13d00000000004a004b000006f20000813d0007000600a002180000000702500029000000000121034f000000000101043b000002370010009c0000030e0000213d000000000010043f0000000801000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039000a0000000a001d084b08460000040f0000000a0800002900000001002001900000030e0000613d000000000101043b000000000101041a000000ff00100190000007850000c13d00000001010003670000000802100360000000000302043b0000000002000031000000060420006a000000230440008a00000251054001970000025106300197000000000756013f000000000056004b00000000050000190000025105004041000000000043004b00000000040000190000025104008041000002510070009c000000000504c019000000000005004b0000030e0000c13d0000000904300029000000000341034f000000000303043b000002170030009c0000030e0000213d00000006053002100000000005520049000000200240003900000251045001970000025106200197000000000746013f000000000046004b00000000040000190000025104004041000000000052004b00000000050000190000025105002041000002510070009c000000000405c019000000000004004b0000030e0000c13d000000000038004b000006e70000813d0000000702200029000000000121034f000000000101043b000002370010009c0000030e0000213d000000000010043f0000000801000039000000200010043f0000000001000414000002140010009c0000021401008041000000c00110021000000239011001c70000801002000039084b08460000040f0000000a0a00002900000001002001900000030e0000613d000000000101043b000000000201041a000002720220019700000001022001bf000000000021041b00000001010003670000000802100360000000000202043b0000000003000031000000050430002900000251054001970000025106200197000000000756013f000000000056004b00000000050000190000025105002041000000000042004b00000000040000190000025104004041000002510070009c000000000504c019000000010aa00039000000000005004b000006630000c13d0000030e0000013d0000025801000041000000000010043f0000003201000039000000040010043f0000023b010000410000084d000104300000024b01000041000000800010043f00000241010000410000084d00010430000000000001042f00000002060000390000000407000029000000000076041b0000000306000039000000000706041a000000000046041b000000000074004b000007020000813d000002670640009a000002670770009a000000000076004b000007020000813d000000000006041b0000000106600039000000000076004b000006fe0000413d0000000306000039000000000060043f000000000004004b0000071d0000613d00000261060000410000000007000019000000000851034f000000000808043b000002370080009c0000030e0000213d0000002009500039000000000991034f000000000909043b000002170090009c0000030e0000213d000000000a06041a000002680aa00197000000a00990021000000269099001970000000009a9019f000000000889019f000000000086041b000000010660003900000040055000390000000107700039000000000047004b000007080000413d0000000104100360000000000404043b000002170040009c0000030e0000213d0000000405000039000000000605041a0000026a06600197000000000646019f000000000065041b0000000506000039000000000506041a0000021407500197000002140070009c000004b60000613d0000026b0750019700000001055000390000021405500197000000000775019f000000000076041b000000400600043d000000200760003900000004080000290000000000870435000000020020006c000000000700001900000251070080410000025108200197000000030980014f000000030080006c00000000080000190000025108004041000002510090009c000000000807c01900000020070000390000000000760435000000000008004b0000030e0000c13d00000006072000290000000402700039000000000221034f000000000202043b000002170020009c0000030e0000213d000000240770003900000006082002100000000003830049000000000037004b0000000008000019000002510800204100000251033001970000025109700197000000000a39013f000000000039004b000000000300001900000251030040410000025100a0009c000000000308c019000000000003004b0000030e0000c13d00000080036000390000004008600039000000600900003900000000009804350000000000230435000000a003600039000000000002004b000007710000613d0000000008000019000000000971034f000000000909043b000002370090009c0000030e0000213d0000000009930436000000200a700039000000000aa1034f000000000a0a043b0000021700a0009c0000030e0000213d0000000000a90435000000400770003900000040033000390000000108800039000000000028004b000007610000413d000000600160003900000000004104350000000001630049000002140010009c00000214010080410000006001100210000002140060009c00000214060080410000004002600210000000000121019f0000000002000414000002140020009c0000021402008041000000c002200210000000000121019f0000023d011001c70000800d0200003900000002030000390000026c04000041000002640000013d000000400100043d0000026d020000410000012a0000013d000000a0080000390000000009000019000007a10000013d0000000001bc001900000000000104350000004001a0003900000000010104330000021701100197000000400d70003900000000001d04350000006001a0003900000000010104330000021701100197000000600d70003900000000001d043500000080017000390000008007a00039000000000707043300000000007104350000001f01b00039000002710110019700000000071c00190000000109900039000000000059004b000005a90000813d0000000001270049000001400110008a00000000061604360000002004400039000000000a04043300000000b10a043400000217011001970000000001170436000000000b0b04330000000000810435000000a00170003900000000db0b04340000000000b10435000000c00c70003900000000000b004b0000078b0000613d000000000e0000190000000001ce0019000000000fed0019000000000f0f04330000000000f10435000000200ee000390000000000be004b000007b20000413d0000078b0000013d0000001f0530018f0000021606300198000000400200043d0000000004620019000007c50000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000007c10000c13d000000000005004b000007d20000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000002140020009c00000214020080410000004002200210000000000112019f0000084d00010430000000400100043d00000256020000410000012a0000013d000000400100043d00000254020000410000012a0000013d000000400100043d00000255020000410000012a0000013d0000001f0220003900000271022001970000000001120019000000000021004b00000000020000390000000102004039000002170010009c000007ed0000213d0000000100200190000007ed0000c13d000000400010043f000000000001042d0000025801000041000000000010043f0000004101000039000000040010043f0000023b010000410000084d0001043000000000430104340000000001320436000000000003004b000007ff0000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b000007f80000413d000000000231001900000000000204350000001f0230003900000271022001970000000001210019000000000001042d00000020030000390000000004310436000000000302043300000000003404350000004001100039000000000003004b000008140000613d000000000400001900000020022000390000000005020433000002380550019700000000015104360000000104400039000000000034004b0000080d0000413d000000000001042d0000000003010433000000000023004b0000081c0000a13d000000050220021000000000012100190000002001100039000000000001042d0000025801000041000000000010043f0000003201000039000000040010043f0000023b010000410000084d00010430000000000001042f00000000050100190000000000200443000000050030008c000008310000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b000008290000413d000002140030009c000002140300804100000060013002100000000002000414000002140020009c0000021402008041000000c002200210000000000112019f00000273011001c70000000002050019084b08460000040f0000000100200190000008400000613d000000000101043b000000000001042d000000000001042f00000844002104210000000102000039000000000001042d0000000002000019000000000001042d00000849002104230000000102000039000000000001042d0000000002000019000000000001042d0000084b000004320000084c0001042e0000084d000104300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000002000000000000000000000000000000800000010000000000000000009cf8540c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000009b15e16f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000070a9089d00000000000000000000000000000000000000000000000000000000d881e09100000000000000000000000000000000000000000000000000000000f2fde38a00000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000f8bb876e00000000000000000000000000000000000000000000000000000000d881e09200000000000000000000000000000000000000000000000000000000eaa83ddd000000000000000000000000000000000000000000000000000000008da5cb5a000000000000000000000000000000000000000000000000000000008da5cb5b000000000000000000000000000000000000000000000000000000009a19b3290000000000000000000000000000000000000000000000000000000070a9089e0000000000000000000000000000000000000000000000000000000079ba509700000000000000000000000000000000000000000000000000000000397796f6000000000000000000000000000000000000000000000000000000006509a953000000000000000000000000000000000000000000000000000000006509a954000000000000000000000000000000000000000000000000000000006d2d399300000000000000000000000000000000000000000000000000000000397796f70000000000000000000000000000000000000000000000000000000062eed415000000000000000000000000000000000000000000000000000000001add205e000000000000000000000000000000000000000000000000000000001add205f000000000000000000000000000000000000000000000000000000002cbc26bb00000000000000000000000000000000000000000000000000000000181f5a7700000000000000000000000000000000000000000000000000000000198f0f777fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffff7f00000000000000000000000000000000ffffffffffffffffffffffffffffffff000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000020000000000000000000000000000000000004000000000000000000000000019d5c79b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002400000000000000000000000009addddcec1d7ba6ad726df49aeea3e93fb0c1037d551236841a60c0c883f2c102000000000000000000000000000000000000000000000000000000000000001716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f74ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e0000000000000000000000000000000000000020000000800000000000000000f652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f2b5c74de0000000000000000000000000000000000000000000000000000000009addddcec1d7ba6ad726df49aeea3e93fb0c1037d551236841a60c0c883f2c273281fa1000000000000000000000000000000000000000000000000000000000676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba102b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e059fa4a93000000000000000000000000000000000000000000000000000000009a8a0592ac89c5ad3bc6df8224c17b485976f597df104ee20d0df415241f670b020000020000000000000000000000000000000400000000000000000000000002000002000000000000000000000000000000440000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff5f80000000000000000000000000000000000000000000000000000000000000009651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf530000000000000000000000000000000000000080000000000000000000000000bbe15e7f00000000000000000000000000000000000000000000000000000000aaaa9141000000000000000000000000000000000000000000000000000000008baa579f00000000000000000000000000000000000000000000000000000000ace124bc000000000000000000000000000000000000000000000000000000004e487b71000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000c0000000000000000001000000000000000000000000000000000000000000000000000000000000007dde556524061d0ce70b736a6e842a48e4927608bf87fd31432ced12a03ffeb8010000000000000000000000000000010000000000000000000000000000000070b766b11586b6b505ed3893938b0cc6c6c98bd6f65e969ac311168d34e4f9e200000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000fffffffffffffebfc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b000000000000000000000000000000000000000000000000ffffffffffffffbf4485151700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffffe000000000000000000000000000000000000000000000001fffffffffffffffe3da8a5f161a6c3ff06a60736d0ed24d7963cc6a5c4fafd2fa1dae9bb908e07a63da8a5f161a6c3ff06a60736d0ed24d7963cc6a5c4fafd2fa1dae9bb908e07a5ffffffff0000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000007f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c28cae27d00000000000000000000000000000000000000000000000000000000014c502000000000000000000000000000000000000000000000000000000000524d4e52656d6f746520312e362e302d646576000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0002000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
