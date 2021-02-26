package interceptor

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)
var limiter = rate.NewLimiter(rate.Limit(100), 100)
func RateLimitInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if !limiter.Allow() {
		fmt.Println("限流了")
		return nil, nil
	}
	return handler(ctx, req)
}