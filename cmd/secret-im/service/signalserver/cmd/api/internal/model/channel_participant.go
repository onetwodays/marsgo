package model
import (
	"time"

	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

// 频道成员
type ChannelParticipant struct {
	ChannelID      gocql.UUID `cql:"channel_id"`
	UserID         gocql.UUID `cql:"user_id"`
	Name           string     `cql:"name"`
	Left           bool       `cql:"left"`
	Kicked         bool       `cql:"kicked"`
	Banned         bool       `cql:"banned"`
	AdminRights    int        `cql:"admin_rights"`
	NotifySettings []byte     `cql:"notify_settings"`
	Date           time.Time  `cql:"date"`
}

func (ChannelParticipant) TableName() string {
	return "channel_participants"
}

func (ChannelParticipant) Keys() gocassa.Keys {
	return gocassa.Keys{
		PartitionKeys:     []string{"channel_id"},
		ClusteringColumns: []string{"user_id"},
	}
}

