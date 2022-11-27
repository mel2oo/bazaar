// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package exiftool

import (
	"bazaar/pkg/cmd"
	"bazaar/pkg/utils"
	"strings"
	"unicode"
)

const (
	// Command to invoke exiftool scanner
	bin = "exiftool"
)

// Scan a file using exiftool
// This will execute exigtool command line tool and read the stdout
func Scan(file string) (map[string]string, error) {
	output, err := cmd.RunOutputTimeout(bin, file)
	if err != nil {
		return nil, err
	}

	return ParseOutput(output), nil
}

// ParseOutput convert exiftool output into map of string|string.
func ParseOutput(exifout string) map[string]string {
	var ignoreTags = []string{
		"Directory",
		"File Name",
		"File Permissions",
	}

	lines := strings.Split(exifout, "\n")
	if utils.StringInSlice("File not found", lines) {
		return nil
	}

	datas := make(map[string]string, len(lines))
	for _, line := range lines {
		keyvalue := strings.Split(line, ":")
		if len(keyvalue) != 2 {
			continue
		}
		if !utils.StringInSlice(strings.TrimSpace(keyvalue[0]), ignoreTags) {
			datas[strings.TrimSpace(camelCase(keyvalue[0]))] =
				strings.TrimSpace(keyvalue[1])
		}
	}

	return datas
}

// camelCase convert a string to camelcase
func camelCase(s string) string {
	s = strings.TrimSpace(s)
	buffer := make([]rune, 0, len(s))
	stringIter(s, func(prev, curr, next rune) {
		if !isDelimiter(curr) {
			if isDelimiter(prev) || (prev == 0) {
				buffer = append(buffer, unicode.ToUpper(curr))
			} else if unicode.IsLower(prev) {
				buffer = append(buffer, curr)
			} else {
				buffer = append(buffer, unicode.ToLower(curr))
			}
		}
	})

	return string(buffer)
}

// isDelimiter checks if a character is some kind of whitespace or '_' or '-'.
func isDelimiter(ch rune) bool {
	return ch == '-' || ch == '_' || unicode.IsSpace(ch)
}

// stringIter iterates over a string, invoking the callback for every single rune in the string.
func stringIter(s string, callback func(prev, curr, next rune)) {
	var prev rune
	var curr rune
	for _, next := range s {
		if curr == 0 {
			prev = curr
			curr = next
			continue
		}

		callback(prev, curr, next)

		prev = curr
		curr = next
	}

	if len(s) > 0 {
		callback(prev, curr, 0)
	}
}

// exif extension
const (
	_SysFlag = "native"
	_ScrFlag = "screen saver"
	_ComFlag = "com"
	_OcxFlag = ".ocx"
	_ElfFlag = "elf"
)

type ExifInfo struct {
	FileType         string
	Extension        string
	FileDescription  string
	OriginalFileName string
	Subsystem        string
}

func ScanExt(file string) (string, error) {
	m, err := Scan(file)
	if err != nil {
		return "", err
	}

	e := ExifInfo{
		FileType:         strings.ToLower(m["FileType"]),
		Extension:        strings.ToLower(m["FileTypeExtension"]),
		FileDescription:  strings.ToLower(m["FileDescription"]),
		OriginalFileName: strings.ToLower(m["OriginalFileName"]),
		Subsystem:        strings.ToLower(m["Subsystem"]),
	}

	if len(e.Extension) > 0 {
		switch e.Extension {
		case "exe":
			if e.Subsystem == _SysFlag {
				return "sys", nil
			} else if strings.HasSuffix(e.OriginalFileName, _ComFlag) {
				return "com", nil
			} else if strings.Contains(e.FileDescription, _ScrFlag) {
				return "scr", nil
			}
			return e.Extension, nil

		case "dll":
			if strings.HasSuffix(e.OriginalFileName, _OcxFlag) {
				return "ocx", nil
			}
			return e.Extension, nil
		}
	}

	if len(e.FileType) > 0 {
		if strings.Contains(e.FileType, _ElfFlag) {
			return "elf", nil
		}
	}

	return e.Extension, nil
}
