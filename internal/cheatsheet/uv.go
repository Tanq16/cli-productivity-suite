package cheatsheet

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(12)).PaddingBottom(1)
	headingStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(11))
	cmdStyle     = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))
	noteStyle    = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8)).Italic(true)
	dividerStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
)

func buildUVSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("UV Cheat Sheet") + "\n")

	// --- Python Version Management ---
	b.WriteString(headingStyle.Render("Python Version Management") + "\n")
	b.WriteString(cmdStyle.Render("  uv python install 3.13") + "          Install a Python version\n")
	b.WriteString(cmdStyle.Render("  uv python list [--only-installed]") + "  List available/installed versions\n")
	b.WriteString(divider + "\n")

	// --- Global (Default) Venv ---
	b.WriteString(headingStyle.Render("Global (Default) Venv — ~/shell/py-default") + "\n")
	b.WriteString(noteStyle.Render("  CPS sets VIRTUAL_ENV=$HOME/shell/py-default in .zshrc") + "\n")
	b.WriteString(noteStyle.Render("  Commands below use this venv when not inside a project.") + "\n\n")
	b.WriteString(cmdStyle.Render("  uv pip install <pkg>") + "            Install into default venv\n")
	b.WriteString(cmdStyle.Render("  uv pip install -r reqs.txt") + "      Install from requirements file\n")
	b.WriteString(cmdStyle.Render("  uv pip list") + "                     List packages in default venv\n")
	b.WriteString(noteStyle.Render("  Also: uv pip show, freeze, uninstall") + "\n")
	b.WriteString(divider + "\n")

	// --- Project Venvs ---
	b.WriteString(headingStyle.Render("Project Venvs — Local .venv + pyproject.toml") + "\n")
	b.WriteString(noteStyle.Render("  uv project commands auto-discover pyproject.toml in cwd.") + "\n")
	b.WriteString(noteStyle.Render("  They create/use .venv/ in the project dir, NOT py-default.") + "\n")
	b.WriteString(noteStyle.Render("  Bare 'python' still points to py-default (VIRTUAL_ENV).") + "\n\n")
	b.WriteString(cmdStyle.Render("  uv init") + "                         Scaffold new pyproject.toml\n")
	b.WriteString(cmdStyle.Render("  uv venv [--python 3.12]") + "         Create .venv in current dir\n")
	b.WriteString(cmdStyle.Render("  uv add [--dev] <pkg>") + "            Add dependency + install\n")
	b.WriteString(cmdStyle.Render("  uv sync") + "                         Sync .venv with lockfile\n")
	b.WriteString(cmdStyle.Render("  uv run <cmd>") + "                    Run command in project venv\n")
	b.WriteString(divider + "\n")

	// --- Tool Install ---
	b.WriteString(headingStyle.Render("Tool Install — Isolated CLI Tools") + "\n")
	b.WriteString(noteStyle.Render("  Each tool gets its own hidden venv (~/.local/share/uv/tools/<name>/).") + "\n")
	b.WriteString(noteStyle.Render("  CLI entry points are symlinked to ~/.local/bin/ (or --bin-dir).") + "\n")
	b.WriteString(noteStyle.Render("  Only works for packages with console_scripts entry points.") + "\n\n")
	b.WriteString(cmdStyle.Render("  uv tool install <pkg>") + "           Install from PyPI\n")
	b.WriteString(cmdStyle.Render("  uv tool install <pkg> --bin-dir <path>") + "\n")
	b.WriteString("                                    " + "Install with custom bin dir\n")
	b.WriteString(cmdStyle.Render("  uv tool install git+https://github.com/user/repo.git") + "\n")
	b.WriteString("                                    " + "Install from git (needs entry points)\n")
	b.WriteString(cmdStyle.Render("  uv tool list") + "                    List installed tools\n")
	b.WriteString(divider + "\n")

	// --- One-Off Script Runs ---
	b.WriteString(headingStyle.Render("One-Off Execution") + "\n")
	b.WriteString(cmdStyle.Render("  uv run --with <pkg> script.py") + "   Run script with temp dependency\n")
	b.WriteString(cmdStyle.Render("  uvx <pkg>") + "                       Run a CLI tool without installing\n")
	b.WriteString(divider + "\n")

	// --- Key Concepts ---
	b.WriteString(headingStyle.Render("Key Concepts") + "\n")
	b.WriteString("  • " + cmdStyle.Render("VIRTUAL_ENV") + " env var  →  controls where bare " + cmdStyle.Render("python/pip") + " resolve\n")
	b.WriteString("  • " + cmdStyle.Render("uv pip ...") + "           →  respects VIRTUAL_ENV (uses py-default)\n")
	b.WriteString("  • " + cmdStyle.Render("uv add/sync/run") + "     →  uses project .venv (ignores VIRTUAL_ENV)\n")
	b.WriteString("  • " + cmdStyle.Render("uv tool install") + "     →  isolated per-tool venv, symlinks to bin dir\n")
	b.WriteString("  • " + cmdStyle.Render("uvx") + "                 →  ephemeral run, nothing persisted\n")

	return b.String()
}

var uvSheet = Sheet{
	Name:        "uv",
	Aliases:     []string{"python", "py"},
	Description: "UV package manager and Python environment cheat sheet",
	Content:     buildUVSheet(),
}
