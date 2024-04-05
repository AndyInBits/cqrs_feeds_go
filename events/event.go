package events

import (
	"andyinbites/cqrs/models"
	"context"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	SuscribeToCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error)
	OnCreatedFeed(f func(CreatedFeedMessage)) error
}

var eventStore EventStore

func SetEventStore(es EventStore) {
	eventStore = es
}

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}

func SuscribeToCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SuscribeToCreatedFeed(ctx)
}

func OnCreatedFeed(f func(CreatedFeedMessage)) error {
	return eventStore.OnCreatedFeed(f)
}
