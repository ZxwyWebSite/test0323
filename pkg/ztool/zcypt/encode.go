// 常用编码

package zcypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"github.com/ZxwyWebSite/ztool/x/bytesconv"
)

// base64

func Base64Encode(enc *base64.Encoding, src []byte) []byte {
	dst := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(dst, src)
	return dst
}

// base64.EncodeToString 进行Base64编码并转为字符串
func Base64ToString(enc *base64.Encoding, src []byte) string {
	return bytesconv.BytesToString(Base64Encode(enc, src))
}

func Base64Decode(enc *base64.Encoding, src []byte) ([]byte, error) {
	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	return dst[:n], err
}

// base64.DecodeToString 进行Base64解码并转为字符串
func Base64DTString(enc *base64.Encoding, src []byte) (string, error) {
	dst, err := Base64Decode(enc, src)
	if err != nil {
		return ``, err
	}
	return bytesconv.BytesToString(dst), nil
}

// hex

func HexEncode(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

func HexDecode(src []byte) ([]byte, error) {
	n, err := hex.Decode(src, src)
	return src[:n], err
}

// hex.EncodeToString 进行Hex编码并转为字符串
func HexToString(src []byte) string {
	return bytesconv.BytesToString(HexEncode(src))
}

// md5

func MD5Encode(src []byte) []byte {
	hash := md5.New()
	hash.Write(src)
	return hash.Sum(nil)
}

// 进行MD5编码并使用Hex转为字符串
func CreateMD5(s []byte) string {
	// hash := md5.New()
	// hash.Write(s)
	// return HexToString(hash.Sum(nil))
	return HexToString(MD5Encode(s))
}

func MD5EncStr(src string) string {
	return CreateMD5(bytesconv.StringToBytes(src))
}
