// Command dev is AstraMind's cross-platform task runner - the single
// source of truth for build/test/coverage/regression/clean commands.
//
// Requires no new toolchain on either platform: Go is already a hard
// requirement to build AstraMind at all. Callable identically from
// bash and cmd.exe:
//
//	go run ./cmd/dev -run=build
//	go run ./cmd/dev -run=test
//	go run ./cmd/dev -run=coverage
//	go run ./cmd/dev -run=regression
//	go run ./cmd/dev -run=clean
//
// Every command flag and path lives here exactly once. There is no
// second file to forget to update - the exact class of bug
// (coverage.sh's stale path after a directory move, regression.bat
// containing an entirely different file's content) already found
// live in this project on 2026-07-25.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	task := flag.String("run", "help", "Automation task to execute: build, test, coverage, regression, clean")
	flag.Parse()

	switch *task {
	case "build":
		executeBuild()
	case "test":
		executeTest()
	case "coverage":
		executeCoverage()
	case "regression":
		executeRegression()
	case "clean":
		executeClean()
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("AstraMind Cross-Platform Development Automation Runner")
	fmt.Println("Usage: go run ./cmd/dev -run=[task]")
	fmt.Println()
	fmt.Println("Tasks:")
	fmt.Println("  build       Compile astramind(.exe) to the repo root")
	fmt.Println("  test        Run the full unit test suite")
	fmt.Println("  coverage    Generate coverage reports in tests/output/coverage/")
	fmt.Println("  regression  Run the full regression suite (calls scripts/regression.sh or .bat)")
	fmt.Println("  clean       Remove build artifacts and Go caches")
}

func executeCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Task execution failed during [%s %v]: %v", name, args, err)
	}
}

// runCapture runs a command and returns its combined stdout+stderr as
// a string, instead of streaming it live - used where the output
// needs to be parsed afterward (coverage summaries), not just shown.
func runCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func binaryName() string {
	if runtime.GOOS == "windows" {
		return "astramind.exe"
	}
	return "astramind"
}

func executeBuild() {
	fmt.Println("[Build] Formatting...")
	executeCommand("go", "fmt", "./...")

	fmt.Println("[Build] Static Analysis...")
	executeCommand("go", "vet", "./...")

	fmt.Println("[Build] Building...")
	// "./cmd/astramind" (the package directory), not
	// "./cmd/astramind/main.go" (one file) - naming a single file
	// silently excludes any other .go file added to that directory
	// later, with no error, just missing code in the built binary.
	executeCommand("go", "build", "-ldflags=-s -w", "-o", binaryName(), "./cmd/astramind")

	fmt.Println("[Build] SUCCESS")
}

func executeTest() {
	fmt.Println("[Test] Running Unit & Integration Tests...")
	// No -race here deliberately: the race detector requires CGO,
	// which requires a real C compiler - not guaranteed present on a
	// Windows dev machine, and would defeat the entire point of this
	// runner (avoiding new toolchain requirements). This does NOT set
	// CGO_ENABLED=0 explicitly; it simply doesn't need CGO for a
	// normal test run, since only -race triggers that requirement for
	// a pure-Go codebase like this one.
	executeCommand("go", "test", "-v", "./...")
	fmt.Println("[Test] SUCCESS")
}

// executeCoverage matches scripts/coverage.sh's real, full output -
// not just a .out/.html pair. Text processing is done with Go's
// standard library (bufio/strings), not by shelling out to grep/sed -
// those aren't available on stock Windows, and depending on them here
// would quietly reintroduce the exact kind of platform-specific
// dependency this runner exists to remove.
func executeCoverage() {
	fmt.Println("====================================")
	fmt.Println("AstraMind Coverage Report")
	fmt.Println("====================================")

	coverageDir := "tests/output/coverage"
	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		log.Fatalf("Failed to create %s: %v", coverageDir, err)
	}

	outFile := coverageDir + "/coverage.out"
	htmlFile := coverageDir + "/coverage.html"
	txtFile := coverageDir + "/coverage.txt"
	packageFile := coverageDir + "/package_coverage.txt"

	for _, f := range []string{outFile, htmlFile, txtFile, packageFile} {
		os.Remove(f) // ignore error - fine if it didn't exist
	}

	packageOutput, err := runCapture("go", "test", "-coverprofile="+outFile, "./...")
	if err != nil {
		fmt.Println(packageOutput)
		log.Fatalf("go test failed: %v", err)
	}
	if err := os.WriteFile(packageFile, []byte(packageOutput), 0644); err != nil {
		log.Fatalf("failed to write %s: %v", packageFile, err)
	}

	fmt.Println()
	fmt.Println("Test Coverage Summary")
	printCoverageSummary(packageOutput)

	executeCommand("go", "tool", "cover", "-html="+outFile, "-o", htmlFile)

	funcOutput, err := runCapture("go", "tool", "cover", "-func="+outFile)
	if err != nil {
		fmt.Println(funcOutput)
		log.Fatalf("go tool cover -func failed: %v", err)
	}
	if err := os.WriteFile(txtFile, []byte(funcOutput), 0644); err != nil {
		log.Fatalf("failed to write %s: %v", txtFile, err)
	}

	total := extractTotalCoverage(funcOutput)
	fmt.Println()
	fmt.Println("Overall Coverage :", total)
	fmt.Println()
	fmt.Println("Coverage report generated.")
}

// printCoverageSummary replicates coverage.sh's:
//
//	grep "coverage:" package_coverage.txt | grep -v "cmd/astramind" | sed ...
//
// using bufio/strings instead of shelling out.
func printCoverageSummary(packageOutput string) {
	scanner := bufio.NewScanner(strings.NewReader(packageOutput))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "coverage:") {
			continue
		}
		if strings.Contains(line, "cmd/astramind") {
			continue
		}
		line = strings.TrimPrefix(line, "ok")
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "(cached)", "")
		fmt.Println(line)
	}
}

// extractTotalCoverage replicates:
//
//	grep "^total:" coverage.txt | awk '{print $3}'
func extractTotalCoverage(funcOutput string) string {
	scanner := bufio.NewScanner(strings.NewReader(funcOutput))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "total:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			return fields[2]
		}
	}
	return "unknown"
}

func executeRegression() {
	fmt.Println("Initializing platform-native regression orchestration...")
	if runtime.GOOS == "windows" {
		executeCommand("cmd", "/c", "scripts\\regression.bat")
	} else {
		// Invoked via bash explicitly, not executed directly - this
		// doesn't require the script's own execute bit to be set,
		// which can legitimately be lost on some checkouts/copies.
		// Mirrors the Windows branch above, which already goes
		// through cmd /c rather than trying to run the .bat directly.
		executeCommand("bash", "./scripts/regression.sh")
	}
}

func executeClean() {
	fmt.Println("Purging build and profile artifacts...")
	executeCommand("go", "clean", "-cache", "-testcache")

	artifacts := []string{
		"astramind",
		"astramind.exe",
		"tests/output/coverage/coverage.out",
		"tests/output/coverage/coverage.html",
		"tests/output/coverage/coverage.txt",
		"tests/output/coverage/package_coverage.txt",
	}
	for _, f := range artifacts {
		if err := os.Remove(f); err == nil {
			fmt.Println("Removed:", f)
		}
	}
}
