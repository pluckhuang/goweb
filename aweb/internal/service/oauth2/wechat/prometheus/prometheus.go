package prometheus

import (
	"context"
	"time"

	"github.com/pluckhuang/goweb/aweb/internal/domain"
	"github.com/pluckhuang/goweb/aweb/internal/service/oauth2/wechat"
	"github.com/prometheus/client_golang/prometheus"
)

type Decorator struct {
	wechat.Service
	sum prometheus.Summary
}

func NewDecorator(svc wechat.Service, sum prometheus.Summary) *Decorator {
	return &Decorator{
		Service: svc,
		sum:     sum,
	}
}

func (d *Decorator) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		d.sum.Observe(float64(duration))
	}()
	return d.Service.VerifyCode(ctx, code)
}
