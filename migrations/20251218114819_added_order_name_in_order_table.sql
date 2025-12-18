-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER table "order" ADD COLUMN order_name text
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER table "order" DROP COLUMN order_name text
-- +goose StatementEnd
