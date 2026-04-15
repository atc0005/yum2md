// Copyright 2026 Adam Chalkley
//
// https://github.com/atc0005/check-memory
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package checkupdate

import (
	"fmt"
	"strings"
	// "unicode"
)

// Row represents a row of pending updates data to be written to a
// Markdown-formatted table.
type Row struct {
	Package string
	Release string
	Repo    string
}

// ObsoletePackagesHeader is a header found in `check-updates` output for a
// list of packages that will be removed when updating packages. This header
// and following content is intended to be ignored.
const ObsoletePackagesHeader string = "Obsoleting Packages"

// CollectionContainsPackageNameSeparator checks if all string in the
// collection contains at least one of the known valid package name separators.
func CollectionContainsPackageNameSeparator(ss []string) bool {
	for _, s := range ss {
		if !ContainsPackageNameSeparator(s) {
			return false
		}
	}

	return true
}

// ContainsPackageNameSeparator checks if the string s contains at least one
// of the known package name separator characters.
func ContainsPackageNameSeparator(s string) bool {
	return strings.ContainsAny(s, ".-_")
}

// CollectionContainsCapital checks if any string in the collection contains
// at least one capital letter.
// func CollectionContainsCapital(ss []string) bool {
// 	for _, s := range ss {
// 		if ContainsCapital(s) {
// 			return true
// 		}
// 	}

// 	return false
// }

// ContainsCapital checks if the string s contains at least one capital letter.
// func ContainsCapital(s string) bool {
// 	for _, r := range s {
// 		if unicode.IsUpper(r) {
// 			return true
// 		}
// 	}
// 	return false
// }

// PrintMarkdownTable generates and emits a Markdown-formatted table from
// parsed check-output rows containing relevant package/release/repo details.
func PrintMarkdownTable(rows []Row) {
	if len(rows) == 0 {
		fmt.Println("No updates to apply.")

		return
	}

	headers := []string{"Package", "Release", "Repo"}

	// Compute column widths
	wPkg := len(headers[0])
	wRel := len(headers[1])
	wRepo := len(headers[2])

	for _, r := range rows {
		if len(r.Package) > wPkg {
			wPkg = len(r.Package)
		}
		if len(r.Release) > wRel {
			wRel = len(r.Release)
		}
		if len(r.Repo) > wRepo {
			wRepo = len(r.Repo)
		}
	}

	// Header row
	fmt.Printf("| %-*s | %-*s | %-*s |\n",
		wPkg, headers[0],
		wRel, headers[1],
		wRepo, headers[2],
	)

	// Separator row
	fmt.Printf("|-%s-|-%s-|-%s-|\n",
		strings.Repeat("-", wPkg),
		strings.Repeat("-", wRel),
		strings.Repeat("-", wRepo),
	)

	// Data rows
	for _, r := range rows {
		fmt.Printf("| %-*s | %-*s | %-*s |\n",
			wPkg, r.Package,
			wRel, r.Release,
			wRepo, r.Repo,
		)
	}
}
