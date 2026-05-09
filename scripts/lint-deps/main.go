// scripts/lint-deps/main.go
//
// Validates that package dependencies follow the layer hierarchy.
// Each layer can only import from lower layers.
//
// Layer Map:
//   Layer 0 (lowest): internal/types, internal/model
//   Layer 1:          internal/repository, internal/cache
//   Layer 2:          internal/service
//   Layer 3:          internal/handler, internal/middleware, internal/request, internal/response
//   Layer 4:          cmd/server, cmd/migrate, cmd/seed
//   Layer -1 (infrastructure, highest): pkg/logger, pkg/config, pkg/database
//
// Usage:
//   go run ./scripts/lint-deps              # human-readable output
//   go run ./scripts/lint-deps --json       # machine-parseable JSON output
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/example/dormitory-management/scripts/linttypes"
	"github.com/example/dormitory-management/scripts/sout"
)

const modulePath = "github.com/example/dormitory-management"

// Layer hierarchy: lower index = lower layer
// Layer -1 is special: infrastructure packages that any layer may import.
var layers = [][]string{
	// Layer 0: Core types — no internal dependencies allowed
	{"internal/types", "internal/model", "internal/request"},
	// Layer 1: Data access layer — depends on Layer 0
	{"internal/repository", "internal/cache"},
	// Layer 2: Business logic layer — depends on Layer 0-1
	{"internal/service"},
	// Layer 3: HTTP layer — depends on Layer 0-2
	{"internal/handler", "internal/middleware", "internal/response"},
	// Layer 4: Entry points — depends on all layers
	{"cmd/server", "cmd/migrate", "cmd/seed"},
}

// layer -1 is handled separately via pkgAllowedFromLayer map below
var pkgAllowedFromLayer = map[string]int{
	"github.com/example/dormitory-management/pkg/logger":   -1,
	"github.com/example/dormitory-management/pkg/config":   -1,
	"github.com/example/dormitory-management/pkg/database": -1,
}

var (
	jsonFlag = flag.Bool("json", false, "Output results as JSON")
)

func main() {
	flag.Parse()
	violations := checkDependencies()

	if *jsonFlag {
		outputJSON(violations)
		return
	}

	if len(violations) == 0 {
		sout.Pass("All package dependencies follow the layer hierarchy")
		os.Exit(0)
	}

	fmt.Printf("✗ Found %d dependency violations:\n", len(violations))
	for _, v := range violations {
		fmt.Printf("  %s:\n    Package: %s (Layer %d) → Imports: %s (Layer %d)\n    Fix: Move logic to a lower layer, or pass as a parameter.\n\n",
			v.File, v.Package, v.FromLayer, v.Imports, v.ToLayer)
	}
	sout.Fail("Dependency violations found")
	os.Exit(1)
}

// JSON output types
type JSONReport struct {
	Timestamp  string                  `json:"timestamp"`
	Passed    bool                    `json:"passed"`
	ViolationCount int                `json:"violation_count"`
	Violations []linttypes.Violation  `json:"violations,omitempty"`
}

func outputJSON(violations []linttypes.Violation) {
	report := JSONReport{
		Timestamp:       time.Now().UTC().Format(time.RFC3339),
		Passed:         len(violations) == 0,
		ViolationCount: len(violations),
		Violations:     violations,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "JSON encode error: %v\n", err)
		os.Exit(1)
	}
	if len(violations) > 0 {
		os.Exit(1)
	}
}

func checkDependencies() []linttypes.Violation {
	var violations []linttypes.Violation
	layerMap := buildLayerMap()

	skipDirs := map[string]bool{
		".git":        true,
		"vendor":      true,
		"dist":        true,
		"node_modules": true,
		"scripts":     true,
	}

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		violations = append(violations, checkFile(path, layerMap)...)
		return nil
	})

	return violations
}

func buildLayerMap() map[string]int {
	m := make(map[string]int)
	for idx, pkgs := range layers {
		for _, pkg := range pkgs {
			m[pkg] = idx
		}
	}
	return m
}

func checkFile(path string, layerMap map[string]int) []linttypes.Violation {
	var violations []linttypes.Violation

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return violations
	}

	dir := filepath.ToSlash(filepath.Dir(path))
	dir = strings.TrimPrefix(dir, "./")
	pkgLayer := findLayer(dir, layerMap)
	if pkgLayer < -1 {
		return violations
	}

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		if !strings.HasPrefix(importPath, modulePath) {
			continue
		}

		relImport := strings.TrimPrefix(importPath, modulePath+"/")
		importLayer := findLayer(relImport, layerMap)
		if importLayer < -1 {
			continue
		}

		// Layer -1 packages (pkg/*) can be imported by any layer
		if _, allowed := pkgAllowedFromLayer[importPath]; allowed {
			continue
		}

		if pkgLayer != -1 && importLayer > pkgLayer && !isSameBasePackage(dir, relImport) {
			violations = append(violations, linttypes.Violation{
				File:     path,
				Package:  dir,
				Imports:  relImport,
				FromLayer: pkgLayer,
				ToLayer:   importLayer,
				Message:   fmt.Sprintf("Layer %d cannot import Layer %d", pkgLayer, importLayer),
			})
		}
	}

	return violations
}

func findLayer(pkg string, layerMap map[string]int) int {
	if layer, ok := layerMap[pkg]; ok {
		return layer
	}
	for key, layer := range layerMap {
		if strings.HasPrefix(pkg, key+"/") {
			return layer
		}
	}
	return -2 // Not found
}

// isSameBasePackage returns true if a and b are the same package or
// both belong to the same parent package (e.g., internal/service and internal/service/sub).
func isSameBasePackage(a, b string) bool {
	if a == b {
		return true
	}
	// a is a sub-package of b
	if strings.HasPrefix(a, b+"/") {
		return true
	}
	// b is a sub-package of a
	if strings.HasPrefix(b, a+"/") {
		return true
	}
	return false
}
