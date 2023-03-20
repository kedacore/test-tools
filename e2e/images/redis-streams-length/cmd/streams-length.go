package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
)

func getRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDRESS")
	if len(addr) == 0 {
		addr = os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	}
	pass := os.Getenv("REDIS_PASSWORD")
	opts := redis.Options{
		Addr:     addr,
		Password: pass,
	}
	return redis.NewClient(&opts)
}

func redisStreamConsumer(ctx context.Context) error {
	client := getRedisClient()
	stream := os.Getenv("REDIS_STREAM_NAME")

	consumerGroup := os.Getenv("REDIS_STREAM_CONSUMER_GROUP_NAME")
	_, err := client.XGroupCreate(context.Background(), stream, consumerGroup, "0").Result()
	if err != nil && !strings.Contains(err.Error(), "Consumer Group name already exists") {
		return fmt.Errorf("failed to create consumer group: %s", err.Error())
	}

	for {
		length, err := client.XLen(ctx, stream).Result()
		if err != nil {
			return err
		}
		if length == 0 {
			continue
		}

		x := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: "consumer",
			Streams:  []string{stream, ">"},
			Count:    1,
		})
		if x.Err() != nil {
			return fmt.Errorf("failed to XREADGROUP from redis stream: %s", x.Err().Error())
		}

		res, err := x.Result()
		if err != nil {
			return fmt.Errorf("failed to fetch XREADGROUP result: %v", err)
		}

		log.Printf("read message no %s", res[0].Messages[0].Values["no"])
		client.XDel(ctx, stream, res[0].Messages[0].ID)
		time.Sleep(500 * time.Millisecond)
	}
}

func redisStreamProducer(ctx context.Context) error {
	client := getRedisClient()
	stream := os.Getenv("REDIS_STREAM_NAME")

	count, err := strconv.ParseInt(os.Getenv("NUM_MESSAGES"), 10, 32)
	if err != nil {
		return fmt.Errorf("number of items to write should be a number: %s", err.Error())
	}

	for i := 0; i < int(count); i++ {
		x := client.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			Values: map[string]interface{}{"no": i + 1},
		})
		if x.Err() != nil {
			return fmt.Errorf("failed to write to redis stream: %s", x.Err().Error())
		}
	}
	return nil
}

func main() {
	ctx := context.Background()
	mode := ""
	if len(os.Args) > 0 {
		mode = os.Args[1]
	}
	switch mode {
	case "consumer":
		if err := redisStreamConsumer(ctx); err != nil {
			log.Fatalf("read from redis stream failed: %v\n", err)
		}
		log.Println("read from redis stream is successful")
	case "producer":
		if err := redisStreamProducer(ctx); err != nil {
			log.Fatalf("write to redis stream failed: %v\n", err)
		}
		log.Println("write to redis stream is successful")
	default:
		log.Fatalf("unknown mode: %s\n", mode)
	}
}
