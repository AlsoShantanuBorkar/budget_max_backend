package utils

import (
	"fmt"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/redis"
)

const (
	MaxLoginAttempts = 5
	LoginWindow      = 10 * time.Minute
	LoginLockTime    = 1 * time.Hour
)

func CheckAndTrackLoginAttempts(email string) error {
	attemptKey := fmt.Sprintf("login_attempts:%s", email)
	lockKey := fmt.Sprintf("login_lock:%s", email)

	isLocked, err := redis.Client.Exists(redis.Ctx, lockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check login lock: %w", err)
	}

	if isLocked == 1 {
		return fmt.Errorf("account is locked due to too many login attempts")
	}

	attempts, err := redis.Client.Incr(redis.Ctx, attemptKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get login attempts: %w", err)
	}

	if attempts == 1 {
		redis.Client.Expire(redis.Ctx, attemptKey, LoginWindow)
	}

	if attempts >= MaxLoginAttempts {
		redis.Client.Set(redis.Ctx, lockKey, "1", LoginLockTime)
		return fmt.Errorf("account is locked due to too many login attempts")
	}

	redis.Client.Expire(redis.Ctx, attemptKey, LoginWindow)
	return nil
}

func ResetLoginAttempts(email string) {
	attemptKey := fmt.Sprintf("login_attempts:%s", email)
	lockKey := fmt.Sprintf("login_lock:%s", email)

	redis.Client.Del(redis.Ctx, attemptKey, lockKey)
}
