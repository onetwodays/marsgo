package storage

import (
	"bytes"
	"encoding/json"
	"github.com/gocql/gocql"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"time"
)

// 计算消息bucket
func calcMessageBucket(messageID int64) int {
	return int(messageID/100000 + 1)
}

// 频道消息
type ChannelMessage struct {
	ChannelID       string                    `json:"channelId"`
	MessageID       int64                     `json:"messageId"`
	Type            model.ChannelMessageType `json:"type"`
	Source          *string                   `json:"source"`
	SourceDevice    *int64                    `json:"sourceDevice"`
	Content         utils.Base64Bytes         `json:"content,omitempty"`
	Action          utils.Base64Bytes         `json:"action,omitempty"`
	Relay           *string                   `json:"relay,omitempty"`
	Editor          *ChannelMessageEditor     `json:"editor,omitempty"`
	Deleted         bool                      `json:"deleted,omitempty"`
	Timestamp       int64                     `json:"timestamp"`
	ServerTimestamp int64                     `json:"serverTimestamp"`
}

// 消息编辑信息
type ChannelMessageEditor struct {
	UUID     string `json:"uuid"`
	EditedAt int64  `json:"edited_at"`
}

// 转换为Record,保存到数据库
func (message *ChannelMessage) toRecord() (model.ChannelMessage, error) {
	var source *gocql.UUID
	if message.Source != nil {
		uuid, err := gocql.ParseUUID(*message.Source)
		if err != nil {
			return model.ChannelMessage{}, err
		}
		source = &uuid
	}

	channelID, err := gocql.ParseUUID(message.ChannelID)
	if err != nil {
		return model.ChannelMessage{}, err
	}

	record := model.ChannelMessage{
		ChannelID:       channelID,
		Bucket:          calcMessageBucket(message.MessageID),
		MessageID:       message.MessageID,
		Type:            message.Type,
		Content:         message.Content,
		Action:          message.Action,
		Timestamp:       time.Unix(message.Timestamp, 0),
		ServerTimestamp: time.Unix(message.ServerTimestamp, 0),
	}
	if source != nil {
		record.Source = *source
	}
	if message.SourceDevice != nil {
		record.SourceDevice = *message.SourceDevice
	}
	if message.Relay != nil {
		record.Relay = *message.Relay
	}

	if message.Editor != nil {
		jsb, err := json.Marshal(message.Editor)
		if err == nil {
			record.Editor = jsb
		}
	}
	return record, nil
}

// 从Record创建频道消息
func newChannelMessageFromRecord(record *model.ChannelMessage) ChannelMessage {
	channelMessage := ChannelMessage{
		ChannelID:       record.ChannelID.String(),
		MessageID:       record.MessageID,
		Type:            record.Type,
		Content:         record.Content,
		Action:          record.Action,
		Timestamp:       record.Timestamp.Unix(),
		ServerTimestamp: record.ServerTimestamp.Unix(),
		Deleted:         record.Deleted,
	}
	if bytes.Compare(record.Source.Bytes(), gocql.UUID{}.Bytes()) != 0 {
		s := record.Source.String()
		channelMessage.Source = &s
	}
	if record.SourceDevice > 0 {
		channelMessage.SourceDevice = &record.SourceDevice
	}
	if len(record.Relay) > 0 {
		channelMessage.Relay = &record.Relay
	}

	if record.Editor != nil {
		var editor ChannelMessageEditor
		if err := json.Unmarshal(record.Editor, &editor); err == nil {
			channelMessage.Editor = &editor
		}
	}
	return channelMessage
}

