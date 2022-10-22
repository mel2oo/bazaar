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
	if utils.StringInSlice("Error: found no file(s) to analyze!", lines) {
		return nil
	}
	lines = lines[6:]

	for _, line := range lines {
		if len(strings.TrimSpace(line)) != 0 {
			keeps = append(keeps, strings.TrimSpace(line))
		}
	}

	return keeps
}
