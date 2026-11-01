package rcgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

func generate(t *testing.T, installed ...string) (env, rc string) {
	t.Helper()
	home := t.TempDir()
	st, err := state.Load(filepath.Join(home, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range installed {
		st.SetToolVersion(name, "test")
	}
	p := platform.Platform{OS: platform.Linux, HomeDir: home}
	if err := Generate(p, st); err != nil {
		t.Fatal(err)
	}
	read := func(name string) string {
		b, err := os.ReadFile(filepath.Join(home, "shell", "rc", name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	return read("cps-env.zsh"), read("cps.zsh")
}

func TestEnvFileHoldsOnlyEnvironment(t *testing.T) {
	env, _ := generate(t, "go-sdk", "fnm", "uv", "python", "neo4j")
	for _, want := range []string{"typeset -U path PATH", "brew.zsh", "shell/custom-bin", "GOPATH", "FNM_DIR", "UV_TOOL_DIR", "NEO4J_CONF"} {
		if !strings.Contains(env, want) {
			t.Errorf("cps-env.zsh missing %q", want)
		}
	}
	for _, banned := range []string{"bindkey", "alias ", "setopt", "compinit", "starship", "fnm env", "completions/"} {
		if strings.Contains(env, banned) {
			t.Errorf("cps-env.zsh contains interactive-only %q", banned)
		}
	}
}

func TestRcSourcesEnvFirstAndKeepsInteractiveParts(t *testing.T) {
	_, rc := generate(t, "fnm")
	if !strings.HasPrefix(rc, envSource) {
		t.Error("cps.zsh does not source cps-env.zsh first")
	}
	for _, want := range []string{"starship.zsh", "eval \"$(fnm env)\"", "zsh-syntax-highlighting"} {
		if !strings.Contains(rc, want) {
			t.Errorf("cps.zsh missing %q", want)
		}
	}
	if strings.Contains(rc, "FNM_DIR") {
		t.Error("cps.zsh duplicates runtime exports that belong in cps-env.zsh")
	}
}

// fnm env used to run before bun, uv and python prepended their bin dirs, so those win
// in interactive shells; the split must keep that order
func TestRcKeepsRuntimeOrderAroundFnmEnv(t *testing.T) {
	_, rc := generate(t, "fnm", "bun", "uv", "python")
	fnm := strings.Index(rc, "eval \"$(fnm env)\"")
	for _, later := range []string{"shell/bun/bin", "shell/uv-tool-executables", "shell/py-default/bin"} {
		if i := strings.Index(rc, later); i < fnm {
			t.Errorf("%s is prepended before fnm env (at %d, fnm at %d)", later, i, fnm)
		}
	}
}
