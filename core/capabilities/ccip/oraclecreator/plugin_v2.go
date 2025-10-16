package oraclecreator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	chainsel "github.com/smartcontractkit/chain-selectors"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	commitocr3 "github.com/smartcontractkit/chainlink-ccip/commit"
	"github.com/smartcontractkit/chainlink-ccip/commit/merkleroot/rmn"
	execocr3 "github.com/smartcontractkit/chainlink-ccip/execute"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/promwrapper"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

// pluginOracleCreatorV2 creates oracles that reference plugins running
// in the same process as the chainlink node, i.e not LOOPPs.
type pluginOracleCreatorV2 struct {
	ocrKeyBundles         map[string]ocr2key.KeyBundle
	transmitters          map[types.RelayID][]string
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	externalJobID         uuid.UUID
	jobID                 int32
	isNewlyCreatedJob     bool
	pluginConfig          job.JSONConfig
	db                    ocr3types.Database
	lggr                  logger.SugaredLogger
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	bootstrapperLocators  []commontypes.BootstrapperLocator
	homeChainReader       ccipreaderpkg.HomeChain
	homeChainSelector     cciptypes.ChainSelector
	relayers              map[types.RelayID]loop.Relayer
	addressCodec          ccipcommon.AddressCodec
	p2pID                 p2pkey.KeyV2
}

func NewPluginOracleCreatorV2(
	ocrKeyBundles map[string]ocr2key.KeyBundle,
	transmitters map[types.RelayID][]string,
	relayers map[types.RelayID]loop.Relayer,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	externalJobID uuid.UUID,
	jobID int32,
	isNewlyCreatedJob bool,
	pluginConfig job.JSONConfig,
	db ocr3types.Database,
	lggr logger.Logger,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	bootstrapperLocators []commontypes.BootstrapperLocator,
	homeChainReader ccipreaderpkg.HomeChain,
	homeChainSelector cciptypes.ChainSelector,
	addressCodec ccipcommon.AddressCodec,
	p2pID p2pkey.KeyV2,
) cctypes.OracleCreator {
	return &pluginOracleCreatorV2{
		ocrKeyBundles:         ocrKeyBundles,
		transmitters:          transmitters,
		relayers:              relayers,
		peerWrapper:           peerWrapper,
		externalJobID:         externalJobID,
		jobID:                 jobID,
		isNewlyCreatedJob:     isNewlyCreatedJob,
		pluginConfig:          pluginConfig,
		db:                    db,
		lggr:                  logger.Sugared(lggr),
		monitoringEndpointGen: monitoringEndpointGen,
		bootstrapperLocators:  bootstrapperLocators,
		homeChainReader:       homeChainReader,
		homeChainSelector:     homeChainSelector,
		addressCodec:          addressCodec,
		p2pID:                 p2pID,
	}
}

// Type implements types.OracleCreator.
func (i *pluginOracleCreatorV2) Type() cctypes.OracleType {
	return cctypes.OracleTypePlugin
}

// Create implements types.OracleCreator.
func (i *pluginOracleCreatorV2) Create(ctx context.Context, donID uint32, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	pluginType := cctypes.PluginType(config.Config.PluginType)
	destChainSelector := uint64(config.Config.ChainSelector)
	destChainFamily, err := chainsel.GetSelectorFamily(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family from selector %d: %w", config.Config.ChainSelector, err)
	}

	destChainID, err := chainsel.GetChainIDFromSelector(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector %d: %w", destChainSelector, err)
	}
	destRelayID := types.NewRelayID(destChainFamily, destChainID)

	configTracker, err := ocrimpls.NewConfigTracker(config, i.addressCodec)
	if err != nil {
		return nil, fmt.Errorf("failed to create config tracker: %w, %d", err, destChainSelector)
	}
	publicConfig, err := configTracker.PublicConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get public config from OCR config: %w", err)
	}

	// Construct ChainReader/ChainWriter options for families that still rely on CR/CW
	crOpts, cwOpts, err := i.constructCRCWOpts(destChainID, pluginType, config, publicConfig, destChainFamily)
	if err != nil {
		return nil, fmt.Errorf("failed to construct ChainReader/ChainWriter options: %w", err)
	}

	destFromAccounts, ok := i.transmitters[destRelayID]
	if !ok {
		i.lggr.Infow("no transmitters found for dest chain, will create nil transmitter",
			"destChainID", destChainID,
			"destChainSelector", config.Config.ChainSelector)
	}

	offRampAddrStr, err := i.addressCodec.AddressBytesToString(
		config.Config.OfframpAddress, cciptypes.ChainSelector(destChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to convert offramp address to string using address codec: %w", err)
	}

	// Create CCIPProviderWrappers for each relayer/chain family
	ccipProviderWrappers, err := i.createCCIPProviderWrappers(
		ctx, config, offRampAddrStr, destFromAccounts, pluginType, crOpts, cwOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCIPProviderWrappers: %w", err)
	}

	// TODO: add binding logic for contract readers and writers

	i.lggr.Infow("Creating plugin using OCR3 settings",
		"plugin", pluginType.String(),
		"chainSelector", destChainSelector,
		"chainID", destChainID,
		"deltaProgress", publicConfig.DeltaProgress,
		"deltaResend", publicConfig.DeltaResend,
		"deltaInitial", publicConfig.DeltaInitial,
		"deltaRound", publicConfig.DeltaRound,
		"deltaGrace", publicConfig.DeltaGrace,
		"deltaCertifiedCommitRequest", publicConfig.DeltaCertifiedCommitRequest,
		"deltaStage", publicConfig.DeltaStage,
		"rMax", publicConfig.RMax,
		"s", publicConfig.S,
		"maxDurationInitialization", publicConfig.MaxDurationInitialization,
		"maxDurationQuery", publicConfig.MaxDurationQuery,
		"maxDurationObservation", publicConfig.MaxDurationObservation,
		"maxDurationShouldAcceptAttestedReport", publicConfig.MaxDurationShouldAcceptAttestedReport,
		"maxDurationShouldTransmitAcceptedReport", publicConfig.MaxDurationShouldTransmitAcceptedReport,
	)

	chainAccessors, contractReaders, contractWriters, looppProviderSupported, err := i.getWrapperComponents(ccipProviderWrappers)

	// Build plugin factory
	factory, err := i.createFactory(
		donID,
		config,
		destRelayID,
		ccipProviderWrappers,
		chainAccessors,
		contractReaders,
		contractWriters,
		looppProviderSupported,
		destFromAccounts,
		publicConfig,
		destChainFamily,
		destChainID,
		offRampAddrStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin factory: %w", err)
	}

	// build the onchain keyring. it will be the signing key for the destination chain family.
	keybundle, ok := i.ocrKeyBundles[destChainFamily]
	if !ok {
		return nil, fmt.Errorf("no OCR key bundle found for chain family %s, forgot to create one?", destChainFamily)
	}
	onchainKeyring := ocrimpls.NewOnchainKeyring[[]byte](keybundle, i.lggr)

	telemetryType, err := pluginTypeToTelemetryType(pluginType)
	if err != nil {
		return nil, fmt.Errorf("failed to get telemetry type: %w", err)
	}

	// TODO: probably add some nil checks here, also add case for NewNoOpTransmitter if we don't have a transmitter
	// address for this chain
	transmitter := ccipProviderWrappers[config.Config.ChainSelector].CCIPProvider().ContractTransmitter()
	oracleArgs := libocr3.OCR3OracleArgs[[]byte]{
		BinaryNetworkEndpointFactory: i.peerWrapper.Peer2,
		Database:                     i.db,
		// NOTE: when specifying V2Bootstrappers here we actually do NOT need to run a full bootstrap node!
		// Thus it is vital that the bootstrapper locators are correctly set in the job spec.
		V2Bootstrappers:       i.bootstrapperLocators,
		ContractConfigTracker: configTracker,
		ContractTransmitter:   transmitter,
		LocalConfig:           defaultLocalConfig(),
		Logger: ocrcommon.NewOCRWrapper(
			i.lggr.
				Named(fmt.Sprintf("CCIP%sOCR3", pluginType.String())).
				Named(destRelayID.String()).
				Named(offRampAddrStr),
			false,
			func(ctx context.Context, msg string) {}),
		MetricsRegisterer: prometheus.WrapRegistererWith(map[string]string{"name": fmt.Sprintf("commit-%d", config.Config.ChainSelector)}, prometheus.DefaultRegisterer),
		MonitoringEndpoint: i.monitoringEndpointGen.GenMonitoringEndpoint(
			destChainFamily,
			destRelayID.ChainID,
			offRampAddrStr,
			telemetryType,
		),
		OffchainConfigDigester: ocrimpls.NewConfigDigester(config.ConfigDigest),
		OffchainKeyring:        keybundle,
		OnchainKeyring:         onchainKeyring,
		ReportingPluginFactory: factory,
	}
	oracle, err := libocr3.NewOracle(oracleArgs)
	if err != nil {
		return nil, err
	}

	closers := make([]io.Closer, 0, len(ccipProviderWrappers))
	for _, w := range ccipProviderWrappers {
		closers = append(closers, w.CCIPProvider())
	}
	return newWrappedOracle(oracle, closers), nil
}

func (*pluginOracleCreatorV2) getWrapperComponents(
	ccipProviderWrappers map[cciptypes.ChainSelector]ccipcommon.CCIPProviderWrapper,
) (
	map[cciptypes.ChainSelector]cciptypes.ChainAccessor,
	map[cciptypes.ChainSelector]contractreader.Extended,
	map[cciptypes.ChainSelector]types.ContractWriter,
	map[string]bool,
	error,
) {
	chainAccessors := make(map[cciptypes.ChainSelector]cciptypes.ChainAccessor)
	extendedReaders := make(map[cciptypes.ChainSelector]contractreader.Extended)
	chainWriters := make(map[cciptypes.ChainSelector]types.ContractWriter)
	looppProviderSupported := make(map[string]bool)
	for chainSelector, wrapper := range ccipProviderWrappers {
		if wrapper == nil || wrapper.CCIPProvider() == nil {
			return nil, nil, nil, nil, fmt.Errorf("CCIPProviderWrapper or CCIPProvider for chain selector %d is nil", chainSelector)
		}
		if wrapper.CCIPProvider().ChainAccessor() != nil {
			chainAccessors[chainSelector] = wrapper.CCIPProvider().ChainAccessor()
		}
		if wrapper.ContractReader() != nil {
			extendedReaders[chainSelector] = wrapper.ContractReader()
		}
		if wrapper.ContractWriter() != nil {
			chainWriters[chainSelector] = wrapper.ContractWriter()
		}
		looppProviderSupported[wrapper.ChainFamily()] = wrapper.CCIPProviderSupported()
	}
	return chainAccessors, extendedReaders, chainWriters, looppProviderSupported, nil
}

func (i *pluginOracleCreatorV2) createFactory(
	donID uint32,
	config cctypes.OCR3ConfigWithMeta,
	destRelayID types.RelayID,
	ccipProviderWrappers map[cciptypes.ChainSelector]ccipcommon.CCIPProviderWrapper,
	chainAccessors map[cciptypes.ChainSelector]cciptypes.ChainAccessor,
	extendedReaders map[cciptypes.ChainSelector]contractreader.Extended,
	chainWriters map[cciptypes.ChainSelector]types.ContractWriter,
	looppProviderSupported map[string]bool,
	destFromAccounts []string,
	publicConfig ocr3confighelper.PublicConfig,
	destChainFamily string,
	destChainID string,
	offrampAddrStr string,
) (ocr3types.ReportingPluginFactory[[]byte], error) {
	var factory ocr3types.ReportingPluginFactory[[]byte]
	ccipProviderWrapper, exists := ccipProviderWrappers[config.Config.ChainSelector]
	if !exists {
		return nil, fmt.Errorf("no CCIPProviderWrapper found for chain selector %d", config.Config.ChainSelector)
	}
	codecs := ccipProviderWrapper.CCIPProvider().Codec()

	if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPCommit) {
		if !i.peerWrapper.IsStarted() {
			return nil, errors.New("peer wrapper is not started")
		}

		i.lggr.Infow("creating rmn peer client",
			"bootstrapperLocators", i.bootstrapperLocators,
			"deltaRound", publicConfig.DeltaRound)

		rmnPeerClient := rmn.NewPeerClient(
			i.lggr.Named("RMNPeerClient"),
			i.peerWrapper.PeerGroupFactory,
			i.bootstrapperLocators,
			publicConfig.DeltaRound,
		)

		factory = commitocr3.NewCommitPluginFactory(
			commitocr3.CommitPluginFactoryParams{
				Lggr: i.lggr.
					Named("CCIPCommitPlugin").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)).
					Named(offrampAddrStr),
				DonID:                      donID,
				OcrConfig:                  ccipreaderpkg.OCR3ConfigWithMeta(config),
				CommitCodec:                codecs.CommitPluginCodec,
				MsgHasher:                  codecs.MessageHasher,
				AddrCodec:                  i.addressCodec,
				HomeChainReader:            i.homeChainReader,
				HomeChainSelector:          i.homeChainSelector,
				ChainAccessors:             chainAccessors,
				LOOPPCCIPProviderSupported: looppProviderSupported,
				ExtendedReaders:            extendedReaders,
				ContractWriters:            chainWriters,
				RmnPeerClient:              rmnPeerClient,
				RmnCrypto:                  ccipProviderWrapper.RMNCrypto(),
			},
		)
		factory = promwrapper.NewReportingPluginFactory(
			factory,
			i.lggr,
			destChainFamily,
			destChainID,
			"CCIPCommit",
		)
	} else if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec) {
		ccipProviderWrapper, exists := ccipProviderWrappers[config.Config.ChainSelector]
		if !exists {
			return nil, fmt.Errorf("no CCIPProviderWrapper found for chain selector %d", config.Config.ChainSelector)
		}

		factory = execocr3.NewExecutePluginFactory(
			execocr3.PluginFactoryParams{
				Lggr: i.lggr.
					Named("CCIPExecPlugin").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)).
					Named(offrampAddrStr),
				DonID:                      donID,
				OcrConfig:                  ccipreaderpkg.OCR3ConfigWithMeta(config),
				ExecCodec:                  codecs.ExecutePluginCodec,
				MsgHasher:                  codecs.MessageHasher,
				AddrCodec:                  i.addressCodec,
				HomeChainReader:            i.homeChainReader,
				TokenDataEncoder:           codecs.TokenDataEncoder,
				EstimateProvider:           ccipProviderWrapper.GasEstimateProvider(),
				LOOPPCCIPProviderSupported: looppProviderSupported,
				ChainAccessors:             chainAccessors,
				ExtendedReaders:            extendedReaders,
				ContractWriters:            chainWriters,
			})
		factory = promwrapper.NewReportingPluginFactory(
			factory,
			i.lggr,
			destChainFamily,
			destChainID,
			"CCIPExec",
		)
	} else {
		return nil, fmt.Errorf("unsupported plugin type: %d", config.Config.PluginType)
	}
	return factory, nil
}

func (i *pluginOracleCreatorV2) constructCRCWOpts(
	destChainID string,
	pluginType cctypes.PluginType,
	config cctypes.OCR3ConfigWithMeta,
	publicCfg ocr3confighelper.PublicConfig,
	destChainFamily string,
) (
	map[cciptypes.ChainSelector]ccipcommon.ChainReaderProviderOpts,
	map[cciptypes.ChainSelector]ccipcommon.ChainWriterProviderOpts,
	error,
) {
	ofc, err := decodeAndValidateOffchainConfig(pluginType, publicCfg)
	if err != nil {
		return nil, nil, err
	}

	var execBatchGasLimit uint64
	if !ofc.ExecEmpty() {
		execBatchGasLimit = ofc.Execute.BatchGasLimit
	} else {
		// Set the default here so chain writer config validation doesn't fail.
		// For commit, this won't be used, so its harmless.
		execBatchGasLimit = defaultExecGasLimit
	}

	homeChainID, err := chainsel.GetChainIDFromSelector(uint64(i.homeChainSelector))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain ID from chain selector %d: %w", i.homeChainSelector, err)
	}

	crOpts := make(map[cciptypes.ChainSelector]ccipcommon.ChainReaderProviderOpts)
	cwOpts := make(map[cciptypes.ChainSelector]ccipcommon.ChainWriterProviderOpts)
	for relayID, relayer := range i.relayers {
		chainID := relayID.ChainID
		relayChainFamily := relayID.Network
		chainDetails, err1 := chainsel.GetChainDetailsByChainIDAndFamily(chainID, relayChainFamily)
		if err1 != nil {
			return nil, nil, fmt.Errorf("failed to get chain selector from relay ID %s and family %s: %w", chainID, relayChainFamily, err)
		}

		chainSelector := cciptypes.ChainSelector(chainDetails.ChainSelector)
		crOpts[cciptypes.ChainSelector(chainDetails.ChainSelector)] = ccipcommon.ChainReaderProviderOpts{
			Lggr:            i.lggr,
			Relayer:         relayer,
			ChainID:         chainID,
			DestChainID:     destChainID,
			HomeChainID:     homeChainID,
			Ofc:             ofc,
			ChainSelector:   chainSelector,
			ChainFamily:     relayChainFamily,
			DestChainFamily: destChainFamily,
			Transmitters:    i.transmitters,
		}

		var solanaChainWriterConfigVersion *string
		if ofc.Execute != nil {
			solanaChainWriterConfigVersion = ofc.Execute.SolanaChainWriterConfigVersion
		}
		cwOpts[cciptypes.ChainSelector(chainDetails.ChainSelector)] = ccipcommon.ChainWriterProviderOpts{
			ChainID:                        chainID,
			Relayer:                        relayer,
			Transmitters:                   i.transmitters,
			ExecBatchGasLimit:              execBatchGasLimit,
			ChainFamily:                    relayChainFamily,
			OfframpProgramAddress:          config.Config.OfframpAddress,
			SolanaChainWriterConfigVersion: solanaChainWriterConfigVersion,
		}
	}
	return crOpts, cwOpts, nil
}

func (i *pluginOracleCreatorV2) createCCIPProviderWrappers(
	ctx context.Context,
	config cctypes.OCR3ConfigWithMeta,
	offRampAddrStr string,
	destFromAccounts []string,
	pluginType cctypes.PluginType,
	crOpts map[cciptypes.ChainSelector]ccipcommon.ChainReaderProviderOpts,
	cwOpts map[cciptypes.ChainSelector]ccipcommon.ChainWriterProviderOpts,
) (map[cciptypes.ChainSelector]ccipcommon.CCIPProviderWrapper, error) {
	ccipProviderWrappers := make(map[cciptypes.ChainSelector]ccipcommon.CCIPProviderWrapper)
	extraDataCodec := ccipcommon.GetExtraDataCodecRegistry()
	for relayID, relayer := range i.relayers {
		chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(relayID.ChainID, relayID.Network)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain selector from relay ID %s and family %s: %w", relayID.ChainID, relayID.Network, err)
		}
		chainSelector := cciptypes.ChainSelector(chainDetails.ChainSelector)

		// Get the registered wrapper factory for this chain family
		factory, exists := ccipcommon.GetCCIPProviderWrapperFactory(relayID.Network)
		if !exists {
			i.lggr.Warnw("no CCIPProviderWrapper factory registered for chain family, skipping",
				"chainFamily", relayID.Network,
				"chainSelector", chainSelector)
			continue
		}

		i.lggr.Debugw("creating CCIPProviderWrapper for chain family",
			"chainSelector", chainSelector,
			"chainFamily", relayID.Network)

		// Get transmitter for this relay
		transmitter := i.transmitters[relayID]
		if len(transmitter) == 0 {
			return nil, fmt.Errorf("transmitter list is empty for relay ID %s", relayID)
		}

		// Check if the transmitter string is a valid utf-8 string
		if !utf8.ValidString(transmitter[0]) {
			i.lggr.Errorw("transmitter contains invalid UTF-8",
				"transmitter", transmitter[0],
				"relayID.Network", relayID.Network,
				"chainSelector", chainSelector)
			return nil, fmt.Errorf("transmitter contains invalid UTF-8: %q", transmitter[0])
		}

		// Construct CCIPProviderArgs
		cargs := types.CCIPProviderArgs{
			PluginType:           pluginType,
			OffRampAddress:       config.Config.OfframpAddress,
			TransmitterAddress:   cciptypes.UnknownEncodedAddress(transmitter[0]),
			ExtraDataCodecBundle: extraDataCodec,
		}

		// Construct LegacyPluginServicesArgs
		largs := ccipcommon.LegacyPluginServicesArgs{
			ChainReaderOpts:      crOpts[chainSelector],
			ChainWriterOpts:      cwOpts[chainSelector],
			DestChainSelector:    config.Config.ChainSelector,
			DestFromAccounts:     destFromAccounts,
			OffRampAddressString: offRampAddrStr,
			RelayID:              relayID,
		}

		// Call the factory to create the wrapper
		wrapper, err := factory(ctx, i.lggr, chainSelector, relayer, cargs, largs, extraDataCodec)
		if err != nil {
			return nil, fmt.Errorf("failed to create CCIPProviderWrapper for chain %s (selector %d): %w",
				relayID.Network, chainSelector, err)
		}

		ccipProviderWrappers[chainSelector] = wrapper
	}

	return ccipProviderWrappers, nil
}
