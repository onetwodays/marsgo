package websocket

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/timex"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/api/websocket/factory"
	"sync"
)

//请求客户端
type Clientx struct {
	session    *Session
	requestMap sync.Map
}

// 处理响应消息
func (client *Clientx) handleResponse(message *textsecure.WebSocketResponseMessage) {
	startTime:=timex.Now()
	id := int64(message.GetId())
	value, ok := client.requestMap.Load(id)
	if !ok {
		return
	}
	future := value.(*Future)
	future.setResult(message, nil)
	client.requestMap.Delete(id)

	logx.WithDuration(timex.Since(startTime)).Infof("完成 [%s] wsres处理 (%s)]",
		client.session.sessionName,
		message.String())
    /*
	if client.session.context.Device != nil {
		logx.Infof("[Client] handle websocket response ok, from %s, recv peer ws response msg:[%s] ",
			client.session.context.Device.Number,
			message.String())
	}
     */
}

// 发送请求
func (client *Clientx) SendRequest(verb, path string, headers []string, body []byte) *Future {
	future := newFuture()
	if client.session.IsClosed(){
		err:=errors.New("closed")
		future.setResult(nil,err)
		return future
	}
	requestID := utils.SecureRandInt64()
	client.requestMap.Store(requestID, future)

	requestMessage := factory.CreateRequest(uint64(requestID), verb, path, headers, body)

	//logx.Info("2.[Client]服务器发出的websocket request :",requestMessage.String())
	data, err := proto.Marshal(requestMessage)
	if err != nil {
		client.requestMap.Delete(requestID)
		future.setResult(nil, err)
		return future
	}

	if err = client.session.Send(data); err != nil {
		client.requestMap.Delete(requestID)
		future.setResult(nil, err)
	}

	logx.Infof("[%s] 接收发给他的消息完成 ]", client.session.sessionName)

	/*
	if client.session.context.Device != nil {


			logx.Info("[Client] send request to ",client.session.context.Device.Number,
				" id:",   requestID,
				" verb:", verb,
				" path:", path)


	}
	*/

	return future
}

func (client *Clientx) CancelAll() {
	err := errors.New("canceled")
	client.requestMap.Range(func(key, value interface{}) bool {
		future := value.(*Future)
		future.setResult(nil, err)
		client.requestMap.Delete(key)
		return true
	})
}
