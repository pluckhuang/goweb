package grpc

import (
	"context"
	"testing"

	articlev1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/article/v1"
	"github.com/pluckhuang/goweb/aweb/article/ioc"
	"google.golang.org/grpc"
)

// cleanupBenchDatabase 清理基准测试数据库
func cleanupBenchDatabase() {
	db := ioc.InitDBV2()
	tables := []string{"articles", "published_articles"}

	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}
}

func newArticleClient1(t testing.TB) articlev1.ArticleServiceClient {
	conn, err := grpc.Dial("localhost:8076", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return articlev1.NewArticleServiceClient(conn)
}

// go test -bench=. -run=^$ ./grpc/article_benchmark_test.go
// goos: darwin
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-7820HQ CPU @ 2.90GHz
// BenchmarkArticleService_Save-8         	      93	  11275041 ns/op
// BenchmarkArticleService_GetPubById-8   	     868	   1439900 ns/op
// PASS
// ok  	command-line-arguments	3.591s

// 写操作（Save）平均耗时约 13.46ms/次
// 读操作（GetPubById）平均耗时约 1.33ms/次

func BenchmarkArticleService_Save(b *testing.B) {
	defer cleanupBenchDatabase() // 基准测试结束后清理数据库

	client := newArticleClient1(b)
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := client.Save(ctx, &articlev1.SaveRequest{
			Article: &articlev1.Article{
				Title:    "Bench Title",
				Content:  "Bench Content",
				AuthorId: 1001,
			},
		})
		if err != nil {
			b.Fatalf("Save failed: %v", err)
		}
	}
}

func BenchmarkArticleService_GetPubById(b *testing.B) {
	client := newArticleClient1(b)
	ctx := context.Background()

	// 先写一篇文章，获得ID
	saveResp, err := client.Publish(ctx, &articlev1.PublishRequest{
		Article: &articlev1.Article{
			Title:    "Bench Title",
			Content:  "Bench Content",
			AuthorId: 1001,
		},
	})
	if err != nil {
		b.Fatalf("pre-Save failed: %v", err)
	}
	aid := saveResp.Id

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetPubById(ctx, &articlev1.GetPubByIdRequest{Id: aid, Uid: 1001})
		if err != nil {
			b.Fatalf("GetPubById failed: %v", err)
		}
	}
}
