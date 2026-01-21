package txrunner_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"courier-service/internal/model"
	integration "courier-service/internal/persistence/database/integration"
	courierstorage "courier-service/internal/repository/courier"
	deliverystorage "courier-service/internal/repository/delivery"
	txrunner "courier-service/internal/repository/txrunner"
	logger "courier-service/pkg/logger"
)

type TxRunnerTestSuite struct {
	suite.Suite
	ctx          context.Context
	pool         *pgxpool.Pool
	deliveryRepo *deliverystorage.DeliveryRepository
	courierRepo  *courierstorage.CourierRepository
	txRunner     *txrunner.PgxTxRunner
	pgContainer  *postgres.PostgresContainer
}

func TestTxRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(TxRunnerTestSuite))
}

func (s *TxRunnerTestSuite) SetupSuite() {
	s.ctx = context.Background()

	container, connStr, err := integration.TestWithMigrations()
	s.Require().NoError(err)
	s.pgContainer = container

	pool, err := pgxpool.New(s.ctx, connStr)
	s.Require().NoError(err)
	s.pool = pool

	log, err := logger.New(logger.LogLevelInfo)
	s.Require().NoError(err)

	s.courierRepo = courierstorage.NewCourierRepository(s.pool, log)
	s.deliveryRepo = deliverystorage.NewDeliveryRepository(s.pool)
	s.txRunner = txrunner.NewTxRunner(s.pool)
}

func (s *TxRunnerTestSuite) SetupTest() {
	err := integration.TruncateAll(s.ctx, s.pool)
	s.Require().NoError(err)
}

func (s *TxRunnerTestSuite) TestTxRunner_CommitOnSuccess() {
	ctx := context.Background()

	var createdID int64
	err := s.txRunner.Run(ctx, func(txCtx context.Context) error {
		courier := model.Courier{
			Name:          "John Doe",
			Phone:         "+79991234567",
			Status:        "available",
			TransportType: "car",
		}
		id, err := s.courierRepo.CreateCourier(txCtx, courier)
		if err != nil {
			return err
		}
		createdID = id
		return nil
	})

	s.Require().NoError(err)
	s.Greater(createdID, int64(0))

	result, err := s.courierRepo.GetCourierById(ctx, createdID)
	s.Require().NoError(err)
	s.Equal("John Doe", result.Name)
	s.Equal("+79991234567", result.Phone)
}

func (s *TxRunnerTestSuite) TestTxRunner_RollbackOnError() {
	ctx := context.Background()

	expectedErr := errors.New("test error")

	err := s.txRunner.Run(ctx, func(txCtx context.Context) error {
		tx, ok := txCtx.Value(txrunner.TxKey{}).(pgx.Tx)
		if !ok {
			return errors.New("no transaction in context")
		}

		_, execErr := tx.Exec(txCtx,
			"INSERT INTO couriers (name, phone, status, transport_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			"Jane Doe", "+79991234568", "available", "bike", time.Now(), time.Now())
		if execErr != nil {
			return execErr
		}

		return expectedErr
	})

	s.Require().Error(err)
	s.Equal(expectedErr, err)

	var count int
	err = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM couriers WHERE phone = $1", "+79991234568").Scan(&count)
	s.NoError(err)
	s.Equal(0, count, "courier should not exist after rollback")
}

func (s *TxRunnerTestSuite) TestTxRunner_RollbackOnDatabaseError() {
	ctx := context.Background()

	err := s.txRunner.Run(ctx, func(txCtx context.Context) error {
		tx, ok := txCtx.Value(txrunner.TxKey{}).(pgx.Tx)
		if !ok {
			return errors.New("no transaction in context")
		}

		_, err := tx.Exec(txCtx,
			"INSERT INTO couriers (name, phone, status, transport_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			"Alice", "+79991234569", "available", "car", time.Now(), time.Now())
		if err != nil {
			return err
		}

		_, err = tx.Exec(txCtx,
			"INSERT INTO couriers (name, phone, status, transport_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			"Bob", "+79991234569", "available", "bike", time.Now(), time.Now())
		return err
	})

	s.Require().Error(err)
	s.Contains(err.Error(), "duplicate key", "should fail on duplicate phone")

	var count int
	err = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM couriers WHERE phone = $1", "+79991234569").Scan(&count)
	s.NoError(err)
	s.Equal(0, count, "no couriers should exist after rollback")
}

func (s *TxRunnerTestSuite) TestTxRunner_MultipleOperations() {
	ctx := context.Background()

	var courierID int64
	var deliveryCreated bool

	err := s.txRunner.Run(ctx, func(txCtx context.Context) error {
		courier := model.Courier{
			Name:          "Multi Op Courier",
			Phone:         "+79991234570",
			Status:        "available",
			TransportType: "car",
		}
		id, err := s.courierRepo.CreateCourier(txCtx, courier)
		if err != nil {
			return err
		}
		courierID = id

		now := time.Now()
		delivery := model.Delivery{
			CourierID:  courierID,
			OrderID:    "550e8400-e29b-41d4-a716-446655440001",
			AssignedAt: now,
			Deadline:   now.Add(2 * time.Hour),
		}
		_, err = s.deliveryRepo.CreateDelivery(txCtx, delivery)
		if err != nil {
			return err
		}
		deliveryCreated = true

		return nil
	})

	s.Require().NoError(err)
	s.Greater(courierID, int64(0))
	s.True(deliveryCreated)

	courier, err := s.courierRepo.GetCourierById(ctx, courierID)
	s.Require().NoError(err)
	s.Equal("Multi Op Courier", courier.Name)

	delivery, err := s.deliveryRepo.CouriersDelivery(ctx, "550e8400-e29b-41d4-a716-446655440001")
	s.Require().NoError(err)
	s.Equal(courierID, delivery.CourierID)
}

func (s *TxRunnerTestSuite) TestTxRunner_NestedTransactionAttempt() {
	ctx := context.Background()

	var outerID, innerID int64

	err := s.txRunner.Run(ctx, func(txCtx context.Context) error {
		courier1 := model.Courier{
			Name:          "Outer",
			Phone:         "+79991234571",
			Status:        "available",
			TransportType: "car",
		}
		id, err := s.courierRepo.CreateCourier(txCtx, courier1)
		if err != nil {
			return err
		}
		outerID = id

		err = s.txRunner.Run(txCtx, func(nestedCtx context.Context) error {
			courier2 := model.Courier{
				Name:          "Inner",
				Phone:         "+79991234572",
				Status:        "available",
				TransportType: "bike",
			}
			id, err := s.courierRepo.CreateCourier(nestedCtx, courier2)
			if err != nil {
				return err
			}
			innerID = id
			return nil
		})
		return err
	})

	s.Require().NoError(err)
	s.Greater(outerID, int64(0))
	s.Greater(innerID, int64(0))

	_, err = s.courierRepo.GetCourierById(ctx, outerID)
	s.NoError(err)
	_, err = s.courierRepo.GetCourierById(ctx, innerID)
	s.NoError(err)
}

func (s *TxRunnerTestSuite) TestTxRunner_EmptyTransaction() {
	ctx := context.Background()

	err := s.txRunner.Run(ctx, func(txCtx context.Context) error {
		return nil
	})

	s.NoError(err)
}
