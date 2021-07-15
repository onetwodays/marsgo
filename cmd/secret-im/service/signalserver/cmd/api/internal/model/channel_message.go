package model

import (
	"time"

	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

// 消息类型
type ChannelMessageType int8

const (
	_                         ChannelMessageType = iota
	ChannelMessageTypeNormal  ChannelMessageType = 1
	ChannelMessageTypeService ChannelMessageType = 2
)

func (t ChannelMessageType) String() string {
	switch t {
	case ChannelMessageTypeNormal:
		return "message"
	case ChannelMessageTypeService:
		return "service"
	default:
		return "unknown"
	}
}

// 频道消息
type ChannelMessage struct {
	ChannelID       gocql.UUID         `cql:"channel_id"`
	Bucket          int                `cql:"bucket"`
	MessageID       int64              `cql:"message_id"`
	Type            ChannelMessageType `cql:"type"`
	Source          gocql.UUID         `cql:"source"`
	SourceDevice    int64              `cql:"source_device"`
	Content         []byte             `cql:"content"`
	Action          []byte             `cql:"action"`
	Relay           string             `cql:"relay"`
	Deleted         bool               `cql:"deleted"`
	Editor          []byte             `cql:"editor"`
	Timestamp       time.Time          `cql:"timestamp"`
	ServerTimestamp time.Time          `cql:"server_timestamp"`
}

func (ChannelMessage) TableName() string {
	return "channel_messages"
}

func (ChannelMessage) Keys() gocassa.Keys {
	return gocassa.Keys{
		PartitionKeys:     []string{"channel_id", "bucket"},
		ClusteringColumns: []string{"message_id"},
	}
}

