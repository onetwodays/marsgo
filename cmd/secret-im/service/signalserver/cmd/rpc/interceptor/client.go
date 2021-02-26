package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"

)

func TimeInterceptor(ctx context.Context,
	                 method string,
	                 req, reply interface{},
                     cc *grpc.ClientConn,
                     invoker grpc.UnaryInvoker,
                     opts ...grpc.CallOption) error {
	stime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		return err
	}

	fmt.Printf("调用 %s 方法 耗时: %v\n", method, time.Now().Sub(stime))
	return nil
}