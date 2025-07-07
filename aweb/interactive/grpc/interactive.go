package grpc

import (
	"context"

	interactivev1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/interactive/v1"
	"github.com/pluckhuang/goweb/aweb/interactive/domain"
	"github.com/pluckhuang/goweb/aweb/interactive/service"
	"google.golang.org/grpc"
)

type InteractiveServiceServer struct {
	interactivev1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}

func (i *InteractiveServiceServer) Register(s *grpc.Server) {
	interactivev1.RegisterInteractiveServiceServer(s, i)
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *interactivev1.IncrReadCntRequest) (*interactivev1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &interactivev1.IncrReadCntResponse{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *interactivev1.LikeRequest) (*interactivev1.LikeResponse, error) {
	err := i.svc.Like(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &interactivev1.LikeResponse{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *interactivev1.CancelLikeRequest) (*interactivev1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &interactivev1.CancelLikeResponse{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *interactivev1.CollectRequest) (*interactivev1.CollectResponse, error) {
	err := i.svc.Collect(ctx, request.GetBiz(), request.GetBizId(),
		request.GetCid(), request.GetUid())
	return &interactivev1.CollectResponse{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *interactivev1.GetRequest) (*interactivev1.GetResponse, error) {
	intr, err := i.svc.Get(ctx, request.GetBiz(),
		request.GetBizId(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &interactivev1.GetResponse{
		Intr: i.toDTO(intr),
	}, nil
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *interactivev1.GetByIdsRequest) (*interactivev1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	intrs := make(map[int64]*interactivev1.Interactive, len(res))
	for k, v := range res {
		intrs[k] = i.toDTO(v)
	}
	return &interactivev1.GetByIdsResponse{
		Intrs: intrs,
	}, nil
}

func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *interactivev1.Interactive {
	return &interactivev1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		Liked:      intr.Liked,
		LikeCnt:    intr.LikeCnt,
	}
}
