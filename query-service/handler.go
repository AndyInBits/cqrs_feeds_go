package main

import (
	"andyinbites/cqrs/events"
	"andyinbites/cqrs/models"
	"andyinbites/cqrs/repository"
	"andyinbites/cqrs/search"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func onCreatedFeed(m events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), &feed); err != nil {
		log.Printf("failed to index feed: %v", err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	feeds, err := repository.ListFeed(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(ctx, query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}
