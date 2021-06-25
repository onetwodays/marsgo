package util

import (
	"strconv"
)

// 验证码
type VerificationCode struct {
	VerificationCode        string `json:"verificationCode"`
	VerificationCodeDisplay string `json:"-"`
	VerificationCodeSpeech  string `json:"-"`
}

// 验证码分隔
func DelimitVerificationCode(code string) string {
	var delimited string
	for i := 0; i < len(code); i++ {
		delimited += string(code[i])
		if i != len(code)-1 {
			delimited += ","
		}
	}
	return delimited
}

// 创建验证码
func NewVerificationCode(verificationCode int) VerificationCode {
	return VerificationCode{
		VerificationCode: strconv.FormatInt(int64(verificationCode), 10),
	}
}

// 从字符串创建验证码
func NewVerificationCodeFromString(verificationCode string) VerificationCode {
	verificationCodeDisplay := verificationCode
	if len(verificationCode) > 3 {
		verificationCodeDisplay = verificationCode[0:3] + "-" + verificationCode[3:]
	}
	return VerificationCode{
		VerificationCode:        verificationCode,
		VerificationCodeDisplay: verificationCodeDisplay,
		VerificationCodeSpeech:  DelimitVerificationCode(verificationCode),
	}
}
