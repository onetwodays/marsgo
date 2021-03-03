package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

func DecodeToken(encoded string) (b []byte, err error) {
	encoded = strings.Replace(encoded, "-", "+", -1)
	encoded = strings.Replace(encoded, "_", "/", -1)
	return base64.RawStdEncoding.DecodeString(encoded)
}

// 更新通讯录
func ContactToken(number string) (hs []byte) {
	h := sha1.New()
	h.Write([]byte(number))
	hs = h.Sum(nil)
	var b bytes.Buffer
	b.Write(hs)
	b.Truncate(10)
	hs = b.Bytes()
	return
}