package testhelpers

import (
	"os"

	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	sui_query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	chainreader "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	suicrcwconfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	suicodec "github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	suitestutils "github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
	suideps "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"

	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	chain_reader_types "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	crConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type SuiSendRequest struct {
	Receiver      []byte
	Data          []byte
	ExtraArgs     []byte
	FeeToken      string
	FeeTokenStore string
	TokenAmounts  []SuiTokenAmount
}

type SuiTokenAmount struct {
	Token  string
	Amount uint64
}

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

func baseCCIPConfig(
	ccipPkg string,
	pubKey []byte,
	extra []config.ChainWriterPTBCommand,
) config.ChainWriterConfig {
	// common PTB command 0: create_token_params
	cmds := []config.ChainWriterPTBCommand{{
		Type:      suicodec.SuiPTBCommandMoveCall,
		PackageId: strPtr(ccipPkg),
		ModuleId:  strPtr("onramp_state_helper"),
		Function:  strPtr("create_token_transfer_params"),
		Params: []suicodec.SuiFunctionParam{
			{
				Name:     "token_receiver",
				Type:     "address",
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
	feeTokenPkgId string,
) config.ChainWriterConfig {
	feeTokenType := fmt.Sprintf("%s::link::LINK", feeTokenPkgId)
	extra := []config.ChainWriterPTBCommand{{
		Type:      suicodec.SuiPTBCommandMoveCall,
		PackageId: strPtr(onRampPkg),
		ModuleId:  strPtr("onramp"),
		Function:  strPtr("ccip_send"),
		Params: []suicodec.SuiFunctionParam{
			{Name: "ref", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(true)},
			{Name: "state", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(true)},
			{Name: "clock", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
			{Name: "dest_chain_selector", Type: "u64", Required: true},
			{Name: "receiver", Type: "vector<u8>", Required: true},
			{Name: "data", Type: "vector<u8>", Required: true},
			{Name: "token_params", Type: "ptb_dependency", Required: true,
				PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 0}},
			{Name: "fee_token_metadata", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false), GenericType: strPtr(feeTokenType)},
			{Name: "fee_token", Type: "object_id", Required: true},
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
				{Name: "c", Type: "object_id", Required: true},
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
				{Name: "dest_chain_selector", Type: "u64", Required: true},
				{Name: "receiver", Type: "vector<u8>", Required: true},
				{Name: "data", Type: "vector<u8>", Required: true},
				{Name: "token_params", Type: "ptb_dependency", Required: true,
					PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 1}},
				{Name: "fee_token_metadata", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
				{Name: "fee_token", Type: "object_id", Required: true},
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

func SendSuiRequestViaChainWriter(e cldf.Environment, cfg *ccipclient.CCIPSendReqConfig) (*ccipclient.AnyMsgSentEvent, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[cfg.SourceChain]

	publicKeyBytes, err := cldf_sui.PublicKeyBytes(suiChain.Signer)
	if err != nil {
		return nil, err
	}

	// keyString, err := suiChain.Signer.GetAddress()
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	deps := suideps.SuiDeps{
		SuiChain: sui_ops.OpTxDeps{
			Client: suiChain.Client,
			Signer: suiChain.Signer,
			GetCallOpts: func() *suiBind.CallOpts {
				b := uint64(400_000_000)
				return &suiBind.CallOpts{
					WaitForExecution: true,
					GasBudget:        &b,
				}
			},
		},
	}

	ccipObjectRefId := state.SuiChains[cfg.SourceChain].CCIPObjectRef
	ccipPackageId := state.SuiChains[cfg.SourceChain].CCIPAddress
	onRampPackageId := state.SuiChains[cfg.SourceChain].OnRampAddress
	onRampStateObjectId := state.SuiChains[cfg.SourceChain].OnRampStateObjectId
	linkTokenPkgId := state.SuiChains[cfg.SourceChain].LinkTokenAddress
	linkTokenObjectMetadataId := state.SuiChains[cfg.SourceChain].LinkTokenCoinMetadataId

	bigIntSourceUsdPerToken, ok := new(big.Int).SetString("150000000000000000000000000000", 10)
	if !ok {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed converting SourceUSDPerToken to bigInt")
	}

	bigIntGasUsdPerUnitGas, ok := new(big.Int).SetString("7500000000000", 10)
	if !ok {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed converting GasUsdPerUnitGas to bigInt")
	}

	// Update Prices on FeeQuoter with minted LinkToken
	_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterUpdateTokenPricesOp, deps.SuiChain,
		ccipops.FeeQuoterUpdateTokenPricesInput{
			CCIPPackageId: ccipPackageId,
			CCIPObjectRef: ccipObjectRefId,
			// FeeQuoterCapId:        feeQuoterCapId,
			SourceTokens:          []string{linkTokenObjectMetadataId},
			SourceUsdPerToken:     []*big.Int{bigIntSourceUsdPerToken},
			GasDestChainSelectors: []uint64{cfg.DestChain},
			GasUsdPerUnitGas:      []*big.Int{bigIntGasUsdPerUnitGas},
		})
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to updatePrice for Sui chain %d: %w", cfg.SourceChain, err)
	}

	msg := cfg.Message.(SuiSendRequest)
	baseArgs := map[string]any{
		"ref":                 ccipObjectRefId,
		"state":               onRampStateObjectId,
		"clock":               "0x6",
		"dest_chain_selector": cfg.DestChain,
		"token_receiver":      "0x0000000000000000000000000000000000000000000000000000000000000000", // random sui address
		"receiver": []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0xdd, 0xbb, 0x6f, 0x35,
			0x8f, 0x29, 0x04, 0x08, 0xd7, 0x68, 0x47, 0xb4,
			0xf6, 0x02, 0xf0, 0xfd, 0x59, 0x92, 0x95, 0xfd,
		},
		"data":               []byte("hello evm from sui"),
		"fee_token_metadata": msg.FeeTokenMetadata, // this should be the minted token objectId
		"fee_token":          msg.FeeToken,
		"extra_args":         []byte{},
	}

	var (
		BurnMintTP      string
		BurnMintTPState string
		ptbArgs         suicrcwconfig.Arguments
	)
	if len(msg.TokenAmounts) > 0 {
		// Build PTB for token transfer

		// TOKEN POOL SETUP
		BurnMintTP, BurnMintTPState, err = handleTokenAndPoolDeploymentForSUI(e, cfg, deps)
		if err != nil {
			return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to setup tokenPool on SUI %d: %w", cfg.SourceChain, err)
		}
		fmt.Println("TOKEN AMOUNTS: ", msg.TokenAmounts)
		extra := map[string]any{
			"state": BurnMintTPState,
			"c":     msg.TokenAmounts[0].Token,
		}
		ptbArgs = BuildPTBArgs(baseArgs, linkTokenPkgId+"::link::LINK", extra)
	} else {
		// Build PTB for msg transfer
		ptbArgs = BuildPTBArgs(baseArgs, linkTokenPkgId+"::link::LINK", nil)
	}

	// Setup new PTB client
	keystoreInstance := suitestutils.NewTestKeystore(&testing.T{})
	priv, err := cldf_sui.PrivateKey(suiChain.Signer)
	if err != nil {
		return nil, err
	}
	keystoreInstance.AddKey(priv)

<<<<<<< HEAD
	relayerClient, err := client.NewPTBClient(e.Logger, suiChain.URL, nil, 30*time.Second, keystoreInstance, 5, "WaitForEffectsCert")
=======
	relayerClient, err := client.NewPTBClient(e.Logger, "https://testnet.sui.eu.endpoints.matrixed.link/?auth=CL-DNqCV86SzbDs2m", nil, 30*time.Second, keystoreInstance, 5, "WaitForEffectsCert")
>>>>>>> aa3a11c6d7 (working on feeToken rn)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	e.Logger.Info("relayerClient", relayerClient)

	store := txm.NewTxmStoreImpl(e.Logger)
	conf := txm.DefaultConfigSet

	retryManager := txm.NewDefaultRetryManager(5)
	gasLimit := big.NewInt(30000000)
	gasManager := txm.NewSuiGasManager(e.Logger, relayerClient, *gasLimit, 0)

	txManager, err := txm.NewSuiTxm(e.Logger, relayerClient, keystoreInstance, conf, store, retryManager, gasManager)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("Failed to create SuiTxm: %v", err)
	}

	var chainWriterConfig suicrcwconfig.ChainWriterConfig
	if BurnMintTP != "" {
		chainWriterConfig = configureChainWriterForMultipleTokens(ccipPackageId, onRampPackageId, publicKeyBytes, BurnMintTP)
	} else {
		chainWriterConfig = configureChainWriterForMsg(ccipPackageId, onRampPackageId, publicKeyBytes, linkTokenPkgId)
	}

	chainWriter, err := chainwriter.NewSuiChainWriter(e.Logger, txManager, chainWriterConfig, false)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel() // to ensure other calls associated with this context are released

	err = chainWriter.Start(ctx)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	err = txManager.Start(ctx)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	txId := "ccip_send_msg_transfer"
	err = chainWriter.SubmitTransaction(ctx,
		suicrcwconfig.PTBChainWriterModuleName,
		"ccip_send",
		&ptbArgs,
		txId,
		onRampPackageId, // this is the contract address so onramp in this case
		&commonTypes.TxMeta{GasLimit: big.NewInt(30000000)},
		nil,
	)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	// TODO: find a better way of handling waitForTransaction
	time.Sleep(10 * time.Second)
	status, err := chainWriter.GetTransactionStatus(ctx, txId)

	if status != commonTypes.Finalized {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("tx failed to get finalized")
	}

	e.Logger.Infof("(Sui) CCIP message sent (tx %s) from chain selector %d to chain selector %d", txId, cfg.SourceChain, cfg.DestChain)

	chainWriter.Close()
	txManager.Close()

	// Query the CCIPSend Event via chainReader
	chainReaderConfig := crConfig.ChainReaderConfig{
		IsLoopPlugin: false,
		Modules: map[string]*crConfig.ChainReaderModule{
			"onramp": {
				Name: "onramp",
				Events: map[string]*crConfig.ChainReaderEvent{
					"CCIPMessageSent": {
						Name:      "CCIPMessageSent",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventSelector{
							Package: onRampPackageId,
							Module:  "onramp",
							Event:   "CCIPMessageSent",
						},
					},
				},
			},
		},
	}

	dbURL := os.Getenv("CL_DATABASE_URL")

	err = sqltest.RegisterTxDB(dbURL)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	db, err := sqlx.Open(pg.DriverTxWrappedPostgres, uuid.New().String())
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)

	// attempt to connect
	_, err = db.Connx(ctx)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	// Create the indexers
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		e.Logger,
		relayerClient,
		10*time.Second,
		10*time.Second,
		// start without any configs, they will be set when ChainReader is initialized and gets a reference
		// to the transaction indexer to avoid having to reading ChainReader configs here as well
		map[string]*crConfig.ChainReaderEvent{},
	)
	evIndexer := indexer.NewEventIndexer(
		db,
		e.Logger,
		relayerClient,
		// start without any selectors, they will be added during .Bind() calls on ChainReader
		[]*client.EventSelector{},
		10*time.Second,
		10*time.Second,
	)
	indexerInstance := indexer.NewIndexer(
		e.Logger,
		evIndexer,
		txnIndexer,
	)

	chainReader, err := chainreader.NewChainReader(ctx, e.Logger, relayerClient, chainReaderConfig, db, indexerInstance)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	err = chainReader.Start(ctx)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	err = indexerInstance.Start(ctx)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	err = chainReader.Bind(context.Background(), []chain_reader_types.BoundContract{{
		Name:    "onramp",
		Address: onRampPackageId, // Package ID of the deployed counter contract
	}})
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to bind onramp contract with chainReader")
	}

	// TODO handle this better, maybe retrieve it from the bindings when we do binding ccip_send
	e.Logger.Debugw("Querying for ccip_send events",
		"filter", "CCIPMessageSent",
		"limit", 50,
		"packageId", onRampPackageId,
		"contract", "onramp",
		"eventType", "CCIPMessageSent")

	var ccipSendEvent CCIPMessageSent
	sequences, err := chainReader.QueryKey(
		ctx,
		chain_reader_types.BoundContract{
			Name:    "onramp",
			Address: onRampPackageId, // Package ID of the deployed counter contract
		},
		sui_query.KeyFilter{
			Key: "CCIPMessageSent",
		},
		sui_query.LimitAndSort{
			Limit: sui_query.Limit{
				Count:  50,
				Cursor: "",
			},
		},
		&ccipSendEvent,
	)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to query events", "error", err)
	}

	if len(sequences) < 1 {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to fetch event sequence")
	}
	e.Logger.Debugw("Query results", "sequences", sequences)
	rawevent := sequences[0].Data.(*CCIPMessageSent)

	chainReader.Close()
	indexerInstance.Close()

	return &ccipclient.AnyMsgSentEvent{
		SequenceNumber: rawevent.SequenceNumber,
		RawEvent:       rawevent,
	}, nil
}
