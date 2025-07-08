package grpc

import (
	"context"

	"github.com/ecodeclub/ekit/slice"
	searchv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/search/v1"
	"github.com/pluckhuang/goweb/aweb/search/domain"
	"github.com/pluckhuang/goweb/aweb/search/service"
	"google.golang.org/grpc"
)

type SearchServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	svc service.SearchService
}

func NewSearchService(svc service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{svc: svc}
}

func (s *SearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func (s *SearchServiceServer) Search(ctx context.Context, request *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	resp, err := s.svc.Search(ctx, request.Uid, request.Expression)
	if err != nil {
		return nil, err
	}
	return &searchv1.SearchResponse{
		Article: &searchv1.ArticleResult{
			Articles: slice.Map(resp.Articles, func(idx int, src domain.Article) *searchv1.Article {
				return &searchv1.Article{
					Id:      src.Id,
					Title:   src.Title,
					Status:  src.Status,
					Content: src.Content,
				}
			}),
		},
	}, nil
}
