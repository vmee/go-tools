package xerr

// 成功返回
const OK string = "ok"

const ParamError string = "P_Param_Err"
const BizError string = "B_Biz_Err"
const SysError string = "S_Sys_Err"
const AuthError string = "A_Auth_Err"
const AuthForbiddenError string = "A_Forbidden_Err"

const DbError string = "D_Db_Err"
const DbUpdateAffectedZeroError string = "D_UpdateAffectedZero_Err"
const DataNoExistError string = "D_DataNoExistError_Err"

// 用于客户端跳转
const RealNameError string = "300001" //未实名认证

// const BadRequestError uint32 = 400
// const UnauthorizedError uint32 = 401
// const ForbiddenError uint32 = 403
// const NotFoundError uint32 = 404

// /**(前3位代表业务,后三位代表具体功能)**/

// /**全局错误码*/
// //服务器开小差
// const ServerCommonError uint32 = 100001

// //请求参数错误
// const ReuqestParamError uint32 = 100002

// //token过期
// const TokenExpireError uint32 = 100003

// //生成token失败
// const TokenGenerateError uint32 = 100004

// //数据库繁忙,请稍后再试
// const DbError uint32 = 100005

// //更新数据影响行数为0
// const DbUpdateAffectedZeroError uint32 = 100006

// //数据不存在
// const DataNoExistError uint32 = 100007

//用户服务

//订单服务

//商品服务

//支付服务
