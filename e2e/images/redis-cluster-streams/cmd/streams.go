package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v9"
)

func splitAndTrim(s, sep, toTrim string) []string {
	x := strings.Split(s, sep)
	for i := range x {
		x[i] = strings.Trim(x[i], toTrim)
	}
	return x
}

func parseAddress() []string {
	addrString := os.Getenv("REDIS_ADDRESSES")
	if len(addrString) != 0 {
		return splitAndTrim(addrString, ",", " ")
	}
	hostString := os.Getenv("REDIS_HOSTS")
	portString := os.Getenv("REDIS_PORTS")
	hosts := splitAndTrim(hostString, ",", " ")
	ports := splitAndTrim(portString, ",", " ")
	addrs := []string{}
	if len(hosts) != len(ports) {
		return addrs
	}
	for i := range hosts {
		addrs = append(addrs, fmt.Sprintf("%s:%s", hosts[i], ports[i]))
	}
	return addrs
}

func checkPendingEntries(c *redis.ClusterClient, stream, consumer string, ids *[]string) {
	for {
		pe, err := c.XPending(context.Background(), stream, consumer).Result()
		if err != nil {
			log.Printf("failed to get pending entries: %s", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}

		if pe.Count == 100 {
			// wait for other consumers to read pending entries.
			time.Sleep(20 * time.Second)
			log.Printf("ACKing %d entries...\n", len(*ids))
			c.XAck(context.Background(), stream, consumer, *ids...)
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func redisStreamConsumer() error {
	addrs := parseAddress()
	pass := os.Getenv("REDIS_PASSWORD")
	opts := redis.ClusterOptions{
		Addrs:    addrs,
		Password: pass,
	}
	client := redis.NewClusterClient(&opts)
	stream := os.Getenv("REDIS_STREAM_NAME")

	consumerGroup := os.Getenv("REDIS_STREAM_CONSUMER_GROUP_NAME")
	_, err := client.XGroupCreate(context.Background(), stream, consumerGroup, "0").Result()
	if err != nil && !strings.Contains(err.Error(), "Consumer Group name already exists") {
		return fmt.Errorf("failed to create consumer group: %s", err.Error())
	}

	pendingEntries := []string{}
	go checkPendingEntries(client, stream, consumerGroup, &pendingEntries)

	msgCount := 0
	for {
		length, err := client.XLen(context.Background(), stream).Result()
		if err != nil {
			return err
		}
		if length > 0 {
			x := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: "damn-you",
				Streams:  []string{stream, ">"},
				Count:    1,
				Block:    0,
			})
			if x.Err() != nil {
				return fmt.Errorf("failed to create consumer group to redis stream: %s", x.Err().Error())
			}

			res, err := x.Result()
			if err != nil {
				return fmt.Errorf("failed to read from redis stream: %v", err)
			}

			msgCount++
			log.Printf("read %d messages from stream\n", msgCount)
			pendingEntries = append(pendingEntries, res[0].Messages[0].ID)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func redisStreamProducer() error {
	addrs := parseAddress()
	pass := os.Getenv("REDIS_PASSWORD")
	opts := redis.ClusterOptions{
		Addrs:    addrs,
		Password: pass,
	}
	client := redis.NewClusterClient(&opts)
	stream := os.Getenv("REDIS_STREAM_NAME")

	count, err := strconv.ParseInt(os.Getenv("NUM_MESSAGES"), 10, 32)
	if err != nil {
		return fmt.Errorf("number of items to write should be a number: %s", err.Error())
	}

	for i := 0; i < int(count); i++ {
		x := client.XAdd(context.Background(), &redis.XAddArgs{
			Stream: stream,
			Values: map[string]interface{}{"key": "value"},
		})
		if x.Err() != nil {
			return fmt.Errorf("failed to write to redis stream: %s", x.Err().Error())
		}
	}
	return nil
}

func main() {
	mode := ""
	if len(os.Args) > 0 {
		mode = os.Args[1]
	}
	if mode == "consumer" {
		if err := redisStreamConsumer(); err != nil {
			log.Fatalf("read from redis stream failed: %v\n", err)
		}
		log.Println("read from redis stream is successful")
	} else if mode == "producer" {
		if err := redisStreamProducer(); err != nil {
			log.Fatalf("write to redis stream failed: %v\n", err)
		}
		log.Println("write to redis stream is successful")
	} else {
		log.Printf("unknown mode: %s\n", mode)
	}
}
