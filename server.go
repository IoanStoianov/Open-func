package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/headers", headers).Methods("GET")

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
