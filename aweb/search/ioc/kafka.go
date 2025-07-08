package ioc

import (
	"github.com/IBM/sarama"
	"github.com/pluckhuang/goweb/aweb/search/events"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
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

// NewConsumers 所有的 Consumer 要在这里注册一下
func NewConsumers(articleConsumer *events.ArticleConsumer, interactiveConsumer *events.InteractiveConsumer) []events.Consumer {
	return []events.Consumer{
		articleConsumer,
		interactiveConsumer,
	}
}
