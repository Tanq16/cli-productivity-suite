package cheatsheet

import "strings"

func buildFnmSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("fnm / Node.js Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  CPS paths: FNM_DIR=~/shell/fnm") + "\n")
	b.WriteString(noteStyle.Render("  Node versions and global packages live under FNM_DIR") + "\n\n")

	// --- Version Management ---
	b.WriteString(headingStyle.Render("Version Management") + "\n")
	b.WriteString(cmdStyle.Render("  fnm install --lts") + "               Install latest LTS version\n")
	b.WriteString(cmdStyle.Render("  fnm install <version>") + "           Install a specific version\n")
	b.WriteString(cmdStyle.Render("  fnm use <version>") + "               Switch to a version (current shell)\n")
	b.WriteString(cmdStyle.Render("  fnm default <version>") + "           Set the default version\n")
	b.WriteString(cmdStyle.Render("  fnm current") + "                     Show active version\n")
	b.WriteString(cmdStyle.Render("  fnm list") + "                        List installed versions\n")
	b.WriteString(cmdStyle.Render("  fnm list-remote") + "                 List available versions\n")
	b.WriteString(cmdStyle.Render("  fnm uninstall <version>") + "         Remove a version\n")
	b.WriteString(divider + "\n")

	// --- Per-Project Versions ---
	b.WriteString(headingStyle.Render("Per-Project Versions") + "\n")
	b.WriteString(noteStyle.Render("  fnm auto-detects .node-version or .nvmrc in the project dir.") + "\n")
	b.WriteString(noteStyle.Render("  eval \"$(fnm env)\" in .zshrc enables auto-switching.") + "\n\n")
	b.WriteString(cmdStyle.Render("  echo 20 > .node-version") + "        Pin project to Node 20.x\n")
	b.WriteString(cmdStyle.Render("  fnm use") + "                         Switch to project's pinned version\n")
	b.WriteString(divider + "\n")

	// --- npm Basics ---
	b.WriteString(headingStyle.Render("npm — Package Management") + "\n")
	b.WriteString(cmdStyle.Render("  npm init -y") + "                     Create package.json\n")
	b.WriteString(cmdStyle.Render("  npm install") + "                     Install project dependencies\n")
	b.WriteString(cmdStyle.Render("  npm install <pkg>") + "               Add a dependency\n")
	b.WriteString(cmdStyle.Render("  npm install -D <pkg>") + "            Add a dev dependency\n")
	b.WriteString(cmdStyle.Render("  npm install -g <pkg>") + "            Install globally\n")
	b.WriteString(cmdStyle.Render("  npm list -g --depth=0") + "           List global packages\n")
	b.WriteString(cmdStyle.Render("  npm update") + "                      Update project dependencies\n")
	b.WriteString(divider + "\n")

	// --- npx ---
	b.WriteString(headingStyle.Render("npx — Run Without Installing") + "\n")
	b.WriteString(cmdStyle.Render("  npx <pkg>") + "                       Run a package directly\n")
	b.WriteString(cmdStyle.Render("  npx <pkg>@<version>") + "             Run a specific version\n")
	b.WriteString(divider + "\n")

	// --- Key Concepts ---
	b.WriteString(headingStyle.Render("Key Concepts") + "\n")
	b.WriteString("  • " + cmdStyle.Render("fnm") + " manages Node versions — each version is isolated under FNM_DIR\n")
	b.WriteString("  • " + cmdStyle.Render("npm install -g") + "  →  global packages are per-Node-version\n")
	b.WriteString("  • " + cmdStyle.Render(".node-version") + "   →  auto-switches Node when entering a project dir\n")
	b.WriteString("  • " + cmdStyle.Render("cps extend runtimes node") + " →  updates fnm + Node LTS, preserves global npm packages\n")

	return b.String()
}

var fnmSheet = Sheet{
	Name:        "fnm",
	Aliases:     []string{"node", "npm"},
	Description: "fnm version manager, npm, and Node.js cheat sheet",
	Content:     buildFnmSheet(),
}
