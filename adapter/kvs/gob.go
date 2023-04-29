package kvs

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/takokun778/fs-go-redis/domain/cache"
)

type GobFactory[T any] interface {
	Gob(string, *redis.Client) (*Gob[T], error)
}

var _ cache.Cache[any] = (*Gob[any])(nil)

type Gob[T any] struct {
	Prefix string
	Client *redis.Client
}

func (gb *Gob[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	key = fmt.Sprintf(prefix, gb.Prefix, key)

	str, err := gb.Client.Get(ctx, key).Result()
	if err != nil {
		return value, fmt.Errorf("failed to get cache: %w", err)
	}

	buf := bytes.NewBuffer([]byte(str))

	_ = gob.NewDecoder(buf).Decode(&value)

	return value, nil
}

func (gb *Gob[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	buf := bytes.NewBuffer(nil)

	_ = gob.NewEncoder(buf).Encode(&value)

	key = fmt.Sprintf(prefix, gb.Prefix, key)

	if err := gb.Client.Set(ctx, key, buf.Bytes(), ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (gb *Gob[T]) Del(ctx context.Context, key string) error {
	key = fmt.Sprintf(prefix, gb.Prefix, key)

	if _, err := gb.Client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to del cache: %w", err)
	}

	return nil
}
