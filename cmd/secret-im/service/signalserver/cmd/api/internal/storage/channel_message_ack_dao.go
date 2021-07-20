package storage

import (
	"fmt"
	"secret-im/pkg/driver/cassa"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"strings"
)

// 频道最后确认消息
type ChannelMessageAckDao struct {
}

// 获取最后确认消息
func (ChannelMessageAckDao) GetMessages(userID string, channelIDs []string) (map[string]int64, error) {
	tableName := internal.cassa.TableName(model.ChannelMessageAck{})
	cql := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND channel_id IN (%s)",
		tableName, strings.Join(channelIDs, ","))

	iter := internal.cassa.Query(cql, userID).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	var records []model.ChannelMessageAck
	scanner, err := cassa.NewScanner(rowData, model.ChannelMessageAck{})
	if err != nil {
		return nil, err
	}
	if err = scanner.ScanRows(iter, rowData, &records); err != nil {
		return nil, err
	}

	mapper := make(map[string]int64)
	for _, record := range records {
		mapper[record.ChannelID.String()] = record.LastAckMessage
	}
	return mapper, nil
}

// 更新最后确认消息
func (ChannelMessageAckDao) UpdateLastAckMessage(userID, channelID string, messageID int64) error {
	var lastAckMessage int64
	tableName := internal.cassa.TableName(model.ChannelMessageAck{})
	cql := fmt.Sprintf(
		"UPDATE %s SET last_ack_message=? WHERE user_id=? AND channel_id=? IF last_ack_message<?", tableName)
	applied, err := internal.cassa.Query(cql, messageID, userID, channelID, messageID).ScanCAS(&lastAckMessage)
	if err != nil {
		return err
	}
	if applied || lastAckMessage > 0 {
		return nil
	}

	cql = fmt.Sprintf(
		"UPDATE %s SET last_ack_message=? WHERE user_id=? AND channel_id=? IF last_ack_message=null", tableName)
	applied, err = internal.cassa.Query(cql, messageID, userID, channelID).ScanCAS(&lastAckMessage)
	return err
}
