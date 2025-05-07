package ccip

import (
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
)

func TestMain(m *testing.M) {
	os.Setenv(string(env.SolanaPlugin.Cmd), env.SolanaPlugin.CmdDefault)
	os.Exit(m.Run())
}
