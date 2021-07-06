package jobs

import (
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"io"
	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/storage/message/operation"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"sync/atomic"
	"time"
)

// 消息持久化作业
type MessagePersistJob struct {
	running       int32
	stopChan      chan struct{}
	pushSender    *push.Sender
	pubSubManager *pubsub.Manager
}

// 创建消息持久化作业
func NewMessagePersistJob(pushSender *push.Sender, pubSubManager *pubsub.Manager) *MessagePersistJob {
	return &MessagePersistJob{
		pushSender:    pushSender,
		pubSubManager: pubSubManager,
	}
}

// 开始作业
func (job *MessagePersistJob) Start() error {
	if !atomic.CompareAndSwapInt32(&job.running, 0, 1) {
		return nil
	}

	job.stopChan = make(chan struct{})
	go job.run()
	logx.Info("[MessagePersistJob] job started")
	return nil
}

// 停止作业
func (job *MessagePersistJob) Stop() {
	if !atomic.CompareAndSwapInt32(&job.running, 1, 0) {
		return
	}
	<-job.stopChan
	logx.Info("[MessagePersistJob] job stopped")
}

// 运行服务
func (job *MessagePersistJob) run() {
	defer func() {
		atomic.StoreInt32(&job.running, 0)
		close(job.stopChan)
	}()

	delayTime := config.AppConfig.MessageCache.PersistDelayMinutes
	delayTimeMillis := (delayTime * int64(time.Minute)) / int64(time.Millisecond)

	for {
		if atomic.LoadInt32(&job.running) == 0 {
			break
		}

		queuesToPersist, err := storage.MessagesCache{}.GetQueuesToPersist(delayTimeMillis)
		if err != nil {
			if err != io.EOF {
				logx.Error("[MessagePersistJob] failed to get queues to persist,reason:", err)
			}
			continue
		}

		for _, queue := range queuesToPersist {
			key, err := operation.Key{}.FromUserMessageQueue(queue)
			if err != nil {
				logx.Error("[MessagePersistJob] failed to get key",
					" queue:", queue,
					" reason:", err)
				continue
			}

			messagesPersistedCount, err := job.persistQueue(key)
			if err != nil {
				logx.Error("[MessagePersistJob] failed to persist queue",
					" queue:", queue,
					" reason:", err)

				continue
			}

			job.notifyClients(key)
			logx.Info("[MessagePersistJob] messages persisted", " queue:", queue, " count:", messagesPersistedCount)

		}

		if len(queuesToPersist) == 0 {
			if atomic.LoadInt32(&job.running) == 0 {
				break
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}
}

// 通知客户端
func (job *MessagePersistJob) notifyClients(key operation.Key) {

	message := &textsecure.PubSubMessage{Type: textsecure.PubSubMessage_QUERY_DB}
	address := push.Address{Number: key.Address, DeviceID: key.DeviceID}
	n, err := job.pubSubManager.Publish(address.Serialize(), message)
	logx.Info("[MessagePersistJob] notify clients:",
		" count:", n,
		" reason:", err)

	if n > 0 {
		return
	}

	account, err := storage.AccountManager{}.GetByNumber(key.Address)
	if err != nil {
		logx.Error("[MessagePersistJob] number not found",
			" number:", key.Address,
			" reason:", err)
		return
	}

	device, ok := account.GetDevice(key.DeviceID)
	if !ok {
		return
	}

	err = job.pushSender.SendQueuedNotification(account.Number, &device.Device)
	if err != nil {
		logx.Error("[MessagePersistJob] failed to send queued notification",
			" number:", key.Address,
			" device_id:", key.DeviceID,
			" reason:", err)
		return
	}
}

// 持久化队列
func (job *MessagePersistJob) persistQueue(key operation.Key) (int, error) {
	messagesPersistedCount := 0
	for {
		err := storage.MessagesCache{}.QueueLock(key.Address, key.DeviceID)
		if err != nil {
			logx.Info("[MessagePersistJob] failed to lock queue", " number:", key.Address, " device_id:", key.DeviceID, " reason:", err)
			return 0, err
		}

		messages, err := storage.MessagesCache{}.GetFromQueue(key.Address, key.DeviceID)
		if err != nil {

			logx.Info("[MessagePersistJob] failed to get message from queue", " number:", key.Address, " device_id:", key.DeviceID, " reason:", err)
			return 0, err
		}

		for _, message := range messages {
			err = job.persistMessage(key, message.ID, message.Data)
			if err != nil {
				return 0, err
			}
			messagesPersistedCount++
		}

		if len(messages) < storage.ResultSetChunkSize {
			storage.MessagesCache{}.QueueUnlock(key.Address, key.DeviceID)
			return messagesPersistedCount, nil
		}
	}
}

// 持久化消息
func (job *MessagePersistJob) persistMessage(key operation.Key, id int64, message []byte) error {
	var envelope textsecure.Envelope
	err := proto.Unmarshal(message, &envelope)
	if err != nil {
		logx.Error("[MessagePersistJob] error parsing envelope ",
			"id:", id,
			"number:", key.Address,
			"device_id:", key.DeviceID,
			"reason:", err)

		return nil
	}

	guid := envelope.GetServerGuid()
	envelope.ServerGuid = ""

	_, err = storage.MessagesManager{}.Store(guid, &envelope, key.Address, key.DeviceID)
	if err != nil {
		logx.Error("[MessagePersistJob] failed to persist message ",
			" id:", id,
			" number:", key.Address,
			" device_id:", key.DeviceID,
			" reason:", err)

		return err
	}

	err = storage.MessagesCache{}.Remove(key.Address, key.DeviceID, id)
	if err != nil {
		logx.Error("[MessagePersistJob] failed to remove message cache",
			" id", id,
			" number", key.Address,
			" device_id", key.DeviceID,
			" reason", err)
		return err
	}
	return nil
}
