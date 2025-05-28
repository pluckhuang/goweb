package service

import (
	"context"
	"fmt"

	"github.com/pluckhuang/goweb/oauth2/domain"
)

type oauth2Service struct {
	handlerMap map[string]Handler
}

func NewOauth2Service(handlerMap map[string]Handler) Oauth2Service {
	return &oauth2Service{
		handlerMap: handlerMap,
	}
}

func (f *oauth2Service) registerService(typ string, handler Handler) {
	f.handlerMap[typ] = handler
}

func (s *oauth2Service) GetAuthURL(ctx context.Context, platform string) (string, string, error) {
	handler, ok := s.handlerMap[platform]
	if !ok {
		return "", "", fmt.Errorf("unsupported platform: %s", platform)
	}

	// 生成授权 URL
	authURL, state, err := handler.GetAuthURL(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get auth URL for platform %s: %w", platform, err)
	}

	return authURL, state, nil
}

func (s *oauth2Service) HandleCallback(ctx context.Context, platform string, code string, state string) (domain.Oauth2Info, error) {
	handler, ok := s.handlerMap[platform]
	if !ok {
		return domain.Oauth2Info{}, fmt.Errorf("unsupported platform: %s", platform)
	}

	// 处理回调
	info, err := handler.HandleCallback(ctx, code, state)
	if err != nil {
		return domain.Oauth2Info{}, fmt.Errorf("failed to handle callback for platform %s: %w", platform, err)
	}

	return domain.Oauth2Info{AccessToken: info.AccessToken}, nil
}
