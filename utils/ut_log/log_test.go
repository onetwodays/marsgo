package ut_log

import "testing"



func TestInfo(t *testing.T) {
    appName:="test"
    version:="1234"
    InitLog(appName,version)
    defer Flush()
    Info("hello")

}
