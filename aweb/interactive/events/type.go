package events

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/pluckhuang/goweb/aweb/interactive/repository"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/pluckhuang/goweb/aweb/pkg/saramax"
)

const topicReadEvent = "article_read_event"

type ReadEvent struct {
	Aid int64
	Uid int64
}

var _ saramax.Consumer = &InteractiveReadEventConsumer{}

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveReadEventConsumer(
	client sarama.Client,
	l logger.LoggerV1,
	repo repository.InteractiveRepository) *InteractiveReadEventConsumer {
	ic := &InteractiveReadEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
	return ic
}

func (r *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicReadEvent},
			saramax.NewHandler[ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("InteractiveReadEvent 退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *InteractiveReadEventConsumer) StartBatch() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicReadEvent},
			saramax.NewBatchHandler[ReadEvent](r.l, r.BatchConsume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (r *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := r.repo.IncrReadCnt(ctx, "article", evt.Aid)
	return err
}

func (r *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	evts []ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	bizs := make([]string, 0, len(msgs))
	ids := make([]int64, 0, len(msgs))
	for _, evt := range evts {
		bizs = append(bizs, "article")
		ids = append(ids, evt.Uid)
	}
	return r.repo.BatchIncrReadCnt(ctx, bizs, ids)
}
