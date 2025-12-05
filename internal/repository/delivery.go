package repository

import (
	"context"
	"courier-service/internal/model"
	"errors"
	"fmt"
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

func (r *DeliveryRepository) CreateDelivery(ctx context.Context, delivery *model.DeliveryDB) (*model.Delivery, error) {
	queryBuilder := sq.
		Insert(deliveryTable).
		Columns(orderIdColumn, courierIdColumn, assignedAtColumn, deadlineColumn).
		Values(delivery.OrderID, delivery.CourierID, delivery.AssignedAt, delivery.Deadline).
		Suffix(fmt.Sprintf("RETURNING %s, %s, %s, %s", idColumn, courierIdColumn, orderIdColumn, deadlineColumn)).
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
		Select(deliveryOrderID, courierID).
		From(deliveryTable).
		Join(fmt.Sprintf("%s ON %s = %s", courierTable, deliveryCourierID, courierID)).
		Where(sq.Eq{deliveryOrderID: orderID}).
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

func (r *DeliveryRepository) DeleteDelivery(ctx context.Context, orderID string) error {
	queryBuilder := sq.
		Delete(deliveryTable).
		Where(sq.Eq{orderIdColumn: orderID}).
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
