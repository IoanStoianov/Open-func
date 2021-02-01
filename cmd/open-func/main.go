package main

import (
	"log"

	"github.com/IoanStoianov/Open-func/pkg/openserver"
)

func main() {
	var port uint = 8090
	server := openserver.NewServer(port)

	done := make(chan bool)
	go func() {
		log.Printf("Starting server on port %d...\n", port)

		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}

		done <- true
	}()

	server.WaitShutdown()

	<-done

	log.Println("Done.")
}
