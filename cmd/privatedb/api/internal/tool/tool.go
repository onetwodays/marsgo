package tool

import (
    "crypto/sha256"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"


    eos "github.com/eoscanada/eos-go"
)


func HttpPost() {
    resp, err := http.Post("http://zhongyingying.qicp.io:38000/v1/chain/get_account",
        "application/json",
        strings.NewReader(`{"account_name": "zhouhao"}`))
    if err != nil {
        fmt.Println(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // handle error
    }

    fmt.Println(string(body))

    Test()
}

func Sha256Str(str string) []byte {

    h := sha256.New()

    h.Write([]byte(str))

    return h.Sum(nil)

}

func Test(){
    api := eos.New("http://zhongyingying.qicp.io:38000")


    accountResp, _ := api.GetAccount("otcexchange")
    fmt.Println("Public key for otcexchange:", accountResp.Permissions[0].RequiredAuth.Keys[0].PublicKey.String())

}