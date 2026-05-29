package global

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var Config *AppConfig
var DB *gorm.DB
var Ctx = context.Background()
var RDB *redis.Client
