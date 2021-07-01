package websocket

import (
	"errors"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"sync/atomic"
	"time"
)

// 异步任务结果
type Future struct {
	done   int32
	result chan futureResult
}

// 创建Future
func newFuture() *Future {
	return &Future{result: make(chan futureResult, 1)}
}

// 任务结果
type futureResult struct {
	err      error
	response *textsecure.WebSocketResponseMessage
}

// 获取结果
func (future *Future) Get(timeout time.Duration) (*textsecure.WebSocketResponseMessage, error) {
	if timeout <= 0 {
		result, ok := <-future.result
		if !ok {
			return nil, errors.New("canceled")
		}
		return result.response, result.err
	}

	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		future.cancel()
		return nil, errors.New(" ")
	case result, ok := <-future.result:
		timer.Stop()
		if !ok {
			return nil, errors.New("canceled")
		}
		return result.response, result.err
	}
}

// 取消任务
func (future *Future) cancel() {
	if atomic.CompareAndSwapInt32(&future.done, 0, 1) {
		close(future.result)
	}
}

// 设置结果
func (future *Future) setResult(response *textsecure.WebSocketResponseMessage, err error) {
	if atomic.CompareAndSwapInt32(&future.done, 0, 1) {
		future.result <- futureResult{
			err:      err,
			response: response,
		}
		close(future.result)
	}
}
