package configs

import _ "embed"

//go:embed tmux.conf
var tmuxConf []byte

//go:embed linux.kittyconf
var linuxKittyConf []byte

//go:embed macos.kittyconf
var macosKittyConf []byte

//go:embed rcfile
var rcfile []byte

//go:embed macos.aerospaceconf
var macosAerospaceConf []byte

func TmuxConf() []byte           { return tmuxConf }
func LinuxKittyConf() []byte     { return linuxKittyConf }
func MacosKittyConf() []byte     { return macosKittyConf }
func Rcfile() []byte             { return rcfile }
func MacosAerospaceConf() []byte { return macosAerospaceConf }
