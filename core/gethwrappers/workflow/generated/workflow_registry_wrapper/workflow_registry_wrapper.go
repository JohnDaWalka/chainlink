// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package workflow_registry_wrapper

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

type WorkflowRegistryWorkflowMetadata struct {
	WorkflowID   [32]byte
	Owner        common.Address
	DonID        uint32
	Status       uint8
	WorkflowName string
	BinaryURL    string
	ConfigURL    string
	SecretsURL   string
}

var WorkflowRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AddressNotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotWorkflowOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"}],\"name\":\"DONNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWorkflowID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryLocked\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"maxAllowedLength\",\"type\":\"uint8\"}],\"name\":\"URLTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowAlreadyInDesiredStatus\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowContentNotUpdated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowIDAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"maxAllowedLength\",\"type\":\"uint8\"}],\"name\":\"WorkflowNameTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"donIDs\",\"type\":\"uint32[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AllowedDONsUpdatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AuthorizedAddressesUpdatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"lockedBy\",\"type\":\"address\"}],\"name\":\"RegistryLockedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"unlockedBy\",\"type\":\"address\"}],\"name\":\"RegistryUnlockedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowActivatedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowDeletedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"secretsURLHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowForceUpdateSecretsRequestedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"WorkflowPausedV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"WorkflowRegisteredV1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"oldWorkflowID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newWorkflowID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"WorkflowUpdatedV1\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"activateWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"computeHashKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"deleteWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedDONs\",\"outputs\":[{\"internalType\":\"uint32[]\",\"name\":\"allowedDONs\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedAddresses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"authorizedAddresses\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"}],\"name\":\"getWorkflowMetadata\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getWorkflowMetadataListByDON\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata[]\",\"name\":\"workflowMetadataList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getWorkflowMetadataListByOwner\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"internalType\":\"structWorkflowRegistry.WorkflowMetadata[]\",\"name\":\"workflowMetadataList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isRegistryLocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"}],\"name\":\"pauseWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"workflowName\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"workflowID\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"donID\",\"type\":\"uint32\"},{\"internalType\":\"enumWorkflowRegistry.WorkflowStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"registerWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"requestForceUpdateSecrets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"donIDs\",\"type\":\"uint32[]\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"updateAllowedDONs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"updateAuthorizedAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"workflowKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newWorkflowID\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"binaryURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"configURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"secretsURL\",\"type\":\"string\"}],\"name\":\"updateWorkflow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461004a57331561003b57600180546001600160a01b03191633179055600b805460ff1916905560405161344890816100508239f35b639b15e16f60e01b8152600490fd5b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806308e7f63a146120f0578063181f5a77146120615780632303348a14611f245780632b596f6d14611e965780633ccd14ff14611553578063695e1340146113835780636f35177114611298578063724c13dd146111a15780637497066b1461108657806379ba509714610fb05780637ec0846d14610f255780638da5cb5b14610ed35780639f4cb53414610eb2578063b87a019414610e5c578063d4b89c74146106a4578063db80009214610609578063e3dce080146104e2578063e690f33214610362578063f2fde38b14610284578063f794bdeb146101495763f99ecb6b1461010357600080fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457602060ff600b54166040519015158152f35b600080fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576007805461018581612472565b61019260405191826122f9565b81815261019e82612472565b916020937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401940136853760005b82811061023257505050906040519283926020840190602085525180915260408401929160005b82811061020557505050500390f35b835173ffffffffffffffffffffffffffffffffffffffff16855286955093810193928101926001016101f6565b6001908260005273ffffffffffffffffffffffffffffffffffffffff817fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c68801541661027d82876125a4565b52016101cf565b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576102bb6123e0565b6102c3612ba3565b73ffffffffffffffffffffffffffffffffffffffff8091169033821461033857817fffffffffffffffffffffffff00000000000000000000000000000000000000006000541617600055600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b60046040517fdad89dca000000000000000000000000000000000000000000000000000000008152fd5b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043560ff600b54166104b8576103a88133612e00565b90606082015160028110156104895760011461045f576000526004602052600160406000200178010000000000000000000000000000000000000000000000007fffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffff8254161790558051608063ffffffff604084015116920151907f6a0ed88e9cf3cb493ab4028fcb1dc7d18f0130fcdfba096edde0aadbfbf5e99f604051602081528061045a33956020830190612173565b0390a4005b60046040517f6f861db1000000000000000000000000000000000000000000000000000000008152fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60046040517f78a4e7d9000000000000000000000000000000000000000000000000000000008152fd5b34610144576104f036612368565b916104f9612ba3565b60ff600b54166104b85760005b828110610595575060405191806040840160408552526060830191906000905b80821061055d5785151560208601527f509460cccbb176edde6cac28895a4415a24961b8f3a0bd2617b9bb7b4e166c9b85850386a1005b90919283359073ffffffffffffffffffffffffffffffffffffffff821680920361014457600191815260208091019401920190610526565b60019084156105d7576105cf73ffffffffffffffffffffffffffffffffffffffff6105c96105c4848888612a43565b612b82565b16612ff1565b505b01610506565b61060373ffffffffffffffffffffffffffffffffffffffff6105fd6105c4848888612a43565b16613222565b506105d1565b346101445761062961061a36612403565b9161062361248a565b50612a64565b6000526004602052604060002073ffffffffffffffffffffffffffffffffffffffff6001820154161561067a57610662610676916126fa565b6040519182916020835260208301906121b6565b0390f35b60046040517f871e01b2000000000000000000000000000000000000000000000000000000008152fd5b346101445760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760443567ffffffffffffffff8111610144576106f390369060040161233a565b60643567ffffffffffffffff81116101445761071390369060040161233a565b92909160843567ffffffffffffffff81116101445761073690369060040161233a565b909360ff600b54166104b85761074d828785612c98565b61075960043533612e00565b9561076e3363ffffffff60408a015116612bee565b86519461078a60a08901516107843688856128b5565b90612e85565b61079d60c08a01516107843686886128b5565b6107b060e08b015161078436898d6128b5565b908280610e55575b80610e4e575b610e24578a926107cf602435612d58565b89600052600660205260406000207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0081541690556004356000526004602052604060002090602435825515610cc9575b5015610b74575b156108b7575b7f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad7353959361045a9361089961088b9794608063ffffffff60406108a7970151169d015197604051998a996024358b5260a060208c015260a08b0190612173565b9189830360408b0152612978565b918683036060880152612978565b9083820360808501523397612978565b60e001518051610b10575b506004356000526004602052600560406000200167ffffffffffffffff8511610ae15787816108fc876108f68d95546125f3565b83612931565b6000601f88116001146109d9578794608063ffffffff60406108a7977f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad73539d9b9761045a9b9761098661088b9f806108999a6000916109ce575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b8a806109a5575b5050975050505094975050935093955061082c565b6109af9133612a64565b60005260056020526109c660043583600020613043565b508f8a610990565b905084013538610955565b8160005260206000209060005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08a168110610ac4575094608063ffffffff60406108a7977f41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad73539d9b9761045a9b976108999761088b9f9c8d807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610610a8b575b905060018092501b019055610989565b60f87fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9160031b161c199084013516905538808d610a7b565b9381013582558c94508b93600190920191602091820191016109e6565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610b5281610b26602082019433866129b7565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018352826122f9565b5190206000526005602052610b6d60043560406000206132e9565b50886108c2565b9050600435600052600460205260046040600020019067ffffffffffffffff8411610ae15781610ba9856108f68d95546125f3565b846000601f8211600114610c07579080610bf592600091610bfc57507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055610826565b90508801358e610955565b60008381526020812092505b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe088168110610cae5750867fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610610c76575b5050600185811b019055610826565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19908801351690558b80610c67565b8882013583558d955060019092019160209182019101610c13565b925067ffffffffffffffff8811610ae15782610cf689610ced60038f9701546125f3565b60038401612931565b886000601f8211600114610d595781600392610d4492600091610d4e57507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9101555b8b61081f565b905088013538610955565b90506003820160005260206000209060005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08c168110610e0957508a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610610dd1575b5050600360018a811b01910155610d48565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88c60031b161c19908701351690558c80610dbf565b8782013583558e965060019092019160209182019101610d6b565b60046040517f6b4a810d000000000000000000000000000000000000000000000000000000008152fd5b50816107be565b50806107b8565b346101445760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457610676610ea6610e996123e0565b6044359060243590612abf565b6040519182918261225a565b34610144576020610ecb610ec536612403565b91612a64565b604051908152f35b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457610f5c612ba3565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00600b5416600b557f11a03e25ee25bf1459f9e1cb293ea03707d84917f54a65e32c9a7be2f2edd68a6020604051338152a1005b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760005473ffffffffffffffffffffffffffffffffffffffff808216330361105c57600154917fffffffffffffffffffffffff0000000000000000000000000000000000000000903382851617600155166000553391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b60046040517f02b543c6000000000000000000000000000000000000000000000000000000008152fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457600980546110c281612472565b6110cf60405191826122f9565b8181526110db82612472565b916020937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401940136853760005b82811061115f57505050906040519283926020840190602085525180915260408401929160005b82811061114257505050500390f35b835163ffffffff1685528695509381019392810192600101611133565b6001908260005263ffffffff817f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af01541661119a82876125a4565b520161110c565b34610144576111af36612368565b916111b8612ba3565b60ff600b54166104b85760005b828110611244575060405191806040840160408552526060830191906000905b80821061121c5785151560208601527fcab63bf31d1e656baa23cebef64e12033ea0ffbd44b1278c3747beec2d2f618c85850386a1005b90919283359063ffffffff8216809203610144576001918152602080910194019201906111e5565b60019084156112765761126e63ffffffff611268611263848888612a43565b612a53565b16612f38565b505b016111c5565b61129263ffffffff61128c611263848888612a43565b166130cf565b50611270565b346101445760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043560ff600b54166104b8576112de8133612e00565b9060608201516002811015610489571561045f576080604083019163ffffffff9061130c3383865116612bee565b600052600460205260016040600020017fffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffff81541690558351925116920151907f17b2d730bb5e064df3fbc6165c8aceb3b0d62c524c196c0bc1012209280bc9a6604051602081528061045a33956020830190612173565b34610144576020807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610144576004359060ff600b54166104b8576113cb8233612e00565b916113e3336000526008602052604060002054151590565b156115235782600493516000526006835260406000207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00815416905533600052600283526114358260406000206132e9565b5060408101608063ffffffff808351166000526003865261145a8560406000206132e9565b5060e084015180516114ef575b508351925116920151907f76ee2dfcae10cb8522e62e713e62660e09ecfaab08db15d9404de19141322571604051868152806114a7339589830190612173565b0390a4600052526114ed6005604060002060008155600060018201556114cf600282016129fa565b6114db600382016129fa565b6114e7600482016129fa565b016129fa565b005b60405161150481610b268a82019433866129b7565b5190206000526005865261151c8560406000206132e9565b5087611467565b60246040517f85982a00000000000000000000000000000000000000000000000000000000008152336004820152fd5b346101445760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043567ffffffffffffffff8111610144576115a290369060040161233a565b6044359163ffffffff8316830361014457600260643510156101445760843567ffffffffffffffff8111610144576115de90369060040161233a565b91909260a43567ffffffffffffffff81116101445761160190369060040161233a565b60c43567ffffffffffffffff81116101445761162190369060040161233a565b96909560ff600b54166104b857611638338a612bee565b60408511611e5e5761164b888483612c98565b611656858733612a64565b80600052600460205273ffffffffffffffffffffffffffffffffffffffff60016040600020015416611e345761168d602435612d58565b6040519061169a826122dc565b602435825233602083015263ffffffff8b1660408301526116c0606435606084016125e7565b6116cb36888a6128b5565b60808301526116db3684866128b5565b60a08301526116eb3686886128b5565b60c08301526116fb368b8b6128b5565b60e0830152806000526004602052604060002091805183556001830173ffffffffffffffffffffffffffffffffffffffff60208301511681549077ffffffff0000000000000000000000000000000000000000604085015160a01b1690606085015160028110156104895778ff0000000000000000000000000000000000000000000000007fffffffffffffff000000000000000000000000000000000000000000000000009160c01b1693161717179055608081015180519067ffffffffffffffff8211610ae1576117de826117d560028801546125f3565b60028801612931565b602090601f8311600114611d685761182c929160009183611c91575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028401555b60a081015180519067ffffffffffffffff8211610ae1576118638261185a60038801546125f3565b60038801612931565b602090601f8311600114611c9c576118b0929160009183611c915750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038401555b60c081015180519067ffffffffffffffff8211610ae1576118e7826118de60048801546125f3565b60048801612931565b602090601f8311600114611bc45791806119389260e09594600092611a9f5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60048501555b015180519267ffffffffffffffff8411610ae157838d926119708e9661196760058601546125f3565b60058601612931565b602090601f8311600114611aaa579463ffffffff61089995819a957fc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb9f9e9d99946119fc8761045a9f9b98600593611a609f9a600092611a9f5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9101555b336000526002602052611a17836040600020613043565b50166000526003602052611a2f816040600020613043565b508d82611a76575b50505061088b6040519a8b9a611a4f8c606435612143565b60a060208d015260a08c0191612978565b9783890360808501521696339660243596612978565b611a9692611a849133612a64565b60005260056020526040600020613043565b508c8f8d611a37565b0151905038806117fa565b906005840160005260206000209160005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085168110611b9a57506108999563ffffffff9a957fc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb9f9e9d999460018761045a9f9b96928f9693611a609f9a94837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06005971610611b63575b505050811b01910155611a00565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055388080611b55565b939550918194969750600160209291839285015181550194019201918f9492918f97969492611abb565b906004860160005260206000209160005b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085168110611c795750918391600193837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060e098971610611c42575b505050811b01600485015561193e565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558f8080611c32565b91926020600181928685015181550194019201611bd5565b015190508f806117fa565b9190600386016000526020600020906000935b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168510611d4d5760019450837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610611d16575b505050811b0160038401556118b6565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080611d06565b81810151835560209485019460019093019290910190611caf565b9190600286016000526020600020906000935b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168510611e195760019450837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0811610611de2575b505050811b016002840155611832565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080611dd2565b81810151835560209485019460019093019290910190611d7b565b60046040517fa0677dd0000000000000000000000000000000000000000000000000000000008152fd5b604485604051907f36a7c503000000000000000000000000000000000000000000000000000000008252600482015260406024820152fd5b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457611ecd612ba3565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00600b541617600b557f2789711f6fd67d131ad68378617b5d1d21a2c92b34d7c3745d70b3957c08096c6020604051338152a1005b34610144576020807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043567ffffffffffffffff811161014457611f7490369060040161233a565b60ff600b54166104b857611f889133612a64565b90816000526005602052604060002091825491821561067a5760005b838110611fad57005b80611fba60019287612f20565b90549060031b1c60005260048352604060002063ffffffff8382015460a01c16600052600a8452604060002054151580612044575b611ffb575b5001611fa4565b7f95d94f817db4971aa99ba35d0fe019bd8cc39866fbe02b6d47b5f0f3727fb6736040518681526040868201528061203b33946002604084019101612646565b0390a286611ff4565b5061205c336000526008602052604060002054151590565b611fef565b346101445760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014457604051604081019080821067ffffffffffffffff831117610ae15761067691604052601a81527f576f726b666c6f77526567697374727920312e302e302d6465760000000000006020820152604051918291602083526020830190612173565b346101445760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101445760043563ffffffff8116810361014457610ea66106769160443590602435906127c0565b9060028210156104895752565b60005b8381106121635750506000910152565b8181015183820152602001612153565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f6020936121af81518092818752878088019101612150565b0116010190565b6122579160e06122466122346122226101008651865273ffffffffffffffffffffffffffffffffffffffff602088015116602087015263ffffffff604088015116604087015261220e60608801516060880190612143565b608087015190806080880152860190612173565b60a086015185820360a0870152612173565b60c085015184820360c0860152612173565b9201519060e0818403910152612173565b90565b6020808201906020835283518092526040830192602060408460051b8301019501936000915b8483106122905750505050505090565b90919293949584806122cc837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc086600196030187528a516121b6565b9801930193019194939290612280565b610100810190811067ffffffffffffffff821117610ae157604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117610ae157604052565b9181601f840112156101445782359167ffffffffffffffff8311610144576020838186019501011161014457565b9060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126101445760043567ffffffffffffffff9283821161014457806023830112156101445781600401359384116101445760248460051b8301011161014457602401919060243580151581036101445790565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361014457565b9060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126101445760043573ffffffffffffffffffffffffffffffffffffffff8116810361014457916024359067ffffffffffffffff82116101445761246e9160040161233a565b9091565b67ffffffffffffffff8111610ae15760051b60200190565b60405190612497826122dc565b606060e0836000815260006020820152600060408201526000838201528260808201528260a08201528260c08201520152565b6040516020810181811067ffffffffffffffff821117610ae1576040526000815290565b906124f882612472565b61250560405191826122f9565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06125338294612472565b019060005b82811061254457505050565b60209061254f61248a565b82828501015201612538565b9190820180921161256857565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9190820391821161256857565b80518210156125b85760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60028210156104895752565b90600182811c9216801561263c575b602083101461260d57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691612602565b805460009392612655826125f3565b918282526020936001916001811690816000146126bd575060011461267c575b5050505050565b90939495506000929192528360002092846000945b8386106126a957505050500101903880808080612675565b805485870183015294019385908201612691565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168685015250505090151560051b010191503880808080612675565b90600560e06040936127bc855191612711836122dc565b6127b583978254855261275b60ff600185015473ffffffffffffffffffffffffffffffffffffffff8116602089015263ffffffff8160a01c168489015260c01c16606087016125e7565b80516127758161276e8160028801612646565b03826122f9565b6080860152805161278d8161276e8160038801612646565b60a086015280516127a58161276e8160048801612646565b60c0860152518096819301612646565b03846122f9565b0152565b63ffffffff1691600083815260036020906003602052604093604084205490818710156128a55761281491816064899311801561289d575b612895575b81612808828561255b565b11156128855750612597565b9461281e866124ee565b96845b87811061283357505050505050505090565b60019082875284865261285288882061284c838761255b565b90612f20565b905490861b1c8752600486526128698888206126fa565b612873828c6125a4565b5261287e818b6125a4565b5001612821565b61289091508261255b565b612597565b5060646127fd565b5080156127f8565b50505050505050506122576124ca565b92919267ffffffffffffffff8211610ae157604051916128fd60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601846122f9565b829481845281830111610144578281602093846000960137010152565b818110612925575050565b6000815560010161291a565b9190601f811161294057505050565b61296c926000526020600020906020601f840160051c8301931061296e575b601f0160051c019061291a565b565b909150819061295f565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000006014939260601b1681526129f58251809360208685019101612150565b010190565b612a0481546125f3565b9081612a0e575050565b81601f60009311600114612a20575055565b908083918252612a3f601f60208420940160051c84016001850161291a565b5555565b91908110156125b85760051b0190565b3563ffffffff811681036101445790565b91906034612ab991836040519485927fffffffffffffffffffffffffffffffffffffffff000000000000000000000000602085019860601b1688528484013781016000838201520360148101845201826122f9565b51902090565b73ffffffffffffffffffffffffffffffffffffffff169160008381526002926020906002602052604093604084205490818310156128a557612b1691816064859311801561289d576128955781612808828561255b565b94612b20866124ee565b96845b878110612b3557505050505050505090565b600190828752838652612b4e88882061284c838861255b565b90549060031b1c875260048652612b668888206126fa565b612b70828c6125a4565b52612b7b818b6125a4565b5001612b23565b3573ffffffffffffffffffffffffffffffffffffffff811681036101445790565b73ffffffffffffffffffffffffffffffffffffffff600154163303612bc457565b60046040517f2b5c74de000000000000000000000000000000000000000000000000000000008152fd5b63ffffffff1680600052600a60205260406000205415612c67575073ffffffffffffffffffffffffffffffffffffffff1680600052600860205260406000205415612c365750565b602490604051907f85982a000000000000000000000000000000000000000000000000000000000082526004820152fd5b602490604051907f8fe6d7e10000000000000000000000000000000000000000000000000000000082526004820152fd5b9060c891828111612d225750818111612ced5750808211612cb7575050565b60449250604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b604491604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b60449083604051917ecd56a800000000000000000000000000000000000000000000000000000000835260048301526024820152fd5b8015612dd65780600052600660205260ff60406000205416612dac576000526006602052604060002060017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055565b60046040517f4cb050e4000000000000000000000000000000000000000000000000000000008152fd5b60046040517f7dc2f4e1000000000000000000000000000000000000000000000000000000008152fd5b90612e0961248a565b506000526004602052612e1f60406000206126fa565b906020820173ffffffffffffffffffffffffffffffffffffffff90818151161561067a57819051169116809103612e54575090565b602490604051907f31ee6dc70000000000000000000000000000000000000000000000000000000082526004820152fd5b9081518151908181149384612e9c575b5050505090565b6020929394508201209201201438808080612e95565b6009548110156125b85760096000527f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af0190600090565b6007548110156125b85760076000527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c6880190600090565b80548210156125b85760005260206000200190600090565b6000818152600a6020526040812054612fec5760095468010000000000000000811015612fbf579082612fab612f7684600160409601600955612eb2565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b9055600954928152600a6020522055600190565b6024827f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b905090565b600081815260086020526040812054612fec5760075468010000000000000000811015612fbf57908261302f612f7684600160409601600755612ee9565b905560075492815260086020522055600190565b919060018301600090828252806020526040822054156000146130c9578454946801000000000000000086101561309c578361308c612f76886001604098999a01855584612f20565b9055549382526020522055600190565b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b50925050565b6000818152600a6020526040812054909190801561321d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff908181018181116131f057600954908382019182116131c35781810361318f575b50505060095480156131625781019061314182612eb2565b909182549160031b1b191690556009558152600a6020526040812055600190565b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526031600452fd5b6131ad61319e612f7693612eb2565b90549060031b1c928392612eb2565b90558452600a6020526040842055388080613129565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b505090565b600081815260086020526040812054909190801561321d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff908181018181116131f057600754908382019182116131c3578181036132b5575b50505060075480156131625781019061329482612ee9565b909182549160031b1b19169055600755815260086020526040812055600190565b6132d36132c4612f7693612ee9565b90549060031b1c928392612ee9565b905584526008602052604084205538808061327c565b90600182019060009281845282602052604084205490811515600014612e95577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9182810181811161340e578254908482019182116133e1578181036133ac575b5050508054801561337f578201916133628383612f20565b909182549160031b1b191690555582526020526040812055600190565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526031600452fd5b6133cc6133bc612f769386612f20565b90549060031b1c92839286612f20565b9055865284602052604086205538808061334a565b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fdfea164736f6c6343000818000a",
}

var WorkflowRegistryABI = WorkflowRegistryMetaData.ABI

var WorkflowRegistryBin = WorkflowRegistryMetaData.Bin

func DeployWorkflowRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WorkflowRegistry, error) {
	parsed, err := WorkflowRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WorkflowRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WorkflowRegistry{address: address, abi: *parsed, WorkflowRegistryCaller: WorkflowRegistryCaller{contract: contract}, WorkflowRegistryTransactor: WorkflowRegistryTransactor{contract: contract}, WorkflowRegistryFilterer: WorkflowRegistryFilterer{contract: contract}}, nil
}

type WorkflowRegistry struct {
	address common.Address
	abi     abi.ABI
	WorkflowRegistryCaller
	WorkflowRegistryTransactor
	WorkflowRegistryFilterer
}

type WorkflowRegistryCaller struct {
	contract *bind.BoundContract
}

type WorkflowRegistryTransactor struct {
	contract *bind.BoundContract
}

type WorkflowRegistryFilterer struct {
	contract *bind.BoundContract
}

type WorkflowRegistrySession struct {
	Contract     *WorkflowRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type WorkflowRegistryCallerSession struct {
	Contract *WorkflowRegistryCaller
	CallOpts bind.CallOpts
}

type WorkflowRegistryTransactorSession struct {
	Contract     *WorkflowRegistryTransactor
	TransactOpts bind.TransactOpts
}

type WorkflowRegistryRaw struct {
	Contract *WorkflowRegistry
}

type WorkflowRegistryCallerRaw struct {
	Contract *WorkflowRegistryCaller
}

type WorkflowRegistryTransactorRaw struct {
	Contract *WorkflowRegistryTransactor
}

func NewWorkflowRegistry(address common.Address, backend bind.ContractBackend) (*WorkflowRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(WorkflowRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindWorkflowRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistry{address: address, abi: abi, WorkflowRegistryCaller: WorkflowRegistryCaller{contract: contract}, WorkflowRegistryTransactor: WorkflowRegistryTransactor{contract: contract}, WorkflowRegistryFilterer: WorkflowRegistryFilterer{contract: contract}}, nil
}

func NewWorkflowRegistryCaller(address common.Address, caller bind.ContractCaller) (*WorkflowRegistryCaller, error) {
	contract, err := bindWorkflowRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryCaller{contract: contract}, nil
}

func NewWorkflowRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*WorkflowRegistryTransactor, error) {
	contract, err := bindWorkflowRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryTransactor{contract: contract}, nil
}

func NewWorkflowRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*WorkflowRegistryFilterer, error) {
	contract, err := bindWorkflowRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryFilterer{contract: contract}, nil
}

func bindWorkflowRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WorkflowRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_WorkflowRegistry *WorkflowRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorkflowRegistry.Contract.WorkflowRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.WorkflowRegistryTransactor.contract.Transfer(opts)
}

func (_WorkflowRegistry *WorkflowRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.WorkflowRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorkflowRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.contract.Transfer(opts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) ComputeHashKey(opts *bind.CallOpts, owner common.Address, field string) ([32]byte, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "computeHashKey", owner, field)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) ComputeHashKey(owner common.Address, field string) ([32]byte, error) {
	return _WorkflowRegistry.Contract.ComputeHashKey(&_WorkflowRegistry.CallOpts, owner, field)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) ComputeHashKey(owner common.Address, field string) ([32]byte, error) {
	return _WorkflowRegistry.Contract.ComputeHashKey(&_WorkflowRegistry.CallOpts, owner, field)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetAllAllowedDONs(opts *bind.CallOpts) ([]uint32, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getAllAllowedDONs")

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetAllAllowedDONs() ([]uint32, error) {
	return _WorkflowRegistry.Contract.GetAllAllowedDONs(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetAllAllowedDONs() ([]uint32, error) {
	return _WorkflowRegistry.Contract.GetAllAllowedDONs(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetAllAuthorizedAddresses(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getAllAuthorizedAddresses")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetAllAuthorizedAddresses() ([]common.Address, error) {
	return _WorkflowRegistry.Contract.GetAllAuthorizedAddresses(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetAllAuthorizedAddresses() ([]common.Address, error) {
	return _WorkflowRegistry.Contract.GetAllAuthorizedAddresses(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetWorkflowMetadata(opts *bind.CallOpts, workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getWorkflowMetadata", workflowOwner, workflowName)

	if err != nil {
		return *new(WorkflowRegistryWorkflowMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(WorkflowRegistryWorkflowMetadata)).(*WorkflowRegistryWorkflowMetadata)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetWorkflowMetadata(workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadata(&_WorkflowRegistry.CallOpts, workflowOwner, workflowName)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetWorkflowMetadata(workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadata(&_WorkflowRegistry.CallOpts, workflowOwner, workflowName)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetWorkflowMetadataListByDON(opts *bind.CallOpts, donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getWorkflowMetadataListByDON", donID, start, limit)

	if err != nil {
		return *new([]WorkflowRegistryWorkflowMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new([]WorkflowRegistryWorkflowMetadata)).(*[]WorkflowRegistryWorkflowMetadata)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetWorkflowMetadataListByDON(donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByDON(&_WorkflowRegistry.CallOpts, donID, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetWorkflowMetadataListByDON(donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByDON(&_WorkflowRegistry.CallOpts, donID, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) GetWorkflowMetadataListByOwner(opts *bind.CallOpts, workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "getWorkflowMetadataListByOwner", workflowOwner, start, limit)

	if err != nil {
		return *new([]WorkflowRegistryWorkflowMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new([]WorkflowRegistryWorkflowMetadata)).(*[]WorkflowRegistryWorkflowMetadata)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) GetWorkflowMetadataListByOwner(workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByOwner(&_WorkflowRegistry.CallOpts, workflowOwner, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) GetWorkflowMetadataListByOwner(workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error) {
	return _WorkflowRegistry.Contract.GetWorkflowMetadataListByOwner(&_WorkflowRegistry.CallOpts, workflowOwner, start, limit)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) IsRegistryLocked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "isRegistryLocked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) IsRegistryLocked() (bool, error) {
	return _WorkflowRegistry.Contract.IsRegistryLocked(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) IsRegistryLocked() (bool, error) {
	return _WorkflowRegistry.Contract.IsRegistryLocked(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) Owner() (common.Address, error) {
	return _WorkflowRegistry.Contract.Owner(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) Owner() (common.Address, error) {
	return _WorkflowRegistry.Contract.Owner(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WorkflowRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WorkflowRegistry *WorkflowRegistrySession) TypeAndVersion() (string, error) {
	return _WorkflowRegistry.Contract.TypeAndVersion(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryCallerSession) TypeAndVersion() (string, error) {
	return _WorkflowRegistry.Contract.TypeAndVersion(&_WorkflowRegistry.CallOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_WorkflowRegistry *WorkflowRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.AcceptOwnership(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.AcceptOwnership(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) ActivateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "activateWorkflow", workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistrySession) ActivateWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.ActivateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) ActivateWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.ActivateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) DeleteWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "deleteWorkflow", workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistrySession) DeleteWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.DeleteWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) DeleteWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.DeleteWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) LockRegistry(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "lockRegistry")
}

func (_WorkflowRegistry *WorkflowRegistrySession) LockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.LockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) LockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.LockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) PauseWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "pauseWorkflow", workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistrySession) PauseWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.PauseWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) PauseWorkflow(workflowKey [32]byte) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.PauseWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) RegisterWorkflow(opts *bind.TransactOpts, workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "registerWorkflow", workflowName, workflowID, donID, status, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistrySession) RegisterWorkflow(workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RegisterWorkflow(&_WorkflowRegistry.TransactOpts, workflowName, workflowID, donID, status, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) RegisterWorkflow(workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RegisterWorkflow(&_WorkflowRegistry.TransactOpts, workflowName, workflowID, donID, status, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) RequestForceUpdateSecrets(opts *bind.TransactOpts, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "requestForceUpdateSecrets", secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistrySession) RequestForceUpdateSecrets(secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RequestForceUpdateSecrets(&_WorkflowRegistry.TransactOpts, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) RequestForceUpdateSecrets(secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.RequestForceUpdateSecrets(&_WorkflowRegistry.TransactOpts, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_WorkflowRegistry *WorkflowRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.TransferOwnership(&_WorkflowRegistry.TransactOpts, to)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.TransferOwnership(&_WorkflowRegistry.TransactOpts, to)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UnlockRegistry(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "unlockRegistry")
}

func (_WorkflowRegistry *WorkflowRegistrySession) UnlockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UnlockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UnlockRegistry() (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UnlockRegistry(&_WorkflowRegistry.TransactOpts)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UpdateAllowedDONs(opts *bind.TransactOpts, donIDs []uint32, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "updateAllowedDONs", donIDs, allowed)
}

func (_WorkflowRegistry *WorkflowRegistrySession) UpdateAllowedDONs(donIDs []uint32, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAllowedDONs(&_WorkflowRegistry.TransactOpts, donIDs, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UpdateAllowedDONs(donIDs []uint32, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAllowedDONs(&_WorkflowRegistry.TransactOpts, donIDs, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UpdateAuthorizedAddresses(opts *bind.TransactOpts, addresses []common.Address, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "updateAuthorizedAddresses", addresses, allowed)
}

func (_WorkflowRegistry *WorkflowRegistrySession) UpdateAuthorizedAddresses(addresses []common.Address, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAuthorizedAddresses(&_WorkflowRegistry.TransactOpts, addresses, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UpdateAuthorizedAddresses(addresses []common.Address, allowed bool) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateAuthorizedAddresses(&_WorkflowRegistry.TransactOpts, addresses, allowed)
}

func (_WorkflowRegistry *WorkflowRegistryTransactor) UpdateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.contract.Transact(opts, "updateWorkflow", workflowKey, newWorkflowID, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistrySession) UpdateWorkflow(workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey, newWorkflowID, binaryURL, configURL, secretsURL)
}

func (_WorkflowRegistry *WorkflowRegistryTransactorSession) UpdateWorkflow(workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error) {
	return _WorkflowRegistry.Contract.UpdateWorkflow(&_WorkflowRegistry.TransactOpts, workflowKey, newWorkflowID, binaryURL, configURL, secretsURL)
}

type WorkflowRegistryAllowedDONsUpdatedV1Iterator struct {
	Event *WorkflowRegistryAllowedDONsUpdatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryAllowedDONsUpdatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryAllowedDONsUpdatedV1)
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
		it.Event = new(WorkflowRegistryAllowedDONsUpdatedV1)
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

func (it *WorkflowRegistryAllowedDONsUpdatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryAllowedDONsUpdatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryAllowedDONsUpdatedV1 struct {
	DonIDs  []uint32
	Allowed bool
	Raw     types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterAllowedDONsUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAllowedDONsUpdatedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "AllowedDONsUpdatedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryAllowedDONsUpdatedV1Iterator{contract: _WorkflowRegistry.contract, event: "AllowedDONsUpdatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchAllowedDONsUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAllowedDONsUpdatedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "AllowedDONsUpdatedV1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryAllowedDONsUpdatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "AllowedDONsUpdatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseAllowedDONsUpdatedV1(log types.Log) (*WorkflowRegistryAllowedDONsUpdatedV1, error) {
	event := new(WorkflowRegistryAllowedDONsUpdatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "AllowedDONsUpdatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator struct {
	Event *WorkflowRegistryAuthorizedAddressesUpdatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
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
		it.Event = new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
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

func (it *WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryAuthorizedAddressesUpdatedV1 struct {
	Addresses []common.Address
	Allowed   bool
	Raw       types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterAuthorizedAddressesUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "AuthorizedAddressesUpdatedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator{contract: _WorkflowRegistry.contract, event: "AuthorizedAddressesUpdatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchAuthorizedAddressesUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAuthorizedAddressesUpdatedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "AuthorizedAddressesUpdatedV1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "AuthorizedAddressesUpdatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseAuthorizedAddressesUpdatedV1(log types.Log) (*WorkflowRegistryAuthorizedAddressesUpdatedV1, error) {
	event := new(WorkflowRegistryAuthorizedAddressesUpdatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "AuthorizedAddressesUpdatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryOwnershipTransferRequestedIterator struct {
	Event *WorkflowRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryOwnershipTransferRequested)
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
		it.Event = new(WorkflowRegistryOwnershipTransferRequested)
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

func (it *WorkflowRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryOwnershipTransferRequestedIterator{contract: _WorkflowRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryOwnershipTransferRequested)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*WorkflowRegistryOwnershipTransferRequested, error) {
	event := new(WorkflowRegistryOwnershipTransferRequested)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryOwnershipTransferredIterator struct {
	Event *WorkflowRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryOwnershipTransferred)
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
		it.Event = new(WorkflowRegistryOwnershipTransferred)
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

func (it *WorkflowRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryOwnershipTransferredIterator{contract: _WorkflowRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryOwnershipTransferred)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*WorkflowRegistryOwnershipTransferred, error) {
	event := new(WorkflowRegistryOwnershipTransferred)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryRegistryLockedV1Iterator struct {
	Event *WorkflowRegistryRegistryLockedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryRegistryLockedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryRegistryLockedV1)
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
		it.Event = new(WorkflowRegistryRegistryLockedV1)
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

func (it *WorkflowRegistryRegistryLockedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryRegistryLockedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryRegistryLockedV1 struct {
	LockedBy common.Address
	Raw      types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterRegistryLockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryLockedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "RegistryLockedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryRegistryLockedV1Iterator{contract: _WorkflowRegistry.contract, event: "RegistryLockedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchRegistryLockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryLockedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "RegistryLockedV1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryRegistryLockedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryLockedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseRegistryLockedV1(log types.Log) (*WorkflowRegistryRegistryLockedV1, error) {
	event := new(WorkflowRegistryRegistryLockedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryLockedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryRegistryUnlockedV1Iterator struct {
	Event *WorkflowRegistryRegistryUnlockedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryRegistryUnlockedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryRegistryUnlockedV1)
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
		it.Event = new(WorkflowRegistryRegistryUnlockedV1)
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

func (it *WorkflowRegistryRegistryUnlockedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryRegistryUnlockedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryRegistryUnlockedV1 struct {
	UnlockedBy common.Address
	Raw        types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterRegistryUnlockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryUnlockedV1Iterator, error) {

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "RegistryUnlockedV1")
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryRegistryUnlockedV1Iterator{contract: _WorkflowRegistry.contract, event: "RegistryUnlockedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchRegistryUnlockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryUnlockedV1) (event.Subscription, error) {

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "RegistryUnlockedV1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryRegistryUnlockedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryUnlockedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseRegistryUnlockedV1(log types.Log) (*WorkflowRegistryRegistryUnlockedV1, error) {
	event := new(WorkflowRegistryRegistryUnlockedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "RegistryUnlockedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowActivatedV1Iterator struct {
	Event *WorkflowRegistryWorkflowActivatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowActivatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowActivatedV1)
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
		it.Event = new(WorkflowRegistryWorkflowActivatedV1)
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

func (it *WorkflowRegistryWorkflowActivatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowActivatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowActivatedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	WorkflowName  string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowActivatedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowActivatedV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowActivatedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowActivatedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowActivatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowActivatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowActivatedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowActivatedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowActivatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowActivatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowActivatedV1(log types.Log) (*WorkflowRegistryWorkflowActivatedV1, error) {
	event := new(WorkflowRegistryWorkflowActivatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowActivatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowDeletedV1Iterator struct {
	Event *WorkflowRegistryWorkflowDeletedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowDeletedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowDeletedV1)
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
		it.Event = new(WorkflowRegistryWorkflowDeletedV1)
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

func (it *WorkflowRegistryWorkflowDeletedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowDeletedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowDeletedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	WorkflowName  string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowDeletedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowDeletedV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowDeletedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowDeletedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowDeletedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowDeletedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowDeletedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowDeletedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowDeletedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowDeletedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowDeletedV1(log types.Log) (*WorkflowRegistryWorkflowDeletedV1, error) {
	event := new(WorkflowRegistryWorkflowDeletedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowDeletedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator struct {
	Event *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
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
		it.Event = new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
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

func (it *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1 struct {
	Owner          common.Address
	SecretsURLHash [32]byte
	WorkflowName   string
	Raw            types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowForceUpdateSecretsRequestedV1(opts *bind.FilterOpts, owner []common.Address) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowForceUpdateSecretsRequestedV1", ownerRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowForceUpdateSecretsRequestedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowForceUpdateSecretsRequestedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowForceUpdateSecretsRequestedV1", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowForceUpdateSecretsRequestedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowForceUpdateSecretsRequestedV1(log types.Log) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, error) {
	event := new(WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowForceUpdateSecretsRequestedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowPausedV1Iterator struct {
	Event *WorkflowRegistryWorkflowPausedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowPausedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowPausedV1)
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
		it.Event = new(WorkflowRegistryWorkflowPausedV1)
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

func (it *WorkflowRegistryWorkflowPausedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowPausedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowPausedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	WorkflowName  string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowPausedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowPausedV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowPausedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowPausedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowPausedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowPausedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowPausedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowPausedV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowPausedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowPausedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowPausedV1(log types.Log) (*WorkflowRegistryWorkflowPausedV1, error) {
	event := new(WorkflowRegistryWorkflowPausedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowPausedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowRegisteredV1Iterator struct {
	Event *WorkflowRegistryWorkflowRegisteredV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowRegisteredV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowRegisteredV1)
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
		it.Event = new(WorkflowRegistryWorkflowRegisteredV1)
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

func (it *WorkflowRegistryWorkflowRegisteredV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowRegisteredV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowRegisteredV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	Status        uint8
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowRegisteredV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowRegisteredV1Iterator, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowRegisteredV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowRegisteredV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowRegisteredV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowRegisteredV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowRegisteredV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var workflowIDRule []interface{}
	for _, workflowIDItem := range workflowID {
		workflowIDRule = append(workflowIDRule, workflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowRegisteredV1", workflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowRegisteredV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowRegisteredV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowRegisteredV1(log types.Log) (*WorkflowRegistryWorkflowRegisteredV1, error) {
	event := new(WorkflowRegistryWorkflowRegisteredV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowRegisteredV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WorkflowRegistryWorkflowUpdatedV1Iterator struct {
	Event *WorkflowRegistryWorkflowUpdatedV1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WorkflowRegistryWorkflowUpdatedV1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorkflowRegistryWorkflowUpdatedV1)
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
		it.Event = new(WorkflowRegistryWorkflowUpdatedV1)
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

func (it *WorkflowRegistryWorkflowUpdatedV1Iterator) Error() error {
	return it.fail
}

func (it *WorkflowRegistryWorkflowUpdatedV1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WorkflowRegistryWorkflowUpdatedV1 struct {
	OldWorkflowID [32]byte
	WorkflowOwner common.Address
	DonID         uint32
	NewWorkflowID [32]byte
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
	Raw           types.Log
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) FilterWorkflowUpdatedV1(opts *bind.FilterOpts, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowUpdatedV1Iterator, error) {

	var oldWorkflowIDRule []interface{}
	for _, oldWorkflowIDItem := range oldWorkflowID {
		oldWorkflowIDRule = append(oldWorkflowIDRule, oldWorkflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.FilterLogs(opts, "WorkflowUpdatedV1", oldWorkflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryWorkflowUpdatedV1Iterator{contract: _WorkflowRegistry.contract, event: "WorkflowUpdatedV1", logs: logs, sub: sub}, nil
}

func (_WorkflowRegistry *WorkflowRegistryFilterer) WatchWorkflowUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowUpdatedV1, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error) {

	var oldWorkflowIDRule []interface{}
	for _, oldWorkflowIDItem := range oldWorkflowID {
		oldWorkflowIDRule = append(oldWorkflowIDRule, oldWorkflowIDItem)
	}
	var workflowOwnerRule []interface{}
	for _, workflowOwnerItem := range workflowOwner {
		workflowOwnerRule = append(workflowOwnerRule, workflowOwnerItem)
	}
	var donIDRule []interface{}
	for _, donIDItem := range donID {
		donIDRule = append(donIDRule, donIDItem)
	}

	logs, sub, err := _WorkflowRegistry.contract.WatchLogs(opts, "WorkflowUpdatedV1", oldWorkflowIDRule, workflowOwnerRule, donIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WorkflowRegistryWorkflowUpdatedV1)
				if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowUpdatedV1", log); err != nil {
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

func (_WorkflowRegistry *WorkflowRegistryFilterer) ParseWorkflowUpdatedV1(log types.Log) (*WorkflowRegistryWorkflowUpdatedV1, error) {
	event := new(WorkflowRegistryWorkflowUpdatedV1)
	if err := _WorkflowRegistry.contract.UnpackLog(event, "WorkflowUpdatedV1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_WorkflowRegistry *WorkflowRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _WorkflowRegistry.abi.Events["AllowedDONsUpdatedV1"].ID:
		return _WorkflowRegistry.ParseAllowedDONsUpdatedV1(log)
	case _WorkflowRegistry.abi.Events["AuthorizedAddressesUpdatedV1"].ID:
		return _WorkflowRegistry.ParseAuthorizedAddressesUpdatedV1(log)
	case _WorkflowRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _WorkflowRegistry.ParseOwnershipTransferRequested(log)
	case _WorkflowRegistry.abi.Events["OwnershipTransferred"].ID:
		return _WorkflowRegistry.ParseOwnershipTransferred(log)
	case _WorkflowRegistry.abi.Events["RegistryLockedV1"].ID:
		return _WorkflowRegistry.ParseRegistryLockedV1(log)
	case _WorkflowRegistry.abi.Events["RegistryUnlockedV1"].ID:
		return _WorkflowRegistry.ParseRegistryUnlockedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowActivatedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowActivatedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowDeletedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowDeletedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowForceUpdateSecretsRequestedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowForceUpdateSecretsRequestedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowPausedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowPausedV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowRegisteredV1"].ID:
		return _WorkflowRegistry.ParseWorkflowRegisteredV1(log)
	case _WorkflowRegistry.abi.Events["WorkflowUpdatedV1"].ID:
		return _WorkflowRegistry.ParseWorkflowUpdatedV1(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (WorkflowRegistryAllowedDONsUpdatedV1) Topic() common.Hash {
	return common.HexToHash("0xcab63bf31d1e656baa23cebef64e12033ea0ffbd44b1278c3747beec2d2f618c")
}

func (WorkflowRegistryAuthorizedAddressesUpdatedV1) Topic() common.Hash {
	return common.HexToHash("0x509460cccbb176edde6cac28895a4415a24961b8f3a0bd2617b9bb7b4e166c9b")
}

func (WorkflowRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (WorkflowRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (WorkflowRegistryRegistryLockedV1) Topic() common.Hash {
	return common.HexToHash("0x2789711f6fd67d131ad68378617b5d1d21a2c92b34d7c3745d70b3957c08096c")
}

func (WorkflowRegistryRegistryUnlockedV1) Topic() common.Hash {
	return common.HexToHash("0x11a03e25ee25bf1459f9e1cb293ea03707d84917f54a65e32c9a7be2f2edd68a")
}

func (WorkflowRegistryWorkflowActivatedV1) Topic() common.Hash {
	return common.HexToHash("0x17b2d730bb5e064df3fbc6165c8aceb3b0d62c524c196c0bc1012209280bc9a6")
}

func (WorkflowRegistryWorkflowDeletedV1) Topic() common.Hash {
	return common.HexToHash("0x76ee2dfcae10cb8522e62e713e62660e09ecfaab08db15d9404de19141322571")
}

func (WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1) Topic() common.Hash {
	return common.HexToHash("0x95d94f817db4971aa99ba35d0fe019bd8cc39866fbe02b6d47b5f0f3727fb673")
}

func (WorkflowRegistryWorkflowPausedV1) Topic() common.Hash {
	return common.HexToHash("0x6a0ed88e9cf3cb493ab4028fcb1dc7d18f0130fcdfba096edde0aadbfbf5e99f")
}

func (WorkflowRegistryWorkflowRegisteredV1) Topic() common.Hash {
	return common.HexToHash("0xc4399022965bad9b2b468bbd8c758a7e80cdde36ff3088ddbb7f93bdfb5623cb")
}

func (WorkflowRegistryWorkflowUpdatedV1) Topic() common.Hash {
	return common.HexToHash("0x41161473ce2ed633d9f902aab9702d16a5531da27ec84e1939abeffe54ad7353")
}

func (_WorkflowRegistry *WorkflowRegistry) Address() common.Address {
	return _WorkflowRegistry.address
}

type WorkflowRegistryInterface interface {
	ComputeHashKey(opts *bind.CallOpts, owner common.Address, field string) ([32]byte, error)

	GetAllAllowedDONs(opts *bind.CallOpts) ([]uint32, error)

	GetAllAuthorizedAddresses(opts *bind.CallOpts) ([]common.Address, error)

	GetWorkflowMetadata(opts *bind.CallOpts, workflowOwner common.Address, workflowName string) (WorkflowRegistryWorkflowMetadata, error)

	GetWorkflowMetadataListByDON(opts *bind.CallOpts, donID uint32, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error)

	GetWorkflowMetadataListByOwner(opts *bind.CallOpts, workflowOwner common.Address, start *big.Int, limit *big.Int) ([]WorkflowRegistryWorkflowMetadata, error)

	IsRegistryLocked(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ActivateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error)

	DeleteWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error)

	LockRegistry(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseWorkflow(opts *bind.TransactOpts, workflowKey [32]byte) (*types.Transaction, error)

	RegisterWorkflow(opts *bind.TransactOpts, workflowName string, workflowID [32]byte, donID uint32, status uint8, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error)

	RequestForceUpdateSecrets(opts *bind.TransactOpts, secretsURL string) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UnlockRegistry(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateAllowedDONs(opts *bind.TransactOpts, donIDs []uint32, allowed bool) (*types.Transaction, error)

	UpdateAuthorizedAddresses(opts *bind.TransactOpts, addresses []common.Address, allowed bool) (*types.Transaction, error)

	UpdateWorkflow(opts *bind.TransactOpts, workflowKey [32]byte, newWorkflowID [32]byte, binaryURL string, configURL string, secretsURL string) (*types.Transaction, error)

	FilterAllowedDONsUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAllowedDONsUpdatedV1Iterator, error)

	WatchAllowedDONsUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAllowedDONsUpdatedV1) (event.Subscription, error)

	ParseAllowedDONsUpdatedV1(log types.Log) (*WorkflowRegistryAllowedDONsUpdatedV1, error)

	FilterAuthorizedAddressesUpdatedV1(opts *bind.FilterOpts) (*WorkflowRegistryAuthorizedAddressesUpdatedV1Iterator, error)

	WatchAuthorizedAddressesUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryAuthorizedAddressesUpdatedV1) (event.Subscription, error)

	ParseAuthorizedAddressesUpdatedV1(log types.Log) (*WorkflowRegistryAuthorizedAddressesUpdatedV1, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*WorkflowRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorkflowRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*WorkflowRegistryOwnershipTransferred, error)

	FilterRegistryLockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryLockedV1Iterator, error)

	WatchRegistryLockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryLockedV1) (event.Subscription, error)

	ParseRegistryLockedV1(log types.Log) (*WorkflowRegistryRegistryLockedV1, error)

	FilterRegistryUnlockedV1(opts *bind.FilterOpts) (*WorkflowRegistryRegistryUnlockedV1Iterator, error)

	WatchRegistryUnlockedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryRegistryUnlockedV1) (event.Subscription, error)

	ParseRegistryUnlockedV1(log types.Log) (*WorkflowRegistryRegistryUnlockedV1, error)

	FilterWorkflowActivatedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowActivatedV1Iterator, error)

	WatchWorkflowActivatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowActivatedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowActivatedV1(log types.Log) (*WorkflowRegistryWorkflowActivatedV1, error)

	FilterWorkflowDeletedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowDeletedV1Iterator, error)

	WatchWorkflowDeletedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowDeletedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowDeletedV1(log types.Log) (*WorkflowRegistryWorkflowDeletedV1, error)

	FilterWorkflowForceUpdateSecretsRequestedV1(opts *bind.FilterOpts, owner []common.Address) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1Iterator, error)

	WatchWorkflowForceUpdateSecretsRequestedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, owner []common.Address) (event.Subscription, error)

	ParseWorkflowForceUpdateSecretsRequestedV1(log types.Log) (*WorkflowRegistryWorkflowForceUpdateSecretsRequestedV1, error)

	FilterWorkflowPausedV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowPausedV1Iterator, error)

	WatchWorkflowPausedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowPausedV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowPausedV1(log types.Log) (*WorkflowRegistryWorkflowPausedV1, error)

	FilterWorkflowRegisteredV1(opts *bind.FilterOpts, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowRegisteredV1Iterator, error)

	WatchWorkflowRegisteredV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowRegisteredV1, workflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowRegisteredV1(log types.Log) (*WorkflowRegistryWorkflowRegisteredV1, error)

	FilterWorkflowUpdatedV1(opts *bind.FilterOpts, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (*WorkflowRegistryWorkflowUpdatedV1Iterator, error)

	WatchWorkflowUpdatedV1(opts *bind.WatchOpts, sink chan<- *WorkflowRegistryWorkflowUpdatedV1, oldWorkflowID [][32]byte, workflowOwner []common.Address, donID []uint32) (event.Subscription, error)

	ParseWorkflowUpdatedV1(log types.Log) (*WorkflowRegistryWorkflowUpdatedV1, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
