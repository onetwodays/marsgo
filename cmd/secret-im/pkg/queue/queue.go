package queue
import (
	"sync"

	equeue "gopkg.in/eapache/queue.v1"
)

// 同步队列
type SyncQueue struct {
	closed bool
	mutex  sync.Mutex
	cond   *sync.Cond
	queue  *equeue.Queue
	ch     chan interface{}
}

// 创建同步队列
func NewSyncQueue() *SyncQueue {
	q := SyncQueue{
		queue: equeue.New(),
		ch:    make(chan interface{}),
	}
	q.cond = sync.NewCond(&q.mutex)
	return &q
}

// 队列长度
func (q *SyncQueue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.queue.Length()
}

// 放入元素
func (q *SyncQueue) Push(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if !q.closed {
		q.queue.Add(v)
		q.cond.Signal()
	}
}

// 取出元素
func (q *SyncQueue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var v interface{}
	for q.queue.Length() == 0 && !q.closed {
		q.cond.Wait()
	}

	if q.queue.Length() > 0 {
		v = q.queue.Peek()
		q.queue.Remove()
	}
	return v
}

// 关闭队列
func (q *SyncQueue) Close() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if !q.closed {
		q.closed = true
		q.cond.Signal()
	}
}