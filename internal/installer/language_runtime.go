package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type LanguageRuntimeInstaller struct{}

func (l *LanguageRuntimeInstaller) Check(tool *registry.Tool, _ platform.Platform, _ *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)
	return current, "check-manually", nil
}

func (l *LanguageRuntimeInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	switch tool.Name {
	case "neovim":
		return l.installNeovim(p, st)
	case "go-sdk":
		return l.installGo(p, st)
	case "python":
		return l.installPython(p, st)
	default:
		return Result{Tool: tool.Name, Err: fmt.Errorf("unknown runtime: %s", tool.Name)}
	}
}

func (l *LanguageRuntimeInstaller) installNeovim(p platform.Platform, st *state.State) Result {
	utils.PrintInfo("installing Neovim")

	if p.OS == platform.Darwin {
		cmd := exec.Command("brew", "install", "neovim")
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: "neovim", Err: fmt.Errorf("brew install neovim failed: %w", err)}
		}
		st.SetToolVersion("neovim", "brew-managed")
		return Result{Tool: "neovim", Version: "brew-managed"}
	}

	// Linux: download pre-built tar.gz from GitHub stable release
	var archStr string
	switch p.Arch {
	case platform.AMD64:
		archStr = "x86_64"
	case platform.ARM64:
		archStr = "arm64"
	}

	url := fmt.Sprintf("https://github.com/neovim/neovim/releases/download/stable/nvim-linux-%s.tar.gz", archStr)

	tmpDir, err := os.MkdirTemp("", "cps-neovim-*")
	if err != nil {
		return Result{Tool: "neovim", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "nvim.tar.gz")
	if err := downloadToFile(url, tarPath); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("download failed: %w", err)}
	}

	// Remove existing installation and extract
	nvimDir := filepath.Join(p.ShellDir(), "nvim")
	os.RemoveAll(nvimDir)

	if err := ExtractTarGz(tarPath, tmpDir); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("extract failed: %w", err)}
	}

	// Move extracted directory to ~/shell/nvim
	extractedDir := filepath.Join(tmpDir, fmt.Sprintf("nvim-linux-%s", archStr))
	if err := os.Rename(extractedDir, nvimDir); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("move to %s failed: %w", nvimDir, err)}
	}

	// Symlink ~/shell/executables/nvim → ~/shell/nvim/bin/nvim
	symlinkPath := filepath.Join(p.ShellExecDir(), "nvim")
	os.Remove(symlinkPath)
	if err := os.Symlink(filepath.Join(nvimDir, "bin", "nvim"), symlinkPath); err != nil {
		return Result{Tool: "neovim", Err: fmt.Errorf("symlink failed: %w", err)}
	}

	st.SetToolVersion("neovim", "stable")
	return Result{Tool: "neovim", Version: "stable"}
}

func (l *LanguageRuntimeInstaller) installGo(p platform.Platform, st *state.State) Result {
	utils.PrintInfo("installing Go SDK")

	if p.OS == platform.Darwin {
		cmd := exec.Command("brew", "install", "go")
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: "go-sdk", Err: fmt.Errorf("brew install go failed: %w", err)}
		}
		st.SetToolVersion("go-sdk", "brew-managed")
		return Result{Tool: "go-sdk", Version: "brew-managed"}
	}

	// Linux: download from go.dev
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

	resp, err := http.Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}
	defer resp.Body.Close()

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
			if f.OS == "linux" && f.Arch == p.Arch.String() && f.Kind == "archive" {
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
		return Result{Tool: "go-sdk", Err: fmt.Errorf("no Go download found for linux/%s", p.Arch)}
	}

	tmpDir, err := os.MkdirTemp("", "cps-go-*")
	if err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "go.tar.gz")
	if err := downloadToFile(downloadURL, tarPath); err != nil {
		return Result{Tool: "go-sdk", Err: err}
	}

	// Verify tar file exists before extracting (overwrites existing installation)
	if _, err := os.Stat(tarPath); err != nil {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("downloaded tar not found: %w", err)}
	}

	rmCmd := exec.Command("sudo", "rm", "-rf", "/usr/local/go")
	utils.RunCmd(rmCmd) // ignore error — may not exist

	extractCmd := exec.Command("sudo", "tar", "-C", "/usr/local", "-xzf", tarPath)
	if err := utils.RunCmd(extractCmd); err != nil {
		return Result{Tool: "go-sdk", Err: fmt.Errorf("extract Go SDK failed: %w", err)}
	}

	st.SetToolVersion("go-sdk", version)
	return Result{Tool: "go-sdk", Version: version}
}

func (l *LanguageRuntimeInstaller) installPython(p platform.Platform, st *state.State) Result {
	utils.PrintInfo("installing Python 3.14 via uv")

	uvPath := filepath.Join(p.ShellExecDir(), "uv")

	// Install Python 3.14
	cmd := exec.Command(uvPath, "python", "install", "3.14")
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: "python", Err: fmt.Errorf("uv python install failed: %w", err)}
	}

	// Create py-default venv
	venvPath := filepath.Join(p.HomeDir, "shell", "py-default")
	venvCmd := exec.Command(uvPath, "venv", "--python", "3.14", "--clear", venvPath)
	if err := utils.RunCmd(venvCmd); err != nil {
		return Result{Tool: "python", Err: fmt.Errorf("uv venv create failed: %w", err)}
	}

	st.SetToolVersion("python", "3.14")
	return Result{Tool: "python", Version: "3.14"}
}
