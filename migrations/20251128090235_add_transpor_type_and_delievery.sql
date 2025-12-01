-- +goose Up
-- +goose StatementBegin
ALTER TABLE couriers
    ADD COLUMN IF NOT EXISTS transport_type TEXT NOT NULL DEFAULT 'on_foot';
ALTER TABLE couriers
    ALTER COLUMN status SET DEFAULT 'available';
CREATE TABLE IF NOT EXISTS delivery (
    id BIGSERIAL PRIMARY KEY,
    courier_id BIGINT NOT NULL,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deadline TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (courier_id) REFERENCES couriers(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE couriers DROP COLUMN IF EXISTS transport_type;
ALTER TABLE couriers ALTER COLUMN status SET DEFAULT 'inactive';
DROP TABLE IF EXISTS delivery;
-- +goose StatementEnd
