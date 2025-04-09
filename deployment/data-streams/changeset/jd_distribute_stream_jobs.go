package changeset

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

var _ deployment.ChangeSetV2[CsDistributeStreamJobSpecsConfig] = CsDistributeStreamJobSpecs{}

type CsDistributeStreamJobSpecsConfig struct {
	ChainSelectorEVM uint64
	Filter           *jd.ListFilter
	Streams          []StreamSpecConfig
}

type StreamSpecConfig struct {
	StreamID   string
	Name       string
	StreamType jobs.StreamType
	// ReportFields should be QuoteReportFields, MedianReportFields, etc., based on the stream type.
	ReportFields    jobs.ReportFields
	EARequestParams EARequestParams
	APIs            []string
	AllowedFaults   int
}

type EARequestParams struct {
	Endpoint string `json:"endpoint"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type CsDistributeStreamJobSpecs struct{}

func (CsDistributeStreamJobSpecs) Apply(e deployment.Environment, cfg CsDistributeStreamJobSpecsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultJobSpecsTimeout)
	defer cancel()

	// Add a label to the job spec to identify the related DON
	labels := append([]*ptypes.Label(nil),
		&ptypes.Label{
			Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
		})

	oracleNodes, err := jd.FetchDONOraclesFromJD(ctx, e.Offchain, cfg.Filter)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get workflow don nodes: %w", err)
	}

	var proposals []*jobv1.ProposeJobRequest
	for _, s := range cfg.Streams {
		for _, n := range oracleNodes {
			spec, err := generateJobSpec(s)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create stream job spec: %w", err)
			}
			renderedSpec, err := spec.MarshalTOML()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal stream job spec: %w", err)
			}

			proposals = append(proposals, &jobv1.ProposeJobRequest{
				NodeId: n.Id,
				Spec:   string(renderedSpec),
				Labels: labels,
			})
		}
	}

	proposedJobs, err := proposeAllOrNothing(ctx, e.Offchain, proposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose all oracle jobs: %w", err)
	}

	return deployment.ChangesetOutput{
		Jobs: proposedJobs,
	}, nil
}

func generateJobSpec(cc StreamSpecConfig) (spec *jobs.StreamJobSpec, err error) {
	spec = &jobs.StreamJobSpec{
		Base: jobs.Base{
			Name:          fmt.Sprintf("%s | %s", cc.Name, cc.StreamID),
			Type:          jobs.JobSpecTypeStream,
			SchemaVersion: 1,
			ExternalJobID: uuid.New(),
		},
		StreamID: cc.StreamID,
	}

	base := jobs.BaseObservationSource{
		Datasources:   generateDatasources(cc),
		AllowedFaults: cc.AllowedFaults,
	}

	err = spec.SetObservationSource(base, cc.ReportFields)

	return spec, err
}

func generateDatasources(cc StreamSpecConfig) []jobs.Datasource {
	dss := make([]jobs.Datasource, len(cc.APIs))
	params := cc.EARequestParams
	for i, api := range cc.APIs {
		dss[i] = jobs.Datasource{
			BridgeName: fmt.Sprintf("bridge-%s", api),
			ReqData:    fmt.Sprintf(`{"data":{"endpoint":"%s","from":"%s","to":"%s"}}`, params.Endpoint, params.From, params.To),
		}
	}
	return dss
}

func (f CsDistributeStreamJobSpecs) VerifyPreconditions(_ deployment.Environment, config CsDistributeStreamJobSpecsConfig) error {
	if config.ChainSelectorEVM == 0 {
		return errors.New("chain selector is required")
	}
	if config.Filter == nil {
		return errors.New("filter is required")
	}
	if config.Streams == nil || len(config.Streams) == 0 {
		return errors.New("streams are required")
	}

	return nil
}
