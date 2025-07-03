package ioc

import (
	grpc2 "github.com/pluckhuang/goweb/aweb/follow/grpc"
	"github.com/pluckhuang/goweb/aweb/pkg/grpcx"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(l logger.LoggerV1, ecli *clientv3.Client, followService *grpc2.FollowServiceServer) *grpcx.Server {
	type Config struct {
		Port    int   `yaml:"port"`
		EtcdTTL int64 `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	followService.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       cfg.Port,
		Name:       "FollowService",
		L:          l,
		EtcdTTL:    cfg.EtcdTTL,
		EtcdClient: ecli,
	}
}
