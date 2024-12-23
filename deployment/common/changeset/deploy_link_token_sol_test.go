package changeset_test

import (
	"context"
	"testing"
	"time"

	"bytes"
	"fmt"
	"os/exec"

	// "github.com/stretchr/testify/require"
	// "go.uber.org/zap/zapcore"
	// "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	// "github.com/smartcontractkit/chainlink/deployment/environment/memory"
	// "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/utils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/external_program_cpi_stub"
	"github.com/stretchr/testify/require"
)

var (
	DefaultCommitment = rpc.CommitmentConfirmed
	StubProgram       = "EQPCTRibpsPcQNb464QVBkS1PkFfuK8kYdpd5Y17HaGh"
)

// deployProgram deploys a Solana program using the Solana CLI.
func deployProgram(programFile string, keypairPath string) (string, error) {

	programKeyPair := "/Users/yashvardhan/chainlink-internal-integrations/solana/contracts/target/deploy/external_program_cpi_stub-keypair.json"
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

// TestDeployProgram is a test for deploying the Solana program.
func TestDeployProgram(t *testing.T) {
	// Path to your .so file and keypair file
	programFile := "/Users/yashvardhan/chainlink-internal-integrations/solana/contracts/target/deploy/external_program_cpi_stub.so"
	keypairPath := "/Users/yashvardhan/.config/solana/id.json" //wallet

	ExternalCpiStubProgram := solana.MustPublicKeyFromBase58("EQPCTRibpsPcQNb464QVBkS1PkFfuK8kYdpd5Y17HaGh")
	solanaGoClient := rpc.New("http://127.0.0.1:8899")

	// Fetch account info for the program ID
	account, _ := solanaGoClient.GetAccountInfo(context.Background(), ExternalCpiStubProgram)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	// Deploy program if it doesn't exist
	if account != nil && account.Value.Executable {
		fmt.Println("Program exists and is executable.")
	} else {
		fmt.Println("Program does not exist or is not executable.")
		// Deploy the program
		programID, err := deployProgram(programFile, keypairPath)
		if err != nil {
			t.Fatalf("Failed to deploy program: %v", err)
		}
		time.Sleep(5 * time.Second) // obviously need to do this better
		// Verify the program ID (simple check for non-empty string)
		if programID == "" {
			t.Fatalf("Program ID is empty")
		}

		t.Logf("programID %s", programID)
	}

	// program should exist by now (either already deployed, or deployed and waited for confirmation)
	external_program_cpi_stub.SetProgramID(ExternalCpiStubProgram)

	// wallet keys
	privateKey, _ := solana.PrivateKeyFromSolanaKeygenFile(keypairPath)
	publicKey := privateKey.PublicKey()

	// this is a PDA that gets initialised when you call init on the programID
	StubAccountPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("u8_value")}, ExternalCpiStubProgram)
	t.Logf("StubAccountPDA %s", StubAccountPDA)

	// check if the PDA is already initialised
	accountInfo, _ := solanaGoClient.GetAccountInfo(context.Background(), ExternalCpiStubProgram)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	var ix *external_program_cpi_stub.Instruction
	var err error

	if accountInfo != nil {
		fmt.Println("Account initialized successfully.")
		fmt.Printf("Account data: %v\n", accountInfo.Value.Data)
		ix, err = external_program_cpi_stub.NewEmptyInstruction().ValidateAndBuild()
	} else {
		fmt.Println("Account does not exist or has no data.")
		ix, err = external_program_cpi_stub.NewInitializeInstruction(
			StubAccountPDA,
			publicKey,
			solana.SystemProgramID, // 1111111
		).ValidateAndBuild()
	}

	utils.SendAndConfirm(context.Background(), t, solanaGoClient, []solana.Instruction{ix}, privateKey, config.DefaultCommitment)

	require.NoError(t, err)
}

// func spinUpDevNet(t *testing.T) (string, string) {
// 	t.Helper()
// 	port := "8899"
// 	portInt, _ := strconv.Atoi(port)

// 	faucetPort := "8877"
// 	url := "http://127.0.0.1:" + port
// 	wsURL := "ws://127.0.0.1:" + strconv.Itoa(portInt+1)

// 	args := []string{
// 		"--reset",
// 		"--rpc-port", port,
// 		"--faucet-port", faucetPort,
// 		"--ledger", t.TempDir(),
// 	}

// 	cmd := exec.Command("solana-test-validator", args...)

// 	var stdErr bytes.Buffer
// 	cmd.Stderr = &stdErr
// 	var stdOut bytes.Buffer
// 	cmd.Stdout = &stdOut
// 	require.NoError(t, cmd.Start())
// 	t.Cleanup(func() {
// 		assert.NoError(t, cmd.Process.Kill())
// 		if err2 := cmd.Wait(); assert.Error(t, err2) {
// 			if !assert.Contains(t, err2.Error(), "signal: killed", cmd.ProcessState.String()) {
// 				t.Logf("solana-test-validator\n stdout: %s\n stderr: %s", stdOut.String(), stdErr.String())
// 			}
// 		}
// 	})

// 	// Wait for api server to boot
// 	var ready bool
// 	for i := 0; i < 30; i++ {
// 		time.Sleep(time.Second)
// 		client := rpc.New(url)
// 		out, err := client.GetHealth(tests.Context(t))
// 		if err != nil || out != rpc.HealthOk {
// 			t.Logf("API server not ready yet (attempt %d)\n", i+1)
// 			continue
// 		}
// 		ready = true
// 		break
// 	}
// 	if !ready {
// 		t.Logf("Cmd output: %s\nCmd error: %s\n", stdOut.String(), stdErr.String())
// 	}
// 	require.True(t, ready)

// 	return url, wsURL
// }

// func getRpcClient(t *testing.T) *rpc.Client {
// 	url, _ := spinUpDevNet(t)
// 	return rpc.New(url)
// }

// func TestTokenDeploy(t *testing.T) {
// 	keypairPath := "/Users/yashvardhan/.config/solana/id.json" //wallet
// 	adminPrivateKey, _ := solana.PrivateKeyFromSolanaKeygenFile(keypairPath)
// 	adminPublicKey := adminPrivateKey.PublicKey()
// 	decimals := uint8(0)
// 	// amount := uint64(1000)
// 	// solanaGoClient := rpc.New("http://127.0.0.1:8899")
// 	solanaGoClient := getRpcClient(t)
// 	mint, _ := solana.NewRandomPrivateKey()
// 	mintPublicKey := mint.PublicKey()
// 	instructions, err := utils.CreateToken(context.Background(), config.Token2022Program, mintPublicKey, adminPublicKey, decimals, solanaGoClient, DefaultCommitment)
// 	utils.SendAndConfirm(context.Background(), t, solanaGoClient, instructions, adminPrivateKey, DefaultCommitment, utils.AddSigners(mint))
// 	require.NoError(t, err)
// }
