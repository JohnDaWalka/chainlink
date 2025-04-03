package pricegetter

import (
	"context"
	"math/big"
	"sort"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/parseutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// Deprecated: not used, should be removed
type PipelineGetter struct {
	source        string
	runner        pipeline.Runner
	jobID         int32
	externalJobID uuid.UUID
	name          string
	lggr          logger.Logger
	destChainID   int64
}

func NewPipelineGetter(source string, runner pipeline.Runner, jobID int32, externalJobID uuid.UUID, name string, lggr logger.Logger, destChainID int64) (*PipelineGetter, error) {
	_, err := pipeline.Parse(source)
	if err != nil {
		return nil, err
	}

	return &PipelineGetter{
		source:        source,
		runner:        runner,
		jobID:         jobID,
		externalJobID: externalJobID,
		name:          name,
		lggr:          lggr,
		destChainID:   destChainID,
	}, nil
}

// FilterForConfiguredTokens implements the PriceGetter interface.
// It filters a list of token addresses for only those that have a pipeline job configured on the TokenPricesUSDPipeline
func (d *PipelineGetter) FilterConfiguredTokens(ctx context.Context, tokens []cciptypes.Address) (configured []cciptypes.Address, unconfigured []cciptypes.Address, err error) {
	lcSource := strings.ToLower(d.source)
	for _, tk := range tokens {
		lcToken := strings.ToLower(string(tk))
		if strings.Contains(lcSource, lcToken) {
			configured = append(configured, tk)
		} else {
			unconfigured = append(unconfigured, tk)
		}
	}
	return configured, unconfigured, nil
}

// GetJobSpecTokenPricesUSD gets all the tokens listed in the results.
// DEPRECATED: it does not support tokens with the same address on different chains.
// It treats every token as destination chain token.
func (d *PipelineGetter) GetJobSpecTokenPricesUSD(ctx context.Context) (map[TokenID]*big.Int, error) {
	prices, err := d.getPricesFromRunner(ctx)
	if err != nil {
		return nil, err
	}

	tokenPrices := make(map[TokenID]*big.Int)
	for tokenAddressStr, rawPrice := range prices {
		tokenAddress := ccipcalc.HexToAddress(tokenAddressStr)

		castedPrice, err := parseutil.ParseBigIntFromAny(rawPrice)
		if err != nil {
			return nil, err
		}

		tokenPrices[NewTokenID(tokenAddress, uint64(d.destChainID))] = castedPrice
	}

	return tokenPrices, nil
}

// GetTokenPrices is Deprecated: not used, should be removed.
// NOTE: Does not support tokens with the same address on different chains.
func (d *PipelineGetter) GetTokenPrices(ctx context.Context, tokenIDs []TokenID) (map[TokenID]*big.Int, error) {
	tokenChains := make(map[cciptypes.Address]uint64)
	tokenSet := mapset.NewSet[cciptypes.Address]()
	for _, tokenID := range tokenIDs {
		tokenSet.Add(tokenID.Address)
		tokenChains[tokenID.Address] = tokenID.ChainID
	}
	tokens := tokenSet.ToSlice()
	sort.Slice(tokens, func(i, j int) bool { return strings.ToLower(string(tokens[i])) < strings.ToLower(string(tokens[j])) })

	prices, err := d.getPricesFromRunner(ctx)
	if err != nil {
		return nil, err
	}

	providedTokensSet := mapset.NewSet(tokens...)
	tokenPrices := make(map[TokenID]*big.Int)
	for tokenAddressStr, rawPrice := range prices {
		tokenAddress := ccipcalc.HexToAddress(tokenAddressStr)
		tokenChain, ok := tokenChains[tokenAddress]
		if !ok {
			return nil, errors.Errorf("token %s not found in tokenChains map", tokenAddressStr)
		}

		castedPrice, err := parseutil.ParseBigIntFromAny(rawPrice)
		if err != nil {
			return nil, err
		}

		if providedTokensSet.Contains(tokenAddress) {
			tokenPrices[NewTokenID(tokenAddress, tokenChain)] = castedPrice
		}
	}

	// The mapping of token address to source of token price has to live offchain.
	// Best we can do is sanity check that the token price spec covers all our desired execution token prices.
	for _, token := range tokens {
		tokenChain, ok := tokenChains[token]
		if !ok {
			return nil, errors.Errorf("token %s not found in tokenChains map", token)
		}
		if _, ok := tokenPrices[NewTokenID(token, tokenChain)]; !ok {
			return nil, errors.Errorf("missing token %s from tokensForFeeCoin spec, got %v", token, prices)
		}
	}

	return tokenPrices, nil
}

func (d *PipelineGetter) getPricesFromRunner(ctx context.Context) (map[string]interface{}, error) {
	_, trrs, err := d.runner.ExecuteRun(ctx, pipeline.Spec{
		ID:           d.jobID,
		DotDagSource: d.source,
		CreatedAt:    time.Now(),
		JobID:        d.jobID,
		JobName:      d.name,
		JobType:      "",
	}, pipeline.NewVarsFrom(map[string]interface{}{}))
	if err != nil {
		return nil, err
	}
	finalResult := trrs.FinalResult()
	if finalResult.HasErrors() {
		return nil, errors.Errorf("error getting prices %v", finalResult.AllErrors)
	}
	if len(finalResult.Values) != 1 {
		return nil, errors.Errorf("invalid number of price results, expected 1 got %v", len(finalResult.Values))
	}
	prices, ok := finalResult.Values[0].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("expected map output of price pipeline, got %T", finalResult.Values[0])
	}

	return prices, nil
}

func (d *PipelineGetter) Close() error {
	return d.runner.Close()
}
