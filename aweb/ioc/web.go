package ioc

import (
	"context"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pluckhuang/goweb/aweb/internal/web"
	ijwt "github.com/pluckhuang/goweb/aweb/internal/web/jwt"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/pluckhuang/goweb/aweb/internal/web/middleware"
	"github.com/pluckhuang/goweb/aweb/pkg/ginx"
	"github.com/pluckhuang/goweb/aweb/pkg/ginx/middleware/prometheus"
	"github.com/pluckhuang/goweb/aweb/pkg/ginx/middleware/ratelimit"
	"github.com/pluckhuang/goweb/aweb/pkg/limiter"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

//func InitWebServerV1(mdls []gin.HandlerFunc, hdls []web.Handler) *gin.Engine {
//	server := gin.Default()
//	server.Use(mdls...)
//	for _, hdl := range hdls {
//		hdl.RegisterRoutes(server)
//	}
//	//userHdl.RegisterRoutes(server)
//	return server
//}

// func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
// 	server := gin.Default()
// 	server.Use(mdls...)
// 	userHdl.RegisterRoutes(server)
// 	return server
// }

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

// 为什么每一个中间件 Build() 返回一个带有 ctx 参数的匿名函数？
// Gin 的中间件规范要求每个中间件都是如下形式的函数：
// func(ctx *gin.Context)
// 也就是说，Gin 在处理每个 HTTP 请求时，会依次调用所有注册的中间件，每个中间件都需要接收当前请求的上下文（ctx），以便：

// 读取请求信息（如 IP、Header、Path 等）
// 执行中间件逻辑
// 决定是否继续处理请求（ctx.Abort() 或 ctx.Next()）
// 返回错误响应或通过请求
// 因此，Build() 的最终产物必须是一个带有 ctx *gin.Context 参数的函数，这样才能被 Gin 框架识别和调用，完成中间件逻辑。
// -------------------------
// ctx 包含了什么
// 在 Gin 框架中，ctx 指的是 *gin.Context 类型的对象。它是每个 HTTP 请求的上下文，包含了请求处理过程中所需的所有信息和操作方法。具体包括：

// 请求信息

// 请求的 URL、方法、Header、Body、参数等。
// 客户端 IP、Cookie、Query 参数等。
// 响应操作

// 设置响应状态码、Header、返回 JSON、字符串、文件等。
// 终止请求（Abort）、重定向等。
// 中间件控制

// 控制请求流程（如 Next() 继续下一个中间件，Abort() 停止后续处理）。
// 在中间件间传递数据（Set/Get）。
// 上下文扩展

// 可以挂载自定义数据，方便业务逻辑和中间件之间共享信息。
// 错误处理

// 记录和处理请求中的错误。
// 总结：
// ctx 是 Gin 处理每个 HTTP 请求时的核心对象，包含了请求、响应、流程控制和数据共享等所有上下文信息。

func InitGinMiddlewares(redisClient redis.Cmdable, hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	pb := &prometheus.Builder{
		Namespace: "pluckh.com",
		Subsystem: "aweb",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}
	ginx.InitCounter(prometheus2.CounterOpts{
		Namespace: "pluckh.com",
		Subsystem: "aweb",
		Name:      "biz_code",
		Help:      "统计业务错误码",
	})
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,

			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token"},
			//AllowHeaders: []string{"content-type"},
			//AllowMethods: []string{"POST"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		func(ctx *gin.Context) {
			println("这是我的 Middleware")
		},
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		otelgin.Middleware("aweb"),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
			l.Debug("", logger.Field{Key: "req", Val: al})
		}).AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
	}
}
