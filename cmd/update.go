package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/display"
	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var updateFlags struct {
	public      bool
	private     bool
	system      bool
	all         bool
	includeConf bool
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update installed tools to latest versions",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&updateFlags.public, "public", false, "Update public tools only")
	updateCmd.Flags().BoolVar(&updateFlags.private, "private", false, "Update private tools only")
	updateCmd.Flags().BoolVar(&updateFlags.system, "system", false, "Update system packages only")
	updateCmd.Flags().BoolVar(&updateFlags.all, "all", false, "Update all tools")
	updateCmd.Flags().BoolVar(&updateFlags.includeConf, "include-conf", false, "Also overwrite deployed config files")
}

func runUpdate(cmd *cobra.Command, args []string) error {
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

	if !updateFlags.public && !updateFlags.private && !updateFlags.system {
		updateFlags.all = true
	}

	var tools []registry.Tool
	if updateFlags.all {
		tools = reg.ForPlatform(p.OS.String())
	} else {
		if updateFlags.public {
			tools = append(tools, reg.ByCategory(registry.Public)...)
		}
		if updateFlags.private {
			tools = append(tools, reg.ByCategory(registry.Private)...)
		}
		if updateFlags.system {
			tools = append(tools, reg.ByCategory(registry.System)...)
		}
	}

	// Filter to only installed tools; skip ConfigFile unless --include-conf
	var updatable []registry.Tool
	for _, t := range tools {
		if st.ToolVersion(t.Name) == "" {
			continue
		}
		if t.Kind == registry.ConfigFile && !updateFlags.includeConf {
			continue
		}
		updatable = append(updatable, t)
	}

	if len(updatable) == 0 {
		utils.PrintWarn("no updatable tools found", nil)
		return nil
	}

	// Check if any updatable tool needs sudo and pre-prompt
	needsSudo := false
	for _, t := range updatable {
		if toolNeedsSudo(t, p) {
			needsSudo = true
			break
		}
	}
	if needsSudo {
		utils.PrintInfo("some tools require sudo — authenticating")
		if err := ensureSudo(); err != nil {
			utils.PrintError("sudo authentication failed", err)
			return fmt.Errorf("sudo authentication failed: %w", err)
		}
	}

	utils.PrintInfo(fmt.Sprintf("updating %d tools", len(updatable)))

	// Use single worker when sudo tools are present to avoid concurrent sudo calls
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
	return nil
}
