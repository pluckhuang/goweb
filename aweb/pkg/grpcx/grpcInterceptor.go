package grpcx

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetMyGrpcClient 创建带有熔断、限流、重试功能的 gRPC 客户端
func GetMyGrpcClient(target string) (*grpc.ClientConn, error) {
	cfg := MiddlewareConfig{
		BreakerSettings: gobreaker.Settings{
			Name:        "MyServiceBreaker",
			MaxRequests: 5,                // 半开状态下允许的最大请求数
			Timeout:     30 * time.Second, // 熔断器从开启到半开的等待时间
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				// 考虑到 gRPC 重试机制，提高熔断阈值
				// 连续失败次数超过5次才触发熔断（避免与重试策略冲突）
				return counts.ConsecutiveFailures > 5
			},
			OnStateChange: func(name string, from, to gobreaker.State) {
				// 记录熔断器状态变化日志
				log.Printf("Circuit breaker %s changed from %s to %s", name, from, to)
			},
		},
		Limiter: rate.NewLimiter(rate.Every(100*time.Millisecond), 10), // 限流：每秒10个请求，突发允许10个
	}

	// 使用 grpc.NewClient 替代已废弃的 grpc.Dial
	conn, err := grpc.NewClient(
		target,
		grpc.WithDefaultServiceConfig(serviceConfig),              // 配置 gRPC 内置重试策略
		grpc.WithUnaryInterceptor(NewUnaryClientInterceptor(cfg)), // 应用自定义拦截器
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// 配置 gRPC 内置重试策略
// gRPC 会自动处理重试，拦截器只需要处理最终结果
var serviceConfig = `{
    "methodConfig": [{
        "name": [{"service": "MyService"}],
        "retryPolicy": {
            "maxAttempts": 4,                                    // 最多重试4次（包含首次请求）
            "initialBackoff": "0.1s",                           // 初始退避时间
            "maxBackoff": "1s",                                 // 最大退避时间
            "backoffMultiplier": 2,                             // 退避时间倍数
            "retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]  // 可重试的错误码
        }
    }]
}`

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	BreakerSettings gobreaker.Settings // 熔断器配置
	Limiter         *rate.Limiter      // 限流器
}

// NewUnaryClientInterceptor 创建统一的客户端拦截器
// 集成限流和熔断功能，与 gRPC 内置重试策略协调工作
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
		// 第一步：限流检查
		// 在请求发送前进行限流，避免对下游服务造成压力
		if limiter != nil {
			if err := limiter.Wait(ctx); err != nil {
				return status.Error(codes.ResourceExhausted, "rate limited")
			}
		}

		// 第二步：熔断器包装
		// 熔断器在 gRPC 重试之外工作，只对最终结果进行判断
		_, err := cb.Execute(func() (interface{}, error) {
			// 调用真实的 gRPC 方法
			// gRPC 内部会根据 serviceConfig 自动处理重试
			finalErr := invoker(ctx, method, req, reply, cc, opts...)

			// 对最终错误进行分类处理
			if finalErr != nil {
				st := status.Convert(finalErr)

				// 区分客户端错误和服务端错误
				// 客户端错误不应该触发熔断
				switch st.Code() {
				case codes.InvalidArgument, // 无效参数
					codes.NotFound,         // 资源未找到
					codes.PermissionDenied, // 权限拒绝
					codes.Unauthenticated:  // 未认证
					// 这些是客户端错误，不计入熔断器失败统计
					// 但仍然返回错误给调用方
					return nil, finalErr

				case codes.ResourceExhausted: // 资源耗尽（可能是限流）
					// 如果是我们自己的限流，不计入熔断统计
					if st.Message() == "rate limited" {
						return nil, finalErr
					}
					// 如果是服务端限流，计入熔断统计
					fallthrough

				default:
					// 服务端错误、网络错误等，计入熔断器失败统计
					// 包括：Unavailable, DeadlineExceeded, Internal, etc.
					return nil, finalErr
				}
			}

			// 成功的请求
			return nil, nil
		})

		// 处理熔断器错误
		if err != nil {
			if errors.Is(err, gobreaker.ErrOpenState) {
				// 熔断器开启时，返回标准的 gRPC 错误
				return status.Error(codes.Unavailable, "service circuit breaker is open")
			}
			if errors.Is(err, gobreaker.ErrTooManyRequests) {
				// 半开状态下请求过多时
				return status.Error(codes.ResourceExhausted, "service circuit breaker: too many requests")
			}
		}

		return err
	}
}

// GetCircuitBreakerState 获取熔断器状态（可选的辅助方法）
func GetCircuitBreakerState(cb *gobreaker.CircuitBreaker) gobreaker.State {
	return cb.State()
}

// IsCircuitBreakerOpen 检查熔断器是否开启（可选的辅助方法）
func IsCircuitBreakerOpen(cb *gobreaker.CircuitBreaker) bool {
	return cb.State() == gobreaker.StateOpen
}
