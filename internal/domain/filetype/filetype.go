package filetype

import (
	"bazaar/pkg/exiftool"
	"bazaar/pkg/trid"
	"os"
	"strconv"
	"time"
)

const _TypeUnkonw = "unknow"

func Scan(file string) (ext string, err error) {
	ext, err = exiftool.ScanExt(file)
	if err != nil && ExtClass(ext) != TYPE_UNDEFINE {
		return ext, nil
	}

	ext, err = trid.ScanExt(file)
	if err != nil || len(ext) == 0 {
		return _TypeUnkonw, nil
	}

	return ext, nil
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
