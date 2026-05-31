package cheatsheet

import "strings"

func buildBunSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("Bun Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  CPS paths: BUN_INSTALL=~/shell/bun  (binary + global packages)") + "\n")
	b.WriteString(noteStyle.Render("  Bun is a runtime + package manager + bundler + test runner") + "\n\n")

	b.WriteString(headingStyle.Render("Package Management") + "\n")
	b.WriteString(cmdStyle.Render("  bun init") + "                        Scaffold a new project (package.json + tsconfig)\n")
	b.WriteString(cmdStyle.Render("  bun install") + "                     Install project dependencies (uses bun.lockb)\n")
	b.WriteString(cmdStyle.Render("  bun add <pkg>") + "                   Add a dependency\n")
	b.WriteString(cmdStyle.Render("  bun add -d <pkg>") + "                Add a dev dependency\n")
	b.WriteString(cmdStyle.Render("  bun add -g <pkg>") + "                Install globally (under BUN_INSTALL/bin)\n")
	b.WriteString(cmdStyle.Render("  bun remove <pkg>") + "                Remove a dependency\n")
	b.WriteString(cmdStyle.Render("  bun update") + "                      Update dependencies to latest allowed by ranges\n")
	b.WriteString(cmdStyle.Render("  bun outdated") + "                    Show outdated dependencies\n")
	b.WriteString(cmdStyle.Render("  bun pm ls") + "                       List installed packages (tree)\n")
	b.WriteString(cmdStyle.Render("  bun pm cache rm") + "                 Clear the global module cache\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Running Scripts & Files") + "\n")
	b.WriteString(cmdStyle.Render("  bun run <script>") + "                Run a package.json script\n")
	b.WriteString(cmdStyle.Render("  bun <file.ts>") + "                   Execute a TS/JS file directly (no compile step)\n")
	b.WriteString(cmdStyle.Render("  bun --hot <file>") + "                Run with hot reload (preserves state)\n")
	b.WriteString(cmdStyle.Render("  bun --watch <file>") + "              Run with file-watcher restart\n")
	b.WriteString(cmdStyle.Render("  bun x <pkg>") + "                     Run a package without installing (npx equivalent)\n")
	b.WriteString(cmdStyle.Render("  bunx <pkg>") + "                      Alias for 'bun x'\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Test Runner (Jest-compatible API)") + "\n")
	b.WriteString(cmdStyle.Render("  bun test") + "                        Run all tests (*.test.ts, *.spec.ts, etc.)\n")
	b.WriteString(cmdStyle.Render("  bun test <pattern>") + "              Filter by file pattern\n")
	b.WriteString(cmdStyle.Render("  bun test --watch") + "                Re-run on file change\n")
	b.WriteString(cmdStyle.Render("  bun test --coverage") + "             Print line coverage\n")
	b.WriteString(cmdStyle.Render("  bun test -t \"<name>\"") + "            Only run tests whose name matches\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Bundler") + "\n")
	b.WriteString(cmdStyle.Render("  bun build <entry>") + "               Bundle to stdout\n")
	b.WriteString(cmdStyle.Render("  bun build <entry> --outdir <dir>") + " Bundle to a directory\n")
	b.WriteString(cmdStyle.Render("  bun build <entry> --target=node") + "  Target Node.js instead of bun/browser\n")
	b.WriteString(cmdStyle.Render("  bun build <entry> --minify") + "      Minify output\n")
	b.WriteString(cmdStyle.Render("  bun build --compile <entry>") + "     Produce a standalone executable\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Misc") + "\n")
	b.WriteString(cmdStyle.Render("  bun upgrade") + "                     Update bun itself to the latest release\n")
	b.WriteString(cmdStyle.Render("  bun --version") + "                   Show installed bun version\n")
	b.WriteString(cmdStyle.Render("  bun repl") + "                        Start an interactive REPL\n")
	b.WriteString(divider + "\n")

	b.WriteString(headingStyle.Render("Key Concepts") + "\n")
	b.WriteString("  • " + cmdStyle.Render("bun") + " runs .ts/.tsx/.jsx natively — no tsc/ts-node needed\n")
	b.WriteString("  • " + cmdStyle.Render("bun install") + " is a drop-in for npm/yarn/pnpm; reads package.json, writes bun.lockb\n")
	b.WriteString("  • " + cmdStyle.Render("bun.lockb") + " is binary — diff with 'bun pm ls' or set bun.lockb=binary in .gitattributes\n")
	b.WriteString("  • " + cmdStyle.Render("Node compat") + " is good but not 100% — some native modules and rare APIs still differ\n")
	b.WriteString("  • " + cmdStyle.Render("Global installs") + " land in $BUN_INSTALL/bin, which CPS adds to PATH\n")

	return b.String()
}

var bunSheet = Sheet{
	Name:        "bun",
	Aliases:     []string{},
	Description: "Bun runtime, package manager, bundler, and test runner",
	Content:     buildBunSheet(),
}
