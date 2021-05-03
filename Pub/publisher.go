package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type Service struct {
	pool *redis.Pool
	conn redis.Conn
}

// NewInput input for constructor
type NewInput struct {
	RedisURL string
}

// New return new service
func New(input *NewInput) *Service {
	if input == nil {
		log.Fatal("input is required")
	}
	var redispool *redis.Pool
	redispool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", input.RedisURL)
		},
	}

	// Get a connection
	conn := redispool.Get()
	defer conn.Close()
	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		log.Fatalf("can't connect to the redis database, got error:\n%v", err)
	}

	return &Service{
		pool: redispool,
		conn: conn,
	}
}

func (s *Service) Publish(key string, value string) error {
	conn := s.pool.Get()
	conn.Do("PUBLISH", key, value)
	return nil
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

	svc := New(&NewInput{
		RedisURL: "127.0.0.1:6379",
	})

	errs := svc.Publish("test/foo", string(data))

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
