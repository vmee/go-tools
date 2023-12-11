package xerr

var message map[string]string

func init() {
	message = make(map[string]string)
	message[OK] = "操作成功"
	message[SysError] = "服务器开小差啦,稍后再来试一试"
	message[BizError] = "业务错误"
	message[ParamError] = "参数错误"
	message[AuthError] = "认证失败"
	message[AuthForbiddenError] = "无权限访问"
	message[DbError] = "数据库繁忙,请稍后再试"
	message[DbUpdateAffectedZeroError] = "更新数据影响行数为0"
	message[DataNoExistError] = "数据不存在"

	// message[BadRequestError] = "请求失败"
	// message[UnauthorizedError] = "认证失败"
	// message[ForbiddenError] = "无权限访问"
	// message[ServerCommonError] = "服务器开小差啦,稍后再来试一试"
	// message[ReuqestParamError] = "参数错误"
	// message[TokenExpireError] = "token失效，请重新登陆"
	// message[TokenGenerateError] = "生成token失败"
	// message[DbError] = "数据库繁忙,请稍后再试"
	// message[DbUpdateAffectedZeroError] = "更新数据影响行数为0"
	// message[DataNoExistError] = "数据不存在"
}

func MapErrMsg(errcode string) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return message[SysError]
	}
}

func IsCodeErr(errcode string) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
