package storage

import (
	"bazaar/pkg/fileutil"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	_StartNo = 1
	_MaxSize = 10000
)

type Seeks struct {
	rootpath string
	seekfile string
	seekinfo map[string]Info
}

type Info struct {
	path string
	size int
}

func InitSeek(rootpath, seekfile string) (*Seeks, error) {
	skinfo := make(map[string]Info)
	skpath := filepath.Join(rootpath, seekfile)

	if fileutil.Exist(skpath) {
		data, err := os.ReadFile(skpath)
		if err != nil {
			return nil, err
		}

		current := make(map[string]string)
		if err := json.Unmarshal(data, &current); err != nil {
			return nil, err
		}

		for ext, path := range current {
			size, err := fileutil.FileSize(path)
			if err != nil {
				return nil, err
			}
			skinfo[ext] = Info{
				path: path,
				size: size,
			}
		}
	}

	return &Seeks{
		rootpath: rootpath,
		seekfile: seekfile,
		seekinfo: skinfo,
	}, nil
}

func (s *Seeks) Check(root, ext string) (string, error) {
	var newpath string

	info, ok := s.seekinfo[ext]
	if ok {
		if info.size < _MaxSize {
			s.seekinfo[ext] = Info{
				path: info.path,
				size: info.size + 1,
			}
			return info.path, nil
		}

		// 如果达到单目录存储上限, 则切换目录
		_, sno := filepath.Split(info.path)
		no, err := strconv.Atoi(sno)
		if err != nil {
			return info.path, err
		}
		newpath = filepath.Join(root, ext, fmt.Sprintf("%03d", no+1))
	}

	// 无同类型样本目录则首次创建
	if !ok {
		newpath = filepath.Join(root, ext, fmt.Sprintf("%03d", _StartNo))
	}

	if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
		return newpath, err
	}

	s.seekinfo[ext] = Info{
		path: newpath,
		size: 0,
	}
	return newpath, s.Save()
}

func (s *Seeks) Save() error {
	current := make(map[string]string)
	for ext, info := range s.seekinfo {
		current[ext] = info.path
	}

	data, err := json.Marshal(&current)
	if err != nil {
		return err
	}

	fi, err := os.Create(filepath.Join(s.rootpath, s.seekfile))
	if err != nil {
		return err
	}
	defer fi.Close()

	_, err = fi.Write(data)
	return err
}
