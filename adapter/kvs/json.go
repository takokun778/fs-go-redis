package kvs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/takokun778/fs-go-redis/domain/cache"
)

type JSONFactory[T any] interface {
	JSON(string, *redis.Client) (*JSON[T], error)
}

var _ cache.Cache[any] = (*JSON[any])(nil)

type JSON[T any] struct {
	Prefix string
	Client *redis.Client
}

func (js *JSON[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	key = fmt.Sprintf(prefix, js.Prefix, key)

	str, err := js.Client.Get(ctx, key).Result()
	if err != nil {
		return value, fmt.Errorf("failed to get cache: %w", err)
	}

	if err := json.Unmarshal([]byte(str), &value); err != nil {
		return value, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return value, nil
}

func (js *JSON[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	key = fmt.Sprintf(prefix, js.Prefix, key)

	if err := js.Client.Set(ctx, key, val, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (js *JSON[T]) Del(ctx context.Context, key string) error {
	key = fmt.Sprintf(prefix, js.Prefix, key)

	if _, err := js.Client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to del cache: %w", err)
	}

	return nil
}
