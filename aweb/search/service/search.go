package service

import (
	"context"
	"strings"

	"github.com/pluckhuang/goweb/aweb/search/domain"
	"github.com/pluckhuang/goweb/aweb/search/repository"
	"golang.org/x/sync/errgroup"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type searchService struct {
	articleRepo repository.ArticleRepository
}

func NewSearchService(articleRepo repository.ArticleRepository) SearchService {
	return &searchService{articleRepo: articleRepo}
}

func (s *searchService) Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error) {
	// 要对 expression 进行解析，生成查询计划
	// 输入预处理
	// 清除掉空格，切割;',.
	keywords := strings.Split(expression, " ")
	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, uid, keywords)
		res.Articles = arts
		return err
	})
	return res, eg.Wait()
}
