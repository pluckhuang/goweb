package service

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/search/domain"
	"github.com/pluckhuang/goweb/aweb/search/repository"
)

type SyncService interface {
	InputArticle(ctx context.Context, article domain.Article) error
	InputAny(ctx context.Context, idxName, docID, data string) error
	Delete(ctx context.Context, index, docId string) error
}

type syncService struct {
	articleRepo repository.ArticleRepository
	anyRepo     repository.AnyRepository
}

func (s *syncService) Delete(ctx context.Context, index, docId string) error {
	return s.anyRepo.Delete(ctx, index, docId)
}

func (s *syncService) InputArticle(ctx context.Context, article domain.Article) error {
	return s.articleRepo.InputArticle(ctx, article)
}

func (s *syncService) InputAny(ctx context.Context, index, docID, data string) error {
	return s.anyRepo.Input(ctx, index, docID, data)
}

func NewSyncService(
	anyRepo repository.AnyRepository,
	articleRepo repository.ArticleRepository) SyncService {
	return &syncService{
		articleRepo: articleRepo,
		anyRepo:     anyRepo,
	}
}
