package delivery

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"courier-service/internal/model"
	db "courier-service/internal/repository/utils/database"
)

type DeliveryRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryRepository(pool *pgxpool.Pool) *DeliveryRepository {
	return &DeliveryRepository{pool: pool}
}

func (r *DeliveryRepository) CreateDelivery(ctx context.Context, delivery model.Delivery) (model.Delivery, error) {
	queryBuilder := sq.
		Insert(db.DeliveryTable).
		Columns(db.OrderIDColumn, db.CourierIDColumn, db.AssignedAtColumn, db.DeadlineColumn).
		Values(delivery.OrderID, delivery.CourierID, delivery.AssignedAt, delivery.Deadline).
		Suffix(db.BuildReturningStatement(db.IDColumn, db.CourierIDColumn, db.OrderIDColumn, db.DeadlineColumn)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return model.Delivery{}, err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&delivery.ID, &delivery.CourierID, &delivery.OrderID, &delivery.Deadline,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return model.Delivery{}, ErrOrderIDExists
		}
		return model.Delivery{}, err
	}

	return delivery, nil
}

func (r *DeliveryRepository) CouriersDelivery(ctx context.Context, orderID string) (model.Delivery, error) {
	queryBuilder := sq.
		Select(db.DeliveryOrderID, db.CourierID).
		From(db.DeliveryTable).
		Join(fmt.Sprintf("%s ON %s = %s", db.CourierTable, db.DeliveryCourierID, db.CourierID)).
		Where(sq.Eq{db.DeliveryOrderID: orderID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return model.Delivery{}, err
	}

	var delivery model.Delivery
	err = r.pool.QueryRow(ctx, query, args...).Scan(&delivery.OrderID, &delivery.CourierID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Delivery{}, ErrOrderIDNotFound
		}
		return model.Delivery{}, err
	}

	return delivery, nil
}

func (r *DeliveryRepository) DeleteDelivery(ctx context.Context, orderID string) error {
	queryBuilder := sq.
		Delete(db.DeliveryTable).
		Where(sq.Eq{db.OrderIDColumn: orderID}).
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
