package integration

import (
	"context"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

func TruncateAll(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, 
		`
		TRUNCATE TABLE couriers, delivery
		RESTART IDENTITY
		CASCADE
	`)
	return err
}

