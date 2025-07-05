package grpc

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/article/domain"
	"github.com/pluckhuang/goweb/aweb/article/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	articlev1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/article/v1"
)

type ArticleServiceServer struct {
	// 组合
	articlev1.UnimplementedArticleServiceServer
	svc service.ArticleService
}

func NewGrpcServer(svc service.ArticleService) *ArticleServiceServer {
	return &ArticleServiceServer{
		svc: svc,
	}
}
func (c *ArticleServiceServer) Register(server grpc.ServiceRegistrar) {
	articlev1.RegisterArticleServiceServer(server, c)
}

func (c *ArticleServiceServer) Save(ctx context.Context, request *articlev1.SaveRequest) (*articlev1.SaveResponse, error) {
	art := convertToDomain(request.GetArticle())
	id, err := c.svc.Save(ctx, art)
	if err != nil {
		return nil, err
	}
	return &articlev1.SaveResponse{Id: id}, nil
}

func convertToDomain(art *articlev1.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
	}
}

func (c *ArticleServiceServer) Publish(ctx context.Context, request *articlev1.PublishRequest) (*articlev1.PublishResponse, error) {
	art := convertToDomain(request.GetArticle())
	id, err := c.svc.Publish(ctx, art)
	if err != nil {
		return nil, err
	}
	return &articlev1.PublishResponse{Id: id}, nil
}

func (c *ArticleServiceServer) Withdraw(ctx context.Context, request *articlev1.WithdrawRequest) (*articlev1.WithdrawResponse, error) {
	err := c.svc.Withdraw(ctx, request.Uid, request.Id)
	if err != nil {
		return nil, err
	}
	return &articlev1.WithdrawResponse{}, nil
}

func (c *ArticleServiceServer) GetByAuthor(ctx context.Context, request *articlev1.GetByAuthorRequest) (*articlev1.GetByAuthorResponse, error) {
	arts, err := c.svc.GetByAuthor(ctx, request.Uid, int(request.Offset), int(request.Limit))
	if err != nil {
		return nil, err
	}
	return &articlev1.GetByAuthorResponse{
		Articles: convertToProtoList(arts),
	}, nil
}

func convertToProto(art domain.Article) *articlev1.Article {
	return &articlev1.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   int32(art.Status.ToUint8()),
		Ctime:    timestamppb.New(art.Ctime),
		Utime:    timestamppb.New(art.Utime),
	}
}

func convertToProtoList(arts []domain.Article) []*articlev1.Article {
	resp := make([]*articlev1.Article, 0, len(arts))
	for _, art := range arts {
		rpcArticle := convertToProto(art)
		resp = append(resp, rpcArticle)
	}
	return resp
}

func (c *ArticleServiceServer) GetById(ctx context.Context, request *articlev1.GetByIdRequest) (*articlev1.GetByIdResponse, error) {
	art, err := c.svc.GetById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &articlev1.GetByIdResponse{Article: convertToProto(art)}, nil
}

func (c *ArticleServiceServer) GetPubById(ctx context.Context, request *articlev1.GetPubByIdRequest) (*articlev1.GetPubByIdResponse, error) {
	art, err := c.svc.GetPubById(ctx, request.Id, request.Uid)
	if err != nil {
		return nil, err
	}
	return &articlev1.GetPubByIdResponse{Article: convertToProto(art)}, nil
}

func (c *ArticleServiceServer) ListPub(ctx context.Context, request *articlev1.ListPubRequest) (*articlev1.ListPubResponse, error) {
	arts, err := c.svc.ListPub(ctx, request.GetStart().AsTime(), int(request.Offset), int(request.Limit))
	if err != nil {
		return nil, err
	}
	resp := &articlev1.ListPubResponse{}
	for _, art := range arts {
		resp.Articles = append(resp.Articles, convertToProto(art))
	}
	return resp, nil
}
