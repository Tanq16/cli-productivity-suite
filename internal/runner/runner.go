package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func Init(ghToken string) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}

	gh := github.NewClient(ghToken)
	reg := registry.New()

	// Phase 1: Prerequisites
	utils.PrintRunning("(Running) Phase 1: Checking prerequisites")
	omzDir := filepath.Join(p.HomeDir, ".oh-my-zsh")
	if _, err := os.Stat(omzDir); os.IsNotExist(err) {
		msg := fmt.Sprintf("Oh My Zsh not found at %s\nInstall it first: sh -c \"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\"", omzDir)
		utils.PrintFatal(msg, nil)
	}
	if _, err := exec.LookPath("git"); err != nil {
		utils.PrintFatal("git not found in PATH", err)
	}
	if _, err := exec.LookPath("brew"); err != nil {
		msg := "Homebrew not found in PATH\nInstall it first: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
		utils.PrintFatal(msg, err)
	}
	for _, dir := range []string{
		p.ShellDir(),
		p.ShellExecDir(),
		filepath.Join(p.ShellDir(), "rc"),
		filepath.Join(p.ShellDir(), "rc", "custom"),
		filepath.Join(p.ShellDir(), "custom"),
		filepath.Join(p.ConfigDir(), "extensions"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			utils.PrintFatal(fmt.Sprintf("failed to create %s", dir), err)
		}
	}
	utils.ClearLines(1)
	utils.PrintInfo("Phase 1: Checking prerequisites")

	var hadErrors bool

	sysPkgs := filterBaseTools(filterPlatformTools(reg.ByKind(registry.SystemPackage), p))

	// Phase 2: System packages
	if runPhase("Phase 2: System packages", sysPkgs, p, gh, st) {
		hadErrors = true
	}
	st.Save()

	// Phase 3: Core GitHub releases
	coreGH := filterGitHubCore(reg.ByKind(registry.GitHubRelease), false)
	coreGH = filterPlatformTools(coreGH, p)
	if runPhase("Phase 3: Core GitHub releases", coreGH, p, gh, st) {
		hadErrors = true
	}
	st.Save()

	// Phase 4: Own core tools
	ownCore := filterOwnCore(reg.ByKind(registry.GitHubRelease))
	if runPhase("Phase 4: Own core tools", ownCore, p, gh, st) {
		hadErrors = true
	}
	st.Save()

	// Phase 5: Shell plugins (base only, exclude extension tools)
	disableOMZSlowPaste(p)
	if runPhase("Phase 5: Shell plugins", filterBaseTools(reg.ByKind(registry.ShellPlugin)), p, gh, st) {
		hadErrors = true
	}
	st.Save()

	// Phase 6: Config files
	if runPhase("Phase 6: Config files", filterPlatformTools(reg.ByKind(registry.ConfigFile), p), p, gh, st) {
		hadErrors = true
	}
	st.Save()

	// Phase 7: Post-install tasks
	runPostInstall(p)

	st.LastInit = time.Now()
	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	if hadErrors {
		utils.PrintWarn("init finished with errors", nil)
	} else {
		utils.PrintSuccess("init complete!")
	}
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
			utils.PrintIndentedError(t.Name+": no installer for kind: "+t.Kind.String(), nil)
			errors = append(errors, jobResult{name: t.Name, err: fmt.Errorf("no installer for kind: %s", t.Kind)})
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

	utils.ClearLines(lineCount + 1) // tool lines + running header
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

func runPostInstall(p platform.Platform) {
	utils.PrintRunning("(Running) Phase 7: Post-install tasks")
	var lineCount int
	var errors []jobResult

	tpmInstall := filepath.Join(p.HomeDir, ".tmux", "plugins", "tpm", "bin", "install_plugins")
	if _, err := os.Stat(tpmInstall); err == nil {
		utils.PrintIndentedRunning("tpm-install: running")
		lineCount++
		tpmCmd := exec.Command("bash", tpmInstall)
		tpmCmd.Env = append(os.Environ(), fmt.Sprintf("TMUX_PLUGIN_MANAGER_PATH=%s", filepath.Join(p.HomeDir, ".tmux", "plugins")))
		if err := utils.RunCmd(tpmCmd); err != nil {
			errors = append(errors, jobResult{name: "tpm-install", err: err})
		}
	}

	if nvimBin, err := exec.LookPath("nvim"); err == nil {
		utils.PrintIndentedRunning("nvchad-setup: running")
		lineCount++
		nvimCmd := exec.Command(nvimBin, "--headless", "+MasonInstallAll", "+Lazy sync", "+qa")
		if err := utils.RunCmd(nvimCmd); err != nil {
			errors = append(errors, jobResult{name: "nvchad-setup", err: err})
		}
	}

	generateCompletions(p, &errors, &lineCount)

	utils.ClearLines(lineCount + 1) // sub-lines + running header
	if len(errors) > 0 {
		utils.PrintError("Phase 7: partially completed with errors", nil)
		for _, e := range errors {
			utils.PrintIndentedError(e.name, e.err)
		}
	} else {
		utils.PrintInfo("Phase 7: Post-install tasks")
	}
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
		{"fzf", "fzf", p.ShellExecDir(), []string{"--zsh"}, "fzf.zsh"},
		{"uv", "uv", p.ShellExtDir(), []string{"generate-shell-completion", "zsh"}, "uv.zsh"},
		{"fnm", "fnm", p.ShellExtDir(), []string{"completions", "--shell", "zsh"}, "fnm.zsh"},
		{"zoxide", "zoxide", p.ShellExecDir(), []string{"init", "zsh"}, "zoxide.zsh"},
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
		if err != nil {
			if detail := strings.TrimSpace(stderr.String()); detail != "" {
				err = fmt.Errorf("%s: %w", detail, err)
			}
			*errors = append(*errors, jobResult{name: "completions-" + d.name, err: err})
			continue
		}
		if err := os.WriteFile(filepath.Join(compDir, d.outFile), out, 0644); err != nil {
			*errors = append(*errors, jobResult{name: "completions-" + d.name, err: err})
		}
	}
}

// --- Helpers ---

func toolForPlatform(tool registry.Tool, p platform.Platform) bool {
	if len(tool.Platforms) == 0 {
		return true
	}
	for _, plat := range tool.Platforms {
		if plat == p.OS.String() {
			return true
		}
	}
	return false
}

func filterGitHubCore(tools []registry.Tool, includeOwn bool) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if t.Category == registry.Core && !t.IsPrivate {
			if !includeOwn && isOwnTool(t.Repo) {
				continue
			}
			result = append(result, t)
		}
	}
	return result
}

func filterOwnCore(tools []registry.Tool) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if t.Category == registry.Core && isOwnTool(t.Repo) {
			result = append(result, t)
		}
	}
	return result
}

func filterPlatformTools(tools []registry.Tool, p platform.Platform) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if toolForPlatform(t, p) {
			result = append(result, t)
		}
	}
	return result
}

func filterBaseTools(tools []registry.Tool) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if !t.Extension {
			result = append(result, t)
		}
	}
	return result
}

func isOwnTool(repo string) bool {
	return strings.HasPrefix(repo, "Tanq16/")
}

func disableOMZSlowPaste(p platform.Platform) {
	miscPath := filepath.Join(p.HomeDir, ".oh-my-zsh", "lib", "misc.zsh")
	data, err := os.ReadFile(miscPath)
	if err != nil {
		return
	}
	content := string(data)
	lines := strings.Split(content, "\n")
	targets := []string{
		"autoload -Uz bracketed-paste-magic",
		"zle -N bracketed-paste bracketed-paste-magic",
		"autoload -Uz url-quote-magic",
		"zle -N self-insert url-quote-magic",
	}
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		for _, target := range targets {
			if strings.Contains(line, target) {
				lines[i] = "#" + line
				break
			}
		}
	}
	if err := os.WriteFile(miscPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		utils.PrintWarn("failed to disable OMZ slow-paste", err)
	}
}
