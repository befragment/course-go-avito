package integration

import (
	"context"
	"database/sql"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	syncOnce  sync.Once
	connStr   string
	container *postgres.PostgresContainer
)

func TestWithMigrations() (*postgres.PostgresContainer, string, error) {
	ctx := context.Background()
	var creationError error

	syncOnce.Do(func() {
		container, creationError = postgres.RunContainer(ctx,
			testcontainers.WithImage("postgres:15-alpine"),
			postgres.WithDatabase("test_db"),
			postgres.WithUsername("testuser"),
			postgres.WithPassword("testpass"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Minute)),
		)
		if creationError != nil {
			return
		}

		connStr, creationError = container.ConnectionString(ctx, "sslmode=disable")
		if creationError != nil {
			return
		}

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			creationError = err
		}
		defer func() {
			if cerr := db.Close(); cerr != nil && creationError == nil {
				creationError = cerr
			}
		}()
		if err := goose.Up(db, "../../../migrations"); err != nil {
			creationError = err
		}
	})

	if creationError != nil {
		return nil, "", creationError
	}

	return container, connStr, nil
}
