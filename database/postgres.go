package database

import (
	"andyinbites/cqrs/models"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}

func (repo *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO feeds (id, title, description) VALUES ($1, $2, $3)", feed.ID, feed.Title, feed.Description)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) ListFeed(ctx context.Context) ([]*models.Feed, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, title, description, created_at FROM feeds")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	feeds := []*models.Feed{}
	for rows.Next() {
		feed := &models.Feed{}
		if err := rows.Scan(&feed.ID, &feed.Title, &feed.Description, &feed.CreatedAt); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}
