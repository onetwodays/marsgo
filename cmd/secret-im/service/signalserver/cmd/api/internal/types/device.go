package types


type Device struct {
	Id              int64        `gorm:"primary_key" json:"id"`
	Salt            string       `gorm:"column:salt" json:"salt,omitempty"`
	AuthToken       string       `gorm:"column:authToken" json:"authToken,omitempty"`
	GcmId           string       `gorm:"column:gcmId" json:"gcmId,omitempty"`
	ApnId           string       `gorm:"column:apnId" json:"apnId,omitempty"`
	SignedPreKey   SignedPreKey `gorm:"column:signedPreKey" json:"signedPreKey"`
	LastSeen       int64        `gorm:"column:lastSeen" json:"lastSeen,omitempty"`
	VoipApnId      string       `gorm:"column:voipApnId" json:"voipApnId,omitempty"`
	UerAgent       string       `gorm:"column:userAgent" json:"userAgent,omitempty"`
	Created        int64        `gorm:"column:created" json:"created,omitempty"`
	AccountAttributes AccountAttributes `json:"accountAttributes"`
}



type Account struct {
	Number       string   `gorm:"column:number" json:"number"`
	Devices      []Device `gorm:"type:json;column:devices" json:"devices"`
	IdentityKey  string   `gorm:"column:identityKey" json:"identityKey"`
	Name         string   `gorm:"column:name" json:"name"`
	Avatar       string   `gorm:"column:avatar" json:"avatar"`
	AvatarDigest string   `gorm:"column:avatarDigest" json:"avatarDigest"`
	Pin          string   `gorm:"column:pin" json:"pin"`
}
