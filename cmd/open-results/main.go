package main

import (
	"log"

	"github.com/IoanStoianov/Open-func/pkg/results"
)

func main() {
	srv, err := results.NewServer(9000)
	if err != nil {
		log.Println(err)
	}

	go srv.SubscribeToRedis()

	log.Println("Staring on port 9000...")

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
