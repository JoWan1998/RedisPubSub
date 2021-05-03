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

type server struct {
	nc *redis.Conn
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.TODO()

	errs := rdb.Publish(ctx, "mychannel1", []byte(data)).Err()
	if errs != nil {
		panic(err)
	}
	duration := time.Since(requestAt)

	fmt.Fprintf(w, "Task scheduled in %+v\nResponse: %v\n", duration, string(response.Data))
}

func main() {
	http.HandleFunc("/", createTask)
	fmt.Println("Server listening on port 8080...")
	if errors := http.ListenAndServe(":8080", nil); errors != nil {
		log.Fatal(errors)
	}
}
