package repository

import (
	"context"

	"github.com/pluckhuang/goweb/aweb/search/repository/dao"
)

type anyRepository struct {
	dao dao.AnyDAO
}

func NewAnyRepository(dao dao.AnyDAO) AnyRepository {
	return &anyRepository{dao: dao}
}

func (repo *anyRepository) Input(ctx context.Context, index string, docID string, data string) error {
	return repo.dao.Input(ctx, index, docID, data)
}

func (repo *anyRepository) Delete(ctx context.Context, index string, docID string) error {
	return repo.dao.Delete(ctx, index, docID)
}
