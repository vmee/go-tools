package result

import (
	"fmt"
	"net/http"

	"github.com/vmee/go-tools/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

//http response
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		//return success
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		errStatus := http.StatusBadRequest
		errCode := xerr.SysError
		//default error msg
		errMsg := xerr.MapErrMsg(xerr.SysError)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*xerr.CodeError); ok {
			//custom error
			errStatus = int(e.GetStatus())
			errCode = e.GetErrCode()
			errMsg = e.GetErrMsg()

			if errCode == xerr.SysError {
				logx.Errorf("【APP-ERR】 : %T,  %+v ", errors.Cause(err), err)
			} else {
				logx.Infof("【APP-ERR-INFO】 : %T,  %+v ", errors.Cause(err), err)
			}
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				//grpc error by uint32 convert
				grpcCode := fmt.Sprint(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) {
					errCode = grpcCode
					errMsg = gstatus.Message()
				}
			} else {
				errMsg = err.Error()
			}
			logx.WithContext(r.Context()).Errorf("【APP-ERR】 : %T, %+v ", err, err)
		}

		httpx.WriteJson(w, errStatus, Error(uint32(errStatus), errCode, errMsg))
	}
}

//http auth error
func AuthHttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		//return success
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//return error
		errCode := xerr.AuthError
		errMsg := xerr.MapErrMsg(xerr.AuthError)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*xerr.CodeError); ok {
			errCode = e.GetErrCode()
			errMsg = e.GetErrMsg()
			logx.WithContext(r.Context()).Infof("【AUTH-ERR】 : %T, %+v ", err, err)
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				grpcCode := fmt.Sprint(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) {
					errCode = grpcCode
					errMsg = gstatus.Message()
				}
			}
			logx.WithContext(r.Context()).Errorf("【AUTH-ERR】 : %T, %+v ", err, err)
		}

		httpx.WriteJson(w, http.StatusUnauthorized, Error(http.StatusUnauthorized, errCode, errMsg))
	}
}

//http param error
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s ,%s", xerr.MapErrMsg(xerr.ParamError), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Error(http.StatusBadRequest, xerr.ParamError, errMsg))
}

//
func GlobalErrorHandler(err error) (int, interface{}) {

	errCode := xerr.SysError
	//default error msg
	var errMsg string

	errStatus := http.StatusInternalServerError

	switch e := err.(type) {
	case *xerr.CodeError:
		errStatus = int(e.GetStatus())
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
		if errCode == xerr.SysError {
			logx.Errorf("【GLOBAL-ERR】 : %T,  %+v ", errors.Cause(err), err)
		} else {
			logx.Infof("【GLOBAL-ERR】 : %T,  %+v ", errors.Cause(err), err)
		}
	default:
		errMsg = err.Error()
		logx.Errorf("【GLOBAL-ERR】 : %T,  %+v ", errors.Cause(err), err)
	}

	return errStatus, Error(uint32(errStatus), errCode, errMsg)
}
