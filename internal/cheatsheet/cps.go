package cheatsheet

import "strings"

func buildCPSSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("CPS Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  Tools: ~/shell/executables/   Extensions: ~/shell/extensions/") + "\n")
	b.WriteString(noteStyle.Render("  State: ~/.config/cps/state.json") + "\n\n")

	// --- Core Commands ---
	b.WriteString(headingStyle.Render("Core Commands") + "\n")
	b.WriteString(cmdStyle.Render("  cps init") + "                        Full environment setup\n")
	b.WriteString(cmdStyle.Render("  cps check") + "                       Check all tools for updates\n")
	b.WriteString(cmdStyle.Render("  cps self-update") + "                 Update cps binary itself\n")
	b.WriteString(cmdStyle.Render("  cps clean") + "                       Remove all CPS-managed dirs\n")
	b.WriteString(divider + "\n")

	// --- Install ---
	b.WriteString(headingStyle.Render("Install — Individual Tools or Categories") + "\n")
	b.WriteString(cmdStyle.Render("  cps install <tool>") + "              Install a single tool by name\n")
	b.WriteString(cmdStyle.Render("  cps install <tool1> <tool2>") + "     Install multiple tools\n")
	b.WriteString(cmdStyle.Render("  cps install core") + "                All core binary tools\n")
	b.WriteString(cmdStyle.Render("  cps install system") + "              System packages (apt/brew)\n")
	b.WriteString(cmdStyle.Render("  cps install cloud") + "               Cloud CLIs (AWS, Azure, gcloud)\n")
	b.WriteString(cmdStyle.Render("  cps install runtimes") + "            Language runtimes (Go, Python, Rust, etc.)\n")
	b.WriteString(cmdStyle.Render("  cps install configs") + "             Config files + shell plugins\n")
	b.WriteString(divider + "\n")

	// --- Extend ---
	b.WriteString(headingStyle.Render("Extend — Extension Packs") + "\n")
	b.WriteString(cmdStyle.Render("  cps extend list") + "                 List available packs\n")
	b.WriteString(cmdStyle.Render("  cps extend <pack>") + "               Install a pack (security, cloudsec, appsec, misc, private)\n")
	b.WriteString(cmdStyle.Render("  cps extend --check <pack>") + "       Check extension pack for updates\n")
	b.WriteString(noteStyle.Render("  Extension tools also work with: cps install <tool-name>") + "\n")
	b.WriteString(divider + "\n")

	// --- Other ---
	b.WriteString(headingStyle.Render("Other") + "\n")
	b.WriteString(cmdStyle.Render("  cps cheat <topic>") + "               Cheat sheets (uv, rust, tmux, nvim, fzf, regex)\n")
	b.WriteString(cmdStyle.Render("  --gh-token <token>") + "              GitHub PAT for private repos\n")
	b.WriteString(cmdStyle.Render("  --debug") + "                         Verbose debug logging\n")
	b.WriteString(cmdStyle.Render("  --for-ai") + "                        AI-friendly output (no color)\n")

	return b.String()
}

var cpsSheet = Sheet{
	Name:        "cps",
	Aliases:     []string{},
	Description: "CPS commands, install categories, extensions, and init phases",
	Content:     buildCPSSheet(),
}
