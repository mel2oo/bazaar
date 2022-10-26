package filetype

import (
	"bazaar/pkg/exiftool"
	"os"
	"strconv"
	"time"
)

const (
	DefaultType = "unkonw"
)

func Scan(file string) (ext string, err error) {
	maps, err := exiftool.Scan(file)
	if err != nil {
		return ext, err
	}

	ext, ok := maps["FileTypeExtension"]
	if ok {
		return ext, nil
	}

	return DefaultType, nil
}

func ScanData(data []byte) (ext string, err error) {
	tmpfile, err := os.CreateTemp("", strconv.Itoa(int(time.Now().Unix())))
	if err != nil {
		return ext, err
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(data)
	if err != nil {
		return ext, err
	}

	return Scan(tmpfile.Name())
}
