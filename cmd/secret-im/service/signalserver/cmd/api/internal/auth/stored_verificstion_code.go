package auth
import (

	"secret-im/pkg/utils-tools"
	"time"
)

// 验证码
type StoredVerificationCode struct {
	Code      string  `json:"code"`
	Timestamp int64   `json:"timestamp"`
	PushCode  string `json:"pushCode,omitempty"`
}

// 验证码是否有效
func (verificationCode *StoredVerificationCode) IsValid(theirCodeString string) bool {
	duration := int64((time.Minute * 10) / time.Millisecond)
	if verificationCode.Timestamp+duration < utils.CurrentTimeMillis() {
		return false
	}

	if len(verificationCode.Code) == 0 || len(theirCodeString) == 0 {
		return false
	}
	return verificationCode.Code == theirCodeString
}