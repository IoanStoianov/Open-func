package golang

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var action func(io.ReadCloser) string

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func triggerFunc(w http.ResponseWriter, r *http.Request) {
	resp := action(r.Body)
	fmt.Fprintf(w, resp)
}

//TriggerListener starts server and waits for triggers
func TriggerListener(f func(io.ReadCloser) string) {
	port := os.Getenv("OPEN_FUNC_PORT")
	if port == "" {
		port = "3014"
	}
	action = f

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/triggerHttp", triggerFunc)

	log.Printf("Starting server on port %s...\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
	}
}
