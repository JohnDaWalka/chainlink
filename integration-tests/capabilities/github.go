package capabilities_test

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/google/go-github/v41/github"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"golang.org/x/oauth2"
)

func downloadGHAssetFromLatestRelease(owner, repository, releaseType, assetName, ghToken string) ([]byte, error) {
	var content []byte
	if ghToken == "" {
		return content, errors.New("no github token provided")
	}

	if (releaseType == test_env.AUTOMATIC_LATEST_TAG) || (releaseType == test_env.AUTOMATIC_STABLE_LATEST_TAG) {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		ghClient := github.NewClient(tc)

		latestTags, _, err := ghClient.Repositories.ListReleases(context.Background(), owner, repository, &github.ListOptions{PerPage: 20})
		if err != nil {
			return content, errors.Wrapf(err, "failed to list releases for %s", repository)
		}

		var latestRelease *github.RepositoryRelease
		for _, tag := range latestTags {
			if releaseType == test_env.AUTOMATIC_STABLE_LATEST_TAG {
				if tag.Prerelease != nil && *tag.Prerelease {
					continue
				}
				if tag.Draft != nil && *tag.Draft {
					continue
				}
			}
			if tag.TagName != nil {
				latestRelease = tag
				break
			}
		}

		if latestRelease == nil {
			return content, errors.New("failed to find latest release with automatic tag: " + releaseType)
		}

		var assetId int64
		for _, asset := range latestRelease.Assets {
			if strings.Contains(asset.GetName(), assetName) {
				assetId = asset.GetID()
				break
			}
		}

		if assetId == 0 {
			return content, fmt.Errorf("failed to find asset %s for %s", assetName, *latestRelease.TagName)
		}

		asset, _, err := ghClient.Repositories.DownloadReleaseAsset(context.Background(), owner, repository, assetId, tc)
		if err != nil {
			return content, errors.Wrapf(err, "failed to download asset %s for %s", assetName, *latestRelease.TagName)
		}

		content, err = io.ReadAll(asset)
		if err != nil {
			return content, err
		}

		return content, nil

	}

	return content, errors.New("no automatic tag provided")
}
