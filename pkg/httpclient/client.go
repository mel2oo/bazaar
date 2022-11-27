package httpclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	httpURL "net/url"
	"time"

	json "github.com/json-iterator/go"
)

const (
	// DefaultTTL 一次http请求最长执行1分钟
	DefaultTTL = time.Minute
)

// Get get 请求
func Get(url string, form httpURL.Values, options ...Option) (body []byte, err error) {
	return withoutBody(http.MethodGet, url, form, options...)
}

// Delete delete 请求
func Delete(url string, form httpURL.Values, options ...Option) (body []byte, err error) {
	return withoutBody(http.MethodDelete, url, form, options...)
}

func withoutBody(method, url string, form httpURL.Values, options ...Option) (body []byte, err error) {
	if len(url) == 0 {
		return nil, errors.New("url required")
	}

	if len(form) > 0 {
		if url, err = addFormValuesIntoURL(url, form); err != nil {
			return
		}
	}

	opt := getOption()
	defer releaseOption(opt)

	for _, f := range options {
		f(opt)
	}
	opt.header["Content-Type"] = []string{"application/x-www-form-urlencoded; charset=utf-8"}

	ttl := opt.ttl
	if ttl <= 0 {
		ttl = DefaultTTL
	}

	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()

	retryTimes := opt.retryTimes
	if retryTimes <= 0 {
		retryTimes = DefaultRetryTimes
	}

	retryDelay := opt.retryDelay
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelay
	}

	var httpCode int

	for k := 0; k < retryTimes; k++ {
		body, httpCode, err = doHTTP(ctx, method, url, nil, opt)
		if shouldRetry(ctx, httpCode) || (opt.retryVerify != nil && opt.retryVerify(body)) {
			time.Sleep(retryDelay)
			continue
		}

		return
	}
	return
}

// PostForm post form 请求
func PostForm(url string, form httpURL.Values, options ...Option) (body []byte, err error) {
	return withFormBody(http.MethodPost, url, form, options...)
}

func PostFormFile(url string, form httpURL.Values, files FileValues, options ...Option) (body []byte, err error) {
	return withFormFileBody(http.MethodPost, url, form, files, options...)
}

// PostJSON post json 请求
func PostJSON(url string, raw json.RawMessage, options ...Option) (body []byte, err error) {
	return withJSONBody(http.MethodPost, url, raw, options...)
}

// PutForm put form 请求
func PutForm(url string, form httpURL.Values, options ...Option) (body []byte, err error) {
	return withFormBody(http.MethodPut, url, form, options...)
}

// PutJSON put json 请求
func PutJSON(url string, raw json.RawMessage, options ...Option) (body []byte, err error) {
	return withJSONBody(http.MethodPut, url, raw, options...)
}

// PatchFrom patch form 请求
func PatchFrom(url string, form httpURL.Values, options ...Option) (body []byte, err error) {
	return withFormBody(http.MethodPatch, url, form, options...)
}

// PatchJSON patch json 请求
func PatchJSON(url string, raw json.RawMessage, options ...Option) (body []byte, err error) {
	return withJSONBody(http.MethodPatch, url, raw, options...)
}

func withFormBody(method, url string, form httpURL.Values, options ...Option) (body []byte, err error) {
	if len(url) == 0 {
		return nil, errors.New("url required")
	}

	opt := getOption()
	defer releaseOption(opt)

	for _, f := range options {
		f(opt)
	}
	opt.header["Content-Type"] = []string{"application/x-www-form-urlencoded; charset=utf-8"}

	ttl := opt.ttl
	if ttl <= 0 {
		ttl = DefaultTTL
	}

	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()

	var formValue string

	if len(form) > 0 {
		formValue = form.Encode()
	}

	retryTimes := opt.retryTimes
	if retryTimes <= 0 {
		retryTimes = DefaultRetryTimes
	}

	retryDelay := opt.retryDelay
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelay
	}

	var httpCode int

	for k := 0; k < retryTimes; k++ {
		body, httpCode, err = doHTTP(ctx, method, url, bytes.NewReader([]byte(formValue)), opt)
		if shouldRetry(ctx, httpCode) || (opt.retryVerify != nil && opt.retryVerify(body)) {
			time.Sleep(retryDelay)
			continue
		}

		return
	}
	return
}

func withFormFileBody(method, url string, form httpURL.Values, files FileValues, options ...Option) (body []byte, err error) {
	if len(url) == 0 {
		return nil, errors.New("url required")
	}

	opt := getOption()
	defer releaseOption(opt)

	for _, f := range options {
		f(opt)
	}

	ttl := opt.ttl
	if ttl <= 0 {
		ttl = DefaultTTL
	}

	retryTimes := opt.retryTimes
	if retryTimes <= 0 {
		retryTimes = DefaultRetryTimes
	}

	retryDelay := opt.retryDelay
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelay
	}

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	for field, file := range files {
		fw, err := bodyWriter.CreateFormFile(field, file.Name)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fw, file.Reader)
		if err != nil {
			return nil, err
		}
	}

	if len(form) > 0 {
		for k, v := range form {
			bodyWriter.WriteField(k, v[0])
		}
	}

	opt.header["Content-Type"] = []string{bodyWriter.FormDataContentType()}
	bodyWriter.Close()

	var httpCode int

	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()

	for k := 0; k < retryTimes; k++ {
		body, httpCode, err = doHTTP(ctx, method, url, bodyBuffer, opt)
		if shouldRetry(ctx, httpCode) || (opt.retryVerify != nil && opt.retryVerify(body)) {
			time.Sleep(retryDelay)
			continue
		}
		return
	}
	return
}

func withJSONBody(method, url string, raw json.RawMessage, options ...Option) (body []byte, err error) {
	if len(url) == 0 {
		return nil, errors.New("url required")
	}
	if len(raw) == 0 {
		return nil, errors.New("raw required")
	}

	opt := getOption()
	defer releaseOption(opt)

	for _, f := range options {
		f(opt)
	}
	opt.header["Content-Type"] = []string{"application/json; charset=utf-8"}

	ttl := opt.ttl
	if ttl <= 0 {
		ttl = DefaultTTL
	}

	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()

	retryTimes := opt.retryTimes
	if retryTimes <= 0 {
		retryTimes = DefaultRetryTimes
	}

	retryDelay := opt.retryDelay
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelay
	}

	var httpCode int

	for k := 0; k < retryTimes; k++ {
		body, httpCode, err = doHTTP(ctx, method, url, bytes.NewBuffer(raw), opt)
		if shouldRetry(ctx, httpCode) || (opt.retryVerify != nil && opt.retryVerify(body)) {
			time.Sleep(retryDelay)
			continue
		}

		return
	}
	return
}
