package results

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IoanStoianov/Open-func/pkg/results/repo"
	"github.com/IoanStoianov/Open-func/pkg/types"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// Server - TODO
type Server struct {
	http.Server
	redisClient *redis.Client
	repo        repo.Results
}

// NewServer - ResultsServer factory
func NewServer(addr uint, redisHost string, mongoHost string) (*Server, error) {
	resultsRepo, err := repo.CreateMongoClient(mongoHost)
	if err != nil {
		return nil, err
	}

	redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", redisHost),
	})

	s := &Server{
		Server: http.Server{
			Addr: fmt.Sprintf(":%d", addr),
		},
		redisClient: redis,
		repo:        resultsRepo,
	}

	r := mux.NewRouter().PathPrefix("/results").Subrouter() // TODO: shouldn't be hardcoded

	r.HandleFunc("/ping", pong).Methods("GET")
	r.HandleFunc("/latest", s.getLatest).Methods("GET")

	s.Handler = r

	return s, nil
}

func (s *Server) getLatest(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["funcName"]
	if !ok {
		log.Println("Missing 'funcName' URL param")
		http.Error(w, "Missing 'funcName' param", http.StatusBadRequest)
		return
	}

	result, err := s.repo.GetRecords(name[0], 1)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading records", http.StatusInternalServerError)
		return
	}

	if len(result) == 0 {
		http.Error(w, "No records found", http.StatusNotFound)
		return
	}

	serialized, err := json.Marshal(result[0])
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(serialized)
}

// SubscribeToRedis initiates the service
func (s *Server) SubscribeToRedis() {
	pubsub := s.redisClient.Subscribe(context.Background(), "results")

	ch := pubsub.Channel()

	for msg := range ch {
		log.Println(msg)

		var result types.FuncResult
		if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
			log.Println(err)
			continue
		}

		if err := s.repo.AddRecord(&result); err != nil {
			log.Println(err)
		}
	}
}

func pong(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
