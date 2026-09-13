package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func Custom(args []string) {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	script := filepath.Join(p.ConfigDir(), "custom.sh")
	info, err := os.Stat(script)
	if err != nil {
		utils.PrintFatal(fmt.Sprintf("no custom script at %s", script), err)
	}
	if info.Mode()&0111 == 0 {
		utils.PrintFatal(fmt.Sprintf("%s is not executable, run: chmod +x %s", script, script), nil)
	}

	cmd := exec.Command(script, args...)
	cmd.Env = p.CustomScriptEnv()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(exitErr.ExitCode())
		}
		utils.PrintFatal("custom script failed", err)
	}
}
