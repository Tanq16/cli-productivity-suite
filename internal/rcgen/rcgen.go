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

var runtimePackages = []string{"go-sdk", "java-sdk", "rust", "fnm", "node", "bun", "uv", "python"}

func path(p platform.Platform) string {
	return filepath.Join(p.ShellDir(), "rc", "cps.zsh")
}

func Generate(p platform.Platform, st *state.State) error {
	snippets := []struct {
		when    bool
		content []byte
	}{
		{true, configs.RcBase()},
		{st.Installed("anbu"), []byte(anbuAlias)},
		{slices.ContainsFunc(runtimePackages, st.Installed), configs.RcRuntimes()},
		{st.Installed("aws-cli") || st.Installed("gcloud-cli"), configs.RcCloud()},
		{st.Installed("nuclei-templates"), []byte(nucleiTemplates)},
		{st.Installed("neo4j"), []byte(neo4jConf)},
		{true, []byte(syntaxHighlighting)},
	}

	var b strings.Builder
	for _, s := range snippets {
		if !s.when {
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
