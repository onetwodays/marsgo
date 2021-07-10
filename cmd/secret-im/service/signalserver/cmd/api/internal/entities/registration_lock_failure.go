package entities

import (
	"secret-im/service/signalserver/cmd/api/internal/types"
)

// 注册锁定失败
type RegistrationLockFailure struct {
	TimeRemaining     int64                            `json:"timeRemaining"`
	BackupCredentials *types.ExternalServiceCredentials `json:"backupCredentials"`
}
