package bbr

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"marsgo/pkg/container/group"
	"marsgo/pkg/ecode"
	"marsgo/pkg/log"
	limit "marsgo/pkg/ratelimit"
	"marsgo/pkg/stat/metric"
	cpustat "marsgo/pkg/stat/sys/cpu"
)



var (
	cpu         int64 //250ms 更新一次
	decay       = 0.95
	initTime    = time.Now()
	defaultConf = &Config{
		Window:       time.Second * 10,//窗口时间是10s
		WinBucket:    100,//10s里有10个桶
		CPUThreshold: 800, //cpu阀是800
	}
)

type cpuGetter func() int64

//定时器每隔250ms 去采cpu的使用率
func init() {
	go cpuproc()
}

// cpu = cpuᵗ⁻¹ * decay + cpuᵗ * (1 - decay)
func cpuproc() {
	ticker := time.NewTicker(time.Millisecond * 250)//每隔250豪秒做一件事情
	defer func() {
		ticker.Stop()
		if err := recover(); err != nil {
			log.Error("rate.limit.cpuproc() err(%+v)", err)
			go cpuproc()
		}
	}()

	// EMA algorithm: https://blog.csdn.net/m0_38106113/article/details/81542863
	for range ticker.C { //阻塞在这里.
		stat := &cpustat.Stat{} //用来保存数据的.
		cpustat.ReadStat(stat)
		prevCpu := atomic.LoadInt64(&cpu)
		curCpu := int64(float64(prevCpu)*decay + float64(stat.Usage)*(1.0-decay))
		atomic.StoreInt64(&cpu, curCpu)
	}
}

// Stats contains the metrics's snapshot of bbr.
type Stat struct {
	Cpu         int64
	InFlight    int64  //飞行
	MaxInFlight int64
	MinRt       int64
	MaxPass     int64
}

// BBR implements bbr-like limiter.
// It is inspired by sentinel.
// https://github.com/alibaba/Sentinel/wiki/%E7%B3%BB%E7%BB%9F%E8%87%AA%E9%80%82%E5%BA%94%E9%99%90%E6%B5%81
/*
Sentinel 系统自适应限流从整体维度对应用入口流量进行控制，
结合应用的 Load、CPU 使用率、总体平均 RT、入口 QPS 和并发线程数等几个维度的监控指标，
通过自适应的流控策略，让系统的入口流量和系统的负载达到一个平衡，让系统尽可能跑在最大吞吐量的同时保证系统整体的稳定性。

判断是否丢弃当前请求的算法如下：

`cpu > 800 AND (Now - PrevDrop) < 1s AND (MaxPass * MinRt * windows / 1000) < InFlight`

MaxPass 表示最近 5s 内，单个采样窗口中最大的请求数。
MinRt 表示最近 5s 内，单个采样窗口中最小的响应时间。
windows 表示一秒内采样窗口的数量，默认配置中是 5s 50 个采样，那么 windows 的值为 10。

*/
type BBR struct {
	cpu             cpuGetter //是个函数.
	passStat        metric.RollingCounter
	rtStat          metric.RollingCounter
	inFlight        int64
	winBucketPerSec int64
	conf            *Config
	prevDrop        atomic.Value
	prevDropHit     int32
	rawMaxPASS      int64
	rawMinRt        int64
}

// Config contains configs of bbr limiter.
type Config struct {
	Enabled      bool
	Window       time.Duration  //窗口时长
	WinBucket    int            //一个窗口几个桶
	Rule         string
	Debug        bool
	CPUThreshold int64         //cpu阀
}

func (l *BBR) maxPASS() int64 {
	rawMaxPass := atomic.LoadInt64(&l.rawMaxPASS)
	if rawMaxPass > 0 && l.passStat.Timespan() < 1 {
		return rawMaxPass
	}
	rawMaxPass = int64(l.passStat.Reduce(func(iterator metric.Iterator) float64 {
		var result = 1.0
		for i := 1; iterator.Next() && i < l.conf.WinBucket; i++ {
			bucket := iterator.Bucket()
			count := 0.0
			for _, p := range bucket.Points {
				count += p
			}
			result = math.Max(result, count)
		}
		return result
	}))
	if rawMaxPass == 0 {
		rawMaxPass = 1
	}
	atomic.StoreInt64(&l.rawMaxPASS, rawMaxPass)
	return rawMaxPass
}

func (l *BBR) minRT() int64 {
	rawMinRT := atomic.LoadInt64(&l.rawMinRt)
	if rawMinRT > 0 && l.rtStat.Timespan() < 1 {
		return rawMinRT
	}
	rawMinRT = int64(math.Ceil(l.rtStat.Reduce(func(iterator metric.Iterator) float64 {
		var result = math.MaxFloat64
		for i := 1; iterator.Next() && i < l.conf.WinBucket; i++ {
			bucket := iterator.Bucket()
			if len(bucket.Points) == 0 {
				continue
			}
			total := 0.0
			for _, p := range bucket.Points {
				total += p
			}
			avg := total / float64(bucket.Count)
			result = math.Min(result, avg)
		}
		return result
	})))
	if rawMinRT <= 0 {
		rawMinRT = 1
	}
	atomic.StoreInt64(&l.rawMinRt, rawMinRT)
	return rawMinRT
}

func (l *BBR) maxFlight() int64 {
	return int64(math.Floor(float64(l.maxPASS()*l.minRT()*l.winBucketPerSec)/1000.0 + 0.5))
}

func (l *BBR) shouldDrop() bool {
	if l.cpu() < l.conf.CPUThreshold {
		prevDrop, _ := l.prevDrop.Load().(time.Duration)
		if prevDrop == 0 {
			return false
		}
		if time.Since(initTime)-prevDrop <= time.Second {
			if atomic.LoadInt32(&l.prevDropHit) == 0 {
				atomic.StoreInt32(&l.prevDropHit, 1)
			}
			inFlight := atomic.LoadInt64(&l.inFlight)
			return inFlight > 1 && inFlight > l.maxFlight()
		}
		l.prevDrop.Store(time.Duration(0))
		return false
	}
	inFlight := atomic.LoadInt64(&l.inFlight)
	drop := inFlight > 1 && inFlight > l.maxFlight()
	if drop {
		prevDrop, _ := l.prevDrop.Load().(time.Duration)
		if prevDrop != 0 {
			return drop
		}
		l.prevDrop.Store(time.Since(initTime))
	}
	return drop
}

// Stat tasks a snapshot of the bbr limiter.
func (l *BBR) Stat() Stat {
	return Stat{
		Cpu:         l.cpu(),
		InFlight:    atomic.LoadInt64(&l.inFlight),
		MinRt:       l.minRT(),
		MaxPass:     l.maxPASS(),
		MaxInFlight: l.maxFlight(),
	}
}

// Allow checks all inbound traffic.
// Once overload is detected, it raises ecode.LimitExceed error.
func (l *BBR) Allow(ctx context.Context, opts ...limit.AllowOption) (func(info limit.DoneInfo), error) {
	allowOpts := limit.DefaultAllowOpts() //空结构体
	for _, opt := range opts {
		opt.Apply(&allowOpts)
	}
	if l.shouldDrop() {
		return nil, ecode.LimitExceed
	}
	atomic.AddInt64(&l.inFlight, 1)
	stime := time.Since(initTime)
	return func(do limit.DoneInfo) {
		rt := int64((time.Since(initTime) - stime) / time.Millisecond)
		l.rtStat.Add(rt)
		atomic.AddInt64(&l.inFlight, -1)
		switch do.Op {
		case limit.Success:
			l.passStat.Add(1)
			return
		default:
			return
		}
	}, nil
}

func newLimiter(conf *Config) limit.Limiter {
	if conf == nil {
		conf = defaultConf
	}
	size := conf.WinBucket
	bucketDuration := conf.Window / time.Duration(conf.WinBucket)
	passStat := metric.NewRollingCounter(metric.RollingCounterOpts{Size: size, BucketDuration: bucketDuration})
	rtStat := metric.NewRollingCounter(metric.RollingCounterOpts{Size: size, BucketDuration: bucketDuration})
	cpu := func() int64 {
		return atomic.LoadInt64(&cpu)
	}
	limiter := &BBR{
		cpu:             cpu,
		conf:            conf,
		passStat:        passStat,
		rtStat:          rtStat,
		winBucketPerSec: int64(time.Second) / (int64(conf.Window) / int64(conf.WinBucket)),
	}
	return limiter
}

// Group represents a class of BBRLimiter and forms a namespace in which
// units of BBRLimiter.
type Group struct {
	group *group.Group
}

// NewGroup new a limiter group container, if conf nil use default conf.
func NewGroup(conf *Config) *Group {
	if conf == nil {
		conf = defaultConf
	}
	//使用的池化创建对象
	group := group.NewGroup(func() interface{} {
		return newLimiter(conf)
	})
	return &Group{
		group: group,
	}
}

// Get get a limiter by a specified key, if limiter not exists then make a new one.
func (g *Group) Get(key string) limit.Limiter {
	limiter := g.group.Get(key)
	return limiter.(limit.Limiter)
}
