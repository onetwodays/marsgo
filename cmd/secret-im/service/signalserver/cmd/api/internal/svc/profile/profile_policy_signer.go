package profile

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// 策略签名器
type PolicySigner struct {
	region          string
	awsAccessSecret string
}

// 创建策略签名器
func NewPolicySigner(awsAccessSecret, region string) *PolicySigner {
	return &PolicySigner{
		region:          region,
		awsAccessSecret: awsAccessSecret,
	}
}

// 获取签名
func (signer *PolicySigner) GetSignature(now time.Time, policy string) string {
	mac := hmac.New(sha256.New, []byte("AWS4"+signer.awsAccessSecret))
	mac.Write([]byte(now.UTC().Format("20060102")))
	dateKey := mac.Sum(nil)

	mac = hmac.New(sha256.New, dateKey)
	mac.Write([]byte(signer.region))
	dateRegionKey := mac.Sum(nil)

	mac = hmac.New(sha256.New, dateRegionKey)
	mac.Write([]byte("s3"))
	dateRegionServiceKey := mac.Sum(nil)

	mac = hmac.New(sha256.New, dateRegionServiceKey)
	mac.Write([]byte("aws4_request"))
	signingKey := mac.Sum(nil)

	mac = hmac.New(sha256.New, signingKey)
	mac.Write([]byte(policy))
	return hex.EncodeToString(mac.Sum(nil))
}