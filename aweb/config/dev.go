//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:password@tcp(localhost:3306)/aweb",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
