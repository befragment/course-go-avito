package core

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustInitPool() *pgxpool.Pool {
	appCfg, _ := LoadConfig()
	cfg, err := pgxpool.ParseConfig(appCfg.DBConnString())
	if err != nil {
		log.Fatal(err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = time.Minute * 30

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	pingAttemptLimit := 3
	var pingErr error
	for i := 0; i < pingAttemptLimit; i++ {
		pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
		pingErr = pool.Ping(pingCtx)
		pingCancel()
		if pingErr == nil {
			break
		}
		log.Printf("db ping attempt %d failed: %v", i+1, pingErr)
		if i < pingAttemptLimit {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if pingErr != nil {
		log.Fatalf("Unable to ping database")
	}
	log.Println("Database connection pool established")

	return pool
}