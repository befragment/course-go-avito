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

func (r *CourierRepository) GetCourierById(ctx context.Context, id int64) (*model.CourierDB, error) {
	queryBuilder := sq.
		Select("id", "name", "phone", "status", "transport_type").
		From("couriers").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return nil, err
	}

	var c model.CourierDB

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCourierNotFound
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CourierRepository) GetAllCouriers(ctx context.Context) ([]model.CourierDB, error) {
	queryBuilder := sq.
		Select("id", "name", "phone", "status", "transport_type").
		From("couriers").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return []model.CourierDB{}, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return []model.CourierDB{}, err
	}

	defer rows.Close()

	var couriers []model.CourierDB
	for rows.Next() {
		var c model.CourierDB
		err = rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)
		if err != nil {
			return []model.CourierDB{}, err
		}
		couriers = append(couriers, c)
	}
	return couriers, nil
}

func (r *CourierRepository) CreateCourier(ctx context.Context, courier *model.CourierDB) (int64, error) {
	var id int64
	queryBuilder := sq.
		Insert("couriers").
		Columns("name", "phone", "status", "transport_type", "created_at", "updated_at").
		Values(courier.Name, courier.Phone, courier.Status, courier.TransportType, time.Now(), time.Now()).
		Suffix("RETURNING id").
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

func (r *CourierRepository) UpdateCourier(ctx context.Context, courier *model.CourierDB) error {
	sets := sq.Eq{"updated_at": time.Now()}
	if courier.Name != "" {
		sets["name"] = courier.Name
	}
	if courier.Phone != "" {
		sets["phone"] = courier.Phone
	}
	if courier.Status != "" {
		sets["status"] = courier.Status
	}
	if courier.TransportType != "" {
		sets["transport_type"] = courier.TransportType
	}
	if len(sets) == 1 {
		return ErrNothingToUpdate
	}
	queryBuilder := sq.
		Update("couriers").
		SetMap(sets).
		Where(sq.Eq{"id": courier.ID}).
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

func (r *CourierRepository) FindAvailableCourier(ctx context.Context) (*model.CourierDB, error) {
	queryBuilder := sq.
		Select("c.id", "c.name", "c.phone", "c.status", "c.transport_type").
		From("couriers c").
		LeftJoin("delivery d on c.id = d.courier_id").
		Where(sq.Eq{"c.status": "available"}).
		GroupBy("c.id").
		OrderBy("COUNT(d.id) asc").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var c model.CourierDB
	err = r.pool.QueryRow(ctx, query, args...).
		Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCouriersBusy
		}
		return nil, err
	}

	return &c, nil
}

func (r *CourierRepository) FreeCouriersWithInterval(ctx context.Context) error {
	sm := sq.Eq{
		"updated_at": time.Now(),
		"status":     "available",
	}

	queryBuilder := sq.
		Update("couriers c").
		SetMap(sm).
		Where(sq.Expr("c.status = 'busy'")).
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
		Select("count(*)").
		From("couriers").
		Where(sq.Eq{"phone": phone}).
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
