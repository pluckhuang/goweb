package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/pluckhuang/goweb/oauth2/domain"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

type twitterService struct {
	oauthConfig *oauth2.Config
	logger      logger.LoggerV1
	redisClient redis.Cmdable
	stateTTL    time.Duration
	verifierTTL time.Duration
}

func NewTwitterService(oauthConfig *oauth2.Config,
	cmd redis.Cmdable,
	logger logger.LoggerV1) Handler {
	return &twitterService{
		oauthConfig: oauthConfig,
		redisClient: cmd,
		logger:      logger,
		stateTTL:    10 * time.Minute,
		verifierTTL: 10 * time.Minute,
	}
}

func (s *twitterService) GetAuthURL(ctx context.Context) (string, string, error) {
	// 生成 state 和 verifier
	state := uuid.New().String()
	verifier := oauth2.GenerateVerifier()
	challenge := oauth2.S256ChallengeOption(verifier)

	// 存储 state 和 verifier 到 Redis
	err := s.redisClient.Set(ctx, "oauth:state:"+state, state, s.stateTTL).Err()
	if err != nil {
		return "", "", fmt.Errorf("failed to store state: %v", err)
	}
	err = s.redisClient.Set(ctx, "oauth:verifier:"+state, verifier, s.verifierTTL).Err()
	if err != nil {
		return "", "", fmt.Errorf("failed to store verifier: %v", err)
	}
	// 生成授权 URL
	authURL := s.oauthConfig.AuthCodeURL(state, challenge)
	return authURL, state, nil
}

func (s *twitterService) HandleCallback(ctx context.Context, code string, state string) (domain.Oauth2Info, error) {
	storedState, err := s.redisClient.Get(ctx, "oauth:state:"+state).Result()
	if err != nil || storedState != state {
		return domain.Oauth2Info{}, fmt.Errorf("invalid state")
	}

	// 获取 verifier
	verifier, err := s.redisClient.Get(ctx, "oauth:verifier:"+state).Result()
	if err != nil {
		return domain.Oauth2Info{}, fmt.Errorf("invalid verifier")
	}

	// 删除 Redis 中的 state 和 verifier
	s.redisClient.Del(ctx, "oauth:state:"+state)
	s.redisClient.Del(ctx, "oauth:verifier:"+state)
	// 交换 token
	token, err := s.oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return domain.Oauth2Info{}, fmt.Errorf("failed to exchange token: %w", err)
	}
	return domain.Oauth2Info{AccessToken: token.AccessToken}, nil
}
