package runner

import (
	"os/exec"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/utils"
)

const kittyInstallerURL = "https://sw.kovidgoyal.net/kitty/installer.sh"

var kittyInstallScript = `set -o pipefail
curl -fsSL ` + kittyInstallerURL + ` | sh /dev/stdin launch=n`

var SystemVerbs = []string{"update", "kitty"}

func System(verb string) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	switch verb {
	case "update":
		systemUpdate(p)
	case "kitty":
		systemKitty()
	}
}

func systemUpdate(p platform.Platform) {
	if p.OS == platform.Linux {
		if !isDebianFamily() {
			utils.PrintFatal("cps system update only handles Debian and Ubuntu", nil)
		}
		for _, args := range [][]string{
			{"apt-get", "update"},
			{"apt-get", "upgrade", "-y"},
			{"apt-get", "dist-upgrade", "-y"},
			{"apt-get", "autoremove", "-y"},
			{"apt-get", "autoclean"},
		} {
			label := "sudo " + strings.Join(args, " ")
			utils.PrintInfo(label)
			if err := streamCmd(exec.Command("sudo", args...)); err != nil {
				utils.PrintFatal(label+" failed", err)
			}
		}
	}

	brew := platform.BrewPath()
	if brew == "" {
		utils.PrintSuccess("system update complete!")
		return
	}
	for _, args := range [][]string{
		{"update"},
		{"upgrade"},
		{"upgrade", "--cask"},
		{"cleanup"},
	} {
		label := "brew " + strings.Join(args, " ")
		utils.PrintInfo(label)
		if err := streamCmd(exec.Command(brew, args...)); err != nil {
			utils.PrintFatal(label+" failed", err)
		}
	}
	utils.PrintSuccess("system update complete!")
}

func systemKitty() {
	utils.PrintInfo("running the kitty installer from " + kittyInstallerURL)
	if err := streamCmd(exec.Command("/bin/bash", "-c", kittyInstallScript)); err != nil {
		utils.PrintFatal("kitty install failed", err)
	}
	utils.PrintSuccess("kitty installed!")
}
