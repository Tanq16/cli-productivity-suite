package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/display"
	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/highway"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize full dev environment setup",
	RunE:  runInit,
}


func runInit(cmd *cobra.Command, args []string) error {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintError("platform detection failed", err)
		return fmt.Errorf("platform detection failed: %w", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintError("failed to load state", err)
		return fmt.Errorf("failed to load state: %w", err)
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer cancel()

	gh := github.NewClient(ghToken)
	reg := registry.New()

	// Phase 1: Prerequisites
	utils.PrintInfo("Phase 1: Checking prerequisites")
	omzDir := filepath.Join(p.HomeDir, ".oh-my-zsh")
	if _, err := os.Stat(omzDir); os.IsNotExist(err) {
		msg := fmt.Sprintf("Oh My Zsh not found at %s\nInstall it first: sh -c \"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\"", omzDir)
		utils.PrintError(msg, nil)
		return fmt.Errorf("%s", msg)
	}
	if _, err := exec.LookPath("git"); err != nil {
		utils.PrintError("git not found in PATH", nil)
		return fmt.Errorf("git not found in PATH")
	}
	os.MkdirAll(p.ShellDir(), 0755)
	os.MkdirAll(p.ShellExecDir(), 0755)
	utils.PrintSuccess("prerequisites OK")

	// --- Sudo phases first (2-4) ---
	if phaseNeedsSudo(p, registry.SystemPackage, registry.CloudCLI, registry.LanguageRuntime) {
		utils.PrintInfo("some phases require sudo — authenticating")
		if err := ensureSudo(); err != nil {
			utils.PrintError("sudo authentication failed", err)
			return fmt.Errorf("sudo authentication failed: %w", err)
		}
	}

	disp := display.New()

	// Phase 2: System Packages (sudo on Linux)
	runPhase(ctx, disp, "Phase 2: System packages", filterPlatformTools(reg.ByKind(registry.SystemPackage), p), 1, p, gh, st)
	st.Save()

	// Phase 3: Cloud CLIs (sudo on Linux for aws/azure)
	runPhase(ctx, disp, "Phase 3: Cloud CLIs", reg.ByKind(registry.CloudCLI), 1, p, gh, st)
	st.Save()

	// Phase 4: Go SDK (sudo on both platforms — installs to /usr/local/go)
	goSDK := filterByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	runPhase(ctx, disp, "Phase 4: Go SDK", goSDK, 1, p, gh, st)
	st.Save()

	// --- Sudo no longer needed past this point ---

	// Phase 5: Public GitHub Release Binaries (concurrent)
	publicGH := filterGitHubPublic(reg.ByKind(registry.GitHubRelease), false)
	publicGH = filterPlatformTools(publicGH, p)
	runPhase(ctx, disp, "Phase 5: Public GitHub releases", publicGH, 2, p, gh, st)
	st.Save()

	// Phase 6: Direct Downloads
	runPhase(ctx, disp, "Phase 6: Direct downloads", reg.ByKind(registry.DirectDownload), 1, p, gh, st)
	st.Save()

	// Phase 7: Own Public Tools (concurrent)
	ownPublic := filterOwnPublic(reg.ByKind(registry.GitHubRelease))
	runPhase(ctx, disp, "Phase 7: Own public tools", ownPublic, 2, p, gh, st)
	st.Save()

	// Phase 8: Own Private Tools
	if ghToken != "" {
		privateTools := reg.ByCategory(registry.Private)
		runPhase(ctx, disp, "Phase 8: Private tools", privateTools, 2, p, gh, st)
	} else {
		utils.PrintWarn("Phase 8: Skipping private tools (no --gh-token)", nil)
	}
	st.Save()

	// Phase 9: Neovim + Python (remaining runtimes, no sudo)
	nonSudoRuntimes := excludeByName(reg.ByKind(registry.LanguageRuntime), "go-sdk")
	runPhase(ctx, disp, "Phase 9: Language runtimes", nonSudoRuntimes, 1, p, gh, st)
	st.Save()

	// Phase 10: Shell Plugins (sequential)
	disableOMZSlowPaste(p)
	runPhase(ctx, disp, "Phase 10: Shell plugins", reg.ByKind(registry.ShellPlugin), 1, p, gh, st)
	st.Save()

	// Phase 11: Config Deployment
	runPhase(ctx, disp, "Phase 11: Config files", filterPlatformTools(reg.ByKind(registry.ConfigFile), p), 1, p, gh, st)
	st.Save()

	// Phase 12: Post-Install
	utils.PrintInfo("Phase 12: Running post-install tasks")

	// TPM install plugins
	tpmInstall := filepath.Join(p.HomeDir, ".tmux", "plugins", "tpm", "bin", "install_plugins")
	if _, err := os.Stat(tpmInstall); err == nil {
		utils.PrintInfo("running TPM install_plugins")
		tpmCmd := exec.Command("bash", tpmInstall)
		tpmCmd.Env = append(os.Environ(), fmt.Sprintf("TMUX_PLUGIN_MANAGER_PATH=%s", filepath.Join(p.HomeDir, ".tmux", "plugins")))
		if err := utils.RunCmd(tpmCmd); err != nil {
			utils.PrintWarn("TPM install_plugins failed", err)
		}
	}

	// NvChad headless setup
	if _, err := exec.LookPath("nvim"); err == nil {
		utils.PrintInfo("running NvChad headless setup")
		nvimCmd := exec.Command("nvim", "--headless", "+MasonInstallAll", "+Lazy sync", "+qa")
		if err := utils.RunCmd(nvimCmd); err != nil {
			utils.PrintWarn("NvChad headless setup failed", err)
		}
	}

	// Save state
	st.LastInit = time.Now()
	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	utils.PrintSuccess("init complete!")
	return nil
}

func runPhase(ctx context.Context, disp *display.Display, phaseName string, tools []registry.Tool, workers int, p platform.Platform, gh *github.Client, st *state.State) {
	if len(tools) == 0 {
		return
	}
	hw := highway.New(workers)
	jobs := make([]highway.Job, len(tools))
	for i, t := range tools {
		jobs[i] = NewInstallJob(t, p, gh, st)
	}
	hw.Submit(jobs...)
	go hw.Run(ctx)
	disp.StartPhase(phaseName, hw.Progress())
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
	return len(repo) > 7 && repo[:7] == "Tanq16/"
}

func disableOMZSlowPaste(p platform.Platform) {
	miscPath := filepath.Join(p.HomeDir, ".oh-my-zsh", "lib", "misc.zsh")
	data, err := os.ReadFile(miscPath)
	if err != nil {
		return
	}
	content := string(data)
	// Only comment out lines that aren't already commented
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
	os.WriteFile(miscPath, []byte(strings.Join(lines, "\n")), 0644)
}
