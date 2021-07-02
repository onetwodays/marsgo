package dispatch

import (
	"container/list"
	"hash/crc32"
	"secret-im/pkg/utils-tools"
	"sync"
	"sync/atomic"
	"time"
)

// 任务池
type Pool struct {
	workers []*Worker
}

// 创建任务池
func NewPool(size int) *Pool {
	if size <= 0 {
		size = 64
	}

	workers := make([]*Worker, 0, size)
	for i := 0; i < size; i++ {
		worker := newWorker()
		workers = append(workers, worker)
	}
	return &Pool{workers: workers}
}

// 停止工作
func (pool *Pool) Stop() {
	for _, worker := range pool.workers {
		worker.stop()
	}
}

// 等待到空闲
func (pool *Pool) WaitForIdle() {
	for _, worker := range pool.workers {
		worker.waitForIdle()
	}
}

// 添加任务
func (pool *Pool) Add(tag string, task func()) {
	sun := crc32.ChecksumIEEE([]byte(tag))
	idx := int(sun % uint32(len(pool.workers)))
	pool.workers[idx].add(task)
}

// 工作者
type Worker struct {
	stopped int32
	cond    *sync.Cond
	queue   *list.List
}

// 创建工作者
func newWorker() *Worker {
	worker := Worker{
		queue: list.New(),
		cond:  sync.NewCond(new(sync.Mutex)),
	}
	go worker.running()
	return &worker
}

// 停止工作
func (worker *Worker) stop() {
	if atomic.CompareAndSwapInt32(&worker.stopped, 0, 1) {
		worker.cond.Signal()
	}
}

// 添加任务
func (worker *Worker) add(task func()) {
	worker.cond.L.Lock()
	worker.queue.PushBack(task)
	worker.cond.L.Unlock()
	worker.cond.Signal()
}

// 等待到空闲
func (worker *Worker) waitForIdle() {
	for {
		worker.cond.L.Lock()
		size := worker.queue.Len()
		worker.cond.L.Unlock()
		if size == 0 {
			break
		}
		time.Sleep(time.Millisecond * 1)
	}
}

// 执行工作
func (worker *Worker) running() {
	for {
		for {
			worker.cond.L.Lock()
			if worker.queue.Len() > 0 || atomic.LoadInt32(&worker.stopped) == 1 {
				worker.cond.L.Unlock()
				break
			}
			worker.cond.Wait()
			worker.cond.L.Unlock()
		}

		worker.cond.L.Lock()
		if worker.queue.Len() == 0 {
			worker.cond.L.Unlock()
			break
		}
		front := worker.queue.Front()
		worker.cond.L.Unlock()

		utils.NoPanic(front.Value.(func()))

		worker.cond.L.Lock()
		worker.queue.Remove(front)
		worker.cond.L.Unlock()
	}
}
