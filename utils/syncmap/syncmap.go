package syncmap
import (
    "sync"
)

//安全的Map
type SyncMap struct {
    rw *sync.RWMutex
    data map[interface{}]interface{}
}
//存储操作
func (sm *SyncMap) Put(k,v interface{}){
    sm.rw.Lock()
    defer sm.rw.Unlock()

    sm.data[k]=v
}
//获取操作
func (sm *SyncMap) Get(k interface{}) interface{}{
    sm.rw.RLock()
    defer sm.rw.RUnlock()

    return sm.data[k]
}

//删除操作
func (sm *SyncMap) Delete(k interface{}) {
    sm.rw.Lock()
    defer sm.rw.Unlock()

    delete(sm.data,k)
}

//遍历Map，并且把遍历的值给回调函数，可以让调用者控制做任何事情
func (sm *SyncMap) Each(cb func (interface{},interface{})){
    sm.rw.RLock()
    defer sm.rw.RUnlock()

    for k, v := range sm.data {
        cb(k,v)
    }
}

//生成初始化一个SyncMap
func NewSyncMap() *SyncMap{
    return &SyncMap{
        rw:new(sync.RWMutex),
        data:make(map[interface{}]interface{}),
    }
}