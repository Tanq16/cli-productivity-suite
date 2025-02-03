package system

import (
	"runtime"
)

type SystemInfo struct {
	OS   string
	Arch string
}

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
