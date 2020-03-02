//+build !jsoniter

package json

import(
    stdjson "encoding/json"
    "fmt"
)

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
    fmt.Println("Use [encoding/json] package")
    return stdjson.MarshalIndent(v,prefix,indent)
}
