package cheatsheet

import "strings"

func buildRegexSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("String Processing Cheat Sheet — grep, ripgrep, awk") + "\n")

	// --- Ripgrep ---
	b.WriteString(headingStyle.Render("Ripgrep (rg)") + "\n")
	b.WriteString(cmdStyle.Render("  rg <pattern>") + "                    Search recursively\n")
	b.WriteString(cmdStyle.Render("  rg -i <pattern>") + "                 Case insensitive\n")
	b.WriteString(cmdStyle.Render("  rg -w <pattern>") + "                 Whole word match\n")
	b.WriteString(cmdStyle.Render("  rg -l <pattern>") + "                 List matching files only\n")
	b.WriteString(cmdStyle.Render("  rg -c <pattern>") + "                 Count matches per file\n")
	b.WriteString(cmdStyle.Render("  rg -v <pattern>") + "                 Invert match\n")
	b.WriteString(cmdStyle.Render("  rg -n <pattern>") + "                 Show line numbers (default)\n")
	b.WriteString(cmdStyle.Render("  rg -A3 -B3 <pattern>") + "            Show 3 lines context\n")
	b.WriteString(cmdStyle.Render("  rg -C3 <pattern>") + "                Same as -A3 -B3\n")
	b.WriteString(cmdStyle.Render("  rg -t py <pattern>") + "              Search only Python files\n")
	b.WriteString(cmdStyle.Render("  rg -g '*.go' <pattern>") + "          Search matching glob\n")
	b.WriteString(cmdStyle.Render("  rg -g '!vendor' <pattern>") + "       Exclude paths matching glob\n")
	b.WriteString(cmdStyle.Render("  rg --json <pattern>") + "             Output as JSON\n")
	b.WriteString(cmdStyle.Render("  rg -o <pattern>") + "                 Print only matched parts\n")
	b.WriteString(cmdStyle.Render("  rg -r '<replacement>' <pattern>") + "  Search and replace (stdout)\n")
	b.WriteString(cmdStyle.Render("  rg -U <pattern>") + "                 Multiline search\n")
	b.WriteString(cmdStyle.Render("  rg --files") + "                      List files rg would search\n")
	b.WriteString(divider + "\n")

	// --- Grep ---
	b.WriteString(headingStyle.Render("Grep — When rg Isn't Available") + "\n")
	b.WriteString(noteStyle.Render("  rg replaces grep for daily use. Grep for remote boxes or PCRE.") + "\n\n")
	b.WriteString(cmdStyle.Render("  grep -P <pattern> <file>") + "        Perl regex (backrefs \\1, lookbehind)\n")
	b.WriteString(cmdStyle.Render("  grep -r <pattern> <dir>") + "         Recursive (no .gitignore filtering)\n")
	b.WriteString(cmdStyle.Render("  <cmd> | grep <pattern>") + "          Stdin filtering in pipes\n")
	b.WriteString(divider + "\n")

	// --- Regex Syntax ---
	b.WriteString(headingStyle.Render("Regex Quick Reference") + "\n")
	b.WriteString(cmdStyle.Render("  .") + "    Any char        " + cmdStyle.Render("*") + "    Zero or more    " + cmdStyle.Render("+") + "    One or more\n")
	b.WriteString(cmdStyle.Render("  ?") + "    Zero or one     " + cmdStyle.Render("^") + "    Start of line   " + cmdStyle.Render("$") + "    End of line\n")
	b.WriteString(cmdStyle.Render("  \\d") + "   Digit [0-9]     " + cmdStyle.Render("\\w") + "   Word char       " + cmdStyle.Render("\\s") + "   Whitespace\n")
	b.WriteString(cmdStyle.Render("  \\b") + "   Word boundary   " + cmdStyle.Render("{n}") + "  Exactly n       " + cmdStyle.Render("{n,m}") + " n to m\n")
	b.WriteString(cmdStyle.Render("  [abc]") + " Char class      " + cmdStyle.Render("[^ab]") + " Negated class   " + cmdStyle.Render("[a-z]") + " Range\n")
	b.WriteString(cmdStyle.Render("  (a|b)") + " Alternation     " + cmdStyle.Render("(?:)") + "  Non-capturing   " + cmdStyle.Render("(?=)") + "  Lookahead\n")
	b.WriteString(divider + "\n")

	// --- Awk ---
	b.WriteString(headingStyle.Render("Awk") + "\n")
	b.WriteString(noteStyle.Render("  awk '<pattern> { <action> }' — runs action on matching lines") + "\n\n")
	b.WriteString(cmdStyle.Render("  awk '{print}' file") + "              Print all lines\n")
	b.WriteString(cmdStyle.Render("  awk '{print $1}' file") + "           Print first field\n")
	b.WriteString(cmdStyle.Render("  awk '{print $1, $3}' file") + "       Print fields 1 and 3\n")
	b.WriteString(cmdStyle.Render("  awk -F: '{print $1}' file") + "       Use : as delimiter\n")
	b.WriteString(cmdStyle.Render("  awk -F'\\t' '{print $2}' file") + "    Use tab as delimiter\n")
	b.WriteString(cmdStyle.Render("  awk '/regex/ {print}' file") + "      Print matching lines\n")
	b.WriteString(cmdStyle.Render("  awk '$3 > 100' file") + "             Filter by field value\n")
	b.WriteString(cmdStyle.Render("  awk 'NR==5' file") + "                Print line 5\n")
	b.WriteString(cmdStyle.Render("  awk 'NR>=5 && NR<=10' file") + "      Print lines 5-10\n")
	b.WriteString(cmdStyle.Render("  awk '{sum+=$1} END{print sum}'") + "   Sum first field\n")
	b.WriteString(cmdStyle.Render("  awk '{print NR, $0}' file") + "       Add line numbers\n")
	b.WriteString(cmdStyle.Render("  awk '!seen[$0]++' file") + "          Remove duplicate lines\n")
	b.WriteString(cmdStyle.Render("  awk '{gsub(/old/,\"new\"); print}'") + "  Find and replace\n")
	b.WriteString(divider + "\n")

	// --- Awk Variables ---
	b.WriteString(headingStyle.Render("Awk Built-in Variables") + "\n")
	b.WriteString(cmdStyle.Render("  $0") + "        Entire line             ")
	b.WriteString(cmdStyle.Render("  $1..$n") + "    Field n\n")
	b.WriteString(cmdStyle.Render("  NR") + "        Current line number     ")
	b.WriteString(cmdStyle.Render("  NF") + "        Number of fields\n")
	b.WriteString(cmdStyle.Render("  FS") + "        Field separator         ")
	b.WriteString(cmdStyle.Render("  OFS") + "       Output field separator\n")
	b.WriteString(cmdStyle.Render("  RS") + "        Record separator        ")
	b.WriteString(cmdStyle.Render("  FILENAME") + "  Current filename\n")

	return b.String()
}

var regexSheet = Sheet{
	Name:        "regex",
	Aliases:     []string{"grep", "rg", "ripgrep", "awk", "strings"},
	Description: "String processing with grep, ripgrep, awk, and regex syntax cheat sheet",
	Content:     buildRegexSheet(),
}
