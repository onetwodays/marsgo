package paladin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify" // 监控文件变化
)

const (
	defaultChSize = 10  // channel长度,change event size
)

var _ Client = &file{} //验证一下file 是否实现了Client接口

// file is file config client.
type file struct {
	values *Map                  // 一大块原子内存
	rawVal map[string]*Value     // key -value的形式

	watchChs map[string][]chan Event // 每个key对应多个事件队列.因为有好几类事件,一类事件可以有多个事件保存在队列里面
	mx       sync.Mutex
	wg       sync.WaitGroup

	base string  //配置文件所在的父目录
	done chan struct{} //协程退出信号
}

func isHiddenFile(name string) bool {
	// TODO: support windows.
	return strings.HasPrefix(filepath.Base(name), ".")
}

func readAllPaths(base string) ([]string, error) {
	fi, err := os.Stat(base)
	if err != nil {
		return nil, fmt.Errorf("check local config file fail! error: %s", err)
	}
	// dirs or file to paths
	var paths []string
	if fi.IsDir() {
		files, err := ioutil.ReadDir(base)
		if err != nil {
			return nil, fmt.Errorf("read dir %s error: %s", base, err)
		}
		for _, file := range files {
			if !file.IsDir() && !isHiddenFile(file.Name()) {
				paths = append(paths, path.Join(base, file.Name()))
			}
		}
	} else {
		paths = append(paths, base)
	}
	return paths, nil
}

func loadValuesFromPaths(paths []string) (map[string]*Value, error) {
	// laod config file to values
	var err error
	values := make(map[string]*Value, len(paths))
	for _, fpath := range paths {
		if values[path.Base(fpath)], err = loadValue(fpath); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func loadValue(fpath string) (*Value, error) {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	content := string(data)
	return &Value{val: content, raw: content}, nil
}

// NewFile new a config file client.
// conf = /data/conf/app/
// conf = /data/conf/app/xxx.toml
func NewFile(base string) (Client, error) {
	// paltform slash
	base = filepath.FromSlash(base)

	paths, err := readAllPaths(base)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("empty config path")
	}

	rawVal, err := loadValuesFromPaths(paths)
	if err != nil {
		return nil, err
	}

	valMap := &Map{}      // 原子性的key-value字典
	valMap.Store(rawVal)  // key - value 的字典.
	fc := &file{
		values:   valMap, // 他俩保存的是同一个变量,文件名-文件内容字符串
		rawVal:   rawVal, // (map[string]*Value
		watchChs: make(map[string][]chan Event), //监控每个文件的变化事件,key是文件名

		base: base,
		done: make(chan struct{}, 1),
	}

	fc.wg.Add(1)
	go fc.daemon()

	return fc, nil
}

// Get return value by key.
func (f *file) Get(key string) *Value {
	return f.values.Get(key)
}

// GetAll return value map.
func (f *file) GetAll() *Map {
	return f.values
}

// WatchEvent watch multi key. 返回一个事件队列,key多是一个文件名
func (f *file) WatchEvent(ctx context.Context, keys ...string) <-chan Event {
	f.mx.Lock()
	defer f.mx.Unlock()
	ch := make(chan Event, defaultChSize)
	for _, key := range keys {
		f.watchChs[key] = append(f.watchChs[key], ch) //才放入一个事件队列.
	}
	return ch
}

// Close close watcher.
func (f *file) Close() error {
	f.done <- struct{}{}
	f.wg.Wait()
	return nil
}

// file config daemon to watch file modification
func (f *file) daemon() {
	defer f.wg.Done()
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("create file watcher fail! reload function will lose efficacy error: %s", err)
		return
	}
	if err = fswatcher.Add(f.base); err != nil {
		log.Printf("create fsnotify for base path %s fail %s, reload function will lose efficacy", f.base, err)
		return
	}
	log.Printf("start watch filepath: %s", f.base)
	for event := range fswatcher.Events {
		switch event.Op {
		// use vim edit config will trigger rename
		case fsnotify.Write, fsnotify.Create:
			f.reloadFile(event.Name) // Relative path to the file or directory
		case fsnotify.Chmod:
		default:
			log.Printf("unsupport event %s ingored", event)
		}
	}
}

func (f *file) reloadFile(name string) {
	if isHiddenFile(name) {
		return
	}
	// NOTE: in some case immediately read file content after receive event
	// will get old content, sleep 100ms make sure get correct content.
	time.Sleep(100 * time.Millisecond)
	key := filepath.Base(name)
	val, err := loadValue(name)
	if err != nil {
		log.Printf("load file %s error: %s, skipped", name, err)
		return
	}
	f.rawVal[key] = val
	f.values.Store(f.rawVal)

	f.mx.Lock()
	chs := f.watchChs[key] //key是文件名.,找到这个文件名的对应的事件队列切片.
	f.mx.Unlock()

	for _, ch := range chs { //key 对应多个通道,对于每个通道
		select {
		case ch <- Event{Event: EventUpdate, Value: val.raw}:
		default:
			log.Printf("event channel full discard file %s update event", name)
		}
	}
}
