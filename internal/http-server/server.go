package serverhttp

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

type serverHTTP struct {
	server  *http.Server
	service Service
}

// Configuration is settings for http server
type Configuration struct {
	Port    string
	Service Service
	// TODO:
	// Logger
}

// New creates new instance of http server
func New(ctx context.Context, conf Configuration) (*serverHTTP, error) {
	if conf.Service == nil {
		return nil, errors.New("Service is not set for HTTP server")
	}

	baseContext := func(net.Listener) context.Context {
		return ctx
	}

	s := &http.Server{
		Addr:         ":" + conf.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		BaseContext:  baseContext,
	}

	server := serverHTTP{
		server:  s,
		service: conf.Service,
	}

	server.registerHandlers()

	return &server, nil
}

// Launch starts http server
func (s *serverHTTP) Launch(ctx context.Context) error {
	log.Println("HTTP server listening port " + s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panic(err)
		return err
	}

	return nil
}

// Close shut down http server
func (s *serverHTTP) Close(ctx context.Context) error {
	log.Println("HTTP server is shutting down...")

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("HTTP server is off.")
	return nil
}

func (s *serverHTTP) registerHandlers() {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth", s.auth)
	mux.HandleFunc("/token_refresh", s.checkToken(s.refresh))
	mux.HandleFunc("/token_refresh_remove", s.checkToken(s.removeRefreshToken))
	mux.HandleFunc("/token_refresh_remove_all", s.checkToken(s.removeAllResfreshToken))

	s.server.Handler = mux
}
