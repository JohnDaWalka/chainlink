package por

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Config interface {
	FeedIDGetter
	FeedIDSetter
}

type FeedIDGetter interface {
	GetFeedID() string
}

type FeedIDSetter interface {
	SetFeedID(string)
}

func CreateConfigFile(
	cfg Config,
) (*os.File, error) {
	fID, err := handleFeedID(cfg.GetFeedID())
	if err != nil {
		return nil, err
	}

	cfg.SetFeedID(fID)
	b, err := marshalConfig(cfg)
	if err != nil {
		return nil, err
	}

	return writeConfigFile(b)
}

func handleFeedID(feedID string) (string, error) {
	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedID)

	if feedLength < 32 {
		return "", errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedLength)
	}

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	feedIDToUse := "0x" + cleanFeedID
	return feedIDToUse, nil
}

func marshalConfig[T any](input T) ([]byte, error) {
	configMarshalled, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal workflow config")
	}
	return configMarshalled, nil
}

func writeConfigFile(configMarshalled []byte) (*os.File, error) {
	configFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create workflow config file")
	}

	_, err = configFile.Write(configMarshalled)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write workflow config file")
	}

	return configFile, nil
}
