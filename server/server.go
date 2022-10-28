package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"ws-go/database"
	"ws-go/repository"
	"ws-go/websocket"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.DatabaseURL == "" {
		return nil, errors.New("DatabaseURL is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWTSecret is required")
	}

	broker := &Broker{
		config: &config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}

	go b.Hub().Run()

	repository.SetRepository(repo)
	log.Println("Starting server on port ", b.Config().Port)

	err = http.ListenAndServe(b.config.Port, b.router)

	if err != nil {
		log.Fatal("Listen and server: ", err)
	}
}
