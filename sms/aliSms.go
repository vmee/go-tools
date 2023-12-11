package sms

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zeromicro/go-zero/core/logx"
)

type AliSms struct {
	accessKeyId        *string
	accessKeySecret    *string
	signName           *string
	verifyTemplateCode *string
	client             *dysmsapi20170525.Client
}

func NewAliSms(accessKeyId, accessKeySecret, signName, verifyTemplateCode string) *AliSms {
	aliSms := &AliSms{
		accessKeyId:        &accessKeyId,
		accessKeySecret:    &accessKeySecret,
		signName:           &signName,
		verifyTemplateCode: &verifyTemplateCode,
	}

	c, err := aliSms.createClient()
	if err != nil {
		logx.Errorf("阿里短信连接出错, %+v", err)
	}

	aliSms.client = c

	return aliSms
}

func (ali *AliSms) Send(mobile, content string, param map[string]string) error {
	if ali.client == nil {
		return fmt.Errorf("阿里短信连接出错")
	}

	return nil
}

func (ali *AliSms) SendVerifyCode(mobile, code string) error {
	if ali.client == nil {
		return fmt.Errorf("阿里短信连接出错")
	}

	codeJson := fmt.Sprintf("{\"code\": \"%s\"}", code)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &mobile,
		SignName:      ali.signName,
		TemplateCode:  ali.verifyTemplateCode,
		TemplateParam: tea.String(codeJson),
	}

	runtime := &util.RuntimeOptions{}
	resp, _err := ali.client.SendSmsWithOptions(sendSmsRequest, runtime)
	if _err != nil {
		return _err
	}

	if *resp.Body.Code != "OK" {
		logx.Infof("短信发送失败: %+v", resp)
	}

	return nil
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func (ali *AliSms) createClient() (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: ali.accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: ali.accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
