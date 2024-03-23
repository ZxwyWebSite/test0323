package zcypt

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
)

// AES解密 (aes-128-ecb)
func AesDecrypt(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err //fmt.Errorf(`aes.NewCipher: %s`, err)
	}
	decrypter := NewECBDecrypter(block)
	out := make([]byte, len(src))
	decrypter.CryptBlocks(out, src)
	return out, nil
}

// RSA加密 (RSA_PKCS1_OAEP_PADDING) PKCS#8 (-----BEGIN PUBLIC KEY-----)
func RsaEncypt(data []byte, key []byte) ([]byte, error) {
	pblock, _ := pem.Decode(key)
	pubKey, err := x509.ParsePKIXPublicKey(pblock.Bytes)
	if err != nil {
		return nil, err //fmt.Errorf(`x509.ParsePKIXPublicKey: %s`, err)
	}
	encData, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pubKey.(*rsa.PublicKey), data, nil)
	if err != nil {
		return nil, err //fmt.Errorf(`rsa.EncryptOAEP: %s`, err)
	}
	return encData, nil
}

// RSA加密 PKCS#1 (-----BEGIN RSA PUBLIC KEY-----)

// RSA加密 (RSA_PKCS1_OAEP_PADDING) PKCS#8 (-----BEGIN PUBLIC KEY-----)
/*
 填充模式：
  none: Api标准
  sha1: CryptoJS标准
  sha256: Golang标准
*/
// func RsaEncyptOAEP(data []byte, key []byte, padding string) ([]byte, error) {
// 	pblock, _ := pem.Decode(key)
// 	pubKey, err := x509.ParsePKIXPublicKey(pblock.Bytes)
// 	if err != nil {
// 		return nil, err //fmt.Errorf(`x509.ParsePKIXPublicKey: %s`, err)
// 	}
// 	var hasher hash.Hash
// 	switch padding {
// 	case `none`:
// 		hasher = nil
// 	case `sha1`:
// 		hasher = sha1.New()
// 	case `sha256`:
// 		hasher = sha256.New()
// 	}
// 	encData, err := rsa.EncryptOAEP(hasher, rand.Reader, pubKey.(*rsa.PublicKey), data, nil)
// 	if err != nil {
// 		return nil, err //fmt.Errorf(`rsa.EncryptOAEP: %s`, err)
// 	}
// 	return encData, nil
// }
