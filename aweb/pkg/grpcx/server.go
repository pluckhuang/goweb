package grpcx

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/pluckhuang/goweb/aweb/pkg/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Port int
	// ETCD 服务注册租约 TTL
	EtcdTTL     int64
	EtcdClient  *clientv3.Client
	etcdManager endpoints.Manager
	etcdKey     string
	cancel      func()
	Name        string
	L           logger.LoggerV1
}

// Serve 启动服务器并且阻塞
func (s *Server) Serve() error {
	// 初始化一个控制整个过程的 ctx
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	port := strconv.Itoa(s.Port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	// 要先确保启动成功，再注册服务
	err = s.register(ctx, port)
	if err != nil {
		return err
	}
	return s.Server.Serve(l)
}

func (s *Server) register(ctx context.Context, port string) error {
	cli := s.EtcdClient
	serviceName := "service/" + s.Name
	em, err := endpoints.NewManager(cli,
		serviceName)
	if err != nil {
		return err
	}
	s.etcdManager = em
	ip := netx.GetOutboundIP()
	s.etcdKey = serviceName + "/" + ip
	addr := ip + ":" + port
	leaseResp, err := cli.Grant(ctx, s.EtcdTTL)
	// 开启续约
	ch, err := cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}
	go func() {
		// 可以预期，当我们的 cancel 被调用的时候，就会退出这个循环
		for chResp := range ch {
			s.L.Debug("续约：", logger.String("resp", chResp.String()))
		}
	}()
	return em.AddEndpoint(ctx, s.etcdKey,
		endpoints.Endpoint{Addr: addr}, clientv3.WithLease(leaseResp.ID))
}

func (s *Server) Close() error {
	s.cancel()
	if s.etcdManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.etcdManager.DeleteEndpoint(ctx, s.etcdKey)
		if err != nil {
			return err
		}
	}
	err := s.EtcdClient.Close()
	if err != nil {
		return err
	}
	s.Server.GracefulStop()
	return nil
}
