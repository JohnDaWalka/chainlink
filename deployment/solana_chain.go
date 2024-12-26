package deployment

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// SolChain represents a Solana chain.
type SolChain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	Client   *rpc.Client
	Confirm  func() (uint64, error)
}

func (c SolChain) String() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	return fmt.Sprintf("%s (%d)", chainInfo.ChainName, chainInfo.ChainSelector)
}

func (c SolChain) Name() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	if chainInfo.ChainName == "" {
		return fmt.Sprintf("%d", c.Selector)
	}
	return chainInfo.ChainName
}

type SolClient interface {
}

type ContractDeploySolana struct {
	ProgramID *solana.PublicKey // We leave this incase a Go binding doesn't have Address()
	Tv        TypeAndVersion
	Err       error
}

func DeploySolContract(
	lggr logger.Logger,
	chain SolChain,
	addressBook AddressBook,
	deploy func(chain SolChain) ContractDeploySolana,
) (*ContractDeploySolana, error) {
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "chain", chain.String(), "err", contractDeploy.Err)
		return nil, contractDeploy.Err
	}
	_, err := chain.Confirm()
	if err != nil {
		lggr.Errorw("Failed to confirm deployment", "chain", chain.String(), "Contract", contractDeploy.Tv.String(), "err", err)
		return nil, err
	}
	lggr.Infow("Deployed contract", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.ProgramID, "chain", chain.String())
	err = addressBook.Save(chain.Selector, "fill in address", contractDeploy.Tv)
	if err != nil {
		lggr.Errorw("Failed to save contract address", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.ProgramID, "chain", chain.String(), "err", err)
		return nil, err
	}
	return &contractDeploy, nil
}

func DeploySolProgramCLI(programFile, keypairPath, programKeyPair string) (string, error) {
	// Construct the CLI command: solana program deploy
	cmd := exec.Command("solana", "program", "deploy", programFile, "--keypair", keypairPath, "--program-id", programKeyPair)

	// Capture the command output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error deploying program: %s: %s", err.Error(), stderr.String())
	}

	// Parse and return the program ID
	output := stdout.String()

	time.Sleep(5 * time.Second) // obviously need to do this better
	return parseProgramID(output)
}

// parseProgramID parses the program ID from the deploy output.
func parseProgramID(output string) (string, error) {
	// Look for the program ID in the CLI output
	// Example output: "Program Id: <PROGRAM_ID>"
	const prefix = "Program Id: "
	startIdx := bytes.Index([]byte(output), []byte(prefix))
	if startIdx == -1 {
		return "", fmt.Errorf("failed to find program ID in output")
	}
	startIdx += len(prefix)
	endIdx := bytes.Index([]byte(output[startIdx:]), []byte("\n"))
	if endIdx == -1 {
		endIdx = len(output)
	}
	return output[startIdx : startIdx+endIdx], nil
}
