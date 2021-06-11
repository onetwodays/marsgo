// Code generated by goctl. DO NOT EDIT.
package types

type IndexReply struct {
	Resp string `json:"resp"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserReply struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	JwtToken
}

type JwtToken struct {
	AccessToken  string `json:"accessToken,omitempty"`
	AccessExpire int64  `json:"accessExpire,omitempty"`
	RefreshAfter int64  `json:"refreshAfter,omitempty"`
}

type AddReq struct {
	Book  string `form:"book"`
	Price int64  `form:"price"`
}

type AddResp struct {
	Ok bool `json:"ok"`
}

type CheckReq struct {
	Book string `form:"book"`
}

type CheckResp struct {
	Found bool  `json:"found"`
	Price int64 `json:"price"`
}

type JwtTokenAdx struct {
	AccessToken  string `json:"accessToken,omitempty"`
	AccessExpire int64  `json:"accessExpire,omitempty"`
	RefreshAfter int64  `json:"refreshAfter,omitempty"`
}

type AdxUserLoginReq struct {
	Name string `json:"name"` //eos chain username,保证unique
	Sign string `json:"sign"` //  eos 用户用自己的私钥对name的签名
}

type AdxUserLoginRes struct {
	JwtTokenAdx
}

type IncomingMessagex struct {
	Content                   string `json:"content"`
	Type                      int    `json:"type"`
	DestinationDeviceId       int    `json:"destinationDeviceId,default=1"` //发到哪一个设备
	DestinationRegistrationId int    `json:"destinationRegistrationId"`
	Destination               string `json:"destination,optional"`
	Body                      string `json:"body,optional"`
	Relay                     string `json:"relay,optional"`
}

type PutMessagesReq struct {
	Destination string             `path:"destination"`
	Online      bool               `json:"online"`
	Timestamp   int64              `json:"timestamp"`
	Messages    []IncomingMessagex `json:"messages"`
}

type PutMessagesRes struct {
	NeedsSync   bool     `json:"needsSync"`
	DestContent [][]byte `json:"destContent,optional"`
}

type OutcomingMessagex struct {
	Id              int64  `json:"id"`
	Cached          bool   `json:"cached"`
	Guid            string `json:"guid"`
	Type            int    `json:"type"`
	Relay           string `json:"relay"`
	Timestamp       int64  `json:"timestamp"`
	Source          string `json:"source"`
	SourceUuid      string `json:"sourceUuid"`
	SourceDevice    int64  `json:"sourceDevice"`
	Message         string `json:"message"`
	Content         string `json:"content"`
	ServerTimestamp int64  `json:"serverTimestamp"`
}

type GetPendingMsgsReq struct {
}

type GetPendingMsgsRes struct {
	List []OutcomingMessagex `json:"list"`
	More bool                `json:"more"`
}

type Envelope struct {
	Xtype           int    `json:"type"`
	Source          string `json:"source"`
	SourceUuid      string `json:"sourceUuid"`
	SourceDevice    int    `json:"sourceDevice"`
	Relay           string `json:"relay"`
	Timestamp       uint64 `json:"timestamp"`
	LegacyMessage   string `json:"legacyMessage"`
	Content         string `json:"content"`
	ServerGuid      string `json:"guid"`
	ServerTimestamp uint64 `json:"serverTimestamp"`
}

type PubsubMessage struct {
	Xtype   int      `json:"type"`
	Content Envelope `json:"envelop"`
}

type PreKeyx struct {
	KeyId     int64  `json:"keyId"`
	PublicKey string `json:"publickey"`
}

type SignedPrekeyx struct {
	Signature string  `json:"signature"`
	PreKey    PreKeyx `json:"prekey"`
}

type PutKeysReqx struct {
	IdentityKey  string        `json:"identityKey"`
	SignedPreKey SignedPrekeyx `json:"signedPreKey"`
	PreKeys      []PreKeyx     `json:"prekeys"`
}

type GetKeysReq struct {
	Identifier string `path:"identifier"`
	DeviceId   int64  `path:"deviceId"`
}

type PreKeyResponseItem struct {
	DeviceId       int64         `json:"deviceId"`
	RegistrationId int64         `json:"registrationId"`
	PreKey         PreKeyx       `json:"preKey"`
	SignedPrekey   SignedPrekeyx `json:"signedPreKey"`
}

type GetKeysResx struct {
	IdentityKey string               `json:"identityKey"`
	Devices     []PreKeyResponseItem `json:"devices"`
}

type WriteWsConnReq struct {
	Login    string `form:"login"`
	Password string `form:"password"`
}
