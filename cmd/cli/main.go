package main

import (
	"bazaar/pkg/httpclient"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var path string
	flag.StringVar(&path, "sample", "testdata/samples", "sample dir")
	flag.Parse()

	uri := "http://127.0.0.1:18528/bazaar/v1/upload"

	filepath.Walk(path, func(path string, info fs.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		res, err := httpclient.PostFormFile(
			uri,
			nil,
			httpclient.FileValues{
				"file": {
					Name:   filepath.Base(path),
					Reader: f,
				},
			},
			httpclient.WithTTL(time.Second*30),
		)
		if err != nil {
			fmt.Println("upload file error", path, err)
			return nil
		}

		_, err = Parse(res)
		if err != nil {
			fmt.Println("upload file error", path, err)
			return nil
		}

		fmt.Println("upload file success", path)

		return nil
	})
}

// http reply
const StatusOk = iota

type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func Parse(body []byte) (*Reply, error) {
	var res Reply
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Code != StatusOk {
		return &res, fmt.Errorf("code: %d, msg: %s", res.Code, res.Msg)
	}

	return &res, nil
}
