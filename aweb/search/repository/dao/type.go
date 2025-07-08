package dao

import (
	"context"
)

const ArticleIndexName = "article_index"
const CollectIndexName = "collect_index"
const LikeIndexName = "like_index"

type ArticleDAO interface {
	InputArticle(ctx context.Context, article Article) error
	// Search artIds 命中了索引的 article id
	Search(ctx context.Context, req SearchReq, keywords []string) ([]Article, error)
}

type LikeDAO interface {
	Search(ctx context.Context, uid int64, biz string) ([]int64, error)
}

type CollectDAO interface {
	Search(ctx context.Context, uid int64, biz string) ([]int64, error)
}

type AnyDAO interface {
	Input(ctx context.Context, index, docID, data string) error
	Delete(ctx context.Context, index string, docID string) error
}

type SearchReq struct {
	LikeIds    []int64
	CollectIds []int64
}

type Biz struct {
	Uid   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
}
