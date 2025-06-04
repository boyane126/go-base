package bootstrap

import (
	"fmt"

	"github.com/boyane126/go-common/config"
	"github.com/boyane126/go-common/redis"
)

// SetupRedis 初始化 Redis
func SetupRedis() {
	redis.ConnectRedis(
		fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		config.GetString("redis.username"),
		config.Get("redis.password"),
		config.GetInt("redis.database"),
	)
}
