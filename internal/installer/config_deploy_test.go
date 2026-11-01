package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tanq16/cli-productivity-suite/internal/configs"
)

func TestEnsureProfileLoaderKeepsExistingLinesAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zprofile")
	existing := "eval \"$(/opt/homebrew/bin/brew shellenv)\""
	if err := os.WriteFile(path, []byte(existing), 0644); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := ensureProfileLoader(path); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(got), existing+"\n") {
		t.Errorf("existing line not preserved:\n%s", got)
	}
	if n := strings.Count(string(got), string(configs.RcProfileLoader())); n != 1 {
		t.Errorf("loader appended %d times, want 1", n)
	}
}

func TestEnsureProfileLoaderCreatesMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zprofile")
	if err := ensureProfileLoader(path); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(configs.RcProfileLoader()) {
		t.Errorf("got %q", got)
	}
}

func TestEnsureProfileLoaderRecognisesAnAlteredLoader(t *testing.T) {
	for name, existing := range map[string]string{
		"crlf":             strings.ReplaceAll(string(configs.RcProfileLoader()), "\n", "\r\n"),
		"edited comment":   "# my note\n[ -f \"$HOME/shell/rc/cps-env.zsh\" ] && source \"$HOME/shell/rc/cps-env.zsh\"",
		"no final newline": strings.TrimSuffix(string(configs.RcProfileLoader()), "\n"),
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".zprofile")
			if err := os.WriteFile(path, []byte(existing), 0644); err != nil {
				t.Fatal(err)
			}
			if err := ensureProfileLoader(path); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != existing {
				t.Errorf("file changed:\n%q", got)
			}
		})
	}
}

func TestEnsureProfileLoaderDoesNotStackBlankLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zprofile")
	if err := os.WriteFile(path, []byte("export A=1\n\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ensureProfileLoader(path); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := "export A=1\n\n" + string(configs.RcProfileLoader()); string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
