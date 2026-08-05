package configs

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

const DefaultTheme = "mocha"

//go:embed tmux.conf
var tmuxConf []byte

//go:embed linux.kittyconf
var linuxKittyConf []byte

//go:embed macos.kittyconf
var macosKittyConf []byte

//go:embed themes/*.kittyconf
var themeFS embed.FS

//go:embed macos.aerospaceconf
var macosAerospaceConf []byte

//go:embed rc-loader
var rcLoader []byte

//go:embed rc-base.zsh
var rcBase []byte

//go:embed rc-runtimes.zsh
var rcRuntimes []byte

//go:embed rc-cloud.zsh
var rcCloud []byte

//go:embed starship.toml
var starshipToml []byte

//go:embed code-server-config.yaml
var codeServerConfig []byte

//go:embed code-server-settings.json
var codeServerSettings []byte

//go:embed nvim-init.lua
var nvimInit []byte

func TmuxConf() []byte           { return tmuxConf }
func LinuxKittyConf() []byte     { return linuxKittyConf }
func MacosKittyConf() []byte     { return macosKittyConf }
func MacosAerospaceConf() []byte { return macosAerospaceConf }
func RcLoader() []byte           { return rcLoader }
func RcBase() []byte             { return rcBase }
func RcRuntimes() []byte         { return rcRuntimes }
func RcCloud() []byte            { return rcCloud }
func StarshipToml() []byte       { return starshipToml }
func CodeServerConfig() []byte   { return codeServerConfig }
func CodeServerSettings() []byte { return codeServerSettings }
func NvimInit() []byte           { return nvimInit }

func Theme(name string) ([]byte, error) {
	content, err := themeFS.ReadFile("themes/" + name + ".kittyconf")
	if err != nil {
		return nil, fmt.Errorf("unknown theme %q, pick one of: %s", name, strings.Join(ThemeNames(), ", "))
	}
	return content, nil
}

func ThemeNames() []string {
	entries, err := themeFS.ReadDir("themes")
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, strings.TrimSuffix(e.Name(), ".kittyconf"))
	}
	sort.Strings(names)
	return names
}
