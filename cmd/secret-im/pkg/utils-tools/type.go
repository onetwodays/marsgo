package utils

import (
	"encoding/base64"
	"errors"
)

// Base64字节数组
type Base64Bytes []byte

func (b Base64Bytes) MarshalJSON() ([]byte, error) {
	s := base64.StdEncoding.EncodeToString(b)
	return []byte(`"` + s + `"`), nil
}

func (b *Base64Bytes) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid base64 string")
	}

	var err error
	s := string(data[1 : len(data)-1])
	*b, err = base64.StdEncoding.DecodeString(s)
	return err
}

func (b Base64Bytes) MarshalYAML() (interface{}, error) {
	return base64.StdEncoding.EncodeToString(b), nil
}

func (b *Base64Bytes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	var err error
	if err = unmarshal(&s); err != nil {
		return err
	}
	*b, err = base64.StdEncoding.DecodeString(s)
	return err
}
