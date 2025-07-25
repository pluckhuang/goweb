package saramax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
)

type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	l  logger.LoggerV1
}

func NewBatchHandler[T any](l logger.LoggerV1, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{fn: fn, l: l}
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10

	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var done = false

		// 收集消息，直到达到批量大小或超时
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 超时退出
				done = true
			case msg, ok := <-msgs:
				if !ok {
					// 消息通道关闭，说明分区重新平衡或消费结束
					cancel()
					return nil
				}
				batch = append(batch, msg)

				// 反序列化消息
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列消息体失败",
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset),
						logger.Error(err))
					continue
				}
				ts = append(ts, t)
			}
		}
		cancel()

		// 如果没有收集到任何消息，直接跳过处理
		if len(batch) == 0 {
			continue
		}

		// 批量处理消息
		err := b.fn(batch, ts)
		if err != nil {
			// 记录批量处理错误
			b.l.Error("处理消息失败",
				logger.Error(err))
			// 返回错误，通知 Sarama 消费失败
			return err
		}

		// 标记消息为已消费
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
