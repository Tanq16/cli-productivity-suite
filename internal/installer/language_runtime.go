package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type LanguageRuntimeInstaller struct{}

func (l *LanguageRuntimeInstaller) Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result {
	switch tool.Name {
	case "go-sdk":
		return l.installGo(p, st)
	case "uv":
		return l.installUVOnly(p, gh, st)
	case "python":
		return l.installPython(p, gh, st)
	case "rust":
		return l.installRust(p, st)
	case "fnm":
		return l.installFnmOnly(p, gh, st)
	case "node":
		return l.installNode(p, gh, st)
	default:
		return Result{Tool: tool.Name, Err: fmt.Errorf("unknown runtime: %s", tool.Name)}
	}
}

func (l *LanguageRuntimeInstaller) installGo(p platform.Platform, st *state.State) Result {
	type goDL struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
		Files   []struct {
			Filename string `json:"filename"`
			OS       string `json:"os"`
			Arch     string `json:"arch"`
			Kind     string `json:"kind"`
		} `json:"files"`
	}

	resp, err := httpGet("https://go.dev/dl/?mode=json")
	if err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("go.dev/dl API returned HTTP %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}

	var releases []goDL
	if err := json.Unmarshal(body, &releases); err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}

	var downloadURL, version string
	for _, r := range releases {
		if !r.Stable {
			continue
		}
		for _, f := range r.Files {
			if f.OS == p.OS.String() && f.Arch == p.Arch.String() && f.Kind == "archive" {
				downloadURL = fmt.Sprintf("https://go.dev/dl/%s", f.Filename)
				version = r.Version
				break
			}
		}
		if downloadURL != "" {
			break
		}
	}

	if downloadURL == "" {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("no Go download found for %s/%s", p.OS, p.Arch)}
	}

	// Temp dir inside p.ShellDir() so os.Rename stays on the same filesystem (avoids EXDEV on Linux tmpfs /tmp).
	tmpDir, err := os.MkdirTemp(p.ShellDir(), "cps-go-*")
	if err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "go.tar.gz")
	if err := DownloadToFile(downloadURL, tarPath); err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}

	if err := ExtractTarGz(tarPath, tmpDir); err != nil {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("extract Go SDK failed: %w", err)}
	}

	goDir := filepath.Join(p.ShellDir(), "go-sdk")
	os.RemoveAll(goDir)
	if err := os.Rename(filepath.Join(tmpDir, "go"), goDir); err != nil {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("move Go SDK failed: %w", err)}
	}

	st.SetToolVersion("go-sdk", version)
	return Result{Tool: "go-sdk", Version: version}
}

func (l *LanguageRuntimeInstaller) installRust(p platform.Platform, st *state.State) Result {
	var archStr string
	switch p.Arch {
	case platform.AMD64:
		archStr = "x86_64"
	case platform.ARM64:
		archStr = "aarch64"
	}

	var target string
	switch p.OS {
	case platform.Darwin:
		target = archStr + "-apple-darwin"
	default:
		target = archStr + "-unknown-linux-gnu"
	}

	rustDir := filepath.Join(p.ShellDir(), "rust")
	rustupHome := filepath.Join(rustDir, ".rustup")
	cargoHome := filepath.Join(rustDir, ".cargo")

	if err := os.MkdirAll(rustDir, 0755); err != nil {
		return Result{Tool: "rust", Err: fmt.Errorf("failed to create rust dir: %w", err)}
	}

	url := fmt.Sprintf("https://static.rust-lang.org/rustup/dist/%s/rustup-init", target)

	tmpDir, err := os.MkdirTemp("", "cps-rust-*")
	if err != nil {
		return Result{Tool: "rust", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	initPath := filepath.Join(tmpDir, "rustup-init")
	if err := DownloadToFile(url, initPath); err != nil {
		return Result{Tool: "rust", Err: fmt.Errorf("download rustup-init failed: %w", err)}
	}
	if err := os.Chmod(initPath, 0755); err != nil {
		return Result{Tool: "rust", Err: err}
	}

	cmd := exec.Command(initPath, "-y", "--no-modify-path")
	cmd.Env = append(os.Environ(),
		"RUSTUP_HOME="+rustupHome,
		"CARGO_HOME="+cargoHome,
	)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: "rust", Err: fmt.Errorf("rustup-init failed: %w", err)}
	}

	rustcPath := filepath.Join(cargoHome, "bin", "rustc")
	verCmd := exec.Command(rustcPath, "--version")
	verCmd.Env = append(os.Environ(),
		"RUSTUP_HOME="+rustupHome,
		"CARGO_HOME="+cargoHome,
	)
	out, err := verCmd.Output()
	version := "stable"
	if err == nil {
		parts := strings.Fields(strings.TrimSpace(string(out)))
		if len(parts) >= 2 {
			version = parts[1]
		}
	}

	st.SetToolVersion("rust", version)
	return Result{Tool: "rust", Version: version}
}

func (l *LanguageRuntimeInstaller) latestPythonCycle() (string, error) {
	resp, err := httpGet("https://endoflife.date/api/python.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var releases []struct {
		Cycle string `json:"cycle"`
		EOL   any    `json:"eol"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}
	now := time.Now().Format("2006-01-02")
	for _, r := range releases {
		switch eol := r.EOL.(type) {
		case bool:
			if !eol {
				return r.Cycle, nil
			}
		case string:
			if eol > now {
				return r.Cycle, nil
			}
		}
	}
	return "", fmt.Errorf("no active Python release found")
}

func (l *LanguageRuntimeInstaller) findUVAsset(assets []github.Asset, p platform.Platform) string {
	var osStr, archStr string
	switch p.OS {
	case platform.Darwin:
		osStr = "apple"
	default:
		osStr = "linux"
	}
	switch p.Arch {
	case platform.AMD64:
		archStr = "x86_64"
	case platform.ARM64:
		archStr = "aarch64"
	}
	for _, a := range assets {
		if strings.Contains(a.Name, osStr) && strings.Contains(a.Name, archStr) &&
			strings.HasSuffix(a.Name, ".tar.gz") && !strings.Contains(a.Name, "musl") {
			return a.BrowserDownloadURL
		}
	}
	return ""
}

func (l *LanguageRuntimeInstaller) downloadUV(p platform.Platform, gh *github.Client) (uvPath string, tagName string, err error) {
	release, err := gh.LatestRelease("astral-sh/uv")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch uv release: %w", err)
	}

	downloadURL := l.findUVAsset(release.Assets, p)
	if downloadURL == "" {
		return "", "", fmt.Errorf("no uv asset found for %s/%s", p.OS, p.Arch)
	}

	tmpDir, err := os.MkdirTemp("", "cps-uv-*")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "uv.tar.gz")
	if err := DownloadToFile(downloadURL, tarPath); err != nil {
		return "", "", fmt.Errorf("download uv failed: %w", err)
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return "", "", err
	}
	if err := ExtractTarGz(tarPath, extractDir); err != nil {
		return "", "", fmt.Errorf("extract uv failed: %w", err)
	}

	uvExtracted, err := FindBinary(extractDir, "*/uv")
	if err != nil {
		return "", "", fmt.Errorf("uv binary not found in archive: %w", err)
	}

	destPath := filepath.Join(p.ShellExtDir(), "uv")
	if err := AtomicInstallBinary(uvExtracted, destPath); err != nil {
		return "", "", fmt.Errorf("install uv binary failed: %w", err)
	}

	return destPath, release.TagName, nil
}

func (l *LanguageRuntimeInstaller) installUVOnly(p platform.Platform, gh *github.Client, st *state.State) Result {
	_, tagName, err := l.downloadUV(p, gh)
	if err != nil {
		return Result{Tool: "uv", Err: err}
	}
	st.SetToolVersion("uv", tagName)
	return Result{Tool: "uv", Version: tagName}
}

func (l *LanguageRuntimeInstaller) installPython(p platform.Platform, gh *github.Client, st *state.State) Result {
	uvPath, uvTag, err := l.downloadUV(p, gh)
	if err != nil {
		return Result{Tool: "python", Err: err}
	}

	version, err := l.latestPythonCycle()
	if err != nil {
		return Result{Tool: "python", Err: err}
	}

	venvPath := filepath.Join(p.ShellDir(), "py-default")
	oldPythonPath := filepath.Join(venvPath, "bin", "python")
	var oldFullVersion string
	if out, err := exec.Command(oldPythonPath, "--version").Output(); err == nil {
		parts := strings.Fields(strings.TrimSpace(string(out)))
		if len(parts) >= 2 {
			oldFullVersion = parts[1]
		}
	}

	cmd := exec.Command(uvPath, "python", "install", version)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: "python", Err: fmt.Errorf("uv python install failed: %w", err)}
	}

	pipPath := filepath.Join(venvPath, "bin", "pip")
	var frozen []byte
	if _, err := os.Stat(pipPath); err == nil {
		frozen, _ = exec.Command(pipPath, "freeze").Output()
	}

	venvCmd := exec.Command(uvPath, "venv", "--python", version, "--clear", venvPath)
	if err := utils.RunCmd(venvCmd); err != nil {
		return Result{Tool: "python", Err: fmt.Errorf("uv venv create failed: %w", err)}
	}

	if len(frozen) > 0 {
		pipNew := filepath.Join(venvPath, "bin", "pip")
		installCmd := exec.Command(pipNew, "install", "--quiet", "-r", "/dev/stdin")
		installCmd.Stdin = strings.NewReader(string(frozen))
		if err := utils.RunCmd(installCmd); err != nil {
			return Result{Tool: "python", Err: fmt.Errorf("pip package restore failed: %w", err)}
		}
	}

	var newFullVersion string
	newPythonPath := filepath.Join(venvPath, "bin", "python")
	if out, err := exec.Command(newPythonPath, "--version").Output(); err == nil {
		parts := strings.Fields(strings.TrimSpace(string(out)))
		if len(parts) >= 2 {
			newFullVersion = parts[1]
		}
	}

	if oldFullVersion != "" && oldFullVersion != newFullVersion {
		uninstallCmd := exec.Command(uvPath, "python", "uninstall", "--quiet", oldFullVersion)
		utils.RunCmd(uninstallCmd)
	}

	st.SetToolVersion("python", uvTag+"/"+version)
	return Result{Tool: "python", Version: uvTag + "/" + version}
}

func (l *LanguageRuntimeInstaller) fnmAssetName(p platform.Platform) string {
	if p.OS == platform.Darwin {
		return "fnm-macos.zip"
	}
	if p.Arch == platform.ARM64 {
		return "fnm-arm64.zip"
	}
	return "fnm-linux.zip"
}

func (l *LanguageRuntimeInstaller) downloadFnm(p platform.Platform, gh *github.Client) (fnmPath string, tagName string, err error) {
	release, err := gh.LatestRelease("Schniz/fnm")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch fnm release: %w", err)
	}

	assetName := l.fnmAssetName(p)
	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return "", "", fmt.Errorf("no fnm asset found for %s/%s", p.OS, p.Arch)
	}

	tmpDir, err := os.MkdirTemp("", "cps-fnm-*")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "fnm.zip")
	if err := DownloadToFile(downloadURL, zipPath); err != nil {
		return "", "", fmt.Errorf("download fnm failed: %w", err)
	}

	if err := ExtractZip(zipPath, tmpDir); err != nil {
		return "", "", fmt.Errorf("extract fnm failed: %w", err)
	}

	destPath := filepath.Join(p.ShellExtDir(), "fnm")
	if err := AtomicInstallBinary(filepath.Join(tmpDir, "fnm"), destPath); err != nil {
		return "", "", fmt.Errorf("install fnm binary failed: %w", err)
	}

	return destPath, release.TagName, nil
}

func (l *LanguageRuntimeInstaller) installFnmOnly(p platform.Platform, gh *github.Client, st *state.State) Result {
	_, tagName, err := l.downloadFnm(p, gh)
	if err != nil {
		return Result{Tool: "fnm", Err: err}
	}
	st.SetToolVersion("fnm", tagName)
	return Result{Tool: "fnm", Version: tagName}
}

func (l *LanguageRuntimeInstaller) installNode(p platform.Platform, gh *github.Client, st *state.State) Result {
	fnmPath, fnmTag, err := l.downloadFnm(p, gh)
	if err != nil {
		return Result{Tool: "node", Err: err}
	}

	fnmDir := filepath.Join(p.ShellDir(), "fnm")
	if err := os.MkdirAll(fnmDir, 0755); err != nil {
		return Result{Tool: "node", Err: err}
	}

	env := append(os.Environ(), "FNM_DIR="+fnmDir)

	var oldNodeVersion string
	oldVersionCmd := exec.Command(fnmPath, "exec", "--using=lts-latest", "--", "node", "--version")
	oldVersionCmd.Env = env
	if out, err := oldVersionCmd.Output(); err == nil {
		oldNodeVersion = strings.TrimSpace(string(out))
	}

	var frozen []byte
	listCmd := exec.Command(fnmPath, "exec", "--using=lts-latest", "--", "npm", "list", "-g", "--depth=0", "--json")
	listCmd.Env = env
	frozen, _ = listCmd.Output()

	installCmd := exec.Command(fnmPath, "install", "--lts")
	installCmd.Env = env
	if err := utils.RunCmd(installCmd); err != nil {
		return Result{Tool: "node", Err: fmt.Errorf("fnm install --lts failed: %w", err)}
	}

	defaultCmd := exec.Command(fnmPath, "default", "lts-latest")
	defaultCmd.Env = env
	if err := utils.RunCmd(defaultCmd); err != nil {
		return Result{Tool: "node", Err: fmt.Errorf("fnm default failed: %w", err)}
	}

	if len(frozen) > 0 {
		var pkgList struct {
			Dependencies map[string]struct {
				Version string `json:"version"`
			} `json:"dependencies"`
		}
		if json.Unmarshal(frozen, &pkgList) == nil {
			var pkgs []string
			for name := range pkgList.Dependencies {
				if name != "npm" && name != "corepack" {
					pkgs = append(pkgs, name)
				}
			}
			if len(pkgs) > 0 {
				args := append([]string{"exec", "--using=lts-latest", "--", "npm", "install", "-g"}, pkgs...)
				reinstallCmd := exec.Command(fnmPath, args...)
				reinstallCmd.Env = env
				if err := utils.RunCmd(reinstallCmd); err != nil {
					return Result{Tool: "node", Err: fmt.Errorf("npm global package restore failed: %w", err)}
				}
			}
		}
	}

	versionCmd := exec.Command(fnmPath, "exec", "--using=lts-latest", "--", "node", "--version")
	versionCmd.Env = env
	out, err := versionCmd.Output()
	nodeVersion := "lts"
	if err == nil {
		nodeVersion = strings.TrimSpace(string(out))
	}

	if oldNodeVersion != "" && oldNodeVersion != nodeVersion {
		uninstallCmd := exec.Command(fnmPath, "uninstall", oldNodeVersion)
		uninstallCmd.Env = env
		utils.RunCmd(uninstallCmd)
	}

	st.SetToolVersion("node", fnmTag+"/"+nodeVersion)
	return Result{Tool: "node", Version: fnmTag + "/" + nodeVersion}
}
