package dispatch

import (
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub/channel"

	"secret-im/service/signalserver/cmd/api/textsecure"

	"io"
	"sync"
)

// Redis调度管理器
type RedisDispatchManager struct {
	running           bool
	client            *redis.Client
	pubSub            *redis.PubSub
	subscriptions     sync.Map
	pool              *Pool
	deadLetterChannel channel.DispatchChannel //默认的信道
}

// 创建Redis调度管理器
func NewRedisDispatchManager(
	cli *redis.Client, poolSize int, deadLetterChannel ...channel.DispatchChannel) *RedisDispatchManager {

	pool := NewPool(poolSize)
	pubSub := cli.Subscribe()
	manager := RedisDispatchManager{
		client: cli,
		pool:   pool,
		pubSub: pubSub,
	}
	if len(deadLetterChannel) > 0 {
		manager.deadLetterChannel = deadLetterChannel[0]
	}
	go manager.startPolling()
	return &manager
}

// 关闭服务
func (r *RedisDispatchManager) Shutdown() {
	r.subscriptions.Range(func(key, value interface{}) bool {
		r.Unsubscribe(key.(string), value.(channel.DispatchChannel))
		return true
	})
	r.pool.Stop()
	r.pool.WaitForIdle()
}

// 是否订阅
func (r *RedisDispatchManager) HasSubscription(name string) bool {
	_, ok := r.subscriptions.Load(name)
	return ok
}

// 发布消息
func (r *RedisDispatchManager) Publish(name string, message *textsecure.PubSubMessage) (int64, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		logx.Error("[RedisDispatchManager.Publish]",err)
		return 0, err
	}

	cmd := r.client.Publish(name, data)
	if cmd.Err() != nil {
		logx.Error("[RedisDispatchManager.Publish]",cmd.Err())
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

// 订阅消息
func (r *RedisDispatchManager) Subscribe(name string, dispatchChannel channel.DispatchChannel) error {
	previous, ok := r.subscriptions.Load(name)
	r.subscriptions.Store(name, dispatchChannel)
	defer func() {
		if ok && previous != nil {
			r.pool.Add(name, func() {
				previous.(channel.DispatchChannel).OnDispatchUnsubscribed(name)
			})
		}
	}()
	err := r.pubSub.Subscribe(name)
	if err == nil {
		r.pool.Add(name, func() { dispatchChannel.OnDispatchSubscribed(name) })
	}
	return err
}

// 取消订阅
func (r *RedisDispatchManager) Unsubscribe(name string, dispatchChannel channel.DispatchChannel) error {
	value, ok := r.subscriptions.Load(name)
	if !ok || value == nil {
		return nil
	}

	subscription := value.(channel.DispatchChannel)
	if subscription != dispatchChannel {
		return nil
	}

	r.subscriptions.Delete(name)
	defer func() {
		r.pool.Add(name, func() { subscription.OnDispatchUnsubscribed(name) })
	}()
	return r.pubSub.Unsubscribe(name)
}

// 轮询消息
func (r *RedisDispatchManager) startPolling() {
	r.running = true
	defer r.pubSub.Close()

	for {
		message, err := r.pubSub.ReceiveMessage()
		if err != nil {

			if err != io.EOF {
				logx.Error("[PubSub] failed to receive message，推出for",err)
			}
			break
		}

		var pubSubMessage textsecure.PubSubMessage
		err = proto.Unmarshal([]byte(message.Payload), &pubSubMessage)
		if err != nil {
			logx.Errorf("[PubSub] invalid pubsub message,%s",err.Error())

			continue
		}
		//logx.Info("[redis_dispatch]从redis读取到1条textsecure.PubSubMessage格式消息:" ,pubSubMessage.String())

		value, ok := r.subscriptions.Load(message.Channel)
		r.pool.Add(message.Channel, func() {
			if !ok {
				logx.Error("信道",message.Channel,"找不到，启用默认信道")
				if r.deadLetterChannel != nil {
					r.deadLetterChannel.OnDispatchMessage(message.Channel, &pubSubMessage)
				}
			} else {
				//logx.Info("[redis_dispatch]从redis收到1条消息交给信道的 ",message.Channel,"OnDispatchMessage函数处理")
				value.(channel.DispatchChannel).OnDispatchMessage(message.Channel, &pubSubMessage)
			}
		})
	}

	r.running = false
	logx.Info("[PubSub] dispatch manager shutting down...")
}
