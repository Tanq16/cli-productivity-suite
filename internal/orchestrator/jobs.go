package orchestrator

import (
	"context"
	"fmt"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/highway"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type InstallJob struct {
	tool registry.Tool
	p    platform.Platform
	gh   *github.Client
	st   *state.State
}

func NewInstallJob(tool registry.Tool, p platform.Platform, gh *github.Client, st *state.State) *InstallJob {
	return &InstallJob{tool: tool, p: p, gh: gh, st: st}
}

func (j *InstallJob) ID() string { return j.tool.Name }

func (j *InstallJob) Run(ctx context.Context, progress chan<- highway.Progress) error {
	progress <- highway.Progress{
		JobID:     j.tool.Name,
		SubStatus: "installing",
		Type:      highway.ProgressTypeSubStatus,
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	inst := installer.Dispatch(j.tool.Kind)
	if inst == nil {
		progress <- highway.Progress{
			JobID:   j.tool.Name,
			Message: fmt.Sprintf("no installer for kind %s", j.tool.Kind),
			Done:    true,
			Error:   fmt.Errorf("no installer for tool kind: %s", j.tool.Kind),
		}
		return nil
	}
	result := inst.Install(&j.tool, j.p, j.gh, j.st)

	var msg string
	if result.Err != nil {
		msg = result.Err.Error()
	} else if result.Skipped {
		msg = fmt.Sprintf("already at %s", result.Version)
	} else if result.WasUpdated {
		msg = fmt.Sprintf("updated to %s", result.Version)
	} else {
		msg = fmt.Sprintf("installed %s", result.Version)
	}

	progress <- highway.Progress{
		JobID:   j.tool.Name,
		Message: msg,
		Done:    true,
		Error:   result.Err,
	}
	return nil
}
