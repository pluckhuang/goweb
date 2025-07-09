package test

import (
	"testing"

	feedv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/feed/v1"
	followMocks "github.com/pluckhuang/goweb/aweb/api/proto/gen/follow/v1/mocks"
	"github.com/pluckhuang/goweb/aweb/feed/grpc"
	"github.com/pluckhuang/goweb/aweb/feed/ioc"
	"github.com/pluckhuang/goweb/aweb/feed/repository"
	"github.com/pluckhuang/goweb/aweb/feed/repository/dao"
	"github.com/pluckhuang/goweb/aweb/feed/service"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func InitGrpcServer(t *testing.T) (feedv1.FeedSvcServer, *followMocks.MockFollowServiceClient, *gorm.DB) {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	feedPullEventDAO := dao.NewFeedPullEventDAO(db)
	feedPushEventDAO := dao.NewFeedPushEventDAO(db)
	feedEventRepo := repository.NewFeedEventRepo(feedPullEventDAO, feedPushEventDAO)
	mockCtrl := gomock.NewController(t)
	followClient := followMocks.NewMockFollowServiceClient(mockCtrl)
	v := ioc.RegisterHandler(feedEventRepo, followClient)
	feedService := service.NewFeedService(feedEventRepo, v)
	feedEventGrpcSvc := grpc.NewFeedEventGrpcSvc(feedService)
	return feedEventGrpcSvc, followClient, db
}
