package cheatsheet

import "strings"

func buildJqSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("jq Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  jq is a streaming JSON processor — filters are composed with | pipes") + "\n")
	b.WriteString(noteStyle.Render("  Reads JSON from stdin or files, writes JSON to stdout") + "\n\n")

	// --- Invocation Flags ---
	b.WriteString(headingStyle.Render("Invocation Flags") + "\n")
	b.WriteString(cmdStyle.Render("  jq '.' file.json") + "                Pretty-print JSON\n")
	b.WriteString(cmdStyle.Render("  jq -r '.'") + "                       Raw output — strip quotes from strings\n")
	b.WriteString(cmdStyle.Render("  jq -c '.'") + "                       Compact output (one line per value)\n")
	b.WriteString(cmdStyle.Render("  jq -s '.'") + "                       Slurp — read entire input into a single array\n")
	b.WriteString(cmdStyle.Render("  jq -n '<expr>'") + "                  No input — useful with --arg for building JSON\n")
	b.WriteString(cmdStyle.Render("  jq -R '.'") + "                       Read raw input (lines as strings, not JSON)\n")
	b.WriteString(cmdStyle.Render("  jq -e '<filter>'") + "                Exit non-zero if result is null/false (scripting)\n")
	b.WriteString(cmdStyle.Render("  jq --arg name val '<expr>'") + "      Inject string variable $name\n")
	b.WriteString(cmdStyle.Render("  jq --argjson name val '<expr>'") + "   Inject JSON variable $name (number/object/etc.)\n")
	b.WriteString(divider + "\n")

	// --- Basic Selection ---
	b.WriteString(headingStyle.Render("Selection") + "\n")
	b.WriteString(cmdStyle.Render("  .") + "                               The whole input\n")
	b.WriteString(cmdStyle.Render("  .foo") + "                            Field access\n")
	b.WriteString(cmdStyle.Render("  .foo.bar") + "                        Nested field\n")
	b.WriteString(cmdStyle.Render("  .foo?") + "                           Optional — null instead of error if missing\n")
	b.WriteString(cmdStyle.Render("  .[\"foo bar\"]") + "                    Field name with spaces / special chars\n")
	b.WriteString(cmdStyle.Render("  .[0]") + "                            Array index (negative allowed: .[-1])\n")
	b.WriteString(cmdStyle.Render("  .[]") + "                             Iterate array or object values (stream)\n")
	b.WriteString(cmdStyle.Render("  .[2:5]") + "                          Array slice\n")
	b.WriteString(cmdStyle.Render("  .foo // \"default\"") + "               Default if null/false\n")
	b.WriteString(divider + "\n")

	// --- Filtering & Mapping ---
	b.WriteString(headingStyle.Render("Filtering & Mapping") + "\n")
	b.WriteString(cmdStyle.Render("  map(.name)") + "                      Apply filter to each array element\n")
	b.WriteString(cmdStyle.Render("  .users[] | .email") + "               Stream array, project field\n")
	b.WriteString(cmdStyle.Render("  .users | map(select(.active))") + "    Filter array by predicate\n")
	b.WriteString(cmdStyle.Render("  select(.age > 21)") + "               Keep values matching predicate\n")
	b.WriteString(cmdStyle.Render("  select(.tags | contains([\"x\"]))") + "  Predicate on nested array\n")
	b.WriteString(cmdStyle.Render("  recurse | select(.id?)") + "          Recursive descent — find any node with .id\n")
	b.WriteString(cmdStyle.Render("  .. | .name? // empty") + "            Same idea, idiomatic deep scan for .name\n")
	b.WriteString(divider + "\n")

	// --- Constructing Output ---
	b.WriteString(headingStyle.Render("Constructing Output") + "\n")
	b.WriteString(cmdStyle.Render("  {name: .name, id: .id}") + "          Build an object\n")
	b.WriteString(cmdStyle.Render("  {name, id}") + "                      Shorthand — same as above\n")
	b.WriteString(cmdStyle.Render("  [.[] | .id]") + "                     Collect stream into array\n")
	b.WriteString(cmdStyle.Render("  to_entries") + "                      Object → array of {key, value}\n")
	b.WriteString(cmdStyle.Render("  from_entries") + "                    Array of {key, value} → object\n")
	b.WriteString(cmdStyle.Render("  with_entries(.value += 1)") + "       Transform each entry of an object\n")
	b.WriteString(cmdStyle.Render("  paths") + "                           Stream every path in the input\n")
	b.WriteString(cmdStyle.Render("  leaf_paths") + "                      Stream paths to scalar values only\n")
	b.WriteString(divider + "\n")

	// --- Aggregations ---
	b.WriteString(headingStyle.Render("Aggregation") + "\n")
	b.WriteString(cmdStyle.Render("  length") + "                          Length of array/string/object\n")
	b.WriteString(cmdStyle.Render("  keys") + "                            Sorted keys (object) / indices (array)\n")
	b.WriteString(cmdStyle.Render("  values") + "                          Values of an object\n")
	b.WriteString(cmdStyle.Render("  add") + "                             Sum / concat / merge across stream\n")
	b.WriteString(cmdStyle.Render("  min, max, min_by(.x), max_by(.x)") + "  Numeric / keyed min/max\n")
	b.WriteString(cmdStyle.Render("  unique, unique_by(.x)") + "           Dedupe array\n")
	b.WriteString(cmdStyle.Render("  sort, sort_by(.x)") + "               Sort array\n")
	b.WriteString(cmdStyle.Render("  group_by(.x)") + "                    Bucket array by key\n")
	b.WriteString(cmdStyle.Render("  reduce .[] as $x (0; . + $x.cost)") + " Fold a stream\n")
	b.WriteString(divider + "\n")

	// --- Strings & Types ---
	b.WriteString(headingStyle.Render("Strings & Type Helpers") + "\n")
	b.WriteString(cmdStyle.Render("  tostring, tonumber") + "              Type coercion\n")
	b.WriteString(cmdStyle.Render("  type") + "                            \"string\" | \"number\" | \"object\" | ...\n")
	b.WriteString(cmdStyle.Render("  ascii_downcase, ascii_upcase") + "    Case conversion\n")
	b.WriteString(cmdStyle.Render("  split(\",\"), join(\",\")") + "          String ↔ array\n")
	b.WriteString(cmdStyle.Render("  test(\"^abc\")") + "                    Regex test → bool\n")
	b.WriteString(cmdStyle.Render("  match(\"<re>\")") + "                   Regex match → object with offsets\n")
	b.WriteString(cmdStyle.Render("  sub(\"<re>\"; \"<repl>\")") + "          Regex replace (first match)\n")
	b.WriteString(cmdStyle.Render("  gsub(\"<re>\"; \"<repl>\")") + "         Regex replace (all matches)\n")
	b.WriteString(cmdStyle.Render("  @csv, @tsv, @json, @sh, @uri") + "    Format output (use with -r)\n")
	b.WriteString(divider + "\n")

	// --- Common Recipes ---
	b.WriteString(headingStyle.Render("Common Recipes") + "\n")
	b.WriteString(cmdStyle.Render("  jq -r '.items[].name'") + "                       Extract a flat list of names\n")
	b.WriteString(cmdStyle.Render("  jq '.items | length'") + "                        Count items\n")
	b.WriteString(cmdStyle.Render("  jq '.items | map(select(.tag==\"a\"))'") + "       Filter array by field\n")
	b.WriteString(cmdStyle.Render("  jq -r '.[] | [.id,.name] | @csv'") + "            Emit CSV\n")
	b.WriteString(cmdStyle.Render("  jq -s 'add'") + "                                 Merge a stream of arrays/objects\n")
	b.WriteString(cmdStyle.Render("  jq 'del(.password)'") + "                         Drop a field\n")
	b.WriteString(cmdStyle.Render("  jq '. + {added: true}'") + "                      Merge / patch an object\n")
	b.WriteString(cmdStyle.Render("  jq --arg q \"$Q\" '.[] | select(.name==$q)'") + "  Shell-var injection (safe)\n")
	b.WriteString(cmdStyle.Render("  curl -s URL | jq '.data[] | {id,name}'") + "      Pipe from curl\n")
	b.WriteString(divider + "\n")

	// --- Key Concepts ---
	b.WriteString(headingStyle.Render("Key Concepts") + "\n")
	b.WriteString("  • Filters are " + cmdStyle.Render("composed with |") + " — each filter consumes the previous output\n")
	b.WriteString("  • " + cmdStyle.Render(".[]") + " produces a stream; wrap in " + cmdStyle.Render("[ ... ]") + " to collect back into an array\n")
	b.WriteString("  • " + cmdStyle.Render("-r") + " strips JSON quoting — almost always what you want for shell pipelines\n")
	b.WriteString("  • Use " + cmdStyle.Render("--arg") + " / " + cmdStyle.Render("--argjson") + " instead of string-interpolating shell vars into the filter\n")
	b.WriteString("  • " + cmdStyle.Render("?") + " suppresses errors on missing fields; " + cmdStyle.Render("//") + " supplies defaults on null\n")

	return b.String()
}

var jqSheet = Sheet{
	Name:        "jq",
	Aliases:     []string{},
	Description: "jq JSON processor — filters, selection, recipes",
	Content:     buildJqSheet(),
}
