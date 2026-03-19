package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type DirectDownloadInstaller struct{}

func (d *DirectDownloadInstaller) Check(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)
	latest, err = d.fetchVersion(tool, gh)
	if err != nil {
		return current, "", err
	}
	return current, latest, nil
}

func (d *DirectDownloadInstaller) Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result {
	log.Debug().Str("package", "installer").Msgf("installing %s via direct download", tool.Name)

	version, err := d.fetchVersion(tool, gh)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	currentVersion := st.ToolVersion(tool.Name)
	if currentVersion == version && version != "" {
		return Result{Tool: tool.Name, Version: version, Skipped: true}
	}

	versionBare := strings.TrimPrefix(version, "v")
	url := tool.URL
	url = strings.ReplaceAll(url, "{version}", version)
	url = strings.ReplaceAll(url, "{version_bare}", versionBare)
	url = strings.ReplaceAll(url, "{os}", p.OS.String())
	url = strings.ReplaceAll(url, "{arch}", p.Arch.String())

	destDir := p.ShellExecDir()
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	archiveFormat := tool.Asset.ArchiveFormat
	if archiveFormat != "" && archiveFormat != "none" {
		return d.installArchive(tool, url, version, currentVersion, destDir, archiveFormat, st)
	}

	resp, err := httpGet(url)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{Tool: tool.Name, Err: fmt.Errorf("download failed: HTTP %d from %s", resp.StatusCode, url)}
	}

	destPath := filepath.Join(destDir, tool.BinaryName)
	f, err := os.Create(destPath)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	if err := os.Chmod(destPath, 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	st.SetToolVersion(tool.Name, version)
	wasUpdated := currentVersion != "" && currentVersion != version
	return Result{Tool: tool.Name, Version: version, WasUpdated: wasUpdated}
}

func (d *DirectDownloadInstaller) fetchVersion(tool *registry.Tool, gh *github.Client) (string, error) {
	if tool.StableURL != "" {
		resp, err := httpGet(tool.StableURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("fetching version from %s: HTTP %d", tool.StableURL, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(body)), nil
	}
	if tool.Repo != "" {
		release, err := gh.LatestRelease(tool.Repo)
		if err != nil {
			return "", err
		}
		return release.TagName, nil
	}
	return "", fmt.Errorf("no version source configured for %s", tool.Name)
}

func (d *DirectDownloadInstaller) installArchive(tool *registry.Tool, url, version, currentVersion, destDir, archiveFormat string, st *state.State) Result {
	tmpDir, err := os.MkdirTemp("", "cps-"+tool.Name+"-*")
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, tool.Name+"."+archiveFormat)
	if err := DownloadToFile(url, archivePath); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("download failed: %w", err)}
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	switch archiveFormat {
	case "tar.gz", "tgz":
		err = ExtractTarGz(archivePath, extractDir)
	case "tar.xz":
		err = ExtractTarXz(archivePath, extractDir)
	case "zip":
		err = ExtractZip(archivePath, extractDir)
	default:
		err = fmt.Errorf("unknown archive format: %s", archiveFormat)
	}
	if err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("extract failed: %w", err)}
	}

	pattern := tool.Asset.BinaryPathInArchive
	if pattern == "" {
		pattern = tool.BinaryName
	}
	binaryPath, err := FindBinary(extractDir, pattern)
	if err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("binary not found in archive: %w", err)}
	}

	destPath := filepath.Join(destDir, tool.BinaryName)
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	if err := os.WriteFile(destPath, data, 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	st.SetToolVersion(tool.Name, version)
	wasUpdated := currentVersion != "" && currentVersion != version
	return Result{Tool: tool.Name, Version: version, WasUpdated: wasUpdated}
}
