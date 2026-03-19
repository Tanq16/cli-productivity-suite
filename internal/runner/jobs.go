package runner

import "github.com/tanq16/cli-productivity-suite/internal/registry"

type jobResult struct {
	name string
	err  error
}

type CheckResult struct {
	Tool    registry.Tool
	Current string
	Latest  string
	Status  string // "update", "error", "config-differs", "not-deployed"
	Err     error
}
