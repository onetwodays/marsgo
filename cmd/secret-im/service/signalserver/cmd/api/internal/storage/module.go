package storage

import (
	"github.com/go-redis/redis"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/storage/message/operation"

)

// 定义了一个全局变量
var internal struct {
	accountDB model.TAccountsModel
	msgDB     model.TMessagesModel
	client          *redis.Client
	// 查询消息操作
	getOperation *operation.GetOperation
	// 插入消息操作
	insertOperation *operation.InsertOperation
	// 删除消息操作
	removeOperation *operation.RemoveOperation
}
func InitStorage(client *redis.Client,accountDB model.TAccountsModel,msgDB model.TMessagesModel)  {
	syncx.Once(func() {
		internal.accountDB =accountDB
		internal.msgDB=msgDB
		internal.client=client
		var err error
		internal.getOperation,err = operation.NewGetOperation(client)
		internal.insertOperation,err = operation.NewInsertOperation(client)
		internal.removeOperation,err = operation.NewRemoveOperation(client)
		if err!=nil{
			logx.Error("[Storage] failed to init storage module,reason:",err)
			panic(err)
		}


	})()
}