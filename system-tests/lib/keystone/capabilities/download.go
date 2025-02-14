package capabilities

import (
	"os"

	libgithub "github.com/smartcontractkit/chainlink/system-tests/lib/github"
)

func DownloadCapabilityFromRelease(ghToken, version, assetFileName string) (string, error) {
	content, err := libgithub.DownloadGHAssetFromRelease("smartcontractkit", "capabilities", version, assetFileName, ghToken)
	if err != nil {
		return "", err
	}

	fileName := assetFileName
	file, err := os.Create(assetFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return "", err
	}

	return fileName, nil
}
