package platform

import (
	"os"
	"os/exec"
	"syscall"
)

func NeedsSudo() bool {
	return os.Getuid() != 0
}

func RunWithSudo(args ...string) error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	cmdArgs := append([]string{self}, args...)
	cmd := exec.Command("sudo", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	return cmd.Run()
}
