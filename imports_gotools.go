// This file contains some code copied from the Go x/tools project.

package gopoet

import (
	"path"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// This file deals with how packages are assumed to be named based on their
// import path.
//
// The package name used for lexically matching a qualified identifier (e.g.,
// the "foo" is the package name in "foo.bar") cannot be determined without
// knowing the package name used in the files that define an imported package.
// Most Go tools require package names follow a convention so that the import
// path can be used to infer the package name.
//
// Go Poet has used path.Base for this historically. Below is an alternative
// implementation that matches what's used by most go tools. The API should
// probably change to allow specifying this function by the user.
//
// Go spec: https://golang.org/ref/spec#Import_declarations.

// assumedPackageNameForImport uses path.Base to assume a package name
// from an import.
func assumedPackageNameForImport(importPath string) string {
	return path.Base(importPath)
}

// goToolsImportConventions uses the same rules as the Go tools for inferring
// the name of a package from an import path.
//
// AssumedPackageName returns the assumed package name of an import
// path. It does this using only string parsing of the import path.
//
// It picks the last element of the path that does not look like a major
// version, and then picks the valid identifier off the start of that element.
// It is used to determine if a local rename should be added to an import for
// clarity.
//
// This package is copied from
// https://pkg.go.dev/golang.org/x/tools/internal/imports#ImportPathToAssumedName.
func goToolsAssumedPackageNameForImport(importPath string) string {
	base := path.Base(importPath)
	if strings.HasPrefix(base, "v") {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			dir := path.Dir(importPath)
			if dir != "." {
				base = path.Base(dir)
			}
		}
	}
	base = strings.TrimPrefix(base, "go-")
	if i := strings.IndexFunc(base, notIdentifier); i >= 0 {
		base = base[:i]
	}
	return base
}

// notIdentifier reports whether ch is an invalid identifier character.
func notIdentifier(ch rune) bool {
	return !('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_' ||
		ch >= utf8.RuneSelf && (unicode.IsLetter(ch) || unicode.IsDigit(ch)))
}
