package events

import (
	"context"
	"errors"
	"time"

	"github.com/IBM/sarama"
	"github.com/pluckhuang/goweb/aweb/interactive/repository"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/pluckhuang/goweb/aweb/pkg/saramax"
)

var _ saramax.Consumer = &InteractiveSyncEventConsumer{}

type InteractiveSyncEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveSyncEventConsumer(
	client sarama.Client,
	l logger.LoggerV1,
	repo repository.InteractiveRepository) *InteractiveSyncEventConsumer {
	ic := &InteractiveSyncEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
	return ic
}

func (r *InteractiveSyncEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("general_interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{interactiveTopic},
			saramax.NewHandler[InteractiveEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("InteractiveSyncEvent 退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *InteractiveSyncEventConsumer) Consume(msg *sarama.ConsumerMessage,
	event InteractiveEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch event.Type {
	case LikeEventType:
		return r.repo.IncrLike(ctx, event.Biz, event.BizId, event.Uid)
	case CancelLikeEventType:
		return r.repo.DecrLike(ctx, event.Biz, event.BizId, event.Uid)
	case CollectEventType:
		return r.repo.AddCollectionItem(ctx, event.Biz, event.BizId, 0, event.Uid)
	default:
		r.l.Warn("未知的互动事件类型", logger.Int64("type: ", int64(event.Type)))
		return errors.New("unknown interactive event type")
	}
}
