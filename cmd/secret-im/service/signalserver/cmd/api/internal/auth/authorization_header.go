package auth

import (
	"encoding/base64"
	"errors"
	"github.com/tal-tech/go-zero/core/logx"
	"strconv"
	"strings"
)

// 鉴权头部
type AuthorizationHeader struct {
	Identifier AmbiguousIdentifier
	DeviceID   int64
	Password   string
}

// 根据请求头创建
func (authorizationHeader *AuthorizationHeader) FromFullHeader(header string) (*AuthorizationHeader, error) {
	headerParts := strings.Split(header, " ")
	if len(headerParts) < 2 {
		return nil, errors.New("invalid authorization header")
	}

	if headerParts[0] != "Basic" {
		return nil, errors.New("unsupported authorization method")
	}

	concatenatedValues, err := base64.StdEncoding.DecodeString(headerParts[1])
	if err != nil || len(concatenatedValues) == 0 {
		return nil, errors.New("invalid authorization header")
	}
	logx.Info("basic auth valus is ",string(concatenatedValues))

	credentialParts := strings.Split(string(concatenatedValues), ":")
	if len(credentialParts) < 2 {
		return nil, errors.New("badly formated credentials")
	}
	return authorizationHeader.FromUserAndPassword(credentialParts[0], credentialParts[1])
}

// 根据用户密码创建,没有设备id，就默认是1
func (authorizationHeader *AuthorizationHeader) FromUserAndPassword(user, pwd string) (*AuthorizationHeader, error) {
	deviceID := int64(1)

	numberAndID := strings.Split(user, ".")
	if len(numberAndID) > 1 {

		var err error
		deviceID, err = strconv.ParseInt(numberAndID[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}
	authorizationHeader.Identifier = *NewAmbiguousIdentifier(numberAndID[0])
	authorizationHeader.DeviceID = deviceID
	authorizationHeader.Password = pwd
	return authorizationHeader, nil
}

