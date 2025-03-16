package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// CommandResult stores the result of command execution
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

// ExecuteCommand runs a command and returns its result
func ExecuteCommand(command string, args ...string) (*CommandResult, error) {
	cmd := exec.Command(command, args...)

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Collect output
	stdoutStr, err := readOutput(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdout: %w", err)
	}

	stderrStr, err := readOutput(stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to read stderr: %w", err)
	}

	// Wait for completion
	err = cmd.Wait()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to wait for command: %w", err)
		}
	}

	return &CommandResult{
		ExitCode: exitCode,
		Stdout:   stdoutStr,
		Stderr:   stderrStr,
	}, nil
}

// ExecuteInteractiveCommand runs a command that requires user interaction
func ExecuteInteractiveCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ExecuteWithProgress runs a command and reports progress
func ExecuteWithProgress(command string, args ...string) (*CommandResult, error) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Create channels for output collection
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	// Start goroutines to collect output
	go collectOutput(stdout, stdoutChan)
	go collectOutput(stderr, stderrChan)

	// Collect all output
	var stdoutBuilder, stderrBuilder strings.Builder
	stdoutDone, stderrDone := false, false

	for {
		if stdoutDone && stderrDone {
			break
		}

		select {
		case line, ok := <-stdoutChan:
			if !ok {
				stdoutDone = true
				continue
			}
			stdoutBuilder.WriteString(line + "\n")
			log.Debug().Str("stdout", line).Msg("Command output")

		case line, ok := <-stderrChan:
			if !ok {
				stderrDone = true
				continue
			}
			stderrBuilder.WriteString(line + "\n")
			log.Debug().Str("stderr", line).Msg("Command error output")
		}
	}

	// Wait for completion
	err = cmd.Wait()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to wait for command: %w", err)
		}
	}

	return &CommandResult{
		ExitCode: exitCode,
		Stdout:   stdoutBuilder.String(),
		Stderr:   stderrBuilder.String(),
	}, nil
}

// Helper function to read command output
func readOutput(r io.Reader) (string, error) {
	var builder strings.Builder
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		builder.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return builder.String(), nil
}

// Helper function to collect command output
func collectOutput(r io.Reader, ch chan<- string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	close(ch)
}

// CheckCommandExists verifies if a command exists in the system
func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
