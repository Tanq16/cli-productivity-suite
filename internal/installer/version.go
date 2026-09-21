package installer

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
)

var ErrNotCheckable = errors.New("version is reported by its package manager")

func LatestVersion(tool *registry.Tool, p platform.Platform, gh *github.Client) (string, error) {
	switch tool.Kind {
	case registry.GitHubRelease:
		return ghTag(gh, tool.Repo)
	case registry.DirectDownload, registry.AppBundle:
		return resolveVersion(tool, gh)
	case registry.RepoSnapshot:
		return snapshotHead(tool.Repo, gh)
	case registry.LanguageRuntime:
		return runtimeLatest(tool.Name, p, gh)
	case registry.PythonTool:
		return pypiLatest(tool.PyTool)
	default:
		return "", ErrNotCheckable
	}
}

func pypiLatest(pkg string) (string, error) {
	resp, err := httpGet("https://pypi.org/pypi/" + pkg + "/json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("PyPI returned HTTP %d for %s", resp.StatusCode, pkg)
	}
	var project struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	if err := json.UnmarshalRead(resp.Body, &project); err != nil {
		return "", err
	}
	if project.Info.Version == "" {
		return "", fmt.Errorf("no version in the PyPI response for %s", pkg)
	}
	return project.Info.Version, nil
}

func ghTag(gh *github.Client, repo string) (string, error) {
	release, err := gh.LatestRelease(repo)
	if err != nil {
		return "", err
	}
	return release.TagName, nil
}

func snapshotHead(repo string, gh *github.Client) (string, error) {
	resp, err := gh.Get(fmt.Sprintf("https://api.github.com/repos/%s/commits/HEAD", repo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error %d for %s: %s", resp.StatusCode, repo, string(body))
	}
	var commit struct {
		SHA string `json:"sha"`
	}
	if err := json.UnmarshalRead(resp.Body, &commit); err != nil {
		return "", err
	}
	if len(commit.SHA) < 7 {
		return "", fmt.Errorf("no commit SHA returned for %s", repo)
	}
	return commit.SHA[:7], nil
}

func runtimeLatest(name string, p platform.Platform, gh *github.Client) (string, error) {
	switch name {
	case "go-sdk":
		version, _, err := goLatest(p)
		return version, err
	case "java-sdk":
		version, _, err := javaLatest(p)
		return version, err
	case "rust":
		return rustLatest()
	case "uv":
		return ghTag(gh, uvRepo)
	case "fnm":
		return ghTag(gh, fnmRepo)
	case "node":
		tag, err := ghTag(gh, fnmRepo)
		if err != nil {
			return "", err
		}
		lts, err := nodeLTS()
		if err != nil {
			return "", err
		}
		return tag + "/" + lts, nil
	case "python":
		tag, err := ghTag(gh, uvRepo)
		if err != nil {
			return "", err
		}
		cycle, err := latestPythonCycle()
		if err != nil {
			return "", err
		}
		return tag + "/" + cycle, nil
	default:
		return "", ErrNotCheckable
	}
}

var rustStableVersion = regexp.MustCompile(`\[pkg\.rust\]\s*\nversion = "([^ "]+)`)

func rustLatest() (string, error) {
	resp, err := httpGet("https://static.rust-lang.org/dist/channel-rust-stable.toml")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("rust channel manifest returned HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	m := rustStableVersion.FindSubmatch(body)
	if len(m) < 2 {
		return "", errors.New("no stable version in the rust channel manifest")
	}
	return string(m[1]), nil
}

func nodeLTS() (string, error) {
	resp, err := httpGet("https://nodejs.org/dist/index.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("nodejs.org release index returned HTTP %d", resp.StatusCode)
	}
	var releases []struct {
		Version string `json:"version"`
		LTS     any    `json:"lts"`
	}
	if err := json.UnmarshalRead(resp.Body, &releases); err != nil {
		return "", err
	}
	for _, r := range releases {
		if codename, ok := r.LTS.(string); ok && codename != "" {
			return r.Version, nil
		}
	}
	return "", errors.New("no Node LTS release found")
}
