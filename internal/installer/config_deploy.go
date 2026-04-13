package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/configs"
	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type ConfigDeployInstaller struct{}

func (c *ConfigDeployInstaller) resolveConfig(tool *registry.Tool, p platform.Platform) (content []byte, destPath string, err error) {
	switch tool.Name {
	case "tmux-config":
		content = configs.TmuxConf()
		destPath = filepath.Join(p.HomeDir, ".tmux.conf")
		if p.OS == platform.Darwin {
			content = []byte(strings.ReplaceAll(
				string(content),
				`MouseDragEnd1Pane send -X copy-selection`,
				`MouseDragEnd1Pane send-keys -X copy-pipe "pbcopy"`,
			))
		}

	case "kitty-config":
		if p.OS == platform.Linux {
			content = configs.LinuxKittyConf()
		} else {
			content = configs.MacosKittyConf()
		}
		destPath = filepath.Join(p.HomeDir, ".config", "kitty", "kitty.conf")

	case "kitty-theme":
		destPath = filepath.Join(p.HomeDir, ".config", "kitty", "current-theme.conf")
		content = configs.MochaKittyConf()

	case "aerospace-config":
		if p.OS != platform.Darwin {
			return nil, "", nil // skip on non-macOS
		}
		content = configs.MacosAerospaceConf()
		destPath = filepath.Join(p.HomeDir, ".aerospace.toml")

	case "rcfile":
		// Loader goes to ~/.zshrc; base fragment deployed separately in Install
		content = configs.RcLoader()
		destPath = filepath.Join(p.HomeDir, ".zshrc")

	default:
		return nil, "", fmt.Errorf("unknown config: %s", tool.Name)
	}

	return content, destPath, nil
}

func (c *ConfigDeployInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	content, destPath, err := c.resolveConfig(tool, p)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	if content == nil {
		return Result{Tool: tool.Name, Skipped: true}
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	if tool.Name == "rcfile" {
		baseFragPath := filepath.Join(p.ShellDir(), "rc", "00-base.zsh")
		if err := os.MkdirAll(filepath.Dir(baseFragPath), 0755); err != nil {
			return Result{Tool: tool.Name, Err: fmt.Errorf("create rc dir: %w", err)}
		}
		if err := os.WriteFile(baseFragPath, configs.RcBase(), 0644); err != nil {
			return Result{Tool: tool.Name, Err: fmt.Errorf("write base fragment: %w", err)}
		}
	}

	st.SetToolVersion(tool.Name, "deployed")
	return Result{Tool: tool.Name, Version: "deployed"}
}
