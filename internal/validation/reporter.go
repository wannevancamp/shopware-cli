package validation

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

func DetectDefaultReporter() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return "github"
	}

	return "summary"
}

func DoCheckReport(result Check, reportingFormat string) error {
	switch reportingFormat {
	case "summary":
		return doSummaryReport(result)
	case "json":
		return doJSONReport(result)
	case "github":
		return doGitHubReport(result)
	case "markdown":
		return doMarkdownReport(result)
	case "junit":
		return doJUnitReport(result)
	}

	return nil
}

func doSummaryReport(result Check) error {
	// Group results by file
	fileGroups := make(map[string][]CheckResult)
	for _, r := range result.GetResults() {
		fileGroups[r.Path] = append(fileGroups[r.Path], r)
	}

	// Print results grouped by file
	totalProblems := 0
	for path, results := range fileGroups {
		totalProblems += len(results)
		if len(results) == 0 {
			continue
		}

		fmt.Printf("❌ %s (%d problems)\n", path, len(results))
		for _, r := range results {
			severity := "⚠️"
			if r.Severity == SeverityError {
				severity = "❌"
			}

			location := ""
			if r.Line > 0 {
				location = fmt.Sprintf(":%d", r.Line)
			}

			fmt.Printf("  %s %s%s: %s (%s)\n", severity, r.Path, location, r.Message, r.Identifier)
		}
		fmt.Println()
	}

	if totalProblems > 0 {
		fmt.Printf("Found %d problems\n", totalProblems)
	} else {
		fmt.Println("✅ No problems found")
	}

	if result.HasErrors() {
		return fmt.Errorf("found errors")
	}

	return nil
}

func doJSONReport(result Check) error {
	data := map[string]interface{}{
		"results": result.GetResults(),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func doGitHubReport(result Check) error {
	for _, r := range result.GetResults() {
		file := r.Path
		if file == "" {
			file = "."
		}

		level := SeverityWarning
		if r.Severity == SeverityError {
			level = SeverityError
		}

		line := ""
		if r.Line > 0 {
			line = fmt.Sprintf(",line=%d", r.Line)
		}

		message := strings.ReplaceAll(r.Message, "\n", "%0A")
		message = strings.ReplaceAll(message, "\r", "%0D")

		fmt.Printf("::%s file=%s%s,title=%s::%s\n", level, file, line, r.Identifier, message)
	}

	return nil
}

func doMarkdownReport(result Check) error {
	// Group results by file
	fileGroups := make(map[string][]CheckResult)
	for _, r := range result.GetResults() {
		if r.Path == "" {
			r.Path = "general"
		}

		fileGroups[r.Path] = append(fileGroups[r.Path], r)
	}

	fmt.Println("# Validation Report")
	fmt.Println()

	totalProblems := 0
	for path, results := range fileGroups {
		totalProblems += len(results)
		if len(results) == 0 {
			continue
		}

		fmt.Printf("## %s (%d problems)\n\n", path, len(results))
		for _, r := range results {
			severity := "⚠️ Warning"
			if r.Severity == "error" {
				severity = "❌ Error"
			}

			location := ""
			if r.Line > 0 {
				location = fmt.Sprintf(":%d", r.Line)
			}

			fmt.Printf("- **%s** %s%s: %s (`%s`)\n", severity, r.Path, location, r.Message, r.Identifier)
		}
		fmt.Println()
	}

	if totalProblems == 0 {
		fmt.Println("✅ No problems found")
	}

	return nil
}

type JUnitTestSuite struct {
	XMLName  xml.Name        `xml:"testsuite"`
	Name     string          `xml:"name,attr"`
	Tests    int             `xml:"tests,attr"`
	Failures int             `xml:"failures,attr"`
	Errors   int             `xml:"errors,attr"`
	TestCase []JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	Name      string            `xml:"name,attr"`
	ClassName string            `xml:"classname,attr"`
	Failure   *JUnitTestFailure `xml:"failure,omitempty"`
	Error     *JUnitTestError   `xml:"error,omitempty"`
}

type JUnitTestFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

type JUnitTestError struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

func doJUnitReport(result Check) error {
	var testCases []JUnitTestCase
	errors := 0
	failures := 0

	for _, r := range result.GetResults() {
		testCase := JUnitTestCase{
			Name:      r.Identifier,
			ClassName: r.Path,
		}

		if r.Severity == SeverityError {
			testCase.Error = &JUnitTestError{
				Message: r.Message,
				Type:    r.Identifier,
				Content: r.Message,
			}
			errors++
		} else {
			testCase.Failure = &JUnitTestFailure{
				Message: r.Message,
				Type:    r.Identifier,
				Content: r.Message,
			}
			failures++
		}

		testCases = append(testCases, testCase)
	}

	suite := JUnitTestSuite{
		Name:     "shopware-cli-validation",
		Tests:    len(testCases),
		Failures: failures,
		Errors:   errors,
		TestCase: testCases,
	}

	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("", "  ")
	return encoder.Encode(suite)
}
