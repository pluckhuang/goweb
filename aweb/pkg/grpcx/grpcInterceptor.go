package grpcx

import (
	"context"
	"math"
	"time"

	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MiddlewareConfig 配置
type MiddlewareConfig struct {
	BreakerSettings gobreaker.Settings
	Limiter         *rate.Limiter
	RetryMax        int
	RetryBaseDelay  time.Duration
	RetryMaxDelay   time.Duration
}

// NewUnaryClientInterceptor 返回统一的拦截器
func NewUnaryClientInterceptor(cfg MiddlewareConfig) grpc.UnaryClientInterceptor {
	cb := gobreaker.NewCircuitBreaker(cfg.BreakerSettings)
	limiter := cfg.Limiter

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// 限流
		if limiter != nil {
			if err := limiter.Wait(ctx); err != nil {
				return status.Error(codes.ResourceExhausted, "rate limited")
			}
		}
		var lastErr error
		for i := 0; i <= cfg.RetryMax; i++ {
			_, err := cb.Execute(func() (interface{}, error) {
				return nil, invoker(ctx, method, req, reply, cc, opts...)
			})
			if err == nil {
				return nil
			}
			lastErr = err
			// 只对可重试错误重试
			st, _ := status.FromError(err)
			if st != nil && st.Code() != codes.Unavailable && st.Code() != codes.DeadlineExceeded {
				break
			}
			// 指数退避
			sleep := time.Duration(math.Min(
				float64(cfg.RetryBaseDelay)*math.Pow(2, float64(i)),
				float64(cfg.RetryMaxDelay),
			))
			time.Sleep(sleep)
		}
		return lastErr
	}
}
