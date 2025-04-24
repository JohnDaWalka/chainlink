package por

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

type ConfigFileInput struct {
	FundedAddress        common.Address
	BalanceReaderAddress common.Address
	FeedsConsumerAddress common.Address
	FeedID               string
	DataURL              string
	WriteTargetName      string
	ReadTargetName       string
	ContractReaderConfig string
	ContractName         string
	ContractMethod       string
}

func CreateConfigFile(in ConfigFileInput) (*os.File, error) {
	configFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create workflow config file")
	}

	cleanFeedID := strings.TrimPrefix(in.FeedID, "0x")
	feedLength := len(cleanFeedID)

	if feedLength < 32 {
		return nil, errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedLength)
	}

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	feedIDToUse := "0x" + cleanFeedID

	workflowConfig := libcrecli.PoRWorkflowConfig{
		FeedID:               feedIDToUse,
		URL:                  in.DataURL,
		ConsumerAddress:      in.FeedsConsumerAddress.Hex(),
		FundedAddress:        in.FundedAddress.Hex(),
		ContractAddress:      in.BalanceReaderAddress.Hex(),
		ContractName:         in.ContractName,
		ContractMethod:       in.ContractMethod,
		ContractReaderConfig: in.ContractReaderConfig,
		ReadTargetName:       in.ReadTargetName,
		WriteTargetName:      in.WriteTargetName,
	}

	configMarshalled, err := json.Marshal(workflowConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal workflow config")
	}

	_, err = configFile.Write(configMarshalled)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write workflow config file")
	}

	return configFile, nil
}
