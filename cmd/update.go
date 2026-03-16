package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var updateFlags struct {
	public  bool
	private bool
	system  bool
	all     bool
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
	rootCmd.AddCommand(updateCmd)
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

	// Filter to only installed tools and GitHub release / direct download types
	var updatable []registry.Tool
	for _, t := range tools {
		if st.ToolVersion(t.Name) == "" {
			continue
		}
		if t.Kind == registry.GitHubRelease || t.Kind == registry.DirectDownload {
			updatable = append(updatable, t)
		}
	}

	if len(updatable) == 0 {
		utils.PrintWarn("no updatable tools found", nil)
		return nil
	}

	utils.PrintInfo(fmt.Sprintf("updating %d tools", len(updatable)))

	g, ctx := errgroup.WithContext(cmd.Context())
	g.SetLimit(5)
	results := make([]installer.Result, len(updatable))

	for i, tool := range updatable {
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

	var updated, skipped, errored int
	for _, r := range results {
		printResult(r)
		if r.Err != nil {
			errored++
		} else if r.Skipped {
			skipped++
		} else {
			updated++
		}
	}

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	utils.PrintSuccess(fmt.Sprintf("update complete: %d updated, %d up-to-date, %d errors", updated, skipped, errored))
	return nil
}
