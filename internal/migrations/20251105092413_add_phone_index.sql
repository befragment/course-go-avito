-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_couriers_phone ON couriers (phone);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_couriers_phone ON couriers;
-- +goose StatementEnd
