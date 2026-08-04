package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type DirectDownloadInstaller struct{}

func (d *DirectDownloadInstaller) Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result {
	version, err := resolveVersion(tool, gh)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	destDir := p.ShellExtDir()

	currentVersion := st.ToolVersion(tool.Name)
	if currentVersion == version && version != "" {
		if _, statErr := os.Stat(filepath.Join(destDir, tool.BinaryName)); statErr == nil {
			return Result{Tool: tool.Name, Version: version, Skipped: true}
		}
	}

	url := expandURL(tool.URL, version, p.OS.String(), p.Arch.String())

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
	tmp, err := os.CreateTemp(destDir, ".cps-tmp-*")
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	tmpPath := tmp.Name()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return Result{Tool: tool.Name, Err: err}
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return Result{Tool: tool.Name, Err: err}
	}
	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return Result{Tool: tool.Name, Err: err}
	}
	if err := os.Rename(tmpPath, destPath); err != nil {
		os.Remove(tmpPath)
		return Result{Tool: tool.Name, Err: err}
	}

	st.SetToolVersion(tool.Name, version)
	wasUpdated := currentVersion != "" && currentVersion != version
	return Result{Tool: tool.Name, Version: version, WasUpdated: wasUpdated}
}

func resolveVersion(tool *registry.Tool, gh *github.Client) (string, error) {
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
		if tool.VersionRegex == "" {
			return strings.TrimSpace(string(body)), nil
		}
		re, err := regexp.Compile(tool.VersionRegex)
		if err != nil {
			return "", fmt.Errorf("invalid version pattern for %s: %w", tool.Name, err)
		}
		m := re.FindSubmatch(body)
		if len(m) < 2 {
			return "", fmt.Errorf("version pattern %q matched nothing at %s", tool.VersionRegex, tool.StableURL)
		}
		return strings.TrimSpace(string(m[1])), nil
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

func expandURL(tmpl, version, osStr, archStr string) string {
	r := strings.NewReplacer(
		"{version}", version,
		"{version_bare}", strings.TrimPrefix(version, "v"),
		"{os}", osStr,
		"{arch}", archStr,
	)
	return r.Replace(tmpl)
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
	if err := AtomicInstallBinary(binaryPath, destPath); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	st.SetToolVersion(tool.Name, version)
	wasUpdated := currentVersion != "" && currentVersion != version
	return Result{Tool: tool.Name, Version: version, WasUpdated: wasUpdated}
}
