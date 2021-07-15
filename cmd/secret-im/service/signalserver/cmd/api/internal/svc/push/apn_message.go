package push

import (
	"fmt"
	"secret-im/service/signalserver/cmd/api/queue"
	"time"
)

var (
	ApnNotificationPayload = `{"aps":{"sound":"default","alert":{"loc-key":"APN_Message"}}}`
	ApnChallengePayload    = `{"aps":{"sound":"default","alert":{"loc-key":"APN_Message"}},"challenge":"%s"}`
)

// APN消息
type ApnMessage struct {
	ApnID         string
	Number        string
	DeviceID      int64
	IsVoip        bool
	ChallengeData *string
}

// 获取消息内容
func (msg ApnMessage) GetMessage() string {
	if msg.ChallengeData == nil {
		return ApnNotificationPayload
	}
	return fmt.Sprintf(ApnChallengePayload, *msg.ChallengeData)
}

// 获取到期时间
func (msg ApnMessage) GetExpirationTime() time.Time {
	return time.Now().AddDate(1, 0, 0)
}

// 添加到推送队列
func AddToApnMessageQueue(msg ApnMessage) error {
	topic := "apntest" //conf.GetServer().APN.Topic
	if msg.IsVoip {
		topic = topic + ".voip"
	}

	message := msg.GetMessage()
	expirationTime := msg.GetExpirationTime().Unix()
	return queue.Publish(queue.SendApnMessageTopic, &queue.SendApnMessage{
		ApnId:          msg.ApnID,
		Topic:          topic,
		Number:         msg.Number,
		DeviceId:       msg.DeviceID,
		IsVoip:         msg.IsVoip,
		Message:        message,
		ExpirationTime: expirationTime,
	})
}

