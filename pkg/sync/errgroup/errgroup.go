package errgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall(总体) task.
//
// A zero Group is valid and does not cancel on error.
type Group struct {
	err     error  //保存错误
	wg      sync.WaitGroup  //
	errOnce sync.Once  //只被执行一次

	workerOnce sync.Once
	ch         chan func(ctx context.Context) error //ch 用来保存函数
	chs        []func(ctx context.Context) error    //函数切片,当ch满了时候,f放到chs里面,否则放在ch里面

	ctx    context.Context
	cancel func()
}

// WithContext create a Group.
// given function from Go will receive this context,
func WithContext(ctx context.Context) *Group {
	return &Group{ctx: ctx}
}

// WithCancel create a new Group and an associated Context derived from ctx.
//
// given function from Go will receive context derived from this ctx,
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithCancel(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{ctx: ctx, cancel: cancel}
}

func (g *Group) do(f func(ctx context.Context) error) {
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		g.wg.Done()
	}()
	err = f(ctx)
}

// GOMAXPROCS set max goroutine to work.
func (g *Group) GOMAXPROCS(n int) {
	if n <= 0 {
		panic("errgroup: GOMAXPROCS must great than 0")
	}
	g.workerOnce.Do(func() {
		g.ch = make(chan func(context.Context) error, n) //这里才创建channel的内存.
		for i := 0; i < n; i++ {
			go func() {
				for f := range g.ch {
					g.do(f) //调用f
				}
			}()
		}
	})
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(f func(ctx context.Context) error) {
	g.wg.Add(1)
	if g.ch != nil {
		select {
		case g.ch <- f:
		default:
			g.chs = append(g.chs, f)
		}
		return
	}
	go g.do(f)
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	if g.ch != nil {
		for _, f := range g.chs {
			g.ch <- f
		}
	}
	g.wg.Wait()
	if g.ch != nil {
		close(g.ch) // let all receiver exit
	}
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
