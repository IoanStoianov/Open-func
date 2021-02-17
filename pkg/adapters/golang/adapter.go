package golang

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
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

//ColdTriggerListener executes function immediately
func ColdTriggerListener(action func(io.ReadCloser) string) {
	payload := os.Getenv("PAYLOAD")
	redisURL := os.Getenv("REDIS_URL")

	data := ioutil.NopCloser(strings.NewReader(payload))
	resp := action(data)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Publish(context.TODO(), "mychannel1", "hello").Err()
	if err != nil {
		log.Println(err)
	}

	err = rdb.Publish(context.TODO(), "ketap", resp).Err()
	if err != nil {
		log.Println(err)
	}
}
