package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/blang/semver"
)

type githubSource struct {
	repo             string
	artifactPrefix   string
	releaseTagPrefix string
}

type githubRelease struct {
	Name       string        `json:"name"`
	TagName    string        `json:"tag_name"`
	Body       string        `json:"body"`
	Prerelease bool          `json:"prerelease"`
	Assets     []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name string `json:"name"`
}

// NewGitHubSource creates a new update source which consumes releases from GitHub.
func NewGitHubSource(repo, releaseTagPrefix, artifactPrefix string) Source {
	return &githubSource{
		repo:             repo,
		releaseTagPrefix: releaseTagPrefix,
		artifactPrefix:   artifactPrefix,
	}
}

func (s *githubSource) Releases() ([]Release, error) {
	releases, err := s.getReleases()

	if err != nil {
		return nil, err
	}

	out := []Release{}

	for _, release := range releases {
		if !strings.HasPrefix(release.TagName, s.releaseTagPrefix) {
			continue
		}

		version, err := semver.Parse(release.TagName[len(s.releaseTagPrefix):])
		if err != nil {
			return nil, err
		}

		out = append(out, Release{
			ID:        release.TagName,
			Changelog: release.Body,
			Version:   version,
			Variants:  s.getVariants(&release),
		})
	}

	return out, nil
}

func (s *githubSource) Download(release *Release, variant *Variant) (io.ReadCloser, error) {
	if variant == nil {
		return nil, fmt.Errorf("github: you must provide a variant to download")
	}

	variant = release.GetVariant(variant)
	if variant == nil {
		return nil, fmt.Errorf("github: requested variant was not available in this release")
	}

	res, err := http.DefaultClient.Get(fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", s.repo, release.ID, variant.ID))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, fmt.Errorf("github: Failed to download release artifact: %d %s", res.StatusCode, res.Status)
	}

	return res.Body, nil
}

func (s *githubSource) getReleases() ([]githubRelease, error) {
	res, err := http.DefaultClient.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", s.repo))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("github: %s", res.Status)
	}

	var releases []githubRelease
	err = json.NewDecoder(res.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (s *githubSource) getVariants(r *githubRelease) []Variant {
	out := []Variant{}
	for _, v := range r.Assets {
		if strings.HasPrefix(v.Name, s.artifactPrefix) {
			variant := Variant{
				ID: v.Name,
			}

			variantName := v.Name[len(s.artifactPrefix):]
			if strings.ContainsRune(variantName, '.') {
				variantName = variantName[:strings.IndexRune(variantName, '.')]
			}

			parts := strings.Split(variantName, "-")
			if len(parts) == 2 {
				variant.Platform = parts[0]
				variant.Arch = parts[1]
			} else {
				variant.Platform = runtime.GOOS
				variant.Arch = runtime.GOARCH
			}

			out = append(out, variant)
		}
	}

	return out
}
