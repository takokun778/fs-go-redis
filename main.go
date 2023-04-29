package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/takokun778/fs-go-redis/domain/model"
	"github.com/takokun778/fs-go-redis/driver/redis"
)

const ttl = 60 * time.Second

func main() {
	// rds, err := redis.NewRedis("docker", "localhost:6379")
	rds, err := redis.NewRedis("upstash", os.Getenv("UPSTASH_URL"))
	if err != nil {
		panic(err)
	}

	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	key := model.Key{
		ID:        uuid.New().String(),
		PublicKey: &k.PublicKey,
	}

	ctx := context.Background()

	fmt.Println("---------- kvs ----------")

	mdl := model.Model{
		ID:   uuid.New().String(),
		Name: "model",
	}

	kvs, err := redis.New[model.Model]().KVS("kvs", rds)
	if err != nil {
		panic(err)
	}

	if err := kvs.Set(ctx, mdl.ID, mdl, ttl); err != nil {
		panic(err)
	}

	got, err := kvs.Get(ctx, mdl.ID)
	if err != nil {
		panic(err)
	}

	log.Printf("got: %+v", got)

	fmt.Println("---------- json ----------")

	js, err := redis.New[model.Key]().JSON("js", rds)
	if err != nil {
		panic(err)
	}

	if err := js.Set(ctx, key.ID, key, ttl); err != nil {
		panic(err)
	}

	jsgot, err := js.Get(ctx, key.ID)
	if err != nil {
		panic(err)
	}

	log.Printf("jsgot: %+v", jsgot.PublicKey)

	fmt.Println("---------- base64 ----------")

	bs, err := redis.New[model.Key]().Base64("bs", rds)
	if err != nil {
		panic(err)
	}

	if err := bs.Set(ctx, key.ID, key, ttl); err != nil {
		panic(err)
	}

	bsgot, err := bs.Get(ctx, key.ID)
	if err != nil {
		panic(err)
	}

	log.Printf("bsgot: %+v", bsgot.PublicKey)

	fmt.Println("---------- gob ----------")

	gb, err := redis.New[model.Key]().Gob("gb", rds)
	if err != nil {
		panic(err)
	}

	if err := gb.Set(ctx, key.ID, key, ttl); err != nil {
		panic(err)
	}

	gbgot, err := gb.Get(ctx, key.ID)
	if err != nil {
		panic(err)
	}

	log.Printf("gbgot: %+v", gbgot.PublicKey)
}
