package pkg

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type DistributedLock interface {
	// Lock 基础加锁【SET NX】
	Lock(key, requestId string, expire time.Duration) (bool, error)
	// TryLock 带重试的抢锁【核心：指数退避 + 随机抖动】
	TryLock(key, requestId string, expire, waitTimeout time.Duration) (bool, error)
	// Unlock 解锁【Lua脚本保证原子性】
	Unlock(key, requestId string) (bool, error)
}

type RedisDistributedLock struct {
	client   *redis.Client
	ctx      context.Context
	watchDog map[string]context.CancelFunc
	mutex    sync.Mutex
}

func NewRedisDistributedLock(client *redis.Client, ctx context.Context) *RedisDistributedLock {
	return &RedisDistributedLock{
		client:   client,
		ctx:      ctx,
		watchDog: make(map[string]context.CancelFunc),
	}
}

// Lock 基础加锁

// SET NX：key不存在才设置，保证互斥

// expire：锁过期时间，防止宕机死锁

func (r *RedisDistributedLock) Lock(key, requestId string, expire time.Duration) (bool, error) {

	ok, err := r.client.SetNX(r.ctx, key, requestId, expire).Result()
	if err != nil {
		return false, nil
	}

	if ok {
		// 加锁成功，启动看门狗自动续期
		r.startWatchDog(key, requestId, expire)
	}
	return ok, err
}

// TryLock 带重试的抢锁【核心：指数退避 + 随机抖动】

// expire：锁过期时长

// waitTimeout：最大等待多久拿不到锁就放弃

func (r *RedisDistributedLock) TryLock(key, requestId string, expire, waitTimeout time.Duration) (bool, error) {

	timeout := time.Now().Add(waitTimeout) // 最大等待时间
	baseDelay := 10 * time.Millisecond     // 初始重试等待间隔
	maxDelay := 500 * time.Millisecond     // 最大等待间隔，防止无限等太久
	for time.Now().Before(timeout) {
		ok, err := r.Lock(key, requestId, expire)
		if ok || err != nil {
			return ok, err
		}

		baseDelay *= 2
		// 1. 指数退避，防止等待时间过长
		if baseDelay > maxDelay {
			baseDelay = maxDelay
		}
		// 2. 再加上随机抖动
		jitter := time.Duration(rand.Int63n(int64(baseDelay) / 2))
		// 3. 等待随机时间，避免所有重试同时触发
		time.Sleep(baseDelay + jitter)

		baseDelay *= 2
	}

	return false, errors.New("等待锁超时")
}

// Unlock 解锁【Lua脚本保证原子性】
// 逻辑：先判断当前锁的值是不是自己的requestId，是才删除，防止误删别人锁
func (r *RedisDistributedLock) Unlock(key, requestId string) (bool, error) {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	// KEYS[1] = 锁key，ARGV[1] = 当前请求唯一ID
	result, err := r.client.Eval(r.ctx, script, []string{key}, requestId).Result()
	if err != nil {
		return false, err
	}

	// 关闭当前锁对应的看门狗协程
	r.mutex.Lock()
	if cancel, ok := r.watchDog[key]; ok {
		cancel()
		delete(r.watchDog, key)
	}
	r.mutex.Unlock()

	// 返回1=解锁成功，0=不是自己的锁/锁已过期
	return result == 1, nil
}

// startWatchDog 看门狗（自动续期）
// 作用：业务执行时间 > 锁过期时间 时，自动延长锁时间，避免锁提前失效
func (r *RedisDistributedLock) startWatchDog(key, requestId string, expire time.Duration) {
	// 可取消上下文，解锁时关闭协程
	ctx, cancel := context.WithCancel(r.ctx)

	// 加锁保护map并发写入
	r.mutex.Lock()
	r.watchDog[key] = cancel
	r.mutex.Unlock()

	// 续期间隔：过期时间的1/3，高频续期更安全
	interval := expire / 3

	// 异步协程持续续期
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				// 解锁触发取消, 退出看门狗
				return
			case <-ticker.C:
				// 刷新key过期时间, 完成过期
				_ = r.client.Expire(r.ctx, key, expire).Err()
			}
		}
	}()
}
