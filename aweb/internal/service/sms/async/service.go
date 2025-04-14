package async

import (
	"context"
	"time"

	"github.com/pluckhuang/goweb/aweb/internal/domain"
	"github.com/pluckhuang/goweb/aweb/internal/repository"
	"github.com/pluckhuang/goweb/aweb/internal/service/sms"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
)

type Service struct {
	svc sms.Service
	// 转异步，存储发短信请求的 repository
	repo repository.AsyncSmsRepository
	l    logger.LoggerV1
}

func NewService(svc sms.Service,
	repo repository.AsyncSmsRepository,
	l logger.LoggerV1) *Service {
	res := &Service{
		svc:  svc,
		repo: repo,
		l:    l,
	}
	go func() {
		res.StartAsyncCycle()
	}()
	return res
}

// StartAsyncCycle 异步发送消息
// 这里我们没有设计退出机制，是因为没啥必要
// 因为程序停止的时候，它自然就停止了
// 原理：这是最简单的抢占式调度
func (s *Service) StartAsyncCycle() {
	// 这个是我为了测试而引入的，防止你在运行测试的时候，会出现偶发性的失败
	time.Sleep(time.Second * 3)
	for {
		s.AsyncSend()
	}
}

func (s *Service) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 抢占一个异步发送的消息，确保在非常多个实例
	// 比如 k8s 部署了三个 pod，一个请求，只有一个实例能拿到
	as, err := s.repo.PreemptWaitingSMS(ctx)
	cancel()
	switch err {
	case nil:
		// 执行发送
		// 这个也可以做成配置的
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = s.svc.Send(ctx, as.TplId, as.Args, as.Numbers...)
		if err != nil {
			// 啥也不需要干
			s.l.Error("执行异步发送短信失败",
				logger.Error(err),
				logger.Int64("id", as.Id))
		}
		res := err == nil
		// 通知 repository 我这一次的执行结果
		err = s.repo.ReportScheduleResult(ctx, as.Id, res)
		if err != nil {
			s.l.Error("执行异步发送短信成功，但是标记数据库失败",
				logger.Error(err),
				logger.Bool("res", res),
				logger.Int64("id", as.Id))
		}
	case repository.ErrWaitingSMSNotFound:
		time.Sleep(time.Second)
	default:
		s.l.Error("抢占异步发送短信任务失败",
			logger.Error(err))
		time.Sleep(time.Second)
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	if s.needAsync() {
		// 需要异步发送，直接转储到数据库
		err := s.repo.Add(ctx, domain.AsyncSms{
			TplId:   tplId,
			Args:    args,
			Numbers: numbers,
			// 设置可以重试三次
			RetryMax: 3,
		})
		return err
	}
	return s.svc.Send(ctx, tplId, args, numbers...)
}

func (s *Service) needAsync() bool {
	// todo
	return true
}
