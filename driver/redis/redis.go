package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/takokun778/fs-go-redis/adapter/kvs"
)

var (
	_ kvs.Factory[any]       = (*Redis[any])(nil)
	_ kvs.JSONFactory[any]   = (*Redis[any])(nil)
	_ kvs.Base64Factory[any] = (*Redis[any])(nil)
	_ kvs.GobFactory[any]    = (*Redis[any])(nil)
)

type Redis[T any] struct{}

func New[T any]() *Redis[T] {
	return &Redis[T]{}
}

func (rds *Redis[T]) KVS(
	prefix string,
	client *redis.Client,
) (*kvs.KVS[T], error) {
	return &kvs.KVS[T]{
		Prefix: prefix,
		Client: client,
	}, nil
}

func (rds *Redis[T]) Base64(
	prefix string,
	client *redis.Client,
) (*kvs.Base64[T], error) {
	return &kvs.Base64[T]{
		Prefix: prefix,
		Client: client,
	}, nil
}

func (rds *Redis[T]) JSON(
	prefix string,
	client *redis.Client,
) (*kvs.JSON[T], error) {
	return &kvs.JSON[T]{
		Prefix: prefix,
		Client: client,
	}, nil
}

func (rds *Redis[T]) Gob(
	prefix string,
	client *redis.Client,
) (*kvs.Gob[T], error) {
	return &kvs.Gob[T]{
		Prefix: prefix,
		Client: client,
	}, nil
}

func NewRedis(env string, url string) (*redis.Client, error) {
	var opt *redis.Options

	if env == "upstash" {
		var err error

		opt, err = redis.ParseURL(url)

		if err != nil {
			return nil, fmt.Errorf("failed to parse redis url: %w", err)
		}
	} else {
		opt = &redis.Options{
			Addr:     url,
			Password: "",
			DB:       0,
		}
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
