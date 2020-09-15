package trace

import (
	"context"
	"encoding/binary"
	"math/rand"
	"time"

	"marsgo/pkg/conf/env"
	"marsgo/pkg/net/ip"

	"github.com/pkg/errors"
)

var _hostHash byte // 0~256

func init() {
	rand.Seed(time.Now().UnixNano())
	_hostHash = byte(oneAtTimeHash(env.Hostname)) //主机名称.初始化
}

//收集一些信息而已.
func extendTag() (tags []Tag) {
	tags = append(tags,
		TagString("region", env.Region),
		TagString("zone", env.Zone),
		TagString("hostname", env.Hostname),
		TagString("ip", ip.InternalIP()),
	)
	return
}

func genID() uint64 {
	var b [8]byte  //8字节数组
	// i think this code will not survive(生存) to 2106-02-07
	binary.BigEndian.PutUint32(b[4:], uint32(time.Now().Unix())>>8)
	b[4] = _hostHash
	binary.BigEndian.PutUint32(b[:4], uint32(rand.Int31()))
	return binary.BigEndian.Uint64(b[:])
}

type stackTracer interface {
	StackTrace() errors.StackTrace //调用堆栈
}

type ctxKey string

var _ctxkey ctxKey = "marsgo/pkg/net/trace.trace" //上下文key,默认值

// FromContext returns the trace bound to the context, if any. 把trace_key 和trace_value 传递给上下文.
func FromContext(ctx context.Context) (t Trace, ok bool) {
	t, ok = ctx.Value(_ctxkey).(Trace)
	return
}

// NewContext new a trace context.
// NOTE: This method is not thread safe.
func NewContext(ctx context.Context, t Trace) context.Context {
	return context.WithValue(ctx, _ctxkey, t)
}
