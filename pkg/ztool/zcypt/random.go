package zcypt

import (
	crand "crypto/rand"
	mrand "math/rand"
)

// 生成随机字节
func RandomBytes(size int) []byte {
	buf := make([]byte, size)
	_, err := crand.Read(buf)
	if err != nil {
		for i := 0; i < size; i++ {
			buf[i] = byte(mrand.Intn(256))
		}
	}
	return buf
}
