package jobs

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/queue"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/api/websocket"
	"sync/atomic"
	"time"
)
// 消息投递作业
type MessageDeliverJob struct {
	running        int32
	subscription   *nats.Subscription
	sessionManager *websocket.SessionManager
}

// 创建消息投递作业
func NewMessageDeliverJob(nc *nats.Conn, sessionManager *websocket.SessionManager) (*MessageDeliverJob, error) {
	job := MessageDeliverJob{
		sessionManager: sessionManager,
	}

	var err error
	partition := config.AppConfig.PartitionID
	job.subscription, err = nc.QueueSubscribe(
		queue.SendToDeviceTopic(partition), "workers", job.handleMessageDeliver)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// 开始作业
func (job *MessageDeliverJob) Start() error {
	if !atomic.CompareAndSwapInt32(&job.running, 0, 1) {
		return nil
	}
	logx.Info("[MessageDeliverJob] job started")
	return nil
}

// 停止作业
func (job *MessageDeliverJob) Stop() {
	if !atomic.CompareAndSwapInt32(&job.running, 1, 0) {
		return
	}
	job.subscription.Unsubscribe()
	logx.Info("[MessageDeliverJob] job stopped")
}

// 处理投递消息
func (job *MessageDeliverJob) handleMessageDeliver(message *nats.Msg) {
	var request queue.SendMessageToDevice
	err := proto.Unmarshal(message.Data, &request)
	if err != nil {
		return
	}

	t := time.Now()
	count, err := job.sendMessageToDevice(&request)

	if err != nil {
		logx.Error("[MessageDeliverJob] failed to send message to devices"," reason:",err)
		return
	}
	logx.Info("[MessageDeliverJob] send message to devices successful"," consume:", time.Since(t),
		" count:",  count)


}

// 发送消息到设备
func (job *MessageDeliverJob) sendMessageToDevice(request *queue.SendMessageToDevice) (int, error) {
	// 解析消息
	var message interface{}
	switch request.GetType() {
	case queue.SendMessageToDevice_CHANNEL_MESSAGE:
		var envelope textsecure.ChannelEnvelope
		if err := proto.Unmarshal(request.GetContent(), &envelope); err != nil {
			return 0, err
		}
		message = envelope
	default:
		return 0, errors.New("unknown message type")
	}

	// 发送在线设备
	count := 0
	for _, device := range request.GetDevices() {
		session, ok := job.sessionManager.GetByDevice(device.GetUuid(), device.GetId())
		if !ok {
			continue
		}
		session.DeliverMessage(message)
		count++
	}
	return count, nil
}
