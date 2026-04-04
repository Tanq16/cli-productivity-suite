package cheatsheet

import "strings"

func buildRustSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("Rust Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  CPS paths: RUSTUP_HOME=~/shell/rust/.rustup") + "\n")
	b.WriteString(noteStyle.Render("  CARGO_HOME=~/shell/rust/.cargo (binaries in .cargo/bin/)") + "\n\n")

	// --- Rustup ---
	b.WriteString(headingStyle.Render("Rustup — Toolchain Management") + "\n")
	b.WriteString(cmdStyle.Render("  rustup update") + "                   Update all toolchains\n")
	b.WriteString(cmdStyle.Render("  rustup show") + "                     Show installed toolchains\n")
	b.WriteString(cmdStyle.Render("  rustup target add <triple>") + "      Add compilation target\n")
	b.WriteString(cmdStyle.Render("  rustup component add clippy") + "     Add a component\n")
	b.WriteString(divider + "\n")

	// --- Cargo Basics ---
	b.WriteString(headingStyle.Render("Cargo — Project Basics") + "\n")
	b.WriteString(cmdStyle.Render("  cargo new [--lib] <name>") + "        Create binary or library project\n")
	b.WriteString(cmdStyle.Render("  cargo build [--release]") + "         Build debug or release binary\n")
	b.WriteString(cmdStyle.Render("  cargo run [-- args...]") + "          Build and run\n")
	b.WriteString(cmdStyle.Render("  cargo check") + "                     Type-check without building\n")
	b.WriteString(divider + "\n")

	// --- Dependencies ---
	b.WriteString(headingStyle.Render("Cargo — Dependencies") + "\n")
	b.WriteString(cmdStyle.Render("  cargo add [--dev] [--features f1,f2] <crate>") + "\n")
	b.WriteString("                                    " + "Add dependency\n")
	b.WriteString(cmdStyle.Render("  cargo update") + "                    Update Cargo.lock\n")
	b.WriteString(cmdStyle.Render("  cargo tree [-d]") + "                 Show dependency tree [-d for dupes]\n")
	b.WriteString(divider + "\n")

	// --- Testing ---
	b.WriteString(headingStyle.Render("Cargo — Testing & Linting") + "\n")
	b.WriteString(cmdStyle.Render("  cargo test [<name>]") + "             Run all or matching tests\n")
	b.WriteString(cmdStyle.Render("  cargo test -- --nocapture") + "       Show stdout in tests\n")
	b.WriteString(cmdStyle.Render("  cargo clippy") + "                    Run linter\n")
	b.WriteString(cmdStyle.Render("  cargo fmt [--check]") + "             Format code [or check only]\n")
	b.WriteString(divider + "\n")

	// --- Install ---
	b.WriteString(headingStyle.Render("Cargo — Install") + "\n")
	b.WriteString(cmdStyle.Render("  cargo install <crate>") + "           Install binary from crates.io\n")
	b.WriteString(cmdStyle.Render("  cargo install --git <url>") + "       Install from git repo\n")
	b.WriteString(divider + "\n")

	// --- Cross Compilation ---
	b.WriteString(headingStyle.Render("Cross Compilation") + "\n")
	b.WriteString(cmdStyle.Render("  rustup target add x86_64-unknown-linux-musl") + "\n")
	b.WriteString("                                    " + "Add musl target for static binaries\n")
	b.WriteString(cmdStyle.Render("  cargo build --target <triple> --release") + "\n")
	b.WriteString("                                    " + "Cross-compile release binary\n")
	b.WriteString(noteStyle.Render("  Output at: target/<triple>/release/<binary>") + "\n")

	return b.String()
}

var rustSheet = Sheet{
	Name:        "rust",
	Aliases:     []string{"cargo", "rustup"},
	Description: "Rust toolchain, Cargo, and cross-compilation cheat sheet",
	Content:     buildRustSheet(),
}
