package grpc

import (
	"context"

	searchv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/search/v1"
	"github.com/pluckhuang/goweb/aweb/search/domain"
	"github.com/pluckhuang/goweb/aweb/search/service"
	"google.golang.org/grpc"
)

type SyncServiceServer struct {
	searchv1.UnimplementedSyncServiceServer
	syncSvc service.SyncService
}

func NewSyncServiceServer(syncSvc service.SyncService) *SyncServiceServer {
	return &SyncServiceServer{
		syncSvc: syncSvc,
	}
}

func (s *SyncServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSyncServiceServer(server, s)
}

func (s *SyncServiceServer) InputArticle(ctx context.Context, request *searchv1.InputArticleRequest) (*searchv1.InputArticleResponse, error) {
	err := s.syncSvc.InputArticle(ctx, s.toDomainArticle(request.GetArticle()))
	return &searchv1.InputArticleResponse{}, err
}

func (s *SyncServiceServer) InputAny(ctx context.Context, req *searchv1.InputAnyRequest) (*searchv1.InputAnyResponse, error) {
	err := s.syncSvc.InputAny(ctx, req.IndexName, req.DocId, req.Data)
	return &searchv1.InputAnyResponse{}, err
}

func (s *SyncServiceServer) toDomainArticle(art *searchv1.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  art.Status,
		Content: art.Content,
	}
}
