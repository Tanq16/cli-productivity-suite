package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tanq16/cli-productivity-suite/utils"
)

func httpGet(url string) (*http.Response, error) {
	return utils.HTTPClient.Get(url)
}

func DownloadToFile(url, destPath string) error {
	resp, err := httpGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d from %s", resp.StatusCode, url)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
