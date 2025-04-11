package ioc

import (
	"log"
	"os"
	"time"

	"github.com/pluckhuang/goweb/aweb/internal/repository/dao"
	"github.com/pluckhuang/goweb/aweb/pkg/gormx"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
)

// 引入 glogger
func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:3316)/webook",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到标准输出
			logger.Config{
				SlowThreshold: 300 * time.Millisecond, // 设置慢查询阈值为 300ms
				LogLevel:      logger.Warn,            // 设置日志级别为 Warn
				Colorful:      true,                   // 启用彩色日志输出
			},
		),
	})
	if err != nil {
		panic(err)
	}

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "aweb",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "pluckh.com",
		Subsystem: "aweb",
		Name:      "gorm_db",
		Help:      "统计 GORM 的数据库查询",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	err = db.Use(cb)
	if err != nil {
		panic(err)
	}

	err = db.Use(tracing.NewPlugin(
		tracing.WithoutMetrics(),
		tracing.WithDBName("aweb"),
	))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
