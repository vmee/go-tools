package limitx

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Limiter struct {
	store *redis.Redis

	periodLimits map[string]*limit.PeriodLimit
}

// 获取限流器
func NewLimiter(redis *redis.Redis) *Limiter {
	return &Limiter{
		store:        redis,
		periodLimits: make(map[string]*limit.PeriodLimit),
	}
}

func (l *Limiter) PeriodLimit(period, quota int) *limit.PeriodLimit {

	mk := fmt.Sprintf("%d:%d", period, quota)
	if p, ok := l.periodLimits[mk]; ok {
		return p
	}

	p := limit.NewPeriodLimit(period, quota, l.store, "starts:")
	l.periodLimits[mk] = p

	return p
}

func (l *Limiter) Take(period, quota int, key string) bool {
	n, err := l.PeriodLimit(period, quota).Take(key)
	if err != nil {
		return true
	}

	return n == 1
}

func (l *Limiter) TakeCtx(ctx context.Context, period, quota int, key string) bool {
	n, err := l.PeriodLimit(period, quota).TakeCtx(ctx, key)
	if err != nil {
		logx.Error("limitx take err:", err)
		return true
	}

	return n == limit.Allowed || n == limit.HitQuota
}
