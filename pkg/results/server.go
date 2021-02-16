package results

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type ResultsServer struct {
	http.Server
	redisClient *redis.Client
}

func NewServer(addr uint) (*ResultsServer, error) {
	r := mux.NewRouter().PathPrefix("/results").Subrouter() // TODO: shouldn't be hardcoded

	r.HandleFunc("/ping", pong).Methods("GET")

	redis := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	s := &ResultsServer{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", addr),
			Handler: r,
		},
		redisClient: redis,
	}

	return s, nil
}

func (s *ResultsServer) SubscribeToRedis() {
	pubsub := s.redisClient.Subscribe(context.Background(), "results")

	ch := pubsub.Channel()

	for msg := range ch {
		log.Println(msg)
	}
}

func pong(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
