package init

import (
	"fmt"
	"user-srv/global"

	"github.com/go-redis/redis/v8"
)

func RedisInit() {
	data := global.ConfigData.Redis
	addr := fmt.Sprintf("%s:%d", data.Host, data.Port)
	global.RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: data.Password, // no password set
		DB:       data.Database, // use default DB
	})
	err := global.RDB.Ping(global.Ctx).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis连接成功")
}
