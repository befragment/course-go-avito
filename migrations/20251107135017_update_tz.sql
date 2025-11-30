-- +goose Up
-- +goose StatementBegin
ALTER DATABASE test_db SET TIMEZONE TO 'Europe/Moscow';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER DATABASE test_db SET TIMEZONE TO 'UTC';
-- +goose StatementEnd