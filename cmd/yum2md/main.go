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
	"strings"

	"github.com/atc0005/yum2md/internal/checkupdate"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var incompleteFields []string
	var rows []checkupdate.Row

LoopLabel:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		fields, action := checkupdate.ProcessLine(line, &incompleteFields)
		switch {
		case action.SkipLine:
			continue

		case action.StopProcessing:
			break LoopLabel
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
