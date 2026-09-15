package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func platformAndState() (platform.Platform, *state.State) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}
	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}
	return p, st
}

func SelfUpdate(appVersion string) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	gh := github.NewClient("")

	utils.PrintRunning("checking latest version")
	release, err := gh.LatestRelease("Tanq16/cli-productivity-suite")
	utils.ClearLines(1)
	if err != nil {
		utils.PrintFatal("failed to check latest version", err)
	}

	if appVersion == release.TagName {
		utils.PrintSuccess(fmt.Sprintf("already at latest version %s", appVersion))
		return
	}

	assetName := fmt.Sprintf("cps-%s-%s", p.OS.String(), p.Arch.String())
	var downloadURL string
	var assetSize int64
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			assetSize = a.Size
			break
		}
	}
	if downloadURL == "" {
		utils.PrintFatal(fmt.Sprintf("no release asset found for %s", assetName), nil)
	}

	destPath, err := os.Executable()
	if err != nil {
		utils.PrintFatal("failed to locate current cps binary", err)
	}
	if resolved, err := filepath.EvalSymlinks(destPath); err == nil {
		destPath = resolved
	}

	tmpBinary := destPath + ".new"
	m := utils.NewMeter("downloading", release.TagName, assetSize, utils.UnitBytes)
	if err := installer.DownloadToFile(downloadURL, tmpBinary, m); err != nil {
		m.Fail(err)
		os.Exit(1)
	}
	m.Done()

	if err := os.Chmod(tmpBinary, 0755); err != nil {
		os.Remove(tmpBinary)
		utils.PrintFatal("chmod failed", err)
	}
	if err := os.Rename(tmpBinary, destPath); err != nil {
		os.Remove(tmpBinary)
		utils.PrintFatal(fmt.Sprintf("failed to install binary at %s", destPath), err)
	}

	utils.PrintSuccess(fmt.Sprintf("updated cps: %s → %s", appVersion, release.TagName))
}

func runPhase(phaseName string, tools []registry.Tool, p platform.Platform, gh *github.Client, st *state.State) bool {
	if len(tools) == 0 {
		return false
	}

	m := utils.NewMeter("", phaseName, int64(len(tools)), packagesUnit)
	var failed int
	for _, t := range tools {
		m.Item(t.Name)
		inst := installer.Dispatch(t.Kind)
		if inst == nil {
			m.ItemFailed(t.Name, fmt.Errorf("no installer for kind: %s", t.Kind))
			failed++
			continue
		}
		if result := inst.Install(&t, p, gh, st); result.Err != nil {
			m.ItemFailed(t.Name, result.Err)
			failed++
			continue
		}
		m.Add(1)
	}
	m.Done()

	return failed > 0
}

func generateShellEnv(p platform.Platform, errors *[]jobResult, lineCount *int) {
	envDir := filepath.Join(p.ShellDir(), "env")
	if err := os.MkdirAll(envDir, 0755); err != nil {
		*errors = append(*errors, jobResult{name: "shell-env", err: err})
		return
	}

	brewBin := platform.BrewPath()
	if brewBin == "" {
		return
	}

	utils.PrintIndentedRunning("shell-env: brew")
	*lineCount++
	cmd := exec.Command(brewBin, "shellenv")
	// brew shellenv returns empty if its own bin is at the front of PATH, so $HOME is prepended to hide brew's own prefix from the subprocess.
	env := os.Environ()
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = "PATH=" + p.HomeDir + ":" + kv[len("PATH="):]
			break
		}
	}
	cmd.Env = env
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err == nil {
		err = os.WriteFile(filepath.Join(envDir, "brew.zsh"), out, 0644)
	} else if detail := strings.TrimSpace(stderr.String()); detail != "" {
		err = fmt.Errorf("%s: %w", detail, err)
	}

	utils.ClearPreviousLine()
	if err != nil {
		utils.PrintIndentedError("shell-env-brew", err)
		*errors = append(*errors, jobResult{name: "shell-env-brew", err: err})
		return
	}
	utils.PrintIndentedSuccess("shell-env: brew")
}

func generateCompletions(p platform.Platform, errors *[]jobResult, lineCount *int) {
	compDir := filepath.Join(p.ShellDir(), "completions")
	if err := os.MkdirAll(compDir, 0755); err != nil {
		*errors = append(*errors, jobResult{name: "completions", err: err})
		return
	}

	type compDef struct {
		name    string
		binary  string
		dir     string
		args    []string
		outFile string
	}

	defs := []compDef{
		{"fzf", "fzf", p.ShellExtDir(), []string{"--zsh"}, "fzf.zsh"},
		{"uv", "uv", p.ShellExtDir(), []string{"generate-shell-completion", "zsh"}, "uv.zsh"},
		{"fnm", "fnm", p.ShellExtDir(), []string{"completions", "--shell", "zsh"}, "fnm.zsh"},
		{"zoxide", "zoxide", p.ShellExtDir(), []string{"init", "zsh"}, "zoxide.zsh"},
		{"starship", "starship", p.ShellExtDir(), []string{"init", "zsh"}, "starship.zsh"},
	}

	for _, d := range defs {
		binPath := filepath.Join(d.dir, d.binary)
		if _, err := os.Stat(binPath); err != nil {
			continue
		}
		utils.PrintIndentedRunning("completions: " + d.name)
		*lineCount++
		cmd := exec.Command(binPath, d.args...)
		var stderr strings.Builder
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err == nil {
			err = os.WriteFile(filepath.Join(compDir, d.outFile), out, 0644)
		} else if detail := strings.TrimSpace(stderr.String()); detail != "" {
			err = fmt.Errorf("%s: %w", detail, err)
		}

		utils.ClearPreviousLine()
		if err != nil {
			utils.PrintIndentedError("completions-"+d.name, err)
			*errors = append(*errors, jobResult{name: "completions-" + d.name, err: err})
			continue
		}
		utils.PrintIndentedSuccess("completions: " + d.name)
	}
}

func runPostInstall(phaseName string, p platform.Platform, withShellEnv bool) bool {
	utils.PrintRunning(phaseName)
	var lineCount int
	var errors []jobResult

	if withShellEnv {
		generateShellEnv(p, &errors, &lineCount)
	}
	generateCompletions(p, &errors, &lineCount)

	utils.ClearLines(lineCount + 1)
	if len(errors) > 0 {
		utils.PrintError(phaseName+": partially completed with errors", nil)
		for _, e := range errors {
			utils.PrintIndentedError(e.name, e.err)
		}
		return true
	}
	utils.PrintInfo(phaseName)
	return false
}

func filterPlatformTools(tools []registry.Tool, p platform.Platform) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if t.SupportsPlatform(p.OS.String()) {
			result = append(result, t)
		}
	}
	return result
}

func filterKind(tools []registry.Tool, kind registry.ToolKind) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if t.Kind == kind {
			result = append(result, t)
		}
	}
	return result
}
