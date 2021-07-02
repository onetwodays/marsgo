package operation

import "fmt"

var (
	PendingNotificationsKey = "pending_apn"
)

func GetEndpointKey(number string, deviceID int64) string {
	return fmt.Sprintf("apn_device::{%s}::%d", number, deviceID)
}

