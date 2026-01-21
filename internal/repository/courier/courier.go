package courier

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"courier-service/internal/model"
	entity "courier-service/internal/repository/entity"
	db "courier-service/internal/repository/utils/database"
	logger "courier-service/pkg/logger"
)

type CourierRepository struct {
	pool   *pgxpool.Pool
	logger logger.LoggerInterface
}

func NewCourierRepository(pool *pgxpool.Pool, logger logger.LoggerInterface) *CourierRepository {
	return &CourierRepository{pool: pool, logger: logger}
}

func (r *CourierRepository) GetCourierById(ctx context.Context, id int64) (model.Courier, error) {
	queryBuilder := sq.
		Select(db.IDColumn, db.NameColumn, db.PhoneColumn, db.StatusColumn, db.TransportTypeColumn, db.CreatedAtColumn, db.UpdatedAtColumn).
		From(db.CourierTable).
		Where(sq.Eq{db.IDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return model.Courier{}, err
	}

	var c entity.CourierDB

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType, &c.CreatedAt, &c.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Courier{}, ErrCourierNotFound
	}
	if err != nil {
		return model.Courier{}, err
	}

	return c.ToModel(), nil
}

func (r *CourierRepository) GetAllCouriers(ctx context.Context) ([]model.Courier, error) {
	queryBuilder := sq.
		Select(db.IDColumn, db.NameColumn, db.PhoneColumn, db.StatusColumn, db.TransportTypeColumn, db.CreatedAtColumn, db.UpdatedAtColumn).
		From(db.CourierTable).
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
		var c entity.CourierDB
		err = rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, c.ToModel())
	}
	return couriers, nil
}

func (r *CourierRepository) CreateCourier(ctx context.Context, courier model.Courier) (int64, error) {
	var id int64
	queryBuilder := sq.
		Insert(db.CourierTable).
		Columns(db.NameColumn, db.PhoneColumn, db.StatusColumn, db.TransportTypeColumn, db.CreatedAtColumn, db.UpdatedAtColumn).
		Values(courier.Name, courier.Phone, courier.Status, courier.TransportType, time.Now(), time.Now()).
		Suffix(db.BuildReturningStatement(db.IDColumn)).
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
	sets := sq.Eq{db.UpdatedAtColumn: time.Now()}
	if courier.Name != "" {
		sets[db.NameColumn] = courier.Name
	}
	if courier.Phone != "" {
		sets[db.PhoneColumn] = courier.Phone
	}
	if courier.Status != "" {
		sets[db.StatusColumn] = courier.Status
	}
	if courier.TransportType != "" {
		sets[db.TransportTypeColumn] = courier.TransportType
	}
	if len(sets) == 1 {
		return ErrNothingToUpdate
	}

	queryBuilder := sq.
		Update(db.CourierTable).
		SetMap(sets).
		Where(sq.Eq{db.IDColumn: courier.ID}).
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
	// Great idea to add index here on delivery.courier_id
	queryBuilder := sq.
		Select(db.CourierID, db.CourierName, db.CourierPhone, db.CourierStatus, db.CourierTransportType).
		From(db.CourierTable).
		LeftJoin(fmt.Sprintf("%s ON %s = %s", db.DeliveryTable, db.CourierID, db.DeliveryCourierID)).
		Where(sq.Eq{db.CourierStatus: db.StatusAvailable}).
		GroupBy(db.CourierID).
		OrderBy(fmt.Sprintf("COUNT(%s) asc", db.DeliveryID)).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return model.Courier{}, err
	}

	var c entity.CourierDB
	err = r.pool.QueryRow(ctx, query, args...).
		Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Courier{}, ErrCouriersBusy
		}
		return model.Courier{}, err
	}

	return c.ToModel(), nil
}

func (r *CourierRepository) FreeCouriersWithInterval(ctx context.Context) error {
	sm := sq.Eq{
		db.UpdatedAtColumn: time.Now(),
		db.StatusColumn:    db.StatusAvailable,
	}

	queryBuilder := sq.
		Update(db.CourierTable + " c").
		SetMap(sm).
		Where(sq.Eq{db.StatusColumn: db.StatusBusy}).
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

	if ct.RowsAffected() > 0 {
		r.logger.Debugf("FreeCouriers updated rows: %d", ct.RowsAffected())
	}

	return nil
}

func (r *CourierRepository) ExistsCourierByPhone(ctx context.Context, phone string) (bool, error) {
	queryBuilder := sq.
		Select(db.CountAll).
		From(db.CourierTable).
		Where(sq.Eq{db.CourierPhone: phone}).
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

func (r *CourierRepository) GetCourierIDByOrderID(ctx context.Context, orderID string) (int64, error) {
	queryBuilder := sq.
		Select(db.DeliveryCourierID).
		From(db.DeliveryTable).
		Where(sq.Eq{db.DeliveryOrderID: orderID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return 0, err
	}

	var courierID int64
	err = r.pool.QueryRow(ctx, query, args...).Scan(&courierID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrOrderNotFound
		}
		return 0, fmt.Errorf("database error: %w", err)
	}
	return courierID, nil
}
