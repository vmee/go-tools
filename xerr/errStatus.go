package xerr

var statusMap map[string]uint32

func init() {
	statusMap = make(map[string]uint32)
	statusMap[OK] = 200
	statusMap[SysError] = 500
	statusMap[BizError] = 400
	statusMap[ParamError] = 400
	statusMap[AuthError] = 401
	statusMap[AuthForbiddenError] = 403

}

func MapErrStatus(errCode string) uint32 {
	if msg, ok := statusMap[errCode]; ok {
		return msg
	} else {
		return statusMap[SysError]
	}
}

func IsErrStatus(errCode string) bool {
	if _, ok := statusMap[errCode]; ok {
		return true
	} else {
		return false
	}
}
