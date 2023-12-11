package alioss

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliYunOssBucket struct {
	c      *AliOssConf
	Bucket *oss.Bucket
}

func NewAliYunOssBucket(options ...AliOssOption) (ab *AliYunOssBucket, err error) {
	c := &AliOssConf{}
	for _, f := range options {
		f(c)
	}

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(c.Endpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	// 填写存储空间名称，例如examplebucket。
	b, err := client.Bucket(c.Bucket)
	if err != nil {
		return nil, err
	}

	return &AliYunOssBucket{
		Bucket: b,
		c:      c,
	}, nil

}

// 上传结构体转存json saveFile只可使用相对路径
func (b *AliYunOssBucket) UploadStruct(saveFile string, s interface{}) (url string, err error) {

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return
	}

	return b.UploadBytes(saveFile, jsonBytes)
}

// 上传Bytes数组 saveFile只可使用相对路径
func (b *AliYunOssBucket) UploadBytes(saveFile string, bs []byte) (url string, err error) {

	err = b.Bucket.PutObject(saveFile, bytes.NewReader(bs))
	if err != nil {
		return
	}

	url = b.getFileUrl(saveFile)
	return
}

// 上传字符串 saveFile只可使用相对路径
func (b *AliYunOssBucket) UploadString(saveFile, str string) (url string, err error) {

	err = b.Bucket.PutObject(saveFile, strings.NewReader(str))
	if err != nil {
		return
	}

	url = b.getFileUrl(saveFile)
	return
}

// 上传本地文件 saveFile只可使用相对路径
func (b *AliYunOssBucket) UploadLocalFile(saveFile, filePath string) (url string, err error) {

	err = b.Bucket.PutObjectFromFile(saveFile, filePath)
	if err != nil {
		return
	}

	url = b.getFileUrl(saveFile)
	return
}

// 上传web文件 saveFile只可使用相对路径
func (b *AliYunOssBucket) UploadWebFile(saveFile, originUrl string) (url string, err error) {

	res, err := http.Get(originUrl)
	if err != nil {
		return
	}

	err = b.Bucket.PutObject(saveFile, io.Reader(res.Body))
	if err != nil {
		return
	}

	url = b.getFileUrl(saveFile)
	return
}

func (b *AliYunOssBucket) getFileUrl(saveFile string) (url string) {
	return b.c.BucketBaseUrl + "/" + saveFile
}
