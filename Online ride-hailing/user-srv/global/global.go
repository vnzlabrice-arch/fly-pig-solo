package global

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

var ConfigData configData
var DB *gorm.DB
var Ctx = context.Background()
var RDB *redis.Client
var Esc *elastic.Client
