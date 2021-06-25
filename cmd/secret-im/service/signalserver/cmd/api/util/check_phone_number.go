package util

import (
	"crypto/sha1"
	"encoding/base64"
	"regexp"

	"github.com/nyaruka/phonenumbers"
)

// 是否有效号码
func IsValidNumber(number string) bool {
	matched, err := regexp.Match("^\\+[0-9]+$", []byte(number))
	if err != nil {
		panic(err)
	}
	if !matched {
		return false
	}

	phoneNumber, err := phonenumbers.Parse(number, "US")
	if err != nil {
		return false
	}
	return phonenumbers.IsPossibleNumber(phoneNumber)
}

// 获取号码令牌
func GetContactToken(number string) string {
	hash := sha1.New()
	hash.Write([]byte(number))
	digest := hash.Sum(nil)
	if len(digest) > 10 {
		digest = digest[:10]
	}
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(digest)
}
