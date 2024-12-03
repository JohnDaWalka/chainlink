package ccip_reader_tester

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

type OffRampSourceChainConfig struct {
	Router    common.Address
	IsEnabled bool
	MinSeqNr  uint64
	OnRamp    []byte
}

var CCIPReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeValueJuels\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPMessageSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeValueJuels\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"emitCCIPMessageSent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"}],\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"emitCommitReportAccepted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"emitExecutionStateChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"getInboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"setDestChainSeqNr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"testNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"setInboundNonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"seqNr\",\"type\":\"uint64\"}],\"name\":\"setLatestPriceSequenceNumber\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceChainConfig\",\"type\":\"tuple\"}],\"name\":\"setSourceChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611960806100206000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c8063bfc9b78911610081578063c92236251161005b578063c9223625146101f9578063e83eabba1461020c578063e9d68a8e1461021f57600080fd5b8063bfc9b7891461017e578063c1a5a35514610191578063c7c1cba1146101e657600080fd5b806369600bca116100b257806369600bca1461010f5780639041be3d1461015857806393df28671461016b57600080fd5b80633f4b04aa146100ce5780634bf78697146100fa575b600080fd5b60035467ffffffffffffffff165b60405167ffffffffffffffff90911681526020015b60405180910390f35b61010d610108366004610a4e565b61023f565b005b61010d61011d366004610b89565b600380547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216919091179055565b6100dc610166366004610b89565b610298565b61010d610179366004610bf4565b6102c8565b61010d61018c366004610ea2565b610343565b61010d61019f36600461102d565b67ffffffffffffffff918216600090815260016020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001691909216179055565b61010d6101f4366004611060565b610385565b6100dc6102073660046110f2565b6103e2565b61010d61021a366004611145565b61042e565b61023261022d366004610b89565b610518565b6040516100f19190611265565b80600001516060015167ffffffffffffffff168267ffffffffffffffff167f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f328360405161028c91906113ab565b60405180910390a35050565b67ffffffffffffffff808216600090815260016020819052604082205491926102c2921690611503565b92915050565b67ffffffffffffffff84166000908152600260205260409081902090518491906102f59085908590611552565b908152604051908190036020019020805467ffffffffffffffff929092167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090921691909117905550505050565b602081015181516040517f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e49261037a929091611656565b60405180910390a150565b848667ffffffffffffffff168867ffffffffffffffff167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b878787876040516103d1949392919061172e565b60405180910390a450505050505050565b67ffffffffffffffff8316600090815260026020526040808220905161040b9085908590611552565b9081526040519081900360200190205467ffffffffffffffff1690509392505050565b67ffffffffffffffff808316600090815260208181526040918290208451815492860151938601519094167501000000000000000000000000000000000000000000027fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff93151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090931673ffffffffffffffffffffffffffffffffffffffff909516949094179190911791909116919091178155606082015182919060018201906105119082611839565b5050505050565b604080516080808201835260008083526020808401829052838501829052606080850181905267ffffffffffffffff87811684528383529286902086519485018752805473ffffffffffffffffffffffffffffffffffffffff8116865274010000000000000000000000000000000000000000810460ff1615159386019390935275010000000000000000000000000000000000000000009092049092169483019490945260018401805493949293918401916105d490611795565b80601f016020809104026020016040519081016040528092919081815260200182805461060090611795565b801561064d5780601f106106225761010080835404028352916020019161064d565b820191906000526020600020905b81548152906001019060200180831161063057829003601f168201915b5050505050815250509050919050565b803567ffffffffffffffff8116811461067557600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff811182821017156106cc576106cc61067a565b60405290565b604051610120810167ffffffffffffffff811182821017156106cc576106cc61067a565b6040805190810167ffffffffffffffff811182821017156106cc576106cc61067a565b6040516060810167ffffffffffffffff811182821017156106cc576106cc61067a565b6040516080810167ffffffffffffffff811182821017156106cc576106cc61067a565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156107a6576107a661067a565b604052919050565b600060a082840312156107c057600080fd5b6107c86106a9565b9050813581526107da6020830161065d565b60208201526107eb6040830161065d565b60408201526107fc6060830161065d565b606082015261080d6080830161065d565b608082015292915050565b73ffffffffffffffffffffffffffffffffffffffff8116811461083a57600080fd5b50565b803561067581610818565b600082601f83011261085957600080fd5b813567ffffffffffffffff8111156108735761087361067a565b6108a460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161075f565b8181528460208386010111156108b957600080fd5b816020850160208301376000918101602001919091529392505050565b600067ffffffffffffffff8211156108f0576108f061067a565b5060051b60200190565b600082601f83011261090b57600080fd5b8135602061092061091b836108d6565b61075f565b82815260059290921b8401810191818101908684111561093f57600080fd5b8286015b84811015610a4357803567ffffffffffffffff808211156109645760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d0301121561099d5760008081fd5b6109a56106a9565b6109b088850161083d565b8152604080850135848111156109c65760008081fd5b6109d48e8b83890101610848565b8a84015250606080860135858111156109ed5760008081fd5b6109fb8f8c838a0101610848565b83850152506080915081860135818401525082850135925083831115610a215760008081fd5b610a2f8d8a85880101610848565b908201528652505050918301918301610943565b509695505050505050565b60008060408385031215610a6157600080fd5b610a6a8361065d565b9150602083013567ffffffffffffffff80821115610a8757600080fd5b908401906101a08287031215610a9c57600080fd5b610aa46106d2565b610aae87846107ae565b8152610abc60a0840161083d565b602082015260c083013582811115610ad357600080fd5b610adf88828601610848565b60408301525060e083013582811115610af757600080fd5b610b0388828601610848565b6060830152506101008084013583811115610b1d57600080fd5b610b2989828701610848565b608084015250610b3c610120850161083d565b60a083015261014084013560c083015261016084013560e083015261018084013583811115610b6a57600080fd5b610b76898287016108fa565b8284015250508093505050509250929050565b600060208284031215610b9b57600080fd5b610ba48261065d565b9392505050565b60008083601f840112610bbd57600080fd5b50813567ffffffffffffffff811115610bd557600080fd5b602083019150836020828501011115610bed57600080fd5b9250929050565b60008060008060608587031215610c0a57600080fd5b610c138561065d565b9350610c216020860161065d565b9250604085013567ffffffffffffffff811115610c3d57600080fd5b610c4987828801610bab565b95989497509550505050565b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461067557600080fd5b600082601f830112610c9257600080fd5b81356020610ca261091b836108d6565b82815260069290921b84018101918181019086841115610cc157600080fd5b8286015b84811015610a435760408189031215610cde5760008081fd5b610ce66106f6565b610cef8261065d565b8152610cfc858301610c55565b81860152835291830191604001610cc5565b600082601f830112610d1f57600080fd5b81356020610d2f61091b836108d6565b82815260059290921b84018101918181019086841115610d4e57600080fd5b8286015b84811015610a4357803567ffffffffffffffff80821115610d735760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610dac5760008081fd5b610db46106a9565b610dbf88850161065d565b815260408085013584811115610dd55760008081fd5b610de38e8b83890101610848565b8a8401525060609350610df784860161065d565b908201526080610e0885820161065d565b93820193909352920135908201528352918301918301610d52565b600082601f830112610e3457600080fd5b81356020610e4461091b836108d6565b82815260069290921b84018101918181019086841115610e6357600080fd5b8286015b84811015610a435760408189031215610e805760008081fd5b610e886106f6565b813581528482013585820152835291830191604001610e67565b60006020808385031215610eb557600080fd5b823567ffffffffffffffff80821115610ecd57600080fd5b9084019060608287031215610ee157600080fd5b610ee9610719565b823582811115610ef857600080fd5b83016040818903811315610f0b57600080fd5b610f136106f6565b823585811115610f2257600080fd5b8301601f81018b13610f3357600080fd5b8035610f4161091b826108d6565b81815260069190911b8201890190898101908d831115610f6057600080fd5b928a01925b82841015610fb05785848f031215610f7d5760008081fd5b610f856106f6565b8435610f9081610818565b8152610f9d858d01610c55565b818d0152825292850192908a0190610f65565b845250505082870135915084821115610fc857600080fd5b610fd48a838501610c81565b81880152835250508284013582811115610fed57600080fd5b610ff988828601610d0e565b8583015250604083013593508184111561101257600080fd5b61101e87858501610e23565b60408201529695505050505050565b6000806040838503121561104057600080fd5b6110498361065d565b91506110576020840161065d565b90509250929050565b600080600080600080600060e0888a03121561107b57600080fd5b6110848861065d565b96506110926020890161065d565b955060408801359450606088013593506080880135600481106110b457600080fd5b925060a088013567ffffffffffffffff8111156110d057600080fd5b6110dc8a828b01610848565b92505060c0880135905092959891949750929550565b60008060006040848603121561110757600080fd5b6111108461065d565b9250602084013567ffffffffffffffff81111561112c57600080fd5b61113886828701610bab565b9497909650939450505050565b6000806040838503121561115857600080fd5b6111618361065d565b9150602083013567ffffffffffffffff8082111561117e57600080fd5b908401906080828703121561119257600080fd5b61119a61073c565b82356111a581610818565b8152602083013580151581146111ba57600080fd5b60208201526111cb6040840161065d565b60408201526060830135828111156111e257600080fd5b6111ee88828601610848565b6060830152508093505050509250929050565b6000815180845260005b818110156112275760208185018101518683018201520161120b565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815273ffffffffffffffffffffffffffffffffffffffff825116602082015260208201511515604082015267ffffffffffffffff6040830151166060820152600060608301516080808401526112c060a0840182611201565b949350505050565b600082825180855260208086019550808260051b84010181860160005b8481101561139e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a073ffffffffffffffffffffffffffffffffffffffff825116855285820151818787015261134782870182611201565b915050604080830151868303828801526113618382611201565b9250505060608083015181870152506080808301519250858203818701525061138a8183611201565b9a86019a94505050908301906001016112e5565b5090979650505050505050565b602081526113fc60208201835180518252602081015167ffffffffffffffff808216602085015280604084015116604085015280606084015116606085015280608084015116608085015250505050565b6000602083015161142560c084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060408301516101a08060e08501526114426101c0850183611201565b915060608501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06101008187860301818801526114808584611201565b945060808801519250818786030161012088015261149e8584611201565b945060a088015192506114ca61014088018473ffffffffffffffffffffffffffffffffffffffff169052565b60c088015161016088015260e08801516101808801528701518685039091018387015290506114f983826112c8565b9695505050505050565b67ffffffffffffffff81811683821601908082111561154b577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5092915050565b8183823760009101908152919050565b805160408084528151848201819052600092602091908201906060870190855b818110156115db578351805173ffffffffffffffffffffffffffffffffffffffff1684528501517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16858401529284019291850191600101611582565b50508583015187820388850152805180835290840192506000918401905b8083101561164a578351805167ffffffffffffffff1683528501517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16858301529284019260019290920191908501906115f9565b50979650505050505050565b60006040808301604084528086518083526060925060608601915060608160051b8701016020808a0160005b8481101561170e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08a8503018652815160a067ffffffffffffffff8083511687528583015182878901526116d983890182611201565b848d01518316898e01528b8501519092168b890152506080928301519290960191909152509482019490820190600101611682565b5050878203908801526117218189611562565b9998505050505050505050565b84815260006004851061176a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b846020830152608060408301526117846080830185611201565b905082606083015295945050505050565b600181811c908216806117a957607f821691505b6020821081036117e2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115611834576000816000526020600020601f850160051c810160208610156118115750805b601f850160051c820191505b818110156118305782815560010161181d565b5050505b505050565b815167ffffffffffffffff8111156118535761185361067a565b611867816118618454611795565b846117e8565b602080601f8311600181146118ba57600084156118845750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611830565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015611907578886015182559484019460019091019084016118e8565b508582101561194357878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000818000a",
}

var CCIPReaderTesterABI = CCIPReaderTesterMetaData.ABI

var CCIPReaderTesterBin = CCIPReaderTesterMetaData.Bin

func DeployCCIPReaderTester(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *generated_zks.Transaction, *CCIPReaderTester, error) {
	parsed, err := CCIPReaderTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated_zks.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated_zks.DeployContract(auth, parsed, common.FromHex(CCIPReaderTesterZKBin), backend)
		contractReturn := &CCIPReaderTester{address: address, abi: *parsed, CCIPReaderTesterCaller: CCIPReaderTesterCaller{contract: contractBind}, CCIPReaderTesterTransactor: CCIPReaderTesterTransactor{contract: contractBind}, CCIPReaderTesterFilterer: CCIPReaderTesterFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CCIPReaderTesterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated_zks.Transaction{Transaction: tx, Hash_zks: tx.Hash()}, &CCIPReaderTester{address: address, abi: *parsed, CCIPReaderTesterCaller: CCIPReaderTesterCaller{contract: contract}, CCIPReaderTesterTransactor: CCIPReaderTesterTransactor{contract: contract}, CCIPReaderTesterFilterer: CCIPReaderTesterFilterer{contract: contract}}, nil
}

type CCIPReaderTester struct {
	address common.Address
	abi     abi.ABI
	CCIPReaderTesterCaller
	CCIPReaderTesterTransactor
	CCIPReaderTesterFilterer
}

type CCIPReaderTesterCaller struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterTransactor struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterFilterer struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterSession struct {
	Contract     *CCIPReaderTester
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CCIPReaderTesterCallerSession struct {
	Contract *CCIPReaderTesterCaller
	CallOpts bind.CallOpts
}

type CCIPReaderTesterTransactorSession struct {
	Contract     *CCIPReaderTesterTransactor
	TransactOpts bind.TransactOpts
}

type CCIPReaderTesterRaw struct {
	Contract *CCIPReaderTester
}

type CCIPReaderTesterCallerRaw struct {
	Contract *CCIPReaderTesterCaller
}

type CCIPReaderTesterTransactorRaw struct {
	Contract *CCIPReaderTesterTransactor
}

func NewCCIPReaderTester(address common.Address, backend bind.ContractBackend) (*CCIPReaderTester, error) {
	abi, err := abi.JSON(strings.NewReader(CCIPReaderTesterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCCIPReaderTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTester{address: address, abi: abi, CCIPReaderTesterCaller: CCIPReaderTesterCaller{contract: contract}, CCIPReaderTesterTransactor: CCIPReaderTesterTransactor{contract: contract}, CCIPReaderTesterFilterer: CCIPReaderTesterFilterer{contract: contract}}, nil
}

func NewCCIPReaderTesterCaller(address common.Address, caller bind.ContractCaller) (*CCIPReaderTesterCaller, error) {
	contract, err := bindCCIPReaderTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCaller{contract: contract}, nil
}

func NewCCIPReaderTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*CCIPReaderTesterTransactor, error) {
	contract, err := bindCCIPReaderTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterTransactor{contract: contract}, nil
}

func NewCCIPReaderTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*CCIPReaderTesterFilterer, error) {
	contract, err := bindCCIPReaderTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterFilterer{contract: contract}, nil
}

func bindCCIPReaderTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CCIPReaderTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPReaderTester.Contract.CCIPReaderTesterCaller.contract.Call(opts, result, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.CCIPReaderTesterTransactor.contract.Transfer(opts)
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.CCIPReaderTesterTransactor.contract.Transact(opts, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPReaderTester.Contract.contract.Call(opts, result, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.contract.Transfer(opts)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.contract.Transact(opts, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _CCIPReaderTester.Contract.GetExpectedNextSequenceNumber(&_CCIPReaderTester.CallOpts, destChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _CCIPReaderTester.Contract.GetExpectedNextSequenceNumber(&_CCIPReaderTester.CallOpts, destChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetInboundNonce(opts *bind.CallOpts, sourceChainSelector uint64, sender []byte) (uint64, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getInboundNonce", sourceChainSelector, sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetInboundNonce(sourceChainSelector uint64, sender []byte) (uint64, error) {
	return _CCIPReaderTester.Contract.GetInboundNonce(&_CCIPReaderTester.CallOpts, sourceChainSelector, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetInboundNonce(sourceChainSelector uint64, sender []byte) (uint64, error) {
	return _CCIPReaderTester.Contract.GetInboundNonce(&_CCIPReaderTester.CallOpts, sourceChainSelector, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getLatestPriceSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _CCIPReaderTester.Contract.GetLatestPriceSequenceNumber(&_CCIPReaderTester.CallOpts)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _CCIPReaderTester.Contract.GetLatestPriceSequenceNumber(&_CCIPReaderTester.CallOpts)
}

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampSourceChainConfig)).(*OffRampSourceChainConfig)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _CCIPReaderTester.Contract.GetSourceChainConfig(&_CCIPReaderTester.CallOpts, sourceChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _CCIPReaderTester.Contract.GetSourceChainConfig(&_CCIPReaderTester.CallOpts, sourceChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitCCIPMessageSent(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitCCIPMessageSent", destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitCCIPMessageSent(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPMessageSent(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitCCIPMessageSent(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPMessageSent(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitCommitReportAccepted(opts *bind.TransactOpts, report OffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitCommitReportAccepted", report)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitCommitReportAccepted(report OffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCommitReportAccepted(&_CCIPReaderTester.TransactOpts, report)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitCommitReportAccepted(report OffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCommitReportAccepted(&_CCIPReaderTester.TransactOpts, report)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitExecutionStateChanged", sourceChainSelector, sequenceNumber, messageId, messageHash, state, returnData, gasUsed)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, messageHash, state, returnData, gasUsed)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, messageHash, state, returnData, gasUsed)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetDestChainSeqNr(opts *bind.TransactOpts, destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setDestChainSeqNr", destChainSelector, sequenceNumber)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetDestChainSeqNr(destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetDestChainSeqNr(&_CCIPReaderTester.TransactOpts, destChainSelector, sequenceNumber)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetDestChainSeqNr(destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetDestChainSeqNr(&_CCIPReaderTester.TransactOpts, destChainSelector, sequenceNumber)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetInboundNonce(opts *bind.TransactOpts, sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setInboundNonce", sourceChainSelector, testNonce, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetInboundNonce(sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetInboundNonce(&_CCIPReaderTester.TransactOpts, sourceChainSelector, testNonce, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetInboundNonce(sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetInboundNonce(&_CCIPReaderTester.TransactOpts, sourceChainSelector, testNonce, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetLatestPriceSequenceNumber(opts *bind.TransactOpts, seqNr uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setLatestPriceSequenceNumber", seqNr)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetLatestPriceSequenceNumber(seqNr uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetLatestPriceSequenceNumber(&_CCIPReaderTester.TransactOpts, seqNr)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetLatestPriceSequenceNumber(seqNr uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetLatestPriceSequenceNumber(&_CCIPReaderTester.TransactOpts, seqNr)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetSourceChainConfig(opts *bind.TransactOpts, sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setSourceChainConfig", sourceChainSelector, sourceChainConfig)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetSourceChainConfig(sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetSourceChainConfig(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sourceChainConfig)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetSourceChainConfig(sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetSourceChainConfig(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sourceChainConfig)
}

type CCIPReaderTesterCCIPMessageSentIterator struct {
	Event *CCIPReaderTesterCCIPMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterCCIPMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterCCIPMessageSent)
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
		it.Event = new(CCIPReaderTesterCCIPMessageSent)
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

func (it *CCIPReaderTesterCCIPMessageSentIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterCCIPMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*CCIPReaderTesterCCIPMessageSentIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCCIPMessageSentIterator{contract: _CCIPReaderTester.contract, event: "CCIPMessageSent", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterCCIPMessageSent)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseCCIPMessageSent(log types.Log) (*CCIPReaderTesterCCIPMessageSent, error) {
	event := new(CCIPReaderTesterCCIPMessageSent)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPReaderTesterCommitReportAcceptedIterator struct {
	Event *CCIPReaderTesterCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterCommitReportAccepted)
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
		it.Event = new(CCIPReaderTesterCommitReportAccepted)
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

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterCommitReportAccepted struct {
	MerkleRoots  []InternalMerkleRoot
	PriceUpdates InternalPriceUpdates
	Raw          types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*CCIPReaderTesterCommitReportAcceptedIterator, error) {

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCommitReportAcceptedIterator{contract: _CCIPReaderTester.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterCommitReportAccepted)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseCommitReportAccepted(log types.Log) (*CCIPReaderTesterCommitReportAccepted, error) {
	event := new(CCIPReaderTesterCommitReportAccepted)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPReaderTesterExecutionStateChangedIterator struct {
	Event *CCIPReaderTesterExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterExecutionStateChanged)
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
		it.Event = new(CCIPReaderTesterExecutionStateChanged)
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

func (it *CCIPReaderTesterExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	MessageHash         [32]byte
	State               uint8
	ReturnData          []byte
	GasUsed             *big.Int
	Raw                 types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*CCIPReaderTesterExecutionStateChangedIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterExecutionStateChangedIterator{contract: _CCIPReaderTester.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterExecutionStateChanged)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseExecutionStateChanged(log types.Log) (*CCIPReaderTesterExecutionStateChanged, error) {
	event := new(CCIPReaderTesterExecutionStateChanged)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CCIPReaderTester *CCIPReaderTester) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CCIPReaderTester.abi.Events["CCIPMessageSent"].ID:
		return _CCIPReaderTester.ParseCCIPMessageSent(log)
	case _CCIPReaderTester.abi.Events["CommitReportAccepted"].ID:
		return _CCIPReaderTester.ParseCommitReportAccepted(log)
	case _CCIPReaderTester.abi.Events["ExecutionStateChanged"].ID:
		return _CCIPReaderTester.ParseExecutionStateChanged(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CCIPReaderTesterCCIPMessageSent) Topic() common.Hash {
	return common.HexToHash("0x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32")
}

func (CCIPReaderTesterCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4")
}

func (CCIPReaderTesterExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b")
}

func (_CCIPReaderTester *CCIPReaderTester) Address() common.Address {
	return _CCIPReaderTester.address
}

type CCIPReaderTesterInterface interface {
	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetInboundNonce(opts *bind.CallOpts, sourceChainSelector uint64, sender []byte) (uint64, error)

	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error)

	EmitCCIPMessageSent(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error)

	EmitCommitReportAccepted(opts *bind.TransactOpts, report OffRampCommitReport) (*types.Transaction, error)

	EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error)

	SetDestChainSeqNr(opts *bind.TransactOpts, destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error)

	SetInboundNonce(opts *bind.TransactOpts, sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error)

	SetLatestPriceSequenceNumber(opts *bind.TransactOpts, seqNr uint64) (*types.Transaction, error)

	SetSourceChainConfig(opts *bind.TransactOpts, sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error)

	FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*CCIPReaderTesterCCIPMessageSentIterator, error)

	WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error)

	ParseCCIPMessageSent(log types.Log) (*CCIPReaderTesterCCIPMessageSent, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*CCIPReaderTesterCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*CCIPReaderTesterCommitReportAccepted, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*CCIPReaderTesterExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*CCIPReaderTesterExecutionStateChanged, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var CCIPReaderTesterZKBin string = ("0x000100000000000200080000000000020000000000010355000000800b0000390000004000b0043f0000000100200190000000950000c13d00000060021002700000022402200197000000040020008c000005b10000413d000000000301043b000000e003300270000002260030009c0000009d0000a13d000002270030009c000000e70000213d0000022b0030009c000003dd0000613d0000022c0030009c000002f10000613d0000022d0030009c000005b10000c13d000000e40020008c000005b10000413d0000000003000416000000000003004b000005b10000c13d0000000403100370000000000303043b000002340030009c000005b10000213d0000002404100370000000000404043b000002340040009c000005b10000213d0000008405100370000000000605043b000000030060008c000005b10000213d000000a405100370000000000805043b000002340080009c000005b10000213d0000002305800039000000000025004b000005b10000813d0000000409800039000000000591034f000000000705043b000002340070009c0000010b0000213d0000001f0a700039000002500aa001970000003f0aa00039000002500aa001970000023600a0009c0000010b0000213d000000800aa000390000004000a0043f000000800070043f00000000087800190000002408800039000000000028004b000005b10000213d000000000e0b00190000002002900039000000000821034f00000250097001980000001f0a70018f000000a0029000390000004e0000613d000000a00b000039000000000c08034f00000000cd0c043c000000000bdb043600000000002b004b0000004a0000c13d00000000000a004b0000005b0000613d000000000898034f0000000309a00210000000000a020433000000000a9a01cf000000000a9a022f000000000808043b0000010009900089000000000898022f00000000089801cf0000000008a8019f0000000000820435000000a00270003900000000000204350000006402100370000000000702043b000000400200043d00000040082000390000000000e804350000002008200039000000000068043500000000007204350000008007200039000000800600043d0000000000670435000000a007200039000000000006004b000000730000613d00000000080000190000000009780019000000a00a800039000000000a0a04330000000000a904350000002008800039000000000068004b0000006c0000413d00000000077600190000000000070435000000c407100370000000000707043b000000600820003900000000007804350000001f066000390000025005600197000000a005500039000002240050009c00000224050080410000006005500210000002240020009c00000224020080410000004002200210000000000225019f0000004401100370000000000701043b0000000001000414000002240010009c0000022401008041000000c001100210000000000121019f0000023405300197000002340640019700000242011001c70000800d0200003900000004030000390000024304000041088c08820000040f0000000100200190000005b10000613d00000000010000190000088d0001042e0000000001000416000000000001004b000005b10000c13d00000020010000390000010000100443000001200000044300000225010000410000088d0001042e0000022e0030009c000001110000a13d0000022f0030009c000003cd0000613d000002300030009c000002d10000613d000002310030009c000005b10000c13d000000640020008c000005b10000413d0000000003000416000000000003004b000005b10000c13d0000000403100370000000000303043b000800000003001d000002340030009c000005b10000213d0000002403100370000000000303043b000700000003001d000002340030009c000005b10000213d0000004401100370000000000101043b000002340010009c000005b10000213d0000000401100039088c08410000040f0000000803000029000000000030043f0000000203000039000000200030043f000800000001001d000600000002001d00000040020000390000000001000019088c086d0000040f000000060a0000290000001f04a0018f0000025005a00198000000080200002900000000062003670000000002010019000000400100043d0000000003510019000000d20000613d000000000706034f0000000008010019000000007907043c0000000008980436000000000038004b000000ce0000c13d000000000004004b000000df0000613d000000000556034f0000000304400210000000000603043300000000064601cf000000000646022f000000000505043b0000010004400089000000000545022f00000000044501cf000000000464019f00000000004304350000000003a1001900000000002304350000002002a00039088c086d0000040f000000000201041a00000244022001970000000703000029000003080000013d000002280030009c000005a30000613d000002290030009c0000030c0000613d0000022a0030009c000005b10000c13d000000240020008c000005b10000413d0000000002000416000000000002004b000005b10000c13d0000000401100370000000000101043b000002340010009c000005b10000213d0000010002000039000000400020043f000000800000043f000000a00000043f000000c00000043f0000006002000039000000e00020043f000000000010043f000000200000043f0000000001000414000002240010009c0000022401008041000000c00110021000000235011001c70000801002000039088c08870000040f0000000100200190000005b10000613d000000400500043d000002360050009c000005e90000a13d0000023f01000041000000000010043f0000004101000039000000040010043f00000240010000410000088e00010430000002320030009c000002c80000613d000002330030009c000005b10000c13d000000440020008c000005b10000413d0000000003000416000000000003004b000005b10000c13d0000000403100370000000000303043b000002340030009c000005b10000213d0000002404100370000000000604043b000002340060009c000005b10000213d00000000046200490000023a0040009c000005b10000213d000001a40040008c000005b10000413d0000024007000039000000400070043f0000000404600039000000000441034f000000000404043b000001a00040043f0000002404600039000000000541034f000000000505043b000002340050009c000005b10000213d000001c00050043f0000002004400039000000000541034f000000000505043b000002340050009c000005b10000213d000001e00050043f0000002004400039000000000541034f000000000505043b000002340050009c000005b10000213d000002000050043f0000002004400039000000000541034f000000000505043b000002340050009c000005b10000213d000002200050043f000001a005000039000000800050043f0000002004400039000000000841034f000000000808043b000002370080009c000005b10000213d000000a00080043f0000002008400039000000000481034f000000000404043b000002340040009c000005b10000213d000000000a6400190000002304a00039000000000024004b000005b10000813d000000040ba000390000000004b1034f000000000904043b0000024c0090009c0000010b0000813d0000001f0c900039000002500cc001970000003f0cc00039000002500cc001970000024d00c0009c0000010b0000213d000002400cc000390000004000c0043f000002400090043f000000000a9a0019000000240aa0003900000000002a004b000005b10000213d000000200ab00039000000000ba1034f000002500c9001980000001f0d90018f000002600ac00039000001740000613d000002600e000039000000000f0b034f00000000f40f043c000000000e4e04360000000000ae004b000001700000c13d00000000000d004b000001810000613d0000000004cb034f000000030bd00210000000000c0a0433000000000cbc01cf000000000cbc022f000000000404043b000001000bb000890000000004b4022f0000000004b401cf0000000004c4019f00000000004a043500000260049000390000000000040435000000c00070043f0000002007800039000000000471034f000000000804043b000002340080009c000005b10000213d000000000b6800190000002304b00039000000000024004b000005b10000813d000000040cb000390000000004c1034f000000000804043b000002340080009c0000010b0000213d0000001f0480003900000250044001970000003f044000390000025004400197000000400900043d000000000a49001900000000009a004b000000000d000039000000010d0040390000023400a0009c0000010b0000213d0000000100d001900000010b0000c13d0000004000a0043f000000000d89043600000000048b00190000002404400039000000000024004b000005b10000213d0000002004c000390008000000410353000002500c8001980000001f0e80018f000000000bcd0019000001b10000613d000000080f00035f000000000a0d001900000000f40f043c000000000a4a04360000000000ba004b000001ad0000c13d00000000000e004b000001be0000613d0000000804c0035f000000030ae00210000000000c0b0433000000000cac01cf000000000cac022f000000000404043b000001000aa000890000000004a4022f0000000004a401cf0000000004c4019f00000000004b043500000000048d00190000000000040435000000e00090043f0000002007700039000000000471034f000000000804043b000002340080009c000005b10000213d000000000b6800190000002304b00039000000000024004b000005b10000813d000000040cb000390000000004c1034f000000000804043b000002340080009c0000010b0000213d0000001f0480003900000250044001970000003f044000390000025004400197000000400900043d000000000a49001900000000009a004b000000000d000039000000010d0040390000023400a0009c0000010b0000213d0000000100d001900000010b0000c13d0000004000a0043f000000000d89043600000000048b00190000002404400039000000000024004b000005b10000213d0000002004c000390008000000410353000002500c8001980000001f0e80018f000000000bcd0019000001ee0000613d000000080f00035f000000000a0d001900000000f40f043c000000000a4a04360000000000ba004b000001ea0000c13d00000000000e004b000001fb0000613d0000000804c0035f000000030ae00210000000000c0b0433000000000cac01cf000000000cac022f000000000404043b000001000aa000890000000004a4022f0000000004a401cf0000000004c4019f00000000004b043500000000048d00190000000000040435000001000090043f0000002007700039000000000471034f000000000804043b000002370080009c000005b10000213d000001200080043f0000002004700039000000000441034f000000000404043b000001400040043f0000004004700039000000000441034f000000000404043b000001600040043f0000006004700039000000000441034f000000000704043b000002340070009c000005b10000213d0000000004670019000400000004001d0000002304400039000000000024004b000005b10000813d00000004040000290000000404400039000000000441034f000000000604043b000002340060009c0000010b0000213d00000005076002100000003f047000390000024504400197000000400900043d0000000008490019000100000009001d000000000098004b00000000090000390000000109004039000002340080009c0000010b0000213d00000001009001900000010b0000c13d000000400080043f000000010400002900000000006404350000000404000029000600240040003d0000000604700029000300000004001d000000000024004b000005b10000213d000000000006004b000006da0000c13d0000000101000029000001800010043f000000800100043d00000060011000390000000002010433000000400800043d00000020010000390000000001180436000000800400043d0000000076040434000000000061043500000000010704330000023401100197000000400680003900000000001604350000004001400039000000000101043300000234011001970000006006800039000000000016043500000060014000390000000001010433000002340110019700000080068000390000000000160435000000800140003900000000010104330000023401100197000000a0048000390000000000140435000000a00100043d0000023701100197000000c0048000390000000000140435000000e001800039000000c00400043d0000000000510435000001c00180003900000000750404340000000000510435000800000008001d000001e006800039000000000005004b000002680000613d000000000800001900000000016800190000000004870019000000000404043300000000004104350000002008800039000000000058004b000002610000413d000000000156001900000000000104350000001f015000390000025001100197000000000116001900000008050000290000000004510049000000200440008a0000010005500039000000e00600043d000000000045043500000000760604340000000005610436000000000006004b0000027f0000613d000000000800001900000000015800190000000004870019000000000404043300000000004104350000002008800039000000000068004b000002780000413d000000000165001900000000000104350000001f016000390000025001100197000000000115001900000008050000290000000004510049000000200440008a0000012005500039000001000600043d000000000045043500000000670604340000000005710436000000000007004b000002960000613d000000000800001900000000015800190000000004860019000000000404043300000000004104350000002008800039000000000078004b0000028f0000413d000702340020019b00000000017500190000000000010435000001200100043d00000237011001970000000804000029000001400240003900000000001204350000016001400039000001400200043d00000000002104350000018001400039000001600200043d00000000002104350000001f01700039000002500110019700000000021500190000000001420049000000200110008a000001a004400039000001800500043d00000000001404350000000007050433000000000072043500000005017002100000000001120019000000200b100039000000000007004b000007d90000c13d000000080200002900000000012b0049000002240010009c00000224010080410000006001100210000002240020009c00000224020080410000004002200210000000000121019f0000000002000414000002240020009c0000022402008041000000c002200210000000000121019f00000242011001c700000234053001970000800d0200003900000003030000390000024e0400004100000007060000290000059f0000013d0000000001000416000000000001004b000005b10000c13d0000000301000039000000000101041a0000023401100197000000800010043f0000024f010000410000088d0001042e000000240020008c000005b10000413d0000000002000416000000000002004b000005b10000c13d0000000401100370000000000101043b000002340010009c000005b10000213d000000000010043f0000000101000039000000200010043f0000000001000414000002240010009c0000022401008041000000c00110021000000235011001c70000801002000039088c08870000040f0000000100200190000005b10000613d000000000101043b000000000101041a0000023401100197000002340010009c000005e10000c13d0000023f01000041000000000010043f0000001101000039000000040010043f00000240010000410000088e00010430000000440020008c000005b10000413d0000000002000416000000000002004b000005b10000c13d0000000402100370000000000202043b000002340020009c000005b10000213d0000002401100370000000000101043b000800000001001d000002340010009c000005b10000213d000000000020043f0000000101000039000000200010043f00000040020000390000000001000019088c086d0000040f000000000201041a00000244022001970000000803000029000000000232019f000000000021041b00000000010000190000088d0001042e000000440020008c000005b10000413d0000000003000416000000000003004b000005b10000c13d0000000403100370000000000303043b000002340030009c000005b10000213d0000002404100370000000000504043b000002340050009c000005b10000213d00000000045200490000023a0040009c000005b10000213d000000840040008c000005b10000413d0000010004000039000000400040043f0000000406500039000000000761034f000000000707043b000002370070009c000005b10000213d000000800070043f0000002006600039000000000761034f000000000707043b000000000007004b0000000008000039000000010800c039000000000087004b000005b10000c13d000000a00070043f0000002006600039000000000761034f000000000707043b000002340070009c000005b10000213d000000c00070043f0000002006600039000000000661034f000000000606043b000002340060009c000005b10000213d00000000065600190000002305600039000000000025004b000005b10000813d0000000407600039000000000571034f000000000505043b000002340050009c0000010b0000213d0000001f0850003900000250088001970000003f0880003900000250088001970000023b0080009c0000010b0000213d0000010008800039000000400080043f000001000050043f00000000065600190000002406600039000000000026004b000005b10000213d0000002002700039000000000221034f00000250065001980000001f0750018f00000120016000390000035c0000613d0000012008000039000000000902034f000000009a09043c0000000008a80436000000000018004b000003580000c13d000000000007004b000003690000613d000000000262034f0000000306700210000000000701043300000000076701cf000000000767022f000000000202043b0000010006600089000000000262022f00000000026201cf000000000272019f000000000021043500000120015000390000000000010435000000e00040043f000000000030043f000000200000043f0000000001000414000002240010009c0000022401008041000000c00110021000000235011001c70000801002000039088c08870000040f0000000100200190000005b10000613d000000800200043d0000023702200197000000000101043b000000000301041a0000023c03300197000000000223019f000000a00300043d000000000003004b0000023d030000410000000003006019000000000232019f000000c00300043d000000a8033002100000023e03300197000000000232019f000000000021041b000000e00200043d000600000002001d0000000032020434000700000003001d000800000002001d000002340020009c0000010b0000213d0000000101100039000500000001001d000000000101041a000000010210019000000001011002700000007f0110618f000400000001001d0000001f0010008c00000000010000390000000101002039000000000012004b000006030000c13d0000000401000029000000200010008c000003b90000413d0000000501000029000000000010043f0000000001000414000002240010009c0000022401008041000000c00110021000000239011001c70000801002000039088c08870000040f0000000100200190000005b10000613d00000008030000290000001f023000390000000502200270000000200030008c0000000002004019000000000301043b00000004010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b000003b90000813d000000000002041b0000000102200039000000000012004b000003b50000413d00000008010000290000001f0010008c0000067c0000a13d0000000501000029000000000010043f0000000001000414000002240010009c0000022401008041000000c00110021000000239011001c70000801002000039088c08870000040f0000000100200190000005b10000613d000000200200008a0000000802200180000000000101043b000006880000c13d0000002003000039000006950000013d000000240020008c000005b10000413d0000000002000416000000000002004b000005b10000c13d0000000401100370000000000101043b000002340010009c000005b10000213d0000000302000039000000000302041a0000024403300197000000000113019f000000000012041b00000000010000190000088d0001042e000000240320008c000005b10000413d0000000004000416000000000004004b000005b10000c13d0000000404100370000000000404043b000002340040009c000005b10000213d00000000054200490000023a0050009c000005b10000213d000000640050008c000005b10000413d000000e005000039000000400050043f0000000406400039000000000661034f000000000606043b000002340060009c000005b10000213d000000000646001900000000076200490000023a0070009c000005b10000213d000000440070008c000005b10000413d0000012007000039000000400070043f0000000408600039000000000981034f000000000909043b000002340090009c000005b10000213d0000000009690019000000230a90003900000000002a004b000005b10000813d000000040a900039000000000aa1034f000000000b0a043b0000023400b0009c0000010b0000213d000000050ab002100000003f0aa00039000002450aa001970000024600a0009c0000010b0000213d000001200aa000390000004000a0043f0000012000b0043f0000002409900039000000060ab00210000000000a9a001900000000002a004b000005b10000213d00000000000b004b000006610000c13d000000e00070043f0000002007800039000000000771034f000000000707043b000002340070009c000005b10000213d00000000076700190000002306700039000000000026004b000000000800001900000249080080410000024906600197000000000006004b00000000090000190000024909004041000002490060009c000000000908c019000000000009004b000005b10000c13d0000000406700039000000000661034f000000000806043b000002340080009c0000010b0000213d00000005068002100000003f066000390000024509600197000000400600043d0000000009960019000000000069004b000000000a000039000000010a004039000002340090009c0000010b0000213d0000000100a001900000010b0000c13d000000400090043f0000000000860435000000240770003900000006088002100000000008780019000000000028004b000005b10000213d000000000087004b0000045f0000813d0000000009060019000000000a7200490000023a00a0009c000005b10000213d0000004000a0008c000005b10000413d000000400a00043d0000024700a0009c0000010b0000213d000000400ba000390000004000b0043f000000000b71034f000000000b0b043b0000023400b0009c000005b10000213d000000000bba0436000000200c700039000000000cc1034f000000000c0c043b0000024800c0009c000005b10000213d00000020099000390000000000cb04350000000000a904350000004007700039000000000087004b000004450000413d000001000060043f000000800050043f0000002405400039000000000551034f000000000505043b000002340050009c000005b10000213d0000000005450019000500000005001d0000002305500039000000000025004b000000000600001900000249060080410000024905500197000000000005004b00000000070000190000024907004041000002490050009c000000000706c019000000000007004b000005b10000c13d00000005050000290000000405500039000000000551034f000000000505043b000002340050009c0000010b0000213d00000005065002100000003f076000390000024507700197000000400800043d0000000007780019000300000008001d000000000087004b00000000080000390000000108004039000002340070009c0000010b0000213d00000001008001900000010b0000c13d000000400070043f0000000307000029000000000057043500000005050000290000002407500039000400000076001d000000040020006b000005b10000213d0000000005070019000000040070006c0000050b0000813d000800030000002d000700000005001d000000000551034f000000000505043b000002340050009c000005b10000213d000000050550002900000000065300490000023a0060009c000005b10000213d000000a00060008c000005b10000413d000000400b00043d0000024a00b0009c0000010b0000213d000000a006b00039000000400060043f0000002406500039000000000861034f000000000808043b000002340080009c000005b10000213d00000000078b0436000600000007001d000000200c6000390000000006c1034f000000000606043b000002340060009c000005b10000213d00000000055600190000004306500039000000000026004b000000000800001900000249080080410000024906600197000000000006004b00000000090000190000024909004041000002490060009c000000000908c019000000000009004b000005b10000c13d0000002408500039000000000681034f000000000e06043b0000023400e0009c0000010b0000213d0000001f06e0003900000250066001970000003f066000390000025006600197000000400f00043d00000000066f00190000000000f6004b00000000090000390000000109004039000002340060009c0000010b0000213d00000001009001900000010b0000c13d000000400060043f000000000aef04360000000005e500190000004405500039000000000025004b000005b10000213d0000002005800039000000000d51034f0000025009e0019800000000059a0019000004df0000613d00000000080d034f00000000060a0019000000008708043c0000000006760436000000000056004b000004db0000c13d0000001f06e00190000004ec0000613d00000000079d034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f00000000006504350000000005ea0019000000000005043500000006050000290000000000f504350000002005c00039000000000651034f000000000606043b000002340060009c000005b10000213d0000004007b0003900000000006704350000002005500039000000000651034f000000000606043b000002340060009c000005b10000213d000000080800002900000020088000390000006007b0003900000000006704350000002005500039000000000551034f000000000505043b0000008006b000390000000000560435000800000008001d0000000000b8043500000007050000290000002005500039000000040050006c000004930000413d0000000303000029000000a00030043f0000004403400039000000000331034f000000000303043b000002340030009c000005b10000213d00000000044300190000002303400039000000000023004b000000000500001900000249050080410000024903300197000000000003004b00000000060000190000024906004041000002490030009c000000000605c019000000000006004b000005b10000c13d0000000403400039000000000331034f000000000503043b000002340050009c0000010b0000213d00000005035002100000003f033000390000024506300197000000400300043d0000000006630019000000000036004b00000000070000390000000107004039000002340060009c0000010b0000213d00000001007001900000010b0000c13d000000400060043f0000000000530435000000240440003900000006055002100000000005450019000000000025004b000005b10000213d000000000054004b000005500000813d000000000603001900000000074200490000023a0070009c000005b10000213d000000400070008c000005b10000413d000000400700043d000002470070009c0000010b0000213d00000020066000390000004008700039000000400080043f000000000841034f000000000808043b00000000088704360000002009400039000000000991034f000000000909043b000000000098043500000000007604350000004004400039000000000054004b0000053a0000413d000000c00030043f000000800100043d000800000001001d000000a00600043d000000400100043d00000040020000390000000002210436000700000002001d000000400210003900000000070604330000000000720435000000600810003900000005027002100000000002820019000000000007004b000006a80000c13d00000000031200490000000704000029000000000034043500000008030000290000000056030434000000400300003900000000033204360000004004200039000000000706043300000000007404350000006004200039000000000007004b0000057a0000613d00000000080000190000002006600039000000000906043300000000a909043400000237099001970000000009940436000000000a0a0433000002480aa001970000000000a9043500000040044000390000000108800039000000000078004b0000056e0000413d00000000050504330000000002240049000000000023043500000000030504330000000002340436000000000003004b0000058e0000613d00000000040000190000002005500039000000000605043300000000760604340000023406600197000000000662043600000000070704330000024807700197000000000076043500000040022000390000000104400039000000000034004b000005820000413d0000000002120049000002240020009c00000224020080410000006002200210000002240010009c00000224010080410000004001100210000000000112019f0000000002000414000002240020009c0000022402008041000000c002200210000000000121019f00000242011001c70000800d0200003900000001030000390000024b04000041088c08820000040f0000000100200190000000930000c13d000005b10000013d000000440020008c000005b10000413d0000000003000416000000000003004b000005b10000c13d0000000403100370000000000303043b000800000003001d000002340030009c000005b10000213d0000002401100370000000000101043b000002340010009c000005b30000a13d00000000010000190000088e000104300000000401100039088c08410000040f0000000803000029000000000030043f0000000203000039000000200030043f000800000001001d000700000002001d00000040020000390000000001000019088c086d0000040f000000070a0000290000001f04a0018f0000025005a00198000000080200002900000000062003670000000002010019000000400100043d0000000003510019000005cd0000613d000000000706034f0000000008010019000000007907043c0000000008980436000000000038004b000005c90000c13d000000000004004b000005da0000613d000000000556034f0000000304400210000000000603043300000000064601cf000000000646022f000000000505043b0000010004400089000000000545022f00000000044501cf000000000464019f00000000004304350000000003a1001900000000002304350000002002a00039088c086d0000040f000000000101041a0000023401100197000005e20000013d0000000101100039000000400200043d0000000000120435000002240020009c0000022402008041000000400120021000000241011001c70000088d0001042e000000000101043b0000008002500039000000400020043f000000000201041a00000237032001970000000007350436000000a80320027000000234033001970000004006500039000000000036043500000238002001980000000002000039000000010200c03900000000002704350000000102100039000000000902041a000000010490019000000001089002700000007f0380018f00000000080360190000001f0080008c00000000010000390000000101002039000000000119013f0000000100100190000006090000613d0000023f01000041000000000010043f0000002201000039000000040010043f00000240010000410000088e00010430000500000007001d000600000006001d000700000005001d000000400100043d0000000005810436000000000004004b000800000001001d000006170000c13d00000251029001970000000000250435000000000003004b00000020020000390000000002006039000006370000013d000400000009001d000300000008001d000000000020043f0000000001000414000002240010009c0000022401008041000000c00110021000000239011001c70000801002000039088c08870000040f0000000100200190000005b10000613d0000000402000029000000020020008c000006290000813d00000020050000390000000801000029000006360000013d000000000601043b000000000200001900000008010000290000000307000029000000000302001900000020022000390000000004120019000000000506041a00000000005404350000000106600039000000000072004b0000062d0000413d0000004005300039000000000201001900000000021200490000000002520019088c082f0000040f00000007030000290000006001300039000000080200002900000000002104350000002002000039000000400400043d000800000004001d000000000224043600000000030304330000023703300197000000000032043500000005020000290000000002020433000000000002004b0000000002000039000000010200c03900000040034000390000000000230435000000060200002900000000020204330000023402200197000000600340003900000000002304350000000001010433000000800240003900000080030000390000000000320435000000a002400039088c085b0000040f00000008020000290000000001210049000002240020009c00000224020080410000004002200210000002240010009c00000224010080410000006001100210000000000121019f0000088d0001042e000001400b000039000000000c9200490000023a00c0009c000005b10000213d0000004000c0008c000005b10000413d000000400c00043d0000024700c0009c0000010b0000213d000000400dc000390000004000d0043f000000000d91034f000000000d0d043b0000023700d0009c000005b10000213d000000000ddc0436000000200e900039000000000ee1034f000000000e0e043b0000024800e0009c000005b10000213d0000000000ed0435000000000bcb043600000040099000390000000000a9004b000006620000413d000004170000013d000000080000006b0000000001000019000006810000613d0000000701000029000000000101043300000008040000290000000302400210000002520220027f0000025202200167000000000221016f0000000101400210000006a30000013d000000010320008a0000000503300270000000000431001900000020030000390000000104400039000000060600002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b0000068e0000c13d000000080020006c000006a00000813d00000008020000290000000302200210000000f80220018f000002520220027f000002520220016700000006033000290000000003030433000000000223016f000000000021041b000000010100003900000008020000290000000102200210000000000112019f0000000502000029000000000012041b00000000010000190000088d0001042e000000a009000039000000000b000019000006c10000013d0000000003de001900000000000304350000004003c0003900000000030304330000023403300197000000400420003900000000003404350000006003c00039000000000303043300000234033001970000006004200039000000000034043500000080022000390000008003c00039000000000303043300000000003204350000001f02d00039000002500220019700000000022e0019000000010bb0003900000000007b004b000005600000813d0000000004120049000000600440008a00000000084804360000002006600039000000000c06043300000000d40c043400000234044001970000000004420436000000000d0d04330000000000940435000000a00420003900000000fd0d04340000000000d40435000000c00e20003900000000000d004b000006ab0000613d00000000040000190000000003e4001900000000054f00190000000005050433000000000053043500000020044000390000000000d4004b000006d20000413d000006ab0000013d0002002400200092000500010000002d000006eb0000013d000000050800002900000020088000390000000004d600190000000000040435000000070600002900000080046000390000000000740435000500000008001d000000000068043500000006040000290000002004400039000600000004001d000000030040006c000002340000813d0000000604100360000000000604043b000002340060009c000005b10000213d0000000406600029000800000006001d00000002066000690000023a0060009c000005b10000213d000000a00060008c000005b10000413d000000400400043d000700000004001d0000024a0040009c0000010b0000213d0000000704000029000000a004400039000000400040043f00000008040000290000002406400039000000000461034f000000000704043b000002370070009c000005b10000213d0000000704000029000000000f740436000000200e6000390000000004e1034f000000000604043b000002340060009c000005b10000213d000000080a6000290000004304a00039000000000024004b000000000600001900000249060080410000024904400197000000000004004b00000000070000190000024907004041000002490040009c000000000706c019000000000007004b000005b10000c13d000000240ba000390000000004b1034f000000000704043b000002340070009c0000010b0000213d0000001f0470003900000250044001970000003f044000390000025004400197000000400600043d0000000008460019000000000068004b00000000090000390000000109004039000002340080009c0000010b0000213d00000001009001900000010b0000c13d000000400080043f000000000d76043600000000047a00190000004404400039000000000024004b000005b10000213d0000002004b00039000000000c41034f0000025008700198000000000a8d00190000073a0000613d000000000b0c034f00000000090d001900000000b40b043c00000000094904360000000000a9004b000007360000c13d0000001f09700190000007470000613d00000000048c034f000000030890021000000000090a043300000000098901cf000000000989022f000000000404043b0000010008800089000000000484022f00000000048401cf000000000494019f00000000004a043500000000047d0019000000000004043500000000006f0435000000200ee000390000000004e1034f000000000604043b000002340060009c000005b10000213d00000008096000290000004304900039000000000024004b000000000600001900000249060080410000024904400197000000000004004b00000000070000190000024907004041000002490040009c000000000706c019000000000007004b000005b10000c13d000000240a9000390000000004a1034f000000000704043b000002340070009c0000010b0000213d0000001f0470003900000250044001970000003f044000390000025004400197000000400f00043d00000000064f00190000000000f6004b00000000080000390000000108004039000002340060009c0000010b0000213d00000001008001900000010b0000c13d000000400060043f00000000067f043600000000047900190000004404400039000000000024004b000005b10000213d0000002004a00039000000000a41034f000002500870019800000000098600190000077f0000613d000000000b0a034f000000000c06001900000000b40b043c000000000c4c043600000000009c004b0000077b0000c13d0000001f0b7001900000078c0000613d00000000048a034f0000000308b00210000000000a090433000000000a8a01cf000000000a8a022f000000000404043b0000010008800089000000000484022f00000000048401cf0000000004a4019f000000000049043500000000047600190000000000040435000000070600002900000040046000390000000000f404350000002004e00039000000000441034f000000000404043b000000600660003900000000004604350000004004e00039000000000441034f000000000604043b000002340060009c000005b10000213d00000008096000290000004304900039000000000024004b000000000600001900000249060080410000024904400197000000000004004b00000000070000190000024907004041000002490040009c000000000706c019000000000007004b000005b10000c13d000000240a9000390000000004a1034f000000000d04043b0000023400d0009c0000010b0000213d0000001f04d0003900000250044001970000003f044000390000025004400197000000400700043d0000000006470019000000000076004b00000000080000390000000108004039000002340060009c0000010b0000213d00000001008001900000010b0000c13d000000400060043f0000000006d704360000000004d900190000004404400039000000000024004b000005b10000213d0000002004a00039000000000a41034f0000025008d001980000000009860019000007cb0000613d000000000b0a034f000000000c06001900000000b40b043c000000000c4c043600000000009c004b000007c70000c13d0000001f0bd00190000006dd0000613d00000000048a034f0000000308b00210000000000a090433000000000a8a01cf000000000a8a022f000000000404043b0000010008800089000000000484022f00000000048401cf0000000004a4019f0000000000490435000006dd0000013d000000a0060000390000000009000019000000000a020019000007e50000013d0000000001cb001900000000000104350000001f01c000390000025001100197000000000b1b00190000000109900039000000000079004b000002b30000813d00000000012b0049000000200110008a000000200aa0003900000000001a04350000002005500039000000000c05043300000000410c0434000002370110019700000000011b043600000000040404330000000000610435000000a001b0003900000000fd0404340000000000d10435000000c00eb0003900000000000d004b000007fe0000613d00000000010000190000000004e1001900000000081f00190000000008080433000000000084043500000020011000390000000000d1004b000007f70000413d0000000001de001900000000000104350000001f01d00039000002500110019700000000011e00190000004004c0003900000000040404330000000008b10049000000400db0003900000000008d043500000000fe040434000000000de1043600000000000e004b000008140000613d00000000010000190000000004d1001900000000081f00190000000008080433000000000084043500000020011000390000000000e1004b0000080d0000413d0000000001ed001900000000000104350000006001c0003900000000010104330000006004b0003900000000001404350000001f01e00039000002500110019700000000011d00190000008004c0003900000000040404330000000008b10049000000800bb0003900000000008b043500000000dc040434000000000bc1043600000000000c004b000007dd0000613d00000000010000190000000004b1001900000000081d00190000000008080433000000000084043500000020011000390000000000c1004b000008270000413d000007dd0000013d0000001f0220003900000250022001970000000001120019000000000021004b00000000020000390000000102004039000002340010009c0000083b0000213d00000001002001900000083b0000c13d000000400010043f000000000001042d0000023f01000041000000000010043f0000004101000039000000040010043f00000240010000410000088e000104300000001f03100039000000000023004b0000000004000019000002490400404100000249052001970000024903300197000000000653013f000000000053004b00000000030000190000024903002041000002490060009c000000000304c019000000000003004b000008590000613d0000000003100367000000000303043b000002340030009c000008590000213d00000020011000390000000004310019000000000024004b000008590000213d0000000002030019000000000001042d00000000010000190000088e0001043000000000430104340000000001320436000000000003004b000008670000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b000008600000413d000000000231001900000000000204350000001f0230003900000250022001970000000001210019000000000001042d000002240010009c00000224010080410000004001100210000002240020009c00000224020080410000006002200210000000000112019f0000000002000414000002240020009c0000022402008041000000c002200210000000000112019f00000242011001c70000801002000039088c08870000040f0000000100200190000008800000613d000000000101043b000000000001042d00000000010000190000088e0001043000000885002104210000000102000039000000000001042d0000000002000019000000000001042d0000088a002104230000000102000039000000000001042d0000000002000019000000000001042d0000088c000004320000088d0001042e0000088e00010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff000000020000000000000000000000000000004000000100000000000000000000000000000000000000000000000000000000000000000000000000bfc9b78800000000000000000000000000000000000000000000000000000000c922362400000000000000000000000000000000000000000000000000000000c922362500000000000000000000000000000000000000000000000000000000e83eabba00000000000000000000000000000000000000000000000000000000e9d68a8e00000000000000000000000000000000000000000000000000000000bfc9b78900000000000000000000000000000000000000000000000000000000c1a5a35500000000000000000000000000000000000000000000000000000000c7c1cba10000000000000000000000000000000000000000000000000000000069600bc90000000000000000000000000000000000000000000000000000000069600bca000000000000000000000000000000000000000000000000000000009041be3d0000000000000000000000000000000000000000000000000000000093df2867000000000000000000000000000000000000000000000000000000003f4b04aa000000000000000000000000000000000000000000000000000000004bf78697000000000000000000000000000000000000000000000000ffffffffffffffff0200000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff7f000000000000000000000000ffffffffffffffffffffffffffffffffffffffff0000000000000000000000ff000000000000000000000000000000000000000002000000000000000000000000000000000000200000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000fffffffffffffeffffffff00000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000000000000004e487b710000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000020000000000000000000000000020000000000000000000000000000000000000000000000000000000000000005665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48bffffffffffffffffffffffffffffffffffffffffffffffff00000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000fffffffffffffedf000000000000000000000000000000000000000000000000ffffffffffffffbf00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e40000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000fffffffffffffdbf192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f320000000000000000000000000000000000000020000000800000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
