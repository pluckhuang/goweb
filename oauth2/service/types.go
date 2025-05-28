package service

import (
	"context"

	"github.com/pluckhuang/goweb/oauth2/domain"
)

type Oauth2Service interface {
	GetAuthURL(ctx context.Context, platform string) (string, string, error)
	HandleCallback(ctx context.Context, platform string, code string, state string) (domain.Oauth2Info, error)
}

// Handler 具体业务处理逻辑
type Handler interface {
	GetAuthURL(ctx context.Context) (string, string, error)
	HandleCallback(ctx context.Context, code string, state string) (domain.Oauth2Info, error)
}
