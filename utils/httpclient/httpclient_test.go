package httpclient

import (
    "fmt"
    "testing"
)

func TestHttpClient_Get(t *testing.T) {
    httpclient:= NewHttpClient("https://blockchainwhispers.com/bitmex-position-calculator/")
    response,_:=httpclient.Get("#")
    fmt.Println(response.String())


}
