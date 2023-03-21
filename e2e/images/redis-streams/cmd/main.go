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

func redisStreamConsumer(ctx context.Context) error {
	client := GetRedisClient()
	stream := os.Getenv("REDIS_STREAM_NAME")
	del, _ := strconv.ParseBool(os.Getenv("DELETE_MESSAGES"))

	consumerGroup := os.Getenv("REDIS_STREAM_CONSUMER_GROUP_NAME")
	_, err := client.XGroupCreate(context.Background(), stream, consumerGroup, "0").Result()
	if err != nil && !strings.Contains(err.Error(), "Consumer Group name already exists") {
		return fmt.Errorf("failed to create consumer group: %s", err.Error())
	}

	var ids []string

	go func(ids *[]string) {
		var pendingQuantity int64

		for {
			pending, err := client.XPending(ctx, stream, consumerGroup).Result()
			if err != nil {
				log.Printf("failed to get pending entries: %s", err.Error())
				time.Sleep(2 * time.Second)
				continue
			}

			if pending.Count == 0 {
				time.Sleep(2 * time.Second)
				continue
			}

			if pending.Count == pendingQuantity {
				time.Sleep(10 * time.Second)
				log.Printf("Acking %d entries...\n", len(*ids))
				client.XAck(ctx, stream, consumerGroup, *ids...)
				if del {
					log.Printf("Deleting %d entries...\n", len(*ids))
					client.XDel(ctx, stream, *ids...)
				}
			}

			pendingQuantity = pending.Count
			time.Sleep(2 * time.Second)
		}
	}(&ids)

	for {
		read := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: "consumer",
			Streams:  []string{stream, ">"},
			Count:    1,
		})
		res, err := read.Result()
		if err != nil {
			return fmt.Errorf("failed to XREADGROUP from redis stream: %s", err)
		}
		if len(res[0].Messages) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		msg := res[0].Messages[0]
		log.Printf("read message no %s", msg.Values["no"])
		ids = append(ids, msg.ID)
		time.Sleep(500 * time.Millisecond)
	}
}

func redisStreamProducer(ctx context.Context) error {
	client := GetRedisClient()
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
