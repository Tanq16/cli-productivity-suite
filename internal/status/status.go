package status

import (
	"errors"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type Entry struct {
	Name       string `json:"name"`
	Group      string `json:"group"`
	Suite      string `json:"suite,omitempty"`
	Kind       string `json:"kind"`
	Installed  bool   `json:"installed"`
	Version    string `json:"version,omitempty"`
	Latest     string `json:"latest,omitempty"`
	Outdated   bool   `json:"outdated"`
	SkipReason string `json:"skip_reason,omitempty"`
	CheckError string `json:"check_error,omitempty"`
}

func Collect(p platform.Platform, st *state.State) []Entry {
	tools := registry.All()
	entries := make([]Entry, 0, len(tools))
	for _, t := range tools {
		if !t.SupportsPlatform(p.OS.String()) {
			continue
		}
		version := st.ToolVersion(t.Name)
		entries = append(entries, Entry{
			Name:      t.Name,
			Group:     t.Group,
			Suite:     t.Suite,
			Kind:      t.Kind.String(),
			Installed: version != "",
			Version:   version,
		})
	}
	return entries
}

func CheckEntry(e *Entry, p platform.Platform, gh *github.Client) {
	if !e.Installed {
		return
	}
	tool, ok := registry.ByName(e.Name)
	if !ok {
		e.SkipReason = "not in the registry"
		return
	}
	if tool.IsPrivate && !gh.HasToken() {
		e.SkipReason = "needs --gh-token"
		return
	}
	latest, err := installer.LatestVersion(&tool, p, gh)
	if errors.Is(err, installer.ErrNotCheckable) {
		e.SkipReason = "no upstream version to compare"
		return
	}
	if err != nil {
		e.CheckError = err.Error()
		return
	}
	e.Latest = latest
	e.Outdated = latest != e.Version
}

func GroupOrder() []string {
	groups := registry.Groups()
	order := make([]string, 0, len(groups)+1)
	order = append(order, registry.GroupShell)
	for _, g := range groups {
		order = append(order, g.Name)
	}
	return order
}

func ByGroup(entries []Entry, group string) []Entry {
	var result []Entry
	for _, e := range entries {
		if e.Group == group {
			result = append(result, e)
		}
	}
	return result
}
