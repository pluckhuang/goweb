package grpc

import (
	"context"
	"log"
	"testing"
	"time"

	articlev1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/article/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newArticleClient(t testing.TB) articlev1.ArticleServiceClient {
	conn, err := grpc.Dial("localhost:8076", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return articlev1.NewArticleServiceClient(conn)
}

func TestArticleService_EndToEnd(t *testing.T) {
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
