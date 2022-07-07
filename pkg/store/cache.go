package store

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var (
	Redis *redis.Client
)

func initd() {
	Redis = redis.NewClient(&redis.Options{
		//Addr:     "host.docker.internal:6379",
		Addr:     "vlin2-dev.redis.rds.aliyuncs.com:6379",
		Password: "Aa12345!", // no password set
		DB:       30,         // use default DB
	})

	pong, err := Redis.Ping().Result()
	fmt.Println(pong, err)
}

func Scan(match string, deal func(key string)) {
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = Redis.Scan(cursor, match, 10).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		if cursor == 0 {
			break
		}
		for _, k := range keys {
			log.Println(k)
			deal(k)
		}
	}
}

func Setnx(key string, second time.Duration) bool {
	cmd := Redis.SetNX(key, "1", second)
	b, err := cmd.Result()
	if err != nil {
		log.Println("setnx err:", err)
	}
	return b
}

func Get(key string) string {
	cmd := Redis.Get(key)
	s, err := cmd.Result()
	if err != nil && err.Error() != "redis: nil" {
		log.Println("Get err:", err)
	}
	return s
}

func Push(value string) {
	cmd := Redis.LPush("info_result", value)
	_, err := cmd.Result()
	if err != nil {
		log.Println("Push err:", err)
	}
}

func Pop() string {
	cmd := Redis.RPop("info_result")
	s, err := cmd.Result()
	if err != nil && err.Error() != "redis: nil" {
		log.Println("Pop err:", err)
	}
	return s
}
