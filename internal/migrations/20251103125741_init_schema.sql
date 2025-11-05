-- +goose Up
-- +goose StatementBegin
CREATE TYPE courier_status AS ENUM ('active', 'inactive');

CREATE TABLE IF NOT EXISTS couriers (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    phone       TEXT NOT NULL UNIQUE,
    status      courier_status NOT NULL,
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS courier_status;
DROP TABLE IF EXISTS couriers;
-- +goose StatementEnd