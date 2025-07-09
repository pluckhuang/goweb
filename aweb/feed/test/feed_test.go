package test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/pluckhuang/goweb/aweb/feed/repository/dao"
	"github.com/pluckhuang/goweb/aweb/feed/service"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	feedv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/feed/v1"
	followv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/follow/v1"
	followv1Mock "github.com/pluckhuang/goweb/aweb/api/proto/gen/follow/v1/mocks"
)

type FeedTestSuite struct {
	suite.Suite
}

func (f *FeedTestSuite) SetupSuite() {
	// 初始化配置文件
	viper.SetConfigFile("config.yaml")
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func (f *FeedTestSuite) Test_Feed() {
	// 初始化
	server, mockFollowClient, db := InitGrpcServer(f.T())
	defer func() {
		db.Table("feed_push_events").Where("id > ? ", 0).Delete(&dao.FeedPushEvent{})
		db.Table("feed_pull_events").Where("id > ? ", 0).Delete(&dao.FeedPullEvent{})
	}()
	// 设置followmock的值
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Minute)
	defer cancel()

	// 创建事件
	err := f.setupEvent(ctx, mockFollowClient, server)
	require.NoError(f.T(), err)

	// 获取feed流事件
	wantEvents := f.getFeedEventWant(ctx, mockFollowClient, server)

	resp, err := server.FindFeedEvents(ctx, &feedv1.FindFeedEventsRequest{
		Uid:       1,
		Limit:     20,
		Timestamp: time.Now().Unix() + 10,
	})
	require.NoError(f.T(), err)
	assert.Equal(f.T(), len(wantEvents), len(resp.FeedEvents))
	checkerMap := map[string]EventCheck{
		service.ArticleEventName: ArticleEvent{},
		service.LikeEventName:    LikeEvent{},
		service.FollowEventName:  FollowEvent{},
	}
	for i := 0; i < len(wantEvents); i++ {
		wantEvent, actualEvent := wantEvents[i], resp.FeedEvents[i]
		checker := checkerMap[wantEvent.Type]
		wantContent, actualContent := checker.Check(wantEvent.Content, actualEvent.Content)
		assert.Equal(f.T(), wantContent, actualContent)
	}
}

func (f *FeedTestSuite) setupEvent(ctx context.Context, mockFollowClient *followv1Mock.MockFollowServiceClient, server feedv1.FeedSvcServer) error {
	// 发表文章事件：用户2发表了4篇文章,用户3发表了3篇文章
	articleEvents := []ArticleEvent{
		{
			Uid:   "2",
			Aid:   "1",
			Title: "用户2发表了文章1",
		},
		{
			Uid:   "2",
			Aid:   "2",
			Title: "用户2发表了文章2",
		},
		{
			Uid:   "2",
			Aid:   "3",
			Title: "用户2发表了文章3",
		},
		{
			Uid:   "2",
			Aid:   "4",
			Title: "用户2发表了文章4",
		},
	}

	// 设置用户2的关注统计, 用户2有5个粉丝, 大于阈值4 , 是个大v
	mockFollowClient.EXPECT().GetFollowStatics(ctx, &followv1.GetFollowStaticsRequest{
		UserId: 2,
	}).Return(&followv1.GetFollowStaticsResponse{
		FollowStatics: &followv1.FollowStatics{
			Followers: 5,
		},
	}, nil).Times(len(articleEvents))

	for _, event := range articleEvents {
		content, _ := json.Marshal(event)
		// 保证事件顺序
		time.Sleep(1 * time.Second)
		_, err := server.CreateFeedEvent(ctx, &feedv1.CreateFeedEventRequest{
			FeedEvent: &feedv1.FeedEvent{
				Type:    service.ArticleEventName,
				Content: string(content),
			},
		})
		if err != nil {
			return err
		}
	}

	articleEvents2 := []ArticleEvent{
		{
			Uid:   "3",
			Aid:   "5",
			Title: "用户3发表了文章5",
		},
		{
			Uid:   "3",
			Aid:   "6",
			Title: "用户3发表了文章6",
		},
		{
			Uid:   "3",
			Aid:   "7",
			Title: "用户3发表了文章7",
		},
	}

	// 设置用户3的关注统计, 用户3有2个粉丝, 小于阈值4 , 不是大v
	mockFollowClient.EXPECT().GetFollowStatics(gomock.Any(), &followv1.GetFollowStaticsRequest{
		UserId: 3,
	}).Return(&followv1.GetFollowStaticsResponse{
		FollowStatics: &followv1.FollowStatics{
			Followers: 2,
		},
	}, nil).Times(len(articleEvents2))

	// 设置用户3的关注者列表, 用户3有2个粉丝, 用户1和用户4
	mockFollowClient.EXPECT().GetFollower(gomock.Any(), &followv1.GetFollowerRequest{
		Followee: 3,
	}).Return(&followv1.GetFollowerResponse{
		FollowRelations: []*followv1.FollowRelation{
			{
				Id:       6,
				Follower: 1,
				Followee: 3,
			},
			{
				Id:       7,
				Follower: 4,
				Followee: 3,
			},
		},
	}, nil).AnyTimes()

	for _, event := range articleEvents2 {
		content, _ := json.Marshal(event)
		// 保证事件顺序
		time.Sleep(1 * time.Second)
		_, err := server.CreateFeedEvent(ctx, &feedv1.CreateFeedEventRequest{
			FeedEvent: &feedv1.FeedEvent{
				Type:    service.ArticleEventName,
				Content: string(content),
			},
		})
		if err != nil {
			return err
		}
	}

	// 创建点赞事件, 用户1被点赞
	likeEvents := []LikeEvent{
		{
			Liked: "1",
			Liker: "10",
			BizID: "8",
			Biz:   "article",
		},
		{
			Liked: "1",
			BizID: "9",
			Biz:   "article",
			Liker: "11",
		},
		{
			Liked: "1",
			BizID: "10",
			Biz:   "article",
			Liker: "12",
		},
	}
	for _, event := range likeEvents {
		content, _ := json.Marshal(event)
		time.Sleep(1 * time.Second)
		_, err := server.CreateFeedEvent(ctx, &feedv1.CreateFeedEventRequest{
			FeedEvent: &feedv1.FeedEvent{
				Type:    service.LikeEventName,
				Content: string(content),
			},
		})
		if err != nil {
			return err
		}
	}
	// 创建关注事件 用户1被关注
	followEvents := []FollowEvent{
		{
			Followee: "1",
			Follower: "2",
		},
		{
			Followee: "1",
			Follower: "3",
		},
		{
			Followee: "1",
			Follower: "4",
		},
	}

	for _, event := range followEvents {
		content, _ := json.Marshal(event)
		time.Sleep(1 * time.Second)
		_, err := server.CreateFeedEvent(ctx, &feedv1.CreateFeedEventRequest{
			FeedEvent: &feedv1.FeedEvent{
				Type:    service.FollowEventName,
				Content: string(content),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FeedTestSuite) getFeedEventWant(ctx context.Context, mockFollowClient *followv1Mock.MockFollowServiceClient, server feedv1.FeedSvcServer) []*feedv1.FeedEvent {
	// 获得用户1的关注列表
	// 用户1关注了用户2,3,4,5,6
	mockFollowClient.EXPECT().GetFollowee(gomock.Any(), &followv1.GetFolloweeRequest{
		Follower: 1,
		Offset:   0,
		Limit:    200,
	}).Return(&followv1.GetFolloweeResponse{
		FollowRelations: []*followv1.FollowRelation{
			{
				Id:       1,
				Follower: 1,
				Followee: 2,
			},
			{
				Id:       6,
				Follower: 1,
				Followee: 3,
			},
			{
				Id:       8,
				Follower: 1,
				Followee: 4,
			},
			{
				Id:       9,
				Follower: 1,
				Followee: 5,
			},
			{
				Id:       10,
				Follower: 1,
				Followee: 6,
			},
		},
	}, nil).AnyTimes()

	wantArticleEvents1 := []ArticleEvent{
		{
			Uid:   "2",
			Aid:   "1",
			Title: "用户2发表了文章1",
		},
		{
			Uid:   "2",
			Aid:   "2",
			Title: "用户2发表了文章2",
		},
		{
			Uid:   "2",
			Aid:   "3",
			Title: "用户2发表了文章3",
		},
		{
			Uid:   "2",
			Aid:   "4",
			Title: "用户2发表了文章4",
		},
	}
	wantArticleEvents2 := []ArticleEvent{
		{
			Uid:   "3",
			Aid:   "5",
			Title: "用户3发表了文章5",
		},
		{
			Uid:   "3",
			Aid:   "6",
			Title: "用户3发表了文章6",
		},
		{
			Uid:   "3",
			Aid:   "7",
			Title: "用户3发表了文章7",
		},
	}

	// 用户1被点赞
	// 用户1被用户10,11,12点赞
	// 业务类型是 article
	// 被点赞的对象是文章
	wantLikeEvents := []LikeEvent{
		{
			Liked: "1",
			Liker: "10",
			BizID: "8",
			Biz:   "article",
		},
		{
			Liked: "1",
			BizID: "9",
			Biz:   "article",
			Liker: "11",
		},
		{
			Liked: "1",
			BizID: "10",
			Biz:   "article",
			Liker: "12",
		},
	}
	// 用户1被关注
	// 用户1被用户2,3,4关注
	wantFollowEvents := []FollowEvent{
		{
			Followee: "1",
			Follower: "2",
		},
		{
			Followee: "1",
			Follower: "3",
		},
		{
			Followee: "1",
			Follower: "4",
		},
	}
	events := make([]*feedv1.FeedEvent, 0, 32)
	for i := len(wantFollowEvents) - 1; i >= 0; i-- {
		e := wantFollowEvents[i]
		content, _ := json.Marshal(e)
		events = append(events, &feedv1.FeedEvent{
			User: &feedv1.User{
				Id: 1,
			},
			Type:    service.FollowEventName,
			Content: string(content),
		})
	}
	for i := len(wantLikeEvents) - 1; i >= 0; i-- {
		e := wantLikeEvents[i]
		content, _ := json.Marshal(e)
		events = append(events, &feedv1.FeedEvent{
			User: &feedv1.User{
				Id: 1,
			},
			Type:    service.LikeEventName,
			Content: string(content),
		})
	}
	for i := len(wantArticleEvents2) - 1; i >= 0; i-- {
		e := wantArticleEvents2[i]
		content, _ := json.Marshal(e)
		events = append(events, &feedv1.FeedEvent{
			User: &feedv1.User{
				Id: 1,
			},
			Type:    service.ArticleEventName,
			Content: string(content),
		})
	}
	for i := len(wantArticleEvents1) - 1; i >= 0; i-- {
		e := wantArticleEvents1[i]
		content, _ := json.Marshal(e)
		uid, _ := strconv.ParseInt(e.Uid, 10, 64)
		events = append(events, &feedv1.FeedEvent{
			User: &feedv1.User{
				Id: uid,
			},
			Type:    service.ArticleEventName,
			Content: string(content),
		})
	}

	return events
}

func TestFeedTestSuite(t *testing.T) {
	suite.Run(t, new(FeedTestSuite))
}
