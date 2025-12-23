package distributedLock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type DistributedLock struct {
	redisCli *redis.Client
	name     string

	// value 用于标识锁的持有者，释放时必须匹配才删除
	value string
	// ttl 锁的过期时间，防止死锁
	ttl time.Duration

	mu          sync.Mutex
	renewCancel context.CancelFunc
	renewDone   chan struct{}
}

func NewDistributedLock(cli *redis.Client, name string) *DistributedLock {
	return &DistributedLock{
		redisCli: cli,
		name:     name,
		ttl:      10 * time.Second,
	}
}

func (d *DistributedLock) WithTTL(ttl time.Duration) *DistributedLock {
	if ttl > 0 {
		d.ttl = ttl
	}
	return d
}

func (d *DistributedLock) Lock(ctx context.Context, retryInterval time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if retryInterval <= 0 {
		retryInterval = 50 * time.Millisecond
	}

	// 每次尝试获取锁前生成新 token
	token := newToken()

	for {
		ok, err := d.redisCli.SetNX(ctx, d.name, token, d.ttl).Result()
		if err != nil {
			return err
		}
		if ok {
			d.mu.Lock()
			d.value = token
			// 启动内置续约（watchdog）
			d.startRenewLocked(ctx)
			d.mu.Unlock()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryInterval):
		}
	}
}

var unlockScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
	  return redis.call("DEL", KEYS[1])
	else
	  return 0
	end
`)

func (d *DistributedLock) Unlock(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	// 先停续约，避免 Unlock 时还在续期
	d.mu.Lock()
	if d.value == "" {
		d.mu.Unlock()
		return
	}
	token := d.value
	d.stopRenewLocked()
	d.mu.Unlock()

	_, err := unlockScript.Run(ctx, d.redisCli, []string{d.name}, token).Result()
	if err != nil {
		return
	}

	d.mu.Lock()
	// 无论 DEL 是否删除成功（可能锁已过期），都清理本地状态
	d.value = ""
	d.mu.Unlock()

}

var renewScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
  return 0
end
`)

func (d *DistributedLock) startRenewLocked(parent context.Context) {
	// 已有续约在跑就先停（防止多开）
	d.stopRenewLocked()

	if d.value == "" {
		return
	}

	interval := d.ttl / 3
	if interval <= 0 {
		interval = 1 * time.Second
	}

	ctx, cancel := context.WithCancel(parent)
	done := make(chan struct{})

	lockName := d.name
	lockValue := d.value
	ttlMs := int64(d.ttl / time.Millisecond)

	d.renewCancel = cancel
	d.renewDone = done

	go func() {
		defer close(done)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 仅当 value 匹配时才续期；不匹配说明锁已丢失,过期后被他人拿到
				r, e := renewScript.Run(ctx, d.redisCli, []string{lockName}, lockValue, ttlMs).Int64()
				if e != nil || r == 0 {
					return
				}
			}
		}
	}()
}

func (d *DistributedLock) stopRenewLocked() {
	if d.renewCancel != nil {
		d.renewCancel()
		d.renewCancel = nil
	}
	if d.renewDone != nil {
		ch := d.renewDone
		d.renewDone = nil
		// 在持锁的调用栈里等待 goroutine 退出，确保不再续期
		d.mu.Unlock()
		<-ch
		d.mu.Lock()
	}
}

func newToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
