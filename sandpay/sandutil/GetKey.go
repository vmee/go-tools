package sandutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
)

func LoadPrivateKey(pemPath string) *rsa.PrivateKey {

	key, _ := ioutil.ReadFile(pemPath)
	block, _ := pem.Decode(key)
	if block == nil {
		return nil
	}
	p, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		pk := p.(*rsa.PrivateKey)
		// keyBytes, err := x509.MarshalPKCS8PrivateKey(pk)
		// if err != nil {
		// 	return pk
		// }
		// pemBlock := pem.Block{
		// 	Type:  "PRIVATE KEY",
		// 	Bytes: keyBytes,
		// }
		// keyBody := string(pem.EncodeToMemory(&pemBlock))

		// fmt.Println(keyBody)
		// fmt.Println("============================== PRIVATE KEY  私钥===================")
		return pk
	}
	return nil
}

//格式转化
func ChunkSplit(body string, chunklen uint, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, erunes := []rune(body), []rune(end)
	l := uint(len(runes))
	if l <= 1 || l < chunklen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(erunes))
	var i uint
	for i = 0; i < l; i += chunklen {
		if i+chunklen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunklen]...)
		}
		ns = append(ns, erunes...)
	}
	return string(ns)
}

func LoadPublicKey(pemPath string) *rsa.PublicKey {

	key, err := ioutil.ReadFile(pemPath)
	if err != nil {
		fmt.Sprintf("read public key file: %s", err)
		return nil
	}

	base64 := base64.StdEncoding.EncodeToString(key)
	// base64 := Base64Encode(data)
	cert := ChunkSplit(base64, 64, "\n")
	cert = "-----BEGIN CERTIFICATE-----\n" + cert + "-----END CERTIFICATE-----\n"

	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		return nil
	}

	certBody, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil
	}
	// publicKeyDer, err := x509.MarshalPKIXPublicKey(certBody.PublicKey)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	return nil
	// }
	// publickeyBlock := pem.Block{
	// 	Type:  "PUBLIC KEY",
	// 	Bytes: publicKeyDer,
	// }
	// publicKeyPem := string(pem.EncodeToMemory(&publickeyBlock))
	// fmt.Println(publicKeyPem)
	// fmt.Println("============================== PUBLIC KEY ===================")
	pb := certBody.PublicKey.(*rsa.PublicKey)
	return pb
}

func SignStr(values url.Values) string {
	if len(values) == 0 {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}

		buf.WriteString(k)
		buf.WriteByte('=')
		vv := values.Get(k)
		if vv == "" {
			continue
		}
		buf.WriteString(vv)
	}
	return buf.String()
}

// 签名
func SignSand(privateKey *rsa.PrivateKey, data string) (string, error) {

	// fmt.Println("============================== sign string ===================")
	// fmt.Println(data)
	// fmt.Println("============================== sign string ===================")

	h := crypto.SHA1.New()
	h.Write([]byte(data))

	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, h.Sum(nil))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), err
}

func SandVerification(data, signature []byte, publicKey *rsa.PublicKey) error {

	hash := crypto.SHA1
	if !hash.Available() {
		return fmt.Errorf("crypto: requested hash function (%s) is unavailable", hash.String())
	}

	h := hash.New()
	h.Write(data)

	return rsa.VerifyPKCS1v15(publicKey, hash, h.Sum(nil), signature)
}

//验签
func Verification(data, signStr string, PublickKeyP *rsa.PublicKey) error {

	sign, err := base64.StdEncoding.DecodeString(strings.Replace(signStr, " ", "+", -1))
	if err != nil {
		return err
	}
	return SandVerification([]byte(data), sign, PublickKeyP)
}
