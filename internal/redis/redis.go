package redis

import (
	"context"
	"time"

	"ptcg_trader/internal/config"
	"ptcg_trader/internal/errors"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

// Redis ...
type Redis interface {
	redis.Cmdable

	RedisLock(ctx context.Context, key, lockerID string, expireTime time.Duration) (bool, error)
	RedisUnlock(ctx context.Context, key, lockerID string) error
}

type _redis struct {
	// ctx context.Context // field `ctx` is unused (unused)

	*redis.Client
}

// RedisLock ...
func (r *_redis) RedisLock(ctx context.Context, key, lockerID string, expireTime time.Duration) (bool, error) {
	timeout := time.NewTicker(5 * time.Second)

	for {
		select {
		default:
			ok, err := r.SetNX(ctx, key, lockerID, expireTime).Result()
			if err != nil {
				return false, errors.Wrapf(errors.ErrInternalError, "RedisLock SetNX Error: %v", err.Error())
			}
			if !ok {
				time.Sleep(1 * time.Microsecond)
				continue
			}
			return true, nil

		case <-timeout.C:
			return false, errors.Wrapf(errors.ErrDataConflict, "RedisLock SetNX with key %v timeout", key)
		}
	}
}

// RedisUnlock ...
func (r *_redis) RedisUnlock(ctx context.Context, key, lockerID string) error {
	result, err := r.Get(ctx, key).Result()
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "RedisUnlock Get %v Error: %v", key, err.Error())
	}
	if result == lockerID {
		err := r.Del(ctx, key).Err()
		if err != nil {
			return errors.Wrapf(errors.ErrInternalError, "RedisUnlock Del %v Error: %v", key, err.Error())
		}
	}

	return nil
}

// NewRedis ...
func NewRedis(cfg config.RedisConfig) (Redis, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	if len(cfg.Addresses) == 0 {
		return nil, errors.Wrap(errors.ErrInternalError, "redis config address is empty")
	}

	var client *redis.Client
	err := backoff.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr:       cfg.Addresses[0],
			Password:   cfg.Password,
			MaxRetries: cfg.MaxRetries,
			PoolSize:   cfg.PoolSize,
			DB:         cfg.DB,
		})
		err := client.Ping(context.Background()).Err()
		if err != nil {
			log.Warn().Msgf("ping occurs error after connecting to redis: %s", err)
			return errors.Wrapf(errors.ErrInternalError, "ping occurs error after connecting to redis: %s", err)
		}
		return nil
	}, bo)
	if err != nil {
		return nil, err
	}

	return &_redis{Client: client}, nil
}
