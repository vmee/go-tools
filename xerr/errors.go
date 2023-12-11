package xerr

import (
	"fmt"
)

type CodeError struct {
	status  uint32
	errCode string
	errMsg  string
}

// status
func (e *CodeError) GetStatus() uint32 {
	return e.status
}

//errorCode
func (e *CodeError) GetErrCode() string {
	return e.errCode
}

//errMsg
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("Status:%d, ErrCode:%s，ErrMsg:%s", e.status, e.errCode, e.errMsg)
}

func NewErr(status uint32, errCode string, errMsg string) *CodeError {
	return &CodeError{status: status, errCode: errCode, errMsg: errMsg}
}
func NewErrCode(errCode string) *CodeError {
	return &CodeError{status: MapErrStatus(errCode), errCode: errCode, errMsg: MapErrMsg(errCode)}
}

func NewSysErr(errMsg string) *CodeError {
	return &CodeError{status: MapErrStatus(SysError), errCode: SysError, errMsg: errMsg}
}

func NewBizErr(errMsg string) *CodeError {
	return &CodeError{status: MapErrStatus(BizError), errCode: BizError, errMsg: errMsg}
}

func NewParamErr(errMsg string) *CodeError {
	return &CodeError{status: MapErrStatus(ParamError), errCode: ParamError, errMsg: errMsg}
}

func NewAuthErr(errMsg string) *CodeError {
	return &CodeError{status: MapErrStatus(AuthError), errCode: AuthError, errMsg: errMsg}
}

func NewAuthForbiddenErr(errMsg string) *CodeError {
	return &CodeError{status: MapErrStatus(AuthForbiddenError), errCode: AuthForbiddenError, errMsg: errMsg}
}
