package orchestrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/tanq16/cli-productivity-suite/internal/display"
	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/highway"
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

func RunInit(parentCtx context.Context, ghToken string) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}

	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt)
	defer cancel()

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
	}

	disp := display.New()

	runPhase(ctx, disp, "Phase 2: System packages", filterPlatformTools(reg.ByKind(registry.SystemPackage), p), 1, p, gh, st)
	st.Save()

	runPhase(ctx, disp, "Phase 3: Cloud CLIs", reg.ByKind(registry.CloudCLI), 1, p, gh, st)
	st.Save()

	goSDK := filterByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	runPhase(ctx, disp, "Phase 4: Go SDK", goSDK, 1, p, gh, st)
	st.Save()

	publicGH := filterGitHubPublic(reg.ByKind(registry.GitHubRelease), false)
	publicGH = filterPlatformTools(publicGH, p)
	runPhase(ctx, disp, "Phase 5: Public GitHub releases", publicGH, 2, p, gh, st)
	st.Save()

	runPhase(ctx, disp, "Phase 6: Direct downloads", reg.ByKind(registry.DirectDownload), 1, p, gh, st)
	st.Save()

	ownPublic := filterOwnPublic(reg.ByKind(registry.GitHubRelease))
	runPhase(ctx, disp, "Phase 7: Own public tools", ownPublic, 2, p, gh, st)
	st.Save()

	if ghToken != "" {
		privateTools := reg.ByCategory(registry.Private)
		runPhase(ctx, disp, "Phase 8: Private tools", privateTools, 2, p, gh, st)
	} else {
		utils.PrintWarn("Phase 8: Skipping private tools (no --gh-token)", nil)
	}
	st.Save()

	nonSudoRuntimes := excludeByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	runPhase(ctx, disp, "Phase 9: Language runtimes", nonSudoRuntimes, 1, p, gh, st)
	st.Save()

	disableOMZSlowPaste(p)
	runPhase(ctx, disp, "Phase 10: Shell plugins", reg.ByKind(registry.ShellPlugin), 1, p, gh, st)
	st.Save()

	runPhase(ctx, disp, "Phase 11: Config files", filterPlatformTools(reg.ByKind(registry.ConfigFile), p), 1, p, gh, st)
	st.Save()

	utils.PrintInfo("Phase 12: Running post-install tasks")

	tpmInstall := filepath.Join(p.HomeDir, ".tmux", "plugins", "tpm", "bin", "install_plugins")
	if _, err := os.Stat(tpmInstall); err == nil {
		utils.PrintInfo("running TPM install_plugins")
		tpmCmd := exec.Command("bash", tpmInstall)
		tpmCmd.Env = append(os.Environ(), fmt.Sprintf("TMUX_PLUGIN_MANAGER_PATH=%s", filepath.Join(p.HomeDir, ".tmux", "plugins")))
		if err := utils.RunCmd(tpmCmd); err != nil {
			utils.PrintWarn("TPM install_plugins failed", err)
		}
	}

	if _, err := exec.LookPath("nvim"); err == nil {
		utils.PrintInfo("running NvChad headless setup")
		nvimCmd := exec.Command("nvim", "--headless", "+MasonInstallAll", "+Lazy sync", "+qa")
		if err := utils.RunCmd(nvimCmd); err != nil {
			utils.PrintWarn("NvChad headless setup failed", err)
		}
	}

	st.LastInit = time.Now()
	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	utils.PrintSuccess("init complete!")
}

func RunCheck(ghToken string, appVersion string, cfg CheckConfig) {
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

	headers := []string{"Tool", "Current", "Latest", "Status"}
	var rows [][]string

	for _, tool := range tools {
		current := st.ToolVersion(tool.Name)
		if current == "" {
			continue
		}

		t := tool
		inst := installer.Dispatch(t.Kind)
		if inst == nil {
			continue
		}

		cur, lat, err := inst.Check(&t, p, gh, st)
		if err != nil {
			rows = append(rows, []string{t.Name, cur, "error", err.Error()})
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
		rows = append(rows, []string{t.Name, cur, lat, status})
	}

	cpsRow := []string{"cps", appVersion, "unknown", "up-to-date"}
	if rel, err := gh.LatestRelease("Tanq16/cli-productivity-suite"); err == nil {
		cpsRow[2] = rel.TagName
		if appVersion == "dev-build" {
			cpsRow[3] = "dev build"
		} else if appVersion != rel.TagName {
			cpsRow[3] = "update available"
		}
	}
	rows = append([][]string{cpsRow}, rows...)

	if len(rows) == 0 {
		utils.PrintWarn("no installed tools found in state", nil)
		return
	}

	utils.PrintTable(headers, rows)
}

func RunUpdate(parentCtx context.Context, ghToken string, cfg UpdateConfig) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}

	ctx, cancel := signal.NotifyContext(parentCtx, os.Interrupt)
	defer cancel()

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
	}

	utils.PrintInfo(fmt.Sprintf("updating %d tools", len(updatable)))

	workers := 2
	if needsSudo {
		workers = 1
	}

	disp := display.New()
	runPhase(ctx, disp, "Updating tools", updatable, workers, p, gh, st)

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	utils.PrintSuccess("update complete!")
}

func RunInstall(toolName string, ghToken string) {
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

	inst := installer.Dispatch(tool.Kind)
	if inst == nil {
		utils.PrintFatal(fmt.Sprintf("no installer for tool kind: %s", tool.Kind), nil)
	}

	if ToolNeedsSudo(*tool, p) {
		utils.PrintInfo("this tool requires sudo — authenticating")
		if err := EnsureSudo(); err != nil {
			utils.PrintFatal("sudo authentication failed", err)
		}
	}

	result := inst.Install(tool, p, gh, st)
	printResult(result)

	if result.Err != nil {
		utils.PrintFatal(fmt.Sprintf("install failed for %s", toolName), nil)
	}

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}
}

func RunClean() {
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

func runPhase(ctx context.Context, disp *display.Display, phaseName string, tools []registry.Tool, workers int, p platform.Platform, gh *github.Client, st *state.State) {
	if len(tools) == 0 {
		return
	}
	hw := highway.New(workers, "")
	jobs := make([]highway.Job, len(tools))
	for i, t := range tools {
		jobs[i] = NewInstallJob(t, p, gh, st)
	}
	hw.Submit(jobs...)
	errCh := make(chan error, 1)
	go func() {
		errCh <- hw.Run(ctx)
	}()
	disp.StartPhase(phaseName, hw.Progress())
	if err := <-errCh; err != nil {
		utils.PrintWarn("phase interrupted", err)
	}
}

func printResult(r installer.Result) {
	if r.Err != nil {
		utils.PrintError(r.Tool, r.Err)
	} else if r.Skipped {
		utils.PrintInfo(fmt.Sprintf("%s: already at %s", r.Tool, r.Version))
	} else if r.WasUpdated {
		utils.PrintSuccess(fmt.Sprintf("%s: updated to %s", r.Tool, r.Version))
	} else {
		utils.PrintSuccess(fmt.Sprintf("%s: installed %s", r.Tool, r.Version))
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
