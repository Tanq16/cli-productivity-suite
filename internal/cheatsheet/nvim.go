package cheatsheet

import "strings"

func buildNvimSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))
	leader := cmdStyle.Render("<Space>")

	b.WriteString(titleStyle.Render("Neovim Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  NvChad + catppuccin, leader: Space, aliased: vim=nvim") + "\n\n")

	// --- Modes ---
	b.WriteString(headingStyle.Render("Modes") + "\n")
	b.WriteString(cmdStyle.Render("  i") + "  insert   " + cmdStyle.Render("v") + "  visual   " + cmdStyle.Render("V") + "  visual line   " + cmdStyle.Render("Ctrl+v") + "  visual block\n")
	b.WriteString(cmdStyle.Render("  :") + "  command  " + cmdStyle.Render("R") + "  replace  " + cmdStyle.Render("Esc") + "  normal\n")
	b.WriteString(divider + "\n")

	// --- Navigation ---
	b.WriteString(headingStyle.Render("Navigation") + "\n")
	b.WriteString(cmdStyle.Render("  h j k l") + "                        Left, down, up, right\n")
	b.WriteString(cmdStyle.Render("  w / b / e") + "                      Next / prev / end of word\n")
	b.WriteString(cmdStyle.Render("  0 / $ / ^") + "                      Start / end / first non-blank\n")
	b.WriteString(cmdStyle.Render("  gg / G") + "                         Top / bottom of file\n")
	b.WriteString(cmdStyle.Render("  Ctrl+d / Ctrl+u") + "                Half-page down / up\n")
	b.WriteString(cmdStyle.Render("  Ctrl+f / Ctrl+b") + "                Full-page down / up\n")
	b.WriteString(noteStyle.Render("  In tmux: Ctrl+b is prefix — press C-b C-b to send to nvim") + "\n")
	b.WriteString(cmdStyle.Render("  { / }") + "                          Prev / next paragraph\n")
	b.WriteString(cmdStyle.Render("  % / :<number>") + "                  Matching bracket / go to line\n")
	b.WriteString(divider + "\n")

	// --- Editing ---
	b.WriteString(headingStyle.Render("Editing") + "\n")
	b.WriteString(cmdStyle.Render("  dd / yy / p") + "                    Delete / yank / paste line\n")
	b.WriteString(cmdStyle.Render("  x") + "                              Delete char under cursor\n")
	b.WriteString(cmdStyle.Render("  u / Ctrl+r") + "                     Undo / redo\n")
	b.WriteString(cmdStyle.Render("  ciw / ci\" / ci(") + "                Change inner word/quotes/parens\n")
	b.WriteString(cmdStyle.Render("  diw / di\" / di(") + "                Delete inner word/quotes/parens\n")
	b.WriteString(cmdStyle.Render("  >> / <<") + "                        Indent / outdent line\n")
	b.WriteString(cmdStyle.Render("  o / O") + "                          New line below / above\n")
	b.WriteString(cmdStyle.Render("  A / J / .") + "                      Append EOL / join line / repeat\n")
	b.WriteString(divider + "\n")

	// --- Search ---
	b.WriteString(headingStyle.Render("Search & Replace") + "\n")
	b.WriteString(cmdStyle.Render("  / / ?") + "                          Search forward / backward\n")
	b.WriteString(cmdStyle.Render("  n / N / *") + "                      Next / prev / word under cursor\n")
	b.WriteString(cmdStyle.Render("  :noh") + "                           Clear search highlight\n")
	b.WriteString(cmdStyle.Render("  :%s/old/new/g[c]") + "               Replace all [c = confirm each]\n")
	b.WriteString(divider + "\n")

	// --- NvChad Specific ---
	b.WriteString(headingStyle.Render("NvChad — File & Buffer") + "\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("ff / fw") + "                  Find file / live grep (Telescope)\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("fb / fo") + "                  Find buffer / recent files\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("e") + "                       Toggle NvimTree\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("x") + "                       Close buffer\n")
	b.WriteString(cmdStyle.Render("  Tab / Shift+Tab") + "                Next / prev buffer\n")
	b.WriteString(divider + "\n")

	// --- NvChad LSP ---
	b.WriteString(headingStyle.Render("NvChad — LSP") + "\n")
	b.WriteString(cmdStyle.Render("  gd / gr") + "                        Go to definition / references\n")
	b.WriteString(cmdStyle.Render("  K") + "                              Hover documentation\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("ra / ca") + "                  Rename / code actions\n")
	b.WriteString(cmdStyle.Render("  ]d / [d") + "                        Next / prev diagnostic\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("fm") + "                      Format file\n")
	b.WriteString(divider + "\n")

	// --- File Ops ---
	b.WriteString(headingStyle.Render("File Operations") + "\n")
	b.WriteString(cmdStyle.Render("  :w / :q / :wq / :q!") + "            Save / quit / save+quit / force quit\n")
	b.WriteString(cmdStyle.Render("  :e <file>") + "                      Open file\n")
	b.WriteString(cmdStyle.Render("  :vs / :sp") + "                      Vertical / horizontal split\n")
	b.WriteString(cmdStyle.Render("  Ctrl+w h/j/k/l") + "                Navigate splits\n")

	return b.String()
}

var nvimSheet = Sheet{
	Name:        "nvim",
	Aliases:     []string{"neovim", "vim", "vi"},
	Description: "Neovim/NvChad navigation, editing, LSP, and Telescope cheat sheet",
	Content:     buildNvimSheet(),
}
