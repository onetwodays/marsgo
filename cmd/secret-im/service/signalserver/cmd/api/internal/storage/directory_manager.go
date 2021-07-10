package storage

import (
	"encoding/json"
	"secret-im/service/signalserver/cmd/api/internal/types"
)

// 令牌值
type tokenValue struct {
	Number   string `json:"n"`
	Relay    string `json:"r"`
	Voice    bool   `json:"v"`
	Video    bool   `json:"w"`
	Inactive bool   `json:"a"`
}

// 联系人管理器
type DirectoryManager struct {
}

// 获取缓存键
func (DirectoryManager) directoryKey() string {
	return "directory"
}

// 获取号码
func (manager DirectoryManager) Get(tokens []string) ([]types.ClientContact, error) {
	cmd := internal.client.HMGet(manager.directoryKey(), tokens...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	contacts := make([]types.ClientContact, 0)
	for idx, item := range cmd.Val() {
		if item == nil {
			continue
		}

		js := item.(string)
		var value tokenValue
		if err := json.Unmarshal([]byte(js), &value); err != nil {
			return nil, err
		}
		contacts = append(contacts, types.ClientContact{
			Token:    tokens[idx],
			Voice:    value.Voice,
			Video:    value.Video,
			Relay:    value.Relay,
			Inactive: value.Inactive,
		})
	}
	return contacts, nil
}

// 添加号码
func (manager DirectoryManager) Add(number string, contact types.ClientContact) error {
	value := tokenValue{
		Number:   number,
		Relay:    contact.Relay,
		Voice:    contact.Voice,
		Video:    contact.Video,
		Inactive: contact.Inactive,
	}
	jsb, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return internal.client.HSet(manager.directoryKey(), contact.Token, string(jsb)).Err()
}

// 删除号码
func (manager DirectoryManager) Remove(token string) error {
	return internal.client.HDel(manager.directoryKey(), token).Err()
}

