package storage

import (
	"fmt"
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
	"secret-im/pkg/driver/cassa"
	"secret-im/service/signalserver/cmd/api/internal/model"

)

// 已加入频道管理
type ChannelJoinedDao struct {
}

// 获取频道列表
func (m ChannelJoinedDao) GetChannels(userID string) ([]string, error) {
	tableName := internal.cassa.TableName(model.ChannelJoined{})
	cql := fmt.Sprintf("SELECT channel_id FROM %s WHERE user_id = ?", tableName)
	cql += " AND kicked=false AND left=false ALLOW FILTERING"

	iter := internal.cassa.Query(cql, userID).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	channels := make([]string, 0, iter.NumRows())
	for iter.Scan(rowData.Values...) {
		if rowData.Values[0] == nil {
			continue
		}
		channelID, ok := rowData.Values[0].(*gocql.UUID)
		if !ok || channelID == nil {
			continue
		}
		channels = append(channels, channelID.String())
	}
	return channels, nil
}

// 筛选频道列表
func (m ChannelJoinedDao) Filter(userID string, maxChannelID *string, limit int) ([]ChannelJoined, error) {
	tableName := internal.cassa.TableName(model.ChannelJoined{})
	cql := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND left=false AND kicked=false", tableName)
	if maxChannelID != nil {
		cql += " AND channel_id>" + *maxChannelID
	}
	cql += " LIMIT ? ALLOW FILTERING"

	iter := internal.cassa.Query(cql, userID, limit).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	var records []model.ChannelJoined
	scanner, err := cassa.NewScanner(rowData, model.ChannelJoined{})
	if err != nil {
		return nil, err
	}
	if err = scanner.ScanRows(iter, rowData, &records); err != nil {
		return nil, err
	}

	joinedChannel := make([]ChannelJoined, 0, len(records))
	for _, record := range records {
		joinedChannel = append(joinedChannel, ChannelJoined{
			UserID:    record.UserID.String(),
			ChannelID: record.ChannelID.String(),
			Left:      record.Left,
			Kicked:    record.Kicked,
		})
	}
	return joinedChannel, nil
}

// 更新频道已停用
func (m ChannelJoinedDao) UpdateDeactivated(userID string, channelIDs []string) error {
	data := make([]ChannelJoined, 0, len(channelIDs))
	for _, channelID := range channelIDs {
		data = append(data, ChannelJoined{
			UserID:    userID,
			ChannelID: channelID,
			Left:      true,
			Kicked:    false,
		})
	}
	op, err := m.batchUpdateOp(data)
	if err != nil {
		return err
	}
	return op.Run()
}

// 批量更新操作
func (ChannelJoinedDao) batchUpdateOp(data []ChannelJoined) (gocassa.Op, error) {
	var op gocassa.Op
	for _, channelJoined := range data {
		relations := []gocassa.Relation{
			gocassa.Eq("user_id", channelJoined.UserID),
			gocassa.Eq("channel_id", channelJoined.ChannelID),
		}
		updates := map[string]interface{}{
			"left":   channelJoined.Left,
			"kicked": channelJoined.Kicked,
		}
		if op == nil {
			op = internal.cassa.Model(model.ChannelJoined{}).Where(relations...).Update(updates)
		} else {
			op = op.Add(internal.cassa.Model(model.ChannelJoined{}).Where(relations...).Update(updates))
		}
	}
	return op, nil
}

