package push


import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/sideshow/apns2"
)

// APN消息发送器
type ApnSender struct {
	client *apns2.Client
}

// 创建消息发送器
func NewApnSender(cert tls.Certificate, sandbox bool) *ApnSender {
	client := apns2.NewClient(cert).Production()
	if sandbox {
		client = client.Development()
	}
	return &ApnSender{client: client}
}

// 发送消息
func (sender *ApnSender) SendMessage(ctx context.Context,
	apnID, topic string, payload []byte, expiration time.Time) (*apns2.Response, error) {

	notification := apns2.Notification{
		DeviceToken: apnID,
		Topic:       topic,
		Payload:     payload,
		Expiration:  expiration,
		Priority:    10,
	}
	if strings.HasSuffix(topic, ".voip") {
		notification.PushType = apns2.PushTypeVOIP
	}
	return sender.client.PushWithContext(ctx, &notification)
}
