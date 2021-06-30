package queue


import (
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

var (
	once sync.Once
	nc   *nats.Conn
)

// 初始化模块
func InitModule(c *nats.Conn) {
	once.Do(func() {
		nc = c
	})
}

// 发布消息
func Publish(topic string, body proto.Message) error {
	if body == nil {
		return nil
	}

	data, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	return nc.Publish(topic, data)
}

