package repository

import (
	"context"
	"courier-service/internal/model"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryRepository(pool *pgxpool.Pool) *DeliveryRepository {
	return &DeliveryRepository{pool: pool}
}

func (r *DeliveryRepository) Create(ctx context.Context, delivery *model.DeliveryDB) (*model.Delivery, error) {
	queryBuilder := sq.
		Insert("delivery").
		Columns("order_id", "courier_id", "assigned_at", "deadline").
		Values(delivery.OrderID, delivery.CourierID, delivery.AssignedAt, delivery.Deadline).
		Suffix("RETURNING id, courier_id, order_id, deadline").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&delivery.ID, &delivery.CourierID, &delivery.OrderID, &delivery.Deadline,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, ErrOrderIDExists
		}
		return nil, err
	}

	out := model.Delivery(*delivery)
	return &out, nil
}

func (r *DeliveryRepository) CouriersDelivery(ctx context.Context, orderID string) (*model.DeliveryDB, error) {
	queryBuilder := sq.
		Select("d.order_id", "c.id").
		From("delivery d").
		Join("couriers c on d.courier_id = c.id").
		Where(sq.Eq{"d.order_id": orderID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var delivery model.DeliveryDB
	err = r.pool.QueryRow(ctx, query, args...).Scan(&delivery.OrderID, &delivery.CourierID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderIDNotFound
		}
		return nil, err
	}

	return &delivery, nil
}

func (r *DeliveryRepository) Delete(ctx context.Context, orderID string) error {
	queryBuilder := sq.
		Delete("delivery").
		Where(sq.Eq{"order_id": orderID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrOrderIDNotFound
	}

	return nil
}
