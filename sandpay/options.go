package sandpay

import "net/http"

// ClientOption 客户端配置项
type ClientOption func(c *client)

// WithHTTPClient 自定义http.Client
func WithHTTPClient(cli *http.Client) ClientOption {
	return func(c *client) {
		c.cli = cli
	}
}

// HeadOption 报文头配置项
type HeadOption func(h X)

// WithVersion 设置版本号：默认：1.0；功能产品号为微信小程序或支付宝生活号，对账单需获取营销优惠金额字段传：3.0
func WithVersion(v string) HeadOption {
	return func(h X) {
		h["version"] = v
	}
}

// WithPLMid 设置平台ID：接入类型为2时必填，在担保支付模式下填写核心商户号；在杉德宝平台终端模式下填写平台商户号
func WithPLMid(id string) HeadOption {
	return func(h X) {
		h["plMid"] = id
	}
}

// WithAccessType 设置接入类型：1 - 普通商户接入（默认）；2 - 平台商户接入
func WithAccessType(at string) HeadOption {
	return func(h X) {
		h["accessType"] = at
	}
}

// WithChannelType 设置渠道类型：07 - 互联网（默认）；08 - 移动端
func WithChannelType(ct string) HeadOption {
	return func(h X) {
		h["channelType"] = ct
	}
}
