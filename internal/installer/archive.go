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

func ExtractArchive(archivePath, destDir, format string) error {
	switch format {
	case "tar.gz", "tgz":
		return ExtractTarGz(archivePath, destDir)
	case "tar.xz":
		return ExtractTarXz(archivePath, destDir)
	case "zip":
		return ExtractZip(archivePath, destDir)
	default:
		return fmt.Errorf("unknown archive format: %s", format)
	}
}

// Bundle archive roots are named per tool and version, so unwrap to keep the install path stable.
func unwrapSingleDir(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 || !entries[0].IsDir() {
		return dir
	}
	return filepath.Join(dir, entries[0].Name())
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
	return stageAndSwapPreserving(srcDir, destDir, nil)
}

// Some bundles keep user state inside the install root, so those paths move across from the outgoing install before it is discarded; on failure the old tree is left at destDir.old.
func stageAndSwapPreserving(srcDir, destDir string, preserve []string) error {
	oldDir := destDir + ".old"
	os.RemoveAll(oldDir)
	hadOld := false
	if _, err := os.Stat(destDir); err == nil {
		if err := os.Rename(destDir, oldDir); err != nil {
			return fmt.Errorf("backup existing %s: %w", destDir, err)
		}
		hadOld = true
	}
	if err := os.Rename(srcDir, destDir); err != nil {
		if hadOld {
			os.Rename(oldDir, destDir)
		}
		return fmt.Errorf("install %s: %w", destDir, err)
	}
	if hadOld {
		for _, rel := range preserve {
			src := filepath.Join(oldDir, rel)
			if _, err := os.Stat(src); err != nil {
				continue
			}
			dst := filepath.Join(destDir, rel)
			if err := os.RemoveAll(dst); err != nil {
				return fmt.Errorf("preserve %s: %w", rel, err)
			}
			if err := os.Rename(src, dst); err != nil {
				return fmt.Errorf("preserve %s: %w", rel, err)
			}
		}
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
