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
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	task := flag.String("run", "help", "Automation task to execute: build, test, coverage, junit, regression-report, regression, clean")
	buildStatus := flag.String("build", "SKIPPED", "Build step status: PASS, FAIL, or SKIPPED")
	testStatus := flag.String("tests", "SKIPPED", "Tests step status: PASS, FAIL, or SKIPPED")
	coverageStatus := flag.String("coverage-status", "SKIPPED", "Coverage step status: PASS, FAIL, or SKIPPED")
	kbragStatus := flag.String("kbrag", "SKIPPED", "KB & RAG step status: PASS, FAIL, or SKIPPED")
	elapsed := flag.String("elapsed", "0", "Total elapsed seconds for the regression run")
	flag.Parse()

	switch *task {
	case "build":
		executeBuild()
	case "test":
		executeTest()
	case "coverage":
		executeCoverage()
	case "junit":
		executeJUnit()
	case "regression-report":
		executeRegressionReport(*buildStatus, *testStatus, *coverageStatus, *kbragStatus, *elapsed)
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
	fmt.Println("  junit       Generate reports/junit.xml for CI test-result display")
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
		os.Remove(f) //nolint:errcheck // ignore error - fine if it didn't exist
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

// goTestEvent mirrors the JSON Lines format `go test -json` emits -
// one line per event (a test starting, its output, and its final
// pass/fail/skip). This is Go's own, already-structured test output;
// parsing it directly avoids needing an external tool
// (go-junit-report, gotestsum) that isn't part of the Go toolchain
// itself - the same reasoning already applied to coverage.sh's
// grep/sed/awk pipeline, rewritten in Go for the same reason.
type goTestEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
}

// JUnit XML schema structs - the format GitHub Actions' test-reporter
// actions (and most CI systems) already know how to render natively.
type junitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []junitTestSuite `xml:"testsuite"`
}

type junitTestSuite struct {
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Time      string          `xml:"time,attr"`
	TestCases []junitTestCase `xml:"testcase"`
}

type junitTestCase struct {
	ClassName string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *junitFailure `xml:"failure,omitempty"`
	Skipped   *junitSkipped `xml:"skipped,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Content string `xml:",chardata"`
}

type junitSkipped struct {
	Message string `xml:"message,attr"`
}

func executeJUnit() {
	fmt.Println("Generating reports/junit.xml from `go test -json` output...")

	if err := os.MkdirAll("reports", 0755); err != nil {
		log.Fatalf("failed to create reports/ directory: %v", err)
	}

	cmd := exec.Command("go", "test", "-json", "./...")
	output, _ := cmd.CombinedOutput() // test failures are expected; report them, don't abort on them

	type testResult struct {
		name    string
		elapsed float64
		failed  bool
		output  strings.Builder
	}

	suites := map[string][]*testResult{}  // package -> ordered test results
	testIndex := map[string]*testResult{} // "package/test" -> result, for fast lookup while streaming

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		var ev goTestEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue // non-JSON line (shouldn't happen with -json, but don't crash on it)
		}
		if ev.Test == "" {
			continue // package-level event, not an individual test
		}
		key := ev.Package + "/" + ev.Test

		switch ev.Action {
		case "run":
			r := &testResult{name: ev.Test}
			testIndex[key] = r
			suites[ev.Package] = append(suites[ev.Package], r)
		case "output":
			if r, ok := testIndex[key]; ok {
				r.output.WriteString(ev.Output)
			}
		case "pass":
			if r, ok := testIndex[key]; ok {
				r.elapsed = ev.Elapsed
			}
		case "fail":
			if r, ok := testIndex[key]; ok {
				r.elapsed = ev.Elapsed
				r.failed = true
			}
		}
	}

	var report junitTestSuites
	for pkg, results := range suites {
		suite := junitTestSuite{Name: pkg}
		var suiteTime float64
		for _, r := range results {
			tc := junitTestCase{
				ClassName: pkg,
				Name:      r.name,
				Time:      fmt.Sprintf("%.3f", r.elapsed),
			}
			suiteTime += r.elapsed
			if r.failed {
				suite.Failures++
				tc.Failure = &junitFailure{
					Message: "Test failed",
					Content: r.output.String(),
				}
			}
			suite.Tests++
			suite.TestCases = append(suite.TestCases, tc)
		}
		suite.Time = fmt.Sprintf("%.3f", suiteTime)
		report.Suites = append(report.Suites, suite)
	}

	xmlBytes, err := xml.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal JUnit XML: %v", err)
	}

	fullXML := []byte(xml.Header + string(xmlBytes) + "\n")
	if err := os.WriteFile("reports/junit.xml", fullXML, 0644); err != nil {
		log.Fatalf("failed to write reports/junit.xml: %v", err)
	}

	totalTests, totalFailures := 0, 0
	for _, s := range report.Suites {
		totalTests += s.Tests
		totalFailures += s.Failures
	}
	fmt.Printf("reports/junit.xml written: %d tests, %d failures across %d packages\n", totalTests, totalFailures, len(report.Suites))

	if totalFailures > 0 {
		os.Exit(1)
	}
}

// executeRegressionReport writes reports/regression.xml - a small
// JUnit-format summary of the 4-step regression pipeline itself
// (Build/Tests/Coverage/KB & RAG), distinct from junit.xml (every
// individual Go test). Takes each step's already-captured real
// status as input rather than re-deriving it, so this logic lives
// in exactly one place: regression.sh and regression.bat both call
// this with their real, honestly-tracked status variables, instead
// of each separately reimplementing XML-writing in bash and batch -
// the same duplication this whole task runner exists to eliminate.
func executeRegressionReport(buildStatus, testStatus, coverageStatus, kbragStatus, elapsedSeconds string) {
	fmt.Println("Generating reports/regression.xml...")

	if err := os.MkdirAll("reports", 0755); err != nil {
		log.Fatalf("failed to create reports/ directory: %v", err)
	}

	steps := []struct {
		name   string
		status string
	}{
		{"Build", buildStatus},
		{"Tests", testStatus},
		{"Coverage", coverageStatus},
		{"KB & RAG", kbragStatus},
	}

	suite := junitTestSuite{
		Name: "AstraMind Regression Pipeline",
		Time: elapsedSeconds,
	}

	failures := 0
	for _, step := range steps {
		tc := junitTestCase{
			ClassName: "regression",
			Name:      step.name,
		}
		switch step.status {
		case "FAIL":
			failures++
			tc.Failure = &junitFailure{Message: "Step failed", Content: step.name + " reported a non-zero exit code."}
		case "SKIPPED":
			tc.Skipped = &junitSkipped{Message: "Not run - an earlier step failed"}
		}
		suite.Tests++
		suite.TestCases = append(suite.TestCases, tc)
	}
	suite.Failures = failures

	report := junitTestSuites{Suites: []junitTestSuite{suite}}

	xmlBytes, err := xml.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal regression.xml: %v", err)
	}

	fullXML := []byte(xml.Header + string(xmlBytes) + "\n")
	if err := os.WriteFile("reports/regression.xml", fullXML, 0644); err != nil {
		log.Fatalf("failed to write reports/regression.xml: %v", err)
	}

	fmt.Printf("reports/regression.xml written: %d steps, %d failures\n", suite.Tests, failures)
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
