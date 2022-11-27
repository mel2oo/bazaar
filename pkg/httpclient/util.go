package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const (
	// _StatusReadRespErr read resp body err, should re-call doHTTP again.
	_StatusReadRespErr = -204
	// _StatusDoReqErr do req err, should re-call doHTTP again.
	_StatusDoReqErr = -500
)

var (
	defaultClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:  true,
			DisableCompression: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        100,
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 100,
		},
	}

	bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 1024))
		},
	}
)

type FileValues map[string]FileValue

type FileValue struct {
	Name   string
	Reader io.Reader
}

func doHTTP(ctx context.Context, method, url string, payload io.Reader, opt *option) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return nil, -1, fmt.Errorf("new request [%s %s] err: %s", method, url, err.Error())
	}

	for key, value := range opt.header {
		req.Header.Set(key, value[0])
	}

	req.Close = true

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, _StatusDoReqErr, err
	}

	var body []byte

	if resp.Body != nil {
		defer resp.Body.Close()

		buffer := bufferPool.Get().(*bytes.Buffer)
		buffer.Reset()
		defer func() {
			if buffer != nil {
				buffer.Reset()
				bufferPool.Put(buffer)
				buffer = nil
			}
		}()

		if _, err := io.Copy(buffer, resp.Body); err != nil {
			return nil, _StatusReadRespErr, fmt.Errorf("read resp body from [%s %s] err: %s", method, url, err.Error())
		}

		body = buffer.Bytes()

		if resp.StatusCode != http.StatusOK {
			return nil, resp.StatusCode, newReplyErr(
				resp.StatusCode,
				body,
				fmt.Errorf("do [%s %s] return code: %d message: %s", method, url, resp.StatusCode, string(body)),
			)
		}

	}

	return body, http.StatusOK, nil
}

// addFormValuesIntoURL append url.Values into url string
func addFormValuesIntoURL(rawURL string, form url.Values) (string, error) {
	if len(rawURL) == 0 {
		return "", errors.New("rawURL required")
	}
	if len(form) == 0 {
		return "", errors.New("form required")
	}

	target, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse rawURL `%s` err", rawURL)
	}

	urlValues := target.Query()
	for key, values := range form {
		for _, value := range values {
			urlValues.Add(key, value)
		}
	}

	target.RawQuery = urlValues.Encode()
	return target.String(), nil
}
