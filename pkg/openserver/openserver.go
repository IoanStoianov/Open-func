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

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"

	"github.com/gorilla/mux"
)

type resource struct {
	DeploymentName string `json:"deploymentName"`
	ServiceName    string `json:"serviceName"`
}

// OpenServer is the core of Open-func
type OpenServer struct {
	http.Server
	shutdowReq chan bool // can be used for remote shutdown
	resources  sync.Map
	client     *kubernetes.Clientset
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
		client:     client.InCluster(),
	}

	r := mux.NewRouter()

	r.HandleFunc("/ping", pong).Methods("GET")
	r.HandleFunc("/listSpecs", s.ListFuncSpecs).Methods("GET")
	r.HandleFunc("/prepare", s.PrepareFunc).Methods("POST")
	r.HandleFunc("/delete", s.DeleteFunc).Methods("DELETE")
	r.HandleFunc("/trigger", triggers.HTTPTriggerRedirect).Methods("POST")
	r.HandleFunc("/coldTrigger", s.ColdTriggerSpawn).Methods("POST")

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
	var payload models.FuncSpecs

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	deployName, err := k8s.CreateDeployment(s.client, payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	serviceName, err := k8s.CreateService(s.client, payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	s.resources.Store(uuid.NewUUID(), resource{deployName, serviceName})

	w.WriteHeader(204)
}

//ColdTriggerSpawn creates new container and executes it immediately
func (s *OpenServer) ColdTriggerSpawn(w http.ResponseWriter, r *http.Request) {
	var trigger models.ColdTriggerEvent

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&trigger); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	_, err := k8s.CreateJob(s.client, trigger)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(204)
}

//ListFuncSpecs returns list of funcSpecs
func (s *OpenServer) ListFuncSpecs(w http.ResponseWriter, r *http.Request) {

	m := map[string]resource{}
	s.resources.Range(func(key, value interface{}) bool {
		m[fmt.Sprint(key)] = value.(resource)
		return true
	})

	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(b)+"\n")
}

// DeleteFunc deletes a deployment and service ands stops hot containers
func (s *OpenServer) DeleteFunc(w http.ResponseWriter, r *http.Request) {

	idParam, ok := r.URL.Query()["id"]
	id := types.UID(idParam[0])

	if !ok {
		log.Println("Url Param 'id' is missing")
		return
	}

	pair, ok := s.resources.Load(id)
	if ok {
		specs := pair.(resource)

		if err := k8s.DeleteDeployment(s.client, specs.DeploymentName); err != nil {
			log.Printf("Error deleting deployment %s:\n%v", specs.DeploymentName, err)
		}

		if err := k8s.DeleteService(s.client, specs.ServiceName); err != nil {
			log.Printf("Error deleting service %s:\n%v", specs.ServiceName, err)
		}

		s.resources.Delete(id)

	} else {
		err := fmt.Sprintf("Error loading resources for id %s", id)
		http.Error(w, err, 500)
		return
	}

	w.WriteHeader(204)
}

func (s *OpenServer) cleanUp() {

	s.resources.Range(func(key, value interface{}) bool {
		r := value.(resource)

		log.Printf("Deleting pair %s %s ...", r.DeploymentName, r.ServiceName)

		if err := k8s.DeleteDeployment(s.client, r.DeploymentName); err != nil {
			log.Printf("Error deleting deployment %s:\n%v", r.DeploymentName, err)
		}

		if err := k8s.DeleteService(s.client, r.ServiceName); err != nil {
			log.Printf("Error deleting service %s:\n%v", r.ServiceName, err)
		}

		return true
	})
}

func pong(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
