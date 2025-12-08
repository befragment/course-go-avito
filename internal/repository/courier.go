package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"courier-service/internal/model"
)

type CourierRepository struct {
	pool *pgxpool.Pool
}

func NewCourierRepository(pool *pgxpool.Pool) *CourierRepository {
	return &CourierRepository{pool: pool}
}

func (r *CourierRepository) GetCourierById(ctx context.Context, id int64) (model.Courier, error) {
	queryBuilder := sq.
		Select(idColumn, nameColumn, phoneColumn, statusColumn, transportTypeColumn).
		From(courierTable).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return model.Courier{}, err
	}

	var c CourierDB

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Courier{}, ErrCourierNotFound
	}
	if err != nil {
		return model.Courier{}, err
	}

	return model.Courier(c), nil
}

func (r *CourierRepository) GetAllCouriers(ctx context.Context) ([]model.Courier, error) {
	queryBuilder := sq.
		Select(idColumn, nameColumn, phoneColumn, statusColumn, transportTypeColumn).
		From(courierTable).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var couriers []model.Courier
	for rows.Next() {
		var c CourierDB
		err = rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, model.Courier(c))
	}
	return couriers, nil
}

func (r *CourierRepository) CreateCourier(ctx context.Context, courier model.Courier) (int64, error) {
	var id int64
	queryBuilder := sq.
		Insert(courierTable).
		Columns(nameColumn, phoneColumn, statusColumn, transportTypeColumn, createdAtColumn, updatedAtColumn).
		Values(courier.Name, courier.Phone, courier.Status, courier.TransportType, time.Now(), time.Now()).
		Suffix(buildReturningStatement(idColumn)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return 0, ErrPhoneNumberExists
		}
		return 0, fmt.Errorf("database error: %w", err)
	}

	return id, nil
}

func (r *CourierRepository) UpdateCourier(ctx context.Context, courier model.Courier) error {
	sets := sq.Eq{updatedAtColumn: time.Now()}
	if courier.Name != "" {
		sets[nameColumn] = courier.Name
	}
	if courier.Phone != "" {
		sets[phoneColumn] = courier.Phone
	}
	if courier.Status != "" {
		sets[statusColumn] = courier.Status
	}
	if courier.TransportType != "" {
		sets[transportTypeColumn] = courier.TransportType
	}
	if len(sets) == 1 {
		return ErrNothingToUpdate
	}

	queryBuilder := sq.
		Update(courierTable).
		SetMap(sets).
		Where(sq.Eq{idColumn: courier.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return err
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrCourierNotFound
	}

	return nil
}

func (r *CourierRepository) FindAvailableCourier(ctx context.Context) (model.Courier, error) {
	queryBuilder := sq.
		Select(courierID, courierName, courierPhone, courierStatus, courierTransportType).
		From(courierTable).
		LeftJoin(fmt.Sprintf("%s ON %s = %s", deliveryTable, courierID, deliveryCourierID)).
		Where(sq.Eq{courierStatus: statusAvailable}).
		GroupBy(courierID).
		OrderBy(fmt.Sprintf("COUNT(%s) asc", deliveryID)).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return model.Courier{}, err
	}

	var c CourierDB
	err = r.pool.QueryRow(ctx, query, args...).
		Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Courier{}, ErrCouriersBusy
		}
		return model.Courier{}, err
	}	

	return model.Courier(c), nil
}

func (r *CourierRepository) FreeCouriersWithInterval(ctx context.Context) error {
	sm := sq.Eq{
		updatedAtColumn: time.Now(),
		statusColumn:    statusAvailable,
	}

	queryBuilder := sq.
		Update(courierTable + " c").
		SetMap(sm).
		Where(sq.Eq{statusColumn: statusBusy}).
		Where(sq.Expr(`
			EXISTS (
				SELECT 1
				FROM delivery d
				WHERE d.courier_id = c.id
				AND d.assigned_at = (
					SELECT MAX(d2.assigned_at)
					FROM delivery d2
					WHERE d2.courier_id = c.id
				)
				AND d.deadline < NOW()
			)
		`)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	ct, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	log.Printf("FreeCouriers updated rows: %d", ct.RowsAffected())
	return nil
}

func (r *CourierRepository) ExistsCourierByPhone(ctx context.Context, phone string) (bool, error) {
	queryBuilder := sq.
		Select(countAll).
		From(courierTable).
		Where(sq.Eq{courierPhone: phone}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return false, err
	}

	var count int64
	err = r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
