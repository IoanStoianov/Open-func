package openserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IoanStoianov/Open-func/pkg/k8s"
	"github.com/IoanStoianov/Open-func/pkg/k8s/client"
	"github.com/IoanStoianov/Open-func/pkg/triggers"
	models "github.com/IoanStoianov/Open-func/pkg/types"

	"k8s.io/apimachinery/pkg/util/uuid"

	"github.com/gorilla/mux"
)

type resource struct {
	deploymentName string
	serviceName    string
}

// OpenServer is the core of Open-func
type OpenServer struct {
	http.Server
	shutdowReq chan bool // can be used for remote shutdown
	resources  sync.Map
}

// NewServer - server factory
func NewServer(addr uint) (*OpenServer, error) {
	if addr < 1001 {
		return nil, errors.New("Invalid port number")
	}

	s := &OpenServer{
		Server: http.Server{
			Addr:         fmt.Sprintf(":%d", addr),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdowReq: make(chan bool),
	}

	r := mux.NewRouter()

	r.HandleFunc("/ping", pong).Methods("GET")
	r.HandleFunc("/prepare", s.PrepareFunc).Methods("POST")
	r.HandleFunc("/trigger", triggers.HTTPTriggerRedirect).Methods("POST")
	r.HandleFunc("/coldTrigger", triggers.HTTPColdTrigger).Methods("POST")

	// TODO: frontend should be extracted as a standalone service
	staticFileDirectory := http.Dir("./web/open-func/build/")
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	s.Handler = r

	return s, nil
}

// WaitShutdown blocks the main thread until an interrupt is received to
// initiate graceful shutdown.
func (s *OpenServer) WaitShutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.cleanUp()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Error on shutdown: %v\n", err)
	}
}

// PrepareFunc readies a deployment and service for hot execution
func (s *OpenServer) PrepareFunc(w http.ResponseWriter, r *http.Request) {
	var payload models.FuncTrigger

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	client := client.InCluster()

	deployName, err := k8s.CreateDeployment(client, payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	serviceName, err := k8s.CreateService(client, payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	s.resources.Store(uuid.NewUUID(), resource{deployName, serviceName})

	w.WriteHeader(204)
}

func (s *OpenServer) cleanUp() {
	client := client.InCluster()

	s.resources.Range(func(key, value interface{}) bool {
		r := value.(resource)

		log.Printf("Deleting pair %s %s ...", r.deploymentName, r.serviceName)

		if err := k8s.DeleteDeployment(client, r.deploymentName); err != nil {
			log.Printf("Error deleting deployment %s:\n%v", r.deploymentName, err)
		}

		if err := k8s.DeleteService(client, r.serviceName); err != nil {
			log.Printf("Error deleting service %s:\n%v", r.serviceName, err)
		}

		return true
	})
}

func pong(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
