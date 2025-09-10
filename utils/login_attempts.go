package utils

import (
	"fmt"
	"time"

	"github.com/AlsoShantanuBorkar/budget_max/redis"
)

const (
	MaxLoginAttempts = 5
	LoginWindow      = 10 * time.Minute
	LoginLockTime    = 5 * time.Minute
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

	// Use a pipeline to make INCR and EXPIRE atomic
	pipe := redis.Client.Pipeline()
	attemptsCmd := pipe.Incr(redis.Ctx, attemptKey)
	pipe.Expire(redis.Ctx, attemptKey, LoginWindow)
	_, err = pipe.Exec(redis.Ctx)

	if err != nil {
		return fmt.Errorf("failed to track login attempts: %w", err)
	}

	attempts := attemptsCmd.Val()

	if attempts >= MaxLoginAttempts {
		redis.Client.Set(redis.Ctx, lockKey, "1", LoginLockTime)
		return fmt.Errorf("account is locked due to too many login attempts")
	}

	return nil
}

func ResetLoginAttempts(email string) {
	attemptKey := fmt.Sprintf("login_attempts:%s", email)
	lockKey := fmt.Sprintf("login_lock:%s", email)

	redis.Client.Del(redis.Ctx, attemptKey, lockKey)
}
