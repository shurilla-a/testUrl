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

// делать проверку на сбодность Ключа
func GenerateShortIDurl(lenghIDurl int) string {
	rdbc := RedisConnect()
newgenerate:
	hd := hashids.NewData()
	hd.MinLength = lenghIDurl
	hesh, err := hashids.NewWithData(hd)
	if err != nil {
		log.Println(err)
	}
	timeNow := time.Now()
	key, err := hesh.Encode([]int{int(timeNow.Unix())})
	if err != nil {
		log.Println(err)
	}
	value, err := rdbc.Get(key).Result()
	if err == redis.Nil {
		goto finish
	} else if err != nil {
		log.Println(err)
	} else {
		log.Println(value, "Зачение существует ")
		goto newgenerate
	}

finish:
	defer rdbc.Close()
	return key
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
	defer rdbc.Close()
}

func Create(w http.ResponseWriter, req *http.Request) {
	rdbc := RedisConnect()
	req.ParseForm()
	url := req.Form["url"][0]
	key := GenerateShortIDurl(5)
	rdbc.Set(key, url, 96*time.Hour).Result()
	//fmt.Println(key, url)
	// дописать отдачу в короткой ссылки в curl
	fmt.Fprintln(w, "http://127.0.0.1:3128/"+key)
	defer rdbc.Close()
}

func main11() {
	router := mux.NewRouter()
	router.HandleFunc("/{key}", Redirect).Methods("GET")
	router.HandleFunc("/create", Create).Methods("POST")
	http.ListenAndServe(":3128", router)
}
