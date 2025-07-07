package events

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"localhost:9094"}, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSaramaSyncProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}

func initProducer() InteractiveProducer {
	// Inline the Sarama client and producer initialization to avoid import cycle
	client := InitSaramaClient()
	syncProducer := InitSaramaSyncProducer(client)
	return NewInteractiveProducer(syncProducer)
}

func TestInteractiveProducer_EndToEnd(t *testing.T) {
	interactiveProducer := initProducer()
	evt := InteractiveEvent{
		Type:  CollectEventType,
		Biz:   "article",
		BizId: 1,
		Uid:   1,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := interactiveProducer.ProduceInteractiveEvent(ctx, evt)
	if err != nil {
		t.Fatalf("Failed to produce like event: %v", err)
	}
}
