package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
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

func getRedisClient() *redis.Client {
	addrs := parseAddress()
	pass := os.Getenv("REDIS_PASSWORD")
	sentinelPass := os.Getenv("REDIS_SENTINEL_PASSWORD")
	sentinelMaster := os.Getenv("REDIS_SENTINEL_MASTER")
	opts := redis.FailoverOptions{
		SentinelAddrs:    addrs,
		Password:         pass,
		MasterName:       sentinelMaster,
		SentinelPassword: sentinelPass,
	}
	return redis.NewFailoverClient(&opts)
}

func writeToRedisList() error {
	client := getRedisClient()
	list := os.Getenv("LIST_NAME")
	itemCount, err := strconv.ParseInt(os.Getenv("NO_LIST_ITEMS_TO_WRITE"), 10, 32)
	if err != nil {
		return fmt.Errorf("number of items to write should be a number: %s", err.Error())
	}

	for i := 0; i < int(itemCount); i++ {
		x := client.LPush(context.Background(), list, i)
		if x.Err() != nil {
			return fmt.Errorf("failed to write to redis list: %s", x.Err().Error())
		}
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}

func readFromRedisList() error {
	client := getRedisClient()
	list := os.Getenv("LIST_NAME")

	waitTime, err := strconv.ParseInt(os.Getenv("READ_PROCESS_TIME"), 10, 32)
	if err != nil {
		return fmt.Errorf("read process time should be a number: %s", err.Error())
	}

	for {
		len, err := client.LLen(context.Background(), list).Result()
		if err != nil {
			return err
		}
		if len > 0 {
			x := client.LPop(context.Background(), list)
			if x.Err() != nil {
				return fmt.Errorf("failed to read from redis list: %s", x.Err().Error())
			}
		}
		time.Sleep(time.Millisecond * time.Duration(waitTime))
	}
}

func main() {
	action := ""
	if len(os.Args) > 0 {
		action = os.Args[1]
	}
	if action == "write" {
		err := writeToRedisList()
		if err != nil {
			log.Fatalf("write to redis list failed: %v\n", err)
		}
		log.Println("write to redis list is successful")
	} else if action == "read" {
		err := readFromRedisList()
		if err != nil {
			log.Fatalf("read from redis list failed: %v\n", err)
		}
		log.Println("read from redis list is successful")
	} else {
		log.Printf("unknown action: %s\n", action)
	}
}
