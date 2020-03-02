package main

import (
    "fmt"

    log "marsgo/utils/ut_log"

    "marsgo/utils/mapiter"

)

const (
    APPNAME = "matchengine-go.exe"
    VERSION = "0.0.1"
)

func init()  {
    log.InitLog(APPNAME,VERSION)

}

func selfInfo(k, v interface{}) {
    fmt.Printf("大家好,我叫%s,今年%d岁\n", k, v)
}


func main()  {
    defer func() {
        log.Flush()
    }()
    log.Info("start run .....")
    persons := make(map[interface{}]interface{})
    persons["张三"] = 20
    persons["李四"] = 23
    persons["王五"] = 26
    mapiter.EachFunc(persons, selfInfo)
}