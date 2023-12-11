package sandpay

import (
	"context"
	"crypto"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client 杉德支付客户端
type Client interface {
	// Do 请求杉德API
	Do(ctx context.Context, reqURL string, form url.Values) (*Data, error)

	// Form 生成统一的POST表单（用于API请求或前端表单提交）
	Form(method, productID string, body X, options ...HeadOption) (url.Values, error)

	// Verify 验证并解析杉德API结果或回调通知
	Verify(form url.Values) (*Data, error)
}

type client struct {
	mid    string
	prvKey *PrivateKey
	pubKey *PublicKey
	cli    *http.Client
}

func (c *client) Do(ctx context.Context, reqURL string, form url.Values) (*Data, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.cli.Do(req)

	if err != nil {
		// If the context has been canceled, the context's error is probably more useful.
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}

		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	query, err := url.QueryUnescape(string(b))

	if err != nil {
		return nil, err
	}

	v, err := url.ParseQuery(query)

	if err != nil {
		return nil, err
	}

	return c.Verify(v)
}

func (c *client) Form(method, productID string, body X, options ...HeadOption) (url.Values, error) {
	data := &Data{
		Head: c.head(method, productID, options...),
		Body: body,
	}

	b, err := MarshalNoEscapeHTML(data)

	if err != nil {
		return nil, err
	}

	sign, err := c.prvKey.Sign(crypto.SHA1, b)

	if err != nil {
		return nil, err
	}

	form := url.Values{}

	form.Set("charset", "utf-8")
	form.Set("data", string(b))
	form.Set("signType", "01")
	form.Set("sign", base64.StdEncoding.EncodeToString(sign))

	return form, nil
}

func (c *client) Verify(form url.Values) (*Data, error) {
	sign, err := base64.StdEncoding.DecodeString(strings.Replace(form.Get("sign"), " ", "+", -1))

	if err != nil {
		return nil, err
	}

	if err = c.pubKey.Verify(crypto.SHA1, []byte(form.Get("data")), sign); err != nil {
		return nil, err
	}

	data := new(Data)

	if err := json.Unmarshal([]byte(form.Get("data")), data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *client) head(method, productID string, options ...HeadOption) X {
	head := X{
		"version":     "1.0",
		"method":      method,
		"productId":   productID,
		"accessType":  "1",
		"mid":         c.mid,
		"channelType": "07",
		"reqTime":     time.Now().Format("20060102150405"),
	}

	for _, f := range options {
		f(head)
	}

	return head
}

// Config 客户端配置
type Config struct {
	MID      string // 商户ID
	KeyFile  string // 商户私钥（PEM格式）
	CertFile string // 杉德公钥（PEM格式）
}

func NewClient(cfg *Config, options ...ClientOption) (Client, error) {
	prvKey, err := NewPrivateKeyFromPemFile(cfg.KeyFile)

	if err != nil {
		return nil, err
	}

	pubKey, err := NewPublicKeyFromDerFile(cfg.CertFile)

	if err != nil {
		return nil, err
	}

	c := &client{
		mid:    cfg.MID,
		prvKey: prvKey,
		pubKey: pubKey,
		cli: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				MaxIdleConns:          0,
				MaxIdleConnsPerHost:   1000,
				MaxConnsPerHost:       1000,
				IdleConnTimeout:       60 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}

	for _, f := range options {
		f(c)
	}

	return c, nil
}
