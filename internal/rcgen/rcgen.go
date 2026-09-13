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

func path(p platform.Platform) string {
	return filepath.Join(p.ShellDir(), "rc", "cps.zsh")
}

func runtimeBlock(st *state.State) []byte {
	var b strings.Builder
	for _, r := range platform.RuntimeEnvs() {
		if !slices.ContainsFunc(r.Pkgs, st.Installed) {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		for _, v := range r.Vars {
			b.WriteString("export " + v[0] + `="$HOME/` + v[1] + "\"\n")
		}
		if len(r.Paths) > 0 {
			quoted := make([]string, len(r.Paths))
			for i, rel := range r.Paths {
				quoted[i] = "$HOME/" + rel
			}
			b.WriteString(`export PATH="` + strings.Join(quoted, ":") + ":$PATH\"\n")
		}
		if r.Shell != "" {
			b.WriteString(r.Shell + "\n")
		}
	}
	return []byte(b.String())
}

func Generate(p platform.Platform, st *state.State) error {
	snippets := []struct {
		when    bool
		content []byte
	}{
		{true, configs.RcBase()},
		{st.Installed("anbu"), []byte(anbuAlias)},
		{true, runtimeBlock(st)},
		{st.Installed("aws-cli") || st.Installed("gcloud-cli"), configs.RcCloud()},
		{st.Installed("nuclei-templates"), []byte(nucleiTemplates)},
		{st.Installed("neo4j"), []byte(neo4jConf)},
		{true, []byte(syntaxHighlighting)},
	}

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

	dest := path(p)
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	return os.WriteFile(dest, []byte(b.String()), 0644)
}
