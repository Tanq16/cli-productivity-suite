package github

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/registry"
)

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	URL                string `json:"url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

func (c *Client) LatestRelease(repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error %d for %s: %s", resp.StatusCode, repo, string(body))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release for %s: %w", repo, err)
	}
	return &release, nil
}

func MatchAsset(release *Release, pattern registry.AssetPattern, osName, archName string) (*Asset, error) {
	osPattern, ok := pattern.OSPatterns[osName]
	if !ok {
		return nil, fmt.Errorf("no OS pattern for %s", osName)
	}
	archPattern, ok := pattern.ArchPatterns[archName]
	if !ok {
		return nil, fmt.Errorf("no arch pattern for %s", archName)
	}

	for _, asset := range release.Assets {
		name := asset.Name

		if !strings.Contains(name, osPattern) {
			continue
		}
		if !strings.Contains(name, archPattern) {
			continue
		}

		matchRequired := true
		for _, req := range pattern.RequiredSubstrings {
			if !strings.Contains(name, req) {
				matchRequired = false
				break
			}
		}
		if !matchRequired {
			continue
		}

		excluded := false
		for _, exc := range pattern.ExcludeSubstrings {
			if strings.Contains(name, exc) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		archiveFmt := pattern.ArchiveFormat
		if f, ok := pattern.OSArchiveFormats[osName]; ok {
			archiveFmt = f
		}
		if archiveFmt != "none" && archiveFmt != "" {
			validFormat := false
			switch archiveFmt {
			case "tar.gz":
				validFormat = strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz")
			case "tgz":
				validFormat = strings.HasSuffix(name, ".tgz") || strings.HasSuffix(name, ".tar.gz")
			case "tar.xz":
				validFormat = strings.HasSuffix(name, ".tar.xz")
			case "zip":
				validFormat = strings.HasSuffix(name, ".zip")
			default:
				validFormat = strings.HasSuffix(name, "."+archiveFmt)
			}
			if !validFormat {
				continue
			}
		}

		return &asset, nil
	}

	return nil, fmt.Errorf("no matching asset found for %s/%s in release %s", osName, archName, release.TagName)
}
