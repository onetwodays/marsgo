package main

import (
    "fmt"
    redis "github.com/go-redis/redis/v7"
)

func main()  {

    /*
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // use default Addr
        Password: "",               // no password set
        DB:       0,                // use default DB
    })

    pong, err := rdb.Ping().Result()
    fmt.Println(pong, err)
     */

    /*
    type FailoverOptions struct {
        // The master name.
        MasterName string
        // A seed list of host:port addresses of sentinel nodes.
        SentinelAddrs    []string
        SentinelPassword string

        Dialer    func(ctx context.Context, network, addr string) (net.Conn, error)
        OnConnect func(*Conn) error

        Password string
        DB       int

        MaxRetries      int
        MinRetryBackoff time.Duration
        MaxRetryBackoff time.Duration

        DialTimeout  time.Duration
        ReadTimeout  time.Duration
        WriteTimeout time.Duration

        PoolSize           int
        MinIdleConns       int
        MaxConnAge         time.Duration
        PoolTimeout        time.Duration
        IdleTimeout        time.Duration
        IdleCheckFrequency time.Duration

        TLSConfig *tls.Config
    }
     */
    rdb := redis.NewFailoverClient(&redis.FailoverOptions{
        MasterName:    "mymaster",
        SentinelAddrs: []string{"192.168.70.131:26379"},
    })
    fmt.Println(rdb.Ping().String())
    fmt.Printf("%+v",(rdb.Options()))


}