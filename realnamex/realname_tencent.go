package realnamex

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	gourl "net/url"

	"github.com/go-pay/gopay/pkg/xhttp"
	jsoniter "github.com/json-iterator/go"
	"github.com/vmee/go-tools/xerr"
)

// https://market.cloud.tencent.com/products/40639
// 望为科技二要素

type tencentRealNameResponse struct {
	Status        string `json:"status"`
	State         int64  `json:"state"`
	ResultMessage string `json:"result_message,omitempty"`
	RequestId     string `json:"request_id"`
}

func calcAuthorization(secretId string, secretKey string) (auth string, datetime string, err error) {
	timeLocation, _ := time.LoadLocation("Etc/GMT")
	datetime = time.Now().In(timeLocation).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s", datetime)

	// hmac-sha1
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	auth = fmt.Sprintf("{\"id\":\"%s\", \"x-date\":\"%s\", \"signature\":\"%s\"}",
		secretId, datetime, sign)

	return auth, datetime, nil
}

func urlencode(params map[string]string) string {
	var p = gourl.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return p.Encode()
}

func genUniqueID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(100000))
}

// tencentRealName 用户实名认证
func TencentRealName(ctx context.Context, secretId string, secretKey string, realName, idCard string) (b bool, err error) {

	// 签名
	auth, _, _ := calcAuthorization(secretId, secretKey)

	// 请求方法
	reqID := genUniqueID()

	// 查询参数
	// queryParams := make(map[string]string)

	// body参数
	bodyParams := make(map[string]interface{})
	bodyParams["id_number"] = idCard
	bodyParams["name"] = realName
	// bodyParamStr := urlencode(bodyParams)
	// url参数拼接
	url := "https://ap-beijing.cloudmarket-apigw.com/service-rjw71k75/verify_id_name"

	// if len(queryParams) > 0 {
	// 	url = fmt.Sprintf("%s?%s", url, urlencode(queryParams))
	// }

	// var body io.Reader = nil
	// body := strings.NewReader(bodyParamStr)
	// headers["Content-Type"] = "application/x-www-form-urlencoded"

	httpClient := xhttp.NewClient()
	httpClient.Header.Add("Authorization", auth)
	httpClient.Header.Add("request-id", reqID)
	//  headers["Content-Type"] = "application/x-www-form-urlencoded"
	httpClient.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, bs, err := httpClient.Post(url).SendBodyMap(bodyParams).EndBytes(ctx)
	response := &tencentRealNameResponse{}
	// return false, fmt.Errorf("test error: %s", string(bs))
	jsonErr := jsoniter.UnmarshalFromString(string(bs), response)

	if err != nil || jsonErr != nil || res.StatusCode != http.StatusOK {
		return false, xerr.NewBizErr("实名认证请求失败")
	}

	return response.State == 1, nil
}
