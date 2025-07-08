package repository

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/search/domain"
)

type ArticleRepository interface {
	InputArticle(ctx context.Context, msg domain.Article) error
	SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error)
}

type AnyRepository interface {
	Input(ctx context.Context, index string, docID string, data string) error
	Delete(ctx context.Context, index, docID string) error
}
