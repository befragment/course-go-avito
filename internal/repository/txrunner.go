package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxTxRunner struct {
    pool *pgxpool.Pool
}

func NewTxRunner(pool *pgxpool.Pool) *pgxTxRunner {
    return &pgxTxRunner{pool: pool}
}

func (r *pgxTxRunner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return err
    }

    ctxTx := context.WithValue(ctx, txKey{}, tx)

    if err := fn(ctxTx); err != nil {
        _ = tx.Rollback(ctx)
        return err
    }
    return tx.Commit(ctx)
}

type txKey struct{}