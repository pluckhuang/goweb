package ioc

import (
	"fmt"

	"github.com/pluckhuang/goweb/aweb/feed/service"
	"github.com/spf13/viper"
)

func InitGlobalVal() *service.ArticleEventConfig {
	var cfg service.ArticleEventConfig
	err := viper.UnmarshalKey("service", &cfg)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败: %w", err))
	}
	return &cfg
}
