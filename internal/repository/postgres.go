package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type URLRepository interface {
	SaveURL(ctx context.Context, originalURL string) (uint64, error)
	FindByShortCode(ctx context.Context, code string) (string, error)
	IncrementClickCount(ctx context.Context, code string) error
}

type PostgresRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresRepository(pool *pgxpool.Pool, logger *zap.Logger) *PostgresRepository {
	return &PostgresRepository{pool: pool, logger: logger}
}

func (r *PostgresRepository) SaveURL(ctx context.Context, originalURL string) (uint64, error) {
	var id uint64
	err := r.pool.QueryRow(ctx,
		"INSERT INTO short_urls(original_url) VALUES($1) RETURNING id",
		originalURL,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to save URL: %w", err)
	}
	return id, nil
}

func (r *PostgresRepository) FindByShortCode(ctx context.Context, code string) (string, error) {
	var originalURL string
	err := r.pool.QueryRow(ctx,
		"SELECT original_url FROM short_urls WHERE short_code = $1",
		code,
	).Scan(&originalURL)

	if err != nil {
		return "", fmt.Errorf("URL not found: %w", err)
	}
	return originalURL, nil
}

func (r *PostgresRepository) IncrementClickCount(ctx context.Context, code string) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE short_urls SET click_count = click_count + 1 WHERE short_code = $1",
		code,
	)
	return err
}
