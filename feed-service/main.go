package main

import (
	"andyinbites/cqrs/database"
	"andyinbites/cqrs/events"
	"andyinbites/cqrs/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB   string `envconfig:"POSTGRES_DB" required:"true"`
	PostgresUser string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPass string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	NatAddress   string `envconfig:"NATS_ADDRESS" required:"true"`
}

func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", createdFeedHandler).Methods(http.MethodPost)
	return
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB)

	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatalf("Failed to create database repository: %v", err)
	}
	repository.SetRepository(repo)
	n, err := events.NewNatsEventStore(fmt.Sprintf("nats://%s", cfg.NatAddress))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	events.SetEventStore(n)

	defer events.Close()

	router := NewRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start Web Server: %v", err)
	}

}
