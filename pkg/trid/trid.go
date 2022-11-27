// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trid

import (
	"bazaar/pkg/cmd"
	"bazaar/pkg/utils"
	"strings"
)

const (
	// Command to invoke TriD scanner.
	bin = "trid"
)

// Scan a file using TRiD Scanner
func Scan(file string) ([]string, error) {
	output, err := cmd.RunOutputTimeout(bin, file)
	if err != nil {
		return nil, err
	}
	return parseOutput(output), nil
}

// parseOutput parse TriD stdout, returns an array of strings
func parseOutput(tridout string) []string {
	keeps := make([]string, 0)

	lines := strings.Split(tridout, "\n")
	if utils.SliceContainsString("Error: found no file(s) to analyze!", lines) {
		return keeps
	}

	if len(lines) > 7 {
		lines = lines[6:]

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			if strings.HasPrefix(line, "Warning:") ||
				strings.Contains(line, "TrID is best suited to analyze binary files!") {
				continue
			}
			keeps = append(keeps, line)
		}
	}

	return keeps
}

func ScanExt(file string) (ext string, err error) {
	res, err := Scan(file)
	if err != nil {
		return ext, err
	}

	for _, s := range res {
		left := strings.Index(s, "(")
		right := strings.Index(s, ")")
		if left+2 < right {
			ext = strings.ToLower(s[left+2 : right])
			break
		}
	}
	return
}
