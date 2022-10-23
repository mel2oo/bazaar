package router

import (
	"encoding/json"
	"fmt"
)

const StatusOk = iota

const (
	ErrParamVerify = 1000 + iota
	ErrSampleUpload
)

var codeText = map[int]string{
	StatusOk: "成功",

	ErrParamVerify:  "参数验证失败",
	ErrSampleUpload: "上传样本失败",
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
