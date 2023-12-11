package realnamex

import (
	"context"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xhttp"
	jsoniter "github.com/json-iterator/go"
	"github.com/vmee/go-tools/xerr"
	"net/http"
)

type realNameResponse struct {
	Status        string `json:"status"`
	State         int64  `json:"state"`
	ResultMessage string `json:"result_message,omitempty"`
	RequestId     string `json:"request_id"`
}

// AliRealName 用户实名认证
func AliRealName(ctx context.Context, appCode, realName, idCard string) (b bool, err error) {

	postUrl := "https://verifyid.market.alicloudapi.com/id_name"
	bm := make(gopay.BodyMap)
	bm.Set("id", idCard).Set("name", realName)

	httpClient := xhttp.NewClient()
	httpClient.Header.Add("Authorization", "APPCODE "+appCode)
	res, bs, err := httpClient.Type(xhttp.TypeFormData).
		Post(postUrl).SendBodyMap(bm).EndBytes(ctx)
	response := &realNameResponse{}
	jsonErr := jsoniter.UnmarshalFromString(string(bs), response)

	if err != nil || jsonErr != nil || res.StatusCode != http.StatusOK {
		return false, xerr.NewBizErr("实名认证请求失败")
	}

	return response.State == 1, nil
}
