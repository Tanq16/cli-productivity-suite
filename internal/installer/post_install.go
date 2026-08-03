package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/configs"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func runPostInstall(tool *registry.Tool, p platform.Platform, destDir string) error {
	switch tool.PostInstall {
	case "":
		return nil
	case "neo4j":
		return seedNeo4jConfig(p, destDir)
	case "code-server":
		return seedCodeServerConfig(p, destDir)
	default:
		return fmt.Errorf("unknown post-install hook %q for %s", tool.PostInstall, tool.Name)
	}
}

// Neo4j defaults these under its install root, which an upgrade replaces wholesale; they are the directory settings that accept an absolute path.
var neo4jRelocatedDirs = []string{"data", "plugins", "import", "logs", "run", "licenses"}

func seedNeo4jConfig(p platform.Platform, bundleDir string) error {
	stateRoot := filepath.Join(p.HomeDir, ".config", "neo4j")
	for _, d := range neo4jRelocatedDirs {
		if err := os.MkdirAll(filepath.Join(stateRoot, d), 0755); err != nil {
			return err
		}
	}
	// Certificates have no server.directories setting, so the directory is only a convention for the TLS hint written below.
	certDir := filepath.Join(stateRoot, "certificates")
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return err
	}

	confDir := filepath.Join(stateRoot, "conf")
	confPath := filepath.Join(confDir, "neo4j.conf")
	if _, err := os.Stat(confPath); err == nil {
		return nil
	}
	if err := os.CopyFS(confDir, os.DirFS(filepath.Join(bundleDir, "conf"))); err != nil {
		return fmt.Errorf("seed neo4j conf: %w", err)
	}

	body, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	conf := string(body)
	var relocations strings.Builder
	// A duplicate key is fatal to Neo4j and server.directories.import ships uncommented, so existing keys are disabled rather than overridden.
	for _, d := range neo4jRelocatedDirs {
		re := regexp.MustCompile(`(?m)^server\.directories\.` + regexp.QuoteMeta(d) + `=`)
		conf = re.ReplaceAllString(conf, "#server.directories."+d+"=")
		fmt.Fprintf(&relocations, "server.directories.%s=%s\n", d, filepath.Join(stateRoot, d))
	}

	var out strings.Builder
	out.WriteString(conf)
	if !strings.HasSuffix(conf, "\n") {
		out.WriteString("\n")
	}
	out.WriteString("\n# Managed by cps: state lives here so upgrades can replace ~/shell/apps/neo4j wholesale.\n")
	out.WriteString(relocations.String())
	out.WriteString("# For TLS, point dbms.ssl.policy.<scope>.base_directory at " + certDir + "\n")

	return os.WriteFile(confPath, []byte(out.String()), 0644)
}

var codeServerExtensions = []string{
	"Catppuccin.catppuccin-vsc",
	"Catppuccin.catppuccin-vsc-icons",
	"EditorConfig.EditorConfig",
	"usernamehw.errorlens",
	"streetsidesoftware.code-spell-checker",
	"redhat.vscode-yaml",
	"golang.Go",
	"ms-python.python",
	"charliermarsh.ruff",
}

func seedCodeServerConfig(p platform.Platform, bundleDir string) error {
	userDataDir := filepath.Join(p.HomeDir, ".local", "share", "code-server")
	extensionsDir := filepath.Join(userDataDir, "extensions")

	// code-server writes its own config with a generated password on any invocation, including --version, so this must land before the binary first runs.
	if err := writeIfAbsent(filepath.Join(p.HomeDir, ".config", "code-server", "config.yaml"), configs.CodeServerConfig()); err != nil {
		return err
	}
	if err := writeIfAbsent(filepath.Join(userDataDir, "User", "settings.json"), configs.CodeServerSettings()); err != nil {
		return err
	}

	// An empty extensions dir is also the retry signal: a failed install leaves it empty and the next run tries again.
	if entries, err := os.ReadDir(extensionsDir); err == nil && len(entries) > 0 {
		return nil
	}
	args := []string{"--user-data-dir", userDataDir, "--extensions-dir", extensionsDir}
	for _, ext := range codeServerExtensions {
		args = append(args, "--install-extension", ext)
	}
	cmd := exec.Command(filepath.Join(bundleDir, "bin", "code-server"), args...)
	if err := utils.RunCmd(cmd); err != nil {
		return fmt.Errorf("install code-server extensions: %w", err)
	}
	return nil
}

func writeIfAbsent(path string, content []byte) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}
