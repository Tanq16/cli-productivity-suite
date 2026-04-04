package cheatsheet

import "strings"

func buildFzfSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("FZF Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  CPS config: reverse layout, multi on, bat preview, fd as source") + "\n")
	b.WriteString(noteStyle.Render("  Alias: f=fzf") + "\n\n")

	// --- Basic Usage ---
	b.WriteString(headingStyle.Render("Basic Usage") + "\n")
	b.WriteString(cmdStyle.Render("  fzf") + "                             Fuzzy find files\n")
	b.WriteString(cmdStyle.Render("  <cmd> | fzf") + "                     Pipe anything into fzf\n")
	b.WriteString(divider + "\n")

	// --- Search Syntax ---
	b.WriteString(headingStyle.Render("Search Syntax") + "\n")
	b.WriteString(cmdStyle.Render("  foo") + "       fuzzy match          " + cmdStyle.Render("'foo") + "      exact match\n")
	b.WriteString(cmdStyle.Render("  ^foo") + "      starts with          " + cmdStyle.Render("foo$") + "      ends with\n")
	b.WriteString(cmdStyle.Render("  !foo") + "      does not match       " + cmdStyle.Render("!'foo") + "     not exact match\n")
	b.WriteString(cmdStyle.Render("  foo bar") + "   AND                   " + cmdStyle.Render("foo | bar") + " OR\n")
	b.WriteString(divider + "\n")

	// --- Keybindings inside fzf ---
	b.WriteString(headingStyle.Render("Keybindings (Inside FZF)") + "\n")
	b.WriteString(cmdStyle.Render("  Ctrl+j / Ctrl+k") + "                Scroll preview down / up\n")
	b.WriteString(cmdStyle.Render("  Enter") + "                           Select item\n")
	b.WriteString(cmdStyle.Render("  Tab / Shift+Tab") + "                Mark / unmark (multi)\n")
	b.WriteString(cmdStyle.Render("  Ctrl+c / Esc") + "                   Cancel\n")
	b.WriteString(divider + "\n")

	// --- Shell Integration ---
	b.WriteString(headingStyle.Render("Shell Integration (Zsh)") + "\n")
	b.WriteString(cmdStyle.Render("  Ctrl+t") + "                         Paste selected file path\n")
	b.WriteString(cmdStyle.Render("  Ctrl+r") + "                         Search command history\n")
	b.WriteString(cmdStyle.Render("  Alt+c") + "                          cd into selected directory\n")
	b.WriteString(divider + "\n")

	// --- Common Patterns ---
	b.WriteString(headingStyle.Render("Common Patterns") + "\n")
	b.WriteString(cmdStyle.Render("  vim $(fzf)") + "                     Open selected file in vim\n")
	b.WriteString(cmdStyle.Render("  cd $(fd --type d | fzf)") + "        cd into selected directory\n")
	b.WriteString(cmdStyle.Render("  git log --oneline | fzf") + "        Pick a commit\n")
	b.WriteString(cmdStyle.Render("  git branch | fzf | xargs git checkout") + "\n")
	b.WriteString("                                    " + "Switch git branch\n")
	b.WriteString(cmdStyle.Render("  ps aux | fzf") + "                   Find a process\n")

	return b.String()
}

var fzfSheet = Sheet{
	Name:        "fzf",
	Aliases:     []string{},
	Description: "FZF fuzzy finder usage, search syntax, and shell integration cheat sheet",
	Content:     buildFzfSheet(),
}
