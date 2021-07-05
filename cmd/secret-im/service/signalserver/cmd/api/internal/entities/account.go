package entities

import (
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/auth"
)

// 帐号信息
type Account struct {
	UUID                           string        `json:"uuid,omitempty"`
	Number                         string        `json:"number,omitempty"`
	Devices                        []*DeviceFull `json:"devices,omitempty"`
	IdentityKey                    string        `json:"identityKey,omitempty"`
	Name                           string        `json:"name,omitempty"`
	Avatar                         string        `json:"avatar,omitempty"`
	AvatarDigest                   string        `json:"avatarDigest,omitempty"`
	Pin                            string        `json:"pin,omitempty"`
	RegistrationLock               string        `json:"registrationLock,omitempty"`
	RegistrationLockSalt           string        `json:"registrationLockSalt,omitempty"`
	UnidentifiedAccessKey          string        `json:"uak,omitempty"`
	UnrestrictedUnidentifiedAccess bool          `json:"uua,omitempty"`
	AuthenticatedDevice            *DeviceFull   `json:"-"`
}

// 上次活动日期
func (account *Account) GetLastSeen() int64 {
	var lastSeen int64
	for _, device := range account.Devices {
		if device.LastSeen > lastSeen {
			lastSeen = device.LastSeen
		}
	}
	return lastSeen
}

// 添加设备
func (account *Account) AddDevice(device *DeviceFull) {
	account.RemoveDevice(device.ID)
	account.Devices = append(account.Devices, device)
}

// 删除设备
func (account *Account) RemoveDevice(deviceID int64) {
	for idx := 0; idx < len(account.Devices); {
		if deviceID != account.Devices[idx].ID {
			idx++
			continue
		}
		account.Devices[idx] = account.Devices[len(account.Devices)-1]
		account.Devices = account.Devices[:len(account.Devices)-1]
	}
}


// 获取设备
func (account *Account) GetDevice(deviceID int64) (*DeviceFull, bool) {
	for _, device := range account.Devices {
		if deviceID == device.ID {
			return device, true
		}
	}
	return nil, false
}

// 获取启动设备数
func (account *Account) GetEnabledDeviceCount() int {
	var count int
	for _, device := range account.Devices {
		if device.IsEnabled() {
			count++
		}
	}
	return count
}


// 获取主设备
func (account *Account) GetMasterDevice() (*DeviceFull, bool) {
	return account.GetDevice(DeviceMasterID)
}

// 是否已启用
func (account *Account) IsEnabled() bool {
	masterDevice, ok := account.GetMasterDevice()
	if !ok {
		return false
	}
	return masterDevice.IsEnabled() &&
		masterDevice.LastSeen > (utils.CurrentTimeMillis()-utils.DaysToMillis(365))
}

// 是否自己
func (account *Account) IsFor(identifier *auth.AmbiguousIdentifier) bool {
	if len(identifier.UUID) != 0 {
		return account.UUID == identifier.UUID
	}
	if len(identifier.Number) != 0 {
		return account.Number == identifier.Number
	}
	return false
}