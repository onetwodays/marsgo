package logic

import (
	"github.com/golang/protobuf/proto"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/textsecure"
)


// 创建频道信息
func newChannelEntity(channel storage.Channel, participant, left, kicked bool) *types.Channel {
	return &types.Channel{
		ID:            channel.ChannelID,
		Title:         channel.Profile.Title,
		Photo:         channel.Profile.Photo,
		About:         channel.Profile.About,
		Creator:       channel.Creator,
		Public:        channel.Public,
		IsParticipant: participant,
		Left:          left,
		Kicked:        kicked,
		Deactivated:   channel.Deactivated,
		Date:          channel.Date,
	}
}

// 填充频道消息属性
func fillChannelMessageAttributes(channelSlice []types.Channel, userID ...string) error {
	if len(channelSlice) == 0 {
		return nil
	}

	channelIDs := make([]string, 0)
	for _, channel := range channelSlice {
		channelIDs = append(channelIDs, channel.ID)
	}

	var err error
	ackMessageMapper := make(map[string]int64, 0)
	if len(userID) > 0 {
		ackMessageMapper, err = storage.ChannelMessageAckDao{}.GetMessages(userID[0], channelIDs)
		if err != nil {
			return err
		}
	}

	messageMapper, err := storage.ChannelMessagesManager{}.GetLatestMessages(channelIDs)
	if err != nil {
		return err
	}

	for idx, channel := range channelSlice {
		ackMessage, ok := ackMessageMapper[channel.ID]
		if ok {
			channelSlice[idx].LastAckMessage = ackMessage
		}

		latestMessage, ok := messageMapper[channel.ID]
		if ok {
			message := newOutgoingChannelMessage(&latestMessage)
			channelSlice[idx].LatestMessage = &message
			channelSlice[idx].Unread = int(latestMessage.MessageID - ackMessage)
		}
	}
	return nil
}



func newOutgoingChannelMessage(message *storage.ChannelMessage) types.OutgoingChannelMessage {
	entity := types.OutgoingChannelMessage{
		ID:              message.MessageID,
		ChannelID:       message.ChannelID,
		Type:            message.Type.String(),
		Source:          *message.Source,
		SourceDeviceID:  *message.SourceDevice,
		Relay:           *message.Relay,
		Deleted:         message.Deleted,
		Timestamp:       message.Timestamp,
		ServerTimestamp: message.ServerTimestamp,
	}

	if entity.Deleted {
		return entity
	}

	if message.Editor != nil {
		entity.Editor = message.Editor.UUID
		entity.EditedAt = message.Editor.EditedAt
	}

	if message.Type == model.ChannelMessageTypeNormal {
		if len(message.Content) > 0 {
			content := string(message.Content)
			entity.Content = content
		}
		return entity
	}

	if message.Type == model.ChannelMessageTypeService {
		if len(message.Action) > 0 {
			var action textsecure.MessageAction
			if err := proto.Unmarshal(message.Action, &action); err != nil {
				return entity
			}
			entity.Action = &types.ChannelMessageAction{
				Action:       action.GetAction().String(),
				Title:        action.GetTitle(),
				Photo:        action.GetPhoto(),
				Participants: action.GetParticipants(),
				MessageID:    action.GetMessageId(),
				Operator:     action.GetOperator(),
			}
		}
		return entity
	}
	return entity
}

