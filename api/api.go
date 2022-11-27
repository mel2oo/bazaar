package api

import (
	"bazaar/internal/router"
	"bazaar/pkg/httpclient"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	_MalwareUpload = "/bazaar/v1/upload"
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
		httpclient.WithTTL(time.Second*5),
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
