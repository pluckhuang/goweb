package grpc

import (
	"context"
	"log"
	"testing"
	"time"

	followv1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/follow/v1"
	"google.golang.org/grpc"
)

func newClient(t *testing.T) followv1.FollowServiceClient {
	conn, err := grpc.Dial("localhost:8076", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return followv1.NewFollowServiceClient(conn)
}

func TestFollowServiceClient_Follow(t *testing.T) {
	client := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &followv1.FollowRequest{
		Follower: 1,
		Followee: 2,
	}
	resp, err := client.Follow(ctx, req)
	if err != nil {
		t.Fatalf("Follow failed: %v", err)
	}
	log.Printf("Follow response: %+v", resp)
}

func TestFollowServiceClient_GetFollowee(t *testing.T) {
	client := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &followv1.GetFolloweeRequest{
		Follower: 1,
		Offset:   0,
		Limit:    10,
	}
	resp, err := client.GetFollowee(ctx, req)
	if err != nil {
		t.Fatalf("GetFollowee failed: %v", err)
	}
	log.Printf("GetFollowee response: %+v", resp)
}

func TestFollowServiceClient_GetFollower(t *testing.T) {
	client := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &followv1.GetFollowerRequest{
		Followee: 2,
		Offset:   0,
		Limit:    10,
	}
	resp, err := client.GetFollower(ctx, req)
	if err != nil {
		t.Fatalf("GetFollower failed: %v", err)
	}
	log.Printf("GetFollower response: %+v", resp)
}

func TestFollowServiceClient_FollowInfo(t *testing.T) {
	client := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &followv1.FollowInfoRequest{
		Follower: 1,
		Followee: 2,
	}
	resp, err := client.FollowInfo(ctx, req)
	if err != nil {
		t.Fatalf("FollowInfo failed: %v", err)
	}
	log.Printf("FollowInfo response: %+v", resp)
}

func TestFollowServiceClient_GetFollowStatics(t *testing.T) {
	client := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &followv1.GetFollowStaticsRequest{
		UserId: 1,
	}
	resp, err := client.GetFollowStatics(ctx, req)
	if err != nil {
		t.Fatalf("GetFollowStatics failed: %v", err)
	}
	log.Printf("GetFollowStatics response: %+v", resp)
}

func TestFollowServiceClient_CancelFollow(t *testing.T) {
	client := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &followv1.CancelFollowRequest{
		Follower: 1,
		Followee: 2,
	}
	resp, err := client.CancelFollow(ctx, req)
	if err != nil {
		t.Fatalf("CancelFollow failed: %v", err)
	}
	log.Printf("CancelFollow response: %+v", resp)
}
