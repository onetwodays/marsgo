package auth

import (
	"encoding/base64"
	"errors"
)



// 匿名访问密钥
type Anonymous struct {
	UnidentifiedSenderAccessKey []byte
}

// 创建匿名访问密钥
func NewAnonymous(header string) (*Anonymous, error) {
	if len(header) == 0 {
		return nil, errors.New("invalid header")
	}
	unidentifiedSenderAccessKey, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return nil, err
	}
	return &Anonymous{UnidentifiedSenderAccessKey: unidentifiedSenderAccessKey}, nil
}
