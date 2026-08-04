package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type AppBundleInstaller struct{}

type bundleSource struct {
	version string
	fetch   func(tmpDir string) (string, error)
}

func (a *AppBundleInstaller) Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result {
	src, err := bundleSourceFor(tool, p, gh)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	destDir := filepath.Join(p.ShellAppsDir(), tool.Name)

	currentVersion := st.ToolVersion(tool.Name)
	if currentVersion == src.version {
		if _, statErr := os.Stat(destDir); statErr == nil {
			if err := runPostInstall(tool, p, destDir); err != nil {
				return Result{Tool: tool.Name, Err: err}
			}
			return Result{Tool: tool.Name, Version: src.version, Skipped: true}
		}
	}

	if err := os.MkdirAll(p.ShellAppsDir(), 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	// Staged inside ~/shell so the swap into place is a same-filesystem rename; /tmp is a separate mount on most Linux hosts.
	tmpDir, err := os.MkdirTemp(p.ShellDir(), "cps-"+tool.Name+"-*")
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	defer os.RemoveAll(tmpDir)

	archivePath, err := src.fetch(tmpDir)
	if err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("download failed: %w", err)}
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	archiveFmt := tool.Asset.ArchiveFormat
	if f, ok := tool.Asset.OSArchiveFormats[p.OS.String()]; ok {
		archiveFmt = f
	}
	if err := ExtractArchive(archivePath, extractDir, archiveFmt); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("extract failed: %w", err)}
	}

	if err := stageAndSwap(unwrapSingleDir(extractDir), destDir); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	// Version is recorded before the hook so a seeding failure retries on the next run instead of re-downloading the bundle.
	st.SetToolVersion(tool.Name, src.version)
	if err := runPostInstall(tool, p, destDir); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	wasUpdated := currentVersion != "" && currentVersion != src.version
	return Result{Tool: tool.Name, Version: src.version, WasUpdated: wasUpdated}
}

func bundleSourceFor(tool *registry.Tool, p platform.Platform, gh *github.Client) (bundleSource, error) {
	if tool.URL != "" {
		version, err := resolveVersion(tool, gh)
		if err != nil {
			return bundleSource{}, err
		}
		url := expandURL(tool.URL, version, p.OS.String(), p.Arch.String())
		return bundleSource{
			version: version,
			fetch: func(tmpDir string) (string, error) {
				archivePath := filepath.Join(tmpDir, tool.Name+"-archive")
				return archivePath, DownloadToFile(url, archivePath)
			},
		}, nil
	}

	release, err := gh.LatestRelease(tool.Repo)
	if err != nil {
		return bundleSource{}, fmt.Errorf("failed to fetch release: %w", err)
	}
	asset, err := github.MatchAsset(release, tool.Asset, p.OS.String(), p.Arch.String())
	if err != nil {
		return bundleSource{}, fmt.Errorf("no matching asset: %w", err)
	}
	downloadURL := asset.BrowserDownloadURL
	if tool.IsPrivate {
		downloadURL = asset.URL
	}
	return bundleSource{
		version: release.TagName,
		fetch: func(tmpDir string) (string, error) {
			return gh.DownloadFile(downloadURL, tmpDir, asset.Name)
		},
	}, nil
}
