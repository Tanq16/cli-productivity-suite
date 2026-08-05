package cheatsheet

import "strings"

func buildNvimSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))
	leader := cmdStyle.Render("<Space>")

	b.WriteString(titleStyle.Render("Neovim Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  cps config, leader: Space, colors inherited from terminal") + "\n\n")

	b.WriteString(headingStyle.Render("Modes") + "\n")
	b.WriteString(cmdStyle.Render("  i") + "  insert   " + cmdStyle.Render("v") + "  visual   " + cmdStyle.Render("V") + "  visual line   " + cmdStyle.Render("Ctrl+v") + "  visual block\n")
	b.WriteString(cmdStyle.Render("  :") + "  command  " + cmdStyle.Render("R") + "  replace  " + cmdStyle.Render("Esc") + "  normal\n")
	b.WriteString(divider + "\n")

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

	b.WriteString(headingStyle.Render("Insert Mode") + "\n")
	b.WriteString(cmdStyle.Render("  Ctrl+a / Ctrl+e") + "                Start / end of line\n")
	b.WriteString(cmdStyle.Render("  Ctrl+Left / Ctrl+Right") + "         Back / forward one word\n")
	b.WriteString(cmdStyle.Render("  Left / Right") + "                   Wrap across line boundaries\n")
	b.WriteString(cmdStyle.Render("  Ctrl+w / Ctrl+u") + "                Delete word / line before cursor\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Editing") + "\n")
	b.WriteString(cmdStyle.Render("  dd / yy / p") + "                    Delete / yank / paste line\n")
	b.WriteString(cmdStyle.Render("  x") + "                              Delete char under cursor\n")
	b.WriteString(cmdStyle.Render("  u / Ctrl+r") + "                     Undo / redo\n")
	b.WriteString(cmdStyle.Render("  ciw / ci\" / ci(") + "                Change inner word/quotes/parens\n")
	b.WriteString(cmdStyle.Render("  diw / di\" / di(") + "                Delete inner word/quotes/parens\n")
	b.WriteString(cmdStyle.Render("  >> / <<") + "                        Indent / outdent line\n")
	b.WriteString(cmdStyle.Render("  o / O") + "                          New line below / above\n")
	b.WriteString(cmdStyle.Render("  A / J / .") + "                      Append EOL / join line / repeat\n")
	b.WriteString(cmdStyle.Render("  gcc / gc") + "                       Toggle comment line / selection\n")
	b.WriteString(cmdStyle.Render("  \"+y / \"+p") + "                      Yank / paste via system clipboard\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Search & Replace") + "\n")
	b.WriteString(cmdStyle.Render("  / / ?") + "                          Search forward / backward\n")
	b.WriteString(cmdStyle.Render("  n / N / *") + "                      Next / prev / word under cursor\n")
	b.WriteString(cmdStyle.Render("  Esc") + "                            Clear search highlight\n")
	b.WriteString(cmdStyle.Render("  :%s/old/new/g[c]") + "               Replace all [c = confirm each]\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Files & Buffers") + "\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("ff") + "                      Find files (fzf-lua)\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("fw") + "                      Live grep across project\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("fa") + "                      Find all files (hidden + ignored)\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("fb") + "                      Find buffers\n")
	b.WriteString(cmdStyle.Render("  Ctrl+n") + "                         Toggle file tree\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("e") + "                       Focus file tree\n")
	b.WriteString(cmdStyle.Render("  ]b / [b") + "                        Next / prev buffer\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("x") + "                       Close buffer\n")
	b.WriteString(cmdStyle.Render("  Ctrl+s") + "                         Save file\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("LSP (built-in)") + "\n")
	b.WriteString(cmdStyle.Render("  gd / gD") + "                        Go to definition / declaration\n")
	b.WriteString(cmdStyle.Render("  grr / gri") + "                      References / implementations\n")
	b.WriteString(cmdStyle.Render("  grn / gra") + "                      Rename / code action\n")
	b.WriteString(cmdStyle.Render("  grt / gO") + "                       Type definition / document symbols\n")
	b.WriteString(cmdStyle.Render("  K") + "                              Hover documentation\n")
	b.WriteString(cmdStyle.Render("  ]d / [d") + "                        Next / prev diagnostic\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("ds") + "                      Diagnostics to loclist\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("fm") + "                      Format buffer\n")
	b.WriteString(noteStyle.Render("  Servers: gopls, ruff, ts_ls — enabled only when on PATH") + "\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Git (gitsigns)") + "\n")
	b.WriteString(cmdStyle.Render("  ]c / [c") + "                        Next / prev hunk\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("hp / hb") + "                  Preview hunk / blame line\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("hs / hr") + "                  Stage / reset hunk\n")
	b.WriteString("  " + leader + " " + cmdStyle.Render("hd") + "                      Diff this file\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Windows & Files") + "\n")
	b.WriteString(cmdStyle.Render("  :w / :q / :wq / :q!") + "            Save / quit / save+quit / force quit\n")
	b.WriteString(cmdStyle.Render("  :e <file>") + "                      Open file\n")
	b.WriteString(cmdStyle.Render("  :vs / :sp") + "                      Vertical / horizontal split\n")
	b.WriteString(cmdStyle.Render("  Ctrl+h/j/k/l") + "                   Navigate splits\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Plugins") + "\n")
	b.WriteString(cmdStyle.Render("  :lua vim.pack.update()") + "         Update plugins (review, :w to apply)\n")
	b.WriteString(cmdStyle.Render("  :TSUpdate") + "                      Update treesitter parsers\n")
	b.WriteString(cmdStyle.Render("  :checkhealth") + "                   Diagnose config problems\n")

	return b.String()
}

var nvimSheet = Sheet{
	Name:        "nvim",
	Aliases:     []string{"neovim", "vim", "vi"},
	Description: "Neovim navigation, editing, LSP, git, and fuzzy-find cheat sheet",
	Content:     buildNvimSheet(),
}
