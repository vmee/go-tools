package result

import "github.com/vmee/go-tools/xerr"

type ResponseSuccess struct {
	Status uint32      `json:"status"`
	Code   string      `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
type NullJson struct{}

func Success(data interface{}) *ResponseSuccess {
	return &ResponseSuccess{200, xerr.OK, xerr.MapErrMsg(xerr.OK), data}
}

type ResponseError struct {
	Status uint32 `json:"status"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
}

func Error(status uint32, errCode string, errMsg string) *ResponseError {
	return &ResponseError{status, errCode, errMsg}
}
