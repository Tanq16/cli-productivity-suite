package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type CloudCLIInstaller struct{}

func (c *CloudCLIInstaller) Check(tool *registry.Tool, _ platform.Platform, _ *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)
	return current, "check-manually", nil
}

func (c *CloudCLIInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	switch tool.Name {
	case "aws-cli":
		return c.installAWSCLI(p, st)
	case "azure-cli":
		return c.installAzureCLI(p, st)
	case "gcloud-cli":
		return c.installGcloudCLI(p, st)
	default:
		return Result{Tool: tool.Name, Err: fmt.Errorf("unknown cloud CLI: %s", tool.Name)}
	}
}

func (c *CloudCLIInstaller) installAWSCLI(p platform.Platform, st *state.State) Result {
	if p.OS == platform.Darwin {
		cmd := exec.Command("brew", "install", "awscli")
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: "aws-cli", Err: fmt.Errorf("brew install awscli failed: %w", err)}
		}
	} else {
		var url string
		switch p.Arch {
		case platform.AMD64:
			url = "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip"
		case platform.ARM64:
			url = "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip"
		}

		tmpDir, err := os.MkdirTemp("", "cps-awscli-*")
		if err != nil {
			return Result{Tool: "aws-cli", Err: err}
		}
		defer os.RemoveAll(tmpDir)

		zipPath := filepath.Join(tmpDir, "awscliv2.zip")
		if err := DownloadToFile(url, zipPath); err != nil {
			return Result{Tool: "aws-cli", Err: err}
		}
		if err := ExtractZip(zipPath, tmpDir); err != nil {
			return Result{Tool: "aws-cli", Err: err}
		}
		installScript := filepath.Join(tmpDir, "aws", "install")
		cmd := exec.Command("sudo", installScript, "--install-dir", "/usr/local/aws-cli", "--bin-dir", "/usr/local/bin", "--update")
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: "aws-cli", Err: fmt.Errorf("aws cli install script failed: %w", err)}
		}
	}

	st.SetToolVersion("aws-cli", "v2-latest")
	return Result{Tool: "aws-cli", Version: "v2-latest"}
}

func (c *CloudCLIInstaller) installAzureCLI(p platform.Platform, st *state.State) Result {
	if p.OS == platform.Darwin {
		cmd := exec.Command("brew", "install", "azure-cli")
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: "azure-cli", Err: fmt.Errorf("brew install azure-cli failed: %w", err)}
		}
	} else {
		cmd := exec.Command("bash", "-c", "curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash")
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: "azure-cli", Err: fmt.Errorf("azure cli install failed: %w", err)}
		}
	}

	st.SetToolVersion("azure-cli", "latest")
	return Result{Tool: "azure-cli", Version: "latest"}
}

func (c *CloudCLIInstaller) installGcloudCLI(p platform.Platform, st *state.State) Result {
	tmpDir, err := os.MkdirTemp("", "cps-gcloud-*")
	if err != nil {
		return Result{Tool: "gcloud-cli", Err: err}
	}
	defer os.RemoveAll(tmpDir)

	var archStr string
	switch p.Arch {
	case platform.AMD64:
		archStr = "x86_64"
	case platform.ARM64:
		archStr = "arm"
	}

	url := fmt.Sprintf("https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-%s-%s.tar.gz", p.OS.String(), archStr)

	tarPath := filepath.Join(tmpDir, "gcloud.tar.gz")
	if err := DownloadToFile(url, tarPath); err != nil {
		return Result{Tool: "gcloud-cli", Err: err}
	}

	destDir := filepath.Join(p.HomeDir, "google-cloud-sdk")
	os.RemoveAll(destDir)

	if err := ExtractTarGz(tarPath, p.HomeDir); err != nil {
		return Result{Tool: "gcloud-cli", Err: fmt.Errorf("extract failed: %w", err)}
	}

	installScript := filepath.Join(destDir, "install.sh")
	cmd := exec.Command("bash", installScript, "--quiet", "--path-update=false", "--command-completion=false")
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: "gcloud-cli", Err: fmt.Errorf("gcloud install script failed: %w", err)}
	}

	st.SetToolVersion("gcloud-cli", "latest")
	return Result{Tool: "gcloud-cli", Version: "latest"}
}

