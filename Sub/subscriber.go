package main

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
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

func (s *Service) Subscribe(key string, msg chan []byte) error {
	rc := s.pool.Get()
	psc := redis.PubSubConn{Conn: rc}
	if err := psc.PSubscribe(key); err != nil {
		return err
	}

	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.PMessage:
				msg <- v.Data
			}
		}
	}()
	return nil
}

func main() {

	svc := New(&NewInput{
		RedisURL: "127.0.0.1:6379",
	})
	channel := fmt.Sprintf("mensaje", time.Now().Add(10*time.Second).String())
	val := time.Now().String()

	reply := make(chan []byte)
	err := svc.Subscribe(channel, reply)
	if err != nil {
		log.Fatal(err)
	}

	msg := <-reply
	if string(msg) != val {
		log.Fatal("expected correct reply message")
	}
	log.Printf("recieved %q", string(msg))

}
