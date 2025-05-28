package ioc

import (
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/pluckhuang/goweb/oauth2/service"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitPrometheus(cmd redis.Cmdable, logv1 logger.LoggerV1) service.Oauth2Service {
	handlerMap := InitService(cmd, logv1)
	type Config struct {
		NameSpace  string `yaml:"nameSpace"`
		Subsystem  string `yaml:"subsystem"`
		InstanceID string `yaml:"instanceId"`
		Name       string `yaml:"name"`
	}
	var cfg Config
	err := viper.UnmarshalKey("prometheus", &cfg)
	if err != nil {
		panic(err)
	}
	oauth2service := service.NewOauth2Service(handlerMap)
	return service.NewPrometheusDecorator(oauth2service, cfg.NameSpace, cfg.Subsystem, cfg.InstanceID, cfg.Name)
}
