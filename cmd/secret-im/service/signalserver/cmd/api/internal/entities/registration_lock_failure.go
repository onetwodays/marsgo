package entities

import "secret-im/service/signalserver/cmd/api/internal/auth"

// 注册锁定失败
type RegistrationLockFailure struct {
	TimeRemaining     int64                            `json:"timeRemaining"`
	BackupCredentials *auth.ExternalServiceCredentials `json:"backupCredentials"`
}
