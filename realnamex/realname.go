package realnamex

import (
	"context"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	faceId "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
)

type realnameClient struct {
	c *faceId.Client
}

func NewClient(secretId, secretKey, region string) (*realnameClient, error) {

	credential := common.NewCredential(secretId, secretKey)
	clientProfile := profile.NewClientProfile()

	client, err := faceId.NewClient(credential, region, clientProfile)
	if err != nil {
		return nil, err
	}

	return &realnameClient{
		c: client,
	}, nil
}

// 二要素认证 姓名身份证号
func (rc realnameClient) IdCardOCRVerification(ctx context.Context, realName, idCard string) (b bool, err error) {

	req := faceId.NewIdCardOCRVerificationRequest()
	req.IdCard = &idCard
	req.Name = &realName

	resp, err := rc.c.IdCardOCRVerificationWithContext(ctx, req)
	if err != nil {
		return false, err
	}

	// fmt.Print(resp.ToJsonString())

	if *resp.Response.Result == "0" {
		return true, nil
	}

	return false, fmt.Errorf(*resp.Response.Description)
}
