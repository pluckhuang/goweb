package ioc

import (
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/pluckhuang/goweb/aweb/search/repository/dao"
	"github.com/spf13/viper"
)

// InitESClient 读取配置文件，进行初始化ES客户端
func InitESClient() *elastic.Client {
	type Config struct {
		Url   string `yaml:"url"`
		Sniff bool   `yaml:"sniff"`
	}
	var cfg Config
	err := viper.UnmarshalKey("es", &cfg)
	if err != nil {
		panic(fmt.Errorf("读取 ES 配置失败 %w", err))
	}
	const timeout = 10 * time.Second
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(cfg.Url),
		elastic.SetSniff(cfg.Sniff),
		elastic.SetHealthcheckTimeoutStartup(timeout),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	err = dao.InitES(client)
	if err != nil {
		panic(err)
	}
	return client
}
