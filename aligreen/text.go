package aligreen

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

type TextScanner struct {
	c       *green.Client
	bizType string
}

func (ts *TextScanner) BindClient(c *green.Client, bizType string) {
	ts.c = c
	ts.bizType = bizType
}

func (ts *TextScanner) Sync(texts []*ContentScanData) (results []*ContentScanResult, err error) {

	if len(texts) <= 0 {
		return
	}

	tasks := []map[string]interface{}{}
	for _, text := range texts {
		tasks = append(tasks,
			map[string]interface{}{"dataId": text.DataId, "content": text.Content})
	}

	// task := map[string]interface{}{"content": "待检测文本内容"}
	// scenes：检测场景，唯一取值：antispam。
	content, _ := json.Marshal(
		map[string]interface{}{
			"scenes":  [...]string{"antispam"},
			"tasks":   tasks,
			"bizType": ts.bizType,
		},
	)

	textScanRequest := green.CreateTextScanRequest()
	textScanRequest.SetContent(content)
	response, err := ts.c.TextScan(textScanRequest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if response.GetHttpStatus() != 200 {
		fmt.Println("response not success. status:" + strconv.Itoa(response.GetHttpStatus()))
	}
	fmt.Println(response.GetHttpContentString())

	return parseResultsResponse(response.GetHttpContentBytes())
}

func (vs *TextScanner) Async(videos []*ContentScanData) (results []*ContentScanResult, err error) {
	return vs.Sync(videos)
}

func (vs *TextScanner) AsyncResults(scans []*ContentScanResult) (results []*ContentScanResult, err error) {
	return
}
