package entities


// 频道资料
type ChannelProfile struct {
	Title string  `json:"title"`
	Photo string `json:"photo,omitempty"`
	About string  `json:"about,omitempty"`
}

// 频道信息
type Channel struct {
	ChannelID   string         `json:"channelId"`
	Creator     string         `json:"creator"`
	Profile     ChannelProfile `json:"profile"`
	Public      bool           `json:"public"`
	Deactivated bool           `json:"deactivated"`
	Date        int64          `json:"timestamp"`
}

