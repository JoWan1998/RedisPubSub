package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.TODO()
	pubsub := rdb.Subscribe(ctx, "mensajes")

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}
