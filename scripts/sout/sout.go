// Package sout provides structured output for scripts.
// This package intentionally does NOT use structured logging (pkg/logger)
// so scripts can output consistently without being flagged by lint-quality.
// Only use for CLI output; never use in production application code.
package sout

import (
	"fmt"
	"os"
)

// Pass prints a success message and exits with code 0.
func Pass(msg string) {
	fmt.Println("✓ " + msg)
	os.Exit(0)
}

// Fail prints an error message and exits with code 1.
func Fail(msg string) {
	fmt.Println("✗ " + msg)
	os.Exit(1)
}

// Failf prints a formatted error message and exits with code 1.
func Failf(format string, args ...interface{}) {
	fmt.Printf("✗ "+format+"\n", args...)
	os.Exit(1)
}

// Info prints an info message.
func Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
