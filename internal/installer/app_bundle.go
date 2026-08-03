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

func (a *AppBundleInstaller) Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result {
	release, err := gh.LatestRelease(tool.Repo)
	if err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("failed to fetch release: %w", err)}
	}

	destDir := filepath.Join(p.ShellAppsDir(), tool.Name)

	currentVersion := st.ToolVersion(tool.Name)
	if currentVersion == release.TagName {
		if _, statErr := os.Stat(destDir); statErr == nil {
			return Result{Tool: tool.Name, Version: release.TagName, Skipped: true}
		}
	}

	asset, err := github.MatchAsset(release, tool.Asset, p.OS.String(), p.Arch.String())
	if err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("no matching asset: %w", err)}
	}

	tmpDir, err := os.MkdirTemp("", "cps-"+tool.Name+"-*")
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	defer os.RemoveAll(tmpDir)

	downloadURL := asset.BrowserDownloadURL
	if tool.IsPrivate {
		downloadURL = asset.URL
	}

	archivePath, err := gh.DownloadFile(downloadURL, tmpDir, asset.Name)
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

	if err := os.MkdirAll(p.ShellAppsDir(), 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	if err := stageAndSwap(unwrapSingleDir(extractDir), destDir); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	st.SetToolVersion(tool.Name, release.TagName)

	wasUpdated := currentVersion != "" && currentVersion != release.TagName
	return Result{Tool: tool.Name, Version: release.TagName, WasUpdated: wasUpdated}
}
