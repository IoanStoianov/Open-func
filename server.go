package main

import (
	"fmt"
	"net/http"
	"open-func/k8s"
	"open-func/k8s/client"
	"open-func/triggers"
	"open-func/types"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func deployNodejs(w http.ResponseWriter, req *http.Request) {
	client := client.OutCluster()
	k8s.CreateDeployment(client, types.FuncTrigger{FuncPort: 3000})
	k8s.CreateService(client, types.FuncTrigger{FuncPort: 3000})
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/deploy", deployNodejs).Methods("GET")
	r.HandleFunc("/test", triggers.HTTPTriggerRedirect).Methods("POST")

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
