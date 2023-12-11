package sms

type SmsSender interface {
	Send(mobile, content string, param map[string]string) error
	SendVerifyCode(mobile, code string) error
}
