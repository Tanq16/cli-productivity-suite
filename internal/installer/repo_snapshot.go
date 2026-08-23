package installer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type RepoSnapshotInstaller struct{}

func (r *RepoSnapshotInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	if tool.Repo == "" || tool.Dest == "" {
		return Result{Tool: tool.Name, Err: errors.New("repo snapshot needs both Repo and Dest")}
	}

	// Staged inside ~/shell so the swap into place is a same-filesystem rename; /tmp is a separate mount on most Linux hosts.
	tmpDir, err := os.MkdirTemp(p.ShellDir(), "cps-"+tool.Name+"-*")
	if err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	defer os.RemoveAll(tmpDir)

	// GitHub resolves HEAD to the repo's own default branch, so no branch name is pinned here.
	url := fmt.Sprintf("https://github.com/%s/archive/HEAD.tar.gz", tool.Repo)
	archivePath := filepath.Join(tmpDir, "snapshot.tar.gz")
	if err := DownloadToFile(url, archivePath); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("download failed: %w", err)}
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	if err := ExtractArchive(archivePath, extractDir, "tar.gz"); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("extract failed: %w", err)}
	}

	srcDir := unwrapSingleDir(extractDir)
	version := snapshotVersion(srcDir)

	dest := expandHome(tool.Dest, p.HomeDir)
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}
	if err := stageAndSwap(srcDir, dest); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	currentVersion := st.ToolVersion(tool.Name)
	st.SetToolVersion(tool.Name, version)
	return Result{Tool: tool.Name, Version: version, WasUpdated: currentVersion != "" && currentVersion != version}
}

// A GitHub branch archive roots at <repo>-<full sha>, which is the only place the snapshot's commit is recorded.
func snapshotVersion(srcDir string) string {
	base := filepath.Base(srcDir)
	i := strings.LastIndex(base, "-")
	if i < 0 {
		return "snapshot"
	}
	sha := base[i+1:]
	if len(sha) < 7 {
		return "snapshot"
	}
	return sha[:7]
}

func expandHome(path, homeDir string) string {
	if after, ok := strings.CutPrefix(path, "~/"); ok {
		return filepath.Join(homeDir, after)
	}
	return path
}
