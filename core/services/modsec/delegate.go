package modsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecexecutor"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecverifier"
)

type Config interface {
	Feature() config.Feature
}

type RelayGetter interface {
	Get(types.RelayID) (loop.Relayer, error)
	GetIDToRelayerMap() map[types.RelayID]loop.Relayer
}

type Keystore[K keystore.Key] interface {
	GetAll() ([]K, error)
}

type Delegate struct {
	cfg          Config
	lggr         logger.Logger
	keystore     keystore.Master
	ds           sqlutil.DataSource
	evmConfigs   toml.EVMConfigs
	legacyChains legacyevm.LegacyChainContainer

	isNewlyCreatedJob bool
}

func NewDelegate(
	cfg Config,
	lggr logger.Logger,
	keystore keystore.Master,
	ds sqlutil.DataSource,
	evmConfigs toml.EVMConfigs,
	legacyChains legacyevm.LegacyChainContainer,
) *Delegate {
	return &Delegate{
		cfg:          cfg,
		lggr:         lggr,
		keystore:     keystore,
		ds:           ds,
		evmConfigs:   evmConfigs,
		legacyChains: legacyChains,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Modsec
}

func (d *Delegate) BeforeJobCreated(job.Job) {
	// This is only called first time the job is created
	d.isNewlyCreatedJob = true
}

func validate(spec *job.ModsecSpec) error {
	sourceChainID := spec.SourceChainID
	sourceChainFamily := spec.SourceChainFamily

	if sourceChainID == "" || sourceChainFamily == "" {
		return fmt.Errorf("source chain id (%s) or family (%s) is empty", sourceChainID, sourceChainFamily)
	}

	if sourceChainFamily != string(chaintype.EVM) {
		return fmt.Errorf("source chain family (%s) is not an EVM chain", sourceChainFamily)
	}

	destChainID := spec.DestChainID
	destChainFamily := spec.DestChainFamily

	if destChainID == "" || destChainFamily == "" {
		return fmt.Errorf("dest chain id (%s) or family (%s) is empty", destChainID, destChainFamily)
	}

	if destChainFamily != string(chaintype.EVM) {
		return fmt.Errorf("dest chain family (%s) is not an EVM chain", destChainFamily)
	}

	if sourceChainID == destChainID {
		return fmt.Errorf("source chain id (%s) and dest chain id (%s) are the same", sourceChainID, destChainID)
	}

	if spec.OnRampAddress == "" {
		return fmt.Errorf("on ramp address is empty")
	}

	if !common.IsHexAddress(spec.OnRampAddress) {
		return fmt.Errorf("on ramp address (%s) is not a valid address", spec.OnRampAddress)
	}

	if spec.OffRampAddress == "" {
		return fmt.Errorf("off ramp address is empty")
	}

	if !common.IsHexAddress(spec.OffRampAddress) {
		return fmt.Errorf("off ramp address (%s) is not a valid address", spec.OffRampAddress)
	}

	if spec.CCIPMessageSentEventSig == "" {
		return fmt.Errorf("ccip message sent event sig is empty")
	}

	if len(common.FromHex(spec.CCIPMessageSentEventSig)) != 32 {
		return fmt.Errorf("ccip message sent event sig is not 32 bytes")
	}

	if spec.StorageEndpoint == "" {
		return fmt.Errorf("storage endpoint is empty")
	}

	if !strings.HasPrefix(spec.StorageEndpoint, "http") {
		return fmt.Errorf("storage endpoint (%s) is not a valid http endpoint", spec.StorageEndpoint)
	}

	if spec.StorageType == "" {
		return fmt.Errorf("storage type is empty")
	}

	if spec.StorageType != "std" {
		return fmt.Errorf("storage type (%s) is not supported", spec.StorageType)
	}

	return nil
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	if spec.ModsecSpec == nil {
		return nil, errors.New("modsec spec is nil")
	}

	marshalledJob, err := json.MarshalIndent(spec.ModsecSpec, "", " ")
	if err != nil {
		return nil, err
	}
	d.lggr.Debugw("Creating services for modsec job spec", "job", string(marshalledJob))

	if !d.cfg.Feature().LogPoller() {
		return nil, errors.New("log poller must be enabled to run modsec")
	}

	if err = validate(spec.ModsecSpec); err != nil {
		return nil, err
	}

	sourceChain, err := d.legacyChains.Get(spec.ModsecSpec.SourceChainID)
	if err != nil {
		return nil, err
	}

	legacySourceChain, ok := sourceChain.(legacyevm.Chain)
	if !ok {
		return nil, fmt.Errorf("source chain (%s) is not an EVM chain", spec.ModsecSpec.SourceChainID)
	}

	destChain, err := d.legacyChains.Get(spec.ModsecSpec.DestChainID)
	if err != nil {
		return nil, err
	}

	legacyDestChain, ok := destChain.(legacyevm.Chain)
	if !ok {
		return nil, fmt.Errorf("dest chain (%s) is not an EVM chain", spec.ModsecSpec.DestChainID)
	}

	storageClient := modsecstorage.NewStdClient(spec.ModsecSpec.StorageEndpoint)

	verifier := modsecverifier.New(
		d.lggr,
		legacySourceChain.LogPoller(),
		spec.ModsecSpec.CCIPMessageSentEventSig,
		spec.ModsecSpec.OnRampAddress,
		storageClient,
	)

	relayer := modsecexecutor.New(
		d.lggr,
		legacyDestChain.LogPoller(),
		legacyDestChain.TxManager(),
		spec.ModsecSpec.CCIPMessageSentEventSig,
		spec.ModsecSpec.OnRampAddress,
		spec.ModsecSpec.OffRampAddress,
		storageClient,
	)

	services = append(services, verifier, relayer)

	return services, nil
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, spec job.Job) error {
	// TODO: shut down needed services?
	return nil
}
