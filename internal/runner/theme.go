package runner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tanq16/cli-productivity-suite/internal/configs"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func ThemeNames() []string {
	return configs.ThemeNames()
}

func ThemeList() {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("Failed to detect platform", err)
	}
	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("Failed to load state", err)
	}
	current := st.CurrentTheme()
	if current == "" {
		current = configs.DefaultTheme
	}

	for _, name := range configs.ThemeNames() {
		if name == current {
			utils.PrintSuccess(name + " (active)")
			continue
		}
		utils.PrintGeneric("  " + name)
	}
}

func Theme(name string) {
	content, err := configs.Theme(name)
	if err != nil {
		utils.PrintFatal(fmt.Sprintf("unknown theme %q — run `cps theme` to list them", name), err)
	}
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("Failed to detect platform", err)
	}
	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("Failed to load state", err)
	}

	utils.PrintRunning("applying " + name)
	destPath := filepath.Join(p.HomeDir, ".config", "kitty", "current-theme.conf")
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		utils.PrintFatal("Failed to create kitty config directory", err)
	}
	if err := os.WriteFile(destPath, content, 0644); err != nil {
		utils.PrintFatal("Failed to write theme", err)
	}
	st.SetTheme(name)
	if err := st.Save(); err != nil {
		utils.PrintFatal("Failed to record theme", err)
	}
	utils.ClearLines(1)

	utils.PrintSuccess(fmt.Sprintf("theme set to %s", name))
	utils.PrintGeneric("  kitty reloads within a second; running tmux and nvim follow on next redraw")
}
