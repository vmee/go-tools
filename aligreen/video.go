package aligreen

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

type VideoScanner struct {
	c       *green.Client
	bizType string
}

func (vs *VideoScanner) BindClient(c *green.Client, bizType string) {
	vs.c = c
	vs.bizType = bizType
}

func (vs *VideoScanner) Sync(images []*ContentScanData) (results []*ContentScanResult, err error) {
	return
}

func (vs *VideoScanner) Async(videos []*ContentScanData) (results []*ContentScanResult, err error) {

	if len(videos) <= 0 {
		return
	}

	tasks := []map[string]interface{}{}
	for _, vidoe := range videos {
		tasks = append(tasks,
			map[string]interface{}{"dataId": vidoe.DataId, "url": vidoe.Url})
	}

	// task := map[string]interface{}{"dataId": "检测数据ID", "url": "待检测视频链接地址"}
	// scenes：检测场景，支持指定多个场景。
	// callback、seed用于回调通知，可选参数。
	content, _ := json.Marshal(
		map[string]interface{}{
			"tasks": tasks, "scenes": [...]string{"porn", "terrorism", "ad"},
			"bizType": vs.bizType,
			// "bizType": "业务场景", "callback": "回调地址", "seed": "随机字符串",
		},
	)

	request := green.CreateVideoAsyncScanRequest()
	request.SetContent(content)
	response, _err := vs.c.VideoAsyncScan(request)
	if _err != nil {
		fmt.Println(_err.Error())
		return
	}
	if response.GetHttpStatus() != 200 {
		fmt.Println("response not success. status:" + strconv.Itoa(response.GetHttpStatus()))
	}
	fmt.Println(response.GetHttpContentString())

	return parseAsyncResponse(response.GetHttpContentBytes())

}

func (vs *VideoScanner) AsyncResults(scans []*ContentScanResult) (results []*ContentScanResult, err error) {
	// 请替换成您的AccessKey ID、AccessKey Secret。

	taskIds := []string{}
	for _, t := range scans {
		taskIds = append(taskIds, t.TaskId)
	}

	content, _ := json.Marshal(taskIds)

	request := green.CreateVideoAsyncScanResultsRequest()
	request.SetContent(content)
	response, err := vs.c.VideoAsyncScanResults(request)
	if err != nil {
		return
	}
	if response.GetHttpStatus() != 200 {
		fmt.Println("response not success. status:" + strconv.Itoa(response.GetHttpStatus()))
	}
	fmt.Println(response.GetHttpContentString())

	return parseResultsResponse(response.GetHttpContentBytes())

}
