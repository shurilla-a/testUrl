package main

import (
	"github.com/speps/go-hashids"
	"log"
	"time"
)

func main() {

	hd := hashids.NewData()
	hd.MinLength = 5
	//fmt.Println(hd)
	h, err := hashids.NewWithData(hd)
	if err != nil {
		log.Println(h, err)
	}
	//fmt.Println(h)
	timeNow := time.Now()
	urlId, err := h.Encode([]int{int(timeNow.Unix())})
	if err != nil {
		log.Println(urlId, err)
	}
}
