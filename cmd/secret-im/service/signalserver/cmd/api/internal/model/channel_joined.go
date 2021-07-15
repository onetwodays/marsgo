package model


import (
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

// 已加入频道
type ChannelJoined struct {
	UserID    gocql.UUID `cql:"user_id"`
	ChannelID gocql.UUID `cql:"channel_id"`
	Left      bool       `cql:"left"`
	Kicked    bool       `cql:"kicked"`
}

func (ChannelJoined) TableName() string {
	return "channel_joined"
}

func (ChannelJoined) Keys() gocassa.Keys {
	return gocassa.Keys{
		PartitionKeys:     []string{"user_id"},
		ClusteringColumns: []string{"channel_id"},
	}
}

