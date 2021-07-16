package storage

import (
	"encoding/json"
	"fmt"
	"secret-im/pkg/driver/cassa"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"time"
	"github.com/gocassa/gocassa"
	"errors"
)

// 频道消息管理
type ChannelMessages struct {
}

// 插入消息
func (m ChannelMessages) Insert(message *ChannelMessage) error {
	var err error
	message.MessageID, err = ChannelMessageCounter{}.Incr(message.ChannelID)
	if err != nil {
		return err
	}

	record, err := message.toRecord()
	if err != nil {
		return err
	}
	return internal.cassa.Model(record).Set(record).Run()
}

// 批量插入
func (m ChannelMessages) BatchInsert(channelID string, messages []*ChannelMessage) error {
	if len(messages) == 0 {
		return nil
	}

	maxMessageID, err := ChannelMessageCounter{}.IncrBy(channelID, int64(len(messages)))
	if err != nil {
		return err
	}

	var op gocassa.Op
	for idx, message := range messages {
		message.MessageID = maxMessageID - int64(len(messages)-idx-1)
		record, err := message.toRecord()
		if err != nil {
			return err
		}

		if op == nil {
			op = internal.cassa.Model(model.ChannelMessage{}).Set(record)
		} else {
			op = op.Add(internal.cassa.Model(model.ChannelMessage{}).Set(record))
		}
	}
	return op.Run()
}

// 获取指定消息
func (m ChannelMessages) Get(channelID string, messageID int64) (ChannelMessage, error) {
	var record model.ChannelMessage
	relations := []gocassa.Relation{
		gocassa.Eq("channel_id", channelID),
		gocassa.Eq("bucket", calcMessageBucket(messageID)),
		gocassa.Eq("message_id", messageID),
	}
	err := internal.cassa.Model(model.ChannelMessage{}).Where(relations...).ReadOne(&record).Run()
	if err != nil {
		return ChannelMessage{}, err
	}
	return newChannelMessageFromRecord(&record), nil
}

// 获取范围消息
func (m ChannelMessages) GetMessages(channelID string, firstMessageID, lastMessageID int64) ([]ChannelMessage, error) {
	if lastMessageID-firstMessageID > 100 {
		return nil, errors.New("filter range too large")
	}

	buckets := make([]interface{}, 0)
	for id := firstMessageID; id <= lastMessageID; id++ {
		bucket := calcMessageBucket(id)
		if len(buckets) == 0 {
			buckets = append(buckets, bucket)
			continue
		}
		if bucket != buckets[len(buckets)-1] {
			buckets = append(buckets, bucket)
		}
	}

	tableName := internal.cassa.TableName(model.ChannelMessage{})
	format := "SELECT * FROM %s WHERE channel_id=? AND bucket IN ? AND message_id >= ? AND message_id <= ?"
	cql := fmt.Sprintf(format, tableName)

	iter := internal.cassa.Query(cql, channelID, buckets, firstMessageID, lastMessageID).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	var records []model.ChannelMessage
	scanner, err := cassa.NewScanner(rowData, model.ChannelMessage{})
	if err != nil {
		return nil, err
	}
	if err = scanner.ScanRows(iter, rowData, &records); err != nil {
		return nil, err
	}

	messages := make([]ChannelMessage, 0, len(records))
	for idx := range records {
		messages = append(messages, newChannelMessageFromRecord(&records[idx]))
	}
	return messages, nil
}

// 编辑消息内容
func (m ChannelMessages) EditContent(channelID string, messageID int64, editor, content string) error {
	jsb, err := json.Marshal(ChannelMessageEditor{
		UUID:     editor,
		EditedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"editor":  jsb,
		"content": content,
	}
	relations := []gocassa.Relation{
		gocassa.Eq("channel_id", channelID),
		gocassa.Eq("bucket", calcMessageBucket(messageID)),
		gocassa.Eq("message_id", messageID),
	}
	return internal.cassa.Model(model.ChannelMessage{}).Where(relations...).Update(updates).Run()
}

// 删除指定消息
func (m ChannelMessages) Delete(channelID string, messageID int64) error {
	updates := map[string]interface{}{
		"deleted": true,
	}
	relations := []gocassa.Relation{
		gocassa.Eq("channel_id", channelID),
		gocassa.Eq("bucket", calcMessageBucket(messageID)),
		gocassa.Eq("message_id", messageID),
	}
	return internal.cassa.Model(model.ChannelMessage{}).Where(relations...).Update(updates).Run()
}

