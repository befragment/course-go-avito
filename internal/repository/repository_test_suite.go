package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"courier-service/internal/model"
)

type RepositoryTestSuite struct {
	suite.Suite
	pool         *pgxpool.Pool
	courierRepo  *CourierRepository
	deliveryRepo *DeliveryRepository
	pgContainer  *postgres.PostgresContainer
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		s.T().Fatalf("Failed to start postgres container: %v", err)
	}
	s.pgContainer = pgContainer

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		s.T().Fatalf("Failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}
	s.pool = pool

	if err := s.runMigrations(ctx); err != nil {
		s.T().Fatalf("Failed to run migrations: %v", err)
	}

	s.courierRepo = NewCourierRepository(pool)
	s.deliveryRepo = NewDeliveryRepository(pool)
}

func (s *RepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()

	if s.pool != nil {
		s.pool.Close()
	}

	if s.pgContainer != nil {
		if err := s.pgContainer.Terminate(ctx); err != nil {
			s.T().Logf("Failed to terminate container: %v", err)
		}
	}
}

func (s *RepositoryTestSuite) TearDownTest() {
	_, err := s.pool.Exec(context.Background(),
		"TRUNCATE TABLE delivery, couriers RESTART IDENTITY CASCADE")
	if err != nil {
		s.T().Fatalf("Failed to truncate tables: %v", err)
	}
}

func (s *RepositoryTestSuite) SetupTest() {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `
		TRUNCATE delivery, couriers
		RESTART IDENTITY
		CASCADE
	`)
	s.Require().NoError(err)
}

func (s *RepositoryTestSuite) createTestCourier(name, phone, transportType string) int64 {
	courier := model.Courier{
		Name:          name,
		Phone:         phone,
		Status:        "available",
		TransportType: transportType,
	}
	id, err := s.courierRepo.CreateCourier(context.Background(), courier)
	if err != nil {
		s.T().Fatalf("Failed to create test courier: %v", err)
	}
	return id
}

func (s *RepositoryTestSuite) runMigrations(ctx context.Context) error {
	migrationsDir := filepath.Join("..", "..", "migrations")

	migrations, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to find migrations: %w", err)
	}

	for _, migration := range migrations {
		sqlContent, err := os.ReadFile(migration)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		upSQL := extractUpMigration(string(sqlContent))
		if upSQL == "" {
			continue
		}

		fmt.Println("Executing migration: ", migration)
		if _, err := s.pool.Exec(ctx, upSQL); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}
	}

	return nil
}

func extractUpMigration(content string) string {
	lines := strings.Split(content, "\n")
	var upSQL []string
	inUp := false
	inStatement := false

	for _, line := range lines {
		if strings.Contains(line, "-- +goose Up") {
			inUp = true
			continue
		}
		if strings.Contains(line, "-- +goose Down") {
			break
		}
		if strings.Contains(line, "-- +goose StatementBegin") {
			inStatement = true
			continue
		}
		if strings.Contains(line, "-- +goose StatementEnd") {
			inStatement = false
			continue
		}
		if inUp && inStatement {
			upSQL = append(upSQL, line)
		}
	}

	return strings.TrimSpace(strings.Join(upSQL, "\n"))
}
