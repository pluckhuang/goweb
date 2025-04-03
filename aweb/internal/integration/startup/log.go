package startup

import "github.com/pluckhuang/goweb/aweb/pkg/logger"

func InitLogger() logger.LoggerV1 {
	return logger.NewNopLogger()
}
