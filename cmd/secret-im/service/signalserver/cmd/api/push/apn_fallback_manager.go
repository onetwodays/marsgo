package push

import (
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/push/operation"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Redis操作
type RedisOperation struct {
	client          *redis.Client
	getOperation    *operation.GetOperation
	insertOperation *operation.InsertOperation
	removeOperation *operation.RemoveOperation
}

// 创建RedisOperation
func NewRedisOperation(client *redis.Client) (*RedisOperation, error) {
	getOperation, err := operation.NewGetOperation(client)
	if err != nil {
		return nil, err
	}
	insertOperation, err := operation.NewInsertOperation(client)
	if err != nil {
		return nil, err
	}
	removeOperation, err := operation.NewRemoveOperation(client)
	if err != nil {
		return nil, err
	}

	redisOperation := RedisOperation{
		client:          client,
		getOperation:    getOperation,
		insertOperation: insertOperation,
		removeOperation: removeOperation,
	}
	return &redisOperation, nil
}

// 取消推送计划
func (redisOperation *RedisOperation) Cancel(number string, deviceID int64) error {
	_, err := redisOperation.removeOperation.Remove(number, deviceID)
	return err
}

// 列入推送计划
func (redisOperation *RedisOperation) Schedule(number string, deviceID int64) error {
	timestamp := utils.CurrentTimeMillis() + (15 * 1000)
	return redisOperation.insertOperation.Insert(number, deviceID, timestamp, 15*1000)
}

// 是否安排计划
func (redisOperation *RedisOperation) IsScheduled(number string, deviceID int64) (bool, error) {
	endpoint := operation.GetEndpointKey(number, deviceID)
	cmd := redisOperation.client.ZScore(operation.PendingNotificationsKey, endpoint)
	if cmd.Err() == redis.Nil {
		return false, nil
	}

	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return true, nil
}

// Fallback管理器
type ApnFallbackManager struct {
	running        int32
	stopChan       chan struct{}
	redisOperation *RedisOperation
}

// 创建Fallback管理器
func NewApnFallbackManager(redisOperation *RedisOperation) *ApnFallbackManager {
	manager := ApnFallbackManager{
		redisOperation: redisOperation,
	}
	return &manager
}

// 开始服务
func (manager *ApnFallbackManager) Start() error {
	if !atomic.CompareAndSwapInt32(&manager.running, 0, 1) {
		return nil
	}

	manager.stopChan = make(chan struct{})
	go manager.run()
	log.Info("[ApnFallbackManager] job started")
	return nil
}

// 停止服务
func (manager *ApnFallbackManager) Stop() {
	if !atomic.CompareAndSwapInt32(&manager.running, 1, 0) {
		return
	}
	<-manager.stopChan
	log.Info("[ApnFallbackManager] job stopped")
}

// 运行服务
func (manager *ApnFallbackManager) run() {
	defer func() {
		atomic.StoreInt32(&manager.running, 0)
		close(manager.stopChan)
	}()

	for {
		if atomic.LoadInt32(&manager.running) == 0 {
			break
		}

		err := manager.handlePending()
		if err != nil {
			log.Infof("[ApnFallbackManager] exception while operating,reason:%s",err.Error())
		}
		time.Sleep(time.Second * 1)
	}
}

// 处理挂起计划
func (manager *ApnFallbackManager) handlePending() error {
	return nil

	/*
	removeOperation := manager.redisOperation.removeOperation
	pendingNotifications, err := manager.redisOperation.getOperation.GetPending(100)
	if err != nil {
		return err
	}

	for _, pendingNotification := range pendingNotifications {
		numberAndDevice := pendingNotification
		number, deviceID, ok := manager.getSeparated(numberAndDevice)

		if !ok {
			if _, err = removeOperation.RemoveByEndpoint(numberAndDevice); err != nil {
				return err
			}
			continue
		}




		endpoint := operation.GetEndpointKey(number, deviceID)


		account, err := storage.AccountsManager{}.GetByNumber(number)
		if err != nil {
			if _, err = removeOperation.RemoveByEndpoint(endpoint); err != nil {
				return err
			}
			continue
		}


		device, ok := account.GetDevice(deviceID)
		if !ok {
			if _, err = removeOperation.RemoveByEndpoint(endpoint); err != nil {
				return err
			}
			continue
		}

		if device.VoipApnID == nil {
			if _, err = removeOperation.Remove(number, deviceID); err != nil {
				return err
			}
			continue
		}

		if device.LastSeen < utils.CurrentTimeMillis()-utils.DaysToMillis(90) {
			if _, err = removeOperation.Remove(number, deviceID); err != nil {
				return err
			}
			continue
		}

		apnMessage := ApnMessage{
			ApnID:    *device.VoipApnID,
			Number:   number,
			DeviceID: deviceID,
			IsVoip:   true,
		}
		AddToApnMessageQueue(apnMessage)
	}



	return nil

	 */
}

// 获取设备信息
func (manager *ApnFallbackManager) getSeparated(encoded string) (number string, deviceID int64, ok bool) {
	parts := strings.Split(encoded, ":")
	if len(parts) != 2 {
		logx.Info("[ApnFallbackManager] got strange encoded number,value=",encoded)
		return "", 0, false
	}

	deviceID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		logx.Errorf("[ApnFallbackManager] badly formatted: %s,error=%s", encoded,err.Error())
		return "", 0, false
	}

	return parts[0], deviceID, true


}

