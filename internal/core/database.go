package core

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool(ctx context.Context) *pgxpool.Pool {
	appCfg, _ := LoadConfig()
	cfg, err := pgxpool.ParseConfig(appCfg.DBConnString())
	if err != nil {
		log.Fatal(err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour
	cfg.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Database connection pool established")

	return pool
}
