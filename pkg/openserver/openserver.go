package openserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IoanStoianov/Open-func/pkg/triggers"

	"github.com/gorilla/mux"
)

// OpenServer is the core of Open-func
type OpenServer struct {
	http.Server
	shutdowReq chan bool
	reqCount   uint32
}

// NewServer - server factory
func NewServer(addr uint) (*OpenServer, error) {
	router := newRouter()

	if addr < 1001 {
		return nil, errors.New("Invalid port number")
	}

	s := &OpenServer{
		Server: http.Server{
			Addr:         fmt.Sprintf(":%d", addr),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      router,
		},
		shutdowReq: make(chan bool),
	}

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

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Error on shutdown: %v\n", err)
	}
}

func pong(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "pong\n")
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ping", pong).Methods("GET")
	r.HandleFunc("/prepare", triggers.DeployFunc).Methods("POST")
	r.HandleFunc("/test", triggers.HTTPTriggerRedirect).Methods("POST")

	staticFileDirectory := http.Dir("./web/open-func/build/")
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	return r
}
