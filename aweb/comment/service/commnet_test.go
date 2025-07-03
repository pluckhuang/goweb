package service

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	commentv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/comment/v1"
)

func TestCommentServiceClient(t *testing.T) {
	addr := "localhost:8076"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := commentv1.NewCommentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. 获取一级评论
	listResp, err := client.GetCommentList(ctx, &commentv1.CommentListRequest{
		Biz:   "article",
		Bizid: 123,
		MinId: 0,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("GetCommentList err: %v", err)
	}
	t.Logf("一级评论: %+v, 数量: %d, hasMore: %+v", listResp.Comments, len(listResp.Comments), listResp.HasMore)

	// 2. 创建1级评论
	createResp, err := client.CreateComment(ctx, &commentv1.CreateCommentRequest{
		Comment: &commentv1.Comment{
			Uid:     1,
			Biz:     "article",
			Bizid:   123,
			Content: "hello world",
			Ctime:   timestamppb.Now(),
			Utime:   timestamppb.Now(),
		},
	})
	if err != nil {
		t.Fatalf("CreateComment err: %v", err)
	}
	t.Logf("创建评论: %+v", createResp)

	// 3. 获取一级评论 again
	listResp, err = client.GetCommentList(ctx, &commentv1.CommentListRequest{
		Biz:   "article",
		Bizid: 123,
		MinId: 0,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("GetCommentList err: %v", err)
	}
	t.Logf("再次获取一级评论: %+v, 数量: %d, hasMore: %+v", listResp.Comments, len(listResp.Comments), listResp.HasMore)

	// 4. 创建2级评论
	createResp, err = client.CreateComment(ctx, &commentv1.CreateCommentRequest{
		Comment: &commentv1.Comment{
			Uid:     1,
			Biz:     "article",
			Bizid:   123,
			Content: "abc",
			RootComment: &commentv1.Comment{
				Id: 1, // 根评论ID
			},
			ParentComment: &commentv1.Comment{
				Id: 1, // 父评论ID
			},
			Ctime: timestamppb.Now(),
			Utime: timestamppb.Now(),
		},
	})
	if err != nil {
		t.Fatalf("CreateComment err: %v", err)
	}
	t.Logf("创建2级评论: %+v", createResp)

	// 5. 获取更多回复
	moreResp, err := client.GetMoreReplies(ctx, &commentv1.GetMoreRepliesRequest{
		Rid:   1,
		MaxId: 0,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("GetMoreReplies err: %v", err)
	}
	t.Logf("更多回复: %+v, 数量: %d, hasMore: %+v", moreResp.Replies, len(moreResp.Replies), moreResp.HasMore)

	// 6. 删除评论
	_, err = client.DeleteComment(ctx, &commentv1.DeleteCommentRequest{
		Id: 1,
	})
	if err != nil {
		t.Fatalf("DeleteComment err: %v", err)
	}
	t.Log("删除评论成功")

}
