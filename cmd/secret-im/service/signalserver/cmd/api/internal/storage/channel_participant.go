package storage

import "secret-im/service/signalserver/cmd/api/internal/entities"

// 成员ID
type ParticipantID struct {
	UserID    string `json:"userId"`
	ChannelID string `json:"channelId"`
}

// 管理权限
type ChannelAdminRights int

const (
	ChannelAdminRightChangeInfo ChannelAdminRights = 1 << iota
	ChannelAdminRightEditMessages
	ChannelAdminRightDeleteMessages
	ChannelAdminRightBanUsers
	ChannelAdminRightInviteUsers
	ChannelAdminRightPinMessages
	ChannelAdminRightAddAdmins
)

func (rights ChannelAdminRights) ToEntity() entities.ChannelAdminRights {
	var entity entities.ChannelAdminRights
	entity.ChangeInfo = rights&ChannelAdminRightChangeInfo > 0
	entity.EditMessages = rights&ChannelAdminRightEditMessages > 0
	entity.DeleteMessages = rights&ChannelAdminRightDeleteMessages > 0
	entity.BanUsers = rights&ChannelAdminRightBanUsers > 0
	entity.InviteUsers = rights&ChannelAdminRightInviteUsers > 0
	entity.PinMessages = rights&ChannelAdminRightPinMessages > 0
	entity.AddAdmins = rights&ChannelAdminRightAddAdmins > 0
	return entity
}

func NewChannelAdminRightsFromEntity(rights entities.ChannelAdminRights) ChannelAdminRights {
	var right ChannelAdminRights
	if rights.ChangeInfo {
		right |= ChannelAdminRightChangeInfo
	}
	if rights.EditMessages {
		right |= ChannelAdminRightEditMessages
	}
	if rights.DeleteMessages {
		right |= ChannelAdminRightDeleteMessages
	}
	if rights.BanUsers {
		right |= ChannelAdminRightBanUsers
	}
	if rights.InviteUsers {
		right |= ChannelAdminRightInviteUsers
	}
	if rights.PinMessages {
		right |= ChannelAdminRightPinMessages
	}
	if rights.AddAdmins {
		right |= ChannelAdminRightAddAdmins
	}
	return right
}

// 默认管理员权限
func DefaultChannelAdminRights() ChannelAdminRights {
	return ChannelAdminRightChangeInfo | ChannelAdminRightEditMessages |
		ChannelAdminRightDeleteMessages | ChannelAdminRightBanUsers | ChannelAdminRightInviteUsers |
		ChannelAdminRightPinMessages | ChannelAdminRightAddAdmins
}

// 频道成员
type ChannelParticipant struct {
	ParticipantID
	Name                  string                `json:"name"`
	Left                  bool                  `json:"left"`
	Kicked                bool                  `json:"kicked"`
	Banned                bool                  `json:"banned"`
	AdminRights           ChannelAdminRights    `json:"adminRights,omitempty"`
	ChannelNotifySettings ChannelNotifySettings `json:"notifySettings"`
	Date                  int64                 `json:"date"`
}

// 频道通知设置
type ChannelNotifySettings struct {
	Silent bool `json:"silent"`
}

