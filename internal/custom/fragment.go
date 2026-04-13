package custom

import (
	"fmt"
	"strings"
)

func RenderFragment(pf PackFile) []byte {
	if len(pf.Shell.Env) == 0 && len(pf.Shell.PathPrepend) == 0 && len(pf.Shell.Source) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("# CPS custom extension: %s\n", pf.Name))

	for k, v := range pf.Shell.Env {
		b.WriteString(fmt.Sprintf("export %s=\"%s\"\n", k, v))
	}

	if len(pf.Shell.PathPrepend) > 0 {
		b.WriteString("export PATH=\"" + strings.Join(pf.Shell.PathPrepend, ":") + ":$PATH\"\n")
	}

	for _, src := range pf.Shell.Source {
		b.WriteString(fmt.Sprintf("[ -f \"%s\" ] && source \"%s\"\n", src, src))
	}

	return []byte(b.String())
}
