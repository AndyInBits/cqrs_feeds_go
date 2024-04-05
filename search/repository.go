package search

import (
	"andyinbites/cqrs/models"
	"context"
)

type SearchRepository interface {
	Close()
	IndexFeed(ctx context.Context, feed *models.Feed) error
	SearchFeed(ctx context.Context, query string) ([]models.Feed, error)
}

var searchRepo SearchRepository

func SetSearchRepository(repo SearchRepository) {
	searchRepo = repo
}

func Close() {
	searchRepo.Close()
}

func IndexFeed(ctx context.Context, feed *models.Feed) error {
	return searchRepo.IndexFeed(ctx, feed)
}

func SearchFeed(ctx context.Context, query string) ([]models.Feed, error) {
	return searchRepo.SearchFeed(ctx, query)
}
