package installer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// exec.Command ignores cmd.Env when resolving the program, so resolve against the augmented PATH here.
func envCommand(env []string, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(lookPathIn(env, name), args...)
	cmd.Env = env
	return cmd
}

func lookPathIn(env []string, name string) string {
	if strings.ContainsRune(name, os.PathSeparator) {
		return name
	}
	var path string
	for _, kv := range env {
		if after, ok := strings.CutPrefix(kv, "PATH="); ok {
			path = after
		}
	}
	for dir := range strings.SplitSeq(path, string(os.PathListSeparator)) {
		if dir == "" {
			continue
		}
		candidate := filepath.Join(dir, name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
			return candidate
		}
	}
	return name
}
