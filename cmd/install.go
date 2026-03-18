package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var installCmd = &cobra.Command{
	Use:               "install <tool>",
	Short:             "Install a single tool by name",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeToolNames,
	RunE:              runInstall,
}


func runInstall(cmd *cobra.Command, args []string) error {
	toolName := args[0]

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

	tool := reg.Get(toolName)
	if tool == nil {
		utils.PrintError(fmt.Sprintf("unknown tool: %s", toolName), nil)
		return fmt.Errorf("unknown tool: %s", toolName)
	}

	if tool.IsPrivate && ghToken == "" {
		msg := fmt.Sprintf("tool %s is private — provide --gh-token or set CPS_GITHUB_PAT", toolName)
		utils.PrintError(msg, nil)
		return fmt.Errorf("%s", msg)
	}

	inst := installer.Dispatch(tool.Kind)
	if inst == nil {
		utils.PrintError(fmt.Sprintf("no installer for tool kind: %s", tool.Kind), nil)
		return fmt.Errorf("no installer for tool kind: %s", tool.Kind)
	}

	if toolNeedsSudo(*tool, p) {
		utils.PrintInfo("this tool requires sudo — authenticating")
		if err := ensureSudo(); err != nil {
			utils.PrintError("sudo authentication failed", err)
			return fmt.Errorf("sudo authentication failed: %w", err)
		}
	}

	result := inst.Install(tool, p, gh, st)
	printResult(result)

	if result.Err != nil {
		return fmt.Errorf("install failed for %s", toolName)
	}

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	return nil
}

func completeToolNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	reg := registry.New()
	return reg.Names(), cobra.ShellCompDirectiveNoFileComp
}
