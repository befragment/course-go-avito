package database

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"courier-service/internal/core"
	"courier-service/internal/models"
)

func GetCourier(ctx context.Context, id int64) (models.Courier, error) {
	db := core.InitPool(ctx)
	defer db.Close()

	queryBuilder := sq.
		Select("id", "name", "phone", "status").
		From("couriers").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return models.Courier{}, err
	}

	var c models.Courier

	err = db.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Courier{}, pgx.ErrNoRows
	}
	if err != nil {
		return models.Courier{}, err
	}

	return c, nil
}

func GetCouriers(ctx context.Context) ([]models.Courier, error) {
	db := core.InitPool(ctx)
	defer db.Close()

	queryBuilder := sq.
		Select("id", "name", "phone", "status").
		From("couriers").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return []models.Courier{}, err
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return []models.Courier{}, err
	}

	defer rows.Close()

	var couriers []models.Courier
	for rows.Next() {
		var c models.Courier
		err = rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status)
		if err != nil {
			return []models.Courier{}, err
		}
		couriers = append(couriers, c)
	}
	return couriers, nil
}

func CreateCourier(ctx context.Context, courier models.Courier) (models.Courier, error) {
	db := core.InitPool(ctx)
	defer db.Close()

	queryBuilder := sq.
		Insert("couriers").
		Columns("name", "phone", "status").
		Values(courier.Name, courier.Phone, courier.Status).
		Suffix("RETURNING id, name, phone, status").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return models.Courier{}, err
	}

	err = db.QueryRow(ctx, query, args...).Scan(
		&courier.ID, &courier.Name, &courier.Phone, &courier.Status,
	)
	if err != nil {
		return models.Courier{}, err
	}

	return courier, nil
}

func UpdateCourier(ctx context.Context, courier models.Courier) (models.Courier, error) {
	db := core.InitPool(ctx)
	defer db.Close()

	queryBuilder := sq.
		Update("couriers").
		Set("name", courier.Name).
		Set("phone", courier.Phone).
		Set("status", courier.Status).
		Where(sq.Eq{"id": courier.ID}).
		Suffix("RETURNING id, name, phone, status").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return models.Courier{}, err
	}

	err = db.QueryRow(ctx, query, args...).Scan(
		&courier.ID, &courier.Name, &courier.Phone, &courier.Status,
	)
	if err != nil {
		return models.Courier{}, err
	}

	return courier, nil
}

func PhoneNumberExists(ctx context.Context, phone string) (bool, error) {
	db := core.InitPool(ctx)
	defer db.Close()

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
	err = db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
