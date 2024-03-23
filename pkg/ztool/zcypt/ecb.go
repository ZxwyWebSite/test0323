package zcypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	"github.com/ZxwyWebSite/ztool/x/bytesconv"
)

// 填充模式
const Pkcs5Padding = "PKCS5" // PKCS5填充模式
const Pkcs7Padding = "PKCS7" // PKCS7填充模式
const ZEROSPadding = "ZEROS" // ZEROS填充模式

// 报错信息
var ErrUnPadding = errors.New(`UnPadding error`)

// func main() {
//         plaintext := "真的爱你goalngcodes.com"

//         //　【特别注意】：这里密钥只能是16位或24位或32位，实际业务中，密钥不能写死在代码中，要从其他配置读取
//         key1 := "1234567890123456" // 16位

//         // 【加密】：密钥16位，即AES-128，PKCS7填充模式
//         ciphertext1, _ := AesECBEncrypt([]byte(plaintext), []byte(key1), Pkcs7Padding)
//         fmt.Printf("ct1:%s\n", ciphertext1)
//         // 注意：这里出来的加密后内容，不像CBC算法一样每次都变化，因为我们算法中没有引入动态iv
//         // ct1:20a3d76128e5c38381d1ca46727c3c3df9b5d4160e7eb5235ad6481e9783804d

//         // 【解密】：密钥16位，即AES-128，PKCS7填充模式
//         ciphertext1Bytes, _ := hex.DecodeString(ciphertext1)
//         plaintext1, _ := AesECBDecrypt(ciphertext1Bytes, []byte(key1), Pkcs7Padding)
//         fmt.Printf("pt1:%s\n", plaintext1)
//         // pt1:真的爱你goalngcodes.com

//         key2 := "123456789012345678901234" // 24位

//         // 【加密】：密钥24位，即AES-192，PKCS7填充模式
//         ciphertext2, _ := AesECBEncrypt([]byte(plaintext), []byte(key2), Pkcs7Padding)
//         fmt.Printf("ct2:%s\n", ciphertext2)
//         // 注意：这里出来的加密后内容，不像CBC算法一样每次都变化，因为我们算法中没有引入动态iv
//         // ct2:8c555a31fac6f90db4a03fdbdb3ebcd820c3b0fbb44c5fa43eb9c5dc4402e707

//         // 【解密】：密钥24位，即AES-192，PKCS7填充模式
//         ciphertext2Bytes, _ := hex.DecodeString(ciphertext2)
//         plaintext2, _ := AesECBDecrypt(ciphertext2Bytes, []byte(key2), Pkcs7Padding)
//         fmt.Printf("pt2:%s\n", plaintext2)
//         // pt1:真的爱你goalngcodes.com

//         key3 := "12345678901234567890123456789012" // 32位

//         // 【加密】：密钥32位，即AES-256，PKCS7填充模式
//         ciphertext3, _ := AesECBEncrypt([]byte(plaintext), []byte(key3), Pkcs7Padding)
//         fmt.Printf("ct3:%s\n", ciphertext3)
//         // 注意：这里出来的加密后内容，不像CBC算法一样每次都变化，因为我们算法中没有引入动态iv
//         // ct3:56f85fdc89c25efcf49c044aa4846ddf7ef18995114a193bfeaab30298d8c767

//         // 【解密】：密钥32位，即AES-256，PKCS7填充模式
//         ciphertext3Bytes, _ := hex.DecodeString(ciphertext3)
//         plaintext3, _ := AesECBDecrypt(ciphertext3Bytes, []byte(key3), Pkcs7Padding)
//         fmt.Printf("pt3:%s\n", plaintext3)
//         // pt3:真的爱你goalngcodes.com
// }

// AesECBEncrypt aes加密算法，ECB模式，可以指定三种填充模式
func AesECBEncrypt(plaintext, key []byte, padding string) (string, error) {
	// 要加密的内容长度必须是快长度的整数倍，首先进行填充
	plaintext = Padding(padding, plaintext, aes.BlockSize)
	if len(plaintext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(plaintext))
	mode := NewECBEncrypter(block)
	mode.CryptBlocks(ciphertext, plaintext)

	return fmt.Sprintf("%x", ciphertext), nil
}

// AesECBDecrypt aes解密算法，ECB模式，可以指定三种填充模式
func AesECBDecrypt(ciphertext, key []byte, padding string) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	// ECB mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := NewECBDecrypter(block)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	ciphertext, err = UnPadding(padding, ciphertext)
	if err != nil {
		return "", ErrUnPadding
	}

	return bytesconv.BytesToString(ciphertext), nil
}

// /////////////////////////////////////////////////////
type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// /////////////////////////
func Padding(padding string, src []byte, blockSize int) []byte {
	switch padding {
	case Pkcs5Padding:
		src = PKCS5Padding(src, blockSize)
	case Pkcs7Padding:
		src = PKCS7Padding(src, blockSize)
	case ZEROSPadding:
		src = ZerosPadding(src, blockSize)
	}
	return src
}

func UnPadding(padding string, src []byte) ([]byte, error) {
	switch padding {
	case Pkcs5Padding:
		return PKCS5Unpadding(src)
	case Pkcs7Padding:
		return PKCS7UnPadding(src)
	case ZEROSPadding:
		return ZerosUnPadding(src)
	}
	return src, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	return PKCS7Padding(src, blockSize)
}

func PKCS5Unpadding(src []byte) ([]byte, error) {
	return PKCS7UnPadding(src)
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return src, ErrUnPadding
	}
	unpadding := int(src[length-1])
	if length < unpadding {
		return src, ErrUnPadding
	}
	return src[:(length - unpadding)], nil
}

func ZerosPadding(src []byte, blockSize int) []byte {
	paddingCount := blockSize - len(src)%blockSize
	if paddingCount == 0 {
		return src
	} else {
		return append(src, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
	}
}

func ZerosUnPadding(src []byte) ([]byte, error) {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1], nil
		}
	}
}
