package github

import (
	"cmp"
	"os"
	"os/exec"
	"strings"
)

func TokenFromEnv() string {
	return cmp.Or(os.Getenv("GITHUB_TOKEN"), os.Getenv("GH_TOKEN"))
}

func ResolveToken(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
