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
	StreamID        string
	Name            string
	StreamType      jobs.StreamType
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
		spec, err := generateJobSpec(s)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create stream spec: %w", err)
		}

		for _, n := range oracleNodes {
			renderedSpec, err := spec.MarshalTOML()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal stream spec: %w", err)
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

func generateJobSpec(ssc StreamSpecConfig) (spec *jobs.StreamJobSpec, err error) {
	spec = &jobs.StreamJobSpec{
		Base: jobs.Base{
			Name:          fmt.Sprintf("%s | %d", ssc.Name, ssc.StreamID),
			Type:          jobs.JobSpecTypeStream,
			SchemaVersion: 1,
			ExternalJobID: uuid.New(),
		},
		StreamID: ssc.StreamID,
	}

	base := jobs.BaseObservationSource{
		Datasources:   generateDatasources(ssc),
		AllowedFaults: ssc.AllowedFaults,
	}

	switch ssc.StreamType {
	case jobs.StreamTypeQuote:
		err = spec.SetObservationSource(jobs.QuoteObservationSource{
			BaseObservationSource: base,
			Bid: jobs.ReportFieldLLO{
				ResultPath: "data,bid", // TODO maybe "data,result" for all?
			},
			Benchmark: jobs.ReportFieldLLO{
				ResultPath: "data,mid",
			},
			Ask: jobs.ReportFieldLLO{
				ResultPath: "data,ask",
			},
		})
	case jobs.StreamTypeMedian:
		err = spec.SetObservationSource(jobs.MedianObservationSource{
			BaseObservationSource: base,
			Benchmark: jobs.ReportFieldLLO{
				ResultPath: "data,mid",
			},
		})
		// TODO Add the rest of the stream types.
	default:
		return nil, fmt.Errorf("unsupported stream type: %s", ssc.StreamType)
	}

	return spec, err
}

func generateDatasources(ssc StreamSpecConfig) []jobs.Datasource {
	dss := make([]jobs.Datasource, len(ssc.APIs))
	params := ssc.EARequestParams
	for i, api := range ssc.APIs {
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
