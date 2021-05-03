package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
)

func createTask(w http.ResponseWriter, r *http.Request) {

	requestAt := time.Now()
	w.Header().Set("Content-Type", "application/json")
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	log.Println("Error Parseando JSON: ", err)
	data, err := json.Marshal(body)
	log.Println("Error Reading Body: ", err)
	fmt.Println(string(data))

	svc := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	errs := svc.Publish("mensaje", string(data))

	if errs != nil {
		log.Fatal(errs)
	}

	duration := time.Since(requestAt)
	fmt.Fprintf(w, "Task scheduled in %+v", duration)
}

func main() {
	http.HandleFunc("/", createTask)
	fmt.Println("Server listening on port 8080...")
	if errors := http.ListenAndServe(":8080", nil); errors != nil {
		log.Fatal(errors)
	}
}
