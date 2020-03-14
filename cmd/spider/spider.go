package main

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/extensions"
    "strings"
    "time"

    "database/sql"
    //"github.com/gocolly/colly/debug"
    "io"
    "log"
    "math/rand"
    "os"
)
const(
   MozillaUA = `Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0`
   ChromeUA  =`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36`
)


const (

    USERNAME = "hopexdev"

    PASSWORD = "Dabaicai123!"

    NETWORK  = "tcp"

    SERVER  = "192.168.70.131"

    PORT    = 3306

    DATABASE = "cbbc_base"

)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
    b := make([]byte, rand.Intn(10)+10)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}


//使用标准库的日志库
func init()  {
    errFile,err:=os.OpenFile("spider.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
    if err!=nil{
        log.Fatalln("打开日志文件失败：",err)
    }else{
        log.SetOutput(io.MultiWriter(os.Stderr,errFile))
    }

    log.SetPrefix("spider ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile )
    //log.Panic("1")
    log.Println("start")
    //os.Exit(-1)

}

type LongShort struct {
    Long string
    Short string
}





func main()  {
    dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
    conPool, err := sql.Open("mysql", dsn)



    if err != nil {

        log.Println("Open mysql failed,err:", err)

        return

    }

    sql:="insert into cbbc_base.t_long_short (exchange_name,coincode,longs,shorts,tips) values "

    exchange_map := make(map[string]LongShort,4)
    exchange_map["Total"]    = LongShort{Short:"0",Long:"0"}
    exchange_map["Binance"]  = LongShort{Short:"0",Long:"0"}
    exchange_map["BitMex"]   = LongShort{Short:"0",Long:"0"}
    exchange_map["Bitfinex"] = LongShort{Short:"0",Long:"0"}



    c := colly.NewCollector()
    //c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

    c.UserAgent = ChromeUA
    c.AllowURLRevisit = true

    extensions.RandomUserAgent(c)
    extensions.Referer(c)

    //Called before a request
    c.OnRequest(func(r *colly.Request) {
        //r.Headers.Set("User-Agent",RandomString())
        log.Println("Visiting", r.URL)
    })

    //Called if error occured during the request
    c.OnError(func(r *colly.Response, err error) {
        log.Println("Something went wrong:", err)
        //time.Sleep(3*time.Second)
        r.Request.Visit(r.Request.URL.String())
    })

    // Called after response received
    c.OnResponse(func(r *colly.Response) {
        log.Println("Visited", r.Request.URL," ",r.StatusCode)
        //log.Println(string(r.Body))

    })

    //Called right after OnResponse if the received content is HTML
    /*
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        //log.Println(e.Attr("href"))
    })
    */
                             //body > div.main-content > div.container-fluid > div > div > div > div > div.col-md-12.visible-xs.hidden > div > div.owl-stage-outer > div > div:nth-child(1) > div > span.portfolio-profit-wrapper
                             //body > div.main-content > div:nth-child(2) > div:nth-child(3) > div > div > span
    c.OnHTML( `body > div.main-content > div:nth-child(2) >  div:nth-child(3) > div > div `,
        func(e *colly.HTMLElement) {


            var builder strings.Builder
            builder.WriteString(sql)

            tip:=e.ChildText(`span`)

            e.ForEach(`div:nth-child(4)>div.col-md-3`, func(i int, element *colly.HTMLElement) {
                h3:= element.ChildText("h3")
                long:=element.ChildText(`div>div:nth-child(2) > div.field.long>small`)
                short:=element.ChildText(`div>div:nth-child(3) > div.field.short>small`)
                h3s:=strings.Split(h3," ")
                if len(h3s)>0{
                    h3=h3s[0]
                }


                fmt.Println(i,":",tip,h3,long,short)

                long  = strings.TrimSuffix(long,"%")
                short = strings.TrimSuffix(short,"%")

                mapValue:= exchange_map[h3]
                if mapVaue.Long==long && mapValue.Short==short{
                    log.Println("本次采集数据跟上次一样没变化")
                    return

                }else{
                    mapValue.Long = long
                    mapValue.Short=short
                    log.Println("本次采集数据跟上次发生变化,更新缓存并写入数据库")
                }


                fmt.Println("-------------")

                builder.WriteString(" (")

                builder.WriteString("'")
                builder.WriteString(h3)
                builder.WriteString("',")

                builder.WriteString("'BTC',")



                builder.WriteString(long)
                builder.WriteString(",")

                builder.WriteString(short)
                builder.WriteString(",")

                builder.WriteString("'")
                builder.WriteString(tip)
                builder.WriteString("'),")
            })
            sqlex:=builder.String()
            sqlex= strings.TrimRight(sqlex,",")
            log.Println(sqlex)
            _,err:=conPool.Exec(sqlex)
            if err!=nil{
                log.Println(sqlex," ",err)
            }




        

    })

    /*
    c.OnHTML(`body > div.main-content > div.container-fluid > div > div > div > div > div.col-md-12.visible-xs.hidden > div > div.owl-stage-outer > div > div:nth-child(1) > div > span.portfolio-profit-wrapper`, func(element *colly.HTMLElement) {
        fmt.Println("2->",element.Text)
    })
    */

    for{
        c.Visit("https://blockchainwhispers.com/bitmex-position-calculator/")
        time.Sleep(3*time.Second)

    }
}

