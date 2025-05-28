package grpc

import (
	"context"

	oauth2v1 "github.com/pluckhuang/goweb/aweb/api/proto/gen/oauth2/v1"
	"github.com/pluckhuang/goweb/oauth2/service"
	"google.golang.org/grpc"
)

type Oauth2ServiceServer struct {
	oauth2v1.UnimplementedOAuth2ServiceServer
	service service.Oauth2Service
}

func NewOauth2ServiceServer(svc service.Oauth2Service) *Oauth2ServiceServer {
	return &Oauth2ServiceServer{
		service: svc,
	}
}

func (o *Oauth2ServiceServer) Register(server grpc.ServiceRegistrar) {
	oauth2v1.RegisterOAuth2ServiceServer(server, o)
}

func (o *Oauth2ServiceServer) GetAuthURL(ctx context.Context, req *oauth2v1.GetAuthURLRequest) (*oauth2v1.GetAuthURLResponse, error) {
	url, state, err := o.service.GetAuthURL(ctx, req.GetPlatform())
	if err != nil {
		return nil, err
	}
	return &oauth2v1.GetAuthURLResponse{
		AuthUrl: url,
		State:   state,
	}, nil
}

func (o *Oauth2ServiceServer) HandleCallback(ctx context.Context, req *oauth2v1.HandleCallbackRequest) (*oauth2v1.HandleCallbackResponse, error) {
	info, err := o.service.HandleCallback(ctx, req.GetPlatform(), req.GetCode(), req.GetState())
	if err != nil {
		return nil, err
	}
	return &oauth2v1.HandleCallbackResponse{
		AccessToken: info.AccessToken,
		Error:       "",
	}, nil
}
