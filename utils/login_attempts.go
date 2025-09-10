package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	MaxLoginAttempts = 5
	LoginWindow      = 10 * time.Minute
	LoginLockTime    = 5 * time.Minute
)

func CheckAndTrackLoginAttempts(email string, redisClient *redis.Client, ctx context.Context) error {
	attemptKey := fmt.Sprintf("login_attempts:%s", email)
	lockKey := fmt.Sprintf("login_lock:%s", email)

	isLocked, err := redisClient.Exists(ctx, lockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check login lock: %w", err)
	}

	if isLocked == 1 {
		return fmt.Errorf("account is locked due to too many login attempts")
	}

	// Use a pipeline to make INCR and EXPIRE atomic
	pipe := redisClient.Pipeline()
	attemptsCmd := pipe.Incr(ctx, attemptKey)
	pipe.Expire(ctx, attemptKey, LoginWindow)
	_, err = pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to track login attempts: %w", err)
	}

	attempts := attemptsCmd.Val()

	if attempts >= MaxLoginAttempts {
		redisClient.Set(ctx, lockKey, "1", LoginLockTime)
		return fmt.Errorf("account is locked due to too many login attempts")
	}

	return nil
}

func ResetLoginAttempts(email string, redisClient *redis.Client, ctx context.Context) {
	attemptKey := fmt.Sprintf("login_attempts:%s", email)
	lockKey := fmt.Sprintf("login_lock:%s", email)

	redisClient.Del(ctx, attemptKey, lockKey)
}
