package security

import (
	"context"
	"fmt"
	"time"

	"admin-srv/global"
	"github.com/redis/go-redis/v9"
)

const (
	// LoginFailKey Redis key前缀: login_fail:{username}
	LoginFailKey = "login_fail:%s"
	// MaxFailCount 最大允许失败次数
	MaxFailCount = 5
	// LockDuration 锁定时长
	LockDuration = 15 * time.Minute
)

// CheckLoginLimit 检查用户是否被限制登录
// 返回: 是否被锁定, 剩余锁定时间(秒), 错误信息
func CheckLoginLimit(ctx context.Context, username string) (bool, int64, error) {
	key := fmt.Sprintf(LoginFailKey, username)

	// 获取当前失败次数
	failCount, err := global.RDB.Get(ctx, key).Int()
	if err == redis.Nil {
		// 没有记录，允许登录
		return false, 0, nil
	}
	if err != nil {
		// Redis错误，不限制但记录日志
		return false, 0, err
	}

	if failCount >= MaxFailCount {
		// 获取剩余TTL
		ttl := global.RDB.TTL(ctx, key).Val()
		remainingSec := int64(ttl.Seconds())
		if remainingSec < 0 {
			remainingSec = int64(LockDuration.Seconds())
		}
		return true, remainingSec, nil
	}

	return false, 0, nil
}

// RecordLoginFail 记录一次登录失败
// 返回: 当前累计失败次数, 是否已达到上限被锁定
func RecordLoginFail(ctx context.Context, username string) (int, bool) {
	key := fmt.Sprintf(LoginFailKey, username)

	pipe := global.RDB.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, LockDuration)
	_, _ = pipe.Exec(ctx)

	count := int(incrCmd.Val())
	isLocked := count >= MaxFailCount

	return count, isLocked
}

// ClearLoginFail 清除登录失败记录（登录成功后调用）
func ClearLoginFail(ctx context.Context, username string) {
	key := fmt.Sprintf(LoginFailKey, username)
	global.RDB.Del(ctx, key)
}

// GetRemainingAttempts 获取剩余尝试次数
func GetRemainingAttempts(ctx context.Context, username string) int {
	key := fmt.Sprintf(LoginFailKey, username)
	failCount, err := global.RDB.Get(ctx, key).Int()
	if err != nil {
		return MaxFailCount
	}
	remaining := MaxFailCount - failCount
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

// FormatLockTime 格式化锁定时间为人类可读格式
func FormatLockTime(seconds int64) string {
	minutes := seconds / 60
	secs := seconds % 60
	if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, secs)
	}
	return fmt.Sprintf("%d秒", secs)
}
