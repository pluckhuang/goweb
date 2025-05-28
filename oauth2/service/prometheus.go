package service

import (
	"context"
	"time"

	"github.com/pluckhuang/goweb/oauth2/domain"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusDecorator 利用组合来避免需要实现所有的接口
type PrometheusDecorator struct {
	Oauth2Service
	sum prometheus.Summary
}

func NewPrometheusDecorator(svc Oauth2Service,
	namespace string,
	subsystem string,
	instanceId string,
	name string) *PrometheusDecorator {
	sum := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      name,
		Namespace: namespace,
		Subsystem: subsystem,
		ConstLabels: map[string]string{
			"instance_id": instanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.95:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	prometheus.MustRegister(sum)
	return &PrometheusDecorator{
		Oauth2Service: svc,
		sum:           sum,
	}
}

func (p *PrometheusDecorator) HandleCallback(ctx context.Context, platform string, code string, state string) (domain.Oauth2Info, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		p.sum.Observe(float64(duration.Milliseconds()))
	}()
	return p.Oauth2Service.HandleCallback(ctx, platform, code, state)
}
