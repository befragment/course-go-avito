package repository

import (
	"context"
	"errors"
	"log"
	"strings"
	"fmt"
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

func (r *CourierRepository) GetById(ctx context.Context, id int64) (*model.CourierDB, error) {
	queryBuilder := sq.
		Select("id", "name", "phone", "status").
		From("couriers").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return nil, err
	}

	var c model.CourierDB

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCourierNotFound
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CourierRepository) GetAll(ctx context.Context) ([]model.CourierDB, error) {
	queryBuilder := sq.
		Select("id", "name", "phone", "status").
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
		err = rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status)
		if err != nil {
			return []model.CourierDB{}, err
		}
		couriers = append(couriers, c)
	}
	return couriers, nil
}

func (r *CourierRepository) Create(ctx context.Context, courier *model.CourierDB) (int64, error) {
	var id int64
	queryBuilder := sq.
		Insert("couriers").
		Columns("name", "phone", "status", "created_at", "updated_at").
		Values(courier.Name, courier.Phone, courier.Status, time.Now(), time.Now()).
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

func (r *CourierRepository) Update(ctx context.Context, courier *model.CourierDB) error {
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
	if len(sets) == 1 {
		return ErrNothingToUpdate
	}
	queryBuilder := sq.
		Update("couriers").
		SetMap(sets).
		Where(sq.Eq{"id": courier.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	log.Printf("Query: %s, Args: %v", query, args)
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

func (r *CourierRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
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
