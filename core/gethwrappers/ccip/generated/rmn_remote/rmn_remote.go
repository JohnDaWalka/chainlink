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

type IRMNTaggedRoot struct {
	CommitStore common.Address
	Root        [32]byte
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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMN\",\"name\":\"legacyRMN\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"AlreadyCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOnchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignerOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IsBlessedNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"NotCursed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfOrderSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ThresholdNotMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroValueNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Cursed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"Uncursed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"curse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCursedSubjects\",\"outputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLocalChainSelector\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"localChainSelector\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReportDigestHeader\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"digestHeader\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersionedConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMN.TaggedRoot\",\"name\":\"taggedRoot\",\"type\":\"tuple\"}],\"name\":\"isBlessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isCursed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"rmnHomeContractConfigDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"onchainPublicKey\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"nodeIndex\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Signer[]\",\"name\":\"signers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"}],\"internalType\":\"structRMNRemote.Config\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"subject\",\"type\":\"bytes16\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"subjects\",\"type\":\"bytes16[]\"}],\"name\":\"uncurse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offrampAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"verify\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b506040516200230438038062002304833981016040819052620000349162000150565b336000816200005657604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038481169190911790915581161562000089576200008981620000d6565b5050816001600160401b0316600003620000b65760405163273e150360e21b815260040160405180910390fd5b6001600160401b039091166080526001600160a01b031660a052620001a5565b336001600160a01b038216036200010057604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080604083850312156200016457600080fd5b82516001600160401b03811681146200017c57600080fd5b60208401519092506001600160a01b03811681146200019a57600080fd5b809150509250929050565b60805160a05161212b620001d9600039600081816108c5015261096d0152600081816102a80152610b7c015261212b6000f3fe608060405234801561001057600080fd5b506004361061011b5760003560e01c80636d2d3993116100b25780639a19b32911610081578063eaa83ddd11610066578063eaa83ddd1461029a578063f2fde38b146102d2578063f8bb876e146102e557600080fd5b80639a19b32914610272578063d881e0921461028557600080fd5b80636d2d39931461021c57806370a9089e1461022f57806379ba5097146102425780638da5cb5b1461024a57600080fd5b8063397796f7116100ee578063397796f7146101c05780634d616771146101c857806362eed415146101db5780636509a954146101ee57600080fd5b8063181f5a7714610120578063198f0f77146101725780631add205f146101875780632cbc26bb1461019d575b600080fd5b61015c6040518060400160405280601381526020017f524d4e52656d6f746520312e362e302d6465760000000000000000000000000081525081565b60405161016991906114d9565b60405180910390f35b6101856101803660046114ec565b6102f8565b005b61018f6106f2565b604051610169929190611527565b6101b06101ab366004611605565b6107ea565b6040519015158152602001610169565b6101b0610847565b6101b06101d6366004611620565b6108c1565b6101856101e9366004611605565b6109e3565b6040517f9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf538152602001610169565b61018561022a366004611605565b610a57565b61018561023d3660046116a6565b610ac7565b610185610e22565b60015460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610169565b610185610280366004611825565b610ef0565b61028d610ff6565b60405161016991906118c2565b60405167ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610169565b6101856102e0366004611928565b611002565b6101856102f3366004611825565b611016565b610300611108565b8035610338576040517f9cf8540c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015b6103486020830183611945565b90508110156104185761035e6020830183611945565b8281811061036e5761036e6119ad565b905060400201602001602081019061038691906119fd565b67ffffffffffffffff1661039d6020840184611945565b6103a8600185611a49565b8181106103b7576103b76119ad565b90506040020160200160208101906103cf91906119fd565b67ffffffffffffffff1610610410576040517f4485151700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60010161033b565b5061042960608201604083016119fd565b610434906002611a5c565b61043f906001611a88565b67ffffffffffffffff166104566020830183611945565b90501015610490576040517f014c502000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6003545b8015610522576008600060036104ab600185611a49565b815481106104bb576104bb6119ad565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905561051b81611aa9565b9050610494565b5060005b6105336020830183611945565b9050811015610668576008600061054d6020850185611945565b8481811061055d5761055d6119ad565b6105739260206040909202019081019150611928565b73ffffffffffffffffffffffffffffffffffffffff16815260208101919091526040016000205460ff16156105d4576040517f28cae27d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600860006105e76020860186611945565b858181106105f7576105f76119ad565b61060d9260206040909202019081019150611928565b73ffffffffffffffffffffffffffffffffffffffff168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055600101610526565b508060026106768282611b97565b5050600580546000919082906106919063ffffffff16611cd2565b91906101000a81548163ffffffff021916908363ffffffff160217905590508063ffffffff167f7f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c836040516106e69190611cf5565b60405180910390a25050565b6040805160608082018352600080835260208301919091529181018290526005546040805160608101825260028054825260038054845160208281028201810190965281815263ffffffff9096169592948593818601939092909160009084015b828210156107c1576000848152602090819020604080518082019091529084015473ffffffffffffffffffffffffffffffffffffffff8116825274010000000000000000000000000000000000000000900467ffffffffffffffff1681830152825260019092019101610753565b505050908252506002919091015467ffffffffffffffff16602090910152919491935090915050565b60006107f6600661115b565b60000361080557506000919050565b610810600683611165565b80610841575061084160067f0100000000000000000000000000000100000000000000000000000000000000611165565b92915050565b6000610853600661115b565b6000036108605750600090565b61088b60067f0100000000000000000000000000000000000000000000000000000000000000611165565b806108bc57506108bc60067f0100000000000000000000000000000100000000000000000000000000000000611165565b905090565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610930576040517f0a7c4edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f4d61677100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690634d616771906109a2908590600401611dff565b602060405180830381865afa1580156109bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108419190611e38565b604080516001808252818301909252600091602080830190803683370190505090508181600081518110610a1957610a196119ad565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660209283029190910190910152610a5381611016565b5050565b604080516001808252818301909252600091602080830190803683370190505090508181600081518110610a8d57610a8d6119ad565b7fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921660209283029190910190910152610a5381610ef0565b60055463ffffffff16600003610b09576040517face124bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600454610b219067ffffffffffffffff166001611a88565b67ffffffffffffffff16811015610b64576040517f59fa4a9300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160c08101825246815267ffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166020820152309181019190915273ffffffffffffffffffffffffffffffffffffffff8616606082015260025460808201526000907f9651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf539060a08101610c008789611e5a565b9052604051610c13929190602001611fba565b60405160208183030381529060405280519060200120905060008060005b84811015610e1757600184601b888885818110610c5057610c506119ad565b90506040020160000135898986818110610c6c57610c6c6119ad565b9050604002016020013560405160008152602001604052604051610cac949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610cce573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015192505073ffffffffffffffffffffffffffffffffffffffff8216610d46576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1610610dab576040517fbbe15e7f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526008602052604090205460ff16610e0a576040517faaaa914100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9091508190600101610c31565b505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e73576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610ef8611108565b60005b8151811015610fbb57610f31828281518110610f1957610f196119ad565b602002602001015160066111a390919063ffffffff16565b610fb357818181518110610f4757610f476119ad565b60200260200101516040517f73281fa1000000000000000000000000000000000000000000000000000000008152600401610faa91907fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b60405180910390fd5b600101610efb565b507f0676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba181604051610feb91906118c2565b60405180910390a150565b60606108bc60066111d1565b61100a611108565b611013816111de565b50565b61101e611108565b60005b81518110156110d85761105782828151811061103f5761103f6119ad565b602002602001015160066112a290919063ffffffff16565b6110d05781818151811061106d5761106d6119ad565b60200260200101516040517f19d5c79b000000000000000000000000000000000000000000000000000000008152600401610faa91907fffffffffffffffffffffffffffffffff0000000000000000000000000000000091909116815260200190565b600101611021565b507f1716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f7481604051610feb91906118c2565b60015473ffffffffffffffffffffffffffffffffffffffff163314611159576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000610841825490565b7fffffffffffffffffffffffffffffffff000000000000000000000000000000008116600090815260018301602052604081205415155b9392505050565b600061119c837fffffffffffffffffffffffffffffffff0000000000000000000000000000000084166112d0565b6060600061119c836113ca565b3373ffffffffffffffffffffffffffffffffffffffff82160361122d576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061119c837fffffffffffffffffffffffffffffffff000000000000000000000000000000008416611426565b600081815260018301602052604081205480156113b95760006112f4600183611a49565b855490915060009061130890600190611a49565b905080821461136d576000866000018281548110611328576113286119ad565b906000526020600020015490508087600001848154811061134b5761134b6119ad565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061137e5761137e6120ef565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610841565b6000915050610841565b5092915050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561141a57602002820191906000526020600020905b815481526020019060010190808311611406575b50505050509050919050565b600081815260018301602052604081205461146d57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610841565b506000610841565b6000815180845260005b8181101561149b5760208185018101518683018201520161147f565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152600061119c6020830184611475565b6000602082840312156114fe57600080fd5b813567ffffffffffffffff81111561151557600080fd5b82016060818503121561119c57600080fd5b63ffffffff831681526040602080830182905283518383015283810151606080850152805160a085018190526000939291820190849060c08701905b808310156115ac578351805173ffffffffffffffffffffffffffffffffffffffff16835285015167ffffffffffffffff1685830152928401926001929092019190850190611563565b50604088015167ffffffffffffffff81166080890152945098975050505050505050565b80357fffffffffffffffffffffffffffffffff000000000000000000000000000000008116811461160057600080fd5b919050565b60006020828403121561161757600080fd5b61119c826115d0565b60006040828403121561163257600080fd5b50919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461101357600080fd5b60008083601f84011261166c57600080fd5b50813567ffffffffffffffff81111561168457600080fd5b6020830191508360208260061b850101111561169f57600080fd5b9250929050565b6000806000806000606086880312156116be57600080fd5b85356116c981611638565b9450602086013567ffffffffffffffff808211156116e657600080fd5b818801915088601f8301126116fa57600080fd5b81358181111561170957600080fd5b8960208260051b850101111561171e57600080fd5b60208301965080955050604088013591508082111561173c57600080fd5b506117498882890161165a565b969995985093965092949392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff811182821017156117ac576117ac61175a565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156117f9576117f961175a565b604052919050565b600067ffffffffffffffff82111561181b5761181b61175a565b5060051b60200190565b6000602080838503121561183857600080fd5b823567ffffffffffffffff81111561184f57600080fd5b8301601f8101851361186057600080fd5b803561187361186e82611801565b6117b2565b81815260059190911b8201830190838101908783111561189257600080fd5b928401925b828410156118b7576118a8846115d0565b82529284019290840190611897565b979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561191c5783517fffffffffffffffffffffffffffffffff0000000000000000000000000000000016835292840192918401916001016118de565b50909695505050505050565b60006020828403121561193a57600080fd5b813561119c81611638565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261197a57600080fd5b83018035915067ffffffffffffffff82111561199557600080fd5b6020019150600681901b360382131561169f57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b67ffffffffffffffff8116811461101357600080fd5b8035611600816119dc565b600060208284031215611a0f57600080fd5b813561119c816119dc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561084157610841611a1a565b67ffffffffffffffff818116838216028082169190828114611a8057611a80611a1a565b505092915050565b67ffffffffffffffff8181168382160190808211156113c3576113c3611a1a565b600081611ab857611ab8611a1a565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b60008135610841816119dc565b8135611af681611638565b73ffffffffffffffffffffffffffffffffffffffff811690508154817fffffffffffffffffffffffff000000000000000000000000000000000000000082161783556020840135611b46816119dc565b7bffffffffffffffff00000000000000000000000000000000000000008160a01b16837fffffffff000000000000000000000000000000000000000000000000000000008416171784555050505050565b81358155600180820160208401357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1853603018112611bd557600080fd5b8401803567ffffffffffffffff811115611bee57600080fd5b6020820191508060061b3603821315611c0657600080fd5b68010000000000000000811115611c1f57611c1f61175a565b825481845580821015611c54576000848152602081208381019083015b80821015611c505782825590870190611c3c565b5050505b50600092835260208320925b81811015611c8457611c728385611aeb565b92840192604092909201918401611c60565b5050505050610a53611c9860408401611ade565b6002830167ffffffffffffffff82167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008254161781555050565b600063ffffffff808316818103611ceb57611ceb611a1a565b6001019392505050565b6000602080835260808301843582850152818501357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1863603018112611d3a57600080fd5b8501828101903567ffffffffffffffff80821115611d5757600080fd5b8160061b3603831315611d6957600080fd5b6040606060408901528483865260a089019050849550600094505b83851015611dd4578535611d9781611638565b73ffffffffffffffffffffffffffffffffffffffff16815285870135611dbc816119dc565b83168188015294810194600194909401938101611d84565b611de060408b016119f2565b67ffffffffffffffff811660608b015296509998505050505050505050565b604081018235611e0e81611638565b73ffffffffffffffffffffffffffffffffffffffff81168352506020830135602083015292915050565b600060208284031215611e4a57600080fd5b8151801515811461119c57600080fd5b6000611e6861186e84611801565b80848252602080830192508560051b850136811115611e8657600080fd5b855b81811015611fae57803567ffffffffffffffff80821115611ea95760008081fd5b818901915060a08236031215611ebf5760008081fd5b611ec7611789565b8235611ed2816119dc565b81528286013582811115611ee65760008081fd5b8301601f3681830112611ef95760008081fd5b813584811115611f0b57611f0b61175a565b611f3a897fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084840116016117b2565b94508085523689828501011115611f5357600091508182fd5b808984018a8701376000898287010152505050818682015260409150611f7a8284016119f2565b8282015260609150611f8d8284016119f2565b91810191909152608091820135918101919091528552938201938201611e88565b50919695505050505050565b60006040848352602060408185015261010084018551604086015281860151606067ffffffffffffffff808316606089015260408901519250608073ffffffffffffffffffffffffffffffffffffffff80851660808b015260608b0151945060a081861660a08c015260808c015160c08c015260a08c0151955060c060e08c015286915085518088526101209750878c019250878160051b8d01019750888701965060005b818110156120dc577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee08d8a030184528751868151168a528a810151848c8c01526120ab858c0182611475565b828e015189168c8f01528983015189168a8d0152918701519a87019a909a529850968901969289019260010161205f565b50969d9c50505050505050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var RMNRemoteABI = RMNRemoteMetaData.ABI

var RMNRemoteBin = RMNRemoteMetaData.Bin

func DeployRMNRemote(auth *bind.TransactOpts, backend bind.ContractBackend, localChainSelector uint64, legacyRMN common.Address) (common.Address, *generated.Transaction, *RMNRemote, error) {
	parsed, err := RMNRemoteMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(RMNRemoteZKBin), backend, localChainSelector, legacyRMN)
		contractReturn := &RMNRemote{address: address, abi: *parsed, RMNRemoteCaller: RMNRemoteCaller{contract: contractBind}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contractBind}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNRemoteBin), backend, localChainSelector, legacyRMN)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, HashZks: tx.Hash()}, &RMNRemote{address: address, abi: *parsed, RMNRemoteCaller: RMNRemoteCaller{contract: contract}, RMNRemoteTransactor: RMNRemoteTransactor{contract: contract}, RMNRemoteFilterer: RMNRemoteFilterer{contract: contract}}, nil
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

func (_RMNRemote *RMNRemoteCaller) IsBlessed(opts *bind.CallOpts, taggedRoot IRMNTaggedRoot) (bool, error) {
	var out []interface{}
	err := _RMNRemote.contract.Call(opts, &out, "isBlessed", taggedRoot)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_RMNRemote *RMNRemoteSession) IsBlessed(taggedRoot IRMNTaggedRoot) (bool, error) {
	return _RMNRemote.Contract.IsBlessed(&_RMNRemote.CallOpts, taggedRoot)
}

func (_RMNRemote *RMNRemoteCallerSession) IsBlessed(taggedRoot IRMNTaggedRoot) (bool, error) {
	return _RMNRemote.Contract.IsBlessed(&_RMNRemote.CallOpts, taggedRoot)
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

	IsBlessed(opts *bind.CallOpts, taggedRoot IRMNTaggedRoot) (bool, error)

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

var RMNRemoteZKBin = ("0x0x0003000000000002000c000000000002000200000001035500000000030100190000006003300270000002390030019d00000239033001970000000100200190000000670000c13d0000008002000039000000400020043f000000040030008c0000008d0000413d000000000401043b000000e004400270000002430040009c0000008f0000a13d000002440040009c000000b60000a13d000002450040009c0000011b0000a13d000002460040009c0000016b0000613d000002470040009c0000017c0000613d000002480040009c0000008d0000c13d000000240030008c0000008d0000413d0000000004000416000000000004004b0000008d0000c13d0000000404100370000000000404043b0000023c0040009c0000008d0000213d0000002305400039000000000035004b0000008d0000813d0000000405400039000000000551034f000000000605043b0000023c0060009c000002a90000213d00000005056002100000003f075000390000025b077001970000025c0070009c000002a90000213d0000008007700039000000400070043f000000800060043f00000024044000390000000005450019000000000035004b0000008d0000213d000000000006004b000000430000613d000000000341034f000000000303043b0000025d003001980000008d0000c13d000000200220003900000000003204350000002004400039000000000054004b0000003a0000413d0000000101000039000000000101041a0000023d011001970000000002000411000000000012004b0000048e0000c13d000000800100043d000000000001004b000004af0000c13d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b0000005d0000613d00000080040000390000000005000019000000200440003900000000060404330000025e0660019700000000026204360000000105500039000000000035004b000000560000413d0000000002120049000002390020009c00000239020080410000006002200210000002390010009c00000239010080410000004001100210000000000112019f0000000002000414000003300000013d0000000002000416000000000002004b0000008d0000c13d0000001f023000390000023a02200197000000c002200039000000400020043f0000001f0430018f0000023b05300198000000c002500039000000780000613d000000c006000039000000000701034f000000007807043c0000000006860436000000000026004b000000740000c13d000000000004004b000000850000613d000000000151034f0000000304400210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000400030008c0000008d0000413d000000c00100043d0000023c0010009c0000008d0000213d000000e00200043d0000023d0020009c000001150000a13d0000000001000019000008e000010430000002500040009c000000f00000213d000002560040009c0000012a0000213d000002590040009c000001910000613d0000025a0040009c0000008d0000c13d000000240030008c0000008d0000413d0000000002000416000000000002004b0000008d0000c13d0000000402100370000000000202043b000600000002001d0000023c0020009c0000008d0000213d000000060530006a000002750050009c0000008d0000213d000000640050008c0000008d0000413d0000000104000039000000000204041a0000023d022001970000000006000411000000000026004b000003680000c13d00000006020000290000000406200039000000000261034f000000000202043b000000000002004b000003fa0000c13d0000024001000041000000800010043f0000026701000041000008e0000104300000024b0040009c000001440000213d0000024e0040009c000001a50000613d0000024f0040009c0000008d0000c13d000000640030008c0000008d0000413d0000000002000416000000000002004b0000008d0000c13d0000000402100370000000000202043b000a00000002001d0000023d0020009c0000008d0000213d0000002402100370000000000202043b0000023c0020009c0000008d0000213d0000002304200039000000000034004b0000008d0000813d0000000404200039000000000441034f000000000404043b0000023c0040009c0000008d0000213d000000050440021000000000024200190000002402200039000000000032004b0000008d0000213d0000004402100370000000000202043b0000023c0020009c0000008d0000213d0000002304200039000000000034004b0000008d0000813d000600040020003d0000000601100360000000000101043b0000023c0010009c0000008d0000213d000500240020003d00000006021002100000000502200029000000000032004b0000008d0000213d0000000502000039000000000202041a0000023900200198000005540000c13d0000027d01000041000000800010043f0000026701000041000008e000010430000002510040009c000001510000213d000002540040009c000002250000613d000002550040009c0000008d0000c13d000000440030008c0000008d0000413d0000000001000416000000000001004b0000008d0000c13d0000026801000041000000000010044300000000010004120000000400100443000000200100003900000024001004430000000001000414000002390010009c0000023901008041000000c00110021000000274011001c7000080050200003908de08d90000040f0000000100200190000006d00000613d000000400300043d000000000101043b0000023d021001980000036c0000c13d00000283010000410000000000130435000002390030009c0000023903008041000000400130021000000241011001c7000008e0000104300000000003000411000000000003004b0000015c0000c13d000000400100043d0000024202000041000001650000013d000002490040009c0000023d0000613d0000024a0040009c0000008d0000c13d0000000001000416000000000001004b0000008d0000c13d0000000602000039000000000102041a000000800010043f000000000020043f000000000001004b000003450000c13d0000002002000039000003500000013d000002570040009c000002920000613d000002580040009c0000008d0000c13d000000240030008c0000008d0000413d0000000002000416000000000002004b0000008d0000c13d0000000401100370000000000101043b0000025d001001980000008d0000c13d0000000602000039000000000202041a000000000002004b0000000002000019000003830000c13d000000010120018f000000400200043d0000000000120435000002390020009c0000023902008041000000400120021000000282011001c7000008df0001042e0000024c0040009c000002af0000613d0000024d0040009c0000008d0000c13d0000000001000416000000000001004b0000008d0000c13d0000000101000039000000000101041a0000023d01100197000000800010043f0000026901000041000008df0001042e000002520040009c000002cc0000613d000002530040009c0000008d0000c13d0000000001000416000000000001004b0000008d0000c13d0000027801000041000000800010043f0000026901000041000008df0001042e0000023c001001980000000104000039000000000504041a0000023e05500197000000000335019f000000000034041b000003390000c13d000000400100043d00000240020000410000000000210435000002390010009c0000023901008041000000400110021000000241011001c7000008e0000104300000000001000416000000000001004b0000008d0000c13d0000000001000412000c00000001001d000b00000000003d0000800501000039000000440300003900000000040004150000000c0440008a0000000504400210000002680200004108de08b60000040f0000023c01100197000000800010043f0000026901000041000008df0001042e000000240030008c0000008d0000413d0000000002000416000000000002004b0000008d0000c13d0000000401100370000000000601043b0000023d0060009c0000008d0000213d0000000101000039000000000101041a0000023d011001970000000005000411000000000015004b000003680000c13d000000000056004b000003ed0000c13d0000026601000041000000800010043f0000026701000041000008e0000104300000000001000416000000000001004b0000008d0000c13d000000c001000039000000400010043f0000001301000039000000800010043f0000029801000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e00200003908de08860000040f000000c00110008a000002390010009c0000023901008041000000600110021000000299011001c7000008df0001042e000000240030008c0000008d0000413d0000000002000416000000000002004b0000008d0000c13d0000000401100370000000000101043b0000025d001001980000008d0000c13d000000c002000039000000400020043f0000000102000039000000800020043f0000025e01100197000000a00010043f000000000102041a0000023d011001970000000002000411000000000012004b000003640000c13d0000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000025e01100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000301041a000000000003004b0000053f0000613d0000000601000039000000000201041a000000000002004b000005590000613d000000010130008a000000000023004b000001f20000613d000000000012004b0000078d0000a13d0000026c0130009a0000026c0220009a000000000202041a000000000021041b000000000020043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c70000801002000039000800000003001d08de08d90000040f00000001002001900000008d0000613d000000000101043b0000000802000029000000000021041b0000000601000039000000000301041a000000000003004b0000054e0000613d000000010130008a0000026c0230009a000000000002041b0000000602000039000000000012041b0000000901000029000000000010043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000001041b0000000a020000290000000102200039000000800100043d000000000012004b000001ba0000413d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b0000021b0000613d00000080040000390000000005000019000000200440003900000000060404330000025e0660019700000000026204360000000105500039000000000035004b000002140000413d0000000002120049000002390020009c00000239020080410000006002200210000002390010009c00000239010080410000004001100210000000000112019f0000000002000414000002890000013d0000000001000416000000000001004b0000008d0000c13d0000000601000039000000000101041a000000000001004b00000000010000190000039f0000613d0000028401000041000000000010043f0000000701000039000000200010043f0000028501000041000000000101041a000000000001004b0000039e0000c13d0000028601000041000000000010043f0000028701000041000000000101041a000000000001004b0000000001000039000000010100c0390000039f0000013d000000240030008c0000008d0000413d0000000002000416000000000002004b0000008d0000c13d0000000402100370000000000202043b0000023c0020009c0000008d0000213d0000002304200039000000000034004b0000008d0000813d0000000404200039000000000441034f000000000504043b0000023c0050009c000002a90000213d00000005045002100000003f064000390000025b066001970000025c0060009c000002a90000213d0000008006600039000000400060043f000000800050043f00000024022000390000000004240019000000000034004b0000008d0000213d000000000005004b000002660000613d0000008003000039000000000521034f000000000505043b0000025d005001980000008d0000c13d000000200330003900000000005304350000002002200039000000000042004b0000025d0000413d0000000101000039000000000101041a0000023d011001970000000002000411000000000012004b0000048e0000c13d000000800100043d000000000001004b000004ed0000c13d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b000002800000613d00000080040000390000000005000019000000200440003900000000060404330000025e0660019700000000026204360000000105500039000000000035004b000002790000413d0000000002120049000002390020009c00000239020080410000006002200210000002390010009c00000239010080410000004001100210000000000112019f0000000002000414000002390020009c0000023902008041000000c002200210000000000121019f00000263011001c70000800d0200003900000001030000390000026e04000041000003380000013d0000000001000416000000000001004b0000008d0000c13d000000800000043f0000006002000039000000a00020043f000000c00000043f0000000501000039000000000301041a0000014004000039000000400040043f0000000201000039000000000101041a000000e00010043f0000000306000039000000000506041a000002880050009c000002a90000813d00000005015002100000003f011000390000025b01100197000002890010009c000003a30000a13d0000027e01000041000000000010043f0000004101000039000000040010043f0000026101000041000008e0000104300000000001000416000000000001004b0000008d0000c13d000000000100041a0000023d021001970000000006000411000000000026004b000003600000c13d0000000102000039000000000302041a0000023e04300197000000000464019f000000000042041b0000023e01100197000000000010041b00000000010004140000023d05300197000002390010009c0000023901008041000000c00110021000000263011001c70000800d020000390000000303000039000002700400004108de08d40000040f00000001002001900000008d0000613d0000000001000019000008df0001042e000000240030008c0000008d0000413d0000000002000416000000000002004b0000008d0000c13d0000000401100370000000000101043b0000025d001001980000008d0000c13d000000c002000039000000400020043f0000000102000039000000800020043f0000025e01100197000000a00010043f000000000102041a0000023d011001970000000002000411000000000012004b000003640000c13d0000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000025e01100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000101041a000000000001004b000004e60000c13d0000000603000039000000000103041a0000023c0010009c000002a90000213d0000000102100039000000000023041b000002620110009a0000000902000029000000000021041b000000000103041a000800000001001d000000000020043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b0000000802000029000000000021041b0000000a020000290000000102200039000000800100043d000000000012004b000002e10000413d000000400100043d00000020020000390000000002210436000000800300043d00000000003204350000004002100039000000000003004b000003270000613d00000080040000390000000005000019000000200440003900000000060404330000025e0660019700000000026204360000000105500039000000000035004b000003200000413d0000000002120049000002390020009c00000239020080410000006002200210000002390010009c00000239010080410000004001100210000000000112019f0000000002000414000002390020009c0000023902008041000000c002200210000000000121019f00000263011001c70000800d0200003900000001030000390000026404000041000002c70000013d000000800010043f000000a00020043f0000014000000443000001600010044300000020010000390000018000100443000001a0002004430000010000100443000000020100003900000120001004430000023f01000041000008df0001042e000000a0040000390000026a0200004100000000030000190000000005040019000000000402041a000000000445043600000001022000390000000103300039000000000013004b000003480000413d000000600250008a000000800100003908de08740000040f000000400100043d000a00000001001d000000800200003908de08980000040f0000000a020000290000000001210049000002390010009c00000239010080410000006001100210000002390020009c00000239020080410000004002200210000000000121019f000008df0001042e0000026f01000041000000800010043f0000026701000041000008e0000104300000026b01000041000000c00010043f0000027f01000041000008e0000104300000026b01000041000000800010043f0000026701000041000008e0000104300000028001000041000a00000003001d000000000013043500000002010003670000000403100370000000000303043b0000023d0030009c0000008d0000213d0000000a0b0000290000000404b0003900000000003404350000002403b000390000002401100370000000000101043b00000000001304350000000001000414000000040020008c000004470000c13d0000000103000031000000200030008c00000020040000390000000004034019000004720000013d0000025e01100197000000000010043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000101041a000000000001004b0000048c0000c13d0000028601000041000000000010043f0000000701000039000000200010043f0000028701000041000000000101041a000000000001004b0000000002000039000000010200c0390000013c0000013d0000000101000039000000010110018f000000800010043f0000026901000041000008df0001042e0000014001100039000000400010043f000001400050043f000000000060043f000000000005004b000003bd0000613d00000160060000390000028a0700004100000000080000190000028b0010009c000002a90000213d0000004009100039000000400090043f000000000907041a000000a00a9002700000023c0aa00197000000200b1000390000000000ab04350000023d09900197000000000091043500000000061604360000000107700039000000400100043d0000000108800039000000000058004b000003ac0000413d0000023903300197000001000040043f0000000404000039000000000404041a0000023c04400197000001200040043f00000020041000390000004005000039000000000054043500000000003104350000004003100039000000e00400043d00000000004304350000006004100039000001000300043d0000000000240435000000a00210003900000000040304330000000000420435000000c002100039000000000004004b000003e00000613d00000000050000190000002003300039000000000603043300000000760604340000023d06600197000000000662043600000000070704330000023c07700197000000000076043500000040022000390000000105500039000000000045004b000003d40000413d000001200300043d0000023c03300197000000800410003900000000003404350000000002120049000002390020009c00000239020080410000006002200210000002390010009c00000239010080410000004001100210000000000112019f000008df0001042e000000000100041a0000023e01100197000000000161019f000000000010041b0000000001000414000002390010009c0000023901008041000000c00110021000000263011001c70000800d0200003900000003030000390000026504000041000002c70000013d000400000002001d000900000006001d000800200060003d0000000802100360000000000202043b000000230850008a0000027705200197000302770080019b000000030650014f000000030050006c00000000050000190000027705004041000200000008001d000000000082004b00000000070000190000027707008041000002770060009c000000000507c019000000000005004b0000008d0000c13d0000000906200029000000000561034f000000000505043b0000023c0050009c0000008d0000213d000000060750021000000000077300490000002008600039000000000078004b0000000009000019000002770900204100000277077001970000027708800197000000000a78013f000000000078004b000000000700001900000277070040410000027700a0009c000000000709c019000000000007004b0000008d0000c13d000000020050008c000004350000413d00000006074002100000000008670019000000000781034f000000000707043b0000023c0070009c0000008d0000213d0000004008800039000000000881034f000000000808043b0000023c0080009c0000008d0000213d000000000087004b0000055f0000813d0000000104400039000000000054004b000004250000413d00000008040000290000002004400039000000000441034f000000000404043b0000023c0040009c0000008d0000213d00000001064002100000028d046001970000028e06600197000000000064004b000005590000c13d00000001044001bf000000000045004b000006d10000813d0000029701000041000000800010043f0000026701000041000008e0000104300000023900b0009c000002390300004100000000030b40190000004003300210000002390010009c0000023901008041000000c001100210000000000131019f00000281011001c708de08d90000040f0000000a0b000029000000000301001900000060033002700000023903300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b0019000004620000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b0000045e0000c13d000000000006004b0000046f0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000004910000613d0000001f01400039000000600210018f0000000001b20019000000000021004b000000000200003900000001020040390000023c0010009c000002a90000213d0000000100200190000002a90000c13d000000400010043f000000200030008c0000008d0000413d0000000a020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b0000008d0000c13d0000000000210435000002390010009c0000023901008041000000400110021000000282011001c7000008df0001042e00000001020000390000013c0000013d000000400100043d0000026b02000041000001650000013d0000001f0530018f0000023b06300198000000400200043d00000000046200190000049c0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000004980000c13d000000000005004b000004a90000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000002390020009c00000239020080410000004002200210000000000112019f000008e0000104300000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000025e01100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000101041a000000000001004b000004e60000c13d0000000603000039000000000103041a0000023c0010009c000002a90000213d0000000102100039000000000023041b000002620110009a0000000902000029000000000021041b000000000103041a000800000001001d000000000020043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b0000000802000029000000000021041b0000000a020000290000000102200039000000800100043d000000000012004b000004b00000413d0000004c0000013d00000080010000390000000a0200002908de08a80000040f0000000001010433000000400200043d0000026003000041000005450000013d0000000002000019000a00000002001d0000000501200210000000a00110003900000000010104330000025e01100197000900000001001d000000000010043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000301041a000000000003004b0000053f0000613d0000000601000039000000000201041a000000000002004b000005590000613d000000010130008a000000000023004b000005260000613d000000000012004b0000078d0000a13d0000026c0130009a0000026c0220009a000000000202041a000000000021041b000000000020043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c70000801002000039000800000003001d08de08d90000040f00000001002001900000008d0000613d000000000101043b0000000802000029000000000021041b0000000601000039000000000301041a000000000003004b0000054e0000613d000000010130008a0000026c0230009a000000000002041b0000000602000039000000000012041b0000000901000029000000000010043f0000000701000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000001041b0000000a020000290000000102200039000000800100043d000000000012004b000004ee0000413d0000026f0000013d00000080010000390000000a0200002908de08a80000040f0000000001010433000000400200043d0000026d0300004100000000003204350000025e0110019700000004032000390000000000130435000002390020009c0000023902008041000000400120021000000261011001c7000008e0000104300000027e01000041000000000010043f0000003101000039000000040010043f0000026101000041000008e0000104300000000402000039000000000202041a0000023c022001970000023c0020009c000005630000c13d0000027e01000041000000000010043f0000001101000039000000040010043f0000026101000041000008e0000104300000028c01000041000000800010043f0000026701000041000008e000010430000000000021004b000007930000a13d0000000201000039000000000101041a000900000001001d0000014001000039000000400010043f000002720100004100000000001004430000000001000414000002390010009c0000023901008041000000c00110021000000273011001c70000800b0200003908de08d90000040f0000000100200190000006d00000613d000000000101043b000000800010043f000002680100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000002390010009c0000023901008041000000c00110021000000274011001c7000080050200003908de08d90000040f0000000100200190000006d00000613d000000000101043b0000023c01100197000000a00010043f0000000001000410000000c00010043f0000000a01000029000000e00010043f0000000901000029000001000010043f00000002020003670000002401200370000000000101043b000900000001001d0000000401100039000000000112034f000000000101043b0000023c0010009c000002a90000213d00000005031002100000003f043000390000025b04400197000000400500043d0000000004450019000400000005001d000000000054004b000000000500003900000001050040390000023c0040009c000002a90000213d0000000100500190000002a90000c13d000000400040043f0000000404000029000000000014043500000009010000290000002405100039000800000053001d000000080050006c000006250000813d0000000001000031000a00000001001d00070024001000920000000409000029000000000152034f000000000101043b0000023c0010009c0000008d0000213d00000009041000290000000701400069000002750010009c0000008d0000213d000000a00010008c0000008d0000413d000000400a00043d0000027600a0009c000002a90000213d000000a001a00039000000400010043f0000002401400039000000000312034f000000000303043b0000023c0030009c0000008d0000213d000000000c3a0436000000200b1000390000000001b2034f000000000101043b0000023c0010009c0000008d0000213d000000000441001900000043014000390000000a06000029000000000061004b0000000003000019000002770300804100000277011001970000027706600197000000000861013f000000000061004b00000000010000190000027701004041000002770080009c000000000103c019000000000001004b0000008d0000c13d0000002403400039000000000132034f000000000d01043b0000023c00d0009c000002a90000213d0000001f01d000390000029a011001970000003f011000390000029a01100197000000400e00043d00000000061e00190000000000e6004b000000000100003900000001010040390000023c0060009c000002a90000213d0000000100100190000002a90000c13d000000400060043f0000000001de04360000000004d4001900000044044000390000000a0040006c0000008d0000213d0000002003300039000000000832034f0000029a06d001980000000004610019000005fd0000613d000000000308034f000000000f010019000000003703043c000000000f7f043600000000004f004b000005f90000c13d0000001f03d001900000060a0000613d000000000668034f0000000303300210000000000704043300000000073701cf000000000737022f000000000606043b0000010003300089000000000636022f00000000033601cf000000000373019f00000000003404350000000001d1001900000000000104350000000000ec04350000002001b00039000000000312034f000000000303043b0000023c0030009c0000008d0000213d0000004004a0003900000000003404350000002001100039000000000312034f000000000303043b0000023c0030009c0000008d0000213d00000020099000390000006004a0003900000000003404350000002001100039000000000112034f000000000101043b0000008003a0003900000000001304350000000000a904350000002005500039000000080050006c000005b00000413d0000000401000029000001200010043f000000400200043d0000004001200039000000400300003900000000003104350000002003200039000002780100004100000000001304350000006001200039000000800400043d0000000000410435000000a00100043d0000023c0110019700000080042000390000000000140435000000c00100043d0000023d01100197000000a0042000390000000000140435000000e00100043d0000023d01100197000000c0042000390000000000140435000000e001200039000001000400043d00000000004104350000010001200039000000c005000039000001200400043d0000000000510435000001200120003900000000050404330000000000510435000001400620003900000005015002100000000007610019000000000005004b0000082d0000c13d0000000001270049000000200410008a00000000004204350000001f011000390000029a041001970000000001240019000000000041004b000000000400003900000001040040390000023c0010009c000002a90000213d0000000100400190000002a90000c13d000000400010043f000002390030009c000002390300804100000040013002100000000002020433000002390020009c00000239020080410000006002200210000000000112019f0000000002000414000002390020009c0000023902008041000000c002200210000000000112019f00000263011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d00000002020003670000000603200360000000000101043b000800000001001d000000000103043b000000000001004b000002ca0000613d000900000000001d000a00000000001d0000000a01000029000000060110021000000005011000290000002003100039000000000332034f000000000112034f000000000101043b000000000203043b000000400300043d000000600430003900000000002404350000004002300039000000000012043500000020013000390000001b02000039000000000021043500000008010000290000000000130435000000000000043f000002390030009c000002390300804100000040013002100000000002000414000002390020009c0000023902008041000000c002200210000000000112019f00000279011001c7000000010200003908de08d90000040f000000000301001900000060033002700000023903300197000000200030008c000000200500003900000000050340190000002004500190000006a10000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b0000069d0000c13d0000001f05500190000006ae0000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f0000000000540435000100000003001f00000001002001900000085f0000613d000000000100043d0000023d021001980000086b0000613d000000090020006b0000086e0000813d000900000002001d000000000020043f0000000801000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000101041a000000ff00100190000008710000613d00000002020003670000000601200360000000000101043b0000000a03000029000a00010030003d0000000a0010006b000900090000002d000006750000413d000002ca0000013d000000000001042f0000000304000039000000000404041a000000000004004b000006f60000613d0000000601000029000700240010003d000a0001004000920000000301000039000000000101041a0000000a0010006b0000078d0000813d0000028f0140009a000000000101041a0000023d01100197000000000010043f0000000801000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f00000001002001900000008d0000613d000000000101043b000000000201041a0000029b02200197000000000021041b0000000a04000029000000000004004b000006d70000c13d000000000300003100000002010003670000000702100360000000000202043b000000230400008a00050006004000720000000504300029000000000042004b0000000005000019000002770500804100000277044001970000027706200197000000000746013f000000000046004b00000000040000190000027704004041000002770070009c000000000405c019000000000004004b0000008d0000c13d0000000604000029000100440040003d000000000a0000190000000905200029000000000451034f000000000404043b0000023c0040009c0000008d0000213d00000006064002100000000006630049000000200550003900000277076001970000027708500197000000000978013f000000000078004b00000000070000190000027707004041000000000065004b00000000060000190000027706002041000002770090009c000000000706c019000000000007004b0000008d0000c13d00000000004a004b000007970000813d0007000600a002180000000702500029000000000121034f000000000101043b0000023d0010009c0000008d0000213d000000000010043f0000000801000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c70000801002000039000a0000000a001d08de08d90000040f0000000a0800002900000001002001900000008d0000613d000000000101043b000000000101041a000000ff001001900000082a0000c13d00000002010003670000000802100360000000000302043b0000000002000031000000060420006a000000230440008a00000277054001970000027706300197000000000756013f000000000056004b00000000050000190000027705004041000000000043004b00000000040000190000027704008041000002770070009c000000000504c019000000000005004b0000008d0000c13d0000000904300029000000000341034f000000000303043b0000023c0030009c0000008d0000213d00000006053002100000000005520049000000200240003900000277045001970000027706200197000000000746013f000000000046004b00000000040000190000027704004041000000000052004b00000000050000190000027705002041000002770070009c000000000405c019000000000004004b0000008d0000c13d000000000038004b0000078d0000813d0000000702200029000000000121034f000000000101043b0000023d0010009c0000008d0000213d000000000010043f0000000801000039000000200010043f0000000001000414000002390010009c0000023901008041000000c0011002100000025f011001c7000080100200003908de08d90000040f0000000a0a00002900000001002001900000008d0000613d000000000101043b000000000201041a0000029b0220019700000001022001bf000000000021041b00000002010003670000000802100360000000000202043b0000000003000031000000050430002900000277054001970000027706200197000000000756013f000000000056004b00000000050000190000027705002041000000000042004b00000000040000190000027704004041000002770070009c000000000504c019000000010aa00039000000000005004b000007090000c13d0000008d0000013d0000027e01000041000000000010043f0000003201000039000000040010043f0000026101000041000008e0000104300000027101000041000000800010043f0000026701000041000008e00001043000000002060000390000000407000029000000000076041b0000000306000039000000000706041a000000000046041b000000000074004b000007a70000813d000002900640009a000002900770009a000000000076004b000007a70000813d000000000006041b0000000106600039000000000076004b000007a30000413d0000000306000039000000000060043f000000000004004b000007c20000613d0000028a060000410000000007000019000000000851034f000000000808043b0000023d0080009c0000008d0000213d0000002009500039000000000991034f000000000909043b0000023c0090009c0000008d0000213d000000000a06041a000002910aa00197000000a00990021000000292099001970000000009a9019f000000000889019f000000000086041b000000010660003900000040055000390000000107700039000000000047004b000007ad0000413d0000000104100360000000000404043b0000023c0040009c0000008d0000213d0000000405000039000000000605041a0000029306600197000000000646019f000000000065041b0000000506000039000000000506041a0000023907500197000002390070009c000005590000613d000002940750019700000001055000390000023905500197000000000775019f000000000076041b000000400600043d000000200760003900000004080000290000000000870435000000020020006c000000000700001900000277070080410000027708200197000000030980014f000000030080006c00000000080000190000027708004041000002770090009c000000000807c01900000020070000390000000000760435000000000008004b0000008d0000c13d00000006072000290000000402700039000000000221034f000000000202043b0000023c0020009c0000008d0000213d000000240770003900000006082002100000000003830049000000000037004b0000000008000019000002770800204100000277033001970000027709700197000000000a39013f000000000039004b000000000300001900000277030040410000027700a0009c000000000308c019000000000003004b0000008d0000c13d00000080036000390000004008600039000000600900003900000000009804350000000000230435000000a003600039000000000002004b000008160000613d0000000008000019000000000971034f000000000909043b0000023d0090009c0000008d0000213d0000000009930436000000200a700039000000000aa1034f000000000a0a043b0000023c00a0009c0000008d0000213d0000000000a90435000000400770003900000040033000390000000108800039000000000028004b000008060000413d000000600160003900000000004104350000000001630049000002390010009c00000239010080410000006001100210000002390060009c00000239060080410000004002600210000000000121019f0000000002000414000002390020009c0000023902008041000000c002200210000000000121019f00000263011001c70000800d0200003900000002030000390000029504000041000003380000013d000000400100043d0000029602000041000001650000013d000000a0080000390000000009000019000008460000013d0000000001bc001900000000000104350000004001a0003900000000010104330000023c01100197000000400d70003900000000001d04350000006001a0003900000000010104330000023c01100197000000600d70003900000000001d043500000080017000390000008007a00039000000000707043300000000007104350000001f01b000390000029a0110019700000000071c00190000000109900039000000000059004b0000064c0000813d0000000001270049000001400110008a00000000061604360000002004400039000000000a04043300000000b10a04340000023c011001970000000001170436000000000b0b04330000000000810435000000a00170003900000000db0b04340000000000b10435000000c00c70003900000000000b004b000008300000613d000000000e0000190000000001ce0019000000000fed0019000000000f0f04330000000000f10435000000200ee000390000000000be004b000008570000413d000008300000013d0000001f0530018f0000023b06300198000000400200043d00000000046200190000049c0000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000008660000c13d0000049c0000013d000000400100043d0000027c02000041000001650000013d000000400100043d0000027a02000041000001650000013d000000400100043d0000027b02000041000001650000013d0000001f022000390000029a022001970000000001120019000000000021004b000000000200003900000001020040390000023c0010009c000008800000213d0000000100200190000008800000c13d000000400010043f000000000001042d0000027e01000041000000000010043f0000004101000039000000040010043f0000026101000041000008e00001043000000000430104340000000001320436000000000003004b000008920000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b0000088b0000413d000000000231001900000000000204350000001f023000390000029a022001970000000001210019000000000001042d00000020030000390000000004310436000000000302043300000000003404350000004001100039000000000003004b000008a70000613d0000000004000019000000200220003900000000050204330000025e0550019700000000015104360000000104400039000000000034004b000008a00000413d000000000001042d0000000003010433000000000023004b000008af0000a13d000000050220021000000000012100190000002001100039000000000001042d0000027e01000041000000000010043f0000003201000039000000040010043f0000026101000041000008e000010430000000000001042f00000000050100190000000000200443000000050030008c000008c40000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b000008bc0000413d000002390030009c000002390300804100000060013002100000000002000414000002390020009c0000023902008041000000c002200210000000000112019f0000029c011001c7000000000205001908de08d90000040f0000000100200190000008d30000613d000000000101043b000000000001042d000000000001042f000008d7002104210000000102000039000000000001042d0000000002000019000000000001042d000008dc002104230000000102000039000000000001042d0000000002000019000000000001042d000008de00000432000008df0001042e000008e00001043000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000000000000000000000000000ffffffffffffffff000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000002000000000000000000000000000000c00000010000000000000000009cf8540c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000009b15e16f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006d2d3992000000000000000000000000000000000000000000000000000000009a19b32800000000000000000000000000000000000000000000000000000000eaa83ddc00000000000000000000000000000000000000000000000000000000eaa83ddd00000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000f8bb876e000000000000000000000000000000000000000000000000000000009a19b32900000000000000000000000000000000000000000000000000000000d881e0920000000000000000000000000000000000000000000000000000000079ba50960000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000008da5cb5b000000000000000000000000000000000000000000000000000000006d2d39930000000000000000000000000000000000000000000000000000000070a9089e00000000000000000000000000000000000000000000000000000000397796f60000000000000000000000000000000000000000000000000000000062eed4140000000000000000000000000000000000000000000000000000000062eed415000000000000000000000000000000000000000000000000000000006509a95400000000000000000000000000000000000000000000000000000000397796f7000000000000000000000000000000000000000000000000000000004d616771000000000000000000000000000000000000000000000000000000001add205e000000000000000000000000000000000000000000000000000000001add205f000000000000000000000000000000000000000000000000000000002cbc26bb00000000000000000000000000000000000000000000000000000000181f5a7700000000000000000000000000000000000000000000000000000000198f0f777fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffff7f00000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000020000000000000000000000000000000000004000000000000000000000000019d5c79b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002400000000000000000000000009addddcec1d7ba6ad726df49aeea3e93fb0c1037d551236841a60c0c883f2c102000000000000000000000000000000000000000000000000000000000000001716e663a90a76d3b6c7e5f680673d1b051454c19c627e184c8daf28f3104f74ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e0000000000000000000000000000000000000020000000800000000000000000f652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f2b5c74de0000000000000000000000000000000000000000000000000000000009addddcec1d7ba6ad726df49aeea3e93fb0c1037d551236841a60c0c883f2c273281fa1000000000000000000000000000000000000000000000000000000000676e709c9cc74fa0519fd78f7c33be0f1b2b0bae0507c724aef7229379c6ba102b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e059fa4a93000000000000000000000000000000000000000000000000000000009a8a0592ac89c5ad3bc6df8224c17b485976f597df104ee20d0df415241f670b020000020000000000000000000000000000000400000000000000000000000002000002000000000000000000000000000000440000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff5f80000000000000000000000000000000000000000000000000000000000000009651943783dbf81935a60e98f218a9d9b5b28823fb2228bbd91320d632facf530000000000000000000000000000000000000080000000000000000000000000bbe15e7f00000000000000000000000000000000000000000000000000000000aaaa9141000000000000000000000000000000000000000000000000000000008baa579f00000000000000000000000000000000000000000000000000000000ace124bc000000000000000000000000000000000000000000000000000000004e487b71000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000c000000000000000004d61677100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004400000000000000000000000000000000000000000000000000000000000000200000000000000000000000000a7c4edd0000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000007dde556524061d0ce70b736a6e842a48e4927608bf87fd31432ced12a03ffeb8010000000000000000000000000000010000000000000000000000000000000070b766b11586b6b505ed3893938b0cc6c6c98bd6f65e969ac311168d34e4f9e20000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000fffffffffffffebfc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b000000000000000000000000000000000000000000000000ffffffffffffffbf4485151700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffffe000000000000000000000000000000000000000000000001fffffffffffffffe3da8a5f161a6c3ff06a60736d0ed24d7963cc6a5c4fafd2fa1dae9bb908e07a63da8a5f161a6c3ff06a60736d0ed24d7963cc6a5c4fafd2fa1dae9bb908e07a5ffffffff0000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000007f22bf988149dbe8de8fb879c6b97a4e56e68b2bd57421ce1a4e79d4ef6b496c28cae27d00000000000000000000000000000000000000000000000000000000014c502000000000000000000000000000000000000000000000000000000000524d4e52656d6f746520312e362e302d646576000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000200000200000000000000000000000000000000000000000000000000000000")
