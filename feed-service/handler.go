package main

import (
	"andyinbites/cqrs/events"
	"andyinbites/cqrs/models"
	"andyinbites/cqrs/repository"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
)

type createdFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createdFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req createdFeedRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed := models.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}

	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("Failed publish created new feed: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)

}
