package cheatsheet

import "strings"

func buildTmuxSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("Tmux Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  CPS config: mouse on, emacs copy-mode, catppuccin theme, status top") + "\n")
	b.WriteString(noteStyle.Render("  Windows/panes start at 1, auto-rename on, history 99999 lines") + "\n\n")

	// --- Shell Aliases ---
	b.WriteString(headingStyle.Render("Shell Aliases") + "\n")
	b.WriteString(cmdStyle.Render("  tt") + "                              New session named \"default\"\n")
	b.WriteString(cmdStyle.Render("  t") + "                               Attach to \"default\" session\n")
	b.WriteString(cmdStyle.Render("  ts <name>") + "                       New session with custom name\n")
	b.WriteString(cmdStyle.Render("  ta <name>") + "                       Attach to named session\n")
	b.WriteString(cmdStyle.Render("  tls") + "                             List sessions\n")
	b.WriteString(divider + "\n")

	// --- Custom Keybindings (no prefix) ---
	b.WriteString(headingStyle.Render("Window & Pane Navigation (No Prefix)") + "\n")
	b.WriteString(cmdStyle.Render("  Alt+=") + "                           Next window\n")
	b.WriteString(cmdStyle.Render("  Alt+-") + "                           Previous window\n")
	b.WriteString(cmdStyle.Render("  Shift+Arrow") + "                     Move between panes\n")
	b.WriteString(cmdStyle.Render("  Alt+\\") + "                           Split pane horizontal\n")
	b.WriteString(cmdStyle.Render("  Alt+|") + "                           Split pane vertical\n")
	b.WriteString(cmdStyle.Render("  Alt+v") + "                           Paste buffer\n")
	b.WriteString(divider + "\n")

	// --- Copy Mode (emacs) ---
	b.WriteString(headingStyle.Render("Copy Mode (Emacs) — Mouse & Keyboard") + "\n")
	b.WriteString(noteStyle.Render("  Scroll up with mouse to enter copy mode automatically.") + "\n")
	b.WriteString(noteStyle.Render("  Mouse drag selects text, releasing auto-copies it.") + "\n\n")

	b.WriteString("  " + headingStyle.Render("Keyboard selection flow:") + "\n")
	b.WriteString(cmdStyle.Render("  C-Space") + "                         Start selection\n")
	b.WriteString(cmdStyle.Render("  Alt+w") + "                           Copy selection\n")
	b.WriteString(cmdStyle.Render("  Alt+v") + "                           Paste " + noteStyle.Render("(outside copy mode)") + "\n\n")

	b.WriteString("  " + headingStyle.Render("Navigation in copy mode:") + "\n")
	b.WriteString(cmdStyle.Render("  Space") + "                           Page down\n")
	b.WriteString(cmdStyle.Render("  Alt+v") + "                           Page up " + noteStyle.Render("(in copy mode)") + "\n")
	b.WriteString(cmdStyle.Render("  Alt+>") + "                           Jump to bottom of buffer\n")
	b.WriteString(cmdStyle.Render("  Alt+<") + "                           Jump to top of buffer\n")
	b.WriteString(cmdStyle.Render("  Arrow keys") + "                      Move cursor\n")
	b.WriteString(cmdStyle.Render("  q") + "                               Exit copy mode\n\n")

	b.WriteString("  " + headingStyle.Render("Select-all-to-bottom flow:") + "\n")
	b.WriteString(noteStyle.Render("  1. Scroll up with mouse to enter copy mode") + "\n")
	b.WriteString(noteStyle.Render("  2. C-Space to start selection at cursor") + "\n")
	b.WriteString(noteStyle.Render("  3. Alt+> to jump to the absolute bottom") + "\n")
	b.WriteString(noteStyle.Render("  4. Alt+w to copy the selection") + "\n")
	b.WriteString(divider + "\n")

	// --- Prefix Commands ---
	b.WriteString(headingStyle.Render("Prefix Commands (C-b)") + "\n")
	b.WriteString(noteStyle.Render("  Rarely needed — most ops have direct bindings above.") + "\n\n")
	b.WriteString(cmdStyle.Render("  C-b c") + "                           Create new window\n")
	b.WriteString(cmdStyle.Render("  C-b ,") + "                           Rename window\n")
	b.WriteString(cmdStyle.Render("  C-b w") + "                           List all windows\n")
	b.WriteString(cmdStyle.Render("  C-b s") + "                           List and switch sessions\n")
	b.WriteString(cmdStyle.Render("  C-b d") + "                           Detach from session\n")
	b.WriteString(cmdStyle.Render("  C-b $") + "                           Rename session\n")
	b.WriteString(cmdStyle.Render("  C-b x") + "                           Kill current pane\n")
	b.WriteString(cmdStyle.Render("  C-b &") + "                           Kill current window\n")
	b.WriteString(cmdStyle.Render("  C-b z") + "                           Toggle pane zoom (fullscreen)\n")
	b.WriteString(cmdStyle.Render("  C-b q") + "                           Show pane numbers\n")
	b.WriteString(cmdStyle.Render("  C-b :") + "                           Command prompt\n")

	return b.String()
}

var tmuxSheet = Sheet{
	Name:        "tmux",
	Aliases:     []string{},
	Description: "Tmux cheat sheet tailored to CPS config",
	Content:     buildTmuxSheet(),
}
