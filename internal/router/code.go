package router

import (
	"encoding/json"
	"fmt"
)

const StatusOk = iota

const (
	ErrParamVerify = 1000 + iota
	ErrSampleUpload
	ErrSampleQuery
	ErrSampleCount
	ErrSampleDownload
	ErrSampleExists
	ErrSampleReader
)

var codeText = map[int]string{
	StatusOk: "success",

	ErrParamVerify:    "parameter validation failed",
	ErrSampleUpload:   "upload malware failed",
	ErrSampleQuery:    "query malware failed",
	ErrSampleCount:    "count malware failed",
	ErrSampleDownload: "download malware failed",
	ErrSampleExists:   "malware file does not exist",
	ErrSampleReader:   "malware file read failed",
}

type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func NewReply(code int) *Reply {
	return &Reply{
		Code: code,
		Msg:  codeText[code],
		Data: []string{},
	}
}

func (r *Reply) WithData(data interface{}) *Reply {
	if data != nil {
		d, ok := data.([]byte)
		if ok {
			json.Unmarshal(d, &r.Data)
		} else {
			r.Data = data
		}
	}
	return r
}

func (r *Reply) WithErr(err error) *Reply {
	if err != nil {
		if len(r.Msg) > 0 {
			r.Msg = fmt.Sprintf("%s: %s", r.Msg, err.Error())
		} else {
			r.Msg = err.Error()
		}
	}
	return r
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
