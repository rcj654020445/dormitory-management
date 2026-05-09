// Package linttypes defines shared types for lint scripts.
package linttypes

// Violation represents a code rule violation found during linting.
type Violation struct {
	File     string
	Package  string
	Imports  string
	FromLayer int
	ToLayer   int
	Message   string
}