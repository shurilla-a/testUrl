package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/speps/go-hashids"
	"log"
	"time"
)

func RedisConnect() *redis.Client {
	rdbc := redis.NewClient(&redis.Options{
		DB:       0,
		Addr:     "172.XX.XXX.XXX:6379",
		Password: "",
	})
	pong, err := rdbc.Ping().Result()
	if err != nil {
		log.Panicln(pong, err)
	}
	return rdbc
}

func GenerateShortIDurl(lenghIDurl int) string {
	hd := hashids.NewData()
	hd.MinLength = lenghIDurl
	h, err := hashids.NewWithData(hd)
	if err != nil {
		log.Println(h, err)
	}
	timeNow := time.Now()
	urlId, err := h.Encode([]int{int(timeNow.Unix())})
	if err != nil {
		log.Println(urlId, err)
	}
	return urlId
}

func main() {
	rdbc := RedisConnect()
	key := GenerateShortIDurl(5)
	url, err := rdbc.Get(key).Result()
	if err != nil {
		log.Println(err)
	}

	if len(url) < 1 {
		fmt.Println("Ничего не нашел", url)
		fmt.Println(key)
		rdbc.Set(key, "https://yandex.ru", 0).Result()
	}
}
