package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func CreateRedisClient() {
	opt, err := redis.ParseURL("redis://localhost:6364/0")
	if err != nil {
		panic(err)
	}

	redis := redis.NewClient(opt)
	Redis = redis
	log.Println("Create connection...")
}

func publishMessage(message []byte) {
	err := Redis.Publish(context.Background(), "mensajes", message).Err()

	if err != nil {
		log.Println(err)
	}
}

func subscribeMessages() {
	pubsub := Redis.Subscribe(context.Background(), "mensajes")
	log.Println("subscriber listen on... ")
	ch := pubsub.Channel()

	for msg := range ch {
		log.Println("Mensaje: ", []byte(msg.Payload))
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {

	requestAt := time.Now()
	w.Header().Set("Content-Type", "application/json")
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	log.Println("Error Parseando JSON: ", err)
	data, err := json.Marshal(body)
	log.Println("Error Reading Body: ", err)
	fmt.Println(string(data))
	publishMessage(data)
	duration := time.Since(requestAt)
	fmt.Fprintf(w, "Task scheduled in %+v", duration)
}

func main() {
	CreateRedisClient()
	http.HandleFunc("/", createTask)
	fmt.Println("Server listening on port 8080...")
	if errors := http.ListenAndServe(":8080", nil); errors != nil {
		log.Fatal(errors)
	}
	subscribeMessages()
}
