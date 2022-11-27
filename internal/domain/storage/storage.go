package storage

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrPathInvalid = errors.New("root path invalid")
)

// Config represents the storage config.
type Config struct {
	Root string `mapstructure:"root-path"`
	Seek string `mapstructure:"seek-file"`
}

type Client struct {
	root  string
	seeks *Seeks
	mutex *sync.Mutex
}

func New(c Config) (*Client, error) {
	stat, err := os.Stat(c.Root)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(c.Root, os.ModePerm); err != nil {
				return nil, err
			}
		}
	} else {
		if !stat.IsDir() {
			return nil, ErrPathInvalid
		}
	}

	seeks, err := InitSeek(c.Root, c.Seek)
	if err != nil {
		return nil, err
	}

	return &Client{
		root:  c.Root,
		seeks: seeks,
		mutex: new(sync.Mutex),
	}, nil
}

func (c *Client) Create(data []byte, md5, ext string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	path, err := c.seeks.Check(c.root, ext)
	if err != nil {
		return path, err
	}

	full := filepath.Join(path, md5)
	file, err := os.Create(full)
	if err != nil {
		return full, err
	}
	defer file.Close()

	_, err = file.Write(data)
	return full, err
}
