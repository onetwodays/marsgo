package trace

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	slotLength = 2048
)
//采样???

var ignoreds = []string{"/metrics", "/ping"} // NOTE: add YOUR URL PATH that want ignore

func init() {
	rand.Seed(time.Now().UnixNano())
}

func oneAtTimeHash(s string) (hash uint32) {
	b := []byte(s)
	for i := range b {
		hash += uint32(b[i])
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15
	return
}

// sampler decides whether a new trace should be sampled or not.
type sampler interface {
	IsSampled(traceID uint64, operationName string) (bool, float32)
	Close() error
}

// 下面定义了一个类型和实现了sampler接口.
type probabilitySampling struct {
	probability float32
	slot        [slotLength]int64 //2048个槽?什么意思呢?是个数组啊,保存的是个时间
}

func (p *probabilitySampling) IsSampled(traceID uint64, operationName string) (bool, float32) {
	for _, ignored := range ignoreds {
		if operationName == ignored {
			return false, 0
		}
	}
	now := time.Now().Unix() //时间戳
	idx := oneAtTimeHash(operationName) % slotLength
	old := atomic.LoadInt64(&p.slot[idx])
	if old != now {
		atomic.SwapInt64(&p.slot[idx], now)
		return true, 1
	}
	return rand.Float32() < float32(p.probability), float32(p.probability) //是否sample?
}

func (p *probabilitySampling) Close() error { return nil }

// newSampler new probability sampler
func newSampler(probability float32) sampler {
	if probability <= 0 || probability > 1 {
		panic("probability P ∈ (0, 1]")
	}
	return &probabilitySampling{probability: probability}
}
