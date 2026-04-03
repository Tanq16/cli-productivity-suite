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
	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "version = ") {
			v := strings.Trim(strings.TrimPrefix(strings.TrimSpace(line), "version = "), "\"")
			parts := strings.Fields(v)
			if len(parts) >= 1 {
				return current, parts[0], nil
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
	resp, err := httpGet("https://endoflife.date/api/python.json")
	if err != nil {
		return current, "", err
	}
	defer resp.Body.Close()

	var releases []struct {
		Cycle  string `json:"cycle"`
		Latest string `json:"latest"`
		EOL    any    `json:"eol"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return current, "", err
	}
	now := time.Now().Format("2006-01-02")
	for _, r := range releases {
		switch eol := r.EOL.(type) {
		case bool:
			if !eol {
				return current, r.Cycle, nil
			}
		case string:
			if eol > now {
				return current, r.Cycle, nil
			}
		}
	}
	return current, "", fmt.Errorf("no active Python release found")
}

func (l *LanguageRuntimeInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	switch tool.Name {
	case "neovim":
		return l.installNeovim(p, st)
	case "go-sdk":
		return l.installGo(p, st)
	case "python":
		return l.installPython(p, st)
	case "rust":
		return l.installRust(p, st)
	default:
		return Result{Tool: tool.Name, Err: fmt.Errorf("unknown runtime: %s", tool.Name)}
	}
}

func (l *LanguageRuntimeInstaller) installNeovim(p platform.Platform, st *state.State) Result {
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

	st.SetToolVersion("neovim", "stable")
	return Result{Tool: "neovim", Version: "stable"}
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

	rmCmd := exec.Command("sudo", "rm", "-rf", "/usr/local/go")
	utils.RunCmd(rmCmd)

	extractCmd := exec.Command("sudo", "tar", "-C", "/usr/local", "-xzf", tarPath)
	if err := utils.RunCmd(extractCmd); err != nil {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("extract Go SDK failed: %w", err)}
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

func (l *LanguageRuntimeInstaller) installPython(p platform.Platform, st *state.State) Result {
	uvPath := filepath.Join(p.ShellExecDir(), "uv")

	cmd := exec.Command(uvPath, "python", "install", "3.14")
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: "python", Err: fmt.Errorf("uv python install failed: %w", err)}
	}

	venvPath := filepath.Join(p.HomeDir, "shell", "py-default")
	venvCmd := exec.Command(uvPath, "venv", "--python", "3.14", "--clear", venvPath)
	if err := utils.RunCmd(venvCmd); err != nil {
		return Result{Tool: "python", Err: fmt.Errorf("uv venv create failed: %w", err)}
	}

	st.SetToolVersion("python", "3.14")
	return Result{Tool: "python", Version: "3.14"}
}
