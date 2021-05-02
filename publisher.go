package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

const key = "myJobQueue"

func main() {

	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	fmt.Println("Waiting for jobs on jobQueue: ", key)

	go func() {
		for {
			result, err := c.BLPop(0*time.Second, key).Result()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Executing job: ", result[1])
		}
	}()

	select {}
}
