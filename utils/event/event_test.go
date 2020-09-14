package event

import (
    "fmt"
    "sync"
    "testing"
)


func TestEvent(t *testing.T){
    t.Run("event", func(t *testing.T) {

        var wg sync.WaitGroup
        //前置依赖
        a1:=New()
        b1:=New()
        c1:=New()
        wg.Add(3)
        go func() {
            defer wg.Done()
            defer a1.Fire()
            fmt.Println("a1")
        }()

        go func() {
            defer wg.Done()
            defer c1.Fire()
            fmt.Println("c1")
        }()

        go func() {
            defer wg.Done()
            defer b1.Fire()
            <- a1.Done()
            <- c1.Done()
            fmt.Println("b1")
        }()
        wg.Wait()
        fmt.Println("over")

    })





}
