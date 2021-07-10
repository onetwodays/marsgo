package storage

import (
	"fmt"
	"secret-im/service/signalserver/cmd/api/internal/model"

	"github.com/go-redis/redis"
)

// 用户名管理器
type UsernamesManager struct {
}

// 放入用户名
func (m UsernamesManager) Put(uuid, username string) (bool, error) {
	/*
	ok, err := ReservedUsernames{}.IsReserved(username, uuid)
	if ok {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	 */
	record :=&model.TUsernames{
		Uuid: uuid,
		Username: username,
	}
	_,err:=internal.userNameDB.Insert(*record)

	if err != nil {
		return false, err
	}

	m.redisSet(uuid, username)
	return true, nil
}

// 根据UUID获取用户名
func (m UsernamesManager) GetUsernameForUUID(uuid string) (string, error) {
	username := m.redisGetUsernameForUUID(uuid)
	if username != nil {
		return *username, nil
	}

	recode,err:=internal.userNameDB.FindOneByUuid(uuid)
	if err!=nil{
		return "", err
	}
	m.redisSet(uuid, recode.Username)
	return recode.Username, nil
}

// 根据用户名获取UUID
func (m UsernamesManager) GetUUIDForUsername(username string) (string, error) {
	uuid := m.redisGetUUIDForUsername(username)
	if uuid != nil {
		return *uuid, nil
	}
	record,err:=internal.userNameDB.FindOneByUsername(username)
	if err!=nil{
		return "", err
	}


	m.redisSet(record.Uuid, username)
	return record.Uuid, nil
}

// 删除用户名
func (m UsernamesManager) Delete(uuid string) error {
	if err := m.redisDelete(uuid); err != nil {
		return err
	}

	return internal.userNameDB.DeleteByUuid(uuid)
}

// UUID->用户名
func (UsernamesManager) uuidMapKey(uuid string) string {
	return fmt.Sprintf("username_by_uuid::{%s}", uuid)
}

// 用户名->UUID
func (UsernamesManager) usernameMapKey(username string) string {
	return fmt.Sprintf("uuid_by_username::{%s}", username)
}

// 设置缓存
func (m UsernamesManager) redisSet(uuid string, username string) error {
	cmd := internal.client.Get(m.uuidMapKey(uuid))
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return cmd.Err()
	}
	if cmd.Err() == nil {
		username := cmd.Val()
		if err := internal.client.Del(m.usernameMapKey(username)).Err(); err != nil {
			return err
		}
	}

	err := internal.client.Set(m.uuidMapKey(uuid), username, 0).Err()
	if err != nil {
		return err
	}
	return internal.client.Set(m.usernameMapKey(username), uuid, 0).Err()
}

// 根据UUID获取缓存
func (m UsernamesManager) redisGetUsernameForUUID(uuid string) *string {
	cmd := internal.client.Get(m.uuidMapKey(uuid))
	if cmd.Err() != nil {
		return nil
	}
	username := cmd.Val()
	return &username
}

// 根据用户名获取缓存
func (m UsernamesManager) redisGetUUIDForUsername(username string) *string {
	cmd := internal.client.Get(m.usernameMapKey(username))
	if cmd.Err() != nil {
		return nil
	}
	uuid := cmd.Val()
	return &uuid
}

// 删除缓存
func (m UsernamesManager) redisDelete(uuid string) error {
	username := m.redisGetUsernameForUUID(uuid)
	if username == nil {
		return nil
	}
	cmd := internal.client.Del(m.uuidMapKey(uuid))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return internal.client.Del(m.usernameMapKey(*username)).Err()
}
