package metric

import (
	"sync"
	"time"
)

// RollingPolicy is a policy for ring window(滑动窗口) based on time duration.
// RollingPolicy moves bucket offset with time duration.
// e.g. If the last point is appended one bucket duration ago,
// RollingPolicy will increment current offset.
// 基于时间间隔的滑动窗口,
type RollingPolicy struct {
	mu     sync.RWMutex
	size   int //窗口的桶数量
	window *Window //窗口
	offset int

	bucketDuration time.Duration //桶时间
	lastAppendTime time.Time
}

// RollingPolicyOpts contains the arguments for creating RollingPolicy.
type RollingPolicyOpts struct {
	BucketDuration time.Duration
}

// NewRollingPolicy creates a new RollingPolicy based on the given window and RollingPolicyOpts.
func NewRollingPolicy(window *Window, opts RollingPolicyOpts) *RollingPolicy {
	return &RollingPolicy{
		window: window,
		size:   window.Size(),
		offset: 0,

		bucketDuration: opts.BucketDuration,
		lastAppendTime: time.Now(),
	}
}

// 上一次采样到现在的时间跨度,需要多少个桶
func (r *RollingPolicy) timespan() int {
	v := int(time.Since(r.lastAppendTime) / r.bucketDuration)
	if v > -1 { // maybe time backwards
		return v
	}
	return r.size
}


func (r *RollingPolicy) add(f func(offset int, val float64), val float64) {
	r.mu.Lock()
	//算一下f的入参offset
	timespan := r.timespan() //最后一次添加时间距现在的时间/桶持续时间,多少个桶
	if timespan > 0 {
		r.lastAppendTime = r.lastAppendTime.Add(time.Duration(timespan * int(r.bucketDuration)))
		offset := r.offset
		// reset the expired buckets
		s := offset + 1
		if timespan > r.size {
			timespan = r.size
		}
		e, e1 := s+timespan, 0 // e: reset offset must start from offset+1
		if e > r.size {
			e1 = e - r.size // r.offset+1+timespan-r.size
			e = r.size
		}
		//重置过期桶 e = offset+1+timespan
		for i := s; i < e; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		//e1>=0
		for i := 0; i < e1; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		r.offset = offset
	}
	f(r.offset, val)  //找到对应的桶
	r.mu.Unlock()
}

// Append appends the given points to the window.
func (r *RollingPolicy) Append(val float64) {
	r.add(r.window.Append, val)
}

// Add adds the given value to the latest point within bucket.
func (r *RollingPolicy) Add(val float64) {
	r.add(r.window.Add, val)
}

// Reduce applies the reduction function to all buckets within the window.
func (r *RollingPolicy) Reduce(f func(Iterator) float64) (val float64) {
	r.mu.RLock()
	timespan := r.timespan()
	if count := r.size - timespan; count > 0 { //不是for循环
		offset := r.offset + timespan + 1
		if offset >= r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	r.mu.RUnlock()
	return val
}
