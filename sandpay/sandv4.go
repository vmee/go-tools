package sandpay

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/vmee/go-tools/orderno"
	"github.com/vmee/go-tools/sandpay/sandutil"
)

type sandV4Client struct {
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

func NewSandV4Client(options ...SandV4Option) (*sandV4Client, error) {

	c := &sandV4Client{}
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

func (sd *sandV4Client) GetSandUserId(uid string) string {
	return sd.userFlag + uid
}

func (sd *sandV4Client) VerifyNotificationSign(data, signStr string) bool {

	sign, err := base64.StdEncoding.DecodeString(strings.Replace(signStr, " ", "+", -1))
	if err != nil {
		return false
	}

	err = sandutil.SandVerification([]byte(data), sign, sd.pubKey)
	if err != nil {
		fmt.Println("验签失败", err)
		return false
	}

	return true
}

type PostData struct {
	Charset  string `json:"charset"`
	SignType string `json:"signType"`
	Data     string `json:"data"`
	Sign     string `json:"sign"`
}

type Header struct {
	Version    string `json:"version"`
	Method     string `json:"method"`
	Mid        string `json:"mid"`
	AccessType string `json:"accessType"`
	//PlMid       string `json:"plMid"`
	ChannelType string `json:"channelType"`
	ReqTime     string `json:"reqTime"`
	ProductId   string `json:"productId"`
}

type OrderQueryBody struct {
	//商户订单号
	OrderCode string `json:"orderCode"`
	//19. 扩展域
	Extends string `json:"extends"`
}

type Response struct {
	Charset  string `json:"charset"`
	Data     string `json:"data"`
	SignType string `json:"signType"`
	Body     string `json:"body"`
}

type OrderQueryResponseData struct {
	Head *OrderQueryResponseHeadData `json:"head"`
	Body *OrderQueryResponseBodyData `json:"body"`
}

type OrderQueryResponseHeadData struct {
	RespCode string `json"respCode"`
}

type OrderQueryResponseBodyData struct {
	OrderStatus string `json"orderStatus"`
}

// 订单查询接口
func (sd *sandV4Client) OrderQueryIsSuccess(orderNo string, extend string) (b bool, err error) {
	vals, err := sd.OrderQuery(orderNo, extend)
	if err != nil {
		return
	}

	data := vals.Get("data")
	sign := vals.Get("sign")

	if !sd.VerifyNotificationSign(data, sign) {
		err = fmt.Errorf("签名验证失败")
	}

	responseData := &OrderQueryResponseData{}
	err = json.Unmarshal([]byte(data), &responseData)
	if err != nil {
		return
	}

	if responseData.Head.RespCode == "000000" && responseData.Body.OrderStatus == "00" {
		b = true
	} else {
		b = false
	}

	return

}

// 订单查询接口
func (sd *sandV4Client) OrderQuery(orderNo string, extend string) (vals url.Values, err error) {
	timeString := time.Now().Format("20060102150405")

	header := Header{
		Method:      "sandpay.trade.query",
		Mid:         sd.merNo,
		Version:     "1.0",
		AccessType:  "1",
		ProductId:   "00000016",
		ChannelType: "08",
		ReqTime:     timeString,
	}

	body := OrderQueryBody{
		OrderCode: orderNo,
		Extends:   extend,
	}

	signDataJsonString := GenerateSignString(body, header)
	sign, _ := sandutil.SignSand(sd.prvKey, signDataJsonString)

	// postData := PostData{
	// 	Charset:  "utf-8",
	// 	SignType: "01",
	// 	Data:     signDataJsonString,
	// 	Sign:     sign,
	// }

	// postDataString, _ := json.Marshal(postData)

	// postData := make(map[string]string, 4)
	// postData["chart"] = `utf-8`
	// postData["signType"] = `01`
	// postData["data"] = signDataJsonString
	// postData["sign"] = sign

	postData := url.Values{}
	postData.Add("chart", "utf-8")
	postData.Add("SignType", "01")
	postData.Add("data", signDataJsonString)
	postData.Add("sign", sign)

	resp, err := sandutil.DoForm("https://cashier.sandpay.com.cn/gateway/api/order/query", postData)

	dataString, _ := url.QueryUnescape(string(resp[:]))

	var fields []string
	fields = strings.Split(dataString, "&")
	// fmt.Println(fields)

	vals = url.Values{}
	for _, field := range fields {
		f := strings.SplitN(field, "=", 2)
		if len(f) >= 2 {
			key, val := f[0], f[1]
			vals.Set(key, val)
		}
	}

	// fmt.Println(vals.Encode())
	return

	// fmt.Println(err)
	// d := make(map[string]interface{})
	// if err := json.Unmarshal(resp, &d); err != nil {
	// 	fmt.Println(err.Error())
	// }
	fmt.Println("杉德回调解析结果:" + string(resp))
	return

	// fd := Params()
	// sanDe := util.SandAES{}
	// key := sanDe.RandStr(16)

	// fdata, _ := FormData(fd, key)
	// fd["encryptKey"], _ = FormEncryptKey(key, pubKey)
	// fd["sign"], _ = FormSign(fdata, prvKey)
	// fd["data"] = fdata
	// //display(fd)

	// DataByte, _ := json.Marshal(fd)
	// api := "https://ceas-uat01.sand.com.cn/v4/electrans/ceas.elec.trans.corp.transfer"

	// resp, err := util.Do(api, string(DataByte))
	// fmt.Println(err)
	// d := make(map[string]interface{})
	// if err := json.Unmarshal(resp, &d); err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println("杉德回调解析结果:" + string(resp))

	// return/
}

func GenerateSignString(body interface{}, header interface{}) (str string) {
	signData := make(map[string]interface{}, 2)
	signData["body"] = body
	signData["head"] = header
	signDataJson, err := json.Marshal(signData)
	if err != nil {
		return
	}
	return string(signDataJson[:])
}

func GeneratePostData(signDataJsonString string, sign string) map[string]string {
	postData := make(map[string]string, 4)
	postData["chart"] = `utf-8`
	postData["signType"] = `01`
	postData["data"] = signDataJsonString
	postData["sign"] = sign
	return postData
}

type AccountInfoQueryReq struct {
	Version         string `json:"version"`
	Mid             string `json:"mid"`
	CustomerOrderNo string `json:"customerOrderNo"`
	SignType        string `json:"signType"`
	EncryptType     string `json:"encryptType"`
	Timestamp       string `json:"timestamp,omitempty"`
	BizUserNo       string `json:"bizUserNo,omitempty"`
	EncryptKey      string `json:"encryptKey,omitempty"`
	Data            string `json:"data,omitempty"`
	Sign            string `json:"sign,omitempty"`
}

type AccountInfoQueryResponse struct {
	ResponseDesc      string `json:"responseDesc"`
	ResponseStatus    string `json:"responseStatus"`
	OpenAccountStatus string `json:"openAccountStatus"`
	FaceStatus        string `json:"faceStatus"`
}

func (sd *sandV4Client) AccountIsOpen(uid string) (b bool, err error) {
	data, err := sd.AccountInfoQuery(uid)
	if err != nil {
		return
	}

	resp := &AccountInfoQueryResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return
	}

	if resp.ResponseStatus != "00" {
		b = false
		return
	}

	return resp.OpenAccountStatus == "01" && resp.FaceStatus == "01", nil

}

// 订单查询接口
func (sd *sandV4Client) AccountInfoQuery(uid string) (responseData []byte, err error) {

	req := AccountInfoQueryReq{
		Version:         "1.0",
		Mid:             sd.merNo,
		CustomerOrderNo: orderno.Generate(),
		SignType:        "SHA1WithRSA",
		EncryptType:     "AES",
		Timestamp:       time.Now().Format("2006-01-02 15:04:05"),
		BizUserNo:       sd.GetSandUserId(uid),
	}

	sanDe := sandutil.SandAES{}
	key := sanDe.RandStr(16)

	dataJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	aes := sandutil.SandAES{}
	aes.Key = []byte(key)

	data := aes.Encypt5(dataJson)

	req.EncryptKey, _ = sandutil.RsaEncrypt(key, sd.pubKey)
	req.Sign, _ = sandutil.SignSand(sd.prvKey, data)
	req.Data = data

	reqJson, _ := json.Marshal(req)

	resp, err := sandutil.Do("https://cap.sandpay.com.cn/v4/elecaccount/ceas.elec.member.info.query", string(reqJson))

	if err != nil {
		return
	}

	d := make(map[string]interface{})
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, err
	}
	fmt.Println("杉德回调解析结果:" + string(resp))

	da := d["data"].(string)
	si := d["sign"].(string)

	err = sandutil.Verification(da, si, sd.pubKey)
	if err != nil {
		return nil, fmt.Errorf("验签失败, %s", err)
	}

	encryptKey := d["encryptKey"].(string)

	decryptAESKey, err := sandutil.RsaDecrypt(encryptKey, sd.prvKey)
	if err != nil {
		return
	}

	aes1 := &sandutil.SandAES{
		Key: []byte(decryptAESKey),
	}

	// daa, err := base64.StdEncoding.DecodeString(da)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// text := sandutil.AesDecrypt([]byte(daa), []byte(decryptAESKey))

	responseData, err = aes1.AesECBDecrypt(da)

	// fmt.Println(string(responseData))
	return

	//
	// {"responseDesc":"会员信息不存在","responseTime":"20230207115137","mid":"6888802118774","sandSerialNo":"CEAS2023020711513755401578","responseStatus":"01","version":"1.0","customerOrderNo":"2302071151350545960001","responseCode":"05017"}
	// String decryptKey = dataJson.getString("encryptKey");
	//     dataJson.remove("encryptKey");
	//     byte[] decryptKeyBytes = Base64.decodeBase64(decryptKey);
	//     decryptKey = new String(RSAUtils.decrypt(decryptKeyBytes, CertCache.getCertCache().getPrivateKey()));
	//     logger.info("RSA解密后随机数：{}", decryptKey);
	//     String encryptValue = dataJson.getString("data");
	//     logger.info("AES解密前值：{}", encryptValue);
	//     byte[] decryptDataBase64 = Base64.decodeBase64(encryptValue);
	//     byte[] decryptDataBytes = AESUtils.decrypt(decryptDataBase64, decryptKey.getBytes(StandardCharsets.UTF_8), (String) null);
	//     String decryptData = new String(decryptDataBytes);
	//     logger.info("AES解密后值：{}", decryptData);
	//     return JSON.parseObject(decryptData);
	// fd := Params()
	// sanDe := util.SandAES{}
	// key := sanDe.RandStr(16)

	// fdata, _ := FormData(fd, key)
	// fd["encryptKey"], _ = FormEncryptKey(key, pubKey)
	// fd["sign"], _ = FormSign(fdata, prvKey)
	// fd["data"] = fdata
	// //display(fd)

	// DataByte, _ := json.Marshal(fd)
	// api := "https://ceas-uat01.sand.com.cn/v4/electrans/ceas.elec.trans.corp.transfer"

	// resp, err := util.Do(api, string(DataByte))
	// fmt.Println(err)
	// d := make(map[string]interface{})
	// if err := json.Unmarshal(resp, &d); err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println("杉德回调解析结果:" + string(resp))

	// return/
}

type B2cReq struct {
	Version         string      `json:"version"`
	Mid             string      `json:"mid"`
	CustomerOrderNo string      `json:"customerOrderNo"`
	SignType        string      `json:"signType"`
	EncryptType     string      `json:"encryptType"`
	Timestamp       string      `json:"timestamp,omitempty"`
	EncryptKey      string      `json:"encryptKey,omitempty"`
	Data            string      `json:"data,omitempty"`
	Sign            string      `json:"sign,omitempty"`
	AccountType     string      `json:"accountType,omitempty"`
	OrderAmt        float64     `json:"orderAmt,omitempty"`
	Payee           B2cPayeeReq `json:"payee,omitempty"`
	Postscript      string      `json:"postscript,omitempty"`
	Remark          string      `json:"remark,omitempty"`
}

type B2cPayeeReq struct {
	BizUserNo string `json:"bizUserNo"`
	Name      string `json:"name"`
}

type B2cResponse struct {
	ResponseDesc   string `json:"responseDesc"`
	ResponseStatus string `json:"responseStatus"`
	OrderStatus    string `json:"OrderStatus"`
}

// b2c查询接口
func (sd *sandV4Client) B2cIsSuccess(orderNo, uid, realName string, money float64, postscript, remark string) (isSuccess bool, err error) {
	data, err := sd.B2c(orderNo, uid, realName, money, postscript, remark)
	if err != nil {
		return
	}

	resp := &B2cResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return
	}

	return resp.ResponseStatus == "00" && resp.OrderStatus == "00", nil
}

// b2c查询接口
func (sd *sandV4Client) B2c(orderNo, uid, nickName string, money float64, postscript, remark string) (responseData []byte, err error) {

	req := B2cReq{
		Version:         "1.0",
		Mid:             sd.merNo,
		CustomerOrderNo: orderNo,
		SignType:        "SHA1WithRSA",
		EncryptType:     "AES",
		Timestamp:       time.Now().Format("2006-01-02 15:04:05"),
		AccountType:     "01",
		OrderAmt:        money,
		// OrderAmt:        fmt.Sprintf("%.2f", money),
		Payee: B2cPayeeReq{
			BizUserNo: sd.GetSandUserId(uid),
			Name:      nickName,
		},
		Postscript: postscript,
		Remark:     remark,
	}

	sanDe := sandutil.SandAES{}
	key := sanDe.RandStr(16)

	dataJson, err := json.Marshal(req)
	fmt.Println(string(dataJson))
	if err != nil {
		return nil, err
	}
	aes := sandutil.SandAES{}
	aes.Key = []byte(key)

	data := aes.Encypt5(dataJson)

	req.EncryptKey, _ = sandutil.RsaEncrypt(key, sd.pubKey)
	req.Sign, _ = sandutil.SignSand(sd.prvKey, data)
	req.Data = data

	reqJson, _ := json.Marshal(req)
	// fmt.Println("请求参数:", string(reqJson))
	resp, err := sandutil.Do("https://cap.sandpay.com.cn/v4/electrans/ceas.elec.trans.corp.transfer", string(reqJson))

	if err != nil {
		return
	}

	d := make(map[string]interface{})
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, err
	}
	fmt.Println("杉德回调解析结果:" + string(resp))

	da := d["data"].(string)
	si := d["sign"].(string)

	err = sandutil.Verification(da, si, sd.pubKey)
	if err != nil {
		return nil, fmt.Errorf("验签失败, %s", err)
	}

	encryptKey := d["encryptKey"].(string)

	decryptAESKey, err := sandutil.RsaDecrypt(encryptKey, sd.prvKey)
	if err != nil {
		return
	}

	aes1 := &sandutil.SandAES{
		Key: []byte(decryptAESKey),
	}

	// daa, err := base64.StdEncoding.DecodeString(da)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// text := sandutil.AesDecrypt([]byte(daa), []byte(decryptAESKey))

	responseData, err = aes1.AesECBDecrypt(da)

	// fmt.Println(string(responseData))
	return

}

type CloseAccountReq struct {
	Version         string  `json:"version"`
	Mid             string  `json:"mid"`
	CustomerOrderNo string  `json:"customerOrderNo"`
	SignType        string  `json:"signType"`
	EncryptType     string  `json:"encryptType"`
	Timestamp       string  `json:"timestamp,omitempty"`
	EncryptKey      string  `json:"encryptKey,omitempty"`
	Data            string  `json:"data,omitempty"`
	Sign            string  `json:"sign,omitempty"`
	AccountType     string  `json:"accountType,omitempty"`
	OrderAmt        float64 `json:"orderAmt,omitempty"`
	BizUserNo       string  `json:"bizUserNo,omitempty"`
	BizType         string  `json:"bizType,omitempty"`
	Postscript      string  `json:"postscript,omitempty"`
	Remark          string  `json:"remark,omitempty"`
	NotifyUrl       string  `json:"notifyUrl,omitempty"`
	FrontUrl        string  `json:"returnUrl,omitempty"`
}

// 注销接口
func (sd *sandV4Client) CloseAccount(uid string, postscript, remark string) (responseData []byte, err error) {

	req := CloseAccountReq{
		Version:         "1.0",
		Mid:             sd.merNo,
		CustomerOrderNo: orderno.Generate(),
		SignType:        "SHA1WithRSA",
		EncryptType:     "AES",
		Timestamp:       time.Now().Format("2006-01-02 15:04:05"),
		AccountType:     "01",
		OrderAmt:        0.11,
		BizUserNo:       uid,
		BizType:         "CLOSE",
		// OrderAmt:        fmt.Sprintf("%.2f", money),

		Postscript: postscript,
		NotifyUrl:  sd.notifyUrl,
		FrontUrl:   sd.returnUrl,

		Remark: remark,
	}

	sanDe := sandutil.SandAES{}
	key := sanDe.RandStr(16)

	dataJson, err := json.Marshal(req)
	fmt.Println(string(dataJson))
	if err != nil {
		return nil, err
	}
	aes := sandutil.SandAES{}
	aes.Key = []byte(key)

	data := aes.Encypt5(dataJson)

	req.EncryptKey, _ = sandutil.RsaEncrypt(key, sd.pubKey)
	req.Sign, _ = sandutil.SignSand(sd.prvKey, data)
	req.Data = data

	reqJson, _ := json.Marshal(req)
	//fmt.Println("请求参数:", string(reqJson))
	resp, err := sandutil.Do("https://cap.sandpay.com.cn/v4/elecaccount/ceas.elec.account.member.status.modify", string(reqJson))

	if err != nil {
		return
	}

	d := make(map[string]interface{})
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, err
	}
	fmt.Println("杉德回调解析结果:" + string(resp))

	da := d["data"].(string)
	si := d["sign"].(string)

	err = sandutil.Verification(da, si, sd.pubKey)
	if err != nil {
		return nil, fmt.Errorf("验签失败, %s", err)
	}

	encryptKey := d["encryptKey"].(string)

	decryptAESKey, err := sandutil.RsaDecrypt(encryptKey, sd.prvKey)
	if err != nil {
		return
	}

	aes1 := &sandutil.SandAES{
		Key: []byte(decryptAESKey),
	}

	responseData, err = aes1.AesECBDecrypt(da)

	fmt.Println(string(responseData))
	return

}

type CloseAccountConfirmReq struct {
	Version            string `json:"version"`
	Mid                string `json:"mid"`
	CustomerOrderNo    string `json:"customerOrderNo"`
	OriCustomerOrderNo string `json:"oriCustomerOrderNo"`
	SignType           string `json:"signType"`
	EncryptType        string `json:"encryptType"`
	Timestamp          string `json:"timestamp,omitempty"`
	EncryptKey         string `json:"encryptKey,omitempty"`
	Data               string `json:"data,omitempty"`
	Sign               string `json:"sign,omitempty"`
	AccountType        string `json:"accountType,omitempty"`
	BizUserNo          string `json:"bizUserNo,omitempty"`
	SmsCode            string `json:"smsCode,omitempty"`
}

// 注销接口
func (sd *sandV4Client) CloseAccountConfirm(uid, smsCode, oriCustomerOrderNo, remark string) (responseData []byte, err error) {

	req := CloseAccountConfirmReq{
		Version:            "1.0",
		Mid:                sd.merNo,
		CustomerOrderNo:    orderno.Generate(),
		OriCustomerOrderNo: oriCustomerOrderNo,
		SignType:           "SHA1WithRSA",
		EncryptType:        "AES",
		Timestamp:          time.Now().Format("2006-01-02 15:04:05"),
		AccountType:        "01",
		BizUserNo:          uid,
		SmsCode:            smsCode,
		// OrderAmt:        fmt.Sprintf("%.2f", money),

		// Postscript: postscript,
	}

	sanDe := sandutil.SandAES{}
	key := sanDe.RandStr(16)

	dataJson, err := json.Marshal(req)
	fmt.Println(string(dataJson))
	if err != nil {
		return nil, err
	}
	aes := sandutil.SandAES{}
	aes.Key = []byte(key)

	data := aes.Encypt5(dataJson)

	req.EncryptKey, _ = sandutil.RsaEncrypt(key, sd.pubKey)
	req.Sign, _ = sandutil.SignSand(sd.prvKey, data)
	req.Data = data

	reqJson, _ := json.Marshal(req)
	fmt.Println("请求参数:", string(reqJson))
	resp, err := sandutil.Do("https://cap.sandpay.com.cn/v4/elecaccount/ceas.elec.account.member.modify.confirm", string(reqJson))

	if err != nil {
		return
	}

	d := make(map[string]interface{})
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, err
	}
	// fmt.Println("杉德回调解析结果:" + string(resp))

	da := d["data"].(string)
	si := d["sign"].(string)

	err = sandutil.Verification(da, si, sd.pubKey)
	if err != nil {
		return nil, fmt.Errorf("验签失败, %s", err)
	}

	encryptKey := d["encryptKey"].(string)

	decryptAESKey, err := sandutil.RsaDecrypt(encryptKey, sd.prvKey)
	if err != nil {
		return
	}

	aes1 := &sandutil.SandAES{
		Key: []byte(decryptAESKey),
	}

	responseData, err = aes1.AesECBDecrypt(da)

	fmt.Println(string(responseData))
	return

}
