package operation
import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 消息键
type Key struct {
	Address                           string
	DeviceID                          int64
	UserMessageQueue                  string
	UserMessageQueueMetadata          string
	UserMessageQueuePersistInProgress string
}

// 创建消息键
func NewKey(address string, deviceID int64) Key {
	return Key{
		Address:                           address,
		DeviceID:                          deviceID,
		UserMessageQueue:                  fmt.Sprintf("user_queue::{%s}::%d", address, deviceID),
		UserMessageQueueMetadata:          fmt.Sprintf("user_queue_metadata::{%s}::%d", address, deviceID),
		UserMessageQueuePersistInProgress: fmt.Sprintf("user_queue_persisting::{%s}::%d", address, deviceID),
	}
}

// 用户消息队列索引
func (p Key) UserMessageQueueIndex() string {
	return "user_queue_index"
}

// 根据用户消息队列名创建
func (p Key) FromUserMessageQueue(userMessageQueue string) (Key, error) {
	parts := strings.Split(userMessageQueue, "::")
	if len(parts) != 3 {
		return Key{}, errors.New("malformed key")
	}

	address := parts[1]
	if len(address) > 1 && address[0] == '{' && address[len(address)-1] == '}' {
		address = address[1 : len(address)-1]
	}
	deviceID, _ := strconv.ParseInt(parts[2], 10, 64)
	return NewKey(address, deviceID), nil
}

