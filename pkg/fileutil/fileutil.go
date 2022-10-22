package fileutil

import "os"

func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FileSize(path string) (size int, err error) {
	fss, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	for _, fs := range fss {
		if fs.IsDir() {
			continue
		}
		size += 1
	}

	return size, nil
}
