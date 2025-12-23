package txrunner

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxTxRunner struct {
	pool *pgxpool.Pool
}

func NewTxRunner(pool *pgxpool.Pool) *PgxTxRunner {
	return &PgxTxRunner{pool: pool}
}

func (r *PgxTxRunner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctxTx := context.WithValue(ctx, TxKey{}, tx)

	if err := fn(ctxTx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

type TxKey struct{}
