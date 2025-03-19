//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(aweb-record-mysql:13306)/aweb",
	},
	Redis: RedisConfig{
		Addr: "aweb-record-redis:6380",
	},
}
