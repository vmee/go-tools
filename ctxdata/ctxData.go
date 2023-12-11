package ctxdata

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

// CtxKeyJwtUserId get uid from ctx
var CtxKeyJwtUserId = "uid"

func GetUidFromCtx(ctx context.Context) uint64 {
	var uid uint64
	if jsonUid, ok := ctx.Value(CtxKeyJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = uint64(int64Uid)
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}
	return uid
}

func IsApp(ctx context.Context) bool {
	isApp, ok := ctx.Value("x-is-app").(string)
	if !ok {
		return false
	}

	return isApp == "1"
}
