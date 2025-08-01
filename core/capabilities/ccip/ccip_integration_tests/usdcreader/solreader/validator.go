package solreader

import (
	"bytes"
	"net/netip"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
)

const SolanaValidatorCmd = "solana-test-validator"

type SolanaValidator struct {
	rpcAddrPort netip.AddrPort
	wssAddrPort netip.AddrPort
	fctAddrPort netip.AddrPort
	flags       map[string][]string // args that can only appear zero or exactly once go here
	opts        []string            // args that can be specified more than once go here
}

func NewSolanaValidator() *SolanaValidator {
	localhost := netip.MustParseAddr("127.0.0.1")
	return &SolanaValidator{
		// default values were obtained by running: `solana-test-validator --help`
		// solana requires the websocket port to be 1 + rpc port
		wssAddrPort: netip.AddrPortFrom(localhost, 8900),
		rpcAddrPort: netip.AddrPortFrom(localhost, 8899),
		fctAddrPort: netip.AddrPortFrom(localhost, 9900),
		flags:       map[string][]string{},
		opts:        []string{},
	}
}

func (v *SolanaValidator) WithTestDefaults(t *testing.T) *SolanaValidator {
	v.flags["--ticks-per-slot"] = []string{"8"} // value in mainnet: 64
	v.flags["--ledger"] = []string{t.TempDir()}
	v.flags["--reset"] = []string{}

	// account data direct mapping feature is disabled on mainnet,
	// so we disable it here to make the local cluster more similar to mainnet
	v.opts = append(v.opts,
		"--deactivate-feature", "EenyoWx9UMXYKpR8mW5Jmfmy2fRjzUtM7NduYMY8bx33",
	)

	return v
}

func (v *SolanaValidator) AddProgram(programID solana.PublicKey, programFilePath string, deployer solana.PublicKey) *SolanaValidator {
	v.opts = append(v.opts, "--upgradeable-program", programID.String(), programFilePath, deployer.String())
	return v
}

func (v *SolanaValidator) SetFaucetPort(port uint16) *SolanaValidator {
	v.fctAddrPort = netip.AddrPortFrom(v.fctAddrPort.Addr(), port)
	v.flags["--faucet-port"] = []string{strconv.FormatUint(uint64(v.fctAddrPort.Port()), 10)}
	return v
}

func (v *SolanaValidator) SetRpcPort(port uint16) *SolanaValidator {
	v.wssAddrPort = netip.AddrPortFrom(v.rpcAddrPort.Addr(), port+1)
	v.rpcAddrPort = netip.AddrPortFrom(v.rpcAddrPort.Addr(), port)
	v.flags["--rpc-port"] = []string{strconv.FormatUint(uint64(v.rpcAddrPort.Port()), 10)}
	return v
}

func (v *SolanaValidator) RpcUrlString() string {
	return "http://" + v.rpcAddrPort.String()
}

func (v *SolanaValidator) WssUrlString() string {
	return "ws://" + v.wssAddrPort.String()
}

func (v *SolanaValidator) args() []string {
	args := []string{}
	for k, v := range v.flags {
		args = append(args, append([]string{k}, v...)...)
	}
	return append(args, v.opts...)
}

func (v *SolanaValidator) Run(t *testing.T) *SolanaValidator {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	args := v.args()

	t.Logf("%s %s", SolanaValidatorCmd, strings.Join(args, " "))
	cmd := exec.CommandContext(t.Context(), SolanaValidatorCmd, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("failed to kill solana validator: %v", err)
		}
	})

	client := rpc.New(v.RpcUrlString())
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		if stdout.Len() > 0 {
			t.Logf("[solana-test-validator | stdout]:\n%s", stdout.String())
		}
		if stderr.Len() > 0 {
			t.Logf("[solana-test-validator | stderr]:\n%s", stderr.String())
		}

		if out, err := client.GetHealth(t.Context()); err != nil || out != rpc.HealthOk {
			t.Logf("API server not ready yet (attempt %d)\n", i+1)
			continue
		} else {
			t.Logf("API server ready (attempt %d)\n", i+1)
			return v
		}
	}

	t.Fatalf(
		"Cmd output: %s\nCmd error: %s\n",
		stdout.String(),
		stderr.String(),
	)

	return v
}
