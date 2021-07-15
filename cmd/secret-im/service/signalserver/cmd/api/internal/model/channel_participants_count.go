package model
import (
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

// 频道成员数量
type ChannelParticipantsCount struct {
	ChannelID         gocql.UUID      `cql:"channel_id"`
	ParticipantsCount gocassa.Counter `cql:"participants_count"`
}

func (ChannelParticipantsCount) TableName() string {
	return "channel_participants_count"
}

func (ChannelParticipantsCount) Keys() gocassa.Keys {
	return gocassa.Keys{
		PartitionKeys: []string{"channel_id"},
	}
}

