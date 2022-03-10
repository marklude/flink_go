package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/marklude/flink_go/logger"
	"github.com/marklude/flink_go/middleware"
)

const defaultAddr = ":8080"

type server struct {
	http.Server
}

type options func(*server)

func NewServer(opts ...options) *server {
	svr := &server{}
	svr.WithOptions(opts...)
	return svr
}

func (s *server) WithOptions(opts ...options) {
	for _, opt := range opts {
		opt(s)
	}
}

func WithAddr(addr string) options {
	return func(s *server) {
		s.Addr = addr
	}
}

func WithTimeout(t time.Duration) options {
	return func(s *server) {
		s.WriteTimeout = t
	}
}

func main() {
	// Get address from env
	addr := os.Getenv("HISTORY_SERVER_LISTEN_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	// start application with some fancy logs
	logger.InfoMessage("=============================================")
	logger.InfoMessage(fmt.Sprintf("Starting Flink on %s...", addr))
	logger.InfoMessage("=============================================")

	r := mux.NewRouter()
	r.HandleFunc("/location/{order_id}/now", middleware.PostLocation).Methods("POST")
	r.HandleFunc("/location/{order_id}", middleware.GetLocation).Methods("GET")
	r.HandleFunc("/location/{order_id}", middleware.DeleteLocation).Methods("DELETE")

	svr := NewServer(WithAddr(addr), WithTimeout(30*time.Second))

	svr.Handler = r

	if err := svr.ListenAndServe(); err != nil {
		logger.FatalMessage("Could not start server", err)
	}

}
