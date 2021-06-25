package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"secret-im/pkg/utils-tools"
	"strconv"
	"strings"
)

// 外部服务证书生成器
type ExternalServiceCredentialGenerator struct {
	key                []byte
	userIDKey          []byte
	usernameDerivation bool  //用户名字推导
}

// 创建外部服务证书生成器
func NewExternalServiceCredentialGenerator(
	key, userIDKey []byte, usernameDerivation bool) *ExternalServiceCredentialGenerator {
	return &ExternalServiceCredentialGenerator{
		key:                key,
		userIDKey:          userIDKey,
		usernameDerivation: usernameDerivation,
	}
}

// 生成证书
func (generator *ExternalServiceCredentialGenerator) GenerateFor(number string) *ExternalServiceCredentials {
	username := generator.getUserID(number, generator.usernameDerivation)
	currentTimeSeconds := utils.CurrentTimeMillis() / 1000
	prefix := username + ":" + strconv.FormatInt(currentTimeSeconds, 10)

	mac := hmac.New(sha256.New, generator.key)
	mac.Write([]byte(prefix))
	sum := mac.Sum(nil)
	if len(sum) > 10 {
		sum = sum[:10]
	}
	output := hex.EncodeToString(sum)

	token := prefix + ":" + output
	return &ExternalServiceCredentials{Username: username, Password: token}
}

// 是否有效
func (generator *ExternalServiceCredentialGenerator) IsValid(token, number string, currentTimeMillis int64) bool {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return false
	}

	if generator.getUserID(number, generator.usernameDerivation) != parts[0] {
		return false
	}

	if !generator.isValidTime(parts[1], currentTimeMillis) {
		return false
	}
	return generator.isValidSignature(parts[0]+":"+parts[1], parts[2])
}

// 获取用户ID
func (generator *ExternalServiceCredentialGenerator) getUserID(number string, usernameDerivation bool) string {
	if !usernameDerivation {
		return number
	}
	mac := hmac.New(sha256.New, generator.userIDKey)
	_, err := mac.Write([]byte(number))
	if err != nil {
		return ""
	}
	sum := mac.Sum(nil)
	if len(sum) > 10 {
		sum = sum[:10]
	}
	return hex.EncodeToString(sum)
}

// 时间是否有效
func (generator *ExternalServiceCredentialGenerator) isValidTime(timeString string, currentTimeMillis int64) bool {
	tokenTime, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return false
	}
	ourTime := currentTimeMillis / 1000
	return int64(math.Abs(float64(ourTime-tokenTime)))/60/60 < 24
}

// 签名是否有效
func (generator *ExternalServiceCredentialGenerator) isValidSignature(prefix, suffix string) bool {
	mac := hmac.New(sha256.New, generator.key)
	_, err := mac.Write([]byte(prefix))
	if err != nil {
		return false
	}
	ourSuffix := mac.Sum(nil)
	if len(ourSuffix) > 10 {
		ourSuffix = ourSuffix[:10]
	}

	theirSuffix, err := hex.DecodeString(suffix)
	if err != nil {
		return false
	}
	return bytes.Equal(ourSuffix, theirSuffix)
}
