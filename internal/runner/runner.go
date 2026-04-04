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

	utils.PrintRunning("(Running) Phase 1: Checking prerequisites")
	omzDir := filepath.Join(p.HomeDir, ".oh-my-zsh")
	if _, err := os.Stat(omzDir); os.IsNotExist(err) {
		msg := fmt.Sprintf("Oh My Zsh not found at %s\nInstall it first: sh -c \"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\"", omzDir)
		utils.PrintFatal(msg, nil)
	}
	if _, err := exec.LookPath("git"); err != nil {
		utils.PrintFatal("git not found in PATH", err)
	}
	if err := os.MkdirAll(p.ShellDir(), 0755); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to create %s", p.ShellDir()), err)
	}
	if err := os.MkdirAll(p.ShellExecDir(), 0755); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to create %s", p.ShellExecDir()), err)
	}
	utils.PrintIndentedSuccess("prerequisites OK")

	var sudoDone chan struct{}
	if PhaseNeedsSudo(p, registry.SystemPackage, registry.CloudCLI, registry.LanguageRuntime) {
		cached := exec.Command("sudo", "-n", "-v").Run() == nil
		utils.ClearLines(2)
		utils.PrintRunning("(Running) Phase 1: Authenticating sudo")
		if err := EnsureSudo(); err != nil {
			utils.PrintFatal("sudo authentication failed", err)
		}
		sudoDone = StartSudoRefresh()
		if cached {
			utils.ClearLines(1)
		} else {
			utils.ClearLines(2)
		}
	} else {
		utils.ClearLines(2)
	}
	utils.PrintInfo("Phase 1: Checking prerequisites")

	var hadErrors bool

	if p.OS == platform.Linux {
		cmd := exec.Command("sudo", "apt-get", "update", "-qq")
		utils.RunCmd(cmd)
	}

	if runPhase("Phase 2: System packages", filterPlatformTools(reg.ByKind(registry.SystemPackage), p), p, gh, st) {
		hadErrors = true
	}
	st.Save()

	if runPhase("Phase 3: Cloud CLIs", reg.ByKind(registry.CloudCLI), p, gh, st) {
		hadErrors = true
	}
	st.Save()

	goSDK := filterByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	if runPhase("Phase 4: Go SDK", goSDK, p, gh, st) {
		hadErrors = true
	}
	st.Save()

	publicGH := filterGitHubPublic(reg.ByKind(registry.GitHubRelease), false)
	publicGH = filterPlatformTools(publicGH, p)
	if runPhase("Phase 5: Public GitHub releases", publicGH, p, gh, st) {
		hadErrors = true
	}
	st.Save()

	if runPhase("Phase 6: Direct downloads", filterNonExtension(reg.ByKind(registry.DirectDownload)), p, gh, st) {
		hadErrors = true
	}
	st.Save()

	ownPublic := filterOwnPublic(reg.ByKind(registry.GitHubRelease))
	if runPhase("Phase 7: Own public tools", ownPublic, p, gh, st) {
		hadErrors = true
	}
	st.Save()

	nonSudoRuntimes := excludeByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	if runPhase("Phase 9: Language runtimes", nonSudoRuntimes, p, gh, st) {
		hadErrors = true
	}
	st.Save()

	disableOMZSlowPaste(p)
	if runPhase("Phase 10: Shell plugins", reg.ByKind(registry.ShellPlugin), p, gh, st) {
		hadErrors = true
	}
	st.Save()

	if runPhase("Phase 11: Config files", filterPlatformTools(reg.ByKind(registry.ConfigFile), p), p, gh, st) {
		hadErrors = true
	}
	st.Save()

	runPostInstall(p)

	if sudoDone != nil {
		close(sudoDone)
	}

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

func Check(ghToken string, appVersion string, skipPrivate bool) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}

	gh := github.NewClient(ghToken)
	tools := registry.New().ForPlatform(p.OS.String())

	var checkable []registry.Tool
	for _, t := range tools {
		if t.IsPrivate && (ghToken == "" || skipPrivate) {
			continue
		}
		checkable = append(checkable, t)
	}

	utils.PrintRunning("Checking tools")

	// Check CPS version inline
	var cpsResult *CheckResult
	if rel, err := gh.LatestRelease("Tanq16/cli-productivity-suite"); err == nil {
		if appVersion == "dev-build" {
			cpsResult = &CheckResult{Current: appVersion, Latest: rel.TagName, Status: "update"}
		} else if appVersion != rel.TagName {
			cpsResult = &CheckResult{Current: appVersion, Latest: rel.TagName, Status: "update"}
		}
	}

	results := checkTools(checkable, p, gh, st)
	utils.ClearLines(1)

	hasResults := len(results) > 0 || cpsResult != nil
	if !hasResults {
		utils.PrintSuccess("everything is up to date")
		return
	}

	utils.PrintInfo("Check complete")
	if cpsResult != nil {
		if appVersion == "dev-build" {
			utils.PrintIndentedWarn(fmt.Sprintf("cps: dev build (latest: %s)", cpsResult.Latest), nil)
		} else {
			utils.PrintIndentedWarn(fmt.Sprintf("cps: update available (%s → %s)", cpsResult.Current, cpsResult.Latest), nil)
		}
	}
	for _, r := range results {
		switch r.Status {
		case "update":
			utils.PrintIndentedWarn(fmt.Sprintf("%s: update available (%s → %s)", r.Tool.Name, r.Current, r.Latest), nil)
		case "error":
			utils.PrintIndentedError(fmt.Sprintf("%s: check failed", r.Tool.Name), r.Err)
		case "config-differs":
			utils.PrintIndentedWarn(fmt.Sprintf("%s: config differs", r.Tool.Name), nil)
		case "not-deployed":
			utils.PrintIndentedWarn(fmt.Sprintf("%s: not deployed", r.Tool.Name), nil)
		}
	}
}

var categoryAliases = map[string]func(*registry.Registry, platform.Platform) []registry.Tool{
	"public": func(reg *registry.Registry, p platform.Platform) []registry.Tool {
		ghPublic := filterNonExtension(filterGitHubPublic(reg.ByKind(registry.GitHubRelease), false))
		ownPublic := filterNonExtension(filterOwnPublic(reg.ByKind(registry.GitHubRelease)))
		dd := filterNonExtension(reg.ByKind(registry.DirectDownload))
		combined := append(ghPublic, ownPublic...)
		combined = append(combined, dd...)
		return filterPlatformTools(combined, p)
	},
	"system": func(reg *registry.Registry, p platform.Platform) []registry.Tool {
		return filterPlatformTools(reg.ByCategory(registry.System), p)
	},
	"cloud": func(reg *registry.Registry, p platform.Platform) []registry.Tool {
		return filterPlatformTools(reg.ByCategory(registry.CloudCLICat), p)
	},
	"runtimes": func(reg *registry.Registry, p platform.Platform) []registry.Tool {
		return filterPlatformTools(reg.ByCategory(registry.Runtime), p)
	},
	"configs": func(reg *registry.Registry, p platform.Platform) []registry.Tool {
		configs := reg.ByCategory(registry.Config)
		shell := reg.ByCategory(registry.Shell)
		combined := append(configs, shell...)
		return filterPlatformTools(combined, p)
	},
}

func CategoryAliasNames() []string {
	return []string{"public", "system", "cloud", "runtimes", "configs"}
}

func Install(args []string, ghToken string) {
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

	var resolved []registry.Tool
	for _, arg := range args {
		if resolver, ok := categoryAliases[arg]; ok {
			resolved = append(resolved, resolver(reg, p)...)
		} else if tool := reg.Get(arg); tool != nil {
			if !toolForPlatform(*tool, p) {
				utils.PrintFatal(fmt.Sprintf("tool %s is not available on %s", arg, p.OS), nil)
			}
			resolved = append(resolved, *tool)
		} else {
			utils.PrintFatal(fmt.Sprintf("unknown tool or category: %s", arg), nil)
		}
	}

	seen := map[string]bool{}
	var tools []registry.Tool
	for _, t := range resolved {
		if !seen[t.Name] {
			seen[t.Name] = true
			tools = append(tools, t)
		}
	}

	if len(tools) == 0 {
		utils.PrintSuccess("no tools to install for this platform")
		return
	}

	for _, t := range tools {
		if t.IsPrivate && ghToken == "" {
			utils.PrintFatal(fmt.Sprintf("tool %s is private — provide --gh-token or set CPS_GITHUB_PAT", t.Name), nil)
		}
	}

	var sudoDone chan struct{}
	needsSudo := false
	for _, t := range tools {
		if ToolNeedsSudo(t, p) {
			needsSudo = true
			break
		}
	}
	if needsSudo {
		cached := exec.Command("sudo", "-n", "-v").Run() == nil
		utils.PrintRunning("authenticating sudo")
		if err := EnsureSudo(); err != nil {
			utils.PrintFatal("sudo authentication failed", err)
		}
		sudoDone = StartSudoRefresh()
		if cached {
			utils.ClearLines(1)
		} else {
			utils.ClearLines(2)
		}
	}

	var hasErrors bool
	if len(tools) == 1 {
		tool := tools[0]
		inst := installer.Dispatch(tool.Kind)
		if inst == nil {
			utils.PrintFatal(fmt.Sprintf("no installer for kind: %s", tool.Kind), nil)
		}
		utils.PrintRunning("installing " + tool.Name)
		st.Remove(tool.Name)
		result := inst.Install(&tool, p, gh, st)
		utils.ClearLines(1)
		if result.Err != nil {
			utils.PrintFatal(fmt.Sprintf("%s: install failed", tool.Name), result.Err)
		}
		utils.PrintSuccess(fmt.Sprintf("%s: installed %s", tool.Name, result.Version))
	} else {
		hasErrors = runPhase("Installing tools", tools, p, gh, st)
	}

	installPostTasks(tools, p)

	if sudoDone != nil {
		close(sudoDone)
	}

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	if len(tools) > 1 {
		if hasErrors {
			utils.PrintWarn("install finished with errors", nil)
		} else {
			utils.PrintSuccess("install complete!")
		}
	}
}

func installPostTasks(tools []registry.Tool, p platform.Platform) {
	var hasConfigs, hasBinaries bool
	for _, t := range tools {
		switch t.Category {
		case registry.Config, registry.Shell:
			hasConfigs = true
		}
		switch t.Kind {
		case registry.GitHubRelease, registry.DirectDownload:
			hasBinaries = true
		}
	}

	if hasConfigs {
		tpmInstall := filepath.Join(p.HomeDir, ".tmux", "plugins", "tpm", "bin", "install_plugins")
		if _, err := os.Stat(tpmInstall); err == nil {
			utils.PrintRunning("running tpm install")
			tpmCmd := exec.Command("bash", tpmInstall)
			tpmCmd.Env = append(os.Environ(), fmt.Sprintf("TMUX_PLUGIN_MANAGER_PATH=%s", filepath.Join(p.HomeDir, ".tmux", "plugins")))
			utils.RunCmd(tpmCmd)
			utils.ClearLines(1)
		}

		if _, err := exec.LookPath("nvim"); err == nil {
			utils.PrintRunning("running nvchad setup")
			nvimCmd := exec.Command("nvim", "--headless", "+MasonInstallAll", "+Lazy sync", "+qa")
			utils.RunCmd(nvimCmd)
			utils.ClearLines(1)
		}
	}

	if hasBinaries {
		var compErrors []jobResult
		var compLines int
		utils.PrintRunning("(Running) Regenerating completions")
		generateCompletions(p, &compErrors, &compLines)
		utils.ClearLines(compLines + 1)
		if len(compErrors) > 0 {
			utils.PrintError("Regenerating completions: partially completed with errors", nil)
			for _, e := range compErrors {
				utils.PrintIndentedError(e.name, e.err)
			}
		} else {
			utils.PrintInfo("Regenerating completions")
		}
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

	cached := exec.Command("sudo", "-n", "-v").Run() == nil
	utils.PrintRunning("authenticating sudo")
	if err := EnsureSudo(); err != nil {
		utils.PrintFatal("sudo authentication failed", err)
	}
	sudoDone := StartSudoRefresh()
	if cached {
		utils.ClearLines(1)
	} else {
		utils.ClearLines(2)
	}

	utils.PrintRunning(fmt.Sprintf("downloading %s", release.TagName))
	tmpDir, err := os.MkdirTemp("", "cps-self-update-*")
	if err != nil {
		utils.PrintFatal("failed to create temp dir", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpBinary := filepath.Join(tmpDir, "cps")
	if err := installer.DownloadToFile(downloadURL, tmpBinary); err != nil {
		utils.PrintFatal("download failed", err)
	}
	if err := os.Chmod(tmpBinary, 0755); err != nil {
		utils.PrintFatal("chmod failed", err)
	}

	destPath, err := exec.LookPath("cps")
	if err != nil {
		destPath = "/usr/local/bin/cps"
	}
	rmCmd := exec.Command("sudo", "rm", "-f", destPath)
	if err := rmCmd.Run(); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to remove old binary at %s", destPath), err)
	}
	cpCmd := exec.Command("sudo", "cp", tmpBinary, destPath)
	var stderr strings.Builder
	cpCmd.Stderr = &stderr
	if err := cpCmd.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail != "" {
			err = fmt.Errorf("%s: %w", detail, err)
		}
		utils.PrintFatal(fmt.Sprintf("failed to install binary at %s", destPath), err)
	}
	utils.ClearLines(1)

	close(sudoDone)
	utils.PrintSuccess(fmt.Sprintf("updated cps: %s → %s", appVersion, release.TagName))
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
		utils.PrintInfo("  " + d)
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

func checkTools(tools []registry.Tool, p platform.Platform, gh *github.Client, st *state.State) []CheckResult {
	var results []CheckResult
	for _, t := range tools {
		if st.ToolVersion(t.Name) == "" {
			continue
		}
		inst := installer.Dispatch(t.Kind)
		if inst == nil {
			continue
		}
		cur, lat, err := inst.Check(&t, p, gh, st)
		if err != nil {
			results = append(results, CheckResult{Tool: t, Current: cur, Latest: "", Status: "error", Err: err})
			continue
		}
		switch lat {
		case "skipped", "system-managed", "check-manually", "git-managed", "up-to-date":
			continue
		case "not-deployed":
			results = append(results, CheckResult{Tool: t, Current: cur, Latest: lat, Status: "not-deployed"})
		case "config-differs":
			results = append(results, CheckResult{Tool: t, Current: cur, Latest: lat, Status: "config-differs"})
		default:
			if cur != lat && lat != "" {
				results = append(results, CheckResult{Tool: t, Current: cur, Latest: lat, Status: "update"})
			}
		}
	}
	return results
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
	utils.PrintRunning("(Running) Phase 12: Post-install tasks")
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

	if _, err := exec.LookPath("nvim"); err == nil {
		utils.PrintIndentedRunning("nvchad-setup: running")
		lineCount++
		nvimCmd := exec.Command("nvim", "--headless", "+MasonInstallAll", "+Lazy sync", "+qa")
		if err := utils.RunCmd(nvimCmd); err != nil {
			errors = append(errors, jobResult{name: "nvchad-setup", err: err})
		}
	}

	generateCompletions(p, &errors, &lineCount)

	utils.ClearLines(lineCount + 1) // sub-lines + running header
	if len(errors) > 0 {
		utils.PrintError("Phase 12: partially completed with errors", nil)
		for _, e := range errors {
			utils.PrintIndentedError(e.name, e.err)
		}
	} else {
		utils.PrintInfo("Phase 12: Post-install tasks")
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
		args    []string
		outFile string
	}

	defs := []compDef{
		{"fzf", "fzf", []string{"--zsh"}, "fzf.zsh"},
		{"uv", "uv", []string{"generate-shell-completion", "zsh"}, "uv.zsh"},
		{"zoxide", "zoxide", []string{"init", "zsh"}, "zoxide.zsh"},
	}

	for _, d := range defs {
		binPath := filepath.Join(p.ShellExecDir(), d.binary)
		if _, err := os.Stat(binPath); err != nil {
			continue
		}
		utils.PrintIndentedRunning("completions: " + d.name)
		*lineCount++
		out, err := exec.Command(binPath, d.args...).Output()
		if err != nil {
			*errors = append(*errors, jobResult{name: "completions-" + d.name, err: err})
			continue
		}
		if err := os.WriteFile(filepath.Join(compDir, d.outFile), out, 0644); err != nil {
			*errors = append(*errors, jobResult{name: "completions-" + d.name, err: err})
		}
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

func filterNonExtension(tools []registry.Tool) []registry.Tool {
	var result []registry.Tool
	for _, t := range tools {
		if !t.Extension {
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
