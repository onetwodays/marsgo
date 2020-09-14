package event


import (
    "sync"
    "sync/atomic"
)
type Event struct{
    fired int32
    c chan struct{}
    once sync.Once
}


func New() *Event{
    return &Event{c:make(chan struct{})}
}

func (ev *Event) Fire() int32{
    atomic.AddInt32(&ev.fired,1)
    ev.once.Do(func(){
        close(ev.c)
    })
    return ev.fired
}

func (ev *Event) Done() <- chan struct{} {
    return ev.c
}

func (ev *Event) HasFired() bool {
    return atomic.LoadInt32(&ev.fired)>0
}

