package helper

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "127.0.0.1:6379",
	DB:   0,
})

func SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	data, err := json.Marshal(value)

	if err != nil {
		log.Println("ERROR: Json parsing error:", err)
		return err
	}

	err = redisClient.Set(ctx, key, data, expiration).Err()
	if err != nil {
		log.Println("ERROR: write redis error:", err)
		return err
	}

	log.Println("INFO: write redis - Key:", key, "TTL:", expiration)

	return nil
}

func GetCache(key string, dest interface{}) error {
	ctx := context.Background()
	data, err := redisClient.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		log.Println("WARNING: key cannot found:", key)
		return err
	} else if err != nil {
		log.Println("ERROR: key cannot read:", err)
		return err
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		log.Println("ERROR: Json parsing error:", err)
		return err
	}

	log.Println("INFO: read from cache - Key:", key)

	return nil
}

func DeleteCache(key string) error {
	ctx := context.Background()

	log.Println("DEBUG: delete from cache - Key:", key)

	err := redisClient.Del(ctx, key).Err()

	if err != nil {
		log.Println("ERROR: cache cannot delete - Key:", key, "Error:", err)
	} else {
		log.Println("INFO: cache delete properly - Key:", key)
	}

	return err
}
