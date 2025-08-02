package testhelpers

import (
	"crypto/ed25519"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	sui_bind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	suicodec "github.com/smartcontractkit/chainlink-sui/relayer/codec"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	suideps "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type suiCtx struct {
	Deps                suideps.SuiDeps
	CCIPObjectRefID     string
	CCIPPackageID       string
	OnRampPackageID     string
	OnRampStateObjectID string
	LinkTokenPkgID      string
	LinkTokenMetaID     string
	LinkTokenCapID      string
	SignerAddr          string
	PubKeyBytes         []byte
}

type RampMessageHeader struct {
	MessageId           []byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type Sui2AnyRampMessage struct {
	Header         RampMessageHeader
	Sender         string
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       string
	FeeTokenAmount uint64
	FeeValueJuels  uint64
}

type CCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           Sui2AnyRampMessage
}

func NewSuiCtx(e cldf.Environment, src uint64) (*suiCtx, error) {
	st, err := stateview.LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	sc := e.BlockChains.SuiChains()[src]
	signer := rel.NewPrivateKeySigner(sc.DeployerKey)
	addr, err := signer.GetAddress()
	if err != nil {
		return nil, err
	}
	pub := sc.DeployerKey.Public().(ed25519.PublicKey)

	deps := suideps.SuiDeps{
		SuiChain: sui_ops.OpTxDeps{
			Client: sc.Client,
			Signer: signer,
			GetCallOpts: func() *sui_bind.CallOpts {
				b := uint64(400_000_000)
				return &sui_bind.CallOpts{
					WaitForExecution: true,
					GasBudget:        &b,
				}
			},
		},
	}

	return &suiCtx{
		Deps:                deps,
		CCIPObjectRefID:     st.SuiChains[src].CCIPObjectRef,
		CCIPPackageID:       st.SuiChains[src].CCIPAddress,
		OnRampPackageID:     st.SuiChains[src].OnRampAddress,
		OnRampStateObjectID: st.SuiChains[src].OnRampStateObjectId,
		LinkTokenPkgID:      st.SuiChains[src].LinkTokenAddress,
		LinkTokenMetaID:     st.SuiChains[src].LinkTokenCoinMetadataId,
		LinkTokenCapID:      st.SuiChains[src].LinkTokenTreasuryCapId,
		SignerAddr:          addr,
		PubKeyBytes:         []byte(pub),
	}, nil
}

func baseCCIPConfig(
	ccipPkg string,
	pubKey []byte,
	extra []config.ChainWriterPTBCommand,
) config.ChainWriterConfig {
	// common PTB command 0: create_token_params
	cmds := []config.ChainWriterPTBCommand{{
		Type:      suicodec.SuiPTBCommandMoveCall,
		PackageId: strPtr(ccipPkg),
		ModuleId:  strPtr("dynamic_dispatcher"),
		Function:  strPtr("create_token_params"),
		Params: []suicodec.SuiFunctionParam{
			{
				Name:     "destination_chain_selector",
				Type:     "u64",
				Required: true,
			},
			{
				Name:     "receiver",
				Type:     "vector<u8>",
				Required: true,
			},
		},
	}}
	// append the variant commands
	cmds = append(cmds, extra...)

	return config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			config.PTBChainWriterModuleName: {
				Name:     config.PTBChainWriterModuleName,
				ModuleID: "0x123",
				Functions: map[string]*config.ChainWriterFunction{
					"ccip_send": {
						Name:        "ccip_send",
						PublicKey:   pubKey,
						Params:      []suicodec.SuiFunctionParam{},
						PTBCommands: cmds,
					},
				},
			},
		},
	}
}

// 2a) Simple Message → EVM
func configureChainWriterForMsg(
	ccipPkg, onRampPkg string,
	pubKey []byte,
) config.ChainWriterConfig {
	extra := []config.ChainWriterPTBCommand{{
		Type:      suicodec.SuiPTBCommandMoveCall,
		PackageId: strPtr(onRampPkg),
		ModuleId:  strPtr("onramp"),
		Function:  strPtr("ccip_send"),
		Params: []suicodec.SuiFunctionParam{
			{Name: "ref", Type: "object_id", Required: true},
			{Name: "state", Type: "object_id", Required: true},
			{Name: "clock", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
			{Name: "data", Type: "vector<u8>", Required: true},
			{Name: "token_params", Type: "ptb_dependency", Required: true,
				PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 0}},
			{Name: "fee_token_metadata", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
			{Name: "fee_token", Type: "object_id", Required: true, IsGeneric: true},
			{Name: "extra_args", Type: "vector<u8>", Required: true},
		},
	}}
	return baseCCIPConfig(ccipPkg, pubKey, extra)
}

// 2b) Message + BurnMintTP → EVM
func configureChainWriterForMultipleTokens(
	ccipPkg, onRampPkg string,
	pubKey []byte,
	tokenPool string,
) config.ChainWriterConfig {
	extra := []config.ChainWriterPTBCommand{
		// lock-or-burn command
		{
			Type:      suicodec.SuiPTBCommandMoveCall,
			PackageId: strPtr(tokenPool),
			ModuleId:  strPtr("burn_mint_token_pool"),
			Function:  strPtr("lock_or_burn"),
			Params: []suicodec.SuiFunctionParam{
				{Name: "ref", Type: "object_id", Required: true},
				{Name: "clock", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
				{Name: "state", Type: "object_id", Required: true},
				{Name: "c", Type: "object_id", Required: true, IsGeneric: true},
				{Name: "token_params", Type: "ptb_dependency", Required: true,
					PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 0}},
			},
		},
		// the same onramp send
		{
			Type:      suicodec.SuiPTBCommandMoveCall,
			PackageId: strPtr(onRampPkg),
			ModuleId:  strPtr("onramp"),
			Function:  strPtr("ccip_send"),
			Params: []suicodec.SuiFunctionParam{
				{Name: "ref", Type: "object_id", Required: true},
				{Name: "state", Type: "object_id", Required: true},
				{Name: "clock", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
				{Name: "data", Type: "vector<u8>", Required: true},
				{Name: "token_params", Type: "ptb_dependency", Required: true,
					PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 1}},
				{Name: "fee_token_metadata", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
				{Name: "fee_token", Type: "object_id", Required: true, IsGeneric: true},
				{Name: "extra_args", Type: "vector<u8>", Required: true},
			},
		},
	}
	return baseCCIPConfig(ccipPkg, pubKey, extra)
}

func BuildPTBArgs(baseArgs map[string]any, coinType string, extraArgs map[string]any) config.Arguments {
	args := make(map[string]any, len(baseArgs)+len(extraArgs))
	for k, v := range baseArgs {
		args[k] = v
	}
	for k, v := range extraArgs {
		args[k] = v
	}

	argTypes := map[string]string{
		"fee_token": coinType,
	}
	if _, ok := args["c"]; ok {
		argTypes["c"] = coinType
	}

	return config.Arguments{
		Args:     args,
		ArgTypes: argTypes,
	}
}
