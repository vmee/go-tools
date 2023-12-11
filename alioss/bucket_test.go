package alioss_test

import (
	"testing"

	"github.com/vmee/go-tools/alioss"
)

func getNewBucket() (*alioss.AliYunOssBucket, error) {

	return alioss.NewAliYunOssBucket(
		alioss.WithAliOssEndpoint(""),
		alioss.WithAliOssAccessKeyId(""),
		alioss.WithAliOssAccessKeySecret(""),
		alioss.WithAliOssBucket(""),
		alioss.WithAliOssBucketBaseUrl(""),
	)

}

func TestUploadLocalFile(t *testing.T) {
	b, err := getNewBucket()
	if err != nil {
		t.Error(err)
	}

	url, err := b.UploadLocalFile("local1.txt", "./local.txt")
	if err != nil {
		t.Error(err)
	}

	t.Log(url)
}

func TestUploadWebFile(t *testing.T) {
	b, err := getNewBucket()
	if err != nil {
		t.Error(err)
	}

	url, err := b.UploadWebFile("local1.png", "https://wenxin.baidu.com/younger/file/ERNIE-ViLG/4a41f39f62022ac900ee11a3dcd25e3aex")
	if err != nil {
		t.Error(err)
	}

	t.Log(url)

}

func TestUploadString(t *testing.T) {
	b, err := getNewBucket()
	if err != nil {
		t.Error(err)
	}

	url, err := b.UploadString("up.text", "string1111")
	if err != nil {
		t.Error(err)
	}

	t.Log(url)
}

func TestUploadStruct(t *testing.T) {
	b, err := getNewBucket()
	if err != nil {
		t.Error(err)
	}

	type Metadata struct {
		Id   uint64
		Name string
	}

	d := Metadata{
		Id:   1,
		Name: "我名字",
	}

	url, err := b.UploadStruct("metadata.json", d)
	if err != nil {
		t.Error(err)
	}

	t.Log(url)
}
