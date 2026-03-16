package registry

type Registry struct {
	tools   []Tool
	byName  map[string]*Tool
}

func New() *Registry {
	r := &Registry{
		tools:  AllTools,
		byName: make(map[string]*Tool),
	}
	for i := range r.tools {
		r.byName[r.tools[i].Name] = &r.tools[i]
	}
	return r
}

func (r *Registry) Get(name string) *Tool {
	return r.byName[name]
}

func (r *Registry) All() []Tool {
	return r.tools
}

func (r *Registry) ByKind(kind ToolKind) []Tool {
	var result []Tool
	for _, t := range r.tools {
		if t.Kind == kind {
			result = append(result, t)
		}
	}
	return result
}

func (r *Registry) ByCategory(cat ToolCategory) []Tool {
	var result []Tool
	for _, t := range r.tools {
		if t.Category == cat {
			result = append(result, t)
		}
	}
	return result
}

func (r *Registry) ForPlatform(os string) []Tool {
	var result []Tool
	for _, t := range r.tools {
		if len(t.Platforms) == 0 {
			result = append(result, t)
			continue
		}
		for _, p := range t.Platforms {
			if p == os {
				result = append(result, t)
				break
			}
		}
	}
	return result
}

func (r *Registry) Names() []string {
	names := make([]string, len(r.tools))
	for i, t := range r.tools {
		names[i] = t.Name
	}
	return names
}
