// scripts/lint-quality/main.go
//
// Validates golden principles:
// - No raw log.Printf (use structured logging from pkg/logger)
// - File size limits (max 500 lines)
// - No fmt.Printf in production code
//
// Usage: go run ./scripts/lint-quality
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/example/dormitory-management/scripts/sout"
)

const maxFileLines = 500

// Patterns to flag as violations
var rawLogPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\blog\.Printf\b`),
	regexp.MustCompile(`\blog\.Println\b`),
	regexp.MustCompile(`\blog\.Fatalf\b`),
	regexp.MustCompile(`\blog\.Fatal\b`),
	regexp.MustCompile(`\bfmt\.Print\b`),
	regexp.MustCompile(`\bfmt\.Printf\b`),
	regexp.MustCompile(`\bfmt\.Println\b`),
}

// Directories to skip
var skipDirs = map[string]bool{
	".git":        true,
	".svn":        true,
	"dist":        true,
	"vendor":      true,
	"node_modules": true,
	"scripts":     true, // scripts are developer tools, not production code
}

type QualityViolation struct {
	File    string
	Line    int
	Rule    string
	Message string
}

func main() {
	var violations []QualityViolation

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			if info != nil && info.IsDir() && skipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		isTest := strings.HasSuffix(path, "_test.go")

		// Check for raw logging (skip test files)
		if !isTest {
			for lineNum, line := range lines {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "//") {
					continue
				}
				for _, pattern := range rawLogPatterns {
					if pattern.MatchString(line) {
						violations = append(violations, QualityViolation{
							File:    path,
							Line:    lineNum + 1,
							Rule:    "structured-logging",
							Message: fmt.Sprintf("Use structured logging (pkg/logger) instead of raw log calls"),
						})
					}
				}
			}
		}

		// Check file size
		lineCount := len(lines)
		if lineCount > maxFileLines {
			violations = append(violations, QualityViolation{
				File:    path,
				Line:    0,
				Rule:    "file-size",
				Message: fmt.Sprintf("File has %d lines (max %d). Split into smaller, focused modules.", lineCount, maxFileLines),
			})
		}

		return nil
	})

	if len(violations) == 0 {
		sout.Pass("All quality checks passed")
	}

	fmt.Printf("✗ Found %d quality violations:\n\n", len(violations))
	for _, v := range violations {
		if v.Line > 0 {
			fmt.Printf("%s:%d [%s]: %s\n", v.File, v.Line, v.Rule, v.Message)
		} else {
			fmt.Printf("%s [%s]: %s\n", v.File, v.Rule, v.Message)
		}
	}
	os.Exit(1)
}