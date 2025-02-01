package system

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func CheckRoot() bool {
	return os.Geteuid() == 0
}

func EnsureRoot() error {
	if !CheckRoot() {
		log.Info().Msg("Requesting root privileges...")
		cmd := exec.Command("sudo", os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}
