package main

import (
	"net/http"

	"fmt"

	"context"
	"log"

	oauth2v1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/oauth2/v1"

	"github.com/pluckhuang/goweb/aweb/pkg/grpcx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const GRPC_ServerPort = "8076"

func main() {
	initViperV2Watch()
	app := Init()
	go func() {
		err := app.server.Serve()
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}()
	initHttpServer()
	// 启动 HTTP 服务器
	log.Println("HTTP server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initHttpServer() {
	http.HandleFunc("/callback/", func(w http.ResponseWriter, r *http.Request) {
		platform := r.URL.Path[len("/callback/"):] // 提取平台：twitter, discord, telegram
		code := r.FormValue("code")
		state := r.FormValue("state")

		// 连接 gRPC 服务器
		conn, err := grpc.Dial("localhost:"+GRPC_ServerPort, grpc.WithInsecure())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to connect to gRPC: %v", err), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := oauth2v1.NewOAuth2ServiceClient(conn)
		resp, err := client.HandleCallback(context.Background(), &oauth2v1.HandleCallbackRequest{
			Platform: platform,
			Code:     code,
			State:    state,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("gRPC error: %v", err), http.StatusInternalServerError)
			return
		}

		if resp.Error != "" {
			http.Error(w, resp.Error, http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Access Token: %s", resp.AccessToken)
	})

	// 示例：获取授权 URL 的 HTTP 端点
	http.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		platform := r.URL.Path[len("/auth/"):] // 提取平台：twitter, discord, telegram

		// 连接 gRPC 服务器
		conn, err := grpc.Dial("localhost:"+GRPC_ServerPort, grpc.WithInsecure())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to connect to gRPC: %v", err), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := oauth2v1.NewOAuth2ServiceClient(conn)
		resp, err := client.GetAuthURL(context.Background(), &oauth2v1.GetAuthURLRequest{
			Platform: platform,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("gRPC error: %v", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, resp.AuthUrl, http.StatusFound)
	})
}

func initViperV2Watch() {
	cfile := pflag.String("config",
		"config/dev.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

type App struct {
	server *grpcx.Server
}
