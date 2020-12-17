package middleware

import (
    "errors"
    "fmt"
    eos "github.com/eoscanada/eos-go"
    "github.com/eoscanada/eos-go/ecc"
    "github.com/tal-tech/go-zero/core/logx"
    "privatedb/api/internal/tool"

    "github.com/tal-tech/go-zero/rest/httpx"

    "net/http"

    //owcrypt "github.com/blocktree/go-owcrypt"

    //base58 "github.com/btcsuite/btcutil/base58"
)

var (
    errorUserNameMiss = errors.New("x-user-name head  miss or empty")
    errorUserNameFormat = errors.New("x-user-name format invalid")
    errorUserSignMiss          = errors.New("x-user-sign head miss  or  empty")
    errorUserSign          = errors.New("x-user-sign  invalid")
)

const (
    usernameKey = `x-user-name`
    userSignKey = `x-user-sign`
)


type EOSUserCheckMiddleware struct {
    API  *eos.API
}

func NewEOSUserCheckMiddleware(url string) *EOSUserCheckMiddleware {
    api:=eos.New(url)
    api.EnableKeepAlives()
    api.Debug=true

    return &EOSUserCheckMiddleware{
        API: api,

    }
}

func (m *EOSUserCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // TODO generate middleware implement function, delete after code implementation
        logx.Info("====EOSUsercheckMiddleware =====:")

        userName := r.Header.Get(usernameKey)
        if len(userName)==0 {
            httpx.Error(w, errorUserNameMiss)
            return

        }
        if len(userName)>12 {
            httpx.Error(w, errorUserNameFormat)
            return
        }


        userSign := r.Header.Get(userSignKey)
        if len(userSign)==0 {
            httpx.Error(w, errorUserSignMiss)
            return

        }

        active_public_key:=""
        ownwer_public_key:=""
        accountResp,_err:=m.API.GetAccount(eos.AccountName(userName))
        if _err!=nil{
            httpx.Error(w, _err)
            return
        }
        for index :=range accountResp.Permissions{
            if accountResp.Permissions[index].PermName=="active"{
                active_public_key =accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
            }
            if accountResp.Permissions[index].PermName=="owner"{
                ownwer_public_key =accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
            }
        }

        fmt.Println(active_public_key,ownwer_public_key)

        //验证签名
        pubkey,_err:=ecc.NewPublicKey(active_public_key)
        if _err!=nil{
            httpx.Error(w,errors.New("create real publickey error:"+_err.Error()))
            return
        }
        digest:=tool.Sha256Str(userName)

        s,_err:=ecc.NewSignature(userSign)
        if _err!=nil{
            httpx.Error(w,errors.New("create real signature error:"+_err.Error()))
            return
        }
        if s.Verify(digest,pubkey)==false{
            httpx.Error(w,errorUserSign)
            return
        }


        /*
       private_key:="5K8bmn8AMNewSzgB3VNnz7pahVVLTF7LaksnF8tjoSPVcvS2xDw"
       public_key:= "EOS8GnxsyzChJF4pobmcvYan3qv1ovw66ypYNiJZatSWKkejfEqg4"

       privkey,_:=ecc.NewPrivateKey(private_key)
       pubkey,_:=ecc.NewPublicKey(public_key)
       digest:=tool.Sha256Str(userName)
       signature,_:=privkey.Sign(digest) //对hash值签名
       fmt.Println("签名的内容是：",signature.String(),"|")

       s,_:=ecc.NewSignature(signature.String())
       fmt.Println("验签结果:",s.Verify(digest,pubkey))

         */
        logx.Info("====EOSUsercheckMiddleware =====:")


        // Passthrough to next handler if need
        next(w, r)
    }
}
