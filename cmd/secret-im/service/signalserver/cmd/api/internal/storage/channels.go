package storage

import (
	"encoding/json"
	"fmt"
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/pkg/driver/cassa"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"strings"
	"time"
)

// 频道信息管理
type Channels struct {
}

// 插入频道信息
func (Channels) Insert(channel *Channel, users []ChannelParticipant) error {
	// 插入基本信息
	profile, err := json.Marshal(channel.Profile)
	if err != nil {
		return err
	}
	creator, err := gocql.ParseUUID(channel.Creator)
	if err != nil {
		return err
	}
	channelID, err := gocql.ParseUUID(channel.ChannelID)
	if err != nil {
		return err
	}

	record := model.Channel{
		ChannelID:   channelID,
		Creator:     creator,
		Profile:     profile,
		Deactivated: false,
		Public:      channel.Public,
		Date:        time.Unix(channel.Date, 0),
	}
	batch := internal.cassa.Model(record).Set(record)

	// 插入已加入频道
	ops := make([]gocassa.Op, 0)
	ids := make([]string, 0, len(users))
	data := make([]ChannelJoined, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.UserID)
		data = append(data, ChannelJoined{
			UserID:    user.UserID,
			ChannelID: channel.ChannelID,
			Left:      user.Left,
			Kicked:    user.Kicked,
		})
	}
	updateJoinedOp, err := ChannelJoinedDao{}.batchUpdateOp(data)
	if err != nil {
		return err
	}
	if updateJoinedOp != nil {
		ops = append(ops, updateJoinedOp)
	}

	// 插入频道成员列表
	updateParticipantsOp, err := ChannelParticipants{}.batchInsertOp(users)
	if err != nil {
		return err
	}
	if updateParticipantsOp != nil {
		ops = append(ops, updateParticipantsOp)
	}

	// 更新频道成员数量
	updateParticipantsCountOp := ChannelParticipantsCountDao{}.incrOp(channel.ChannelID, len(users))
	if updateParticipantsCountOp != nil {
		ops = append(ops, updateParticipantsCountOp)
	}

	// 执行数据库操作
	err = batch.Add(ops...).Run()
	if err != nil {
		return err
	}

	// 更新频道在线用户列表
	devices, err := DevicesManager{}.GetOnlineDevices(ids)
	if err == nil {
		ChannelParticipantsManager{}.join(channel.ChannelID, devices)
	}
	return nil
}

// 获取频道
func (Channels) Get(channelID string) (Channel, error) {
	var record model.Channel
	op := internal.cassa.Model(model.Channel{}).Where(
		gocassa.Eq("channel_id", channelID)).ReadOne(&record)
	if err := op.Run(); err != nil {
		return Channel{}, err
	}

	channel := Channel{
		ChannelID:   record.ChannelID.String(),
		Creator:     record.Creator.String(),
		Public:      record.Public,
		Deactivated: record.Deactivated,
		Date:        record.Date.Unix(),
	}
	json.Unmarshal(record.Profile, &channel.Profile)
	return channel, nil
}

// 获取频道列表
func (Channels) GetList(channelIDs []string) ([]Channel, error) {
	tableName := internal.cassa.TableName(model.Channel{})
	cql := fmt.Sprintf("SELECT * FROM %s WHERE channel_id IN (%s)",
		tableName, strings.Join(channelIDs, ","))

	iter := internal.cassa.Query(cql).Iter()
	rowData, err := iter.RowData()
	if err != nil {
		return nil, err
	}

	var records []model.Channel
	scanner, err := cassa.NewScanner(rowData, model.Channel{})
	if err != nil {
		return nil, err
	}
	if err = scanner.ScanRows(iter, rowData, &records); err != nil {
		return nil, err
	}

	channels := make([]Channel, 0, len(records))
	for _, record := range records {
		channel := Channel{
			ChannelID:   record.ChannelID.String(),
			Creator:     record.Creator.String(),
			Public:      record.Public,
			Deactivated: record.Deactivated,
			Date:        record.Date.Unix(),
		}
		json.Unmarshal(record.Profile, &channel.Profile)
		channels = append(channels, channel)
	}
	return channels, nil
}

// 获取频道资料
func (Channels) GetProfile(channelID string) (ChannelProfile, error) {
	var result struct {
		Profile []byte
	}
	op := internal.cassa.Model(model.Channel{}).Where(gocassa.Eq("channel_id", channelID)).ReadOne(&result)
	err := op.WithOptions(gocassa.Options{Select: []string{"profile"}}).Run()
	if err != nil {
		return ChannelProfile{}, err
	}

	var profile ChannelProfile
	err = json.Unmarshal(result.Profile, &profile)
	if err != nil {
		return ChannelProfile{}, err
	}
	return profile, err
}

// 更新频道资料
func (Channels) UpdateProfile(channelID string, profile ChannelProfile) error {
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{"profile": data}
	op := internal.cassa.Model(model.Channel{}).Where(gocassa.Eq("channel_id", channelID)).Update(updates)
	return op.Run()
}

// 添加成员
func (Channels) AddParticipants(channelID string, nameMapper map[string]string) error {
	if len(nameMapper) == 0 {
		return nil
	}

	// 插入频道成员
	var batchOp gocassa.Op
	userIDs := make([]string, 0, len(nameMapper))
	for userID, name := range nameMapper {
		// 更新频道成员
		participant := ChannelParticipant{
			ParticipantID: ParticipantID{
				UserID:    userID,
				ChannelID: channelID,
			},
			Name: name,
			Date: time.Now().Unix(),
		}
		updateOp := ChannelParticipants{}.updateOp(participant)

		// 更新已加入频道
		joined := ChannelJoined{
			UserID:    userID,
			ChannelID: channelID,
		}
		updateJoinedOp, err := ChannelJoinedDao{}.batchUpdateOp([]ChannelJoined{joined})
		if err != nil {
			return err
		}

		if batchOp == nil {
			batchOp = updateOp.Add(updateJoinedOp)
		} else {
			batchOp = batchOp.Add(updateOp, updateJoinedOp)
		}
		userIDs = append(userIDs, userID)
	}

	// 更新频道成员数量
	updateParticipantsCountOp := ChannelParticipantsCountDao{}.incrOp(channelID, len(nameMapper))

	// 执行数据库操作
	if err := batchOp.Add(updateParticipantsCountOp).Run(); err != nil {
		return err
	}

	// 更新频道在线用户列表
	devices, err := DevicesManager{}.GetOnlineDevices(userIDs)
	if err == nil {
		ChannelParticipantsManager{}.join(channelID, devices)
	} else {
		logx.Error("[Storage] failed to add channel participants to cache",
			" channel:", channelID,
			" users:",   userIDs,
			" reason:",  err,
			)

	}
	return nil
}

// 移除成员
func (Channels) RemoveParticipant(channelID, userID string, kicked bool) error {
	// 删除在线缓存
	err := ChannelParticipantsManager{}.leave(channelID, userID)
	if err != nil {
		logx.Error("[Storage] failed to remove channel participant in cache",
			" channel:", channelID,
			" users:",   userID,
			" reason:",  err,
		)

	}

	// 更新频道成员
	batchOp := ChannelParticipants{}.updateOp(ChannelParticipant{
		ParticipantID: ParticipantID{
			UserID:    userID,
			ChannelID: channelID,
		},
		Left:   !kicked,
		Kicked: kicked,
		Date:   time.Now().Unix(),
	})

	// 更新已加入频道
	joined := ChannelJoined{
		UserID:    userID,
		ChannelID: channelID,
		Left:      !kicked,
		Kicked:    kicked,
	}
	updateJoinedOp, err := ChannelJoinedDao{}.batchUpdateOp([]ChannelJoined{joined})
	if err != nil {
		return err
	}

	// 更新频道成员数量
	updateParticipantsCountOp := ChannelParticipantsCountDao{}.incrOp(channelID, -1)
	return batchOp.Add(updateJoinedOp, updateParticipantsCountOp).Run()
}

// 停用频道
func (Channels) Deactivate(channelID string) error {
	err := ChannelParticipantsManager{}.deactivate(channelID)
	if err != nil {

		logx.Error("[Storage] failed to remove channel participants in cache",
			" channel:", channelID,
			" reason:",  err,
		)

	}

	updates := map[string]interface{}{"deactivated": true}
	op := internal.cassa.Model(model.Channel{}).Where(gocassa.Eq("channel_id", channelID)).Update(updates)
	return op.Run()
}

