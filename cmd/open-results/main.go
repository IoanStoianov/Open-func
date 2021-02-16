package main

import (
	"log"
	"os"

	"github.com/IoanStoianov/Open-func/pkg/results"
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "127.0.0.1"
	}

	srv, err := results.NewServer(9000, redisHost)
	if err != nil {
		log.Fatalln(err)
	}

	go srv.SubscribeToRedis()

	log.Println("Staring on port 9000...")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
