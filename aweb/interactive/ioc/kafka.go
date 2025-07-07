package ioc

import (
	"github.com/IBM/sarama"
	"github.com/pluckhuang/goweb/aweb/interactive/events"
	"github.com/pluckhuang/goweb/aweb/pkg/saramax"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
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

func InitConsumers(c1 *events.InteractiveReadEventConsumer, c2 *events.InteractiveSyncEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{c1, c2}
}
