package ioc

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitOTEL 返回一个关闭函数，并且让调用者关闭的时候来决定这个 ctx
func InitOTEL() func(ctx context.Context) {
	res, err := newResource("aweb", "v0.0.1")
	if err != nil {
		panic(err)
	}
	prop := newPropagator()
	// 在客户端和服务端之间传递 tracing 的相关信息
	otel.SetTextMapPropagator(prop)

	// 初始化 trace provider
	// 这个 provider 就是用来在打点的时候构建 trace 的
	tp, err := newTraceProvider(res)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	return func(ctx context.Context) {
		_ = tp.Shutdown(ctx)
	}
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	batchProcessor := trace.NewBatchSpanProcessor(
		exporter,
		trace.WithBatchTimeout(5*time.Second), // 批量导出的超时时间
		trace.WithMaxExportBatchSize(100),     // 每次导出的最大 Span 数量
	)
	filteringProcessor := NewFilteringSpanProcessor(batchProcessor)
	// 创建 TracerProvider，并注册 FilteringSpanProcessor
	traceProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(filteringProcessor), // 确保它是第一个处理器
		trace.WithResource(res),
	)

	return traceProvider, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
