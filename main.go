package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"log"
	"net/http"
	"time"
)

func RedisConnect() *redis.Client {
	rdbc := redis.NewClient(&redis.Options{
		DB:       0,
		Addr:     "172.31.201.78:6379",
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
func Redirect(w http.ResponseWriter, req *http.Request) {
	rdbc := RedisConnect()
	params := mux.Vars(req)
	key := params["key"]
	url, err := rdbc.Get(key).Result()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(key, url)
	http.Redirect(w, req, url, 301)
}

func Create(w http.ResponseWriter, req *http.Request) {
	rdbc := RedisConnect()
	req.ParseForm()
	url := req.Form["url"][0]
	key := GenerateShortIDurl(5)
	rdbc.Set(key, url, 0).Result()
	fmt.Println(key, url)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{key}", Redirect).Methods("GET")
	router.HandleFunc("/create", Create).Methods("POST")
	http.ListenAndServe(":3128", router)
}
