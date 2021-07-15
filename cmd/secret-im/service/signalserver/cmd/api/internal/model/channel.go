package model
import (
	"time"

	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

// 频道信息
type Channel struct {
	ChannelID   gocql.UUID `cql:"channel_id"`
	Creator     gocql.UUID `cql:"creator"`
	Profile     []byte     `cql:"profile"`
	Public      bool       `cql:"public"`
	Deactivated bool       `cql:"deactivated"`
	Date        time.Time  `cql:"date"`
}

func (Channel) TableName() string {
	return "channels"
}

func (Channel) Keys() gocassa.Keys {
	return gocassa.Keys{
		PartitionKeys: []string{"channel_id"},
	}
}
