package registry

type Registry struct {
	tools []Tool
}

func New() *Registry {
	combined := make([]Tool, len(AllTools))
	copy(combined, AllTools)
	combined = append(combined, AllExtensionTools()...)
	return &Registry{tools: combined}
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
