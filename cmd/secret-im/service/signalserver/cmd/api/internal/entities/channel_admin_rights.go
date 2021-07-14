package entities
// 频道管理权限
type ChannelAdminRights struct {
	ChangeInfo     bool `json:"changeInfo,omitempty"`
	EditMessages   bool `json:"editMessages,omitempty"`
	DeleteMessages bool `json:"deleteMessages,omitempty"`
	BanUsers       bool `json:"banUsers,omitempty"`
	InviteUsers    bool `json:"inviteUsers,omitempty"`
	PinMessages    bool `json:"pinMessages,omitempty"`
	AddAdmins      bool `json:"addAdmins,omitempty"`
}

