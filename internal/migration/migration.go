package migration

import (
	"strconv"
	"strings"
)

type Migration struct {
	From  string
	Notes []string
}

var migrations = []Migration{
	{
		From: "v1.10.0",
		Notes: []string{
			"rm ~/shell/rc/40-misc.zsh   (the misc pack is gone; NEO4J_CONF now ships in 50-homelab.zsh)",
			"cps extend security nuclei-templates   (swaps the old git clone, which grows past 800MB, for a ~78MB snapshot)",
			"cps extend homelab neo4j   (neo4j moved out of misc, and this redeploys the fragment its config needs)",
			"appsec is gone and its tools are in security; gowitness moved there too. Installed binaries are unaffected.",
			"Directories cps snapshots (~/shell/plugins/*, ~/shell/nuclei-templates) are replaced wholesale on update. Keep your own files elsewhere.",
		},
	},
}

func Pending(stateVersion, appVersion string) []Migration {
	app, ok := parse(appVersion)
	if !ok || stateVersion == appVersion {
		return nil
	}
	state, tracked := parse(stateVersion)

	var pending []Migration
	for _, m := range migrations {
		from, ok := parse(m.From)
		if !ok || compare(from, app) >= 0 {
			continue
		}
		if tracked && compare(from, state) < 0 {
			continue
		}
		pending = append(pending, m)
	}
	return pending
}

func parse(v string) ([3]int, bool) {
	var out [3]int
	parts := strings.Split(strings.TrimPrefix(v, "v"), ".")
	if len(parts) != 3 {
		return out, false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return out, false
		}
		out[i] = n
	}
	return out, true
}

func compare(a, b [3]int) int {
	for i := range a {
		if a[i] != b[i] {
			if a[i] < b[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}
