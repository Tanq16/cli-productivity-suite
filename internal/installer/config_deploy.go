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

// resolveConfig returns the embedded config content and destination path for
// the given tool. Returns (nil, "", nil) when the config should be skipped on
// the current platform (e.g. aerospace on Linux).
func (c *ConfigDeployInstaller) resolveConfig(tool *registry.Tool, p platform.Platform) (content []byte, destPath string, err error) {
	switch tool.Name {
	case "tmux-config":
		content = configs.TmuxConf()
		destPath = filepath.Join(p.HomeDir, ".tmux.conf")
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
			return nil, "", nil // skip on non-macOS
		}
		content = configs.MacosAerospaceConf()
		destPath = filepath.Join(p.HomeDir, ".aerospace.toml")

	case "rcfile":
		content = configs.Rcfile()
		destPath = filepath.Join(p.HomeDir, ".zshrc")

	default:
		return nil, "", fmt.Errorf("unknown config: %s", tool.Name)
	}

	return content, destPath, nil
}

// configLinesMatch performs a prefix comparison: the deployed file must start
// with the same lines as the embedded content. Extra user lines appended after
// the embedded portion are allowed.
func configLinesMatch(embedded, deployed []byte) bool {
	embeddedLines := strings.Split(string(embedded), "\n")
	deployedLines := strings.Split(string(deployed), "\n")
	if len(deployedLines) < len(embeddedLines) {
		return false
	}
	for i := range embeddedLines {
		if embeddedLines[i] != deployedLines[i] {
			return false
		}
	}
	return true
}

func (c *ConfigDeployInstaller) Check(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)

	embedded, destPath, err := c.resolveConfig(tool, p)
	if err != nil {
		return current, "error", err
	}
	if embedded == nil {
		return current, "skipped", nil
	}

	deployed, err := os.ReadFile(destPath)
	if os.IsNotExist(err) {
		return current, "not-deployed", nil
	}
	if err != nil {
		return current, "error", err
	}

	if configLinesMatch(embedded, deployed) {
		return current, "up-to-date", nil
	}
	return current, "config-differs", nil
}

func (c *ConfigDeployInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	utils.PrintInfo(fmt.Sprintf("deploying config: %s", tool.Name))

	content, destPath, err := c.resolveConfig(tool, p)
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	if content == nil {
		return Result{Tool: tool.Name, Skipped: true}
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
