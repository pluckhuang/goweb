package ioc

import (
	followv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/follow/v1"
	"github.com/pluckhuang/goweb/aweb/feed/repository"
	"github.com/pluckhuang/goweb/aweb/feed/service"
)

func RegisterHandler(repo repository.FeedEventRepo, followClient followv1.FollowServiceClient) map[string]service.Handler {
	articleHandler := service.NewArticleEventHandler(repo, followClient)
	followHandler := service.NewFollowEventHandler(repo)
	likeHandler := service.NewLikeEventHandler(repo)
	return map[string]service.Handler{
		service.ArticleEventName: articleHandler,
		service.FollowEventName:  followHandler,
		service.LikeEventName:    likeHandler,
	}
}
