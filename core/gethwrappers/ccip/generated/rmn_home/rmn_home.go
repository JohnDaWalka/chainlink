package rmn_home

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

type RMNHomeDynamicConfig struct {
	SourceChains   []RMNHomeSourceChain
	OffchainConfig []byte
}

type RMNHomeNode struct {
	PeerId            [32]byte
	OffchainPublicKey [32]byte
}

type RMNHomeSourceChain struct {
	ChainSelector       uint64
	F                   uint64
	ObserverNodesBitmap *big.Int
}

type RMNHomeStaticConfig struct {
	Nodes          []RMNHomeNode
	OffchainConfig []byte
}

type RMNHomeVersionedConfig struct {
	Version       uint32
	ConfigDigest  [32]byte
	StaticConfig  RMNHomeStaticConfig
	DynamicConfig RMNHomeDynamicConfig
}

var RMNHomeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expectedConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"gotConfigDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"DigestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateOffchainPublicKey\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicatePeerId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSourceChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoOpStateTransitionNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughObservers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsNodesLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfBoundsObserverNodeIndex\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RevokingZeroDigestNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ActiveConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"CandidateConfigRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"ConfigPromoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllConfigs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"activeConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"candidateConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCandidateDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"internalType\":\"structRMNHome.VersionedConfig\",\"name\":\"versionedConfig\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"ok\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDigests\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"activeConfigDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"candidateConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"digestToPromote\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digestToRevoke\",\"type\":\"bytes32\"}],\"name\":\"promoteCandidateAndRevokeActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"revokeCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"peerId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"offchainPublicKey\",\"type\":\"bytes32\"}],\"internalType\":\"structRMNHome.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"digestToOverwrite\",\"type\":\"bytes32\"}],\"name\":\"setCandidate\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"newConfigDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"f\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"observerNodesBitmap\",\"type\":\"uint256\"}],\"internalType\":\"structRMNHome.SourceChain[]\",\"name\":\"sourceChains\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structRMNHome.DynamicConfig\",\"name\":\"newDynamicConfig\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"currentDigest\",\"type\":\"bytes32\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e80546001600160401b03191690553480156200002157600080fd5b50336000816200004457604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b0384811691909117909155811615620000775762000077816200007f565b5050620000f9565b336001600160a01b03821603620000a957604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6128cb80620001096000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636dd5b69d1161008c5780638c76967f116100665780638c76967f146101d45780638da5cb5b146101e7578063f2fde38b1461020f578063fb4022d41461022257600080fd5b80636dd5b69d14610196578063736be802146101b757806379ba5097146101cc57600080fd5b80633567e6b4116100bd5780633567e6b41461015b57806338354c5c14610178578063635079561461018057600080fd5b8063118dbac5146100e4578063123e65db1461010a578063181f5a7714610112575b600080fd5b6100f76100f236600461186a565b610235565b6040519081526020015b60405180910390f35b6100f7610418565b61014e6040518060400160405280601181526020017f524d4e486f6d6520312e362e302d64657600000000000000000000000000000081525081565b6040516101019190611945565b610163610457565b60408051928352602083019190915201610101565b6100f76104d8565b6101886104f7565b604051610101929190611ab0565b6101a96101a4366004611ad5565b610a79565b604051610101929190611aee565b6101ca6101c5366004611b12565b610d5d565b005b6101ca610e79565b6101ca6101e2366004611b57565b610f47565b60015460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610101565b6101ca61021d366004611b79565b61115a565b6101ca610230366004611ad5565b61116e565b600061023f61128a565b61025961024b85611d2b565b61025485611e2c565b6112dd565b60006102636104d8565b90508281146102ad576040517f93df584c00000000000000000000000000000000000000000000000000000000815260048101829052602481018490526044015b60405180910390fd5b80156102df5760405183907f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b90600090a25b600e80546000919082906102f89063ffffffff16611f41565b91906101000a81548163ffffffff021916908363ffffffff160217905590506103408660405160200161032b91906120ec565b60405160208183030381529060405282611455565b600e54909350600090600290640100000000900463ffffffff1660011863ffffffff1660028110610373576103736120ff565b600602016001810185905580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8416178155905086600282016103bd828261234c565b50869050600482016103cf828261254b565b905050837ff6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf08389896040516104069392919061278f565b60405180910390a25050509392505050565b60006002610434600e5463ffffffff6401000000009091041690565b63ffffffff166002811061044a5761044a6120ff565b6006020160010154905090565b6000806002610474600e5463ffffffff6401000000009091041690565b63ffffffff166002811061048a5761048a6120ff565b600602016001015460026104b2600e54600163ffffffff640100000000909204919091161890565b63ffffffff16600281106104c8576104c86120ff565b6006020160010154915091509091565b600e54600090600290640100000000900463ffffffff16600118610434565b6104ff6117ec565b6105076117ec565b60006002610523600e5463ffffffff6401000000009091041690565b63ffffffff1660028110610539576105396120ff565b6040805160808101825260069290920292909201805463ffffffff16825260018101546020808401919091528351600283018054606093810283018401875282870181815295969495948701949293919284929091849160009085015b828210156105dc57838290600052602060002090600202016040518060400160405290816000820154815260200160018201548152505081526020019060010190610596565b5050505081526020016001820180546105f490612193565b80601f016020809104026020016040519081016040528092919081815260200182805461062090612193565b801561066d5780601f106106425761010080835404028352916020019161066d565b820191906000526020600020905b81548152906001019060200180831161065057829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b8282101561070f5760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff80821685526801000000000000000090910416838501526001908101549183019190915290835290920191016106af565b50505050815260200160018201805461072790612193565b80601f016020809104026020016040519081016040528092919081815260200182805461075390612193565b80156107a05780601f10610775576101008083540402835291602001916107a0565b820191906000526020600020905b81548152906001019060200180831161078357829003601f168201915b505050919092525050509052506020810151909150156107be578092505b600e54600090600290640100000000900463ffffffff1660011863ffffffff16600281106107ee576107ee6120ff565b6040805160808101825260069290920292909201805463ffffffff16825260018101546020808401919091528351600283018054606093810283018401875282870181815295969495948701949293919284929091849160009085015b828210156108915783829060005260206000209060020201604051806040016040529081600082015481526020016001820154815250508152602001906001019061084b565b5050505081526020016001820180546108a990612193565b80601f01602080910402602001604051908101604052809291908181526020018280546108d590612193565b80156109225780601f106108f757610100808354040283529160200191610922565b820191906000526020600020905b81548152906001019060200180831161090557829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b828210156109c45760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff8082168552680100000000000000009091041683850152600190810154918301919091529083529092019101610964565b5050505081526020016001820180546109dc90612193565b80601f0160208091040260200160405190810160405280929190818152602001828054610a0890612193565b8015610a555780601f10610a2a57610100808354040283529160200191610a55565b820191906000526020600020905b815481529060010190602001808311610a3857829003601f168201915b50505091909252505050905250602081015190915015610a73578092505b50509091565b610a816117ec565b6000805b6002811015610d52578360028260028110610aa257610aa26120ff565b6006020160010154148015610ab657508315155b15610d4a5760028160028110610ace57610ace6120ff565b6040805160808101825260069290920292909201805463ffffffff16825260018082015460208085019190915284516002840180546060938102830184018852828801818152959794969588958701948492849160009085015b82821015610b6e57838290600052602060002090600202016040518060400160405290816000820154815260200160018201548152505081526020019060010190610b28565b505050508152602001600182018054610b8690612193565b80601f0160208091040260200160405190810160405280929190818152602001828054610bb290612193565b8015610bff5780601f10610bd457610100808354040283529160200191610bff565b820191906000526020600020905b815481529060010190602001808311610be257829003601f168201915b50505050508152505081526020016004820160405180604001604052908160008201805480602002602001604051908101604052809291908181526020016000905b82821015610ca15760008481526020908190206040805160608101825260028602909201805467ffffffffffffffff8082168552680100000000000000009091041683850152600190810154918301919091529083529092019101610c41565b505050508152602001600182018054610cb990612193565b80601f0160208091040260200160405190810160405280929190818152602001828054610ce590612193565b8015610d325780601f10610d0757610100808354040283529160200191610d32565b820191906000526020600020905b815481529060010190602001808311610d1557829003601f168201915b50505091909252505050905250969095509350505050565b600101610a85565b509092600092509050565b610d6561128a565b60005b6002811015610e3f578160028260028110610d8557610d856120ff565b6006020160010154148015610d9957508115155b15610e3757610dd0610daa84611e2c565b60028360028110610dbd57610dbd6120ff565b600602016002016000018054905061155d565b8260028260028110610de457610de46120ff565b600602016004018181610df7919061254b565b905050817f1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e629484604051610e2a91906127ca565b60405180910390a2505050565b600101610d68565b506040517fd0b2c031000000000000000000000000000000000000000000000000000000008152600481018290526024016102a4565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610eca576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610f4f61128a565b81158015610f5b575080155b15610f92576040517f7b4d1e4f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54600163ffffffff6401000000009092048216181682600282818110610fbc57610fbc6120ff565b6006020160010154146110225760028160028110610fdc57610fdc6120ff565b6006020160010154836040517f93df584c0000000000000000000000000000000000000000000000000000000081526004016102a4929190918252602082015260400190565b6000600261103e600e5463ffffffff6401000000009091041690565b63ffffffff1660028110611054576110546120ff565b600602019050828160010154146110a75760018101546040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018490526044016102a4565b6000600180830191909155600e805463ffffffff6401000000008083048216909418169092027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff90921691909117905582156111295760405183907f0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b890600090a25b60405184907ffc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e90600090a250505050565b61116261128a565b61116b81611728565b50565b61117661128a565b806111ad576040517f0849d8cc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600e54600163ffffffff64010000000090920482161816816002828181106111d7576111d76120ff565b60060201600101541461123d57600281600281106111f7576111f76120ff565b6006020160010154826040517f93df584c0000000000000000000000000000000000000000000000000000000081526004016102a4929190918252602082015260400190565b60405182907f53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b90600090a26002816002811061127b5761127b6120ff565b60060201600101600090555050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146112db576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b815151610100101561131b576040517faf26d5e300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251518110156114455760006113358260016127dd565b90505b83515181101561143c578351805182908110611356576113566120ff565b60200260200101516000015184600001518381518110611378576113786120ff565b602002602001015160000151036113bb576040517f221a8ae800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83518051829081106113cf576113cf6120ff565b602002602001015160200151846000015183815181106113f1576113f16120ff565b60200260200101516020015103611434576040517fae00651d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600101611338565b5060010161131e565b50610e758183600001515161155d565b604080517f45564d00000000000000000000000000000000000000000000000000000000006020820152469181019190915230606082015263ffffffff821660808201526000907dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060a001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526114fc9186906020016127f0565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190528051602090910120167e0b0000000000000000000000000000000000000000000000000000000000001790505b92915050565b81515160005b8181101561172257600084600001518281518110611583576115836120ff565b60200260200101519050600082600161159c91906127dd565b90505b8381101561161f5785518051829081106115bb576115bb6120ff565b60200260200101516000015167ffffffffffffffff16826000015167ffffffffffffffff1603611617576040517f3857f84d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60010161159f565b506040810151806116328661010061281f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff901c82161461168e576040517f2847b60600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81156116b6576116a260018361281f565b909116906116af81612832565b9050611691565b60208301516116c690600261286a565b6116d1906001612896565b67ffffffffffffffff16811015611714576040517fa804bcb300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050806001019050611563565b50505050565b3373ffffffffffffffffffffffffffffffffffffffff821603611777576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040518060800160405280600063ffffffff1681526020016000801916815260200161182b604051806040016040528060608152602001606081525090565b815260200161184d604051806040016040528060608152602001606081525090565b905290565b60006040828403121561186457600080fd5b50919050565b60008060006060848603121561187f57600080fd5b833567ffffffffffffffff8082111561189757600080fd5b6118a387838801611852565b945060208601359150808211156118b957600080fd5b506118c686828701611852565b925050604084013590509250925092565b60005b838110156118f25781810151838201526020016118da565b50506000910152565b600081518084526119138160208601602086016118d7565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061195860208301846118fb565b9392505050565b8051604080845281518482018190526000926060916020918201918388019190865b828110156119bb578451805167ffffffffffffffff9081168652838201511683860152870151878501529381019392850192600101611981565b50808801519550888303818a015250506119d581856118fb565b979650505050505050565b63ffffffff81511682526000602080830151818501526040808401516080604087015260c0860181516040608089015281815180845260e08a0191508683019350600092505b80831015611a4f5783518051835287015187830152928601926001929092019190850190611a26565b50948301518886037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff800160a08a015294611a8981876118fb565b9550505050505060608301518482036060860152611aa7828261195f565b95945050505050565b604081526000611ac360408301856119e0565b8281036020840152611aa781856119e0565b600060208284031215611ae757600080fd5b5035919050565b604081526000611b0160408301856119e0565b905082151560208301529392505050565b60008060408385031215611b2557600080fd5b823567ffffffffffffffff811115611b3c57600080fd5b611b4885828601611852565b95602094909401359450505050565b60008060408385031215611b6a57600080fd5b50508035926020909101359150565b600060208284031215611b8b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461195857600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611c0157611c01611baf565b60405290565b6040516060810167ffffffffffffffff81118282101715611c0157611c01611baf565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611c7157611c71611baf565b604052919050565b600067ffffffffffffffff821115611c9357611c93611baf565b5060051b60200190565b600082601f830112611cae57600080fd5b813567ffffffffffffffff811115611cc857611cc8611baf565b611cf960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611c2a565b818152846020838601011115611d0e57600080fd5b816020850160208301376000918101602001919091529392505050565b60006040808336031215611d3e57600080fd5b611d46611bde565b833567ffffffffffffffff80821115611d5e57600080fd5b9085019036601f830112611d7157600080fd5b81356020611d86611d8183611c79565b611c2a565b82815260069290921b84018101918181019036841115611da557600080fd5b948201945b83861015611de557878636031215611dc25760008081fd5b611dca611bde565b86358152838701358482015282529487019490820190611daa565b86525087810135955082861115611dfb57600080fd5b611e0736878a01611c9d565b90850152509195945050505050565b67ffffffffffffffff8116811461116b57600080fd5b60006040808336031215611e3f57600080fd5b611e47611bde565b833567ffffffffffffffff80821115611e5f57600080fd5b9085019036601f830112611e7257600080fd5b81356020611e82611d8183611c79565b82815260609283028501820192828201919036851115611ea157600080fd5b958301955b84871015611efb57808736031215611ebe5760008081fd5b611ec6611c07565b8735611ed181611e16565b815287850135611ee081611e16565b81860152878a01358a82015283529586019591830191611ea6565b5086525087810135955082861115611dfb57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff808316818103611f5a57611f5a611f12565b6001019392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611f9957600080fd5b830160208101925035905067ffffffffffffffff811115611fb957600080fd5b803603821315611fc857600080fd5b9250929050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000604080840183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe185360301811261205157600080fd5b8401602081810191359067ffffffffffffffff82111561207057600080fd5b8160061b360383131561208257600080fd5b6040885292819052909160009190606088015b828410156120bb5784358152818501358282015293850193600193909301928501612095565b6120c86020890189611f64565b9650945088810360208a01526120df818787611fcf565b9998505050505050505050565b6020815260006119586020830184612018565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261216357600080fd5b83018035915067ffffffffffffffff82111561217e57600080fd5b602001915036819003821315611fc857600080fd5b600181811c908216806121a757607f821691505b602082108103611864577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b601f82111561222c576000816000526020600020601f850160051c810160208610156122095750805b601f850160051c820191505b8181101561222857828155600101612215565b5050505b505050565b67ffffffffffffffff83111561224957612249611baf565b61225d836122578354612193565b836121e0565b6000601f8411600181146122af57600085156122795750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355612345565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156122fe57868501358255602094850194600190920191016122de565b5086821015612339577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b81357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe183360301811261237e57600080fd5b8201803567ffffffffffffffff81111561239757600080fd5b6020820191508060061b36038213156123af57600080fd5b680100000000000000008111156123c8576123c8611baf565b8254818455808210156124555760017f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808316831461240957612409611f12565b808416841461241a5761241a611f12565b5060008560005260206000208360011b81018560011b820191505b80821015612450578282558284830155600282019150612435565b505050505b5060008381526020902060005b8281101561248e5783358255602084013560018301556040939093019260029190910190600101612462565b5050505061249f602083018361212e565b611722818360018601612231565b81356124b881611e16565b67ffffffffffffffff811690508154817fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000821617835560208401356124fc81611e16565b6fffffffffffffffff00000000000000008160401b16837fffffffffffffffffffffffffffffffff00000000000000000000000000000000841617178455505050604082013560018201555050565b81357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe183360301811261257d57600080fd5b8201803567ffffffffffffffff81111561259657600080fd5b602082019150606080820236038313156125af57600080fd5b680100000000000000008211156125c8576125c8611baf565b8354828555808310156126555760017f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808316831461260957612609611f12565b808516851461261a5761261a611f12565b5060008660005260206000208360011b81018660011b820191505b80821015612650578282558284830155600282019150612635565b505050505b5060008481526020902060005b838110156126875761267485836124ad565b9382019360029190910190600101612662565b505050505061249f602083018361212e565b6000604080840183357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18536030181126126d257600080fd5b8401602081810191359067ffffffffffffffff808311156126f257600080fd5b6060808402360385131561270557600080fd5b60408a529483905292936000939060608a015b8486101561275c57863561272b81611e16565b831681528684013561273c81611e16565b831681850152868801358882015295810195600195909501948101612718565b61276960208b018b611f64565b985096508a810360208c0152612780818989611fcf565b9b9a5050505050505050505050565b63ffffffff841681526060602082015260006127ae6060830185612018565b82810360408401526127c08185612699565b9695505050505050565b6020815260006119586020830184612699565b8082018082111561155757611557611f12565b600083516128028184602088016118d7565b8351908301906128168183602088016118d7565b01949350505050565b8181038181111561155757611557611f12565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361286357612863611f12565b5060010190565b67ffffffffffffffff81811683821602808216919082811461288e5761288e611f12565b505092915050565b67ffffffffffffffff8181168382160190808211156128b7576128b7611f12565b509291505056fea164736f6c6343000818000a",
}

var RMNHomeABI = RMNHomeMetaData.ABI

var RMNHomeBin = RMNHomeMetaData.Bin

func DeployRMNHome(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RMNHome, error) {
	parsed, err := RMNHomeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RMNHomeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RMNHome{address: address, abi: *parsed, RMNHomeCaller: RMNHomeCaller{contract: contract}, RMNHomeTransactor: RMNHomeTransactor{contract: contract}, RMNHomeFilterer: RMNHomeFilterer{contract: contract}}, nil
}

type RMNHome struct {
	address common.Address
	abi     abi.ABI
	RMNHomeCaller
	RMNHomeTransactor
	RMNHomeFilterer
}

type RMNHomeCaller struct {
	contract *bind.BoundContract
}

type RMNHomeTransactor struct {
	contract *bind.BoundContract
}

type RMNHomeFilterer struct {
	contract *bind.BoundContract
}

type RMNHomeSession struct {
	Contract     *RMNHome
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type RMNHomeCallerSession struct {
	Contract *RMNHomeCaller
	CallOpts bind.CallOpts
}

type RMNHomeTransactorSession struct {
	Contract     *RMNHomeTransactor
	TransactOpts bind.TransactOpts
}

type RMNHomeRaw struct {
	Contract *RMNHome
}

type RMNHomeCallerRaw struct {
	Contract *RMNHomeCaller
}

type RMNHomeTransactorRaw struct {
	Contract *RMNHomeTransactor
}

func NewRMNHome(address common.Address, backend bind.ContractBackend) (*RMNHome, error) {
	abi, err := abi.JSON(strings.NewReader(RMNHomeABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindRMNHome(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMNHome{address: address, abi: abi, RMNHomeCaller: RMNHomeCaller{contract: contract}, RMNHomeTransactor: RMNHomeTransactor{contract: contract}, RMNHomeFilterer: RMNHomeFilterer{contract: contract}}, nil
}

func NewRMNHomeCaller(address common.Address, caller bind.ContractCaller) (*RMNHomeCaller, error) {
	contract, err := bindRMNHome(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMNHomeCaller{contract: contract}, nil
}

func NewRMNHomeTransactor(address common.Address, transactor bind.ContractTransactor) (*RMNHomeTransactor, error) {
	contract, err := bindRMNHome(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMNHomeTransactor{contract: contract}, nil
}

func NewRMNHomeFilterer(address common.Address, filterer bind.ContractFilterer) (*RMNHomeFilterer, error) {
	contract, err := bindRMNHome(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMNHomeFilterer{contract: contract}, nil
}

func bindRMNHome(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RMNHomeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_RMNHome *RMNHomeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNHome.Contract.RMNHomeCaller.contract.Call(opts, result, method, params...)
}

func (_RMNHome *RMNHomeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.Contract.RMNHomeTransactor.contract.Transfer(opts)
}

func (_RMNHome *RMNHomeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNHome.Contract.RMNHomeTransactor.contract.Transact(opts, method, params...)
}

func (_RMNHome *RMNHomeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMNHome.Contract.contract.Call(opts, result, method, params...)
}

func (_RMNHome *RMNHomeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.Contract.contract.Transfer(opts)
}

func (_RMNHome *RMNHomeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMNHome.Contract.contract.Transact(opts, method, params...)
}

func (_RMNHome *RMNHomeCaller) GetActiveDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getActiveDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNHome *RMNHomeSession) GetActiveDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetActiveDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetActiveDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetActiveDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetAllConfigs(opts *bind.CallOpts) (GetAllConfigs,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getAllConfigs")

	outstruct := new(GetAllConfigs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfig = *abi.ConvertType(out[0], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)
	outstruct.CandidateConfig = *abi.ConvertType(out[1], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetAllConfigs() (GetAllConfigs,

	error) {
	return _RMNHome.Contract.GetAllConfigs(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetAllConfigs() (GetAllConfigs,

	error) {
	return _RMNHome.Contract.GetAllConfigs(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetCandidateDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getCandidateDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_RMNHome *RMNHomeSession) GetCandidateDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetCandidateDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetCandidateDigest() ([32]byte, error) {
	return _RMNHome.Contract.GetCandidateDigest(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) GetConfig(opts *bind.CallOpts, configDigest [32]byte) (GetConfig,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getConfig", configDigest)

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.VersionedConfig = *abi.ConvertType(out[0], new(RMNHomeVersionedConfig)).(*RMNHomeVersionedConfig)
	outstruct.Ok = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetConfig(configDigest [32]byte) (GetConfig,

	error) {
	return _RMNHome.Contract.GetConfig(&_RMNHome.CallOpts, configDigest)
}

func (_RMNHome *RMNHomeCallerSession) GetConfig(configDigest [32]byte) (GetConfig,

	error) {
	return _RMNHome.Contract.GetConfig(&_RMNHome.CallOpts, configDigest)
}

func (_RMNHome *RMNHomeCaller) GetConfigDigests(opts *bind.CallOpts) (GetConfigDigests,

	error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "getConfigDigests")

	outstruct := new(GetConfigDigests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ActiveConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.CandidateConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_RMNHome *RMNHomeSession) GetConfigDigests() (GetConfigDigests,

	error) {
	return _RMNHome.Contract.GetConfigDigests(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) GetConfigDigests() (GetConfigDigests,

	error) {
	return _RMNHome.Contract.GetConfigDigests(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_RMNHome *RMNHomeSession) Owner() (common.Address, error) {
	return _RMNHome.Contract.Owner(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) Owner() (common.Address, error) {
	return _RMNHome.Contract.Owner(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RMNHome.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_RMNHome *RMNHomeSession) TypeAndVersion() (string, error) {
	return _RMNHome.Contract.TypeAndVersion(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeCallerSession) TypeAndVersion() (string, error) {
	return _RMNHome.Contract.TypeAndVersion(&_RMNHome.CallOpts)
}

func (_RMNHome *RMNHomeTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "acceptOwnership")
}

func (_RMNHome *RMNHomeSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNHome.Contract.AcceptOwnership(&_RMNHome.TransactOpts)
}

func (_RMNHome *RMNHomeTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _RMNHome.Contract.AcceptOwnership(&_RMNHome.TransactOpts)
}

func (_RMNHome *RMNHomeTransactor) PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "promoteCandidateAndRevokeActive", digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeSession) PromoteCandidateAndRevokeActive(digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.PromoteCandidateAndRevokeActive(&_RMNHome.TransactOpts, digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeTransactorSession) PromoteCandidateAndRevokeActive(digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.PromoteCandidateAndRevokeActive(&_RMNHome.TransactOpts, digestToPromote, digestToRevoke)
}

func (_RMNHome *RMNHomeTransactor) RevokeCandidate(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "revokeCandidate", configDigest)
}

func (_RMNHome *RMNHomeSession) RevokeCandidate(configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.RevokeCandidate(&_RMNHome.TransactOpts, configDigest)
}

func (_RMNHome *RMNHomeTransactorSession) RevokeCandidate(configDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.RevokeCandidate(&_RMNHome.TransactOpts, configDigest)
}

func (_RMNHome *RMNHomeTransactor) SetCandidate(opts *bind.TransactOpts, staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "setCandidate", staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeSession) SetCandidate(staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetCandidate(&_RMNHome.TransactOpts, staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeTransactorSession) SetCandidate(staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetCandidate(&_RMNHome.TransactOpts, staticConfig, dynamicConfig, digestToOverwrite)
}

func (_RMNHome *RMNHomeTransactor) SetDynamicConfig(opts *bind.TransactOpts, newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "setDynamicConfig", newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeSession) SetDynamicConfig(newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetDynamicConfig(&_RMNHome.TransactOpts, newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeTransactorSession) SetDynamicConfig(newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error) {
	return _RMNHome.Contract.SetDynamicConfig(&_RMNHome.TransactOpts, newDynamicConfig, currentDigest)
}

func (_RMNHome *RMNHomeTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _RMNHome.contract.Transact(opts, "transferOwnership", to)
}

func (_RMNHome *RMNHomeSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNHome.Contract.TransferOwnership(&_RMNHome.TransactOpts, to)
}

func (_RMNHome *RMNHomeTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _RMNHome.Contract.TransferOwnership(&_RMNHome.TransactOpts, to)
}

type RMNHomeActiveConfigRevokedIterator struct {
	Event *RMNHomeActiveConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeActiveConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeActiveConfigRevoked)
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
		it.Event = new(RMNHomeActiveConfigRevoked)
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

func (it *RMNHomeActiveConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeActiveConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeActiveConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeActiveConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeActiveConfigRevokedIterator{contract: _RMNHome.contract, event: "ActiveConfigRevoked", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ActiveConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeActiveConfigRevoked)
				if err := _RMNHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseActiveConfigRevoked(log types.Log) (*RMNHomeActiveConfigRevoked, error) {
	event := new(RMNHomeActiveConfigRevoked)
	if err := _RMNHome.contract.UnpackLog(event, "ActiveConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeCandidateConfigRevokedIterator struct {
	Event *RMNHomeCandidateConfigRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeCandidateConfigRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeCandidateConfigRevoked)
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
		it.Event = new(RMNHomeCandidateConfigRevoked)
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

func (it *RMNHomeCandidateConfigRevokedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeCandidateConfigRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeCandidateConfigRevoked struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeCandidateConfigRevokedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeCandidateConfigRevokedIterator{contract: _RMNHome.contract, event: "CandidateConfigRevoked", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "CandidateConfigRevoked", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeCandidateConfigRevoked)
				if err := _RMNHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseCandidateConfigRevoked(log types.Log) (*RMNHomeCandidateConfigRevoked, error) {
	event := new(RMNHomeCandidateConfigRevoked)
	if err := _RMNHome.contract.UnpackLog(event, "CandidateConfigRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeConfigPromotedIterator struct {
	Event *RMNHomeConfigPromoted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeConfigPromotedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeConfigPromoted)
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
		it.Event = new(RMNHomeConfigPromoted)
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

func (it *RMNHomeConfigPromotedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeConfigPromotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeConfigPromoted struct {
	ConfigDigest [32]byte
	Raw          types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigPromotedIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeConfigPromotedIterator{contract: _RMNHome.contract, event: "ConfigPromoted", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ConfigPromoted", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeConfigPromoted)
				if err := _RMNHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseConfigPromoted(log types.Log) (*RMNHomeConfigPromoted, error) {
	event := new(RMNHomeConfigPromoted)
	if err := _RMNHome.contract.UnpackLog(event, "ConfigPromoted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeConfigSetIterator struct {
	Event *RMNHomeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeConfigSet)
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
		it.Event = new(RMNHomeConfigSet)
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

func (it *RMNHomeConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNHomeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeConfigSet struct {
	ConfigDigest  [32]byte
	Version       uint32
	StaticConfig  RMNHomeStaticConfig
	DynamicConfig RMNHomeDynamicConfig
	Raw           types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeConfigSetIterator{contract: _RMNHome.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "ConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeConfigSet)
				if err := _RMNHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseConfigSet(log types.Log) (*RMNHomeConfigSet, error) {
	event := new(RMNHomeConfigSet)
	if err := _RMNHome.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeDynamicConfigSetIterator struct {
	Event *RMNHomeDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeDynamicConfigSet)
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
		it.Event = new(RMNHomeDynamicConfigSet)
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

func (it *RMNHomeDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *RMNHomeDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeDynamicConfigSet struct {
	ConfigDigest  [32]byte
	DynamicConfig RMNHomeDynamicConfig
	Raw           types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeDynamicConfigSetIterator, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "DynamicConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeDynamicConfigSetIterator{contract: _RMNHome.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeDynamicConfigSet, configDigest [][32]byte) (event.Subscription, error) {

	var configDigestRule []interface{}
	for _, configDigestItem := range configDigest {
		configDigestRule = append(configDigestRule, configDigestItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "DynamicConfigSet", configDigestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeDynamicConfigSet)
				if err := _RMNHome.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseDynamicConfigSet(log types.Log) (*RMNHomeDynamicConfigSet, error) {
	event := new(RMNHomeDynamicConfigSet)
	if err := _RMNHome.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeOwnershipTransferRequestedIterator struct {
	Event *RMNHomeOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeOwnershipTransferRequested)
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
		it.Event = new(RMNHomeOwnershipTransferRequested)
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

func (it *RMNHomeOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *RMNHomeOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeOwnershipTransferRequestedIterator{contract: _RMNHome.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeOwnershipTransferRequested)
				if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseOwnershipTransferRequested(log types.Log) (*RMNHomeOwnershipTransferRequested, error) {
	event := new(RMNHomeOwnershipTransferRequested)
	if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type RMNHomeOwnershipTransferredIterator struct {
	Event *RMNHomeOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *RMNHomeOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMNHomeOwnershipTransferred)
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
		it.Event = new(RMNHomeOwnershipTransferred)
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

func (it *RMNHomeOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *RMNHomeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type RMNHomeOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_RMNHome *RMNHomeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RMNHomeOwnershipTransferredIterator{contract: _RMNHome.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_RMNHome *RMNHomeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _RMNHome.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(RMNHomeOwnershipTransferred)
				if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_RMNHome *RMNHomeFilterer) ParseOwnershipTransferred(log types.Log) (*RMNHomeOwnershipTransferred, error) {
	event := new(RMNHomeOwnershipTransferred)
	if err := _RMNHome.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllConfigs struct {
	ActiveConfig    RMNHomeVersionedConfig
	CandidateConfig RMNHomeVersionedConfig
}
type GetConfig struct {
	VersionedConfig RMNHomeVersionedConfig
	Ok              bool
}
type GetConfigDigests struct {
	ActiveConfigDigest    [32]byte
	CandidateConfigDigest [32]byte
}

func (_RMNHome *RMNHome) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _RMNHome.abi.Events["ActiveConfigRevoked"].ID:
		return _RMNHome.ParseActiveConfigRevoked(log)
	case _RMNHome.abi.Events["CandidateConfigRevoked"].ID:
		return _RMNHome.ParseCandidateConfigRevoked(log)
	case _RMNHome.abi.Events["ConfigPromoted"].ID:
		return _RMNHome.ParseConfigPromoted(log)
	case _RMNHome.abi.Events["ConfigSet"].ID:
		return _RMNHome.ParseConfigSet(log)
	case _RMNHome.abi.Events["DynamicConfigSet"].ID:
		return _RMNHome.ParseDynamicConfigSet(log)
	case _RMNHome.abi.Events["OwnershipTransferRequested"].ID:
		return _RMNHome.ParseOwnershipTransferRequested(log)
	case _RMNHome.abi.Events["OwnershipTransferred"].ID:
		return _RMNHome.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (RMNHomeActiveConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x0b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8")
}

func (RMNHomeCandidateConfigRevoked) Topic() common.Hash {
	return common.HexToHash("0x53f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b")
}

func (RMNHomeConfigPromoted) Topic() common.Hash {
	return common.HexToHash("0xfc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e")
}

func (RMNHomeConfigSet) Topic() common.Hash {
	return common.HexToHash("0xf6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf0")
}

func (RMNHomeDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e6294")
}

func (RMNHomeOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (RMNHomeOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_RMNHome *RMNHome) Address() common.Address {
	return _RMNHome.address
}

type RMNHomeInterface interface {
	GetActiveDigest(opts *bind.CallOpts) ([32]byte, error)

	GetAllConfigs(opts *bind.CallOpts) (GetAllConfigs,

		error)

	GetCandidateDigest(opts *bind.CallOpts) ([32]byte, error)

	GetConfig(opts *bind.CallOpts, configDigest [32]byte) (GetConfig,

		error)

	GetConfigDigests(opts *bind.CallOpts) (GetConfigDigests,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	PromoteCandidateAndRevokeActive(opts *bind.TransactOpts, digestToPromote [32]byte, digestToRevoke [32]byte) (*types.Transaction, error)

	RevokeCandidate(opts *bind.TransactOpts, configDigest [32]byte) (*types.Transaction, error)

	SetCandidate(opts *bind.TransactOpts, staticConfig RMNHomeStaticConfig, dynamicConfig RMNHomeDynamicConfig, digestToOverwrite [32]byte) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, newDynamicConfig RMNHomeDynamicConfig, currentDigest [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterActiveConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeActiveConfigRevokedIterator, error)

	WatchActiveConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeActiveConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseActiveConfigRevoked(log types.Log) (*RMNHomeActiveConfigRevoked, error)

	FilterCandidateConfigRevoked(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeCandidateConfigRevokedIterator, error)

	WatchCandidateConfigRevoked(opts *bind.WatchOpts, sink chan<- *RMNHomeCandidateConfigRevoked, configDigest [][32]byte) (event.Subscription, error)

	ParseCandidateConfigRevoked(log types.Log) (*RMNHomeCandidateConfigRevoked, error)

	FilterConfigPromoted(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigPromotedIterator, error)

	WatchConfigPromoted(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigPromoted, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigPromoted(log types.Log) (*RMNHomeConfigPromoted, error)

	FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*RMNHomeConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (*RMNHomeDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *RMNHomeDynamicConfigSet, configDigest [][32]byte) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*RMNHomeDynamicConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*RMNHomeOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RMNHomeOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RMNHomeOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*RMNHomeOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var RMNHomeZKBin = ("0x0002000000000002000d00000000000200010000000103550000006004100270000003720040019d0000008003000039000000400030043f0000000100200190000000210000c13d0000037202400197000000040020008c0000094b0000413d000000000401043b000000e004400270000003780040009c000000390000a13d000003790040009c000000550000a13d0000037a0040009c000000c70000213d0000037d0040009c000000de0000613d0000037e0040009c0000094b0000c13d0000000001000416000000000001004b0000094b0000c13d0000000101000039000000000101041a0000038901100197000000800010043f000003910100004100000dc50001042e0000000001000416000000000001004b0000094b0000c13d0000000001000411000000000001004b0000002b0000c13d0000037601000041000000800010043f000003770100004100000dc6000104300000000102000039000000000302041a0000037303300197000000000113019f000000000012041b0000000e01000039000000000201041a0000037402200197000000000021041b000000200100003900000100001004430000012000000443000003750100004100000dc50001042e000003820040009c000000740000213d000003860040009c000001590000613d000003870040009c000002110000613d000003880040009c0000094b0000c13d0000000001000416000000000001004b0000094b0000c13d000000c001000039000000400010043f0000001101000039000000800010043f000003aa01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e0020000390dc40b690000040f000000c00110008a000003720010009c00000372010080410000006001100210000003ab011001c700000dc50001042e0000037f0040009c000001090000613d000003800040009c000001350000613d000003810040009c0000094b0000c13d0000000001000416000000000001004b0000094b0000c13d000000000100041a00000389021001970000000006000411000000000026004b0000023f0000c13d0000000102000039000000000302041a0000037304300197000000000464019f000000000042041b0000037301100197000000000010041b00000000010004140000038905300197000003720010009c0000037201008041000000c0011002100000038c011001c70000800d0200003900000003030000390000039704000041000002b60000013d000003830040009c000001fd0000613d000003840040009c000002200000613d000003850040009c0000094b0000c13d0000000001000416000000000001004b0000094b0000c13d000000800000043f000000a00000043f0000006001000039000001000010043f000001200010043f0000010002000039000000c00020043f000001400010043f000001600010043f0000014002000039000000e00020043f000001800000043f000001a00000043f000002000010043f000002200010043f0000020002000039000001c00020043f0000028002000039000000400020043f000002400010043f000002600010043f0000024001000039000001e00010043f0000000e01000039000000000101041a00000020011002700000037201100197000000010010008c0000021a0000213d000b00000001001d00000006011000c900000002011000390dc40c940000040f000800000001001d00000020021000390000000001020433000700000001001d0000000b01000029000000010110015f00000006011000c900000002011000390dc40c940000040f000a00000001001d00000020021000390000000001020433000900000001001d0000004001000039000000400200043d000b00000002001d0000000001120436000600000001001d000000070000006b0000000801000029000000800100603900000040022000390dc40b7b0000040f00000000020100190000000b0120006a00000006030000290000000000130435000000090000006b0000000a0100002900000180010060390dc40b7b0000040f0000000b020000290000000001210049000003720020009c00000372020080410000004002200210000003720010009c00000372010080410000006001100210000000000121019f00000dc50001042e0000037b0040009c000000f30000613d0000037c0040009c0000094b0000c13d000000240020008c0000094b0000413d0000000002000416000000000002004b0000094b0000c13d0000000401100370000000000501043b0000000101000039000000000101041a00000389011001970000000002000411000000000012004b0000022b0000c13d000000000005004b0000028e0000c13d0000038e01000041000000800010043f000003770100004100000dc600010430000000440020008c0000094b0000413d0000000002000416000000000002004b0000094b0000c13d0000002402100370000000000502043b0000000401100370000000000401043b0000000101000039000000000101041a00000389011001970000000002000411000000000012004b0000022b0000c13d00000000005401a0000002430000c13d0000039501000041000000800010043f000003770100004100000dc600010430000000240020008c0000094b0000413d0000000002000416000000000002004b0000094b0000c13d0000000401100370000000000101043b000003890010009c0000094b0000213d0000000102000039000000000202041a00000389022001970000000005000411000000000025004b0000022b0000c13d0000038906100197000000000056004b000002aa0000c13d0000039001000041000000800010043f000003770100004100000dc600010430000000240020008c0000094b0000413d0000000002000416000000000002004b0000094b0000c13d0000000401100370000000000101043b000000800000043f000000a00000043f0000006002000039000001000020043f000001200020043f0000010004000039000000c00040043f0000018004000039000000400040043f000001400020043f000001600020043f0000014002000039000000e00020043f000000000001004b0000022f0000c13d0000000002000019000900000004001d000a00000002001d00000040010000390000000001140436000b00000001001d000000400240003900000000010300190dc40b7b0000040f0000000b020000290000000a03000029000000000032043500000009020000290000000001210049000003720010009c00000372010080410000006001100210000003720020009c00000372020080410000004002200210000000000121019f00000dc50001042e000000440020008c0000094b0000413d0000000003000416000000000003004b0000094b0000c13d0000000403100370000000000303043b000b00000003001d000003980030009c0000094b0000213d0000000b0220006a000003990020009c0000094b0000213d000000440020008c0000094b0000413d0000000102000039000000000202041a00000389022001970000000003000411000000000023004b0000022b0000c13d0000002401100370000000000301043b0000000002000415000000000003004b000002bb0000c13d000000400100043d000003a502000041000000000021043500000004021000390000000000320435000003720010009c00000372010080410000004001100210000003a6011001c700000dc600010430000000640020008c0000094b0000413d0000000003000416000000000003004b0000094b0000c13d0000000403100370000000000403043b000003980040009c0000094b0000213d0000000003420049000003990030009c0000094b0000213d000000440030008c0000094b0000413d0000002403100370000000000303043b000003980030009c0000094b0000213d0000000005320049000003990050009c0000094b0000213d000000440050008c0000094b0000413d0000000105000039000000000505041a00000389055001970000000006000411000000000056004b0000022b0000c13d000000c005000039000000400050043f0000000406400039000000000761034f000000000707043b000003980070009c0000094b0000213d00000000074700190000002308700039000000000028004b0000094b0000813d000000040d7000390000000008d1034f000000000808043b0000039e0080009c000002880000813d00000005098002100000003f099000390000039c09900197000003ad0090009c000002880000213d000000c009900039000000400090043f000000c00080043f000b00240070003d00000006078002100000000b07700029000000000027004b0000094b0000213d000000000008004b000001ac0000613d000000e0080000390000000b09000029000000000a9200490000039900a0009c0000094b0000213d0000004000a0008c0000094b0000413d000000400a00043d0000039a00a0009c000002880000213d000000400ba000390000004000b0043f000000000b91034f000000000b0b043b000000000bba0436000000200c900039000000000cc1034f000000000c0c043b0000000000cb04350000000008a804360000004009900039000000000079004b000001970000413d000000800050043f0000002005600039000000000551034f000000000505043b000003980050009c0000094b0000213d00000000074500190000002304700039000000000024004b00000000050000190000039b050080410000039b04400197000000000004004b00000000060000190000039b060040410000039b0040009c000000000605c019000000000006004b0000094b0000c13d0000000408700039000000000481034f000000000404043b000003980040009c000002880000213d0000001f05400039000003bf055001970000003f05500039000003bf06500197000000400500043d0000000006650019000000000056004b00000000090000390000000109004039000003980060009c000002880000213d0000000100900190000002880000c13d000000400060043f000000000645043600000000074700190000002407700039000000000027004b0000094b0000213d00080000000d001d0000002007800039000000000771034f000003bf084001980000001f0940018f0000000001860019000001e40000613d000000000a07034f000000000b06001900000000ac0a043c000000000bcb043600000000001b004b000001e00000c13d000000000009004b000001f10000613d000000000787034f0000000308900210000000000901043300000000098901cf000000000989022f000000000707043b0000010008800089000000000787022f00000000078701cf000000000797019f000000000071043500000000014600190000000000010435000000a00050043f00000004013000390dc40be30000040f000000800200043d0000000032020434000001000020008c000004b70000a13d000000400100043d000003be02000041000005100000013d0000000001000416000000000001004b0000094b0000c13d0000000e01000039000000000101041a00000020011002700000037201100197000000010010008c0000021a0000213d000000010210015f00000006022000c90000000302200039000000000202041a00000006011000c90000000301100039000000000101041a000000800010043f000000a00020043f000003a90100004100000dc50001042e0000000001000416000000000001004b0000094b0000c13d0000000e01000039000000000101041a00000020011002700000037201100197000000010010008c000002390000a13d000003ba01000041000000000010043f0000003201000039000000040010043f000003a60100004100000dc6000104300000000001000416000000000001004b0000094b0000c13d0dc40c830000040f000000400200043d0000000000120435000003720020009c00000372020080410000004001200210000003a8011001c700000dc50001042e000003ac01000041000000800010043f000003770100004100000dc6000104300000000302000039000000000202041a000000000012004b000002700000c13d000000070400003900000006050000390000000506000039000000040200003900000002030000390000027a0000013d00000006011000c90000000301100039000000000101041a000000800010043f000003910100004100000dc50001042e0000039601000041000000800010043f000003770100004100000dc6000104300000000e01000039000000000201041a00000020022002700000037202200197000000010020008c0000021a0000213d000000010320015f00000006033000c90000000303300039000000000303041a000000000043004b000003050000c13d00000006022000c90000000303200039000000000203041a000000000052004b000003d80000c13d000b00000004001d000000000003041b000000000201041a0000039202200167000000000021041b000000000005004b000002660000613d0000000001000414000003720010009c0000037201008041000000c0011002100000038c011001c70000800d02000039000000020300003900000393040000410dc40dba0000040f00000001002001900000094b0000613d0000000001000414000003720010009c0000037201008041000000c0011002100000038c011001c70000800d02000039000000020300003900000394040000410000000b05000029000002b60000013d0000000902000039000000000202041a000000000012004b0000000002000019000001200000c13d0000000d040000390000000c050000390000000b060000390000000a020000390000000803000039000000000303041a0000037203300197000001800030043f000001a00010043f0000024001000039000000400010043f000000000302041a000003980030009c000002880000213d00000005013002100000003f011000390000039c01100197000003a70010009c000002ca0000a13d000003ba01000041000000000010043f0000004101000039000000040010043f000003a60100004100000dc6000104300000000e01000039000000000101041a00000020011002700000037201100197000000010010008c0000021a0000213d000000010110015f00000006011000c90000000302100039000000000102041a000000000051004b000003d20000c13d000b00000002001d0000000001000414000003720010009c0000037201008041000000c0011002100000038c011001c70000800d0200003900000002030000390000038d040000410dc40dba0000040f00000001002001900000094b0000613d0000000b01000029000000000001041b000000000100001900000dc50001042e000000000100041a0000037301100197000000000161019f000000000010041b0000000001000414000003720010009c0000037201008041000000c0011002100000038c011001c70000800d0200003900000003030000390000038f040000410dc40dba0000040f00000001002001900000094b0000613d000000000100001900000dc50001042e00000000040004150000000d0440008a0000000504400210000d00030000003d0000000301000039000000000101041a000000000031004b0000030b0000c13d000500000004001d000700000003001d000800000002001d000600070000003d000900060000003d0000000403000039000003190000013d000a00000006001d000b00000005001d000800000004001d0000024001100039000000400010043f000900000003001d000002400030043f000000000020043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d0000000909000029000000000009004b0000000b060000290000000a070000290000024008000039000002f30000613d000000000101043b0000000002000019000000400300043d0000039a0030009c000002880000213d0000004004300039000000400040043f000000000401041a00000000044304360000000105100039000000000505041a00000000005404350000002008800039000000000038043500000002011000390000000102200039000000000092004b000002e30000413d0000024001000039000002000010043f000000000107041a000000010210019000000001041002700000007f0440618f0000001f0040008c00000000030000390000000103002039000000000331013f0000000100300190000003de0000613d000003ba01000041000000000010043f0000002201000039000000040010043f000003a60100004100000dc6000104300000038a01000041000000800010043f000000840030043f000000a40040043f0000038b0100004100000dc60001043000000000040004150000000c0440008a0000000504400210000c00030000003d0000000901000039000000000101041a000000000031004b0000014f0000c13d000500000004001d000700000003001d000800000002001d0006000d0000003d0009000c0000003d0000000a030000390000000b01000029000a00040010003d00000000010000310000000a0210006a000003990020009c0000094b0000213d000000400020008c0000094b0000413d000000400500043d0000039a0050009c000002880000213d000000000303041a0000004007500039000000400070043f00000001060003670000000a04600360000000000404043b000003980040009c0000094b0000213d0000000a094000290000001f08900039000000000018004b000000000a0000190000039b0a0080410000039b0b8001970000039b08100197000000000c8b013f00000000008b004b000000000b0000190000039b0b0040410000039b00c0009c000000000b0ac01900000000000b004b0000094b0000c13d000000000a96034f000000000a0a043b00040000000a001d0000039800a0009c000002880000213d000000040a000029000000050aa002100000003f0aa000390000039c0aa00197000000000a7a00190000039800a0009c000002880000213d0000004000a0043f000000040a0000290000000000a70435000300200090003d0000006009a000c90000000309900029000000000019004b0000094b0000213d000000040000006b0000047c0000c13d00000000077504360000000a090000290000002009900039000000000996034f000000000909043b000003980090009c0000094b0000213d0000000a0b9000290000001f09b00039000000000019004b000000000a0000190000039b0a0080410000039b09900197000000000c89013f000000000089004b00000000080000190000039b080040410000039b00c0009c00000000080ac019000000000008004b0000094b0000c13d0000000008b6034f000000000808043b000003980080009c000002880000213d0000001f09800039000003bf099001970000003f09900039000003bf0a900197000000400900043d000000000aa9001900000000009a004b000000000c000039000000010c0040390000039800a0009c000002880000213d0000000100c00190000002880000c13d0000004000a0043f000000000a890436000000200bb00039000000000cb8001900000000001c004b0000094b0000213d0002000000b60353000003bf0c8001980000001f0d80018f0000000006ca0019000003880000613d000000020e00035f000000000f0a001900000000eb0e043c000000000fbf043600000000006f004b000003840000c13d00000000000d004b000003950000613d000000020bc0035f000000030cd00210000000000d060433000000000dcd01cf000000000dcd022f000000000b0b043b000001000cc00089000000000bcb022f000000000bcb01cf000000000bdb019f0000000000b604350000000b06000029000200240060003d00000000068a0019000000000006043500000000009704350000010007300089000003c00670027f000000ff0070008c00000000060020190000000007050433000000200570003900000000070704330000000008000019000000000078004b000005430000813d0000000509800210000000000959001900000000090904330000000108800039000000000078004b000003b60000813d000000000a090433000000000b080019000000050cb00210000000000c5c0019000000000c0c0433000000000c0c0433000000000cac013f0000039800c001980000050b0000613d000000010bb0003900000000007b004b000003ac0000413d000001000030008c0000053d0000213d000000400a900039000000000b0a0433000000000a6b016f0000000000ba004b000005160000c13d00000000000b004b000003c60000613d000000000a000019000000010aa0003a0000053d0000613d000000010cb0008a000000000bbc0170000003c00000c13d000003c70000013d000000000a00001900000020099000390000000009090433000000010b900210000003a309b00197000003a40bb001970000000000b9004b0000053d0000c13d00000001099001bf00000000009a004b000003a20000813d000004ee0000013d0000038a02000041000000800020043f000000840010043f000000a40050043f0000038b0100004100000dc6000104300000038a01000041000000800010043f000000840020043f000000a40050043f0000038b0100004100000dc600010430000000400500043d0000000003450436000000000002004b000004000000613d000600000003001d000700000004001d000900000005001d000000000070043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d0000000707000029000000000007004b00000000020000190000000b0600002900000009050000290000000608000029000004050000613d000000000101043b00000000020000190000000003280019000000000401041a000000000043043500000001011000390000002002200039000000000072004b000003f80000413d000004050000013d000003c1011001970000000000130435000000000004004b000000200200003900000000020060390000003f01200039000003bf011001970000000002510019000000000012004b00000000010000390000000101004039000003980020009c000002880000213d0000000100100190000002880000c13d000000400020043f000002200050043f0000020001000039000001c00010043f0000039a0020009c000002880000213d0000004001200039000b00000001001d000000400010043f000000000406041a000003980040009c000002880000213d00000005014002100000003f011000390000039c011001970000000b01100029000003980010009c000002880000213d000a00000002001d000000400010043f0000000b01000029000900000004001d0000000000410435000000000060043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d000000090a00002900000000000a004b00000008080000290000000a090000290000044d0000613d0000006002900039000000000101043b0000000003000019000000400400043d0000039d0040009c000002880000213d0000006005400039000000400050043f000000000501041a0000004006500270000003980660019700000020074000390000000000670435000003980550019700000000005404350000000105100039000000000505041a000000400640003900000000005604350000000002420436000000020110003900000001033000390000000000a3004b000004380000413d0000000b010000290000000001190436000900000001001d000000000108041a000000010210019000000001031002700000007f0330618f000b00000003001d0000001f0030008c00000000030000390000000103002039000000000331013f0000000100300190000002ff0000c13d000000400300043d000700000003001d0000000b040000290000000003430436000600000003001d000000000002004b0000049d0000613d0000000801000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d0000000b05000029000000000005004b00000000020000190000000606000029000004a30000613d000000000101043b00000000020000190000000003260019000000000401041a000000000043043500000001011000390000002002200039000000000052004b000004740000413d000004a30000013d000000600a500039000000030b000029000000000cb100490000039900c0009c0000094b0000213d0000006000c0008c0000094b0000413d000000400c00043d0000039d00c0009c000002880000213d000000600dc000390000004000d0043f000000000db6034f000000000d0d043b0000039800d0009c0000094b0000213d000000000edc0436000000200db00039000000000fd6034f000000000f0f043b0000039800f0009c0000094b0000213d0000000000fe0435000000200dd00039000000000dd6034f000000000d0d043b000000400ec000390000000000de0435000000000aca0436000000600bb0003900000000009b004b0000047e0000413d000003510000013d000003c101100197000000060200002900000000001204350000000b0000006b000000200200003900000000020060390000003f01200039000003bf021001970000000701200029000000000021004b00000000020000390000000102004039000003980010009c0000000a03000029000002880000213d0000000100200190000002880000c13d000000400010043f000000090100002900000007020000290000000000210435000001e00030043f00000001020000390000018003000039000000400400043d000001200000013d000000000002004b000004f10000c13d0000010003200089000003c00230027f000000ff0030008c00000000020020190000000003010433000000200130003900000000030304330000000004000019000000000034004b000005190000813d0000000505400210000000000515001900000000050504330000000104400039000000000034004b000004d50000813d000000000605043300000000070400190000000508700210000000000818001900000000080804330000000008080433000000000868013f00000398008001980000050b0000613d0000000107700039000000000037004b000004cb0000413d00000040065000390000000007060433000000000627016f000000000076004b000005160000c13d000000000007004b000004e30000613d0000000006000019000000010660003a0000053d0000613d000000010870008a0000000007780170000004dd0000c13d000004e40000013d0000000006000019000000200550003900000000050504330000000107500210000003a305700197000003a407700197000000000075004b0000053d0000c13d00000001055001bf000000000056004b000004c10000813d000000400100043d000003bc02000041000005100000013d0000000004000019000004f60000013d0000000104400039000000000024004b000004b90000813d0000000505400210000000000535001900000000060400190000000106600039000000000026004b000004f30000813d00000005076002100000000007370019000000000707043300000000790704340000000008050433000000008a08043400000000009a004b0000050e0000613d00000000070704330000000008080433000000000078004b000004f90000c13d000000400100043d000003ae02000041000005100000013d000000400100043d000003bd02000041000005100000013d000000400100043d000003b0020000410000000000210435000003720010009c00000372010080410000004001100210000003af011001c700000dc600010430000000400100043d000003bb02000041000005100000013d0000000e01000039000000000101041a000700000001001d00000020011002700000037202100197000000020020008c0000021a0000813d00000044010000390000000101100367000000000101043b000000010220015f00000006022000c90000000302200039000000000502041a000000000015004b0000061c0000c13d000000000005004b000005390000613d0000000001000414000003720010009c0000037201008041000000c0011002100000038c011001c70000800d0200003900000002030000390000038d040000410dc40dba0000040f00000001002001900000094b0000613d0000000e01000039000000000101041a000700000001001d00000007010000290000037201100197000003720010009c000006280000c13d000003ba01000041000000000010043f0000001101000039000000040010043f000003a60100004100000dc6000104300000001f0220008a000000000024004b00000000030000190000039b030080410000039b044001970000039b02200197000000000524013f000000000024004b00000000020000190000039b020040410000039b0050009c000000000203c019000000000002004b0000094b0000c13d000000600300008a00000004023000b900000000011200190000000303000029000000000013004b00000000020000190000039b020020410000039b011001970000039b03300197000000000413013f000000000013004b00000000010000190000039b010040410000039b0040009c000000000102c019000000000001004b0000094b0000c13d00000004010000290000039e0010009c000002880000213d0000000902000029000000000302041a0000000401000029000000000012041b000100000003001d000000000031004b000005890000813d0000000101000029000003990010009c0000053d0000213d0000000901000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d00000004020000290000000102200210000000000301043b0000000001230019000000010200002900000001022002100000000002230019000000000021004b000005890000813d000000000001041b0000000103100039000000000003041b0000000201100039000000000021004b000005830000413d0000000901000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d0000000102000367000000040000006b000005b50000613d000000000101043b00000000030000190000000304200360000000000404043b000003980040009c0000094b0000213d00000003050000290000002005500039000000000652034f000000000606043b000003980060009c0000094b0000213d000000000701041a000003a0077001970000004006600210000003a106600197000000000676019f000000000446019f000000000041041b0000002004500039000000000442034f000000000404043b0000000105100039000000000045041b00000002011000390000000304000029000300600040003d0000000103300039000000040030006c000005990000413d0000000201200360000000000301043b00000000010000310000000b0410006a000000230440008a0000039b054001970000039b06300197000000000756013f000000000056004b00000000050000190000039b05004041000000000043004b00000000040000190000039b040080410000039b0070009c000000000504c019000000000005004b0000094b0000c13d0000000a03300029000000000232034f000000000202043b000900000002001d000003980020009c0000094b0000213d000000090110006a00000020053000390000039b021001970000039b03500197000000000423013f000000000023004b00000000020000190000039b02004041000b00000005001d000000000015004b00000000010000190000039b010020410000039b0040009c000000000201c019000000000002004b0000094b0000c13d0000000601000029000000000101041a000000010010019000000001021002700000007f0220618f000400000002001d0000001f0020008c00000000020000390000000102002039000000000121013f0000000100100190000002ff0000c13d0000000401000029000000200010008c000006080000413d0000000601000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d00000009030000290000001f023000390000000502200270000000200030008c0000000002004019000000000301043b00000004010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b000006080000813d000000000002041b0000000102200039000000000012004b000006040000413d00000009010000290000001f0010008c0000081a0000a13d0000000601000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d000000200200008a0000000903200180000000000101043b0000082f0000c13d0000000002000019000008390000013d000000400200043d000000240320003900000000001304350000038a01000041000000000012043500000004012000390000000000510435000003720020009c00000372020080410000004001200210000003b1011001c700000dc6000104300000000702000029000003b2012001970000000102200039000503720020019b00000005011001af0000000e02000039000000000012041b000000400100043d000600000001001d00000020021000390000002001000039000a00000002001d000000000012043500000001010003670000000402100370000000000202043b0000000403200039000000000231034f000000000602043b000000000200003100000000043200490000001f0440008a0000039b054001970000039b07600197000000000857013f000000000057004b00000000070000190000039b07004041000000000046004b00000000090000190000039b090080410000039b0080009c000000000709c019000000000007004b0000094b0000c13d0000000006360019000000000761034f000000000807043b000003980080009c0000094b0000213d000000200960003900000006068002100000000006620049000000000069004b00000000070000190000039b070020410000039b066001970000039b0a900197000000000b6a013f00000000006a004b00000000060000190000039b060040410000039b00b0009c000000000607c019000000000006004b0000094b0000c13d000000060b0000290000008007b000390000004006b00039000000400a0000390000000000a604350000000000870435000000a007b00039000000000008004b000006760000613d000000000a000019000000000b91034f000000000b0b043b000000000bb70436000000200c900039000000000cc1034f000000000c0c043b0000000000cb043500000040099000390000004007700039000000010aa0003900000000008a004b0000066a0000413d0000002008300039000000000881034f000000000808043b0000039b09800197000000000a59013f000000000059004b00000000050000190000039b05004041000000000048004b00000000040000190000039b040080410000039b00a0009c000000000504c019000000000005004b0000094b0000c13d0000000004380019000000000341034f000000000303043b000003980030009c0000094b0000213d00000020044000390000000002320049000000000024004b00000000050000190000039b050020410000039b022001970000039b08400197000000000928013f000000000028004b00000000020000190000039b020040410000039b0090009c000000000205c019000000000002004b0000094b0000c13d0000000002670049000000060500002900000060055000390000000000250435000000000441034f0000000001370436000003bf053001980000001f0630018f0000000002510019000006a90000613d000000000704034f0000000008010019000000007907043c0000000008980436000000000028004b000006a50000c13d000000000006004b000006b60000613d000000000454034f0000000305600210000000000602043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000042043500000000021300190000000000020435000000060400002900000000014100490000001f02300039000003bf022001970000000001120019000000200210008a00000000002404350000001f01100039000003bf011001970000000002410019000000000012004b00000000010000390000000101004039000400000002001d000003980020009c000002880000213d0000000100100190000002880000c13d0000000401000029000000400010043f0000002002100039000003b301000041000900000002001d0000000000120435000003b40100004100000000001004430000000001000414000003720010009c0000037201008041000000c001100210000003b5011001c70000800b020000390dc40dbf0000040f0000000100200190000008190000613d000000000101043b00000004040000290000008002400039000000050300002900000000003204350000006002400039000000000300041000000000003204350000004002400039000000000012043500000080010000390000000000140435000003b60040009c000002880000213d0000000403000029000000c002300039000000a001300039000000400010043f0000000003030433000000000003004b000006f80000613d000000000400001900000000052400190000000906400029000000000606043300000000006504350000002004400039000000000034004b000006f10000413d0000000004230019000000000004043500000006050000290000000005050433000000000005004b000007060000613d000000000600001900000000074600190000000a08600029000000000808043300000000008704350000002006600039000000000056004b000006ff0000413d00000000044500190000000000040435000000000335001900000000003104350000003f03300039000003bf043001970000000003140019000000000043004b00000000040000390000000104004039000003980030009c000002880000213d0000000100400190000002880000c13d000000400030043f000003720020009c000003720200804100000040022002100000000001010433000003720010009c00000372010080410000006001100210000000000121019f0000000002000414000003720020009c0000037202008041000000c002200210000000000112019f0000038c011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d000000000201043b000000070100002900000020011002700000037201100197000000010010008c0000021a0000213d000003b702200197000003b8022001c7000000010110015f00000006031000c90000000301300039000900000002001d000000000021041b000a00000003001d0000000201300039000000000201041a000003b20220019700000005022001af000000000021041b00000001020003670000000401200370000000000301043b00000000010000310000000004310049000000230440008a0000000403300039000000000332034f000000000303043b000000000043004b00000000050000190000039b050080410000039b044001970000039b03300197000000000643013f000000000043004b00000000030000190000039b030040410000039b0060009c000000000305c019000000000003004b00000008030000290000094b0000c13d000000000232034f000000000202043b000003980020009c0000094b0000213d000000060320021000000000013100490000000b04000029000000000014004b00000000030000190000039b030020410000039b011001970000039b04400197000000000514013f000000000014004b00000000010000190000039b010040410000039b0050009c000000000103c019000000000001004b0000094b0000c13d0000000a010000290000000401100039000000000301041a000700000001001d000000000021041b000600000003001d000000000032004b0000078c0000813d0000000601000029000003990010009c0000053d0000213d0000000701000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f000000010020019000000008040000290000094b0000613d00000006020000290000000102200210000000000301043b00000000012300190000000102400367000000000202043b00000001022002100000000002320019000000000012004b0000078c0000813d000000000002041b0000000103200039000000000003041b0000000202200039000000000012004b000007860000413d0000000701000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f000000010020019000000008030000290000094b0000613d0000000102000367000000000332034f000000000303043b000000000003004b000007ad0000613d000000000101043b00000000040000190000000b07000029000000000572034f000000000505043b000000000051041b0000002005700039000000000552034f000000000505043b0000000106100039000000000056041b0000000201100039000b00400070003d0000000104400039000000000034004b0000079f0000413d0000000401200370000000000301043b0000002401300039000000000112034f000000000401043b00000000010000310000000005310049000000230550008a0000039b065001970000039b07400197000000000867013f000000000067004b00000000060000190000039b06004041000000000054004b00000000050000190000039b050080410000039b0080009c000000000605c019000000000006004b0000094b0000c13d00000000034300190000000403300039000000000232034f000000000202043b000800000002001d000003980020009c0000094b0000213d000000080110006a00000020053000390000039b021001970000039b03500197000000000423013f000000000023004b00000000020000190000039b02004041000b00000005001d000000000015004b00000000010000190000039b010020410000039b0040009c000000000201c019000000000002004b0000094b0000c13d0000000a010000290000000501100039000700000001001d000000000101041a000000010210019000000001011002700000007f0110618f000600000001001d0000001f0010008c00000000010000390000000101002039000000000012004b000002ff0000c13d0000000601000029000000200010008c000008050000413d0000000701000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d00000008030000290000001f023000390000000502200270000000200030008c0000000002004019000000000301043b00000006010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b000008050000813d000000000002041b0000000102200039000000000012004b000008010000413d00000008010000290000001f0010008c000008f50000a13d0000000701000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d000000200200008a0000000802200180000000000101043b000009030000c13d00000000030000190000090d0000013d000000000001042f000000090000006b0000000001000019000008200000613d0000000b010000290000000101100367000000000101043b000000050200002900000005022002700000000002020031000000090500002900000000032501cf000000ff0020008c0000000003002019000000010400008a000000000234022f000000ff0030008c000000000442a13f000000000114016f0000000102500210000000000121019f0000084c0000013d000000010400036700000000020000190000000b05200029000000000554034f000000000505043b000000000051041b00000001011000390000002002200039000000000032004b000008310000413d000000090030006c000008490000813d00000005030000290000000503300270000000000303003100000009043001ef000000f80440018f000003c00440027f000003c004400167000000ff0030008c00000000040020190000000b022000290000000102200367000000000202043b000000000224016f000000000021041b0000000901000029000000010110021000000001011001bf0000000602000029000000000012041b0000002002000039000000400100043d00000000032104360000000202000029000000200420008a0000000102000367000000000442034f000000000704043b00000000040000310000000a0540006a0000001f0550008a0000039b065001970000039b08700197000000000968013f000000000068004b00000000080000190000039b08004041000000000057004b000000000a0000190000039b0a0080410000039b0090009c00000000080ac019000000000008004b0000094b0000c13d0000000a07700029000000000872034f000000000808043b000003980080009c0000094b0000213d0000002009700039000003c2078000d10000000007470019000000000079004b000000000a0000190000039b0a0020410000039b077001970000039b0b900197000000000c7b013f00000000007b004b00000000070000190000039b070040410000039b00c0009c00000000070ac019000000000007004b0000094b0000c13d0000006007100039000000400a0000390000000000a3043500000000008704350000008007100039000000000008004b000008980000613d000000000a000019000000000b92034f000000000b0b043b0000039800b0009c0000094b0000213d000000000cb70436000000200b900039000000000db2034f000000000d0d043b0000039800d0009c0000094b0000213d0000000000dc0435000000200bb00039000000000bb2034f000000000b0b043b000000400c7000390000000000bc043500000060099000390000006007700039000000010aa0003900000000008a004b000008830000413d0000000208200360000000000808043b0000039b09800197000000000a69013f000000000069004b00000000060000190000039b06004041000000000058004b00000000050000190000039b050080410000039b00a0009c000000000605c019000000000006004b0000094b0000c13d0000000a06800029000000000562034f000000000505043b000003980050009c0000094b0000213d00000020066000390000000004540049000000000046004b00000000080000190000039b080020410000039b044001970000039b09600197000000000a49013f000000000049004b00000000040000190000039b040040410000039b00a0009c000000000408c019000000000004004b0000094b0000c13d000000000337004900000040041000390000000000340435000000000462034f0000000002570436000003bf065001980000001f0750018f0000000003620019000008c90000613d000000000804034f0000000009020019000000008a08043c0000000009a90436000000000039004b000008c50000c13d000000000007004b000008d60000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f0000000000430435000000000325001900000000000304350000001f03500039000003bf0330019700000000021200490000000002320019000003720020009c00000372020080410000006002200210000003720010009c00000372010080410000004001100210000000000112019f0000000002000414000003720020009c0000037202008041000000c002200210000000000112019f0000038c011001c70000800d020000390000000203000039000003a20400004100000007050000290dc40dba0000040f00000001002001900000094b0000613d000000000100041500000008011000690000000001000002000000000100001900000dc50001042e000000080000006b0000000001000019000008fb0000613d0000000b010000290000000101100367000000000101043b00000008040000290000000302400210000003c00220027f000003c002200167000000000121016f0000000102400210000000000121019f0000091c0000013d000000010400036700000000030000190000000b05300029000000000554034f000000000505043b000000000051041b00000001011000390000002003300039000000000023004b000009050000413d000000080020006c000009190000813d00000008020000290000000302200210000000f80220018f000003c00220027f000003c0022001670000000b033000290000000103300367000000000303043b000000000223016f000000000021041b0000000801000029000000010110021000000001011001bf0000000702000029000000000012041b00000001020003670000002401200370000000000301043b0000000401300039000000000112034f000000000401043b00000000010000310000000005310049000000230550008a0000039b065001970000039b07400197000000000867013f000000000067004b00000000060000190000039b06002041000000000054004b00000000050000190000039b050040410000039b0080009c000000000605c019000000000006004b0000094b0000613d0000000003340019000800040030003d0000000802200360000000000202043b000003980020009c0000094b0000213d000003c2042000d1000000000114001900000024063000390000039b031001970000039b04600197000000000534013f000000000034004b00000000030000190000039b03004041000b00000006001d000000000016004b00000000010000190000039b010020410000039b0050009c000000000301c019000000000003004b0000094d0000613d000000000100001900000dc6000104300000000a010000290000000601100039000000000301041a000700000001001d000000000021041b000600000003001d000000000032004b000009740000813d0000000601000029000003990010009c0000053d0000213d0000000701000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d00000006020000290000000102200210000000000301043b000000000123001900000008020000290000000102200367000000000202043b00000001022002100000000002320019000000000012004b000009740000813d000000000002041b0000000103200039000000000003041b0000000202200039000000000012004b0000096e0000413d0000000701000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d00000001020003670000000803200360000000000303043b000000000003004b000009a20000613d000000000101043b00000000040000190000000b05200360000000000505043b000003980050009c0000094b0000213d0000000b060000290000002006600039000000000762034f000000000707043b000003980070009c0000094b0000213d000000000801041a000003a0088001970000004007700210000003a107700197000000000787019f000000000557019f000000000051041b0000002005600039000000000552034f000000000505043b0000000106100039000000000056041b00000002011000390000000b05000029000b00600050003d0000000104400039000000000034004b000009860000413d0000002401200370000000000301043b0000002401300039000000000112034f000000000401043b00000000010000310000000005310049000000230550008a0000039b065001970000039b07400197000000000867013f000000000067004b00000000060000190000039b06004041000000000054004b00000000050000190000039b050080410000039b0080009c000000000605c019000000000006004b0000094b0000c13d00000000034300190000000403300039000000000232034f000000000202043b000800000002001d000003980020009c0000094b0000213d000000080110006a00000020053000390000039b021001970000039b03500197000000000423013f000000000023004b00000000020000190000039b02004041000b00000005001d000000000015004b00000000010000190000039b010020410000039b0040009c000000000201c019000000000002004b0000094b0000c13d0000000a010000290000000701100039000a00000001001d000000000101041a000000010010019000000001021002700000007f0220618f000700000002001d0000001f0020008c00000000020000390000000102002039000000000121013f0000000100100190000002ff0000c13d0000000701000029000000200010008c000009fb0000413d0000000a01000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d00000008030000290000001f023000390000000502200270000000200030008c0000000002004019000000000301043b00000007010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b000009fb0000813d000000000002041b0000000102200039000000000012004b000009f70000413d00000008010000290000001f0010008c0007000100100218000600030010021800000a110000a13d0000000a01000029000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f00000001002001900000094b0000613d000000200200008a0000000802200180000000000101043b00000a1d0000c13d000000000300001900000a270000013d000000080000006b000000000100001900000a170000613d0000000b010000290000000101100367000000000101043b000000010300008a0000000602300250000000000232013f000000000121016f00000007011001af00000a340000013d000000010400036700000000030000190000000b05300029000000000554034f000000000505043b000000000051041b00000001011000390000002003300039000000000023004b00000a1f0000413d000000080020006c00000a320000813d0000000602000029000000f80220018f000003c00220027f000003c0022001670000000b033000290000000103300367000000000303043b000000000223016f000000000021041b000000070100002900000001011001bf0000000a02000029000000000012041b000000400100043d0000002002100039000000600300003900000000003204350000000502000029000000000021043500000001020003670000000403200370000000000303043b0000000404300039000000000342034f000000000703043b000000000300003100000000054300490000001f0550008a0000039b065001970000039b08700197000000000968013f000000000068004b00000000080000190000039b08004041000000000057004b000000000a0000190000039b0a0080410000039b0090009c00000000080ac019000000000008004b0000094b0000c13d0000000007470019000000000872034f000000000908043b000003980090009c0000094b0000213d000000200a7000390000000607900210000000000773004900000000007a004b00000000080000190000039b080020410000039b077001970000039b0ba00197000000000c7b013f00000000007b004b00000000070000190000039b070040410000039b00c0009c000000000708c019000000000007004b0000094b0000c13d000000a0071000390000006008100039000000400b0000390000000000b804350000000000970435000000c007100039000000000009004b00000a7c0000613d000000000b000019000000000ca2034f000000000c0c043b000000000cc70436000000200da00039000000000dd2034f000000000d0d043b0000000000dc0435000000400aa000390000004007700039000000010bb0003900000000009b004b00000a700000413d0000002009400039000000000992034f000000000909043b0000039b0a900197000000000b6a013f00000000006a004b00000000060000190000039b06004041000000000059004b00000000050000190000039b050080410000039b00b0009c000000000605c019000000000006004b0000094b0000c13d0000000005490019000000000452034f000000000404043b000003980040009c0000094b0000213d00000020055000390000000006430049000000000065004b00000000090000190000039b090020410000039b066001970000039b0a500197000000000b6a013f00000000006a004b00000000060000190000039b060040410000039b00b0009c000000000609c019000000000006004b0000094b0000c13d000000000687004900000080081000390000000000680435000000000852034f0000000005470436000003bf074001980000001f0940018f000000000675001900000aae0000613d000000000a08034f000000000b05001900000000ac0a043c000000000bcb043600000000006b004b00000aaa0000c13d000000000009004b00000abb0000613d000000000778034f0000000308900210000000000906043300000000098901cf000000000989022f000000000707043b0000010008800089000000000787022f00000000078701cf000000000797019f0000000000760435000000000654001900000000000604350000001f04400039000003bf0440019700000000045400190000000005140049000000400610003900000000005604350000002405200370000000000505043b0000000405500039000000000652034f000000000806043b00000000065300490000001f0660008a0000039b076001970000039b09800197000000000a79013f000000000079004b00000000090000190000039b09004041000000000068004b000000000b0000190000039b0b0080410000039b00a0009c00000000090bc019000000000009004b0000094b0000c13d0000000008580019000000000982034f000000000a09043b0000039800a0009c0000094b0000213d000000200b800039000003c208a000d1000000000838001900000000008b004b00000000090000190000039b090020410000039b088001970000039b0cb00197000000000d8c013f00000000008c004b00000000080000190000039b080040410000039b00d0009c000000000809c019000000000008004b0000094b0000c13d0000004009400039000000400800003900000000088404360000000000a90435000000600940003900000000000a004b00000b090000613d000000000c000019000000000db2034f000000000d0d043b0000039800d0009c0000094b0000213d000000000ed90436000000200db00039000000000fd2034f000000000f0f043b0000039800f0009c0000094b0000213d0000000000fe0435000000200dd00039000000000dd2034f000000000d0d043b000000400e9000390000000000de0435000000600bb000390000006009900039000000010cc000390000000000ac004b00000af40000413d000000200a500039000000000aa2034f000000000a0a043b0000039b0ba00197000000000c7b013f00000000007b004b00000000070000190000039b0700404100000000006a004b00000000060000190000039b060080410000039b00c0009c000000000706c019000000000007004b0000094b0000c13d00000000065a0019000000000562034f000000000505043b000003980050009c0000094b0000213d00000020066000390000000003530049000000000036004b00000000070000190000039b070020410000039b033001970000039b0a600197000000000b3a013f00000000003a004b00000000030000190000039b030040410000039b00b0009c000000000307c019000000000003004b0000094b0000c13d00000000034900490000000000380435000000000462034f0000000002590436000003bf065001980000001f0750018f000000000362001900000b3a0000613d000000000804034f0000000009020019000000008a08043c0000000009a90436000000000039004b00000b360000c13d000000000007004b00000b470000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f0000000000430435000000000325001900000000000304350000001f03500039000003bf0330019700000000021200490000000002320019000003720020009c00000372020080410000006002200210000003720010009c00000372010080410000004001100210000000000112019f0000000002000414000003720020009c0000037202008041000000c002200210000000000112019f0000038c011001c70000800d020000390000000203000039000003b90400004100000009050000290dc40dba0000040f00000001002001900000094b0000613d000000400100043d00000009020000290000000000210435000003720010009c00000372010080410000004001100210000003a8011001c700000dc50001042e00000000430104340000000001320436000000000003004b00000b750000613d000000000200001900000000051200190000000006240019000000000606043300000000006504350000002002200039000000000032004b00000b6e0000413d000000000213001900000000000204350000001f02300039000003bf022001970000000001210019000000000001042d00000000430104340000037203300197000000000332043600000000040404330000000000430435000000400310003900000000030304330000004004200039000000800500003900000000005404350000008005200039000000400400003900000000360304340000000000450435000000c00520003900000000070604330000000000750435000000e005200039000000000007004b00000b9a0000613d00000000080000190000002006600039000000000906043300000000a90904340000000009950436000000000a0a04330000000000a9043500000040055000390000000108800039000000000078004b00000b900000413d0000000006250049000000800660008a0000000003030433000000a007200039000000000067043500000000630304340000000005350436000000000003004b00000bab0000613d000000000700001900000000085700190000000009760019000000000909043300000000009804350000002007700039000000000037004b00000ba40000413d000000000653001900000000000604350000001f06300039000003bf06600197000000000556001900000060011000390000000001010433000000000625004900000060022000390000000000620435000000002601043400000000014504360000004004500039000000000706043300000000007404350000006004500039000000000007004b00000bce0000613d00000000080000190000002006600039000000000906043300000000ba090434000003980aa00197000000000aa40436000000000b0b0433000003980bb001970000000000ba043500000040099000390000000009090433000000400a40003900000000009a043500000060044000390000000108800039000000000078004b00000bbe0000413d00000000020204330000000005540049000000000051043500000000520204340000000001240436000000000002004b00000bdd0000613d000000000400001900000000061400190000000007450019000000000707043300000000007604350000002004400039000000000024004b00000bd60000413d000000000412001900000000000404350000001f02200039000003bf022001970000000001120019000000000001042d00000000030100190000000001120049000003990010009c00000c7b0000213d0000003f0010008c00000c7b0000a13d000000400100043d000003c30010009c00000c7d0000813d0000004005100039000000400050043f0000000104000367000000000634034f000000000606043b000003980060009c00000c7b0000213d00000000073600190000001f06700039000000000026004b00000000080000190000039b080080410000039b096001970000039b06200197000000000a69013f000000000069004b00000000090000190000039b090040410000039b00a0009c000000000908c019000000000009004b00000c7b0000c13d000000000874034f000000000808043b000003980080009c00000c7d0000213d00000005098002100000003f099000390000039c099001970000000009590019000003980090009c00000c7d0000213d000000400090043f000000000085043500000060088000c900000020077000390000000008870019000000000028004b00000c7b0000213d000000000087004b00000c340000813d0000006009100039000000000a7200490000039900a0009c00000c7b0000213d0000006000a0008c00000c7b0000413d000000400a00043d0000039d00a0009c00000c7d0000213d000000600ba000390000004000b0043f000000000b74034f000000000b0b043b0000039800b0009c00000c7b0000213d000000000cba0436000000200b700039000000000db4034f000000000d0d043b0000039800d0009c00000c7b0000213d0000000000dc0435000000200bb00039000000000bb4034f000000000b0b043b000000400ca000390000000000bc04350000000009a904360000006007700039000000000087004b00000c160000413d00000000055104360000002007300039000000000774034f000000000707043b000003980070009c00000c7b0000213d00000000083700190000001f03800039000000000023004b00000000070000190000039b070080410000039b03300197000000000963013f000000000063004b00000000030000190000039b030040410000039b0090009c000000000307c019000000000003004b00000c7b0000c13d000000000384034f000000000303043b000003980030009c00000c7d0000213d0000001f06300039000003bf066001970000003f06600039000003bf07600197000000400600043d0000000007760019000000000067004b000000000a000039000000010a004039000003980070009c00000c7d0000213d0000000100a0019000000c7d0000c13d000000400070043f00000000073604360000002008800039000000000a83001900000000002a004b00000c7b0000213d000000000484034f000003bf083001980000001f0930018f000000000287001900000c6a0000613d000000000a04034f000000000b07001900000000ac0a043c000000000bcb043600000000002b004b00000c660000c13d000000000009004b00000c770000613d000000000484034f0000000308900210000000000902043300000000098901cf000000000989022f000000000404043b0000010008800089000000000484022f00000000048401cf000000000494019f0000000000420435000000000237001900000000000204350000000000650435000000000001042d000000000100001900000dc600010430000003ba01000041000000000010043f0000004101000039000000040010043f000003a60100004100000dc6000104300000000e01000039000000000101041a00000020011002700000037201100197000000020010008c00000c8e0000813d000000010110015f00000006011000c90000000301100039000000000101041a000000000001042d000003ba01000041000000000010043f0000003201000039000000040010043f000003a60100004100000dc6000104300007000000000002000000400200043d000700000002001d000003c40020009c00000dab0000813d00000007020000290000008003200039000000400030043f000000000301041a000003720330019700000000033204360000000102100039000000000202041a0000000000230435000000400400043d0000039a0040009c00000dab0000213d00000002031000390000004005400039000000400050043f000000000603041a000003980060009c00000dab0000213d00000005026002100000003f022000390000039c022001970000000002520019000003980020009c00000dab0000213d000500000004001d000600000001001d000000400020043f000400000005001d0000000000650435000000000030043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c70000801002000039000300000006001d0dc40dbf0000040f000000010020019000000db10000613d000000030a00002900000000000a004b00000006070000290000000508000029000000040900002900000cda0000613d000000000101043b00000000020000190000000003090019000000400400043d0000039a0040009c00000dab0000213d0000004005400039000000400050043f000000000501041a00000000055404360000000106100039000000000606041a000000000065043500000020033000390000000000430435000000020110003900000001022000390000000000a2004b00000cca0000413d00000000059804360000000301700039000000000201041a000000010320019000000001092002700000007f0990618f0000001f0090008c00000000040000390000000104002039000000000043004b00000db30000c13d000000400600043d0000000004960436000000000003004b00000d090000613d000100000004001d000200000009001d000300000006001d000400000005001d000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f000000010020019000000db10000613d0000000209000029000000000009004b00000d0f0000613d000000000201043b00000000010000190000000607000029000000050800002900000004050000290000000306000029000000010a00002900000000031a0019000000000402041a000000000043043500000001022000390000002001100039000000000091004b00000d010000413d00000d140000013d000003c1012001970000000000140435000000000009004b0000002001000039000000000100603900000d140000013d000000000100001900000006070000290000000508000029000000040500002900000003060000290000003f01100039000003bf021001970000000001620019000000000021004b00000000020000390000000102004039000003980010009c00000dab0000213d000000010020019000000dab0000c13d000000400010043f0000000000650435000000070100002900000040011000390000000000810435000000400300043d0000039a0030009c00000dab0000213d00000004017000390000004004300039000000400040043f000000000501041a000003980050009c00000dab0000213d00000005025002100000003f022000390000039c022001970000000002420019000003980020009c00000dab0000213d000500000003001d000000400020043f000400000004001d0000000000540435000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c70000801002000039000300000005001d0dc40dbf0000040f000000010020019000000db10000613d000000030b00002900000000000b004b00000006080000290000000509000029000000040a00002900000d5f0000613d0000006002900039000000000101043b0000000003000019000000400400043d0000039d0040009c00000dab0000213d0000006005400039000000400050043f000000000501041a0000004006500270000003980660019700000020074000390000000000670435000003980550019700000000005404350000000105100039000000000505041a000000400640003900000000005604350000000002420436000000020110003900000001033000390000000000b3004b00000d4a0000413d0000000006a904360000000501800039000000000201041a000000010320019000000001052002700000007f0550618f0000001f0050008c00000000040000390000000104002039000000000442013f000000010040019000000db30000c13d000000400700043d0000000004570436000000000003004b00000d8f0000613d000200000004001d000600000005001d000300000007001d000400000006001d000000000010043f0000000001000414000003720010009c0000037201008041000000c0011002100000039f011001c700008010020000390dc40dbf0000040f000000010020019000000db10000613d0000000608000029000000000008004b00000d960000613d000000000201043b0000000001000019000000200500008a000000050900002900000004060000290000000307000029000000020a00002900000000031a0019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00000d870000413d00000d9b0000013d000003c1012001970000000000140435000000000005004b00000020010000390000000001006039000000200500008a00000d9b0000013d0000000001000019000000200500008a0000000509000029000000040600002900000003070000290000003f01100039000000000251016f0000000001720019000000000021004b00000000020000390000000102004039000003980010009c00000dab0000213d000000010020019000000dab0000c13d000000400010043f0000000000760435000000070100002900000060021000390000000000920435000000000001042d000003ba01000041000000000010043f0000004101000039000000040010043f000003a60100004100000dc600010430000000000100001900000dc600010430000003ba01000041000000000010043f0000002201000039000000040010043f000003a60100004100000dc600010430000000000001042f00000dbd002104210000000102000039000000000001042d0000000002000019000000000001042d00000dc2002104230000000102000039000000000001042d0000000002000019000000000001042d00000dc40000043200000dc50001042e00000dc600010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000002000000000000000000000000000000400000010000000000000000009b15e16f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000000000000000000000000000000000000000000000000000000000006dd5b69c000000000000000000000000000000000000000000000000000000008c76967e00000000000000000000000000000000000000000000000000000000f2fde38a00000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000fb4022d4000000000000000000000000000000000000000000000000000000008c76967f000000000000000000000000000000000000000000000000000000008da5cb5b000000000000000000000000000000000000000000000000000000006dd5b69d00000000000000000000000000000000000000000000000000000000736be8020000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000003567e6b3000000000000000000000000000000000000000000000000000000003567e6b40000000000000000000000000000000000000000000000000000000038354c5c000000000000000000000000000000000000000000000000000000006350795600000000000000000000000000000000000000000000000000000000118dbac500000000000000000000000000000000000000000000000000000000123e65db00000000000000000000000000000000000000000000000000000000181f5a77000000000000000000000000ffffffffffffffffffffffffffffffffffffffff93df584c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000800000000000000000020000000000000000000000000000000000000000000000000000000000000053f5d9228f0a4173bea6e5931c9b3afe6eeb6692ede1d182952970f152534e3b0849d8cc00000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000080000000000000000000000000000000000000000000000000000000000000000000000001000000000b31c0055e2d464bef7781994b98c4ff9ef4ae0d05f59feb6a68c42de5e201b8fc3e98dbbd47c3fa7c1c05b6ec711caeaf70eca4554192b9ada8fc11a37f298e7b4d1e4f0000000000000000000000000000000000000000000000000000000002b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0000000000000000000000000000000000000000000000000ffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffbf80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffff9f00000000000000000000000000000000000000000000000100000000000000000200000000000000000000000000000000000020000000000000000000000000ffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff00000000000000001f69d1a2edb327babc986b3deb80091f101b9105d42a6c30db4d99c31d7e6294000000000000000000000000000000000000000000000000fffffffffffffffe000000000000000000000000000000000000000000000001fffffffffffffffed0b2c031000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000000000000000fffffffffffffdbf00000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000040000000800000000000000000524d4e486f6d6520312e362e302d6465760000000000000000000000000000000000000000000000000000000000000000000000000000c000000000000000002b5c74de00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff3fae00651d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000221a8ae8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000045564d00000000000000000000000000000000000000000000000000000000009a8a0592ac89c5ad3bc6df8224c17b485976f597df104ee20d0df415241f670b0200000200000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5f0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000b000000000000000000000000000000000000000000000000000000000000f6c6d1be15ba0acc8ee645c1ec613c360ef786d2d3200eb8e695b6dec757dbf04e487b71000000000000000000000000000000000000000000000000000000002847b60600000000000000000000000000000000000000000000000000000000a804bcb3000000000000000000000000000000000000000000000000000000003857f84d00000000000000000000000000000000000000000000000000000000af26d5e300000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0000000000000000000000000000000000000000000000000ffffffffffffffc0000000000000000000000000000000000000000000000000ffffffffffffff80")

func DeployRMNHomeZK(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *generated.Transaction,

	*RMNHome,
	error) {
	parsed,
		err := RMNHomeMetaData.
		GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed ==
		nil {
		return common.Address{}, nil, nil, errors.
			New("GetABI returned nil")
	}
	address,
		ethTx, contract, err := generated.DeployContract(auth, parsed, common.FromHex(RMNHomeZKBin), backend)
	if err != nil {
		return common.Address{}, nil,

			nil,
			err
	}
	return address,
		ethTx, &RMNHome{address: address,
			abi: *parsed, RMNHomeCaller: RMNHomeCaller{contract: contract}, RMNHomeTransactor: RMNHomeTransactor{contract: contract}, RMNHomeFilterer: RMNHomeFilterer{contract: contract}}, nil
}
