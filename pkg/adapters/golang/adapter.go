package golang

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var action func(interface{}) string

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func triggerFunc(w http.ResponseWriter, r *http.Request) {

	var payload interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	resp := action(payload)
	fmt.Fprintf(w, resp)
}

//TriggerListener starts server and waits for triggers
func TriggerListener(f func(interface{}) string) {
	port := os.Getenv("OPEN_FUNC_PORT")
	if port == "" {
		port = "3014"
	}
	action = f

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/triggerHttp", triggerFunc)

	http.ListenAndServe(":"+port, nil)
}
