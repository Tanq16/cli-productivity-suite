package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/tanq16/cli-productivity-suite/internal/github"
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

func init() {
	rootCmd.AddCommand(initCmd)
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

	// Phase 2: System Packages
	utils.PrintInfo("Phase 2: Installing system packages")
	for _, tool := range reg.ByKind(registry.SystemPackage) {
		if !toolForPlatform(tool, p) {
			continue
		}
		t := tool
		result := installer.Dispatch(t.Kind).Install(&t, p, gh, st)
		printResult(result)
	}

	st.Save()

	// Phase 3: Public GitHub Release Binaries (concurrent)
	utils.PrintInfo("Phase 3: Installing public GitHub release binaries")
	publicGH := filterGitHubPublic(reg.ByKind(registry.GitHubRelease), false)
	publicGH = filterPlatformTools(publicGH, p)
	runConcurrentInstalls(cmd.Context(), publicGH, p, gh, st, 5)
	st.Save()

	// Phase 4: Direct Downloads
	utils.PrintInfo("Phase 4: Installing direct download tools")
	for _, tool := range reg.ByKind(registry.DirectDownload) {
		t := tool
		result := installer.Dispatch(t.Kind).Install(&t, p, gh, st)
		printResult(result)
	}

	st.Save()

	// Phase 5: Own Public Tools (concurrent)
	utils.PrintInfo("Phase 5: Installing own public tools")
	ownPublic := filterOwnPublic(reg.ByKind(registry.GitHubRelease))
	runConcurrentInstalls(cmd.Context(), ownPublic, p, gh, st, 5)
	st.Save()

	// Phase 6: Own Private Tools
	if ghToken != "" {
		utils.PrintInfo("Phase 6: Installing private tools")
		privateTools := reg.ByCategory(registry.Private)
		runConcurrentInstalls(cmd.Context(), privateTools, p, gh, st, 5)
	} else {
		utils.PrintWarn("Phase 6: Skipping private tools (no --gh-token)", nil)
	}
	st.Save()

	// Phase 7: Cloud CLIs
	utils.PrintInfo("Phase 7: Installing cloud CLIs")
	for _, tool := range reg.ByKind(registry.CloudCLI) {
		t := tool
		result := installer.Dispatch(t.Kind).Install(&t, p, gh, st)
		printResult(result)
	}

	st.Save()

	// Phase 8: Language Runtimes
	utils.PrintInfo("Phase 8: Installing language runtimes")
	for _, tool := range reg.ByKind(registry.LanguageRuntime) {
		t := tool
		result := installer.Dispatch(t.Kind).Install(&t, p, gh, st)
		printResult(result)
	}

	st.Save()

	// Phase 9: Shell Plugins (sequential)
	utils.PrintInfo("Phase 9: Installing shell plugins")
	// Disable OMZ slow paste magic
	disableOMZSlowPaste(p)
	for _, tool := range reg.ByKind(registry.ShellPlugin) {
		t := tool
		result := installer.Dispatch(t.Kind).Install(&t, p, gh, st)
		printResult(result)
	}

	st.Save()

	// Phase 10: Config Deployment
	utils.PrintInfo("Phase 10: Deploying configuration files")
	for _, tool := range reg.ByKind(registry.ConfigFile) {
		if !toolForPlatform(tool, p) {
			continue
		}
		t := tool
		result := installer.Dispatch(t.Kind).Install(&t, p, gh, st)
		printResult(result)
	}

	st.Save()

	// Phase 11: Post-Install
	utils.PrintInfo("Phase 11: Running post-install tasks")

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

func runConcurrentInstalls(ctx context.Context, tools []registry.Tool, p platform.Platform, gh *github.Client, st *state.State, limit int) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)
	results := make([]installer.Result, len(tools))

	for i, tool := range tools {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			results[i] = installer.Dispatch(tool.Kind).Install(&tool, p, gh, st)
			return nil
		})
	}
	g.Wait()

	for _, r := range results {
		printResult(r)
	}
}

func printResult(r installer.Result) {
	if r.Err != nil {
		utils.PrintError(r.Tool, r.Err)
	} else if r.Skipped {
		utils.PrintDebug(fmt.Sprintf("%s: already at %s", r.Tool, r.Version))
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
