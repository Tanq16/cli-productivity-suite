package cheatsheet

import "strings"

func buildCPSSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("CPS Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  Tools: ~/shell/extensions/       User binaries: ~/shell/custom-bin/") + "\n")
	b.WriteString(noteStyle.Render("  RC fragments: ~/shell/rc/") + "\n")
	b.WriteString(noteStyle.Render("  App bundles: ~/shell/apps/ (not on PATH — run by full path)") + "\n\n")

	b.WriteString(headingStyle.Render("Core Commands") + "\n")
	b.WriteString(cmdStyle.Render("  cps init") + "                        Base shell environment setup\n")
	b.WriteString(cmdStyle.Render("  cps self-update") + "                 Update cps binary itself\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Extend — Extension Packs") + "\n")
	b.WriteString(cmdStyle.Render("  cps extend list") + "                 List available packs\n")
	b.WriteString(cmdStyle.Render("  cps extend <pack>") + "               Install entire extension pack\n")
	b.WriteString(cmdStyle.Render("  cps extend <pack> <tool> ...") + "    Install specific tools from a pack\n")
	b.WriteString(noteStyle.Render("  Packs: essentials, core, cloud, runtimes, security, cloudsec, appsec,") + "\n")
	b.WriteString(noteStyle.Render("         misc, private, ai-tools, homelab") + "\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Shell Integration — RC Fragments") + "\n")
	b.WriteString(noteStyle.Render("  ~/.zshrc sources ~/shell/rc/*.zsh and ~/shell/rc/custom/*.zsh") + "\n")
	b.WriteString(noteStyle.Render("  00-base.zsh        deployed by cps init") + "\n")
	b.WriteString(noteStyle.Render("  10-runtimes.zsh    deployed by cps extend runtimes") + "\n")
	b.WriteString(noteStyle.Render("  20-cloud.zsh       deployed by cps extend cloud") + "\n")
	b.WriteString(noteStyle.Render("  30-security.zsh    deployed by cps extend security") + "\n")
	b.WriteString(noteStyle.Render("  40-misc.zsh        deployed by cps extend misc") + "\n")
	b.WriteString(noteStyle.Render("  custom/*.zsh       user-managed fragments") + "\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Adding Your Own Stuff (no extension pack needed)") + "\n")
	b.WriteString(noteStyle.Render("  Aliases / exports / funcs   →  drop a *.zsh file in ~/shell/rc/custom/") + "\n")
	b.WriteString(noteStyle.Render("  Your own binaries           →  drop them in ~/shell/custom-bin/ (on PATH)") + "\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Other") + "\n")
	b.WriteString(cmdStyle.Render("  cps cheat <topic>") + "               Cheat sheets (go, java, uv, fnm, bun, rust, tmux, nvim, fzf, jq, regex)\n")
	b.WriteString(cmdStyle.Render("  --gh-token <token>") + "              GitHub PAT for private repos\n")
	b.WriteString(cmdStyle.Render("  --debug") + "                         Verbose debug logging\n")
	b.WriteString(cmdStyle.Render("  --for-ai") + "                        AI-friendly output (no color)\n")

	return b.String()
}

var cpsSheet = Sheet{
	Name:        "cps",
	Aliases:     []string{},
	Description: "CPS commands, extension packs, and shell integration",
	Content:     buildCPSSheet(),
}
