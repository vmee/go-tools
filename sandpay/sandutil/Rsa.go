package sandutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

func RsaEncrypt(value string, rsaPublicKey *rsa.PublicKey) (string, error) {
	// fmt.Println("RSA 加密数据", value)
	buffer, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(value))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buffer), nil
}

//RSA解密
func RsaDecrypt(value string, privateKey *rsa.PrivateKey) (string, error) {

	valueBytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	buffer, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, valueBytes)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}
