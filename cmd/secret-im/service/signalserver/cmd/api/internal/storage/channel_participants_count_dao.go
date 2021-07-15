package storage

import (
	"fmt"
	"github.com/gocassa/gocassa"
	"secret-im/pkg/driver/cassa"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"strings"
)


// 频道成员数量管理
type ChannelParticipantsCountDao struct {
}

// 自增成员数量
func (ChannelParticipantsCountDao) incrOp(channelID string, nums int) gocassa.Op {
	table := internal.cassa.Model(model.ChannelParticipantsCount{})
	filter := table.Where(gocassa.Eq("channel_id", channelID))
	return filter.Update(map[string]interface{}{"participants_count": gocassa.CounterIncrement(nums)})
}

// 获取成员数量
func (ChannelParticipantsCountDao) GetChannels(channelIDs []string) (map[string]int, error) {
	tableName := internal.cassa.TableName(model.ChannelParticipantsCount{})
	cql := fmt.Sprintf("SELECT * FROM %s WHERE channel_id IN (%s)",
		tableName, strings.Join(channelIDs, ","))

	iter := internal.cassa.Query(cql).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	var records []model.ChannelParticipantsCount
	scanner, err := cassa.NewScanner(rowData, model.ChannelParticipantsCount{})
	if err != nil {
		return nil, err
	}
	if err = scanner.ScanRows(iter, rowData, &records); err != nil {
		return nil, err
	}

	mapper := make(map[string]int)
	for _, record := range records {
		mapper[record.ChannelID.String()] = int(record.ParticipantsCount)
	}
	return mapper, nil
}
