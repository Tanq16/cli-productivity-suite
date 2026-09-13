package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var aptPackages = []string{
	"git", "curl", "zsh", "build-essential", "wget", "zip", "unzip",
	"file", "tmux", "htop", "cmake", "ninja-build", "gettext",
}

var brewPackages = []string{"wget", "zip", "unzip", "file", "tmux", "htop"}

const brewInstallScript = `set -o pipefail
curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh | /bin/bash`

func Prereq(excludeBrew bool) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	if p.OS == platform.Darwin {
		prereqDarwin()
		return
	}
	prereqLinux(excludeBrew)
}

func prereqDarwin() {
	utils.PrintInfo("Xcode Command Line Tools")
	if err := streamCmd(exec.Command("xcode-select", "--install")); err != nil {
		utils.PrintWarn("xcode-select reported an error, which is what it does when the tools are already installed", err)
	}

	installBrew()

	brew := platform.BrewPath()
	if brew == "" {
		utils.PrintFatal("Homebrew is not installed and the install script did not leave a brew binary behind", nil)
	}
	utils.PrintInfo("Homebrew packages")
	args := append([]string{"install"}, brewPackages...)
	if err := streamCmd(exec.Command(brew, args...)); err != nil {
		utils.PrintFatal("brew install failed", err)
	}
	utils.PrintSuccess("prerequisites complete!")
}

func prereqLinux(excludeBrew bool) {
	if !isDebianFamily() {
		utils.PrintFatal("cps prereq only handles Debian and Ubuntu; install these with your package manager, then re-run `cps shell`: "+strings.Join(aptPackages, " "), nil)
	}

	utils.PrintInfo("apt-get update")
	if err := streamCmd(exec.Command("sudo", "apt-get", "update")); err != nil {
		utils.PrintFatal("apt-get update failed", err)
	}

	utils.PrintInfo("apt-get install")
	args := append([]string{"apt-get", "install", "-y"}, aptPackages...)
	if err := streamCmd(exec.Command("sudo", args...)); err != nil {
		utils.PrintFatal("apt-get install failed", err)
	}

	if excludeBrew {
		utils.PrintSuccess("prerequisites complete, Homebrew excluded!")
		return
	}
	installBrew()
	utils.PrintSuccess("prerequisites complete!")
}

func installBrew() {
	if platform.BrewPath() != "" {
		utils.PrintInfo("Homebrew is already installed")
		return
	}
	utils.PrintInfo("Homebrew")
	cmd := exec.Command("/bin/bash", "-c", brewInstallScript)
	cmd.Env = append(os.Environ(), "NONINTERACTIVE=1")
	if err := streamCmd(cmd); err != nil {
		utils.PrintFatal("Homebrew install failed", err)
	}
}

func isDebianFamily() bool {
	body, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}
	for line := range strings.SplitSeq(string(body), "\n") {
		key, value, found := strings.Cut(strings.TrimSpace(line), "=")
		if !found {
			continue
		}
		value = strings.Trim(value, `"`)
		if key == "ID" && (value == "debian" || value == "ubuntu") {
			return true
		}
		if key == "ID_LIKE" && strings.Contains(value, "debian") {
			return true
		}
	}
	return false
}

func streamCmd(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", cmd.Path, err)
	}
	return nil
}
