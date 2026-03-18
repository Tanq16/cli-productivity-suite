package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// RunCmd executes a command with captured stdout/stderr.
// On success, output is discarded (logged in debug mode).
// On error, captured stderr is included in the returned error.
func RunCmd(cmd *exec.Cmd) error {
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	output := strings.TrimSpace(buf.String())
	if GlobalDebugFlag && output != "" {
		log.Debug().Str("package", "utils").Str("cmd", cmd.String()).Msg(output)
	}
	if err != nil {
		if output != "" {
			return fmt.Errorf("%w: %s", err, output)
		}
		return err
	}
	return nil
}
