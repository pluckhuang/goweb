package ioc

import (
	"github.com/pluckhuang/goweb/aweb/internal/service/oauth2/wechat"
	"github.com/pluckhuang/goweb/aweb/pkg/logger"
	// "os"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	// appID, ok := os.LookupEnv("WECHAT_APP_ID")
	// if !ok {
	// 	panic("找不到环境变量 WECHAT_APP_ID")
	// }
	// appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	// if !ok {
	// 	panic("找不到环境变量 WECHAT_APP_SECRET")
	// }
	return wechat.NewService("appID", "appSecret", l)
}
