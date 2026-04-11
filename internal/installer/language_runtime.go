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

func (l *LanguageRuntimeInstaller) Check(tool *registry.Tool, _ platform.Platform, gh *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)
	switch tool.Name {
	case "go-sdk":
		return l.checkGo(current)
	case "rust":
		return l.checkRust(current)
	case "neovim":
		return l.checkNeovim(current, gh)
	case "python":
		return l.checkPython(current)
	case "node":
		return l.checkNode(current, gh)
	default:
		return current, "check-manually", nil
	}
}

func (l *LanguageRuntimeInstaller) checkGo(current string) (string, string, error) {
	resp, err := httpGet("https://go.dev/dl/?mode=json")
	if err != nil {
		return current, "", err
	}
	defer resp.Body.Close()

	var releases []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return current, "", err
	}
	for _, r := range releases {
		if r.Stable {
			return current, r.Version, nil
		}
	}
	return current, "", fmt.Errorf("no stable Go release found")
}

func (l *LanguageRuntimeInstaller) checkRust(current string) (string, string, error) {
	resp, err := httpGet("https://static.rust-lang.org/dist/channel-rust-stable.toml")
	if err != nil {
		return current, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return current, "", err
	}
	lines := strings.Split(string(body), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "[pkg.rustc]" && i+1 < len(lines) {
			next := strings.TrimSpace(lines[i+1])
			if strings.HasPrefix(next, "version = ") {
				v := strings.Trim(strings.TrimPrefix(next, "version = "), "\"")
				parts := strings.Fields(v)
				if len(parts) >= 1 {
					return current, parts[0], nil
				}
			}
		}
	}
	return current, "", fmt.Errorf("could not parse rust stable version")
}

func (l *LanguageRuntimeInstaller) checkNeovim(current string, gh *github.Client) (string, string, error) {
	release, err := gh.LatestRelease("neovim/neovim")
	if err != nil {
		return current, "", err
	}
	return current, release.TagName, nil
}

func (l *LanguageRuntimeInstaller) checkPython(current string) (string, string, error) {
	latest, err := l.latestPythonCycle()
	if err != nil {
		return current, "", err
	}
	return current, latest, nil
}

func (l *LanguageRuntimeInstaller) Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result {
	switch tool.Name {
	case "neovim":
		return l.installNeovim(p, gh, st)
	case "go-sdk":
		return l.installGo(p, st)
	case "python":
		return l.installPython(p, st)
	case "rust":
		return l.installRust(p, st)
	case "node":
		return l.installNode(p, gh, st)
	default:
		return Result{Tool: tool.Name, Err: fmt.Errorf("unknown runtime: %s", tool.Name)}
	}
}

func (l *LanguageRuntimeInstaller) installNeovim(p platform.Platform, gh *github.Client, st *state.State) Result {
	var archStr string
	switch p.Arch {
	case platform.AMD64:
		archStr = "x86_64"
	case platform.ARM64:
		archStr = "arm64"
	}

	var osStr string
	switch p.OS {
	case platform.Darwin:
		osStr = "macos"
	default:
		osStr = "linux"
	}

	url := fmt.Sprintf("https://github.com/neovim/neovim/releases/download/stable/nvim-%s-%s.tar.gz", osStr, archStr)

	tmpDir, err := os.MkdirTemp("", "cps-neovim-*")
	if err != nil {
		return Result{Tool: "neovim", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "nvim.tar.gz")
	if err := DownloadToFile(url, tarPath); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("download failed: %w", err)}
	}

	if p.OS == platform.Darwin {
		xattrCmd := exec.Command("xattr", "-c", tarPath)
		utils.RunCmd(xattrCmd)
	}

	nvimDir := filepath.Join(p.ShellDir(), "nvim")
	os.RemoveAll(nvimDir)

	if err := ExtractTarGz(tarPath, tmpDir); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("extract failed: %w", err)}
	}

	extractedDir := filepath.Join(tmpDir, fmt.Sprintf("nvim-%s-%s", osStr, archStr))
	if err := os.Rename(extractedDir, nvimDir); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("move to %s failed: %w", nvimDir, err)}
	}

	symlinkPath := filepath.Join(p.ShellExecDir(), "nvim")
	os.Remove(symlinkPath)
	if err := os.Symlink(filepath.Join(nvimDir, "bin", "nvim"), symlinkPath); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("symlink failed: %w", err)}
	}

	version := "stable"
	if release, err := gh.LatestRelease("neovim/neovim"); err == nil {
		version = release.TagName
	}
	st.SetToolVersion("neovim", version)
	return Result{Tool: "neovim", Version: version}
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

	tmpDir, err := os.MkdirTemp("", "cps-go-*")
	if err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "go.tar.gz")
	if err := DownloadToFile(downloadURL, tarPath); err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}

	goDir := filepath.Join(p.ShellDir(), "go-sdk")
	os.RemoveAll(goDir)

	if err := ExtractTarGz(tarPath, tmpDir); err != nil {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("extract Go SDK failed: %w", err)}
	}

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

func (l *LanguageRuntimeInstaller) installPython(p platform.Platform, st *state.State) Result {
	uvPath := filepath.Join(p.ShellExecDir(), "uv")

	version, err := l.latestPythonCycle()
	if err != nil {
		return Result{Tool: "python", Err: err}
	}

	cmd := exec.Command(uvPath, "python", "install", version)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: "python", Err: fmt.Errorf("uv python install failed: %w", err)}
	}

	venvPath := filepath.Join(p.HomeDir, "shell", "py-default")
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
		installCmd.Run()
	}

	st.SetToolVersion("python", version)
	return Result{Tool: "python", Version: version}
}

func (l *LanguageRuntimeInstaller) checkNode(current string, gh *github.Client) (string, string, error) {
	release, err := gh.LatestRelease("Schniz/fnm")
	if err != nil {
		return current, "", err
	}
	// State stores "fnmVersion/nodeVersion"; if fnm version changed, signal update
	if strings.HasPrefix(current, release.TagName+"/") {
		return current, current, nil
	}
	return current, release.TagName, nil
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

func (l *LanguageRuntimeInstaller) installNode(p platform.Platform, gh *github.Client, st *state.State) Result {
	release, err := gh.LatestRelease("Schniz/fnm")
	if err != nil {
		return Result{Tool: "node", Err: fmt.Errorf("failed to fetch fnm release: %w", err)}
	}

	fnmPath := filepath.Join(p.ShellExecDir(), "fnm")
	fnmDir := filepath.Join(p.ShellDir(), "fnm")

	if err := os.MkdirAll(fnmDir, 0755); err != nil {
		return Result{Tool: "node", Err: err}
	}

	// Download and install fnm binary
	assetName := l.fnmAssetName(p)
	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return Result{Tool: "node", Err: fmt.Errorf("no fnm asset found for %s/%s", p.OS, p.Arch)}
	}

	tmpDir, err := os.MkdirTemp("", "cps-fnm-*")
	if err != nil {
		return Result{Tool: "node", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "fnm.zip")
	if err := DownloadToFile(downloadURL, zipPath); err != nil {
		return Result{Tool: "node", Err: fmt.Errorf("download fnm failed: %w", err)}
	}

	if err := ExtractZip(zipPath, tmpDir); err != nil {
		return Result{Tool: "node", Err: fmt.Errorf("extract fnm failed: %w", err)}
	}

	extractedBin := filepath.Join(tmpDir, "fnm")
	if err := AtomicInstallBinary(extractedBin, fnmPath); err != nil {
		return Result{Tool: "node", Err: fmt.Errorf("install fnm binary failed: %w", err)}
	}

	// Capture existing global npm packages before installing new Node
	var frozen []byte
	npmPath := filepath.Join(fnmDir, "aliases", "lts-latest", "installation", "bin", "npm")
	if _, err := os.Stat(npmPath); err == nil {
		env := append(os.Environ(), "FNM_DIR="+fnmDir)
		listCmd := exec.Command(npmPath, "list", "-g", "--depth=0", "--json")
		listCmd.Env = env
		frozen, _ = listCmd.Output()
	}

	// Install Node LTS via fnm
	env := append(os.Environ(), "FNM_DIR="+fnmDir)

	installCmd := exec.Command(fnmPath, "install", "--lts")
	installCmd.Env = env
	if err := utils.RunCmd(installCmd); err != nil {
		return Result{Tool: "node", Err: fmt.Errorf("fnm install --lts failed: %w", err)}
	}

	defaultCmd := exec.Command(fnmPath, "default", "lts-latest")
	defaultCmd.Env = env
	utils.RunCmd(defaultCmd)

	// Reinstall global npm packages if any were captured
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
				newNpmPath := filepath.Join(fnmDir, "aliases", "lts-latest", "installation", "bin", "npm")
				if _, err := os.Stat(newNpmPath); err == nil {
					reinstallCmd := exec.Command(newNpmPath, append([]string{"install", "-g"}, pkgs...)...)
					reinstallCmd.Env = env
					reinstallCmd.Run()
				}
			}
		}
	}

	// Get the installed Node version via fnm exec (doesn't need fnm env)
	versionCmd := exec.Command(fnmPath, "exec", "--using=lts-latest", "--", "node", "--version")
	versionCmd.Env = env
	out, err := versionCmd.Output()
	nodeVersion := "lts"
	if err == nil {
		nodeVersion = strings.TrimSpace(string(out))
	}

	st.SetToolVersion("node", release.TagName+"/"+nodeVersion)
	return Result{Tool: "node", Version: release.TagName + "/" + nodeVersion}
}
