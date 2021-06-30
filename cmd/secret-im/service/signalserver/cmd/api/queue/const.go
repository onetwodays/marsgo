package queue

import (
	"fmt"
)

var (
	// 发送短信
	SendSmsTopic = "sms_topic"
	// 发送语音短信
	SendVoiceSmsTopic = "voice_sms_topic"
	// 删除文件对象
	DeleteObjectTopic = "del_object_topic"
	// 发送APN消息
	SendApnMessageTopic = "apn_message_topic"
)

// 发送到设备
func SendToDeviceTopic(partition int) string {
	return fmt.Sprintf("send_to_device_%d_topic", partition)
}

