package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"

	"marsgo/pkg/naming"
)

func TestNew(t *testing.T) {

	config := &clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 3,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}
	builder, err := New(config)

	if err != nil {
		fmt.Println("etcd 连接失败")
		return
	}
	app1 := builder.Build("app1") //服务发现

	//一直监控这个节点是否发生变化
	go func() {
		fmt.Printf("看下app1 是否有事件到来 \n")
		for {
			select {
			case <-app1.Watch():
				//fmt.Printf("app1 节点发生变化 \n")
			}

		}

	}()
	time.Sleep(time.Second)

	//注册2个实例
	app1Cancel, err := builder.Register(context.Background(), &naming.Instance{
		AppID:    "app1",
		Hostname: "h1",
		Zone:     "z1",
	})

	app2Cancel, err := builder.Register(context.Background(), &naming.Instance{
		AppID:    "app2",
		Hostname: "h5",
		Zone:     "z3",
	})

	if err != nil {
		fmt.Println(err)
	}

	//服务发现app2
	app2 := builder.Build("app2")

	go func() {
		fmt.Println("节点列表")
		for {



			r1, _ := app1.Fetch(context.Background())
			if r1 != nil {
				app1_info:="app1 info:"
				for z, ins := range r1.Instances {

					app1_info+=fmt.Sprintf("zone:%s",z)
					for _, in := range ins {
						app1_info+=fmt.Sprintf(" app: %s host %s \n", in.AppID, in.Hostname)
					}
				}
			} else {
				fmt.Println("app1 empty")
			}
			fmt.Printf("\n")

			r2, _ := app2.Fetch(context.Background())
			if r2 != nil {
				app1_info:="app2 info:"
				for z, ins := range r2.Instances {
					app1_info+=fmt.Sprintf("zone:%s",z)
					for _, in := range ins {
						app1_info+=fmt.Sprintf(" app: %s host %s \n", in.AppID, in.Hostname)
					}
				}
			} else {
				fmt.Println("app2 empty")
			}
			fmt.Printf("\n")
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 5)
	fmt.Println("取消app1")
	app1Cancel()

	time.Sleep(time.Second * 10)
	fmt.Println("取消app2")
	app2Cancel()

	time.Sleep(30 * time.Second)
}
