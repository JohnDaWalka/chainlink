package github

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func DownloadGHAssetFromRelease(owner, repository, releaseTag, assetName, ghToken string) ([]byte, error) {
	var content []byte
	if ghToken == "" {
		return content, errors.New("no github token provided")
	}

	// assuming 180s is enough to fetch releases, find the asset we need and download it
	// some assets might be 30+ MB, so we need to give it some time (for really slow connections)
	ctx, cancelFn := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancelFn()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)

	ghReleases, _, err := ghClient.Repositories.ListReleases(ctx, owner, repository, &github.ListOptions{PerPage: 20})
	if err != nil {
		return content, errors.Wrapf(err, "failed to list releases for %s", repository)
	}

	var ghRelease *github.RepositoryRelease
	for _, release := range ghReleases {
		if release.TagName == nil {
			continue
		}

		if *release.TagName == releaseTag {
			ghRelease = release
			break
		}
	}

	if ghRelease == nil {
		return content, errors.New("failed to find release with tag: " + releaseTag)
	}

	var assetID int64
	for _, asset := range ghRelease.Assets {
		if strings.Contains(asset.GetName(), assetName) {
			assetID = asset.GetID()
			break
		}
	}

	if assetID == 0 {
		return content, fmt.Errorf("failed to find asset %s for %s", assetName, *ghRelease.TagName)
	}

	asset, _, err := ghClient.Repositories.DownloadReleaseAsset(ctx, owner, repository, assetID, tc)
	if err != nil {
		return content, errors.Wrapf(err, "failed to download asset %s for %s", assetName, *ghRelease.TagName)
	}

	content, err = io.ReadAll(asset)
	if err != nil {
		return content, err
	}

	return content, nil
}
