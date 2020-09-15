package breaker

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"marsgo/pkg/ecode"
	"marsgo/pkg/log"
	"marsgo/pkg/stat/metric"
)

// sreBreaker is a sre CircuitBreaker pattern. Circuit(电路图)
// Circuit-Breaker的作用就好似可能失败操作的代理。
//代理会监控最近发生的错误，然后依据这一信息来决定是否允许操作的继续执行，或者直接立刻返回异常信息
type sreBreaker struct {
	stat metric.RollingCounter //指标数据
	r    *rand.Rand
	// rand.New(...) returns a non thread safe object
	randLock sync.Mutex

	k       float64
	request int64

	state int32
}

func newSRE(c *Config) Breaker {
	counterOpts := metric.RollingCounterOpts{
		Size:           c.Bucket,                                         //滑动窗口桶的数目
		BucketDuration: time.Duration(int64(c.Window) / int64(c.Bucket)), //每个桶的时间
	}
	stat := metric.NewRollingCounter(counterOpts)
	return &sreBreaker{
		stat: stat,
		r:    rand.New(rand.NewSource(time.Now().UnixNano())),

		request: c.Request,
		k:       c.K,
		state:   StateClosed, //状态,默认值
	}
}

func (b *sreBreaker) summary() (success int64, total int64) {
	b.stat.Reduce(func(iterator metric.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			total += bucket.Count //是个计数
			for _, p := range bucket.Points {
				success += int64(p)
			}
		}
		return 0
	})
	return
}

func (b *sreBreaker) Allow() error {
	success, total := b.summary()
	k := b.k * float64(success)
	if log.V(5) {
		log.Info("breaker: request: %d, succee: %d, fail: %d", total, success, total-success)
	}
	// check overflow requests = K * success
	//关闭断路器,准许访问后端
	if total < b.request || float64(total) < k {
		if atomic.LoadInt32(&b.state) == StateOpen {
			atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClosed)
		}
		return nil
	}
	//打开断路器,不许访问后端
	if atomic.LoadInt32(&b.state) == StateClosed {
		atomic.CompareAndSwapInt32(&b.state, StateClosed, StateOpen)
	}
	dr := math.Max(0, (float64(total)-k)/float64(total+1))
	drop := b.trueOnProba(dr)
	if log.V(5) {
		log.Info("breaker: drop ratio: %f, drop: %t", dr, drop)
	}
	if drop {
		return ecode.ServiceUnavailable
	}
	return nil
}

func (b *sreBreaker) MarkSuccess() {
	b.stat.Add(1)
}

func (b *sreBreaker) MarkFailed() {
	// NOTE: when client reject requets locally, continue add counter let the
	// drop ratio higher.
	b.stat.Add(0)
}

func (b *sreBreaker) trueOnProba(proba float64) (truth bool) {
	b.randLock.Lock()
	truth = b.r.Float64() < proba
	b.randLock.Unlock()
	return
}
