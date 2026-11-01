package rcgen

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/configs"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

const anbuAlias = "alias a=anbu\n"

const nucleiTemplates = `export NUCLEI_TEMPLATES_DIR="$HOME/shell/nuclei-templates"
`

const neo4jConf = `export NEO4J_CONF="$HOME/.config/neo4j/conf"
`

const syntaxHighlighting = `[ -f "$ZSH_PLUGINS/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" ] && source "$ZSH_PLUGINS/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
`

// Interactive non-login shells (Linux terminals, tmux panes) never read ~/.zprofile
const envSource = `[ -f "$HOME/shell/rc/cps-env.zsh" ] && source "$HOME/shell/rc/cps-env.zsh"
`

func rcPath(p platform.Platform) string {
	return filepath.Join(p.ShellDir(), "rc", "cps.zsh")
}

func envPath(p platform.Platform) string {
	return filepath.Join(p.ShellDir(), "rc", "cps-env.zsh")
}

func installedRuntimes(st *state.State) []platform.RuntimeEnv {
	var out []platform.RuntimeEnv
	for _, r := range platform.RuntimeEnvs() {
		if slices.ContainsFunc(r.Pkgs, st.Installed) {
			out = append(out, r)
		}
	}
	return out
}

// runtimeEnvBlock holds the exports and PATH entries, which non-interactive shells need too
func runtimeEnvBlock(st *state.State) []byte {
	var b strings.Builder
	for _, r := range installedRuntimes(st) {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		for _, v := range r.Vars {
			b.WriteString("export " + v[0] + `="$HOME/` + v[1] + "\"\n")
		}
		if len(r.Paths) > 0 {
			b.WriteString(pathExport(r.Paths))
		}
	}
	return []byte(b.String())
}

// runtimeShellBlock holds completions and fnm env, which only an interactive shell wants.
// Each runtime's PATH entries are prepended again beside its snippet so the interactive
// precedence matches the old single-file order; typeset -U drops the duplicates.
func runtimeShellBlock(st *state.State) []byte {
	var b strings.Builder
	for _, r := range installedRuntimes(st) {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		if len(r.Paths) > 0 {
			b.WriteString(pathExport(r.Paths))
		}
		if r.Shell != "" {
			b.WriteString(r.Shell + "\n")
		}
	}
	return []byte(b.String())
}

func pathExport(paths []string) string {
	quoted := make([]string, len(paths))
	for i, rel := range paths {
		quoted[i] = "$HOME/" + rel
	}
	return `export PATH="` + strings.Join(quoted, ":") + ":$PATH\"\n"
}

type snippet struct {
	when    bool
	content []byte
}

func join(snippets []snippet) []byte {
	var b strings.Builder
	for _, s := range snippets {
		if !s.when || len(s.content) == 0 {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.Write(s.content)
	}
	return []byte(b.String())
}

func write(dest string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	return os.WriteFile(dest, content, 0644)
}

func Generate(p platform.Platform, st *state.State) error {
	env := join([]snippet{
		{true, configs.RcEnv()},
		{true, runtimeEnvBlock(st)},
		{st.Installed("nuclei-templates"), []byte(nucleiTemplates)},
		{st.Installed("neo4j"), []byte(neo4jConf)},
	})
	if err := write(envPath(p), env); err != nil {
		return err
	}

	rc := join([]snippet{
		{true, []byte(envSource)},
		{true, configs.RcBase()},
		{st.Installed("anbu"), []byte(anbuAlias)},
		{true, runtimeShellBlock(st)},
		{st.Installed("aws-cli") || st.Installed("gcloud-cli"), configs.RcCloud()},
		{true, []byte(syntaxHighlighting)},
	})
	return write(rcPath(p), rc)
}
