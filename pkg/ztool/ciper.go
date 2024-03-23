// 数据加密 (beta); 分类: Cyp_(加密)

package ztool

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/ZxwyWebSite/ztool/x/bytesconv"
	"github.com/ZxwyWebSite/ztool/zcypt"
)

// AES解密 [base64, 密钥] [数据, 错误] (aes-128-ecb)
func Cyp_aesDecrypt(data string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return ``, fmt.Errorf(`aes.NewCipher: %s`, err)
	}
	decrypter := zcypt.NewECBDecrypter(block)
	msg, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ``, fmt.Errorf(`base64.DecodeString: %s`, err)
	}
	out := make([]byte, len(msg))
	decrypter.CryptBlocks(out, msg)
	return bytesconv.BytesToString(out), nil
}

// RSA加密 [数据, 公钥] [base64, 错误] (RSA_PKCS1_OAEP_PADDING)
func Cyp_rsaEncypt(data []byte, key string) (string, error) {
	pblock, _ := pem.Decode(bytesconv.StringToBytes(key))
	pubKey, err := x509.ParsePKIXPublicKey(pblock.Bytes)
	if err != nil {
		return ``, fmt.Errorf(`x509.ParsePKIXPublicKey: %s`, err)
	}
	encData, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pubKey.(*rsa.PublicKey), data, nil)
	if err != nil {
		return ``, fmt.Errorf(`rsa.EncryptOAEP: %s`, err)
	}
	return base64.StdEncoding.EncodeToString(encData), nil
}
