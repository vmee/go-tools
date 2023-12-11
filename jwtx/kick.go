package jwtx

import (
	"strconv"

	"github.com/vmee/go-tools/tool"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	cacheTokenKey = "cache:login:token"
)

type Kicker struct {
	store *redis.Redis
}

// 获取互踢控制器
func NewKicker(redis *redis.Redis) *Kicker {
	return &Kicker{
		store: redis,
	}
}

func (l *Kicker) getKey(client string) string {
	return cacheTokenKey + ":" + client
}

func (l *Kicker) Save(uid int64, token, client string) error {

	// 缓存最新token
	md5, _ := tool.Md5ByString(token)
	return l.store.Hset(l.getKey(client), strconv.FormatInt(uid, 10), md5)

}

func (l *Kicker) Verify(uid int64, token, client string) bool {

	md5, err := tool.Md5ByString(token)
	if err != nil || md5 == "" {
		return true
	}

	tokMd5, err := l.store.Hget(l.getKey(client), strconv.FormatInt(uid, 10))
	if err != nil || tokMd5 == "" {
		return true
	}

	return md5 == tokMd5
}
