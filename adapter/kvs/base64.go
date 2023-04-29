package kvs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/takokun778/fs-go-redis/domain/cache"
)

type Base64Factory[T any] interface {
	Base64(string, *redis.Client) (*Base64[T], error)
}

var _ cache.Cache[any] = (*Base64[any])(nil)

type Base64[T any] struct {
	Prefix string
	Client *redis.Client
}

func (bs *Base64[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	key = fmt.Sprintf(prefix, bs.Prefix, key)

	str, err := bs.Client.Get(ctx, key).Result()
	if err != nil {
		return value, fmt.Errorf("failed to get cache: %w", err)
	}

	dec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return value, fmt.Errorf("failed to decode base64: %w", err)
	}

	if err := json.Unmarshal(dec, &value); err != nil {
		return value, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return value, nil
}

func (bs *Base64[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	key = fmt.Sprintf(prefix, bs.Prefix, key)

	enc := base64.StdEncoding.EncodeToString(val)

	if err := bs.Client.Set(ctx, key, enc, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (bs *Base64[T]) Del(ctx context.Context, key string) error {
	key = fmt.Sprintf(prefix, bs.Prefix, key)

	if _, err := bs.Client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to del cache: %w", err)
	}

	return nil
}
