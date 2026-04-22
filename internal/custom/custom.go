package custom

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"gopkg.in/yaml.v3"
)

type PackFile struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Shell       ShellBlock `yaml:"shell"`
	Tools       []ToolDef  `yaml:"tools"`
}

type ShellBlock struct {
	Env         map[string]string `yaml:"env"`
	PathPrepend []string          `yaml:"path_prepend"`
	Source      []string          `yaml:"source"`
}

type ToolDef struct {
	Name    string `yaml:"name"`
	Install string `yaml:"install"`
	Remove  string `yaml:"remove"`
}

var (
	validName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
	cache     = map[string]*PackFile{}
)

func LoadDir(dir string, builtinNames map[string]bool) ([]registry.ExtensionPack, []string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	var packs []registry.ExtensionPack
	var warnings []string
	seen := map[string]bool{}

	for _, e := range entries {
		if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
			continue
		}

		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("skipping %s: %v", e.Name(), err))
			continue
		}

		var pf PackFile
		if err := yaml.Unmarshal(data, &pf); err != nil {
			warnings = append(warnings, fmt.Sprintf("skipping %s: invalid YAML: %v", e.Name(), err))
			continue
		}

		if warn := validate(pf, e.Name(), builtinNames, seen); warn != "" {
			warnings = append(warnings, warn)
			continue
		}

		seen[pf.Name] = true
		cache[pf.Name] = &pf
		packs = append(packs, toExtensionPack(pf))
	}

	return packs, warnings
}

func GetPackFile(name string) *PackFile {
	return cache[name]
}

func validate(pf PackFile, filename string, builtinNames, seen map[string]bool) string {
	if pf.Name == "" {
		return fmt.Sprintf("skipping %s: missing 'name' field", filename)
	}
	if !validName.MatchString(pf.Name) {
		return fmt.Sprintf("skipping %s: name %q must match [a-z0-9-]", filename, pf.Name)
	}
	if builtinNames[pf.Name] {
		return fmt.Sprintf("skipping %s: name %q conflicts with built-in pack", filename, pf.Name)
	}
	if seen[pf.Name] {
		return fmt.Sprintf("skipping %s: duplicate pack name %q", filename, pf.Name)
	}
	if len(pf.Tools) == 0 {
		return fmt.Sprintf("skipping %s: no tools defined", filename)
	}

	toolNames := map[string]bool{}
	for i, t := range pf.Tools {
		if t.Name == "" {
			return fmt.Sprintf("skipping %s: tool #%d has no name", filename, i+1)
		}
		if t.Install == "" {
			return fmt.Sprintf("skipping %s: tool %q has no install command", filename, t.Name)
		}
		if toolNames[t.Name] {
			return fmt.Sprintf("skipping %s: duplicate tool name %q", filename, t.Name)
		}
		toolNames[t.Name] = true
	}

	return ""
}

func toExtensionPack(pf PackFile) registry.ExtensionPack {
	desc := pf.Description
	if desc == "" {
		desc = "Custom extension pack"
	}

	tools := make([]registry.Tool, len(pf.Tools))
	for i, td := range pf.Tools {
		tools[i] = registry.Tool{
			Name:        td.Name,
			Kind:        registry.CustomScript,
			Category:    registry.ExtCustom,
			Extension:   true,
			InstallCmd:  td.Install,
			RemoveCmd:   td.Remove,
			Description: "custom tool",
		}
	}

	return registry.ExtensionPack{
		Name:        pf.Name,
		Description: desc,
		Category:    registry.ExtCustom,
		Tools:       tools,
	}
}
