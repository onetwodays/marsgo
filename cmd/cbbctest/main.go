package main

import (
    "marsgo/utils/httpclient"
    prom_runtime "marsgo/utils/runtime"
    log "marsgo/utils/ut_log"
)
const (
    APPNAME = "cbbc.exe"
    VERSION = "0.0.1"
)
func init() {
    log.InitLog(APPNAME, VERSION)


}


var (
    baseUrl = "http://192.168.70.131:7000"
    cbbc_deals=`
    {
        "head": {
        "method": "cbbc.deals",
        "msgType": "request",
        "packType": "1",
        "lang": "cn",
        "version": "1.0.0",
        "timestamps": "1439261904",
        "serialNumber": "57"
        },
        "param": {
        "code":"60000BTCEX0107NA", 
        "pageSize":"2",
        "lastId":"1"
        }
    }

`

)

func main()  {
    defer log.Flush()

    log.Info("host_details ", prom_runtime.Uname())
    log.Info("fd_limits ", prom_runtime.FdLimits())
    log.Info("vm_limits ", prom_runtime.VmLimits())

    client:=httpclient.NewHttpClient(baseUrl)
    r,err:=client.Post("/api/v1/cbbc",&cbbc_deals)
    if err==nil{
        log.Info(r.String())
    }else{
        log.Error("deal list :",err)
    }








}
