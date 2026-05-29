package pkg

import "user-srv/global"

const ServiceListCachePrefix = "service:list:"

// ClearServiceListCache 清理服务列表缓存
func ClearServiceListCache() error {
	if global.RDB == nil {
		return nil
	}

	keys, err := global.RDB.Keys(global.Ctx, ServiceListCachePrefix+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	return global.RDB.Del(global.Ctx, keys...).Err()
}
