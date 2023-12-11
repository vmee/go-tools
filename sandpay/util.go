package sandpay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const OK = "000000"

type X map[string]string

type Data struct {
	Head X `json:"head"`
	Body X `json:"body"`
}

// PemBlockType pem block type which taken from the preamble.
type PemBlockType string

const (
	// RSAPKCS1 private key in PKCS#1
	RSAPKCS1 PemBlockType = "RSA PRIVATE KEY"
	// RSAPKCS8 private key in PKCS#8
	RSAPKCS8 PemBlockType = "PRIVATE KEY"
)

// PrivateKey RSA private key
type PrivateKey struct {
	key *rsa.PrivateKey
}

// Decrypt rsa decrypt with PKCS #1 v1.5
func (pk *PrivateKey) Decrypt(cipherText []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, pk.key, cipherText)
}

// Sign returns sha-with-rsa signature.
func (pk *PrivateKey) Sign(hash crypto.Hash, data []byte) ([]byte, error) {
	if !hash.Available() {
		return nil, fmt.Errorf("crypto: requested hash function (%s) is unavailable", hash.String())
	}

	h := hash.New()
	h.Write(data)

	signature, err := rsa.SignPKCS1v15(rand.Reader, pk.key, hash, h.Sum(nil))

	if err != nil {
		return nil, err
	}

	return signature, nil
}

// NewPrivateKeyFromPemFile returns new private key with pem file.
func NewPrivateKeyFromPemFile(pemFile string) (*PrivateKey, error) {
	keyPath, err := filepath.Abs(pemFile)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(keyPath)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)

	if block == nil {
		return nil, errors.New("no PEM data is found")
	}

	var pk interface{}

	switch PemBlockType(block.Type) {
	case RSAPKCS1:
		pk, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case RSAPKCS8:
		pk, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	if err != nil {
		return nil, err
	}

	return &PrivateKey{key: pk.(*rsa.PrivateKey)}, nil
}

// PublicKey RSA public key
type PublicKey struct {
	key *rsa.PublicKey
}

// Encrypt rsa encrypt with PKCS #1 v1.5
func (pk *PublicKey) Encrypt(plainText []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pk.key, plainText)
}

// Verify verifies the sha-with-rsa signature.
func (pk *PublicKey) Verify(hash crypto.Hash, data, signature []byte) error {
	if !hash.Available() {
		return fmt.Errorf("crypto: requested hash function (%s) is unavailable", hash.String())
	}

	h := hash.New()
	h.Write(data)

	return rsa.VerifyPKCS1v15(pk.key, hash, h.Sum(nil), signature)
}

// NewPublicKeyFromDerFile returns public key with DER file.
// NOTE: PEM format with -----BEGIN CERTIFICATE----- | -----END CERTIFICATE-----
// openssl x509 -inform der -in cert.cer -out cert.pem
func NewPublicKeyFromDerFile(pemFile string) (*PublicKey, error) {
	keyPath, err := filepath.Abs(pemFile)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(keyPath)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)

	if block == nil {
		return nil, errors.New("no PEM data is found")
	}

	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		return nil, err
	}

	return &PublicKey{key: cert.PublicKey.(*rsa.PublicKey)}, nil
}

// MarshalNoEscapeHTML marshal with no escape HTML
func MarshalNoEscapeHTML(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)

	if err := jsonEncoder.Encode(v); err != nil {
		return nil, err
	}

	b := buf.Bytes()

	// 去掉 go std 给末尾加的 '\n'
	// @see https://github.com/golang/go/issues/7767
	if l := len(b); l != 0 && b[l-1] == '\n' {
		b = b[:l-1]
	}

	return b, nil
}
