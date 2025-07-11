package dao

import (
	"context"
	_ "embed"
	"time"

	"github.com/olivere/elastic/v7"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed article_index.json
	articleIndex string
	//go:embed like_index.json
	likeIndex string
	//go:embed collect_index.json
	collectIndex string
)

// InitES 创建索引
func InitES(client *elastic.Client) error {
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var eg errgroup.Group
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, ArticleIndexName, articleIndex)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, LikeIndexName, likeIndex)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, CollectIndexName, collectIndex)
	})

	return eg.Wait()
}

func tryCreateIndex(ctx context.Context,
	client *elastic.Client,
	idxName, idxCfg string,
) error {
	// 索引可能已经建好了
	ok, err := client.IndexExists(idxName).Do(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	_, err = client.CreateIndex(idxName).Body(idxCfg).Do(ctx)
	return err
}
