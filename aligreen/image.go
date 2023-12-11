package aligreen

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

type ImageScanner struct {
	c       *green.Client
	bizType string
}

func (ims *ImageScanner) BindClient(c *green.Client, bizType string) {
	ims.c = c
	ims.bizType = bizType
}

func (ims *ImageScanner) Sync(images []*ContentScanData) (results []*ContentScanResult, err error) {
	if len(images) <= 0 {
		return
	}

	tasks := []map[string]interface{}{}
	for _, image := range images {
		tasks = append(tasks,
			map[string]interface{}{"dataId": image.DataId, "url": image.Url})
	}

	// scenes：检测场景，支持指定多个场景。
	content, _ := json.Marshal(
		map[string]interface{}{
			"tasks":   tasks,
			"scenes":  [...]string{"porn", "terrorism", "ad", "qrcode"},
			"bizType": ims.bizType,
		},
	)

	request := green.CreateImageSyncScanRequest()
	request.SetContent(content)
	response, err := ims.c.ImageSyncScan(request)
	if err != nil {
		return
	}
	if response.GetHttpStatus() != 200 {
		fmt.Println("response not success. status:" + strconv.Itoa(response.GetHttpStatus()))
	}
	fmt.Println(response.GetHttpContentString())

	return parseResultsResponse(response.GetHttpContentBytes())
}

func (ims *ImageScanner) Async(images []*ContentScanData) (results []*ContentScanResult, err error) {

	if len(images) <= 0 {
		return
	}

	tasks := []map[string]interface{}{}
	for _, image := range images {
		tasks = append(tasks,
			map[string]interface{}{"dataId": image.DataId, "url": image.Url})
	}

	// scenes：检测场景，支持指定多个场景。
	content, _ := json.Marshal(
		map[string]interface{}{
			"tasks":   tasks,
			"scenes":  [...]string{"porn", "terrorism", "ad", "qrcode"},
			"bizType": ims.bizType,
		},
	)

	request := green.CreateImageAsyncScanRequest()
	request.SetContent(content)
	response, err := ims.c.ImageAsyncScan(request)
	if err != nil {
		return
	}
	if response.GetHttpStatus() != 200 {
		return nil, fmt.Errorf("response not success. status:%d", response.GetHttpStatus())
	}
	fmt.Println(response.GetHttpContentString())

	return parseAsyncResponse(response.GetHttpContentBytes())
}

func (ims *ImageScanner) AsyncResults(scans []*ContentScanResult) (results []*ContentScanResult, err error) {
	if len(scans) <= 0 {
		return
	}

	taskIds := []string{}
	for _, t := range scans {
		taskIds = append(taskIds, t.TaskId)
	}

	content, _ := json.Marshal(taskIds)

	request := green.CreateImageAsyncScanResultsRequest()
	request.SetContent(content)
	response, err := ims.c.ImageAsyncScanResults(request)
	if err != nil {
		return
	}
	if response.GetHttpStatus() != 200 {
		fmt.Println("response not success. status:" + strconv.Itoa(response.GetHttpStatus()))
	}
	fmt.Println(response.GetHttpContentString())

	return parseResultsResponse(response.GetHttpContentBytes())

}
