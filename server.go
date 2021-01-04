package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"open-func/k8s"
	"open-func/k8s/client"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func deployNodejs(w http.ResponseWriter, req *http.Request) {
	client := client.OutCluster()
	k8s.Deploy(client)
	k8s.CreateService(client)
}

type Message struct {
	ID   int64       `json:"id"`
	Name interface{} `json:"name"`
}

// curl localhost:8000 -d '{"name":"Hello"}'
func Cleaner(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/deploy", deployNodejs).Methods("GET")
	r.HandleFunc("/test", Cleaner).Methods("POST")

	staticFileDirectory := http.Dir("./assets/open-func/build/")
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	return r
}

func main() {
	// client.OutCluster()
	router := newRouter()
	fmt.Printf("Starting server at port 8090\n")
	http.ListenAndServe(":8090", router)
}
