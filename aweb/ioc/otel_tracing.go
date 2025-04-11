package ioc

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/sdk/trace"
)

var excludedSpans = []string{"gorm.Row", "gorm.Query"}

func shouldExcludeSpan(spanName string) bool {
	for _, excluded := range excludedSpans {
		if strings.Contains(spanName, excluded) {
			return true
		}
	}
	return false
}

type FilteringSpanProcessor struct {
	processor trace.SpanProcessor
}

func NewFilteringSpanProcessor(processor trace.SpanProcessor) *FilteringSpanProcessor {
	return &FilteringSpanProcessor{processor: processor}
}

func (fsp *FilteringSpanProcessor) OnStart(ctx context.Context, s trace.ReadWriteSpan) {
	// fmt.Printf("OnStart: Span Name: %s:%s\n", s.Name(), s.SpanContext().SpanID())
	if shouldExcludeSpan(s.Name()) {
		// fmt.Printf("Excluded Span: %s\n", s.Name())
		return
	}
	fsp.processor.OnStart(ctx, s)
}

func (fsp *FilteringSpanProcessor) OnEnd(s trace.ReadOnlySpan) {
	// 打印 Span 的名称和上下文信息
	// fmt.Printf("OnEnd: Span Name: %s:%s\n", s.Name(), s.SpanContext().SpanID())
	if shouldExcludeSpan(s.Name()) {
		// fmt.Printf("Excluded Span: %s\n", s.Name())
		return
	}
	fsp.processor.OnEnd(s)
}

func (fsp *FilteringSpanProcessor) Shutdown(ctx context.Context) error {
	return fsp.processor.Shutdown(ctx)
}

func (fsp *FilteringSpanProcessor) ForceFlush(ctx context.Context) error {
	return fsp.processor.ForceFlush(ctx)
}
