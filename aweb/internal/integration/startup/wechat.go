package startup

import (
	"github.com/pluckhuang/goweb/aweb/internal/service/oauth2/wechat"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
