package breaker

import (
	"sync"
	"time"

	xtime "marsgo/pkg/time"
)

// Config broker config.
type Config struct {
	SwitchOff bool // breaker switch,default off.

	// Google
	K float64 //触发熔断的错误率（K = 1 - 1/错误率）

	Window  xtime.Duration //统计桶窗口时间
	Bucket  int            //统计桶大小
	Request int64          //触发熔断的最少请求数量（请求少于该值时不会触发熔断）
}

func (conf *Config) fix() {
	if conf.K == 0 {
		conf.K = 1.5
	}
	if conf.Request == 0 {
		conf.Request = 100 //达到100个请求才触发熔断
	}
	if conf.Bucket == 0 {
		conf.Bucket = 10 //10个桶
	}
	if conf.Window == 0 {
		conf.Window = xtime.Duration(3 * time.Second) // 3s
	}
}

// Breaker is a CircuitBreaker pattern.
// FIXME on int32 atomic.LoadInt32(&b.on) == _switchOn
type Breaker interface {
	Allow() error //判断是否有错误,来表示是否发生熔断
	MarkSuccess()
	MarkFailed()
}

// Group represents a class of CircuitBreaker and forms a namespace in which
// units of CircuitBreaker.
type Group struct {
	mu   sync.RWMutex
	brks map[string]Breaker
	conf *Config
}

const (
	// StateOpen when circuit breaker open, request not allowed, after sleep
	// some duration, allow one single request for testing the health, if ok
	// then state reset to closed, if not continue the step.
	StateOpen int32 = iota
	// StateClosed when circuit breaker closed, request allowed, the breaker
	// calc the succeed ratio, if request num greater request setting and
	// ratio lower than the setting ratio, then reset state to open.
	StateClosed
	// StateHalfopen when circuit breaker open, after slepp some duration, allow
	// one request, but not state closed.
	StateHalfopen

	//_switchOn int32 = iota
	// _switchOff
)

var (
	_mu   sync.RWMutex //保护全局变量_conf
	_conf = &Config{
		Window:  xtime.Duration(3 * time.Second),//窗口时间是3s
		Bucket:  10,
		Request: 100,

		// Percentage of failures must be lower than 33.33%
		K: 1.5,

		// Pattern: "",
	}
	_group = NewGroup(_conf)
)

// Init init global breaker config, also can reload config after first time call.
func Init(conf *Config) {
	if conf == nil {
		return
	}
	_mu.Lock()
	_conf = conf
	_mu.Unlock()
}
// 全局函数
// Go runs your function while tracking the breaker state of default group.
func Go(name string, run, fallback func() error) error {
	breaker := _group.Get(name)
	if err := breaker.Allow(); err != nil {
		return fallback()
	}
	return run()
}

// newBreaker new a breaker.
func newBreaker(c *Config) (b Breaker) {
	// factory
	return newSRE(c)
}

// NewGroup new a breaker group container, if conf nil use default conf.
func NewGroup(conf *Config) *Group {
	if conf == nil {
		_mu.RLock()
		conf = _conf
		_mu.RUnlock()
	} else {
		conf.fix()
	}
	return &Group{
		conf: conf,
		brks: make(map[string]Breaker),
	}
}

// Get get a breaker by a specified key, if breaker not exists then make a new one.
func (g *Group) Get(key string) Breaker {
	g.mu.RLock()
	brk, ok := g.brks[key]
	conf := g.conf
	g.mu.RUnlock()
	if ok {
		return brk
	}
	// NOTE here may new multi breaker for rarely case, let gc drop it.
	brk = newBreaker(conf)
	g.mu.Lock()
	if _, ok = g.brks[key]; !ok {
		g.brks[key] = brk
	}
	g.mu.Unlock()
	return brk
}

// Reload reload the group by specified config, this may let all inner breaker
// reset to a new one.
func (g *Group) Reload(conf *Config) {
	if conf == nil {
		return
	}
	conf.fix()
	g.mu.Lock()
	g.conf = conf
	g.brks = make(map[string]Breaker, len(g.brks))
	g.mu.Unlock()
}

// Go runs your function while tracking the breaker state of group.
func (g *Group) Go(name string, run, fallback func() error) error {
	breaker := g.Get(name)
	if err := breaker.Allow(); err != nil {
		return fallback()
	}
	return run()
}
