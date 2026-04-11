package cheatsheet

import "strings"

func buildGoSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("Go Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  CPS paths: GOROOT=~/shell/go-sdk  GOPATH=~/shell/go") + "\n")
	b.WriteString(noteStyle.Render("  Binaries from go install land in $GOPATH/bin/") + "\n\n")

	// --- Build & Run ---
	b.WriteString(headingStyle.Render("Build & Run") + "\n")
	b.WriteString(cmdStyle.Render("  go build ./...") + "                  Build all packages\n")
	b.WriteString(cmdStyle.Render("  go build -o <name> .") + "            Build with custom output name\n")
	b.WriteString(cmdStyle.Render("  go run .") + "                        Build and run main package\n")
	b.WriteString(divider + "\n")

	// --- Testing ---
	b.WriteString(headingStyle.Render("Testing") + "\n")
	b.WriteString(cmdStyle.Render("  go test ./...") + "                   Run all tests\n")
	b.WriteString(cmdStyle.Render("  go test -v -run <regex> ./...") + "   Verbose + filter by name\n")
	b.WriteString(cmdStyle.Render("  go test -cover -race ./...") + "      Coverage + race detector\n")
	b.WriteString(divider + "\n")

	// --- Tools ---
	b.WriteString(headingStyle.Render("Tools") + "\n")
	b.WriteString(cmdStyle.Render("  go vet ./...") + "                    Report likely mistakes\n")
	b.WriteString(cmdStyle.Render("  gofmt -w .") + "                      Format all files in place\n")
	b.WriteString(cmdStyle.Render("  go install <pkg>@latest") + "         Install a Go binary tool\n")
	b.WriteString(divider + "\n")

	// --- Go Tool ---
	b.WriteString(headingStyle.Render("Go Tool — Built-in Analyzers") + "\n")
	b.WriteString(cmdStyle.Render("  go tool pprof <profile>") + "         Analyze CPU/memory profiles\n")
	b.WriteString(cmdStyle.Render("  go tool nm <binary>") + "             List symbols in a binary\n")
	b.WriteString(cmdStyle.Render("  go tool objdump <binary>") + "        Disassemble a binary\n")
	b.WriteString(cmdStyle.Render("  go tool bisect") + "                  Binary search for failure cause\n")
	b.WriteString(noteStyle.Render("  External tools (install first):") + "\n")
	b.WriteString(cmdStyle.Render("  go install golang.org/x/tools/cmd/deadcode@latest") + "\n")
	b.WriteString(cmdStyle.Render("  go tool deadcode ./...") + "          Find unreachable functions\n")
	b.WriteString(divider + "\n")

	// --- Cross Compilation ---
	b.WriteString(headingStyle.Render("Cross Compilation") + "\n")
	b.WriteString(cmdStyle.Render("  GOOS=linux GOARCH=amd64 go build -o app .") + "\n")
	b.WriteString(cmdStyle.Render("  GOOS=darwin GOARCH=arm64 go build -o app .") + "\n")
	b.WriteString(noteStyle.Render("  Common: GOOS=linux|darwin|windows  GOARCH=amd64|arm64") + "\n")
	b.WriteString(divider + "\n")

	// --- Build Flags ---
	b.WriteString(headingStyle.Render("Build Flags") + "\n")
	b.WriteString(cmdStyle.Render("  -ldflags \"-s -w\"") + "                Strip debug info (smaller binary)\n")
	b.WriteString(cmdStyle.Render("  -ldflags \"-X main.ver=v1\"") + "      Inject version at build time\n")
	b.WriteString(cmdStyle.Render("  CGO_ENABLED=0") + "                   Static binary (no C deps)\n")

	return b.String()
}

var goSheet = Sheet{
	Name:        "go",
	Aliases:     []string{"golang"},
	Description: "Go build, test, tools, and cross-compilation cheat sheet",
	Content:     buildGoSheet(),
}
