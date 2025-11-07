-- +goose Up
-- +goose StatementBegin
-- Create b-tree index, because we search exact equality of phone number
CREATE INDEX IF NOT EXISTS idx_couriers_phone ON couriers (phone);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_couriers_phone;
-- +goose StatementEnd
