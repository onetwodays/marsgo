package entities

import (
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/types"
)

// 默认ID
const DeviceMasterID int64 = 1


// 预共享密钥
type PreKey struct {
	KeyID     int64  `json:"keyId"`
	PublicKey string `json:"publicKey"`
}

// 预签名密钥
type SignedPreKey struct {
	PreKey
	Signature string `json:"signature"`
}


// 设备信息
type Device struct {
	ID              int64   `json:"id,omitempty"`
	SignalingKey    string  `json:"signalingKey,omitempty"`
	GcmID           string  `json:"gcmId,omitempty"`
	ApnID           string  `json:"apnId,omitempty"`
	VoipApnID       string  `json:"voipApnId,omitempty"`
	FetchesMessages bool    `json:"fetchesMessages,omitempty"`
}

// 完整设备信息
type DeviceFull struct {
	Device
	Name                string                   `json:"name,omitempty"`
	AuthToken           string                   `json:"authToken,omitempty"`
	Salt                string                   `json:"salt,omitempty"`
	PushTimestamp       int64                    `json:"pushTimestamp,omitempty"`
	UninstalledFeedback int64                    `json:"uninstalledFeedback,omitempty"`
	RegistrationID      int                      `json:"registrationId,omitempty"`
	SignedPreKey        *types.SignedPrekey            `json:"signedPreKey,omitempty"`
	LastSeen            int64                    `json:"lastSeen,omitempty"`
	Created             int64                    `json:"created,omitempty"`
	UserAgent           string                   `json:"userAgent,omitempty"`
	Capabilities        types.DeviceCapabilities `json:"capabilities,omitempty"`
}

// 设备所在区域
type DevicePartition struct {
	UUID      string `json:"uuid"`
	DeviceID  int64  `json:"device_id"`
	Partition int    `json:"partition"`
}

// 是否已启用
func (device *DeviceFull) IsEnabled() bool {
	hasChannel := device.FetchesMessages || len(device.ApnID) != 0 ||  len(device.GcmID)!=0
	return (device.ID == DeviceMasterID && hasChannel && device.SignedPreKey != nil) ||
		(device.ID != DeviceMasterID && hasChannel && device.SignedPreKey != nil) &&
			device.LastSeen > (utils.CurrentTimeMillis()-utils.DaysToMillis(30))
}

// 设置APN ID
func (device *DeviceFull) SetApnID(apnID string) {
	device.ApnID = apnID
	if len(apnID) !=0 {
		device.PushTimestamp = utils.CurrentTimeMillis()
	}
}

// 设置GCM ID
func (device *DeviceFull) SetGcmID(gcmID string) {
	device.GcmID = gcmID
	if len(gcmID) != 0 {
		device.PushTimestamp = utils.CurrentTimeMillis()
	}
}

// 设置认证凭据
func (device *DeviceFull) SetAuthenticationCredentials(credentials auth.AuthenticationCredentials) {
	device.Salt = credentials.Salt
	device.AuthToken = credentials.HashedAuthenticationToken
}

// 获取认证凭据
func (device *DeviceFull) GetAuthenticationCredentials() auth.AuthenticationCredentials {
	return auth.AuthenticationCredentials{Salt: device.Salt, HashedAuthenticationToken: device.AuthToken}
}

