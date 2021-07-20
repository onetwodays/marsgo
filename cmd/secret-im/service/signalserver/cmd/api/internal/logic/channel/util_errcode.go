package logic

import (
	"fmt"
)

// 错误信息
type Error struct {
	Code        int         `json:"code" example:"1001"`
	Description string      `json:"description" example:"account not found"`
	Data        interface{} `json:"data"`
}

func (e Error) String() string {
	return fmt.Sprintf("code: %d, description: %s", e.Code, e.Description)
}

// 账号不存在
func ErrAccountNotFound(uuid string) Error {
	return Error{
		Code:        1001,
		Description: "account not found",
		Data: map[string]string{
			"uuid": uuid,
		},
	}
}

// 频道不存在
func ErrChannelNotFound(channelID string) Error {
	return Error{
		Code:        1002,
		Description: "channel not found",
		Data: map[string]string{
			"channel_id": channelID,
		},
	}
}

// 非频道成员
func ErrNotChannelParticipant(channelID, userID string) Error {
	return Error{
		Code:        1003,
		Description: "not a channel participant",
		Data: map[string]string{
			"channel_id": channelID,
			"user_id":    userID,
		},
	}
}

// 没有操作权限
func ErrNoOperationPermission(channelID, userID string) Error {
	return Error{
		Code:        1004,
		Description: "no operation permission",
		Data: map[string]string{
			"channel_id": channelID,
			"user_id":    userID,
		},
	}
}

// 不能修改服务消息
func ErrCannotModifyServiceMessage(channelID string, messageID int64) Error {
	return Error{
		Code:        1005,
		Description: "cannot modify service message",
		Data: map[string]interface{}{
			"channel_id": channelID,
			"message_id": messageID,
		},
	}
}

// 用户被禁止
func ErrUserIsBanned(uuid string) Error {
	return Error{
		Code:        1006,
		Description: "user is banned",
		Data: map[string]interface{}{
			"user_id": uuid,
		},
	}
}

// 私有频道
func ErrPrivateChannel(channelID string) Error {
	return Error{
		Code:        1007,
		Description: "channel is private",
		Data: map[string]interface{}{
			"channel_id": channelID,
		},
	}
}

