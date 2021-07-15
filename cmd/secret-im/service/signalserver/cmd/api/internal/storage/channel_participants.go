package storage

import (
	"encoding/json"
	"fmt"
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
	"secret-im/pkg/driver/cassa"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"strings"
	"time"
)

// 频道成员管理
type ChannelParticipants struct {
}

// 获取成员信息
func (ChannelParticipants) Get(channelID, userID string) (ChannelParticipant, error) {
	var result model.ChannelParticipant
	relations := []gocassa.Relation{
		gocassa.Eq("channel_id", channelID),
		gocassa.Eq("user_id", userID),
	}
	filter := internal.cassa.Model(model.ChannelParticipant{}).Where(relations...)
	err := filter.ReadOne(&result).Run()
	if err != nil {
		return ChannelParticipant{}, err
	}

	participant := ChannelParticipant{
		ParticipantID: ParticipantID{
			UserID:    result.UserID.String(),
			ChannelID: result.ChannelID.String(),
		},
		Name:        result.Name,
		Left:        result.Left,
		Kicked:      result.Kicked,
		Banned:      result.Banned,
		AdminRights: ChannelAdminRights(result.AdminRights),
		Date:        result.Date.Unix(),
	}
	json.Unmarshal(result.NotifySettings, &participant.ChannelNotifySettings)
	return participant, nil
}

// 获取成员列表
func (ChannelParticipants) GetList(channelID string, userIDs []string) ([]ChannelParticipant, error) {
	tableName := internal.cassa.TableName(model.ChannelParticipant{})
	cql := fmt.Sprintf("SELECT * FROM %s WHERE channel_id=? AND user_id IN (%s)",
		tableName, strings.Join(userIDs, ","))

	iter := internal.cassa.Query(cql, channelID).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	var records []model.ChannelParticipant
	scanner, err := cassa.NewScanner(rowData, model.ChannelParticipant{})
	if err != nil {
		return nil, err
	}
	if err = scanner.ScanRows(iter, rowData, &records); err != nil {
		return nil, err
	}

	participants := make([]ChannelParticipant, 0, len(records))
	for _, record := range records {
		participant := ChannelParticipant{
			ParticipantID: ParticipantID{
				UserID:    record.UserID.String(),
				ChannelID: record.ChannelID.String(),
			},
			Name:        record.Name,
			Left:        record.Left,
			Kicked:      record.Kicked,
			Banned:      record.Banned,
			AdminRights: ChannelAdminRights(record.AdminRights),
			Date:        record.Date.Unix(),
		}
		json.Unmarshal(record.NotifySettings, &participant.ChannelNotifySettings)
		participants = append(participants, participant)
	}
	return participants, nil
}

// 获取成员数量
func (ChannelParticipants) GetCount(channelID string) (int, error) {
	var result struct {
		Count int
	}
	relations := []gocassa.Relation{
		gocassa.Eq("channel_id", channelID),
		gocassa.Eq("left", false),
		gocassa.Eq("kicked", false),
	}
	filter := internal.cassa.Model(model.ChannelParticipant{}).Where(relations...)
	op := filter.ReadOne(&result).WithOptions(gocassa.Options{Select: []string{"COUNT(channel_id)"}})
	if err := op.Run(); err != nil {
		return 0, nil
	}
	return result.Count, nil
}

// 获取成员列表
func (ChannelParticipants) GetUsers(channelID string, includedLeft bool) ([]string, error) {
	tableName := internal.cassa.TableName(model.ChannelParticipant{})
	cql := fmt.Sprintf("SELECT user_id FROM %s WHERE channel_id = ?", tableName)
	if !includedLeft {
		cql += " AND kicked=false AND left=false ALLOW FILTERING"
	}

	iter := internal.cassa.Query(cql, channelID).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	participants := make([]string, 0, iter.NumRows())
	for iter.Scan(rowData.Values...) {
		if rowData.Values[0] == nil {
			continue
		}
		userID, ok := rowData.Values[0].(*gocql.UUID)
		if !ok || userID == nil {
			continue
		}
		participants = append(participants, userID.String())
	}
	return participants, nil
}

// 更新成员名称
func (ChannelParticipants) UpdateName(channelID, userID string, name string) error {
	tableName := internal.cassa.TableName(model.ChannelParticipant{})
	stmt := "UPDATE %s SET name = ? WHERE user_id = ? AND channel_id = ? IF EXISTS"
	return internal.cassa.Query(fmt.Sprintf(stmt, tableName), name, userID, channelID).Exec()
}

// 更新管理权限
func (ChannelParticipants) UpdateAdminRights(channelID, userID string, rights ChannelAdminRights) error {
	tableName := internal.cassa.TableName(model.ChannelParticipant{})
	stmt := "UPDATE %s SET admin_rights = ? WHERE user_id = ? AND channel_id = ? IF EXISTS"
	return internal.cassa.Query(fmt.Sprintf(stmt, tableName), int(rights), userID, channelID).Exec()
}

// 更新通知设置
func (ChannelParticipants) UpdateNotifySettings(channelID, userID string, notifySettings ChannelNotifySettings) error {
	data, err := json.Marshal(notifySettings)
	if err != nil {
		return err
	}

	tableName := internal.cassa.TableName(model.ChannelParticipant{})
	stmt := "UPDATE %s SET notify_settings = ? WHERE user_id = ? AND channel_id = ? IF EXISTS"
	return internal.cassa.Query(fmt.Sprintf(stmt, tableName), data, userID, channelID).Exec()
}

// 更新成员操作
func (ChannelParticipants) updateOp(participant ChannelParticipant) gocassa.Op {
	relations := []gocassa.Relation{
		gocassa.Eq("channel_id", participant.ChannelID),
		gocassa.Eq("user_id", participant.UserID),
	}

	notifySettings, _ := json.Marshal(participant.ChannelNotifySettings)
	filter := internal.cassa.Model(model.ChannelParticipant{}).Where(relations...)
	op := filter.Update(map[string]interface{}{
		"name":            participant.Name,
		"left":            participant.Left,
		"kicked":          participant.Kicked,
		"admin_rights":    int(participant.AdminRights),
		"notify_settings": notifySettings,
		"date":            time.Unix(participant.Date, 0),
	})
	return op
}

// 批量插入操作
func (m ChannelParticipants) batchInsertOp(participants []ChannelParticipant) (gocassa.Op, error) {
	var batchOp gocassa.Op
	for _, participant := range participants {
		op := m.updateOp(participant)
		if batchOp == nil {
			batchOp = op
		} else {
			batchOp = batchOp.Add(op)
		}
	}
	return batchOp, nil
}
