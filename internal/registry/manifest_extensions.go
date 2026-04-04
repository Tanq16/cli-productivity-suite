package registry

// ExtensionPack defines a named collection of extension tools.
type ExtensionPack struct {
	Name        string
	Description string
	Category    ToolCategory
	Tools       []Tool
}

var extensionPacks = []ExtensionPack{
	{
		Name:        "cloud-sec",
		Description: "Cloud security tools",
		Category:    ExtCloudSec,
		Tools:       []Tool{
			// Add cloud security tools here (GitHubRelease / DirectDownload only)
		},
	},
	{
		Name:        "app-sec",
		Description: "Application security tools",
		Category:    ExtAppSec,
		Tools:       []Tool{
			// Add application security tools here (GitHubRelease / DirectDownload only)
		},
	},
}

// AllExtensionPacks returns the list of available extension packs.
func AllExtensionPacks() []ExtensionPack {
	return extensionPacks
}

// ExtensionPackByName returns the pack with the given name, or nil.
func ExtensionPackByName(name string) *ExtensionPack {
	for i := range extensionPacks {
		if extensionPacks[i].Name == name {
			return &extensionPacks[i]
		}
	}
	return nil
}

// AllExtensionTools returns every tool across all extension packs.
func AllExtensionTools() []Tool {
	var all []Tool
	for _, pack := range extensionPacks {
		all = append(all, pack.Tools...)
	}
	return all
}
