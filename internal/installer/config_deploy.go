package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tanq16/cli-productivity-suite/configs"
	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type ConfigDeployInstaller struct{}

func (c *ConfigDeployInstaller) Check(tool *registry.Tool, _ platform.Platform, _ *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)
	return current, "embedded", nil
}

func (c *ConfigDeployInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	utils.PrintInfo(fmt.Sprintf("deploying config: %s", tool.Name))

	var content []byte
	var destPath string

	switch tool.Name {
	case "tmux-config":
		content = configs.TmuxConf()
		destPath = filepath.Join(p.HomeDir, ".tmux.conf")
		// macOS patch: replace copy-selection with copy-pipe pbcopy
		if p.OS == platform.Darwin {
			content = []byte(strings.ReplaceAll(
				string(content),
				`bind-key -T copy-mode MouseDragEnd1Pane send -X copy-selection`,
				`bind-key -T copy-mode MouseDragEnd1Pane send-keys -X copy-pipe "pbcopy"`,
			))
		}

	case "kitty-config":
		if p.OS == platform.Linux {
			content = configs.LinuxKittyConf()
		} else {
			content = configs.MacosKittyConf()
		}
		destPath = filepath.Join(p.HomeDir, ".config", "kitty", "kitty.conf")

	case "aerospace-config":
		if p.OS != platform.Darwin {
			return Result{Tool: tool.Name, Skipped: true}
		}
		content = configs.MacosAerospaceConf()
		destPath = filepath.Join(p.HomeDir, ".aerospace.toml")

	case "rcfile":
		content = configs.Rcfile()
		destPath = filepath.Join(p.HomeDir, ".zshrc")

	default:
		return Result{Tool: tool.Name, Err: fmt.Errorf("unknown config: %s", tool.Name)}
	}

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	// Write config file (overwrite)
	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	st.SetToolVersion(tool.Name, "deployed")
	return Result{Tool: tool.Name, Version: "deployed"}
}
