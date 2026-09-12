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
	if err != nil {
		utils.PrintFatal("failed to check latest version", err)
	}
	utils.ClearLines(1)

	if appVersion == release.TagName {
		utils.PrintSuccess(fmt.Sprintf("already at latest version %s", appVersion))
		return
	}

	assetName := fmt.Sprintf("cps-%s-%s", p.OS.String(), p.Arch.String())
	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
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

	utils.PrintRunning(fmt.Sprintf("downloading %s", release.TagName))
	tmpBinary := destPath + ".new"
	if err := installer.DownloadToFile(downloadURL, tmpBinary); err != nil {
		utils.PrintFatal("download failed", err)
	}
	if err := os.Chmod(tmpBinary, 0755); err != nil {
		os.Remove(tmpBinary)
		utils.PrintFatal("chmod failed", err)
	}
	if err := os.Rename(tmpBinary, destPath); err != nil {
		os.Remove(tmpBinary)
		utils.PrintFatal(fmt.Sprintf("failed to install binary at %s", destPath), err)
	}
	utils.ClearLines(1)

	utils.PrintSuccess(fmt.Sprintf("updated cps: %s → %s", appVersion, release.TagName))
}

func runPhase(phaseName string, tools []registry.Tool, p platform.Platform, gh *github.Client, st *state.State) bool {
	if len(tools) == 0 {
		return false
	}
	utils.PrintRunning("(Running) " + phaseName)

	var lineCount int
	var errors []jobResult

	for _, t := range tools {
		inst := installer.Dispatch(t.Kind)
		if inst == nil {
			kindErr := fmt.Errorf("no installer for kind: %s", t.Kind)
			utils.PrintIndentedError(t.Name, kindErr)
			errors = append(errors, jobResult{name: t.Name, err: kindErr})
			lineCount++
			continue
		}
		result := inst.Install(&t, p, gh, st)
		if result.Err != nil {
			utils.PrintIndentedError(t.Name, result.Err)
			errors = append(errors, jobResult{name: t.Name, err: result.Err})
		} else if result.Skipped {
			utils.PrintIndentedSuccess(fmt.Sprintf("%s: already at %s", t.Name, result.Version))
		} else if result.WasUpdated {
			utils.PrintIndentedSuccess(fmt.Sprintf("%s: updated to %s", t.Name, result.Version))
		} else {
			utils.PrintIndentedSuccess(fmt.Sprintf("%s: installed %s", t.Name, result.Version))
		}
		lineCount++
	}

	utils.ClearLines(lineCount + 1)
	if len(errors) > 0 {
		utils.PrintError(phaseName+": partially completed with errors", nil)
		for _, e := range errors {
			utils.PrintIndentedError(e.name, e.err)
		}
	} else {
		utils.PrintInfo(phaseName)
	}

	return len(errors) > 0
}

func generateShellEnv(p platform.Platform, errors *[]jobResult, lineCount *int) {
	envDir := filepath.Join(p.ShellDir(), "env")
	if err := os.MkdirAll(envDir, 0755); err != nil {
		*errors = append(*errors, jobResult{name: "shell-env", err: err})
		return
	}

	brewBin, err := exec.LookPath("brew")
	if err != nil {
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

func runPostInstall(phaseName string, p platform.Platform, withShellEnv bool) {
	utils.PrintRunning("(Running) " + phaseName)
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
		return
	}
	utils.PrintInfo(phaseName)
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
