package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type GitHubReleaseInstaller struct{}

func (g *GitHubReleaseInstaller) Check(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)
	release, err := gh.LatestRelease(tool.Repo)
	if err != nil {
		return current, "", err
	}
	return current, release.TagName, nil
}

func (g *GitHubReleaseInstaller) Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result {
	log.Debug().Str("package", "installer").Msgf("installing %s from %s", tool.Name, tool.Repo)

	release, err := gh.LatestRelease(tool.Repo)
	if err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("failed to fetch release: %w", err)}
	}

	currentVersion := st.ToolVersion(tool.Name)
	if currentVersion == release.TagName {
		return Result{Tool: tool.Name, Version: release.TagName, Skipped: true}
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

	destDir := p.ShellExecDir()
	if tool.Extension {
		destDir = p.ShellExtDir()
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	var binaryPath string

	archiveFmt := tool.Asset.ArchiveFormat
	if f, ok := tool.Asset.OSArchiveFormats[p.OS.String()]; ok {
		archiveFmt = f
	}
	if archiveFmt == "none" {
		binaryPath = archivePath
	} else {
		extractDir := filepath.Join(tmpDir, "extracted")
		if err := os.MkdirAll(extractDir, 0755); err != nil {
			return Result{Tool: tool.Name, Err: err}
		}

		switch archiveFmt {
		case "tar.gz":
			err = ExtractTarGz(archivePath, extractDir)
		case "tgz":
			err = ExtractTarGz(archivePath, extractDir)
		case "tar.xz":
			err = ExtractTarXz(archivePath, extractDir)
		case "zip":
			err = ExtractZip(archivePath, extractDir)
		default:
			err = fmt.Errorf("unknown archive format: %s", archiveFmt)
		}
		if err != nil {
			return Result{Tool: tool.Name, Err: fmt.Errorf("extract failed: %w", err)}
		}

		pattern := tool.Asset.BinaryPathInArchive
		if pattern == "" {
			pattern = tool.BinaryName
		}
		binaryPath, err = FindBinary(extractDir, pattern)
		if err != nil {
			return Result{Tool: tool.Name, Err: fmt.Errorf("binary not found in archive: %w", err)}
		}
	}

	destPath := filepath.Join(destDir, tool.BinaryName)
	if err := AtomicInstallBinary(binaryPath, destPath); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	st.SetToolVersion(tool.Name, release.TagName)

	wasUpdated := currentVersion != "" && currentVersion != release.TagName
	return Result{Tool: tool.Name, Version: release.TagName, WasUpdated: wasUpdated}
}
