package configs

import _ "embed"

//go:embed tmux.conf
var tmuxConf []byte

//go:embed linux.kittyconf
var linuxKittyConf []byte

//go:embed macos.kittyconf
var macosKittyConf []byte

//go:embed mocha.kittyconf
var mochaKittyConf []byte

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

func TmuxConf() []byte           { return tmuxConf }
func LinuxKittyConf() []byte     { return linuxKittyConf }
func MacosKittyConf() []byte     { return macosKittyConf }
func MochaKittyConf() []byte     { return mochaKittyConf }
func MacosAerospaceConf() []byte { return macosAerospaceConf }
func RcLoader() []byte           { return rcLoader }
func RcBase() []byte             { return rcBase }
func RcRuntimes() []byte         { return rcRuntimes }
func RcCloud() []byte            { return rcCloud }
