package api

import (
	"bazaar/internal/router"
	"bazaar/pkg/httpclient"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	_MalwareUpload   = "/bazaar/v1/upload"
	_MalwareQuery    = "/bazaar/v1/query"
	_MalwareDownload = "/bazaar/v1/download"
)

func MalwareUpload(host, file, ext string, tags []string) error {
	uri := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   _MalwareUpload,
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	res, err := httpclient.PostFormFile(
		uri.String(),
		url.Values{
			"type": []string{ext},
			"tags": tags,
		},
		httpclient.FileValues{
			"file": {
				Name:   filepath.Base(file),
				Reader: f,
			},
		},
		httpclient.WithTTL(time.Second*20),
	)
	if err != nil {
		return err
	}

	_, err = router.Parse(res)
	if err != nil {
		return err
	}

	return nil
}

type QueryResult struct {
	Date   string   `json:"date,omitempty"`
	Name   string   `json:"name,omitempty"`
	MD5    string   `json:"md5,omitempty"`
	SHA256 string   `json:"sha256,omitempty"`
	Type   string   `json:"type,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

func MalwareQuery(host, hash, ext, tag, size string) ([]QueryResult, error) {
	uri := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   _MalwareQuery,
	}

	res, err := httpclient.Get(
		uri.String(),
		url.Values{
			"hash": []string{hash},
			"type": []string{ext},
			"tags": []string{tag},
			"size": []string{size},
		},
		httpclient.WithTTL(time.Second*20),
	)
	if err != nil {
		return nil, err
	}

	rep, err := router.Parse(res)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(rep.Data)
	if err != nil {
		return nil, err
	}

	list := make([]QueryResult, 0)
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func MalwareDownload(host, md5 string) ([]byte, error) {
	uri := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   _MalwareDownload,
	}

	return httpclient.Get(
		uri.String(),
		url.Values{
			"md5": []string{md5},
		},
		httpclient.WithTTL(time.Second*20),
	)
}
