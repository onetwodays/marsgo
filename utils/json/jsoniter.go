//+build jsoniter

package json

import (
    "fmt"

    jsoniter "github.com/json-iterator/go"

)

var (
    json = jsoniter.ConfigCompatibleWithStandardLibrary

)

// 生成json形式的字节组
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {

    fmt.Println("Use [jsoniter] package")
    return json.MarshalIndent(v,prefix,indent)
}
