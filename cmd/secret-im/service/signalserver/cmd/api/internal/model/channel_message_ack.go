package model

import (
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

// 频道消息确认
type ChannelMessageAck struct {
	UserID         gocql.UUID `cql:"user_id"`
	ChannelID      gocql.UUID `cql:"channel_id"`
	LastAckMessage int64      `cql:"last_ack_message"`
}

func (ChannelMessageAck) TableName() string {
	return "channel_message_ack"
}

func (ChannelMessageAck) Keys() gocassa.Keys {
	return gocassa.Keys{
		PartitionKeys:     []string{"user_id"},
		ClusteringColumns: []string{"channel_id"},
	}
}

