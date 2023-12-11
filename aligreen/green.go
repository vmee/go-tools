package aligreen

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

type Scanner interface {
	//bind client bizType业务场景,默认default
	BindClient(c *green.Client, bizType string)

	Sync(images []*ContentScanData) (results []*ContentScanResult, err error)

	Async(videos []*ContentScanData) (results []*ContentScanResult, err error)

	AsyncResults(scans []*ContentScanResult) (results []*ContentScanResult, err error)
}
type AuditStatus string

const (
	AuditStatusDoing AuditStatus = "doing"
	AuditStatusPass  AuditStatus = "pass"
	AuditStatusFail  AuditStatus = "fail"
)

type scanResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data []*scanData `json:"data"`
}

type scanData struct {
	Code    int               `json:"code"`
	Msg     string            `json:"msg"`
	Url     string            `json:"url"`
	DataId  string            `json:"dataId"`
	TaskId  string            `json:"taskId"`
	Results []*scanDataResult `json:"results"`
}

type scanDataResult struct {
	Suggestion string `json:"suggestion"`
	Scene      string `json:"scene"`
}

type ContentScanData struct {
	DataId  string
	Url     string
	Content string
	TaskId  string // 检查结果需要带个值
}

type ContentScanResult struct {
	DataId      string
	TaskId      string
	AuditStatus AuditStatus // doing审核中 pass审核通过 fail审核未通过
}

type contentFormat []string

var (

	// 图片格式 PNG","JPG","JPEG、BMP、GIF、WEBP
	// 视频格式 AVI、FLV、MP4、MPG、ASF、WMV、MOV、WMA、RMVB、RM、FLASH、TS
	// 音频格式 MP3、WAV、AAC、WMA、OGG、M4A、AMR、AUDIO、M3U8
	imageFormats contentFormat = []string{"PNG", "JPG", "JPEG", "BMP", "GIF", "WEBP"}
	videoFormats contentFormat = []string{"AVI", "FLV", "MP4", "MPG", "ASF", "WMV", "MOV", "WMA", "RMVB", "RM", "FLASH", "TS"}
	audioFormats contentFormat = []string{"MP3", "WAV", "AAC", "WMA", "OGG", "M4A", "AMR", "AUDIO", "M3U8"}
)

func (cf contentFormat) valid(format string) bool {

	format = strings.ToUpper(format)
	for _, v := range cf {
		if v == format {
			return true
		}
	}

	return false
}

type GreenClient struct {
	c       *green.Client
	bizType string
}

func NewGreenClient(regionId, accessKeyId, accessKeySecret string) *GreenClient {
	client, err := green.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}

	return &GreenClient{
		c:       client,
		bizType: "default",
	}
}

func NewBizTypeGreenClient(regionId, accessKeyId, accessKeySecret, bizType string) *GreenClient {
	client, err := green.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}

	return &GreenClient{
		c:       client,
		bizType: bizType,
	}
}

type ContentScanDataType []*ContentScanData

func (cas ContentScanDataType) filter(cf contentFormat) (res []*ContentScanData) {
	for _, v := range cas {
		if v.Url == "" {
			continue
		}

		ext := path.Ext(v.Url)
		if ext == "" {
			continue
		}
		if cf.valid(ext[1:]) {
			res = append(res, v)
		}
	}

	return
}

func (g *GreenClient) GetBizType() string {
	if g.bizType != "" {
		return g.bizType
	}

	return "default"
}

// 检查是否可以审核
func (g *GreenClient) CanScan(url string) bool {
	if url == "" {
		return false
	}

	ext := path.Ext(url)
	if ext == "" {
		return false
	}

	extStr := ext[1:]

	return imageFormats.valid(extStr) || videoFormats.valid(extStr) || audioFormats.valid(extStr)
}

// 同步scan 仅支持图片文本 有一个失败就返回失败
func (g *GreenClient) SyncScanIsPass(data []*ContentScanData) (isPass bool, err error) {
	results, err := g.SyncScan(data)
	if err != nil {
		return
	}

	isPass = true
	for _, v := range results {
		if v.AuditStatus != AuditStatusPass {
			isPass = false
			return
		}
	}

	return
}

func (g *GreenClient) getScannerContents(data []*ContentScanData) map[Scanner][]*ContentScanData {

	contents := ContentScanDataType(data)

	commonData := map[Scanner][]*ContentScanData{}

	commonData[&ImageScanner{}] = contents.filter(imageFormats)
	commonData[&AudioScanner{}] = contents.filter(audioFormats)
	commonData[&VideoScanner{}] = contents.filter(videoFormats)
	texts := []*ContentScanData{}
	for _, v := range contents {
		if v.Content != "" {
			texts = append(texts, v)
		}
	}

	commonData[&TextScanner{}] = texts

	return commonData
}

// 同步scan 仅支持图片文本
func (g *GreenClient) SyncScan(data []*ContentScanData) (results []*ContentScanResult, err error) {
	commonData := g.getScannerContents(data)

	for k, v := range commonData {
		if len(v) > 0 {
			r, err := g.Sync(k, v)
			if err != nil {
				return nil, err
			}

			results = append(results, r...)
		}
	}

	return
}

// 有一个失败就返回失败
func (g *GreenClient) CommonScanIsPass(data []*ContentScanData) (isPass bool, err error) {
	results, err := g.CommonScan(data)
	if err != nil {
		return
	}

	isPass = true
	for _, v := range results {
		if v.AuditStatus != AuditStatusPass {
			isPass = false
			return
		}
	}

	return
}

// 聚合dataId 相同的dataId有一个失败则认为失败
func (g *GreenClient) CommonScanResultsAggregation(results []*ContentScanResult) (newResults map[string]*ContentScanResult, err error) {
	newResults = map[string]*ContentScanResult{}
	for _, v := range results {
		result, ok := newResults[v.DataId]
		if !ok {
			newResults[v.DataId] = v
			continue
		}

		if result.AuditStatus == AuditStatusFail {
			continue
		}

		if v.AuditStatus == AuditStatusFail {
			newResults[v.DataId] = v
		}
	}

	return
}

func (g *GreenClient) CommonScan(data []*ContentScanData) (results []*ContentScanResult, err error) {
	commonData := g.getScannerContents(data)

	for k, v := range commonData {
		if len(v) > 0 {
			r, err := g.Async(k, v)
			if err != nil {
				return nil, err
			}

			results = append(results, r...)
		}
	}

	return
}

func (g *GreenClient) CommonScanResults(data []*ContentScanData) (results []*ContentScanResult, err error) {
	commonData := g.getScannerContents(data)

	for k, v := range commonData {
		if len(v) > 0 {

			vs := []*ContentScanResult{}
			for _, vv := range v {
				if vv.TaskId != "" {
					vs = append(vs, &ContentScanResult{
						TaskId: vv.TaskId,
					})
				}
			}

			r, err := g.AsyncResults(k, vs)
			if err != nil {
				return nil, err
			}

			results = append(results, r...)
		}
	}

	return
}

func (g *GreenClient) Sync(s Scanner, data []*ContentScanData) (results []*ContentScanResult, err error) {
	s.BindClient(g.c, g.bizType)
	return s.Sync(data)
}

func (g *GreenClient) Async(s Scanner, data []*ContentScanData) (results []*ContentScanResult, err error) {
	s.BindClient(g.c, g.bizType)
	return s.Async(data)
}

func (g *GreenClient) AsyncResults(s Scanner, data []*ContentScanResult) (results []*ContentScanResult, err error) {
	s.BindClient(g.c, g.bizType)
	return s.AsyncResults(data)
}

func parseResultsResponse(responseContent []byte) (results []*ContentScanResult, err error) {
	resp := scanResponse{}
	err = json.Unmarshal(responseContent, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("%s", resp.Msg)
	}

	for _, data := range resp.Data {

		b := AuditStatusPass
		for _, r := range data.Results {
			if r.Suggestion != "pass" {
				b = AuditStatusFail
			}
		}

		results = append(results, &ContentScanResult{
			DataId:      data.DataId,
			AuditStatus: b,
			TaskId:      data.TaskId,
		})
	}

	return
}

func parseAsyncResponse(responseContent []byte) (results []*ContentScanResult, err error) {
	resp := scanResponse{}
	err = json.Unmarshal(responseContent, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("%s", resp.Msg)
	}

	for _, data := range resp.Data {
		if data.Code != 200 {
			results = append(results, &ContentScanResult{
				DataId:      data.DataId,
				AuditStatus: AuditStatusFail,
				TaskId:      data.TaskId,
			})
		}

		results = append(results, &ContentScanResult{
			DataId:      data.DataId,
			AuditStatus: AuditStatusDoing,
			TaskId:      data.TaskId,
		})
	}

	return

}
