package registry

import "slices"

func All() []Tool {
	all := make([]Tool, 0, len(shellTools)+len(packages))
	all = append(all, shellTools...)
	all = append(all, packages...)
	return all
}

func ShellTools() []Tool {
	return slices.Clone(shellTools)
}

func Groups() []Group {
	return slices.Clone(groups)
}

func Suites() []Suite {
	return slices.Clone(suites)
}

func ByName(name string) (Tool, bool) {
	for _, set := range [][]Tool{shellTools, packages} {
		if i := slices.IndexFunc(set, func(t Tool) bool { return t.Name == name }); i >= 0 {
			return set[i], true
		}
	}
	return Tool{}, false
}

func ByGroup(group string) []Tool {
	if group == GroupShell {
		return ShellTools()
	}
	var result []Tool
	for _, t := range packages {
		if t.Group == group {
			result = append(result, t)
		}
	}
	return result
}

func BySuite(suite string) []Tool {
	var result []Tool
	for _, t := range packages {
		if t.Suite == suite {
			result = append(result, t)
		}
	}
	return result
}

func Resolve(name string) ([]Tool, bool) {
	if i := slices.IndexFunc(packages, func(t Tool) bool { return t.Name == name }); i >= 0 {
		return []Tool{packages[i]}, true
	}
	if slices.ContainsFunc(groups, func(g Group) bool { return g.Name == name }) {
		return ByGroup(name), true
	}
	if slices.ContainsFunc(suites, func(s Suite) bool { return s.Name == name }) {
		return BySuite(name), true
	}
	return nil, false
}

func InstallTargets() []string {
	names := make([]string, 0, len(groups)+len(suites)+len(packages))
	for _, g := range groups {
		names = append(names, g.Name)
	}
	for _, s := range suites {
		names = append(names, s.Name)
	}
	for _, t := range packages {
		names = append(names, t.Name)
	}
	return names
}

func Order(tools []Tool, installed func(name string) bool) []Tool {
	requested := make(map[string]bool, len(tools))
	for _, t := range tools {
		requested[t.Name] = true
	}

	ordered := make([]Tool, 0, len(tools))
	seen := make(map[string]bool, len(tools))

	var visit func(t Tool)
	visit = func(t Tool) {
		if seen[t.Name] {
			return
		}
		seen[t.Name] = true
		for _, name := range t.Requires {
			if !requested[name] && installed(name) {
				continue
			}
			if dep, ok := ByName(name); ok {
				visit(dep)
			}
		}
		ordered = append(ordered, t)
	}

	for _, t := range tools {
		visit(t)
	}
	return ordered
}
