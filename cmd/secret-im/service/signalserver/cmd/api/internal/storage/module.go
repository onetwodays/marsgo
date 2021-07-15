package storage

import (
	"github.com/go-redis/redis"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
	"secret-im/pkg/driver/cassa"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/storage/message/operation"

)

// 定义了一个全局变量
var internal struct {
	accountDB model.TAccountsModel
	msgDB     model.TMessagesModel
	profileDB model.TProfilesModel
	userNameDB model.TUsernamesModel
	client          *redis.Client
	cassa       *cassa.Conn
	// 查询消息操作
	getOperation *operation.GetOperation
	// 插入消息操作
	insertOperation *operation.InsertOperation
	// 删除消息操作
	removeOperation *operation.RemoveOperation
}

// 注册cassandra表
var channelTables = []cassa.Table{
	new(model.Channel),
	new(model.ChannelJoined),
	new(model.ChannelMessage),
	new(model.ChannelMessageAck),
	new(model.ChannelParticipant),
	new(model.ChannelParticipantsCount),
}
func InitStorage(client *redis.Client,accountDB model.TAccountsModel,
	                                  msgDB model.TMessagesModel,
	                                  profileDB model.TProfilesModel,
	                                  userNameDB model.TUsernamesModel,
                                      cassandraClient *cassa.Conn)  {
	syncx.Once(func() {
		internal.accountDB =accountDB
		internal.msgDB=msgDB
		internal.profileDB=profileDB
		internal.userNameDB=userNameDB
		internal.client=client
		internal.cassa=cassandraClient
		var err error
		internal.getOperation,err = operation.NewGetOperation(client)
		internal.insertOperation,err = operation.NewInsertOperation(client)
		internal.removeOperation,err = operation.NewRemoveOperation(client)
		if err!=nil{
			logx.Error("[Storage] failed to init storage module,reason:",err)
			panic(err)
		}
		if internal.cassa!=nil{
			err:=internal.cassa.Ensure(channelTables)
			if err!=nil{
				logx.Error("[Storage] failed to initialize cassandra table,reason:",err)
				panic(err)
			}

		}




	})()
}