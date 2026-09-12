package cheatsheet

import "strings"

func buildCPSSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("CPS Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  Tools: ~/shell/extensions/       User binaries: ~/shell/custom-bin/") + "\n")
	b.WriteString(noteStyle.Render("  Generated rc file: ~/shell/rc/cps.zsh") + "\n")
	b.WriteString(noteStyle.Render("  App bundles: ~/shell/apps/ (not on PATH - run by full path)") + "\n\n")

	b.WriteString(headingStyle.Render("First Run") + "\n")
	b.WriteString(cmdStyle.Render("  cps prereq") + "                      System packages and Homebrew (uses sudo)\n")
	b.WriteString(cmdStyle.Render("  cps shell") + "                       Shell binaries, Neovim, plugins and configs\n")
	b.WriteString(cmdStyle.Render("  cps package install <group>") + "     Everything past the shell\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Packages") + "\n")
	b.WriteString(cmdStyle.Render("  cps package list") + "                Every group, suite and package\n")
	b.WriteString(cmdStyle.Render("  cps package install <name>") + "      One package, a whole group, or a suite\n")
	b.WriteString(noteStyle.Render("  Groups: brew, macos-desktop, runtime, ai, binaries, homelab, private") + "\n")
	b.WriteString(noteStyle.Render("  Suites: go-suite, js-suite, python-suite, java-suite, rust-suite") + "\n")
	b.WriteString(noteStyle.Render("  Missing dependencies are installed first; re-running updates") + "\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Status") + "\n")
	b.WriteString(cmdStyle.Render("  cps status") + "                      Installed packages and recorded versions\n")
	b.WriteString(cmdStyle.Render("  cps status --check") + "              Compare against upstream (one API call per package)\n")
	b.WriteString(cmdStyle.Render("  cps status --json") + "               The same data as JSON\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("System (uses sudo)") + "\n")
	b.WriteString(cmdStyle.Render("  cps system update") + "               apt and brew upgrade, autoremove, cleanup\n")
	b.WriteString(cmdStyle.Render("  cps system kitty") + "                Run the kitty terminal installer\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Adding Your Own Stuff") + "\n")
	b.WriteString(noteStyle.Render("  Aliases / exports / funcs   →  drop a *.zsh file in ~/shell/rc/custom/") + "\n")
	b.WriteString(noteStyle.Render("  Your own binaries           →  drop them in ~/shell/custom-bin/ (on PATH)") + "\n")
	b.WriteString(cmdStyle.Render("  cps custom <args>") + "               Run ~/.config/cps/custom.sh with the cps PATH\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Theme") + "\n")
	b.WriteString(cmdStyle.Render("  cps theme") + "                       List themes, marking the active one\n")
	b.WriteString(cmdStyle.Render("  cps theme <name>") + "                Switch the terminal palette\n")
	b.WriteString(noteStyle.Render("  Repaints kitty, tmux and nvim at once - nothing restarts") + "\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Other") + "\n")
	b.WriteString(cmdStyle.Render("  cps cheat <topic>") + "               Cheat sheets (go, java, uv, fnm, bun, rust, tmux, nvim, fzf, jq, regex)\n")
	b.WriteString(cmdStyle.Render("  cps self-update") + "                 Update cps itself to the latest release\n")
	b.WriteString(cmdStyle.Render("  --gh-token <token>") + "              GitHub PAT, on package install and status\n")
	b.WriteString(cmdStyle.Render("  --debug") + "                         Verbose debug logging\n")

	return b.String()
}

var cpsSheet = Sheet{
	Name:        "cps",
	Aliases:     []string{},
	Description: "CPS commands, package groups, and shell integration",
	Content:     buildCPSSheet(),
}
