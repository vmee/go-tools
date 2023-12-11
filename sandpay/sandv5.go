package sandpay

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/vmee/go-tools/orderno"
	"github.com/vmee/go-tools/sandpay/sandutil"
)

type sandV5Client struct {
	userFlag   string
	merNo      string // 商户号
	notifyUrl  string
	returnUrl  string
	prvKeyFile string
	pubKeyFile string
	prvKey     *rsa.PrivateKey
	pubKey     *rsa.PublicKey
	// cli    *http.Client
}

func (sd *sandV5Client) GetSandUserId(uid string) string {
	return sd.userFlag + uid
}

// 开通云账户
func (sd *sandV5Client) GetC2BUrl(orderNo, goodsName, money, userId, nickname string) string {
	baseUrl := "https://faspay-oss.sandpay.com.cn/pay/h5/cloud"

	productCode := "04010001"
	payExtra := `{"userId":"` + sd.GetSandUserId(userId) + `","nickName":"` + nickname + `","accountType":"1"}`

	// 请求参数
	params := sd.GlobalParams(productCode, orderNo, goodsName, money, payExtra)

	// 参与签名的参数
	signParams := url.Values{}
	for k, v := range params {
		signParams[k] = v
	}

	// 去除不需要签名的字符
	signParams.Del("expire_time")
	signParams.Del("goods_name")
	signParams.Del("product_code")
	signParams.Del("clear_cycle")
	signParams.Del("jump_scheme")
	signParams.Del("meta_option")
	signParams.Del("limit_pay")
	signParams.Del("extend_params")

	// 签名
	sign, err := sandutil.SignSand(sd.prvKey, sandutil.SignStr(signParams))
	if err != nil {
		log.Fatal(err)
	}
	params.Set("sign", sign)

	// 拼装URL
	return baseUrl + "?" + params.Encode()

}

// 开通云账户
func (sd *sandV5Client) GetCloudUrlByUserId(userId, nickname, realName, idCard string) string {
	baseUrl := "https://faspay-oss.sandpay.com.cn/pay/h5/cloud"

	productCode := "00000001"
	payExtra := `{"userId":"` + userId + `","nickName":"` + nickname + `","accountType":"1"}`

	orderNo := orderno.Generate()
	// 请求参数
	params := sd.GlobalParams(productCode, orderNo, "开通杉德云账户", "0", payExtra)

	// 参与签名的参数
	signParams := url.Values{}
	for k, v := range params {
		signParams[k] = v
	}

	// 去除不需要签名的字符
	signParams.Del("expire_time")
	signParams.Del("goods_name")
	signParams.Del("product_code")
	signParams.Del("clear_cycle")
	signParams.Del("jump_scheme")
	signParams.Del("meta_option")
	signParams.Del("limit_pay")
	signParams.Del("extend_params")

	// 签名
	sign, err := sandutil.SignSand(sd.prvKey, sandutil.SignStr(signParams))
	if err != nil {
		log.Fatal(err)
	}
	params.Set("sign", sign)

	// 拼装URL
	return baseUrl + "?" + params.Encode()

}

// 开通云账户
func (sd *sandV5Client) GetCloudUrl(userId, nickname, realName, idCard string) string {
	baseUrl := "https://faspay-oss.sandpay.com.cn/pay/h5/cloud"

	productCode := "00000001"
	payExtra := `{"userId":"` + sd.GetSandUserId(userId) + `","nickName":"` + nickname + `|` + realName + `|` + idCard + `","accountType":"1"}`

	orderNo := orderno.Generate()
	// 请求参数
	params := sd.GlobalParams(productCode, orderNo, "开通杉德云账户", "0", payExtra)

	// 参与签名的参数
	signParams := url.Values{}
	for k, v := range params {
		signParams[k] = v
	}

	// 去除不需要签名的字符
	signParams.Del("expire_time")
	signParams.Del("goods_name")
	signParams.Del("product_code")
	signParams.Del("clear_cycle")
	signParams.Del("jump_scheme")
	signParams.Del("meta_option")
	signParams.Del("limit_pay")
	signParams.Del("extend_params")

	// 签名
	sign, err := sandutil.SignSand(sd.prvKey, sandutil.SignStr(signParams))
	if err != nil {
		log.Fatal(err)
	}
	params.Set("sign", sign)

	// 拼装URL
	return baseUrl + "?" + params.Encode()

}

func (sd *sandV5Client) GetFastPaymentUrl(orderNo, goodsName, money, userId string) string {
	baseUrl := "https://sandcash.mixienet.com.cn/pay/h5/fastpayment"

	productCode := "05030001"
	payExtra := `{"userId":"` + sd.GetSandUserId(userId) + `"}`

	// 请求参数
	params := sd.GlobalParams(productCode, orderNo, goodsName, money, payExtra)

	// 参与签名的参数
	signParams := url.Values{}
	for k, v := range params {
		signParams[k] = v
	}

	// 去除不需要签名的字符
	signParams.Del("expire_time")
	signParams.Del("goods_name")
	signParams.Del("product_code")
	signParams.Del("clear_cycle")
	signParams.Del("jump_scheme")
	signParams.Del("meta_option")
	signParams.Del("limit_pay")
	signParams.Del("extend_params")

	// 签名
	sign, err := sandutil.SignSand(sd.prvKey, sandutil.SignStr(signParams))
	if err != nil {
		log.Fatal(err)
	}
	params.Set("sign", sign)

	// 拼装URL
	return baseUrl + "?" + params.Encode()

}

func (sd *sandV5Client) GlobalParams(productCode, orderNo, goodName, money, payExtra string) url.Values {
	params := url.Values{}

	params.Set("version", "10")
	params.Set("mer_no", sd.merNo)      //商户号
	params.Set("mer_order_no", orderNo) // 商户唯一订单号
	params.Set("create_time", time.Now().Format("20060102150405"))
	params.Set("expire_time", time.Now().Add(5*time.Minute).Format("20060102150405"))
	params.Set("order_amt", money)                                                                           //订单支付金额
	params.Set("notify_url", sd.notifyUrl)                                                                   //订单支付异步通知
	params.Set("return_url", sd.returnUrl)                                                                   //订单前端页面跳转地址
	params.Set("create_ip", "1_1_1_1")                                                                       //客户端IP
	params.Set("goods_name", goodName)                                                                       //订单前端页面跳转地址
	params.Set("store_id", "000000")                                                                         //门店号
	params.Set("product_code", productCode)                                                                  //产品编码
	params.Set("clear_cycle", "3")                                                                           //清算模式
	params.Set("pay_extra", payExtra)                                                                        //支付扩展域 json string
	params.Set("accsplit_flag", "NO")                                                                        //分账标识
	params.Set("jump_scheme", "sandcash://scpay")                                                            //跳转scheme
	params.Set("meta_option", `[{"s":"Android","n":"","id":"","sc":""},{"s":"IOS","n":"","id":"","sc":""}]`) //终端/网站参数
	params.Set("sign_type", "RSA")

	return params
}

func NewSandV5Client(options ...SandV5Option) (*sandV5Client, error) {

	c := &sandV5Client{}
	for _, f := range options {
		f(c)
	}

	if c.prvKeyFile != "" {
		prvKey := sandutil.LoadPrivateKey(c.prvKeyFile)
		if prvKey == nil {
			return nil, fmt.Errorf("私钥未取到")
		}
		c.prvKey = prvKey

	}

	if c.pubKeyFile != "" {
		pubKey := sandutil.LoadPublicKey(c.pubKeyFile)
		if pubKey == nil {
			return nil, fmt.Errorf("公钥未取到")
		}
		c.pubKey = pubKey
	}

	return c, nil

}
