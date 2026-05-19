// Copyright 2026 Adam Chalkley
//
// https://github.com/atc0005/check-memory
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package checkupdate

import (
	"fmt"
	"slices"
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

// ProcessAction indicates what action should be taken by the caller for the
// given line of output from `yum check-update`.
type ProcessAction struct {
	SkipLine       bool
	StopProcessing bool
}

// ProcessLine evaluates a given line of output from `yum check-update` and
// indicates to the caller what action should be taken for that line.
//
// The caller is responsible for providing a *[]string to permit saving
// incomplete fields retrieved from a line of `yum check-update` output for a
// later pass. RHEL 7 in particular is known to break long lines of `yum
// check-update` when it (incorrectly) detects that the output sink is too
// narrow to properly hold the entire line without breaking it across multiple
// lines.
func ProcessLine(line string, incompleteFields *[]string) ([]string, ProcessAction) {
	if line == "" {
		return nil, ProcessAction{SkipLine: true}
	}

	// Skip over header for obsolete packages section and all lines
	// following it.
	if strings.TrimSpace(line) == ObsoletePackagesHeader {
		return nil, ProcessAction{StopProcessing: true}
	}

	fields := strings.Fields(line)

	// RHEL 7+ yum/dnf safety rule:
	// Only accept lines that can form exactly:
	//   package | release | repo
	//
	// We attempt to exclude known problematic patterns and hotfix any
	// check-update output that may be incorrectly broken over multiple
	// lines. This may occur from `yum` incorrectly detecting that it is
	// writing output to a reduced width terminal.
	switch {
	case len(fields) == 0:
		// Skip blank lines.
		return nil, ProcessAction{SkipLine: true}

	case len(fields) < 3 && len(*incompleteFields) == 0:
		// We have less than three fields after splitting the input and
		// don't currently have a saved copy of the available fields from
		// a previous line which didn't match the required count.
		//
		// Save current fields and attempt to reconstruct on the next pass
		// with the following line. We rely on later validation to reject
		// any invalid input lines that happens to also contain 3 fields.
		*incompleteFields = slices.Clone(fields)

		return nil, ProcessAction{SkipLine: true}

	case len(fields)+len(*incompleteFields) == 3:
		// Attempt to detect and hotfix unexpected line wrapping.
		//
		// Example 1:
		//
		//    ca-certificates.noarch           2025.2.80_v9.0.304-71.el7_9
		//                                                              rhel-7-server-els-rpms
		//
		// Example 2:
		//
		//    device-mapper-persistent-data.x86_64
		//                                     0.8.5-3.el7_9.2          rhel-7-server-els-rpms
		fields = append(*incompleteFields, fields...)

		// CRITICAL step; we need to wipe the buffer after retrieving its
		// contents to prevent unexpected reuse.
		*incompleteFields = nil

		return fields, ProcessAction{SkipLine: false}

	case len(fields) < 3 && len(*incompleteFields) != 0:
		// fmt.Println("DEBUG: Failed to match previous incomplete line and this one")
		// fmt.Printf("DEBUG: Previous incomplete line: %q, this one: %q\n", incompleteFields, fields)
		// fmt.Printf("DEBUG: len(fields)==%d, len(incompleteFields)==%d\n", len(fields), len(incompleteFields))

		// We attempted to capture the previous incomplete line and match
		// that up with the current line and that didn't work out. Let's
		// clear the previous saved line and save the current one instead.
		*incompleteFields = slices.Clone(fields)

		// Let's try the next line.
		return nil, ProcessAction{SkipLine: true}

	}

	return fields, ProcessAction{SkipLine: false}
}

// PrintMarkdownTable generates and emits a Markdown-formatted table from
// parsed check-output rows containing relevant package/release/repo details.
func PrintMarkdownTable(rows []Row) {
	switch {
	case len(rows) == 0:
		fmt.Println("No updates to apply.")

		return

	case len(rows) == 1:
		fmt.Printf("%d pending update to apply.\n", len(rows))

	default:
		fmt.Printf("%d pending updates to apply.\n", len(rows))
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
