package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
}

func InitRlockClient(client redis.Cmdable) *rlock.Client {
	return rlock.NewClient(client)
}

func InitRedisV1() redis.Cmdable {
	v := viper.New()
	v.SetConfigType("conf")
	v.SetConfigFile("config/redis.conf")
	addr := v.GetString("addr")
	return redis.NewClient(&redis.Options{
		//Addr: viper.GetString("redis.addr"),
		Addr: addr,
	})
}
