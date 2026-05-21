package runner

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/utils"
)

const knownExtensionsBaseURL = "https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/main/custom-extensions"

var knownExtensions = []string{
	"ai-tools",
	"additional-cloud-tools",
	"database",
	"praetorian",
}

func DownloadKnownExtensions() {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	extDir := filepath.Join(p.ConfigDir(), "extensions")
	if err := os.MkdirAll(extDir, 0755); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to create %s", extDir), err)
	}

	var hadErrors bool

	for _, name := range knownExtensions {
		url := fmt.Sprintf("%s/%s.yaml", knownExtensionsBaseURL, name)
		dest := filepath.Join(extDir, name+".yaml")

		utils.PrintRunning("fetching " + name + ".yaml")
		err := downloadKnownYAML(utils.HTTPClient, url, dest)
		utils.ClearLines(1)
		if err != nil {
			utils.PrintError(name+".yaml: download failed", err)
			hadErrors = true
			continue
		}
		utils.PrintSuccess(name + ".yaml")
	}

	if hadErrors {
		utils.PrintWarn(fmt.Sprintf("some downloads failed — packs are written to %s", extDir), nil)
	} else {
		utils.PrintSuccess(fmt.Sprintf("downloaded %d reference packs to %s", len(knownExtensions), extDir))
	}
}

func downloadKnownYAML(client *http.Client, url, dest string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, body, 0644)
}
