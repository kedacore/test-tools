package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v9"
)

type redisMode string

const (
	standalone redisMode = "STANDALONE"
	cluster    redisMode = "CLUSTER"
	sentinel   redisMode = "SENTINEL"
)

const (
	cutset    = " "
	separator = ","
)

func GetRedisClient() redis.Cmdable {
	addrs := parseAddress()
	pass := os.Getenv("REDIS_PASSWORD")
	switch strings.ToUpper(os.Getenv("REDIS_MODE")) {
	case string(cluster):
		opts := redis.ClusterOptions{
			Addrs:    addrs,
			Password: pass,
		}
		return redis.NewClusterClient(&opts)
	case string(sentinel):
		sentinelPass := os.Getenv("REDIS_SENTINEL_PASSWORD")
		sentinelMaster := os.Getenv("REDIS_SENTINEL_MASTER")
		opts := redis.FailoverOptions{
			SentinelAddrs:    addrs,
			Password:         pass,
			MasterName:       sentinelMaster,
			SentinelPassword: sentinelPass,
		}
		return redis.NewFailoverClient(&opts)
	default:
		var addr string
		if len(addrs) > 0 {
			addr = addrs[0]
		}
		opts := redis.Options{
			Addr:     addr,
			Password: pass,
		}
		return redis.NewClient(&opts)
	}
}

func parseAddress() []string {
	addrsString := os.Getenv("REDIS_ADDRESS")
	if len(addrsString) == 0 {
		addrsString = os.Getenv("REDIS_ADDRESSES")
	}
	if len(addrsString) != 0 {
		return splitAndTrim(addrsString)
	}

	hostString := os.Getenv("REDIS_HOST")
	if len(hostString) == 0 {
		hostString = os.Getenv("REDIS_HOSTS")
	}
	hosts := splitAndTrim(hostString)

	portString := os.Getenv("REDIS_PORT")
	if len(portString) == 0 {
		portString = os.Getenv("REDIS_PORTS")
	}
	ports := splitAndTrim(portString)

	var addrs []string
	if len(hosts) != len(ports) {
		return addrs
	}
	for i := range hosts {
		addrs = append(addrs, fmt.Sprintf("%s:%s", hosts[i], ports[i]))
	}
	return addrs
}

func splitAndTrim(s string) []string {
	result := strings.Split(s, separator)
	for i := range result {
		result[i] = strings.Trim(result[i], cutset)
	}
	return result
}
