package captchax

import (
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

//configJsonBody json request body.
type Base64CaptchaConfig struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var store = base64Captcha.DefaultMemStore

type Base64CaptchaX struct {
	Store base64Captcha.Store
}

func NewBase64CaptchaX(redis *redis.Redis, expiration time.Duration) *Base64CaptchaX {
	return &Base64CaptchaX{
		Store: NewRedisStore(redis, expiration),
	}
}

// base64Captcha create http handler
func (bc *Base64CaptchaX) GenerateBase64Captcha(param Base64CaptchaConfig) (id, b64s string, err error) {
	//parse request parameters

	var driver base64Captcha.Driver

	//create base64 encoding captcha
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}
	c := base64Captcha.NewCaptcha(driver, bc.Store)
	return c.Generate()
}

type redisStore struct {
	BizRedis   *redis.Redis
	Expiration time.Duration
}

func NewRedisStore(redis *redis.Redis, expiration time.Duration) *redisStore {
	return &redisStore{
		BizRedis:   redis,
		Expiration: expiration,
	}
}

func (r *redisStore) Set(id string, value string) error {
	return r.BizRedis.Setex("cache:captcha:img:"+id, value, int(r.Expiration/time.Second))
}

func (r *redisStore) Get(id string, clear bool) string {
	v, err := r.BizRedis.Get("cache:captcha:img:" + id)
	if err != nil {
		logx.Error(err)
	}
	if clear {
		r.BizRedis.Del("cache:captcha:img:" + id)
	}

	return v
}

func (r *redisStore) Verify(id, answer string, clear bool) bool {

	if id == "" || answer == "" {
		return false
	}

	v := r.Get(id, clear)
	return strings.EqualFold(v, answer)
}
