package grpc

import (
	"context"
	"log"
	"testing"
	"time"

	articlev1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/article/v1"
	"github.com/pluckhuang/goweb/aweb/article/ioc"
	"github.com/pluckhuang/goweb/aweb/pkg/grpcx"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// cleanupDatabase 清理数据库表
func cleanupDatabase(t testing.TB) {
	db := ioc.InitDBV2()

	// 清空文章相关表
	tables := []string{
		"articles",
		"published_articles",
	}

	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table).Error; err != nil {
			t.Logf("清理 %s 表失败: %v", table, err)
		}
	}

	log.Println("数据库清理完成")
}

func newArticleClient(t testing.TB) articlev1.ArticleServiceClient {
	mwCfg := grpcx.MiddlewareConfig{
		BreakerSettings: gobreaker.Settings{
			Name:        "intr-grpc",
			MaxRequests: 3,
			Interval:    60 * time.Second,
			Timeout:     10 * time.Second,
		},
		Limiter:        rate.NewLimiter(10, 2), // 每秒10次，突发2次
		RetryMax:       2,
		RetryBaseDelay: 100 * time.Millisecond,
		RetryMaxDelay:  2 * time.Second,
	}
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcx.NewUnaryClientInterceptor(mwCfg)),
	}
	conn, err := grpc.Dial("localhost:8076", opts...)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return articlev1.NewArticleServiceClient(conn)
}

func TestArticleService_EndToEnd(t *testing.T) {
	// 测试前清理数据库
	defer cleanupDatabase(t) // 测试结束后清理数据库

	client := newArticleClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Save
	saveResp, err := client.Save(ctx, &articlev1.SaveRequest{
		Article: &articlev1.Article{
			Title:    "Test Title",
			Content:  "Test Content",
			AuthorId: 1001,
		},
	})
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	aid := saveResp.Id
	log.Printf("Saved article id: %d", aid)

	// 2. Publish
	pubResp, err := client.Publish(ctx, &articlev1.PublishRequest{
		Article: &articlev1.Article{
			Id:       aid,
			Title:    "Test Title",
			Content:  "Test Content",
			AuthorId: 1001,
		},
	})
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
	log.Printf("Published article id: %d", pubResp.Id)

	// 3. GetById
	getResp, err := client.GetById(ctx, &articlev1.GetByIdRequest{Id: aid})
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	log.Printf("GetById: %+v", getResp.Article)

	// 4. GetByAuthor
	authorResp, err := client.GetByAuthor(ctx, &articlev1.GetByAuthorRequest{
		Uid:    1001,
		Offset: 0,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("GetByAuthor failed: %v", err)
	}
	log.Printf("GetByAuthor: %d articles", len(authorResp.Articles))

	// 5. GetPubById
	pubByIdResp, err := client.GetPubById(ctx, &articlev1.GetPubByIdRequest{
		Id:  aid,
		Uid: 1001,
	})
	if err != nil {
		t.Fatalf("GetPubById failed: %v", err)
	}
	log.Printf("GetPubById: %+v", pubByIdResp.Article)

	// 6. ListPub
	listResp, err := client.ListPub(ctx, &articlev1.ListPubRequest{
		Start:  timestamppb.New(time.Now().Add(-time.Hour)),
		Offset: 0,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("ListPub failed: %v", err)
	}
	log.Printf("ListPub: %d articles", len(listResp.Articles))

	// 7. Withdraw
	_, err = client.Withdraw(ctx, &articlev1.WithdrawRequest{
		Uid: 1001,
		Id:  aid,
	})
	if err != nil {
		t.Fatalf("Withdraw failed: %v", err)
	}
	log.Printf("Withdraw success")
}
