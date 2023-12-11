package captchax

import (
	"strings"
	"time"

	"github.com/vmee/go-tools/sms"
	"github.com/vmee/go-tools/tool"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type SmsCaptchaX struct {
	smsSender  sms.SmsSender
	bizRedis   *redis.Redis
	expiration time.Duration
}

func NewSmsCaptchaX(sender sms.SmsSender, redis *redis.Redis, expiration time.Duration) *SmsCaptchaX {
	return &SmsCaptchaX{
		smsSender:  sender,
		bizRedis:   redis,
		expiration: expiration,
	}
}

func (sc *SmsCaptchaX) Send(mobile string) error {

	code := tool.Krand(6, tool.KC_RAND_KIND_NUM)

	if err := sc.smsSender.SendVerifyCode(mobile, code); err != nil {
		return err
	}

	if err := sc.bizRedis.Setex("cache:captcha:sms:"+mobile, code, int(sc.expiration/time.Second)); err != nil {
		return err
	}

	return nil
}

func (sc *SmsCaptchaX) Verify(mobile, answer string, clear bool) bool {

	if mobile == "" || answer == "" {
		return false
	}

	v, err := sc.bizRedis.Get("cache:captcha:sms:" + mobile)
	if err != nil {
		logx.Error(err)
	}
	if clear {
		sc.bizRedis.Del("cache:captcha:sms:" + mobile)
	}

	return strings.EqualFold(v, answer)
}
