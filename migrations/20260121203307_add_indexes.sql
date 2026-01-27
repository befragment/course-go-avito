-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_delivery_courier_id ON delivery(courier_id);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_delivery_courier_id;
