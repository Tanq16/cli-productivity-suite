package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tanq16/cli-productivity-suite/utils"
)

func ExtractTarGz(archivePath, destDir string) error {
	cmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir)
	return utils.RunCmd(cmd)
}

func ExtractTarXz(archivePath, destDir string) error {
	cmd := exec.Command("tar", "-xJf", archivePath, "-C", destDir)
	return utils.RunCmd(cmd)
}

func ExtractZip(archivePath, destDir string) error {
	cmd := exec.Command("unzip", "-o", "-q", archivePath, "-d", destDir)
	return utils.RunCmd(cmd)
}

func FindBinary(dir, pattern string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return "", err
	}
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			return m, nil
		}
	}
	return "", fmt.Errorf("binary not found matching pattern %s in %s", pattern, dir)
}
