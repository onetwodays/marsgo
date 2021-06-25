package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"secret-im/pkg/utils-tools"
	"strconv"
)

// 认证凭据
type AuthenticationCredentials struct {
	Salt                      string
	HashedAuthenticationToken string
}

// 创建认证凭据
func NewAuthenticationCredentials(authenticationToken string) AuthenticationCredentials {
	salt := strconv.Itoa(utils.SecureRandInt())
	credentials := AuthenticationCredentials{Salt: salt}
	credentials.HashedAuthenticationToken = credentials.GetHashedValue(salt, authenticationToken)
	return credentials
}

// 验证凭据
func (credentials AuthenticationCredentials) Verify(authenticationToken string) bool {
	theirValue := credentials.GetHashedValue(credentials.Salt, authenticationToken)
	return theirValue == credentials.HashedAuthenticationToken
}

// 获取哈希值
func (AuthenticationCredentials) GetHashedValue(salt, token string) string {
	hash := sha1.New()
	hash.Write([]byte(salt + token))
	return hex.EncodeToString(hash.Sum(nil))
}
