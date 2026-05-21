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

func stageAndSwap(srcDir, destDir string) error {
	oldDir := destDir + ".old"
	os.RemoveAll(oldDir)
	if _, err := os.Stat(destDir); err == nil {
		if err := os.Rename(destDir, oldDir); err != nil {
			return fmt.Errorf("backup existing %s: %w", destDir, err)
		}
	}
	if err := os.Rename(srcDir, destDir); err != nil {
		if _, statErr := os.Stat(oldDir); statErr == nil {
			os.Rename(oldDir, destDir)
		}
		return fmt.Errorf("install %s: %w", destDir, err)
	}
	os.RemoveAll(oldDir)
	return nil
}

func AtomicInstallBinary(srcPath, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	dir := filepath.Dir(destPath)
	tmp, err := os.CreateTemp(dir, ".cps-tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, destPath)
}
