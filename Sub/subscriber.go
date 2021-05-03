package main

import (
	"log"

	"github.com/go-redis/redis/v7"
)

func main() {

	svc := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	sub := svc.pubsub()
	sub.subscribe("mensaje")
	for {
		msg := sub.get_message()
		log.Println("Message: ", msg)
	}

}
