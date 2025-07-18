package validation

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportingOutputIsDeterministic(t *testing.T) {
	// Create test results in non-alphabetical order
	testResults := []CheckResult{
		{
			Path:       "z_file.go",
			Line:       5,
			Identifier: "test.rule2",
			Message:    "Second message",
			Severity:   SeverityError,
		},
		{
			Path:       "a_file.go",
			Line:       10,
			Identifier: "test.rule1",
			Message:    "First message",
			Severity:   SeverityWarning,
		},
		{
			Path:       "a_file.go",
			Line:       5,
			Identifier: "test.rule3",
			Message:    "Third message",
			Severity:   SeverityError,
		},
		{
			Path:       "b_file.go",
			Line:       1,
			Identifier: "test.rule1",
			Message:    "Fourth message",
			Severity:   SeverityWarning,
		},
	}

	check := &testCheck{Results: testResults}

	// Test summary report multiple times to ensure deterministic output
	for range 5 {
		output := captureOutput(func() {
			// Ignore the error since we're testing output format, not validation logic
			_ = doSummaryReport(check)
		})

		// Check that files are sorted alphabetically
		lines := strings.Split(output, "\n")
		var fileHeaderLines []string
		for _, line := range lines {
			// File headers are lines that contain .go but not indented (no leading spaces)
			if strings.Contains(line, ".go") && !strings.HasPrefix(line, " ") {
				fileHeaderLines = append(fileHeaderLines, strings.TrimSpace(line))
			}
		}

		assert.Len(t, fileHeaderLines, 3) // 3 files
		assert.Equal(t, "a_file.go", fileHeaderLines[0])
		assert.Equal(t, "b_file.go", fileHeaderLines[1])
		assert.Equal(t, "z_file.go", fileHeaderLines[2])
	}

	// Test GitHub report multiple times to ensure deterministic output
	for range 5 {
		output := captureOutput(func() {
			_ = doGitHubReport(check)
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, 4) // 4 results

		// Check that results are sorted by path first, then by line
		assert.Contains(t, lines[0], "a_file.go")
		assert.Contains(t, lines[0], "line=5")
		assert.Contains(t, lines[1], "a_file.go")
		assert.Contains(t, lines[1], "line=10")
		assert.Contains(t, lines[2], "b_file.go")
		assert.Contains(t, lines[2], "line=1")
		assert.Contains(t, lines[3], "z_file.go")
		assert.Contains(t, lines[3], "line=5")
	}
}

func TestMarkdownReportIsDeterministic(t *testing.T) {
	testResults := []CheckResult{
		{
			Path:       "z_file.go",
			Line:       5,
			Identifier: "test.rule2",
			Message:    "Second message",
			Severity:   SeverityError,
		},
		{
			Path:       "a_file.go",
			Line:       10,
			Identifier: "test.rule1",
			Message:    "First message",
			Severity:   SeverityWarning,
		},
	}

	check := &testCheck{Results: testResults}

	// Test markdown report multiple times to ensure deterministic output
	for range 5 {
		output := captureOutput(func() {
			_ = doMarkdownReport(check)
		})

		// Check that files are sorted alphabetically
		lines := strings.Split(output, "\n")
		var headerLines []string
		for _, line := range lines {
			if strings.HasPrefix(line, "## ") {
				headerLines = append(headerLines, line)
			}
		}

		assert.Len(t, headerLines, 2) // 2 files
		assert.Contains(t, headerLines[0], "a_file.go")
		assert.Contains(t, headerLines[1], "z_file.go")
	}
}

// testCheck is a simple implementation of Check interface for testing
type testCheck struct {
	Results []CheckResult
}

func (c *testCheck) AddResult(result CheckResult) {
	c.Results = append(c.Results, result)
}

func (c *testCheck) GetResults() []CheckResult {
	return c.Results
}

func (c *testCheck) HasErrors() bool {
	for _, r := range c.Results {
		if r.Severity == SeverityError {
			return true
		}
	}
	return false
}

func (c *testCheck) RemoveByIdentifier(ignores []ToolConfigIgnore) Check {
	// Simple implementation for testing
	return c
}

// captureOutput captures stdout during function execution
func captureOutput(fn func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		panic(err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}
	return buf.String()
}
