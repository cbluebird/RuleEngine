package startConfig

import (
	"engine/config/config"
	"engine/config/database"
	"engine/config/redis"
)

func Init() {
	config.InitConfig()
	database.Init()
	redis.Init()
}
