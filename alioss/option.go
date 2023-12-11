package alioss

import "strings"

type AliOssConf struct {
	Endpoint, AccessKeyId, AccessKeySecret, Bucket, BucketBaseUrl string
}

// AliOssOption 配置项
type AliOssOption func(c *AliOssConf)

// 地域节点
func WithAliOssEndpoint(endpoint string) AliOssOption {
	return func(c *AliOssConf) {
		c.Endpoint = endpoint
	}
}

func WithAliOssAccessKeyId(accessKeyId string) AliOssOption {
	return func(c *AliOssConf) {
		c.AccessKeyId = accessKeyId
	}
}

func WithAliOssAccessKeySecret(accessKeySecret string) AliOssOption {
	return func(c *AliOssConf) {
		c.AccessKeySecret = accessKeySecret
	}
}

func WithAliOssBucket(bucket string) AliOssOption {
	return func(c *AliOssConf) {
		c.Bucket = bucket
	}
}

// bucket访问域名 请使用http/https完整URL
func WithAliOssBucketBaseUrl(baseUrl string) AliOssOption {
	return func(c *AliOssConf) {
		c.BucketBaseUrl = strings.TrimLeft(baseUrl, "/")
	}
}
