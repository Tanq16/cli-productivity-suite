package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type OS int

const (
	Linux OS = iota
	Darwin
)

func (o OS) String() string {
	switch o {
	case Linux:
		return "linux"
	case Darwin:
		return "darwin"
	default:
		return "unknown"
	}
}

type Arch int

const (
	AMD64 Arch = iota
	ARM64
)

func (a Arch) String() string {
	switch a {
	case AMD64:
		return "amd64"
	case ARM64:
		return "arm64"
	default:
		return "unknown"
	}
}

type Platform struct {
	OS      OS
	Arch    Arch
	HomeDir string
}

func Detect() (Platform, error) {
	var p Platform

	switch runtime.GOOS {
	case "linux":
		p.OS = Linux
	case "darwin":
		p.OS = Darwin
	default:
		return p, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "amd64":
		p.Arch = AMD64
	case "arm64":
		p.Arch = ARM64
	default:
		return p, fmt.Errorf("unsupported arch: %s", runtime.GOARCH)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return p, fmt.Errorf("cannot determine home directory: %w", err)
	}
	p.HomeDir = home

	return p, nil
}

func (p Platform) ShellDir() string {
	return filepath.Join(p.HomeDir, "shell")
}

func (p Platform) ShellExtDir() string {
	return filepath.Join(p.HomeDir, "shell", "extensions")
}

func (p Platform) ShellAppsDir() string {
	return filepath.Join(p.HomeDir, "shell", "apps")
}

func (p Platform) ConfigDir() string {
	return filepath.Join(p.HomeDir, ".config", "cps")
}

func (p Platform) StatePath() string {
	return filepath.Join(p.ConfigDir(), "state.json")
}

var brewPaths = []string{
	"/opt/homebrew/bin/brew",
	"/usr/local/bin/brew",
	"/home/linuxbrew/.linuxbrew/bin/brew",
}

func BrewPath() string {
	if path, err := exec.LookPath("brew"); err == nil {
		return path
	}
	for _, path := range brewPaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}
	return ""
}

type RuntimeEnv struct {
	Pkgs  []string
	Vars  [][2]string
	Paths []string
	Shell string
}

var runtimeEnvs = []RuntimeEnv{
	{
		Pkgs:  []string{"go-sdk"},
		Vars:  [][2]string{{"GOROOT", "shell/go-sdk"}, {"GOPATH", "shell/go"}, {"GOCACHE", "shell/go/cache"}},
		Paths: []string{"shell/go-sdk/bin", "shell/go/bin"},
	},
	{
		Pkgs:  []string{"java-sdk"},
		Vars:  [][2]string{{"JAVA_HOME", "shell/java-sdk"}},
		Paths: []string{"shell/java-sdk/bin"},
	},
	{
		Pkgs:  []string{"rust"},
		Vars:  [][2]string{{"RUSTUP_HOME", "shell/rust/.rustup"}, {"CARGO_HOME", "shell/rust/.cargo"}},
		Paths: []string{"shell/rust/.cargo/bin"},
	},
	{
		Pkgs:  []string{"fnm", "node"},
		Vars:  [][2]string{{"FNM_DIR", "shell/fnm"}, {"npm_config_cache", "shell/npm-cache"}},
		Paths: []string{"shell/fnm/aliases/lts-latest/bin"},
		Shell: `[ -f "$HOME/shell/completions/fnm.zsh" ] && source "$HOME/shell/completions/fnm.zsh"
command -v fnm &>/dev/null && eval "$(fnm env)"`,
	},
	{
		Pkgs:  []string{"bun"},
		Vars:  [][2]string{{"BUN_INSTALL", "shell/bun"}},
		Paths: []string{"shell/bun/bin"},
	},
	{
		Pkgs: []string{"uv", "python", "ruff"},
		Vars: [][2]string{
			{"UV_TOOL_DIR", "shell/uv-tools"},
			{"UV_TOOL_BIN_DIR", "shell/uv-tool-executables"},
			{"UV_PYTHON_INSTALL_DIR", "shell/uv-python"},
		},
		Paths: []string{"shell/uv-tool-executables"},
		Shell: `[ -f "$HOME/shell/completions/uv.zsh" ] && source "$HOME/shell/completions/uv.zsh"`,
	},
	{
		Pkgs:  []string{"python"},
		Vars:  [][2]string{{"VIRTUAL_ENV", "shell/py-default"}},
		Paths: []string{"shell/py-default/bin"},
	},
}

func RuntimeEnvs() []RuntimeEnv {
	return runtimeEnvs
}

func (p Platform) CustomScriptEnv() []string {
	env := os.Environ()
	paths := []string{
		filepath.Join(p.HomeDir, "shell", "custom-bin"),
		p.ShellExtDir(),
	}
	for _, r := range runtimeEnvs {
		for _, v := range r.Vars {
			env = append(env, v[0]+"="+filepath.Join(p.HomeDir, filepath.FromSlash(v[1])))
		}
		for _, rel := range r.Paths {
			paths = append(paths, filepath.Join(p.HomeDir, filepath.FromSlash(rel)))
		}
	}

	prefix := strings.Join(paths, ":")
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = "PATH=" + prefix + ":" + kv[len("PATH="):]
			return env
		}
	}
	return append(env, "PATH="+prefix)
}
