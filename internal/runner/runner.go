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

type CheckConfig struct {
	Public  bool
	Private bool
	System  bool
	All     bool
}

type UpdateConfig struct {
	Public      bool
	Private     bool
	System      bool
	All         bool
	IncludeConf bool
}

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

	utils.PrintInfo("Phase 1: Checking prerequisites")
	omzDir := filepath.Join(p.HomeDir, ".oh-my-zsh")
	if _, err := os.Stat(omzDir); os.IsNotExist(err) {
		msg := fmt.Sprintf("Oh My Zsh not found at %s\nInstall it first: sh -c \"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\"", omzDir)
		utils.PrintFatal(msg, nil)
	}
	if _, err := exec.LookPath("git"); err != nil {
		utils.PrintFatal("git not found in PATH", nil)
	}
	if err := os.MkdirAll(p.ShellDir(), 0755); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to create %s", p.ShellDir()), err)
	}
	if err := os.MkdirAll(p.ShellExecDir(), 0755); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to create %s", p.ShellExecDir()), err)
	}
	utils.PrintSuccess("prerequisites OK")

	if PhaseNeedsSudo(p, registry.SystemPackage, registry.CloudCLI, registry.LanguageRuntime) {
		utils.PrintInfo("some phases require sudo — authenticating")
		if err := EnsureSudo(); err != nil {
			utils.PrintFatal("sudo authentication failed", err)
		}
		utils.ClearPreviousLine()
	}

	runPhase("Phase 2: System packages", filterPlatformTools(reg.ByKind(registry.SystemPackage), p), p, gh, st)
	st.Save()

	runPhase("Phase 3: Cloud CLIs", reg.ByKind(registry.CloudCLI), p, gh, st)
	st.Save()

	goSDK := filterByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	runPhase("Phase 4: Go SDK", goSDK, p, gh, st)
	st.Save()

	publicGH := filterGitHubPublic(reg.ByKind(registry.GitHubRelease), false)
	publicGH = filterPlatformTools(publicGH, p)
	runPhase("Phase 5: Public GitHub releases", publicGH, p, gh, st)
	st.Save()

	runPhase("Phase 6: Direct downloads", reg.ByKind(registry.DirectDownload), p, gh, st)
	st.Save()

	ownPublic := filterOwnPublic(reg.ByKind(registry.GitHubRelease))
	runPhase("Phase 7: Own public tools", ownPublic, p, gh, st)
	st.Save()

	if ghToken != "" {
		privateTools := reg.ByCategory(registry.Private)
		runPhase("Phase 8: Private tools", privateTools, p, gh, st)
	} else {
		utils.PrintWarn("Phase 8: Skipping private tools (no --gh-token)", nil)
	}
	st.Save()

	nonSudoRuntimes := excludeByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	runPhase("Phase 9: Language runtimes", nonSudoRuntimes, p, gh, st)
	st.Save()

	disableOMZSlowPaste(p)
	runPhase("Phase 10: Shell plugins", reg.ByKind(registry.ShellPlugin), p, gh, st)
	st.Save()

	runPhase("Phase 11: Config files", filterPlatformTools(reg.ByKind(registry.ConfigFile), p), p, gh, st)
	st.Save()

	runPostInstall(p)

	st.LastInit = time.Now()
	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	utils.PrintSuccess("init complete!")
}

func Check(ghToken string, appVersion string, cfg CheckConfig) {
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

	if !cfg.Public && !cfg.Private && !cfg.System {
		cfg.All = true
	}

	var tools []registry.Tool
	if cfg.All {
		tools = reg.ForPlatform(p.OS.String())
	} else {
		if cfg.Public {
			tools = append(tools, reg.ByCategory(registry.Public)...)
		}
		if cfg.Private {
			tools = append(tools, reg.ByCategory(registry.Private)...)
		}
		if cfg.System {
			tools = append(tools, reg.ByCategory(registry.System)...)
		}
	}

	if len(tools) == 0 {
		utils.PrintWarn("no tools to check", nil)
		return
	}

	results := runCheckPhase("Checking tools", tools, p, gh, st)

	cpsRow := []string{"cps", appVersion, "unknown", "up-to-date"}
	if rel, err := gh.LatestRelease("Tanq16/cli-productivity-suite"); err == nil {
		cpsRow[2] = rel.TagName
		if appVersion == "dev-build" {
			cpsRow[3] = "dev build"
		} else if appVersion != rel.TagName {
			cpsRow[3] = "update available"
		}
	}

	headers := []string{"Tool", "Current", "Latest", "Status"}
	rows := [][]string{cpsRow}
	for _, r := range results {
		rows = append(rows, []string{r.name, r.current, r.latest, r.status})
	}

	if len(rows) == 0 {
		utils.PrintWarn("no installed tools found in state", nil)
		return
	}

	utils.PrintTable(headers, rows)
}

func Update(ghToken string, cfg UpdateConfig) {
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

	if !cfg.Public && !cfg.Private && !cfg.System {
		cfg.All = true
	}

	var tools []registry.Tool
	if cfg.All {
		tools = reg.ForPlatform(p.OS.String())
	} else {
		if cfg.Public {
			tools = append(tools, reg.ByCategory(registry.Public)...)
		}
		if cfg.Private {
			tools = append(tools, reg.ByCategory(registry.Private)...)
		}
		if cfg.System {
			tools = append(tools, reg.ByCategory(registry.System)...)
		}
	}

	var updatable []registry.Tool
	for _, t := range tools {
		if st.ToolVersion(t.Name) == "" {
			continue
		}
		if t.Kind == registry.ConfigFile && !cfg.IncludeConf {
			continue
		}
		updatable = append(updatable, t)
	}

	if len(updatable) == 0 {
		utils.PrintWarn("no updatable tools found", nil)
		return
	}

	needsSudo := false
	for _, t := range updatable {
		if ToolNeedsSudo(t, p) {
			needsSudo = true
			break
		}
	}
	if needsSudo {
		utils.PrintInfo("some tools require sudo — authenticating")
		if err := EnsureSudo(); err != nil {
			utils.PrintFatal("sudo authentication failed", err)
		}
		utils.ClearPreviousLine()
	}

	runPhase("Updating tools", updatable, p, gh, st)

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	utils.PrintSuccess("update complete!")
}

func Install(toolName string, ghToken string) {
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

	tool := reg.Get(toolName)
	if tool == nil {
		utils.PrintFatal(fmt.Sprintf("unknown tool: %s", toolName), nil)
	}

	if tool.IsPrivate && ghToken == "" {
		utils.PrintFatal(fmt.Sprintf("tool %s is private — provide --gh-token or set CPS_GITHUB_PAT", toolName), nil)
	}

	if ToolNeedsSudo(*tool, p) {
		utils.PrintInfo("this tool requires sudo — authenticating")
		if err := EnsureSudo(); err != nil {
			utils.PrintFatal("sudo authentication failed", err)
		}
		utils.ClearPreviousLine()
	}

	hasErrors := runPhase(fmt.Sprintf("Installing %s", toolName), []registry.Tool{*tool}, p, gh, st)

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	if hasErrors {
		utils.PrintFatal(fmt.Sprintf("install failed for %s", toolName), nil)
	}

	utils.PrintSuccess(fmt.Sprintf("%s installed", toolName))
}

func Clean() {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	dirs := []string{
		filepath.Join(p.HomeDir, "shell"),
		filepath.Join(p.HomeDir, ".tmux"),
		filepath.Join(p.HomeDir, ".config", "nvim"),
		filepath.Join(p.HomeDir, ".nvm"),
		filepath.Join(p.HomeDir, "nuclei-templates"),
		filepath.Join(p.HomeDir, "google-cloud-sdk"),
		filepath.Join(p.HomeDir, ".config", "cps"),
	}

	utils.PrintWarn("this will remove the following directories:", nil)
	for _, d := range dirs {
		utils.PrintGeneric("  " + d)
	}
	answer, err := utils.PromptInput("\nare you sure? (yes/no):", "yes/no")
	if err != nil {
		utils.PrintFatal("failed to read input", err)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "yes" {
		utils.PrintInfo("clean aborted")
		return
	}

	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			continue
		}
		if err := os.RemoveAll(d); err != nil {
			utils.PrintError(fmt.Sprintf("failed to remove %s", d), err)
		} else {
			utils.PrintSuccess(fmt.Sprintf("removed %s", d))
		}
	}

	utils.PrintSuccess("clean complete")
}

func runPhase(phaseName string, tools []registry.Tool, p platform.Platform, gh *github.Client, st *state.State) bool {
	if len(tools) == 0 {
		return false
	}
	utils.PrintInfo(phaseName)

	var lineCount int
	var errors []jobResult

	for _, t := range tools {
		inst := installer.Dispatch(t.Kind)
		if inst == nil {
			utils.PrintError(t.Name, fmt.Errorf("no installer for kind: %s", t.Kind))
			errors = append(errors, jobResult{name: t.Name, err: fmt.Errorf("no installer for kind: %s", t.Kind)})
			lineCount++
			continue
		}
		result := inst.Install(&t, p, gh, st)
		if result.Err != nil {
			utils.PrintError(t.Name, result.Err)
			errors = append(errors, jobResult{name: t.Name, err: result.Err})
		} else if result.Skipped {
			utils.PrintSuccess(fmt.Sprintf("%s: already at %s", t.Name, result.Version))
		} else if result.WasUpdated {
			utils.PrintSuccess(fmt.Sprintf("%s: updated to %s", t.Name, result.Version))
		} else {
			utils.PrintSuccess(fmt.Sprintf("%s: installed %s", t.Name, result.Version))
		}
		lineCount++
	}

	utils.ClearLines(lineCount)
	for _, e := range errors {
		utils.PrintError(e.name, e.err)
	}

	return len(errors) > 0
}

func runCheckPhase(phaseName string, tools []registry.Tool, p platform.Platform, gh *github.Client, st *state.State) []checkResult {
	utils.PrintInfo(phaseName)

	var all []checkResult
	var lineCount int

	for _, t := range tools {
		current := st.ToolVersion(t.Name)
		if current == "" {
			continue
		}
		inst := installer.Dispatch(t.Kind)
		if inst == nil {
			continue
		}
		cur, lat, err := inst.Check(&t, p, gh, st)
		if err != nil {
			all = append(all, checkResult{name: t.Name, current: cur, latest: "error", status: err.Error()})
			utils.PrintInfo(fmt.Sprintf("%s: error", t.Name))
			lineCount++
			continue
		}
		status := "up-to-date"
		switch lat {
		case "system-managed", "check-manually", "git-managed":
			status = lat
		case "not-deployed":
			status = "not deployed"
		case "config-differs":
			status = "config differs"
		case "skipped":
			continue
		default:
			if cur != lat && lat != "" {
				status = "update available"
			}
		}
		all = append(all, checkResult{name: t.Name, current: cur, latest: lat, status: status})
		utils.PrintInfo(fmt.Sprintf("%s: %s", t.Name, status))
		lineCount++
	}

	utils.ClearLines(lineCount)

	return all
}

func runPostInstall(p platform.Platform) {
	utils.PrintInfo("Phase 12: Post-install tasks")
	var lineCount int
	var errors []jobResult

	tpmInstall := filepath.Join(p.HomeDir, ".tmux", "plugins", "tpm", "bin", "install_plugins")
	if _, err := os.Stat(tpmInstall); err == nil {
		utils.PrintSuccess("tpm-install: running")
		lineCount++
		tpmCmd := exec.Command("bash", tpmInstall)
		tpmCmd.Env = append(os.Environ(), fmt.Sprintf("TMUX_PLUGIN_MANAGER_PATH=%s", filepath.Join(p.HomeDir, ".tmux", "plugins")))
		if err := utils.RunCmd(tpmCmd); err != nil {
			errors = append(errors, jobResult{name: "tpm-install", err: err})
		}
	}

	if _, err := exec.LookPath("nvim"); err == nil {
		utils.PrintSuccess("nvchad-setup: running")
		lineCount++
		nvimCmd := exec.Command("nvim", "--headless", "+MasonInstallAll", "+Lazy sync", "+qa")
		if err := utils.RunCmd(nvimCmd); err != nil {
			errors = append(errors, jobResult{name: "nvchad-setup", err: err})
		}
	}

	utils.ClearLines(lineCount)
	for _, e := range errors {
		utils.PrintError(e.name, e.err)
	}
}

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

func filterGitHubPublic(tools []registry.Tool, includeOwn bool) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if t.Category == registry.Public && !t.IsPrivate {
			if !includeOwn && isOwnTool(t.Repo) {
				continue
			}
			result = append(result, t)
		}
	}
	return result
}

func filterOwnPublic(tools []registry.Tool) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if t.Category == registry.Public && isOwnTool(t.Repo) {
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

func filterByName(tools []registry.Tool, names ...string) []registry.Tool {
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}
	var result []registry.Tool
	for _, t := range tools {
		if nameSet[t.Name] {
			result = append(result, t)
		}
	}
	return result
}

func excludeByName(tools []registry.Tool, names ...string) []registry.Tool {
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}
	var result []registry.Tool
	for _, t := range tools {
		if !nameSet[t.Name] {
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
