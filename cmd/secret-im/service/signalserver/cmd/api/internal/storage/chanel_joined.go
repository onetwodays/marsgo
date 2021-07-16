package storage

// 已加入频道
type ChannelJoined struct {
	UserID    string `json:"userId"`
	ChannelID string `json:"channelId"`
	Left      bool   `json:"left"`
	Kicked    bool   `json:"kicked"` // 被踢
}

