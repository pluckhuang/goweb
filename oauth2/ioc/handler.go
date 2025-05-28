package ioc

import (
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	"github.com/pluckhuang/goweb/oauth2/service"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func InitService(cmd redis.Cmdable, logv1 logger.LoggerV1) map[string]service.Handler {
	type TwitterOauthConfig struct {
		ClientID     string `yaml:"TWITTER_CLIENT_ID"`
		ClientSecret string `yaml:"TWITTER_CLIENT_SECRET"`
		RedirectUrl  string `yaml:"REDIRECT_URL"`
		AuthURL      string `yaml:"AUTH_URL"`
		TokenURL     string `yaml:"TOKEN_URL"`
	}
	var cfg TwitterOauthConfig
	err := viper.UnmarshalKey("twitter", &cfg)
	if err != nil {
		panic(err)
	}
	twitterOauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectUrl,
		Scopes:       []string{"tweet.read", "users.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.AuthURL,
			TokenURL: cfg.TokenURL,
		},
	}
	return map[string]service.Handler{
		"twitter": service.NewTwitterService(twitterOauthConfig, cmd, logv1),
	}
}
