package capabilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func DownloadCapabilityFromRelease(ghToken, version, assetFileName string) (string, error) {
	ghClient := client.NewGithubClient(ghToken)
	content, response, err := ghClient.DownloadAssetFromRelease("smartcontractkit", "capabilities", version, assetFileName)
	if err != nil {
		if response != nil && response.StatusCode >= 400 {
			fmt.Printf("Request to GitHub failed with status code: %d\n", response.StatusCode)
			fmt.Printf("Response headers:\n")
			for header, values := range response.Header {
				valuesStr := strings.Join(values, ", ")
				fmt.Printf("Header: %s: %s\n", header, valuesStr)
			}
		}
		return "", err
	}

	fileName := assetFileName
	file, err := os.Create(assetFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func DefaultContainerDirectory(infraType types.InfraType) (string, error) {
	switch infraType {
	case types.CRIB:
		// chainlink user will always have access to this directory
		return "/home/chainlink", nil
	case types.Docker:
		// needs to match what CTFv2 uses by default, we should define a constant there and import it here
		return clnode.DefaultCapabilitiesDir, nil
	default:
		return "", fmt.Errorf("unknown infra type: %s", infraType)
	}
}
