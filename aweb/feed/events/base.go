package events

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/pluckhuang/goweb/aweb/feed/domain"
	"github.com/pluckhuang/goweb/aweb/feed/service"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/pluckhuang/goweb/aweb/pkg/saramax"
)

const topicFeedEvent = "feed_event"

// FeedEvent 通用业务事件
type FeedEvent struct {
	Type     string
	Metadata map[string]string
}

type FeedEventConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.FeedService
}

func NewFeedEventConsumer(
	client sarama.Client,
	l logger.LoggerV1,
	svc service.FeedService) *FeedEventConsumer {
	return &FeedEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
}

func (r *FeedEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("feed_event",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicFeedEvent},
			saramax.NewHandler[FeedEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (r *FeedEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt FeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: evt.Type,
		Ext:  evt.Metadata,
	})
}
