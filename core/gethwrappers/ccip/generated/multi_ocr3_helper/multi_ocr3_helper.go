package multi_ocr3_helper

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

type MultiOCR3BaseConfigInfo struct {
	ConfigDigest                   [32]byte
	F                              uint8
	N                              uint8
	IsSignatureVerificationEnabled bool
}

type MultiOCR3BaseOCRConfig struct {
	ConfigInfo   MultiOCR3BaseConfigInfo
	Signers      []common.Address
	Transmitters []common.Address
}

type MultiOCR3BaseOCRConfigArgs struct {
	ConfigDigest                   [32]byte
	OcrPluginType                  uint8
	F                              uint8
	IsSignatureVerificationEnabled bool
	Signers                        []common.Address
	Transmitters                   []common.Address
}

type MultiOCR3BaseOracle struct {
	Index uint8
	Role  uint8
}

var MultiOCR3HelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"AfterConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"}],\"name\":\"getOracle\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"index\",\"type\":\"uint8\"},{\"internalType\":\"enumMultiOCR3Base.Role\",\"name\":\"role\",\"type\":\"uint8\"}],\"internalType\":\"structMultiOCR3Base.Oracle\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"setTransmitOcrPluginType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmitWithSignatures\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"transmitWithoutSignatures\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b503360008161003257604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b0384811691909117909155811615610062576100628161006d565b5050466080526100e6565b336001600160a01b0382160361009657604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b608051611d256200010960003960008181610edd0152610f290152611d256000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80637ac0aa1a11610076578063c673e5841161005b578063c673e584146101c5578063f2fde38b146101e5578063f716f99f146101f857600080fd5b80637ac0aa1a1461015b5780638da5cb5b1461019d57600080fd5b806334a9c92e116100a757806334a9c92e1461012057806344e65e551461014057806379ba50971461015357600080fd5b8063181f5a77146100c357806326bf9d261461010b575b600080fd5b604080518082018252601981527f4d756c74694f4352334261736548656c70657220312e302e30000000000000006020820152905161010291906114ca565b60405180910390f35b61011e610119366004611591565b61020b565b005b61013361012e36600461161f565b61023a565b6040516101029190611681565b61011e61014e3660046116f4565b6102ca565b61011e61034d565b61011e6101693660046117a7565b600480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff92909216919091179055565b60015460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610102565b6101d86101d33660046117a7565b61041b565b604051610102919061181b565b61011e6101f33660046118ae565b610593565b61011e610206366004611a1a565b6105a7565b604080516000808252602082019092526004549091506102349060ff16858585858060006105e9565b50505050565b6040805180820182526000808252602080830182905260ff86811683526003825284832073ffffffffffffffffffffffffffffffffffffffff871684528252918490208451808601909552805480841686529394939092918401916101009091041660028111156102ad576102ad611652565b60028111156102be576102be611652565b90525090505b92915050565b60045460408051602080880282810182019093528782526103439360ff16928c928c928c928c918c91829185019084908082843760009201919091525050604080516020808d0282810182019093528c82529093508c92508b9182918501908490808284376000920191909152508a92506105e9915050565b5050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461039e576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61045e6040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c08201529485529182018054845181840281018401909552808552929385830193909283018282801561051457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116104e9575b505050505081526020016003820180548060200260200160405190810160405280929190818152602001828054801561058357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610558575b5050505050815250509050919050565b61059b610972565b6105a4816109c5565b50565b6105af610972565b60005b81518110156105e5576105dd8282815181106105d0576105d0611b83565b6020026020010151610a89565b6001016105b2565b5050565b60ff878116600090815260026020908152604080832081516080810183528154815260019091015480861693820193909352610100830485169181019190915262010000909104909216151560608301528735906106488760a4611be1565b9050826060015115610690578451610661906020611bf4565b865161066e906020611bf4565b6106799060a0611be1565b6106839190611be1565b61068d9082611be1565b90505b3681146106d7576040517f8e1192e1000000000000000000000000000000000000000000000000000000008152600481018290523660248201526044015b60405180910390fd5b508151811461071f5781516040517f93df584c0000000000000000000000000000000000000000000000000000000081526004810191909152602481018290526044016106ce565b610727610eda565b60ff808a166000908152600360209081526040808320338452825280832081518083019092528054808616835293949193909284019161010090910416600281111561077557610775611652565b600281111561078657610786611652565b90525090506002816020015160028111156107a3576107a3611652565b1480156108045750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff16815481106107df576107df611b83565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b61083a576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5081606001511561091c576020820151610855906001611c0b565b60ff16855114610891576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83518551146108cc576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600087876040516108de929190611c24565b6040519081900381206108f5918b90602001611c34565b60405160208183030381529060405280519060200120905061091a8a82888888610f5b565b505b6040805182815260208a81013567ffffffffffffffff169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146109c3576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b3373ffffffffffffffffffffffffffffffffffffffff821603610a14576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff16600003610acd5760006040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ce9190611c48565b60208082015160ff80821660009081526002909352604083206001810154929390928392169003610b3a57606084015160018201805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055610b8f565b6060840151600182015460ff6201000090910416151590151514610b8f576040517f87f6037c00000000000000000000000000000000000000000000000000000000815260ff841660048201526024016106ce565b60a084015180516101001015610bd45760016040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ce9190611c48565b8051600003610c125760056040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ce9190611c48565b610c858484600301805480602002602001604051908101604052809291908181526020018280548015610c7b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610c50575b505050505061116b565b846060015115610e2a57610d008484600201805480602002602001604051908101604052809291908181526020018280548015610c7b5760200282019190600052602060002090815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610c5057505050505061116b565b608085015180516101001015610d455760026040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ce9190611c48565b6040860151610d55906003611c62565b60ff16815111610d945760036040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ce9190611c48565b815181511015610dd35760016040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ce9190611c48565b80516001840180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff841602179055610e1b906002860190602084019061142b565b50610e2885826001611203565b505b610e3684826002611203565b8051610e4b906003850190602084019061142b565b506040858101516001840180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff8316179055865180855560a088015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f54793610ec29389939260028a01929190611c85565b60405180910390a1610ed3846113f2565b5050505050565b467f0000000000000000000000000000000000000000000000000000000000000000146109c3576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000060048201524660248201526044016106ce565b8251600090815b81811015610343576000600188868460208110610f8157610f81611b83565b610f8e91901a601b611c0b565b898581518110610fa057610fa0611b83565b6020026020010151898681518110610fba57610fba611b83565b602002602001015160405160008152602001604052604051610ff8949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561101a573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015160ff808e1660009081526003602090815285822073ffffffffffffffffffffffffffffffffffffffff8516835281528582208587019096528554808416865293975090955092939284019161010090041660028111156110a6576110a6611652565b60028111156110b7576110b7611652565b90525090506001816020015160028111156110d4576110d4611652565b1461110b576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600160ff9091161b85161561114e576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000015160ff166001901b851794505050806001019050610f62565b60005b81518110156111fe5760ff8316600090815260036020526040812083519091908490849081106111a0576111a0611b83565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905560010161116e565b505050565b60005b825181101561023457600083828151811061122357611223611b83565b602002602001015190506000600281111561124057611240611652565b60ff808716600090815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152902054610100900416600281111561128c5761128c611652565b146112c65760046040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ce9190611c48565b73ffffffffffffffffffffffffffffffffffffffff8116611313576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052808360ff16815260200184600281111561133957611339611652565b905260ff808716600090815260036020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845282529091208351815493167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841681178255918401519092909183917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016176101008360028111156113de576113de611652565b021790555090505050806001019050611206565b60405160ff821681527f897ac1b2c12867721b284f3eb147bd4ab046d4eef1cf31c1d8988bfcfb962b539060200160405180910390a150565b8280548282559060005260206000209081019282156114a5579160200282015b828111156114a557825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90911617825560209092019160019091019061144b565b506114b19291506114b5565b5090565b5b808211156114b157600081556001016114b6565b60006020808352835180602085015260005b818110156114f8578581018301518582016040015282016114dc565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b80606081018310156102c457600080fd5b60008083601f84011261155a57600080fd5b50813567ffffffffffffffff81111561157257600080fd5b60208301915083602082850101111561158a57600080fd5b9250929050565b6000806000608084860312156115a657600080fd5b6115b08585611537565b9250606084013567ffffffffffffffff8111156115cc57600080fd5b6115d886828701611548565b9497909650939450505050565b803560ff811681146115f657600080fd5b919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146115f657600080fd5b6000806040838503121561163257600080fd5b61163b836115e5565b9150611649602084016115fb565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b815160ff16815260208201516040820190600381106116a2576116a2611652565b8060208401525092915050565b60008083601f8401126116c157600080fd5b50813567ffffffffffffffff8111156116d957600080fd5b6020830191508360208260051b850101111561158a57600080fd5b60008060008060008060008060e0898b03121561171057600080fd5b61171a8a8a611537565b9750606089013567ffffffffffffffff8082111561173757600080fd5b6117438c838d01611548565b909950975060808b013591508082111561175c57600080fd5b6117688c838d016116af565b909750955060a08b013591508082111561178157600080fd5b5061178e8b828c016116af565b999c989b50969995989497949560c00135949350505050565b6000602082840312156117b957600080fd5b6117c2826115e5565b9392505050565b60008151808452602080850194506020840160005b8381101561181057815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016117de565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a084015261186a60e08401826117c9565b905060408401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160c08501526118a582826117c9565b95945050505050565b6000602082840312156118c057600080fd5b6117c2826115fb565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff8111828210171561191b5761191b6118c9565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611968576119686118c9565b604052919050565b600067ffffffffffffffff82111561198a5761198a6118c9565b5060051b60200190565b803580151581146115f657600080fd5b600082601f8301126119b557600080fd5b813560206119ca6119c583611970565b611921565b8083825260208201915060208460051b8701019350868411156119ec57600080fd5b602086015b84811015611a0f57611a02816115fb565b83529183019183016119f1565b509695505050505050565b60006020808385031215611a2d57600080fd5b823567ffffffffffffffff80821115611a4557600080fd5b818501915085601f830112611a5957600080fd5b8135611a676119c582611970565b81815260059190911b83018401908481019088831115611a8657600080fd5b8585015b83811015611b7657803585811115611aa157600080fd5b860160c0818c037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0011215611ad65760008081fd5b611ade6118f8565b8882013581526040611af18184016115e5565b8a8301526060611b028185016115e5565b8284015260809150611b15828501611994565b9083015260a08381013589811115611b2d5760008081fd5b611b3b8f8d838801016119a4565b838501525060c0840135915088821115611b555760008081fd5b611b638e8c848701016119a4565b9083015250845250918601918601611a8a565b5098975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156102c4576102c4611bb2565b80820281158282048414176102c4576102c4611bb2565b60ff81811683821601908111156102c4576102c4611bb2565b8183823760009101908152919050565b828152606082602083013760800192915050565b6020810160068310611c5c57611c5c611652565b91905290565b60ff8181168382160290811690818114611c7e57611c7e611bb2565b5092915050565b600060a0820160ff88168352602087602085015260a0604085015281875480845260c086019150886000526020600020935060005b81811015611cec57845473ffffffffffffffffffffffffffffffffffffffff1683526001948501949284019201611cba565b50508481036060860152611d0081886117c9565b935050505060ff83166080830152969550505050505056fea164736f6c6343000818000a",
}

var MultiOCR3HelperABI = MultiOCR3HelperMetaData.ABI

var MultiOCR3HelperBin = MultiOCR3HelperMetaData.Bin

func DeployMultiOCR3Helper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MultiOCR3Helper, error) {
	parsed, err := MultiOCR3HelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultiOCR3HelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiOCR3Helper{address: address, abi: *parsed, MultiOCR3HelperCaller: MultiOCR3HelperCaller{contract: contract}, MultiOCR3HelperTransactor: MultiOCR3HelperTransactor{contract: contract}, MultiOCR3HelperFilterer: MultiOCR3HelperFilterer{contract: contract}}, nil
}

type MultiOCR3Helper struct {
	address common.Address
	abi     abi.ABI
	MultiOCR3HelperCaller
	MultiOCR3HelperTransactor
	MultiOCR3HelperFilterer
}

type MultiOCR3HelperCaller struct {
	contract *bind.BoundContract
}

type MultiOCR3HelperTransactor struct {
	contract *bind.BoundContract
}

type MultiOCR3HelperFilterer struct {
	contract *bind.BoundContract
}

type MultiOCR3HelperSession struct {
	Contract     *MultiOCR3Helper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MultiOCR3HelperCallerSession struct {
	Contract *MultiOCR3HelperCaller
	CallOpts bind.CallOpts
}

type MultiOCR3HelperTransactorSession struct {
	Contract     *MultiOCR3HelperTransactor
	TransactOpts bind.TransactOpts
}

type MultiOCR3HelperRaw struct {
	Contract *MultiOCR3Helper
}

type MultiOCR3HelperCallerRaw struct {
	Contract *MultiOCR3HelperCaller
}

type MultiOCR3HelperTransactorRaw struct {
	Contract *MultiOCR3HelperTransactor
}

func NewMultiOCR3Helper(address common.Address, backend bind.ContractBackend) (*MultiOCR3Helper, error) {
	abi, err := abi.JSON(strings.NewReader(MultiOCR3HelperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMultiOCR3Helper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3Helper{address: address, abi: abi, MultiOCR3HelperCaller: MultiOCR3HelperCaller{contract: contract}, MultiOCR3HelperTransactor: MultiOCR3HelperTransactor{contract: contract}, MultiOCR3HelperFilterer: MultiOCR3HelperFilterer{contract: contract}}, nil
}

func NewMultiOCR3HelperCaller(address common.Address, caller bind.ContractCaller) (*MultiOCR3HelperCaller, error) {
	contract, err := bindMultiOCR3Helper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperCaller{contract: contract}, nil
}

func NewMultiOCR3HelperTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiOCR3HelperTransactor, error) {
	contract, err := bindMultiOCR3Helper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperTransactor{contract: contract}, nil
}

func NewMultiOCR3HelperFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiOCR3HelperFilterer, error) {
	contract, err := bindMultiOCR3Helper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperFilterer{contract: contract}, nil
}

func bindMultiOCR3Helper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiOCR3HelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MultiOCR3Helper *MultiOCR3HelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiOCR3Helper.Contract.MultiOCR3HelperCaller.contract.Call(opts, result, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.MultiOCR3HelperTransactor.contract.Transfer(opts)
}

func (_MultiOCR3Helper *MultiOCR3HelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.MultiOCR3HelperTransactor.contract.Transact(opts, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiOCR3Helper.Contract.contract.Call(opts, result, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.contract.Transfer(opts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.contract.Transact(opts, method, params...)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) GetOracle(opts *bind.CallOpts, ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "getOracle", ocrPluginType, oracleAddress)

	if err != nil {
		return *new(MultiOCR3BaseOracle), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOracle)).(*MultiOCR3BaseOracle)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) GetOracle(ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error) {
	return _MultiOCR3Helper.Contract.GetOracle(&_MultiOCR3Helper.CallOpts, ocrPluginType, oracleAddress)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) GetOracle(ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error) {
	return _MultiOCR3Helper.Contract.GetOracle(&_MultiOCR3Helper.CallOpts, ocrPluginType, oracleAddress)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "latestConfigDetails", ocrPluginType)

	if err != nil {
		return *new(MultiOCR3BaseOCRConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOCRConfig)).(*MultiOCR3BaseOCRConfig)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _MultiOCR3Helper.Contract.LatestConfigDetails(&_MultiOCR3Helper.CallOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _MultiOCR3Helper.Contract.LatestConfigDetails(&_MultiOCR3Helper.CallOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) Owner() (common.Address, error) {
	return _MultiOCR3Helper.Contract.Owner(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) Owner() (common.Address, error) {
	return _MultiOCR3Helper.Contract.Owner(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MultiOCR3Helper.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TypeAndVersion() (string, error) {
	return _MultiOCR3Helper.Contract.TypeAndVersion(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperCallerSession) TypeAndVersion() (string, error) {
	return _MultiOCR3Helper.Contract.TypeAndVersion(&_MultiOCR3Helper.CallOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "acceptOwnership")
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.AcceptOwnership(&_MultiOCR3Helper.TransactOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.AcceptOwnership(&_MultiOCR3Helper.TransactOpts)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "setOCR3Configs", ocrConfigArgs)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetOCR3Configs(&_MultiOCR3Helper.TransactOpts, ocrConfigArgs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetOCR3Configs(&_MultiOCR3Helper.TransactOpts, ocrConfigArgs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) SetTransmitOcrPluginType(opts *bind.TransactOpts, ocrPluginType uint8) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "setTransmitOcrPluginType", ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) SetTransmitOcrPluginType(ocrPluginType uint8) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetTransmitOcrPluginType(&_MultiOCR3Helper.TransactOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) SetTransmitOcrPluginType(ocrPluginType uint8) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.SetTransmitOcrPluginType(&_MultiOCR3Helper.TransactOpts, ocrPluginType)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transferOwnership", to)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransferOwnership(&_MultiOCR3Helper.TransactOpts, to)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransferOwnership(&_MultiOCR3Helper.TransactOpts, to)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransmitWithSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transmitWithSignatures", reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransmitWithSignatures(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransmitWithSignatures(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactor) TransmitWithoutSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.contract.Transact(opts, "transmitWithoutSignatures", reportContext, report)
}

func (_MultiOCR3Helper *MultiOCR3HelperSession) TransmitWithoutSignatures(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithoutSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report)
}

func (_MultiOCR3Helper *MultiOCR3HelperTransactorSession) TransmitWithoutSignatures(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _MultiOCR3Helper.Contract.TransmitWithoutSignatures(&_MultiOCR3Helper.TransactOpts, reportContext, report)
}

type MultiOCR3HelperAfterConfigSetIterator struct {
	Event *MultiOCR3HelperAfterConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperAfterConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperAfterConfigSet)
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
		it.Event = new(MultiOCR3HelperAfterConfigSet)
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

func (it *MultiOCR3HelperAfterConfigSetIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperAfterConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperAfterConfigSet struct {
	OcrPluginType uint8
	Raw           types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterAfterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperAfterConfigSetIterator, error) {

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "AfterConfigSet")
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperAfterConfigSetIterator{contract: _MultiOCR3Helper.contract, event: "AfterConfigSet", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchAfterConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperAfterConfigSet) (event.Subscription, error) {

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "AfterConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperAfterConfigSet)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "AfterConfigSet", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseAfterConfigSet(log types.Log) (*MultiOCR3HelperAfterConfigSet, error) {
	event := new(MultiOCR3HelperAfterConfigSet)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "AfterConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperConfigSetIterator struct {
	Event *MultiOCR3HelperConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperConfigSet)
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
		it.Event = new(MultiOCR3HelperConfigSet)
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

func (it *MultiOCR3HelperConfigSetIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperConfigSet struct {
	OcrPluginType uint8
	ConfigDigest  [32]byte
	Signers       []common.Address
	Transmitters  []common.Address
	F             uint8
	Raw           types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperConfigSetIterator, error) {

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperConfigSetIterator{contract: _MultiOCR3Helper.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperConfigSet) (event.Subscription, error) {

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperConfigSet)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseConfigSet(log types.Log) (*MultiOCR3HelperConfigSet, error) {
	event := new(MultiOCR3HelperConfigSet)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperOwnershipTransferRequestedIterator struct {
	Event *MultiOCR3HelperOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperOwnershipTransferRequested)
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
		it.Event = new(MultiOCR3HelperOwnershipTransferRequested)
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

func (it *MultiOCR3HelperOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperOwnershipTransferRequestedIterator{contract: _MultiOCR3Helper.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperOwnershipTransferRequested)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseOwnershipTransferRequested(log types.Log) (*MultiOCR3HelperOwnershipTransferRequested, error) {
	event := new(MultiOCR3HelperOwnershipTransferRequested)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperOwnershipTransferredIterator struct {
	Event *MultiOCR3HelperOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperOwnershipTransferred)
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
		it.Event = new(MultiOCR3HelperOwnershipTransferred)
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

func (it *MultiOCR3HelperOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperOwnershipTransferredIterator{contract: _MultiOCR3Helper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperOwnershipTransferred)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseOwnershipTransferred(log types.Log) (*MultiOCR3HelperOwnershipTransferred, error) {
	event := new(MultiOCR3HelperOwnershipTransferred)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MultiOCR3HelperTransmittedIterator struct {
	Event *MultiOCR3HelperTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MultiOCR3HelperTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiOCR3HelperTransmitted)
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
		it.Event = new(MultiOCR3HelperTransmitted)
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

func (it *MultiOCR3HelperTransmittedIterator) Error() error {
	return it.fail
}

func (it *MultiOCR3HelperTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MultiOCR3HelperTransmitted struct {
	OcrPluginType  uint8
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*MultiOCR3HelperTransmittedIterator, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.FilterLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return &MultiOCR3HelperTransmittedIterator{contract: _MultiOCR3Helper.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperTransmitted, ocrPluginType []uint8) (event.Subscription, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _MultiOCR3Helper.contract.WatchLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MultiOCR3HelperTransmitted)
				if err := _MultiOCR3Helper.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_MultiOCR3Helper *MultiOCR3HelperFilterer) ParseTransmitted(log types.Log) (*MultiOCR3HelperTransmitted, error) {
	event := new(MultiOCR3HelperTransmitted)
	if err := _MultiOCR3Helper.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MultiOCR3Helper *MultiOCR3Helper) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MultiOCR3Helper.abi.Events["AfterConfigSet"].ID:
		return _MultiOCR3Helper.ParseAfterConfigSet(log)
	case _MultiOCR3Helper.abi.Events["ConfigSet"].ID:
		return _MultiOCR3Helper.ParseConfigSet(log)
	case _MultiOCR3Helper.abi.Events["OwnershipTransferRequested"].ID:
		return _MultiOCR3Helper.ParseOwnershipTransferRequested(log)
	case _MultiOCR3Helper.abi.Events["OwnershipTransferred"].ID:
		return _MultiOCR3Helper.ParseOwnershipTransferred(log)
	case _MultiOCR3Helper.abi.Events["Transmitted"].ID:
		return _MultiOCR3Helper.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MultiOCR3HelperAfterConfigSet) Topic() common.Hash {
	return common.HexToHash("0x897ac1b2c12867721b284f3eb147bd4ab046d4eef1cf31c1d8988bfcfb962b53")
}

func (MultiOCR3HelperConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (MultiOCR3HelperOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MultiOCR3HelperOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (MultiOCR3HelperTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_MultiOCR3Helper *MultiOCR3Helper) Address() common.Address {
	return _MultiOCR3Helper.address
}

type MultiOCR3HelperInterface interface {
	GetOracle(opts *bind.CallOpts, ocrPluginType uint8, oracleAddress common.Address) (MultiOCR3BaseOracle, error)

	LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error)

	SetTransmitOcrPluginType(opts *bind.TransactOpts, ocrPluginType uint8) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransmitWithSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	TransmitWithoutSignatures(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error)

	FilterAfterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperAfterConfigSetIterator, error)

	WatchAfterConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperAfterConfigSet) (event.Subscription, error)

	ParseAfterConfigSet(log types.Log) (*MultiOCR3HelperAfterConfigSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*MultiOCR3HelperConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*MultiOCR3HelperConfigSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MultiOCR3HelperOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MultiOCR3HelperOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MultiOCR3HelperOwnershipTransferred, error)

	FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*MultiOCR3HelperTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *MultiOCR3HelperTransmitted, ocrPluginType []uint8) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*MultiOCR3HelperTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var MultiOCR3HelperZKBin = ("0x0002000000000002001400000000000200010000000103550000006003100270000002a30030019d0000000100200190000000350000c13d000002a3023001970000008003000039000000400030043f000000040020008c000002e60000413d000000000401043b000000e004400270000002a90040009c000000410000213d000002b00040009c0000007a0000a13d000002b10040009c000001950000613d000002b20040009c000002490000613d000002b30040009c000002e60000c13d0000000001000416000000000001004b000002e60000c13d000000000100041a000002bc021001970000000006000411000000000026004b000002e80000c13d0000000102000039000000000302041a000002a404300197000000000464019f000000000042041b000002a401100197000000000010041b0000000001000414000002bc05300197000002a30010009c000002a301008041000000c001100210000002ca011001c70000800d020000390000000303000039000002d2040000410a860a7c0000040f0000000100200190000002e60000613d000000000100001900000a870001042e000000a001000039000000400010043f0000000001000416000000000001004b000002e60000c13d0000000001000411000000000001004b000000650000c13d000002a701000041000000a00010043f000002a80100004100000a8800010430000002aa0040009c000001680000a13d000002ab0040009c000001da0000613d000002ac0040009c000002dd0000613d000002ad0040009c000002e60000c13d0011002400200094000002e60000413d0000000003000416000000000003004b000002e60000c13d0000000403100370000000000403043b000002b60040009c000002e60000213d0000002303400039000000000023004b000002e60000813d0000000403400039000000000331034f000000000703043b000002b60070009c0000005f0000213d00000005067002100000003f03600039000002b705300197000002b80050009c0000030a0000a13d000002e601000041000000000010043f0000004101000039000000040010043f000002c20100004100000a88000104300000000103000039000000000203041a000002a402200197000000000112019f000000000013041b0000800b0100003900000004030000390000000004000415000000140440008a0000000504400210000002a5020000410a860a5e0000040f000000800010043f000001400000044300000160001004430000002001000039000001000010044300000001010000390000012000100443000002a60100004100000a870001042e000002b40040009c000001750000613d000002b50040009c000002e60000c13d000000840020008c000002e60000413d0000000003000416000000000003004b000002e60000c13d0000006403100370000000000303043b000002b60030009c000002e60000213d0000002304300039000000000024004b000002e60000813d001000040030003d0000001001100360000000000101043b001100000001001d000002b60010009c000002e60000213d0000001101300029000f00240010003d0000000f0020006b000002e60000213d000000a001000039000e00000001001d000000400010043f000000800000043f0000000401000039000000000101041a000000ff0110018f000d00000001001d000000000010043f0000000201000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000400200043d000002e10020009c0000005f0000813d000000000301043b0000008001200039000000400010043f000000000103041a00000000061204360000000103300039000000000303041a0000000804300270000000ff0440018f00000040052000390000000000450435000000ff0430018f00000000004604350000006007200039000002c0003001980000000002000039000000010200c03900000000002704350000001102000029000000a402200039000000d30000613d000000800400043d0000000503400210000000000004004b000000ca0000613d000002e90030009c000002d70000213d00000000044300d9000000200040008c000002d70000c13d000e00a00030003d0000000e043000290000000002240019000000000042004b000000000400003900000001040040390000000e0030002a000002d70000413d0000000100400190000002d70000c13d0000000003000031000000000023004b0000075c0000c13d000e00000007001d000c00000006001d00000004020000390000000102200367000000000202043b000000000021004b000007680000c13d000002d60100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000002a30010009c000002a301008041000000c001100210000002d7011001c700008005020000390a860a810000040f0000000100200190000007740000613d000000000101043b000b00000001001d000002a50100004100000000001004430000000001000414000002a30010009c000002a301008041000000c001100210000002d8011001c70000800b020000390a860a810000040f0000000100200190000007740000613d000000000101043b0000000b0010006b000006990000c13d0000000d01000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000000002000411000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000400200043d000b00000002001d000002c70020009c0000005f0000213d000000000101043b0000000b030000290000004002300039000000400020043f000000000201041a000000ff0120018f00000000011304360000000802200270000000ff0220018f000000030020008c000001d40000813d0000000000210435000000020020008c000007750000c13d0000000d01000029000000000010043f0000000201000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000000b020000290000000002020433000b00ff00200193000000000101043b0000000301100039000000000201041a0000000b0020006b0000099a0000813d000000000010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000000b01100029000000000101041a000002bc011001970000000002000411000000000012004b000007750000c13d0000000e010000290000000001010433000000000001004b000007b40000c13d00000001010003670000000402100370000000000202043b000000400300043d00000000022304360000002401100370000000000101043b000002b6011001970000000000120435000002a30030009c000002a30300804100000040013002100000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002bf011001c70000800d020000390000000203000039000002df040000410000000d05000029000000300000013d000002ae0040009c000001850000613d000002af0040009c000002e60000c13d0000000001000416000000000001004b000002e60000c13d0000000101000039000000000101041a000002bc01100197000000800010043f000002d00100004100000a870001042e0000000001000416000000000001004b000002e60000c13d000000c001000039000000400010043f0000001901000039000000800010043f000002e702000041000000a00020043f0000002003000039000000c00030043f000000e00010043f000001000020043f000001190000043f000002e80100004100000a870001042e000000240020008c000002e60000413d0000000002000416000000000002004b000002e60000c13d0000000401100370000000000101043b000000ff0010008c000002e60000213d0000000402000039000000000302041a000002ea03300197000000000113019f000000000012041b000000000100001900000a870001042e000000440020008c000002e60000413d0000000002000416000000000002004b000002e60000c13d0000000402100370000000000202043b000000ff0020008c000002e60000213d0000002401100370000000000101043b001100000001001d000002bc0010009c000002e60000213d000000c001000039000000400010043f000000800000043f000000a00000043f000000000020043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000001102000029000002bc02200197000000000101043b000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000400200043d000002c70020009c0000005f0000213d000000000101043b0000004003200039000000400030043f000000000101041a000000ff0310018f00000000023204360000000801100270000000ff0110018f000000020010008c000001d40000213d0000000000120435000000400100043d00000000033104360000000002020433000000020020008c000006930000a13d000002e601000041000000000010043f0000002101000039000000040010043f000002c20100004100000a8800010430000000240020008c000002e60000413d0000000002000416000000000002004b000002e60000c13d0000000401100370000000000101043b001100000001001d000000ff0010008c000002e60000213d000000e001000039000000400010043f0a8609ee0000040f000000e00000043f000001000000043f000001200000043f000001400000043f000000e001000039000000800010043f0000006001000039000000a00010043f000000c00010043f0000001101000029000000000010043f0000000201000039000000200010043f000000400200003900000000010000190a860a490000040f001100000001001d000000400100043d000f00000001001d0a8609e30000040f00000011010000290a860a0b0000040f0000000f020000290000000001120436001000000001001d000000400200043d000e00000002001d000000110100002900000002011000390a860a270000040f0000000e0210006a0000000e010000290a8609f90000040f00000010010000290000000e02000029000000000021043500000011010000290000000301100039000000400200043d001100000002001d0a860a270000040f000000110210006a00000011010000290a8609f90000040f0000000f020000290000004001200039000e00000001001d000000110300002900000000003104350000002001000039000000400500043d001100000005001d00000000011504360000000002020433000000004302043400000000003104350000000001040433000000ff0110018f0000004003500039000000000013043500000040012000390000000001010433000000ff0110018f0000006003500039000000000013043500000060012000390000000001010433000000000001004b0000000001000039000000010100c0390000008002500039000000000012043500000010010000290000000001010433000000c002000039000000a0035000390000000000230435000000e0025000390a8609d50000040f000000000201001900000011040000290000000001410049000000200310008a0000000e010000290000000001010433000000c00440003900000000003404350a8609d50000040f00000011020000290000000001210049000002a30020009c000002a3020080410000004002200210000002a30010009c000002a3010080410000006001100210000000000121019f00000a870001042e000000e40020008c000002e60000413d0000000004000416000000000004004b000002e60000c13d0000006404100370000000000404043b000002b60040009c000002e60000213d0000002305400039000000000025004b000002e60000813d0000000405400039000000000551034f000000000505043b001100000005001d000002b60050009c000002e60000213d0000002405400039001000000005001d000f00110050002d0000000f0020006b000002e60000213d0000008404100370000000000404043b000002b60040009c000002e60000213d0000002305400039000000000025004b000002e60000813d0000000405400039000000000551034f000000000905043b000002b60090009c000002e60000213d0000002407400039000000050b90021000000000087b0019000000000028004b000002e60000213d000000a404100370000000000404043b000002b60040009c000002e60000213d0000002305400039000000000025004b000002e60000813d0000000405400039000000000551034f000000000605043b000002b60060009c000002e60000213d0000002404400039000000050a60021000000000054a0019000000000025004b000002e60000213d0000003f02b00039000002b70b200197000002b800b0009c0000005f0000213d0000000402000039000000000202041a000000800bb00039000e0000000b001d0000004000b0043f000000800090043f000000000009004b000002970000613d000000000971034f000000000909043b000000200330003900000000009304350000002007700039000000000087004b0000028e0000413d000000400300043d000e00000003001d0000003f03a00039000002b7033001970000000e033000290000000e0030006c00000000070000390000000107004039000002b60030009c0000005f0000213d00000001007001900000005f0000c13d000d00ff00200193000000400030043f0000000e020000290000000002620436000c00000002001d000000000006004b000002b00000613d0000000e02000029000000000341034f000000000303043b000000200220003900000000003204350000002004400039000000000054004b000002a90000413d0000000d01000029000000000010043f0000000201000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000400200043d000002b80020009c0000005f0000213d000000000301043b0000008001200039000000400010043f000000000103041a00000000061204360000000103300039000000000303041a0000000804300270000000ff0440018f00000040052000390000000000450435000000ff0430018f000a00000006001d00000000004604350000006004200039000002c0033001980000000002000039000000010200c039000b00000004001d0000000000240435000000a50200008a000000110020006b000006b40000a13d000002e601000041000000000010043f0000001101000039000000040010043f000002c20100004100000a8800010430000000240020008c000002e60000413d0000000002000416000000000002004b000002e60000c13d0000000401100370000000000101043b000002bc0010009c000002ec0000a13d000000000100001900000a8800010430000002d101000041000000800010043f000002cd0100004100000a88000104300000000102000039000000000202041a000002bc022001970000000005000411000000000025004b000002f90000c13d000002bc06100197000000000056004b000002fd0000c13d000002cf01000041000000800010043f000002cd0100004100000a8800010430000002bd01000041000000800010043f000002cd0100004100000a8800010430000000000100041a000002a401100197000000000161019f000000000010041b0000000001000414000002a30010009c000002a301008041000000c001100210000002ca011001c70000800d020000390000000303000039000002ce04000041000000300000013d0000008003500039000000400030043f000000800070043f00000024054000390000000006560019000000000026004b000002e60000213d000000000007004b000005ec0000c13d0000000101000039000000000101041a000002bc011001970000000002000411000000000012004b0000068b0000c13d000000800100043d000000000001004b000000330000613d000600000000001d00000006010000290000000501100210000000a00110003900000000020104330000004001200039000400000001001d0000000001010433000000ff00100190000007780000613d000700000002001d00000020012000390000000001010433000000ff0110018f001100000001001d000000000010043f0000000201000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000070300002900000060073000390000000002070433000000000401043b0000000105400039000000000105041a000000ff00100190000003480000613d000000000002004b0000000002000039000000010200c039000002c0001001980000000001000039000000010100c039000000000021004b0000034e0000613d0000078d0000013d000002c301100197000000000002004b000002c4020000410000000002006019000000000112019f000000000015041b000000400100043d000f00000001001d000000a00230003900000000030204330000000061030434000b00000006001d000001000010008c000007a10000213d000d00000007001d000c00000003001d000100000002001d000300000005001d000000000001004b0000077e0000613d000500000004001d0000000301400039000000000301041a0000000f02000029001000000003001d0000000002320436000e00000002001d000200000001001d000000000010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000001005000029000000000005004b0000000e020000290000037c0000613d000000000101043b0000000e020000290000000003000019000000000401041a000002bc04400197000000000242043600000001011000390000000103300039000000000053004b000003750000413d0000000f0120006a0000001f01100039000002eb021001970000000f01200029000000000021004b00000000020000390000000102004039000002b60010009c0000005f0000213d00000001002001900000005f0000c13d000000400010043f0000000f010000290000000001010433000000000001004b000003b90000613d0000000001000019001000000001001d0000001101000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000000f0200002900000000020204330000001003000029000000000032004b0000099a0000a13d00000005023002100000000e022000290000000002020433000002bc02200197000000000101043b000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b000000000001041b000000100200002900000001022000390000000f010000290000000001010433000000000012004b00000000010200190000038d0000413d0000000d010000290000000001010433000000000001004b0000000501000029000800020010003d000004d40000613d0000000801000029000000000301041a000000400200043d000f00000002001d001000000003001d0000000002320436000e00000002001d000000000010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000001005000029000000000005004b0000000e02000029000003de0000613d000000000101043b0000000e020000290000000003000019000000000401041a000002bc04400197000000000242043600000001011000390000000103300039000000000053004b000003d70000413d0000000f0120006a0000001f01100039000002eb021001970000000f01200029000000000021004b00000000020000390000000102004039000002b60010009c0000005f0000213d00000001002001900000005f0000c13d000000400010043f0000000f010000290000000001010433000000000001004b0000041b0000613d0000000001000019001000000001001d0000001101000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000000f0200002900000000020204330000001003000029000000000032004b0000099a0000a13d00000005023002100000000e022000290000000002020433000002bc02200197000000000101043b000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b000000000001041b000000100200002900000001022000390000000f010000290000000001010433000000000012004b0000000001020019000003ef0000413d000000070100002900000080011000390000000001010433000a00000001001d0000000014010434000900000001001d000001000040008c000007930000213d00000004010000290000000001010433000000fe0210018f000000550020008c00000003030000290000000c02000029000002d70000213d00000003011000c9000000ff0110018f000000000014004b000007990000a13d0000000001020433000000000014004b0000079f0000413d000000000103041a000002ec0110019700000008024002100000ff000220018f000000000121019f000000000013041b0000000801000029000000000201041a000000000041041b001000000004001d000000000024004b000004520000813d000f00000002001d0000000801000029000000000010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000201043b0000000f012000290000001002200029000000000012004b000004520000813d000000000002041b0000000102200039000000000012004b0000044e0000413d0000000801000029000000000010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b00000000020000190000000a030000290000001006000029000000000412001900000020033000390000000005030433000002bc05500197000000000054041b0000000102200039000000000062004b000004610000413d0000000a010000290000000001010433000000000001004b000004d40000613d0000000002000019000f00000002001d000000050120021000000009011000290000000001010433001000000001001d0000001101000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000001002000029000002bc02200197000000000101043b001000000002001d000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b000000000101041a0000000801100270000000ff0110018f000000020010008c000001d40000213d000000000001004b000006a60000c13d000000100000006b0000000f02000029000006b10000613d000000400300043d000002c70030009c0000005f0000213d0000004001300039000000400010043f000000ff0120018f000d00000003001d00000000021304360000000101000039000e00000002001d00000000001204350000001101000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000001002000029000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000000e020000290000000002020433000000020020008c000001d40000213d00000008022002100000ff000220018f0000000d030000290000000003030433000000ff0330018f000000000223019f000000000101043b000000000301041a000002c803300197000000000232019f000000000021041b0000000f0200002900000001022000390000000a010000290000000001010433000000000012004b0000046e0000413d0000000c010000290000000001010433000000000001004b000005420000613d0000000002000019000f00000002001d00000005012002100000000b011000290000000001010433001000000001001d0000001101000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000001002000029000002bc02200197000000000101043b001000000002001d000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b000000000101041a0000000801100270000000ff0110018f000000020010008c000001d40000213d000000000001004b000006a60000c13d000000100000006b0000000f02000029000006b10000613d000000400300043d000002c70030009c0000005f0000213d0000004001300039000000400010043f000000ff0120018f000d00000003001d00000000021304360000000201000039000e00000002001d00000000001204350000001101000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000001002000029000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000000e020000290000000002020433000000020020008c000001d40000213d00000008022002100000ff000220018f0000000d030000290000000003030433000000ff0330018f000000000223019f000000000101043b000000000301041a000002c803300197000000000232019f000000000021041b0000000f0200002900000001022000390000000c010000290000000001010433000000000012004b000004d90000413d000002b60010009c000005430000a13d0000005f0000013d00000000010000190000000203000029000000000203041a000000000013041b001000000001001d000000000021004b0000055e0000813d000f00000002001d000000000030043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000201043b0000000f012000290000001002200029000000000012004b00000002030000290000055e0000813d000000000002041b0000000102200039000000000012004b0000055a0000413d000000000030043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000001006000029000000000006004b0000000c05000029000005760000613d0000000002000019000000000312001900000020055000390000000004050433000002bc04400197000000000043041b0000000102200039000000000062004b0000056e0000413d0000000303000029000000000103041a000002ea0110019700000004020000290000000002020433001000ff0020019300000010011001af000000000013041b000000070100002900000000010104330000000502000029000000000012041b00000001020000290000000002020433000f00000002001d000000400400043d0000004002400039000000a00300003900000000003204350000002002400039000000000012043500000011010000290000000000140435000e00000004001d000000a0014000390000000802000029000000000302041a000d00000003001d0000000000310435000000000020043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d0000000e05000029000000c0025000390000000d07000029000000000007004b000005ad0000613d000000000101043b00000000030000190000000f06000029000000000401041a000002bc04400197000000000242043600000001011000390000000103300039000000000073004b000005a50000413d000005ae0000013d0000000f0600002900000000015200490000006003500039000000000013043500000000030604330000000001320436000000000003004b000005bd0000613d000000000200001900000020066000390000000004060433000002bc0440019700000000014104360000000102200039000000000032004b000005b60000413d0000008002500039000000100300002900000000003204350000000001510049000002a30010009c000002a3010080410000006001100210000002a30050009c000002a3050080410000004002500210000000000121019f0000000002000414000002a30020009c000002a302008041000000c002200210000000000121019f000002ca011001c70000800d020000390000000103000039000002cb040000410a860a7c0000040f0000000100200190000002e60000613d000000400100043d00000011020000290000000000210435000002a30010009c000002a30100804100000040011002100000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002c5011001c70000800d020000390000000103000039000002cc040000410a860a7c0000040f0000000100200190000002e60000613d0000000602000029000600010020003d000000800100043d000000060010006b0000031d0000413d000000330000013d000000a007000039000005f40000013d000000a003800039000000000093043500000000078704360000002005500039000000000065004b000003130000813d000000000351034f000000000803043b000002b60080009c000002e60000213d00000000094800190000001108900069000002b90080009c000002e60000213d000000c00080008c000002e60000413d000000400800043d000002ba0080009c0000005f0000213d000000c003800039000000400030043f0000002403900039000000000331034f000000000303043b000000000b380436000000440a9000390000000003a1034f000000000c03043b000000ff00c0008c000002e60000213d0000000000cb0435000000200aa000390000000003a1034f000000000b03043b000000ff00b0008c000002e60000213d00000040038000390000000000b30435000000200aa000390000000003a1034f000000000b03043b00000000000b004b0000000003000039000000010300c03900000000003b004b000002e60000c13d00000060038000390000000000b30435000000200aa000390000000003a1034f000000000b03043b000002b600b0009c000002e60000213d000000000c9b00190000004303c00039000000000023004b000000000b000019000002bb0b008041000002bb03300197000000000003004b000000000d000019000002bb0d004041000002bb0030009c000000000d0bc01900000000000d004b000002e60000c13d0000002403c00039000000000331034f000000000d03043b000002b600d0009c0000005f0000213d000000050ed002100000003f03e00039000002b703300197000000400b00043d000000000f3b00190000000000bf004b00000000030000390000000103004039000002b600f0009c0000005f0000213d00000001003001900000005f0000c13d0000004000f0043f0000000000db0435000000440cc00039000000000dce001900000000002d004b000002e60000213d0000000000dc004b000006530000813d000000000e0b00190000000003c1034f000000000f03043b000002bc00f0009c000002e60000213d000000200ee000390000000000fe0435000000200cc000390000000000dc004b0000064a0000413d00000080038000390000000000b304350000002003a00039000000000331034f000000000a03043b000002b600a0009c000002e60000213d000000000a9a00190000004303a00039000000000023004b0000000009000019000002bb09008041000002bb03300197000000000003004b000000000b000019000002bb0b004041000002bb0030009c000000000b09c01900000000000b004b000002e60000c13d0000002403a00039000000000331034f000000000b03043b000002b600b0009c0000005f0000213d000000050cb002100000003f03c00039000002b703300197000000400900043d000000000d39001900000000009d004b00000000030000390000000103004039000002b600d0009c0000005f0000213d00000001003001900000005f0000c13d0000004000d0043f0000000000b90435000000440aa00039000000000bac001900000000002b004b000002e60000213d0000000000ba004b000005ee0000813d000000000c0900190000000003a1034f000000000d03043b000002bc00d0009c000002e60000213d000000200cc000390000000000dc0435000000200aa000390000000000ba004b000006810000413d000005ee0000013d000000400100043d000002bd020000410000000000210435000002a30010009c000002a3010080410000004001100210000002be011001c700000a88000104300000000000230435000002a30010009c000002a3010080410000004001100210000002e0011001c700000a870001042e000000400200043d00000024032000390000000000130435000002d901000041000000000012043500000004012000390000000b030000290000000000310435000002a30020009c000002a3020080410000004001200210000002d4011001c700000a8800010430000000400100043d000002c6020000410000000000210435000000040210003900000004030000390000000000320435000002a30010009c000002a3010080410000004001100210000002c2011001c700000a8800010430000000400100043d000002c9020000410000068d0000013d0000001102000029000000a402200039000000000003004b000006d30000613d0000000004000415000000130440008a0000000504400210000000800500043d0000000503500210000000000005004b000007840000c13d0000000504400270000000a00430003f0000000e0400002900000000050404330000000504500210000000000005004b000006c90000613d00000000055400d9000000200050008c000002d70000c13d000000a00530003900000000035400190000000002230019000000000032004b00000000030000390000000103004039000000000054001a000002d70000413d0000000100300190000002d70000c13d0000000003000031000000000023004b0000075c0000c13d00000004020000390000000102200367000000000202043b000000000021004b000007680000c13d000002d60100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000002a30010009c000002a301008041000000c001100210000002d7011001c700008005020000390a860a810000040f0000000100200190000007740000613d000000000101043b000900000001001d000002a50100004100000000001004430000000001000414000002a30010009c000002a301008041000000c001100210000002d8011001c70000800b020000390a860a810000040f0000000100200190000007740000613d000000000101043b000000090010006b000007ac0000c13d0000000d01000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000000002000411000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000400200043d000900000002001d000002c70020009c0000005f0000213d000000000101043b00000009030000290000004002300039000000400020043f000000000201041a000000ff0120018f00000000011304360000000802200270000000ff0220018f000000020020008c000001d40000213d0000000000210435000007750000c13d0000000d01000029000000000010043f0000000201000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d00000009020000290000000002020433000900ff00200193000000000101043b0000000301100039000000000201041a000000090020006b0000099a0000813d000000000010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000000901100029000000000101041a000002bc011001970000000002000411000000000012004b000007750000c13d0000000b010000290000000001010433000000000001004b0000089e0000c13d00000001010003670000000402100370000000000202043b000000400300043d00000000022304360000002401100370000000000101043b000002b6011001970000000000120435000002a30030009c000002a303008041000000400130021000000000020004140000015e0000013d000000400100043d00000024041000390000000000340435000002d303000041000000000031043500000004031000390000000000230435000002a30010009c000002a3010080410000004001100210000002d4011001c700000a8800010430000000400300043d00000024043000390000000000240435000002d502000041000000000023043500000004023000390000000000120435000002a30030009c000002a3030080410000004001300210000002d4011001c700000a8800010430000000000001042f000000400100043d000002e2020000410000068d0000013d000000400100043d000002c602000041000000000021043500000004021000390000000000020435000006ac0000013d000002c6010000410000000f03000029000000000013043500000004013000390000000502000039000007a60000013d00000000045300d9000000200040008c000002d70000c13d0000000004000415000000120440008a0000000504400210000002e90030009c000002d70000213d000006bf0000013d000000400100043d000002c102000041000000000021043500000004021000390000001103000029000006ab0000013d000000400100043d000002c602000041000000000021043500000004021000390000000203000039000006ab0000013d000000400100043d000002c602000041000000000021043500000004021000390000000303000039000006ab0000013d000000400100043d000f00000001001d000002c6010000410000000f030000290000000000130435000000040130003900000001020000390000000000210435000002a30030009c000002a3030080410000004001300210000002c2011001c700000a8800010430000000400200043d00000024032000390000000000130435000002d901000041000000000012043500000004012000390000000903000029000006a00000013d0000000c010000290000000001010433000000ff0110018f000000ff0010008c000002d70000613d0000000101100039000000800200043d000000000012004b000009a00000c13d00000011010000290000001f01100039000002eb011001970000003f01100039000002eb02100197000000400100043d0000000002210019000000000012004b00000000040000390000000104004039000002b60020009c0000005f0000213d00000001004001900000005f0000c13d000000400020043f000000110200002900000000022104360000000f05000029000000000050007c000002e60000213d0000001105000029000002eb045001980000001f0550018f0000000003420019000000100600002900000020066000390000000106600367000007df0000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000007db0000c13d000000000005004b000007ec0000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000043043500000011032000290000000000030435000002a30020009c000002a30200804100000040022002100000000001010433000002a30010009c000002a3010080410000006001100210000000000121019f0000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002ca011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000301043b000000400100043d00000020021000390000000000320435000000a003100039000000400410003900000004050000390000000105500367000000005605043c0000000004640436000000000034004b000008080000c13d00000080040000390000000000410435000002db0010009c0000005f0000213d000000400030043f000002a30020009c000002a30200804100000040022002100000000001010433000002a30010009c000002a3010080410000006001100210000000000121019f0000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002ca011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b000e00000001001d000000800100043d000c00000001001d000000000001004b000001510000613d001100000000001d001000000000001d00000011010000290000001f0010008c0000099a0000213d000000800100043d000000110010006c0000099a0000a13d00000011010000290000000501100210000000a0011000390000000001010433000000400200043d000000600320003900000000001304350000004003200039000000000013043500000020012000390000001b0300003900000000003104350000000e010000290000000000120435000000000000043f000002a30020009c000002a30200804100000040012002100000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002dd011001c700000001020000390a860a810000040f0000006003100270000002a303300197000000200030008c000000200500003900000000050340190000002004500190000008580000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b000008540000c13d0000001f05500190000008650000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f00000000005404350000000100200190000009a90000613d000000000100043d000f00000001001d0000000d01000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000000f02000029000002bc02200197000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000400200043d000002c70020009c0000005f0000213d000000000101043b0000004003200039000000400030043f000000000301041a000000ff0130018f00000000021204360000000803300270000000ff0330018f000000020030008c000001d40000213d0000000000320435000000010030008c000009a30000c13d000000010110020f0000001000100180000009a60000c13d00100010001001b300000011020000290000000102200039001100000002001d0000000c0020006c0000082b0000413d000001510000013d0000000a010000290000000001010433000000ff0110018f000000ff0010008c000002d70000613d000000800200043d0000000101100039000000000012004b000009a00000c13d000000400100043d0000000e030000290000000003030433000000000032004b000009c70000c13d00000011020000290000001f02200039000002eb022001970000003f02200039000002eb022001970000000002210019000000000012004b00000000040000390000000104004039000002b60020009c0000005f0000213d00000001004001900000005f0000c13d000000400020043f000000110200002900000000022104360000000f05000029000000000050007c000002e60000213d0000001105000029000002eb045001980000001f0550018f000000100300002900000001063003670000000003420019000008cc0000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000008c80000c13d000000000005004b000008d90000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000043043500000011032000290000000000030435000002a30020009c000002a30200804100000040022002100000000001010433000002a30010009c000002a3010080410000006001100210000000000121019f0000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002ca011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000301043b000000400100043d00000020021000390000000000320435000000a003100039000000400410003900000004050000390000000105500367000000005605043c0000000004640436000000000034004b000008f50000c13d00000080040000390000000000410435000002db0010009c0000005f0000213d000000400030043f000002a30020009c000002a30200804100000040022002100000000001010433000002a30010009c000002a3010080410000006001100210000000000121019f0000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002ca011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b000b00000001001d000000800100043d000a00000001001d000000000001004b0000074e0000613d001100000000001d001000000000001d00000011010000290000001f0010008c0000099a0000213d00000011010000290000000301100210000000c4020000390000000102200367000000000202043b00000000011201cf000002dc0010009c000002d70000213d000000800200043d000000110020006c0000099a0000a13d0000000e020000290000000002020433000000110020006c0000099a0000a13d000000f8011002700000001b01100039000000110200002900000005022002100000000c032000290000000003030433000000a0022000390000000002020433000000400400043d0000006005400039000000000035043500000040034000390000000000230435000000200240003900000000001204350000000b010000290000000000140435000000000000043f000002a30040009c000002a30400804100000040014002100000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002dd011001c700000001020000390a860a810000040f0000006003100270000002a303300197000000200030008c000000200500003900000000050340190000002004500190000009540000613d000000000601034f0000000007000019000000006806043c0000000007870436000000000047004b000009500000c13d0000001f05500190000009610000613d000000000641034f0000000305500210000000000704043300000000075701cf000000000757022f000000000606043b0000010005500089000000000656022f00000000055601cf000000000575019f00000000005404350000000100200190000009c90000613d000000000100043d000f00000001001d0000000d01000029000000000010043f0000000301000039000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000000101043b0000000f02000029000002bc02200197000000000020043f000000200010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002bf011001c700008010020000390a860a810000040f0000000100200190000002e60000613d000000400200043d000002c70020009c0000005f0000213d000000000101043b0000004003200039000000400030043f000000000301041a000000ff0130018f00000000021204360000000803300270000000ff0330018f000000020030008c000001d40000213d0000000000320435000000010030008c000009a30000c13d000000010110020f0000001000100180000009a60000c13d00100010001001b300000011020000290000000102200039001100000002001d0000000a0020006c000009180000413d0000074e0000013d000002e601000041000000000010043f0000003201000039000000040010043f000002c20100004100000a8800010430000000400100043d000002e3020000410000068d0000013d000000400100043d000002e4020000410000068d0000013d000000400100043d000002e5020000410000068d0000013d0000001f0530018f000002de06300198000000400200043d0000000004620019000009b40000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000009b00000c13d000000000005004b000009c10000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000002a30020009c000002a3020080410000004002200210000000000112019f00000a8800010430000002da020000410000068d0000013d0000001f0530018f000002de06300198000000400200043d0000000004620019000009b40000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000009d00000c13d000009b40000013d000000000301001900000000040104330000000001420436000000000004004b000009e20000613d000000000200001900000020033000390000000005030433000002bc0550019700000000015104360000000102200039000000000042004b000009db0000413d000000000001042d000002ed0010009c000009e80000813d0000006001100039000000400010043f000000000001042d000002e601000041000000000010043f0000004101000039000000040010043f000002c20100004100000a8800010430000002e10010009c000009f30000813d0000008001100039000000400010043f000000000001042d000002e601000041000000000010043f0000004101000039000000040010043f000002c20100004100000a88000104300000001f02200039000002eb022001970000000001120019000000000021004b00000000020000390000000102004039000002b60010009c00000a050000213d000000010020019000000a050000c13d000000400010043f000000000001042d000002e601000041000000000010043f0000004101000039000000040010043f000002c20100004100000a88000104300000000002010019000000400100043d000002e10010009c00000a210000813d0000008003100039000000400030043f000000000302041a00000000033104360000000102200039000000000202041a000000ff0420018f0000000000430435000002c0002001980000000003000039000000010300c039000000600410003900000000003404350000000802200270000000ff0220018f00000040031000390000000000230435000000000001042d000002e601000041000000000010043f0000004101000039000000040010043f000002c20100004100000a88000104300002000000000002000000000301041a000100000003001d0000000002320436000200000002001d000000000010043f0000000001000414000002a30010009c000002a301008041000000c001100210000002c5011001c700008010020000390a860a810000040f000000010020019000000a460000613d0000000105000029000000000005004b00000a440000613d000000000401043b00000000020000190000000201000029000000000304041a000002bc03300197000000000131043600000001044000390000000102200039000000000052004b00000a3c0000413d000000000001042d0000000201000029000000000001042d000000000100001900000a8800010430000000000001042f000002a30010009c000002a3010080410000004001100210000002a30020009c000002a3020080410000006002200210000000000112019f0000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002ca011001c700008010020000390a860a810000040f000000010020019000000a5c0000613d000000000101043b000000000001042d000000000100001900000a880001043000000000050100190000000000200443000000050030008c00000a6c0000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b00000a640000413d000002a30030009c000002a30300804100000060013002100000000002000414000002a30020009c000002a302008041000000c002200210000000000112019f000002ee011001c700000000020500190a860a810000040f000000010020019000000a7b0000613d000000000101043b000000000001042d000000000001042f00000a7f002104210000000102000039000000000001042d0000000002000019000000000001042d00000a84002104230000000102000039000000000001042d0000000002000019000000000001042d00000a860000043200000a870001042e00000a880001043000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffff00000000000000000000000000000000000000009a8a0592ac89c5ad3bc6df8224c17b485976f597df104ee20d0df415241f670b00000002000000000000000000000000000000800000010000000000000000009b15e16f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000a00000000000000000000000000000000000000000000000000000000000000000000000007ac0aa1900000000000000000000000000000000000000000000000000000000c673e58300000000000000000000000000000000000000000000000000000000c673e58400000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000f716f99f000000000000000000000000000000000000000000000000000000007ac0aa1a000000000000000000000000000000000000000000000000000000008da5cb5b0000000000000000000000000000000000000000000000000000000034a9c92d0000000000000000000000000000000000000000000000000000000034a9c92e0000000000000000000000000000000000000000000000000000000044e65e550000000000000000000000000000000000000000000000000000000079ba509700000000000000000000000000000000000000000000000000000000181f5a770000000000000000000000000000000000000000000000000000000026bf9d26000000000000000000000000000000000000000000000000ffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffff7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff3f8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffff2b5c74de00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000002000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000ff000087f6037c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff00000000000000000000000000000000000000000000000000000000000100000200000000000000000000000000000000000020000000000000000000000000367f56a200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000d6c62c9b000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000ab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547897ac1b2c12867721b284f3eb147bd4ab046d4eef1cf31c1d8988bfcfb962b530000000000000000000000000000000000000004000000800000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000080000000000000000002b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08e1192e100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004400000000000000000000000093df584c00000000000000000000000000000000000000000000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e020000020000000000000000000000000000004400000000000000000000000002000002000000000000000000000000000000040000000000000000000000000f01ce8500000000000000000000000000000000000000000000000000000000a75d88af00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5fe4ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffe0198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef00000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff80da0f08e80000000000000000000000000000000000000000000000000000000071253a2500000000000000000000000000000000000000000000000000000000ca31867a00000000000000000000000000000000000000000000000000000000f67bc7c4000000000000000000000000000000000000000000000000000000004e487b71000000000000000000000000000000000000000000000000000000004d756c74694f4352334261736548656c70657220312e302e30000000000000000000000000000000000000000000000000000060000000c00000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff000000000000000000000000000000000000000000000000ffffffffffffffa00200000200000000000000000000000000000000000000000000000000000000")

func DeployMultiOCR3HelperZK(auth *bind.TransactOpts, backend bind.ContractBackend) (common.
	Address, *generated.Transaction,

	*MultiOCR3Helper,

	error) {
	parsed, err := MultiOCR3HelperMetaData.GetAbi()
	if err !=
		nil {
		return common.Address{}, nil, nil,
			err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	address, ethTx, contract, err := generated.
		DeployContract(auth, parsed, common.
			FromHex(MultiOCR3HelperZKBin), backend)
	if err != nil {
		return common.
				Address{}, nil,
			nil, err
	}
	return address, ethTx, &MultiOCR3Helper{address: address, abi: *parsed, MultiOCR3HelperCaller: MultiOCR3HelperCaller{contract: contract}, MultiOCR3HelperTransactor: MultiOCR3HelperTransactor{
		contract: contract}, MultiOCR3HelperFilterer: MultiOCR3HelperFilterer{contract: contract}}, nil
}
