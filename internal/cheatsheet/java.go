package cheatsheet

import "strings"

func buildJavaSheet() string {
	var b strings.Builder
	divider := dividerStyle.Render(strings.Repeat("─", 60))

	b.WriteString(titleStyle.Render("Java Cheat Sheet") + "\n")
	b.WriteString(noteStyle.Render("  CPS paths: JAVA_HOME=~/shell/java-sdk  (Eclipse Temurin LTS)") + "\n")
	b.WriteString(noteStyle.Render("  $JAVA_HOME/bin is on PATH — java, javac, jar, jshell, jlink, etc.") + "\n\n")

	// --- Compile & Run ---
	b.WriteString(headingStyle.Render("Compile & Run") + "\n")
	b.WriteString(cmdStyle.Render("  javac Foo.java") + "                  Compile to Foo.class\n")
	b.WriteString(cmdStyle.Render("  java Foo") + "                        Run class (no .class suffix)\n")
	b.WriteString(cmdStyle.Render("  java Foo.java") + "                   Compile + run single file (JDK 11+)\n")
	b.WriteString(cmdStyle.Render("  java -jar app.jar") + "               Run executable JAR\n")
	b.WriteString(cmdStyle.Render("  java -cp lib/*:. Foo") + "            Run with classpath\n")
	b.WriteString(divider + "\n")

	// --- REPL & Scripting ---
	b.WriteString(headingStyle.Render("REPL & Scripting") + "\n")
	b.WriteString(cmdStyle.Render("  jshell") + "                          Interactive Java REPL\n")
	b.WriteString(cmdStyle.Render("  jshell script.jsh") + "               Run a jshell script\n")
	b.WriteString(cmdStyle.Render("  /vars  /methods  /imports") + "       Inspect session inside jshell\n")
	b.WriteString(cmdStyle.Render("  /exit") + "                           Leave jshell\n")
	b.WriteString(divider + "\n")

	// --- Packaging ---
	b.WriteString(headingStyle.Render("Packaging") + "\n")
	b.WriteString(cmdStyle.Render("  jar cf app.jar -C classes/ .") + "    Create JAR from classes/\n")
	b.WriteString(cmdStyle.Render("  jar cfe app.jar Main -C classes/ .") + "  Create + set Main-Class\n")
	b.WriteString(cmdStyle.Render("  jar tf app.jar") + "                  List JAR contents\n")
	b.WriteString(cmdStyle.Render("  jar xf app.jar") + "                  Extract JAR\n")
	b.WriteString(cmdStyle.Render("  jlink --module-path $JAVA_HOME/jmods \\\n         --add-modules java.base --output runtime") + "\n")
	b.WriteString(noteStyle.Render("    jlink builds a minimal custom JRE image") + "\n")
	b.WriteString(divider + "\n")

	// --- Inspection & Diagnostics ---
	b.WriteString(headingStyle.Render("Inspection & Diagnostics") + "\n")
	b.WriteString(cmdStyle.Render("  javap -p -c Foo.class") + "           Disassemble bytecode\n")
	b.WriteString(cmdStyle.Render("  jdeps app.jar") + "                   Show class dependencies\n")
	b.WriteString(cmdStyle.Render("  jps -l") + "                          List running JVM PIDs\n")
	b.WriteString(cmdStyle.Render("  jstack <pid>") + "                    Thread dump\n")
	b.WriteString(cmdStyle.Render("  jmap -histo <pid>") + "               Heap histogram\n")
	b.WriteString(cmdStyle.Render("  jcmd <pid> GC.heap_dump dump.hprof") + "  Heap dump for analysis\n")
	b.WriteString(divider + "\n")

	// --- JVM Flags ---
	b.WriteString(headingStyle.Render("Common JVM Flags") + "\n")
	b.WriteString(cmdStyle.Render("  -Xms512m -Xmx2g") + "                 Initial / max heap\n")
	b.WriteString(cmdStyle.Render("  -XX:+UseZGC") + "                     Use ZGC (low-pause)\n")
	b.WriteString(cmdStyle.Render("  -XX:+HeapDumpOnOutOfMemoryError") + " Dump heap on OOM\n")
	b.WriteString(cmdStyle.Render("  -Dkey=value") + "                     Set system property\n")
	b.WriteString(cmdStyle.Render("  -ea") + "                             Enable assertions\n")
	b.WriteString(divider + "\n")

	// --- Build Tools ---
	b.WriteString(headingStyle.Render("Build Tools (install separately)") + "\n")
	b.WriteString(cmdStyle.Render("  brew install maven") + "              Maven (pom.xml)\n")
	b.WriteString(cmdStyle.Render("  brew install gradle") + "             Gradle (build.gradle[.kts])\n")
	b.WriteString(cmdStyle.Render("  mvn package    mvn test") + "         Maven lifecycle\n")
	b.WriteString(cmdStyle.Render("  gradle build   gradle test") + "      Gradle tasks\n")
	b.WriteString(noteStyle.Render("  Wrapper scripts (mvnw / gradlew) ship with most projects — prefer them") + "\n")

	return b.String()
}

var javaSheet = Sheet{
	Name:        "java",
	Aliases:     []string{"jdk"},
	Description: "Java/JDK compile, run, packaging, and diagnostics cheat sheet",
	Content:     buildJavaSheet(),
}
