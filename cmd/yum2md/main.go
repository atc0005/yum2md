// Copyright 2026 Adam Chalkley
//
// https://github.com/atc0005/yum2md
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/atc0005/yum2md/internal/checkupdate"
)

func processLine(line string, incompleteFields *[]string) ([]string, bool) {
	if line == "" {
		return nil, true
	}

	// Skip over header for obsolete packages section and all lines
	// following it.
	if strings.TrimSpace(line) == checkupdate.ObsoletePackagesHeader {
		return nil, true
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
		return nil, true

	case len(fields) < 3 && len(*incompleteFields) == 0:
		// We have less than three fields after splitting the input and
		// don't currently have a saved copy of the available fields from
		// a previous line which didn't match the required count.
		//
		// Save current fields and attempt to reconstruct on the next pass
		// with the following line. We rely on later validation to reject
		// any invalid input lines that happens to also contain 3 fields.
		*incompleteFields = slices.Clone(fields)

		return nil, true

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

		return fields, false

	case len(fields) < 3 && len(*incompleteFields) != 0:
		// fmt.Println("DEBUG: Failed to match previous incomplete line and this one")
		// fmt.Printf("DEBUG: Previous incomplete line: %q, this one: %q\n", incompleteFields, fields)
		// fmt.Printf("DEBUG: len(fields)==%d, len(incompleteFields)==%d\n", len(fields), len(incompleteFields))

		// We attempted to capture the previous incomplete line and match
		// that up with the current line and that didn't work out. Let's
		// clear the previous saved line and save the current one instead.
		*incompleteFields = slices.Clone(fields)

		// Let's try the next line.
		return nil, true

	}

	return fields, false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var incompleteFields []string
	var rows []checkupdate.Row

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		fields, skipLine := processLine(line, &incompleteFields)
		if skipLine {
			continue
		}

		// Skip over any lines which contain capital letters in any of the fields.
		//
		// FIXME: This won't work since there are valid repo names (usually
		// custom) and valid package names (e.g., perl-XML-Parser) which are
		// caught by this check).
		//
		// if checkupdate.CollectionContainsCapital(fields) {
		// 	fmt.Println("Skipping line due to collectionContainsCapital check")
		// 	fmt.Println(fields)
		// 	continue
		// }

		// Skip over any lines which do not contain a separator character in expected fields.
		if !checkupdate.CollectionContainsPackageNameSeparator(fields) {
			// fmt.Println("Skipping line due to collectionContainsPackageNameSeparator check")
			// fmt.Printf("SKIPPED:%q\n", fields)
			continue
		}

		// Earlier sanity check rejects lines without the required 3 fields.
		pkg := fields[0]
		release := fields[1]
		// repo := fields[len(fields)-1]
		// release := strings.Join(fields[1:len(fields)-1], " ") // TODO: Is this really needed?
		repo := fields[2]

		rows = append(rows, checkupdate.Row{
			Package: pkg,
			Release: release,
			Repo:    repo,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}

	checkupdate.PrintMarkdownTable(rows)
}
